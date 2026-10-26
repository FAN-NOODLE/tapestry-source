package plan_task

import (
	"encoding/json"
	"errors"
	"fmt"
	"micro-scrm/internal/app"
	membersDao "micro-scrm/internal/app/dao/members"
	planDao "micro-scrm/internal/app/dao/micro_scrm"
	tagsDao "micro-scrm/internal/app/dao/tags"
	workDao "micro-scrm/internal/app/dao/work"
	planReq "micro-scrm/internal/app/http/request/plan"
	membersModel "micro-scrm/internal/app/models/members"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	workModels "micro-scrm/internal/app/models/work"
	marketingService "micro-scrm/internal/app/service/marketing"
	planService "micro-scrm/internal/app/service/plan"
	tagsSvc "micro-scrm/internal/app/service/tags"
	"strconv"
	"time"
	"unicode"

	"micro-scrm/internal/config"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm/group_chat"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/pkg/helper"
	micro_redis "micro-scrm/internal/pkg/redis"
	"micro-scrm/internal/pkg/uid"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/wx/wxoffiaccount/cache"
)

// PlanTasksAdd 子任务的添加
func PlanTasksAdd(enterpriseID int64, _ctx ctx.Context, params PlanTasksRequest) error {
	if params.ConditionType == microScrmModel.TaskConditionTypeTag && params.Dimension == microScrmModel.TaskDimensionTag {
		if params.SendType == 0 || params.TaskType == 0 {
			return errors.New("用户画像纬度的任务未添加推送对象或任务类型")
		}
	}
	taskStartTime, taskEndTime := params.TaskStartTime, params.TaskEndTime
	plan, err := orm.Service.PlanTask.PlanFound(params.PlanId, enterpriseID)
	if err != nil || plan.ID == 0 {
		return errors.New("计划不存在")
	}
	if plan.Status != 0 {
		return errors.New("计划已开始或完成,不能操作")
	}
	if len(params.PlanTaskConditions) == 0 {
		return errors.New("任务条件不存在")
	}
	if params.TaskExpireType == 0 {
		taskStartTime, taskEndTime = int64(plan.StartTime), int64(plan.EndTime)
	}
	if params.TaskExpireType == 1 {
		if int64(plan.StartTime) > params.TaskStartTime || int64(plan.EndTime) < params.TaskStartTime {
			return errors.New("子任务的时间范围应该在主任务的时间范围内")
		}
	}
	taskId := app.UID()
	workTagId := params.TagId
	endTime := 0
	addTag := false
	tagType := workModels.WorkExternalTagTypeQy
	if params.GroupSendType == microScrmModel.PLANTASKTAGSEND && params.IsAddTag == 1 {
		switch params.TagExpireType {
		case microScrmModel.TAGTASKFINISHEXPIRE:
			expireTime := params.TagExpireTime * 86400
			endTime = int(taskEndTime) + expireTime
		case microScrmModel.TAGTASKCUSTOMEXPIRE:
			endTime = params.TagExpireTime
		}
		if params.TagName == "" {
			return errors.New("标签名称不能为空")
		}
		addTag = true
	}
	if params.GroupSendType == microScrmModel.PLANTASKQYSEND {
		if params.TaskCategory == microScrmModel.PLANTASKMOMENT && helper.InArray(params.Dimension, []int{3, 4}) {
			timeTag := time.Now().Unix()
			params.TagName = "朋友圈临时标签" + strconv.Itoa(int(timeTag))
			tagType = workModels.MomentTaskTempTag
			params.TagExpireType = microScrmModel.TAG_TASK_SEND_EXPIRE
			addTag = true
		}
	}
	if addTag {
		wxtagId, tagId, err := tagsSvc.NewWorkTags().CreateTaskTag(_ctx, enterpriseID, "任务标签", params.TagName, endTime, tagType)
		if err != nil {
			return err
		}
		workTagId = wxtagId
		// todo 任务标签关联表插入 tagId
		tagRelationId := app.UID()
		tasksTagsRelation := workModels.PlanTasksTagRelationModel{
			ID:           tagRelationId,
			EnterpriseID: enterpriseID,
			PlanID:       taskId,
			TagID:        tagId,
			StartTime:    int(taskStartTime),
			EndTime:      int(taskEndTime),
		}
		relationId, err := workDao.NewPlanTasksTagRelation().CreateData(_ctx, tasksTagsRelation)
		if err != nil {
			return err
		}
		if relationId != tagRelationId {
			return errors.New("任务关联标签失败")
		}
	}

	planTask := micro_scrm_enterprise.PlanTaskModel{
		BranchModel: model.BranchModel{
			BaseModel: model.BaseModel{
				ID: taskId,
			},
		},
		PlanId:          params.PlanId,
		Name:            params.Name,
		Dimension:       params.Dimension,
		ConditionType:   params.ConditionType,
		IsOpen:          1,
		Status:          0,
		Remark:          params.Remark,
		TaskType:        params.TaskType,
		SendType:        params.SendType,
		BirthdayPreTime: params.BirthdayPreTime,
		ConsumerFileId:  params.ConsumerFileId,
		TaskCategory:    params.TaskCategory,
		GroupSendType:   params.GroupSendType,
		IsAddTag:        params.IsAddTag,
		TagExpireTime:   params.TagExpireTime,
		TagName:         params.TagName,
		QyTagId:         workTagId,
		TagExpireType:   params.TagExpireType,
		FinishRate:      params.FinishRate,
		TaskStartTime:   taskStartTime,
		TaskEndTime:     taskEndTime,
		TaskExpireType:  params.TaskExpireType,
		AllowSelect:     params.AllowSelect,
	}
	var planTaskSpecifyCrowds []micro_scrm_enterprise.PlanTaskSpecifyCrowdModel
	isSpecifyCrowd := 1
	if len(params.PlanTaskSpecifyCrowd) > 0 {
		specifyCrowdMap := make(map[any]bool, len(params.PlanTaskSpecifyCrowd))
		for _, value := range params.PlanTaskSpecifyCrowd {
			// 验证必填字段
			if value.UserId == "" && value.ThirdId == 0 {
				return errors.New("指定人群里用户信息有误")
			}

			// 确定去重的key和类型
			var key any
			var crowdType int
			if value.ThirdId > 0 {
				key = value.ThirdId
				crowdType = micro_scrm_enterprise.TaskCrowdTypeDept
			} else {
				key = value.UserId
				crowdType = micro_scrm_enterprise.TaskCrowdTypeSa
			}

			// 去重判断
			if specifyCrowdMap[key] {
				continue
			}
			specifyCrowdMap[key] = true

			// 构建数据模型
			planTaskSpecifyCrowd := micro_scrm_enterprise.PlanTaskSpecifyCrowdModel{
				UserId:   value.UserId,
				UserName: value.UserName,
				IsBa:     2,
				ThirdId:  value.ThirdId,
				Type:     crowdType,
				TaskId:   planTask.ID,
				PlanId:   planTask.PlanId,
			}
			planTaskSpecifyCrowds = append(planTaskSpecifyCrowds, planTaskSpecifyCrowd)
		}
		isSpecifyCrowd = 2
	}
	planTask.IsSpecifyCrowd = isSpecifyCrowd
	var planTaskConditionModels []micro_scrm_enterprise.PlanTaskConditionModel // 需要操作的任务条件内容
	var planTaskConditions []PlanTaskConditions                                // 需要操作的条件
	planTaskConditionsByte, _ := json.Marshal(params.PlanTaskConditions)
	err = json.Unmarshal(planTaskConditionsByte, &planTaskConditions)
	if err != nil {
		return err
	}
	var CardId []int
	for _, value := range planTaskConditions {
		contentType := 0
		contentMap := make(map[string]interface{}, 0)
		if value.Content.Text.Content != "" {
			contentType = 1
			contentMap["text"] = map[string]interface{}{
				"content": value.Content.Text.Content,
			}
		}
		if len(value.Content.Image.ImgList) > 0 && value.Content.Image.ImgList[0].PicUrl != "" {
			contentType = 2
			contentMap["image"] = map[string]interface{}{
				"imgList": value.Content.Image.ImgList,
			}
		}
		if value.Content.Voice.MediaId != "" || value.Content.Voice.Format != "" {
			contentType = 3
			contentMap["voice"] = map[string]interface{}{
				"media_id": value.Content.Voice.MediaId,
				"format":   value.Content.Voice.Format,
			}
		}
		if value.Content.Video.MediaId != "" || value.Content.Video.ThumbMediaId != "" || value.Content.Video.VideoUrl != "" {
			contentType = 4
			contentMap["video"] = map[string]interface{}{
				"media_id":        value.Content.Video.MediaId,
				"thumb_media_id":  value.Content.Video.ThumbMediaId,
				"video_url":       value.Content.Video.VideoUrl,
				"video_image_url": value.Content.Video.VideoImageUrl,
			}
		}
		if value.Content.Location.Location_Y != "" || value.Content.Location.Location_X != "" {
			contentType = 5
			contentMap["location"] = map[string]interface{}{
				"location_x": value.Content.Location.Location_X,
				"location_y": value.Content.Location.Location_Y,
				"scale":      value.Content.Location.Scale,
				"label":      value.Content.Location.Label,
			}
		}
		if value.Content.Link.Url != "" || value.Content.Link.Title != "" || value.Content.Link.Description != "" {
			contentType = 6
			contentMap["link"] = map[string]interface{}{
				"id":          value.Content.Link.Id,
				"title":       value.Content.Link.Title,
				"description": value.Content.Link.Description,
				"url":         value.Content.Link.Url,
				"link_url":    value.Content.Link.LinkUrl,
				"cover_url":   value.Content.Link.PicUrl,
				"pic_url":     value.Content.Link.PicUrl,
			}
		}
		if value.Content.Applets.AppId != "" {
			contentType = 7
			contentMap["applets"] = map[string]interface{}{
				"id":              value.Content.Applets.Id,
				"type":            value.Content.Applets.Type,
				"admin_id":        value.Content.Applets.AdminId,
				"title":           value.Content.Applets.Title,
				"mini_program_id": value.Content.Applets.MiniProgramId,
				"app_id":          value.Content.Applets.AppId,
				"app_name":        value.Content.Applets.AppName,
				"path":            value.Content.Applets.Path,
				"cover_url":       value.Content.Applets.CoverUrl,
			}
		}
		if value.Content.Card.Id != 0 {
			contentType = 8 // 贺卡类型
			contentMap["card"] = map[string]interface{}{
				"id":                   value.Content.Card.Id,
				"template_name":        value.Content.Card.TemplateName,
				"template_picture_url": value.Content.Card.TemplatePictureURL,
				"template_content":     value.Content.Card.TemplateContent,
			}
			CardId = append(CardId, value.Content.Card.Id)
		}
		if contentType == 0 && params.GroupSendType == microScrmModel.PLANTASKQYSEND {
			return errors.New("企业群发任务的内容不能为空")
		}
		conditionMap := make(map[string]interface{}, 0)

		if params.ConditionType == microScrmModel.TaskConditionTypeTag {
			conditionMap["tag"] = value.ConditionDetail.Tag
			for _, v := range value.ConditionDetail.Tag {
				tagRelationId := app.UID()
				tasksTagsRelation := workModels.PlanTasksTagRelationModel{
					ID:           tagRelationId,
					EnterpriseID: enterpriseID,
					PlanID:       taskId,
					TagID:        v.TagID,
					StartTime:    plan.StartTime,
					EndTime:      plan.EndTime,
				}
				relationId, err := workDao.NewPlanTasksTagRelation().CreateData(_ctx, tasksTagsRelation)
				if err != nil {
					return err
				}
				if relationId != tagRelationId {
					return errors.New("任务关联标签失败")
				}
			}
		} else {
			// var logical micro_scrm_enterprise.LogicalInfo
			logical := value.ConditionDetail.Logical
			switch params.ConditionType {
			case 2:
				logical.Remark = "群内总人数"
			case 3:
				logical.Remark = "群内新增人数"
			case 4:
				logical.Remark = "群内累计参与活动次数"
			case 5:
				logical.Remark = "群成员累计购买次数"
			case 6:
				logical.Remark = "群活跃度"
			case 7:
				logical.Remark = "退群人数"
			}
			conditionMap["logical"] = logical
		}
		planTaskConditionModels = append(planTaskConditionModels, micro_scrm_enterprise.PlanTaskConditionModel{
			BranchModel: model.BranchModel{
				BaseModel: model.BaseModel{
					ID: uid.ID(_ctx),
				},
			},
			PlanId:      params.PlanId,
			TaskId:      taskId,
			ContentType: contentType,
			Conditions:  helper.MarshalToString(conditionMap),
			Content:     helper.MarshalToString(contentMap),
		})

	}
	// 新增任务
	taskNum := plan.TaskNum + 1
	if err := orm.Service.PlanTask.PlanTasksAdd(enterpriseID, plan.ID, taskNum, planTask, planTaskConditionModels, planTaskSpecifyCrowds); err != nil {
		log.Info(_ctx).Msgf("SFTP链接创建子任务失败: %v", err)
		return errors.New("新增子任务失败")
	}
	if planTask.IsSpecifyCrowd == 2 { // 任务指定了人群,将任务设置到redis
		setPlanTaskRids(_ctx, enterpriseID, planTask.ID, int64(plan.EndTime), planTask.IsSpecifyCrowd, planTaskSpecifyCrowds)
	}

	if (planTask.Dimension == microScrmModel.TaskDimensionImportUnionId || planTask.Dimension == microScrmModel.TaskDimensionImportSaUnionId) && planTask.ConsumerFileId != 0 {
		// 消息队列处理人群包数据
		data := planReq.ImportDataQueueData{
			EnterpriseId:   enterpriseID,
			ConsumerFileId: int64(planTask.ConsumerFileId),
			SendType:       planTask.SendType,
			TaskName:       planTask.Name,
			TaskId:         taskId,
		}
		planService.NewPlanTaskQueueSvc().ImportDataProduce(_ctx, data)

		// go planService.NewPlanTaskImport().SyncTaskImportData(_ctx, enterpriseID, int64(planTask.ConsumerFileId), planTask.SendType, planTask.Name)
		// 异步处理订单数据获取写入
	}
	// 任务创建通知
	planService.NewTaskNotice().CreateTaskNotice(_ctx, planService.TaskMessageData{
		EnterpriseId: enterpriseID,
		TaskId:       planTask.ID,
		PlanId:       planTask.PlanId,
		TaskTitle:    planTask.Name,
		TaskCategory: planTask.TaskCategory,
		Dimension:    planTask.Dimension,
		StartTime:    plan.StartTime,
		EndTime:      plan.EndTime,
	})
	return nil
}

