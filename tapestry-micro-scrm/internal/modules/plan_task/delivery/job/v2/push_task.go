package v2

import (
	"encoding/json"
	"fmt"
	"micro-scrm/internal/app"
	membersDao "micro-scrm/internal/app/dao/members"
	microBossDao "micro-scrm/internal/app/dao/micro_boss"
	planDao "micro-scrm/internal/app/dao/micro_scrm"
	planReq "micro-scrm/internal/app/http/request/plan"
	"micro-scrm/internal/app/models"
	membersModel "micro-scrm/internal/app/models/members"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	cacheSvc "micro-scrm/internal/app/service/cache"
	memberSvc "micro-scrm/internal/app/service/member"
	planService "micro-scrm/internal/app/service/plan"
	"micro-scrm/internal/modules/groupmsg"
	"micro-scrm/internal/orm/model/boss"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/helper"
	qyWork "micro-scrm/internal/pkg/sdk/qy_work"
	"micro-scrm/internal/pkg/sdk/work"
	"regexp"
	"strings"
	"time"

	"github.com/samber/lo"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"

	"github.com/chrislusf/glow/flow"

	"micro-scrm/internal/config"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/pkg/xjob"

	"github.com/emirpasic/gods/sets/hashset"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
)

var _ runner = taskPushData{}

type taskPushData struct {
	baseJob
	name string
}
type ConditionDataResult struct {
	TagId  int64  `json:"tag_id"`
	UserId string `json:"user_id"`
}

type UserTagsInfo struct {
	ID                 int64  `json:"id"`
	TagId              int64  `json:"tag_id"`
	WorkExternalUserId string `json:"work_external_user_id"`
	Unionid            string `json:"unionid"`
	MemberId           int64  `json:"member_id"`
	UserId             string `json:"user_id"`
	UserName           string `json:"user_name"`
}

type chunkSendQyMsg struct {
	IsChunk bool `json:"is_chunk"`
	work.SendMsgTemplateReq
}

func NewtaskPushData() *taskPushData {
	return &taskPushData{
		name: "推送任务计划",
	}
}

