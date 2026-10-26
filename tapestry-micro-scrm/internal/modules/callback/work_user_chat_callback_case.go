package callback

import (
	"encoding/json"
	"fmt"
	"micro-scrm/internal/modules/customer/respository/elastic"
	customerMysql "micro-scrm/internal/modules/customer/respository/mysql"
	"micro-scrm/internal/modules/drainage"
	drainageMysql "micro-scrm/internal/modules/drainage/respository/mysql"
	enterpriseMysql "micro-scrm/internal/modules/enterprise/respository/mysql"
	"micro-scrm/internal/modules/fission"
	fissionBossApi "micro-scrm/internal/modules/fission/api/bossapi"
	fissionMysql "micro-scrm/internal/modules/fission/respository/mysql"
	membersMysql "micro-scrm/internal/modules/members/respository/mysql"
	userModule "micro-scrm/internal/modules/user"
	userMysql "micro-scrm/internal/modules/user/respository/mysql"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/events/entity"
	"micro-scrm/internal/logic/change_external_contact/community_task"
	"micro-scrm/internal/modules/customer"
	"micro-scrm/internal/modules/enterprise"
	"micro-scrm/internal/modules/members"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm/group_chat"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/sdk/work"
	"micro-scrm/internal/pkg/uid"
)

// ChangeExternalChat 群聊回调事件处理
func (s Service) ChangeExternalChat(ctx ctx.Context, callbackData *work.MemberExternalCallback, enterpriseData *enterprise.Corp, workImpl *work.QyWorkImpl) error {
	// 查询群详情
	workImpl = workImpl.SetDomain(enterpriseData.Domain)
	req := work.ExternalGroupRequest{
		ChatID:      callbackData.ChatID,
		NeedName:    1,
		QyGwRequest: work.QyGwRequest{},
	}
	if req.ChatID == "" {
		log.Info(ctx).Msgf("群聊回调时间处理：chatid为空，callbackData:%v", callbackData)
		return nil
	}
	chatDetail, err := workImpl.ExternalGroupChatGet(ctx, req)
	if err != nil {
		return err
	}
	var adminList []string
	for _, a := range chatDetail.AdminList {
		adminList = append(adminList, a.Userid)
	}
	admin := "[]"
	if len(adminList) > 0 {
		a1, _ := json.Marshal(adminList)
		admin = string(a1)
	}

	// 查询群是否已经在数据库中存在
	chatReq := customer.WorkGroupChat{
		ChatID:       callbackData.ChatID,
		EnterpriseID: enterpriseData.EnterpriseID,
	}
	chat, cErr := s.CustomerService.FindWorkGroupChatDetail(ctx, chatReq)
	if cErr != nil {
		log.Info(ctx).Msgf("查询客户群失败，错误信息：%+v,查询条件：%+v", cErr, chatReq)
		return cErr
	}

	var memberUserIDs []string
	var existsInQwMembers []*customer.WorkGroupChatMember
	// 群存在
	workGroupType := customer.WorkGroupTypeDissolve
	memberUserIDs = append(memberUserIDs, chatDetail.Owner)
	/*
	*人员列表
	 */
	chatMemberlogs := map[string]work.ExternalGroupMemberList{}
	// =====================客户群新增成员====================
	var chatMemberNum = len(chatDetail.MemberList)
	if chatMemberNum > 0 {
		for _, chatMemberItem := range chatDetail.MemberList {
			existsInQwMembers = append(existsInQwMembers, customer.TransferWorkGroupChatMember(&chatMemberItem))
			if chatMemberItem.Type == customer.WorkChatGroupMemberTypeQw {
				memberUserIDs = append(memberUserIDs, chatMemberItem.UserID)
			}
			chatMemberlogs[chatMemberItem.UserID] = chatMemberItem
		}
		workGroupType = customer.WorkGroupTypeExists
	}
	// =====================查询群里的内部群员（企业微信成员）====================
	membersMapFromDb := s.SearchChatMemberInner(ctx, memberUserIDs, enterpriseData.EnterpriseID)
	// =====================群信息处理=========================================
	var workGroupChatID int64
	var hasExistWorkGroupMembersFromDb map[string]customer.WorkGroupChatMember
	// 新增群信息
	if chat == nil {
		insertGroupChat := customer.TransferWorkGroupChat(chatDetail)
		insertGroupChat.ID = uid.ID(ctx)
		insertGroupChat.CreatedAt = time.Now().Unix()
		insertGroupChat.UpdatedAt = time.Now().Unix()
		insertGroupChat.EnterpriseID = enterpriseData.EnterpriseID
		insertGroupChat.Type = customer.WorkGroupChatInnerGroup
		insertGroupChat.Status = workGroupType
		insertGroupChat.MemberNum = chatMemberNum
		insertGroupChat.FormatNotice()
		insertGroupChat.AdminList = admin
		insertGroupChat.Name = fmt.Sprintf("用户[%s]发起群聊", chatDetail.Owner)
		if owner, ok := membersMapFromDb[chatDetail.Owner]; ok {
			insertGroupChat.OwnerMemberID = owner.ID
			insertGroupChat.Name = fmt.Sprintf("用户[%s]发起群聊", owner.Name)
			// insertGroupChat.OwnerMemberID = owner.Userid
		}
		_, aErr := s.CustomerService.AddWorkGroupChat(ctx, *insertGroupChat)
		if aErr != nil {
			return aErr
		}
		workGroupChatID = insertGroupChat.ID
		chat = insertGroupChat
	} else {
		// 更新群信息
		if chatDetail.Name != "" {
			chat.Name = chatDetail.Name
		}
		chat.Notice = chatDetail.Notice
		chat.UpdatedAt = time.Now().Unix()
		chat.Status = workGroupType
		chat.MemberNum = chatMemberNum
		chat.AdminList = admin
		chat.FormatNotice()
		if workGroupType == customer.WorkGroupTypeDissolve && chat.DissolveTime == 0 {
			chat.DissolveTime = chat.UpdatedAt
		}
		if owner, ok := membersMapFromDb[chatDetail.Owner]; ok {
			chat.OwnerMemberID = owner.ID
			// chat.OwnerMemberID = owner.Userid
		}
		uErr := s.CustomerService.UpdateWorkGroupChat(ctx, *chat)
		if uErr != nil {
			log.Error(ctx).Msgf("更新客户群信息失败，更新数据:%+v,错误信息：%+v", chat, uErr)
		}
		workGroupChatID = chat.ID
		hasExistWorkGroupMembersFromDb, _ = searchWorkChatMemberFromDb(ctx, workGroupChatID, enterpriseData.EnterpriseID, s.CustomerService)

		// 同步更新用户行为事件记录表|更新群名
		gormx.GetInstance().Model(&micro_scrm.UserEventModel{}).
			Where("event", "join_work_chat_group").
			Where("object", chat.ChatID).
			Where("enterprise_id", enterpriseData.EnterpriseID).
			UpdateColumn("event_desc", chatDetail.Name)
	}
	// ====================群成员处理========================================
	// 新增群成员
	var batchChatMembers []customer.WorkGroupChatMember
	// 需要添加的群成员
	var batchAddchatMemberlogs []group_chat.WorkGroupChatMemberLog
	// 已经处理的群成员
	hasDealChatMembers := make(map[string]bool)
	for _, chatMemberItem := range existsInQwMembers {
		hasDealChatMembers[chatMemberItem.UserID] = true
		// 是企业内部成员
		if chatMemberItem.Type == customer.WorkChatGroupMemberTypeQw {
			if qwMember, ok := membersMapFromDb[chatMemberItem.UserID]; ok {
				chatMemberItem.MemberID = qwMember.ID
			}
		}
		chatMemberItem.GroupChatID = workGroupChatID
		chatMemberItem.EnterpriseID = enterpriseData.EnterpriseID
		chatMemberItem.Status = 1
		// 查询客户群成员是否存在于数据库
		if _, ok := hasExistWorkGroupMembersFromDb[chatMemberItem.UserID]; !ok {
			chatMemberItem.ID = uid.ID(ctx)
			chatMemberItem.CreatedAt = time.Now().Unix()
			chatMemberItem.UpdatedAt = time.Now().Unix()
			chatMemberItem.Status = customer.WorkGroupChatMemberJoinStatus
			batchChatMembers = append(batchChatMembers, *chatMemberItem)
			// 2021-08-27添加
			if _, ok := chatMemberlogs[chatMemberItem.UserID]; ok {
				batchAddchatMemberlogs = append(batchAddchatMemberlogs, group_chat.WorkGroupChatMemberLog{
					MicroScrmBaseModel: model.MicroScrmBaseModel{
						BaseModel: model.BaseModel{
							ID: uid.ID(ctx),
						},
						EnterpriseID: enterpriseData.EnterpriseID,
					},
					GroupChatId:     chatMemberItem.GroupChatID,
					ExternalUserid:  chatMemberItem.UserID,
					UnionId:         chatMemberItem.UnionID,
					OperationStatus: 1,
					OperationTime:   int(time.Now().Unix()),
					MemberType:      chatMemberItem.Type,
					JoinScene:       chatMemberItem.JoinScene,
					GroupNickname:   chatMemberlogs[chatMemberItem.UserID].GroupNickname,
					Name:            chatMemberlogs[chatMemberItem.UserID].Name,
					InvitorUserid:   chatMemberItem.InviteUserID,
				})
			}
		}
	}
	// 批量新增客户群成员
	if len(batchChatMembers) > 0 {
		_, aErr := s.CustomerService.AddWorkGroupChatMember(ctx, batchChatMembers)

		if aErr != nil {
			log.Error(ctx).Msgf("新增客户群成员数据失败，错误信息：%+v", aErr)
		} else {
			_ = s.batchAddChatEvent(ctx, batchChatMembers, chat, enterpriseData.EnterpriseID, common.EventEnum.JoinWorkChatGroup, callbackData.CreateTime)
		}

		data := entity.ExternalUserInfoEntity{
			EnterpriseID: enterpriseData.EnterpriseID,
			GroupChatID:  workGroupChatID,
			ToUserName:   callbackData.ToUserName,
		}
		if len(batchAddchatMemberlogs) > 0 {
			_ = orm.Service.PlanTask.CreatememberLos(batchAddchatMemberlogs)
		}
		_ = community_task.ChangeExternalChatAddMember(ctx, data)
	}

	// ====================退群群成员处理========================================

	var hasDeletedMembers []int64
	// 需要删除的群成员
	var batchDeletechatMemberlogs []group_chat.WorkGroupChatMemberLog
	var deleteMembers []customer.WorkGroupChatMember
	for key, value := range hasExistWorkGroupMembersFromDb {
		if _, ok := hasDealChatMembers[value.UserID]; !ok {
			hasDeletedMembers = append(hasDeletedMembers, value.ID)
			deleteMembers = append(deleteMembers, hasExistWorkGroupMembersFromDb[key])
			// 2021-08-27添加删除日志
			batchDeletechatMemberlogs = append(batchDeletechatMemberlogs, group_chat.WorkGroupChatMemberLog{
				MicroScrmBaseModel: model.MicroScrmBaseModel{
					BaseModel: model.BaseModel{
						ID: uid.ID(ctx),
					},
					EnterpriseID: enterpriseData.EnterpriseID,
				},
				GroupChatId:     workGroupChatID,
				ExternalUserid:  value.UserID,
				UnionId:         value.UnionID,
				OperationStatus: 2,
				OperationTime:   int(time.Now().Unix()),
				MemberType:      value.Type,
				JoinScene:       value.JoinScene,
				GroupNickname:   chatMemberlogs[value.UserID].GroupNickname,
				Name:            chatMemberlogs[value.UserID].Name,
				InvitorUserid:   value.InviteUserID,
			})
		}
	}
	if len(hasDeletedMembers) > 0 {
		num, dErr := s.CustomerService.BatchUpdateWorkGroupChatMembers(ctx, hasDeletedMembers, customer.WorkGroupChatMember{
			Status:       customer.WorkGroupChatMemberLeaveStatus,
			LeaveTime:    time.Now().Unix(),
			EnterpriseID: enterpriseData.EnterpriseID,
			UpdatedAt:    time.Now().Unix(),
		})

		if dErr != nil {
			log.Error(ctx).Msgf("批量删除已退群成员数据失败，错误信息：%+v", dErr)
		} else if num > 0 {
			_ = s.batchAddChatEvent(ctx, deleteMembers, chat, chat.EnterpriseID, common.EventEnum.LeaveWorkChatGroup, callbackData.CreateTime)
		}
		deletedata := entity.ExternalUserInfoEntity{
			EnterpriseID: enterpriseData.EnterpriseID,
			GroupChatID:  workGroupChatID,
			ToUserName:   callbackData.ToUserName,
		}
		if len(batchDeletechatMemberlogs) > 0 {
			_ = orm.Service.PlanTask.CreatememberLos(batchDeletechatMemberlogs)
		}
		_ = community_task.ChangeExternalChatDeleteMember(ctx, deletedata)
	}
	return nil
}

