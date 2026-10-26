package sms

import (
	"bytes"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/pkg/errs"
	"micro-scrm/internal/pkg/sdk/boss"
	"micro-scrm/internal/pkg/sdk/microauthority"
	"micro-scrm/internal/pkg/sdk/sendcloud"
	"micro-scrm/internal/pkg/uid"
)

// Service 应用服务
type Service struct {
	Repo      Repository
	Marketing MarketingInterface
}

// StatisticInterface 统计接口
type MarketingInterface interface {
	GetSmsCount(c ctx.Context, smsId int64) (int64, error)
	CloseMarketingInnerStatus(c ctx.Context, EnterpriseId int64) error
}

// 创建短信模板
func (s *Service) CreateSmsTemplate(c ctx.Context, request CreateTemplateRequest) (*CreateTemplateResponse, error) {

	smsSignInfoRequest := SmsSignInfoRequest{}
	smsSignInfoRequest.Id = request.TemplateSignId

	// 查询短信签名Id
	smsSignInfo, smsSignInfoErr := s.Repo.SmsSignFindOne(c, smsSignInfoRequest)
	if smsSignInfoErr != nil {
		return nil, smsSignInfoErr
	}
	// 判断模板是否审核通过
	if smsSignInfo.ThirdSignStatus != ThirdSignStatusSuccessful {
		return nil, errs.SmsSignNotReviewed
	}

	var newTime = time.Now().Unix()
	// 添加短信模板
	var smsTemplateRequest = SmsTemplate{
		Id:              uid.ID(ctx.New()),
		EnterpriseId:    request.EnterpriseId,
		AdminId:         request.AdminId,
		TemplateName:    request.TemplateName,
		TemplateSignId:  smsSignInfo.Id,
		MessageTemplate: request.MessageTemplate,
		TemplateContent: request.TemplateContent,
		TemplateHtml:    request.TemplateHtml,
		MessageStr:      request.MessageStr,
		//Status:          TemplateStatusWaitVerify,
		Status:    TemplateStatusHasValue,
		CreatedAt: newTime,
	}

	// 创建模板验证模板名称是否重复
	SmsInfo, SmsInfoErr := s.CreateSmsTemplateCheck(ctx.Wrap(c), smsTemplateRequest)
	if SmsInfoErr != nil {
		return nil, SmsInfoErr
	}
	// 名称已重复
	if SmsInfo != nil && SmsInfo.Id != 0 {
		return nil, errs.SmsTemplateNameRepeat
	}
	//// 请求第三方创建短信模板
	//var createSmsTemplateRequest = sendcloud.CreateSmsTemplatequest{
	//	SmsUser:         sendcloud.SendCloudSmsUser,
	//	TemplateName:    request.TemplateName + "-" + strconv.FormatInt(request.EnterpriseId, 10),
	//	TemplateText:    checkMessageContent(request.MessageTemplate),
	//	SignName:        smsSignInfo.SignName,
	//	SignId:          smsSignInfo.ThirdSignId,
	//	SignPositionStr: "0",
	//	SmsTypeStr:      "1",
	//}
	//
	//// 请求第三方创建签名
	//signString := GetSign(createSmsTemplateRequest, sendcloud.SendCloudSmsKey)
	//createSmsTemplateRequest.Signature = signString
	//createSmsTemplateRes, createSmsTemplateErr := sendcloud.CreateSmsTemplate(c, createSmsTemplateRequest)
	//
	//if createSmsTemplateErr != nil {
	//	return nil, createSmsTemplateErr
	//}

	// 将提交成功的短信模板存入本地数据库
	//smsTemplateRequest.TemplateId = createSmsTemplateRes.Info.TemplateId
	// 创建模板验证模板名称是否重复
	_, CreateSmsTemplateInfoErr := s.Repo.CreateSmsTemplate(c, smsTemplateRequest)
	if CreateSmsTemplateInfoErr != nil {
		return nil, CreateSmsTemplateInfoErr
	}

	var CreateTemplateResponse = CreateTemplateResponse{
		TemplateId: smsTemplateRequest.Id,
	}

	// 判断是否提审
	//if request.IsExamine == 1 {
	//	var smsTemplateInfoRequest = SmsTemplateArraignmentRequest{
	//		Id: smsTemplateRequest.Id,
	//	}
	//	s.SmsTemplateArraignment(c, smsTemplateInfoRequest)
	//}
	return &CreateTemplateResponse, nil
}

// 创建短信模板
func (s *Service) CreateSmsTemplateBak(c ctx.Context, request CreateTemplateRequest) (*CreateTemplateResponse, error) {

	smsSignInfoRequest := SmsSignInfoRequest{}
	smsSignInfoRequest.Id = request.TemplateSignId

	// 查询短信签名Id
	smsSignInfo, smsSignInfoErr := s.Repo.SmsSignFindOne(c, smsSignInfoRequest)
	if smsSignInfoErr != nil {
		return nil, smsSignInfoErr
	}
	// 判断模板是否审核通过
	if smsSignInfo.ThirdSignStatus != ThirdSignStatusSuccessful {
		return nil, errs.SmsSignNotReviewed
	}

	var newTime = time.Now().Unix()
	// 添加短信模板
	var smsTemplateRequest = SmsTemplate{
		Id:              uid.ID(ctx.New()),
		EnterpriseId:    request.EnterpriseId,
		AdminId:         request.AdminId,
		TemplateName:    request.TemplateName,
		TemplateSignId:  smsSignInfo.Id,
		MessageTemplate: request.MessageTemplate,
		TemplateContent: request.TemplateContent,
		TemplateHtml:    request.TemplateHtml,
		MessageStr:      request.MessageStr,
		Status:          TemplateStatusWaitVerify,
		CreatedAt:       newTime,
	}

	// 创建模板验证模板名称是否重复
	SmsInfo, SmsInfoErr := s.CreateSmsTemplateCheck(ctx.Wrap(c), smsTemplateRequest)
	if SmsInfoErr != nil {
		return nil, SmsInfoErr
	}
	// 名称已重复
	if SmsInfo != nil && SmsInfo.Id != 0 {
		return nil, errs.SmsTemplateNameRepeat
	}
	// 请求第三方创建短信模板
	var createSmsTemplateRequest = sendcloud.CreateSmsTemplatequest{
		SmsUser:         sendcloud.SendCloudSmsUser,
		TemplateName:    request.TemplateName + "-" + strconv.FormatInt(request.EnterpriseId, 10),
		TemplateText:    checkMessageContent(request.MessageTemplate),
		SignName:        smsSignInfo.SignName,
		SignId:          smsSignInfo.ThirdSignId,
		SignPositionStr: "0",
		SmsTypeStr:      "1",
	}

	// 请求第三方创建签名
	signString := GetSign(createSmsTemplateRequest, sendcloud.SendCloudSmsKey)
	createSmsTemplateRequest.Signature = signString
	createSmsTemplateRes, createSmsTemplateErr := sendcloud.CreateSmsTemplate(c, createSmsTemplateRequest)

	if createSmsTemplateErr != nil {
		return nil, createSmsTemplateErr
	}

	// 将提交成功的短信模板存入本地数据库
	smsTemplateRequest.TemplateId = createSmsTemplateRes.Info.TemplateId
	// 创建模板验证模板名称是否重复
	_, CreateSmsTemplateInfoErr := s.Repo.CreateSmsTemplate(c, smsTemplateRequest)
	if CreateSmsTemplateInfoErr != nil {
		return nil, CreateSmsTemplateInfoErr
	}

	var CreateTemplateResponse = CreateTemplateResponse{
		TemplateId: smsTemplateRequest.Id,
	}

	// 判断是否提审
	if request.IsExamine == 1 {
		var smsTemplateInfoRequest = SmsTemplateArraignmentRequest{
			Id: smsTemplateRequest.Id,
		}
		s.SmsTemplateArraignment(c, smsTemplateInfoRequest)
	}
	return &CreateTemplateResponse, nil
}

