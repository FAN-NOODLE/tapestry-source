package tag

import (
	"errors"
	"gorm.io/gorm"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/orm/service"
)

// CheckTagGroupName 检测标签组名称是否存在
func (t *TagService) CheckTagGroupName(enterpriseID int64, groupName string) (data micro_scrm.WorkExternalTagGroupModel, isExist bool, err error) {
	err = t.Where(map[string]interface{}{
		"enterprise_id": enterpriseID,
		"group_name":    groupName,
	}).First(&data).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		isExist = false
		err = nil
		return
	}

	isExist = true

	return
}

func (t *TagService) CreateTagGroup(data micro_scrm.WorkExternalTagGroupModel) error {
	return t.Create(&data).Error
}

func (t *TagService) CreateTagGroup2(data micro_scrm.WorkExternalTagGroupModel) (micro_scrm.WorkExternalTagGroupModel, int64, error) {
	result := t.Create(&data)

	return data, result.RowsAffected, result.Error
}

// GetTagGroupDetail 根据企业ID，查询企业标签组详情
func (t *TagService) GetTagGroupDetail(enterpriseID int64) (tagGroupList []micro_scrm.WorkExternalTagGroupModel, err error) {
	result := t.Model([]micro_scrm.WorkExternalTagGroupModel{}).
		Where(map[string]interface{}{"enterprise_id": enterpriseID, "deleted_at": 0}).
		Preload("TagList", func(db *gorm.DB) *gorm.DB {
			return db.Where(map[string]interface{}{"enterprise_id": enterpriseID, "deleted_at": 0})
		}).
		Find(&tagGroupList)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return tagGroupList, result.Error
	}

	return tagGroupList, nil
}

// GetTagGroupInfoByWwTagGroupID 根据企业微信ID，获取本地标签分组信息
func (t *TagService) GetTagGroupInfoByWwTagGroupID(enterpriseID int64, wwTagGroupID string) (isExist bool, data micro_scrm.WorkExternalTagGroupModel, err error) {
	err = t.Scopes(service.ScopeCond_EQDelete).
		Where("enterprise_id = ?", enterpriseID).
		Where("group_id = ?", wwTagGroupID).
		First(&data).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			isExist = false
			err = nil
			return
		}
		return
	} else {
		isExist = true
	}

	return
}
