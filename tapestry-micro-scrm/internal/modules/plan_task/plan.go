package plan_task

import (
	"errors"
	"fmt"
	"micro-scrm/internal/app"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	workDao "micro-scrm/internal/app/dao/work"
	planReq "micro-scrm/internal/app/http/request/plan"
	"micro-scrm/internal/app/models"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	workModels "micro-scrm/internal/app/models/work"
	planService "micro-scrm/internal/app/service/plan"
	tagsSvc "micro-scrm/internal/app/service/tags"
	"micro-scrm/internal/pkg/helper"
	"strconv"
	"time"

	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/pkg/uid"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

const (
	specifyCrowdBaNum = 1000
)

func PlanAdd(c ctx.Context, enterpriseId int64, info ReceivePlan) error {
	var isOpen = 1
	if !info.WhetherOpen {
		isOpen = 0
	}
	planId := app.UID()
	newInfo := micro_scrm_enterprise.PlanModel{}
	newInfo.BranchModel.ID = planId
	newInfo.BranchModel.CreatedAt = time.Now().Unix()
	newInfo.BranchModel.UpdatedAt = time.Now().Unix()
	newInfo.BranchModel.DeletedAt = 0
	newInfo.Name = info.PlanName
	newInfo.StartTime = int(info.StartTime)
	newInfo.EndTime = int(info.EndTime)
	newInfo.Remark = info.Description
	newInfo.IsOpen = isOpen
	newInfo.Status = 0
	newInfo.OperationUserId = info.OperationUserId
	newInfo.Category = microScrmModel.PlanMaTypeNo
	if _, err := orm.Service.PlanTask.AddPlan(newInfo, enterpriseId); err != nil {
		return err
	}
	// 创建一个task_condition_data的分表
	err := planService.NewPlanTaskConditionData().CreateShardTable(c, enterpriseId, planId)
	return err
}

func PlanEdit(c ctx.Context, enterpriseID int64, info ReceivePlanModel) error {
	newInfo := micro_scrm_enterprise.PlanModel{}
	newInfo.ID = info.ID
	newInfo.StartTime = int(info.StartTime)
	newInfo.IsOpen = 1
	if !info.WhetherOpen {
		newInfo.IsOpen = 0
	}
	newInfo.Status = info.Status
	newInfo.Name = info.PlanName
	newInfo.EndTime = int(info.EndTime)
	newInfo.CreatedAt = info.CreatedAt
	newInfo.Remark = info.Description
	newInfo.OperationUserId = info.OperationUserId
	newInfo.TaskNum = info.TaskNum
	newInfo.Category = info.Category
	planTaskInfo, err := orm.Service.PlanTask.PlanFound(info.ID, enterpriseID)
	if err != nil {
		return err
	}
	if planTaskInfo.Status == 1 {
		return errors.New("任务已开始,不可操作")
	}
	taskInfo, err := microScrmDao.NewPlanTaskDao().GetOne(c, enterpriseID, "plan_id=? and task_expire_type=1 and (task_start_time<? or task_end_time>?)", info.ID, info.StartTime, info.EndTime)
	if err != nil {
		return err
	}
	if taskInfo.ID > 0 {
		return errors.New("主任务的有效期时间应包含子任务的有效期时间")
	}
	if _, err := orm.Service.PlanTask.PlanEdit(newInfo, enterpriseID); err != nil {
		return err
	}
	err = microScrmDao.NewPlanTaskDao().UpdateData(c, enterpriseID, map[string]interface{}{
		"plan_id":          info.ID,
		"task_expire_type": 0,
		"deleted_at":       0,
	}, map[string]interface{}{
		"task_start_time": info.StartTime,
		"task_end_time":   info.EndTime,
	})
	return err
}

func PlanRemove(c ctx.Context, planId, enterpriseID int64) error {

	if err := orm.Service.PlanTask.DelPlan(planId, time.Now().Unix(), enterpriseID); err != nil {
		return err
	}
	if data, err := orm.Service.PlanTask.GetPlanTaskId(planId, enterpriseID); err != nil {
		return err
	} else {
		if len(data) == 0 {
			return errors.New("Data length is 0")
		} else {
			if err := orm.Service.PlanTask.DelPlanTask(data, time.Now().Unix(), enterpriseID); err != nil {
				return err
			}
			err = workDao.NewPlanTasksTagRelation().DeleteData(c, map[string]interface{}{
				"enterprise_id": enterpriseID,
				"plan_id":       data,
			})
		}
	}
	return nil
}

func PlanCopy(ctx ctx.Context, enterpriseID int64, info ReceivePlanModel) error {
	planId := info.ID
	newInfo := micro_scrm_enterprise.PlanModel{}
	newInfo.ID = uid.ID(ctx)
	newInfo.StartTime = int(info.StartTime)
	newInfo.IsOpen = 1
	if !info.WhetherOpen {
		newInfo.IsOpen = 0
	}
	newInfo.Status = info.Status
	newInfo.Name = info.PlanName
	newInfo.EndTime = int(info.EndTime)
	planEndTIme := info.EndTime
	newInfo.CreatedAt = info.CreatedAt
	newInfo.Remark = info.Description
	newInfo.OperationUserId = info.OperationUserId
	newInfo.Category = info.Category
	planStartTime := int(info.StartTime)
	newId, err := orm.Service.PlanTask.AddPlan(newInfo, enterpriseID)
	if err != nil {
		return err
	}
	planTaskInfo, err := orm.Service.PlanTask.PlanTaskFound(planId, enterpriseID)
	if err != nil {
		return err
	}
	var tasksInfo = make([]micro_scrm_enterprise.PlanTaskModel, 0)
	var newConditionInfo = make([]micro_scrm_enterprise.PlanTaskConditionModel, 0)
	conditionCopy, err := orm.Service.PlanTask.GetPlanTaskCondition(planId, enterpriseID)
	specifyCrowdBas := make([]micro_scrm_enterprise.PlanTaskSpecifyCrowdModel, 0)
	baUserNum := 0
	nowTime := time.Now().Unix()
	for _, v := range planTaskInfo {
		info := v
		info.ID = app.UID()
		info.PlanId = newId
		info.Status = 0
		info.TaskIsPush = 0
		info.CreatedAt = nowTime
		info.UpdatedAt = nowTime
		info.MomentStatus = 0
		info.MomentJobId = ""
		info.FinishRate = v.FinishRate
		if info.TaskExpireType == 0 {
			info.TaskStartTime = int64(newInfo.StartTime)
			info.TaskEndTime = int64(newInfo.EndTime)
		}
		info.AllowSelect = v.AllowSelect
		redisspecifyCrowdBa := []micro_scrm_enterprise.PlanTaskSpecifyCrowdModel{}
		if info.Dimension == microScrmModel.TaskDimensionTag && info.TaskType == microScrmModel.PLANTCOUNT && info.PlanSendTime > 0 {
			continue
		}
		addTag := false
		tagType := workModels.WorkExternalTagTypeQy
		endTime := 0
		tagName := info.TagName
		//  判断任务是不是要生成标签的任务。
		if info.GroupSendType == microScrmModel.PLANTASKTAGSEND && info.IsAddTag == 1 { // 如果标签任务群发,并且要添加标签
			endTime = int(planEndTIme + 10*86400)
			tagName = fmt.Sprintf("%s_1", info.TagName)
			addTag = true
		}
		if info.TaskCategory == microScrmModel.PLANTASKQYSEND {
			if info.GroupSendType == microScrmModel.PLANTASKMOMENT && helper.InArray(info.Dimension, []int{3, 4}) {
				tagName = "朋友圈临时标签" + strconv.Itoa(int(time.Now().Unix()))
				tagType = workModels.MomentTaskTempTag
				info.TagExpireType = microScrmModel.TAG_TASK_SEND_EXPIRE
				addTag = true
			}
		}
		if addTag {
			wxtagId, tagId, err := tagsSvc.NewWorkTags().CreateTaskTag(ctx, enterpriseID, "任务标签", tagName, endTime, tagType)
			if err != nil {
				continue
			}
			tagRelationId := app.UID()
			// todo 任务标签关联表插入 tagId
			tasksTagsRelation := workModels.PlanTasksTagRelationModel{
				ID:           tagRelationId,
				EnterpriseID: enterpriseID,
				PlanID:       info.ID,
				TagID:        tagId,
				StartTime:    planStartTime,
				EndTime:      endTime,
			}
			relationId, err := workDao.NewPlanTasksTagRelation().CreateData(ctx, tasksTagsRelation)
			if err != nil {
				continue
			}
			if relationId != tagRelationId {
				continue
			}
			info.QyTagId = wxtagId
			info.TagName = tagName
		}

		if v.IsSpecifyCrowd == 2 {
			// 根据任务编号查询指定的导购
			logDates, err := orm.Service.PlanTask.GetListByTaskId(enterpriseID, v.ID)
			if err == nil && len(logDates) > 0 {
				for _, logDate := range logDates {
					baUserNum++
					specifyCrowdBa := micro_scrm_enterprise.PlanTaskSpecifyCrowdModel{
						PlanId:   info.PlanId,
						TaskId:   info.ID,
						UserId:   logDate.UserId,
						ThirdId:  logDate.ThirdId,
						Type:     logDate.Type,
						UserName: logDate.UserName,
						IsBa:     logDate.IsBa,
					}
					redisspecifyCrowdBa = append(redisspecifyCrowdBa, specifyCrowdBa)
					specifyCrowdBas = append(specifyCrowdBas, specifyCrowdBa)
				}
			}
			// 设置缓存
			if v.IsSpecifyCrowd == 2 {
				setPlanTaskRids(ctx, enterpriseID, info.ID, planEndTIme, info.IsSpecifyCrowd, redisspecifyCrowdBa)
			}
			if baUserNum == specifyCrowdBaNum {
				if err := orm.Service.PlanTask.PlanTasksIsSpecifyCrowdUpdateCreate(enterpriseID, specifyCrowdBas); err == nil {
					baUserNum = 0
					specifyCrowdBas = make([]micro_scrm_enterprise.PlanTaskSpecifyCrowdModel, 0)
				}
			}
		}
		if (info.Dimension == microScrmModel.TaskDimensionImportUnionId || info.Dimension == microScrmModel.TaskDimensionImportSaUnionId) && info.ConsumerFileId != 0 {
			// 消息队列处理人群包数据
			fileInfo, err := microScrmDao.NewPlanImportDao().GetImportFile(ctx, enterpriseID, int64(info.ConsumerFileId))
			if err != nil {
				return err
			}
			fileData := microScrmModel.TaskImportFileModel{
				BaseModel: models.BaseModel{
					CreatedAt: time.Now().Unix(),
					UpdatedAt: time.Now().Unix(),
					DeletedAt: 0,
				},
				ID:        0,
				FileName:  fileInfo.FileName,
				FileType:  fileInfo.FileType,
				FileURL:   fileInfo.FileURL,
				IsExtract: 0,
				TaskType:  fileInfo.TaskType,
			}
			fileId, err := microScrmDao.NewPlanImportDao().SaveImportFile(ctx, enterpriseID, fileData)
			if err != nil {
				return err
			}
			info.ConsumerFileId = int(fileId)
			info.IsExtract = 0
		}
		tasksInfo = append(tasksInfo, info)
		for _, model := range conditionCopy {
			if v.ID == model.TaskId {
				newCondition := model
				newCondition.ID = app.UID()
				newCondition.TaskId = info.ID
				newCondition.PlanId = newId
				newConditionInfo = append(newConditionInfo, newCondition)
			}
		}
	}
	if err := orm.Service.PlanTask.PlanTaskCopy(tasksInfo, enterpriseID); err != nil {
		return err
	}
	if err := orm.Service.PlanTask.PlanTaskConditionCopy(newConditionInfo, enterpriseID); err != nil {
		return err
	}
	if len(specifyCrowdBas) > 0 {
		err = orm.Service.PlanTask.PlanTasksIsSpecifyCrowdUpdateCreate(enterpriseID, specifyCrowdBas)
		if err != nil {
			return err
		}
	}
	for _, model := range tasksInfo {
		info := model
		if (info.Dimension == microScrmModel.TaskDimensionImportUnionId || info.Dimension == microScrmModel.TaskDimensionImportSaUnionId) && info.ConsumerFileId != 0 {
			// 消息队列处理人群包数据
			data := planReq.ImportDataQueueData{
				EnterpriseId:   enterpriseID,
				ConsumerFileId: int64(info.ConsumerFileId),
				SendType:       info.SendType,
				TaskName:       info.Name,
				TaskId:         info.ID,
			}
			planService.NewPlanTaskQueueSvc().ImportDataProduce(ctx, data)
			// go planService.NewPlanTaskImport().SyncTaskImportData(_ctx, enterpriseID, int64(planTask.ConsumerFileId), planTask.SendType, planTask.Name)
			// 异步处理订单数据获取写入
		}
	}
	taskNum := len(tasksInfo)
	err = microScrmDao.NewPlanTaskDao().UpdatePlanData(ctx, enterpriseID, map[string]interface{}{"id": newId}, map[string]interface{}{"task_num": taskNum})

	err = planService.NewPlanTaskConditionData().CreateShardTable(ctx, enterpriseID, newInfo.ID)
	// 任务创建通知
	planService.NewTaskNotice().CreateTaskNotice(ctx, planService.TaskMessageData{
		EnterpriseId: enterpriseID,
		PlanId:       newInfo.ID,
		TaskTitle:    fmt.Sprintf("%s(复制主任务%d)", newInfo.Name, info.ID),
		StartTime:    newInfo.StartTime,
		EndTime:      newInfo.EndTime,
	})
	return err
}

func PlanUpdateStatus(c ctx.Context, info PlanStatus, enterpriseID int64) error {
	isOpen := 1
	if !info.IsOpen {
		isOpen = 0
	}
	if err := orm.Service.PlanTask.UpdatePlanStatus(info.ID, int64(isOpen), enterpriseID); err != nil {
		return err
	}
	return nil
}
