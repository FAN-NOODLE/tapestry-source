package file

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/xuri/excelize/v2"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-download/internal/common"
	"micro-download/internal/modules/download/delivery"
)

var filePool sync.Map
var defaultSheetName = "Sheet1"
var defaultRowHeight = 28.3464567

type customFile struct {
	file *excelize.File
	*excelize.StreamWriter
	styles     sync.Map
	rowIndex   *int64
	dataLength int
}

func (c *customFile) Close() error {
	return c.file.Close()
}

func (c *customFile) MergeLine() error {

	startCell, err := excelize.CoordinatesToCellName(1, int(atomic.LoadInt64(c.rowIndex)))
	if err != nil {
		return err
	}
	overCell, err := excelize.CoordinatesToCellName(c.dataLength, int(atomic.LoadInt64(c.rowIndex)))
	if err != nil {
		return err
	}
	return c.MergeCell(startCell, overCell)
}

// 添加一行信息，不用再维护标志位
func (c *customFile) AddRow(values []interface{}, opts ...excelize.RowOpts) error {
	if len(opts) == 0 {
		opts = append(opts, excelize.RowOpts{
			Height: defaultRowHeight,
		})
	}
	r := atomic.AddInt64(c.rowIndex, 1)
	cell := fmt.Sprintf("A%d", r)
	return c.SetRow(cell, values, opts...)
}

func (c *customFile) SetStyle(styleName string, style *excelize.Style) (int, error) {
	if sid, ok := c.styles.Load(styleName); ok {
		return sid.(int), nil
	}
	sid, err := c.file.NewStyle(style)
	if err != nil {
		return 0, err
	}
	c.styles.Store(styleName, sid)
	return sid, nil
}
func (c *customFile) GetStyleID(styleName string) int {
	if sid, ok := c.styles.Load(styleName); ok {
		return sid.(int)
	}
	return 0
}

type XlsxGenerate struct {
}

type cellData struct {
	Value  any
	HMerge int
	VMerge int
}

func (x *XlsxGenerate) Generate(ctx ctx.Context, filename string, startLine int, data [][]any, overview *delivery.SourceOverview, complete bool) error {
	// filename = "test128.xlsx"
	dataLength := len(overview.Title)
	log.Info(ctx).Msgf("获取文件句柄")
	// 没有就新建 有就读取
	file, err := x.getFile(ctx, filename, dataLength)
	if err != nil {
		filePool.Delete(filename)
		file.Close()
		return err
	}
	// 表头
	if startLine == 1 {
		var alias []any
		for _, item := range overview.Title {
			alias = append(alias, overview.TitleAlias[item])
		}
		newData := [][]any{alias}
		newData = append(newData, data...)
		data = newData
	}
	// 计算单元格数据
	startGetMergeDataTime := time.Now()
	log.Info(ctx).Msgf("开始计算合并数据...")
	insertData := x.getMergeData(ctx, overview.CellMerge, data, overview.Title)
	log.Info(ctx).Msgf("合并数据计算结束，耗时%s", time.Since(startGetMergeDataTime))
	startWriteDataTime := time.Now()
	log.Info(ctx).Msgf("开始写入数据")

	if err != nil {
		filePool.Delete(filename)
		file.Close()
		return err
	}

	if err != nil {
		filePool.Delete(filename)
		file.Close()
		return err
	}
	defaultStyle := file.GetStyleID("defaultStyle")
	titleStyle := file.GetStyleID("titleStyle")
	// 判断是否有标题前缀合并的需求
	if len(overview.HeadPrefix) > 0 && startLine == 1 {
		for _, prefixOne := range overview.HeadPrefix {
			cells := make([]interface{}, 0, len(prefixOne))
			for i := 0; i < len(prefixOne); i++ {
				cell := &excelize.Cell{
					Value:   prefixOne[i],
					StyleID: defaultStyle,
				}
				cells = append(cells, cell)
			}
			file.AddRow(cells)
			err := file.MergeLine()
			if err != nil {
				panic(err)
			}
		}
	}

	for rowIndex, rowItem := range insertData {
		cells := make([]interface{}, 0, len(rowItem))
		for _, cellItem := range rowItem {
			stid := defaultStyle
			if startLine == 1 && rowIndex == 0 {
				stid = titleStyle
			}
			cells = append(cells, &excelize.Cell{
				StyleID: stid,
				Value:   cellItem.Value,
			})
		}
		file.AddRow(cells)
	}
	log.Info(ctx).Msgf("结束写入数据,耗时%s", time.Since(startWriteDataTime))
	// 添加后缀
	if complete {
		for _, rowItem := range overview.Suffix {
			cells := make([]interface{}, 0, len(rowItem))
			for _, cellItem := range rowItem {
				cells = append(cells, &excelize.Cell{
					Value:   cellItem,
					StyleID: defaultStyle,
				})
			}
			file.AddRow(cells)
		}
		startSaveFileTime := time.Now()
		log.Info(ctx).Msgf("开始保存文件...")
		err = file.Flush()
		if err != nil {
			log.Error(ctx).Msgf("保存文件失败...%v", err)
			filePool.Delete(filename)
			file.Close()
			return err
		}
		err = file.file.SaveAs(filename)
		if err != nil {
			log.Error(ctx).Msgf("保存文件失败...%v", err)
		}
		log.Info(ctx).Msgf("保存文件耗时%s", time.Since(startSaveFileTime))
		filePool.Delete(filename)
		file.Close()
	}
	return err
}

