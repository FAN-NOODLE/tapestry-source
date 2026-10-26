package activityQRCodeDao

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	activityQRCodeReq "micro-scrm/internal/app/http/request/activity_qrcode"
	activityQRCodeModel "micro-scrm/internal/app/models/activity_qrcode"
	"micro-scrm/internal/pkg/helper"
)

type MainDao struct {
	dao.EnterpriseRepo[activityQRCodeModel.ActivityQRCodeMainModel]
	model     *activityQRCodeModel.ActivityQRCodeMainModel
	tableName string
}

func NewMainDao() *MainDao {
	return &MainDao{
		EnterpriseRepo: dao.EnterpriseRepo[activityQRCodeModel.ActivityQRCodeMainModel]{
			Dao: dao.Dao{DB: app.DB()},
		},
		model:     new(activityQRCodeModel.ActivityQRCodeMainModel),
		tableName: activityQRCodeModel.ActivityQRCodeMainModel{}.TableName(),
	}
}

func (d *MainDao) GetActivityList(c ctx.Context, req activityQRCodeReq.ActivityListReq) ([]activityQRCodeModel.ActivityQRCodeMainModel, int64, error) {
	var (
		result []activityQRCodeModel.ActivityQRCodeMainModel
		total  int64
	)
	query := d.Dao.WithContext(c).Scopes(d.Dao.BranchTable(d.tableName, req.EnterpriseId)).Model(d.model)
	if req.CategoryID > 0 {
		query.Where("category_ids like ? ", "%"+helper.Int64ToString(req.CategoryID)+"%")
	}
	if req.EffectTime.StartTime > 0 {
		query.Where("(end_time >= ? or time_type = ?)", req.EffectTime.StartTime, activityQRCodeModel.TimeTypeLongTerm)
	}
	if req.EffectTime.EndTime > 0 {
		query.Where("(start_time <= ? or time_type = ?)", req.EffectTime.EndTime, activityQRCodeModel.TimeTypeLongTerm)
	}
	if req.CreateTime.StartTime > 0 {
		query.Where("created_at >= ? ", req.CreateTime.StartTime)
	}
	if req.CreateTime.EndTime > 0 {
		query.Where("created_at <= ? ", req.CreateTime.EndTime)
	}
	if req.Name != "" {
		query.Where("name like ? ", "%"+req.Name+"%")
	}
	offset := (req.Page - 1) * req.Size
	err := query.Count(&total).
		Order("activity_qrcode_main.created_at desc").
		Offset(offset).Limit(req.Size).
		Find(&result).Error
	return result, total, err
}

func (d *MainDao) GetStatisticsList(c ctx.Context, req activityQRCodeReq.StatisticsReq) ([]activityQRCodeModel.ActivityQRCodeMainModel, error) {
	var list []activityQRCodeModel.ActivityQRCodeMainModel
	query := d.Dao.WithContext(c).Scopes(d.Dao.BranchTable(d.tableName, req.EnterpriseId)).Model(d.model)
	if req.StartTime > 0 {
		query.Where("(end_time >= ? or time_type = ?)", req.StartTime, activityQRCodeModel.TimeTypeLongTerm)
	}
	if req.EndTime > 0 {
		query.Where("(start_time <= ? or time_type = ?)", req.EndTime, activityQRCodeModel.TimeTypeLongTerm)
	}
	if req.NameStr != "" {
		query.Where("name like ? ", "%"+req.NameStr+"%")
	}
	err := query.Find(&list).Error
	return list, err
}
