package plan_tasks

import (
	"encoding/json"
	"fmt"
	"micro-scrm/internal/app"
	microBossDao "micro-scrm/internal/app/dao/micro_boss"
	planDao "micro-scrm/internal/app/dao/micro_scrm"
	planReq "micro-scrm/internal/app/http/request/plan"
	"micro-scrm/internal/app/models"
	"micro-scrm/internal/app/models/micro_scrm"
	planService "micro-scrm/internal/app/service/plan"
	"micro-scrm/internal/modules/groupmsg"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/boss"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/helper"
	qyWork "micro-scrm/internal/pkg/sdk/qy_work"
	"micro-scrm/internal/pkg/sdk/work"
	"sync"
	"time"

	"github.com/samber/lo"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"

	"micro-scrm/internal/pkg/xjob"

	"github.com/chrislusf/glow/flow"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
)

type templateTaskData struct {
	name string
}

func NewTemplateTaskData() *templateTaskData {
	return &templateTaskData{
		name: "客户导入模板的任务计划",
	}
}

// ChunkSendQyMsg 分块发送企业微信消息,该结构体，在ma推送任务也用到了
type ChunkSendQyMsg struct {
	IsChunk bool `json:"is_chunk"`
	work.SendMsgTemplateReq
}

type templateTask struct {
	taskID         int64
	fileID         int
	IsSpecifyCrowd int
	task           micro_scrm_enterprise.ResponseTemplateCondition
}

func (job templateTaskData) Run(ctx ctx.Context, handler *xjob.JobHandler) error {
	var (
		err            error
		enterpriseList []boss.EnterpriseModel
	)
	n := time.Now()
	lockEight := time.Date(n.Year(), n.Month(), n.Day(), 8, 0, 0, 0, time.Local)
	if n.Before(lockEight) { // 任务推送在每天的8点之后开始推送
		return nil
	}
	enterpriseList, err = orm.BossService.Enterprise.GetAllInfo()
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, enterprise := range enterpriseList {
		// job.getTask(ctx, enterprise.ID, enterprise.Name)
		wg.Add(1)
		go func(enterpriseID int64, enterpriseName string) {
			defer wg.Done()
			job.getTask(ctx, enterpriseID, enterpriseName)
		}(enterprise.ID, enterprise.Name)
	}
	wg.Wait()
	return nil
}

func (job templateTaskData) getTask(ctx ctx.Context, enterpriseID int64, enterpriseName string) error {
	planTasks, err := orm.Service.PlanTask.GetTemplatePushData(enterpriseID, micro_scrm_enterprise.PlanModel{})
	if err != nil {
		return err
	}

	var taskTemplate = make(map[int64]templateTask, 0)
	nowTime := time.Now().Unix()
	for _, task := range planTasks {
		if task.TaskEndTime >= nowTime && task.TaskStartTime <= nowTime && task.IsExtract == 1 && task.TaskIsPush == 0 {
			template := templateTask{taskID: task.TaskId, fileID: task.ConsumerFileID, IsSpecifyCrowd: task.IsSpecifyCrowd, task: task}
			taskTemplate[task.TaskId] = template
		}
	}
	var wg sync.WaitGroup
	for _, template := range taskTemplate {
		// job.operationPlanTask(ctx, enterpriseID, template, enterpriseName)
		wg.Add(1)
		go func(template templateTask) {
			defer wg.Done()
			job.operationPlanTask(ctx, enterpriseID, template, enterpriseName)
		}(template)
	}
	wg.Wait()

	return nil
}

