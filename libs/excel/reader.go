package excel

import (
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"io"
	"reflect"
	"strings"
)

func EasyRead(reader io.Reader, dest any, options ...SheetReaderOption) error {
	fileReader, err := NewReaderV2(reader)
	if err != nil {
		return err
	}
	return EasyReadWithReader(fileReader, dest, options...)
}

func EasyReadWithReader(fileReader *ReaderV2, dest any, options ...SheetReaderOption) error {
	sheetReader, err := fileReader.DefaultSheetReader(options...)
	if err != nil {
		return err
	}
	return sheetReader.ReadAll(dest)
}

type ReaderV2 struct {
	f *excelize.File
}

func (r *ReaderV2) Date1904() bool {
	return r.f.WorkBook.WorkbookPr.Date1904
}

func NewReaderV2(reader io.Reader) (*ReaderV2, error) {
	f, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, err
	}
	return &ReaderV2{
		f: f,
	}, nil
}

type SheetReaderOption func(*SheetReader) error

func SheetReaderWithCustomReadFunc(header string, readFunc ReadCellValueFunc) SheetReaderOption {
	return func(s *SheetReader) error {
		if s.customReadFuncMap == nil {
			s.customReadFuncMap = make(map[string]ReadCellValueFunc)
		}
		s.customReadFuncMap[header] = readFunc
		return nil
	}
}

func (r *ReaderV2) NewSheetReader(sheetName string, options ...SheetReaderOption) (*SheetReader, error) {
	rows, err := r.f.Rows(sheetName)
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		return nil, fmt.Errorf("sheet %s empty", sheetName)
	}
	headers, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	headerMap := make(map[string]int, len(headers))
	for i, v := range headers {
		headerMap[strings.TrimSpace(v)] = i
	}
	sheetReader := SheetReader{
		headers:   headers,
		headerMap: headerMap,
		rows:      rows,
		curRowIdx: 0,
		f:         r.f,
	}
	for _, option := range options {
		if err = option(&sheetReader); err != nil {
			return &sheetReader, err
		}
	}
	return &sheetReader, nil
}

func (r *ReaderV2) DefaultSheetReader(options ...SheetReaderOption) (*SheetReader, error) {
	sheets := r.f.GetSheetList()
	if len(sheets) == 0 {
		return nil, errors.New("no sheet in file")
	}
	return r.NewSheetReader(sheets[0], options...)
}

type SheetReader struct {
	headers           []string
	headerMap         map[string]int
	rows              *excelize.Rows
	customReadFuncMap map[string]ReadCellValueFunc
	curRowIdx         int
	f                 *excelize.File
}

func TitleHeader(reader *SheetReader) error {
	casesT := cases.Title(language.Chinese)
	for i, header := range reader.headers {
		reader.headers[i] = casesT.String(strings.ToLower(header))
	}
	for k, v := range reader.headerMap {
		reader.headerMap[casesT.String(strings.ToLower(k))] = v
	}
	return nil
}

func (r *SheetReader) suggestCap(lens int) int {
	const DefaultCap = 100
	if lens <= 0 {
		return DefaultCap
	}
	return lens
}

func (r *SheetReader) GetHeader() []string {
	return r.headers
}

func (r *SheetReader) getFieldOptionMap(elementType reflect.Type) (map[int]sheetReaderColOption, int, error) {
	fieldOptionMap := make(map[int]sheetReaderColOption, elementType.NumField())
	missingColumns := make([]string, 0, elementType.NumField())
	rawRowIdx := -1
	for i := 0; i < elementType.NumField(); i++ {
		ft := elementType.Field(i)
		if ft.Type == rtReaderRawRow {
			rawRowIdx = i
			continue
		}
		tag := ft.Tag.Get("xlsx")
		if tag == "" {
			continue
		}
		option, err := r.explainTag(ft, tag)
		if err != nil {
			return nil, 0, err
		}
		idx, ok := r.headerMap[option.header]
		if !ok {
			if option.required {
				missingColumns = append(missingColumns, option.header)
			}
			continue
		}
		option.idx = idx
		fieldOptionMap[i] = *option
	}
	if len(missingColumns) != 0 {
		return nil, -1, fmt.Errorf("missing required columns of \"%s\"", strings.Join(missingColumns, ","))
	}
	return fieldOptionMap, rawRowIdx, nil
}

func (r *SheetReader) explainTag(field reflect.StructField, tag string) (*sheetReaderColOption, error) {
	tags := strings.Split(tag, ";")
	key := tags[0]
	var supportOptions struct {
		required   bool
		headerText string
	}
	supportOptions.headerText = key
	for i := 1; i < len(tags); i++ {
		vals := strings.Split(tags[i], "=")
		if len(vals) == 0 {
			return nil, fmt.Errorf("wrong tag of %s", field.Name)
		}
		switch vals[0] {
		case "required":
			supportOptions.required = true
		case "displayheader":
			supportOptions.headerText = vals[1]
		default:
			return nil, fmt.Errorf("unsupported option %s at field %s", tag, field.Name)
		}
	}
	res := newSheetReaderColOption(key, supportOptions.headerText, supportOptions.required, r.customReadFuncMap[key])
	return &res, nil
}

