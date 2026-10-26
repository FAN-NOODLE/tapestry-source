package fission

import (
	"errors"
	"fmt"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/orm/model/micro_scrm/fission"
	fission2 "micro-scrm/internal/orm/service/fission"
	"micro-scrm/internal/pkg/uid"
	"strconv"
	"strings"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/fx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/lock"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/pkg/assist"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/errs"
	pkgUid "micro-scrm/internal/pkg/uid"
)

type Service struct {
	repo       Repository
	Contact    ContactInterface
	User       UserInterface
	BossApi    BossApiInterface
	Enterprise EnterpriseInterface
}

var UserService UserInterface

// GenerateContactFunc 封装GenerateContact方法，方便单元测试
var GenerateContactFunc = generateContact

func generateContact(s *Service, c ctx.Context, enterpriseID int64, channelID int64, domain string, memberIDs []int64, contactType int) (int64, error) {
	return s.GenerateContact(c, enterpriseID, channelID, domain, memberIDs, contactType)
}

// GetActivityByChatID 获取ID
func (s *Service) GetActivityByChatID(c ctx.Context, req *ActivityByChatIDReq) ([]int64, error) {
	return s.repo.GetActivityByChatID(c, req)
}

// FindGroupRecordsByChatID 根据群聊ID获取入群但没有群状态的记录或者退群但有群状态的数据
func (s Service) FindGroupRecordsByChatID(c ctx.Context, activityID []int64, chatID, enterpriseID int64, joinOrLeave bool) ([]*AssistantRecordsWithTime, error) {
	return s.repo.FindGroupRecordsByChatID(c, activityID, chatID, enterpriseID, joinOrLeave)
}

func (s Service) getScoreNum(action string) int {
	switch action {
	case ActionScore:
		return 1
	case ActionReduce:
		return 1
	case ActionRecovery:
		return 0
	default:
		panic("getScoreNum 参数错误")
	}
}

// GetFissionGroupIDs 获取活动下的社群
func (s *Service) GetGroupsIDByActivity(c ctx.Context, activityID int64, enterpriseID int64, withDel bool) ([]int64, error) {
	fissionGroup, err := s.repo.FissionGroupIDs(c, activityID, enterpriseID, withDel)
	if err != nil {
		return nil, err
	}
	var groupIDs []int64
	for _, group := range fissionGroup {
		groupIDs = append(groupIDs, group.ChatID)
	}

	return groupIDs, nil
}

// ExistRelationWithMembers 通过客户和员工的关系查询是否存在助力记录
func (s *Service) ExistRelationWithMembers(c ctx.Context, userID int64, members []int64, enterpriseID int64) (bool, error) {
	req := &ExistRelationReq{
		CustomerID:   userID,
		Members:      members,
		EnterpriseID: enterpriseID,
	}

	data, err := s.GetRecordByRelation(c, req)
	if err != nil {
		return false, err
	}

	if data == nil {
		return false, nil
	}

	return true, nil

}

// GetAddRecordByRelation 获取添加的裂变助力记录
func (s *Service) GetAddRecordByRelation(c ctx.Context, customerID, memberID, enterpriseID int64) (*AssistantRecordsEntity, error) {
	req := &ExistRelationReq{
		CustomerID:     customerID,
		Members:        []int64{memberID},
		EnterpriseID:   enterpriseID,
		RelationStatus: AddMemberStatus,
	}

	return s.GetRecordByRelation(c, req)
}

// GetActivityUserByContactID 通过联系我二维码获取fission_user记录
func (s *Service) GetActivityUserByContactID(c ctx.Context, ContactID int64) (*ActivityUserEntity, error) {
	return s.repo.GetActivityUserByContactID(c, ContactID)
}

// FissionActivityDetail 获取活动详情
func (s *Service) FissionActivityDetail(c ctx.Context, ActivityID int64, EnterpriseID int64) (*ActivityEntity, error) {
	return s.repo.ActivityDetail(c, ActivityID, EnterpriseID)
}

// ActivityWithMembers 通过活动查询助力人
func (s *Service) ActivityWithMembers(c ctx.Context, activityID, enterpriseID int64) ([]int64, error) {
	fissionMemberIDs, err := s.repo.FissionMemberIDs(c, activityID)
	if err != nil {
		return nil, err
	}
	if len(fissionMemberIDs) == 0 {
		return nil, err
	}
	var memberIDs []int64
	for _, member := range fissionMemberIDs {
		memberIDs = append(memberIDs, member.MemberID)
	}

	return memberIDs, nil
}

// RecoveryFissionRecord 通过裂变人ID，客户ID, 员工ID, 来查询是否有删除的时间
func (s *Service) RecoveryFissionRecord(c ctx.Context, customerID, memberID int64) (*AssistantRecordsEntity, error) {
	return s.repo.RecoveryFissionRecord(c, customerID, memberID)
}

// GetRecordByRelation 通过客户和员工的关系查询是否存在助力记录
func (s *Service) GetRecordByRelation(c ctx.Context, req *ExistRelationReq) (*AssistantRecordsEntity, error) {
	if len(req.Members) == 0 {
		return nil, errs.FissionQueryRelationParamErr
	}
	return s.repo.GetRecordByRelation(c, req)
}