func (x *XlsxGenerate) getMergeData(ctx ctx.Context, mergeOperation delivery.SourceExcelCellMerge, data [][]any, title []string) [][]*cellData {
	var mergeCells [][]*cellData
	// 列合并
	// if mergeOperation.Type == delivery.ExcelMergeTypeColumn {
	// 相等合并
	// if mergeOperation.Rule == delivery.ExcelMergeRuleEqual {
	// 获取需要合并的列的索引
	var mergeIndex []int
	// 获取合并前提的列的索引
	var premiseIndex []int
	for index, item := range title {
		if common.StringInArray(item, mergeOperation.Title) {
			mergeIndex = append(mergeIndex, index)
		}
		if common.StringInArray(item, mergeOperation.Premise) {
			premiseIndex = append(premiseIndex, index)
		}
	}

	// 下一行的数据
	var nextRow []any
	var continueRowIndex []int
	var continueRows [][]any
	for rowIndex, row := range data {
		if common.IntInArray(rowIndex, continueRowIndex) {
			continue
		}
		vMerge := 0
		nextRowIndex := rowIndex + 1

		if premiseIndex != nil {
		startCheck:
			if nextRowIndex < len(data)+1 {
				if nextRowIndex != len(data) {
					nextRow = data[nextRowIndex]
				} else {
					nextRow = nil
				}
			} else {
				break
			}
			// 判断本行是否满足合并前提
			rowIsSatisfy := true
			if len(nextRow) > 0 {
				for _, premiseRowIndex := range premiseIndex {
					if row[premiseRowIndex] != nextRow[premiseRowIndex] {
						rowIsSatisfy = false
						break
					}
				}
			} else {
				rowIsSatisfy = false
			}

			// 是否满足合并前提
			if rowIsSatisfy {
				continueRowIndex = append(continueRowIndex, nextRowIndex)
				continueRows = append(continueRows, nextRow)
				vMerge++
				nextRowIndex++
				goto startCheck
			}
		}

		// 追加主行数据
		var rowItem []*cellData
		for cellIndex, cell := range row {
			cellItem := &cellData{
				Value: cell,
			}
			// 这里之前是根据前提条件算出来垂直合并值，然后过滤追加给合并字段
			if common.IntInArray(cellIndex, mergeIndex) {
				cellItem.VMerge = vMerge
			}
			// if mergeOperation.PremiseMap == nil {
			//	if common.IntInArray(cellIndex, mergeIndex) {
			//		cellItem.VMerge = vMerge
			//	}
			// } else {
			//	if common.IntInArray(cellIndex, premiseIndex) {
			//		cellItem.VMerge = vMerge
			//	}
			// }

			rowItem = append(rowItem, cellItem)
		}
		// mergeCells = append(mergeCells, rowItem)
		var mergeRows [][]*cellData
		mergeRows = append(mergeRows, rowItem)
		// 追加合并行数据
		for _, mergeRow := range continueRows {
			var mergeRowItem []*cellData
			for _, cell := range mergeRow {
				mergeRowItem = append(mergeRowItem, &cellData{
					Value: cell,
				})
			}
			// mergeCells = append(mergeCells, mergeRowItem)
			mergeRows = append(mergeRows, mergeRowItem)
		}
		// 新增二次合并逻辑(有二次规则则进行二次合并)
		// ruleMergeMap := map[string]interface{}{//规则测试示例值
		//	"order_no":[]string{
		//		"order_total_amount",
		//	},
		// }
		ruleMergeMap := mergeOperation.PremiseMap
		if ruleMergeMap != nil {
			// 二次加工合并数据
			mergeRecursive(ruleMergeMap, mergeRows, mergeIndex, premiseIndex, title) // 切片长度不变为引用传递不需要返回值
		}
		mergeCells = append(mergeCells, mergeRows...)
		// 清楚标记
		vMerge = 0
		continueRows = nil
	}
	//	}
	// }
	return mergeCells
}

