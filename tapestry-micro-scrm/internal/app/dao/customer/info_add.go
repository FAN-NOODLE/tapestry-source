package customerDao

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	customerReq "micro-scrm/internal/app/http/request/customer"
	customerModels "micro-scrm/internal/app/models/customer"
	"micro-scrm/internal/pkg/db/gormx"
)

type CustomerInfoAdd struct {
	dao       dao.Dao
	tableName string
}

func NewCustomerInfoAdd() *CustomerInfoAdd {
	return &CustomerInfoAdd{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: customerModels.CustomerInfoAddModel{}.TableName(),
	}
}

func (d *CustomerInfoAdd) Create(c app.Ctx, data customerModels.CustomerInfoAddModel) error {
	return d.dao.WithContext(c).Table(d.tableName).Create(&data).Error
}

func (d *CustomerInfoAdd) Save(c app.Ctx, id int64, data customerReq.CustomerInfoAddReq) error {
	return d.dao.WithContext(c).Table(d.tableName).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"address":         data.Address,
			"shoe_size":       data.ShoeSize,
			"month":           data.Month,
			"day":             data.Day,
			"operate_user_id": data.OperateUserId,
			"operate_name":    data.OperateName,
		}).Error
}

func (d *CustomerInfoAdd) GetLogList(c app.Ctx, enterpriseId int64, updatedAt int64) (data []customerModels.CustomerInfoAddModel, err error) {
	err = d.dao.WithContext(c).Table(d.tableName).
		Select("operate_user_id,address,shoe_size,external_unionid,updated_at").
		Where("enterprise_id = ?", enterpriseId).
		Where("updated_at > ?", updatedAt).
		Where("deleted_at = 0").Find(&data).Error
	return data, err
}

func (d *CustomerInfoAdd) GetInfo(c app.Ctx, enterpriseId int64, unionId string) (data customerModels.CustomerInfoAddModel, err error) {
	err = d.dao.WithContext(c).Table(d.tableName).
		Where("enterprise_id = ?", enterpriseId).
		Where("external_unionid = ?", unionId).
		Where("deleted_at = 0").Find(&data).Error
	return data, err
}

func (d *CustomerInfoAdd) GetListByTime(c app.Ctx, enterpriseId int64, time int64) (data []customerModels.CustomerInfoAddModel, err error) {
	err = d.dao.WithContext(c).Table(d.tableName).
		Where("enterprise_id = ?", enterpriseId).
		Where("updated_at > ?", time).
		Where("deleted_at = 0").Find(&data).Error
	return data, err
}
