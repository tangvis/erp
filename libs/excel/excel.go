package excel

import (
	"context"
	"fmt"
	"io"
	"reflect"

	"github.com/xuri/excelize/v2"
)

const (
	DefaultSheetName   = "Sheet1"
	SheetNotExistIndex = -1
)

var (
	DefaultHeaderStyle = &excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
	}
	DefaultCellStyle = &excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "center",
		},
	}
	DefaultHeaderStyleV2 = &excelize.Style{
		Border: []excelize.Border{
			{
				Type:  "right",
				Color: "#cccccc",
				Style: 2,
			},
			{
				Type:  "left",
				Color: "#cccccc",
				Style: 2,
			},
		},
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#666666"},
			Pattern: 1,
			Shading: 0,
		},
	}
)

type WriterV2 struct {
	file   *excelize.File
	sheets map[string]*SheetWriterV2
}

func NewWriterV2(defaultSheetName string, cols []WriterV2ColOption) (*WriterV2, error) {
	file := excelize.NewFile()
	if defaultSheetName != DefaultSheetName {
		file.NewSheet(defaultSheetName)
		file.DeleteSheet(DefaultSheetName)
	}
	w := &WriterV2{
		file:   file,
		sheets: make(map[string]*SheetWriterV2, 1),
	}
	if _, err := w.newSheet(defaultSheetName, cols); err != nil {
		return nil, err
	}
	return w, nil
}

func NewWriterV2ByOpenFile(excelPath string) (*WriterV2, error) {
	file, err := excelize.OpenFile(excelPath)
	w := &WriterV2{
		file:   file,
		sheets: make(map[string]*SheetWriterV2, 1),
	}
	return w, err
}

func (w *WriterV2) Sheet(name string) *SheetWriterV2 {
	return w.sheets[name]
}

// GetRawFile 对于某些场景来说，需要做很多额外的业务特殊的操作，所以我们把原始文件给出去，
// 这样调用方也可以直接对原始文件操作而无需再去在client处做更多别扭的封装。
func (w *WriterV2) GetRawFile() *excelize.File {
	return w.file
}

func (w *WriterV2) AddSheet(name string, cols []WriterV2ColOption) (*SheetWriterV2, error) {
	if exist := w.sheets[name]; exist != nil {
		return nil, fmt.Errorf("already exist %s", name)
	}
	if idx, _ := w.file.GetSheetIndex(name); idx != SheetNotExistIndex {
		return nil, fmt.Errorf("already exist %s", name)
	}
	w.file.NewSheet(name)
	return w.newSheet(name, cols)
}

func (w *WriterV2) SaveAs(filepath string) error {
	return w.file.SaveAs(filepath)
}

func (w *WriterV2) Write(buffer io.Writer) error {
	return w.file.Write(buffer)
}