func PlanTasksUpdate(enterpriseID int64, _ctx ctx.Context, params PlanTasksRequest) error {
	if params.ConditionType == 1 && params.Dimension == 2 {
		if params.SendType == 0 || params.TaskType == 0 {
			return errors.New("用户画像纬度的任务未添加推送对象或任务类型")
		}
	}

	/**
	*获取计划内容
	 */
	plan, err := orm.Service.PlanTask.PlanFound(params.PlanId, enterpriseID)
	if err != nil || plan.ID == 0 {
		return errors.New("计划不存在")
	}
	if plan.Status != 0 && params.TaskExpireType != 1 {
		return errors.New("主计划已开始或完成,不能操作")
	}

	if params.TaskExpireType == 1 {
		if int64(plan.StartTime) > params.TaskStartTime || int64(plan.EndTime) < params.TaskStartTime {
			return errors.New("子任务的时间范围应该在主任务的时间范围内")
		}
	}

	planTask := micro_scrm_enterprise.PlanTaskModel{
		BranchModel: model.BranchModel{
			BaseModel: model.BaseModel{
				ID: params.TaskId,
			},
		},
		PlanId:         params.PlanId,
		Name:           params.Name,
		Dimension:      params.Dimension,
		ConditionType:  params.ConditionType,
		IsOpen:         1,
		Status:         0,
		Remark:         params.Remark,
		TaskType:       params.TaskType,
		SendType:       params.SendType,
		GroupSendType:  params.GroupSendType,
		FinishRate:     params.FinishRate,
		TaskStartTime:  params.TaskStartTime,
		TaskEndTime:    params.TaskEndTime,
		TaskExpireType: params.TaskExpireType,
		ConsumerFileId: params.ConsumerFileId,
		AllowSelect:    params.AllowSelect,
	}
	if params.TaskExpireType == 0 {
		planTask.TaskStartTime = int64(plan.StartTime)
		planTask.TaskEndTime = int64(plan.EndTime)
	}
	_planTask := orm.Service.PlanTask.FindplanTaskInfoById(enterpriseID, planTask.ID)
	if _planTask.BaseModel.ID == 0 {
		return errors.New("暂无该任务信息")
	}
	if _planTask.Status != 0 {
		return errors.New("子任务已开始或已完成不能修改")
	}
	planTaskSpecifyCrowds := []micro_scrm_enterprise.PlanTaskSpecifyCrowdModel{}
	isSpecifyCrowd := 1
	if len(params.PlanTaskSpecifyCrowd) > 0 {
		for _, value := range params.PlanTaskSpecifyCrowd {
			if value.UserId == "" && value.ThirdId == 0 {
				return errors.New("指定人群里用户信息有误")
			}
			planTaskSpecifyCrowd := micro_scrm_enterprise.PlanTaskSpecifyCrowdModel{
				UserId:   value.UserId,
				UserName: value.UserName,
				IsBa:     2,
				ThirdId:  value.ThirdId,
				Type:     helper.RuleThree(value.UserId != "", micro_scrm_enterprise.TaskCrowdTypeSa, micro_scrm_enterprise.TaskCrowdTypeDept),
				TaskId:   planTask.ID,
				PlanId:   planTask.PlanId,
			}
			planTaskSpecifyCrowds = append(planTaskSpecifyCrowds, planTaskSpecifyCrowd)
		}
		isSpecifyCrowd = 2
	}
	planTask.IsSpecifyCrowd = isSpecifyCrowd
	/**
	*当任务发生变化后则修改内容
	 */
	if err := orm.Service.PlanTask.PlanTasksUpdate(enterpriseID, planTask); err != nil {
		log.Info(_ctx).Msgf("SFTP链接修改子任务失败: %v", err)
		return errors.New("子任务修改失败")
	}
	if planTask.GroupSendType != _planTask.GroupSendType {
		err := planDao.NewPlanTaskDao().UpdateData(_ctx, enterpriseID, map[string]interface{}{
			"id": _planTask.ID,
		}, map[string]interface{}{"group_send_type": planTask.GroupSendType})
		if err != nil {
			return errors.New("子任务修改失败")
		}
	}
	if planTask.FinishRate != _planTask.FinishRate {
		err := planDao.NewPlanTaskDao().UpdateData(_ctx, enterpriseID, map[string]interface{}{
			"id": _planTask.ID,
		}, map[string]interface{}{"finish_rate": planTask.FinishRate})
		if err != nil {
			return errors.New("子任务修改失败")
		}
	}
	if planTask.TaskExpireType != _planTask.TaskExpireType {
		err := planDao.NewPlanTaskDao().UpdateData(_ctx, enterpriseID, map[string]interface{}{
			"id": _planTask.ID,
		}, map[string]interface{}{"task_expire_type": planTask.TaskExpireType})
		if err != nil {
			return errors.New("子任务修改失败")
		}
	}
	if planTask.AllowSelect != _planTask.AllowSelect {
		err := planDao.NewPlanTaskDao().UpdateData(_ctx, enterpriseID, map[string]interface{}{
			"id": _planTask.ID,
		}, map[string]interface{}{"allow_select": planTask.AllowSelect})
		if err != nil {
			return errors.New("子任务修改失败")
		}
	}
	addplanTaskConditionModels := []micro_scrm_enterprise.PlanTaskConditionModel{}
	updateplanTaskConditionModels := []micro_scrm_enterprise.PlanTaskConditionModel{}
	CardUse := false
	var CardId []int
	if len(params.PlanTaskConditions) > 0 {
		planTaskConditions := []PlanTaskConditions{}
		planTaskConditionsByte, _ := json.Marshal(params.PlanTaskConditions)
		json.Unmarshal(planTaskConditionsByte, &planTaskConditions)
		for _, value := range planTaskConditions {
			contentType := 0
			contentMap := make(map[string]interface{}, 0)
			if value.Content.Text.Content != "" {
				contentType = 1
				contentMap["text"] = map[string]interface{}{
					"content": value.Content.Text.Content,
				}
			}
			if len(value.Content.Image.ImgList) > 0 && value.Content.Image.ImgList[0].PicUrl != "" {
				contentType = 2
				contentMap["image"] = map[string]interface{}{
					"imgList": value.Content.Image.ImgList,
				}
			}
			if value.Content.Voice.MediaId != "" || value.Content.Voice.Format != "" {
				contentType = 3
				contentMap["voice"] = map[string]interface{}{
					"media_id": value.Content.Voice.MediaId,
					"format":   value.Content.Voice.Format,
				}
			}
			if value.Content.Video.MediaId != "" || value.Content.Video.ThumbMediaId != "" || value.Content.Video.VideoUrl != "" {
				contentType = 4
				contentMap["video"] = map[string]interface{}{
					"media_id":       value.Content.Video.MediaId,
					"thumb_media_id": value.Content.Video.ThumbMediaId,
					"video_url":      value.Content.Video.VideoUrl,
				}
			}
			if value.Content.Location.Location_Y != "" || value.Content.Location.Location_X != "" {
				contentType = 5
				contentMap["location"] = map[string]interface{}{
					"location_x": value.Content.Location.Location_X,
					"location_y": value.Content.Location.Location_Y,
					"scale":      value.Content.Location.Scale,
					"label":      value.Content.Location.Label,
				}
			}
			if value.Content.Link.Url != "" || value.Content.Link.Title != "" || value.Content.Link.Description != "" {
				contentType = 6
				contentMap["link"] = map[string]interface{}{
					"title":       value.Content.Link.Title,
					"description": value.Content.Link.Description,
					"url":         value.Content.Link.Url,
					"link_url":    value.Content.Link.LinkUrl,
					"cover_url":   value.Content.Link.PicUrl,
					"pic_url":     value.Content.Link.PicUrl,
				}
			}
			if value.Content.Applets.AppId != "" {
				contentType = 7
				contentMap["applets"] = map[string]interface{}{
					"id":              value.Content.Applets.Id,
					"type":            value.Content.Applets.Type,
					"admin_id":        value.Content.Applets.AdminId,
					"title":           value.Content.Applets.Title,
					"mini_program_id": value.Content.Applets.MiniProgramId,
					"app_id":          value.Content.Applets.AppId,
					"app_name":        value.Content.Applets.AppName,
					"path":            value.Content.Applets.Path,
					"cover_url":       value.Content.Applets.CoverUrl,
				}
			}
			if value.Content.Card.Id != 0 {
				contentType = 8 // 贺卡类型
				contentMap["card"] = map[string]interface{}{
					"id":                   value.Content.Card.Id,
					"template_name":        value.Content.Card.TemplateName,
					"template_picture_url": value.Content.Card.TemplatePictureURL,
					"template_content":     value.Content.Card.TemplateContent,
				}
				CardUse = true
				CardId = append(CardId, value.Content.Card.Id)
			}

			if contentType == 0 && params.GroupSendType == microScrmModel.PLANTASKQYSEND {
				return errors.New("企业群发任务内容不能为空")
			}
			conditionMap := make(map[string]interface{}, 0)
			if params.ConditionType == 1 && params.Dimension == microScrmModel.TaskDimensionTag {
				conditionMap["tag"] = value.ConditionDetail.Tag
			} else {
				switch params.ConditionType {
				case 2:
					value.ConditionDetail.Logical.Remark = "群内总人数"
				case 3:
					value.ConditionDetail.Logical.Remark = "群内新增人数"
				case 4:
					value.ConditionDetail.Logical.Remark = "群内累计参与活动次数"
				case 5:
					value.ConditionDetail.Logical.Remark = "群成员累计购买次数"
				case 6:
					value.ConditionDetail.Logical.Remark = "群活跃度"
				case 7:
					value.ConditionDetail.Logical.Remark = "退群人数"
				}
				conditionMap["logical"] = value.ConditionDetail.Logical
			}
			if value.Id == 0 { // 表示条件内容为空 新增条件
				addplanTaskConditionModels = append(addplanTaskConditionModels, micro_scrm_enterprise.PlanTaskConditionModel{
					BranchModel: model.BranchModel{
						BaseModel: model.BaseModel{
							ID: uid.ID(_ctx),
						},
					},
					PlanId:      params.PlanId,
					TaskId:      params.TaskId,
					ContentType: contentType,
					Conditions:  helper.MarshalToString(conditionMap),
					Content:     helper.MarshalToString(contentMap),
				})
			} else { // 表示对原有的条件进行修改

				updatePlanTaskCondition := micro_scrm_enterprise.PlanTaskConditionModel{
					BranchModel: model.BranchModel{
						BaseModel: model.BaseModel{
							ID: value.Id,
						},
					},
					PlanId:      params.PlanId,
					TaskId:      params.TaskId,
					ContentType: contentType,
					Conditions:  helper.MarshalToString(conditionMap),
					Content:     helper.MarshalToString(contentMap),
				}
				if planTaskConditionInfo := orm.Service.PlanTask.FindPlanTasksConditionsById(enterpriseID, updatePlanTaskCondition.ID); planTaskConditionInfo.BaseModel.ID != 0 {
					if updatePlanTaskCondition.Conditions != planTaskConditionInfo.Conditions || updatePlanTaskCondition.Content != planTaskConditionInfo.Content || updatePlanTaskCondition.ContentType != planTaskConditionInfo.ContentType {
						updateplanTaskConditionModels = append(updateplanTaskConditionModels, updatePlanTaskCondition)
					}
				}

			}
		}
	}
	if len(params.DeletPlanTaskConditions) > 0 || len(addplanTaskConditionModels) > 0 || len(updateplanTaskConditionModels) > 0 {
		if err := orm.Service.PlanTask.PlanTasksConditionsBatchUpdate(enterpriseID, params.DeletPlanTaskConditions, addplanTaskConditionModels, updateplanTaskConditionModels); err != nil {
			log.Info(_ctx).Msgf("任务修改条件失败: %v", err)
			return errors.New("条件新增失败")
		}
	}
	if planTask.Dimension == microScrmModel.TaskDimensionTag {
		// 查询指定用户是否发生变化
		if planTask.IsSpecifyCrowd == 1 {
			if err := orm.Service.PlanTask.PlanTasksIsSpecifyCrowdUpdateDelete(enterpriseID, planTask.ID); err != nil {
				log.Info(_ctx).Msgf("修改任务推送人群错误: %v", err)
				return errors.New("修改任务推送人群错误")
			}
			setPlanTaskRids(_ctx, enterpriseID, planTask.ID, int64(plan.EndTime), planTask.IsSpecifyCrowd, planTaskSpecifyCrowds)
		} else if planTask.IsSpecifyCrowd == 2 {
			logDate, err := orm.Service.PlanTask.GetListByTaskId(enterpriseID, planTask.ID)
			if err != nil {
				return errors.New("修改任务推送人群错误")
			}
			if len(logDate) > 0 {
				logDateByte, _ := json.Marshal(logDate)
				comparePlanTaskSpecifyCrowds := []micro_scrm_enterprise.PlanTaskSpecifyCrowdData{}
				for _, value := range planTaskSpecifyCrowds {
					comparePlanTaskSpecifyCrowds = append(comparePlanTaskSpecifyCrowds, micro_scrm_enterprise.PlanTaskSpecifyCrowdData{
						UserId:   value.UserId,
						UserName: value.UserName,
						PlanId:   value.PlanId,
						TaskId:   value.TaskId,
						IsBa:     value.IsBa,
					})
				}
				planTaskSpecifyCrowdsByte, _ := json.Marshal(comparePlanTaskSpecifyCrowds)
				if string(logDateByte) != string(planTaskSpecifyCrowdsByte) { // 说明指定的人群发生变化需要修改 否则不修改
					if err := orm.Service.PlanTask.PlanTasksIsSpecifyCrowd(enterpriseID, planTask.ID, planTask.IsSpecifyCrowd, planTaskSpecifyCrowds); err != nil {
						log.Info(_ctx).Msgf("修改任务推送人群错误: %v", err)
						return errors.New("修改任务推送人群错误")
					}
					setPlanTaskRids(_ctx, enterpriseID, _planTask.ID, int64(plan.EndTime), planTask.IsSpecifyCrowd, planTaskSpecifyCrowds)
				}
			} else {
				if err := orm.Service.PlanTask.PlanTasksIsSpecifyCrowdUpdateCreate(enterpriseID, planTaskSpecifyCrowds); err != nil {
					log.Info(_ctx).Msgf("修改任务推送人群错误: %v", err)
					return errors.New("修改任务推送人群错误")
				}
				setPlanTaskRids(_ctx, enterpriseID, planTask.ID, int64(plan.EndTime), planTask.IsSpecifyCrowd, planTaskSpecifyCrowds)
			}
		}
	}
	if CardUse { // 处理贺卡与任务的关联,用来判断贺卡是否可以被删除
		go marketingService.NewCustomCard().FindOrCreateUseCard(_ctx, enterpriseID, plan.ID, CardId)
	}
	// 编辑子任务不处理人群包数据
	//if (planTask.Dimension == microScrmModel.TaskDimensionImportUnionId || planTask.Dimension == microScrmModel.TaskDimensionImportSaUnionId) && planTask.ConsumerFileId != 0 {
	//	// 消息队列处理人群包数据
	//	data := planReq.ImportDataQueueData{
	//		EnterpriseId:   enterpriseID,
	//		ConsumerFileId: int64(planTask.ConsumerFileId),
	//		SendType:       planTask.SendType,
	//		TaskName:       planTask.Name,
	//		TaskId:         planTask.ID,
	//	}
	//	planService.NewPlanTaskQueueSvc().ImportDataProduce(_ctx, data)
	//	// go planService.NewPlanTaskImport().SyncTaskImportData(_ctx, enterpriseID, int64(planTask.ConsumerFileId), planTask.SendType, planTask.Name)
	//	// 异步处理订单数据获取写入
	//}
	// 任务创建通知
	planService.NewTaskNotice().CreateTaskNotice(_ctx, planService.TaskMessageData{
		EnterpriseId: enterpriseID,
		TaskId:       planTask.ID,
		PlanId:       planTask.PlanId,
		TaskTitle:    planTask.Name,
		TaskCategory: planTask.TaskCategory,
		Dimension:    planTask.Dimension,
		StartTime:    plan.StartTime,
		EndTime:      plan.EndTime,
	})
	return nil
}
func PlanTasksInfo(c ctx.Context, enterpriseID int64, params PlanTasksInfoRequst) PlanTasksResponse {
	planTask := orm.Service.PlanTask.FindplanTaskInfoById(enterpriseID, params.TaskId)
	response := PlanTasksResponse{
		PlanTaskModel:   *planTask,
		PlanId:          planTask.PlanId,
		TaskId:          planTask.BranchModel.ID,
		Name:            planTask.Name,
		Dimension:       planTask.Dimension,
		ConditionType:   planTask.ConditionType,
		IsSpecifyCrowd:  planTask.IsSpecifyCrowd,
		IsOpen:          planTask.IsOpen,
		Status:          planTask.Status,
		Remark:          planTask.Remark,
		CreatedAt:       planTask.BranchModel.CreatedAt,
		UpdatedAt:       planTask.BranchModel.UpdatedAt,
		TaskType:        planTask.TaskType,
		SendType:        planTask.SendType,
		BirthdayPreTime: planTask.BirthdayPreTime,
		ConsumerFileId:  planTask.ConsumerFileId,
		IsExtract:       planTask.IsExtract,
		DownLoadUrl:     planTask.DownloadUrl,
		TaskCategory:    planTask.TaskCategory,
		GroupSendType:   planTask.GroupSendType,
	}
	if response.TaskId == 0 {
		return response
	}
	if planTask.ConsumerFileId > 0 {
		importFileInfo, err := planDao.NewPlanImportDao().GetImportFile(c, enterpriseID, int64(planTask.ConsumerFileId))
		if err != nil {
			return response
		}
		response.DownLoadUrl = importFileInfo.DownloadURL
		response.FileName = importFileInfo.FileName
		response.FileUrl = importFileInfo.FileURL
	}

	if response.IsSpecifyCrowd == 2 {
		if users, err := orm.Service.PlanTask.GetListByTaskId(enterpriseID, response.TaskId); err == nil {
			for _, value := range users {
				response.PlanTaskSpecifyCrowd = append(response.PlanTaskSpecifyCrowd, PlanTaskSpecifyCrowd{
					UserId:   value.UserId,
					UserName: value.UserName,
					ThirdId:  value.ThirdId,
				})
			}
		}
	}
	planTaskConditions := orm.Service.PlanTask.PlanTasksConditionsById(enterpriseID, response.TaskId)

	qyTagIdMap := make(map[string]string)
	tags := make([]string, 0)
	many, _ := tagsDao.NewWorkExternalTag().GetMany(c, map[string]interface{}{
		"enterprise_id": enterpriseID,
		"type":          workModels.WorkExternalTagTypeQy,
	})
	for _, item := range many {
		qyTagIdMap[item.TagID] = item.Name
	}
	for _, planTaskCondition := range planTaskConditions {
		conditions := make(map[string]interface{}, 0)
		content := make(map[string]interface{}, 0)
		json.Unmarshal([]byte(planTaskCondition.Conditions), &conditions)
		json.Unmarshal([]byte(planTaskCondition.Content), &content)
		// 判断任务条件是否是标签response.Dimension是2表示是用户画像标签
		if response.Dimension == 2 {
			planTaskTags := PlanTaskTag{}
			json.Unmarshal([]byte(planTaskCondition.Conditions), &planTaskTags)
			tagIds := []int64{}
			for _, tag := range planTaskTags.Tag {
				tagIds = append(tagIds, tag.TagID)
				if value, ok := qyTagIdMap[tag.QwTagID]; ok {
					tags = append(tags, value)
				}
			}
			// 查询标签
			newPlanTaskTags := PlanTaskTag{}
			mysSqlTags, err := orm.Service.Tag.FindWorkTags(enterpriseID, tagIds)
			if err != nil || len(mysSqlTags) == 0 { // 如果没有找到表示标签已被删除
				for _, tag := range planTaskTags.Tag {
					tag.IsDelete = true
					newPlanTaskTags.Tag = append(newPlanTaskTags.Tag, tag)
				}
			} else {
				for _, tag := range planTaskTags.Tag {
					tagFlag := false
					for _, mysqltag := range mysSqlTags {
						if tag.TagID == mysqltag.ID {
							if mysqltag.DeletedAt > 0 {
								tag.IsDelete = true
							}
							tagFlag = true
							break
						}
					}
					if !tagFlag {
						tag.IsDelete = true
					}
					newPlanTaskTags.Tag = append(newPlanTaskTags.Tag, tag)
				}
			}
			if tagsByte, err := json.Marshal(newPlanTaskTags); err == nil {
				newconditions := map[string]interface{}{}
				json.Unmarshal(tagsByte, &newconditions)
				conditions = newconditions
			}
		}
		response.PlanTaskConditions = append(response.PlanTaskConditions, PlanTaskConditionResponse{
			Id:            planTaskCondition.BranchModel.ID,
			PlanId:        planTaskCondition.PlanId,
			TaskId:        planTaskCondition.TaskId,
			Conditions:    conditions,
			Content:       content,
			ConditionType: planTaskCondition.ContentType,
			CreatedAt:     planTaskCondition.BranchModel.CreatedAt,
			UpdatedAt:     planTaskCondition.BranchModel.UpdatedAt,
		})
	}
	if planTask.QyTagId != "" {
		if value, ok := qyTagIdMap[planTask.QyTagId]; ok {
			tags = append(tags, value)
		}
	}
	response.Tags = tags
	return response
}
func PlanTasksDelet(enterpriseID int64, params PlanTasksInfoRequst, _ctx ctx.Context) error {
	planTask := orm.Service.PlanTask.FindplanTaskInfoById(enterpriseID, params.TaskId)
	if planTask.ID == 0 {
		return errors.New("任务不存在")
	}
	if planTask.Status != 0 {
		return errors.New("任务已进行或完成")
	}
	/**
	*获取计划内容
	 */
	plan, err := orm.Service.PlanTask.PlanFound(planTask.PlanId, enterpriseID)
	if err != nil || plan.ID == 0 {
		return errors.New("计划不存在")
	}
	if plan.Status != 0 {
		return errors.New("计划已开始或完成,不能操作")
	}
	if plan.TaskNum == 0 {
		return errors.New("任务编号错误")
	}
	taskNum := plan.TaskNum - 1
	if err := orm.Service.PlanTask.PlanTasksDelet(enterpriseID, planTask.ID, time.Now().Unix(), plan.ID, taskNum); err != nil {
		log.Info(_ctx).Msgf("SFTP链接删除子任务: %v", err)
		return errors.New("子任务修改失败")
	}

	_ = workDao.NewPlanTasksTagRelation().DeleteData(_ctx, map[string]interface{}{
		"enterprise_id": enterpriseID,
		"plan_id":       params.TaskId,
	})
	return nil
}
func PlanTasksConditionsCompletion(enterpriseID int64, params PlanTasksConditionsRequst) ([]micro_scrm_enterprise.PlanTasksConditionsAllotLog, int64) {
	size := params.Size
	if size == 0 {
		size = 10
	}
	page := params.Page
	if page == 0 {
		page = 1
	}
	if params.CompletionStartTime == params.CompletionEndTime && params.CompletionStartTime != 0 {
		params.CompletionEndTime = params.CompletionStartTime + (60*60*24 - 1)
	}
	if params.PlanStartTime == params.PlanEndTime && params.PlanStartTime != 0 {
		params.PlanEndTime = params.PlanStartTime + (60*60*24 - 1)
	}
	result, total := orm.Service.PlanTask.PlanTasksConditionsCompletion(enterpriseID, params.Id, params.DepartmentId, params.PlanStartTime, params.PlanEndTime, params.CompletionStartTime, params.CompletionEndTime, params.BaName, params.IsFinished, page, size)
	userIds := []string{}
	for _, value := range result {
		userIds = append(userIds, value.UserId)
	}
	result1 := []micro_scrm_enterprise.PlanTasksConditionsAllotLog{}
	chatNameMap := make(map[int64]group_chat.WorkGroupChatModel, 0)
	for _, v := range orm.Service.PlanTask.GetGroupChatName(enterpriseID) {
		chatNameMap[int64(v.ID)] = v
	}
	one, err := orm.Service.PlanTask.FindTaskConditionsById(enterpriseID, params.Id)
	if err != nil {
		return result1, total
	}
	planId := one.PlanId
	if len(userIds) > 0 {
		if users, err := orm.Service.PlanTask.PlanTasksConditionsCompletioncustomerNumv1(enterpriseID, params.Id, planId, userIds); err == nil {
			for _, v := range result {
				if val, ok := chatNameMap[v.ChatID]; ok {
					v.ChatName = val.Name
					v.Count = leaveCount(v.StartTime, v.EndTime, val.ID)
				}
				for _, v1 := range users {
					if v.UserId == v1.UserId {
						v.CustomerNum = v1.CustomerNum
						break
					}
				}
				result1 = append(result1, v)
			}

		}
	}
	return result1, total
}

