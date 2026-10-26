package chat

import (
	"encoding/json"
	"errors"
	"fmt"
	"micro-scrm/internal/pkg/helper"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/lock"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/typex"
	"micro-scrm/internal/config"
	"micro-scrm/internal/modules/callback"
	"micro-scrm/internal/modules/customer"
	userEs "micro-scrm/internal/modules/customer/respository/elastic"
	"micro-scrm/internal/modules/customer/respository/mysql"
	"micro-scrm/internal/modules/enterprise"
	enterpriseRepository "micro-scrm/internal/modules/enterprise/respository/mysql"
	"micro-scrm/internal/modules/sensitive"
	sensitiveRepository "micro-scrm/internal/modules/sensitive/respository/mysql"
	"micro-scrm/internal/modules/settings"
	settingRepository "micro-scrm/internal/modules/settings/respository/mysql"
	"micro-scrm/internal/pkg/errorIgnore"
	chat "micro-scrm/internal/pkg/sdk/microchat"
	"micro-scrm/internal/pkg/sdk/work"
	"micro-scrm/internal/pkg/uid"
)

// Service 聊天综合服务
type Service struct {
	repo MsgRepository
}

// MsgTypeList 消息类型列表
func (s Service) MsgTypeList() []string {
	return []string{
		chat.MsgTypeText,
		chat.MsgTypeImage,
		chat.MsgTypeRevoke,
		chat.MsgTypeAgree,
		chat.MsgTypeVoice,
		chat.MsgTypeVideo,
		chat.MsgTypeCard,
		chat.MsgTypeLocation,
		chat.MsgTypeEmotion,
		chat.MsgTypeFile,
		chat.MsgTypeLink,
		chat.MsgTypeWeapp,
		// chat.MsgTypeChatrecord,
		chat.MsgTypeTodo,
		chat.MsgTypeVote,
		chat.MsgTypeCollect,
		chat.MsgTypeRedpacket,
		chat.MsgTypeMeeting,
		chat.MsgTypeDocmsg,
		chat.MsgTypeMarkdown,
		chat.MsgTypeInfo,
		chat.MsgTypeCalendar,
		chat.MsgTypeMixed,
		chat.MsgTypeMeetingVoiceCall,
		chat.MsgTypeVoipDocShare,
	}
}

// MsgPageList 消息分页列表
func (s Service) MsgPageList(c ctx.Context, mediaService MsgMediaService, cond MsgPageListRequest) (MsgPageListResponse, error) {
	var cmp MsgPageListResponse

	// 搜的不是文本消息
	if cond.Type != "" && cond.Type != chat.MsgTypeText && cond.Keywords != "" {
		var l = make([]Msg, 0)
		cmp.List = l
		cmp.Size = cond.Size
		cmp.Page = cond.Page

		return cmp, nil
	}

	offset := (cond.Page - 1) * cond.Size
	pageVo := repos.PageVO{Offset: int(offset), Limit: int(cond.Size)}
	option := &repos.FindPageOption{
		TplPrefix: "t_chat_msg/list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]Msg{}
		},
	}
	total, list, err := s.repo.FindPage(c, pageVo, cond, option)
	if errorIgnore.Ignore(err) != nil {
		return cmp, err
	}

	cmp.Page = cond.Page
	cmp.Size = cond.Size
	cmp.List = s.FormatMsgList(c, mediaService, *list.(*[]Msg))
	cmp.Total = total

	return cmp, nil
}

// GroupMsgPageList 群消息分页列表
func (s Service) GroupMsgPageList(c ctx.Context, mediaService MsgMediaService, cond GroupMsgListRequest) (MsgPageListResponse, error) {
	var cmp MsgPageListResponse

	offset := (cond.Page - 1) * cond.Size
	pageVo := repos.PageVO{Offset: int(offset), Limit: int(cond.Size)}
	option := &repos.FindPageOption{
		TplPrefix: "t_chat_msg/group-list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]Msg{}
		},
	}
	total, list, err := s.repo.FindPage(c, pageVo, cond, option)
	if errorIgnore.Ignore(err) != nil {
		return cmp, err
	}

	cmp.Page = cond.Page
	cmp.Size = cond.Size
	cmp.List = s.FormatMsgList(c, mediaService, *list.(*[]Msg))
	cmp.Total = total

	return cmp, nil
}

// Detail 消息详情
func (s Service) Detail(c ctx.Context, id int64) (Msg, error) {
	var msg Msg
	option := &repos.DetailOption{
		TplPrefix: "t_chat_msg/detail",
		Fn: func(ctx.Context) interface{} {
			return &Msg{}
		},
	}
	thisMsg, err := s.repo.DetailPlus(c, id, option)
	if err != nil {
		return msg, err
	}

	msg = *thisMsg.(*Msg)

	return msg, nil
}

