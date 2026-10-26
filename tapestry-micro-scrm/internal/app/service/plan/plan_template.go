package planService

import (
	"errors"
	"micro-scrm/internal/app"
	microBossDao "micro-scrm/internal/app/dao/micro_boss"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	openReq "micro-scrm/internal/app/http/request/open"
	planReq "micro-scrm/internal/app/http/request/plan"
	"micro-scrm/internal/app/models"
	microBossModel "micro-scrm/internal/app/models/micro_boss"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/helper"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

type PlanTemplateSvc struct {
	planTemplateDao *microScrmDao.PlanTemplateDao
}

func NewPlanTemplateSvc() *PlanTemplateSvc {
	return &PlanTemplateSvc{
		planTemplateDao: microScrmDao.NewPlanTemplateDao(),
	}
}

func (s PlanTemplateSvc) Create(c ctx.Context, request planReq.TemplateAddReq) error {
	if request.TaskCategory == microScrmModel.TaskCategoryMoment {
		request.GroupSendType = microScrmModel.GroupSendTypeGroupMsg
		request.ALLowSelect = 0
		request.FinishRate = 100
	}
	templateId := app.UID()
	data := microScrmModel.PlanTaskTemplateModel{
		BaseModel:           models.BaseModel{ID: templateId},
		TemplateName:        request.TemplateName,
		TemplateDesc:        request.TemplateDesc,
		TemplateContent:     helper.MarshalToString(request.TemplateContent),
		TemplateContentType: request.TemplateContentType,
		AdminID:             request.AdminId,
		GroupSendType:       request.GroupSendType,
		FinishRate:          request.FinishRate,
		ALLowSelect:         request.ALLowSelect,
		OpenBackUp:          request.OpenBackUp,
		TaskCategory:        request.TaskCategory,
	}
	err := s.planTemplateDao.Create(c, request.EnterpriseId, data)
	if err != nil {
		return err
	}
	if request.OpenBackUp == microScrmModel.OpenBackYes {
		templateBackData := microScrmModel.PlanTaskTemplateBackupModel{
			BaseModel:           models.BaseModel{ID: app.UID()},
			TemplateID:          templateId,
			TemplateContent:     helper.MarshalToString(request.TemplateBack.TemplateContent),
			TemplateContentType: request.TemplateBack.TemplateContentType,
			StartTime:           request.TemplateBack.StartTime,
			EndTime:             request.TemplateBack.EndTime,
			GroupSendType:       request.TemplateBack.GroupSendType,
			ALLowSelect:         request.TemplateBack.ALLowSelect,
			FinishRate:          request.TemplateBack.FinishRate,
		}
		err = s.planTemplateDao.CreateBack(c, request.EnterpriseId, templateBackData)
	}
	return err
}

func (s PlanTemplateSvc) Update(c ctx.Context, request planReq.TemplateUpdateReq) error {
	info, _ := s.planTemplateDao.GetTemplateInfo(c, request.EnterpriseId, request.Id)
	if info.ID == 0 {
		return errors.New("模版不存在")
	}
	data := map[string]any{
		"template_name":         request.TemplateName,
		"template_desc":         request.TemplateDesc,
		"template_content":      helper.MarshalToString(request.TemplateContent),
		"template_content_type": request.TemplateContentType,
		"updated_at":            time.Now().Unix(),
		"open_backup":           request.OpenBackUp,
		"task_category":         request.TaskCategory,
	}
	if request.TaskCategory == 0 {
		//沟通任务需要接这些参数，朋友圈不需要
		data["group_send_type"] = request.GroupSendType
		data["allow_select"] = request.ALLowSelect
		data["finish_rate"] = request.FinishRate
	} else if request.TaskCategory == microScrmModel.TaskCategoryMoment {
		data["group_send_type"] = microScrmModel.GroupSendTypeGroupMsg
		data["allow_select"] = 0
		data["finish_rate"] = 100
	}

	backUpInfo, err := s.planTemplateDao.GetTemplateBackUpInfo(c, request.EnterpriseId, request.Id)
	if err != nil {
		return err
	}
	if request.OpenBackUp != microScrmModel.OpenBackYes {
		if backUpInfo.ID > 0 {
			err = s.planTemplateDao.DeleteBack(c, request.EnterpriseId, backUpInfo.ID)
		}
	} else {
		if backUpInfo.ID > 0 {
			updateMap := map[string]any{
				"template_content":      helper.MarshalToString(request.TemplateBack.TemplateContent),
				"template_content_type": request.TemplateBack.TemplateContentType,
				"updated_at":            time.Now().Unix(),
				"start_time":            request.TemplateBack.StartTime,
				"end_time":              request.TemplateBack.EndTime,
			}
			if request.TaskCategory == 0 {
				//沟通任务需要接这些参数，朋友圈不需要
				updateMap["group_send_type"] = request.TemplateBack.GroupSendType
				updateMap["allow_select"] = request.TemplateBack.ALLowSelect
				updateMap["finish_rate"] = request.TemplateBack.FinishRate
			} else if request.TaskCategory == microScrmModel.TaskCategoryMoment {
				updateMap["group_send_type"] = microScrmModel.GroupSendTypeGroupMsg
				updateMap["allow_select"] = 0
				updateMap["finish_rate"] = 100
			}
			err = s.planTemplateDao.UpdateBack(c, request.EnterpriseId, updateMap, map[string]any{"id": backUpInfo.ID})
		} else {
			if info.TaskCategory == microScrmModel.TaskCategoryMoment {
				request.TemplateBack.GroupSendType = microScrmModel.GroupSendTypeGroupMsg
				request.TemplateBack.ALLowSelect = 0
				request.TemplateBack.FinishRate = 100
			}
			templateBackData := microScrmModel.PlanTaskTemplateBackupModel{
				BaseModel:           models.BaseModel{ID: app.UID()},
				TemplateID:          request.Id,
				TemplateContent:     helper.MarshalToString(request.TemplateBack.TemplateContent),
				TemplateContentType: request.TemplateBack.TemplateContentType,
				StartTime:           request.TemplateBack.StartTime,
				EndTime:             request.TemplateBack.EndTime,
				GroupSendType:       request.TemplateBack.GroupSendType,
				ALLowSelect:         request.TemplateBack.ALLowSelect,
				FinishRate:          request.TemplateBack.FinishRate,
			}
			err = s.planTemplateDao.CreateBack(c, request.EnterpriseId, templateBackData)
		}
	}
	err = s.planTemplateDao.UpdateData(c, request.EnterpriseId, map[string]any{"id": request.Id}, data)

	return err
}

func (s PlanTemplateSvc) Lists(c ctx.Context, request planReq.TemplateListsReq) ([]microScrmModel.PlanTaskTemplateModel, int64, error) {
	var adminId int64
	if request.AdminName != "" {
		enterpriseAdmin, _ := microBossDao.NewEnterpriseAdmin().GetInfo(c, microBossModel.EnterpriseAdminModel{
			MicroBaseModel: models.MicroBaseModel{EnterpriseID: request.EnterpriseId},
			Username:       request.AdminName})
		if enterpriseAdmin.ID > 0 {
			adminId = enterpriseAdmin.ID
		}
	}
	result, total, err := s.planTemplateDao.Lists(c, request, adminId)

	if err != nil {
		return nil, 0, err
	}

	if total > 0 {
		var adminIds []int64
		for _, res := range result {
			adminIds = append(adminIds, res.AdminID)
		}
		admins, err := microBossDao.NewEnterpriseAdmin().GetMany(c, map[string]any{"id": adminIds, "enterprise_id": request.EnterpriseId})
		if err != nil {
			return nil, 0, err
		}

		adminNameMap := make(map[int64]string)
		for _, admin := range admins {
			adminNameMap[admin.ID] = admin.Username
		}

		for k, res := range result {
			if res.AdminID > 0 {
				result[k].AdminName = adminNameMap[res.AdminID]
			}
		}
	}

	return result, total, err
}

func (s PlanTemplateSvc) Details(c ctx.Context, req planReq.TemplateDetails) (planReq.TemplateDetailsRes, error) {
	var result planReq.TemplateDetailsRes
	templateInfo, err := s.planTemplateDao.GetTemplateInfo(c, req.EnterpriseId, req.Id)
	if err != nil {
		return result, err
	}
	result.TemplateName = templateInfo.TemplateName
	result.TemplateDesc = templateInfo.TemplateDesc
	result.TemplateContent = templateInfo.TemplateContent
	result.TemplateContentType = templateInfo.TemplateContentType
	result.OpenBackUp = templateInfo.OpenBackUp
	result.Id = templateInfo.ID
	result.AdminId = templateInfo.AdminID
	result.GroupSendType = templateInfo.GroupSendType
	result.ALLowSelect = templateInfo.ALLowSelect
	result.Status = templateInfo.Status
	result.FinishRate = templateInfo.FinishRate
	result.TaskCategory = templateInfo.TaskCategory
	if templateInfo.OpenBackUp == microScrmModel.OpenBackYes {
		templateBackInfo, err := s.planTemplateDao.GetTemplateBackUpInfo(c, req.EnterpriseId, req.Id)
		if err != nil {
			return result, err
		}
		result.TemplateBack.TemplateContent = templateBackInfo.TemplateContent
		result.TemplateBack.TemplateContentType = templateBackInfo.TemplateContentType
		result.TemplateBack.ALLowSelect = templateBackInfo.ALLowSelect
		result.TemplateBack.GroupSendType = templateBackInfo.GroupSendType
		result.TemplateBack.StartTime = templateBackInfo.StartTime
		result.TemplateBack.EndTime = templateBackInfo.EndTime
		result.TemplateBack.FinishRate = templateBackInfo.FinishRate
	}
	return result, err
}

func (s PlanTemplateSvc) MaTemplateLists(c ctx.Context, request openReq.PlanTemplateListsReq) ([]microScrmModel.MaTemplateLists, error) {
	return s.planTemplateDao.MaTemplateLists(c, request)
}
