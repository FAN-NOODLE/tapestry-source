package socialDraw

import (
	"fmt"
	"time"

	"micro-scrm/internal/modules/members"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/pkg/redis"
)

// Service 社群画像服务
type Service struct {
	repo   Repository
	ckrepo Ckrepository
}

// GetGroupChatList 获取群列表数据
func (s *Service) GetGroupChatList(c ctx.Context, params WorkGroupListRequest) ([]*WorkGroupChatResponse, error) {
	return s.repo.FindWorkGroupChatList(c, params)
}

// WorkGroupChatSetting 社群设置
func (s *Service) WorkGroupChatSetting(c ctx.Context, params WorkGroupChatSettingRequest) error {
	return s.repo.WorkGroupChatSetting(c, params)
}

// GroupUserStatistics 获取群人数统计 新增 留存 退群
func (s *Service) GroupUserUvStatistics(c ctx.Context, params GroupInfoReq) (*WorkGroupChatStatistic, error) {
	return s.ckrepo.GroupUserUvStatistics(c, params)
}

// GroupUserStatistics get group chat info
func (s *Service) GroupMemberStatistics(c ctx.Context, params GroupInfoReq) (*WorkGroupChatStatistic, error) {
	return s.ckrepo.GroupMemberStatistics(c, params)
}

func (s *Service) GroupMemberStatisticsForCrm(c ctx.Context, params GroupInfoReq) (*WorkGroupChatStatistic, error) {
	return s.ckrepo.GroupMemberStatisticsForCrm(c, params)
}

// GetGroupChatTypeList 获取社群类型列表数据
func (s *Service) GetGroupChatTypeList(c ctx.Context, params WorkGroupChatType, markUsed bool) ([]*WorkGroupChatType, error) {
	data, err := s.repo.FindWorkGroupTypeList(c, params)
	if markUsed == false {
		return data, err
	}
	countMap, _ := s.repo.CountWorkGroupByType(c, params)
	for _, v := range data {
		if gNum, ok := countMap[v.ID]; ok {
			v.GroupChatNum = gNum
		}
	}
	return data, err
}

// ManageGroupChatType 社群类型管理
func (s *Service) ManageGroupChatType(c ctx.Context, params WorkGroupChatType) (*WorkGroupChatType, error) {
	return s.repo.ManageWorkGroupType(c, params)
}

// GetGroupChatTypeInfo get group chat type info
func (s *Service) GetGroupChatTypeInfo(c ctx.Context, params WorkGroupChatType) (*WorkGroupChatType, error) {
	return s.repo.GetWorkGroupTypeInfo(c, params)
}

// GetWorkGroupFilters 群列表筛选
func (s *Service) GetWorkGroupFilters(c ctx.Context, params WorkGroupFilterRequest) ([]*WorkGroupFilterResponse, error) {
	return s.repo.WorkGroupFilters(c, params)
}

// GetWorkGroupChatMembers 获取社群成员明细
func (s *Service) GetWorkGroupChatMembers(c ctx.Context, params WorkGroupChatMembersRequest) ([]*WorkGroupChatMember, int, error) {
	return s.repo.GetWorkGroupMembers(c, params)
}

// GetWorkGroupChatMembers 获取社群成员明细
func (s *Service) GetWorkGroupChatInfo(c ctx.Context, params WorkGroupChatMembersRequest) ([]*WorkGroupChatMemberList, int, error) {
	return s.repo.GetWorkGroupMembersInfo(c, params)
}

func (s *Service) NewGroupUserChat(c ctx.Context, param GroupInfoReq) ([]*Chart, error) {
	return s.ckrepo.NewGroupUserChat(c, param)
}

func (s *Service) NewGroupUserChatForCrm(c ctx.Context, param GroupInfoReq) ([]*Chart, error) {
	return s.ckrepo.NewGroupUserChatForCrm(c, param)
}
func (s *Service) ValidGroupUserChat(c ctx.Context, param GroupInfoReq) ([]*Chart, error) {
	return s.ckrepo.ValidGroupUserChat(c, param)
}

func (s *Service) ValidGroupUserChatForCrm(c ctx.Context, param GroupInfoReq) ([]*Chart, error) {
	return s.ckrepo.ValidGroupUserChatForCrm(c, param)
}

