package marketingService

import (
	"encoding/json"
	"micro-scrm/internal/app"
	templateDao "micro-scrm/internal/app/dao/template"
	templateReq "micro-scrm/internal/app/http/request/template"
	"micro-scrm/internal/app/models"
	templateModels "micro-scrm/internal/app/models/template"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/uid"
	"time"
)

type CustomCard struct {
	templateDao             *templateDao.Template
	templatePlanRelationDao *templateDao.TemplatePlanRelation
}

func NewCustomCard() *CustomCard {
	return &CustomCard{
		templateDao:             templateDao.NewTemplate(),
		templatePlanRelationDao: templateDao.NewTemplatePlanRelation(),
	}
}

func (c CustomCard) GetTemplateList(ctx app.Ctx, req templateReq.ReqList) (data []templateReq.TemplateInfo, err error) {
	offset := (req.Page - 1) * req.Size
	list, err := c.templateDao.GetLists(ctx, req.TemplateName, req.EnterpriseID, offset, req.Size)
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		for _, item := range list {
			var templateInfo templateReq.TemplateInfo
			templateInfo.ID = item.ID
			templateInfo.TemplateName = item.TemplateName
			templateInfo.TemplateContent = item.TemplateContent
			templateInfo.UpdatedAt = item.UpdatedAt
			if errJson := json.Unmarshal([]byte(item.TemplatePictureURL), &templateInfo.TemplatePictureURL); errJson != nil {
				continue
			}
			data = append(data, templateInfo)
		}
	}
	return data, err
}

func (c CustomCard) GetCount(ctx app.Ctx, req templateReq.ReqList) (int64, error) {
	return c.templateDao.GetCount(ctx, req.TemplateName, req.EnterpriseID)
}

func (c CustomCard) Create(ctx app.Ctx, info templateReq.ReqCreate) error {
	data := templateModels.TemplateModel{
		TemplateName:       info.TemplateName,
		TemplatePictureURL: helper.MarshalToString(info.TemplatePictureURL),
		TemplateContent:    info.TemplateContent,
		MicroBaseModel:     models.CreateMicroBase(info.EnterpriseID, 0),
	}
	return c.templateDao.CreateInfo(ctx, data)
}

func (c CustomCard) UpdateInfo(ctx app.Ctx, info templateReq.ReqCreate) error {
	data := templateModels.TemplateModel{
		TemplateName:       info.TemplateName,
		TemplatePictureURL: helper.MarshalToString(info.TemplatePictureURL),
		TemplateContent:    info.TemplateContent,
	}
	return c.templateDao.UpdateInfo(ctx, data, info.ID)
}

func (c CustomCard) Delete(ctx app.Ctx, Req templateReq.RequestInfo) error {
	nowTime := time.Now().Unix()
	tx := app.DB().WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}
	err := c.templateDao.WithTx(tx).UpdateData(ctx,
		map[string]interface{}{"id": Req.ID},
		map[string]interface{}{"deleted_at": nowTime},
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = c.templatePlanRelationDao.WithTx(tx).UpdateData(ctx,
		map[string]interface{}{"enterprise_id": Req.EnterpriseID, "template_id": Req.ID},
		map[string]interface{}{"deleted_at": nowTime},
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (c CustomCard) Check(ctx app.Ctx, Req templateReq.RequestInfo) (Total int64, err error) {
	return c.templatePlanRelationDao.Check(ctx, Req.EnterpriseID, Req.ID)
}

func (c CustomCard) Detail(ctx app.Ctx, Req templateReq.RequestInfo) (data templateReq.TemplateInfo, err error) {
	Info, err := c.templateDao.Detail(ctx, Req.ID)
	if err != nil {
		return templateReq.TemplateInfo{}, err
	}
	data.ID = Info.ID
	data.TemplateName = Info.TemplateName
	data.TemplateContent = Info.TemplateContent
	if errJson := json.Unmarshal([]byte(Info.TemplatePictureURL), &data.TemplatePictureURL); errJson != nil {
		return templateReq.TemplateInfo{}, errJson
	}
	data.UpdatedAt = Info.UpdatedAt
	return data, nil
}

func (c CustomCard) FindOrCreateUseCard(ctx app.Ctx, enterpriseID int64, planId int64, cardIds []int) {
	for i := 0; i < len(cardIds); i++ {
		data := templateModels.TemplatePlanRelationModel{
			PlanID:         planId,
			PlanStatus:     0,
			TemplateID:     cardIds[i],
			MicroBaseModel: models.CreateMicroBase(enterpriseID, uid.ID(ctx)),
		}
		err := c.templatePlanRelationDao.FindOrCreateUseCard(ctx, enterpriseID, planId, cardIds[i], data)
		if err != nil {
			app.Logger().Info(ctx).Msgf("子任务和贺卡的关联创建失败: %v,失败的任务id: %d,失败的贺卡模板id: %d", err, planId, cardIds[i])
			return
		}
	}
}

func (c CustomCard) UpdatePlanCardInfo(ctx app.Ctx, enterpriseID int64, planIds []int64) error {
	return c.templatePlanRelationDao.UpdatePlanCardInfo(ctx, enterpriseID, planIds)
}