// MsgListByID 根据消息ID找附近的消息
func (s Service) MsgListByID(c ctx.Context, cond MsgListByIDRequest, mediaService MsgMediaService) (MsgListByIDResponse, error) {
	var res MsgListByIDResponse

	list, err := s.repo.ListByID(c, cond)
	if err != nil {
		log.Error(ctx.Wrap(c)).Msgf("根据消息ID找附近的消息错误:" + err.Error())

		return res, err
	}
	res.List = s.FormatMsgList(c, mediaService, list)

	return res, nil
}

// FormatMsg 格式化消息列表
func (s Service) FormatMsgList(c ctx.Context, mediaService MsgMediaService, msgList []Msg) []Msg {
	if len(msgList) == 0 {
		return msgList
	}
	enterpriseID := msgList[0].EnterpriseID

	// 放入所有的userid
	var useridArr []string
	var fromUserRemark string
	var fromUser string
	var toUser string
	var fromUserID int64
	var toUserID int64
	var fromUserIDArr []int64

	for i, msg := range msgList {
		contentMap := make(map[string]interface{})
		_ = json.Unmarshal([]byte(msg.Content), &contentMap)
		var toListArr []string
		_ = json.Unmarshal([]byte(msg.ToList), &toListArr)

		// 放入所有的userid,下面再遍历一次 给每个消息加上user信息
		useridArr = append(useridArr, msg.FromUser)
		if len(toListArr) == 1 {
			for _, toUserid := range toListArr {
				useridArr = append(useridArr, toUserid)
			}
		}

		if msg.FromUserType == UserTypeMember {
			toUser = msg.FromUser
		}

		if msg.FromUserType == UserTypeWorkExternalUser || msg.FromUserType == UserTypeExternalUser {
			fromUser = msg.FromUser
			toUser = toListArr[0]
		}
		msgList[i].ToListArr = toListArr
		msgList[i].ContentMap = contentMap
	}

	var memberUseridArr []string
	var robotUseridArr []string
	var externalUseridArr []string
	for _, userid := range useridArr {
		userType := s.ConversionUserType(userid)

		if userType == UserTypeMember {
			memberUseridArr = append(memberUseridArr, userid)
		}
		if userType == UserTypeRobot {
			robotUseridArr = append(robotUseridArr, userid)
		}
		if userType == UserTypeWorkExternalUser || userType == UserTypeExternalUser {
			externalUseridArr = append(externalUseridArr, userid)
		}
	}

	var userInfoArr []UserInfo
	var memberInfoArr *[]UserInfo
	var robotInfoArr *[]UserInfo
	var externalInfoArr *[]UserInfo
	var userInfoErr error

	if len(memberUseridArr) != 0 {
		memberInfoArr, userInfoErr = s.GetUserInfo(c, UserTypeMember, memberUseridArr, enterpriseID)
		if userInfoErr == nil && memberInfoArr != nil {
			for _, v := range *memberInfoArr {
				userInfoArr = append(userInfoArr, v)

				if v.UserID == toUser {
					toUserID = v.ObjectID
				}
			}
		}
	}
	if len(robotUseridArr) != 0 {
		robotInfoArr, userInfoErr = s.GetUserInfo(c, UserTypeRobot, robotUseridArr, enterpriseID)
		if userInfoErr == nil && robotInfoArr != nil {
			for _, v := range *robotInfoArr {
				userInfoArr = append(userInfoArr, v)
			}
		}
	}
	if len(externalUseridArr) != 0 {
		externalInfoArr, userInfoErr = s.GetUserInfo(c, UserTypeExternalUser, externalUseridArr, enterpriseID)
		if userInfoErr == nil && externalInfoArr != nil {
			for _, v := range *externalInfoArr {
				userInfoArr = append(userInfoArr, v)

				if v.UserID == fromUser {
					fromUserID = v.ObjectID
					fromUserIDArr = append(fromUserIDArr, fromUserID)
				}
			}
		}
	}

	var sdkfileidArr []string

	// 每个消息放入对应的用户信息
	customerService := customer.New(mysql.NewCustomer(), mysql.NewWorkUser(), userEs.NewEsRepo())
	externalList, eternalListErr := customerService.ExternalListRemark(c, enterpriseID, toUserID, fromUserIDArr)
	if eternalListErr != nil {
		log.Error(ctx.Wrap(c)).Msgf("查询联系人列表备注失败,错误信息:%+v", eternalListErr)
	}

	for _, remarkRelation := range externalList {
		if remarkRelation.WorkExternalUserID == fromUserID {
			fromUserRemark = remarkRelation.Remark
		}
	}

	for i, userInfo := range userInfoArr {
		if userInfo.UserID == fromUser {
			userInfoArr[i].Remark = fromUserRemark
		}
	}
	for i, msg := range msgList {
		for _, userInfo := range userInfoArr {
			if msg.FromUser == userInfo.UserID {
				msgList[i].FromUserInfo = userInfo
			}
			if len(msg.ToListArr) != 0 && msg.ToListArr[0] == userInfo.UserID {
				msgList[i].ToUserInfo = userInfo
			}
		}
		for _, v := range msg.GetSdkFileId() {
			sdkfileidArr = append(sdkfileidArr, v)
		}
	}

	// 从媒体表where in出要的记录,再遍历放入对应消息的content_map中
	msgMediaList := mediaService.ListByFileId(c, enterpriseID, sdkfileidArr)
	for i, msg := range msgList {
		if helper.InArray[string](msg.Type, FileType) == false {
			continue
		}

		// 混合消息单独处理
		if msg.Type == chat.MsgTypeMixed {
			var mixedStruct chat.Mixed
			_ = json.Unmarshal([]byte(msg.Content), &mixedStruct)

			contentMapItem := mixedStruct.Item
			for itemIndex, v := range contentMapItem {
				var itemContentMap map[string]interface{}
				_ = json.Unmarshal([]byte(v.Content), &itemContentMap)
				sdkFileid, ok := itemContentMap["sdkfileid"]
				if ok == true {
					itemContentMap["media_info"] = MatchMediaInfo(msgMediaList, sdkFileid.(string))
				}
				contentMapItem[itemIndex].ContentMap = itemContentMap
				msgList[i].ContentMap["item"] = contentMapItem
			}

			continue
		}
		// 没有媒体id忽略
		if _, ok := msg.ContentMap["sdkfileid"]; ok == false {
			continue
		}

		// 最终给格式化完毕的含媒体文件的数据必然存在media_info
		msgList[i].ContentMap["media_info"] = MatchMediaInfo(msgMediaList, msg.ContentMap["sdkfileid"].(string))
	}

	return msgList
}

