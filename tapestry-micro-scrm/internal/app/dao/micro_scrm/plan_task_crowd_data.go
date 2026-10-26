package microScrmDao

import (
	"micro-scrm/internal/app/dao"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

type PlanTaskCrowdDataDao struct {
	dao       dao.Dao
	tableName string
}

func NewPlanTaskCrowdDataDao() *PlanTaskCrowdDataDao {
	return &PlanTaskCrowdDataDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: microScrmModel.PlanTaskCrowdDataModel{}.TableName(),
	}
}

func (d *PlanTaskCrowdDataDao) InsertCrowdData(c ctx.Context, data microScrmModel.PlanTaskCrowdDataModel, enterpriseId int64) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Create(&data).Error
}

func (d *PlanTaskCrowdDataDao) GetCrowdDataCount(c ctx.Context, enterpriseId, taskId int64) (total int64, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where("task_id=?", taskId).Count(&total).Error
	return
}

func (d *PlanTaskCrowdDataDao) GetCrowdDataList(c ctx.Context, enterpriseId int64, taskId int64, page, size int) (data []microScrmModel.PlanTaskCrowdDataModel, err error) {
	offset := (page - 1) * size
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Where("task_id", taskId).Offset(offset).Limit(size).Find(&data).Error
	return
}

func (d *PlanTaskCrowdDataDao) InsertCrowdDataBatch(c ctx.Context, data []microScrmModel.PlanTaskCrowdDataModel, enterpriseId int64) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Create(&data).Error
}

func (d *PlanTaskCrowdDataDao) GetCrowdDataById(c ctx.Context, enterpriseId int64, taskId int64, id int64) (data []microScrmModel.PlanTaskCrowdDataModel, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where("id > ? and task_id=?", id, taskId).Limit(1000).Find(&data).Error
	return
}