func (job taskPushData) Run(ctx ctx.Context, handler *xjob.JobHandler) error {

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

	for _, enterprise := range enterpriseList {
		_ = job.getTask(ctx, enterprise.ID)
	}

	return nil
}
func (job taskPushData) getTask(ctx ctx.Context, enterpriseID int64) error {
	planTasks, err := orm.Service.PlanTask.GetPushData(enterpriseID, micro_scrm_enterprise.PlanModel{})
	if err != nil {
		return err
	}
	var taskIds []int64
	var taskCondMap = make(map[int64]map[int64]string)
	var taskMap = make(map[int64]micro_scrm_enterprise.ResponCondition)
	for _, task := range planTasks {
		if task.TaskIsPush == 1 { // 推送过的就不在推送
			continue
		}
		if task.TaskEndTime >= time.Now().Unix() && task.PlanSendTime <= time.Now().Unix() && task.TaskStartTime <= time.Now().Unix() {
			if _, ok := taskCondMap[task.TaskId]; !ok {
				taskIds = append(taskIds, task.TaskId)
				taskMap[task.TaskId] = task
				taskCondMap[task.TaskId] = make(map[int64]string)
			}
			taskCondMap[task.TaskId][task.ConditionId] = task.Conditions
		}
	}
	var wg threading.RoutineGroup
	for _, task := range taskMap {
		newPlanTask := task
		wg.RunSafe(func() {
			err = job.OperationPlanTask(ctx, enterpriseID, newPlanTask, taskCondMap[newPlanTask.TaskId])
		})
	}
	wg.Wait()

	return nil
}
func (job taskPushData) OperationPlanTask(ctx ctx.Context, enterpriseID int64, task micro_scrm_enterprise.ResponCondition, conditions map[int64]string) error {
	_ = planService.NewPlanTask().UpdatePlanTaskPush(ctx, enterpriseID, task.TaskId)
	resConfig, _ := cacheSvc.NewRedisData().GetEnterpriseConfig(ctx, enterpriseID, "store_name")
	// 判断task是否指定人群 task.IsSpecifyCrowd=2表示指定人群,否则未指定人群
	var specifyCrowdUser []string
	var originTaskId int64        // 持续任务的扩展原始任务id
	if task.IsSpecifyCrowd == 2 { // 指定人群
		originTaskId = task.TaskId
		if task.TaskType == microScrmModel.PLANTCOUNT && task.OriginTaskId != 0 { // 如果是持续任务的指定导购人群，获取的是第一次持续任务的指定导购人群
			originTaskId = task.OriginTaskId
		}
		specifyCrowdUser, _ = orm.Service.PlanTask.PlanTasksIsSpecifyCrowdByTaskId(ctx, enterpriseID, originTaskId)
		if len(specifyCrowdUser) == 0 {
			log.Info(ctx).Msgf("企业信息为%v的用户标签任务%v:指定人群但未获取到", enterpriseID, task.TaskId)
			return nil
		}
	}
	taskInfo, err := planDao.NewPlanTaskDao().GetOne(ctx, enterpriseID, map[string]interface{}{
		"id": task.TaskId,
	})
	if err != nil {
		return err
	}
	qySend := false
	if taskInfo.GroupSendType == microScrmModel.PLANTASKQYSEND {
		qySend = true
	}
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
	// 获取全部的member信息
	memberMap := make(map[string]membersModel.MemberModel)
	memberList, err := membersDao.NewMembers().GetMany(ctx, map[string]interface{}{
		"enterprise_id": enterpriseID,
		"status":        1,
	})
	if err != nil {
		return err
	}
	for _, model := range memberList {
		memberMap[model.Userid] = model
	}
	var tagMap = make(map[int64]micro_scrm_enterprise.TagInfo)
	var qyTagId []string
	memberExternalMap := make(map[string][]string, 0)
	for condid, condsstr := range conditions {
		var tagNames []string
		var tags []int64
		if taskInfo.ConditionType == microScrmModel.TaskConditionTypeTag {
			condition := micro_scrm_enterprise.ConditionDetail{}
			err := json.Unmarshal([]byte(condsstr), &condition)
			if err != nil {
				return err
			}
			for _, tag := range condition.Tag {
				tags = append(tags, tag.TagID)
				qyTagId = append(qyTagId, tag.QwTagID)
				tagMap[tag.TagID] = tag
				tagNames = append(tagNames, tag.TagName)
			}
		}
		cursorId := int64(0)
		size := 5000
		userset := hashset.New()
		userLogIdMap := make(map[string]int64)
		customerNumMap := make(map[string]int64)
		var userIds []string
		var userIdPathMap = make(map[string]string)
		for {
			//lists, err := job.getUserListsByTags(ctx, enterpriseID, tags, specifyCrowdUser, task.SendType, cursorId, size, taskInfo.ConditionType)
			lists, err := job.getNewUserListsByTags(ctx, enterpriseID, tags, specifyCrowdUser, task.SendType, cursorId, size, taskInfo.ConditionType, memberMap)
			if err != nil {
				log.Error(ctx).Msgf("获取标签人员数据列表失败，报错为%v", err)
				return err
			}
			index := len(lists)
			if index == 0 {
				break
			}
			cursorId = lists[index-1].ID
			var conditionList []microScrmModel.TaskConditionDataModel
			for _, data := range lists {
				if data.Unionid == "" { // 如果unionid为空，在跳过
					continue
				}
				if data.UserId == "" {
					continue
				}
				tagName := ""
				wwTagId := ""
				logId, ok := userLogIdMap[data.UserId]
				if !ok {
					allotLog := micro_scrm_enterprise.TaskAllotLogModel{
						PlanId:      task.PlanId,
						ConditionId: condid,
						TaskId:      task.TaskId,
						UserId:      data.UserId,
						StartTime:   int(time.Now().Unix()),
						IsFinished:  0,
					}
					logId, err = orm.Service.PlanTask.CreateAllotLog(enterpriseID, allotLog)
					if err != nil {
						continue
					}
					userLogIdMap[data.UserId] = logId
				}
				if taskInfo.ConditionType == microScrmModel.TaskConditionTypeTag {
					tagName = tagMap[data.TagId].TagName
					wwTagId = tagMap[data.TagId].QwTagID
				}
				taskAllotLogId := logId
				newConditionInfo := microScrmModel.TaskConditionDataModel{
					TagName:            tagName,
					IsUnion:            1,
					UnionId:            data.Unionid,
					UserId:             data.UserId,
					WorkExternalUserId: data.WorkExternalUserId,
					WwTagId:            wwTagId,
					PlanId:             task.PlanId,
					TaskId:             task.TaskId,
					ConditionId:        condid,
					TaskAllotLogId:     taskAllotLogId,
					UserName:           data.UserName,
				}
				conditionList = append(conditionList, newConditionInfo)
				userset.Add(data.UserId)
				if data.WorkExternalUserId != "" {
					if !helper.InArray(data.WorkExternalUserId, memberExternalMap[data.UserId]) {
						memberExternalMap[data.UserId] = append(memberExternalMap[data.UserId], data.WorkExternalUserId)
						customerNumMap[data.UserId]++
						userIds = append(userIds, data.UserId)
					}
				}
			}
			if len(conditionList) > 0 {
				for _, val := range lo.Chunk(conditionList, 1000) {
					if err := planDao.NewPlanTaskConditionDataDao().InsertConditionData(ctx, enterpriseID, task.PlanId, val); err != nil {
						log.Info(ctx).Msgf("插入taskConditionData表失败，错误:%v", err)
						continue
					}
				}
			}
		}
		// 更新task_allot_log中customer_num值
		for k, v := range userLogIdMap {
			logId := v
			customerNum := customerNumMap[k]
			_ = planDao.NewPlanTaskAllotLogDao().UpdateData(ctx, enterpriseID, map[string]interface{}{"id": logId}, map[string]interface{}{"customer_num": customerNum})
		}
		userList := userset.Values()
		var sendUsers []string
		var sendUserIds []string
		users := lo.Chunk(userList, 1000)
		for _, user := range users {
			miniSendUserIds := make([]string, 0)
			for _, userId := range user {
				sendUsers = append(sendUsers, userId.(string))
				sendUserIds = append(sendUserIds, userId.(string))
				miniSendUserIds = append(miniSendUserIds, userId.(string))
			}
			if !qySend { // 企业群发的任务，不推送weClient小程序消息通知
				job.sendPlanTaskMiniMessage(ctx, enterpriseID, tagNames, miniSendUserIds, task, userLogIdMap, customerNumMap, resConfig.Name, pluginInfo.MiniAppID, info.Domain)
			}
		}
		// 企业群发类型
		if qySend {
			var conditionContent planReq.ConditionContent
			err := json.Unmarshal([]byte(task.Content), &conditionContent)
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

					if resConfig.Name == "sw" {
						userIds = helper.ArrayUniq(userIds)
						miniPathMap, err := memberSvc.NewMembersService().SetMiniProgramLinkParamsByTask(ctx, enterpriseID, userIds, conditionContent.Applets.Path)
						if err != nil {
							return err
						}
						for userId, path := range miniPathMap {
							userIdPathMap[userId] = path
						}
					}
				}
				mediaContentStr := helper.MarshalToString(mediaContent)
				groupSendContent, err := qyWorkApi.SetMediaID(ctx, mediaContentStr, enterpriseID)
				if err != nil {
					return err
				}
				job.qyGroupSend(ctx, enterpriseID, task.TaskId, memberExternalMap, groupSendContent, info.Domain, userIdPathMap, taskInfo.AllowSelect)
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
							continue
						}
						mediaContent.Attachments = append(mediaContent.Attachments, work.AddMomentTaskAttachments{MsgType: "image", Image: &work.MomentImage{MediaId: mediaId}})
					}
				}
				if conditionContent.Link.Title != "" {
					imgUrl := conditionContent.Link.PicUrl
					if imgUrl == "" {
						imgUrl = conditionContent.Link.CoverUrl
					}
					mediaId, err := qyNewWorkApi.GetUploadAttachmentMediaId(ctx, imgUrl, "image", 1)
					if err != nil {
						log.Info(ctx).Msgf("上传朋友圈附件图片资源报错,错误为%v", err)
						return err
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
						return err
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
					}{UserList: sendUserIds}}
				if len(qyTagId) > 0 {
					addMomentTaskReq.VisibleRange.ExternalContactList = &struct {
						TagList []string `json:"tag_list"`
					}{TagList: qyTagId}
				}

				response, err := qyWorkApi.AddMomentTask(ctx, addMomentTaskReq)
				if err != nil {
					log.Info(ctx).Msgf("发送朋友圈群发消息失败, err:%s", err.Error())
					return err
				}
				if response.ErrCode != 0 {
					fmt.Printf("发送朋友圈群发消息企微返回错误, response:%v\n", response)
					return err
				}
				time.Sleep(500 * time.Millisecond) // 延迟0.5秒，然后查询朋友圈的moment_id
				log.Info(ctx).Err(err).Interface("params", addMomentTaskReq).Msg("qyGroupSend.AddMomentTask")
				var massHairReq = work.MomentTaskResultRequest{}
				massHairReq.JobId = response.Jobid
				res, err := qyWorkApi.GetMomentTaskResult(ctx, massHairReq)
				if err != nil {
					app.Logger().Error(ctx).Msgf("获取朋友圈的发表结果失败, err:%s", err.Error())
					return err
				}
				if res.ErrCode != 0 {
					app.Logger().Error(ctx).Msgf("获取朋友圈的发表结果返回错误, res:%v", res)
					return err
				}
				momentId := res.Result.MomentId
				updateMap := make(map[string]interface{}, 0)
				updateMap = map[string]interface{}{"moment_status": res.Status, "moment_job_id": response.Jobid}
				if res.Status == 3 {
					updateMap = map[string]interface{}{"moment_status": res.Status, "moment_job_id": momentId}
				}
				err = planDao.NewPlanTaskDao().UpdateData(ctx, enterpriseID, map[string]interface{}{
					"id": task.TaskId,
				}, updateMap)
				if err != nil {
					fmt.Printf("更新朋友圈的任务momentId失败, 错误为:%v\n", err)
				}
			}
		}
		// 任务发送结果监控消息通知 userIds是原始数据筛出的导购id集合，sendUserIds是朋友圈任务发送时使用的导购id集合
		planService.NewTaskNotice().TaskResult(ctx, taskInfo, enterpriseID, userIds, sendUserIds)

	}
	// 如果是持续类型的子任务，当天的子任务完成之后，需要创建一条明天的子任务
	if task.TaskType == microScrmModel.PLANTCOUNT {
		// 根据时间戳获取时分秒
		createdAt := taskInfo.CreatedAt + 86400
		planSendTime := time.Now().Add(time.Hour * 24).Unix()
		taskStartTime := time.Unix(taskInfo.TaskStartTime, 0).Add(time.Hour * 24).Unix()
		if taskStartTime >= task.EndTime || taskStartTime >= taskInfo.TaskEndTime {
			return err
		}
		taskCopyId := app.UID()
		taskCopyInfo := microScrmModel.PlanTasksModel{
			BaseModel: models.BaseModel{
				ID:        taskCopyId,
				CreatedAt: createdAt,
				UpdatedAt: createdAt,
			},
			PlanID:         taskInfo.PlanID,
			Name:           taskInfo.Name,
			Dimension:      taskInfo.Dimension,
			ConditionType:  taskInfo.ConditionType,
			IsSpecifyCrowd: taskInfo.IsSpecifyCrowd,
			IsOpen:         1,
			TaskType:       taskInfo.TaskType,
			SendType:       taskInfo.SendType,
			Status:         taskInfo.Status,
			Remark:         taskInfo.Remark,
			PlanSendTime:   planSendTime,
			TaskCategory:   taskInfo.TaskCategory,
			GroupSendType:  taskInfo.GroupSendType,
			IsAddTag:       taskInfo.IsAddTag,
			TagExpireTime:  taskInfo.TagExpireTime,
			TagName:        taskInfo.TagName,
			QyTagId:        taskInfo.QyTagId,
			TagExpireType:  taskInfo.TagExpireType,
			OriginTaskId:   originTaskId,
			FinishRate:     taskInfo.FinishRate,
			TaskStartTime:  taskStartTime,
			TaskEndTime:    taskInfo.TaskEndTime,
			TaskExpireType: taskInfo.TaskExpireType,
		}
		conditionData := microScrmModel.PlanTaskConditionsModel{
			BaseModel: models.BaseModel{
				ID:        app.UID(),
				CreatedAt: createdAt,
				UpdatedAt: createdAt,
			},
			PlanID:      task.PlanId,
			TaskID:      taskCopyId,
			Conditions:  task.Conditions,
			Content:     task.Content,
			ContentType: int32(task.ContentType),
		}
		err = planDao.NewPlanTaskDao().CreateCopyTask(ctx, enterpriseID, taskCopyInfo, conditionData)
	}
	err = planService.NewPlanTask().UpdatePlanTaskPushStatus(ctx, enterpriseID, task.TaskId)
	if taskInfo.TaskCategory == microScrmModel.PLANTASKMOMENT {

		dbDSN := config.CfgData.DB.MasterDSN
		// 正则表达式匹配
		re := regexp.MustCompile(`^([^:]+):([^@]+)@(?:tcp\()?([^:]+):([0-9]+)\)?(?:/([^?]+))?(?:\?(.*))?$`)
		matches := re.FindStringSubmatch(dbDSN)

		var host, port, username, password string
		if matches != nil && len(matches) >= 5 {
			username = matches[1]
			password = matches[2]
			host = matches[3]
			port = matches[4]
			sql := fmt.Sprintf("INSERT INTO micro_scrm_enterprise_%d_ck.task_condition_data_local "+
				" SELECT id,task_allot_log_id,tag_name,ww_tag_id,union_id,work_external_user_id,is_union,mobile,mars_id,user_id,condition_id,task_id,plan_id,created_at,updated_at,deleted_at,customer_id,is_finished,finished_at,remark,is_contact,is_call,is_add,is_match,is_touch,user_name,group_chat_id,customer_status,forward FROM"+
				" mysql('%s:%s', 'micro_scrm_enterprise_%d', 'task_condition_data_%d', '%s', '%s') where task_id=%d", enterpriseID, host, port, enterpriseID, task.PlanId, username, password, task.TaskId)
			log.Info(ctx).Msgf("ck同步SQl:	%s", sql)
			err = gormx.GetCkInstance().Exec(sql).Error
		}

	}
	return err
}

