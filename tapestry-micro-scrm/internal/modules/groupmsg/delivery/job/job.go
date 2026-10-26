package groupmsgJob

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	appJob "micro-scrm/internal/app/job"
	"micro-scrm/internal/modules/groupmsg"
	"micro-scrm/internal/modules/groupmsg/delivery"
	groupMsgMysql "micro-scrm/internal/modules/groupmsg/respository/mysql"
	"micro-scrm/internal/pkg/xjob"
)

// SendGroupMsgCronJob 发送的群信息任务
func SendGroupMsgCronJob(c ctx.Context, jobHandler *xjob.JobHandler) error {
	enterpriseIds, err := appJob.GetEnterpriseIds(jobHandler)
	if err != nil {
		return err
	}

	groupMsgService := groupmsg.New(groupMsgMysql.New())
	delivery.InitServiceDefault()
	for _, enterpriseId := range enterpriseIds {
		log.Info(c).Msgf("开始执行企业【%d】发送社群消息任务", enterpriseId)
		groupMsgService.SendGroupMsgCron(c, enterpriseId)
	}
	return nil
}

// SendResultCronJob 定时查询群发任务结果
func SendResultCronJob(c ctx.Context, jobHandler *xjob.JobHandler) error {
	enterpriseIds, err := appJob.GetEnterpriseIds(jobHandler)
	if err != nil {
		return err
	}

	groupMsgService := groupmsg.New(groupMsgMysql.New())
	delivery.InitServiceDefault()

	for _, enterpriseId := range enterpriseIds {
		// 执行策略
		err = groupMsgService.CheckEnterpriseTaskSendResult(c, enterpriseId)
		if err != nil {
			log.Info(c).Msgf("开始执行企业【%d】查询群发任务结果,出错：%+v", enterpriseId, err)
		}
	}
	return nil
}
