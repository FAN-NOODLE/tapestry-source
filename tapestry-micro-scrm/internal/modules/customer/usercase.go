package customer

import (
	"database/sql"
	"encoding/json"
	"fmt"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/helper"
	"reflect"
	"strings"
	"sync"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/config"
	"micro-scrm/internal/modules/enterprise"
	"micro-scrm/internal/modules/shop_member"
	"micro-scrm/internal/modules/tag"
	"micro-scrm/internal/pkg/assist"
	"micro-scrm/internal/pkg/errorIgnore"
	"micro-scrm/internal/pkg/errs"
	"micro-scrm/internal/pkg/sdk/work"
	"micro-scrm/internal/pkg/uid"
)

// Service  应用服务
type Service struct {
	repo            Repository
	workRepo        WorkUserRepository
	esRepo          EsRepository
	enterpriseRepo  enterprise.Repository
	tagRepo         tag.Repository
	memberLevelRepo shop_member.Repository
}

var DrainageContact DrainageContactInterface

var GuideService GuideInterface

// DrainageContactInterface 渠道接口
type DrainageContactInterface interface {
	ListChannelByQrCodeIds(ctx ctx.Context, qrCodeIds []int64, enterpriseID int64, page int, size int) (int64, map[int64]ListChannelByQrCodeIdsResponse)
}

type GuideInterface interface {
	FindGuideByUserIDForCustomer(ctx ctx.Context, enterpriseID int64, userID int64) (*PortraitBindGuide, error)
}

func ListFactory(context ctx.Context, i int64) interface{} {
	return &ResponseFootprintList{}
}

const (
	// 后台管理查询信息时，不需要传递member_id. 用改常量判断是否来自哪个端。RequestCustomer中添加ChannelType来判断
	ChannelSuperAdmin = 1
)

// FootprintList 客户足迹
func (s Service) FootprintList(c ctx.Context, params *RequestCustomer) (int64, ResponseFootprintList, error) {
	// 需要查询用户的外部联系人user_id
	enterpriseUser, err := microScrmDao.NewWxUserCorpExternal().GetOne(c, map[string]any{"id": params.UserID}, params.EnterpriseID)

	// 未查询到外部联系人关系数据
	if enterpriseUser.ID == 0 {
		return 0, nil, errs.FindWorkUserParamErr
	}

	params.UserID = enterpriseUser.Euid

	// 先判断当前的外部联系和员工是否有关系
	if params.ChannelType != ChannelSuperAdmin {
		relation, fErr := s.repo.FindCustomerRelation(c, params.MemberID, enterpriseUser.Euid, params.EnterpriseID)
		// 查询外部联系人关系出错
		if fErr != nil {
			return 0, nil, fErr
		}

		// 未查询到外部联系人关系数据
		if relation == nil {
			return 0, nil, errs.ExternalRelationNotExistErr
		}
	}

	pageVo := assist.PageVO(params.Page, params.Size)
	option := &repos.FindPageOption{
		TplPrefix: "t_customer/footprint-list",
		Fn:        ListFactory,
	}
	total, list, err := s.repo.FindPage(c, pageVo, params, option)
	if errorIgnore.Ignore(err) != nil {
		return 0, nil, err
	}

	var eventMap = map[string]string{
		"first_visit_shop":               "首次访问小程序商城",
		"add_external_user":              "添加好友",
		"delete_external_user":           "删除好友",
		"shop_user_join_activity":        "参加活动",
		"pay_order":                      "支付订单",
		"join_work_chat_group":           "加入群聊",
		"leave_work_chat_group":          "退出群聊",
		"subscribe_official":             "关注公众号",
		"unsubscribe_official":           "取消关注公众号",
		"delete_external_user_by_member": "被员工删除",
		"view_moment":                    "查看朋友圈",
		"user_join_fission_activity":     "参与裂变活动",
	}
	resList, ok := list.(*ResponseFootprintList)
	if !ok {
		// 处理类型转换失败
		return 0, nil, nil
	}

	for _, val := range *resList {
		if eventName, found := eventMap[val.Event]; found {
			val.EventTxt = eventName
		}
	}
	return total, *resList, nil
}