// 小程序类型的应用消息提醒
func (job taskPushData) sendPlanTaskMiniMessage(ctx ctx.Context, enterpriseID int64, tagNames []string, sendUsers []string, task micro_scrm_enterprise.ResponCondition, userLogIdMap, customerNumMap map[string]int64, enterpriseName, appId, domain string) {
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
	if len(tagNames) > 0 {
		content = append(content, work.MiniprogramNoticeItem{
			Key:   "客户标签",
			Value: strings.Join(tagNames, ";"),
		})
	}
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
					Title:       "导购用户画像任务计划提醒",
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
		fmt.Printf("taskId:%d----导购id:%s----任务发送小程序应用消息的--msgID:%s\n", task.TaskId, params.Touser, msgId)
		err = planDao.NewPlanTaskAllotLogDao().UpdateData(ctx, enterpriseID, map[string]interface{}{
			"task_id": task.TaskId,
			"user_id": params.Touser,
		}, map[string]interface{}{"agent_msg_id": msgId, "agent_msg_time": time.Now().Unix()})
		log.Info(ctx).Err(err).Interface("params", params).Msg("sendPlanTaskMiniMessage.SendMiniMessageRequest")
	}).Run()
}

// 根据标签获取人群包数据
func (job taskPushData) getUserListsByTags(c ctx.Context, enterpriseId int64, tagIds []int64, userIds []string, sendType int, cursorId int64, size int, conditionType int32) ([]UserTagsInfo, error) {
	var result []UserTagsInfo
	searchAll := false
	if conditionType == microScrmModel.TaskConditionTypeAllCustomer {
		searchAll = true
	}
	userCorpExternalTable := fmt.Sprintf("micro_scrm_enterprise_%d_ck.wx_user_corp_external", enterpriseId)
	memberCk := "micro_scrm_ck.members"
	userRelationTable := "micro_scrm_ck.work_external_user_relation"
	whereMap := map[string]interface{}{
		"a.enterprise_id":   enterpriseId,
		"a.deleted_at":      0,
		"a.relation_status": 0,
		"m.status":          1,
		"m.deleted_at":      0,
	}
	if len(userIds) > 0 {
		whereMap["a.user_id"] = userIds
	}
	selectStr := "a.id as `id`,c.openid as work_external_user_id,a.unionid as `unionid`,a.member_id as `member_id`,a.user_id as `user_id`,c.name as `user_name`"
	query := gormx.GetCkInstance().WithContext(c).Table(userRelationTable + " as a").Joins(" LEFT JOIN " + userCorpExternalTable + " c ON a.work_external_user_id=c.id").
		Joins("LEFT JOIN " + memberCk + " m ON a.member_id=m.id")
	if !searchAll {
		if sendType == 2 {
			whereMap["a.is_key_sa"] = 1
		}
		whereMap["work_external_user_tag.external_tag_id"] = tagIds
		whereMap["work_external_user_tag.deleted_at"] = 0
		query = gormx.GetCkInstance().WithContext(c).Table("micro_scrm_ck.work_external_user_tag").
			Joins("LEFT JOIN " + userRelationTable + " as a ON work_external_user_tag.relation_id = a.id").
			Joins(" LEFT JOIN " + userCorpExternalTable + " c ON a.work_external_user_id=c.id").
			Joins("LEFT JOIN " + memberCk + " m ON a.member_id=m.id")
		selectStr = "a.id as `id`,work_external_user_tag.external_tag_id as tag_id,c.openid as work_external_user_id,a.unionid as `unionid`,a.member_id as `member_id`,a.user_id as `user_id`,c.name as `user_name`"
	}
	err := query.Where(whereMap).Where("a.id>?", cursorId).Select(selectStr).Order("a.id asc").Limit(size).Scan(&result).Error
	return result, err
}

