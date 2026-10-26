package work_external_user

import (
	"fmt"

	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/orm/service"
)

const (
	Initialize    = iota
	MEMBERCHANGED = Initialize + 1 //会员信息是否发生:0 未变化 1 变化'
	ISMEMBER      = Initialize     //是否是会员:0会员,1非会员
)

func (w *WorkExternalUserService) CreateCTMember(enterpriseID int64, newData []micro_scrm_enterprise.SyncMemberCampaignToolModel) error {
	return w.Scopes(model.BranchTable(&micro_scrm_enterprise.SyncMemberCampaignToolModel{}, enterpriseID)).Create(&newData).Error
}
func (w *WorkExternalUserService) DeleteCTMember(enterpriseID int64) error {
	table := micro_scrm_enterprise.SyncMemberCampaignToolModel{}
	return w.Scopes(model.BranchTable(&table, enterpriseID)).Delete(table, "id!=?", 0).Error

}

func (w *WorkExternalUserService) SelectOneCTMember(enterpriseID int64) (data micro_scrm_enterprise.SyncMemberCampaignToolModel, err error) {
	err = w.Scopes(model.BranchTable(&micro_scrm_enterprise.SyncMemberCampaignToolModel{}, enterpriseID), service.ScopeCond_EQDelete).First(&data).Error
	return
}

func (w *WorkExternalUserService) FindCTMemberByUnionIds(enterpriseID int64, unionIDs []string) (data []micro_scrm_enterprise.SyncMemberData, err error) {
	err = w.Scopes(model.BranchTable(&micro_scrm_enterprise.SyncMemberCampaignToolModel{}, enterpriseID), service.ScopeCond_EQDelete).
		Where("union_id IN (?)", unionIDs).Find(&data).Error
	return
}

func (w *WorkExternalUserService) UpdateCTMember(enterpriseID int64, data micro_scrm_enterprise.SyncMemberData) (err error) {
	err = w.Scopes(model.BranchTable(&micro_scrm_enterprise.SyncMemberCampaignToolModel{}, enterpriseID)).
		Where("union_id=? AND deleted_at=? ", data.UnionID, 0).Updates(map[string]interface{}{
		"is_member":     data.IsMember,
		"grade_code":    data.GradeCode,
		"member_code":   data.MemberCode,
		"admin_ba":      data.AdminBa,
		"member_change": 1,
	}).Error
	return
}

func (w *WorkExternalUserService) FindCampaignMembers(enterpriseID int64, page, pageSize int) (data []micro_scrm_enterprise.WxUserCorpExternalEuid, err error) {
	wxUserCorpExternalTable := "micro_scrm_enterprise_" + fmt.Sprintf("%d", enterpriseID) + ".wx_user_corp_external "

	err = w.Scopes(model.BranchTable(&micro_scrm_enterprise.SyncMemberCampaignToolModel{}, enterpriseID)).
		Joins("LEFT JOIN "+wxUserCorpExternalTable+" a ON sync_member_campaign_tool.union_id = a.unionid AND a.deleted_at = 0 ").
		Select("DISTINCT sync_member_campaign_tool.union_id,sync_member_campaign_tool.grade_code,sync_member_campaign_tool.member_code,sync_member_campaign_tool.admin_ba,a.euid").
		Where("sync_member_campaign_tool.member_change=? AND sync_member_campaign_tool.is_member=? AND sync_member_campaign_tool.deleted_at=? AND a.unionid IS NOT NULL", MEMBERCHANGED, ISMEMBER, Initialize).
		Order("sync_member_campaign_tool.union_id").Limit(pageSize).Offset((page - 1) * pageSize).
		Find(&data).Error
	return
}

func (w *WorkExternalUserService) UpdatesyncMemberCampaignchangedByunionId(enterpriseID int64) (err error) {
	err = w.Scopes(model.BranchTable(&micro_scrm_enterprise.SyncMemberCampaignToolModel{}, enterpriseID), service.ScopeCond_EQDelete).
		Update("member_change", Initialize).Error
	return
}
