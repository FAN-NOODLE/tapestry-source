package workexternal

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/pkg/sdk/work"
	"micro-scrm/internal/pkg/uid"
)

// 员工删除客户渠道值变化状态
const (
	AddExternalChange = 1 + iota
	AddExternalNoChange
)

type ExternalTraceService struct {
	repo Repository
}

// ExternalExists 判断外部联系人是否处于加入状态
func (s *ExternalTraceService) ExternalExists(c ctx.Context, userId string, enterpriseId int64) (bool, error) {
	external, err := s.repo.ExistsJoinExternal(c, userId, enterpriseId)
	if err != nil || external == nil {
		return false, err
	}
	return true, err
}

// JoinExternal 加入外部联系人
func (s *ExternalTraceService) JoinExternal(c ctx.Context, data *work.MemberExternalCallback, enterpriseId int64) error {
	// 查询上次离开状态记录
	lastLeaveTime := 0
	trace, err := s.repo.LastExternalTrace(c, data.ExternalUserID, enterpriseId)
	if err != nil {
		return err
	}
	if trace != nil {
		lastLeaveTime = trace.LeaveTime
	}

	// 插入
	req := &ExternalTrace{
		ID:             uid.ID(c),
		ExternalUserID: data.ExternalUserID,
		EnterpriseID:   enterpriseId,
		JoinTime:       int(data.CreateTime),
		LastLeaveTime:  lastLeaveTime,
	}
	err = s.repo.JoinExternal(c, req, trace)
	if err != nil {
		return err
	}

	return nil
}

// LeaveExternal 删除该外部联系人
func (s *ExternalTraceService) LeaveExternal(c ctx.Context, data *work.MemberExternalCallback, enterpriseId int64) error {
	// 删除外部联系人
	externalTrace := &ExternalTrace{
		ExternalUserID: data.ExternalUserID,
		EnterpriseID:   enterpriseId,
		LeaveTime:      int(data.CreateTime),
	}

	// 联系被删除只保留最早的哪个
	err := s.repo.LeaveExternalByUserId(c, externalTrace)
	if err != nil {
		return err
	}

	return nil
}

// ExternalExistsInChannel 判断渠道外部联系人是否处于加入状态
func (s *ExternalTraceService) ExternalExistsInChannel(c ctx.Context, userId string, enterpriseId, channelId int64) (bool, error) {
	external, err := s.repo.ExistsJoinExternalInChannel(c, userId, enterpriseId, channelId)
	if err != nil || external == nil {
		return false, err
	}
	return true, err
}

func (s *ExternalTraceService) JoinExternalInChannel(c ctx.Context, data *work.MemberExternalCallback, enterpriseId, channelId int64) error {
	// 查询该渠道上次加入状态记录
	lastLeaveTime := 0
	trace, err := s.repo.LastExternalTraceInChannel(c, data.ExternalUserID, enterpriseId, channelId)
	if err != nil {
		return err
	}

	if trace != nil {
		lastLeaveTime = trace.LeaveTime
	}

	// 插入
	req := &ExternalChannelTrace{
		ID:             uid.ID(c),
		ExternalUserID: data.ExternalUserID,
		EnterpriseID:   enterpriseId,
		JoinTime:       int(data.CreateTime),
		AddChannel:     channelId,
		LeaveTime:      0,
		LastLeaveTime:  lastLeaveTime,
		NextJoinTime:   0,
	}
	err = s.repo.JoinExternalInChannel(c, req, trace)
	if err != nil {
		return err
	}

	return nil
}

// LeaveExternalInChannel 删除该渠道外部联系人
func (s *ExternalTraceService) LeaveExternalInChannel(c ctx.Context, data *work.MemberExternalCallback, enterpriseId, channelId int64) error {
	trace := &ExternalChannelTrace{
		ExternalUserID: data.ExternalUserID,
		EnterpriseID:   enterpriseId,
		LeaveTime:      int(data.CreateTime),
		AddChannel:     channelId,
	}

	err := s.repo.LeaveChannelExternalByUserId(c, trace)
	if err != nil {
		return err
	}

	return nil
}

func (s *ExternalTraceService) AddDelExternalContact(c ctx.Context, data *work.MemberExternalCallback, enterpriseId int64, channelId int64) error {
	trace := &DelRelationChannelTrace{
		ID:             uid.ID(c),
		ExternalUserID: data.ExternalUserID,
		UserID:         data.UserID,
		EnterpriseID:   enterpriseId,
		DeleteTime:     int(data.CreateTime),
		DeleteChannel:  channelId,
	}

	return s.repo.AddDelExternalContact(c, trace)
}

// ExistDelRelation 是否存在需要恢复的关系
func (s *ExternalTraceService) ExistDelRelation(c ctx.Context, data *work.MemberExternalCallback, enterpriseId int64) (*DelRelationChannelTrace, error) {
	condition := &DelRelationChannelTrace{
		ExternalUserID: data.ExternalUserID,
		UserID:         data.UserID,
		EnterpriseID:   enterpriseId,
	}
	return s.repo.FindDelRelation(c, condition)
}

// RecoveryDelRelation 恢复外部联系人删除员工关系
func (s *ExternalTraceService) RecoveryDelRelation(c ctx.Context, data *work.MemberExternalCallback, existDelRelation *DelRelationChannelTrace, channelId int64) error {
	var change uint8 = AddExternalChange
	if existDelRelation.AddChannel != channelId {
		change = AddExternalNoChange
	}

	update := &DelRelationChannelTrace{
		AddTime:      int(data.CreateTime),
		AddChannel:   channelId,
		Change:       change,
		EnterpriseID: existDelRelation.EnterpriseID,
	}

	return s.repo.RecoveryDelRelation(c, existDelRelation.ID, update)
}

func (s *ExternalTraceService) AddRelationRecord(c ctx.Context, data *work.MemberExternalCallback, addChannel int64) (*WorkRelationAddRecord, error) {
	record := &WorkRelationAddRecord{
		CorpID:         data.ToUserName,
		AddTime:        data.CreateTime,
		UserId:         data.UserID,
		ExternalUserId: data.ExternalUserID,
		AddChannel:     addChannel,
	}

	err := s.repo.AddRelationRecord(c, record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func (s *ExternalTraceService) FindLastAddRelationRecord(c ctx.Context, data *work.MemberExternalCallback) (*WorkRelationAddRecord, error) {
	record := &WorkRelationAddRecord{
		CorpID:         data.ToUserName,
		AddTime:        data.CreateTime,
		UserId:         data.UserID,
		ExternalUserId: data.ExternalUserID,
	}

	record, err := s.repo.FindLastAddRelationRecord(c, record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func (s *ExternalTraceService) FinishAddRelation(c ctx.Context, record *WorkRelationAddRecord) error {
	return s.repo.FinishAddRelation(c, record)
}

// HistoryUcRelation 查询历史数据
// status = 0 查询未处理的; status = 1 查询处理成功的
func (s *ExternalTraceService) HistoryUcRelation(c ctx.Context, page, size int64, status int8) ([]*HistoryUcRelationLog, error) {
	offset := (page - 1) * size
	pageVo := repos.PageVO{Offset: int(offset), Limit: int(size)}
	return s.repo.HistoryUcRelation(c, pageVo, status)
}

// UpdateHistoryStatus 更新状态
func (s *ExternalTraceService) UpdateHistoryStatus(c ctx.Context, id int64, status int8) error {
	return s.repo.UpdateHistoryStatus(c, id, status)
}