// 编辑更新短信模板
func (s *Service) UpdateSmsTemplate(c ctx.Context, request UpdateSmsTemplateStatus) (bool, error) {

	// 获取第三方模板信息
	var smsTemplateInfoRequest = SmsTemplate{
		Id: request.Id,
	}
	smsTemplateInfo, smsSignInfoErr := s.Repo.SmsTemplateFindOne(c, smsTemplateInfoRequest)
	if smsSignInfoErr != nil {
		return false, smsSignInfoErr
	}

	smsSignInfoRequest := SmsSignInfoRequest{}
	smsSignInfoRequest.Id = request.TemplateSignId

	// 查询短信签名Id
	smsSignInfo, smsSignInfoErr := s.Repo.SmsSignFindOne(c, smsSignInfoRequest)
	if smsSignInfoErr != nil {
		return false, smsSignInfoErr
	}

	// 判断模板是否审核通过
	if smsSignInfo.ThirdSignStatus != ThirdSignStatusSuccessful {
		return false, nil
	}

	var thirdTemplateId = strconv.FormatInt(smsTemplateInfo.TemplateId, 10)
	// 请求第三方更新短信模板
	var updateSmsTemplateRequest = sendcloud.UpdateSmsTemplatequest{
		SmsUser:         sendcloud.SendCloudSmsUser,
		TemplateIdStr:   thirdTemplateId,
		TemplateName:    request.TemplateName + "-" + strconv.FormatInt(request.EnterpriseId, 10),
		TemplateText:    checkMessageContent(request.MessageTemplate),
		SignName:        smsSignInfo.SignName,
		SignId:          smsSignInfo.ThirdSignId,
		SignPositionStr: "0",
		SmsTypeStr:      "1",
	}

	// 请求第三方创建签名
	signString := GetSign(updateSmsTemplateRequest, sendcloud.SendCloudSmsKey)
	updateSmsTemplateRequest.Signature = signString
	_, updateSmsTemplateErr := sendcloud.UpdateSmsTemplate(c, updateSmsTemplateRequest)

	if updateSmsTemplateErr != nil {
		return false, updateSmsTemplateErr
	}

	// 修改短信模板状态为未审核
	var TemplateStatus = UpdateSmsTemplateStatus{
		Id:     request.Id,
		Status: TemplateStatusWaitVerify,
	}
	TemplateDeleteErr := s.Repo.SmsTemplateEdit(c, TemplateStatus)

	if TemplateDeleteErr != nil {
		return true, TemplateDeleteErr
	}

	return true, nil
}

// SmsTemplateArraignment 模板提审
func (s *Service) SmsTemplateArraignment(c ctx.Context, request SmsTemplateArraignmentRequest) (*SmsTemplateArraignmentReponse, error) {
	// 获取第三方模板信息
	var smsTemplateInfoRequest = SmsTemplate{
		Id: request.Id,
	}
	smsTemplateInfo, smsSignInfoErr := s.Repo.SmsTemplateFindOne(c, smsTemplateInfoRequest)
	if smsSignInfoErr != nil {
		return nil, smsSignInfoErr
	}

	if smsTemplateInfo.Id == 0 {
		return nil, errs.SmsTemplateNonExistent
	}

	if len(smsTemplateInfo.MessageTemplate) == 0 {
		return nil, errs.TemplateContentFormat
	}

	// 获取第三方短信签名
	smsSignInfoRequest := SmsSignInfoRequest{}
	smsSignInfoRequest.Id = smsTemplateInfo.TemplateSignId

	// 查询短信签名Id
	smsSignInfo, smsSignInfoErr := s.Repo.SmsSignFindOne(c, smsSignInfoRequest)
	if smsSignInfoErr != nil {
		return nil, smsSignInfoErr
	}

	// 请求第三方更新短信模板
	var updateSmsTemplateRequest = sendcloud.UpdateSmsTemplatequest{
		SmsUser:         sendcloud.SendCloudSmsUser,
		TemplateIdStr:   strconv.FormatInt(smsTemplateInfo.TemplateId, 10),
		TemplateName:    smsTemplateInfo.TemplateName + "-" + strconv.FormatInt(smsTemplateInfo.EnterpriseId, 10),
		TemplateText:    checkMessageContent(smsTemplateInfo.MessageTemplate),
		SignId:          smsSignInfo.ThirdSignId,
		SignPositionStr: "0",
		SmsTypeStr:      "1",
	}

	// 请求第三方编辑短信模板
	signString := GetSign(updateSmsTemplateRequest, sendcloud.SendCloudSmsKey)
	updateSmsTemplateRequest.Signature = signString
	_, updateSmsTemplateErr := sendcloud.UpdateSmsTemplate(c, updateSmsTemplateRequest)

	if updateSmsTemplateErr != nil {
		return nil, updateSmsTemplateErr
	}

	var thirdTemplateId = strconv.FormatInt(smsTemplateInfo.TemplateId, 10)

	// 提交模板审核
	var smsTemplateSubmitsms = sendcloud.SmsTemplateSubmitsmsRequest{
		SmsUser:       sendcloud.SendCloudSmsUser,
		TemplateIdStr: thirdTemplateId,
	}
	// 请求第三方创建签名
	submitSmsSignString := GetSign(smsTemplateSubmitsms, sendcloud.SendCloudSmsKey)
	smsTemplateSubmitsms.Signature = submitSmsSignString
	_, createSmsTemplateErr := sendcloud.SmsTemplateSubmitsms(c, smsTemplateSubmitsms)
	if createSmsTemplateErr != nil {
		return nil, createSmsTemplateErr
	}

	var SmsTemplateArraignmentReponse = SmsTemplateArraignmentReponse{
		Id: request.Id,
	}
	var newTime = time.Now().Unix()
	// 修改短信模板状态
	var updateSmsTemplateStatus = UpdateSmsTemplateStatus{
		Id:        request.Id,
		Status:    TemplateStatusUnderReview,
		UpdatedAt: newTime,
	}
	smsSignInfoErr = s.Repo.UpdateTemplateStatus(c, updateSmsTemplateStatus)
	if smsSignInfoErr != nil {
		return nil, smsSignInfoErr
	}

	return &SmsTemplateArraignmentReponse, nil
}