// CustomerDetail 客户详情
func (s Service) CustomerDetail(c ctx.Context, params *RequestCustomer) (*ResponseCustomerDetail, error) {
	userID := params.UserID

	enterpriseUser, err := s.repo.QueryEnterpriseUser(c, params.EnterpriseID, &SourceUserBase{EuID: userID})
	if err != nil {
		return nil, err
	}
	resDetail := &ResponseCustomerDetail{
		UserID:       enterpriseUser.ID,
		UnionID:      enterpriseUser.UnionId,
		UserName:     enterpriseUser.Name,
		Nickname:     enterpriseUser.Nick,
		Avatar:       enterpriseUser.Avatar,
		Gender:       int(enterpriseUser.Sex),
		Country:      enterpriseUser.Country,
		Province:     enterpriseUser.Province,
		City:         enterpriseUser.City,
		Remark:       enterpriseUser.Remark,
		IsWorkUser:   len(enterpriseUser.WxExternal) > 0,
		ShowOfficial: len(enterpriseUser.WxOpen) > 0,
		Birthday:     enterpriseUser.Birthday,
		UnionCode:    enterpriseUser.UnionCode,
		Point:        enterpriseUser.Point,
		Grade:        enterpriseUser.Grade,
		StoreCode:    enterpriseUser.StoreCode,
		BACode:       enterpriseUser.BACode,
		Mobile:       enterpriseUser.Mobile,
	}

	if len(enterpriseUser.Tags) > 0 {
		if tags, err := s.tagRepo.QueryWorkExternalTagList(c, enterpriseUser.Tags, 0); err == nil {
			for _, tag := range tags {
				resDetail.TagList = append(resDetail.TagList, tag.TagName)
			}
		}
	}

	if len(enterpriseUser.Groups) > 0 {
		if groups, err := s.tagRepo.QueryUserTagGroups(c, enterpriseUser.Groups); err == nil {
			for _, group := range groups {
				resDetail.GroupList = append(resDetail.GroupList, group.Name)
			}
		}
	}

	remarkList := strings.Split(resDetail.Remark, ",")
	remarkList = helper.DelItems[string](remarkList, []string{""})
	var officialAccountList = make([]*DetailOfficialAccount, 0)
	if resDetail.ShowOfficial {
		officialAccountList, err = s.repo.GetCustomerOfficialAccounts(c, params.EnterpriseID, enterpriseUser.WxOpen)
		if err != nil {
			log.Info(c).Msgf("获取用户公众号数据失败,err:%s, enterpriseUser.WxOpen:%+v", err.Error(), enterpriseUser.WxOpen)
		}
	}

	resDetail.OfficialList = officialAccountList
	resDetail.RemarkList = remarkList
	shopLevel, err := s.memberLevelRepo.GetShopMemberLevel(ctx.Wrap(c), params.EnterpriseID, enterpriseUser.LevelID)
	if err != nil {
		log.Info(c).Msgf("获取会员登记失败,err:%s, %d", err.Error(), enterpriseUser.LevelID)
	} else {
		resDetail.LevelName = shopLevel.LevelName
	}

	return resDetail, nil
}

// BelongToMembers 所属员工列表
func (s Service) BelongToMembers(c ctx.Context, params *RequestCustomer) (int64, *ResponseMemberList, error) {
	enterpriseUser, err := microScrmDao.NewWxUserCorpExternal().GetOne(c, map[string]any{"id": params.UserID}, params.EnterpriseID)
	list := &ResponseMemberList{}
	if err != nil {
		return 0, list, err
	}

	pageVo := assist.PageVO(params.Page, params.Size)
	total, list, err := s.repo.BelongToMembers(c, []int64{enterpriseUser.ID}, pageVo)
	if errorIgnore.Ignore(err) != nil || total == 0 {
		return 0, list, err
	}

	// 获取数据
	var memberIDs []int64
	var memberUserIDs []string
	for _, vItem := range *list {
		memberIDs = append(memberIDs, vItem.MemberId)
		memberUserIDs = append(memberUserIDs, vItem.MemberUserid)
	}
	// 所属部门
	belongToDepartments, _ := s.repo.MemberMultiDepartments(c, memberIDs)

	// 员工所属群
	pageVo = assist.PageVO(0, 999)
	_, belongToChatList, err := s.repo.BelongToChats(c, params.EnterpriseID, memberUserIDs, memberIDs, WorkGroupChatMemberJoinStatus, pageVo)
	if errorIgnore.Ignore(err) != nil {
		return 0, nil, err
	}

	belongToChatListByMemberId := make(map[int64]BelongToGroupChatList)

	if belongToChatList != nil {
		for _, v := range *belongToChatList {
			memberID := v.MemberId
			belongToChatListByMemberId[memberID] = append(belongToChatListByMemberId[memberID], v)
		}
	}

	for k, v := range *list {
		memberID := v.MemberId
		if chatList, ok := belongToChatListByMemberId[memberID]; ok {
			(*list)[k].BelongToGroupChatList = chatList
			(*list)[k].BelongToGroupChatCount = int64(len(chatList))
		}

		if departments, ok := belongToDepartments[memberID]; ok {
			(*list)[k].DepartmentList = departments
		}
		if v.Remark == "" {
			(*list)[k].RemarkList = []string{}
		} else {
			(*list)[k].RemarkList = helper.DelItems(strings.Split(v.Remark, ","), []string{""})
		}
	}

	return total, list, nil
}