func PlanTasksConditionsCompletionV2(enterpriseID int64, params PlanTasksConditionsRequst, p *PTaskService, cc ctx.Context) (newResult []micro_scrm_enterprise.PlanTasksConditionsAllotLog, count int64) {
	size := params.Size
	if size == 0 {
		size = 10
	}
	page := params.Page
	if page == 0 {
		page = 1
	}
	if params.CompletionStartTime == params.CompletionEndTime && params.CompletionStartTime != 0 {
		params.CompletionEndTime = params.CompletionStartTime + (60*60*24 - 1)
	}
	if params.PlanStartTime == params.PlanEndTime && params.PlanStartTime != 0 {
		params.PlanEndTime = params.PlanStartTime + (60*60*24 - 1)
	}
	// 根据部门ID，获取所有下级部门id集合
	// DepartmentIds := orm.Service.Departments.GetDepartmentsChildrenIds(enterpriseID, params.DepartmentId)
	var searchUserIds []string
	if params.BaName != "" {
		searchUserIds, _ = membersDao.NewMembers().GetUserIdsByName(cc, enterpriseID, params.BaName)
	}
	if params.DepartmentId != 0 {
		memberLists, err := membersDao.NewMembers().GetMemberByDepartmentId(cc, enterpriseID, []int64{params.DepartmentId})
		if err != nil {
			return nil, 0
		}
		for _, v := range memberLists {
			info := v
			if !helper.InArray[string](info.MembersUserid, searchUserIds) {
				searchUserIds = append(searchUserIds, info.MembersUserid)
			}
		}
	}
	info, err := planDao.NewPlanTaskDao().GetConditionInfo(cc, enterpriseID, map[string]interface{}{"id": params.Id})
	if err != nil {
		return nil, 0
	}
	taskInfo, err := planDao.NewPlanTaskDao().GetOne(cc, enterpriseID, map[string]interface{}{"id": info.TaskID})
	if err != nil {
		return nil, 0
	}

	result, total := orm.Service.PlanTask.CustomerTaskListData(enterpriseID, params.Id, params.PlanStartTime, params.PlanEndTime, params.CompletionStartTime, params.CompletionEndTime, searchUserIds, params.IsFinished, page, size)
	userIds := make([]string, 0)
	for _, value := range result {
		userIds = append(userIds, value.UserId)
	}
	memberInfoLists, err := membersDao.NewMembers().GetMemberInfoAndDepartmentByUserIds(cc, enterpriseID, userIds)
	if err != nil {
		return nil, 0
	}
	memberInfoMap := make(map[string]membersModel.MemberDepartment)
	for _, v := range memberInfoLists {
		if _, ok := memberInfoMap[v.Userid]; !ok {
			if len(v.MemberDepartments) > 0 {
				for _, md := range v.MemberDepartments {
					if md.DeletedAt == 0 && md.Department.DeletedAt == 0 {
						memberInfoMap[v.Userid] = membersModel.MemberDepartment{
							DepartmentName: md.Department.Name,
							MembersName:    v.Name,
						}
						break
					}
					memberInfoMap[v.Userid] = membersModel.MemberDepartment{
						DepartmentName: md.Department.Name,
						MembersName:    v.Name,
					}
				}
			}
		}
	}
	count = total
	chatNameMap := make(map[int64]group_chat.WorkGroupChatModel, 0)
	for _, v := range orm.Service.PlanTask.GetGroupChatName(enterpriseID) {
		chatNameMap[v.ID] = v
	}
	if len(userIds) > 0 {
		for _, v := range result {
			v.TaskCategory = taskInfo.TaskCategory
			v.GroupSendType = taskInfo.GroupSendType
			// // 查询统计记录
			// planStatistic, _ := p.GetTaskStatistics(cc, enterpriseID, v.ID, time.Now().Unix())
			// if planStatistic != nil {
			// 	v.TotalReachCustomerNumber = planStatistic.TotalReachCustomerNumber
			// 	v.TotalOrderCustomerNumber = planStatistic.TotalOrderCustomerNumber
			// 	v.TotalOrderNumber = planStatistic.TotalOrderNumber
			// 	v.TotalOrderFee = planStatistic.TotalOrderFee
			// }
			if taskInfo.TaskCategory == microScrmModel.PLANTASKCONTACT {
				v.TotalReachCustomerNumber = v.ContactedCount
			}
			if taskInfo.TaskCategory == microScrmModel.PLANTASKMOMENT {
				v.TotalReachCustomerNumber = v.MomentLookCount
			}
			if memberInfo, ok := memberInfoMap[v.UserId]; ok {
				v.DepartmentName = memberInfo.DepartmentName
				v.Name = memberInfo.MembersName
			}
			v.FinishRate = fmt.Sprintf("%s%%", v.FinishRate)
			newResult = append(newResult, v)
		}
	}
	return
}

