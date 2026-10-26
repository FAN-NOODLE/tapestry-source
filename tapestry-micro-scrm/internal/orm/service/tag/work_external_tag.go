package tag

import (
	"errors"
	"gorm.io/gorm"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/orm/model/micro_scrm/work_external_user"
	"micro-scrm/internal/orm/service"
	"time"
)

// CreateTag 创建标签
func (t *TagService) CreateTag(data []micro_scrm.WorkExternalTagModel) ([]micro_scrm.WorkExternalTagModel, int64, error) {
	result := t.Create(&data)

	return data, result.RowsAffected, result.Error
}

// CheckTagName 检测标签组名称是否存在
func (t *TagService) CheckTagName(enterpriseID, groupID int64, tagName string) (data micro_scrm.WorkExternalTagModel, isExist bool, err error) {
	err = t.Where(map[string]interface{}{
		"enterprise_id": enterpriseID,
		"group_id":      groupID,
		"name":          tagName,
		"deleted_at":    0,
	}).First(&data).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		isExist = false
		err = nil
		return
	}

	isExist = true
	return
}

// GetEnterpriseTagInfoList 获取企业的标签信息集合
func (t *TagService) GetEnterpriseTagInfoList(enterpriseID int64) (data []micro_scrm.WorkExternalTagModel, err error) {
	err = t.Scopes(service.ScopeCond_EQDelete).
		Where("enterprise_id = ?", enterpriseID).
		Find(&data).Error

	return
}

// 根据企微标签分组ID，获取该标签分组下的所有标签集合
func (t *TagService) GetTagListByWwTagGroupID(enterpriseID int64, wwTagGroupID string) (tagList []micro_scrm.WorkExternalTagModel, err error) {
	err = t.Scopes(service.ScopeCond_EQDelete).
		Where("enterprise_id = ?", enterpriseID).
		Where("group_id = ?", wwTagGroupID).
		Find(&tagList).Error

	return
}

// 新增标签
func (t *TagService) BatchesCreate(data []micro_scrm.WorkExternalTagModel) error {
	return t.Create(&data).Error
}

// 根据企微标签ID，获取数据表work_external_tag.id集合
func (t *TagService) GetTagIDListByQYID(enterpriseID int64, wxTagIDList []string) []int64 {
	tadIDList := make([]int64, 0)

	t.Model(&micro_scrm.WorkExternalTagModel{}).
		Scopes(service.ScopeCond_EQDelete).
		Where(map[string]interface{}{
			"enterprise_id": enterpriseID,
			"tag_id":        wxTagIDList,
		}).Pluck("id", &tadIDList)

	return tadIDList
}

// 根据work_external_user_relation.id获取用户的标签集合ID
func (t *TagService) GetTagIDListByRelationID(enterpriseID, relationID int64) (tagIDList []int64, err error) {
	err = t.Model(&work_external_user.WorkExternalUserTagModel{}).
		Scopes(service.ScopeCond_EQDelete).
		Where(map[string]interface{}{"enterprise_id": enterpriseID, "relation_id": relationID}).
		Pluck("external_tag_id", &tagIDList).Error

	return
}

// 根据用户联系人ID，获取用户的标签ID集合
func (t *TagService) GetTagIDListByExternalUserID(enterpriseID, externalUserID int64) (tagIDList []int64, err error) {
	err = t.Model(&work_external_user.WorkExternalUserTagModel{}).
		Scopes(service.ScopeCond_EQDelete).
		Where("enterprise_id = ?", enterpriseID).
		Where("relation_id = ?", func() int64 {
			data := work_external_user.WorkExternalUserRelationModel{}
			t.Scopes(service.ScopeCond_EQDelete).
				Where("enterprise_id = ?", enterpriseID).
				Where("work_external_user_id = ?", externalUserID).
				Last(&data)
			return data.ID
		}()).Pluck("external_tag_id", &tagIDList).Error

	return
}

func (t *TagService) FindWorkTags(enterpriseID int64, tagIds []int64) (data []micro_scrm.WorkExternalTagModel, err error) {
	err = t.Model(&micro_scrm.WorkExternalTagModel{}).
		Scopes(service.ScopeCond_EQDelete).
		Where("enterprise_id=? AND id IN (?)", enterpriseID, tagIds).Find(&data).Error

	return
}

func (t *TagService) FindEffectiveTagBy(enterpriseID int64, tagIds []int64) (data []int64, err error) {
	err = t.Model(&micro_scrm.WorkExternalTagModel{}).Debug().
		Scopes(service.ScopeCond_EQDelete).
		Where("enterprise_id=? AND id IN (?)", enterpriseID, tagIds).
		Select("id").Find(&data).Error
	return
}

// 更新标签信息
func (t *TagService) UpdateTag(data micro_scrm.WorkExternalTagModel) error {
	return t.Model(&micro_scrm.WorkExternalTagModel{}).Scopes(service.ScopeCond_EQDelete).
		Where("enterprise_id = ?", data.EnterpriseID).Where("group_id = ?", data.GroupID).Where("tag_id = ?", data.TagID).UpdateColumns(map[string]interface{}{
		"name":       data.Name,
		"order":      data.Order,
		"updated_at": time.Now().Unix(),
	}).Error
}

// 删除指定标签
func (t *TagService) DeleteTag(data micro_scrm.WorkExternalTagModel) error {
	return t.Model(&micro_scrm.WorkExternalTagModel{}).Scopes(service.ScopeCond_EQDelete).
		Where("enterprise_id = ?", data.EnterpriseID).Where("tag_id = ?", data.TagID).UpdateColumns(map[string]interface{}{
		"deleted_at": time.Now().Unix(),
		"updated_at": time.Now().Unix(),
	}).Error
}

// 重排标签
func (t *TagService) ShuffleTag(data micro_scrm.WorkExternalTagModel) error {
	return t.Model(&micro_scrm.WorkExternalTagModel{}).Scopes(service.ScopeCond_EQDelete).
		Where("enterprise_id = ?", data.EnterpriseID).Where("group_id = ?", data.GroupID).Where("tag_id = ?", data.TagID).UpdateColumns(map[string]interface{}{
		"order":      data.Order,
		"updated_at": time.Now().Unix(),
	}).Error
}
