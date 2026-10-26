package workExternalSvc

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao/micro_scrm"
	workExternalModel "micro-scrm/internal/app/models/micro_scrm"
)

type RepeatCustomer struct {
	repeatCustomerDao *microScrmDao.RepeatCustomer
}

func NewRepeatCustomer() *RepeatCustomer {
	return &RepeatCustomer{
		repeatCustomerDao: microScrmDao.NewRepeatCustomer(),
	}
}

// SyncRepeatCustomerInfo 同步重复的好友
func (s RepeatCustomer) SyncRepeatCustomerInfo(c app.Ctx, enterpriseID int64, page, size int) ([]workExternalModel.WxUserCorpExternalRepeatCustomerModel, error) {
	return s.repeatCustomerDao.GetRepeatCustomerInfo(c, enterpriseID, page, size)
}

func (s RepeatCustomer) SyncRepeatCustomerTotal(c app.Ctx, enterpriseID int64) (int64, error) {
	return s.repeatCustomerDao.SyncRepeatCustomerTotal(c, enterpriseID)
}

// UpdateRepeatCustomer 更新处理重复的好友关联导购数小于2的外部联系人
func (s RepeatCustomer) UpdateRepeatCustomer(c app.Ctx, enterpriseID int64, lists []workExternalModel.WxUserCorpExternalRepeatCustomerModel) error {
	var workExternalIds []int64
	workExternalIdMap := make(map[int64]int, 0)
	for _, v := range lists {
		workExternalIds = append(workExternalIds, v.WorkExternalUserID)
		workExternalIdMap[v.WorkExternalUserID] = 0
	}
	result, err := microScrmDao.NewUserRelation().GetRepeatCustomerRelationInfo(c, enterpriseID, workExternalIds)
	if err != nil {
		return err
	}
	var UpdateExternalIds []int64
	for _, item := range result {
		if _, ok := workExternalIdMap[item.WorkExternalUserId]; ok {
			workExternalIdMap[item.WorkExternalUserId] = 1
		}
		if item.Num < 2 {
			UpdateExternalIds = append(UpdateExternalIds, item.WorkExternalUserId)
		}
	}
	for k := range workExternalIdMap {
		if workExternalIdMap[k] == 0 {
			UpdateExternalIds = append(UpdateExternalIds, k)
		}
	}
	if len(UpdateExternalIds) > 0 {
		return s.repeatCustomerDao.DeleteRepeatCustomerByExternalIds(c, enterpriseID, UpdateExternalIds)
	}
	return nil
}
