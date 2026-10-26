package notification

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

func workbookBytes(t *testing.T, rows [][]interface{}) []byte {
	t.Helper()
	file := excelize.NewFile()
	for index, row := range rows {
		cell, _ := excelize.CoordinatesToCellName(1, index+1)
		if err := file.SetSheetRow("Sheet1", cell, &row); err != nil {
			t.Fatal(err)
		}
	}
	buffer, err := file.WriteToBuffer()
	if err != nil {
		t.Fatal(err)
	}
	return buffer.Bytes()
}

func TestParseCrowdWorkbook(t *testing.T) {
	data := workbookBytes(t, [][]interface{}{{"*门店code"}, {" OC001 "}, {""}, {"OC002"}, {"OC001"}})
	codes, err := parseCrowdWorkbook(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(codes, []string{"OC001", "OC002", "OC001"}) {
		t.Fatalf("codes = %#v", codes)
	}
}

func TestParseCrowdWorkbookRejectsWrongHeader(t *testing.T) {
	data := workbookBytes(t, [][]interface{}{{"门店编码"}, {"OC001"}})
	if _, err := parseCrowdWorkbook(bytes.NewReader(data)); err == nil {
		t.Fatal("expected invalid header error")
	}
}

func TestSummarizeCrowdImport(t *testing.T) {
	summary := summarizeCrowdImport(
		[]string{"OC001", "OC001", "OC002", "OC003"},
		[]string{"OC001", "OC003"},
	)
	if summary.Import != 4 || summary.Success != 2 || summary.Repeat != 1 || summary.Fail != 1 {
		t.Fatalf("summary = %#v", summary)
	}
	if !reflect.DeepEqual(summary.FailStoreCodes, []string{"OC002"}) {
		t.Fatalf("fail codes = %#v", summary.FailStoreCodes)
	}
}
