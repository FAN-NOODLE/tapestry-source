package microScrmDao

import (
	"fmt"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	planReq "micro-scrm/internal/app/http/request/plan"
	excitationModels "micro-scrm/internal/app/models/excitation"
	"micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
)

type PlanTaskAllotLogDao struct {
	dao       dao.Dao
	tableName string
}

func NewPlanTaskAllotLogDao() *PlanTaskAllotLogDao {
	return &PlanTaskAllotLogDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: microScrmModel.TaskAllotLogModel{}.TableName(),
	}
}

// CreateTaskAllotLogData task_allot_log数据插入
func (d *PlanTaskAllotLogDao) CreateTaskAllotLogData(c app.Ctx, enterpriseID int64, info microScrmModel.TaskAllotLogModel) (ID int64, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseID)).Create(&info).Error
	if err != nil {
		return 0, err
	}
	ID = info.ID
	return ID, err
}

// GetPlanOrderMonthLists 获取计划任务的日历列表数据
func (d *PlanTaskAllotLogDao) GetPlanOrderMonthLists(c app.Ctx, request planReq.MonthListsReq, userId string) (data []planReq.PlanAndTaskAllotFinishList, err error) {
	plansTable := d.dao.GetTable("plans", request.EnterpriseId)
	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, request.EnterpriseId)).
		Select("plans.id,task_allot_log.created_at,task_allot_log.is_finished").
		Joins("left join "+plansTable+" on task_allot_log.plan_id=plans.id").
		Where("plans.deleted_at=0 and task_allot_log.deleted_at=0 and plans.is_open = 1 and task_allot_log.user_id=? AND task_allot_log.created_at>=? and task_allot_log.created_at<=?", userId, request.StartTime, request.EndTime)
	if request.PlanId > 0 {
		query = query.Where("plans.id=?", request.PlanId)
	}
	err = query.Order("task_allot_log.created_at asc").Find(&data).Error
	return
}

// GetPlanLists 获取导购任务中心的主任务列表
func (d *PlanTaskAllotLogDao) GetPlanLists(c app.Ctx, request planReq.ListsReq, userId string) (data []planReq.PlansListInfo, err error) {
	plansTable := d.dao.GetTable("plans", request.EnterpriseId)
	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, request.EnterpriseId)).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if request.Name != "" {
				db.Where("plans.name like ?", "%"+request.Name+"%")
			}
			return db
		}).
		Select("plans.id,plans.remark,plans.start_time,plans.end_time,plans.name,task_allot_log.is_finished,plans.category").
		Joins("left join "+plansTable+" on task_allot_log.plan_id=plans.id").
		Where("plans.deleted_at=0 and task_allot_log.deleted_at=0 and plans.is_open = 1 and task_allot_log.user_id=? AND task_allot_log.created_at>=? and task_allot_log.created_at<=? ", userId, request.StartTime, request.EndTime)

	err = query.Order("task_allot_log.is_finished asc,plans.end_time desc,plans.id asc").Find(&data).Error
	return
}

// GetPlanTipsCount 获取指定时间内未完成的主任务数
func (d *PlanTaskAllotLogDao) GetPlanTipsCount(c app.Ctx, request planReq.TipsReq, userId string) (count int64, err error) {
	plansTable := d.dao.GetTable("plans", request.EnterpriseId)
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, request.EnterpriseId)).
		Joins("left join "+plansTable+" on task_allot_log.plan_id=plans.id").
		Where("plans.deleted_at=0 and task_allot_log.deleted_at=0 and plans.is_open = 1 and task_allot_log.is_finished=0 and task_allot_log.user_id=? AND task_allot_log.created_at>=? and task_allot_log.created_at<=?", userId, request.StartTime, request.EndTime).Group("plans.id").Count(&count).Error
	return
}

func (d *PlanTaskAllotLogDao) GetInfoById(c ctx.Context, EnterpriseId int64, Id int64) (data microScrmModel.TaskAllotLogModel, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, EnterpriseId)).Where("id=?", Id).Find(&data).Error
	return
}

func (d *PlanTaskAllotLogDao) GetTaskAllotGroupInfoById(c ctx.Context, EnterpriseId int64, Id int64) (data planReq.TaskAllotGroupInfo, err error) {
	plansTable := d.dao.GetTable("plans", EnterpriseId)
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, EnterpriseId)).Where("task_allot_log.id=?", Id).
		Joins("left join " + plansTable + " on task_allot_log.plan_id=plans.id").
		Select("plans.id as plan_id,plans.start_time,plans.end_time,task_allot_log.user_id").
		Find(&data).Error
	return
}

