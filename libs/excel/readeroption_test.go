package excel

import (
	"reflect"
	"testing"

	"github.com/shopspring/decimal"
)

func TestSheetOptionAutoSet(t *testing.T) {
	var TestStruct struct {
		Int     int
		Float64 float64
		Uint    uint
		String  string
		Decimal decimal.Decimal
	}
	option := &sheetReaderColOption{}
	datas := []string{
		"-12",
		"12.4545",
		"18",
		"哈哈",
		"123.444",
	}
	rv := reflect.ValueOf(&TestStruct).Elem()
	for i := 0; i < rv.NumField(); i++ {
		if err := option.autoSet(datas[i], rv.Field(i)); err != nil {
			t.Fatal(err)
		}
	}
	t.Logf("result is %+v", TestStruct)
}
