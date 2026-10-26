package microScrmDao

import (
	"micro-scrm/internal/app/dao"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

type PlanTaskCrowdDataLogDao struct {
	dao       dao.Dao
	tableName string
}

func NewPlanTaskCrowdDataLogDao() *PlanTaskCrowdDataLogDao {
	return &PlanTaskCrowdDataLogDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: microScrmModel.PlanTaskCrowdDataLogModel{}.TableName(),
	}
}

// InsertICrowdDataLog 插入MA的人群包的记录
func (d *PlanTaskCrowdDataLogDao) InsertICrowdDataLog(c ctx.Context, data microScrmModel.PlanTaskCrowdDataLogModel, enterpriseId int64) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Create(&data).Error
}

// GetCrowdDataResult 获取MA任务人群包的匹配结果
func (d *PlanTaskCrowdDataLogDao) GetCrowdDataResult(c ctx.Context, enterpriseId int64, taskId int64) ([]microScrmModel.PlanTaskCrowdDataLogModel, error) {
	var result []microScrmModel.PlanTaskCrowdDataLogModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where("task_id=?", taskId).Find(&result).Error
	return result, err
}

// GetCrowdDataResultByPage 获取MA任务人群包的匹配结果
func (d *PlanTaskCrowdDataLogDao) GetCrowdDataResultByPage(c ctx.Context, enterpriseId int64, taskId int64, page, size int) ([]microScrmModel.PlanTaskCrowdDataLogModel, error) {
	var result []microScrmModel.PlanTaskCrowdDataLogModel
	offset := (page - 1) * size
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where("task_id=?", taskId).Offset(offset).Limit(size).Find(&result).Error
	return result, err
}

// InsertICrowdDataLogBatch 批量插入MA的人群包的记录
func (d *PlanTaskCrowdDataLogDao) InsertICrowdDataLogBatch(c ctx.Context, data []microScrmModel.PlanTaskCrowdDataLogModel, enterpriseId int64) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Create(&data).Error
}

func (d *PlanTaskCrowdDataLogDao) GetCrowdDataResultById(c ctx.Context, enterpriseId, taskId, id int64) ([]microScrmModel.PlanTaskCrowdDataLogModel, error) {
	var result []microScrmModel.PlanTaskCrowdDataLogModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where("id > ? and task_id=?", id, taskId).Limit(1000).Find(&result).Error
	return result, err
}
