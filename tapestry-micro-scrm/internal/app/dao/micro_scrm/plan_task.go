package microScrmDao

import (
	"encoding/json"
	"fmt"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	openReq "micro-scrm/internal/app/http/request/open"
	planReq "micro-scrm/internal/app/http/request/plan"
	planResp "micro-scrm/internal/app/http/response/plan"
	"micro-scrm/internal/app/models"
	"micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/helper"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
)

type PlanTaskDao struct {
	dao       dao.Dao
	tableName string
}

func NewPlanTaskDao() *PlanTaskDao {
	return &PlanTaskDao{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: microScrmModel.PlanTasksModel{}.TableName(),
	}
}

func (d *PlanTaskDao) GetManyPreload(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]microScrmModel.PlanTasksModel, error) {
	var result []microScrmModel.PlanTasksModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Preload("Plan", func(db *gorm.DB) *gorm.DB {
			return db.Scopes(d.dao.BranchTable(microScrmModel.PlansModel{}.TableName(), enterpriseId))
		}).
		Where(query, args...).
		Find(&result).Error
	return result, err
}

// UpdatePlanTaskISExtract 更新子任务的is_extract状态为1 ,导入的模板内容解析完毕
func (d *PlanTaskDao) UpdatePlanTaskISExtract(c app.Ctx, enterpriseID, fileId int64) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseID)).
		Where("consumer_file_id=?", fileId).Update("is_extract", 1).Error
}

// UpdatePlanTaskPushStatus 更新计划任务的状态为已推送
func (d *PlanTaskDao) UpdatePlanTaskPushStatus(c app.Ctx, enterpriseID, taskID int64) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseID)).
		Where("id=?", taskID).Updates(map[string]interface{}{"task_is_push": microScrmModel.PlanTaskPushSuccess}).Error
}

// UpdatePlanTaskPush 更新plan_task的task_is_push 状态为推送中
func (d *PlanTaskDao) UpdatePlanTaskPush(c app.Ctx, enterpriseID, taskID int64) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseID)).
		Where("id=?", taskID).Updates(map[string]interface{}{"task_is_push": microScrmModel.PlanTaskPushRunning}).Error
}

// UpdatePlanTaskPushFailed 更新plan_task的task_is_push 状态为下发失败
func (d *PlanTaskDao) UpdatePlanTaskPushFailed(c app.Ctx, enterpriseID, taskID int64) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseID)).
		Where("id=?", taskID).Updates(map[string]interface{}{"task_is_push": microScrmModel.PlanTaskPushFailed}).Error
}

// GetPlanDetail 获取主任务详情
func (d *PlanTaskDao) GetPlanDetail(c app.Ctx, request planReq.PlanInfo) (data planReq.PlanInfoRes, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.PlansModel{}.TableName(), request.EnterpriseId)).Where("id=?", request.Id).Find(&data).Error
	return
}

// GetTaskTypeInfo 查询主任务下是否有持续性子任务
func (d *PlanTaskDao) GetTaskTypeInfo(c ctx.Context, request planReq.PlanInfo) (total int64, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, request.EnterpriseId)).Where("plan_id=? and task_type=2 and deleted_at=0", request.Id).Count(&total).Error
	return
}

func (d *PlanTaskDao) GetOne(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) (microScrmModel.PlanTasksModel, error) {
	var result microScrmModel.PlanTasksModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *PlanTaskDao) UpdateData(c app.Ctx, enterpriseId int64, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query).UpdateColumns(data).Error
	return err
}

// CreateMaPlanTask 创建MA的定时任务，返回plan_tasks表id 和plans 表id
func (d *PlanTaskDao) CreateMaPlanTask(c ctx.Context, enterpriseId int64, request openReq.AddPlanReq, info microScrmModel.PlanTaskTemplateModel) (int64, int64, error) {
	planId := app.UID()
	taskId := int64(0)
	startTime := time.Now().Unix()
	endTime := startTime + 10*86400*365
	planData := microScrmModel.PlansModel{
		BaseModel: models.BaseModel{
			ID: planId,
		},
		Name:            request.PlanName,
		TaskNum:         0,
		OperationUserID: info.AdminID,
		StartTime:       int(startTime),
		EndTime:         int(endTime),
		IsOpen:          1,
		Status:          0,
		Remark:          request.PlanRemark,
		IsPush:          0,
		Category:        microScrmModel.PlanMaTypeNewYes,
		PlanTemplateID:  info.ID,
	}
	conditions := make(map[string]interface{})
	conditions["tag"] = []string{}
	err := d.dao.WithContext(c).Transaction(func(tx *gorm.DB) error {
		err := tx.Scopes(d.dao.BranchTable(microScrmModel.PlansModel{}.TableName(), enterpriseId)).Create(&planData).Error
		return err
	})

	return planId, taskId, err
}

