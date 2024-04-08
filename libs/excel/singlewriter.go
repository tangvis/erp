package excel

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"reflect"
)

// SingleExcelSheetData
// @Description:
// PartOne struct
// PartTwo slice
type SingleExcelSheetData struct {
	Name    string
	PartOne interface{}
	PartTwo interface{}
}

// SingleExportExcel
// @Description: ✎ 导出excel
// 比较常用的一种导出内容是：详情 + 列表信息的方式
func SingleExportExcel(data SingleExcelSheetData) (*excelize.File, error) {
	f := excelize.NewFile()
	if data.Name == "" {
		data.Name = DefaultSheetName
	}
	if data.Name != DefaultSheetName {
		f.NewSheet(data.Name)
		f.DeleteSheet(DefaultSheetName)
	}
	f, rowIdx, err := writePartOne(data, f)
	if err != nil {
		return nil, err
	}
	return writePartTwo(data, f, rowIdx)
}

func writePartTwo(data SingleExcelSheetData, f *excelize.File, rowIdx int) (*excelize.File, error) {
	rt := reflect.TypeOf(data.PartTwo)
	if kind := rt.Kind(); kind != reflect.Slice {
		return nil, fmt.Errorf("sheet %s expect Slice, got %s instead", data.Name, kind)
	}
	options, err := generateOptionsByStruct(rt.Elem())
	if err != nil {
		return nil, err
	}
	for i := range options {
		options[i].headerStyle = DefaultHeaderStyleV2
		options[i].cellStyle = DefaultCellStyle
	}
	cols, err := convertOpt(options, f)
	if err != nil {
		return nil, err
	}
	nowRowIdx := rowIdx + 1
	res := &SheetWriterV2{
		cols:   cols,
		sheet:  data.Name,
		rowIdx: nowRowIdx,
		file:   f,
	}
	if err := res.writeHeader(nowRowIdx); err != nil {
		return nil, err
	}
	if err := writeRowBatch(data.PartTwo, res); err != nil {
		return nil, err
	}
	return f, nil
}

func writePartOne(data SingleExcelSheetData, f *excelize.File) (*excelize.File, int, error) {
	options, err := generateOptionsByStruct(reflect.TypeOf(data.PartOne))
	if err != nil {
		return nil, 0, err
	}
	v := reflect.ValueOf(data.PartOne)
	for i, opt := range options {
		axis := fmt.Sprintf("A%d", i+1)
		if err = f.SetCellStr(data.Name, axis, opt.header); err != nil {
			return nil, 0, err
		}
		headerStyleID, _ := f.NewStyle(DefaultHeaderStyle)
		cellStyleID, _ := f.NewStyle(DefaultCellStyle)
		if err = f.SetCellStyle(data.Name, axis, axis, headerStyleID); err != nil {
			return nil, 0, err
		}
		if opt.width > 0 {
			if err = f.SetColWidth(data.Name, "A", "B", opt.width); err != nil {
				return nil, 0, err
			}
		}
		valAxis := fmt.Sprintf("B%d", i+1)
		if err = f.SetCellStyle(data.Name, valAxis, valAxis, cellStyleID); err != nil {
			return nil, 0, err
		}
		if err := f.SetCellValue(data.Name, valAxis, v.Field(i).Interface()); err != nil {
			return nil, 0, err
		}
	}
	return f, len(options) + 1, nil
}
