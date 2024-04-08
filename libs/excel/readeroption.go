package excel

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/shopspring/decimal"
)

var (
	reflectTypeDecimal = reflect.TypeOf(decimal.Decimal{})
)

// sheetReaderColOption 给外部的
type sheetReaderColOption struct {
	idx        int
	header     string
	headerText string
	required   bool
	readFunc   ReadCellValueFunc
}

func newSheetReaderColOption(
	header string,
	headerText string,
	required bool,
	readFunc ReadCellValueFunc,
) sheetReaderColOption {
	return sheetReaderColOption{
		header:     header,
		headerText: headerText,
		required:   required,
		readFunc:   readFunc,
	}
}

func (s *sheetReaderColOption) setValue(data string, rv reflect.Value) error {
	if data == "" {
		if s.required {
			return fmt.Errorf("%s is required", s.headerText)
		}
		return nil
	}
	if s.readFunc == nil {
		return s.autoSet(data, rv)
	}
	v, err := s.readFunc(data)
	if err != nil {
		return err
	}
	rv.Set(reflect.ValueOf(v))
	return nil
}

func (s *sheetReaderColOption) autoSet(data string, rv reflect.Value) error {
	switch kind := rv.Kind(); kind {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		val, err := strconv.ParseInt(data, 10, 64)
		if err != nil {
			return fmt.Errorf("%s: %s is not a valid Integer", s.headerText, data)
		}
		rv.SetInt(val)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		val, err := strconv.ParseUint(data, 10, 64)
		if err != nil {
			return fmt.Errorf("%s: %s is not a valid Integer", s.headerText, data)
		}
		rv.SetUint(val)
	case reflect.Float32, reflect.Float64:
		val, err := strconv.ParseFloat(data, 64)
		if err != nil {
			return fmt.Errorf("%s: %s is not a valid Number", s.headerText, data)
		}
		rv.SetFloat(val)
	case reflect.String:
		rv.SetString(data)
	case reflect.Struct:
		return s.autoSetStruct(data, rv)
	default:
		return fmt.Errorf("unsupported kind %s", kind)
	}
	return nil
}

func (s *sheetReaderColOption) autoSetStruct(data string, rv reflect.Value) error {
	switch typ := rv.Type(); typ {
	case reflectTypeDecimal:
		val, err := decimal.NewFromString(data)
		if err != nil {
			return fmt.Errorf("%s: %s is not a valid Number", s.headerText, data)
		}
		rv.Set(reflect.ValueOf(val))
	default:
		return fmt.Errorf("unsupported struct type %s", typ)
	}
	return nil
}
