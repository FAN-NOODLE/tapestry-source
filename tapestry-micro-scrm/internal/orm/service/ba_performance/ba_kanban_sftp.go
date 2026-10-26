package ba_performance

import (
	"time"

	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
)

func (ba *BaPerformanceService) GetBaInfoBYmarscode(enterpriseID int64, marsId string) (data micro_scrm_enterprise.MembersModel, err error) {
	err = ba.Scopes(model.BranchTable(&data, enterpriseID)).Select("name,mobile").Where("mars_id=?", marsId).Find(&data).Error
	return
}
func (ba *BaPerformanceService) BathCreateBaMonthlyKanbanModel(enterpriseID int64, infos []*micro_scrm_enterprise.BaMonthlyKanbanModel) error {
	return ba.Scopes(model.BranchTable(&micro_scrm_enterprise.BaMonthlyKanbanModel{}, enterpriseID)).Create(&infos).Error
}
func (ba *BaPerformanceService) FindBaMonthlyKan(enterpriseID int64, baName, mobile, month string, salesType int) (infos []micro_scrm_enterprise.BaMonthlyKanbanModel, err error) {
	err = ba.Scopes(model.BranchTable(&micro_scrm_enterprise.BaMonthlyKanbanModel{}, enterpriseID)).Where("month=? AND ba_name=? AND mobile=? AND sales_type=? AND deleted_at=?", month, baName, mobile, salesType, 0).Find(&infos).Error
	return
}
func (ba *BaPerformanceService) DeleteBaMonthlyKan(enterpriseID int64, baName, mobile, month string, salesType int) error {
	return ba.Scopes(model.BranchTable(&micro_scrm_enterprise.BaMonthlyKanbanModel{}, enterpriseID)).Where("month=? AND ba_name=? AND mobile=? AND sales_type=? AND deleted_at=?", month, baName, mobile, salesType, 0).Update("deleted_at", time.Now().Unix()).Error
}

func (ba *BaPerformanceService) UpoladBathCreateBaMonthlyKanbanModel(enterpriseID int64, infos []micro_scrm_enterprise.BaMonthlyKanbanModel) error {
	return ba.Scopes(model.BranchTable(&micro_scrm_enterprise.BaMonthlyKanbanModel{}, enterpriseID)).Create(&infos).Error
}
