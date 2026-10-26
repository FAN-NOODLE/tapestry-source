package overtime

import (
	"encoding/json"
	"errors"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/modules/chat"
	chatRepository "micro-scrm/internal/modules/chat/repository/mysql"
	"micro-scrm/internal/modules/enterprise"
	enterpriseRepository "micro-scrm/internal/modules/enterprise/respository/mysql"
	"micro-scrm/internal/modules/members"
	membersRepository "micro-scrm/internal/modules/members/respository/mysql"
	"micro-scrm/internal/modules/remindconfig"
	remindconfigRepository "micro-scrm/internal/modules/remindconfig/respository/mysql"
	"micro-scrm/internal/modules/settings"
	settingRepository "micro-scrm/internal/modules/settings/respository/mysql"
	"micro-scrm/internal/pkg/errorIgnore"
	"micro-scrm/internal/pkg/sdk/boss"
	"micro-scrm/internal/pkg/sdk/work"
	"micro-scrm/internal/pkg/uid"
	"strconv"
	"strings"
	"time"
)

//
type Service struct {
	repo Repository
}

// SelectList 超时记录列表（分页）
func (s Service) SelectList(c ctx.Context, request ListRequest) (int64, *[]RecordsListResponse, error) {
	// 判断查询条件为员工id或者部门时
	var msgIds []int64
	if request.MemberID > 0 || request.DepartmentID > 0 {
		memlist, err := s.repo.SelectMemberByWhere(c, request)
		if err != nil {
			return 0, nil, err
		}
		var memberUserIds []string
		for _, mem := range *memlist {
			memberUserIds = append(memberUserIds, mem.MemberUserId)
		}
		receiveService := chat.NewReceive(chatRepository.NewChatReceiveUser())
		if memberUserIds != nil {
			// 查询接受不了表里to_user等于memberUserIds
			msgIDs, err := receiveService.SelectMsgIDByUserIds(c, memberUserIds)
			if err != nil {
				return 0, nil, err
			}
			for _, msgID := range msgIDs {
				msgIds = append(msgIds, msgID.ChatMsgId)
			}
		}
		if msgIds == nil {
			return 0, nil, err
		}
	}
	request.ChatMsgIDsNum = 0
	if msgIds != nil {
		request.ChatMsgIDs = msgIds
		request.ChatMsgIDsNum = 1
	}
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

	msgsList, err := s.SelectChatMsgList(c, ids, detailById)
	if err != nil || len(*msgsList) == 0 {
		return 0, nil, err
	}
	return total, msgsList, err
}

// SelectChatMsgList 超时对应聊天记录列表
func (s Service) SelectChatMsgList(c ctx.Context, ids []int64, detailById map[int64]Records) (*[]RecordsListResponse, error) {

	list, err := s.repo.SelectChatMsgList(c, ids)
	if err != nil {
		return nil, err
	}

	var result []RecordsListResponse
	if len(*list) != 0 {
		for _, value := range *list {
			var relation = &ChatMsgMembers{}
			// 发送方为员工
			if value.FromUserType == FromTypeMember {
				var releaionRequest = []string{
					value.FromUser,
					value.ToUser,
				}
				memberList, err := s.repo.SelectMembersList(c, releaionRequest, value.EnterpriseID)
				if errorIgnore.Ignore(err) != nil {
					log.Info(c).Msgf("查询员工与员工关系失败Request: %+v", releaionRequest)
					log.Info(c).Msgf("查询员工与员工关系失败Err: %+v", err)
					continue
				}
				memberById := make(map[string]MemberInfo)
				for _, detail := range *memberList {
					memberById[detail.UserID] = detail
				}
				if fromMember, ok := memberById[value.FromUser]; ok {
					relation.EXternalUserID = fromMember.UserID
					relation.EXternalUserName = fromMember.Name
					relation.EXternalUserAvatar = fromMember.Avatar
					relation.EXternalUserRemark = ""
				}
				if toUserMember, ok := memberById[value.ToUser]; ok {
					relation.MemberID = toUserMember.MemberID
					relation.MemberUserID = toUserMember.UserID
					relation.MemberName = toUserMember.Name
					relation.MemberAvatar = toUserMember.Avatar
				}
			} else {
				var releaionRequest = ChatMsgMembers{
					EnterpriseID:   value.EnterpriseID,
					MemberUserID:   value.ToUser,
					EXternalUserID: value.FromUser,
				}
				relation, err = s.repo.SelectUserRelation(c, releaionRequest)
				if errorIgnore.Ignore(err) != nil || relation == nil {
					relation = &releaionRequest
					log.Info(c).Msgf("查询员工与客户关系失败Request: %+v", releaionRequest)
					log.Info(c).Msgf("查询员工与客户关系失败Err: %+v", err)
					continue
				}
			}
			detail := RecordsListResponse{}
			detail.ChatMsgID = value.ID
			detail.ID = value.ID
			contentMap := make(map[string]interface{})
			_ = json.Unmarshal([]byte(value.Content), &contentMap)
			detail.Content = value.Content
			detail.ContentMap = contentMap
			detail.ContentType = value.ContentType
			detail.MemberID = relation.MemberID
			detail.MemberUserID = relation.MemberUserID
			detail.MemberName = relation.MemberName
			detail.MemberAvatar = relation.MemberAvatar
			detail.EXternalUserID = relation.EXternalUserID
			detail.EXternalUserName = relation.EXternalUserName
			detail.EXternalUserAvatar = relation.EXternalUserAvatar
			detail.EXternalUserRemark = relation.EXternalUserRemark
			detail.MsgAt = detailById[value.ID].CreatedAt
			// 所属部门
			memberIDs := []int64{relation.MemberID}
			belongToDepartments, _ := s.repo.MemberMultiDepartments(c, memberIDs)
			if departments, ok := belongToDepartments[relation.MemberID]; ok {
				detail.MemberDepartsList = departments
			}

			result = append(result, detail)
		}
	}

	return &result, nil
}