func leaveCount(startTime, endTime, groupChatID int64) (count rune) {
	var i rune
	for _, v := range orm.Service.PlanTask.GetLeaveCount() {
		if v.GroupChatID == groupChatID && v.LeaveTime >= startTime && v.LeaveTime <= endTime {
			i += 1
		}
		continue
	}
	return i
}

func PlanTasksConditionsUserCompletion(c ctx.Context, enterpriseID int64, params PlanTasksConditionsBaRequst) (result []PlanTasksConditionsBaRespone, total int64) {
	size := params.Size
	if size == 0 {
		size = 10
	}
	page := params.Page
	if page == 0 {
		page = 1
	}
	one, err := planDao.NewPlanTaskAllotLogDao().GetOne(c, enterpriseID, map[string]interface{}{
		"id": params.Id,
	})
	if err != nil {
		return
	}

	numberInfo, count, err := orm.Service.PlanTask.GetCorrespondCustomerNumber(enterpriseID, params.Id, one.PlanId, params.Size, params.Page, params.Customer_name)
	if err != nil {
		return
	}
	total = count
	qyTagIdMap := make(map[string]string)
	many, _ := tagsDao.NewWorkExternalTag().GetMany(ctx.Wrap(c), map[string]interface{}{
		"enterprise_id": enterpriseID,
		"type":          workModels.WorkExternalTagTypeQy,
	})

	for _, item := range many {
		qyTagIdMap[item.TagID] = item.Name
	}
	for _, value := range numberInfo {
		tagName := value.TagName
		if value.WwTagId != "" {
			tagName = qyTagIdMap[value.WwTagId]
		}
		resultmap := PlanTasksConditionsBaRespone{
			TagName:   tagName,
			Name:      Namesensitive(value.Name),
			Avatar:    value.Avatar,
			LeaveTime: value.CreatedAt,
		}
		if value.Sex == 1 {
			resultmap.Sex = "男"
		} else if value.Sex == 2 {
			resultmap.Sex = "女"
		} else {
			resultmap.Sex = "未知"
		}
		if value.Mobile != "" && len(value.Mobile) >= 8 {
			resultmap.Mobile = value.Mobile[:3] + "****" + value.Mobile[7:]
		}

		result = append(result, resultmap)
	}
	return
}
func Namesensitive(str string) string {
	strLen := len(str)
	newStr := ""
	Hanflag := true
	for i, v := range str {
		if unicode.Is(unicode.Han, v) {
			if i == 0 {
				newStr = str[:3]
			} else if i == strLen-3 {
				newStr += str[strLen-3:]
			} else {
				newStr += "*"
			}
		} else {
			Hanflag = false
			break
			// if i == 0 {
			// 	newStr = str[0:1]
			// } else if i == strLen-1 {
			// 	newStr += str[strLen-1:]
			// } else {
			// 	if i%3 == 0 {
			// 		newStrtidaifu += "*"
			// 	}
			// 	newStr += "*"
			// }
		}
	}
	if !Hanflag {
		return str
	} else {
		if strLen == 6 {
			return str[:3] + "*"
		}
		return newStr
	}
}
func setPlanTaskRids(ctx ctx.Context, enterpriseID, taskId, endTime int64, isSpecifyCrowd int, planTaskSpecifyCrowds []micro_scrm_enterprise.PlanTaskSpecifyCrowdModel) {
	key := fmt.Sprintf("%s::%d::plan_task::specify_crowd::%d", config.CfgData.App, enterpriseID, taskId)
	redisCache := cache.NewRedisCache(micro_redis.Client())
	if isSpecifyCrowd == 1 {
		redisCache.Set(ctx, key, "", 0)
		return
	}
	redisTime := endTime - time.Now().Unix()
	if redisTime <= 0 {
		return
	}
	newplanTaskSpecifyCrowds := []string{}
	for _, value := range planTaskSpecifyCrowds {
		newplanTaskSpecifyCrowds = append(newplanTaskSpecifyCrowds, value.UserId)
	}
	jsonData, jsonErr := json.Marshal(newplanTaskSpecifyCrowds)
	if jsonErr != nil {
		return
	}
	redisCache.Set(ctx, key, jsonData, redisTime)
}