// BelongToChats 所属群工列表
func (s Service) BelongToChats(c ctx.Context, params *RequestCustomerGroup) (int64, *BelongToGroupChatList, error) {
	enterpriseUser, err := microScrmDao.NewWxUserCorpExternal().GetOne(c, map[string]any{"id": params.UserID}, params.EnterpriseID)

	list := &BelongToGroupChatList{}
	if err != nil && err != sql.ErrNoRows || enterpriseUser.ID == 0 {
		return 0, nil, err
	}

	pageVo := assist.PageVO(params.Page, params.Size)
	var memberIDs []int64
	var total int64

	total, list, err = s.repo.BelongToChats(c, params.EnterpriseID, []string{}, memberIDs, params.Status, pageVo)
	if errorIgnore.Ignore(err) != nil {
		return 0, nil, err
	}

	// 查询群成员
	if total > 0 {
		var chatIds []int64
		for _, chat := range *list {
			chatIds = append(chatIds, chat.GroupChatId)
		}
		// 查询客户群成员的userid
		req := BatchWorkGroupMemberRequest{
			GroupChatIds:     chatIds,
			EnterpriseID:     params.EnterpriseID,
			GroupChatIdCount: int64(len(chatIds)),
		}
		cTotal, chatMembers, _ := s.repo.BatchFindWorGroupChatMember(c, req)
		chatMembersMap := make(map[int64][]WorkGroupChatMember)
		if cTotal > 0 {
			for _, value := range chatMembers {
				chatMembersMap[value.GroupChatID] = append(chatMembersMap[value.GroupChatID], value)
			}
		}

		for k, item := range *list {
			(*list)[k].GroupMember = chatMembersMap[item.GroupChatId]
		}
	}

	return total, list, nil
}

// FindExternal 查询外部联系人
func (s Service) FindExternal(c ctx.Context, req WorkExternalUser) (*WorkExternalUser, error) {
	if req.EnterpriseID == 0 {
		return nil, errs.EnterpriseIdEmpty
	}
	if req.ID == 0 && req.ExternalUserID == "" && req.UnionID == "" {
		return nil, errs.FindWorkUserParamErr
	}
	return s.repo.FindExternal(c, req)
}

func (s Service) FindExternalByUserID(c ctx.Context, externalUserID string, enterpriseID int64) (*WorkExternalUser, error) {
	req := WorkExternalUser{
		ExternalUserID: externalUserID,
		EnterpriseID:   enterpriseID,
		Type:           1,
	}

	detail, err := s.FindExternal(c, req)
	switch true {
	case err == sql.ErrNoRows || detail == nil:
		return nil, nil
	case err != nil:
		return nil, err
	default:
		return detail, nil
	}
}

// ExistsRelationByStartTime 查询这些成员在开始时间之前是否别添加过客户
func (s Service) ExistsRelationByStartTime(c ctx.Context, req *ExistsRelationByTimeReq) (bool, error) {
	return s.repo.ExistsRelationByStartTime(c, req)
}