// 扩展新增规则
// item=>interface{}.(type) = 无限极map (map[string]interface{})
// 递归合并--新增规则合并算法
func mergeRecursive(mergeRuleMap map[string]interface{}, rowData [][]*cellData, mergeIndex, premiseIndex []int, title []string) {
	// var mergeCells [][]*cellData
	for key, val := range mergeRuleMap {
		switch v := val.(type) {
		case []string, []interface{}:
			// todo merge logic here.
			// key = merge cell title name.
			totalLen := len(rowData)
			if totalLen <= 1 {
				return
			}
			var (
				cellKeyIndex      = -0xff
				offsetIndex       = 0
				cellMergeIndexArr []int
				temp              []string
			)
			if vType, ok := v.([]interface{}); ok {
				for _, val := range vType {
					if vstr, ok := val.(string); ok {
						temp = append(temp, vstr)
					}
				}
			} else {
				temp = v.([]string)
			}
			for index, item := range title {
				if item == key {
					cellKeyIndex = index
				}
				if common.StringInArray(item, temp) {
					cellMergeIndexArr = append(cellMergeIndexArr, index)
				}
			}
			if cellKeyIndex < 0 {
				return
			}
		recheck:
			nextRowIndex := offsetIndex + 1
			rIndex := offsetIndex
			for i := offsetIndex; i < totalLen-1; i++ {
				if rowData[rIndex][cellKeyIndex].Value == rowData[nextRowIndex][cellKeyIndex].Value {
					if !common.IntInArray(cellKeyIndex, mergeIndex) {
						rowData[rIndex][cellKeyIndex].VMerge++
					}
					for _, cellMergeIndex := range cellMergeIndexArr {
						if cellMergeIndex == cellKeyIndex || common.IntInArray(cellKeyIndex, mergeIndex) {
							continue
						}
						rowData[rIndex][cellMergeIndex].VMerge++
					}
					nextRowIndex++
				} else {
					break
				}
			}
			if nextRowIndex != totalLen-1 && offsetIndex != 0 {
				offsetIndex = nextRowIndex
				goto recheck
			}

		case map[string]interface{}:
			mergeRecursive(v, rowData, mergeIndex, premiseIndex, title)
		default:
			break
		}
	}
}

func (x *XlsxGenerate) getFile(ctx ctx.Context, filename string, dataLength int) (*customFile, error) {
	object, ok := filePool.Load(filename)
	if ok {
		return object.(*customFile), nil
	}
	file := excelize.NewFile()
	fstream, err := file.NewStreamWriter(defaultSheetName)
	if err != nil {
		log.Error(ctx).Msgf("[get file] fail: %v", err)
		return nil, err
	}

	// 设置列表宽度为25
	endCol := dataLength
	if endCol < 1 {
		endCol = 1
	}
	err = fstream.SetColWidth(1, endCol, 40)
	if err != nil {
		log.Error(ctx).Msgf("[set col width] fail: %v", err)
		return nil, err
	}
	cf := &customFile{
		file:         file,
		StreamWriter: fstream,
		styles:       sync.Map{},
		dataLength:   dataLength,
		rowIndex:     new(int64),
	}
	border := make([]excelize.Border, 0, 4)
	border = append(border,
		excelize.Border{Type: "bottom", Color: "#000000", Style: 7},
		excelize.Border{Type: "right", Color: "#000000", Style: 7},
		excelize.Border{Type: "top", Color: "#000000", Style: 7},
		excelize.Border{Type: "left", Color: "#000000", Style: 7})
	_, err = cf.SetStyle("defaultStyle", &excelize.Style{
		Border: border,
		Font:   &excelize.Font{Size: 12, Family: "Verdana"},
	})
	if err != nil {
		cf.Close()
		return nil, err
	}
	_, err = cf.SetStyle("titleStyle", &excelize.Style{
		Border: border,
		Font:   &excelize.Font{Size: 12, Bold: true, Family: "Verdana"},
	})
	if err != nil {
		cf.Close()
		return nil, err
	}
	filePool.Store(filename, cf)
	return cf, nil
}

func (x *XlsxGenerate) getFileFlag(value int) string {
	if value < 0 {
		return ""
	}
	var ans string
	i := value + 1
	for i > 0 {
		ans = string((i-1)%26+65) + ans
		i = (i - 1) / 26
	}
	return ans
}
