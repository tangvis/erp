package excel

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

var (
	rtReaderRawRow = reflect.TypeOf(ReaderRawRow{})
)

type ReaderRawRow struct {
	RowIdx  int
	Columns []string
}

type ReadCellValueFunc func(string) (interface{}, error)

func NewReadValueInt() ReadCellValueFunc {
	return func(s string) (interface{}, error) {
		return strconv.Atoi(s)
	}
}

func NewReadValueTime(loc *time.Location, layout string) ReadCellValueFunc {
	return func(s string) (interface{}, error) {
		return time.ParseInLocation(layout, s, loc)
	}
}

func NewReadValueTime2(loc *time.Location, layouts []string) ReadCellValueFunc {
	const LastSucessNo = -1
	lastSuccess := LastSucessNo
	errFMT := fmt.Sprintf("value '%%s' do not match any of the given layouts, must be one of: [%s]", strings.Join(layouts, ", "))
	return func(s string) (interface{}, error) {
		if lastSuccess != LastSucessNo {
			t, err := time.ParseInLocation(layouts[lastSuccess], s, loc)
			if err == nil {
				return t, nil
			}
		}
		for i, layout := range layouts {
			if i == lastSuccess {
				continue
			}
			t, err := time.ParseInLocation(layout, s, loc)
			if err != nil {
				continue
			}
			lastSuccess = i
			return t, nil
		}
		lastSuccess = LastSucessNo
		return nil, fmt.Errorf(errFMT, s)
	}
}

func NewReadValueTime3(loc *time.Location, layouts []string) ReadCellValueFunc {
	return func(s string) (interface{}, error) {
		for _, layout := range layouts {
			t, err := time.ParseInLocation(layout, s, loc)
			if err != nil {
				continue
			}
			return t, nil
		}
		return nil, errors.New("wrong Date Format")
	}
}

func NewReadValueLowerString() ReadCellValueFunc {
	return func(s string) (interface{}, error) {
		return strings.ToLower(s), nil
	}
}

func NewReadValueUpperString() ReadCellValueFunc {
	return func(s string) (interface{}, error) {
		return strings.ToUpper(s), nil
	}
}

func NewReadValueDecimal() ReadCellValueFunc {
	return func(s string) (interface{}, error) {
		return decimal.NewFromString(s)
	}
}