// batchAddChatEvent 批量设置用户群聊事件
func (s Service) batchAddChatEvent(ctx ctx.Context, members []customer.WorkGroupChatMember, chat *customer.WorkGroupChat, enterpriseID int64, event common.Event, eventTime int64) error {
	log.Info(ctx).Msgf("批量设置用户群聊事件，members:%v, chat:%v, event:%s, enterpriseID:%d", members, *chat, event.String(), enterpriseID)
	for _, v := range members {
		if v.UnionID == "" {
			continue
		}
		userEntity, err := s.CustomerService.QueryEnterpriseUser(ctx, enterpriseID, &customer.SourceUserBase{UnionID: v.UnionID})
		if err != nil {
			log.Info(ctx).Msgf("添加用户群聊事件：获取用户数据失败：err:%s, unionID:%s, enterpriseID:%d", err.Error(), v.UnionID, enterpriseID)
			continue
		}

		err = s.UserService.AddUserEvent(ctx, enterpriseID, userEntity.ID, event.String(), chat.Name, chat.ChatID, chat, eventTime)
		if err != nil {
			log.Info(ctx).Msgf("添加用户事件失败：err:%s, userInfo:%+v, chat:%+v", err.Error(), userEntity, chat)
			continue
		}
	}
	return nil
}

