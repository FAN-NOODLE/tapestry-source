package work_external_user

import (
	"micro-scrm/internal/orm/model/micro_scrm/work_external_user"
	"micro-scrm/internal/orm/service"
)

// GetRelationInfo 根据企业ID，外部联系人ID，员工ID，获取外部联系人和员工的的关系信息
func (w *WorkExternalUserService) GetRelationInfo(enterpriseID, workExternalUserId, memberID int64) (data work_external_user.WorkExternalUserRelationModel) {
	w.Model(&data).
		Scopes(service.ScopeCond_EQDelete).
		Where(map[string]interface{}{"enterprise_id": enterpriseID, "work_external_user_id": workExternalUserId, "member_id": memberID}).
		First(&data)

	return
}
