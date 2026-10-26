package workexternal

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
)

// Repository 存储方法
type Repository interface {
	repos.IRepository
	SelectReUv(c ctx.Context, param SelectKeep) (int64, error)
	SelectAddUv(c ctx.Context, param SelectKeep) (int64, error)
	SelectLossUv(c ctx.Context, param SelectKeep) (int64, error)

	// 查询外部联系人状态
	ExistsJoinExternal(c ctx.Context, userId string, enterpriseId int64) (*ExternalTrace, error)
	// 查询上一条记录
	LastExternalTrace(c ctx.Context, userId string, enterpriseId int64) (*ExternalTrace, error)
	// 插入外部联系人状态
	JoinExternal(c ctx.Context, req *ExternalTrace, lastTrace *ExternalTrace) error
	// 外部联系人离开系统
	LeaveExternalByUserId(c ctx.Context, trace *ExternalTrace) error

	// 查询正常外部联系人是否存在该渠道
	ExistsJoinExternalInChannel(c ctx.Context, userId string, enterpriseId, channelId int64) (*ExternalChannelTrace, error)
	// 查询上一条记录
	LastExternalTraceInChannel(c ctx.Context, userId string, enterpriseId, channelId int64) (*ExternalChannelTrace, error)
	// 加入该渠道
	JoinExternalInChannel(c ctx.Context, req, existsTrace *ExternalChannelTrace) error
	// 离开该渠道
	LeaveChannelExternalByUserId(c ctx.Context, trace *ExternalChannelTrace) error

	// 添加外部联系人删除员工记录
	AddDelExternalContact(c ctx.Context, trace *DelRelationChannelTrace) error
	// 查询是否存在删除记录
	FindDelRelation(c ctx.Context, condition *DelRelationChannelTrace) (*DelRelationChannelTrace, error)
	// 恢复删除的关系
	RecoveryDelRelation(c ctx.Context, id int64, update *DelRelationChannelTrace) error
	AddRelationRecord(c ctx.Context, record *WorkRelationAddRecord) error
	FindLastAddRelationRecord(c ctx.Context, record *WorkRelationAddRecord) (*WorkRelationAddRecord, error)
	FinishAddRelation(c ctx.Context, record *WorkRelationAddRecord) error
	HistoryUcRelation(c ctx.Context, page repos.PageVO, status int8) ([]*HistoryUcRelationLog, error)
	UpdateHistoryStatus(c ctx.Context, id int64, status int8) error
}