// GetPlanListsByPlanID 获取导购任务详情的任务状态
func (d *PlanTaskAllotLogDao) GetPlanListsByPlanID(c app.Ctx, request planReq.PlanInfo, userId string) (data []planReq.PlansListInfo, err error) {
	plansTable := d.dao.GetTable("plans", request.EnterpriseId)
	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, request.EnterpriseId)).
		Select("plans.id,plans.remark,plans.start_time,plans.end_time,plans.name,task_allot_log.is_finished").
		Joins("left join "+plansTable+" on task_allot_log.plan_id=plans.id").
		Where("plans.deleted_at=0 and task_allot_log.deleted_at=0 and plans.is_open = 1 and task_allot_log.user_id=? AND plans.id=? AND task_allot_log.created_at>=? and task_allot_log.created_at<=?", userId, request.Id, request.StartTime, request.EndTime)

	err = query.Order("task_allot_log.is_finished asc,plans.id asc").Find(&data).Error
	return
}

// GetAverage  获取任务相关平均值
func (d *PlanTaskAllotLogDao) GetAverage(c app.Ctx, enterpriseId int64, userIds []string, averageType, startTime, endTime int64) (average int64, err error) {
	var all, tmp int64
	orm := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Where("user_id in ?", userIds).
		Where("end_time < ?  ", time.Now().Unix()). // 结束时间已经到了的任务
		Where("created_at BETWEEN ? AND ?  ", startTime, endTime)
	err = orm.Count(&all).Error
	err = orm.Scopes(func(db *gorm.DB) *gorm.DB {
		switch averageType {
		case excitationModels.RuleCategoryId8:
			db.Where("is_finished=?", microScrmModel.FINISH)
			// 预留其他统计
		}
		return db
	}).Count(&tmp).Error

	average = tmp / all
	return
}

// GetAllMemberAchievement  获取任务相关统计值
func (d *PlanTaskAllotLogDao) GetAllMemberAchievement(c app.Ctx, enterpriseId int64, categoryId int64, startTime, endTime int64) (map[string]float64, error) {
	var (
		resMap   = make(map[string]float64)
		ratioRes []planReq.RatioRes
	)
	orm := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Where("created_at BETWEEN ? AND ?  ", startTime, endTime).
		Group("user_id")
	err := orm.Scopes(func(db *gorm.DB) *gorm.DB {
		switch categoryId {
		case excitationModels.RuleCategoryId8:
			db.Select("user_id,round(count(IF(is_finished=?,true,null))/count(1) * 100) as ratio", microScrmModel.FINISH).
				Having("ratio > 0")
			// 预留其他统计
		}
		return db
	}).Find(&ratioRes).Error
	if err != nil {
		return nil, err
	}

	for _, res := range ratioRes {
		resMap[res.UserId] = res.Ratio
	}
	return resMap, nil
}

// GetTopList 获取导购任务完成人数排行
func (d *PlanTaskAllotLogDao) GetTopList(c app.Ctx, enterpriseId int64, inviteTop int32, startTime, endTime int64) ([]microScrmModel.PlanTopList, error) {
	var result []microScrmModel.PlanTopList
	var lastTopTotal int64
	orm := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Select("COUNT(1) AS total").
		Where("created_at BETWEEN ? AND ? AND end_time > 0", startTime, endTime).
		Order("total desc")
	// 获取top以及和top最后一位数量相同的其他导购
	newOrm := orm.Session(&gorm.Session{NewDB: true}).WithContext(c)
	err := orm.Group("user_id").Offset(int(inviteTop)-1).Limit(1).Pluck("total", &lastTopTotal).Error
	err = newOrm.Select("user_id as member_user_id,COUNT(1) AS total").
		Group("member_user_id").
		Having("total > 0 AND total >= ?", lastTopTotal).Scan(&result).Error
	return result, err
}

// GetTopListByTouch 获取导购任务触达人数排行
func (d *PlanTaskAllotLogDao) GetTopListByTouch(c app.Ctx, enterpriseId int64, inviteTop int, startTime, endTime int64) ([]microScrmModel.PlanTopList, error) {
	var result []microScrmModel.PlanTopList
	var lastTopTotal int64
	orm := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Select("SUM(contacted_count) AS total").
		Where("created_at BETWEEN ? AND ?", startTime, endTime).
		Order("total desc")
	// 获取top以及和top最后一位数量相同的其他导购
	newOrm := orm.Session(&gorm.Session{NewDB: true}).WithContext(c)
	err := orm.Group("user_id").Offset(int(inviteTop)-1).Limit(1).Pluck("total", &lastTopTotal).Error
	err = newOrm.Select("user_id as member_user_id,SUM(contacted_count) AS total").
		Group("member_user_id").
		Having("total > 0 AND total >= ?", lastTopTotal).Scan(&result).Error
	return result, err
}

