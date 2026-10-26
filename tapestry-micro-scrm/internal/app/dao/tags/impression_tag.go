package tagsDao

import (
	"errors"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	tagsReq "micro-scrm/internal/app/http/request/tags"
	tagModels "micro-scrm/internal/app/models/tag"
	"micro-scrm/internal/pkg/db/gormx"
)

type Impression struct {
	dao       dao.Dao
	tableName string
}

func NewImpression() *Impression {
	return &Impression{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: tagModels.ImpressionTagModel{}.TableName(),
	}
}

func (d *Impression) GetList(c app.Ctx, request tagsReq.TagList) ([]tagModels.ImpressionTagModel, int64, error) {
	query := d.dao.WithContext(c).Table(d.tableName).Where("enterprise_id", request.EnterpriseID)
	if request.Keyword != "" {
		query = query.Where("name like ?", "%"+request.Keyword+"%")
	}
	if request.TagGroupId != 0 {
		query = query.Where("group_id", request.TagGroupId)
	}

	var (
		total        int64
		categoryList []tagModels.ImpressionTagModel
		err          error
	)

	err = query.Count(&total).
		Order("created_at DESC").
		Offset((request.Page - 1) * request.Size).
		Limit(request.Size).
		Find(&categoryList).Error

	return categoryList, total, err
}

func (d *Impression) GetListByTagId(c app.Ctx, tagIds []int64) ([]tagModels.ImpressionTagModel, error) {
	var result []tagModels.ImpressionTagModel
	err := d.dao.WithContext(c).Table(d.tableName).Where("id IN ?", tagIds).Find(&result).Error
	return result, err
}

func (d *Impression) GetInfo(c app.Ctx, data tagModels.ImpressionTagModel) (tagModels.ImpressionTagModel, error) {
	var result tagModels.ImpressionTagModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(data).First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return result, err
}

func (d *Impression) Create(c app.Ctx, data tagModels.ImpressionTagModel) error {
	err := d.dao.WithContext(c).Table(d.tableName).Create(&data).Error
	return err
}

func (d *Impression) Update(c app.Ctx, data tagModels.ImpressionTagModel) error {
	return d.dao.WithContext(c).Table(d.tableName).UpdateColumns(data).Error
}

// TagsAndTagGroupList 全部印象标签分组和标签
func (d *Impression) TagsAndTagGroupList(c app.Ctx, enterpriseID int64) (info []tagModels.TagsAndGroupTagListRes, total int64, err error) {

	query := d.dao.WithContext(c).Table(d.tableName).
		Select("impression_tag.id as tag_id,impression_tag.group_id,impression_tag.third_group_id,impression_tag.name,impression_tag.description,itg.group_name").
		Joins("LEFT JOIN impression_tag_group itg ON impression_tag.group_id = itg.id ").
		Where(map[string]interface{}{
			"impression_tag.enterprise_id": enterpriseID,
			"impression_tag.deleted_at":    0,
		})

	err = query.Count(&total).Order("itg.id ASC").
		Find(&info).Error

	return info, total, err

}

func (d *Impression) GetTagList(c app.Ctx, request tagsReq.ReqList) (list []tagModels.ImpressionTagModel, err error) {

	err = d.dao.WithContext(c).Table(d.tableName).
		Where("enterprise_id", request.EnterpriseID).
		Order("created_at DESC").
		Find(&list).Error

	return
}

func (d *Impression) GetMany(c app.Ctx, query interface{}, args ...interface{}) ([]tagModels.ImpressionTagModel, error) {
	var result []tagModels.ImpressionTagModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(query, args...).Find(&result).Error
	return result, err
}
