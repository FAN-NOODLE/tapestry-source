package api

import (
	"encoding/json"
	"errors"
	url2 "net/url"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/modules/groupmsg"
	"micro-scrm/internal/modules/ws_trace"
	"micro-scrm/internal/pkg/errs"
	"micro-scrm/internal/pkg/sdk/work"
)

type WorkApi struct {
}

var workSdk work.Api = &work.QyWorkImpl{}

// SendGroupMsg 发送群发消息
func (w *WorkApi) SendGroupMsg(ctx ctx.Context, sender string, content string, chatType string, domain string, WsTraceService *ws_trace.Service, memberID int64, taskID int64, enterpriseID int64) (string, error) {
	workSdk.SetDomainValue(domain)
	var massHairReq = work.SendMsgTemplateReq{}
	err := json.Unmarshal([]byte(content), &massHairReq)
	if err != nil {
		return "", errs.GroupMsgContentParseErr
	}
	if massHairReq.Text == nil && len(massHairReq.Attachments) == 0 {
		return "", errs.GroupMsgContentEmptyErr
	}

	for _, mass := range massHairReq.Attachments {
		if mass.MiniProgram != nil {
			// 解析小程序链接
			url, err := url2.Parse(mass.MiniProgram.Page)
			if err != nil {
				log.Info(ctx).Msgf("解析小程序page链接失败, err:%s, page:%s", err.Error(), mass.MiniProgram.Page)
				return content, nil
			}

			// 生成员工参数
			var wsTraceID string
			wst := ws_trace.WsTrace{
				MemberID:      memberID,
				QwGroupSendID: taskID,
				SourceType:    ws_trace.SourceTypeGroupChatSend,
			}
			wsTrace, err := WsTraceService.WsAdd(ctx, enterpriseID, wst)
			if err != nil {
				log.Err(ctx, err).Msgf("生成ws_trace_id错误:%+v", wst)
			}
			if wsTrace != nil {
				wsTraceID = wsTrace.WsTraceID
			}
			// url中拼接参数
			if url.RawQuery == "" {
				url.RawQuery = "ws_trace_id=" + wsTraceID + "&fromOrigin=ws"
			} else {
				url.RawQuery = url.RawQuery + "&ws_trace_id=" + wsTraceID + "&fromOrigin=ws"
			}
			mass.MiniProgram.Page = url.String()
		}
	}

	massHairReq.ChatType = chatType
	massHairReq.Sender = sender
	response, err := workSdk.SendMsgTemplate(ctx, massHairReq)
	if err != nil {
		log.Info(ctx).Msgf("发送群发消息失败, err:%s", err.Error())
		return "", err
	}
	if response.ErrCode != 0 {
		log.Info(ctx).Msgf("发送群发消息企微返回错误, response:%v", response)
		return "", errorx.NewErrorMsg(response.ErrCode, response.ErrMsg)
	}
	return response.MsgID, nil
}

// GroupMsgSendResult 查询群发结果
func (w *WorkApi) GroupMsgSendResult(c ctx.Context, sender string, msgID string, domain string) ([]*groupmsg.ChatGroupMsgSendRecord, error) {
	if msgID == "" {
		return nil, errors.New("消息Id为空")
	}

	workSdk.SetDomainValue(domain)
	var limit = 500
	var request = work.GroupMsgSendResultRequest{
		MsgID:  msgID,
		UserID: sender,
		Limit:  limit,
	}
	var responseList []*work.GroupMsgSendResultItem
	// 循环调用企微接口查询msgid下的所有发送结果
	for {
		response, err := workSdk.GetGroupMsgSendResult(c, request)
		if err != nil {
			return nil, err
		}
		if response.ErrCode != 0 {
			log.Info(c).Msgf("企微查询群发结果返回错误, errCode:%d, msg:%s", response.ErrCode, response.ErrMsg)
			return nil, errorx.NewErrorMsg(response.ErrCode, response.ErrMsg)
		}
		responseList = append(responseList, response.SendList...)
		if response.NextCursor == "" {
			break
		}
		request.Cursor = response.NextCursor
	}
	var resultList []*groupmsg.ChatGroupMsgSendRecord
	for _, v := range responseList {
		resultList = append(resultList, &groupmsg.ChatGroupMsgSendRecord{
			ChatId:         v.ChatID,
			Status:         v.Status,
			SendTime:       v.SendTime,
			ExternalUserID: v.ExternalUserID,
		})
	}
	return resultList, nil
}

func (w *WorkApi) SetMediaID(c ctx.Context, content string, enterpriseID int64, domain string) (string, error) {
	workSdk.SetDomainValue(domain)
	return workSdk.SetMediaID(c, content, enterpriseID)
}