func (d *PlanTaskDao) CreatePlanByContent(c ctx.Context, enterpriseId int64, request openReq.CreatePlanByContentReq, startTime, endTime int64) (int64, error) {
	planId := app.UID()
	planData := microScrmModel.PlansModel{
		BaseModel: models.BaseModel{
			ID: planId,
		},
		Name:           request.PlanName,
		TaskNum:        0,
		StartTime:      int(startTime),
		EndTime:        int(endTime),
		IsOpen:         1,
		Status:         microScrmModel.PlanNotStart,
		Remark:         request.PlanRemark,
		IsPush:         0,
		Category:       microScrmModel.PlanMaTypeNewYes,
		PlanTemplateID: 0,
	}
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.PlansModel{}.TableName(), enterpriseId)).Create(&planData).Error
	return planId, err
}

func (d *PlanTaskDao) CreateTaskByCrowdId(c ctx.Context, request openReq.TaskCrowdReq, enterpriseId, taskId int64, info microScrmModel.PlanTaskTemplateModel, backUpInfo microScrmModel.PlanTaskTemplateBackupModel, tagExpireType int, tagName, wxTagId string) error {
	dimension := helper.RuleThree[uint32](request.PlanSendType == 1, 3, 4)
	sendTime := int64(0)
	if request.PlanSendTime != "" {
		sendTime = helper.GetDateToTimes(request.PlanSendTime)
	}
	taskName := "【" + time.Now().Format("01.02") + "】" + request.TaskName
	conditions := make(map[string]interface{})
	conditions["tag"] = []string{}
	// 根据模板判断是否走备用模板内容
	groupSendType := info.GroupSendType
	isBackupTemplate := microScrmModel.BackupTemplateCLose
	allowSelect := info.ALLowSelect
	templateContent := info.TemplateContent
	templateContentType := info.TemplateContentType
	finishRate := info.FinishRate
	if backUpInfo.ID > 0 {
		var conditionContent planReq.ConditionContent
		err := json.Unmarshal([]byte(backUpInfo.TemplateContent), &conditionContent)
		if err != nil {
			return err
		}
		if backUpInfo.EndTime == 0 {
			if backUpInfo.StartTime <= time.Now().Unix() {
				groupSendType = backUpInfo.GroupSendType
				isBackupTemplate = microScrmModel.BackupTemplateOpen
				allowSelect = backUpInfo.ALLowSelect
				finishRate = backUpInfo.FinishRate
				if conditionContent.Text.Content != "" || (backUpInfo.TemplateContentType > 1) {
					templateContent = backUpInfo.TemplateContent
					templateContentType = backUpInfo.TemplateContentType
				}
			}

		} else {
			if backUpInfo.StartTime <= time.Now().Unix() && time.Now().Unix() <= backUpInfo.EndTime {
				groupSendType = backUpInfo.GroupSendType
				isBackupTemplate = microScrmModel.BackupTemplateOpen
				allowSelect = backUpInfo.ALLowSelect
				finishRate = backUpInfo.FinishRate
				if conditionContent.Text.Content != "" || (backUpInfo.TemplateContentType > 1) {
					templateContent = backUpInfo.TemplateContent
					templateContentType = backUpInfo.TemplateContentType
				}
			}
		}
	}

	err := d.dao.WithContext(c).Transaction(func(tx *gorm.DB) error {
		taskData := microScrmModel.PlanTasksModel{
			BaseModel: models.BaseModel{
				ID: taskId,
			},
			PlanID:           request.TaskId,
			Name:             taskName,
			Dimension:        dimension,
			ConditionType:    1,
			IsSpecifyCrowd:   1,
			IsOpen:           1,
			TaskType:         int32(request.PlanType),
			SendType:         int32(request.PlanSendType),
			Status:           0,
			Remark:           request.TaskRemark,
			TaskIsPush:       0,
			BirthdayPreTime:  0,
			ConsumerFileID:   0,
			IsExtract:        0,
			CrowdID:          request.CrowdId,
			PlanSendTime:     sendTime,
			GroupSendType:    groupSendType,
			FinishRate:       finishRate,
			TaskStartTime:    helper.GetDateToTimes(request.TaskStartTime),
			TaskEndTime:      helper.GetDateToTimes(request.TaskEndTime),
			TaskExpireType:   1,
			PlanTemplateId:   info.ID,
			FlowId:           request.FlowId,
			Campaign:         request.Campaign,
			PlanName:         request.TaskName,
			OriginType:       microScrmModel.START_DT_TASK_ORIGIN_TYPE,
			IsBackupTemplate: isBackupTemplate,
			AllowSelect:      allowSelect,
			TaskCategory:     info.TaskCategory,
			TagExpireType:    tagExpireType,
			QyTagId:          wxTagId,
			TagName:          tagName,
		}
		err := tx.Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Create(&taskData).Error
		conditionData := microScrmModel.PlanTaskConditionsModel{
			BaseModel: models.BaseModel{
				ID: app.UID(),
			},
			PlanID:      request.TaskId,
			TaskID:      taskId,
			Conditions:  helper.MarshalToString(conditions),
			Content:     templateContent,
			ContentType: templateContentType,
		}
		err = tx.Scopes(d.dao.BranchTable(microScrmModel.PlanTaskConditionsModel{}.TableName(), enterpriseId)).Create(&conditionData).Error
		return err
	})
	if err != nil {
		return err
	}
	err = d.Increment(c, enterpriseId, map[string]interface{}{"id": request.TaskId}, "task_num", 1)
	if err != nil {
		return err
	}
	return nil
}