// MatchMediaInfo 匹配媒体信息
func MatchMediaInfo(msgMediaList []MsgMedia, sdkFileId string) map[string]interface{} {
	info := map[string]interface{}{
		"path":   "",
		"domain": config.CfgData.Storage.Blob.Domain,
	}
	for _, msgMedia := range msgMediaList {
		if sdkFileId == msgMedia.SdkFileID {
			info["path"] = msgMedia.Path
			break
		}
	}

	return info
}

// GetSdkFileId 获取内容中的媒体id
func (s Service) GetSdkFileId(msgType string, content string) []string {
	var sdkFileIdArr []string

	var contentMap map[string]interface{}
	_ = json.Unmarshal([]byte(content), &contentMap)

	// 混合消息单独处理
	if msgType == chat.MsgTypeMixed {
		var mixedStruct chat.Mixed
		_ = json.Unmarshal([]byte(content), &mixedStruct)

		for _, v := range mixedStruct.Item {
			var itemContentMap map[string]interface{}
			_ = json.Unmarshal([]byte(v.Content), &itemContentMap)
			if sdkFileid, ok := itemContentMap["sdkfileid"]; ok {
				sdkFileIdArr = append(sdkFileIdArr, sdkFileid.(string))
			}
		}
	}

	// 会话记录消息单独处理
	if msgType == chat.MsgTypeChatrecord {
		// 如果含会话记录消息则处理3层,且第3层如果还存在会话记录消息则为"未知的消息"
		var chatrecord chat.Chatrecord
		_ = json.Unmarshal([]byte(content), &chatrecord)

		for _, v := range chatrecord.Item {

			var vContentMap map[string]interface{}
			_ = json.Unmarshal([]byte(v.Content), &vContentMap)
			if sdkFileid, ok := vContentMap["sdkfileid"]; ok {
				sdkFileIdArr = append(sdkFileIdArr, sdkFileid.(string))
			}
		}
	}

	// 观察微信文档基本直接转map取sdkfileid字段即可
	if sdkfileid, ok := contentMap["sdkfileid"]; ok {
		sdkFileIdArr = append(sdkFileIdArr, sdkfileid.(string))
	}

	return sdkFileIdArr
}

