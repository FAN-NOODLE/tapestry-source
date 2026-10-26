package customer

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	customerDao "micro-scrm/internal/app/dao/customer"
	appJob "micro-scrm/internal/app/job"
	customerModels "micro-scrm/internal/app/models/customer"
	cacheSvc "micro-scrm/internal/app/service/cache"
	messageSvc "micro-scrm/internal/app/service/message"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sdk/microwork"
	"micro-scrm/internal/pkg/xjob"
	"time"
)

type TimingSowTask struct {
	name       string
	sowTaskDao *customerDao.SowTask
}

func NewTimingSowTask() *TimingSowTask {
	return &TimingSowTask{
		name:       "定时通知播种任务",
		sowTaskDao: customerDao.NewSowTask(),
	}
}

func (job *TimingSowTask) Run(c ctx.Context, handler *xjob.JobHandler) error {
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
			job.handle(c, tmpId, cacheRes.Name)
		})
	}

	wg.Wait()

	return nil
}

func (job *TimingSowTask) handle(c ctx.Context, enterpriseId int64, enterpriseName string) {
	// 获取半小时后5分钟误差的时间戳
	nowTime := time.Now()
	startTimestamp, endTimestamp := nowTime.Add(30*time.Minute).Unix(), nowTime.Add(35*time.Minute).Unix()

	// 获取企业所有未完成的播种任务
	sowTaskList, err := job.sowTaskDao.GetMany(c,
		"enterprise_id = ? and status = ? and remind_time >= ? and remind_time <= ? and is_push = ?",
		enterpriseId,
		customerModels.SowTaskStatusUnfinished,
		startTimestamp,
		endTimestamp,
		customerModels.SowTaskIsPushFalse)
	if err != nil || len(sowTaskList) == 0 {
		return
	}

	var sowTaskIds []int64
	memberTimeMap := make(map[string]map[int64]int)
	for _, value := range sowTaskList {
		sowTaskIds = append(sowTaskIds, value.ID)

		remindTimeMap, ok := memberTimeMap[value.MemberUserID]
		if !ok {
			memberTimeMap[value.MemberUserID] = map[int64]int{value.RemindTime: 1}
			continue
		}

		remindTimeVal, remindTimeOk := remindTimeMap[value.RemindTime]
		if !remindTimeOk {
			memberTimeMap[value.MemberUserID] = map[int64]int{value.RemindTime: 1}
			continue
		}

		memberTimeMap[value.MemberUserID] = map[int64]int{value.RemindTime: remindTimeVal + 1}
	}

	countTitle, timeTitle, messageTitle := job.getPushTitle(enterpriseName)

	for memberUserId, remindTimeMap := range memberTimeMap {
		for remindTime, count := range remindTimeMap {
			job.PushMessage(c, pushReq{
				EnterpriseId:   enterpriseId,
				MemUserId:      memberUserId,
				CountTitle:     countTitle,
				Count:          count,
				TimeTitle:      timeTitle,
				TimeFormat:     time.Unix(remindTime, 0).Format("2006-01-02 15:04"),
				MessageTitle:   messageTitle,
				StartTimestamp: startTimestamp,
			})
		}
	}

	// 更新为已发送
	_ = job.sowTaskDao.UpdateData(c,
		map[string]any{"enterprise_id": enterpriseId, "id": sowTaskIds},
		map[string]any{"is_push": customerModels.SowTaskIsPushTrue})
}

func (job *TimingSowTask) PushMessage(c ctx.Context, req pushReq) {
	content := make([]map[string]string, 0)
	content = append(content, map[string]string{"key": req.CountTitle, "value": helper.Int64ToString(int64(req.Count))})
	content = append(content, map[string]string{"key": req.TimeTitle, "value": req.TimeFormat})

	// 将时间戳转换为对应的时间对象
	timeObj := time.Unix(req.StartTimestamp, 0)
	// 去除时分秒，设置为零点
	truncatedTime := time.Date(timeObj.Year(), timeObj.Month(), timeObj.Day(), 0, 0, 0, 0, timeObj.Location())

	params := microwork.SendMiniMessageRequest{
		EnterpriseID: req.EnterpriseId,
		MemUserIds:   []string{req.MemUserId},
		Title:        req.MessageTitle,
		Content:      content,
		Page:         fmt.Sprintf("pages/task/task.html?type=1&time=%d", truncatedTime.Unix()),
	}

	_, _ = messageSvc.NewMessageService().SendMiniMessage(c, params)
}

// pushType 0 未完成 1 待完成
func (job *TimingSowTask) getPushTitle(enterpriseName string) (string, string, string) {
	if enterpriseName == "coach" {
		return "待播种数", "播种时间", "有播种任务在30分钟后生效"
	}
	if helper.InArray(enterpriseName, []string{"sw", "ks"}) {
		return "待完成任务数", "提醒时间", "自设任务提醒"
	}

	return "", "", ""
}
