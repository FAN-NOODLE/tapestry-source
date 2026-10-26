package fission

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
)

//go:generate mockgen -source=repository.go -destination=mock_repository.go -package=fission
type Repository interface {
	repos.IRepository
	// ActivitySet 裂变活动设置
	ActivitySet(c ctx.Context, enterpriseID int64, action int, activity *ActivityEntity, members *MembersSet, chatGroups *ChatGroupSet, welcomeContent string, prize *PrizeSet, mutual *MutualSet) (int64, error)
	// ActivityList 裂变活动列表
	ActivityList(c ctx.Context, request *ActivityListRequest) (int64, []*ActivityEntity, error)
	// MutualActivityList 选择互斥活动
	MutualActivityList(c ctx.Context, request *MutualActivityListRequest) (int64, []*MutualActivityEntity, error)
	// ActivityDelete 活动删除
	ActivityDelete(c ctx.Context, activityID int64) error
	// UpdateStatus 更新活动状态
	UpdateStatus(c ctx.Context, activityID int64, status int) error
	// ChatGroupList 社群活动列表
	ChatGroupList(c ctx.Context, request *ChatGroupListRequest) (int64, []*WorkChatGroupEntity, error)
	// ActivityDetail 活动详情
	ActivityDetail(c ctx.Context, activityID int64, enterpriseID int64) (*ActivityEntity, error)
	// 通过欢迎语ID查询活动
	ByContactIdActivityDetail(c ctx.Context, contactID int64, enterpriseID int64) (*ActivityEntity, error)
	FissionMemberIDs(c ctx.Context, activityID int64) ([]*Member, error)
	// MemberList 查询员工名称
	MemberList(c ctx.Context, memberIDs []int64, enterpriseID int64) ([]*Member, error)
	// ActivityList 裂变活动列表
	QwActivityList(c ctx.Context, request *QwActivityListRequest) (int64, *[]QwActivityList, error)
	// MutualActivityIDs 查询互斥活动ID
	MutualActivityIDs(c ctx.Context, activityID int64, enterpriseID int64) ([]*MutualRelationEntity, error)
	// StatisticsInfo 活动详情
	StatisticsInfo(c ctx.Context, request *QwActivityListRequest) (*Statistics, error)
	// FissionGroupIDs 查询活动裂变群ID
	FissionGroupIDs(c ctx.Context, activityID int64, enterpriseID int64, withDel bool) ([]*ChatGroupEntity, error)
	// WelComeInfo 查询活动欢迎语
	WelComeInfo(c ctx.Context, enterpriseID int64, drainageWelcomesID int64) (string, error)
	//AssistantUserList 发起人的助力人列表
	AssistantUserList(c ctx.Context, ids []int64, enterpriseID int64) ([]*AssistantRecordInfo, error)
	// PrizeList 奖品列表
	PrizeList(c ctx.Context, enterpriseID int64, activityID int64, page int64, size int64) (int64, *[]*PrizeEntity, error)
	// PrizeInfo 奖品详情
	PrizeInfo(c ctx.Context, ID int64, withDeleted bool) (*PrizeEntity, error)
	//ExchangeRecord发起人兑换记录
	ExchangeRecord(c ctx.Context, req ExchangeDetailRequest) (int64, *[]ExchangeRecordInfo, error)
	//ActivityDetailList 活动详情 列表
	ActivityDetailList(c ctx.Context, req ActivityDetailRequest) (int64, *[]ActivityUserEntity, error)
	//StoreExchangeWarnConfig 配置兑换奖品提醒
	StoreExchangeWarnConfig(c ctx.Context, request ExchangeWarnRequest) error
	//InviteDetail 发起人邀请明细
	InviteDetail(c ctx.Context, req InviteDetailRequest) (int64, []*AssistantRecordInfo, error)
	//GetExchangeWarnConfig
	GetExchangeWarnConfig(c ctx.Context, enterpriseId int64, key string) (string, error)
	//FissionMemberList 查找所有活动下的裂变成员
	FissionMemberList(c ctx.Context, enterpriseId int64, activityID int64) ([]*Member, error)
	//ActivityListSimple 活动列表 全部有效活动
	ActivityListSimple(c ctx.Context, enterpriseId int64) ([]*ActivityEntity, error)
	GetActivityUser(c ctx.Context, enterpriseID int64, activityID int64, userID int64) (*ActivityUserEntity, error)
	// GetMpUserInfo 获取用户h5授权信息
	GetMpUserInfo(c ctx.Context, ID int64) (*MpUser, error)
	//奖品兑换总数
	ExchangeTotal(c ctx.Context, enterpriseId int64) (int64, error)
	// =======================================fission_assistant_records========================
	// RecoveryFissionRecord 查询是否存在删除的记录
	RecoveryFissionRecord(c ctx.Context, customerID, memberID int64) (*AssistantRecordsEntity, error)
	// 查询是否存在关系
	GetRecordByRelation(c ctx.Context, req *ExistRelationReq) (*AssistantRecordsEntity, error)

	// 编辑并加入好友并进行统计逻辑
	RecoveryRecordAndStatistics(c ctx.Context, assistantRecord *AssistantRecordsEntity, fissionType int, updateData *UpdateAssistantRecord, fissionUserStatistics *UpdateStatisticsData) error
	// 查询裂变记录并进行统计逻辑
	CreateAssistantRecordAndStatistics(c ctx.Context, record *AssistantRecordsEntity, fissionUserStatistics *UpdateStatisticsData) error

	// =======================================fission_user========================
	GetActivityUserByContactID(c ctx.Context, ContactID int64) (*ActivityUserEntity, error)

	ByActivityContactQrCode(c ctx.Context, activityID int64, enterpriseID int64) (*ChatGroupEntity, error)
	// GetActivityUserList 获取用户参与指定活动列表信息
	GetActivityUserList(c ctx.Context, userID int64, activityIDList []int64, enterpriseID int64) ([]*ActivityUserEntity, error)
	AddFissionUser(c ctx.Context, entity *ActivityUserEntity) (int64, error)
	GetExchangeList(c ctx.Context, enterpriseID int64, activityID int64, userID int64, page int64, size int64) (int64, *[]*ExchangePrizeRecordsEntity, error)
	// PrizeListByIDs 奖品列表
	PrizeListByIDs(c ctx.Context, IDList []int64) (*[]*PrizeEntity, error)
	// GetExchangeInfo 兑换详情
	GetExchangeInfo(c ctx.Context, ID int64) (*ExchangePrizeRecordsEntity, error)
	// GetAssistantList 助力列表
	GetAssistantList(c ctx.Context, activityID int64, userID int64, page int64, size int64) (int64, *[]*AssistantRecordsEntity, error)
	// UserExchange 用户兑换奖品
	UserExchange(c ctx.Context, userActivity *ActivityUserEntity, prize *PrizeEntity, record *ExchangePrizeRecordsEntity) (bool, error)
	// GetUserExchangeNum 获取用户兑换奖品数
	GetUserExchangeNum(c ctx.Context, userID int64, prizeID int64) (int64, error)
	// QueryUserLatestExchangeRecord 获取用户最近兑换的信息
	QueryUserLatestExchangeRecord(c ctx.Context, userID int64) (*ExchangePrizeRecordsEntity, error)
	// 获取Activity通过群ID
	GetActivityByChatID(c ctx.Context, req *ActivityByChatIDReq) ([]int64, error)
	// 通过聊天ID查询需要入群的记录
	FindGroupRecordsByChatID(c ctx.Context, activityID []int64, chatID, enterpriseID int64, joinOrLeave bool) ([]*AssistantRecordsWithTime, error)
	// 更新用户活动信息
	UpdateFissionUser(c ctx.Context, entity *ActivityUserEntity) (int64, error)
	// GetActivityDetailWithDeleted 获取活动详情，查询删除的数据
	GetActivityDetailWithDeleted(c ctx.Context, activityID int64, EnterpriseID int64) (*ActivityEntity, error)
	// HasWorkExternalRelation 查询用户是否有外部联系人关系数据
	HasWorkExternalRelation(c ctx.Context, enterpriseID int64, workUserID int64, memberIDs []int64) (bool, error)
	// GetAssistantUserList 获取活动最新的助力数据
	GetAssistantUserList(c ctx.Context, enterpriseID, activityID, index, limit int64, needGroup bool) ([]*AssistantRecordsEntity, error)
	// GetTotalAssistantUser 获取活动总助力数
	GetTotalAssistantUser(c ctx.Context, enterpriseID int64, activityID int64, needGroup bool) (int64, error)
	// ActivityUserNumAdd 添加活动发起人数
	ActivityUserNumAdd(c ctx.Context, activityID int64, num int) (int64, error)
	// MutualDeleteByActivityID 删除该活动互斥活动
	MutualDeleteByActivityID(c ctx.Context, activityID int64) error
	// ChatDeleteByActivityID 删除该活动社群
	ChatDeleteByActivityID(c ctx.Context, activityID int64) error
	// PrizeDeleteByActivityID 删除该活动的奖品
	PrizeDeleteByActivityID(c ctx.Context, activityID int64) error
	// WelcomeDeleteByID 删除该活动的欢迎语
	WelcomeDeleteByID(c ctx.Context, ID int64) error
	// ContactDeleteByID 删除该活动的默认二维码
	ContactDeleteByID(c ctx.Context, ID int64) error
}

// EventRepository 用户事件仓库
type EventRepository interface {
	repos.IRepository
}
