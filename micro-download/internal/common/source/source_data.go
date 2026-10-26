package source

import (
	"encoding/json"
	"errors"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-download/internal/common"
	"micro-download/internal/modes"
	"micro-download/internal/modules/download/delivery"
)

// 数据类型
type dataSource struct{}

func (c *dataSource) Create(ctx ctx.Context, export delivery.CreateExport) (error, int64) {
	if export.SourceData == nil {
		return errors.New("无效的data导出参数"), 0
	}
	if len(export.SourceData) > 100 {
		return errors.New("data导出类型初衷是为了小型简单数据快速生成,大于100长度的列表数据清选择SQL或RESTFUL类型"), 0
	}
	model := modes.DownloadSource{}
	model.SetParams(export)
	return model.SourceCreate(ctx, func(d *modes.DownloadSource) (error, *modes.DownloadSource) {
		data := &delivery.SourceTypeDataData{
			ExcelSourceDataCommonOption: delivery.ExcelSourceDataCommonOption{
				SourceTitleAlias:     export.SourceTitleAlias,
				SourceTitle:          export.SourceTitle,
				SourceTitleValue:     export.SourceTitleValue,
				SourceSuffix:         export.SourceSuffix,
				SourceExcelCellMerge: export.SourceExcelCellMerge,
			},
			Data: export.SourceData,
		}
		bytes, err := json.Marshal(data)
		if err != nil {
			return err, d
		}
		d.SourceData = string(bytes)
		return nil, d
	})
}

func (c *dataSource) Overview(ctx ctx.Context, detail *modes.DownloadSource) (*delivery.SourceOverview, interface{}, error) {
	data := &delivery.SourceTypeDataData{}
	err := json.Unmarshal([]byte(detail.SourceData), data)
	if err != nil {
		return nil, nil, err
	}
	dataLength := int64(len(data.Data))
	result := &delivery.SourceOverview{
		Count:      dataLength,
		Title:      data.SourceTitle,
		TitleAlias: data.SourceTitleAlias,
		PageNum:    common.PageCount(dataLength),
		Suffix:     data.SourceSuffix,
		CellMerge:  data.SourceExcelCellMerge,
	}
	return result, data, err
}

// Read 读取
func (c *dataSource) Read(ctx ctx.Context, page int64, title []string, params interface{}, value []string) (error, [][]any) {
	var data [][]any
	requestParam := params.(*delivery.SourceTypeDataData)
	dataIndex := c.getTitleIndex(requestParam.Data[0], title)
	var dataSource [][]string
	var i = (page - 1) * common.DefaultSplitSize
	if i > 0 {
		i++
	}
	for ; i <= common.DefaultSplitSize; i++ {
		if i < int64(len(requestParam.Data)) {
			dataSource = append(dataSource, requestParam.Data[i])
		} else {
			break
		}
	}
	for lineIndex, line := range dataSource {
		var newLine []any
		for _, itemIndex := range dataIndex {
			newLine = append(newLine, line[itemIndex])
		}
		// 过滤默认表头数据
		if page == 1 && lineIndex == 0 {
			continue
		}
		data = append(data, newLine)
	}
	return nil, data
}

// getTitleIndex 返回标题对应的索引数据
func (c *dataSource) getTitleIndex(rawTitle, title []string) []int {
	var dataIndex []int
	for index, item := range title {
		if common.StringInArray(item, rawTitle) {
			dataIndex = append(dataIndex, index)
		}
	}
	return dataIndex
}