// 获取超时消息
func (s Service) GetOvertime(c ctx.Context) {
	log.Info(c).Msgf("GetOvertime-start: %+v", "GetOvertime-start")
	// 遍历chat_settings表的企业开始处理
	settingService := settings.New(settingRepository.New())
	for _, enterpriseSetting := range settingService.EnterpriseList(c) {
		log.Info(c).Msgf("Start-EnterpriseID: %+v", enterpriseSetting.EnterpriseId)
		EnterpriseID := enterpriseSetting.EnterpriseId
		// 查询该企业超时配置
		remindService := remindconfig.New(remindconfigRepository.New())
		membersService := members.New(membersRepository.New())
		remindRequest := remindconfig.ConfigInfoRequest{
			EnterpriseID: EnterpriseID,
			NoticeType:   remindconfig.OvertimeNotice,
		}
		remindInfo, _ := remindService.ConfigInfo(c, &remindRequest, membersService)
		// 当超时提醒配置信息为空时，为默认配置
		if remindInfo == nil {
			remindInfo = &remindconfig.ConfigInfoResponse{
				OverTime: int32(remindconfig.DefaultOverTime),
			}
		}
		// 获取接收方为员工的消息记录
		toMemberCond := ChatMsg{
			EnterpriseID: EnterpriseID,
			UserChannel:  ReceiveTypeMember,
		}
		sendToMemList, err := s.repo.SelectReplyMsg(c, toMemberCond)
		if err != nil {
			log.Info(c).Msgf("sendToMemList-err: %+v", err)
			continue
		}
		log.Info(c).Msgf("sendToMemList: %+v", *sendToMemList)

		// 查询员工是否回复客户
		for _, value := range *sendToMemList {
			toUserCond := ChatMsg{
				FromUser:     value.ToUser,
				ToUser:       value.FromUser,
				TimeStart:    value.TimeStart,
				TimeEnd:      value.TimeStart + int64(remindInfo.OverTime*60),
				EnterpriseID: EnterpriseID,
			}
			sendToUserList, _ := s.repo.SelectReplyMsg(c, toUserCond)
			log.Info(c).Msgf("value: %+v", value)
			log.Info(c).Msgf("toUserCond: %+v", toUserCond)
			log.Info(c).Msgf("sendToUserList: %+v", *sendToUserList)
			if len(*sendToUserList) == 0 {
				// 判断超时时间是否大于当前时间
				if toUserCond.TimeEnd-time.Now().Unix() > 0 {
					continue
				}
				// 查询该条消息是否已记录到表中
				request := Records{
					EnterpriseID: EnterpriseID,
					ChatMsgID:    value.ID,
				}
				isExists, _ := s.repo.FindOvertime(c, request)

				if isExists == nil {
					// 未回复消息，加入超时记录
					overTimeRecord := AddOverRecord{
						ID:           uid.ID(ctx.Wrap(c)),
						ChatMsgID:    value.ID,
						EnterpriseID: EnterpriseID,
					}
					_, errAdd := s.repo.Add(c, overTimeRecord)
					if errorIgnore.Ignore(errAdd) != nil {
						continue
					}
					// 发送消息提醒
					sendMsgRequest := SendMessageRequest{
						EnterpriseID: EnterpriseID,
						FromUser:     value.ToUser,
						ToUser:       value.FromUser,
					}
					s.SendMessage(c, sendMsgRequest)
				}
			}
		}
		log.Info(c).Msgf("End-EnterpriseID: %+v", enterpriseSetting.EnterpriseId)
	}
}

