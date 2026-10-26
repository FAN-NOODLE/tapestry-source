package callback

import (
	"errors"
	"micro-scrm/internal/logic/common"
	fission2 "micro-scrm/internal/orm/model/micro_scrm/fission"
	"strconv"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/modules/customer"
	"micro-scrm/internal/modules/fission"
	"micro-scrm/internal/pkg/sdk/work"
	"micro-scrm/internal/pkg/uid"
)

// FissionHandle 回调裂变活动处理逻辑
type FissionHandle struct {
	fissionService    FissionService
	memberService     MemberService
	customerService   CustomerService
	enterpriseService EnterpriseService
}

// FissionGroupChange 群变更
func (fh *FissionHandle) FissionGroupChange(c ctx.Context, data *work.MemberExternalCallback) error {
	/**
	1. 查询群关联的活动(对于删除成员需要查询出所有的，对于添加的成员需要查询出增加的)
	2. 查询删除群的成员和添加群的成员
	3. 对添加的成员进行操作
	4. 对删除的成员进行操作
	*/

	// 查询企业详情
	enterpriseID, err := fh.getEnterpriseID(c, data.ToUserName)
	if err != nil {
		return err
	}

	// 查询群聊关联的活动
	chat, err := fh.isRunFissionChat(c, data.ChatID, enterpriseID)
	if err != nil {
		return err
	}

	// 获取该群聊关联的正常的活动
	normalActivity, err := fh.getNormalActivityByChatID(c, chat.ID, enterpriseID)
	if err != nil {
		return err
	}

	// 获取该群聊关联的所有在兑换日期之前的活动
	allActivity, err := fh.withDelActivityByChatID(c, chat.ID, enterpriseID)
	if err != nil {
		return err
	}

	// 没有查询到活动或者是查询出错都返回
	if len(normalActivity) == 0 && len(allActivity) == 0 {
		return nil
	}

	err = fh.addGroupByActivity(c, normalActivity, chat.ID, enterpriseID)
	if err != nil {
		return err
	}
	err = fh.liveGroupByActivity(c, allActivity, chat.ID, enterpriseID)
	if err != nil {
		return err
	}

	return nil
}

// FissionAddUser 裂变活动，添加好友
func (fh *FissionHandle) FissionAddUser(c ctx.Context, data *work.MemberExternalCallback) error {
	/**
	1. 通过 state 获取活动的信息，判断活动是否仍在运行
	2. 查询员工和客户的信息，判断是否为删除后在添加情况
	3. 添加/编辑裂变记录表
	4. 更新统计数据
	*/

	activityUser, activity, err := fh.activityInfo(c, data.State)
	if err != nil || activityUser == nil || activity == nil {
		return err
	}
	enterpriseID := activity.EnterpriseID

	// 如果裂变活动是成为VIP
	if activity.FissionType == fission2.FA_FissionType_VIP {
		log.Info(c).Msgf("裂变活动：vip类型：%v", data)
		// 判断该外部联系人是否为VIP
		isVip, err := common.NewVipUser(c).SetUser(common.VipUserParam{
			EnterpriseID:   enterpriseID,
			ExternalUserID: data.ExternalUserID,
		}).IsVip()
		if err != nil {
			return err
		}

		// 如果已经是vip则不符合条件
		if isVip {
			return nil
		}
	}

	// 查询员工和外部联系人的信息
	customerInfo, err := fh.getCustomerInfo(c, data.ExternalUserID, enterpriseID)
	if err != nil || customerInfo == nil {
		return err
	}
	member, err := fh.getMemberInfo(c, data.UserID, enterpriseID)
	if err != nil || member == 0 {
		return err
	}
	memberID := member
	customerID := customerInfo.ID

	// 判断是否为删除后在添加情况
	recoveryRecord, err := fh.hasRecoveryFissionRecord(c, customerID, memberID)
	if err != nil {
		return err
	}
	if recoveryRecord != nil {
		// 删除在添加操作
		updateData := &fission.UpdateAssistantRecord{
			ID:              recoveryRecord.ID,
			AddMemberStatus: fission.AddMemberStatus,
			AddMemberAt:     data.CreateTime,
		}
		err := fh.fissionService.RecoveryRecordAndStatistics(c, recoveryRecord, activity.FissionType, updateData)
		if err != nil {
			return err
		}
		return nil
	}

	// 判断是否满足插入记录
	ok, err := fh.isRunAddRecord(c, activity, enterpriseID, customerID, memberID)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	// 插入裂变记录并更新统计信息
	err = fh.createAssistantRecordAndStatistics(c, activity, enterpriseID, customerID, memberID, activityUser.ID, data.CreateTime, customerInfo.ExternalUserID)
	if err != nil {
		return err
	}

	return nil
}