// HandleMsg 处理消息
func (s Service) HandleMsg(c ctx.Context, enterpriseInfo enterprise.Corp, plainMsg chat.PlainMsg, receiveUserService ReceiveUserService, objectService ObjectService, memberRelationService MemberService, sensitiveService sensitive.Service, groupChatService GroupChatService, groupChatMemberService GroupChatMemberService) {
	// 切换企业"消息"忽略掉
	if plainMsg.Action == "switch" {
		return
	}
	toList, _ := json.Marshal(plainMsg.Tolist)
	t := time.Now().Unix()

	content := plainMsg.Content()
	var contentMap map[string]interface{}
	_ = json.Unmarshal([]byte(content), &contentMap)

	// 音频存档消息的voiceid返回再外层,这里特殊处理 放入content json中
	if plainMsg.Msgtype == chat.MsgTypeMeetingVoiceCall {
		contentMap["voiceid"] = plainMsg.Voiceid
	}

	buf, _ := json.Marshal(contentMap)
	content = string(buf)

	// 消息表新增
	msg := Msg{
		ID:           uid.ID(c),
		MsgId:        plainMsg.MsgID,
		Action:       plainMsg.Action,
		FromUser:     plainMsg.From,
		FromUserType: s.ConversionUserType(plainMsg.From),
		ToList:       string(toList),
		RoomId:       plainMsg.Roomid,
		MsgAt:        plainMsg.Msgtime / 1000,
		MsgTime:      plainMsg.Msgtime,
		Type:         plainMsg.Msgtype,
		Content:      content,
		Seq:          plainMsg.EncryptDataChatData.Seq,
		EnterpriseID: enterpriseInfo.EnterpriseID,
		CreatedAt:    t,
		UpdatedAt:    t,
		DeletedAt:    0,
	}
	s.Add(c, msg)

	if msg.RoomId == "" && msg.FromUserType == UserTypeMember && (msg.Type == chat.MsgTypeText || msg.Type == chat.MsgTypeCard) {
		var text chat.Text
		contentMap := make(map[string]interface{})
		content := ""
		if msg.Type == chat.MsgTypeText {
			_ = json.Unmarshal([]byte(msg.Content), &contentMap)
			_ = mapstructure.Decode(contentMap, &text)
			content = text.Content
		}
		var toListMap []string

		_ = json.Unmarshal([]byte(msg.ToList), &toListMap)
		if len(toListMap) != 0 {
			// 敏感词检测
			sensitiveReq := sensitive.SendMessageRequest{
				EnterpriseID:  enterpriseInfo.EnterpriseID,
				ChatMessageID: msg.ID,
				Type:          msg.Type,
				ChatRecord:    content,
				FromUser:      msg.FromUser,
				ToUser:        toListMap[0],
			}
			_, _ = sensitiveService.SensitiveMatching(c, &sensitiveReq, enterpriseInfo.Domain)
		}
	}

	// 群聊
	if plainMsg.Roomid != "" {
		// 拉取内部群
		log.Info(c).Msgf("拉取内部群开始:")
		s.HandleInsideGroup(c, enterpriseInfo, msg, groupChatService, groupChatMemberService)
		log.Info(c).Msgf("拉取内部群结束")
		// 消息对象表
		_ = objectService.Add(c, s, memberRelationService, Object{
			FromUserid:   "",
			ToUserid:     plainMsg.Roomid,
			FromChannel:  UserTypeMember,
			ToChannel:    UserTypeGroup,
			EnterpriseID: enterpriseInfo.EnterpriseID,
		})
		return
	}

	// 接收人消息表新增数据
	for _, toUser := range plainMsg.Tolist {
		receiveUserService.Add(c, ReceiveUser{
			ChatMsgId:    msg.ID,
			ToUser:       toUser,
			UserChannel:  s.ConversionUserType(toUser),
			EnterpriseID: enterpriseInfo.EnterpriseID,
			CreatedAt:    t,
			UpdatedAt:    t,
			DeletedAt:    0,
		})

		// 消息对象表
		_ = objectService.Add(c, s, memberRelationService, Object{
			FromUserid:   plainMsg.From,
			ToUserid:     toUser,
			FromChannel:  s.ConversionUserType(plainMsg.From),
			ToChannel:    s.ConversionUserType(toUser),
			EnterpriseID: enterpriseInfo.EnterpriseID,
		})
	}
}

// GetMsg 获取聊天消息
func (s Service) GetMsg(c ctx.Context, receiveUserService ReceiveUserService, objectService ObjectService, memberRelationService MemberService, groupChatService GroupChatService, groupChatMemberService GroupChatMemberService) {
	sensitiveService := sensitive.New(sensitiveRepository.New())

	// 遍历chat_settings表的企业开始处理
	settingService := settings.New(settingRepository.New())
	for _, enterpriseSetting := range settingService.EnterpriseList(c) {
		// 查配置信息
		enterpriseInfo, settingInfo, err := s.GetConfig(c, enterpriseSetting.EnterpriseId)
		if err != nil {
			log.Info(c).Msgf("获取集合配置信息错误:" + err.Error())

			continue
		}
		// 如果企业未配置绕过
		if settingInfo.Secret == "" {
			log.Info(c).Msgf("企业[%d]未配置会话存档,配置信息：%+v", enterpriseSetting.EnterpriseId, settingInfo)

			continue
		}

		// 每个企业单独锁
		mtx, err := lock.Mutex(c, lock.Key("micro-scrm:get-msg:"+strconv.FormatInt(enterpriseInfo.EnterpriseID, 10)), lock.TTL(10*time.Minute), lock.Retry(lock.NoRetry()))
		if err != nil {
			continue
		}
		err = mtx.Lock()
		defer mtx.Unlock()

		// 每个企业各自运行
		//go func() {
		// 查出当前企业最后一次消息的seq
		lastMsg := s.GetLastMsg(c, enterpriseInfo.EnterpriseID)
		chatMsgRequest := chat.ChatMsgRequest{
			Seq:     lastMsg.Seq,
			Limit:   500,
			Proxy:   "",
			Auth:    "",
			Timeout: 3,
			CorpID:  enterpriseInfo.CorpID,
			Secret:  settingInfo.Secret,
		}
		req := &chat.GetChatMsgReq{
			ChatMsgRequest: chatMsgRequest,
			PK:             settingInfo.MessagePrivateKey,
		}
		chatMsgList, chatMsgListErr := chat.NewService(config.CfgData.Service.MicroChatScrm.Domain).GetChatMsg(c, req)
		if chatMsgListErr != nil {
			log.Info(c).Msgf("获取会话消息错误:" + chatMsgListErr.Error())

			continue
		}

		for _, v := range chatMsgList {
			s.HandleMsg(c, *enterpriseInfo, *v, receiveUserService, objectService, memberRelationService, *sensitiveService, groupChatService, groupChatMemberService)
		}
		//}()
	}
}

