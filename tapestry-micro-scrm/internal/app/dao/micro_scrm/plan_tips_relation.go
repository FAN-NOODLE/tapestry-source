package microScrmDao

import (
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
)

type PlanTipsRelationDao struct {
	dao       dao.Dao
	tableName string
}

func NewPlanTipsRelationDao() *PlanTipsRelationDao {
	return &PlanTipsRelationDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: microScrmModel.PlanTipsRelationModel{}.TableName(),
	}
}

// WithTx 切换db连接,用于事物
func (d *PlanTipsRelationDao) WithTx(tx *gorm.DB) *PlanTipsRelationDao {
	d.dao.DB = tx
	return d
}

func (d *PlanTipsRelationDao) GetOne(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) (microScrmModel.PlanTipsRelationModel, error) {
	var result microScrmModel.PlanTipsRelationModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *PlanTipsRelationDao) GetMany(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]microScrmModel.PlanTipsRelationModel, error) {
	var result []microScrmModel.PlanTipsRelationModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *PlanTipsRelationDao) CreateData(c app.Ctx, enterpriseId int64, data microScrmModel.PlanTipsRelationModel) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Create(&data).Error
	return err
}

func (d *PlanTipsRelationDao) GetClosePlanIds(c app.Ctx, enterpriseId int64, userid string, planIds []int64) ([]int64, error) {
	var result []int64
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(map[string]interface{}{
		"user_id": userid,
		"plan_id": planIds,
	}).Pluck("plan_id", &result).Error
	return result, err
}