// FissionDelUser 裂变活动，删除操作
func (fh *FissionHandle) FissionDelUser(c ctx.Context, data *work.MemberExternalCallback) (err error) {
	/**
	1. 查询外部联系人和员工和企业的信息
	2. 查询是否存在添加关系的助力记录
	3. 存在则更新
	4. 更新统计数据
	*/

	// 查询企业ID
	enterpriseID, err := fh.getEnterpriseID(c, data.ToUserName)
	if err != nil {
		return err
	}

	// 查询员工和外部联系人的信息
	customerInfo, err := fh.getCustomerInfo(c, data.ExternalUserID, enterpriseID)
	if err != nil || customerInfo == nil {
		return err
	}
	member, err := fh.getMemberInfo(c, data.UserID, enterpriseID)
	if err != nil || member == 0 {
		return err
	}
	memberID := member
	customerID := customerInfo.ID

	// 查询是否存在添加关系的助力记录
	addRecord, err := fh.fissionService.GetAddRecordByRelation(c, customerID, memberID, enterpriseID)
	if err != nil {
		log.Error(c).Msgf("GetAddRecordByRelation 错误:%s", err)
		return
	}
	if addRecord == nil {
		log.Info(c).Msg("无添加记录无需处理")
		return
	}

	// 查询是否满足删除的条件
	isRunDel, activity, err := fh.isRunDel(c, addRecord.ActivityID, addRecord.EnterpriseID)
	if err != nil {
		log.Error(c).Msgf("查询是否满足删除的条件错误, %v", err)
		return err
	}
	if !isRunDel || activity == nil {
		log.Info(c).Msg("关联的活动无需不满足处理条件或者为查询到活动")
		return nil
	}

	// 更新记录
	updateData := &fission.UpdateAssistantRecord{
		ID:              addRecord.ID,
		AddMemberStatus: fission.DelRelationStatus,
	}
	err = fh.fissionService.RecoveryRecordAndStatistics(c, addRecord, activity.FissionType, updateData)
	if err != nil {
		return err
	}
	return nil
}

// getNormalActivityByChatID 获取正常的群聊ID
func (fh *FissionHandle) getNormalActivityByChatID(c ctx.Context, chatID, enterpriseID int64) ([]int64, error) {
	req := &fission.ActivityByChatIDReq{
		WithDel:       false,
		ChatID:        chatID,
		EnterpriseID:  enterpriseID,
		ActivityRunAt: time.Now().Unix(),
	}

	return fh.fissionService.GetActivityByChatID(c, req)
}

// withDelActivityByChatID 获取含有删除的活动ID
func (fh *FissionHandle) withDelActivityByChatID(c ctx.Context, chatID, enterpriseID int64) ([]int64, error) {
	req := &fission.ActivityByChatIDReq{
		WithDel:               true,
		ChatID:                chatID,
		EnterpriseID:          enterpriseID,
		ActivityExchangeRunAt: time.Now().Unix(),
	}

	return fh.fissionService.GetActivityByChatID(c, req)
}

