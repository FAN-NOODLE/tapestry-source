package customer

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/modules/enterprise"
	"micro-scrm/internal/pkg/errs"
	"micro-scrm/internal/pkg/sdk/boss"
	"micro-scrm/internal/pkg/sdk/work"
	"strings"
)

// Service 应用服务
type MessageService struct {
	repo Repository
}

//
func (m MessageService) AddMessage(c ctx.Context, data *DelMessage) (int64, error) {
	_, err := m.repo.Add(c, data)
	if err != nil {
		return 0, err
	}
	return data.ID, nil
}

func (m MessageService) UpdateMessage(c ctx.Context, data *UpdateMessageRequest) (int64, error) {
	res, err := m.repo.Update(c, data)
	if err != nil {
		return 0, err
	}
	return res, err
}

func (m MessageService) MessageDetail(c ctx.Context, id int64) (*DelMessage, error) {
	detail, err := m.repo.Detail(c, id)
	if err != nil {
		return nil, err
	}
	return detail.(*DelMessage), nil
}

func (m MessageService) MessageList(c ctx.Context, data *MessageListRequest) (*MessageListResponse, error) {
	list, err := m.repo.List(c, data)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return list.(*MessageListResponse), nil
}

func (m MessageService) SendMessage(c ctx.Context, data *SendMessageRequest, enterpriseData *enterprise.Corp) (bool, error) {

	req := boss.WorkPlugin{
		EnterpriseID: data.EnterpriseID,
	}

	// 获取自建应用信息
	pluginInfo, errInfo := boss.WorkPluginInfo(c, req)
	if errInfo != nil {
		return false, errs.WorkPluginErr
	}

	if pluginInfo.Data.AgentID == 0 {
		return false, errs.PluginInfoRrr
	}

	// 获取发消息的人员配置
	MessageParam := &MessageListRequest{
		EnterpriseID: data.EnterpriseID,
		SendType:     data.Type,
	}
	MessageData, errMem := m.MessageList(c, MessageParam)
	if errMem != nil {
		return false, errMem
	}

	if len(*MessageData) == 0 {
		log.Info(c).Msgf("企业%d下未配置删除发送信息", data.EnterpriseID)
		return true, nil
	}

	messMem := &MessageMember{}
	for _, sendData := range *MessageData {
		jsonErr := json.Unmarshal([]byte(sendData.QyStatTxt), messMem)
		if jsonErr != nil {
			return false, jsonErr
		}
		if sendData.Status == 2 {
			log.Info(c).Msgf("企业%d下配置删除发送信息,未开启", data.EnterpriseID)
			return true, nil
		}
		break
	}
	// 被删除的员工或者执行删除操作的员工
	messMem.Member = append(messMem.Member, data.DeletedMemberUserID)
	if len(messMem.Member) == 0 && len(messMem.Department) == 0 {
		log.Info(c).Msgf("企业%d未配置删除消息发送人群", data.EnterpriseID)
		return true, nil
	}
	messUser := strings.Join(messMem.Member, "|")
	messDepartment := strings.Join(messMem.Department, "|")
	// messContent := work.MessageTypeTextCardContent{
	//	Title:       "被删除提醒",
	//	Description: fmt.Sprintf(DelMsgByCustomerStr, data.UserName, data.MemberName),
	//	URL:         CrsDomain,
	//	Btntxt:      "详情",
	// }
	messContent := work.TextData{
		Content: fmt.Sprintf(DelMsgByCustomerStr, data.UserName, data.MemberName),
	}
	if data.Type == SendTypeFromMember {
		// messContent = work.MessageTypeTextCardContent{
		//	Title:       "删除提醒",
		//	Description: fmt.Sprintf(DelMsgByMemberStr, data.MemberName, data.UserName),
		//	URL:         CrsDomain,
		//	Btntxt:      "详情",
		// }
		messContent = work.TextData{
			Content: fmt.Sprintf(DelMsgByMemberStr, data.MemberName, data.UserName),
		}
	}

	workImpl := new(work.QyWorkImpl).SetDomain(enterpriseData.Domain)
	deptReq := work.SendMessageRequest{
		Touser:                 messUser,
		Toparty:                messDepartment,
		Text:                   &messContent,
		Msgtype:                "text",
		Agentid:                pluginInfo.Data.AgentID,
		DuplicateCheckInterval: 1800,
		QyGwRequest: work.QyGwRequest{
			ReqPL: pluginInfo.Data.Identifier,
		},
	}
	res, errSend := workImpl.SendMessage(c, deptReq)
	if errSend != nil {
		log.Info(c).Msgf("企业%d下的发送信息应用[%s]的成员出错, 错误信息：%+v", data.EnterpriseID, pluginInfo.Data.Identifier, errSend)
		return false, errSend
	}
	if res.Errcode == 0 && res.Errmsg == "ok" {
		return true, nil
	}
	return true, nil
}
