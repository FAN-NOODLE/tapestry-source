package tagsDao

import (
	"fmt"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	tagModels "micro-scrm/internal/app/models/tag"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/helper"
)

type ImpressionRelation struct {
	dao       dao.Dao
	tableName string
}

func NewImpressionRelation() *ImpressionRelation {
	return &ImpressionRelation{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: tagModels.ImpressionTagRelationModel{}.TableName(),
	}
}

func (d *ImpressionRelation) GetOne(c app.Ctx, query interface{}, args ...interface{}) (tagModels.ImpressionTagRelationModel, error) {
	var result tagModels.ImpressionTagRelationModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *ImpressionRelation) GetMany(c app.Ctx, query interface{}, args ...interface{}) ([]tagModels.ImpressionTagRelationModel, error) {
	var result []tagModels.ImpressionTagRelationModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *ImpressionRelation) SaveData(c app.Ctx, data tagModels.ImpressionTagRelationModel) error {
	err := d.dao.WithContext(c).Table(d.tableName).Save(&data).Error
	return err
}

func (d *ImpressionRelation) Create(c app.Ctx, data tagModels.ImpressionTagRelationModel) error {
	return d.dao.WithContext(c).Table(d.tableName).Create(&data).Error
}

func (d *ImpressionRelation) BatchCreate(c app.Ctx, data []tagModels.ImpressionTagRelationModel) error {
	return d.dao.WithContext(c).Table(d.tableName).CreateInBatches(data, len(data)).Error
}

func (d *ImpressionRelation) RemoveByExternalUserIdTagIds(c app.Ctx, enterpriseId, externalUserId int64, tagIds []int64) error {
	err := d.dao.WithContext(c).Table(d.tableName).
		Where("enterprise_id", enterpriseId).
		Where("external_user_id", externalUserId).
		Where("tag_id IN ?", tagIds).
		Update("deleted_at", helper.GetNowDateTimestamp()).Error
	return err
}

func (d *ImpressionRelation) GetListByUserIdAndTagId(c app.Ctx, enterpriseId, externalUserId int64, tagIds []int64) ([]tagModels.ImpressionTagRelationModel, error) {
	var result []tagModels.ImpressionTagRelationModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(map[string]interface{}{
		"enterprise_id":    enterpriseId,
		"external_user_id": externalUserId,
		"deleted_at":       0,
	}).Where("tag_id IN ?", tagIds).Find(&result).Error

	return result, err
}

// GetExistedTagIdsByUserIdAndTagId 获取已存在的标签id集合
func (d *ImpressionRelation) GetExistedTagIdsByUserIdAndTagId(c app.Ctx, enterpriseId, externalUserId int64, tagIds []int64, userId string) ([]int64, error) {
	var result []int64
	err := d.dao.WithContext(c).Table(d.tableName).Where(map[string]interface{}{
		"enterprise_id":    enterpriseId,
		"external_user_id": externalUserId,
		"user_id":          userId,
		"deleted_at":       0,
	}).Where("tag_id IN ?", tagIds).Pluck("tag_id", &result).Error

	return result, err
}

// GetListByExternalUserId 导购给用户打的印象标签
func (d *ImpressionRelation) GetListByExternalUserId(c app.Ctx, enterpriseId, externalUserId int64, userIds []string) ([]tagModels.ImpressionTagRelationModel, error) {
	var result []tagModels.ImpressionTagRelationModel
	err := d.dao.WithContext(c).Table(d.tableName).
		Select("id,enterprise_id,external_user_id,tag_id").
		Where(map[string]interface{}{
			"enterprise_id":    enterpriseId,
			"external_user_id": externalUserId,
			"user_id":          userIds,
			"deleted_at":       0,
		}).Find(&result).Order("id DESC").Error

	return result, err
}

func (d *ImpressionRelation) GetTagCount(c app.Ctx, enterpriseId int64, tagIds []int64) (map[int64]int64, error) {
	result := make(map[int64]int64)
	var tagCount []struct {
		TagId int64 `json:"tag_id"`
		Count int64 `json:"count"`
	}

	err := gormx.GetCkInstance().WithContext(c).
		Scopes(func(db *gorm.DB) *gorm.DB {
			return db.Table(fmt.Sprintf("micro_scrm_ck.%s", d.tableName))
		}).
		Select("tag_id,count(DISTINCT external_user_id) as count").
		Where("tag_id", tagIds).
		Where("enterprise_id = ?", enterpriseId).
		Where("deleted_at = 0").
		Group("tag_id").
		Scan(&tagCount).Error

	if err != nil {
		return nil, err
	}

	for _, tag := range tagCount {
		result[tag.TagId] = tag.Count
	}

	return result, nil
}

func (d *ImpressionRelation) GetTagIdsByUserId(c app.Ctx, enterpriseId, externalUserId int64) ([]int64, error) {
	var tagIds []int64
	err := d.dao.WithContext(c).Table(d.tableName).
		Select("tag_id").
		Where(map[string]interface{}{
			"enterprise_id":    enterpriseId,
			"external_user_id": externalUserId,
		}).Scan(&tagIds).Error

	return tagIds, err
}

func (d *ImpressionRelation) GetManyById(c app.Ctx, query interface{}, args ...interface{}) ([]tagModels.ImpressionTagRelationModel, error) {
	var result []tagModels.ImpressionTagRelationModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(query, args...).Order("id asc").Limit(1000).Find(&result).Error
	return result, err
}
