package tag

import (
	"micro-scrm/internal/orm/model/micro_scrm/work_external_user"
	"micro-scrm/internal/orm/service"
)

// CreateUserTag  创建标签
func (t *TagService) CreateUserTag(data work_external_user.WorkExternalUserTagModel) (work_external_user.WorkExternalUserTagModel, int64, error) {
	result := t.Create(&data)

	return data, result.RowsAffected, result.Error
}

// GetEnterpriseTagUserInfo 获取外部联系人的企业标签信息是否已同步
func (t *TagService) GetEnterpriseTagUserInfo(enterpriseID, relationID, externalTagID int64) (data work_external_user.WorkExternalUserTagModel, row int64, err error) {
	result := t.Scopes(service.ScopeCond_EQDelete).
		Where(map[string]interface{}{"enterprise_id": enterpriseID, "relation_id": relationID,"external_tag_id":externalTagID}).
		First(&data)

	return data,result.RowsAffected,result.Error
}
