package plan_tasks

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"micro-scrm/internal/app"
	microBossDao "micro-scrm/internal/app/dao/micro_boss"
	planDao "micro-scrm/internal/app/dao/micro_scrm"
	planReq "micro-scrm/internal/app/http/request/plan"
	"micro-scrm/internal/app/models"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	planService "micro-scrm/internal/app/service/plan"
	"micro-scrm/internal/modules/groupmsg"
	"micro-scrm/internal/modules/marketingPromotion/storage"
	"micro-scrm/internal/modules/marketingPromotion/storage/cos"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/boss"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sdk/macallback"
	"micro-scrm/internal/pkg/sdk/microwork"
	qyWork "micro-scrm/internal/pkg/sdk/qy_work"
	"micro-scrm/internal/pkg/sdk/work"
	"micro-scrm/internal/pkg/xjob"
	"os"
	"strconv"
	"time"

	"github.com/chrislusf/glow/flow"
	"github.com/samber/lo"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
)

type MaTaskData struct {
	Name string
}

func NewMaTaskData() MaTaskData {
	return MaTaskData{
		Name: "MA任务的推送(三方根据获取任务模板创建的任务)",
	}
}

type TaskInfo struct {
	TaskId   int64
	PlanName string
	TaskName string
}

func (job MaTaskData) Run(c ctx.Context, handler *xjob.JobHandler) error {

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
		enterpriseName := v.Name
		err = job.handle(c, enterpriseId, enterpriseName)
		if err != nil {
			return err
		}
	}
	return nil
}

// 为解决任务同时推送带来的导购并发操作，一次只推一条任务
func (job MaTaskData) handle(c ctx.Context, enterpriseId int64, enterpriseName string) error {
	taskInfo, err := planDao.NewPlanTaskDao().GetMaPlanPushInfo(c, enterpriseId)
	if err != nil {
		return err
	}
	if taskInfo.ID == 0 {
		return nil
	}
	// 更新任务人群包数据为已完成,任务为推送中
	err = planDao.NewPlanTaskDao().UpdateData(c, enterpriseId, map[string]interface{}{
		"id": taskInfo.ID,
	}, map[string]interface{}{"task_is_push": microScrmModel.PlanTaskPushRunning})
	if err != nil {
		return err
	}
	return job.operationPlanTask(c, enterpriseId, taskInfo.PlanID, taskInfo.ID, enterpriseName)
}

//func (job MaTaskData) handle(c ctx.Context, enterpriseId int64, enterpriseName string) error {
//	lists, err := planDao.NewPlanTaskDao().GetMaPlanPushLists(c, enterpriseId)
//	if err != nil {
//		return err
//	}
//	if len(lists) == 0 {
//		return nil
//	}
//	var pushTaskIds []int64
//	for _, value := range lists {
//		info := value
//		if info.PlanSendTime <= time.Now().Unix() && info.TaskStartTime <= time.Now().Unix() && info.TaskEndTime >= time.Now().Unix() {
//			pushTaskIds = append(pushTaskIds, info.ID)
//		}
//	}
//	if len(pushTaskIds) == 0 {
//		return nil
//	}
//	// 更新任务人群包数据为已完成,任务为推送中
//	err = planDao.NewPlanTaskDao().UpdateData(c, enterpriseId, map[string]interface{}{
//		"id": pushTaskIds,
//	}, map[string]interface{}{"task_is_push": 2})
//	if err != nil {
//		return err
//	}
//	var wg sync.WaitGroup
//	for _, v := range lists {
//		item := v
//		if !helper.InArray(item.ID, pushTaskIds) {
//			continue
//		}
//		wg.Add(1)
//		go func(planId, taskId int64) {
//			defer wg.Done()
//			err = job.operationPlanTask(c, enterpriseId, planId, taskId, enterpriseName)
//			if err != nil {
//				return
//			}
//		}(item.PlanID, item.ID)
//	}
//	wg.Wait()
//	return nil
//}