// checkMessageContent 验证短信内容中是否存在退订回T存在的话不在拼接
func checkMessageContent(message string) string {
	if find := strings.Contains(message, "退订回T"); find {
		return message
	}
	return message + "退订回T"
}

// CreateSmsTemplateCheck 验证短信模板是否存在
func (s *Service) CreateSmsTemplateCheck(c ctx.Context, request SmsTemplate) (*SmsTemplate, error) {
	// 验证模板名称是否存在
	request.Id = 0
	SmsTemplateInfo, SmsTemplateInfoErr := s.Repo.SmsTemplateFindOne(c, request)
	if SmsTemplateInfoErr != nil {
		return nil, SmsTemplateInfoErr
	}
	// 获取模板信息
	return SmsTemplateInfo, nil
}

// CreateSmsTemplateSign 创建签名
func (s *Service) CreateSmsTemplateSign(c ctx.Context, request SmsSignInfoRequest) (*SmsSignInfo, error) {
	var smsSign, smsSignDataErr = s.Repo.SmsSignFindOne(c, request)
	if smsSignDataErr != nil && smsSignDataErr != sql.ErrNoRows {
		return nil, smsSignDataErr
	}
	// 查询到签名,但不是此企业下的签名，给当前企业新增签名记录
	if smsSign != nil && smsSign.EnterpriseId != request.EnterpriseId {
		// 创建短信签名成功添加到数据表
		smsSign.EnterpriseId = request.EnterpriseId
		smsSign.CreatedAt = time.Now().Unix()
		smsSign.Id = uid.ID(c)
		smsSign.ThirdSignStatus = 3
		_, cErr := s.Repo.CreateSmsSign(c, *smsSign)
		if cErr != nil {
			return nil, smsSignDataErr
		}
		return smsSign, nil
	}
	if smsSignDataErr == sql.ErrNoRows {
		//// 没查到标签调用第三方添加标签
		//var createSmsSignRequest = sendcloud.CreateSmsSignRequest{
		//	SignName: request.SignName,
		//	SmsUser:  sendcloud.SendCloudSmsUser,
		//	SignType: 0,
		//}
		//
		//// 请求第三方创建签名
		//signString := GetSign(createSmsSignRequest, sendcloud.SendCloudSmsKey)
		//createSmsSignRequest.Signature = signString
		//createSmsSignRes, createSmsSignErr := sendcloud.CreateSmsSign(c, createSmsSignRequest)
		//
		//if createSmsSignErr != nil {
		//	return nil, errs.TemplateSignFormat
		//}
		//
		//var thirdSignStatus = sendcloud.CheckSmsSignStatus(createSmsSignRes.Info.SmsSignVo.VerifyStatus)
		//var signTypeDesc = sendcloud.CheckSmsDesc(createSmsSignRes.Info.SmsSignVo.SignTypeDesc)
		var newTime = time.Now().Unix()
		// 创建短信签名成功添加到数据表
		var CreateSmsSign = SmsSignInfo{
			Id:              uid.ID(ctx.New()),
			EnterpriseId:    request.EnterpriseId,
			AdminId:         request.AdminId,
			ThirdSignStatus: 3,
			SignName:        request.SignName,
			CreatedAt:       newTime,
		}
		var _, createSmsSignDataErr = s.Repo.CreateSmsSign(c, CreateSmsSign)
		if createSmsSignDataErr != nil {
			return nil, smsSignDataErr
		}

		return &CreateSmsSign, nil
	} else if smsSignDataErr != nil {
		return nil, smsSignDataErr
	}

	return smsSign, nil

}

// CreateSmsTemplateSign 创建签名
func (s *Service) CreateSmsTemplateSignBak(c ctx.Context, request SmsSignInfoRequest) (*SmsSignInfo, error) {
	var smsSign, smsSignDataErr = s.Repo.SmsSignFindOne(c, request)
	if smsSignDataErr != nil && smsSignDataErr != sql.ErrNoRows {
		return nil, smsSignDataErr
	}
	// 查询到签名,但不是此企业下的签名，给当前企业新增签名记录
	if smsSign != nil && smsSign.EnterpriseId != request.EnterpriseId {
		// 创建短信签名成功添加到数据表
		smsSign.EnterpriseId = request.EnterpriseId
		smsSign.CreatedAt = time.Now().Unix()
		smsSign.Id = uid.ID(c)
		_, cErr := s.Repo.CreateSmsSign(c, *smsSign)
		if cErr != nil {
			return nil, smsSignDataErr
		}
		return smsSign, nil
	}
	if smsSignDataErr == sql.ErrNoRows {
		// 没查到标签调用第三方添加标签
		var createSmsSignRequest = sendcloud.CreateSmsSignRequest{
			SignName: request.SignName,
			SmsUser:  sendcloud.SendCloudSmsUser,
			SignType: 0,
		}

		// 请求第三方创建签名
		signString := GetSign(createSmsSignRequest, sendcloud.SendCloudSmsKey)
		createSmsSignRequest.Signature = signString
		createSmsSignRes, createSmsSignErr := sendcloud.CreateSmsSign(c, createSmsSignRequest)

		if createSmsSignErr != nil {
			return nil, errs.TemplateSignFormat
		}

		var thirdSignStatus = sendcloud.CheckSmsSignStatus(createSmsSignRes.Info.SmsSignVo.VerifyStatus)
		var signTypeDesc = sendcloud.CheckSmsDesc(createSmsSignRes.Info.SmsSignVo.SignTypeDesc)
		var newTime = time.Now().Unix()
		// 创建短信签名成功添加到数据表
		var CreateSmsSign = SmsSignInfo{
			Id:              uid.ID(ctx.New()),
			EnterpriseId:    request.EnterpriseId,
			AdminId:         request.AdminId,
			ThirdSignId:     createSmsSignRes.Info.SmsSignVo.Id,
			ThirdSignStatus: thirdSignStatus,
			ThirdSignType:   signTypeDesc,
			SignName:        createSmsSignRes.Info.SmsSignVo.SignName,
			CreatedAt:       newTime,
		}
		var _, createSmsSignDataErr = s.Repo.CreateSmsSign(c, CreateSmsSign)
		if createSmsSignDataErr != nil {
			return nil, smsSignDataErr
		}

		return &CreateSmsSign, nil
	} else if smsSignDataErr != nil {
		return nil, smsSignDataErr
	}

	return smsSign, nil

}

// SmsTemplateParamsList 短信模板参数列表
func (s *Service) SmsTemplateParamsList(c ctx.Context) (*SmsParamsRequest, error) {

	smsTemplateParamsList, err := s.Repo.SelectSmsTemplateParamsList(c)

	if err != nil {
		return nil, err
	}
	var SmsParamsRequest = SmsParamsRequest{
		List: *smsTemplateParamsList,
	}
	return &SmsParamsRequest, nil
}