// SearchChatMemberInner 查询群里的内部成员
func (s Service) SearchChatMemberInner(ctx ctx.Context, memberUserIDs []string, enterpriseID int64) map[string]members.Member {
	var membersMap = make(map[string]members.Member)
	if len(memberUserIDs) > 0 {
		reqMembers := &members.MemberListRequest{
			UserIDs:      memberUserIDs,
			UserIDCount:  len(memberUserIDs),
			Page:         1,
			Size:         len(memberUserIDs),
			EnterpriseID: enterpriseID,
		}
		total, memberData, _ := s.MemberService.SimpleMemberList(ctx, reqMembers)
		if total > 0 {
			for _, v := range *memberData {
				membersMap[v.Userid] = v
			}
		}
	}
	return membersMap
}

// searchWorkChatMemberFromDb 从数据库查询客户群成员信息
func searchWorkChatMemberFromDb(ctx ctx.Context, workGroupChatID int64, enterpriseID int64, customerService *customer.Service) (map[string]customer.WorkGroupChatMember, map[string]*customer.WorkGroupChatMember) {
	workGroupMembersMap := make(map[string]customer.WorkGroupChatMember)
	leavedMembersMap := make(map[string]*customer.WorkGroupChatMember)
	// 查询已有的群成员
	batchReq := customer.BatchWorkGroupMemberRequest{
		Size:         0,
		EnterpriseID: enterpriseID,
		GroupChatID:  workGroupChatID,
	}
	total, workGroupMembers, _ := customerService.BatchFindWorkGroupChatMembers(ctx, batchReq)
	if total > 0 {
		for i, v := range workGroupMembers {
			if v.Status == WorkChatGroupMemberStatusJoin {
				workGroupMembersMap[v.UserID] = v
			} else {
				leavedMembersMap[v.UserID] = &workGroupMembers[i]
			}
		}
	}
	return workGroupMembersMap, leavedMembersMap
}

