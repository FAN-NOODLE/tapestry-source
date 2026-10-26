package plan_task

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

//go:generate mockgen -source=${GOFILE} -destination=mocks/mock_${GOFILE}

// Repository 存储定义
type Repository interface {
	SelectEverydayStatisticsByDay(c ctx.Context, enterpriseID int64, date string) ([]TaskEverydayStatistics, error)
	GetTaskEverydayStatistics(c ctx.Context, enterpriseID int64, taskID int64, yesterdayTs int64) (*TaskEverydayStatistics, error)
	// 获取指定区间内完成任务的集合
	GetTasksFinishDay(c ctx.Context, enterpriseID int64, endTime int64) ([]TaskAllotLog, error)
	GetTasksPortraitCustomers(c ctx.Context, enterpriseID, planId int64, taskIDs []int64, endTime int64) ([]MatchUnion, error)
	GetTasksChatMembers(c ctx.Context, enterpriseID int64, chatIDs []int64) ([]string, error)
	InsertTaskEveryDayStatistic(c ctx.Context, enterpriseID int64, taskStatistic TaskEverydayStatistics, isUpdateOrder bool) error
	GetTotalOrderStatistic(c ctx.Context, enterpriseID int64, taskID int64, isMain bool, endTime int64) (*OrderStatistic, error)
	GetSubTotalOrderStatistic(c ctx.Context, enterpriseID int64, ukVoucher string, endTime int64) (*OrderStatistic, error)
	GetSubTasksChatCustomerNumber(c ctx.Context, enterpriseID int64, chatID int64) (int64, error)
	GetSubTasksPortraitCustomerNumber(c ctx.Context, enterpriseID int64, subTaskID, planId int64, taskCategory int) (int64, error)
}