// FindExternalRelation 查询外部联系人关系
func (s Service) FindExternalRelation(c ctx.Context, req WorkUserMemberRelation) (*WorkUserMemberRelation, error) {
	relation, err := s.repo.FindExternalRelation(c, req)
	if err != nil {
		return nil, err
	}
	return relation, nil
}

// AddExternalRelation 外部联系人关系新增
func (s Service) AddExternalRelation(c ctx.Context, req []WorkUserMemberRelation) error {
	_, err := s.repo.AddExternalRelation(c, req)
	if err != nil {
		return err
	}
	return nil
}

// AddExternalRelation 外部联系人关系新增
func (s Service) AddExternalRelationLog(c ctx.Context, req []WorkUserMemberRelation) error {
	_, err := s.repo.AddExternalRelationLog(c, req)
	if err != nil {
		return err
	}
	return nil
}

// UpdateExternalRelation 更新外部联系人关系
func (s Service) UpdateExternalRelation(c ctx.Context, req WorkUserMemberRelation) error {
	_, err := s.repo.UpdateExternalRelation(c, req)
	if err != nil {
		return err
	}
	return nil
}

func (s Service) BatchDeleteUpdateExternalReleation(c ctx.Context, enterpriseID int64, lastUpdateAt int64) error {
	_, err := s.repo.BatchDeleteUpdateExternalReleation(c, enterpriseID, lastUpdateAt)
	if err != nil {
		return err
	}
	return nil
}

// BatchTagGroupByName 根据标签组名称批量查询标签组
func (s Service) BatchTagGroupByName(ctx ctx.Context, groupNames []string, enterpriseID int64) ([]WorkExternalTagGroup, error) {
	return s.repo.BatchTagGroupByName(ctx, groupNames, enterpriseID)
}

// FindTagGroupFirst 企业微信查询标签组信息
func (s Service) FindTagGroupFirst(ctx ctx.Context, req WorkExternalTagGroup) (*WorkExternalTagGroup, error) {
	return s.repo.FindTagGroupFirst(ctx, req)
}

// UpdateTagGroup 企业微信标签组更新
func (s Service) UpdateTagGroup(ctx ctx.Context, req WorkExternalTagGroup) error {
	return s.repo.UpdateTagGroup(ctx, req)
}

// FindWorkGroupChatDetail 查询群聊详情
func (s Service) FindWorkGroupChatDetail(ctx ctx.Context, req WorkGroupChat) (*WorkGroupChat, error) {
	return s.repo.FindWorkGroupChatDetail(ctx, req)
}

// AddWorkGroupChat 新增群聊
func (s Service) AddWorkGroupChat(ctx ctx.Context, req WorkGroupChat) (int64, error) {
	return s.repo.AddWorkGroupChat(ctx, req)
}

// UpdateWorkGroupChat 新增群聊
func (s Service) UpdateWorkGroupChat(ctx ctx.Context, req WorkGroupChat) error {
	return s.repo.UpdateWorkGroupChat(ctx, req)
}

// FindWorkGroupChatMember 查询群成员
func (s Service) FindWorkGroupChatMember(ctx ctx.Context, req WorkGroupChatMember) (*WorkGroupChatMember, error) {
	return s.repo.FindWorkGroupChatMember(ctx, req)
}

// AddWorkGroupChatMember 新增群成员
func (s Service) AddWorkGroupChatMember(ctx ctx.Context, req []WorkGroupChatMember) (int64, error) {
	return s.repo.AddWorkGroupChatMember(ctx, req)
}

// WorkExternalInGroup 查询人员是否在群里
func (s Service) WorkExternalInGroup(ctx ctx.Context, externalUserID string, groupIDs []int64, enterpriseID int64) (int64, error) {
	return s.repo.FindFirstJoinTimeInGroups(ctx, externalUserID, groupIDs, enterpriseID)
}

// Externals 外部联系人列表
func (s Service) Externals(c ctx.Context, params ExternalsRequest) (int64, []WorkExternalUser, error) {
	offset := (params.Page - 1) * params.Size
	pageVo := repos.PageVO{Offset: offset, Limit: params.Size}

	option := &repos.FindPageOption{
		TplPrefix: "t_work_external_user/externals-list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]WorkExternalUser{}
		},
	}
	total, list, err := s.repo.FindPage(c, pageVo, params, option)
	if errorIgnore.Ignore(err) != nil {
		return 0, nil, err
	}

	resList := list.([]WorkExternalUser)
	return total, resList, nil
}