func (job templateTaskData) operationPlanTask(ctx ctx.Context, enterpriseID int64, templateTask templateTask, enterpriseName string) error {
	// 开始处理推送数据
	nowTime := time.Now().Unix()
	var userIds []string
	var userIdPathMap = make(map[string]string)

	conditionID := templateTask.task.ConditionId
	taskId := templateTask.taskID
	// if taskId != 3458149634980544 {
	// 	return nil
	// }
	taskInfo, err := planDao.NewPlanTaskDao().GetOne(ctx, enterpriseID, map[string]interface{}{
		"id": taskId,
	})
	if err != nil {
		return err
	}
	qySend := false
	if taskInfo.GroupSendType == microScrmModel.PLANTASKQYSEND {
		qySend = true
	}
	// 开始查询数据
	_ = planService.NewPlanTask().UpdatePlanTaskPush(ctx, enterpriseID, templateTask.taskID)
	customerNumMap := make(map[string]int64, 0) // 导购的客户数
	uniqueUnionids := make(map[string]interface{}, 0)
	memberExternalMap := make(map[string][]string, 0)                                // 导购外部联系人信息
	conditionUserList := make(map[string][]microScrmModel.TaskConditionDataModel, 0) // 导购的task_condition_data表的数据
	conditionUserLogID := make(map[string]int64, 0)                                  // 导购task_allot_log的id
	tagName := ""
	if taskInfo.IsAddTag == 1 {
		tagName = taskInfo.TagName
	}
	page := 0
	pageSize := 1000
	for {
		page++
		datalist, errs := planService.NewPlanTaskImportData().GetImportDataList(ctx, enterpriseID, int64(templateTask.fileID), page, pageSize)
		if errs != nil {
			fmt.Printf("获取%d人群包文件内容报错，错误%v\n", templateTask.fileID, err)
			continue
		}
		if len(datalist) == 0 {
			break
		}
		for _, v := range datalist {
			IsUnion := 1
			if v.Name == "" {
				IsUnion = 0
			}
			condition := microScrmModel.TaskConditionDataModel{
				ConditionId:        conditionID,
				PlanId:             templateTask.task.PlanId,
				TaskId:             templateTask.task.TaskId,
				TagName:            tagName,
				UserId:             v.Said,
				IsUnion:            IsUnion,
				UnionId:            v.UnionId,
				WorkExternalUserId: v.WorkExternalUserId,
				WwTagId:            taskInfo.QyTagId,
				CustomerID:         v.CustomerID,
				UserName:           v.UserName,
			}
			key := v.Said + "-" + v.UnionId + "-" + v.CustomerID
			if v.UnionId != "" && v.CustomerID != "" {
				key = v.Said + "-" + v.UnionId
			}
			if v.WorkExternalUserId != "" {
				memberExternalMap[v.Said] = append(memberExternalMap[v.Said], v.WorkExternalUserId)
			}
			if _, ok := customerNumMap[v.Said]; !ok {
				customerNumMap[v.Said] = 1
				userIds = append(userIds, v.Said)
				TaskAllotLog := microScrmModel.TaskAllotLogModel{
					PlanId:      templateTask.task.PlanId,
					ConditionId: conditionID,
					TaskId:      templateTask.taskID,
					UserId:      v.Said,
					StartTime:   int(nowTime),
					CustomerNum: 0,
					BaseModel:   models.CreateBase(0),
				}
				logID, _ := planService.NewTaskAllotLog().CreateTaskAllotLogData(ctx, enterpriseID, TaskAllotLog)
				condition.TaskAllotLogId = logID
				conditionUserLogID[v.Said] = logID
				uniqueUnionids[key] = v.UnionId + v.CustomerID
				conditionUserList[v.Said] = append(conditionUserList[v.Said], condition)
			} else {
				if _, exist := uniqueUnionids[key]; !exist {
					uniqueUnionids[key] = v.UnionId + v.CustomerID
					customerNumMap[v.Said]++
					condition.TaskAllotLogId = conditionUserLogID[v.Said]
					conditionUserList[v.Said] = append(conditionUserList[v.Said], condition)
				}
			}
		}
	}
	for k, conditionModels := range conditionUserList {
		insertData := conditionModels
		if len(insertData) > 1000 {
			chunkData := lo.Chunk(insertData, 1000)
			for _, data := range chunkData {
				_ = planService.NewPlanTaskConditionData().InsertConditionData(ctx, enterpriseID, taskInfo.PlanID, data)
			}
		} else {
			_ = planService.NewPlanTaskConditionData().InsertConditionData(ctx, enterpriseID, taskInfo.PlanID, insertData)
		}
		if taskAllotLogId, ok := conditionUserLogID[k]; ok { // 更新task_allot_log表中customer_num的值
			customerNum := customerNumMap[k]
			_ = planDao.NewPlanTaskAllotLogDao().UpdateData(ctx, enterpriseID, map[string]interface{}{"id": taskAllotLogId}, map[string]interface{}{"customer_num": customerNum})
		}
	}

	if len(conditionUserLogID) > 0 {

		pluginInfo, err := microBossDao.NewEnterpriseWorkPlugin().GetOne(ctx, map[string]interface{}{
			"enterprise_id": enterpriseID,
			"identifier":    work.QWCrsStoreAssistantSelf,
		})
		if err != nil {
			return err
		}
		info, err := common.GetEnterpriseInfo(ctx, enterpriseID)
		if err != nil {
			return err
		}
		// 企业群发类型的任务
		if qySend {
			// 获取企微API
			conditions, err2 := planDao.NewPlanTaskDao().GetConditionInfo(ctx, enterpriseID, map[string]interface{}{
				"task_id": taskId,
			})
			if err2 != nil {
				return err2
			}
			var conditionContent planReq.ConditionContent
			err := json.Unmarshal([]byte(conditions.Content), &conditionContent)
			if err != nil {
				return err
			}
			if taskInfo.TaskCategory == microScrmModel.PLANTASKCONTACT { // 企业群发类型
				qyWorkApi := new(work.QyWorkImpl).SetDomain(info.Domain)
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
				}
				mediaContentStr := helper.MarshalToString(mediaContent)
				groupSendContent, err := qyWorkApi.SetMediaID(ctx, mediaContentStr, enterpriseID)
				if err != nil {
					fmt.Printf("群发消息获取素材库id报错%v\n", err)
					job.handleMaterialUploadFail(ctx, enterpriseID, taskInfo, "企业群发素材", mediaContentStr, err)
					return nil
				}
				job.qyGroupSend(ctx, enterpriseID, taskId, memberExternalMap, groupSendContent, info.Domain, userIdPathMap, taskInfo.AllowSelect)
			}
			if taskInfo.TaskCategory == microScrmModel.PLANTASKMOMENT { // 朋友圈任务的群发类型
				qyWorkApi := new(work.QyWorkImpl).SetDomain(info.Domain)
				qyNewWorkApi := qyWork.NewQyWork(enterpriseID)
				var mediaContent work.AddMomentTaskMediaContent
				mediaContent.Text = &work.Text{Content: conditionContent.Text.Content}
				if len(conditionContent.Image.ImgList) > 0 {
					for _, v := range conditionContent.Image.ImgList {
						item := v
						mediaId, err := qyNewWorkApi.GetUploadAttachmentMediaId(ctx, item.PicUrl, "image", 1)
						if err != nil {
							log.Info(ctx).Msgf("上传朋友圈附件图片资源报错,错误为%v", err)
							job.handleMaterialUploadFail(ctx, enterpriseID, taskInfo, "朋友圈图片附件", item.PicUrl, err)
							return nil
						}
						mediaContent.Attachments = append(mediaContent.Attachments, work.AddMomentTaskAttachments{MsgType: "image", Image: &work.MomentImage{MediaId: mediaId}})
					}
				}
				if conditionContent.Link.Title != "" {
					mediaId, err := qyNewWorkApi.GetUploadAttachmentMediaId(ctx, conditionContent.Link.PicUrl, "image", 1)
					if err != nil {
						log.Info(ctx).Msgf("上传朋友圈附件图片资源报错,错误为%v", err)
						job.handleMaterialUploadFail(ctx, enterpriseID, taskInfo, "朋友圈链接封面附件", conditionContent.Link.PicUrl, err)
						return nil
					}
					mediaContent.Attachments = append(mediaContent.Attachments, work.AddMomentTaskAttachments{MsgType: "link", Link: &work.MomentLink{
						Title:   conditionContent.Link.Title,
						Url:     conditionContent.Link.LinkUrl,
						MediaId: mediaId,
					}})
				}
				if conditionContent.Video.VideoUrl != "" {
					mediaId, err := qyNewWorkApi.GetUploadAttachmentMediaId(ctx, conditionContent.Video.VideoUrl, "video", 1)
					if err != nil {
						log.Info(ctx).Msgf("上传朋友圈附件视频资源报错,错误为%v", err)
						job.handleMaterialUploadFail(ctx, enterpriseID, taskInfo, "朋友圈视频附件", conditionContent.Video.VideoUrl, err)
						return nil
					}
					mediaContent.Attachments = append(mediaContent.Attachments, work.AddMomentTaskAttachments{MsgType: "video", Video: &work.MomentVideo{
						MediaId: mediaId,
					}})
				}
				var addMomentTaskReq work.AddMomentTaskReq
				addMomentTaskReq.AddMomentTaskMediaContent = mediaContent
				addMomentTaskReq.VisibleRange = work.AddMomentVisibleRange{
					SenderList: &struct {
						UserList       []string `json:"user_list"`
						DepartmentList []int    `json:"department_list,omitempty"`
					}{UserList: userIds}}
				if taskInfo.QyTagId != "" {
					addMomentTaskReq.VisibleRange.ExternalContactList = &struct {
						TagList []string `json:"tag_list"`
					}{TagList: []string{taskInfo.QyTagId}}
				}

				response, err := qyWorkApi.AddMomentTask(ctx, addMomentTaskReq)
				if err != nil {
					log.Info(ctx).Msgf("发送朋友圈群发消息失败, err:%s", err.Error())
					return err
				}
				if response.ErrCode != 0 {
					log.Info(ctx).Msgf("发送朋友圈群发消息企微返回错误, response:%v", response)
				}
				log.Info(ctx).Err(err).Interface("params", addMomentTaskReq).Msg("qyGroupSend.AddMomentTask")
				updateMap := make(map[string]interface{}, 0)
				updateMap = map[string]interface{}{"moment_job_id": response.Jobid}
				err = planDao.NewPlanTaskDao().UpdateData(ctx, enterpriseID, map[string]interface{}{
					"id": taskId,
				}, updateMap)

			}
		} else {
			job.sendPlanTaskMiniMessage(ctx, enterpriseID, userIds, templateTask.task, conditionUserLogID, customerNumMap, enterpriseName, pluginInfo.MiniAppID, info.Domain)

		}
	}
	_ = planService.NewPlanTask().UpdatePlanTaskPushStatus(ctx, enterpriseID, templateTask.taskID)
	// 任务发送结果监控消息通知(为什么两个userIds集合，是非人群包任务那边变量有点混乱，目前还不能调整，后期再优化，所以就先这样兼容吧)
	planService.NewTaskNotice().TaskResult(ctx, taskInfo, enterpriseID, userIds, userIds)

	return nil
}

