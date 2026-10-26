package tagsDao

import (
	"errors"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	tagsReq "micro-scrm/internal/app/http/request/tags"
	tagModels "micro-scrm/internal/app/models/tag"
	"micro-scrm/internal/pkg/db/gormx"

	"gorm.io/gorm"
)

type ImpressionGroup struct {
	dao       dao.Dao
	tableName string
}

func NewImpressionGroup() *ImpressionGroup {
	return &ImpressionGroup{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: tagModels.ImpressionTagGroupModel{}.TableName(),
	}
}

func (d *ImpressionGroup) GetList(c app.Ctx, request tagsReq.GroupSearchList) ([]tagModels.ImpressionTagGroupModel, int64, error) {
	query := d.dao.WithContext(c).Table(d.tableName).
		Select("impression_tag_group.id,impression_tag_group.enterprise_id,impression_tag_group.third_group_id,impression_tag_group.group_name").
		Where("impression_tag_group.enterprise_id", request.EnterpriseID).
		Where("impression_tag_group.deleted_at = 0")

	if request.Keyword != "" {
		if request.SearchType == 1 {
			query = query.Where("impression_tag_group.group_name like ?", "%"+request.Keyword+"%")
		} else {
			query = query.Joins("left join impression_tag on impression_tag.group_id = impression_tag_group.id").
				Where("impression_tag.name like ?", "%"+request.Keyword+"%")
		}
	}

	var (
		total        int64
		categoryList []tagModels.ImpressionTagGroupModel
		err          error
	)

	err = query.Count(&total).
		Order("impression_tag_group.created_at DESC").
		Offset((request.Page - 1) * request.Size).
		Limit(request.Size).
		Find(&categoryList).Error

	return categoryList, total, err
}

func (d *ImpressionGroup) GetInfo(c app.Ctx, data tagModels.ImpressionTagGroupModel) (tagModels.ImpressionTagGroupModel, error) {
	var result tagModels.ImpressionTagGroupModel
	err := d.dao.WithContext(c).Table(d.tableName).Where(data).First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return result, err
}

func (d *ImpressionGroup) Create(c app.Ctx, data tagModels.ImpressionTagGroupModel) (id int64, err error) {
	err = d.dao.WithContext(c).Table(d.tableName).Where(data).Create(&data).Error
	if err != nil {
		return 0, err
	}
	return data.ID, nil
}

func (d *ImpressionGroup) Update(c app.Ctx, data tagModels.ImpressionTagGroupModel) error {
	return d.dao.WithContext(c).Table(d.tableName).UpdateColumns(data).Error
}