func (s *Service) ActivityMembers(c ctx.Context, activityID, enterpriseID int64) ([]int64, error) {
	fissionMemberIDs, err := s.repo.FissionMemberIDs(c, activityID)
	if err != nil {
		return nil, err
	}
	if len(fissionMemberIDs) == 0 {
		return nil, err
	}

	var memberIDs []int64
	for _, member := range fissionMemberIDs {
		memberIDs = append(memberIDs, member.MemberID)
	}

	return memberIDs, nil
}

// PrizeList 获取活动奖品列表
func (s *Service) PrizeList(c ctx.Context, enterpriseID int64, activityID int64, page int64, size int64) (*PrizeList, error) {
	total, list, err := s.repo.PrizeList(c, enterpriseID, activityID, page, size)
	if err != nil {
		log.Info(c).Msgf("查询奖品列表失败, err:%s", err.Error())
		return nil, err
	}
	var priseList = PrizeList{
		PageType: common.PageType{
			Page:  page,
			Size:  size,
			Total: total,
		},
	}
	if list == nil {
		priseList.List = nil
	} else {
		priseList.List = *list
	}
	return &priseList, nil
}

// PrizeInfo 获取奖品详情
func (s *Service) PrizeInfo(c ctx.Context, ID int64) (*PrizeEntity, error) {
	return s.getPrizeInfo(c, ID, false)
}

// PrizeInfoWithDeleted 获取商品详情-查询已删除数据
func (s *Service) PrizeInfoWithDeleted(c ctx.Context, ID int64) (*PrizeEntity, error) {
	return s.getPrizeInfo(c, ID, true)
}

func (s *Service) getPrizeInfo(c ctx.Context, ID int64, withDeleted bool) (*PrizeEntity, error) {
	if ID == 0 {
		return nil, errs.FissionPrizeIDEmpty
	}
	prizeInfo, err := s.repo.PrizeInfo(c, ID, withDeleted)
	if err != nil {
		log.Info(c).Msgf("获取列表活动奖品信息错误, err:%s, ID:%d", err.Error(), ID)
		return nil, err
	}
	if prizeInfo == nil {
		return nil, errs.FissionPrizeNotExist
	}
	return prizeInfo, nil
}

// FissionActivityInfo 获取活动信息
func (s *Service) ByContactIdActivityDetail(c ctx.Context, contactID int64, enterpriseID int64) (*ActivityEntity, error) {
	log.Info(c).Msgf("获取活动信息请求参数, ID:%d, enterpriseID:%d", contactID, enterpriseID)
	if contactID == 0 {
		return nil, errs.FissionActivityNotExist
	}
	fissionActivityInfo, activityDetailErr := s.repo.ByContactIdActivityDetail(c, contactID, enterpriseID)

	if activityDetailErr != nil {
		return nil, activityDetailErr
	}
	if fissionActivityInfo == nil {
		return nil, errs.FissionPrizeNotExist
	}
	return fissionActivityInfo, nil
}

// GetActivityUser 获取活动用户信息
func (s *Service) GetActivityUser(c ctx.Context, enterpriseID int64, activityID int64, userID int64) (*ActivityUserEntity, error) {
	if userID == 0 {
		return nil, nil
	}
	activityUserInfo, err := s.repo.GetActivityUser(c, enterpriseID, activityID, userID)
	if err != nil {
		return nil, err
	}
	if activityUserInfo == nil {
		return nil, errs.FissionActivityUserNotExist
	}
	return activityUserInfo, nil
}

// GetMpUser 获取用户信息
func (s *Service) GetMpUser(c ctx.Context, ID int64) (*MpUser, error) {
	userInfo, err := s.repo.GetMpUserInfo(c, ID)
	if err != nil {
		log.Info(c).Msgf("获取用户信息失败, err:%s", err.Error())
		return nil, err
	}
	return userInfo, nil
}

// GetContactQrCode 通过二维码id获取二维码链接
func (s *Service) GetContactQrCode(c ctx.Context, ID int64, enterpriseID int64) (string, string, error) {
	return s.Contact.GetContactQrCode(c, ID, enterpriseID)
}

// ByActivityContactQrCode 通过活动Id获取二维码链接
func (s *Service) ByActivityContactQrCode(c ctx.Context, activityID int64, enterpriseID int64) (string, error) {
	chatGroupEntity, err := s.repo.ByActivityContactQrCode(c, activityID, enterpriseID)
	if err != nil {
		return "", err
	}
	if chatGroupEntity == nil {
		return "", nil
	}
	return chatGroupEntity.QrCode, nil
}

// IsUserInMutualActivity 检查用户是否参加了互斥的活动
func (s *Service) IsUserInMutualActivity(c ctx.Context, enterpriseID int64, activityID int64, userID int64) (bool, error) {
	if userID == 0 {
		return false, nil
	}
	mutualActivityIDList, err := s.GetMutualActivityIDList(c, enterpriseID, activityID)
	if err != nil {
		return false, err
	}
	if len(mutualActivityIDList) == 0 {
		return false, nil
	}

	userActivityList, err := s.repo.GetActivityUserList(c, userID, mutualActivityIDList, enterpriseID)
	if err != nil {
		log.Info(c).Msgf("查询用户参与活动信息失败, err:%s, userID:%d, activityIDList:%v", err.Error(), userID, mutualActivityIDList)
		return false, err
	}
	if len(userActivityList) > 0 {
		return true, nil
	}
	return false, nil
}