// isRunDel 判断能否运行删除
func (fh *FissionHandle) isRunDel(c ctx.Context, activityID, enterpriseID int64) (bool, *fission.ActivityEntity, error) {
	activity, err := fh.fissionService.ActivityDetail(
		c,
		enterpriseID,
		activityID,
	)
	if err != nil {
		log.Error(c).Msgf("ActivityDetail 错误:%s", err)
		return false, nil, err
	}

	if !activity.RunExchange() {
		return false, nil, nil
	}

	return true, activity, nil
}

// getFirstJoinTime 获取最早的加入活动群的时间
func (fh *FissionHandle) getFirstJoinTime(c ctx.Context, activityID, enterpriseID int64, externalUserID string, withDel bool) (joinTime int64, err error) {
	groupIDs, err := fh.fissionService.GetGroupsIDByActivity(c, activityID, enterpriseID, withDel)
	if err != nil {
		return
	}
	if len(groupIDs) == 0 {
		return
	}
	joinTime, err = fh.customerService.WorkExternalInGroup(c, externalUserID, groupIDs, enterpriseID)

	return joinTime, err
}

// isRunAddRecord 判断是否能运行添加记录
func (fh *FissionHandle) isRunAddRecord(c ctx.Context, activity *fission.ActivityEntity, enterpriseID, customerID, memberID int64) (bool, error) {
	// 获取该活动关联的成员
	membersID, err := fh.fissionService.ActivityWithMembers(c, activity.ID, enterpriseID)
	if err != nil {
		log.Error(c).Msgf("ActivityWithMembers 错误:%s", err)
		return false, err
	}
	if len(membersID) == 0 {
		log.Error(c).Msgf("通过活动未查询到成员:%d", activity.ID)
		return false, errors.New("活动未查询到关联的裂变成员")
	}

	// 查询该活动关联的渠道成员在开始之间是否被添加过
	ok, err := fh.customerService.ExistsRelationByStartTime(c, &customer.ExistsRelationByTimeReq{
		CustomerID:   customerID,
		EnterpriseID: enterpriseID,
		MembersID:    membersID,
		StartTime:    activity.StartTime,
	})
	if err != nil {
		log.Error(c).Msgf("查询该活动关联的渠道成员在开始之间是否被添加过错误:%s", err)
		return false, err
	}
	if ok {
		log.Info(c).Msg("开始之间添加过该活动关联的员工")
		return false, nil
	}

	// 通过客户和员工的关系查询是否存在助力记录
	otherMembers := delItem(membersID, memberID)
	if len(otherMembers) <= 0 {
		return true, nil
	}

	has, err := fh.fissionService.ExistRelationWithMembers(
		c,
		customerID,
		otherMembers,
		enterpriseID,
	)
	if err != nil {
		log.Error(c).Msgf("通过客户和员工的关系查询是否存在助力记录错误:%s", err)
		return false, err
	}

	if has {
		log.Info(c).Msgf("该客户已经与其他裂变人员存在关系:%d", customerID)
		return false, nil
	}

	return true, nil
}