// GetMsgMedia 获取聊天消息媒体文件
func (s Service) GetMsgMedia(c ctx.Context, msgMediaService MsgMediaService) {
	// 遍历chat_settings表的企业开始处理
	settingService := settings.New(settingRepository.New())
	for _, enterpriseSetting := range settingService.EnterpriseList(c) {
		// 查配置信息
		enterpriseInfo, settingInfo, err := s.GetConfig(c, enterpriseSetting.EnterpriseId)
		if err != nil {
			log.Info(c).Msgf("获取集合配置信息错误:" + err.Error())

			continue
		}

		// 查询未转换的消息，只转换近3天的数据
		StartAt := time.Now().Unix() - 259200
		NotMediaList, err := s.repo.GetNotMediaList(c, FileType, enterpriseInfo.EnterpriseID, StartAt)
		if err != nil {
			continue
		}
		var toTransMediaList []MsgMedia
		for _, value := range *NotMediaList {
			detail := MsgMedia{}
			// 混合媒体消息
			if value.Type == chat.MsgTypeMixed {
				// todo
				var mixcontent MixMediaContent
				_ = json.Unmarshal([]byte(value.Content), &mixcontent)
				if mixcontent.Items != nil {
					for k, items := range mixcontent.Items {
						var itemContent ContentEntity
						_ = json.Unmarshal([]byte(items.Content), &itemContent)
						if itemContent.SdkFileid != "" {
							suffix := ContentEntity{
								FileType:    items.Type,
								EmotionType: itemContent.EmotionType,
								FileExt:     itemContent.FileExt,
								FileName:    itemContent.FileName,
							}
							detail.SdkFileID = itemContent.SdkFileid
							detail.EnterpriseID = value.EnterpriseID
							detail.ID = value.ID
							detail.Suffix = s.GetMediaFileSuffix(suffix)
							if value.Type == chat.MsgTypeFile || value.Type == chat.MsgTypeVoipDocShare {
								detail.Name = fmt.Sprintf("work-chat/chat-msg/media/output-%s", itemContent.FileName)
							} else {
								detail.Name = fmt.Sprintf("work-chat/chat-msg/media/output-%d-%d-%d"+detail.Suffix, value.ID, time.Now().Unix(), k)
							}
							toTransMediaList = append(toTransMediaList, detail)
						}
					}
				}
			} else {
				var content ContentEntity

				_ = json.Unmarshal([]byte(value.Content), &content)
				suffix := ContentEntity{
					FileType:    value.Type,
					EmotionType: content.EmotionType,
					FileExt:     content.FileExt,
					FileName:    content.FileName,
				}
				detail.Suffix = s.GetMediaFileSuffix(suffix)
				detail.EnterpriseID = value.EnterpriseID
				if content.SdkFileid != "" {
					detail.SdkFileID = content.SdkFileid
				}
				detail.ID = value.ID
				if value.Type == chat.MsgTypeFile || value.Type == chat.MsgTypeVoipDocShare {
					detail.Name = fmt.Sprintf("work-chat/chat-msg/media/output-%s", content.FileName)
				} else {
					detail.Name = fmt.Sprintf("work-chat/chat-msg/media/output-%d-%d"+detail.Suffix, detail.ID, time.Now().Unix())
				}
				toTransMediaList = append(toTransMediaList, detail)
			}
		}

		log.Info(c).Msgf("toTransMediaList: %+v", toTransMediaList)
		for _, value := range toTransMediaList {
			sdkFileID := value.SdkFileID
			MediaDataRequest := chat.MediaDataRequest{
				SdkFileid: sdkFileID,
				Proxy:     "",
				Auth:      "",
				Timeout:   30,

				CorpID: enterpriseInfo.CorpID,
				Secret: settingInfo.Secret,
			}

			if err != nil {
				log.Info(c).Msgf("初始化 writer错误: %+v", err)
				continue
			}
			transMediaErr := chat.NewService(config.CfgData.Service.MicroChatScrm.Domain).GetMediaData(
				c,
				&chat.GetMediaDataReq{
					MediaDataRequest: MediaDataRequest,
					// 重试次数
					RetryLimit: 3,
					Drive:      "blob",
					Name:       value.Name,
				},
				nil,
				500,
			)
			// 转换成功，添加到表中
			mediaRecord := MsgMedia{
				ID:           uid.ID(ctx.Wrap(c)),
				ChatMsgID:    value.ID,
				SdkFileID:    value.SdkFileID,
				EnterpriseID: value.EnterpriseID,
				Path:         "/" + value.Name,
			}
			errAdd := msgMediaService.Add(c, mediaRecord)
			if errorIgnore.Ignore(errAdd) != nil {
				continue
			}

			if transMediaErr != nil {
				log.Info(c).Msgf("MediaDataRequest: %+v", MediaDataRequest)
				log.Info(c).Msgf("transMediaErr: %+v", transMediaErr)
				continue
			}

		}
	}
}