func (d *PlanTaskDao) CreateTask(c ctx.Context, request openReq.CreateTaskReq, enterpriseId, taskId int64, templateContent string, tagExpireType int, tagName, wxTagId string) error {
	dimension := uint32(microScrmModel.TaskDimensionImportSaUnionId)
	sendTime := helper.GetDateToTimes(request.TaskStartTime)
	finishRate := request.FinishRate
	if finishRate <= 0 {
		finishRate = 100
	}
	taskName := "[" + time.Now().Format("01.02") + "]" + request.TaskName
	conditions := make(map[string]interface{})
	conditions["tag"] = []string{}

	return d.dao.WithContext(c).Transaction(func(tx *gorm.DB) error {
		taskData := microScrmModel.PlanTasksModel{
			BaseModel: models.BaseModel{
				ID: taskId,
			},
			PlanID:           request.PlanId,
			Name:             taskName,
			Dimension:        dimension,
			ConditionType:    1,
			IsSpecifyCrowd:   1,
			IsOpen:           1,
			TaskType:         int32(microScrmModel.PLANTSINGLE),
			SendType:         int32(microScrmModel.SendTypeKeySa),
			Status:           0,
			Remark:           request.TaskRemark,
			TaskIsPush:       0,
			BirthdayPreTime:  0,
			ConsumerFileID:   0,
			IsExtract:        0,
			CrowdID:          request.CrowdId,
			PlanSendTime:     sendTime,
			GroupSendType:    request.GroupSendType,
			FinishRate:       finishRate,
			TaskStartTime:    helper.GetDateToTimes(request.TaskStartTime),
			TaskEndTime:      helper.GetDateToTimes(request.TaskEndTime),
			TaskExpireType:   1,
			PlanTemplateId:   0,
			FlowId:           request.FlowId,
			Campaign:         request.Campaign,
			PlanName:         request.TaskName,
			OriginType:       microScrmModel.SHENYAN_TASK_ORIGIN_TYPE,
			IsBackupTemplate: microScrmModel.BackupTemplateCLose,
			AllowSelect:      request.ALLowSelect,
			TaskCategory:     request.TaskCategory,
			TagExpireType:    tagExpireType,
			QyTagId:          wxTagId,
			TagName:          tagName,
		}
		if err := tx.Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Create(&taskData).Error; err != nil {
			return err
		}
		conditionData := microScrmModel.PlanTaskConditionsModel{
			BaseModel: models.BaseModel{
				ID: app.UID(),
			},
			PlanID:      request.PlanId,
			TaskID:      taskId,
			Conditions:  helper.MarshalToString(conditions),
			Content:     templateContent,
			ContentType: request.TemplateContentType,
		}
		if err := tx.Scopes(d.dao.BranchTable(microScrmModel.PlanTaskConditionsModel{}.TableName(), enterpriseId)).Create(&conditionData).Error; err != nil {
			return err
		}
		return d.Increment(c, enterpriseId, map[string]interface{}{"id": request.PlanId}, "task_num", 1)
	})
}