// LossCustomersList 流失客户列表
func (s Service) LossCustomersList(c ctx.Context, params LossCustomersListRequest) (int64, *[]LossCustomersListItem, error) {
	cropReq := enterprise.Corp{
		EnterpriseID: params.EnterpriseID,
	}
	crop, err := s.enterpriseRepo.FindEnterpriseDetail(c, cropReq)
	if err != nil {
		return 0, nil, err
	}
	params.CorpId = crop.CorpID

	params.ChangeType = "del_external_contact"
	if params.Type == 2 {
		params.ChangeType = "del_follow_user"
	}

	offset := (params.Page - 1) * params.Size
	pageVo := repos.PageVO{Offset: offset, Limit: params.Size}

	log.Info(c).Msgf("LossCustomersList-params %+v", params)
	option := &repos.FindPageOption{
		TplPrefix: "t_customer/loss-customers-list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]LossCustomersListItem{}
		},
	}

	total, list, err := s.repo.FindPage(c, pageVo, params, option)
	if errorIgnore.Ignore(err) != nil {
		return 0, nil, err
	}

	resList := list.(*[]LossCustomersListItem)
	if len(*resList) > 0 {
		for key, val := range *resList {
			memData, err := s.repo.FindMemberByUserid(c, params.EnterpriseID, val.WorkUserID)
			if err != nil {
				log.Error(c).Msgf("LossCustomersList-FindMemberByUserid err=%+v", err)
			} else {
				(*resList)[key].MemberID = memData.ID
				(*resList)[key].MemberName = memData.Name
				(*resList)[key].MemberAvatar = memData.Avatar
			}
		}
	}

	var memberIDs []int64
	for _, res := range *resList {
		memberIDs = append(memberIDs, res.MemberID)
	}
	if len(memberIDs) > 0 {
		departmentListById, _ := s.repo.MemberMultiDepartments(c, memberIDs)

		for k, relationDetail := range *resList {
			memberID := relationDetail.MemberID
			(*resList)[k].MemberDepartsList = departmentListById[memberID]
		}
	}

	return total, resList, nil
}

// BatchFindWorkGroupChatMembers 批量查询客户群成员
func (s Service) BatchFindWorkGroupChatMembers(c ctx.Context, req BatchWorkGroupMemberRequest) (int64, []WorkGroupChatMember, error) {
	return s.repo.BatchFindWorGroupChatMember(c, req)
}

// BatchDeleteWorkGroupChatMembers 批量删除成
func (s Service) BatchDeleteWorkGroupChatMembers(c ctx.Context, ids []int64) error {
	return s.repo.BatchDeleteWorkGroupChatMembers(c, ids)
}

// BatchUpdateWorkGroupChatMembers 批量更新
func (s Service) BatchUpdateWorkGroupChatMembers(c ctx.Context, ids []int64, update WorkGroupChatMember) (int64, error) {
	where := make(map[string]interface{})
	where["ids"] = ids
	return s.repo.BatchUpdateWorkGroupChatMembers(c, where, update)
}

// AddCustomerCallbackRecord 添加企业微信外部联系人回调记录
func (s Service) AddCustomerCallbackRecord(c ctx.Context, req *work.MemberExternalCallback) (int64, error) {
	if req.ID == 0 {
		req.ID = uid.ID(c)
	}
	var channelID = int64(0)
	log.Info(c).Msgf("AddCustomerCallbackRecordReq:%+v,", req)
	_, err := s.repo.AddCustomerCallbackRecord(c, req)
	// 新增完成后更新
	if err == nil {
		req.AddRelationTime = 0
		// 需要添加关系建立时间
		if req.ChangeType == "del_follow_user" || req.ChangeType == "del_external_contact" {
			timeChannId := s.repo.LossCustomerAddRelationTime(c, req)
			mecReq := work.MemberExternalCallback{
				ID:              req.ID,
				ChannelID:       timeChannId.ChannelID,
				AddRelationTime: timeChannId.AddRelationTime,
			}
			log.Info(c).Msgf("LossCustomerAddRelationTimeInfo %+v", timeChannId)
			channelID = timeChannId.ChannelID
			_, err = s.repo.UpdateCustomerCallbackRecord(c, mecReq)
			if err != nil {
				log.Error(c).Msgf("LossCustomerAddRelationTimeErr %+v,%+v", err, mecReq)
			}
		} else if req.ChangeType == "add_external_contact" {
			// 当建立关系时，判断是否有历史删除记录，如果有更新删除记录的关键建立时间，目的时为了让统计删除不重复统计
			err := gormx.GetInstance().
				Model(micro_scrm.RelationCallbackRecordsModel{}).
				Where("user_id", req.UserID).
				Where("external_user_id", req.ExternalUserID).
				Where("change_type", "del_external_contact").
				Update("add_relation_time", time.Now().Unix()).Error
			if err != nil {
				log.Error(c).Msgf("DelFollowUserByAddAgain %+v,%+v", err, req)
			}
		}
	}
	return channelID, err
}