// GetUserInfo 获取用户信息
func (s Service) GetUserInfo(c ctx.Context, userType int64, userId []string, enterpriseId int64) (*[]UserInfo, error) {
	if userType == UserTypeMember {
		// 查询员工
		list, err := s.repo.GetMembersList(c, userId, enterpriseId)
		if err != nil {
			return nil, err
		}

		return list, nil
	}

	if userType == UserTypeExternalUser {
		// 查询外部联系人
		list, err := s.repo.GetExternalUsersList(c, userId, enterpriseId)
		if err != nil {
			return nil, err
		}

		return list, nil
	}

	if userType == UserTypeRobot {
		// 查询机器人
		robot, err := s.GetRobotInfo(c, userId, enterpriseId)
		if err != nil {
			return nil, err
		}
		var list []UserInfo
		list = append(list, *robot)

		return &list, nil
	}

	var userInfo *[]UserInfo

	return userInfo, errors.New("不支持的用户渠道")
}

// GetCheckSingleAgree 获取会话同意情况
func (s Service) GetCheckSingleAgree(c ctx.Context, req AgreeInfoRequest) ([]work.CheckSingleAgreeResponse, error) {
	var res []work.CheckSingleAgreeResponse
	var infoArr []work.CheckSingleAgreeInfoRequest
	for _, v := range req.List {
		infoArr = append(infoArr, work.CheckSingleAgreeInfoRequest{
			Userid:          v.Userid,
			Exteranalopenid: v.ExteranalUserid,
		})
	}

	enterpriseService := enterprise.New(enterpriseRepository.New())
	externalReq := work.CheckSingleAgreeRequest{
		QyGwRequest: work.QyGwRequest{},
		Info:        infoArr,
	}
	enterpriseData, enterpriseErr := enterpriseService.Detail(ctx.Wrap(c), enterprise.Corp{
		EnterpriseID: req.EnterpriseID,
	})
	if enterpriseErr != nil {
		log.Error(ctx.Wrap(c)).Msgf("查询企业信息失败,错误信息:%+v", enterpriseErr)

		return res, enterpriseErr
	}

	workImpl := new(work.QyWorkImpl).SetDomain(enterpriseData.Domain)
	checkSingleAgree, csaErr := workImpl.CheckSingleAgree(c, externalReq)
	if csaErr != nil {
		log.Error(ctx.Wrap(c)).Msgf("调用企业微信,获取会话同意情况错误:%+v", csaErr)

		return res, csaErr
	}

	return checkSingleAgree, nil
}

