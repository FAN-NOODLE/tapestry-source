package workDao

import (
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	workModels "micro-scrm/internal/app/models/work"
	"micro-scrm/internal/pkg/db/gormx"
)

type PlanTasksTagRelation struct {
	dao   dao.Dao
	model *workModels.PlanTasksTagRelationModel
}

func NewPlanTasksTagRelation() *PlanTasksTagRelation {
	return &PlanTasksTagRelation{
		dao:   dao.Dao{DB: gormx.GetInstance()},
		model: new(workModels.PlanTasksTagRelationModel),
	}
}

// WithTx 切换db连接,用于事物
func (d *PlanTasksTagRelation) WithTx(tx *gorm.DB) *PlanTasksTagRelation {
	d.dao.DB = tx
	return d
}

func (d *PlanTasksTagRelation) GetOne(c app.Ctx, query interface{}, args ...interface{}) (workModels.PlanTasksTagRelationModel, error) {
	var result workModels.PlanTasksTagRelationModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *PlanTasksTagRelation) GetMany(c app.Ctx, query interface{}, args ...interface{}) ([]workModels.PlanTasksTagRelationModel, error) {
	var result []workModels.PlanTasksTagRelationModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *PlanTasksTagRelation) CreateData(c app.Ctx, data workModels.PlanTasksTagRelationModel) (int64, error) {
	err := d.dao.WithContext(c).Model(d.model).Create(&data).Error
	return data.ID, err
}

func (d *PlanTasksTagRelation) CreateInBatchesData(c app.Ctx, data []workModels.PlanTasksTagRelationModel) error {
	err := d.dao.WithContext(c).Model(d.model).CreateInBatches(&data, len(data)).Error
	return err
}

func (d *PlanTasksTagRelation) SaveData(c app.Ctx, data workModels.PlanTasksTagRelationModel) error {
	err := d.dao.WithContext(c).Model(d.model).Save(&data).Error
	return err
}

func (d *PlanTasksTagRelation) UpdateData(c app.Ctx, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Model(d.model).Where(query).UpdateColumns(data).Error
	return err
}

func (d *PlanTasksTagRelation) DeleteData(c app.Ctx, query interface{}) error {
	var model workModels.PlanTasksTagRelationModel
	err := d.dao.WithContext(c).Model(d.model).Where(query).Delete(&model).Error
	return err
}

func (d *PlanTasksTagRelation) ForceDelete(c app.Ctx, query interface{}) error {
	var model workModels.PlanTasksTagRelationModel
	err := d.dao.WithContext(c).Model(d.model).Where(query).Unscoped().Delete(&model).Error
	return err
}

func (d *PlanTasksTagRelation) GetOneByTagId(c app.Ctx, enterpriseID, tagId int64) (workModels.PlanTasksTagRelationModel, error) {
	return d.GetOne(c, "enterprise_id = ? and tag_id = ?", enterpriseID, tagId)
}