// SmsTemplateSignList 短信签名列表
func (s *Service) SmsTemplateSignList(c ctx.Context, templateSignRequest TemplateListRequest) (*SmsSignListRequest, error) {
	total, list, err := s.Repo.SmsTemplateSignList(c, templateSignRequest)
	if err != nil {
		return nil, err
	}
	var smsTemplateRequest = &SmsSignListRequest{
		List:       *list,
		Page:       templateSignRequest.Page,
		Size:       templateSignRequest.Size,
		Count:      total,
		UpdateTime: time.Now().Unix(),
	}
	return smsTemplateRequest, nil
}

// TemplateListByTemplateIds 根据模板ID 查询模板
func (s *Service) TemplateListByTemplateIds(cc ctx.Context, enterpriseID int64, templateIds []int64, status int64) (map[int64]*SmsTemplate, error) {
	req := TemplateListRequest{
		EnterpriseId:  enterpriseID,
		TemplateIds:   templateIds,
		TemplateIdCnt: len(templateIds),
		Status:        status,
		Page:          1,
		Size:          int64(len(templateIds)),
	}
	total, list, lErr := s.Repo.SmsTemplateList(cc, req)
	if lErr != nil {
		return nil, lErr
	}
	templateMap := make(map[int64]*SmsTemplate)
	if total > 0 {
		for _, item := range *list {
			templateMap[item.Id] = &item
		}
	}

	return templateMap, nil
}

// SmsTemplateList 短信模板列表
func (s *Service) SmsTemplateList(c ctx.Context, TemplateListRequest TemplateListRequest) (*SmsTemplateListRequest, error) {
	var smsTemplate = SmsTemplate{
		EnterpriseId: TemplateListRequest.EnterpriseId,
	}
	s.CheckDefaultSmsTemplate(c, smsTemplate, true)
	total, list, lErr := s.Repo.SmsTemplateList(c, TemplateListRequest)
	if lErr != nil {
		log.Error(c).Msgf("查询短信模板列表出错，请求参数:%+v,错误信息：%+v", TemplateListRequest, lErr)
		return &SmsTemplateListRequest{
			List:  make([]SmsTemplate, 0),
			Count: 0,
			Page:  TemplateListRequest.Page,
			Size:  TemplateListRequest.Size,
		}, nil
	}
	if total > 0 {
		for k, vItem := range *list {
			// 如果是短信加粉
			if TemplateTypeIsFans == vItem.Type {
				templateNameArr := strings.Split(vItem.TemplateName, "-")
				(*list)[k].TemplateName = templateNameArr[0]
			}
		}
	}
	smsTemplateRequest := &SmsTemplateListRequest{
		List:       *list,
		Page:       TemplateListRequest.Page,
		Size:       TemplateListRequest.Size,
		Count:      total,
		TpdateTime: time.Now().Unix(),
	}

	return smsTemplateRequest, nil
}

// CheckDefaultSmsTemplateByEnterpriseID 根据企业ID 查询默认短信模板
func (s *Service) CheckDefaultSmsTemplateByEnterpriseID(c ctx.Context, enterpriseID int64) (int64, error) {
	templateID, err := s.CheckDefaultSmsTemplate(c, SmsTemplate{
		EnterpriseId: enterpriseID,
	}, true)

	if err != nil {
		return 0, err
	}

	return templateID, nil
}

// CheckDefaultSmsTemplate 检查各种异常提醒短信模板是否存在
func (s *Service) CheckDefaultSmsTemplate(c ctx.Context, smsTemplate SmsTemplate, IsCheckOther bool) (int64, error) {
	var (
		templateId int64
	)
	// 查询内置短信模板
	defaultTemplate, defaultTemplateErr := s.GetDefaultTemplate(c, smsTemplate.EnterpriseId, DefaultName)
	if defaultTemplateErr != nil || defaultTemplate == nil {
		templateId, _ = s.CreateSmsTemplateDefault(c, smsTemplate.EnterpriseId, DefaultName, DefaultTemplateContent, DefaultTemplateContentStr, DefaultTemplateParams)
	} else {
		templateId = defaultTemplate.Id
	}
	// 不检查其他短信模板直接返回
	if !IsCheckOther {
		return templateId, nil
	}

	// 查询是否存在营销策略模板
	abnormalTemplate, abnormalTemplateErr := s.GetDefaultTemplate(c, smsTemplate.EnterpriseId, AbnormalName)
	// 不存在则创建
	if abnormalTemplateErr != nil || abnormalTemplate == nil {
		_, _ = s.CreateSmsTemplateDefault(c, smsTemplate.EnterpriseId, AbnormalName, AbnormalTemplateContent, AbnormalTemplateContentStr, AbnormalTemplateParams)
	}

	// 查询是否存在短信被置为未提审状态同步策略触发短信内容
	ReappearingTemplate, ReappearingTemplateErr := s.GetDefaultTemplate(c, smsTemplate.EnterpriseId, ReappearingName)
	// 不存在则创建
	if ReappearingTemplateErr != nil || ReappearingTemplate == nil {
		_, _ = s.CreateSmsTemplateDefault(c, smsTemplate.EnterpriseId, ReappearingName, ReappearingTemplateContent, ReappearingTemplateContentStr, ReappearingTemplateParams)
	}

	// 查询是否存在短信被置为未提审状态同步策略触发短信内容
	OverdueTemplate, OverdueTemplateErr := s.GetDefaultTemplate(c, smsTemplate.EnterpriseId, OverdueName)
	// 不存在则创建
	if OverdueTemplateErr != nil || OverdueTemplate == nil {
		_, _ = s.CreateSmsTemplateDefault(c, smsTemplate.EnterpriseId, OverdueName, OverdueTemplateContent, OverdueTemplateContentStr, OverdueTemplateParams)
	}

	// 查询是否存在优惠券不足100条内置短信模板
	CouponTemplate, CouponTemplateErr := s.GetDefaultTemplate(c, smsTemplate.EnterpriseId, CouponName)
	// 不存在则创建
	if CouponTemplateErr != nil || CouponTemplate == nil {
		_, _ = s.CreateSmsTemplateDefault(c, smsTemplate.EnterpriseId, CouponName, CouponTemplateContent, CouponTemplateContentStr, CouponTemplateParams)
	}

	return templateId, nil
}