// 企业群发
func (job taskPushData) qyGroupSend(ctx ctx.Context, enterpriseId, taskId int64, memberExternalMap map[string][]string, groupSendContent string, domain string, userIdPathMap map[string]string, allowSelect int) {
	flow.New().Source(func(ch chan chunkSendQyMsg) {
		for k, v := range memberExternalMap {
			userid := k
			externalUserIds := v
			chunkExternalUserIds := lo.Chunk(externalUserIds, 10000)
			sendQyMsg := chunkSendQyMsg{
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
	}, 1).Map(func(params chunkSendQyMsg) {
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

// 根据标签获取人群包数据
func (job taskPushData) getNewUserListsByTags(c ctx.Context, enterpriseId int64, tagIds []int64, userIds []string, sendType int, cursorId int64, size int, conditionType int32, memberMap map[string]membersModel.MemberModel) ([]UserTagsInfo, error) {
	var result []UserTagsInfo
	searchAll := false
	if conditionType == microScrmModel.TaskConditionTypeAllCustomer {
		searchAll = true
	}
	whereMap := map[string]interface{}{
		"enterprise_id":   enterpriseId,
		"deleted_at":      0,
		"relation_status": 0,
	}
	if len(userIds) > 0 {
		whereMap["user_id"] = userIds
	}
	var relationData []microScrmModel.WorkExternalUserRelationModel
	userRelationTable := "micro_scrm_ck.work_external_user_relation"
	tagTable := "micro_scrm_ck.work_external_user_tag"
	query := gormx.GetCkInstance().WithContext(c).Table(userRelationTable).Where(whereMap)
	if !searchAll {
		if sendType == 2 {
			whereMap["is_key_sa"] = 1
		}
		query = query.Where("id in (SELECT relation_id FROM micro_scrm_ck.work_external_user_tag WHERE enterprise_id=? and external_tag_id in (?) and deleted_at=0)", enterpriseId, tagIds)
	}

	err := query.Where("id>?", cursorId).Order("id asc").Limit(size).Scan(&relationData).Error
	relationIds := make([]int64, 0)
	unionIds := make([]string, 0)
	for _, datum := range relationData {
		relationIds = append(relationIds, datum.ID)
		unionIds = append(unionIds, datum.Unionid)
	}
	relationIds = helper.ArrayUniq(relationIds)
	unionIds = helper.ArrayUniq(unionIds)
	if len(unionIds) == 0 {
		return result, err
	}
	userCorpExternalTable := fmt.Sprintf("micro_scrm_enterprise_%d_ck.wx_user_corp_external", enterpriseId)
	var externalList []microScrmModel.WxUserCorpExternalModel
	err = gormx.GetCkInstance().WithContext(c).Table(userCorpExternalTable).Where(map[string]interface{}{
		"unionid": unionIds,
	}).Find(&externalList).Error

	externalMap := make(map[string]string, 0)

	for _, item := range externalList {
		externalMap[item.Unionid] = item.Openid
	}
	tagMap := make(map[int64]int64, 0)
	if !searchAll {
		var tagList []microScrmModel.WorkExternalUserTagModel
		err = gormx.GetCkInstance().WithContext(c).Table(tagTable).Where(map[string]interface{}{
			"enterprise_id":   enterpriseId,
			"external_tag_id": tagIds,
			"relation_id":     relationIds,
			"deleted_at":      0,
		}).Find(&result).Error
		if err != nil {
			return result, err
		}
		for _, item := range tagList {
			tagMap[item.RelationID] = item.ExternalTagID
		}
	}
	for _, datum := range relationData {
		result = append(result, UserTagsInfo{
			ID:                 datum.ID,
			Unionid:            datum.Unionid,
			WorkExternalUserId: externalMap[datum.Unionid],
			MemberId:           datum.MemberID,
			UserId:             datum.UserID,
			TagId:              tagMap[datum.ID],
			UserName:           memberMap[datum.UserID].Name,
		})
	}
	return result, err
}