// GetMutualActivities 获取活动的互斥活动列表
func (s *Service) GetMutualActivityIDList(c ctx.Context, enterpriseID int64, activityID int64) ([]int64, error) {
	list, err := s.repo.MutualActivityIDs(c, activityID, enterpriseID)
	if err != nil {
		return nil, err
	}
	var idList []int64
	for _, v := range list {
		idList = append(idList, v.MutualActivityID)
	}
	return idList, nil
}

// JoinActivity 用户参加活动
func (s *Service) JoinActivity(c ctx.Context, enterpriseID int64, activityID int64, userID int64) (*ActivityUserEntity, error) {
	if enterpriseID == 0 || activityID == 0 || userID == 0 {
		return nil, errs.FissionJoinParamErr
	}
	// 获取活动信息
	activityInfo, err := s.FissionActivityDetail(c, activityID, enterpriseID)
	if err != nil {
		return nil, err
	}
	if activityInfo == nil {
		return nil, errs.FissionActivityNotExistErr
	}
	lockKey := fmt.Sprintf("JoinActivity_%d_%d", activityID, userID)
	mutex, err := lock.Mutex(c, lock.Key(lockKey), lock.TTL(5*time.Second))
	if err != nil {
		log.Info(c).Msgf("初始化锁失败, err:%s", err.Error())
		return nil, err
	}
	err = mutex.Lock()
	if err != nil {
		log.Info(c).Msgf("获取锁失败, err:%s", err.Error())
		return nil, err
	}
	defer mutex.Unlock()
	// 检查活动状态
	dbFissionUser, err := s.checkJoinActivity(c, activityInfo, userID)
	if err != nil {
		return nil, err
	}
	// 用户已参加活动，检查二维码是否正常
	if dbFissionUser != nil {
		if dbFissionUser.ContactID > 0 {
			return dbFissionUser, nil
		}
		newContactID, err := s.generateUserContact(c, activityInfo)
		if err != nil {
			return nil, err
		}
		dbFissionUser.ContactID = newContactID
		_, err = s.repo.UpdateFissionUser(c, dbFissionUser)
		if err != nil {
			log.Info(c).Msgf("更新用户信息失败, err:%s", err.Error())
			return nil, err
		}
		return dbFissionUser, nil
	}
	contactID, err := s.generateUserContact(c, activityInfo)
	if err != nil {
		return nil, err
	}
	fissionUser, err := s.userJoinActivity(c, enterpriseID, activityID, userID, contactID)
	if err != nil {
		return nil, err
	}
	_ = s.addActivityUserNum(c, activityID)

	//羽西参与裂变活动同步用户足迹
	//if isCan, authErr := auth.BrandAuth(c, auth.BrandAuthConfig{
	//	EnterpriseID: enterpriseID,
	//	Business:     auth.BUSINESS_USER_JOIN_FISSION_EVENT,
	//}); isCan && authErr == nil {
	//
	//}
	_, data, err := orm.Service.Fission.GetFissionMpUserById(enterpriseID, strconv.FormatInt(userID, 10))
	if data.UnionID == "" {
		return nil, errors.New("参与裂变用户信息不存在")
	}
	eModel, err := orm.Service.WorkExternalUser.GetWxUserCorpExternalInfoByUnionID(enterpriseID, data.UnionID)
	if eModel.Euid != 0 {
		model := micro_scrm.UserEventModel{
			UserId:    eModel.Euid,
			Event:     common.EventEnum.UserJoinFissionActivity.String(),
			EventDesc: strings.Join([]string{"参与", "裂变活动"}, activityInfo.Name),
			Object:    "",
			Extra:     "",
		}
		model.ID = uid.ID(c)
		model.EnterpriseID = enterpriseID
		err = orm.Service.UserEvent.CreateModel(model)
		if err != nil {
			return nil, err
		}
	}

	return fissionUser, nil
}

// addActivityUserNum 添加活动发起人数
func (s *Service) addActivityUserNum(c ctx.Context, activityID int64) error {
	_, err := s.repo.ActivityUserNumAdd(c, activityID, 1)
	if err != nil {
		log.Info(c).Msgf("更新活动发起人数据失败, err:%s", err.Error())
		return err
	}
	return nil
}

// generateUserContact 生成用户裂变联系我二维码
func (s *Service) generateUserContact(c ctx.Context, activity *ActivityEntity) (int64, error) {
	members, err := s.ActivityWithMembers(c, activity.ID, activity.EnterpriseID)
	if err != nil {
		return 0, err
	}
	domain, _, _, err := s.Enterprise.GetInfo(c, activity.EnterpriseID)
	if err != nil {
		return 0, err
	}
	return GenerateContactFunc(s, c, activity.EnterpriseID, activity.ChannelID, domain, members, ContactTypeFissionUser)
	//return s.GenerateContact(c, activity.EnterpriseID, activity.ChannelID, domain, members, ContactTypeFissionUser)
}