// 获取超时消息
func (s Service) GetOvertimeOne(c ctx.Context, enterpriseSetting settings.ChatSettings) {
	log.Info(c).Msgf("GetOvertime-start: %+v", "GetOvertime-start")
	// 遍历chat_settings表的企业开始处理
	log.Info(c).Msgf("Start-EnterpriseID: %+v", enterpriseSetting.EnterpriseId)
	EnterpriseID := enterpriseSetting.EnterpriseId
	// 查询该企业超时配置
	remindService := remindconfig.New(remindconfigRepository.New())
	membersService := members.New(membersRepository.New())
	remindRequest := remindconfig.ConfigInfoRequest{
		EnterpriseID: EnterpriseID,
		NoticeType:   remindconfig.OvertimeNotice,
	}
	remindInfo, _ := remindService.ConfigInfo(c, &remindRequest, membersService)
	// 当超时提醒配置信息为空时，为默认配置
	if remindInfo == nil {
		remindInfo = &remindconfig.ConfigInfoResponse{
			OverTime: int32(remindconfig.DefaultOverTime),
		}
	}
	// 获取接收方为员工的消息记录
	toMemberCond := ChatMsg{
		EnterpriseID: EnterpriseID,
		UserChannel:  ReceiveTypeMember,
	}
	sendToMemList, err := s.repo.SelectReplyMsg(c, toMemberCond)
	if err != nil {
		log.Info(c).Msgf("sendToMemList-err: %+v", err)
		return
	}
	log.Info(c).Msgf("sendToMemList: %+v", *sendToMemList)

	// 查询员工是否回复客户
	for _, value := range *sendToMemList {
		toUserCond := ChatMsg{
			FromUser:     value.ToUser,
			ToUser:       value.FromUser,
			TimeStart:    value.TimeStart,
			TimeEnd:      value.TimeStart + int64(remindInfo.OverTime*60),
			EnterpriseID: EnterpriseID,
		}
		sendToUserList, _ := s.repo.SelectReplyMsg(c, toUserCond)
		log.Info(c).Msgf("value: %+v", value)
		log.Info(c).Msgf("toUserCond: %+v", toUserCond)
		log.Info(c).Msgf("sendToUserList: %+v", *sendToUserList)
		if len(*sendToUserList) == 0 {
			// 查询该条消息是否已记录到表中
			request := Records{
				EnterpriseID: EnterpriseID,
				ChatMsgID:    value.ID,
			}
			isExists, _ := s.repo.FindOvertime(c, request)

			if isExists == nil {
				// 未回复消息，加入超时记录
				overTimeRecord := AddOverRecord{
					ID:           uid.ID(ctx.Wrap(c)),
					ChatMsgID:    value.ID,
					EnterpriseID: EnterpriseID,
				}
				_, errAdd := s.repo.Add(c, overTimeRecord)
				if errorIgnore.Ignore(errAdd) != nil {
					continue
				}
				// 发送消息提醒
				sendMsgRequest := SendMessageRequest{
					EnterpriseID: EnterpriseID,
					FromUser:     value.ToUser,
					ToUser:       value.FromUser,
				}
				s.SendMessage(c, sendMsgRequest)
			}
		}
	}
	log.Info(c).Msgf("End-EnterpriseID: %+v", enterpriseSetting.EnterpriseId)
}