func (d *PlanTaskDao) Increment(c app.Ctx, enterpriseId int64, query interface{}, column string, num int64) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.PlansModel{}.TableName(), enterpriseId)).Where(query).Model(&microScrmModel.PlansModel{}).Update(column, gorm.Expr(fmt.Sprintf("%s + ?", column), num)).Error
	return err
}

// GetMaPlanPushLists 获取要推送的MA任务
func (d *PlanTaskDao) GetMaPlanPushLists(c ctx.Context, enterpriseId int64) ([]microScrmModel.PlanTasksModel, error) {
	var result []microScrmModel.PlanTasksModel
	planTable := d.dao.GetTable(microScrmModel.PlansModel{}.TableName(), enterpriseId)
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Joins("LEFT JOIN " + planTable + " as p ON p.id=plan_tasks.plan_id").Where(map[string]interface{}{
		"p.status":                microScrmModel.PlanIng,
		"p.category":              []int{microScrmModel.PlanMaTypeYes, microScrmModel.PlanMaTypeNewYes},
		"p.deleted_at":            0,
		"plan_tasks.task_is_push": 0,
		"plan_tasks.is_extract":   0,
		"plan_tasks.is_push":      microScrmModel.FINISHPUSHMA,
	}).Find(&result).Error
	return result, err
}

// GetMaPlanPushInfo  获取一条要推送的MA任务
func (d *PlanTaskDao) GetMaPlanPushInfo(c ctx.Context, enterpriseId int64) (microScrmModel.PlanTasksModel, error) {
	var nowUnix = time.Now().Unix()
	var result microScrmModel.PlanTasksModel
	planTable := d.dao.GetTable(microScrmModel.PlansModel{}.TableName(), enterpriseId)
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Joins("LEFT JOIN "+planTable+" as p ON p.id=plan_tasks.plan_id").Where(map[string]interface{}{
		"p.status":                microScrmModel.PlanIng,
		"p.category":              []int{microScrmModel.PlanMaTypeYes, microScrmModel.PlanMaTypeNewYes},
		"p.deleted_at":            0,
		"plan_tasks.task_is_push": 0,
		"plan_tasks.is_extract":   0,
		"plan_tasks.is_push":      microScrmModel.FINISHPUSHMA,
	}).Where("plan_tasks.plan_send_time <= ?", nowUnix).
		Where("plan_tasks.task_start_time <= ?", nowUnix).
		Where("plan_tasks.task_end_time >= ?", nowUnix).
		Order("plan_tasks.task_start_time asc").
		Find(&result).Error
	return result, err
}

// GetConditionInfo 获取condition信息
func (d *PlanTaskDao) GetConditionInfo(c ctx.Context, enterpriseId int64, query interface{}, args ...interface{}) (info microScrmModel.PlanTaskConditionsModel, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.PlanTaskConditionsModel{}.TableName(), enterpriseId)).Where(query, args...).Find(&info).Error
	return
}

// GetAllTasksByPlanId 根据planid获取所有的子任务信息
func (d *PlanTaskDao) GetAllTasksByPlanId(c ctx.Context, enterpriseId int64, planId int64) ([]microScrmModel.PlanTasksModel, error) {
	var result []microScrmModel.PlanTasksModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where("plan_id=?", planId).Find(&result).Error
	return result, err
}

// GetMaPlanEndLists 获取已结束的未开始推送MA任务结果回执的任务列表
func (d *PlanTaskDao) GetMaPlanEndLists(c ctx.Context, enterpriseId int64) ([]microScrmModel.PlanTasksModel, error) {
	var result []microScrmModel.PlanTasksModel
	planTable := d.dao.GetTable(microScrmModel.PlansModel{}.TableName(), enterpriseId)
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Joins("LEFT JOIN "+planTable+" as p ON p.id=plan_tasks.plan_id").Where(map[string]interface{}{
		"p.category":         []int{microScrmModel.PlanMaTypeYes, microScrmModel.PlanMaTypeNewYes},
		"p.deleted_at":       0,
		"plan_tasks.is_push": microScrmModel.FINISHPUSHMA, // MA的子任务数据已推送
	}).Where("plan_tasks.task_end_time<?", time.Now().Unix()).Find(&result).Error
	return result, err
}

// UpdateMaPlanFinishPush MA任务结束之后的完成回执推送状态
func (d *PlanTaskDao) UpdateMaPlanFinishPush(c ctx.Context, enterpriseId, taskId int64, isPush int) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(&microScrmModel.PlanTasksModel{}).Where("id=?", taskId).Updates(map[string]interface{}{
		"status":  microScrmModel.PlanEnd,
		"is_push": isPush,
	}).Error
}

