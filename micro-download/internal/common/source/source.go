package source

import (
	"errors"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-download/internal/modes"
	"micro-download/internal/modules/download/delivery"
	"sync"
)

type source interface {
	// 创建任务
	Create(ctx ctx.Context, export delivery.CreateExport) (error, int64)
	// 任务概览
	Overview(ctx ctx.Context, detail *modes.DownloadSource) (*delivery.SourceOverview, interface{}, error)
	// 读取任务
	Read(ctx ctx.Context, page int64, title []string, params interface{}, value []string) (error, [][]any)
}

// factory 资源类型创建对象
func factory(sourceType int8) source {
	var once sync.Once
	var sourceData *dataSource
	var sourceSql *sqlSource
	var sourceRestful *restfulSource
	var sourcePack *packSource
	switch sourceType {
	case delivery.SourceTypeData:
		once.Do(func() {
			sourceData = &dataSource{}
		})
		return sourceData
	case delivery.SourceTypeSql:
		once.Do(func() {
			sourceSql = &sqlSource{}
		})
		return sourceSql
	case delivery.SourceTypeRestful:
		once.Do(func() {
			sourceRestful = &restfulSource{}
		})
		return sourceRestful
	case delivery.SourceTypePack:
		once.Do(func() {
			sourcePack = &packSource{}
		})
		return sourcePack
	}
	return nil
}

// SourceCreate 资源创建
func Create(ctx ctx.Context, export delivery.CreateExport) (error, int64) {
	object := factory(export.SourceType)
	if object != nil {
		return object.Create(ctx, export)
	}
	return errors.New("任务无法创建"), 0
}

// Overview 获取资源概览
func Overview(ctx ctx.Context, detail *modes.DownloadSource) (*delivery.SourceOverview, interface{}, error) {
	object := factory(detail.SourceType)
	if object != nil {
		return object.Overview(ctx, detail)
	}
	return nil, nil, errors.New(fmt.Sprintf("无法获取资源概览:%d", detail.ID))
}

// Read 资源读取
func Read(ctx ctx.Context, page, fileId int64, sourceType int8, title []string, params interface{}, value []string) (error, [][]any) {
	// 创建对象
	object := factory(sourceType)
	if object != nil {
		err, data := object.Read(ctx, page, title, params, value)
		return err, data
	}
	return errors.New(fmt.Sprintf("资源数据无法读取:%d", fileId)), [][]any{}
}
