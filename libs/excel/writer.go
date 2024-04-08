package excel

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type SheetData struct {
	Name string
	Data any
}

type excelSupportedOption struct {
	Header        string
	HeaderSuffix  string
	HeaderComment string
	Key           string
	Width         float64
}

func (o excelSupportedOption) FinalHeader() string {
	if o.HeaderSuffix != "" {
		return fmt.Sprintf("%s %s", o.Header, o.HeaderSuffix)
	}
	return o.Header
}

func GenerateOptionsByStruct(data any) ([]WriterV2ColOption, error) {
	return generateOptionsByStruct(reflect.TypeOf(data))
}

func generateOptionsByStruct(rt reflect.Type) ([]WriterV2ColOption, error) {
	// 禁止传什么鬼指针
	if kind := rt.Kind(); kind != reflect.Struct {
		return nil, fmt.Errorf("expect Struct, got %s instead", kind)
	}
	res := make([]WriterV2ColOption, 0, rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		xlsxTag := field.Tag.Get("xlsx")
		if xlsxTag == "" || xlsxTag == "-" {
			continue
		}
		option, err := explainTagToOption(field, xlsxTag)
		if err != nil {
			return nil, err
		}
		res = append(res, option)
	}
	return res, nil
}

func explainTagToOption(
	field reflect.StructField,
	tag string,
) (WriterV2ColOption, error) {

	var emptyResult WriterV2ColOption
	var supportedOption excelSupportedOption
	defines := strings.Split(tag, ";")
	if defines[0] == "" {
		return emptyResult, fmt.Errorf("empty tag at Field %s", field.Name)
	}
	// 设置K初始值
	supportedOption.Key = tag
	supportedOption.Header = defines[0]
	// 设置其它配置
	for i := 1; i < len(defines); i++ {
		tuples := strings.Split(defines[i], "=")
		switch tuples[0] {
		case "width":
			val, err := strconv.ParseFloat(tuples[1], 64)
			if err != nil {
				return emptyResult, fmt.Errorf("wrong conf at tuple(%s) at Field %s: %s", defines[i], field.Name, err)
			}
			supportedOption.Width = val
		case "comment":
			supportedOption.HeaderComment = tuples[1]
		default:
			return emptyResult, fmt.Errorf("unsupported conf(%s) at Field %s", tuples[0], field.Name)
		}
	}
	res := NewWriterV2ColOptionSimple3(
		supportedOption.FinalHeader(),
		supportedOption.Key,
		supportedOption.Width,
		supportedOption.HeaderComment,
	)
	return res, nil
}

func ExportBatchSheet(
	ctx context.Context,
	sheets []SheetData,
) (*WriterV2, error) {
	firstSheet := sheets[0]

	rt := reflect.TypeOf(firstSheet.Data)
	if kind := rt.Kind(); kind != reflect.Slice {
		return nil, fmt.Errorf("sheet 0 expect Slice, got %s instead", kind)
	}
	options, err := generateOptionsByStruct(rt.Elem())
	if err != nil {
		return nil, err
	}
	writer, err := NewWriterV2(firstSheet.Name, options)
	if err != nil {
		return nil, err
	}
	for i := 1; i < len(sheets); i++ {
		rt = reflect.TypeOf(sheets[i].Data)
		if kind := rt.Kind(); kind != reflect.Slice {
			return nil, fmt.Errorf("sheet %d expect Slice, got %s instead", i, kind)
		}
		options, err = generateOptionsByStruct(rt.Elem())
		if err != nil {
			return nil, err
		}
		if _, err = writer.AddSheet(sheets[i].Name, options); err != nil {
			return nil, err
		}
	}
	for i, v := range sheets {
		if err = writer.Sheet(v.Name).WriteRowBatch(ctx, v.Data); err != nil {
			return nil, fmt.Errorf("sheet %d write error: %s", i, err)
		}
	}
	return writer, nil
}

func ExportSingleSheet(
	ctx context.Context,
	dataSlice any,
) (*WriterV2, error) {
	sheets := []SheetData{
		{Name: DefaultSheetName, Data: dataSlice},
	}
	return ExportBatchSheet(ctx, sheets)
}