func (d *PlanTaskDao) GetPlanCountByTemplateId(c ctx.Context, enterpriseId, templateId int64) (int64, error) {
	var total int64
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.PlansModel{}.TableName(), enterpriseId)).Where(map[string]interface{}{
		"status":           []int{0, 1},
		"plan_template_id": templateId,
		"deleted_at":       0,
	}).Count(&total).Error
	return total, err
}

// GetTaskTypePlanIdByPlanIds 根据获取持续任务还是单次任务
func (d *PlanTaskDao) GetTaskTypePlanIdByPlanIds(c ctx.Context, enterpriseId int64, planIds []int64) ([]int64, error) {
	var result []int64
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(map[string]interface{}{
		"plan_id":    planIds,
		"task_type":  2,
		"deleted_at": 0,
	}).Pluck("plan_id", &result).Group("planId").Error
	return result, err
}

// GetSpecifyCrowdUserIdsByTaskId 根据任务id获取指定人群的user_id
func (d *PlanTaskDao) GetSpecifyCrowdUserIdsByTaskId(c ctx.Context, enterpriseId, taskId int64) ([]string, error) {
	var result []string
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.PlanTaskSpecifyCrowdModel{}.TableName(), enterpriseId)).Where("task_id=?", taskId).Pluck("user_id", &result).Error
	return result, err
}

// GetQyGroupListsByPlanId  获取获取有效期内的群发任务类型的子任务列表
func (d *PlanTaskDao) GetQyGroupListsByPlanId(c ctx.Context, enterpriseId int64, groupSendType int, planIds []int64) ([]microScrmModel.PlanTasksModel, error) {
	var result []microScrmModel.PlanTasksModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(map[string]interface{}{
		"plan_id":         planIds,
		"group_send_type": groupSendType,
		"status":          microScrmModel.PlanIng,
		"deleted_at":      0,
	}).Find(&result).Error
	return result, err
}

// GetPlanIdByStatus 获取任务状态下的主任务id
func (d *PlanTaskDao) GetPlanIdByStatus(c ctx.Context, enterpriseId int64, status int) ([]int64, error) {
	var result []int64
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.PlansModel{}.TableName(), enterpriseId)).Where("status=? and deleted_at=0", status).Pluck("id", &result).Error
	return result, err
}

// CreateCopyTask 创建子任务
func (d *PlanTaskDao) CreateCopyTask(c app.Ctx, enterpriseId int64, taskData microScrmModel.PlanTasksModel, conditionData microScrmModel.PlanTaskConditionsModel) error {
	err := d.dao.WithContext(c).Transaction(func(tx *gorm.DB) error {
		err := tx.Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Create(&taskData).Error
		err = tx.Scopes(d.dao.BranchTable(microScrmModel.PlanTaskConditionsModel{}.TableName(), enterpriseId)).Create(&conditionData).Error
		err = tx.Scopes(d.dao.BranchTable(microScrmModel.PlansModel{}.TableName(), enterpriseId)).Where("id=?", taskData.PlanID).UpdateColumn("task_num", gorm.Expr(fmt.Sprintf("%s + ?", "task_num"), 1)).Error
		return err
	})
	return err
}

func (d *PlanTaskDao) GetMany(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]microScrmModel.PlanTasksModel, error) {
	var result []microScrmModel.PlanTasksModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *PlanTaskDao) GetStatLists(c ctx.Context, request planReq.NewStatDetailReq) (int64, []microScrmModel.PlanTasksModel, error) {
	var result []microScrmModel.PlanTasksModel
	var total int64
	whereMap := map[string]interface{}{
		"plan_id":    request.PlanId,
		"deleted_at": 0,
	}
	if request.TaskCategory >= 0 {
		whereMap["task_category"] = request.TaskCategory
	}
	if request.GroupSendType >= 0 {
		whereMap["group_send_type"] = request.GroupSendType
	}
	if request.Dimension > 0 {
		whereMap["dimension"] = request.Dimension
	}
	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, request.EnterpriseId)).Where(whereMap)

	if request.TaskName != "" {
		query.Where("name like ?", "%"+request.TaskName+"%")
	}
	offset := (request.Page - 1) * request.Size
	size := request.Size
	err := query.Count(&total).Offset(offset).Limit(size).Find(&result).Error
	return total, result, err
}