// checkJoinActivity 用户参加活动检查
func (s *Service) checkJoinActivity(c ctx.Context, activity *ActivityEntity, userID int64) (*ActivityUserEntity, error) {
	if activity == nil {
		return nil, errs.FissionActivityNotExistErr
	}
	err := s.checkActivity(activity)
	if err != nil {
		return nil, err
	}
	// 检查用户是否参加过当前活动
	userActivity, err := s.repo.GetActivityUser(c, activity.EnterpriseID, activity.ID, userID)
	if err != nil {
		return nil, err
	}
	if userActivity != nil {
		return userActivity, nil
	}
	// 检查是否参加了互斥活动
	if activity.MutualType == IsMutualType {
		joined, err := s.IsUserInMutualActivity(c, activity.EnterpriseID, activity.ID, userID)
		if err != nil {
			return nil, err
		}
		if joined {
			return nil, errs.FissionUserJoinMutual
		}
	}

	return nil, nil
}

func (s *Service) HasWorkExternalRelation(c ctx.Context, enterpriseID int64, workExternalUserID int64, memberIDs []int64) (bool, error) {
	return s.repo.HasWorkExternalRelation(c, enterpriseID, workExternalUserID, memberIDs)
}

// userJoinActivity 添加用户活动数据
func (s *Service) userJoinActivity(c ctx.Context, enterpriseID int64, activityID int64, userID int64, contactID int64) (*ActivityUserEntity, error) {
	if enterpriseID == 0 || activityID == 0 || userID == 0 {
		return nil, errs.FissionJoinParamErr
	}
	entity := ActivityUserEntity{
		ID:           pkgUid.ID(c),
		EnterpriseID: enterpriseID,
		ActivityID:   activityID,
		UserID:       userID,
		ContactID:    contactID,
		CreatedAt:    time.Now().Unix(),
	}
	_, err := s.repo.AddFissionUser(c, &entity)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// checkActivity 检查活动状态
func (s *Service) checkActivity(activity *ActivityEntity) error {
	// 活动状态检查
	if activity.Status != StatusOnShelf {
		return errs.FissionActivityOff
	}
	now := assist.Now().Unix()
	if activity.StartTime > now {
		return errs.FissionActivityNotStart
	}
	if activity.EndTime < now {
		return errs.FissionActivityEnd
	}
	return nil
}

// ExchangeInfo 兑换列表
func (s *Service) ExchangeList(c ctx.Context, enterpriseID int64, activityID int64, userID int64, page int64, size int64) (int64, []*ExchangeInfo, error) {
	var resList = make([]*ExchangeInfo, 0)
	total, exchangeList, err := s.repo.GetExchangeList(c, enterpriseID, activityID, userID, page, size)
	if err != nil {
		return 0, resList, err
	}
	var prizeIDList []int64
	for _, v := range *exchangeList {
		prizeIDList = append(prizeIDList, v.PrizeID)
	}
	prizeList, err := s.repo.PrizeListByIDs(c, prizeIDList)
	if err != nil {
		return 0, resList, err
	}
	var prizeMap = make(map[int64]*PrizeEntity)
	for _, v := range *prizeList {
		prizeMap[v.ID] = v
	}
	for _, v := range *exchangeList {
		var info = ExchangeInfo{
			ID:               v.ID,
			ExchangeTime:     v.CreatedAt,
			CostAssistantNum: v.CostAssistantNum,
			Num:              v.Num,
		}
		if prize, ok := prizeMap[v.PrizeID]; ok {
			info.PrizeName = prize.Name
			info.PrizePic = prize.Pic
		}
		resList = append(resList, &info)
	}
	return total, resList, nil
}

// ExchangeDetail 兑换详情
func (s *Service) ExchangeDetail(c ctx.Context, ID int64) (*ExchangeDetail, error) {
	if ID == 0 {
		return nil, errs.ParamsErr
	}
	exchangeInfo, err := s.repo.GetExchangeInfo(c, ID)
	if err != nil {
		return nil, err
	}
	prizeIDList := []int64{
		exchangeInfo.PrizeID,
	}
	prizeList, err := s.repo.PrizeListByIDs(c, prizeIDList)
	if err != nil {
		return nil, err
	}
	deliverName, _ := common.AesDecrypt(exchangeInfo.DeliverName)
	deliverMobile, _ := common.AesDecrypt(exchangeInfo.DeliverMobile)
	deliverArea, _ := common.AesDecrypt(exchangeInfo.DeliverArea)
	deliverAddress, _ := common.AesDecrypt(exchangeInfo.DeliverAddress)
	var detail = ExchangeDetail{
		ExchangeInfo: ExchangeInfo{
			ID:               exchangeInfo.ID,
			ExchangeTime:     exchangeInfo.CreatedAt,
			CostAssistantNum: exchangeInfo.CostAssistantNum,
			Num:              exchangeInfo.Num,
		},
		AddressInfo: AddressInfo{
			DeliverName:   deliverName,
			DeliverMobile: deliverMobile,
			Address:       deliverArea + deliverAddress,
		},
	}
	if prizeList != nil && len(*prizeList) >= 0 {
		detail.PrizePic = (*prizeList)[0].Pic
		detail.PrizeName = (*prizeList)[0].Name
	}
	return &detail, nil
}

// AssistantList 用户助力列表
func (s *Service) AssistantList(c ctx.Context, enterpriseID int64, activityID int64, userID int64, page int64, size int64) (*AssistantList, error) {
	var resList = AssistantList{
		List: make([]*AssistantItem, 0),
	}
	if userID == 0 {
		return &resList, nil
	}
	activityUser, err := s.GetActivityUser(c, enterpriseID, activityID, userID)
	if err != nil && err != errs.FissionActivityUserNotExist {
		return &resList, err
	}
	if err == errs.FissionActivityUserNotExist {
		return &resList, nil
	}
	if activityUser == nil {
		return &resList, nil
	}
	total, list, err := s.repo.GetAssistantList(c, activityID, activityUser.ID, page, size)
	if err != nil {
		return &resList, err
	}

	resList.CurAssistantNum = activityUser.CurAssistantNum
	resList.Total = total
	resList.Page = page
	resList.Size = size
	if list == nil {
		return &resList, nil
	}
	for _, v := range *list {
		resList.List = append(resList.List, &AssistantItem{
			AssistantUserID:     v.AssistantUserID,
			AddMemberStatus:     v.AddMemberStatus,
			JoinChatGroupStatus: v.JoinChatGroupStatus,
			Score:               v.Score,
			IsVip:               v.IsVip,
		})
	}
	err = s.setAssistantListWorkUserInfo(c, &resList, activityUser.EnterpriseID)
	return &resList, nil
}

// setAssistantListWorkUserInfo 设置助力列表外部联系人信息
func (s *Service) setAssistantListWorkUserInfo(c ctx.Context, list *AssistantList, enterpriseID int64) error {
	if list == nil {
		return nil
	}
	return s.setAssistantItemWorkUserInfo(c, &list.List, enterpriseID)
}

// setAssistantItemWorkUserInfo 设置助力列表外部联系人信息
func (s *Service) setAssistantItemWorkUserInfo(c ctx.Context, list *[]*AssistantItem, enterpriseID int64) error {
	if list == nil {
		return nil
	}
	var workUserIDs []int64
	for _, v := range *list {
		workUserIDs = append(workUserIDs, v.AssistantUserID)
	}
	workUserList, err := s.getWorkUserInfoByUserIDs(c, enterpriseID, workUserIDs)
	if err != nil {
		return err
	}
	if workUserList == nil {
		return nil
	}
	var userMap = make(map[int64]*WorkUserBasicInfo)
	for _, v := range *workUserList {
		userMap[v.UserID] = v
	}
	for _, v := range *list {
		if userInfo, ok := userMap[v.AssistantUserID]; ok {
			v.Name = userInfo.Name
			v.ImgUrl = userInfo.ImgURL
		}
	}
	return nil
}

// getWorkUserInfoByUserIDs 根据用户id查询外部联系人信息
func (s *Service) getWorkUserInfoByUserIDs(c ctx.Context, enterpriseID int64, userIDs []int64) (*[]*WorkUserBasicInfo, error) {
	return UserService.GetWorkUserInfoByUserIDs(c, enterpriseID, userIDs)
}

// ExchangePrizeCheck 兑换奖品前检查
func (s *Service) ExchangePrizeCheck(c ctx.Context, enterpriseID int64, activityID int64, userID int64, prizeID int64) (bool, error) {
	if activityID == 0 || enterpriseID == 0 || prizeID == 0 {
		return false, errs.ParamsErr
	}

	_, err := s.checkActivityWhenExchange(c, enterpriseID, activityID)
	if err != nil {
		return false, err
	}

	prizeInfo, err := s.PrizeInfo(c, prizeID)
	if err != nil {
		if err == errs.FissionPrizeNotExist {
			err = errs.FissionPrizeOff
		}
		return false, err
	}
	if userID == 0 {
		return false, errorx.NewErrorMsg(errs.FissionExchangeAssistNotEnoughCode, fmt.Sprintf(errs.FissionExchangeAssistNotEnoughMsg, prizeInfo.NeedNum))
	}
	if prizeInfo == nil {
		return false, errs.FissionPrizeOff
	}
	if prizeInfo.Stocks <= prizeInfo.ExchangeNum {
		return false, errs.FissionPrizeExchangeOut
	}
	userActivity, err := s.GetActivityUser(c, enterpriseID, activityID, userID)
	if err == errs.FissionActivityUserNotExist {
		return false, errorx.NewErrorMsg(errs.FissionExchangeAssistNotEnoughCode, fmt.Sprintf(errs.FissionExchangeAssistNotEnoughMsg, prizeInfo.NeedNum))
	}
	if err != nil {
		return false, err
	}
	if userActivity.CurAssistantNum < prizeInfo.NeedNum {
		return false, errorx.NewErrorMsg(errs.FissionExchangeAssistNotEnoughCode, fmt.Sprintf(errs.FissionExchangeAssistNotEnoughMsg, prizeInfo.NeedNum-userActivity.CurAssistantNum))
	}
	err = s.checkUserExchangeNum(c, prizeInfo, userID)
	if err != nil {
		return false, err
	}
	return true, nil
}

// checkActivityWhenExchange 用户兑换奖品时检查活动状态
func (s *Service) checkActivityWhenExchange(c ctx.Context, enterpriseID int64, activityID int64) (*ActivityEntity, error) {
	activityInfo, err := s.FissionActivityDetail(c, activityID, enterpriseID)
	if err != nil {
		return nil, err
	}
	if activityInfo == nil {
		return nil, errs.FissionActivityNotExist
	}
	now := time.Now().Unix()
	if now > activityInfo.ExchangeEndTime {
		return nil, errs.FissionPrizeExchangeEnd
	}
	return activityInfo, nil
}

// ExchangePrize 用户兑换奖品
func (s *Service) ExchangePrize(c ctx.Context, enterpriseID int64, activityID int64, userID int64, prizeID int64, address *AddressInfo) (bool, error) {
	if userID == 0 {
		return false, errs.FissionUserNotLogin
	}
	if activityID == 0 || enterpriseID == 0 || prizeID == 0 {
		return false, errs.ParamsErr
	}
	activityInfo, err := s.checkActivityWhenExchange(c, enterpriseID, activityID)
	if err != nil {
		return false, err
	}

	// 奖品上锁
	prizeMutex, err := s.exchangeLock(c, "prize", strconv.FormatInt(prizeID, 10))
	if err != nil {
		return false, err
	}
	defer prizeMutex.Unlock()
	prizeInfo, err := s.PrizeInfo(c, prizeID)
	if err != nil {
		return false, err
	}
	if prizeInfo == nil {
		return false, errs.FissionPrizeOff
	}
	if prizeInfo.Stocks <= prizeInfo.ExchangeNum {
		return false, errs.FissionPrizeExchangeOut
	}
	// 用户上锁
	userMutex, err := s.exchangeLock(c, "user", fmt.Sprintf("%d_%d", activityID, userID))
	if err != nil {
		return false, err
	}
	defer userMutex.Unlock()
	userActivity, err := s.GetActivityUser(c, enterpriseID, activityID, userID)
	if err != nil {
		return false, err
	}
	if userActivity.CurAssistantNum < prizeInfo.NeedNum {
		return false, errorx.NewErrorMsg(errs.FissionExchangeAssistNotEnoughCode, fmt.Sprintf(errs.FissionExchangeAssistNotEnoughMsg, prizeInfo.NeedNum-userActivity.CurAssistantNum))
	}
	err = s.checkUserExchangeNum(c, prizeInfo, userID)
	if err != nil {
		return false, err
	}
	now := time.Now().Unix()
	userActivity.CurAssistantNum -= prizeInfo.NeedNum
	userActivity.UpdatedAt = now
	prizeInfo.ExchangeNum++
	prizeInfo.UpdatedAt = now
	//加密用户信息
	deliverName, err := common.AesEncrypt(address.DeliverName)
	deliverMobile, err := common.AesEncrypt(address.DeliverMobile)
	addressStr, err := common.AesEncrypt(address.Address)
	area, err := common.AesEncrypt(address.Area)
	areaCode, err := common.AesEncrypt(address.AreaCode)
	if err != nil {
		return false, err
	}

	record := ExchangePrizeRecordsEntity{
		ID:               pkgUid.ID(c),
		EnterpriseID:     enterpriseID,
		ActivityID:       activityID,
		PrizeID:          prizeID,
		FissionUserID:    userActivity.ID,
		UserID:           userID,
		CostAssistantNum: prizeInfo.NeedNum,
		Num:              1,
		CreatedAt:        now,
		DeliverName:      deliverName,
		DeliverMobile:    deliverMobile,
		DeliverAddress:   addressStr,
		DeliverArea:      area,
		DeliverAreaCode:  areaCode,
		//DeliverName:      address.DeliverName,
		//DeliverMobile:    address.DeliverMobile,
		//DeliverAddress:   address.Address,
		//DeliverArea:      address.Area,
		//DeliverAreaCode:  address.AreaCode,
	}
	_, err = s.repo.UserExchange(c, userActivity, prizeInfo, &record)
	if err != nil {
		return false, err
	}
	go func() {
		_, _ = fx.RunSafe(func() error {
			userInfo, err := s.GetMpUser(c, userID)
			if err != nil {
				return err
			}
			if userInfo == nil {
				return nil
			}
			err = s.SendMessage(c, enterpriseID, activityID, activityInfo.Name, userInfo.NickName)
			if err != nil {
				log.Info(c).Msgf("发送兑换消息失败, err:%s", err.Error())
			}
			return nil
		})
	}()

	return true, nil
}

func (s *Service) exchangeLock(c ctx.Context, key string, ID string) (lock.Locker, error) {
	lockKey := fmt.Sprintf("exchangePrize_%s_%s", key, ID)
	mutex, err := lock.Mutex(c, lock.Key(lockKey), lock.TTL(5*time.Second))
	if err != nil {
		log.Info(c).Msgf("初始化锁失败, err:%s", err.Error())
		return nil, err
	}
	err = mutex.Lock()
	if err != nil {
		log.Info(c).Msgf("获取锁失败, err:%s", err.Error())
		return nil, err
	}
	return mutex, nil
}

// checkUserExchangeNum 检查用户是否达到兑换上限
func (s *Service) checkUserExchangeNum(c ctx.Context, prizeInfo *PrizeEntity, userID int64) error {
	if prizeInfo.LimitNum == 0 {
		return nil
	}
	num, err := s.getUserExchangeNum(c, userID, prizeInfo.ID)
	if err != nil {
		return err
	}
	if int(num) >= prizeInfo.LimitNum {
		return errorx.NewErrorMsg(errs.FissionUserExchangeLimitCode, fmt.Sprintf(errs.FissionUserExchangeLimitMsg, prizeInfo.LimitNum))
	}
	return nil
}

// getUserExchangeNum 获取用户兑奖数
func (s *Service) getUserExchangeNum(c ctx.Context, userID int64, prizeID int64) (int64, error) {
	return s.repo.GetUserExchangeNum(c, userID, prizeID)
}

// QueryUserLatestExchangeRecord 获取用户最近兑换的记录
func (s *Service) QueryUserLatestExchangeRecord(c ctx.Context, userID int64) (*ExchangePrizeRecordsEntity, error) {
	if userID == 0 {
		return nil, nil
	}
	return s.repo.QueryUserLatestExchangeRecord(c, userID)
}

// GetEnterpriseAuthInfo 获取企业授权信息
func (s *Service) GetEnterpriseAuthInfo(c ctx.Context, enterpriseID int64) (*ShopAuthInfo, error) {
	return s.BossApi.GetShopAuthInfo(c, enterpriseID, EnterpriseAuthTypeDefault)
}

// UpdateAssistantRecordAndStatistics 更新记录并且更新统计信息
func (s *Service) RecoveryRecordAndStatistics(c ctx.Context, assistantRecord *AssistantRecordsEntity, fissionType int, updateData *UpdateAssistantRecord) error {
	// 1. 更新积分
	scoreRecord := *assistantRecord
	if updateData.AddMemberStatus != 0 {
		scoreRecord.AddMemberStatus = updateData.AddMemberStatus
	}

	if updateData.JoinChatGroupStatus != 0 {
		scoreRecord.JoinChatGroupStatus = updateData.JoinChatGroupStatus
	}

	// 获取积分
	updateData.Score = scoreRecord.GetScore(fissionType)
	updateData.UpdateAt = time.Now().Unix()

	// 2. 更新统计数据
	commonStatistics := new(Statistics)
	// 与原记录进行对比判断是否需要增加裂变数
	if assistantRecord.AddMemberStatus == InitRelationStatus && updateData.AddMemberStatus == AddMemberStatus {
		commonStatistics.FissionNum = +1
	}
	if assistantRecord.JoinChatGroupStatus == InitGroupStatus && updateData.JoinChatGroupStatus == JoinGroupStatus {
		commonStatistics.JoinGroupNum = +1
	}

	// 与原记录进行对比判断是否需要增加/减少成员流失数
	if assistantRecord.AddMemberStatus == AddMemberStatus && updateData.AddMemberStatus == DelRelationStatus {
		commonStatistics.RunAwayNum = +1
	}
	if assistantRecord.AddMemberStatus == DelRelationStatus && updateData.AddMemberStatus == AddMemberStatus {
		commonStatistics.RunAwayNum = -1
	}

	// 与原记录进行对比判断是否需要增加/减少群流失数
	if assistantRecord.JoinChatGroupStatus == JoinGroupStatus && updateData.JoinChatGroupStatus == LeaveGroupStatus {
		commonStatistics.LeaveGroupNum = +1
	}
	if assistantRecord.JoinChatGroupStatus == LeaveGroupStatus && updateData.JoinChatGroupStatus == JoinGroupStatus {
		commonStatistics.LeaveGroupNum = -1
	}

	fissionUserStatistics := &UpdateStatisticsData{
		Statistics: *commonStatistics,
	}

	// 与与原记录进行对比判断积分
	if assistantRecord.Score != updateData.Score {
		if updateData.Score == 0 {
			fissionUserStatistics.TotalAssistantNum = -1
			fissionUserStatistics.CurAssistantNum = -1
		} else {
			fissionUserStatistics.TotalAssistantNum = +1
			fissionUserStatistics.CurAssistantNum = +1
		}
	}

	return s.repo.RecoveryRecordAndStatistics(c, assistantRecord, fissionType, updateData, fissionUserStatistics)
}

// CreateAssistantRecordAndStatistics 添加记录并且更新统计信息
func (s *Service) CreateAssistantRecordAndStatistics(c ctx.Context, record *AssistantRecordsEntity) error {
	// 计算统计信息
	commonStatistics := new(Statistics)
	commonStatistics.FissionNum = 1
	commonStatistics.UpdateAt = time.Now().Unix()
	if record.JoinChatGroupStatus == JoinGroupStatus {
		commonStatistics.JoinGroupNum = 1
	}
	fissionUserStatistics := &UpdateStatisticsData{
		Statistics: *commonStatistics,
	}
	if record.HasScore() {
		fissionUserStatistics.TotalAssistantNum = 1
		fissionUserStatistics.CurAssistantNum = 1
	}
	return s.repo.CreateAssistantRecordAndStatistics(c, record, fissionUserStatistics)
}

// GetScrollList 获取滚动弹幕列表数据
func (s *Service) GetScrollList(c ctx.Context, enterpriseID int64, activityID int64, page int64, size int64) (int64, []*AssistantUserInfo, error) {
	activityInfo, err := s.ActivityDetail(c, enterpriseID, activityID)
	if err != nil {
		return 0, nil, err
	}
	needGroup := activityInfo.FissionType == GroupFissionType
	total, err := s.repo.GetTotalAssistantUser(c, enterpriseID, activityID, needGroup)
	if err != nil {
		return 0, nil, err
	}
	addRobot := (assist.Now().Unix() - activityInfo.StartTime) > 60*ScrollListAddRobotTime
	index := (page - 1) * size
	if !addRobot {
		list, err := s.getScrollAssistantUsers4Page(c, activityID, enterpriseID, index, size, needGroup)
		return total, list, err
	}
	end := page * size
	total += ScrollListRobotNum
	// 全部取真实数据
	if index > ScrollListRobotNum {
		index -= ScrollListRobotNum
		list, err := s.getScrollAssistantUsers4Page(c, activityID, enterpriseID, index, size, needGroup)
		return total, list, err
	}
	robotList, err := UserService.GetRobotUserList(c, ScrollListRobotNum)
	if err != nil {
		return 0, nil, err
	}
	if len(robotList) == 0 {
		total -= ScrollListRobotNum
		list, err := s.getScrollAssistantUsers4Page(c, activityID, enterpriseID, index, size, needGroup)
		return total, list, err
	}
	// 全部取假数据
	if end <= ScrollListRobotNum {
		return total, robotList[index:end], nil
	}
	// 假数据和真实数据混合
	var resList = make([]*AssistantUserInfo, 0)
	resList = append(resList, robotList[index:]...)
	userList, err := s.getScrollAssistantUsers4Page(c, activityID, enterpriseID, 0, end-ScrollListRobotNum, needGroup)
	if err != nil {
		return 0, nil, err
	}
	resList = append(resList, userList...)
	return total, resList, nil
}

func (s *Service) getScrollAssistantUsers4Page(c ctx.Context, activityID int64, enterpriseID int64, index int64, size int64, needGroup bool) ([]*AssistantUserInfo, error) {
	var resList = make([]*AssistantUserInfo, 0)
	list, err := s.repo.GetAssistantUserList(c, enterpriseID, activityID, index, size, needGroup)
	if err != nil {
		return resList, err
	}

	var workUserIDs []int64
	for _, v := range list {
		workUserIDs = append(workUserIDs, v.AssistantUserID)
	}
	workUserList, err := s.getWorkUserInfoByUserIDs(c, enterpriseID, workUserIDs)
	if err != nil {
		return resList, err
	}
	var userMap = make(map[int64]*WorkUserBasicInfo)
	if workUserList != nil {
		for _, v := range *workUserList {
			userMap[v.UserID] = v
		}
	}
	for _, v := range list {
		if userInfo, ok := userMap[v.AssistantUserID]; ok {
			resList = append(resList, &AssistantUserInfo{
				Name:   userInfo.Name,
				ImgUrl: userInfo.ImgURL,
			})
		}
	}
	return resList, nil
}

//GetStyleConfigList 获取配色列表
func (s *Service) GetStyleConfigList(enterpriseID int64) ([]fission2.FissionStyleConfig, error) {
	list, err := orm.Service.Fission.GetStyleConfigList(enterpriseID)
	if nil != err {
		return list, err
	}
	return list, nil
}

//GetBgConfigList
func (s *Service) GetBgConfigList(enterpriseID int64) ([]fission2.FissionBgConfig, error) {
	list, err := orm.Service.Fission.GetBgConfigList(enterpriseID)
	if nil != err {
		return list, err
	}
	return list, nil
}

//GetStyleConfig 获取指定企业当前配色
func (s *Service) GetStyleConfig(enterpriseID int64) (fission2.FissionStyleConfig, error) {
	list, err := orm.Service.Fission.GetCheckStyle(enterpriseID)
	if nil != err {
		return list, err
	}
	return list, nil
}

//GetStyleConfig 获取指定企业当前配色
func (s *Service) GetBgConfig(enterpriseID int64) (fission.FissionBgConfigModel, error) {
	list, err := orm.Service.Fission.GetBgConfig(enterpriseID)
	if nil != err {
		return list, err
	}
	return list, nil
}

//GetStyleConfigList 获取配色列表
func (s *Service) AddStyleConfig(model fission.FissionStyleConfigModel) error {
	err := orm.Service.Fission.AddStyleConfig(model)
	if nil != err {
		return err
	}
	return nil
}

//GetStyleConfigList 获取配色列表
func (s *Service) DelStyleConfig(enterpriseID int64, id []int64) error {
	err := orm.Service.Fission.DelStyleConfig(enterpriseID, id)
	if nil != err {
		return err
	}
	return nil
}

//SetStyleConfig 设置选中配色
func (s *Service) SetStyleConfig(enterpriseID, id int64) error {
	err := orm.Service.Fission.SetCheckStyle(enterpriseID, id)
	if nil != err {
		return err
	}
	return nil
}

//SetBgConfig
func (s *Service) SetBgConfig(enterpriseID int64, bgImg, shareImg string) error {

	model, err := orm.Service.Fission.GetBgConfig(enterpriseID)
	if model.ID != 0 {
		err = orm.Service.Fission.SetBgConfig(enterpriseID, bgImg, shareImg)
		if nil != err {
			return err
		}
	} else {
		m := fission.FissionBgConfigModel{
			EnterpriseID: enterpriseID,
			BgImg:        bgImg,
			ShareImg:     shareImg,
		}
		err = orm.Service.Fission.AddBgConfig(m)
		if nil != err {
			return err
		}
	}
	return nil
}