func (w *WriterV2) convertOptions(colOptions []WriterV2ColOption) (map[string]*WriterV2ColOption, error) {
	res, err := convertOpt(colOptions, w.file)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func convertOpt(colOptions []WriterV2ColOption, f *excelize.File) (map[string]*WriterV2ColOption, error) {
	res := make(map[string]*WriterV2ColOption, len(colOptions))
	var err error
	for i := range colOptions {
		v := colOptions[i]
		if v.headerStyle != nil {
			v.headerStyleID, err = f.NewStyle(v.headerStyle)
			if err != nil {
				return nil, err
			}
		}
		if v.cellStyle != nil {
			v.cellStyleID, err = f.NewStyle(v.cellStyle)
			if err != nil {
				return nil, err
			}
		}
		v.col = indexToCol(i)
		res[v.cellKey] = &v
	}
	return res, nil
}

func (w *WriterV2) newSheet(sheet string, colOptions []WriterV2ColOption) (*SheetWriterV2, error) {
	cols, err := w.convertOptions(colOptions)
	if err != nil {
		return nil, err
	}
	res := &SheetWriterV2{
		cols:   cols,
		sheet:  sheet,
		rowIdx: 1,
		file:   w.file,
	}
	if err = res.writeHeader(1); err != nil {
		return nil, err
	}
	w.sheets[sheet] = res
	return res, nil
}

func (w *WriterV2) WriteCellByXY(ctx context.Context, sheet string, x, y int, val interface{}) error {
	col := indexToCol(x)
	axis := fmt.Sprintf("%s%d", col, y+1)
	if err := w.file.SetCellValue(sheet, axis, val); err != nil {
		return err
	}
	return nil
}

func (w *WriterV2) GetCellByXY(ctx context.Context, sheet string, x, y int) (string, error) {
	col := indexToCol(x)
	axis := fmt.Sprintf("%s%d", col, y+1)
	return w.file.GetCellValue(sheet, axis)
}

func NewWriterV2ColOptionSimple(header, key string) WriterV2ColOption {
	return NewWriterV2ColOption(header, DefaultHeaderStyle, "", key, nil, nil, 0)
}

func NewWriterV2ColOptionSimpleWithStyle(header, key string, headerStyle, cellStyle *excelize.Style) WriterV2ColOption {
	return NewWriterV2ColOption(header, headerStyle, "", key, cellStyle, nil, 0)
}

func NewWriterV2ColOptionSimple3(header, key string, width float64, headerComment string) WriterV2ColOption {
	return NewWriterV2ColOption(header, DefaultHeaderStyle, headerComment, key, nil, nil, width)
}
func NewWriterV2ColOptionSimple2(header, key string, width float64) WriterV2ColOption {
	return NewWriterV2ColOption(header, DefaultHeaderStyle, "", key, nil, nil, width)
}

func NewWriterV2ColOption(
	header string,
	headerStyle *excelize.Style,
	headerComment string,
	cellKey string,
	cellStyle *excelize.Style,
	formatter SheetWriterV2CellFormatter,
	width float64,
) WriterV2ColOption {
	return WriterV2ColOption{
		header:        header,
		headerStyle:   headerStyle,
		headerComment: headerComment,
		cellKey:       cellKey,
		cellStyle:     cellStyle,
		cellFormatter: formatter,
		width:         width,
	}
}

type SheetWriterV2CellFormatter func(key string, row reflect.Value, value reflect.Value) interface{}

type WriterV2ColOption struct {
	// header
	header        string
	headerStyle   *excelize.Style // a json string
	headerComment string
	// 对应的cell Value
	cellKey       string
	cellStyle     *excelize.Style // a json string
	cellFormatter SheetWriterV2CellFormatter
	// 通用的
	width float64
	// 生成的
	col           string
	headerStyleID int
	cellStyleID   int
}

func (w *WriterV2ColOption) Value(row, col reflect.Value) interface{} {
	if w.cellFormatter != nil {
		return w.cellFormatter(w.cellKey, row, col)
	}
	return col.Interface()
}

func indexToCol(index int) string {
	n, _ := excelize.ColumnNumberToName(index + 1)
	return n
}

type SheetWriterV2 struct {
	cols   map[string]*WriterV2ColOption
	sheet  string
	rowIdx int
	file   *excelize.File
}

func (w *SheetWriterV2) markGetRow() int {
	w.rowIdx++
	return w.rowIdx
}

func (w *SheetWriterV2) WriteRow(ctx context.Context, row interface{}) error {
	var err error
	rowIdx := w.markGetRow()
	t := reflect.TypeOf(row)
	v := reflect.ValueOf(row)
	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)
		key := ft.Tag.Get("xlsx")
		col, ok := w.cols[key]
		if !ok {
			continue
		}
		fv := col.Value(v, v.Field(i))
		axis := fmt.Sprintf("%s%d", col.col, rowIdx)
		if err = w.file.SetCellValue(w.sheet, axis, fv); err != nil {
			return err
		}
		if col.cellStyleID != 0 {
			_ = w.file.SetCellStyle(w.sheet, axis, axis, col.cellStyleID)
		}
	}
	return nil
}

