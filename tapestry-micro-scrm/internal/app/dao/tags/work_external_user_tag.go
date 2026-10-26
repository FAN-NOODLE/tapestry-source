package tagsDao

import (
	"fmt"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	tagModels "micro-scrm/internal/app/models/tag"
	"micro-scrm/internal/pkg/db/gormx"
)

type WorkExternalUserTag struct {
	dao   dao.Dao
	model *tagModels.WorkExternalUserTagModel
}

func NewWorkExternalUserTag() *WorkExternalUserTag {
	return &WorkExternalUserTag{
		dao:   dao.Dao{DB: gormx.GetInstance()},
		model: new(tagModels.WorkExternalUserTagModel),
	}
}

func (d *WorkExternalUserTag) GetMany(c app.Ctx, query interface{}) ([]tagModels.WorkExternalUserTagModel, error) {
	var result []tagModels.WorkExternalUserTagModel
	err := d.dao.WithContext(c).Model(d.model).Where(query).Find(&result).Error
	return result, err
}

func (d *WorkExternalUserTag) DeleteData(c app.Ctx, query interface{}) error {
	var model tagModels.WorkExternalUserTagModel
	err := d.dao.WithContext(c).Model(d.model).Where(query).Delete(&model).Error
	return err
}

func (d *WorkExternalUserTag) GetUserTagsByPage(c app.Ctx, page, size int, query interface{}) ([]tagModels.WorkExternalUserTagModel, error) {
	var result []tagModels.WorkExternalUserTagModel
	offset := (page - 1) * size
	err := d.dao.WithContext(c).Model(d.model).Where(query).Offset(offset).Limit(size).Find(&result).Error
	return result, err
}

func (d *WorkExternalUserTag) GetTagCount(c app.Ctx, enterpriseId int64, tagIds []int64) (map[int64]int64, error) {
	result := make(map[int64]int64)
	var tagCount []struct {
		ExternalTagId int64 `json:"external_tag_id"`
		Count         int64 `json:"count"`
	}

	err := gormx.GetCkInstance().WithContext(c).
		Table("micro_scrm_ck.work_external_user_tag as user_tag").
		Select("user_tag.external_tag_id,uniqExact(user_relation.work_external_user_id) as count").
		Joins(fmt.Sprintf("GLOBAL left join (select id,work_external_user_id from micro_scrm_ck.work_external_user_relation where enterprise_id = %d) user_relation on user_relation.id = user_tag.relation_id", enterpriseId)).
		Where("external_tag_id", tagIds).
		Where("user_tag.enterprise_id = ?", enterpriseId).
		Where("user_tag.deleted_at = 0").
		Group("user_tag.external_tag_id").
		Scan(&tagCount).Error

	if err != nil {
		return nil, err
	}

	for _, tag := range tagCount {
		result[tag.ExternalTagId] = tag.Count
	}

	return result, nil
}

func (d *WorkExternalUserTag) GetTagIdsByUserId(c app.Ctx, enterpriseId int64, relationIds []int64) ([]int64, error) {
	var tagIds []int64
	err := d.dao.WithContext(c).Model(d.model).
		Select("external_tag_id").
		Where(map[string]interface{}{
			"enterprise_id": enterpriseId,
			"relation_id":   relationIds,
		}).Scan(&tagIds).Error

	return tagIds, err
}

func (d *WorkExternalUserTag) CreateData(c app.Ctx, data []tagModels.WorkExternalUserTagModel) error {
	err := d.dao.WithContext(c).Model(d.model).Create(&data).Error
	return err
}
