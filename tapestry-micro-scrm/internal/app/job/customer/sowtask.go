package customer

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	customerDao "micro-scrm/internal/app/dao/customer"
	appJob "micro-scrm/internal/app/job"
	"micro-scrm/internal/app/models/customer"
	cacheSvc "micro-scrm/internal/app/service/cache"
	messageSvc "micro-scrm/internal/app/service/message"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sdk/microwork"
	"micro-scrm/internal/pkg/xjob"
	"time"
)

type SowTask struct {
	name       string
	sowTaskDao *customerDao.SowTask
}

func NewSowTask() *SowTask {
	return &SowTask{
		name:       "通知播种任务",
		sowTaskDao: customerDao.NewSowTask(),
	}
}

func (job *SowTask) Run(c ctx.Context, handler *xjob.JobHandler) error {
	enterpriseIds, err := appJob.GetEnterpriseIds(handler)
	if err != nil {
		return err
	}

	var wg threading.RoutineGroup

	for _, id := range enterpriseIds {
		tmpId := id
		cacheRes, err := cacheSvc.NewRedisData().GetEnterpriseConfig(c, tmpId, "store_name")
		if err != nil {
			continue
		}
		wg.Run(func() {
			job.HandleUnFinished(c, tmpId, cacheRes.Name)
		})
		wg.Run(func() {
			job.Handle(c, tmpId, cacheRes.Name)
		})
	}

	wg.Wait()

	return nil
}

func (job *SowTask) HandleUnFinished(c ctx.Context, enterpriseId int64, enterpriseName string) {
	// 获取昨天凌晨时间戳
	startTimestamp, endTimestamp := helper.YesterdayPeriodOfTimestamp(1)

	// 获取企业所有未完成的播种任务
	query := fmt.Sprintf("enterprise_id = %d and status = %d and remind_time >= %d and remind_time <= %d",
		enterpriseId,
		customerModels.SowTaskStatusUnfinished,
		startTimestamp,
		endTimestamp,
	)
	list, err := job.sowTaskDao.GetMany(c, query)
	if err != nil || len(list) == 0 {
		return
	}

	membersMap := make(map[string]int)

	for _, value := range list {
		tmpVal := value
		// 获取昨天未完成任务的导购和数量
		val, ok := membersMap[tmpVal.MemberUserID]
		if !ok {
			membersMap[tmpVal.MemberUserID] = 1
		} else {
			membersMap[tmpVal.MemberUserID] = val + 1
		}
	}

	countTitle, timeTitle, messageTitle := job.getPushTitle(enterpriseName, 0)

	wg := threading.RoutineGroup{}
	for key, val := range membersMap {
		memberUserId := key
		count := val
		wg.Run(func() {
			job.PushMessage(c, pushReq{
				EnterpriseId:   enterpriseId,
				MemUserId:      memberUserId,
				CountTitle:     countTitle,
				Count:          count,
				TimeTitle:      timeTitle,
				TimeFormat:     time.Unix(endTimestamp, 0).Format("2006年01月02日"),
				MessageTitle:   messageTitle,
				StartTimestamp: startTimestamp,
			})
		})
	}
	wg.Wait()
}

func (job *SowTask) Handle(c ctx.Context, enterpriseId int64, enterpriseName string) {
	// 获取今天凌晨时间戳
	startTimestamp, endTimestamp := helper.PeriodOfTimestamp()

	// 获取企业今日所有播种任务
	query := fmt.Sprintf("enterprise_id = %d and status = %d and remind_time >= %d and remind_time <= %d",
		enterpriseId,
		customerModels.SowTaskStatusUnfinished,
		startTimestamp,
		endTimestamp,
	)
	list, err := job.sowTaskDao.GetMany(c, query)
	if err != nil || len(list) == 0 {
		return
	}

	membersMap := make(map[string]int)

	for _, value := range list {
		tmpVal := value
		// 获取今天未完成任务的导购和数量
		val, ok := membersMap[tmpVal.MemberUserID]
		if !ok {
			membersMap[tmpVal.MemberUserID] = 1
		} else {
			membersMap[tmpVal.MemberUserID] = val + 1
		}
	}

	countTitle, timeTitle, messageTitle := job.getPushTitle(enterpriseName, 1)

	wg := threading.RoutineGroup{}
	for key, val := range membersMap {
		memberUserId := key
		count := val
		wg.Run(func() {
			job.PushMessage(c, pushReq{
				EnterpriseId:   enterpriseId,
				MemUserId:      memberUserId,
				CountTitle:     countTitle,
				Count:          count,
				TimeTitle:      timeTitle,
				TimeFormat:     time.Now().Format("2006年01月02日"),
				MessageTitle:   messageTitle,
				StartTimestamp: startTimestamp,
			})
		})
	}
	wg.Wait()
}

type pushReq struct {
	EnterpriseId   int64
	MemUserId      string
	CountTitle     string // 数量标题
	Count          int
	TimeTitle      string // 时间标题
	TimeFormat     string // 时间格式化
	MessageTitle   string // 主题标题
	StartTimestamp int64  // 当前开始时间
}

func (job *SowTask) PushMessage(c ctx.Context, req pushReq) {
	content := make([]map[string]string, 0)
	content = append(content, map[string]string{"key": req.CountTitle, "value": helper.Int64ToString(int64(req.Count))})
	content = append(content, map[string]string{"key": req.TimeTitle, "value": req.TimeFormat})

	params := microwork.SendMiniMessageRequest{
		EnterpriseID: req.EnterpriseId,
		MemUserIds:   []string{req.MemUserId},
		Title:        req.MessageTitle,
		Content:      content,
		Page:         fmt.Sprintf("pages/task/task.html?type=1&time=%d", req.StartTimestamp),
	}

	_, _ = messageSvc.NewMessageService().SendMiniMessage(c, params)

	// _, _ = microwork.SendMiniMessage(c, &params)
}

// pushType 0 未完成 1 待完成
func (job *SowTask) getPushTitle(enterpriseName string, pushType int) (string, string, string) {
	if pushType == 0 {
		if enterpriseName == "coach" {
			return "未完成播种数", "播种时间", "昨日未完成播种任务提醒"
		}
		if helper.InArray(enterpriseName, []string{"sw", "ks"}) {
			return "未完成任务数", "提醒时间", "昨日未完成自设任务提醒"
		}
	}

	if enterpriseName == "coach" {
		return "待播种数", "播种时间", "今日播种任务提醒"
	}
	if helper.InArray(enterpriseName, []string{"sw", "ks"}) {
		return "待完成任务数", "提醒时间", "今日自设任务提醒"
	}

	return "", "", ""
}