func (w *SheetWriterV2) WriteMapRow(ctx context.Context, row map[string]string) error {
	var err error
	rowIdx := w.markGetRow()
	for key, value := range row {
		col, ok := w.cols[key]
		if !ok {
			continue
		}
		axis := fmt.Sprintf("%s%d", col.col, rowIdx)
		if err = w.file.SetCellValue(w.sheet, axis, value); err != nil {
			return err
		}
		if col.cellStyleID != 0 {
			_ = w.file.SetCellStyle(w.sheet, axis, axis, col.cellStyleID)
		}
	}
	return nil
}

func (w *SheetWriterV2) WriteRowString(ctx context.Context, row []string) error {
	var err error
	rowIdx := w.markGetRow()
	for idx, v := range row {
		col := indexToCol(idx)
		axis := fmt.Sprintf("%s%d", col, rowIdx)
		if err = w.file.SetCellValue(w.sheet, axis, v); err != nil {
			return err
		}
	}
	return nil
}

func (w *SheetWriterV2) WriteRowStringBatch(ctx context.Context, rows [][]string) error {
	for _, row := range rows {
		err := w.WriteRowString(ctx, row)
		if err != nil {
			return nil
		}
	}
	return nil
}

// WriteRowBatch 这种传错参数的直接panic就好了，理论上只要有测试过就不会出现这个错误，而那种不测试的情况我们不管
func (w *SheetWriterV2) WriteRowBatch(ctx context.Context, rows interface{}) error {
	return writeRowBatch(rows, w)
}

func writeRowBatch(rows interface{}, w *SheetWriterV2) error {
	sliceType := reflect.TypeOf(rows)
	if kind := sliceType.Kind(); kind != reflect.Slice {
		panic(fmt.Sprintf("must pass slice, but not %s", kind))
	}
	elementType := sliceType.Elem()
	for kind := elementType.Kind(); kind == reflect.Ptr; kind = elementType.Kind() {
		elementType = elementType.Elem()
	}
	if kind := elementType.Kind(); kind != reflect.Struct {
		panic(fmt.Sprintf("element must be struct, but not %s", kind))
	}
	optionMap := make(map[int]*WriterV2ColOption, elementType.NumField())
	for i := 0; i < elementType.NumField(); i++ {
		ft := elementType.Field(i)
		key := ft.Tag.Get("xlsx")
		col, ok := w.cols[key]
		if !ok {
			continue
		}
		optionMap[i] = col
	}

	var err error
	sliceValue := reflect.ValueOf(rows)
	for i := 0; i < sliceValue.Len(); i++ {
		rowIdx := w.markGetRow()
		elementValue := sliceValue.Index(i)
		for fieldIdx, col := range optionMap {
			axis := fmt.Sprintf("%s%d", col.col, rowIdx)
			value := col.Value(elementValue, elementValue.Field(fieldIdx))
			if err = w.file.SetCellValue(w.sheet, axis, value); err != nil {
				return err
			}
			if col.cellStyleID != 0 {
				_ = w.file.SetCellStyle(w.sheet, axis, axis, col.cellStyleID)
			}
		}
	}
	return nil
}

func (w *SheetWriterV2) writeHeader(rowIdx int) error {
	var err error
	for _, v := range w.cols {
		axis := fmt.Sprintf("%s%d", v.col, rowIdx)
		if err = w.file.SetCellStr(w.sheet, axis, v.header); err != nil {
			return err
		}
		if v.headerStyleID != 0 {
			if err = w.file.SetCellStyle(w.sheet, axis, axis, v.headerStyleID); err != nil {
				return err
			}
		}
		if v.headerComment != "" {
			if err = w.file.AddComment(w.sheet, excelize.Comment{
				Author: "System: ",
				Cell:   axis,
				Text:   v.headerComment,
			}); err != nil {
				return err
			}
		}
		if v.width > 0 {
			if err = w.file.SetColWidth(w.sheet, v.col, v.col, v.width); err != nil {
				return err
			}
		}
	}
	return nil
}
