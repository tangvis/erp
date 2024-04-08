package excel

import (
	"os"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

type TestStruct struct {
	Raw        ReaderRawRow
	KeyInt     int             `xlsx:"KeyInt;required"`
	KeyDecimal decimal.Decimal `xlsx:"decimal"`
	KeyString  string          `xlsx:"KeyString;required"`
	KeyTime    time.Time       `xlsx:"KeyTime;required"`
	KeyTime2   time.Time       `xlsx:"KeyTime2"`
}

func TestReadExcel(t *testing.T) {
	fs, err := os.Open("./test_read.xlsx")
	if err != nil {
		t.Fatal(err)
	}
	defer fs.Close()
	reader, err := NewReaderV2(fs)
	if err != nil {
		t.Fatal(err)
	}
	assert := require.New(t)
	var dest []TestStruct
	loc, _ := time.LoadLocation("UTC")
	// required
	sheetReader, err := reader.NewSheetReader("Sheet1")
	if err != nil {
		t.Fatal(err)
	}
	err = sheetReader.ReadSome(&dest, 1)
	t.Logf("dest is %+v, err is %+v", dest, err)
	assert.NotNil(err, "must not nil")
	assert.NotNil(dest, "must not nil")
	t.Logf("[Debug]required err is %s", err)

	// normal read but time failed
	sheetReader, err = reader.NewSheetReader(
		"Sheet1",
	)
	if err != nil {
		t.Fatal(err)
	}
	err = sheetReader.ReadSome(&dest, 1)
	assert.NotNil(err, "must not nil")
	assert.NotNil(dest, "must not nil")
	assert.Equal(0, len(dest), "must not nil")
	t.Logf("[Debug]read err is %s", err)
	// read normal
	sheetReader, err = reader.NewSheetReader(
		"Sheet1",
		SheetReaderWithCustomReadFunc("KeyTime", NewReadValueTime(loc, "2006-01-02")),
	)
	if err != nil {
		t.Fatal(err)
	}
	err = sheetReader.ReadSome(&dest, 1)
	assert.Nil(err, "must nil")
	assert.Equal(1, dest[0].Raw.RowIdx)
	assert.NotEqual(0, len(dest[0].Raw.Columns))
	assert.Equal(1, len(dest), "must same")
	assert.Equal("Tw1", dest[0].KeyString, "must same")
	assert.Equal(1015, dest[0].KeyInt, "must same")
	assert.Equal(int64(1636588800), dest[0].KeyTime.Unix(), "must same")
	// read all
	err = sheetReader.ReadAll(&dest)
	t.Logf("dest is %+v", dest)
	assert.Nil(err, "must nil")
	assert.Equal(2, len(dest), "must same")
}

func TestEasyReadExcel(t *testing.T) {
	fs, err := os.Open("./test_read.xlsx")
	if err != nil {
		t.Fatal(err)
	}
	defer fs.Close()
	var dest []TestStruct
	err = EasyRead(
		fs,
		&dest,
		SheetReaderWithCustomReadFunc("KeyTime", NewReadValueTime(time.Local, "2006-01-02")),
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("dest is %+v\n", dest)
}
