package microScrmDao

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app/dao"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
	"time"
)

type PlanTaskConditionDataCkDao struct {
	dao            dao.Dao
	tableName      string
	localTableName string
}

func NewPlanTaskConditionDataCkDao() *PlanTaskConditionDataCkDao {
	return &PlanTaskConditionDataCkDao{
		dao:            dao.Dao{DB: gormx.GetInstance()},
		tableName:      microScrmModel.TaskConditionDataModel{}.TableName(),
		localTableName: microScrmModel.TaskConditionDataModel{}.TableName() + "_local",
	}
}

func getCkTableName(enterpriseId int64, tableName string) string {
	return fmt.Sprintf("micro_scrm_enterprise_%d_ck.%s", enterpriseId, tableName)
}

func (d *PlanTaskConditionDataCkDao) GetTaskAllotLogIds(c ctx.Context, enterpriseId int64, query interface{}, args ...interface{}) ([]int64, error) {
	var taskAllotLogIds []int64
	lastTime := time.Now().Unix() - 86400*30*3 // 只展示最近三个月的任务
	err := gormx.GetCkInstance().WithContext(c).Table(getCkTableName(enterpriseId, d.tableName)).Where("created_at>=?", lastTime).
		Where(query, args...).Pluck("task_allot_log_id", &taskAllotLogIds).Error
	return taskAllotLogIds, err
}
func (d *PlanTaskConditionDataCkDao) GetTaskAllotLogCount(c ctx.Context, enterpriseId int64, query interface{}, args ...interface{}) (int64, error) {
	var total int64
	lastTime := time.Now().Unix() - 86400*30*3 // 只展示最近三个月的任务
	err := gormx.GetCkInstance().WithContext(c).Table(getCkTableName(enterpriseId, d.tableName)).Where("created_at>=?", lastTime).
		Where(query, args...).Count(&total).Error
	return total, err
}

// CustomerCount 查询子任务的用户画像匹配客户数
func (d *PlanTaskConditionDataCkDao) CustomerCount(c ctx.Context, enterpriseID, taskId, planId int64) (data int64, err error) {
	err = gormx.GetCkInstance().WithContext(c).Table(getCkTableName(enterpriseID, d.tableName)).
		Select("uniq(id)").
		Where("plan_id = ? AND task_id = ? and deleted_at = 0", planId, taskId).
		Scan(&data).Error
	return
}

// CustomerDistinctCount 查询子任务的用户画像匹配客户数(去除重复的union_id)
func (d *PlanTaskConditionDataCkDao) CustomerDistinctCount(c ctx.Context, enterpriseID, planId int64, query interface{}) (data int64, err error) {
	err = gormx.GetCkInstance().WithContext(c).Table(getCkTableName(enterpriseID, d.tableName)).
		Select("uniq(union_id,customer_id)").
		Where("plan_id = ?", planId).
		Where(query).
		Scan(&data).Error
	return
}

// GetCount   查询子任务的用户画像匹配客户数
func (d *PlanTaskConditionDataCkDao) GetCount(c ctx.Context, enterpriseID, planId int64, query interface{}) (data int64, err error) {
	err = gormx.GetCkInstance().WithContext(c).Table(getCkTableName(enterpriseID, d.tableName)).
		Select("uniq(id)").
		Where("plan_id = ?", planId).
		Where(query).
		Scan(&data).Error
	return
}

// GetCountByArg    查询子任务的用户画像匹配客户数
func (d *PlanTaskConditionDataCkDao) GetCountByArg(c ctx.Context, enterpriseID, planId int64, query interface{}, args ...interface{}) (data int64, err error) {
	err = gormx.GetCkInstance().WithContext(c).Table(getCkTableName(enterpriseID, d.tableName)).
		Select("uniq(id)").
		Where("plan_id = ?", planId).
		Where(query, args...).
		Scan(&data).Error
	return
}