// UpdateCustomerCallbackRecord 更新企业微信外部联系人回调记录
func (s Service) UpdateCustomerCallbackRecord(c ctx.Context, req work.MemberExternalCallback) (int64, error) {
	if req.ID == 0 {
		req.ID = uid.ID(c)
	}
	return s.repo.UpdateCustomerCallbackRecord(c, req)
}

// QueryCustomerCallbackRecord 查询企业微信外部联系人回调记录
func (s Service) QueryCustomerCallbackRecord(c ctx.Context, page int, size int, req work.MemberExternalCallback) ([]*work.MemberExternalCallback, error) {
	return s.repo.QueryCustomerCallbackRecord(c, page, size, req)
}

// CallbackRecordDetail 查询详情
func (s Service) CallbackRecordDetail(c ctx.Context, req *work.MemberExternalCallback) (*work.MemberExternalCallback, error) {
	return s.repo.CallbackRecordDetail(c, req)
}

// DelCallbackRecord 查询删除企业微信外部联系人回调记录
func (s Service) DelCallbackRecord(c ctx.Context, page int, size int, corpID string) ([]*work.MemberExternalCallback, error) {
	return s.repo.DelCallbackRecord(c, page, size, corpID)
}

// LastRecordDetail 查询详情
func (s Service) LastRecordDetail(c ctx.Context, req *work.MemberExternalCallback) (*work.MemberExternalCallback, error) {
	return s.repo.LastRecordDetail(c, req)
}

// BatchUpdateUserCategorySplit 批量更新客户分类（标签或分组）
func (s Service) BatchUpdateUserCategorySplit(c ctx.Context, receive *ReceiveCDPUserCategory) (map[string]*EnterpriseUser, error) {
	var users = make(map[string]*EnterpriseUser)
	var addUsers = make(map[string]*EnterpriseUser)
	var updateUsers = make(map[string]*EnterpriseUser)
	var err error

	for _, unionID := range receive.Add {
		user, ok := users[unionID]
		if !ok {
			if user, err = s.repo.QueryEnterpriseUser(c, receive.EnterpriseID, &SourceUserBase{UnionID: unionID, CustomerType: CustomerTypeUnion}); err != nil {
				log.Error(c).Msgf("BatchUpdateUserCategorySplit 查询用户: %v", err)
				continue
			}
			if user == nil {
				continue
			}
			users[unionID] = user
		}

		if user.AddTagOrGroup(receive.Type, receive.ID) {
			addUsers[unionID] = user
			updateUsers[unionID] = user
		}
	}
	for _, unionID := range receive.Del {
		user, ok := users[unionID]
		if !ok {
			if user, err = s.repo.QueryEnterpriseUser(c, receive.EnterpriseID, &SourceUserBase{UnionID: unionID, CustomerType: CustomerTypeUnion}); err != nil {
				log.Error(c).Msgf("BatchUpdateUserCategorySplit 查询用户: %v", err)
				continue
			}

			if user == nil {
				continue
			}

			users[unionID] = user
		}

		if user.DelTagOrGroup(receive.Type, receive.ID) {
			updateUsers[unionID] = user
		}
	}

	for _, user := range updateUsers {
		if err := s.repo.UpdateEnterpriseUserAttr(c, receive.EnterpriseID, user.ID, &EnterpriseUserAttr{Tags: user.Tags, Groups: user.Groups}); err != nil {
			log.Err(c, err).Msgf("更新用户标签分组失败: %d", user.ID)
		}
	}

	return addUsers, nil
}

