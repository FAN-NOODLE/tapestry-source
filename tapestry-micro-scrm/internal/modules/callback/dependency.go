package callback

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/modules/customer"
	"micro-scrm/internal/modules/fission"
	"micro-scrm/internal/modules/members"
	workexternal "micro-scrm/internal/modules/work-external"
	"micro-scrm/internal/pkg/sdk/work"
)

type ExternalService interface {
	// ExternalExists 判断外部联系人是否处于加入状态
	ExternalExists(c ctx.Context, userID string, enterpriseId int64) (bool, error)
	// 加入外部联系人
	JoinExternal(c ctx.Context, data *work.MemberExternalCallback, enterpriseId int64) error
	// 删除该外部联系人
	LeaveExternal(c ctx.Context, data *work.MemberExternalCallback, enterpriseId int64) error

	// 判断渠道外部联系人是否处于加入状态
	ExternalExistsInChannel(c ctx.Context, userID string, enterpriseId, int64 int64) (bool, error)
	// 加入渠道外部联系人
	JoinExternalInChannel(c ctx.Context, data *work.MemberExternalCallback, enterpriseId, int64 int64) error
	// 删除该渠道外部联系人
	LeaveExternalInChannel(c ctx.Context, data *work.MemberExternalCallback, enterpriseId, channelId int64) error

	// 添加外部联系人删除员工
	AddDelExternalContact(c ctx.Context, data *work.MemberExternalCallback, enterpriseId int64, channelId int64) error
	// 是否存在需要恢复的关系
	ExistDelRelation(c ctx.Context, data *work.MemberExternalCallback, enterpriseId int64) (*workexternal.DelRelationChannelTrace, error)
	// 恢复外部联系人删除员工关系
	RecoveryDelRelation(c ctx.Context, data *work.MemberExternalCallback, existDelRelation *workexternal.DelRelationChannelTrace, channelId int64) error

	// 添加
	AddRelationRecord(c ctx.Context, data *work.MemberExternalCallback, channelID int64) (*workexternal.WorkRelationAddRecord, error)
	// 查询上次添加记录
	FindLastAddRelationRecord(c ctx.Context, data *work.MemberExternalCallback) (*workexternal.WorkRelationAddRecord, error)
	// 完成状态
	FinishAddRelation(c ctx.Context, data *workexternal.WorkRelationAddRecord) error
	// 查询历史数据
	HistoryUcRelation(c ctx.Context, page, size int64, status int8) ([]*workexternal.HistoryUcRelationLog, error)
	// 更新历史记录状态
	UpdateHistoryStatus(c ctx.Context, id int64, status int8) error
}

type CustomerService interface {
	// 判断在开始时间之前是否存在关系
	ExistsRelationByStartTime(c ctx.Context, req *customer.ExistsRelationByTimeReq) (bool, error)
	// 获取外部联系人详情
	FindExternalByUserID(c ctx.Context, externalUserID string, enterpriseID int64) (*customer.WorkExternalUser, error)
	// 判断外部联系人是否在群列表中
	WorkExternalInGroup(ctx ctx.Context, externalUserID string, groupIDs []int64, enterpriseID int64) (int64, error)
	// 查询群详情
	FindWorkGroupChatDetail(ctx ctx.Context, req customer.WorkGroupChat) (*customer.WorkGroupChat, error)
}

type FissionService interface {
	// 通过联系我查询裂变人表
	GetActivityUserByContactID(c ctx.Context, ContactID int64) (*fission.ActivityUserEntity, error)
	// 查询活动通过ID
	ActivityDetail(c ctx.Context, enterpriseID int64, id int64) (*fission.ActivityEntity, error)
	// 获取活动关联的成员
	ActivityWithMembers(c ctx.Context, activityID, enterpriseID int64) ([]int64, error)
	// 通过员工查询是否存在助力记录
	ExistRelationWithMembers(c ctx.Context, userID int64, members []int64, enterpriseID int64) (bool, error)
	// 更新 fission_assistant_records 数据和统计数据
	RecoveryRecordAndStatistics(c ctx.Context, assistantRecord *fission.AssistantRecordsEntity, fissionType int, updateData *fission.UpdateAssistantRecord) error
	// 创建 fission_assistant_records 数据和统计数据
	CreateAssistantRecordAndStatistics(c ctx.Context, record *fission.AssistantRecordsEntity) error
	// 是否存在删除的助力记录
	RecoveryFissionRecord(c ctx.Context, customerID, memberID int64) (*fission.AssistantRecordsEntity, error)
	// 查询是否已经助力
	GetAddRecordByRelation(c ctx.Context, customerID, memberID, enterpriseID int64) (*fission.AssistantRecordsEntity, error)

	// 查询这个活动的社群
	GetGroupsIDByActivity(c ctx.Context, activityID int64, enterpriseID int64, withDel bool) ([]int64, error)
	// 通过群聊ID获取活动
	GetActivityByChatID(c ctx.Context, req *fission.ActivityByChatIDReq) ([]int64, error)
	// 查找群和记录关联数据
	FindGroupRecordsByChatID(c ctx.Context, activityID []int64, chatID, enterpriseID int64, joinOrLeave bool) ([]*fission.AssistantRecordsWithTime, error)
}

type MemberService interface {
	GetMemberIDByUserid(c ctx.Context, userID string, enterpriseID int64) (*members.Member, error)
}

type EnterpriseService interface {
	GetEnterpriseID(c ctx.Context, corpID string) (int64, error)
}
