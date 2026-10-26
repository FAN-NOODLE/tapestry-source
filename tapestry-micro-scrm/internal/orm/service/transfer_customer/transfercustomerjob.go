package transfer_customer

import (
	"fmt"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise/transfer_customer"

	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/orm/service"
)

func (t *TransferCustomerService) FindTransferCustomerWelcomelanguage(enterpriseID int64) (data transfer_customer.TransferCustomerWelcomeLanguageModel, err error) {
	err = t.Scopes(model.BranchTable(&data, enterpriseID), service.ScopeCond_EQDelete).First(data).Error
	return
}
func (t *TransferCustomerService) JobFindAdminBaByUnionid(enterpriseID int64, unionid string) (data JobWorkExternalAdminBaUser, err error) {
	membersTable := "micro_scrm_enterprise_" + fmt.Sprintf("%d", enterpriseID) + ".members menber"
	err = t.Scopes(model.BranchTable(&micro_scrm_enterprise.SyncMemberCampaignToolModel{}, enterpriseID)).
		Select("menber.id,menber.userid user_id,menber.name").
		Joins("INNER JOIN "+membersTable+" ON sync_member_campaign_tool.admin_ba=menber.mars_id AND sync_member_campaign_tool.admin_ba!='' AND menber.deleted_at=0").
		Where("sync_member_campaign_tool.union_id=? AND sync_member_campaign_tool.deleted_at=?", unionid, 0).Find(&data).Error
	return
}

func (t *TransferCustomerService) JobFindworkExternalUserRelation(enterpriseID, workExternalUserId int64) (data []JobWorkExternalAdminBaUser, err error) {
	err = t.Table("work_external_user_relation").Joins("INNER JOIN members a ON work_external_user_relation.member_id = a.id AND a.deleted_at = 0 AND a.enterprise_id=?", enterpriseID).
		Select("a.userid,a.name,a.id").Where("work_external_user_relation.work_external_user_id=? AND work_external_user_relation.deleted_at=?", workExternalUserId, 0).
		Group("a.userid,a.name,a.id").Find(&data).Error
	return
}

func (t *TransferCustomerService) JobAdminWxUserCorpExternal(enterpriseID int64, page, pageSize int) (data []JobWorkExternalUser, err error) {
	err = t.Scopes(model.BranchTable(&micro_scrm_enterprise.WxUserCorpExternalModel{}, enterpriseID), service.ScopeCond_EQDelete).
		Select("id,openid,unionid").Where("unionid=?", "oUmNp1sDylOKgjRX2kpubXDQUYvw").Order("id").
		Limit(pageSize).Offset((page - 1) * pageSize).Find(&data).Error
	return
}

func (t *TransferCustomerService) BatchCreatetransferCustomerData(transferCustomerDatas []*transfer_customer.TransferCustomerDataModel) error {
	// t.Scopes(model.BranchTable(&micro_scrm_enterprise.WxUserCorpExternalModel{}
	return nil
}