// BatchExternalRelations 外部联系人关系列表
func (s Service) BatchExternalRelations(c ctx.Context, req BatchWorkUserMemberRelationRequest) ([]BatchWorkUserMemberRelationResponse, error) {
	return s.repo.BatchExternalRelations(c, req)
}

// ExternalListRemark 联系人列表备注
func (s Service) ExternalListRemark(c ctx.Context, enterpriseID int64, memberID int64, externalIDList []int64) ([]ExternalRemark, error) {
	return s.repo.ExternalListRemark(c, enterpriseID, memberID, externalIDList)
}

// UpdateWorkGroupChatMemberInfo 更新群聊成员信息
func (s Service) UpdateWorkGroupChatMemberInfo(c ctx.Context, entity *WorkGroupChatMember) error {
	_, err := s.repo.UpdateWorkChatMemberInfo(c, entity)
	if err != nil {
		log.Info(c).Msgf("更新群成员信息失败：%s, member:%v", err.Error(), entity)
		return err
	}
	return nil
}

// HasNormalRelation 判断外部联系人是否存在正常的关系(渠道有值)
func (s Service) HasNormalRelation(ctx ctx.Context, workExternalUserID int64, enterpriseId, channelId int64) (bool, error) {
	// 判断state能否关联到存在的渠道
	return s.repo.HasNormalRelation(ctx, workExternalUserID, enterpriseId, channelId)
}

// UpdateEnterpriseUserAttr 更新属性
func (s Service) UpdateEnterpriseUserAttr(cc ctx.Context, enterpriseID int64, id int64, attr *EnterpriseUserAttr) error {
	return s.repo.UpdateEnterpriseUserAttr(cc, enterpriseID, id, attr)
}

func (s Service) SelectEUidBYDepartment(cc ctx.Context, enterpriseID int64, departmentIDs []int64, page repos.PageVO) (int64, *[]*EUidList, error) {
	return s.repo.SelectEUidBYDepartment(cc, enterpriseID, departmentIDs, page)
}