func (r *SheetReader) readDestType(dest any) (reflect.Type, reflect.Type, error) {
	ptrType := reflect.TypeOf(dest)
	if kind := ptrType.Kind(); kind != reflect.Ptr {
		return nil, nil, fmt.Errorf("target type must be pointer as we need to change origin data, but got %s", kind)
	}
	sliceType := ptrType.Elem()
	if kind := sliceType.Kind(); kind != reflect.Slice {
		return nil, nil, fmt.Errorf("target type must be pointer to slice, but not pointer to %s", kind)
	}
	elementType := sliceType.Elem()
	// 我们只支持array struct，不支持array struct pointer，这种滥用指针的风气一定要矫正回来
	if kind := elementType.Kind(); kind != reflect.Struct {
		return nil, nil, fmt.Errorf("target element type must be struct, but not %s", kind)
	}
	return sliceType, elementType, nil
}

func (r *SheetReader) IterRows(lens int, readRow func(idx int, columns []string) error) error {
	// 记得r.Next必须放在最后面，因为它是有side effect的
	var readRowErrors ReadRowErrors
	// 因为要支持-1所以需要使用不等于
	for i := 0; (i != lens) && r.Next(); i++ {
		columns, err := r.rows.Columns()
		if err != nil {
			return err
		}
		if r.isEmptyRow(columns) {
			if err := r.emptyRowsCheck(); err != nil {
				return err
			}
			break
		}
		curIdx := r.GetCurIdx()
		columns = r.fixRowLens(columns)
		if err = readRow(curIdx, columns); err != nil {
			readRowErrors = append(readRowErrors, NewReadRowError(curIdx, columns, err))
		}
	}
	if len(readRowErrors) == 0 {
		return nil
	}
	return &readRowErrors
}

func (r *SheetReader) isEmptyRow(columns []string) bool {
	if len(columns) == 0 {
		return true
	}
	for _, v := range columns {
		if v != "" {
			return false
		}
	}
	return true
}

// 有些Excel的写入，对于最后的空行会把value干掉
func (r *SheetReader) fixRowLens(columns []string) []string {
	for i := len(columns); i < len(r.headers); i++ {
		columns = append(columns, "")
	}
	return columns
}

// 有一些软件为了性能或者逻辑简单，会提前插入很多行,
// 所以读取的时候我们要判断遇到空行后，禁止再出现数据行
func (r *SheetReader) emptyRowsCheck() error {
	for r.Next() {
		columns, err := r.rows.Columns()
		if err != nil {
			return err
		}
		if !r.isEmptyRow(columns) {
			return fmt.Errorf("row[%d]: You cannot set row after empty rows", r.GetCurIdx())
		}
	}
	return nil
}

func (r *SheetReader) Next() bool {
	if !r.rows.Next() {
		return false
	}
	r.curRowIdx++
	return true
}

func (r *SheetReader) GetCurIdx() int {
	return r.curRowIdx
}

func (r *SheetReader) File() *excelize.File {
	return r.f
}

// ReadSome 非文件问题情况下，即使有部分列读取失败，我们也会修改dest
func (r *SheetReader) ReadSome(dest any, lens int) error {
	sliceType, elementType, err := r.readDestType(dest)
	if err != nil {
		return err
	}
	fieldOptionMap, rawRowIdx, err := r.getFieldOptionMap(elementType)
	if err != nil {
		return err
	}
	sliceValue := reflect.MakeSlice(sliceType, 0, r.suggestCap(lens))
	err = r.IterRows(lens, func(idx int, columns []string) error {
		elementValue := reflect.New(elementType).Elem()
		for i := 0; i < elementValue.NumField(); i++ {
			fv := elementValue.Field(i)
			if i == rawRowIdx {
				fv.Set(reflect.ValueOf(ReaderRawRow{RowIdx: idx, Columns: columns}))
				continue
			}
			option, ok := fieldOptionMap[i]
			if !ok {
				continue
			}
			if err = option.setValue(columns[option.idx], fv); err != nil {
				return err
			}
		}
		sliceValue = reflect.Append(sliceValue, elementValue)
		return nil
	})
	destValue := reflect.ValueOf(dest)
	destValue.Elem().Set(sliceValue)
	if err != nil {
		return err
	}
	return nil
}

func (r *SheetReader) ReadAll(dest any) error {
	return r.ReadSome(dest, -1)
}
