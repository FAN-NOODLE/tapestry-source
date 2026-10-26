package messageSvc

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	microBossDao "micro-scrm/internal/app/dao/micro_boss"
	"micro-scrm/internal/app/http/request"
	"micro-scrm/internal/app/service"
	cacheSvc "micro-scrm/internal/app/service/cache"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sdk/microwork"
	qyWork "micro-scrm/internal/pkg/sdk/qy_work"
	"micro-scrm/internal/pkg/sdk/work"
	"strconv"
	"strings"
)

type MessageService struct {
	service.Base
}

func NewMessageService() *MessageService {
	return &MessageService{}
}

// SendMiniMessage 发送小程序应用消息
func (s MessageService) SendMiniMessage(c app.Ctx, request microwork.SendMiniMessageRequest) (string, error) {
	if len(request.MemUserIds) == 0 {
		return "", errors.New("员工id不能为空")
	}

	if len(request.MemUserIds) > 1000 {
		return "", errors.New("员工id不能超过1000")
	}

	// 获取企业插件表 小程序appid
	pluginInfo, err := cacheSvc.NewRedisData().GetEnterpriseWorkPlugin(c, request.EnterpriseID, work.QWCrsStoreAssistantSelf)
	if err != nil || pluginInfo.Secret == "" {
		return "", errors.New("获取企业微信自建应用失败")
	}

	workClient := qyWork.NewQyWork(request.EnterpriseID)

	// 组装content消息体
	var contentItems []qyWork.ContentItem
	for _, item := range request.Content {
		var contentItem qyWork.ContentItem
		_ = mapstructure.Decode(item, &contentItem)
		contentItems = append(contentItems, contentItem)
	}

	messageReq := qyWork.SendMiniMessageReq{
		SendMessageBase: qyWork.SendMessageBase{
			ToUser:  strings.Join(request.MemUserIds, "|"),
			MsgType: "miniprogram_notice",
		},
		MiniProgramNotice: qyWork.MiniProgramNotice{
			Appid:       pluginInfo.MiniAppID,
			Page:        request.Page,
			Title:       request.Title,
			ContentItem: contentItems,
		},
		EnableIdTrans:        0,
		EnableDuplicateCheck: 0,
	}

	messageResp, err := workClient.SendMiniMessage(c, messageReq)
	if err != nil {
		s.Log().Error(c).Msgf("发送小程序消息失败")
		return "", err
	}

	return messageResp.MsgId, err
}

// SendH5TemplateCardMessage 发送H5的模板消息
func (s MessageService) SendH5TemplateCardMessage(c ctx.Context, params request.SendCardMessageToH5) (string, error) {
	res, err := microBossDao.NewEnterpriseConfig().GetInfo(c, map[string]interface{}{"enterprise_id": params.EnterpriseId, "key": "web_h5_domain"})
	if err != nil {
		return "", err
	}
	if res.Name == "" {
		return "", err
	}
	urlH5 := res.Name
	enterpriseInfo, err := microBossDao.NewEnterprise().GetOne(c, map[string]interface{}{
		"id":         params.EnterpriseId,
		"deleted_at": 0,
	})
	if err != nil {
		return "", err
	}
	baseUrl := fmt.Sprintf("%s/work/%s/", urlH5, enterpriseInfo.Hash)
	url := fmt.Sprintf("%s%s", baseUrl, params.Uri)
	jumpList := make([]qyWork.JumpItem, 0)
	jumpList = append(jumpList, qyWork.JumpItem{
		Type:  1,
		Title: helper.RuleThree(params.Language == "en", "Details", "详情"),
		URL:   url,
	})
	templateCardMessage := qyWork.SendTemplateCardMessageRequest{
		SendMessageRequest: qyWork.SendMessageRequest{
			EnterpriseID: params.EnterpriseId,
			Title:        params.Title,
			Content:      params.Content,
			MemUserIds:   params.UserIds,
			PartyIds:     params.Party,
		},
		CardType: "text_notice",
		JumpList: jumpList,
		CardAction: qyWork.CardAction{
			Type: 1,
			URL:  url,
		},
	}

	return s.sendTemplateCardMessage(c, params.AppSign, templateCardMessage)

}

// sendTemplateCardMessage 发送模板卡片消息
func (s MessageService) sendTemplateCardMessage(c app.Ctx, identifier string, request qyWork.SendTemplateCardMessageRequest) (string, error) {
	workClient := qyWork.NewQyWork(request.EnterpriseID)
	// 获取企业插件表 小程序appid
	pluginInfo, err := cacheSvc.NewRedisData().GetEnterpriseWorkPlugin(c, request.EnterpriseID, identifier)
	if err != nil || pluginInfo.Secret == "" {
		return "", errors.New("获取企业微信自建应用失败")
	}
	agentID, err := strconv.Atoi(pluginInfo.AgentID)
	if err != nil {
		return "", err
	}
	// 组装content消息体
	var contentItems []qyWork.HorizontalContent
	for _, item := range request.Content {
		var contentItem qyWork.HorizontalContent
		_ = mapstructure.Decode(item, &contentItem)
		contentItems = append(contentItems, contentItem)
	}
	if len(request.MemUserIds) == 0 && len(request.PartyIds) == 0 {
		return "", errors.New("员工id和部门id不能都为空")
	}
	if len(request.MemUserIds) > 1000 || len(request.PartyIds) > 100 {
		return "", errors.New("员工id不能超过1000|部门Id不能超过100")
	}
	sendMessageBase := qyWork.SendMessageBase{
		ToUser:  strings.Join(request.MemUserIds, "|"),
		ToParty: strings.Join(request.PartyIds, "|"),
		MsgType: "template_card",
	}
	messageReq := qyWork.SendTemplateCardMessageReq{
		SendMessageBase: sendMessageBase,
		AgentId:         agentID,
		TemplateCard: qyWork.TemplateCard{
			CardType: request.CardType,
			MainTitle: &qyWork.MainTitle{
				Title: request.Title,
			},
			HorizontalContentList: &contentItems,
			JumpList:              &request.JumpList,
			CardAction:            &request.CardAction,
		},
		EnableIdTrans:        0,
		EnableDuplicateCheck: 0,
	}

	messageResp, err := workClient.SendTemplateCardMessage(c, identifier, messageReq)
	if err != nil {
		s.Log().Error(c).Msgf("发送模板卡片消息失败")
		return "", err
	}

	return messageResp.MsgId, err
}