// GetWorkGroupChatStatics 获取社群统计明细
func (s *Service) GetWorkGroupChatStatics(c ctx.Context, params WorkGroupChatStaticsRequest) (map[int64]*WorkGroupChatStatistic, error) {
	return s.ckrepo.GetWorkGroupStatics(c, params)
}

// CheckWorkGroupChatTypeUsed 检测群类型是否被使用
func (s *Service) CheckWorkGroupChatTypeUsed(c ctx.Context, enterpriseId, chatTypeId int64) (bool, error) {
	return s.repo.CheckWorkGroupChatTypeUsed(c, enterpriseId, chatTypeId)
}

// GetWorkGroupExternalUsers 获取社群外部联系人信息
func (s *Service) GetWorkGroupExternalUsers(c ctx.Context, enterpriseId int64, userIds []string) (map[string]*WxUserCorpExternal, error) {
	return s.repo.GetWorkGroupExternalUsers(c, enterpriseId, userIds)
}

// GetWorkGroupExternalUsers 获取社群外部联系人信息
func (s *Service) GetWorkGroupMembersList(ctx ctx.Context, param WorkGroupChatMembersRequest) ([]*WorkGroupChatMemberList, int, error) {
	return s.ckrepo.GetWorkGroupMembersList(ctx, param)
}

func (s *Service) GetMarsID(ctx ctx.Context, enterpriseId int64, UserID string) (string, error) {
	return s.repo.GetMarsID(ctx, enterpriseId, UserID)
}

// GetWorkGroupChatListFilters 社群列表过滤条件（运营人员，所属部门）
func (s *Service) GetWorkGroupChatListFilters(ctx ctx.Context, params WorkGroupListRequest) []string {
	return s.repo.GetWorkGroupChatListFilters(ctx, params)
}

func (s *Service) getCacheKey(enterpriseId int64) string {
	return fmt.Sprintf("work_group_chat_last_update_time:%+v", enterpriseId)
}

func (s *Service) SetLastUpdateTime(c ctx.Context, enterpriseId int64, val int64) {
	client := redis.Client()
	client.SetEX(ctx.Wrap(c), s.getCacheKey(enterpriseId), val, 0)
	return
}

func (s *Service) GetLastUpdateTime(c ctx.Context, enterpriseId int64) int64 {
	client := redis.Client()
	timestamp, _ := client.Get(ctx.Wrap(c), s.getCacheKey(enterpriseId)).Int64()
	return timestamp
}

func (s *Service) getWorkGroupSyncKey(enterpriseId int64) string {
	return fmt.Sprintf("work_group_chat_sync:%+v", enterpriseId)
}

func (s *Service) SetOrClearWorkGroupSyncFlag(c ctx.Context, enterpriseId int64, exp int64) {
	client := redis.Client()
	client.Set(ctx.Wrap(c), s.getWorkGroupSyncKey(enterpriseId), time.Now().Unix(), time.Second*time.Duration(exp))
	return
}

func (s *Service) GetWorkGroupSyncFlag(c ctx.Context, enterpriseId int64) int64 {
	client := redis.Client()
	timestamp, _ := client.Get(ctx.Wrap(c), s.getWorkGroupSyncKey(enterpriseId)).Int64()
	return timestamp
}

// GetWorkGroupUserActivityInfo 获取用户活动详情
func (s *Service) GetWorkGroupUserActivityInfo(ctx ctx.Context, enterpriseId int64, UserID string, cid int64) ([]*ActInfo, error) {
	return s.repo.GetWorkGroupUserActivityInfo(ctx, enterpriseId, UserID, cid)
}

// GetExternalUserRelatedMembers 获取外部联系人关联的员工数据
func (s *Service) GetExternalUserRelatedMembers(ctx ctx.Context, enterpriseId int64, externalUserIds []string) (map[string][]*members.Member, error) {
	externalUsers, err := s.repo.GetWorkGroupExternalUsers(ctx, enterpriseId, externalUserIds)
	if err != nil {
		return nil, err
	}
	var corpExternalId []int64
	for _, eu := range externalUsers {
		corpExternalId = append(corpExternalId, eu.ID)
	}
	return nil, nil
}
