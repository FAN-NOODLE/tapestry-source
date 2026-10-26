package transfer_customer

import (
	"fmt"

	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/orm/model/micro_scrm/work_external_user"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise/transfer_customer"
	"micro-scrm/internal/orm/service"
)

// InsertTransferCustomerData 插入在职继承数据
func (s *TransferCustomerService) InsertTransferCustomerData(enterpriseID int64, insertDataList []*transfer_customer.TransferCustomerDataModel) error {
	return s.Scopes(model.BranchTable(&transfer_customer.TransferCustomerDataModel{}, enterpriseID)).
		Create(&insertDataList).Error
}

// GetTransferDataList 获取要在职继承的数据
func (s *TransferCustomerService) GetTransferDataList(enterpriseID int64, paginate service.PaginateParam) (dataList []*transfer_customer.TransferCustomerDataModel_GetTransferDataList, err error) {
	tablePrefix := model.BranchModel{}.DBPrefix()

	// 子表
	SMCT := fmt.Sprintf(`%s%d.%s`, tablePrefix, enterpriseID, micro_scrm_enterprise.SyncMemberCampaignToolModel{}.TableName())
	WUCE := fmt.Sprintf(`%s%d.%s`, tablePrefix, enterpriseID, micro_scrm_enterprise.WxUserCorpExternalModel{}.TableName())
	EM := fmt.Sprintf(`%s%d.%s`, tablePrefix, enterpriseID, micro_scrm_enterprise.MembersModel{}.TableName())
	// 关联表
	WEUR := work_external_user.WorkExternalUserRelationModel{}.TableName()
	M := micro_scrm.MemberModel{}.TableName()

	// 查询
	err = s.Table("(?) AS WS",
		s.Table(SMCT+" AS SMCT").
			Joins("LEFT JOIN "+WUCE+" WUCE ON WUCE.`unionid` = SMCT.`union_id`").
			Joins("LEFT JOIN "+EM+" EM ON EM.`mars_id` = SMCT.`admin_ba`").
			Select("WUCE.`id` AS customer_id", "WUCE.`unionid` AS union_id", "WUCE.`openid` AS work_external_user_id", "EM.`userid` AS user_id", "EM.`name` AS takeover_name").
			Where("SMCT.`admin_ba` <> ?", "").
			Where("WUCE.`openid` <> ?", "").
			Where("EM.`userid` <> ?", "").
			Offset((paginate.Page-1)*paginate.PageSize).
			Limit(paginate.PageSize),
	).Joins("LEFT JOIN "+WEUR+" WEUR ON WEUR.`work_external_user_id` = WS.`customer_id` AND WEUR.deleted_at=?  AND WEUR.`enterprise_id` = ?", 0, enterpriseID).
		Joins("LEFT JOIN "+M+" M ON M.`id` = WEUR.`member_id` AND M.`enterprise_id` = ?", enterpriseID).
		Where("WEUR.`deleted_at` = ? AND M.userid !='' ", 0).
		Select("DISTINCT M.`userid` AS handover_user_id", "M.`name` AS handover_name", "WS.`user_id` AS takeover_user_id", "WS.`takeover_name`", "WS.`work_external_user_id`", "WS.`union_id`").
		Find(&dataList).Error

	return
}