func (d *PlanTaskAllotLogDao) UpdateData(c ctx.Context, enterpriseId int64, query map[string]interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Model(&microScrmModel.TaskAllotLogModel{}).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query).Updates(data).Error
	return err
}
func (d *PlanTaskAllotLogDao) GetTaskAllotLogs(c app.Ctx, enterpriseId int64, query planReq.TaskDetailExportReq, userIds []string) ([]microScrmModel.TaskAllotLogModel, error) {
	var result []microScrmModel.TaskAllotLogModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(map[string]interface{}{
		"condition_id": query.Id,
		"deleted_at":   0,
	}).Scopes(func(db *gorm.DB) *gorm.DB {
		if query.CompletionStartTime > 0 && query.CompletionEndTime > 0 {
			db.Where("end_time BETWEEN ? AND ?", query.CompletionStartTime, query.CompletionEndTime)
		}
		if len(userIds) > 0 {
			db.Where(map[string]interface{}{"user_id": userIds})
		}
		if query.IsFinished >= 0 {
			db.Where(map[string]interface{}{"is_finished": query.IsFinished})
		}
		return db
	}).Find(&result).Error
	return result, err
}

// GetAllMemberTaskRate  获取任务完成率统计
func (d *PlanTaskAllotLogDao) GetAllMemberTaskRate(c app.Ctx, enterpriseId int64, userId []string, startTime int64, endTime int64) (map[string]float64, error) {

	var (
		resMap   = make(map[string]float64)
		ratioRes []planReq.TaskRateRes
	)
	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Where("created_at BETWEEN ? AND ?  ", startTime, endTime).
		Group("user_id").
		Select("user_id,round(count(IF(is_finished=1,true,null))/count(1) * 100) as ratio").
		Having("ratio > 0")
	if len(userId) > 0 {
		query.Where("user_id in ?", userId)
	}

	err := query.Find(&ratioRes).Error
	if err != nil {
		return nil, err
	}

	for _, res := range ratioRes {
		resMap[res.UserId] = res.Ratio
	}
	return resMap, nil
}

func (d *PlanTaskAllotLogDao) GetOne(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) (microScrmModel.TaskAllotLogModel, error) {
	var result microScrmModel.TaskAllotLogModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *PlanTaskAllotLogDao) GetFinishNumByPlanId(c ctx.Context, enterpriseId int64, userid string, planId int64) (int64, int64, error) {
	var finishTotal struct {
		Num int64 `json:"num"`
	}
	var allTotal struct {
		Num int64 `json:"num"`
	}
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where("plan_id= ? and user_id= ? and is_finished=1", planId, userid).Group("CONCAT(`plan_id`,`user_id`)").Select("count(1) as num").Scan(&finishTotal).Error
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where("plan_id= ? and user_id= ? ", planId, userid).Group("CONCAT(`plan_id`,`user_id`)").Select("count(1) as num").Scan(&allTotal).Error
	return finishTotal.Num, allTotal.Num, err

}

// GetTaskIdGroupNum 获取子任务已发送数量集合
func (d *PlanTaskAllotLogDao) GetTaskIdGroupNum(c app.Ctx, enterpriseId int64, startTime int64, endTime int64) ([]microScrmModel.TaskIdGroupNum, error) {
	var result []microScrmModel.TaskIdGroupNum
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Select("plan_id,task_id, COUNT(*) as group_msg_count,created_at").
		// Where("agent_msg_time >= ? and agent_msg_time <= ?", startTime, endTime).
		// Where("agent_msg_id <> '' ").
		Where("created_at >= ? and created_at <= ?", startTime, endTime).
		Group("task_id").
		Scan(&result).Error

	return result, err
}

// GetListsByTaskId 根据子任务id获取导购信息
func (d *PlanTaskAllotLogDao) GetListsByTaskId(c ctx.Context, enterpriseId int64, taskId int64) ([]microScrmModel.TaskAllotLogModel, error) {
	var result []microScrmModel.TaskAllotLogModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(map[string]interface{}{
		"task_id":       taskId,
		"is_group_send": microScrmModel.QyGroupMsgSend,
		// "is_finished":   microScrmModel.PLAN_FISISHING,  针对于补推任务的调整，此处不验证任务是否完成，只验证群发消息是否发送
	}).Where("group_msg_id!=''").Find(&result).Error
	return result, err

}

