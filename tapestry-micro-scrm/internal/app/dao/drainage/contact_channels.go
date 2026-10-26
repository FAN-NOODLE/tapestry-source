package drainageContactDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	"micro-scrm/internal/app/models/drainage"
	"micro-scrm/internal/pkg/db/gormx"
)

type DrainageContactChannelDao struct {
	dao       dao.Dao
	tableName string
	model     *drainageContactModel.DrainageContactChannelsModel
}

func NewDrainageContactChannelDao() *DrainageContactChannelDao {
	return &DrainageContactChannelDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: drainageContactModel.DrainageContactChannelsModel{}.TableName(),
		model:     new(drainageContactModel.DrainageContactChannelsModel),
	}
}

func (d *DrainageContactChannelDao) GetOne(c app.Ctx, query interface{}, args ...interface{}) (drainageContactModel.DrainageContactChannelsModel, error) {
	var result drainageContactModel.DrainageContactChannelsModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *DrainageContactChannelDao) GetMany(c app.Ctx, query interface{}, args ...interface{}) ([]drainageContactModel.DrainageContactChannelsModel, error) {
	var result []drainageContactModel.DrainageContactChannelsModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *DrainageContactChannelDao) GetOnePreload(c app.Ctx, query interface{}, args ...interface{}) (drainageContactModel.DrainageContactChannelsModel, error) {
	var result drainageContactModel.DrainageContactChannelsModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Preload("Labels.Tags").Find(&result).Error
	return result, err
}

func (d *DrainageContactChannelDao) GetManyPreload(c app.Ctx, query interface{}, args ...interface{}) ([]drainageContactModel.DrainageContactChannelsModel, error) {
	var result []drainageContactModel.DrainageContactChannelsModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Preload("Labels.Tags").Find(&result).Error
	return result, err
}