func (d *PlanTaskDao) GetCount(c ctx.Context, enterpriseId int64, planId, startTime, endTime int64) (int64, error) {
	var total int64
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where("`plan_id`= ?", planId).Where("`created_at` BETWEEN ? AND ? AND deleted_at=0 ", startTime, endTime).Count(&total).Error
	return total, err
}

func (d *PlanTaskDao) UpdatePlanData(c app.Ctx, enterpriseId int64, query interface{}, data map[string]interface{}) error {
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.PlansModel{}.TableName(), enterpriseId)).Where(query).Model(&microScrmModel.PlansModel{}).Updates(data).Error
	return err
}

// GetPlanMany 获取多条主任务
func (d *PlanTaskDao) GetPlanMany(c app.Ctx, enterpriseId int64, query interface{}, args ...interface{}) ([]microScrmModel.PlansModel, error) {
	var result []microScrmModel.PlansModel
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.PlansModel{}.TableName(), enterpriseId)).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *PlanTaskDao) UpdateTaskStatus(c ctx.Context, enterpriseId int64, status int, query interface{}, args ...interface{}) error {
	return d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Where(query, args...).Model(&microScrmModel.PlanTasksModel{}).Update("status", status).Error

}

func (d *PlanTaskDao) GetPlanNameList(c app.Ctx, enterpriseId, planId int64) ([]string, error) {
	var result []string
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).Model(&microScrmModel.PlanTasksModel{}).Where("plan_id=?", planId).Group("plan_name").Pluck("plan_name", &result).Error
	return result, err
}

func (d *PlanTaskDao) GetTaskList(c ctx.Context, enterpriseId int64, allotIds []int64, req planReq.CustomerTaskReq) ([]planResp.CustomerTaskResp, int64, error) {
	offset := (req.Page - 1) * req.Size
	var list []planResp.CustomerTaskResp
	var total int64
	allotLogTable := d.dao.GetTable(microScrmModel.TaskAllotLogModel{}.TableName(), req.EnterpriseId)
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Joins("LEFT JOIN "+allotLogTable+" as a ON a.task_id=plan_tasks.id").Where("a.id in ?", allotIds).Scopes(func(db *gorm.DB) *gorm.DB {
		if req.TaskCategory >= 0 {
			db.Where("plan_tasks.task_category = ?", req.TaskCategory)
		}
		if req.Status >= 0 {
			db.Where("a.is_finished=?", req.Status)
		}
		return db
	}).Select("a.id,a.task_id,a.is_finished as status,plan_tasks.task_category,plan_tasks.name,plan_tasks.remark,plan_tasks.task_start_time,plan_tasks.task_end_time").Order("a.is_finished asc,plan_tasks.task_end_time asc").Count(&total).Offset(offset).
		Limit(req.Size).Scan(&list).Error
	return list, total, err
}

// GetCustomerTaskNums 顾客详情中，导购和该顾客的待完成任务数
func (d *PlanTaskDao) GetCustomerTaskNums(c ctx.Context, enterpriseId int64, allotIds []int64) (int64, error) {
	var total int64
	allotLogTable := d.dao.GetTable(microScrmModel.TaskAllotLogModel{}.TableName(), enterpriseId)
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Joins("LEFT JOIN "+allotLogTable+" as a ON a.task_id=plan_tasks.id").Where("a.id in ?", allotIds).Scopes(func(db *gorm.DB) *gorm.DB {
		db.Where("a.is_finished=?", microScrmModel.PLAN_FISISHING)
		return db
	}).Count(&total).Error
	return total, err
}

func (d *PlanTaskDao) GetNoFinishTaskNums(c ctx.Context, enterpriseId int64, userId string) ([]int64, error) {
	lastTime := time.Now().Unix() - 86400*30*3 // 只展示最近三个月的任务
	var taskIds []int64
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.TaskAllotLogModel{}.TableName(), enterpriseId)).
		Where("user_id = ?", userId).
		Where("created_at>=?", lastTime).
		Where("is_finished=?", microScrmModel.PLAN_FISISHING).Pluck("task_id", &taskIds).Error
	return taskIds, err
}

// GetTemplateBackUpInfo 获取备份模板的信息
func (d *PlanTaskDao) GetTemplateBackUpInfo(c ctx.Context, enterpriseId, templateId int64) (result microScrmModel.PlanTaskTemplateBackupModel, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(microScrmModel.PlanTaskTemplateBackupModel{}.TableName(), enterpriseId)).Where(map[string]interface{}{
		"template_id": templateId,
	}).Find(&result).Error
	return
}
