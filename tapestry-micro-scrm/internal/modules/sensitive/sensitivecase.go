package sensitive

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/modules/members"
	membersRepository "micro-scrm/internal/modules/members/respository/mysql"
	"micro-scrm/internal/modules/remindconfig"
	remindconfigRepository "micro-scrm/internal/modules/remindconfig/respository/mysql"
	"micro-scrm/internal/pkg/errs"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sdk/boss"
	"micro-scrm/internal/pkg/sdk/work"
	"micro-scrm/internal/pkg/uid"
	"strings"
)

//
type Service struct {
	repo Repository
}

// SelectList 敏感词列表（分页）
func (s Service) SelectList(c ctx.Context, request ListRequest) (int64, *[]Detail, error) {

	total, list, err := s.repo.SelectList(c, request)
	if err != nil {
		return 0, nil, err
	}
	if total != 0 {
		var ids []int64

		for _, detail := range *list {
			ids = append(ids, detail.ID)
		}

		countBySensitive := &[]CountBySensitive{}

		countBySensitive, err = s.repo.CountBySensitive(c, ids)
		if err != nil {
			return 0, nil, err
		}
		if countBySensitive != nil {
			countById := make(map[int64]CountBySensitive)
			for _, count := range *countBySensitive {
				countById[count.SensitiveID] = count
			}

			for key, value := range *list {
				(*list)[key].TodayTriggerNum = countById[value.ID].TodayTriggerNum
			}
		}
	}

	return total, list, err
}

// SensitiveSave 添加敏感词
func (s Service) SensitiveSave(c ctx.Context, request SaveRequest) error {
	res, _ := s.repo.FindSensitive(c, request)

	if res != nil {
		if request.SensitiveType == TypeCart {
			return errs.SensitiveCartIsExistsErr
		}
		if request.SensitiveType == TypeText {
			return errs.SensitiveTextRepeatErr
		}
	}

	_, err := s.repo.Add(c, request)
	if err != nil {
		return err
	}
	return nil
}

// SensitiveDelete 删除敏感词
func (s Service) SensitiveDelete(c ctx.Context, request DeleteRequest) error {
	_, err := s.repo.DeleteSensitive(c, request)
	if err != nil {
		return err
	}

	_, err = s.repo.DeleteSensitiveRecords(c, request)
	if err != nil {
		return err
	}

	return nil
}

// SensitiveDetail 敏感词详情
func (s Service) SensitiveDetail(c ctx.Context, request DetailListRequest) (*Detail, error) {
	var detailRes *Detail

	res, err := s.repo.Detail(c, request.ID)
	if err != nil {
		return detailRes, err
	}

	detail := res.(*Detail)

	ids := []int64{detail.ID}

	countBySensitive := &[]CountBySensitive{}

	countBySensitive, err = s.repo.CountBySensitive(c, ids)
	if err != nil {
		return detailRes, err
	}

	if len(*countBySensitive) != 0 {
		(*detail).TodayTriggerNum = (*countBySensitive)[0].TodayTriggerNum
	}

	return detail, nil
}

// SensitiveDetailList 敏感词明细列表（分页）
func (s Service) SensitiveDetailList(c ctx.Context, request DetailListRequest) (int64, *[]DetailListResponse, error) {

	total, list, err := s.repo.SelectRecordsList(c, request)
	if err != nil || total == 0 {
		return 0, nil, err
	}

	var ids []int64
	detailById := make(map[int64]Records)

	for _, detail := range *list {
		ids = append(ids, detail.ChatMsgID)
		detailById[detail.ChatMsgID] = detail
	}
	detailList, err := s.SelectChatMsgList(c, ids)
	if err != nil || len(*detailList) == 0 {
		return 0, nil, err
	}

	for key, value := range *detailList {
		(*detailList)[key].ID = detailById[value.ChatMsgID].ID
		(*detailList)[key].CreatedAt = detailById[value.ChatMsgID].CreatedAt
	}

	return total, detailList, err
}

// SelectChatMsgList 敏感词对应聊天记录列表
func (s Service) SelectChatMsgList(c ctx.Context, ids []int64) (*[]DetailListResponse, error) {
	list, err := s.repo.SelectChatMsgList(c, ids)
	if err != nil {
		return nil, err
	}

	var result []DetailListResponse
	if len(*list) != 0 {
		for _, value := range *list {
			detail := DetailListResponse{}
			if value.UserChannel == UserChannelMember {
				memberUseridArr := []string{value.FromUser, value.ToUser}
				memberInfo, err := s.repo.GetMembersList(c, memberUseridArr, value.EnterpriseID)
				if err == nil && memberInfo != nil {
					for _, member := range *memberInfo {
						if member.UserID == value.FromUser {
							detail.MemberName = member.Name
							detail.MemberAvatar = member.Avatar
						} else {
							detail.EXternalUserName = member.Name
							detail.EXternalUserAvatar = member.Avatar
						}
					}
				}
			} else {
				releaionRequest := ChatMsgMembers{
					EnterpriseID:   value.EnterpriseID,
					MemberUserID:   value.FromUser,
					EXternalUserID: value.ToUser,
				}
				relation, err := s.repo.SelectUserRelation(c, releaionRequest)
				if err == nil && relation != nil {
					detail = *relation
				}
			}

			detail.ChatMsgID = value.ID
			detail.EnterpriseID = value.EnterpriseID
			detail.MemberUserID = value.FromUser
			detail.EXternalUserID = value.ToUser

			result = append(result, detail)
		}
	}

	return &result, nil
}

