package mysql

import (
	"database/sql"
	"fmt"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"time"

	"micro-scrm/internal/pkg/db"

	"github.com/jmoiron/sqlx"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/plan_task"
)

type PlanTaskRepo struct {
	*repos.SQLRepo
}

func (p *PlanTaskRepo) SelectEverydayStatisticsByDay(c ctx.Context, enterpriseID int64, date string) ([]plan_task.TaskEverydayStatistics, error) {
	var list []plan_task.TaskEverydayStatistics
	query := fmt.Sprintf("select task_id,total_reach_customer_number,total_order_customer_number,total_order_number,total_order_fee from micro_scrm_enterprise_%d.task_everyday_statistics where date=?", enterpriseID)
	err := p.DB.Select(c, &list, query, date)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (p *PlanTaskRepo) GetTaskEverydayStatistics(c ctx.Context, enterpriseID int64, taskID int64, yesterdayTs int64) (*plan_task.TaskEverydayStatistics, error) {
	var info plan_task.TaskEverydayStatistics
	// 查询该距离改日期前一天时间戳最近的记录
	query := fmt.Sprintf("select task_id,total_reach_customer_number,total_order_customer_number,total_order_number,total_order_fee from micro_scrm_enterprise_%d.task_everyday_statistics where task_id=? AND data_time<=? order by data_time desc limit 1", enterpriseID)
	err := p.DB.DB.Get(c, &info, query, taskID, yesterdayTs)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (p *PlanTaskRepo) GetTasksFinishDay(c ctx.Context, enterpriseID int64, endTime int64) ([]plan_task.TaskAllotLog, error) {
	var list []plan_task.TaskAllotLog
	query := fmt.Sprintf("select b.id,b.plan_id,b.task_id,b.group_chat_id as chat_id,b.user_id from micro_scrm_enterprise_%d.task_allot_log b,micro_scrm_enterprise_%d.plans a where b.plan_id=a.id AND b.end_time<=? AND b.is_finished=1 AND b.deleted_at=0 AND a.deleted_at=0", enterpriseID, enterpriseID)
	err := p.DB.Select(c, &list, query, endTime)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (p *PlanTaskRepo) GetTasksChatMembers(c ctx.Context, enterpriseID int64, chatIDs []int64) ([]string, error) {
	var list []string
	query := "select DISTINCT(unionid) as un from work_group_chat_member where enterprise_id = ? AND group_chat_id in (?) AND type=2 AND leave_time=0 AND deleted_at=0"
	query, args, err := sqlx.In(query, enterpriseID, chatIDs)
	if err != nil {
		return nil, err
	}
	err = p.DB.Select(c, &list, query, args...)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (p *PlanTaskRepo) GetSubTasksChatCustomerNumber(c ctx.Context, enterpriseID int64, chatID int64) (int64, error) {
	var num int64
	query := fmt.Sprintf("select count(DISTINCT(unionid)) from work_group_chat_member where enterprise_id = ? AND group_chat_id =? AND type=2 AND leave_time=0 AND deleted_at=0")

	err := p.DB.DB.Get(c, &num, query, enterpriseID, chatID)
	if err != nil {
		return num, err
	}
	return num, nil
}

func (p *PlanTaskRepo) GetSubTasksPortraitCustomerNumber(c ctx.Context, enterpriseID int64, subTaskID, planId int64, taskCategory int) (int64, error) {
	var num int64
	sqlStr := "select count(DISTINCT union_id,customer_id) from micro_scrm_enterprise_%d.task_condition_data_%d where task_allot_log_id = ? and is_finished=1"
	if taskCategory == microScrmModel.PLANTASKMOMENT { // 朋友圈任务
		sqlStr = "select count(DISTINCT union_id,customer_id ) from micro_scrm_enterprise_%d.task_condition_data_%d where task_allot_log_id = ? and is_touch=1"
	}
	query := fmt.Sprintf(sqlStr, enterpriseID, planId)
	err := p.DB.DB.Get(c, &num, query, subTaskID)
	if err != nil {
		return num, err
	}
	return num, nil
}

func (p *PlanTaskRepo) GetTasksPortraitCustomers(c ctx.Context, enterpriseID, planId int64, taskIDs []int64, endTime int64) ([]plan_task.MatchUnion, error) {
	var list []plan_task.MatchUnion
	query := fmt.Sprintf("select a.task_id,a.union_id from micro_scrm_enterprise_%d.task_condition_data_%d a, micro_scrm_enterprise_%d.task_allot_log b where a.task_allot_log_id = b.id AND b.task_id in (?) AND b.end_time<=? AND b.is_finished=1 and a.is_finished=1 AND b.deleted_at=0", enterpriseID, planId, enterpriseID)
	query, args, err := sqlx.In(query, taskIDs, endTime)
	if err != nil {
		return nil, err
	}
	err = p.DB.Select(c, &list, query, args...)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (p *PlanTaskRepo) InsertTaskEveryDayStatistic(c ctx.Context, enterpriseID int64, taskStatistic plan_task.TaskEverydayStatistics, isUpdateOrder bool) error {
	if isUpdateOrder {
		query := db.MakeUpsertSQL(fmt.Sprintf("micro_scrm_enterprise_%d.task_everyday_statistics", enterpriseID), []db.UpsertColumn{
			{Name: "date", Insert: true},
			{Name: "data_time", Insert: true},
			{Name: "task_id", Insert: true},
			{Name: "total_reach_customer_number", Insert: true, Update: true},
			{Name: "total_order_customer_number", Insert: true, Update: true},
			{Name: "total_order_number", Insert: true, Update: true},
			{Name: "total_order_fee", Insert: true, Update: true},
			{Name: "task_type", Insert: true},
			{Name: "created_at", Insert: true},
			{Name: "updated_at", Insert: true, Update: true},
		})
		_, err := p.DB.DB.Exec(c, query, taskStatistic.Date, taskStatistic.DataTime, taskStatistic.TaskID, taskStatistic.TotalReachCustomerNumber, taskStatistic.TotalOrderCustomerNumber, taskStatistic.TotalOrderNumber, taskStatistic.TotalOrderFee, taskStatistic.TaskType, time.Now().Unix(), time.Now().Unix(), taskStatistic.TotalReachCustomerNumber, taskStatistic.TotalOrderCustomerNumber, taskStatistic.TotalOrderNumber, taskStatistic.TotalOrderFee, time.Now().Unix())
		return err
	} else {
		query := db.MakeUpsertSQL(fmt.Sprintf("micro_scrm_enterprise_%d.task_everyday_statistics", enterpriseID), []db.UpsertColumn{
			{Name: "date", Insert: true},
			{Name: "data_time", Insert: true},
			{Name: "task_id", Insert: true},
			{Name: "total_reach_customer_number", Insert: true, Update: true},
			{Name: "total_order_customer_number", Insert: true},
			{Name: "total_order_number", Insert: true},
			{Name: "total_order_fee", Insert: true},
			{Name: "task_type", Insert: true},
			{Name: "created_at", Insert: true},
			{Name: "updated_at", Insert: true, Update: true},
		})
		_, err := p.DB.DB.Exec(c, query, taskStatistic.Date, taskStatistic.DataTime, taskStatistic.TaskID, taskStatistic.TotalReachCustomerNumber, taskStatistic.TotalOrderCustomerNumber, taskStatistic.TotalOrderNumber, taskStatistic.TotalOrderFee, taskStatistic.TaskType, time.Now().Unix(), time.Now().Unix(), taskStatistic.TotalReachCustomerNumber, time.Now().Unix())
		return err
	}

}

func (p *PlanTaskRepo) GetTotalOrderStatistic(c ctx.Context, enterpriseID int64, taskID int64, isMain bool, endTime int64) (*plan_task.OrderStatistic, error) {
	var info plan_task.OrderStatistic
	query := fmt.Sprintf("select count(distinct(a.unionid)) AS total_order_customer_number,count(a.order_no) as total_order_number,IFNULL(sum(a.order_fee),0) AS total_order_fee from micro_scrm_enterprise_%d.orders a,micro_scrm_enterprise_%d.ws_trace b where a.ws_trace_id=b.ws_trace_id AND b.main_task_id = ? AND a.order_pay_time<=?", enterpriseID, enterpriseID)
	if isMain == false {
		query = fmt.Sprintf("select count(distinct(a.unionid)) AS total_order_customer_number,count(a.order_no) as total_order_number,IFNULL(sum(a.order_fee),0) AS total_order_fee from micro_scrm_enterprise_%d.orders a,micro_scrm_enterprise_%d.ws_trace b where a.ws_trace_id=b.ws_trace_id AND b.child_task_id = ? AND a.order_pay_time<=?", enterpriseID, enterpriseID)
	}
	err := p.DB.DB.Get(c, &info, query, taskID, endTime)
	if err != nil && err != sql.ErrNoRows {
		return &info, err
	}
	return &info, nil
}

func (p *PlanTaskRepo) GetSubTotalOrderStatistic(c ctx.Context, enterpriseID int64, ukVoucher string, endTime int64) (*plan_task.OrderStatistic, error) {
	var info plan_task.OrderStatistic
	query := fmt.Sprintf("select count(distinct(a.unionid)) AS total_order_customer_number,count(a.order_no) as total_order_number,IFNULL(sum(a.order_fee),0) AS total_order_fee from micro_scrm_enterprise_%d.orders a,micro_scrm_enterprise_%d.ws_trace b where a.ws_trace_id=b.ws_trace_id AND b.uk_voucher = ? AND a.order_pay_time<=?", enterpriseID, enterpriseID)
	err := p.DB.DB.Get(c, &info, query, ukVoucher, endTime)
	if err != nil && err != sql.ErrNoRows {
		return &info, err
	}
	return &info, nil
}