func (job templateTaskData) handleMaterialUploadFail(ctx ctx.Context, enterpriseId int64, taskInfo microScrmModel.PlanTasksModel, materialType, materialUrl string, uploadErr error) {
	if err := planService.NewPlanTask().UpdatePlanTaskPushFailed(ctx, enterpriseId, taskInfo.ID); err != nil {
		log.Info(ctx).Msgf("ConsumerPush任务素材上传失败后标记子任务下发失败报错, task_id:%d, err:%s", taskInfo.ID, err.Error())
	}
	errMsg := ""
	if uploadErr != nil {
		errMsg = uploadErr.Error()
	}
	planService.NewTaskNotice().TaskMaterialUploadFailNotice(ctx, planService.MaterialUploadFailMessageData{
		EnterpriseId: enterpriseId,
		PlanId:       taskInfo.PlanID,
		TaskId:       taskInfo.ID,
		TaskTitle:    taskInfo.Name,
		TaskCategory: taskInfo.TaskCategory,
		Dimension:    int(taskInfo.Dimension),
		MaterialType: materialType,
		MaterialUrl:  materialUrl,
		ErrorMessage: errMsg,
	})
	log.Info(ctx).Msgf("ConsumerPush任务素材上传失败, task_id:%d, plan_id:%d, material_type:%s, material_url:%s, err:%s", taskInfo.ID, taskInfo.PlanID, materialType, materialUrl, errMsg)
}

