package helper

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"io"
	"micro-scrm/internal/app"
	"micro-scrm/internal/pkg/http"
	"strings"
)

// ReadRemoteUrlData 读取远程Excel文件
func ReadRemoteUrlData(c app.Ctx, fileUrl string, readLine int) ([]map[string]any, error) {
	client := http.NewHTTPClient()
	resp, err := client.Get(c, fileUrl)
	if err != nil {
		return nil, err
	}

	respByte, err := io.ReadAll(resp.Response().Body)
	if err != nil {
		return nil, err
	}

	urlArr := strings.Split(fileUrl, ".")
	ext := urlArr[len(urlArr)-1]

	if !InArray(ext, []string{"xlsx", "xls"}) {
		return nil, fmt.Errorf("不支持 [%s] 类型", ext)
	}

	// 解析 Excel 文件内容到内存中
	file, err := xlsx.OpenBinary(respByte)
	if err != nil {
		return nil, err
	}

	// 提取第一个工作簿（worksheet）中的数据
	sheet := file.Sheets[0]
	rows := sheet.Rows

	if len(rows) < readLine {
		return nil, fmt.Errorf("行号超出范围")
	}

	// 读取标头行和数据行
	headerRow := rows[readLine]
	dataRows := rows[readLine+1:]

	if len(dataRows) < 0 {
		return nil, nil
	}

	// 将每一行的数据解析为 map[string]interface{} 类型
	result := make([]map[string]any, len(dataRows))
	for i, row := range dataRows {
		data := make(map[string]any)
		for j, cell := range row.Cells {
			key := headerRow.Cells[j].Value
			value, _ := cell.FormattedValue()
			if value != "" {
				data[key] = value
			}
		}
		result[i] = data
	}

	return result, nil
}