func (s *Service) CreateSmsTemplateDefault(c ctx.Context, EnterpriseID int64, TemplateName string, Template string, TemplateContentStr string, Content string) (int64, error) {
	// send clound 请求参数
	req := sendcloud.CreateSmsTemplatequest{
		SmsUser:         sendcloud.SendCloudSmsUser,
		TemplateName:    fmt.Sprintf(TemplateName+"-%d", EnterpriseID),
		TemplateText:    Template + "退订回T",
		SignName:        DefaultSignName,
		SignId:          DefaultSignID,
		SignPositionStr: "0",
		SmsTypeStr:      "1",
	}
	// 签名
	req.Signature = GetSign(req, sendcloud.SendCloudSmsKey)
	createSmsTemplateRes, createSmsTemplateErr := sendcloud.CreateSmsTemplate(c, req)
	if createSmsTemplateErr != nil {
		log.Error(c).Msgf("创建短信模板出错，请求信息：%+v,错误信息：%+v", req, createSmsTemplateErr)
		return 0, createSmsTemplateErr
	}

	// 不存在则创建
	var smsTemplateRequest = SmsTemplate{
		Id:              uid.ID(ctx.New()),
		EnterpriseId:    EnterpriseID,
		TemplateName:    req.TemplateName,
		TemplateSignId:  DefaultSignSysID,
		MessageTemplate: Template,
		MessageStr:      TemplateContentStr,
		TemplateContent: Content,
		Status:          TemplateStatusWaitVerify,
		CreatedAt:       time.Now().Unix(),
		TemplateId:      createSmsTemplateRes.Info.TemplateId,
		Type:            TemplateTypeIsFans,
	}

	// 提交审核
	approvalReq := sendcloud.SmsTemplateSubmitsmsRequest{
		SmsUser:       sendcloud.SendCloudSmsUser,
		TemplateIdStr: strconv.FormatInt(smsTemplateRequest.TemplateId, 10),
	}
	submitSmsSignString := GetSign(approvalReq, sendcloud.SendCloudSmsKey)
	approvalReq.Signature = submitSmsSignString
	_, subErr := sendcloud.SmsTemplateSubmitsms(c, approvalReq)
	if subErr == nil {
		smsTemplateRequest.Status = TemplateStatusUnderReview
	}
	_, _ = s.Repo.CreateSmsTemplate(c, smsTemplateRequest)
	return smsTemplateRequest.Id, nil
}

// SmsTemplateParamsList 短信模板参数列表
func (s *Service) SmsTemplateDelete(c ctx.Context, request UpdateSmsTemplateStatus) (*SmsTemplateArraignmentReponse, error) {
	var newTime = time.Now().Unix()
	// 删除短信模板
	var updateSmsTemplateStatus = UpdateSmsTemplateStatus{
		Id:        request.Id,
		DeletedAt: newTime,
	}
	TemplateDeleteErr := s.Repo.SmsTemplateDelete(c, updateSmsTemplateStatus)

	if TemplateDeleteErr != nil {
		return nil, TemplateDeleteErr
	}

	var TemplateDeleteReponse = SmsTemplateArraignmentReponse{
		Id: request.Id,
	}
	return &TemplateDeleteReponse, nil
}

// SmsTemplateEdit 编辑短信模板
func (s *Service) SmsTemplateEdit(c ctx.Context, request UpdateSmsTemplateStatus) (*SmsTemplateArraignmentReponse, error) {
	var newTime = time.Now().Unix()
	// 编辑短信模板
	var updateSmsTemplateStatus = UpdateSmsTemplateStatus{
		Id:              request.Id,
		TemplateName:    request.TemplateName,
		TemplateSignId:  request.TemplateSignId,
		MessageTemplate: request.MessageTemplate,
		TemplateContent: request.TemplateContent,
		TemplateHtml:    request.TemplateHtml,
		MessageStr:      request.MessageStr,
		UpdatedAt:       newTime,
	}

	// 检测短信模板是否全部赋值
	//if request.TemplateContent != "" {
	//	//1.创建json数据
	//	b := []byte(request.TemplateContent)
	//	//2.声明空interface
	//	var i interface{}
	//	//3.解析
	//	err := json.Unmarshal(b, &i)
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	//默认转成了map
	//	fmt.Println(i)
	//	//4.解析到interface的json可以判断类型
	//	m := i.(map[string]interface{})
	//	for _, v := range m {
	//		switch vv := v.(type) {
	//		case string:
	//			if vv != "" {
	//				request.Status = 5
	//			}
	//		default:
	//
	//		}
	//	}
	//}

	smsTemplate := SmsTemplate{
		Id: request.Id,
	}
	res, err := s.SmsTemplateInfo(c, smsTemplate)
	if err != nil {
		return nil, err
	}
	// 审核中不可编辑
	//if res.Status == TemplateStatusUnderReview {
	//	return nil, errs.CannotEditInReview
	//}
	//// 审核成功后短信模板内容不可编辑
	//if res.Status == TemplateStatusSuccessfulReview || res.Status == TemplateStatusHasValue {
	//	if updateSmsTemplateStatus.MessageTemplate != res.MessageTemplate {
	//		return nil, errs.CannotEditNot
	//	}
	//}
	request.Status = 5
	// 如果模板选项发生改变
	if res.TemplateSignId != request.TemplateSignId {
		request.Status = TemplateStatusWaitVerify
		if res.Type == TemplateTypeIsFans {
			s.Marketing.CloseMarketingInnerStatus(c, res.EnterpriseId)
		}
	}

	// 检测短信模板名称如果是“粉丝”在短信名称后面拼接企业ID
	if res.Type == TemplateTypeIsFans {
		updateSmsTemplateStatus.TemplateName = fmt.Sprintf("短信加粉-%d", request.EnterpriseId)
	}

	// 如果已赋值将状态改为已赋值
	if request.Status == TemplateStatusHasValue {
		updateSmsTemplateStatus.Status = request.Status
	}

	TemplateDeleteErr := s.Repo.SmsTemplateEdit(c, updateSmsTemplateStatus)

	if TemplateDeleteErr != nil {
		return nil, TemplateDeleteErr
	}

	// 如果短信模板状态是 4：审核成功、5：已赋值
	if res.Status == TemplateStatusSuccessfulReview || res.Status == TemplateStatusHasValue {
		if res.TemplateSignId != updateSmsTemplateStatus.TemplateSignId {
			// 更新模板
			_, err := s.UpdateSmsTemplate(c, updateSmsTemplateStatus)
			if err != nil {
				return nil, err
			}
		}
	}
	// 赋值
	var TemplateDeleteReponse = SmsTemplateArraignmentReponse{
		Id: request.Id,
	}
	return &TemplateDeleteReponse, nil
}