func (s *Service) SendMessage(c ctx.Context, data SendMessageRequest) {
	enterpriseService := enterprise.New(enterpriseRepository.New())
	// 查企业信息
	enterpriseData, enterpriseInfoErr := enterpriseService.Detail(c, enterprise.Corp{
		EnterpriseID: data.EnterpriseID,
	})
	if enterpriseInfoErr != nil {
		log.Info(c).Msgf("获取企业信息错误:" + enterpriseInfoErr.Error())
		return
	}
	// 获取超时提醒配置信息
	remindService := remindconfig.New(remindconfigRepository.New())
	membersService := members.New(membersRepository.New())

	remindRequest := remindconfig.ConfigInfoRequest{
		EnterpriseID: data.EnterpriseID,
		NoticeType:   remindconfig.OvertimeNotice,
	}
	remindInfo, _ := remindService.ConfigInfo(c, &remindRequest, membersService)
	if remindInfo == nil {
		log.Info(c).Msgf("企业%d下未配置超时提醒设置", data.EnterpriseID)
		remindInfo = &remindconfig.ConfigInfoResponse{
			NoticeStatus: remindconfig.NoticeOpen,
			ReceiveStart: remindconfig.DefaultReceiveStart,
			ReceiveEnd:   remindconfig.DefaultReceiveEnd,
		}
	}
	log.Info(c).Msgf("remindInfo: %+v", remindInfo)
	// 转换时间戳
	cstSh, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Info(ctx.Wrap(c)).Msgf("LoadLocation:%s", err.Error())
		cstSh = time.Local
	}
	startTime := strings.Split(remindInfo.ReceiveStart, ":")
	endTime := strings.Split(remindInfo.ReceiveEnd, ":")
	timeNow := time.Now()
	startH, err := strconv.Atoi(startTime[0]) // 开始时间-时
	startM, err := strconv.Atoi(startTime[1]) // 开始时间-分
	endH, err := strconv.Atoi(endTime[0])     // 结束时间-时
	endM, err := strconv.Atoi(endTime[1])     // 结束时间-分
	log.Info(c).Msgf("time: %+v", time.Now())
	startTimeData := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), startH, startM, 0, 0, cstSh)
	endTimeData := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), endH, endM, 0, 0, cstSh)
	log.Info(c).Msgf("startTimeData: %+v", startTimeData)

	ReceiveStart := startTimeData.Unix()
	log.Info(c).Msgf("ReceiveStart: %+v", ReceiveStart)
	ReceiveEnd := endTimeData.Unix()
	log.Info(c).Msgf("ReceiveEnd: %+v", ReceiveEnd)
	nowTime := timeNow.Unix()
	log.Info(c).Msgf("nowTime: %+v", nowTime)
	messUser := ""
	if len(remindInfo.StatList) > 0 && nowTime > ReceiveStart && nowTime < ReceiveEnd {
		log.Info(c).Msgf("sendStatus: %+v", "ok")
		var userStrIds []string
		for _, v := range remindInfo.StatList {
			userStrIds = append(userStrIds, v.UserId)
		}
		log.Info(c).Msgf("userStrIds: %+v", userStrIds)
		messUser = strings.Join(userStrIds, "|")
		if remindInfo.NoticeStatus == remindconfig.NoticeOpen {
			if !strings.Contains("|"+messUser+"|", "|"+data.FromUser+"|") {
				messUser = messUser + "|" + data.FromUser
			}
		}
	} else if remindInfo.NoticeStatus == remindconfig.NoticeOpen {
		messUser = data.FromUser
	}

	log.Info(c).Msgf("messUser: %+v", messUser)
	if messUser == "" {
		return
	}
	req := boss.WorkPlugin{
		EnterpriseID: data.EnterpriseID,
	}
	// 获取自建应用信息
	pluginInfo, errInfo := boss.WorkPluginInfo(c, req)
	if errInfo != nil {
		return
	}

	if pluginInfo.Data.AgentID == 0 {
		return
	}

	// 获取外部联系人和员工姓名
	releaionRequest := ChatMsgMembers{
		EnterpriseID:   data.EnterpriseID,
		MemberUserID:   data.FromUser,
		EXternalUserID: data.ToUser,
	}

	relation, err := s.repo.SelectUserRelation(c, releaionRequest)
	if err != nil {
		log.Info(c).Msgf("企业%d下超时提醒未找到好友关系", data.EnterpriseID)
		return
	}
	workImpl := new(work.QyWorkImpl).SetDomain(enterpriseData.Domain)
	// 发消息
	// messContent := work.MessageTypeTextCardContent{
	//	Title:       "聊天超时提醒",
	//	Description: "成员" + relation.MemberName + "与客户" + relation.EXternalUserName + "的聊天已超时，请及时回复。",
	//	URL:         customer.CrsDomain,
	//	Btntxt:      "详情",
	// }
	messContent := work.TextData{
		Content: "成员" + relation.MemberName + "与客户" + relation.EXternalUserName + "的聊天已超时，请及时回复。",
	}

	sendReq := work.SendMessageRequest{
		Touser: messUser,
		// Msgtype:                customer.MessageType[customer.MessageTypeTextCard],
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
		log.Info(c).Msgf("企业%d下的发送聊天超时提醒成员出错, 错误信息：%+v", data.EnterpriseID, errSend)
	}
}

// GetConfig 根据企业ID获取需要的配置信息
func (s Service) GetConfig(c ctx.Context, EnterpriseID int64) (*enterprise.Corp, *settings.ChatSettings, error) {
	var enterpriseInfo *enterprise.Corp
	var settingInfo *settings.ChatSettings
	var err error
	enterpriseService := enterprise.New(enterpriseRepository.New())
	settingService := settings.New(settingRepository.New())
	// 查企业信息
	enterpriseInfo, err = enterpriseService.Detail(c, enterprise.Corp{
		EnterpriseID: EnterpriseID,
	})
	if err != nil {
		return enterpriseInfo, settingInfo, errors.New("获取企业信息错误:" + err.Error())
	}

	// 查私钥配置
	settingInfo, err = settingService.SettingsInfo(c, &settings.SettingsInfoRequest{EnterpriseID: EnterpriseID})
	if err != nil {
		return enterpriseInfo, settingInfo, errors.New("获取企业私钥错误:" + err.Error())
	}

	return enterpriseInfo, settingInfo, nil
}
