package plan_tasks

import (
	"encoding/json"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	membersDao "micro-scrm/internal/app/dao/members"
	planDao "micro-scrm/internal/app/dao/micro_scrm"
	planReq "micro-scrm/internal/app/http/request/plan"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	cacheSvc "micro-scrm/internal/app/service/cache"
	memberSvc "micro-scrm/internal/app/service/member"
	planService "micro-scrm/internal/app/service/plan"
	"micro-scrm/internal/modules/groupmsg"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/boss"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sdk/work"
	"micro-scrm/internal/pkg/xjob"
	"sync"
)

type SupplyPushGroupMsg struct {
	Name string
}

func NewSupplyPushGroupMsg() SupplyPushGroupMsg {
	return SupplyPushGroupMsg{
		Name: "企业群发任务，顾客的群发消息当日触达上线，次日补推群发消息，每天16点推送",
	}
}

func (job SupplyPushGroupMsg) Run(c ctx.Context, handler *xjob.JobHandler) error {
	var (
		err            error
		enterpriseList []boss.EnterpriseModel
	)
	enterpriseList, err = orm.BossService.Enterprise.GetAllInfo()
	if err != nil {
		return err
	}

	for _, v := range enterpriseList {
		enterpriseId := v.ID
		err = job.handle(c, enterpriseId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (job SupplyPushGroupMsg) handle(c ctx.Context, enterpriseId int64) error {
	var err error

	resConfig, _ := cacheSvc.NewRedisData().GetEnterpriseConfig(c, enterpriseId, "store_name")
	if resConfig.Name != "coach" {
		return nil
	}
	// 补推逻辑调整，只补推一次，今天补推任务完成时间在昨天的任务，补推的任务群发消息可以选择外部联系人
	startTime, endTime := helper.YesterdayPeriodOfTimestamp(1)
	var planIds []int64
	planIds, err = planDao.NewPlanTaskAllotLogDao().GetSupplyPlanIds(c, enterpriseId, startTime, endTime)
	if err != nil {
		return err
	}
	// 获取进行中的任务
	many, err := planDao.NewPlanTaskDao().GetPlanMany(c, enterpriseId, map[string]interface{}{
		"id":     planIds,
		"status": 1,
	})
	var sendPlanIds []int64 // 真正满足补推发送的主任务ID
	for _, model := range many {
		sendPlanIds = append(sendPlanIds, model.ID)
	}
	taskLists, err := planService.NewPlanTask().GetQyGroupListsByPlanId(c, enterpriseId, microScrmModel.PLANTASKQYSEND, sendPlanIds)
	if err != nil {
		return err
	}
	if len(taskLists) == 0 {
		return err
	}

	enterpriseInfo, err := common.GetEnterpriseInfo(c, enterpriseId)
	if err != nil {
		return err
	}

	memberLists, err := membersDao.NewMembers().GetMany(c, map[string]interface{}{
		"enterprise_id": enterpriseId,
		"deleted_at":    0,
		"status":        1,
	})
	if err != nil {
		return err
	}
	var userIds []string
	for _, list := range memberLists {
		item := list
		userIds = append(userIds, item.Userid)
	}

	p, _ := ants.NewPool(1)
	defer p.Release()
	for _, v := range taskLists {
		taskId := v.ID
		planId := v.PlanID
		var userIdPathMap = make(map[string]string)
		touchLists, err := planService.NewPlanTaskConditionData().GetNoTouchLists(c, enterpriseId, planId, taskId)
		if err != nil {
			continue
		}
		// 获取企微API
		conditions, err2 := planDao.NewPlanTaskDao().GetConditionInfo(c, enterpriseId, map[string]interface{}{
			"task_id": taskId,
		})
		if err2 != nil {
			continue
		}
		var conditionContent planReq.ConditionContent
		err = json.Unmarshal([]byte(conditions.Content), &conditionContent)
		if err != nil {
			continue
		}
		qyWorkApi := new(work.QyWorkImpl).SetDomain(enterpriseInfo.Domain)
		var mediaContent work.MediaContentData
		mediaContent.Text = &work.Text{Content: conditionContent.Text.Content}
		if len(conditionContent.Image.ImgList) > 0 {
			for _, v := range conditionContent.Image.ImgList {
				item := v
				info := work.Image{
					MediaId: "",
					PicUrl:  item.PicUrl,
				}
				mediaContent.Content = append(mediaContent.Content, work.MediaContent{Image: &info})
			}
		}
		if conditionContent.Link.Title != "" {
			mediaContent.Content = append(mediaContent.Content, work.MediaContent{Link: &work.Link{
				Title:    conditionContent.Link.Title,
				WxPicUrl: conditionContent.Link.CoverUrl,
				Desc:     conditionContent.Link.Description,
				Url:      conditionContent.Link.LinkUrl,
				PicUrl:   conditionContent.Link.PicUrl,
			}})
		}
		if conditionContent.Video.VideoUrl != "" {
			mediaContent.Content = append(mediaContent.Content, work.MediaContent{Video: &work.Video{
				MediaId:  "",
				CoverUrl: conditionContent.Video.VideoUrl,
				Url:      conditionContent.Video.VideoUrl,
			}})
		}
		if conditionContent.Applets.AppId != "" {
			mediaContent.Content = append(mediaContent.Content, work.MediaContent{MiniProgram: &work.MiniProgram{
				Title:      conditionContent.Applets.Title,
				PicMediaId: "",
				Appid:      conditionContent.Applets.AppId,
				Page:       conditionContent.Applets.Path,
				CoverUrl:   conditionContent.Applets.CoverUrl,
			}})
			if resConfig.Name == "sw" {
				userIds = helper.ArrayUniq(userIds)
				miniPathMap, err := memberSvc.NewMembersService().SetMiniProgramLinkParamsByTask(c, enterpriseId, userIds, conditionContent.Applets.Path)
				if err != nil {
					return err
				}
				for userId, path := range miniPathMap {
					userIdPathMap[userId] = path
				}
			}
		}
		mediaContentStr := helper.MarshalToString(mediaContent)
		groupSendContent, err := qyWorkApi.SetMediaID(c, mediaContentStr, enterpriseId)
		if err != nil {
			fmt.Printf("群发消息获取素材库id报错%v\n", err)
			continue
		}
		sendUserLists := make(map[string][]string, 0)
		needSendUserLists := make(map[string]int64, 0)
		for _, list := range touchLists {
			item := list
			if !helper.InArray(item.UserId, userIds) {
				continue
			}
			if item.WorkExternalUserId != "" {
				sendUserLists[item.UserId] = append(sendUserLists[item.UserId], list.WorkExternalUserId)
				if _, ok := needSendUserLists[item.UserId]; !ok {
					needSendUserLists[item.UserId] = item.TaskAllotLogId
				}
			}
		}
		taskAllotLogIds := make([]int64, 0)
		for _, v := range needSendUserLists {
			id := v
			taskAllotLogIds = append(taskAllotLogIds, id)
		}
		taskAllotLists, err := planDao.NewPlanTaskAllotLogDao().GetMany(c, enterpriseId, "id in (?) and end_time between ? and ? and is_group_send= ? ", taskAllotLogIds, startTime, endTime, microScrmModel.QyGroupMsgSendFinish)
		if err != nil {
			return err
		}
		sendUsers := make([]string, 0)
		for _, value := range taskAllotLists {
			sendUsers = append(sendUsers, value.UserId)
		}
		wg := sync.WaitGroup{}
		for k, lists := range sendUserLists {
			item := lists
			if len(item) == 0 {
				continue
			}
			if !helper.InArray(k, sendUsers) { // 排除未发送的群发消息导购
				continue
			}
			wg.Add(1)
			fn := func(ctx ctx.Context, enterpriseId, taskId int64, userId string, externalUserIds []string, qyWorkApi *work.QyWorkImpl) func() {
				return func() {
					defer wg.Done()
					job.qyGroupSend(ctx, enterpriseId, taskId, userId, externalUserIds, groupSendContent, enterpriseInfo.Domain, userIdPathMap)
				}
			}
			err := p.Submit(fn(c, enterpriseId, taskId, k, item, qyWorkApi))
			if err != nil {
				fmt.Printf("执行补推企微企业群发失败，err:%v", err)
				continue
			}

		}
		wg.Wait()
	}
	fmt.Printf("执行补推企微群发任务信息完成")
	return nil
}

// 企业群发
func (job SupplyPushGroupMsg) qyGroupSend(ctx ctx.Context, enterpriseId, taskId int64, userId string, externalUserIds []string, groupSendContent string, domain string, userIdPathMap map[string]string) {
	var massHairReq = work.SendMsgTemplateReq{}
	err := json.Unmarshal([]byte(groupSendContent), &massHairReq)
	if err != nil {
		return
	}
	massHairReq.Sender = userId
	massHairReq.ChatType = groupmsg.ChatTypeUser
	massHairReq.ExternalUserid = externalUserIds
	massHairReq.AllowSelect = true
	for k, attachment := range massHairReq.Attachments {
		if attachment.MiniProgram != nil {
			if attachment.MiniProgram.Page != "" {
				if path, ok := userIdPathMap[userId]; ok {
					massHairReq.Attachments[k].MiniProgram.Page = path
				}
			}
		}
	}
	params := massHairReq
	qyWorkApi := new(work.QyWorkImpl).SetDomain(domain)
	response, err := qyWorkApi.SendMsgTemplate(ctx, massHairReq)
	if err != nil {
		log.Info(ctx).Msgf("发送群发消息失败, err:%s", err.Error())
		return
	}
	if response.ErrCode != 0 {
		log.Info(ctx).Msgf("发送群发消息企微返回错误, response:%v", response)
		return
	}
	log.Info(ctx).Err(err).Interface("params", params).Msg("qyGroupSend.qyGroupSend")
	msgId := response.MsgID
	fmt.Printf("taskId:%d----导购id:%s----任务企业群发的--msgID:%s\n", taskId, params.Sender, msgId)
	err = planDao.NewPlanTaskAllotLogDao().UpdateData(ctx, enterpriseId, map[string]interface{}{
		"task_id": taskId,
		"user_id": params.Sender,
	}, map[string]interface{}{"group_msg_id": msgId, "is_group_send": microScrmModel.QyGroupMsgSend})
	fmt.Printf("taskId:%d----导购id:%s----更新发送状态失败:%v\n", taskId, params.Sender, err)
	failList := response.FailList
	one, err := planDao.NewPlanTaskAllotLogDao().GetOne(ctx, enterpriseId, map[string]interface{}{
		"task_id": taskId,
		"user_id": params.Sender,
	})
	if err != nil {
		return
	}
	if len(failList) > 0 {
		err = planDao.NewPlanTaskConditionDataDao().UpdateData(ctx, enterpriseId, one.PlanId, map[string]interface{}{
			"task_allot_log_id":     one.ID,
			"work_external_user_id": failList,
		}, map[string]interface{}{
			"is_touch": microScrmModel.TouchStatusFail,
			"is_union": 0,
		})
	}
}