// SmsTemplateEditBak 编辑短信模板
func (s *Service) SmsTemplateEditBak(c ctx.Context, request UpdateSmsTemplateStatus) (*SmsTemplateArraignmentReponse, error) {
	var newTime = time.Now().Unix()
	// 编辑短信模板
	var updateSmsTemplateStatus = UpdateSmsTemplateStatus{
		Id:              request.Id,
		TemplateName:    request.TemplateName,
		TemplateSignId:  request.TemplateSignId,
		MessageTemplate: request.MessageTemplate,
		TemplateContent: request.TemplateContent,
		TemplateHtml:    request.TemplateHtml,
		MessageStr:      request.MessageStr,
		UpdatedAt:       newTime,
	}

	// 检测短信模板是否全部赋值
	if request.TemplateContent != "" {
		//1.创建json数据
		b := []byte(request.TemplateContent)
		//2.声明空interface
		var i interface{}
		//3.解析
		err := json.Unmarshal(b, &i)
		if err != nil {
			fmt.Println(err)
		}
		//默认转成了map
		fmt.Println(i)
		//4.解析到interface的json可以判断类型
		m := i.(map[string]interface{})
		for _, v := range m {
			switch vv := v.(type) {
			case string:
				if vv != "" {
					request.Status = 5
				}
			default:

			}
		}
	}

	smsTemplate := SmsTemplate{
		Id: request.Id,
	}
	res, err := s.SmsTemplateInfo(c, smsTemplate)
	if err != nil {
		return nil, err
	}
	// 审核中不可编辑
	if res.Status == TemplateStatusUnderReview {
		return nil, errs.CannotEditInReview
	}
	// 审核成功后短信模板内容不可编辑
	if res.Status == TemplateStatusSuccessfulReview || res.Status == TemplateStatusHasValue {
		if updateSmsTemplateStatus.MessageTemplate != res.MessageTemplate {
			return nil, errs.CannotEditNot
		}
	}

	// 如果模板选项发生改变
	if res.TemplateSignId != request.TemplateSignId {
		request.Status = TemplateStatusWaitVerify
		if res.Type == TemplateTypeIsFans {
			s.Marketing.CloseMarketingInnerStatus(c, res.EnterpriseId)
		}
	}

	// 检测短信模板名称如果是“粉丝”在短信名称后面拼接企业ID
	if res.Type == TemplateTypeIsFans {
		updateSmsTemplateStatus.TemplateName = fmt.Sprintf("短信加粉-%d", request.EnterpriseId)
	}

	// 如果已赋值将状态改为已赋值
	if request.Status == TemplateStatusHasValue {
		updateSmsTemplateStatus.Status = request.Status
	}

	TemplateDeleteErr := s.Repo.SmsTemplateEdit(c, updateSmsTemplateStatus)

	if TemplateDeleteErr != nil {
		return nil, TemplateDeleteErr
	}

	// 如果保存并提审
	if request.IsExamine == 1 {
		var arraignmentRequest = SmsTemplateArraignmentRequest{
			Id: request.Id,
		}

		// 更新模板
		_, err := s.UpdateSmsTemplate(c, updateSmsTemplateStatus)

		if err != nil {
			return nil, err
		}

		_, ArraignmentRequestErr := s.SmsTemplateArraignment(c, arraignmentRequest)
		if ArraignmentRequestErr != nil {
			return nil, ArraignmentRequestErr
		}
	} else {
		// 如果短信模板状态是 4：审核成功、5：已赋值
		if res.Status == TemplateStatusSuccessfulReview || res.Status == TemplateStatusHasValue {
			if res.TemplateSignId != updateSmsTemplateStatus.TemplateSignId {
				// 更新模板
				_, err := s.UpdateSmsTemplate(c, updateSmsTemplateStatus)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	// 赋值
	var TemplateDeleteReponse = SmsTemplateArraignmentReponse{
		Id: request.Id,
	}
	return &TemplateDeleteReponse, nil
}

// FindSmsTemplateStatusByID 通过短信模板ID 查询模板状态
func (s *Service) FindSmsTemplateStatusByID(c ctx.Context, id int64) (int64, error) {
	info, err := s.Repo.SmsTemplateFindOne(c, SmsTemplate{
		Id: id,
	})
	if err != nil {
		return 0, err
	}
	return info.Status, nil
}

// SmsTemplateParamsList 短信模板参数列表
func (s *Service) SmsTemplateInfo(c ctx.Context, request SmsTemplate) (*SmsTemplate, error) {
	var smsTemplate = SmsTemplate{
		Id: request.Id,
	}
	smsTemplateInfo, smsSignInfoErr := s.Repo.SmsTemplateFindOne(c, smsTemplate)

	if smsSignInfoErr != nil && smsSignInfoErr != sql.ErrNoRows {
		return nil, smsSignInfoErr
	}

	if smsTemplateInfo != nil {
		smsTemplateInfo.TemplateShortName = smsTemplateInfo.TemplateName
		if smsTemplateInfo.Type == TemplateTypeIsFans {
			templateNameArr := strings.Split(smsTemplateInfo.TemplateName, "-")
			smsTemplateInfo.TemplateShortName = templateNameArr[0]
		}

		total, _ := s.Marketing.GetSmsCount(c, request.Id)

		if int(total) > 0 {
			smsTemplateInfo.Occupancy = 1
		}
	}

	return smsTemplateInfo, nil
}

// GetDefaultTemplate 获取企业默认的短信模板
func (s *Service) GetDefaultTemplate(c ctx.Context, EnterpriseId int64, Name string) (*SmsTemplate, error) {

	smsTemplateInfo, smsSignInfoErr := s.Repo.GetDefaultTemplate(c, EnterpriseId, Name)

	if smsSignInfoErr != nil && smsSignInfoErr != sql.ErrNoRows {
		smsTemplateInfo := &SmsTemplate{}
		return smsTemplateInfo, smsSignInfoErr
	}
	return smsTemplateInfo, nil
}

// SmsTemplateParamsList 短信模板参数列表
func (s *Service) SmsTemplateSignDelete(c ctx.Context, request SmsSignInfoRequest) error {
	var newTime = time.Now().Unix()
	var smsTemplateDelete = SmsSignInfoRequest{
		Id:        request.Id,
		DeletedAt: newTime,
		UpdatedAt: newTime,
	}
	smsSignInfoErr := s.Repo.UpdateSmsSign(c, smsTemplateDelete)

	if smsSignInfoErr != nil && smsSignInfoErr != sql.ErrNoRows {
		return smsSignInfoErr
	}

	return nil
}

// SynTemplateStatus 同步短信模板状态
func (s *Service) SynTemplateStatus(c ctx.Context, EnterpriseId int64) error {
	// 查询企业短信模板总数
	SmsTemplateRequest := SmsTemplate{
		EnterpriseId: EnterpriseId,
	}
	count := s.Repo.GetEnterpriseSmsTemplateCount(c, SmsTemplateRequest)

	var size int64 = 100
	var page int64 = 1
	for size*page < count+size {
		// 通过分页的方式查询企业短信模板
		TemplateListRequest := TemplateListRequest{
			EnterpriseId: EnterpriseId,
			Page:         page,
			Size:         size,
		}
		_, list, lErr := s.Repo.SmsTemplateList(c, TemplateListRequest)
		if lErr != nil {
			return nil
		}
		page++

		smsTemplateListById := make(map[int64]SmsTemplate)

		var TemplateThirdIds bytes.Buffer
		for _, value := range *list {
			if TemplateThirdIds.Len() > 0 {
				TemplateThirdIds.WriteByte(';')
			}
			TemplateThirdIds.WriteString(strconv.FormatInt(value.TemplateId, 10))

			smsTemplateListById[value.TemplateId] = value
		}

		// 根据短信第三方ID批量查询多个模板状态信息
		var TemplateListQuest = sendcloud.TemplateListRequest{
			SmsUser:       sendcloud.SendCloudSmsUser,
			TemplateIdStr: TemplateThirdIds.String(),
			Start:         0,
			Limit:         int64(len(*list)),
		}

		// 请求第三方创建签名
		TemplateListQuest.Signature = GetSign(TemplateListQuest, sendcloud.SendCloudSmsKey)
		TemplateListRes, TemplateListErr := sendcloud.SmsTemplateList(c, TemplateListQuest)

		if TemplateListErr != nil {
			return TemplateListErr
		}

		// 短信模板参数列表
		smsParamsList, _ := s.Repo.SelectSmsTemplateParamsList(c)

		for _, value := range TemplateListRes.Info.List {
			var status int64
			switch value.IsVerify {
			case "未提审":
				if smsTemplateListById[value.TemplateId].Status == TemplateStatusWaitVerify {
					continue
				}
				// 修改模板状态为未提审
				status = TemplateStatusWaitVerify
			case "待审核":
				if smsTemplateListById[value.TemplateId].Status == TemplateStatusUnderReview {
					continue
				}
				// 修改模板状态为未提审
				status = TemplateStatusUnderReview
			case "审核中":
				if smsTemplateListById[value.TemplateId].Status == TemplateStatusUnderReview {
					continue
				}
				// 修改模板状态为审核中
				status = TemplateStatusUnderReview
			case "待确认":
				if smsTemplateListById[value.TemplateId].Status == TemplateStatusUnderReview {
					continue
				}
				// 修改模板状态为审核中
				status = TemplateStatusUnderReview
			case "审核通过":
				if smsTemplateListById[value.TemplateId].Status == TemplateStatusSuccessfulReview || smsTemplateListById[value.TemplateId].Status == TemplateStatusHasValue {
					continue
				}
				// 修改模板状态为审核通过
				status = TemplateStatusSuccessfulReview
			default:
				if smsTemplateListById[value.TemplateId].Status == TemplateStatusAuditFailure {
					continue
				}
				// 默认审核失败
				status = TemplateStatusAuditFailure
			}

			if status == TemplateStatusSuccessfulReview {
				status = TemplateStatusHasValue
				for _, param := range *smsParamsList {
					if param.ParamsType == ParamsTypeStatic && strings.Contains(smsTemplateListById[value.TemplateId].MessageTemplate, "%"+param.Key+"%") {
						status = TemplateStatusSuccessfulReview
						break
					}
				}
			}

			var updateSmsTemplateStatus = UpdateSmsTemplateStatus{
				TemplateId: value.TemplateId,
				Status:     status,
				FailReason: value.IsVerify,
			}
			_ = s.Repo.SmsTemplateEdit(c, updateSmsTemplateStatus)
		}
	}

	// 根据短信第三方ID批量查询多个模板状态信息
	var TemplateSignListQuest = sendcloud.TemplateSignListRequest{
		SmsUser: sendcloud.SendCloudSmsUser,
	}

	// 请求第三方创建签名
	TemplateSignListQuest.Signature = GetSign(TemplateSignListQuest, sendcloud.SendCloudSmsKey)
	TemplateSignListRes, TemplateSignListErr := sendcloud.SmsTemplateSignList(c, TemplateSignListQuest)
	if TemplateSignListErr != nil {
		return TemplateSignListErr
	}

	var signStatus int64
	for _, value := range TemplateSignListRes.Info.VoList {
		switch value.VerifyStatus {
		case "审核通过":
			signStatus = ThirdSignStatusSuccessful
		case "待审核":
			signStatus = ThirdSignStatusUnderReview
		default:
			signStatus = ThirdSignStatusFailure
		}
		var updateSmsTemplateStatus = SmsSignInfoRequest{
			ThirdSignId:     value.Id,
			ThirdSignStatus: signStatus,
		}
		_ = s.Repo.UpdateSmsSign(c, updateSmsTemplateStatus)
	}

	return nil
}

// SmsSend 发送短信
func (s *Service) SmsSend(c ctx.Context, request SendRequest) error {
	// 查询短信模板
	var smsTemplate = SmsTemplate{
		Id: request.TemplateId,
	}

	smsTemplateInfo, smsSignInfoErr := s.Repo.SmsTemplateFindOne(c, smsTemplate)

	if smsSignInfoErr != nil && smsSignInfoErr != sql.ErrNoRows {
		return smsSignInfoErr
	}

	if smsTemplateInfo == nil {
		return errs.SmsSignNotReviewed
	}

	// 组装发送短信参数
	var sendQuest = sendcloud.SendMsgRequest{
		Template:     smsTemplateInfo.TemplateId,
		Mobile:       request.Phone,
		EnterpriseId: request.EnterpriseId,
	}
	content := smsTemplateInfo.TemplateContent
	if content != "" {
		if request.Vars.UserName != "" && strings.Contains(content, "%user_name%") {
			content = strings.Replace(content, "%user_name%", request.Vars.UserName, -1)
		}
		if request.Vars.UserNick != "" && strings.Contains(content, "%user_nick%") {
			content = strings.Replace(content, "%user_nick%", request.Vars.UserNick, -1)
		}
		content = strings.Replace(content, "%user_name%", "", -1)
		content = strings.Replace(content, "%user_nick%", "", -1)
		sendQuest.Content = content
	}
	// 请求第三方发送短信
	sendRes, TemplateSignListErr := microauthority.SendMsg(c, sendQuest)

	if TemplateSignListErr != nil {
		return TemplateSignListErr
	}
	if sendRes.Code == 200 {
		// 发送短信成功存储发送日志记录
		insertErr := s.InsertSendLog(c, request)
		if insertErr != nil {
			return insertErr
		}
	}
	return nil
}

// SmsSendBak 发送短信
func (s *Service) SmsSendBak(c ctx.Context, request SendRequest) error {
	// 查询短信模板
	var smsTemplate = SmsTemplate{
		Id: request.TemplateId,
	}

	// 查询企业配置了多少短信
	req := &boss.EnterpriseDetailRequest{
		EnterpriseID: request.EnterpriseId,
	}
	bossSmsData, err := boss.EnterpriseSmsTotal(c, req)

	if err != nil {
		return err
	}
	if bossSmsData.SmsNumber < 1 {
		return errs.NotSmsNumber
	}

	smsTemplateInfo, smsSignInfoErr := s.Repo.SmsTemplateFindOne(c, smsTemplate)

	if smsSignInfoErr != nil && smsSignInfoErr != sql.ErrNoRows {
		return smsSignInfoErr
	}

	if smsTemplateInfo == nil {
		return errs.SmsSignNotReviewed
	}

	// 组装发送短信参数
	var sendQuest = sendcloud.SendRequest{
		SmsUser:    sendcloud.SendCloudSmsUser,
		TemplateId: smsTemplateInfo.TemplateId,
		Phone:      request.Phone,
	}
	if smsTemplateInfo.TemplateContent != "" {
		var dat map[string]interface{}
		if err := json.Unmarshal([]byte(smsTemplateInfo.TemplateContent), &dat); err == nil {
			for key, value := range dat {
				if key == "user_name" && value == "" && request.Vars.UserName != "" {
					dat[key] = request.Vars.UserName
				}
				//if key == "address_name" && value == "" && request.Vars.AddressName != "" {
				//	dat[key] = request.Vars.AddressName
				//}
				//if key == "order_id" && value == "" && request.Vars.OrderId != "" {
				//	dat[key] = request.Vars.OrderId
				//}
			}
		}
		myson, _ := json.Marshal(dat)
		vars := string(myson)
		sendQuest.Vars = vars
	}
	sendQuest.Signature = GetSign(sendQuest, sendcloud.SendCloudSmsKey)
	// 请求第三方发送短信
	sendRes, TemplateSignListErr := sendcloud.Send(c, sendQuest)

	if TemplateSignListErr != nil {
		return TemplateSignListErr
	}
	if sendRes.StatusCode == 200 {
		// 发送短信成功存储发送日志记录
		insertErr := s.InsertSendLog(c, request)
		if insertErr != nil {
			return insertErr
		}

		// 通知boss短信发送成功累计加一
		req := &boss.EnterpriseDepleteSmsRequest{
			EnterpriseID: request.EnterpriseId,
			Number:       1,
		}
		_ = boss.EnterpriseDeplete(c, req)
	}
	return nil
}

// InsertSendLog 存储短信发送日志
func (s *Service) InsertSendLog(c ctx.Context, request SendRequest) error {
	CreateSendLog := SendLog{
		Id:             uid.ID(ctx.New()),
		EnterpriseId:   request.EnterpriseId,
		UserId:         request.UserId,
		SmsTemplateId:  request.TemplateId,
		ChannelType:    request.ChannelType,
		ChannelValueId: request.ChannelValueId,
		Status:         1,
		CreatedAt:      time.Now().Unix(),
	}
	_, CreateSendLogErr := s.Repo.CreateSmsLog(c, CreateSendLog)
	if CreateSendLogErr != nil {
		return CreateSendLogErr
	}
	return nil
}

// SmsTemplateTotal 获取boss企业可发送短信条数
func (s *Service) SmsTemplateTotal(c ctx.Context, request Request) (*TemplateTotalRequest, error) {
	// 查询当前短信发送条数
	SendCount, _ := s.Repo.SmsCount(c, request)

	// 查询企业配置了多少短信
	req := &boss.EnterpriseDetailRequest{
		EnterpriseID: request.EnterpriseId,
	}
	bossSmsData, err := boss.EnterpriseSmsTotal(c, req)

	if err != nil {
		return nil, err
	}

	var AvailableNum int64 = 0
	if bossSmsData.SmsTotalNumber != 0 {
		AvailableNum = bossSmsData.SmsNumber
	}

	// 查询短信模板
	var TemplateTotalRequest = TemplateTotalRequest{
		AvailableNum: AvailableNum,
		UsedNum:      SendCount.Count,
	}
	return &TemplateTotalRequest, nil
}

// GetSign 获取签名
func GetSign(data interface{}, appSecret string) string {
	md5ctx := md5.New()

	switch v := reflect.ValueOf(data); v.Kind() {
	case reflect.String:
		md5ctx.Write([]byte(v.String() + appSecret))
		return hex.EncodeToString(md5ctx.Sum(nil))
	case reflect.Struct:
		orderStr := StructToMapSing(v.Interface(), appSecret)
		md5ctx.Write([]byte(orderStr))
		return hex.EncodeToString(md5ctx.Sum(nil))
	case reflect.Ptr:
		originType := v.Elem().Type()
		if originType.Kind() != reflect.Struct {
			return ""
		}
		dataType := reflect.TypeOf(data).Elem()
		dataVal := v.Elem()
		orderStr := buildOrderStr(dataType, dataVal, appSecret)
		md5ctx.Write([]byte(orderStr))
		return hex.EncodeToString(md5ctx.Sum(nil))
	default:
		return ""
	}
}

// buildOrderStr 组装字符串
func buildOrderStr(t reflect.Type, v reflect.Value, appSecret string) (returnStr string) {
	keys := make([]string, 0, t.NumField())
	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Tag.Get("json") == "signature" {
			continue
		}
		data[t.Field(i).Tag.Get("json")] = v.Field(i).Interface()

		keys = append(keys, t.Field(i).Tag.Get("json"))
	}
	sort.Sort(sort.StringSlice(keys))
	var buf bytes.Buffer
	buf.WriteString(appSecret)
	for _, k := range keys {
		if data[k] == "" {
			continue
		}
		checkRes := CheckISNotNil(data[k])
		if checkRes != true {
			continue
		}
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(k)
		buf.WriteByte('=')
		getToStr, getToStrBool := InterfaceToStr(data[k])
		if getToStrBool != true {
			continue
		}
		buf.WriteString(getToStr)
	}
	buf.WriteString("&" + appSecret)
	returnStr = buf.String()
	return returnStr
}

// InterfaceToStr 将未知数据转换为字符串
func InterfaceToStr(data interface{}) (string, bool) {
	switch vv := data.(type) {
	case string:
		return vv, true
	case int:
		return strconv.FormatInt(int64(vv), 10), true
	case int64:
		return strconv.FormatInt(int64(vv), 10), true
	case struct{}:
		bytes, _ := json.Marshal(vv)
		return string(bytes), true
	case interface{}:
		bytes, _ := json.Marshal(vv)
		return string(bytes), true
	default:
		return "", false
	}
}

// 判断数据是否为空
func CheckISNotNil(data interface{}) bool {
	switch value := data.(type) {
	case string:
		if value == "" {
			return false
		}
	case struct{}:
		bytes, _ := json.Marshal(value)
		if string(bytes) == "" {
			return false
		}
	case interface{}:
		bytes, _ := json.Marshal(value)
		if string(bytes) == "" {
			return false
		}
	}
	return true
}

func StructToMapSing(content interface{}, appSecret string) (returnStr string) {

	t := reflect.TypeOf(content)
	v := reflect.ValueOf(content)

	returnStr = buildOrderStr(t, v, appSecret)
	return returnStr
}

// SelectTemplateList 查询模版ID
func (s *Service) SelectTemplateListByIDs(c ctx.Context, enterpriseID int64, templateIDs []int64, strategyIDs []int64) (*[]SmsTemplate, error) {
	resp := &[]SmsTemplate{}
	resp, err := s.Repo.SelectTemplateListByIDs(c, enterpriseID, templateIDs, strategyIDs)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