// AddFriendCallbackActivityHandle 添加好友回调裂变活动处理
func (s Service) AddFriendCallbackActivityHandle(ctx ctx.Context, callbackData *work.MemberExternalCallback) error {
	userService := userModule.NewUserService(
		userMysql.NewUser(),
		userMysql.NewEvent(),
		elastic.NewEsRepo(),
		enterpriseMysql.New(),
	)
	enterpriseService := enterprise.New(enterpriseMysql.New())
	memberService := members.New(membersMysql.New())
	customerService := customer.New(customerMysql.NewCustomer(), customerMysql.NewWorkUser(), elastic.NewEsRepo())
	fissionService := fission.NewService(fissionMysql.New(), drainage.NewContact(drainageMysql.NewDrainage()), userService, fissionBossApi.NewApi(), enterpriseService)

	fissionHandle := NewFissionHandle(fissionService, memberService, customerService, enterpriseService)

	log.Info(ctx).Msgf("CallbackHandleFission 获取到的数据：%s", *callbackData)

	var err error
	switch true {
	case callbackData.Event == ExternalCallbackEventOfChangeExternalContact &&
		callbackData.ChangeType == AddExternalContact:
		// 添加好友
		// 通过state 查询关联的活动和渠道
		err = fissionHandle.FissionAddUser(ctx, callbackData)
	case callbackData.Event == ExternalCallbackEventOfChangeExternalContact &&
		callbackData.ChangeType == DelExternalContact:
		// 员工删除好友
		fallthrough
	case callbackData.Event == ExternalCallbackEventOfChangeExternalContact &&
		callbackData.ChangeType == DelFollowUser:
		// 好友删除员工
		err = fissionHandle.FissionDelUser(ctx, callbackData)
	case callbackData.Event == ExternalCallbackEventOfChangeExternalChat && callbackData.ChangeType == "update":
		err = fissionHandle.FissionGroupChange(ctx, callbackData)
	}

	if err == nil {
		log.Info(ctx).Msg("CallbackHandleFission操作成功")
	} else {
		log.Error(ctx).Msgf("CallbackHandleFission 操作失败错误:%s", err)
	}

	return err
}
