package customerSvc

import (
	"errors"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	customerDao "micro-scrm/internal/app/dao/customer"
	membersDao "micro-scrm/internal/app/dao/members"
	customerReq "micro-scrm/internal/app/http/request/customer"
	customerResp "micro-scrm/internal/app/http/response/customer"
	"micro-scrm/internal/app/models"
	customerModels "micro-scrm/internal/app/models/customer"
	"micro-scrm/internal/app/service"
	"time"
)

type SowTemplateService struct {
	service.Base
	sowTemplateDao       *customerDao.SowTemplate
	sowTemplateNoticeDao *customerDao.SowTemplateNotice
}

func NewSowTemplateService() *SowTemplateService {
	return &SowTemplateService{
		sowTemplateDao:       customerDao.NewSowTemplate(),
		sowTemplateNoticeDao: customerDao.NewSowTemplateNotice(),
	}
}

func (s *SowTemplateService) Save(c app.Ctx, params customerReq.SowTemplateSaveReq) error {
	templateId := params.ID
	memberInfo, err := membersDao.NewMembers().GetMemberByMobile(c, params.EnterpriseID, params.Mobile)
	if err != nil {
		return err
	}
	nameCheck, _ := s.sowTemplateDao.GetOne(c, params.EnterpriseID, "name=?", params.Name)
	if templateId > 0 {
		info, err := s.sowTemplateDao.GetOne(c, params.EnterpriseID, "id=?", params.ID)
		if err != nil || info.ID == 0 {
			return errors.New("找不到该数据")
		}
		if nameCheck.ID > 0 && nameCheck.ID != info.ID {
			return errors.New("该模版名称已存在")
		}
		err = s.sowTemplateDao.UpdateData(c, map[string]interface{}{
			"id": templateId,
		}, map[string]interface{}{
			"name":       params.Name,
			"theme":      params.Theme,
			"status":     params.Status,
			"updated_at": time.Now().Unix(),
		})
		if err != nil {
			return err
		}
		// 删除模版提醒数据
		err = s.sowTemplateNoticeDao.UnscopedDeleteData(c, map[string]interface{}{
			"template_id": templateId,
		})
		if err != nil {
			return err
		}
	} else {
		if nameCheck.ID > 0 {
			return errors.New("该模版名称已存在")
		}
		templateId = app.UID()
		err = s.sowTemplateDao.CreateData(c, customerModels.SowTemplateModel{
			MicroBaseModel: models.CreateMicroBase(params.EnterpriseID, templateId),
			Name:           params.Name,
			Theme:          params.Theme,
			Status:         params.Status,
			AdminName:      memberInfo.Name,
		})
		if err != nil {
			return err
		}
	}
	// 重新写入模版提醒数据列表
	var noticeData []customerModels.SowTemplateNoticeModel
	for _, notice := range params.SowTemplateNotice {
		tmpNotice := customerModels.SowTemplateNoticeModel{
			MicroBaseModel: models.CreateMicroBase(params.EnterpriseID, app.UID()),
			TemplateID:     templateId,
			DelayType:      notice.DelayType,
			DelayValue:     notice.DelayValue,
			DelayTime:      notice.DelayTime,
			Content:        notice.Content,
		}
		noticeData = append(noticeData, tmpNotice)
	}
	err = s.sowTemplateNoticeDao.CreateInBatchesData(c, noticeData)
	if err != nil {
		return err
	}
	return nil
}

func (s *SowTemplateService) List(c ctx.Context, params customerReq.SowTemplateListReq) (resList []customerResp.SowTemplateInfo, total int64, err error) {
	list, total, err := s.sowTemplateDao.GetPagination(c, params)
	if err != nil {
		return
	}
	if len(list) == 0 {
		return
	}
	var templateIds []int64
	for _, val := range list {
		templateIds = append(templateIds, val.ID)
	}
	memberIdNum, err := customerDao.NewSowTask().GetMemberIdNumByTemplateIds(c, params.EnterpriseID, templateIds)
	if err != nil {
		return
	}
	for _, val := range list {
		tmpInfo := customerResp.SowTemplateInfo{
			Id:        val.ID,
			Name:      val.Name,
			TaskNum:   len(val.NoticeList),
			Status:    val.Status,
			UseNum:    val.UseNum,
			AdminName: val.AdminName,
			CreatedAt: val.CreatedAt,
			UpdatedAt: val.UpdatedAt,
		}
		for _, v := range memberIdNum {
			if v.TemplateId == val.ID {
				tmpInfo.UseSaNum = v.SaNum
			}
		}
		resList = append(resList, tmpInfo)
	}

	return
}

func (s *SowTemplateService) Detail(c ctx.Context, params customerReq.SowTemplateDetailReq) (detail customerModels.SowTemplateModel, err error) {
	detail, err = s.sowTemplateDao.GetDetail(c, params.EnterpriseID, params.Id)
	if err != nil {
		return
	}
	return
}

func (s *SowTemplateService) UpdateStatus(c ctx.Context, params customerReq.SowTemplateDetailReq) (err error) {
	err = s.sowTemplateDao.UpdateData(c, map[string]interface{}{
		"id":            params.Id,
		"enterprise_id": params.EnterpriseID,
	}, map[string]interface{}{
		"status":     params.Status,
		"updated_at": time.Now().Unix(),
	})
	return
}

func (s *SowTemplateService) Delete(c ctx.Context, params customerReq.SowTemplateDetailReq) error {
	// 使用事务
	err := s.Base.DB().Transaction(func(tx *gorm.DB) error {
		err := customerDao.NewSowTemplate().WithTx(tx).DeleteData(c, "enterprise_id=? AND id=?", params.EnterpriseID, params.Id)
		if err != nil {
			return err
		}
		err = customerDao.NewSowTemplateNotice().WithTx(tx).DeleteData(c, "template_id=?", params.Id)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
