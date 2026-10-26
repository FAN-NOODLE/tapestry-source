package templateDao

import (
	"errors"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	templateModels "micro-scrm/internal/app/models/template"
	"micro-scrm/internal/pkg/db/gormx"
)

type TemplatePlanRelation struct {
	dao       dao.Dao
	tableName string
}

// WithTx 切换db连接,用于事物
func (d *TemplatePlanRelation) WithTx(tx *gorm.DB) *TemplatePlanRelation {
	d.dao.DB = tx
	return d
}

func NewTemplatePlanRelation() *TemplatePlanRelation {
	return &TemplatePlanRelation{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: templateModels.TemplatePlanRelationModel{}.TableName(),
	}
}

func (d *TemplatePlanRelation) Check(c app.Ctx, enterpriseID int64, id int) (Total int64, err error) {
	err = d.dao.WithContext(c).Table(d.tableName).
		Where(map[string]interface{}{"enterprise_id": enterpriseID, "template_id": id}).
		Count(&Total).Error
	return

}

func (d *TemplatePlanRelation) FindOrCreateUseCard(c app.Ctx, enterpriseID int64, planId int64, templateID int, data templateModels.TemplatePlanRelationModel) error {
	var info templateModels.TemplatePlanRelationModel
	err := d.dao.WithContext(c).Table(d.tableName).
		Where(map[string]interface{}{"enterprise_id": enterpriseID, "template_id": templateID, "plan_id": planId}).
		First(&info).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return d.dao.WithContext(c).Table(d.tableName).Create(&data).Error
	}
	return nil
}

func (d *TemplatePlanRelation) UpdatePlanCardInfo(c app.Ctx, enterpriseID int64, planIds []int64) error {
	return d.dao.WithContext(c).Table(d.tableName).
		Where(map[string]interface{}{"enterprise_id": enterpriseID, "plan_id": planIds}).
		Update("plan_status", 1).Error
}

func (d *TemplatePlanRelation) UpdateData(c app.Ctx, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Table(d.tableName).Where(query).UpdateColumns(data).Error
	return err
}
