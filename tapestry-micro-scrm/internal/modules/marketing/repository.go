package marketing

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
)

//go:generate mockgen -source=repository.go -destination=mock_repository.go -package=marketing
// Repository 存储定义
type Repository interface {
	repos.IRepository
	MarketingAdd(c ctx.Context, request *MarketingStrategies) error
	MarketingStrategyConditionAdd(c ctx.Context, request *MarketingStrategyCondition) error
	MarketingStrategyChannelAdd(c ctx.Context, request []ChannelInfo) error
	MarketingStrategyTriggerRelationAdd(c ctx.Context, request []RelationInfo) error
	MarketingDetail(c ctx.Context, request int64) (*MarketingStrategy, error)
	MarketingDetailChannels(c ctx.Context, request MarketingDetailRequest) (*[]ChannelInfo, error)
	MarketingDetailRelations(c ctx.Context, request MarketingDetailRequest) (*[]RelationInfo, error)
	MarketingUpdate(c ctx.Context, request MarketingStrategy) error
	MarketingUpdateStatus(c ctx.Context, request *UpdateStatusRequest) error
	MarketingUpdateTriggerCount2(c ctx.Context, request MarketingStrategy) error
	//逻辑删除
	MarketingDelete(c ctx.Context, request DeleteRequest) error
	//物理删除，新增数据替换
	DeleteMarketingStrategy(c ctx.Context, request DeleteRequest) error
	DeleteMarketingStrategyCondition(c ctx.Context, request DeleteRequest) error
	DeleteMarketingStrategyChannels(c ctx.Context, request DeleteRequest) error
	DeleteMarketingStrategyRelations(c ctx.Context, request DeleteRequest) error

	MarketingStrategyRecordAdd(c ctx.Context, request MarketingStrategy) (int64, error)
	GetSmsNumber(c ctx.Context, request MarketingDataDetailRequest) ([]SmsSendLog, error)
	AddMarketStrategyRecord(c ctx.Context, request *MarketStrategyRecord) error
	BatchAddMarketStrategyRecord(c ctx.Context, request []*MarketStrategyRecord) error
	UpdateMarketStrategyRecord(c ctx.Context, request *MarketStrategyRecord) error
	MarketTrig(cc ctx.Context, data *MarketTrigRequest) (*MarketTrigResponse, error)
	MarketTrigList(cc ctx.Context, data *MarketTrigRequest) (*MarketTrigResponse, error)
	UpdateStrategyTime(c ctx.Context, request *UpdateStatusRequest) error

	// 查询所有进行中的策略列表
	SelectMarketingList(ctx ctx.Context) (*[]MarketingStrategies, error)

	MarketingCron(c ctx.Context, request *RecordCronRequest) ([]*MarketTrigExecRequest, error)
	// 查询客户分组列表
	GetUserGroup(c ctx.Context, request UserGroup) ([]UserGroup, error)
	// 获取客户分组用户总数
	GetUserCategory(c ctx.Context, request UserCategory) int64
	// 获取用户分组中有手机号的用户总数
	GetMobileUserCount(c ctx.Context, request UserCategory) ([]*UserInfo, error)
	// 获取用户分组中有企微的用户总数
	GetQyUserCount(c ctx.Context, request UserCategory) ([]*UserInfo, error)
	// 获取用户分组商城用户总数
	GetShopUserCount(c ctx.Context, request UserCategory) ([]*UserInfo, error)
	// 获取用户分组中公众号用户总数
	GetWeChatOfficialUserCount(c ctx.Context, request UserCategory) ([]*UserInfo, error)
	// 判断用户分组是否被营销策略引用
	IsUserGroupUserInMarket(c ctx.Context, groupID int64) (bool, error)

	GetUserStatisticData(c ctx.Context, enterpriseID int64, userID int64) (*UserStatisticEntity, error)
	AddUserStatisticData(ctx ctx.Context, data *UserStatisticEntity) (int64, error)
	UpdateUserStatisticData(ctx ctx.Context, data *UserStatisticEntity) (int64, error)
	GeSmsCount(c ctx.Context, smsId int64) (int64, error)
	CloseMarketingInnerStatus(c ctx.Context, request *MarketingDataDetailRequest) error

	// 查询策略渠道内容
	MarketTrigChannelDetailByWhere(ctx ctx.Context, data ChannelInfo) (*ChannelInfo, error)
	// 批量更新执行记录
	BatchUpdateMarketStrategyRecord(c ctx.Context, request *BatchUpdateRecordsRequest) error
	// 添加群发结果表
	AddQySendMsgResult(c ctx.Context, request *QyMsgSendResult) error

	SelectMsgIds(ctx ctx.Context) ([]QyMsgSendResult, error)
	UpdateQyMsgSendResult(ctx ctx.Context, data *QyMsgSendResult) (int64, error)
	UpdateMarketStrategyRecordByEcuserid(c ctx.Context, data *QyMsgSendResult) (int64, error)
	MarketingUpdateStatusByTime(c ctx.Context, request *UpdateStatusRequest) error
	UpdateStrategyStatusByIds(c ctx.Context, status int, ids []int64) error
	UpdateStrategyHistoryStatusByIds(c ctx.Context, status int, ids []int64) error

	SimpleStrategyList(cc ctx.Context, request ListRequest) ([]*MarketingStrategy, error)
	ListChannelByStrategyIds(cc ctx.Context, req StrategyChannelListRequest) (*[]*ChannelInfo, error)
	ListConditionByStrategyIds(cc ctx.Context, ids []int64) (*[]*StrategyConditionItem, error)

	//
	GetEnterList(cc ctx.Context) ([]int64, error)

	GetWxChannelList(c ctx.Context) ([]*WxChannelInfo, error)

	GetMarketingStrategyIsUsed(c ctx.Context, strategyId int64) (int64, error)

	GetWxCouponDetail(c ctx.Context, id int64) (*WxCouponDatail, error)

	// 查询所有未开始的策略列表
	SelectMarketingWaitStartList(ctx ctx.Context, enterpriseID int64) (*[]MarketingStrategies, error)

	// 查询所有需要结束的策略列表
	SelectMarketingWaitEndList(ctx ctx.Context, enterpriseID int64) (*[]MarketingStrategies, error)

	//通过策略id查询渠道列表
	MarketingDetailChannelsByID(c ctx.Context, strategyId int64) ([]ChannelInfo, error)

	//获取策略的优惠券
	GetStrategyCouponList(c ctx.Context, strategyId int64) ([]StrategyCoupon, error)

	//获取商品详情
	GetSpuDetail(c ctx.Context, SpuId int64) (SpuDetailRes, error)

	ShortUrlAdd(c ctx.Context, SpuId ShortUrlAddReq) error

	ShortUrlDel(c ctx.Context, Id ShortUrlDelReq) error

	TwoLevelIds(c ctx.Context, EnterpriseId ShortUrlAddReq) (IdsRes []ShortUrlDataReq, err error)

	ShortUrlData(c ctx.Context, SpuId ShortUrlDataReq) (ShortUrlDataRes, error)

	//更新策略改成使用事务方式
	TxUpdateMarketingStrategy(c ctx.Context, request DeleteRequest, strategyInfo *MarketingStrategies, strategyCondition *MarketingStrategyCondition, relationInfos []RelationInfo, channelInfos []ChannelInfo, couponIdList []MarketingStrategyCoupons) error

	// BatchUpdateStrategy
	BatchUpdateStrategy(c ctx.Context, strategyIDs []int64, strategyStatus int8, AbnormalPause int8) error
	// GoOnStrategyTemplateIDs
	GoOnStrategyTemplateIDs(c ctx.Context, enterpriseID int64) (*[]*GetGoOnStrategyResp, error)

	AllShortUrlList(c ctx.Context, request ShortUrlListReq) (*[]*ShortUrlDataRes, error)

	//获取企微群发的发送记录
	GetQyMsgSendLogWithIdr(c ctx.Context, EnterpriseID int64, MsgID, externalUserid string) (int64, error)
}