// 计划任务发送小程序模板消息通知
func (job templateTaskData) sendPlanTaskMiniMessage(ctx ctx.Context, enterpriseId int64, sendUsers []string, task micro_scrm_enterprise.ResponseTemplateCondition, userLogIdMap, customerNumMap map[string]int64, enterpriseName, appId string, domain string) {
	pagePath := "/pagesA/task-detail/task-detail"
	if enterpriseName == "coach" {
		pagePath = "packages/task/pages/task-detail/task-detail"
	}
	var content []work.MiniprogramNoticeItem
	content = append(content, work.MiniprogramNoticeItem{
		Key:   "任务计划名称",
		Value: task.PlanName,
	})
	content = append(content, work.MiniprogramNoticeItem{
		Key:   "子任务名称",
		Value: task.TaskName,
	})
	content = append(content, work.MiniprogramNoticeItem{
		Key:   "说明",
		Value: "你的任务计划已开启，点击查看你的任务吧",
	})
	flow.New().Source(func(ch chan work.SendMessageRequest) {
		for _, userid := range sendUsers {
			if _, found := userLogIdMap[userid]; !found {
				continue
			}
			num := int64(0)
			if _, found := customerNumMap[userid]; found {
				num = customerNumMap[userid]
			}
			ch <- work.SendMessageRequest{
				QyGwRequest: work.QyGwRequest{ReqPL: work.QWCrsStoreAssistantSelf},
				Touser:      userid,
				Msgtype:     "miniprogram_notice",
				MiniprogramNotice: &work.MiniprogramNoticeContent{
					Appid:       appId,
					Page:        fmt.Sprintf("%s?id=%d&name=%s&num=%d", pagePath, userLogIdMap[userid], task.TaskName, num),
					Title:       enterpriseName + " VIP活动提醒", // 导入客户模板任务计划提醒
					ContentItem: content,
				},
			}
		}
	}, 1).Map(func(params work.SendMessageRequest) {
		qyWorkApi := new(work.QyWorkImpl).SetDomain(domain)
		response, err := qyWorkApi.SendMessage(ctx, params)
		if err != nil {
			log.Info(ctx).Msgf("发送小程序应用消息失败, err:%s", err.Error())
			return
		}
		if response.Errcode != 0 {
			log.Info(ctx).Msgf("发送小程序应用消息企微返回错误, response:%v", response)
		}
		log.Info(ctx).Err(err).Interface("params", params).Msg("qyGroupSend.qyGroupSend")
		msgId := response.Msgid
		// fmt.Printf("taskId:%d----导购id:%s----任务发送小程序应用消息的--msgID:%s\n", task.TaskId, params.Touser, msgId)
		err = planDao.NewPlanTaskAllotLogDao().UpdateData(ctx, enterpriseId, map[string]interface{}{
			"task_id": task.TaskId,
			"user_id": params.Touser,
		}, map[string]interface{}{"agent_msg_id": msgId, "agent_msg_time": time.Now().Unix()})
		log.Info(ctx).Err(err).Interface("params", params).Msg("sendPlanTaskMiniMessage.SendMiniMessageRequest")
	}).Run()
}

