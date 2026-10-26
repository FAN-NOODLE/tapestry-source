package ba_performance

import (
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/orm/service"
	"time"
)

//GetPageList BA 业绩指标分页列表
func (g *BaPerformanceService) GetPageList(eid int64, isPage, page, size int, userIds []string, name string) (data []micro_scrm_enterprise.BaPerformanceIndicatorModel, total int64, err error) {
	db := g.Scopes(model.BranchTable(&micro_scrm_enterprise.BaPerformanceIndicatorModel{}, eid)).
		Scopes(service.ScopeCond_EQDelete).
		Where("month = ?", time.Now().Format("200601"))
	if name != "" {
		db = db.Where("ba_name like ?", "%"+name+"%")
	}
	if len(userIds) != 0 {
		db = db.Where("mobile in ?", userIds)
	}
	if isPage == 1 {
		err = db.Count(&total).Limit(size).Offset((page - 1) * size).Find(&data).Error
	} else {
		err = db.Find(&data).Error
	}
	return
}

//Create 支持批量新增
func (g *BaPerformanceService) Create(eid int64, data []micro_scrm_enterprise.BaPerformanceIndicatorModel) (err error) {
	err = g.Scopes(model.BranchTable(&micro_scrm_enterprise.BaPerformanceIndicatorModel{}, eid)).Create(&data).Error
	return
}

//AddLog 导入记录
func (g *BaPerformanceService) AddLog(eid int64, data []micro_scrm_enterprise.BaPerformanceIndicatorModel) (err error) {
	err = g.Scopes(model.BranchTable(&micro_scrm_enterprise.BaPerformanceIndicatorLogModel{}, eid)).Create(&data).Error
	return
}

//Delete 批量删除
func (g *BaPerformanceService) Delete(eid int64) (err error) {
	err = g.Scopes(model.BranchTable(&micro_scrm_enterprise.BaPerformanceIndicatorModel{}, eid)).
		Where("month = ?", time.Now().Format("200601")).
		Update("deleted_at", time.Now().Unix()).Error
	return
}