func (job MaTaskData) operationPlanTask(c ctx.Context, enterpriseId, planId, taskId int64, enterpriseName string) error {

	condition, err := planDao.NewPlanTaskDao().GetConditionInfo(c, enterpriseId, map[string]interface{}{
		"task_id": taskId,
	})
	if err != nil {
		return err
	}
	conditionID := condition.ID
	var userIds []string
	nowTime := time.Now().Unix()
	taskDetail, err := planDao.NewPlanTaskDao().GetOne(c, enterpriseId, map[string]interface{}{"id": taskId})
	if err != nil {
		return err
	}
	qySend := false
	if taskDetail.GroupSendType == microScrmModel.PLANTASKQYSEND {
		qySend = true
	}
	planDetail, err := planDao.NewPlanTaskDao().GetPlanDetail(c, planReq.PlanInfo{
		BaseReq: planReq.BaseReq{
			EnterpriseId: enterpriseId,
		},
		Id: taskDetail.PlanID,
	})
	if err != nil {
		return err
	}
	//page := 0
	//pageSize := 1000

	customerNumMap := make(map[string]int64, 0)
	uniqueUnionids := make(map[string]interface{}, 0)

	conditionUserList := make(map[string][]microScrmModel.TaskConditionDataModel, 0)
	conditionUserLogID := make(map[string]int64, 0)

	memberExternalMap := make(map[string][]string, 0) // 导购外部联系人信息
	crowdDataId := int64(0)
	for {
		//page++
		//datalist, errs := planDao.NewPlanTaskCrowdDataDao().GetCrowdDataList(c, enterpriseId, taskId, page, pageSize)
		datalist, errs := planDao.NewPlanTaskCrowdDataDao().GetCrowdDataById(c, enterpriseId, taskId, crowdDataId)
		if errs != nil {
			fmt.Printf("获取MA人群包匹配结果数据报错,错误%s\n", errs.Error())
			continue
		}
		if len(datalist) == 0 {
			break
		}

		for _, v := range datalist {
			crowdDataId = v.ID
			IsUnion := 1
			if v.Name == "" {
				IsUnion = 0
			}
			Condition := microScrmModel.TaskConditionDataModel{
				ConditionId:        conditionID,
				PlanId:             planId,
				TaskId:             taskId,
				TagName:            "",
				UserId:             v.SaCode,
				IsUnion:            IsUnion,
				UnionId:            v.Unionid,
				WorkExternalUserId: v.WorkExternalUserID,
				WwTagId:            "",
				CustomerID:         v.CustomerID,
				IsMatch:            v.IsMatch,
				UserName:           v.UserName,
			}

			key := v.SaCode + "-" + v.Unionid
			if v.WorkExternalUserID != "" {
				memberExternalMap[v.SaCode] = append(memberExternalMap[v.SaCode], v.WorkExternalUserID)
			}
			if _, ok := customerNumMap[v.SaCode]; !ok {
				customerNumMap[v.SaCode] = 1
				userIds = append(userIds, v.SaCode)
				TaskAllotLog := microScrmModel.TaskAllotLogModel{
					PlanId:      planId,
					ConditionId: conditionID,
					TaskId:      taskId,
					UserId:      v.SaCode,
					StartTime:   int(nowTime),
					CustomerNum: 0,
					BaseModel:   models.CreateBase(0),
				}
				logID, _ := planService.NewTaskAllotLog().CreateTaskAllotLogData(c, enterpriseId, TaskAllotLog)
				Condition.TaskAllotLogId = logID
				conditionUserList[v.SaCode] = append(conditionUserList[v.SaCode], Condition)
				conditionUserLogID[v.SaCode] = logID
				uniqueUnionids[key] = v.Unionid
			} else {
				if _, exist := uniqueUnionids[key]; !exist {
					uniqueUnionids[key] = v.Unionid
					customerNumMap[v.SaCode]++
				}
				Condition.TaskAllotLogId = conditionUserLogID[v.SaCode]
				conditionUserList[v.SaCode] = append(conditionUserList[v.SaCode], Condition)
			}
		}
	}
	for k, conditionModels := range conditionUserList {
		size := 1000
		length := len(conditionModels)
		InsertData := conditionModels
		if length > size {
			maxPageData := int(math.Ceil(float64(length / size)))
			maxPageData += 1
			for i := 1; i <= maxPageData; i++ {
				start := (i - 1) * size
				end := i * size
				InsertData = conditionModels[start:]
				if end < length {
					InsertData = conditionModels[start:end]
				}
				_ = planService.NewPlanTaskConditionData().InsertConditionData(c, enterpriseId, planId, InsertData)
			}
		} else {
			_ = planService.NewPlanTaskConditionData().InsertConditionData(c, enterpriseId, planId, InsertData)
		}
		if taskAllotLogId, ok := conditionUserLogID[k]; ok { // 更新task_allot_log表中customer_num的值
			customerNum := customerNumMap[k]
			_ = planDao.NewPlanTaskAllotLogDao().UpdateData(c, enterpriseId, map[string]interface{}{"id": taskAllotLogId}, map[string]interface{}{"customer_num": customerNum})
		}
	}

	taskInfo := TaskInfo{
		TaskId:   taskId,
		PlanName: taskDetail.Name,
		TaskName: planDetail.Name,
	}
	if len(conditionUserLogID) > 0 {
		pluginInfo, err := microBossDao.NewEnterpriseWorkPlugin().GetOne(c, map[string]interface{}{
			"enterprise_id": enterpriseId,
			"identifier":    work.QWCrsStoreAssistantSelf,
		})
		if err != nil {
			return err
		}
		info, err := common.GetEnterpriseInfo(c, enterpriseId)
		if err != nil {
			return err
		}
		if qySend {
			// 获取企微API
			conditions, err2 := planDao.NewPlanTaskDao().GetConditionInfo(c, enterpriseId, map[string]interface{}{
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
			if taskDetail.TaskCategory == microScrmModel.PLANTASKCONTACT { // 企业群发类型
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
				groupSendContent, err := qyWorkApi.SetMediaID(c, mediaContentStr, enterpriseId)
				if err != nil {
					fmt.Printf("群发消息获取素材库id报错%v\n", err)
					job.handleMaterialUploadFail(c, enterpriseId, taskDetail, "企业群发素材", mediaContentStr, err)
					return nil
				}
				job.qyGroupSend(c, enterpriseId, taskId, memberExternalMap, groupSendContent, info.Domain, nil, taskDetail.AllowSelect)
			}
			if taskDetail.TaskCategory == microScrmModel.PLANTASKMOMENT { // 朋友圈任务的群发类型
				log.Info(c).Msgf("发送MA朋友圈任务开始, task_id:%d", taskId)
				qyWorkApi := new(work.QyWorkImpl).SetDomain(info.Domain)
				qyNewWorkApi := qyWork.NewQyWork(enterpriseId)
				var mediaContent work.AddMomentTaskMediaContent
				mediaContent.Text = &work.Text{Content: conditionContent.Text.Content}
				if len(conditionContent.Image.ImgList) > 0 {
					for _, v := range conditionContent.Image.ImgList {
						item := v
						mediaId, err := qyNewWorkApi.GetUploadAttachmentMediaId(c, item.PicUrl, "image", 1)
						if err != nil {
							log.Info(c).Msgf("上传朋友圈附件图片资源报错,错误为%v", err)
							job.handleMaterialUploadFail(c, enterpriseId, taskDetail, "朋友圈图片附件", item.PicUrl, err)
							return nil
						}
						mediaContent.Attachments = append(mediaContent.Attachments, work.AddMomentTaskAttachments{MsgType: "image", Image: &work.MomentImage{MediaId: mediaId}})
					}
				}
				if conditionContent.Link.Title != "" {
					mediaId, err := qyNewWorkApi.GetUploadAttachmentMediaId(c, conditionContent.Link.PicUrl, "image", 1)
					if err != nil {
						log.Info(c).Msgf("上传朋友圈附件图片资源报错,错误为%v", err)
						job.handleMaterialUploadFail(c, enterpriseId, taskDetail, "朋友圈链接封面附件", conditionContent.Link.PicUrl, err)
						return nil
					}
					mediaContent.Attachments = append(mediaContent.Attachments, work.AddMomentTaskAttachments{MsgType: "link", Link: &work.MomentLink{
						Title:   conditionContent.Link.Title,
						Url:     conditionContent.Link.LinkUrl,
						MediaId: mediaId,
					}})
				}
				if conditionContent.Video.VideoUrl != "" {
					mediaId, err := qyNewWorkApi.GetUploadAttachmentMediaId(c, conditionContent.Video.VideoUrl, "video", 1)
					if err != nil {
						log.Info(c).Msgf("上传朋友圈附件视频资源报错,错误为%v", err)
						job.handleMaterialUploadFail(c, enterpriseId, taskDetail, "朋友圈视频附件", conditionContent.Video.VideoUrl, err)
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
				if taskDetail.QyTagId != "" {
					addMomentTaskReq.VisibleRange.ExternalContactList = &struct {
						TagList []string `json:"tag_list"`
					}{TagList: []string{taskDetail.QyTagId}}
				}

				response, err := qyWorkApi.AddMomentTask(c, addMomentTaskReq)
				if err != nil {
					log.Info(c).Msgf("发送朋友圈群发消息失败, err:%s", err.Error())
					return err
				}
				if response.ErrCode != 0 {
					log.Info(c).Msgf("发送朋友圈群发消息企微返回错误, response:%v", response)
				}
				log.Info(c).Err(err).Interface("params", addMomentTaskReq).Msg("qyGroupSend.AddMomentTask")
				updateMap := make(map[string]interface{}, 0)
				updateMap = map[string]interface{}{"moment_job_id": response.Jobid}
				err = planDao.NewPlanTaskDao().UpdateData(c, enterpriseId, map[string]interface{}{
					"id": taskId,
				}, updateMap)
			}
		} else {
			job.sendPlanTaskMiniMessage(c, enterpriseId, userIds, taskInfo, conditionUserLogID, customerNumMap, enterpriseName, pluginInfo.MiniAppID, info.Domain)
		}
	}
	// 处理MA人群包的匹配结果上传到COS
	uniqueMd5 := helper.Md5V(strconv.FormatInt(taskId, 10))
	fileName := uniqueMd5 + "导入名单的匹配结果.csv"
	create, err := os.Create(fileName)
	defer os.Remove(fileName)
	if err != nil {
		fmt.Printf("创建MA匹配结果文件失败,错误%s\n", err.Error())
		return err
	}
	_, err = create.WriteString("\xEF\xBB\xBF")
	if err != nil {
		fmt.Printf("处理csv的bom头报错,错误%s\n", err.Error())
		return err
	} // 写入UTF-8 BOM
	csvFile := csv.NewWriter(create)
	var header []string
	if taskDetail.Dimension == microScrmModel.TaskDimensionImportUnionId {
		header = []string{
			"unionid", "sacode", "匹配结果",
		}
	}
	if taskDetail.Dimension == microScrmModel.TaskDimensionImportSaUnionId {
		header = []string{
			"unionid", "customer_id", "sacode", "匹配结果",
		}
	}
	err = csvFile.Write(header)
	if err != nil {
		fmt.Printf("写入MA匹配结果文件失败,错误%s\n", err.Error())
		return err
	}
	// 获取数据写入csv文件
	//crowdLogPage := 0

	crowdLogId := int64(0)
	for {
		//crowdLogPage++
		//byPage, err := planDao.NewPlanTaskCrowdDataLogDao().GetCrowdDataResultByPage(c, enterpriseId, taskId, crowdLogPage, pageSize)
		byPage, err := planDao.NewPlanTaskCrowdDataLogDao().GetCrowdDataResultById(c, enterpriseId, taskId, crowdLogId)
		if err != nil {
			return err
		}
		if len(byPage) == 0 {
			break
		}
		fileSlice := make([][]string, 0)
		for _, model := range byPage {
			crowdLogId = model.ID
			if taskDetail.Dimension == microScrmModel.TaskDimensionImportSaUnionId {
				fileSlice = append(fileSlice, []string{
					model.Unionid, model.CustomerID + "\t", model.SaCode, model.ImportRes,
				})
			}
			if taskDetail.Dimension == microScrmModel.TaskDimensionImportUnionId {
				fileSlice = append(fileSlice, []string{
					model.Unionid, model.SaCode, model.ImportRes,
				})
			}
		}
		_ = csvFile.WriteAll(fileSlice)
	}

	csvFile.Flush()
	cs := storage.Init(cos.COS)
	cosPath, err := cs.Upload(c, fileName, "plan_task/"+fileName)
	err = create.Close()
	err = planDao.NewPlanTaskDao().UpdateData(c, enterpriseId, map[string]interface{}{
		"id": taskId,
	}, map[string]interface{}{"is_extract": 1, "download_url": cs.GetURL(cosPath)})
	err = planService.NewPlanTask().UpdatePlanTaskPushStatus(c, enterpriseId, taskId)
	var templateName string
	if taskDetail.PlanTemplateId > 0 {
		templateInfo, _ := planDao.NewPlanTemplateDao().GetTemplateInfo(c, enterpriseId, taskDetail.PlanTemplateId)
		templateName = templateInfo.TemplateName
	}

	crowdPushLogId := int64(0)
	for {
		var sendMaResultDataLists [][]microwork.SendMaResultData
		//lists, err2 := planDao.NewPlanTaskCrowdDataLogDao().GetCrowdDataResultByPage(c, enterpriseId, taskId, logPage, pageSize)
		lists, err2 := planDao.NewPlanTaskCrowdDataLogDao().GetCrowdDataResultById(c, enterpriseId, taskId, crowdPushLogId)
		if err2 != nil {
			fmt.Printf("获取MA的任务匹配结果报错,错误%s\n", err.Error())
			continue
		}
		if len(lists) == 0 {
			break
		}
		var crowdResult []microwork.SendMaResultData
		for _, value := range lists {
			info := value
			crowdPushLogId = info.ID
			if taskDetail.Dimension == microScrmModel.TaskDimensionImportUnionId {
				_ = csvFile.Write([]string{
					info.Unionid, info.SaCode, info.ImportRes,
				})
			}
			if taskDetail.Dimension == microScrmModel.TaskDimensionImportSaUnionId {
				_ = csvFile.Write([]string{
					info.Unionid, info.CustomerID + "\t", info.SaCode, info.ImportRes,
				})
			}
			sendMaResultData := microwork.SendMaResultData{
				Date:              time.Now().UTC().Format("2006-01-02T15:04:05Z"),
				SaName:            info.Name,
				SaCode:            info.SaCode,
				Unionid:           info.Unionid,
				Result:            "success",
				ResultMessage:     "success",
				FlowId:            taskDetail.FlowId,
				Campaign:          taskDetail.Campaign,
				CPlanTemplateId:   strconv.FormatInt(taskDetail.PlanTemplateId, 10),
				CPlanTemplateName: templateName,
				CTaskRemark:       taskDetail.Remark,
				CTaskClassifation: planDetail.Name,
				TaskRemark:        planDetail.Remark,
			}
			if info.ImportRes == "匹配失败" {
				sendMaResultData.Result = "fail"
				sendMaResultData.ResultMessage = "未匹配到好友关系"
			}
			crowdResult = append(crowdResult, sendMaResultData)
		}
		sendMaResultDataLists = append(sendMaResultDataLists, crowdResult)
		// 推送MA的回执数据
		job.sendMAPlanResultData(c, enterpriseId, taskDetail, sendMaResultDataLists)
	}
	return err

}

// sendPlanTaskMiniMessage 发送任务模板消息
func (job MaTaskData) sendPlanTaskMiniMessage(ctx ctx.Context, enterpriseId int64, sendUsers []string, task TaskInfo, userLogIdMap, customerNumMap map[string]int64, enterpriseName, appId, domain string) {
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

	for _, userid := range sendUsers {
		logId, ok := userLogIdMap[userid]
		if !ok {
			continue
		}
		num := int64(0)
		if _, found := customerNumMap[userid]; found {
			num = customerNumMap[userid]
		}
		params := work.SendMessageRequest{
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
		qyWorkApi := new(work.QyWorkImpl).SetDomain(domain)
		response, err := qyWorkApi.SendMessage(ctx, params)
		if err != nil {
			log.Info(ctx).Msgf("发送小程序应用消息失败, err:%s", err.Error())
			continue
		}
		if response.Errcode != 0 {
			log.Info(ctx).Msgf("发送小程序应用消息企微返回错误, response:%v", response)
			continue
		}
		log.Info(ctx).Err(err).Interface("params", params).Msg("qyGroupSend.qyGroupSend")
		msgId := response.Msgid
		fmt.Printf("taskId:%d----导购id:%s----任务发送小程序应用消息的--msgID:%s\n", task.TaskId, params.Touser, msgId)
		err = planDao.NewPlanTaskAllotLogDao().UpdateData(ctx, enterpriseId, map[string]interface{}{
			"id": logId,
		}, map[string]interface{}{"agent_msg_id": msgId, "agent_msg_time": time.Now().Unix()})
		log.Info(ctx).Err(err).Interface("params", params).Msg("sendPlanTaskMiniMessage.SendMiniMessageRequest")
	}

}

func (job MaTaskData) handleMaterialUploadFail(c ctx.Context, enterpriseId int64, taskDetail microScrmModel.PlanTasksModel, materialType, materialUrl string, uploadErr error) {
	if err := planService.NewPlanTask().UpdatePlanTaskPushFailed(c, enterpriseId, taskDetail.ID); err != nil {
		log.Info(c).Msgf("MA任务素材上传失败后标记子任务下发失败报错, task_id:%d, err:%s", taskDetail.ID, err.Error())
	}
	errMsg := ""
	if uploadErr != nil {
		errMsg = uploadErr.Error()
	}
	planService.NewTaskNotice().TaskMaterialUploadFailNotice(c, planService.MaterialUploadFailMessageData{
		EnterpriseId: enterpriseId,
		PlanId:       taskDetail.PlanID,
		TaskId:       taskDetail.ID,
		TaskTitle:    taskDetail.Name,
		TaskCategory: taskDetail.TaskCategory,
		Dimension:    int(taskDetail.Dimension),
		MaterialType: materialType,
		MaterialUrl:  materialUrl,
		ErrorMessage: errMsg,
	})
	log.Info(c).Msgf("MA任务素材上传失败, task_id:%d, plan_id:%d, material_type:%s, material_url:%s, err:%s", taskDetail.ID, taskDetail.PlanID, materialType, materialUrl, errMsg)
}

// sendMAPlanResultData 任务发送将匹配结果回执给到MA任务
func (job MaTaskData) sendMAPlanResultData(c ctx.Context, enterpriseId int64, taskDetail microScrmModel.PlanTasksModel, lists [][]microwork.SendMaResultData) {
	flow.New().Source(func(ch chan microwork.SendMaResultReq) {
		for _, item := range lists {
			list := item
			ch <- microwork.SendMaResultReq{
				EnterpriseID: enterpriseId,
				TaskId:       taskDetail.ID,
				CrowdId:      taskDetail.CrowdID,
				TaskCategory: taskDetail.TaskCategory,
				CrowdResult:  list,
			}
		}
	}, 8).Map(func(params microwork.SendMaResultReq) {
		err := macallback.SendResult(c, taskDetail.OriginType, params)
		log.Info(c).Err(err).Interface("params", params).Msg("请求MA的任务发送匹配结果回执")
	}).Run()
}

// 企业群发
func (job MaTaskData) qyGroupSend(ctx ctx.Context, enterpriseId, taskId int64, memberExternalMap map[string][]string, groupSendContent string, domain string, userIdPathMap map[string]string, allowSelect int) {
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
		one, err := planDao.NewPlanTaskAllotLogDao().GetOne(ctx, enterpriseId, map[string]interface{}{
			"task_id": taskId,
			"user_id": params.Sender,
		})
		fmt.Printf("taskId:%d----导购id:%s----任务企业群发的--msgID:%s\n", taskId, params.Sender, msgId)
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