// 企业群发
func (job templateTaskData) qyGroupSend(ctx ctx.Context, enterpriseId, taskId int64, memberExternalMap map[string][]string, groupSendContent string, domain string, userIdPathMap map[string]string, allowSelect int) {
	flow.New().Source(func(ch chan ChunkSendQyMsg) {
		for k, v := range memberExternalMap {
			userid := k
			externalUserIds := v
			chunkExternalUserIds := lo.Chunk(externalUserIds, 10000)
			sendQyMsg := ChunkSendQyMsg{
				IsChunk: false,
			}
			if len(chunkExternalUserIds) > 1 {
				sendQyMsg.IsChunk = true
			}
			for _, item := range chunkExternalUserIds {
				var massHairReq = work.SendMsgTemplateReq{}
				err := json.Unmarshal([]byte(groupSendContent), &massHairReq)
				if err != nil {
					continue
				}
				massHairReq.Sender = userid
				massHairReq.ChatType = groupmsg.ChatTypeUser
				massHairReq.ExternalUserid = item
				massHairReq.AllowSelect = false
				if allowSelect > 0 {
					massHairReq.AllowSelect = true
				}
				for k, attachment := range massHairReq.Attachments {
					if attachment.MiniProgram != nil {
						if attachment.MiniProgram.Page != "" {
							if path, ok := userIdPathMap[userid]; ok {
								massHairReq.Attachments[k].MiniProgram.Page = path
							}
						}
					}
				}
				sendQyMsg.SendMsgTemplateReq = massHairReq
				ch <- sendQyMsg
			}
		}
	}, 1).Map(func(params ChunkSendQyMsg) {
		qyWorkApi := new(work.QyWorkImpl).SetDomain(domain)
		response, err := qyWorkApi.SendMsgTemplate(ctx, params.SendMsgTemplateReq)
		if err != nil {
			log.Info(ctx).Msgf("发送群发消息失败, err:%s", err.Error())
			return
		}
		if response.ErrCode != 0 {
			log.Info(ctx).Msgf("发送群发消息企微返回错误, response:%v", response)
		}
		log.Info(ctx).Err(err).Interface("params", params).Msg("qyGroupSend.qyGroupSend")
		msgId := response.MsgID
		fmt.Printf("taskId:%d----导购id:%s----任务企业群发的--msgID:%s\n", taskId, params.Sender, msgId)
		one, err := planDao.NewPlanTaskAllotLogDao().GetOne(ctx, enterpriseId, map[string]interface{}{
			"task_id": taskId,
			"user_id": params.Sender,
		})
		if params.IsChunk {
			msgData := microScrmModel.TaskAllotLogMsgModel{
				BaseModel: models.BaseModel{
					ID: app.UID(),
				},
				TaskAllotLogId: one.ID,
				MsgId:          msgId,
				IsSend:         0,
			}
			err = planDao.NewPlanTaskAllotLogDao().CreateTaskSendMsg(ctx, enterpriseId, msgData)
			err = planDao.NewPlanTaskAllotLogDao().UpdateData(ctx, enterpriseId, map[string]interface{}{
				"id": one.ID,
			}, map[string]interface{}{"group_msg_id": "-1"})
		} else {
			err = planDao.NewPlanTaskAllotLogDao().UpdateData(ctx, enterpriseId, map[string]interface{}{
				"id": one.ID,
			}, map[string]interface{}{"group_msg_id": msgId})
		}
		fmt.Printf("taskId:%d----导购id:%s----更新发送状态失败:%v\n", taskId, params.Sender, err)
		failList := response.FailList
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

	}).Run()
}