// SensitiveMatching 敏感词检测
func (s Service) SensitiveMatching(c ctx.Context, data *SendMessageRequest, domain string) (bool, error) {

	// 获取敏感词列表
	sensitiveRequest := ListRequest{
		EnterpriseID: data.EnterpriseID,
		Page:         1,
		Size:         99999,
	}
	total, list, _ := s.SelectList(c, sensitiveRequest)
	if total == 0 {
		log.Info(c).Msgf("企业%d下未配置敏感词", data.EnterpriseID)
		return true, nil
	}

	// 获取敏感词提醒配置信息
	remindService := remindconfig.New(remindconfigRepository.New())
	membersService := members.New(membersRepository.New())

	remindRequest := remindconfig.ConfigInfoRequest{
		EnterpriseID: data.EnterpriseID,
		NoticeType:   remindconfig.SensitiveNotice,
	}

	messUser := ""

	remindInfo, _ := remindService.ConfigInfo(c, &remindRequest, membersService)
	if remindInfo == nil {
		log.Info(c).Msgf("企业%d下未配置敏感词提醒设置", data.EnterpriseID)
		remindInfo = &remindconfig.ConfigInfoResponse{
			NoticeStatus: remindconfig.NoticeOpen,
		}
	}

	var userStrIds []string
	for _, v := range remindInfo.StatList {
		if v.UserId == data.FromUser && remindInfo.NoticeStatus == remindconfig.NoticeClose {
			continue
		}
		userStrIds = append(userStrIds, v.UserId)
	}

	if remindInfo.NoticeStatus == remindconfig.NoticeOpen && helper.InArray[string](data.FromUser, userStrIds) == false {
		userStrIds = append(userStrIds, data.FromUser)
	}
	messUser = strings.Join(userStrIds, "|")

	log.Info(c).Msgf("debug7:%v", messUser)

	workImpl := new(work.QyWorkImpl).SetDomain(domain)

	// 匹配敏感词
	for _, value := range *list {
		if (value.SensitiveType == TypeCart && data.Type == "card") || (value.SensitiveType == TypeText && strings.Contains(data.ChatRecord, value.SensitiveText)) {
			if messUser != "" {
				req := boss.WorkPlugin{
					EnterpriseID: data.EnterpriseID,
				}

				// 获取自建应用信息
				pluginInfo, errInfo := boss.WorkPluginInfo(c, req)
				log.Info(c).Msgf("debug2:%v", pluginInfo)
				log.Info(c).Msgf("debug3:%v", errInfo)
				flag := true
				if errInfo != nil || pluginInfo.Data.AgentID == 0 {
					flag = false
				}

				releaionRequest := ChatMsgMembers{
					EnterpriseID:   data.EnterpriseID,
					MemberUserID:   data.FromUser,
					EXternalUserID: data.ToUser,
				}

				// 获取外部联系人和员工姓名
				relation, err := s.repo.SelectUserRelation(c, releaionRequest)
				if err != nil {
					log.Info(c).Msgf("企业%d下敏感词提醒未找到好友关系", data.EnterpriseID)
					flag = false
				}

				if flag {
					// 发消息
					// messContent := work.MessageTypeTextCardContent{
					//	Title:       "敏感词触发提醒",
					//	Description: "成员" + relation.MemberName + "与客户" + relation.EXternalUserName + "的聊天有敏感信息，请注意敏感行为。",
					//	URL:         customer.CrsDomain,
					//	Btntxt:      "详情",
					// }
					messContent := work.TextData{
						Content: "成员" + relation.MemberName + "与客户" + relation.EXternalUserName + "的聊天有敏感信息，请注意敏感行为。",
					}

					sendReq := work.SendMessageRequest{
						Touser:                 messUser,
						Msgtype:                "text",
						Text:                   &messContent,
						Agentid:                pluginInfo.Data.AgentID,
						DuplicateCheckInterval: 1800,
						QyGwRequest: work.QyGwRequest{
							ReqPL: pluginInfo.Data.Identifier,
						},
					}
					_, errSend := workImpl.SendMessage(c, sendReq)
					if errSend != nil {
						log.Info(c).Msgf("企业%d下的发送敏感词提醒成员出错, 错误信息：%+v", data.EnterpriseID, errSend)
					}
				}
			}

			insertData := Records{
				ID:           uid.ID(c),
				EnterpriseID: data.EnterpriseID,
				ChatMsgID:    data.ChatMessageID,
				SensitiveID:  value.ID,
			}
			eff, _ := s.repo.InsertSensitiveRecords(c, insertData)
			if eff != 0 {
				sensitiveData := SaveRequest{
					ID: value.ID,
				}
				eff, _ = s.repo.Update(c, sensitiveData)
			}
		}
	}

	return true, nil
}