// GetTaskNotFinishNumByPlanId 查询出有未完成的子任务
func (d *PlanTaskAllotLogDao) GetTaskNotFinishNumByPlanId(c ctx.Context, enterpriseId int64, planId int64) ([]microScrmModel.TaskFinish, error) {
	var result []microScrmModel.TaskFinish
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Select("task_id,if(min(is_finished)=1&&max(is_finished)=1,1,0) as num").Where("plan_id= ? ", planId).Group("task_id").Find(&result).Error
	return result, err
}

func (d *PlanTaskAllotLogDao) GetCount(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) (int64, error) {
	var total int64
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query, args...).Count(&total).Error
	return total, err
}

// GetSaCount 获取触达的导购数
func (d *PlanTaskAllotLogDao) GetSaCount(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) (int64, error) {
	var total int64
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query, args...).Distinct("user_id").Count(&total).Error
	return total, err
}

func (d *PlanTaskAllotLogDao) GetMany(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]microScrmModel.TaskAllotLogModel, error) {
	var result []microScrmModel.TaskAllotLogModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query, args...).Find(&result).Error
	return result, err
}

// GetSupplyPlanIds 获取补推的主任务id
func (d *PlanTaskAllotLogDao) GetSupplyPlanIds(c ctx.Context, enterpriseId int64, startTime, endTime int64) ([]int64, error) {
	var result []int64
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where("end_time Between ? and ? and is_finished=1", startTime, endTime).Pluck("plan_id", &result).Error
	return result, err

}

// GetManyByIdPagination 通过分页id分页查询数据
func (d *PlanTaskAllotLogDao) GetManyByIdPagination(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]microScrmModel.TaskAllotLogModel, error) {
	var result []microScrmModel.TaskAllotLogModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query, args...).Limit(1000).Find(&result).Error
	return result, err
}

func (d *PlanTaskAllotLogDao) ExecRowSql(c app.Ctx, enterpriseId, updatedAt, planId, taskId int64) error {
	sql := `update micro_scrm_enterprise_%d.task_allot_log  t left join (select task_allot_log_id,count(DISTINCT(union_id)) as count from micro_scrm_enterprise_%d.task_condition_data_%d where is_finished=1 and task_id=%d GROUP BY task_allot_log_id) tt  on t.id=tt.task_allot_log_id  set t.contacted_count=tt.count,updated_at=%d where t.plan_id=? AND tt.count is not null`
	sql = fmt.Sprintf(sql, enterpriseId, enterpriseId, planId, taskId, updatedAt)
	res := d.dao.WithContext(c).Exec(sql, planId)
	return res.Error
}

func (d *PlanTaskAllotLogDao) ExecMomentRowSql(c app.Ctx, enterpriseId, updatedAt, planId, taskId int64) error {
	sql := `update micro_scrm_enterprise_%d.task_allot_log  t left join (select task_allot_log_id,count(DISTINCT(union_id)) as count from micro_scrm_enterprise_%d.task_condition_data_%d where is_touch=1 and task_id=%d GROUP BY task_allot_log_id) tt  on t.id=tt.task_allot_log_id  set t.moment_look_count=tt.count,updated_at=%d,t.finish_rate=ROUND((tt.count/t.customer_num)*100) where t.plan_id=? AND tt.count is not null`
	sql = fmt.Sprintf(sql, enterpriseId, enterpriseId, planId, taskId, updatedAt)
	res := d.dao.WithContext(c).Exec(sql, planId)
	return res.Error
}

func (d *PlanTaskAllotLogDao) CreateTaskSendMsg(c ctx.Context, enterpriseId int64, data microScrmModel.TaskAllotLogMsgModel) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.TaskAllotLogMsgModel{}.TableName(), enterpriseId)).Model(microScrmModel.TaskAllotLogMsgModel{}).Create(&data).Error
	return err
}

func (d *PlanTaskAllotLogDao) GetTaskAllotMsg(c ctx.Context, enterpriseId int64, query interface{}, args ...interface{}) ([]microScrmModel.TaskAllotLogMsgModel, error) {
	var data []microScrmModel.TaskAllotLogMsgModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.TaskAllotLogMsgModel{}.TableName(), enterpriseId)).Where(query, args...).Find(&data).Error
	return data, err
}
