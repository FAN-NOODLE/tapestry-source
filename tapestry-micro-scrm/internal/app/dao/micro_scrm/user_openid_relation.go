package microScrmDao

import (
	"errors"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	"micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
)

type OpenidRelation struct {
	dao       dao.Dao
	tableName string
}

func NewOpenidRelation() *OpenidRelation {
	return &OpenidRelation{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: microScrmModel.WxUserOpenidRelationModel{}.TableName(),
	}
}

func (d *OpenidRelation) GetInfo(c app.Ctx, enterpriseId int64, cond interface{}) (microScrmModel.WxUserOpenidRelationModel, error) {
	var data microScrmModel.WxUserOpenidRelationModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(cond).Order("created_at desc").First(&data).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return data, err
}

func (d *OpenidRelation) SetInfo(c app.Ctx, enterpriseId int64, data microScrmModel.WxUserOpenidRelationModel) error {
	var result microScrmModel.WorkExternalUserRelationModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(map[string]interface{}{
		"unionid":   data.Unionid,
		"wx_openid": data.WxOpenid,
	}).First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return d.dao.Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Create(&data).Error
	}
	return err
}
