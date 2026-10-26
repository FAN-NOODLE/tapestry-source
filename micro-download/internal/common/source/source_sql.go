package source

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-download/internal/common"
	"micro-download/internal/modes"
	"micro-download/internal/modules/download/delivery"
	"micro-download/internal/pkg/db"
	"strings"
)

// SQL类型
type sqlSource struct{}

func (c *sqlSource) Create(ctx ctx.Context, export delivery.CreateExport) (error, int64) {
	if export.SourceSql == "" || export.SourceSqlCount == 0 || export.SourceSqlDatabase == "" {
		return errors.New("无效的sql导出参数"), 0
	}
	// 检查数据库是否有效
	if db.DB(export.SourceSqlDatabase) == nil {
		return errors.New(fmt.Sprintf("无效的数据库%s", export.SourceSqlDatabase)), 0
	}
	// 检查Sql是否是select语句
	if !strings.HasPrefix(strings.ToLower(export.SourceSql), "select") {
		return errors.New("sql必须以select开头"), 0
	}
	// 如果sql末尾有分号,把分号去掉
	if strings.HasSuffix(export.SourceSql, ";") {
		export.SourceSql = export.SourceSql[0 : len(export.SourceSql)-1]
	}
	model := &modes.DownloadSource{}
	model.SetParams(export)
	return model.SourceCreate(ctx, func(d *modes.DownloadSource) (err error, downloadSource *modes.DownloadSource) {
		data := &delivery.SourceTypeSqlData{
			ExcelSourceDataCommonOption: delivery.ExcelSourceDataCommonOption{
				SourceTitleAlias:     export.SourceTitleAlias,
				SourceTitleValue:     export.SourceTitleValue,
				SourceTitle:          export.SourceTitle,
				SourceSuffix:         export.SourceSuffix,
				SourceExcelCellMerge: export.SourceExcelCellMerge,
			},
			Count:    export.SourceSqlCount,
			Sql:      export.SourceSql,
			Database: export.SourceSqlDatabase,
		}
		bytes, err := json.Marshal(data)
		if err != nil {
			return err, d
		}
		d.SourceData = string(bytes)
		return nil, d
	})
}

func (c *sqlSource) Overview(ctx ctx.Context, detail *modes.DownloadSource) (*delivery.SourceOverview, interface{}, error) {
	data := &delivery.SourceTypeSqlData{}
	err := json.Unmarshal([]byte(detail.SourceData), data)
	if err != nil {
		return nil, nil, err
	}
	result := &delivery.SourceOverview{
		Count:      data.Count,
		Title:      data.SourceTitle,
		TitleAlias: data.SourceTitleAlias,
		PageNum:    common.PageCount(data.Count),
		Suffix:     data.SourceSuffix,
		CellMerge:  data.SourceExcelCellMerge,
	}
	return result, data, nil
}

func (c *sqlSource) Read(ctx ctx.Context, page int64, title []string, params interface{}, value []string) (error, [][]any) {
	data := params.(*delivery.SourceTypeSqlData)
	sql := fmt.Sprintf("%s limit %d offset %d", data.Sql, common.DefaultSplitSize, common.DefaultSplitSize*(page-1))
	rows, err := db.DB(data.Database).Query(ctx, sql)
	if err != nil {
		return err, nil
	}
	columns, err := rows.Columns()
	if err != nil {
		return err, nil
	}
	// 明确请求参数中title的索引
	var titleInColumnsIndex []int
	for key, value := range columns {
		if common.StringInArray(value, title) {
			titleInColumnsIndex = append(titleInColumnsIndex, key)
		}
	}

	// 构造对应columns数量的指针供rows scan
	cache := make([]interface{}, len(columns))
	for i := range cache {
		var tmp interface{}
		cache[i] = &tmp
	}
	var list [][]any
	for rows.Next() {
		err := rows.Scan(cache...)
		if err != nil {
			return err, nil
		}
		var item []string
		for _, d := range cache {
			item = append(item, fmt.Sprintf("%s", *(d.(*interface{}))))
		}
		if item != nil {
			// 只取title中存在的数据列
			var listItem []any
			for _, columnsIndex := range titleInColumnsIndex {
				listItem = append(listItem, item[columnsIndex])
			}
			if listItem != nil {
				list = append(list, listItem)
			}
		}
	}
	return nil, list
}
