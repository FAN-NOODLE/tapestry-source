package sourceJob

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/app"
	membersDao "micro-scrm/internal/app/dao/members"
	microBossDao "micro-scrm/internal/app/dao/micro_boss"
	sourceRequest "micro-scrm/internal/app/http/request/source"
	cacheSvc "micro-scrm/internal/app/service/cache"
	messageSvc "micro-scrm/internal/app/service/message"
	sourceSvc "micro-scrm/internal/app/service/source"
	customerService "micro-scrm/internal/modules/customer"
	"micro-scrm/internal/modules/enterprise"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/boss"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sdk/microwork"
	"micro-scrm/internal/pkg/xjob"
	"sync"
	"time"
)

type PushDailyUpdateSourceData struct {
	name              string
	enterpriseService enterprise.Service
	customerService   customerService.MessageService
}

func NewPushDailyUpdateSourceData() *PushDailyUpdateSourceData {

	return &PushDailyUpdateSourceData{
		name: "推送每日更新素材",
	}
}

func (sourceJob PushDailyUpdateSourceData) Run(ctx ctx.Context, handler *xjob.JobHandler) error {

	var (
		err            error
		enterpriseList []boss.EnterpriseModel
	)
	enterpriseList, err = orm.BossService.Enterprise.GetAllInfo()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, enterpriseInfo := range enterpriseList {
		wg.Add(1)

		go func(enterpriseID int64) {
			defer wg.Done()

			sourceJob.SendMessage(ctx, enterpriseID)

		}(enterpriseInfo.ID)

	}
	wg.Wait()
	return nil

}
func (sourceJob PushDailyUpdateSourceData) SendMessage(ctx ctx.Context, enterpriseId int64) error {
	config, err := cacheSvc.NewRedisData().GetEnterpriseConfig(ctx, enterpriseId, "store_name")
	if err != nil {
		return err
	}
	if config.Name == "coach" || config.Name == "sw" {
		return sourceJob.sendMiniProgram(ctx, enterpriseId)
	} else {
		return sourceJob.SendTextMessage(ctx, enterpriseId)
	}
}

// 发送小程序消息
func (sourceJob PushDailyUpdateSourceData) sendMiniProgram(ctx ctx.Context, enterpriseId int64) error {

	// 需要验证企业配置调用enterprise_config
	res, err := microBossDao.NewEnterpriseConfig().GetInfo(ctx,
		map[string]interface{}{
			"enterprise_id": enterpriseId,
			"key":           "source_send_text_message",
		},
	)
	if err != nil {
		return err
	}
	if res.ID == 0 {
		return nil
	}

	// 昨天的起止时间戳
	startTime, endTime := helper.YesterdayPeriodOfTimestamp(1)
	// 查询需要推送消息的导购 导购可查看素材权限的新素材条数
	sendNewSourceList, err1 := sourceSvc.NewService().GetNewSourceCountForMember(ctx, enterpriseId, startTime, endTime)

	if err1 != nil {
		return err1
	}

	// 处理相同数量的门店
	sendMessageThirdGroupMap := make(map[int64]sourceRequest.SendMessageThirdGroup)
	for _, v := range sendNewSourceList {
		sendMessageThirdGroupMap[v.IdCount] = sourceRequest.SendMessageThirdGroup{
			ThirdIds:  append(sendMessageThirdGroupMap[v.IdCount].ThirdIds, v.ThirdId),
			StartTime: v.StartTime,
			IdCount:   v.IdCount,
		}
	}
	for _, v := range sendMessageThirdGroupMap {
		if len(v.ThirdIds) == 0 {
			log.Info(ctx).Msgf("企业%d未配置文本消息发送人群", enterpriseId)
			continue
		}
		var userIds []string
		// 根据门店获取要发送的员工
		memberList, _ := membersDao.NewMembers().GetMemberByThirdIds(ctx, enterpriseId, v.ThirdIds, []int64{})
		for _, member := range memberList {
			userIds = append(userIds, member.MembersUserid)
		}
		tm := time.Unix(v.StartTime, 0)
		content := make([]map[string]string, 0)
		content = append(content, map[string]string{"key": "素材数量", "value": helper.Int64ToString(v.IdCount)})
		content = append(content, map[string]string{"key": "生效时间", "value": tm.Format("2006-01-02 15:04:05")})
		params := microwork.SendMiniMessageRequest{
			EnterpriseID: enterpriseId,
			MemUserIds:   []string{},
			Title:        "素材更新提醒",
			Content:      content,
			Page:         "/pagesMaterial/pages/list/list",
		}
		var sendUserIds []string
		for k, userId := range userIds {
			sendUserIds = append(sendUserIds, userId)
			if len(sendUserIds) == 100 || len(userIds) == k+1 {
				params.MemUserIds = sendUserIds
				_, err = messageSvc.NewMessageService().SendMiniMessage(ctx, params)
				app.Logger().Info(ctx).Err(err).Interface("params", params).Msg("sendPlanTaskMiniMessage.SendMiniMessageRequest")
				sendUserIds = []string{}
			}
		}
	}

	return nil
}

// SendTextMessage 发送侧边栏文本卡片消息
func (sourceJob PushDailyUpdateSourceData) SendTextMessage(ctx ctx.Context, enterpriseId int64) error {

	corpReq := enterprise.Corp{
		EnterpriseID: enterpriseId,
	}
	// 查询企业信息
	enterpriseData, enterpriseErr := sourceJob.enterpriseService.Detail(ctx, corpReq)
	if enterpriseErr != nil {
		return enterpriseErr
	}

	// 服务发送
	params := sourceRequest.SendRequest{
		EnterpriseID: enterpriseId,
	}
	message, err := sourceSvc.NewService().SendTextMessage(ctx, params, enterpriseData)
	if err != nil {
		return err
	}
	msgRes := fmt.Sprintf("SendTextMessage.res:%v", message)
	app.Logger().Info(ctx).Err(err).Interface("params", params).Msg(msgRes)
	return nil
}