// createAssistantRecord 创建记录，返回有没有加分
func (fh *FissionHandle) createAssistantRecordAndStatistics(
	c ctx.Context,
	activity *fission.ActivityEntity,
	enterpriseID int64,
	customerID int64,
	memberID int64,
	activityUserID int64,
	createRelationTime int64,
	eternalUserID string,
) error {
	now := time.Now().Unix()
	record := &fission.AssistantRecordsEntity{
		ID:              uid.ID(c),
		EnterpriseID:    enterpriseID,
		ActivityID:      activity.ID,
		FissionUserID:   activityUserID,
		AssistantUserID: customerID,
		ExternalUserID:  eternalUserID,
		MemberID:        memberID,
		AddMemberStatus: fission.AddMemberStatus,
		AddMemberAt:     createRelationTime,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	// 当活动类型为群类型的时候需判断有没有加入
	if activity.IsGroupType() {
		// 考虑到存在用户之前就存在这些群的情况，因此在这个时候需要查询下
		firstJoinGroupTime, err := fh.getFirstJoinTime(c, activity.ID, enterpriseID, eternalUserID, false)
		if err != nil {
			log.Error(c).Msgf("获取客户最早入群时间错误:%s", err)
			return err
		}

		// 说明之前已经存在这些群聊中
		if firstJoinGroupTime > 0 {
			record.JoinChatGroupStatus = fission.JoinGroupStatus
			record.JoinChatGroupAt = firstJoinGroupTime
		}
	}

	record.Score = record.GetScore(activity.FissionType)
	if err := fh.fissionService.CreateAssistantRecordAndStatistics(c, record); err != nil {
		log.Error(c).Msgf("CreateAssistantRecord 错误:%s", err)
		return err
	}

	return nil
}

func (fh *FissionHandle) isJoinGroup(fissionType int) bool {
	return fissionType == fission.GroupFissionType
}

// activityInfo 查询活动相关的信息
func (fh *FissionHandle) activityInfo(c ctx.Context, state string) (*fission.ActivityUserEntity, *fission.ActivityEntity, error) {
	//通过 state 查询关联的活动和渠道
	contactsId, err := strconv.ParseInt(state, 10, 64)
	if err != nil {
		return nil, nil, err
	}

	// 通过联系我二维码获取fission_user记录
	activityUser, err := fh.fissionService.GetActivityUserByContactID(c, contactsId)
	if err != nil {
		log.Error(c).Msgf("GetActivityUserByContactID 错误：%d", contactsId)
		return nil, nil, err
	}
	if activityUser == nil {
		log.Error(c).Msgf("GetActivityUserByContactID 未查询到数据：%d", contactsId)
		return nil, nil, errors.New("根据联系我二维码未查询到发起人数据")
	}
	// 查询活动详情
	activity, err := fh.fissionService.ActivityDetail(c, activityUser.EnterpriseID, activityUser.ActivityID)
	if err != nil {
		return nil, nil, err
	}
	if activity == nil {
		log.Error(c).Msgf("根据联系我二维码查询到活动数据:%d", contactsId)
		return nil, nil, nil
	}
	// 判断活动是否在运行
	if !activity.RunTime() {
		log.Error(c).Msgf("二维码所对应的活动不在运行：%d", contactsId)
		return nil, nil, nil
	}

	return activityUser, activity, nil
}

// getCustomerInfo 获取外部联系人信息
func (fh *FissionHandle) getCustomerInfo(c ctx.Context, externalUserID string, enterpriseID int64) (*customer.WorkExternalUser, error) {
	user, err := fh.customerService.FindExternalByUserID(c, externalUserID, enterpriseID)
	if err != nil {
		log.Error(c).Msgf("GetExternalUserID 错误:%s", err)
		return nil, err
	}
	if user == nil {
		log.Error(c).Msgf("通过外部联系人ID未查询到信息:%s", externalUserID)
		return nil, nil
	}

	return user, nil
}

// getMemberInfo 获取员工信息
func (fh *FissionHandle) getMemberInfo(c ctx.Context, memberUserID string, enterpriseID int64) (int64, error) {
	member, err := fh.memberService.GetMemberIDByUserid(c, memberUserID, enterpriseID)
	if err != nil {
		log.Error(c).Msgf("GetMemberIDByUserid 错误:%s", err)
		return 0, nil
	}
	if member == nil {
		log.Error(c).Msgf("通过员工userid未查询到信息:%s", memberUserID)
		return 0, nil
	}

	return member.ID, nil
}

// delItem 删除切片中的元素
func delItem(data []int64, value int64) []int64 {
	valueKey := 0
	find := false
	for k, v := range data {
		if v == value {
			find = true
			valueKey = k
			break
		}
	}

	if find == false {
		return data
	}

	return append(data[0:valueKey], data[(valueKey+1):]...)
}

func (fh *FissionHandle) getEnterpriseID(c ctx.Context, corpID string) (int64, error) {
	return fh.enterpriseService.GetEnterpriseID(c, corpID)
}

// isRunFissionChat 判断群是否计算活动统计
func (fh *FissionHandle) isRunFissionChat(c ctx.Context, chatID string, enterpriseID int64) (*customer.WorkGroupChat, error) {
	req := customer.WorkGroupChat{
		EnterpriseID: enterpriseID,
		ChatID:       chatID,
	}
	detail, err := fh.customerService.FindWorkGroupChatDetail(c, req)
	if err != nil {
		return nil, err
	}

	// 如果是外部群的话才计算
	if detail == nil || !detail.IsOpenChat() {
		log.Info(c).Msgf("该群聊非外部群聊")
		return nil, nil
	}

	return detail, nil
}

// addGroupByActivity 加群操作
func (fh *FissionHandle) addGroupByActivity(c ctx.Context, normalActivity []int64, chatID, enterpriseID int64) error {
	// 获取在在群里的成员状态是未入群的记录
	records, err := fh.fissionService.FindGroupRecordsByChatID(c, normalActivity, chatID, enterpriseID, true)
	if err != nil {
		log.Error(c).Msgf("获取在在群里的成员状态是未入群的记录错误%s", err)
		return err
	}
	if len(records) == 0 {
		return nil
	}

	for _, record := range records {
		// 更新记录
		updateData := &fission.UpdateAssistantRecord{
			ID:                  record.ID,
			JoinChatGroupStatus: fission.JoinGroupStatus,
			JoinChatGroupAt:     record.JoinTime,
		}

		err = fh.fissionService.RecoveryRecordAndStatistics(c, record.AssistantRecordsEntity, fission.GroupFissionType, updateData)
		if err != nil {
			log.Error(c).Msgf("入群更新状态失败,失败原因:%s, 数据%#v", err, updateData)
		}
	}

	return nil
}

// liveGroupByActivity 离开群操作
func (fh *FissionHandle) liveGroupByActivity(c ctx.Context, activity []int64, chatID, enterpriseID int64) error {
	// 获取在未在群里的成员状态是已入群的记录
	records, err := fh.fissionService.FindGroupRecordsByChatID(c, activity, chatID, enterpriseID, false)
	if err != nil {
		log.Error(c).Msgf("获取在在群里的成员状态是未入群的记录错误, %v", err)
		return err
	}
	if len(records) == 0 {
		return nil
	}

	for _, record := range records {
		// 查询该条记录关联的活动是否仍在群里
		joinTime, err := fh.getFirstJoinTime(c, record.ActivityID, record.EnterpriseID, record.ExternalUserID, true)
		if err != nil {
			return err
		}
		// 说明还有别的群
		if joinTime != 0 {
			continue
		}

		// 更新记录
		updateData := &fission.UpdateAssistantRecord{
			ID:                  record.ID,
			JoinChatGroupStatus: fission.LeaveGroupStatus,
			JoinChatGroupAt:     0,
		}

		err = fh.fissionService.RecoveryRecordAndStatistics(c, record.AssistantRecordsEntity, fission.GroupFissionType, updateData)
		if err != nil {
			log.Error(c).Msgf("退出群更新状态失败,失败原因:%s, 数据%#v", err, updateData)
		}
	}

	return nil
}

// hasRecoveryFissionRecord 判断是否为删除后添加的记录
func (fh *FissionHandle) hasRecoveryFissionRecord(c ctx.Context, customerID, memberID int64) (*fission.AssistantRecordsEntity, error) {
	recoveryRecord, err := fh.fissionService.RecoveryFissionRecord(c, customerID, memberID)
	if err != nil {
		log.Error(c).Msgf("查询是否是删除后添加记录错误%s", err)
		return nil, err
	}

	if recoveryRecord != nil && recoveryRecord.AddMemberStatus == fission.AddMemberStatus {
		return nil, errors.New("数据错误查询出添加状态的数据")
	}

	return recoveryRecord, nil
}
