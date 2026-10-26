package ws_file

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

// 常量
const (
	// DefaultCreatSheetName 默认创建的sheet的名称
	DefaultCreatSheetName = "Sheet1"
)

type excelFile struct {
	File       *excelize.File
	titleName  []string // 表头内容
	SheetName  string
	index      int
	TitleStyle int
	DataStyle  int
	Row        int
}

var CollFlagList = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "G", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

func NewExcelFile(sheetName string) (file *excelFile, err error) {
	file.SheetName = sheetName
	file.File = excelize.NewFile()
	file.index = file.File.NewSheet(sheetName)

	if err = file.File.SetColWidth(file.SheetName, "A", "G", 40); err != nil {
		return
	}

	if file.TitleStyle, err = file.File.NewStyle(`{"alignment":{"horizontal":"center","Vertical":"center"},"font":{"bold":true}}`); err != nil {
		return
	}

	if err = file.File.SetCellStyle(sheetName, "A1", "G1", file.TitleStyle); err != nil {
		return
	}

	if file.DataStyle, err = file.File.NewStyle(`{"alignment":{"horizontal":"center","Vertical":"center"}}`); err != nil {
		return
	}

	return
}

// SetTitleName 设置excel的title
func (file *excelFile) SetTitleName(titleList []string) *excelFile {
	for index, title := range titleList {
		_ = file.File.SetCellStr(file.SheetName, fmt.Sprintf("%s%d", CollFlagList[index], file.Row), title)
	}

	file.Row++

	return file
}

// FillData 设置数据
//func (file *excelFile) FillData() *excelFile {
//
//	return file
//}

// Save 设置sheet
func (file *excelFile) Save() *excelFile {
	file.File.SetActiveSheet(file.index)
	return file
}
