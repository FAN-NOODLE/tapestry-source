package file

//
// import (
//	"fmt"
//	"github.com/360EntSecGroup-Skylar/excelize/v2"
//	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
//	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
//	"micro-download/internal/common"
//	"micro-download/internal/modules/download/delivery"
//	"time"
// )
//
//
// type ExcelGenerate struct {}
//
// // Generate 生成excel
// func (e *ExcelGenerate) Generate(ctx ctx.Context, filename string, startLine int, data [][]string, overview *delivery.SourceOverview, complete bool) error {
//	dataLength :=0
//	if len(data) > 0 {
//		dataLength = len(data[0])
//	}
//	log.Info(ctx).Msgf("")
//	// 没有就新建 有就读取
//	file, err := e.getFile(filename, dataLength)
//	if err != nil {
//		return err
//	}
//	// 表头
//	if startLine == 1 {
//		var alias []string
//		for _, item := range overview.Title {
//			alias = append(alias, overview.TitleAlias[item])
//		}
//		newData := [][]string{alias}
//		newData = append(newData, data...)
//		data = newData
//	}
//	writeFileStart := time.Now()
//	log.Info(ctx).Msgf("开始写入文件...")
//	// 写入
//	writeLine := startLine
//	for _, line := range data {
//		file.SetSheetRow("Sheet1", fmt.Sprintf("A%d", writeLine), &line)
//		writeLine ++
//	}
//	log.Info(ctx).Msgf("结束写入文件...耗时:%s", time.Since(writeFileStart))
//	log.Info(ctx).Msgf("开始合并单元格...")
//	mergeCellStart := time.Now()
//	// 合并单元格
//	//e.mergeCell(ctx, file, overview.CellMerge, startLine, data, overview.Title)
//	log.Info(ctx).Msgf("结束合并...耗时:%s", time.Since(mergeCellStart))
//	// 追加后缀
//	if complete {
//		suffixStart := startLine + len(data)
//		for _, item := range overview.Suffix {
//			file.SetSheetRow("Sheet1", fmt.Sprintf("A%d", suffixStart), &item)
//			suffixStart ++
//		}
//	}
//	return file.SaveAs(filename)
// }
//
// // getFile 获取文件
// func (e *ExcelGenerate) getFile(filename string, dataLength int) (*excelize.File, error) {
//	if common.FileExists(filename) {
//		return excelize.OpenFile(filename)
//	}
//	file := excelize.NewFile()
//	// 设置默认列宽
//	file.SetColWidth("Sheet1", "A", e.getFileFlag(dataLength), 25)
//	return file, nil
// }
//
// // 合并单元格
// func (e *ExcelGenerate) mergeCell(ctx ctx.Context, file *excelize.File, mergeOperation delivery.SourceExcelCellMerge, startLine int, data [][]string, title []string) {
//	var mergeCells [][]string
//	if mergeOperation.Type == delivery.ExcelMergeTypeColumn {
//		if mergeOperation.Rule == delivery.ExcelMergeRuleEqual {
//			// 取列索引
//			// 需要合并的列的索引
//			var mergeIndex []int
//			// 合并前提的列的索引
//			var premiseIndex []int
//			for index, item := range title {
//				if common.StringInArray(item, mergeOperation.Title) {
//					mergeIndex = append(mergeIndex, index)
//				}
//				if common.StringInArray(item, mergeOperation.Premise) {
//					premiseIndex = append(premiseIndex, index)
//				}
//			}
//			// 合并本页数据
//			// 上一行的数据
//			var prevRow []string
//			var mergeIndexLength = len(mergeIndex)
//			// 开始合并的cell
//			var startMergeCell = make(map[int]string, mergeIndexLength)
//			// 结束合并的cell
//			var endMergeCell = make(map[int]string, mergeIndexLength)
//			// 列是否可进入合并操作(这个不一定就是有合并的列，这个标识的意思是每次循环是否可以进进入到合并操作的那一步，也就是最后一部，如果没有合并数据也会走这一步)
//			var canMerge = make(map[int]bool, mergeIndexLength)
//			// 列是否存在已找到的合并列
//			var hasMergeData = make(map[int]bool, mergeIndexLength)
//			log.Info(ctx).Msgf("数据长度%d", len(data))
//			loopNum := 0
//			mergeNum := 0
//			for rowIndex, row := range data {
//				// 大于0 说明是在匹数据中
//				if rowIndex > 0 {
//					prevRow = data[rowIndex - 1]
//				} else {
//					// 等于0说明是在本批数据第一行  所以需要判断是否存在上一批(页)数据，与上一批数据最后一行（如果最后一行进行了合并需要取到最后进行合并行的起始行）进行合并
//					// todo 等于0的时候要进行分页merge
//				}
//				if len(prevRow) > 0 {
//					for _, mergeIndexItem := range mergeIndex {
//						// 判断合并前提
//						mergePremise := true
//						for _, premiseIndexItem := range premiseIndex {
//							if prevRow[premiseIndexItem] != row[premiseIndexItem] {
//								mergePremise = false
//								break
//							}
//							loopNum ++
//						}
//						// 判断合并
//						if mergePremise && prevRow[mergeIndexItem] == row[mergeIndexItem] {
//							if !hasMergeData[mergeIndexItem] {
//								startMergeCell[mergeIndexItem] = fmt.Sprintf("%s%d", e.getFileFlag(mergeIndexItem), rowIndex - 1 + startLine)
//								hasMergeData[mergeIndexItem] = true
//							}
//							endMergeCell[mergeIndexItem] = fmt.Sprintf("%s%d", e.getFileFlag(mergeIndexItem), rowIndex + startLine)
//							if rowIndex == len(data) - 1 {
//								canMerge[mergeIndexItem] = true
//							}
//						} else {
//							canMerge[mergeIndexItem] = true
//						}
//						// 进行合并
//						if canMerge[mergeIndexItem] {
//							//file.MergeCell("Sheet1", startMergeCell[mergeIndexItem], endMergeCell[mergeIndexItem])
//							if startMergeCell[mergeIndexItem] != "" && endMergeCell[mergeIndexItem] != "" {
//								mergeCells = append(mergeCells, []string{
//									"Sheet1",
//									startMergeCell[mergeIndexItem],
//									endMergeCell[mergeIndexItem],
//								})
//							}
//							canMerge[mergeIndexItem] = false
//							hasMergeData[mergeIndexItem] = false
//							mergeNum ++
//						}
//						loopNum ++
//					}
//				}
//			}
//			// 先合并一个
//			for _, item := range mergeCells {
//				if item[1] != "" && item[2] != "" {
//					file.MergeCell("Sheet1", item[1], item[2])
//				}
//			}
//			log.Info(ctx).Msgf("循环次数%d,合并次数%d", loopNum, mergeNum)
//		}
//	}
// }
//
// // getFileFlag 这是一个优雅的Excel序列计算
// // 支持范围 A - ZZ => 1 - 676
// func (e *ExcelGenerate) getFileFlag(cellIndex int) string {
//	// 65为unicode A
//	if cellIndex < 26 {
//		return fmt.Sprintf("%c", cellIndex + 65)
//	}
//	return fmt.Sprintf("%c%c", 65 + cellIndex / 26 - 1, cellIndex % 26 + 65)
// }