func (s Service) SyncEterpriseUserToES(c ctx.Context, param *SyncEterpriseUserToEsRequest) {
	pageSize := 50
	euserList := make(chan []EnterpriseUserEntity)
	esIndex := fmt.Sprintf("micro_scrm_enterprise_%d.enterprise_user", param.EnterpriseId)
	prefix := config.CfgData.Env + "."
	esIndex = prefix + esIndex
	if !s.esRepo.ExistsIndex(ctx.New(), esIndex) {
		CreateIndex(s.esRepo, esIndex)
	}
	dtUpdate := time.Now().Format("2006-01-02 15:04:05")
	// 数据日志
	go func() {
		s.repo.AllEterpriseUserForEsList(c, int64(param.EnterpriseId), param.BeginTime, param.EndTime, pageSize, euserList)
	}()
	for {
		list, ok := <-euserList
		if !ok {
			return
		}
		var docList = struct {
			sync.RWMutex
			m map[string]map[string]string
		}{m: make(map[string]map[string]string)}

		eueMap := make(map[string]interface{})
		for _, euser := range list {
			euejson, err := json.Marshal(euser)
			if err != nil {
				continue
			}
			err = json.Unmarshal(euejson, &eueMap)
			if err != nil {
				continue
			}

			t := reflect.TypeOf(euser)
			v := reflect.ValueOf(euser)

			for k := 0; k < t.NumField(); k++ {
				rowkey := t.Field(k).Name
				rowval := v.Field(k).Interface()

				switch t.Field(k).Name {
				case "Remark":
					if _, ok := rowval.([]byte); !ok {
						eueMap["remark"] = ""
					}
					eueMap["remark"] = string(rowval.([]byte))
					if string(rowval.([]byte)) == "[]" {
						eueMap["remark"] = ""
					}
				case "Birthday":
					if _, err := time.Parse("2006-01-02", rowval.(string)); err != nil {
						delete(eueMap, "birthday")
					}
				case "WxUnion":
					if _, ok := rowval.([]byte); !ok {
						rowval = []byte("{}")
					}

					wxUnionList := []WxUnion{}
					rowkey = "wx_union"
					eueMap[rowkey] = wxUnionList
					item := make(map[string]string)
					err = json.Unmarshal(rowval.([]byte), &item)
					if err == nil {
						for k, v := range item {
							wx := WxUnion{
								Pid:     k,
								UnionId: v,
							}
							wxUnionList = append(wxUnionList, wx)
						}
						eueMap[rowkey] = wxUnionList
					}
				case "WxShop", "WxOpen", "WxExternal":
					if _, ok := rowval.([]byte); !ok {
						rowval = []byte("{}")
					}
					if rowkey == "WxShop" {
						rowkey = "wx_shop"
					}
					if rowkey == "WxOpen" {
						rowkey = "wx_open"
					}
					if rowkey == "WxExternal" {
						rowkey = "wx_external"
					}
					wtcs := []WxThirdCommon{}
					eueMap[rowkey] = wtcs
					item := make(map[string]ThirdUserID)
					err = json.Unmarshal(rowval.([]byte), &item)
					if err == nil {
						for k, v := range item {
							wx := WxThirdCommon{
								AppID:     k,
								ThirdUser: v,
							}
							wtcs = append(wtcs, wx)
						}

						eueMap[rowkey] = wtcs
					}

				case "Groups", "Tags":
					if _, ok := rowval.([]byte); !ok {
						rowval = []byte("[]")
					}
					if rowkey == "Groups" {
						rowkey = "groups"
					}
					if rowkey == "Tags" {
						rowkey = "tags"
					}
					item := []int64{}
					eueMap[rowkey] = item
					err = json.Unmarshal(rowval.([]byte), &item)
					if err == nil && len(item) > 0 {
						eueMap[rowkey] = item
					}
				}
			}
			eueMapj, _ := json.Marshal(eueMap)

			euestr := string(eueMapj)
			id := fmt.Sprint(euser.ID)

			bindTo := fmt.Sprint(eueMap["bind_to"])
			if bindTo != "0" {
				err = s.esRepo.DeleteDoc(ctx.New(), esIndex, id)
				log.Info(ctx.New()).Msgf("SyncEterpriseUserToES DeleteDoc -> esIndex : %s, id:%s bindTo:%v err:%v", esIndex, id, bindTo, err)
				tname := strings.TrimLeft(esIndex, prefix)
				s.repo.UpdateEsResult(c, tname, dtUpdate, []string{id})
			} else {
				docList.Lock()
				if _, ok := docList.m[esIndex]; !ok {
					docList.m[esIndex] = make(map[string]string)
				}
				docList.m[esIndex][id] = euestr
				docList.Unlock()
			}
		}

		var body strings.Builder
		docList.RLock()
		if len(docList.m) > 0 {
			for index, dmap := range docList.m {
				if len(dmap) > 0 {
					for id, doc := range dmap {
						body.WriteString(`{"index" : { "_index" : "` + index + `", "_type" : "_doc", "_id" : "` + id + `" }}`)
						body.WriteString("\n" + doc + "\n")

						delete(docList.m[index], id)
					}
					mapres, err := s.esRepo.AddDocsByBulk(ctx.New(), body.String())
					if err != nil {
						log.Error(ctx.New()).Msg("cs.EsRepo.AddDocsByBulk error  " + err.Error())
					}
					if len(mapres) > 0 {
						for k, v := range mapres {
							if len(v) > 0 {
								idList := []string{}
								for m := range v {
									delete(docMap.m[k], m)
									idList = append(idList, m)
								}
								tname := strings.TrimLeft(k, prefix)
								s.repo.UpdateEsResult(c, tname, dtUpdate, idList)
							}
						}
					}
				}
			}
		}
		docList.RUnlock()
	}
}

// 全量初始员工外部联系人的添加时间
func (s Service) LossCustomerAddRelationTimeJob(c ctx.Context, enterpriseId int64, corpID string) {
	log.Info(c).Msgf("LossCustomerAddRelationTimeJobBegin %d,%s", enterpriseId, corpID)
	err := s.repo.LossCustomerAddRelationTimeJob(c, enterpriseId, corpID)
	if err != nil {
		log.Error(c).Msgf("LossCustomerAddRelationTimeJobEnd err %+v", err)
	}
}