// GetRobotInfo 获取机器人信息
func (s *Service) GetRobotInfo(c ctx.Context, req []string, enterpriseId int64) (*UserInfo, error) {
	enterpriseService := enterprise.New(enterpriseRepository.New())

	robotReq := work.RobotInfoRequest{
		QyGwRequest: work.QyGwRequest{},
		RobotID:     req[0],
	}
	enterpriseData, enterpriseErr := enterpriseService.Detail(ctx.Wrap(c), enterprise.Corp{
		EnterpriseID: enterpriseId,
	})
	if enterpriseErr != nil {
		log.Error(ctx.Wrap(c)).Msgf("查询企业信息失败，错误信息：%+v", enterpriseErr)

		return nil, enterpriseErr
	}

	workImpl := new(work.QyWorkImpl).SetDomain(enterpriseData.Domain)
	robotInfo, csaErr := workImpl.GetRobotInfo(c, robotReq)
	if csaErr != nil {
		log.Error(ctx.Wrap(c)).Msgf("调用企业微信,获取会话同意情况错误:%+v", csaErr)

		return nil, csaErr
	}

	userRobot := UserInfo{
		UserID: robotInfo.RobotId,
		Name:   robotInfo.Name,
	}

	return &userRobot, nil
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

// Add 新增消息成功返回1
func (s Service) Add(c ctx.Context, row Msg) typex.ID {
	if s.MsgExists(c, row.EnterpriseID, row.MsgId) {
		// TODO::存在可以考虑更新

		return 1
	}
	if row.ID == 0 {
		row.ID = uid.ID(c)
	}
	num, err := s.repo.Add(c, row)
	if err != nil {
		log.Error(ctx.Wrap(c)).Msgf("新增消息记录错误:" + err.Error())

		return 0
	}

	return typex.ID(num)
}

// MsgExists 判断消息是否已经存在
func (s Service) MsgExists(c ctx.Context, enterpriseId int64, msgId string) bool {
	ex, _ := s.repo.Exists(c, enterpriseId, msgId)

	return ex
}

// GetLastMsg 获取最后一条消息
func (s Service) GetLastMsg(c ctx.Context, enterpriseId int64) Msg {
	var msg Msg
	_ = s.repo.LastRow(c, &msg, enterpriseId)

	return msg
}

// GetUserType 转换用户类型
func (s Service) ConversionUserType(userid string) int64 {
	if len(userid) < 2 {
		return UserTypeMember
	}
	prefix := userid[0:2]

	// 机器人
	if prefix == "wb" {
		return UserTypeRobot
	}

	// 微信外部联系人
	if prefix == "wo" {
		return UserTypeWorkExternalUser
	}

	// 企业微信外部联系人
	if prefix == "wm" {
		return UserTypeExternalUser
	}

	return UserTypeMember
}

// GetMediaFileSuffix 获取媒体文件的后缀
func (s *Service) GetMediaFileSuffix(suff ContentEntity) string {
	var suffix = ""
	switch suff.FileType {
	case "image":
		suffix = ".jpg"
		break
	case "voice":
		suffix = ".amr"
		break
	case "video":
		suffix = ".mp4"
		break
	case "emotion":
		if suff.EmotionType == 1 {
			suffix = ".gif"
		} else if suff.EmotionType == 2 {
			suffix = ".png"
		}
		break
	case "file":
		if suff.FileExt != "" {
			suffix = "." + suff.FileExt
		}
		break
	case "meeting_voice_call":
		suffix = ".amr"
		break
	case "voip_doc_share":
		suffix = "." + suff.FileName
		break
	}
	return suffix
}

// ConversionAgreeStatus 转换同意状态
func (s *Service) ConversionAgreeStatus(status string) int64 {
	if status == "Disagree" {
		return AgreeStatusDisAgree
	}
	if status == "Agree" {
		return AgreeStatusAgree
	}

	return AgreeStatusDefaultAgree
}

// 内部群处理
func (s *Service) HandleInsideGroup(c ctx.Context, enterpriseInfo enterprise.Corp, msg Msg, groupChatService GroupChatService, groupChatMemberService GroupChatMemberService) {
	customerService := customer.New(mysql.NewCustomer(), mysql.NewWorkUser(), userEs.NewEsRepo())
	// 判断该roomId是否已记录到表中
	req := WorkGroupChat{
		ChatID:       msg.RoomId,
		EnterpriseID: msg.EnterpriseID,
		DeletedAt:    0,
	}
	log.Info(c).Msgf("baseGroupInfoRequest: %+v", req)
	group, err := groupChatService.GetGroupInfo(ctx.Wrap(c), req)
	if errorIgnore.Ignore(err) != nil {
		log.Info(c).Msgf("baseGroupInfoErr: %+v", err)
	}

	// 拉取内部群消息
	workImpl := new(work.QyWorkImpl).SetDomain(enterpriseInfo.Domain)
	groupReq := work.InsideGroupRequest{
		QyGwRequest: work.QyGwRequest{},
		RoomID:      msg.RoomId,
	}
	groupChatInfo, groupErr := workImpl.InsideGroupChatGet(ctx.Wrap(c), groupReq)
	log.Info(c).Msgf("WorkGroupChatInfo: %+v", groupChatInfo)
	if groupErr != nil {
		log.Info(c).Msgf("获取内部群消息错误: %+v", groupErr)
		return
	}

	var memberUserIDs []string
	// =====================客户群新增成员====================
	var chatMembersByID = make(map[string]work.GroupMembers)
	for _, chatMemberItem := range groupChatInfo.MemberList {
		memberUserIDs = append(memberUserIDs, chatMemberItem.MemberID)
		chatMembersByID[chatMemberItem.MemberID] = chatMemberItem
	}

	// =====================查询群里的内部群员（企业微信成员）====================
	membersMapFromDb, _ := s.repo.GetMembersList(c, memberUserIDs, msg.EnterpriseID)

	var membersListByID = make(map[string]UserInfo)
	for _, MemberItem := range *membersMapFromDb {
		membersListByID[MemberItem.UserID] = MemberItem
	}
	log.Info(c).Msgf("membersListByID: %+v", membersListByID)
	// =====================群信息处理=========================================
	var workGroupChatID int64
	var hasExistWorkGroupMembersFromDb = make(map[string]customer.WorkGroupChatMember)
	if group == nil {
		// 添加群信息
		insertGroupChat := customer.WorkGroupChat{}
		insertGroupChat.ID = uid.ID(c)
		insertGroupChat.EnterpriseID = msg.EnterpriseID
		insertGroupChat.ChatID = msg.RoomId
		insertGroupChat.Name = groupChatInfo.RoomName
		insertGroupChat.Notice = groupChatInfo.Notice
		insertGroupChat.CreateTime = groupChatInfo.RoomCreateTime
		insertGroupChat.CreatedAt = time.Now().Unix()
		insertGroupChat.UpdatedAt = time.Now().Unix()
		insertGroupChat.Type = GroupTypeInside
		insertGroupChat.Status = customer.WorkGroupTypeExists
		insertGroupChat.MemberNum = len(groupChatInfo.MemberList)
		if owner, ok := membersListByID[groupChatInfo.Creator]; ok {
			insertGroupChat.OwnerMemberID = owner.ObjectID
		}
		_, aErr := groupChatService.repo.Add(c, insertGroupChat)
		if aErr != nil {
			log.Info(c).Msgf("添加内部群消息错误: %+v,错误信息：%+v", insertGroupChat, aErr)
		}
		workGroupChatID = insertGroupChat.ID
	} else {
		// 更新群信息
		group.Name = groupChatInfo.RoomName
		group.CreateTime = groupChatInfo.RoomCreateTime
		group.Notice = groupChatInfo.Notice
		group.UpdatedAt = time.Now().Unix()
		group.MemberNum = len(groupChatInfo.MemberList)
		if owner, ok := membersListByID[groupChatInfo.Creator]; ok {
			group.OwnerMemberID = owner.ObjectID
		}
		_, uErr := groupChatService.repo.Update(c, *group)
		if uErr != nil {
			log.Error(c).Msgf("更新客户群信息失败，更新数据:%+v,错误信息：%+v", group, uErr)
		}
		workGroupChatID = group.ID

		hasExistWorkGroupMembersFromDb = searchWorkChatMemberFromDb(c, workGroupChatID, msg.EnterpriseID, customerService)
	}
	// ====================群成员处理========================================
	var chatMembers WorkGroupChatMember
	var delMembers WorkGroupChatMember
	var chatMembersList = make(map[string]WorkGroupChatMember)
	for _, value := range *membersMapFromDb {
		if value.Status != MemberStatusNormal {
			continue

		}
		chatMembers.ID = uid.ID(c)
		chatMembers.Status = customer.WorkGroupChatMemberJoinStatus
		chatMembers.Type = customer.WorkChatGroupMemberTypeQw
		chatMembers.UserID = value.UserID
		chatMembers.MemberID = value.ObjectID
		chatMembers.EnterpriseID = value.EnterpriseID
		chatMembers.GroupChatID = workGroupChatID
		if Member, ok := chatMembersByID[value.UserID]; ok {
			chatMembers.JoinTime = Member.JoinTime
		}
		chatMembers.CreatedAt = time.Now().Unix()
		chatMembers.UpdatedAt = time.Now().Unix()
		// 查询该成员是否已经在库中
		if _, ok := hasExistWorkGroupMembersFromDb[chatMembers.UserID]; !ok {
			groupChatMemberService.repo.Add(c, chatMembers)
		}
		chatMembersList[value.UserID] = chatMembers
	}
	log.Info(c).Msgf("hasExistWorkGroupMembersFromDb: %+v", hasExistWorkGroupMembersFromDb)
	var deleteIds []int64
	// 循环查询是否有人退群
	if hasExistWorkGroupMembersFromDb != nil {
		for _, exisMember := range hasExistWorkGroupMembersFromDb {
			if _, ok := chatMembersList[exisMember.UserID]; !ok {
				deleteIds = append(deleteIds, exisMember.ID)
			}
		}
		delMembers.Status = customer.WorkGroupChatMemberLeaveStatus
		delMembers.LeaveTime = time.Now().Unix()
		delMembers.UpdatedAt = time.Now().Unix()
		delMembers.EnterpriseID = enterpriseInfo.EnterpriseID
		log.Info(c).Msgf("deleteIds: %+v", deleteIds)
		if deleteIds != nil {
			groupChatMemberService.BatchUpdateWorkGroupChatMembers(c, deleteIds, delMembers)
		}
	}

}

// searchWorkChatMemberFromDb 从数据库查询客户群成员信息
func searchWorkChatMemberFromDb(ctx ctx.Context, workGroupChatID int64, enterpriseID int64, customerService *customer.Service) map[string]customer.WorkGroupChatMember {
	workGroupMembersMap := make(map[string]customer.WorkGroupChatMember)
	// 查询已有的群成员
	batchReq := customer.BatchWorkGroupMemberRequest{
		Size:         0,
		EnterpriseID: enterpriseID,
		GroupChatID:  workGroupChatID,
	}
	total, workGroupMembers, _ := customerService.BatchFindWorkGroupChatMembers(ctx, batchReq)
	if total > 0 {
		for _, v := range workGroupMembers {
			if v.Status == callback.WorkChatGroupMemberStatusJoin {
				workGroupMembersMap[v.UserID] = v
			}
		}
	}
	return workGroupMembersMap
}
