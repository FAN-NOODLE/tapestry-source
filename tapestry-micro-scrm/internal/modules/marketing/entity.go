package marketing

import (
	"strconv"
	"strings"
)

// 渠道类型
const (
	ChannelTypeSMS      = 1 + iota // 短信
	ChannelTypeMassHair            // 企微群发
	ChannelTypeCoupon              // 自动发券
)

// marketing_strategies表 策略状态
const (
	// 未开始
	StrategyStatusInit = 1 + iota
	// 投放中
	StrategyStatusStart
	// 暂停中
	StrategyStatusStop
	// 已结束
	StrategyStatusEnd
)

// marketing_strategies表 策略执行方式
const (
	// 长期执行
	StrategyMethodALongTime = 1 + iota
	// 立即执行
	MethodImmediately
	// 定时执行
	MethodTiming
)

// marketing_strategies表 策略类型
const (
	// 客户策略
	CustomerStrategy = 1 + iota
	// 内置策略
	BuiltInPolicy
)

// 策略执行状态
const (
	// 未执行
	ExecuteStatusWait = 1 + iota
	// 执行成功
	ExecuteStatusSuccess
	// 执行失败
	ExecuteStatusFail
	// 发送中
	ExecuteStatusIng
)

// 分组中用户类型
const (
	// 有手机号码的用户
	StrategyChannelMobile = 1 + iota
	// 企微好友用户
	StrategyChannelQW
	// 商城访客用户
	StrategyChannelShop
	// 分组公众号用户
	StrategyChannelOfficial
)

// 素材类型
const (
	MaterialTypePicture = 1 + iota
	MaterialTypeLink
	MaterialTypeApplet
)

// marketing_strategy_conditions 触发条件
const (
	// 不限
	StrategyConditionTypeUnlimited = 1 + iota
	// 访问
	StrategyConditionTypeLoad
	// 加购
	StrategyConditionTypeAddCart
	// 下单
	StrategyConditionTypeCreateOrder
	// 支付
	StrategyConditionTypePay
)

// marketing_strategy_conditions 关联类型
const (
	// 首页
	StrategyConditionRelationTypeIndex = 1 + iota
	// 商品
	StrategyConditionRelationTypeGoods
	// 直播
	StrategyConditionRelationTypeLive
	// 活动页
	StrategyConditionRelationTypeActivityPage
)
const (
	ChatTypeGroup = "group"  // 企微群发任务的类型，group表示发送给客户群
	ChatTypeUser  = "single" // single 发送给用户
	// 群发偏移量
	MaxSendNum       = 10000
	ResultStatusFail = 4

	StrategyChannelDelivery = 5

	// 存量数据已执行
	HistoryStatusYes = 1

	// RequestSourceIAP 内购请求发券
	RequestSourceIAP = "iap"
	// RequestSourceShop 商城请求发券
	RequestSourceShop = "shop"

	// 触发条件不限
	StrategyConditionConditionNoTrigger = 1
	// 触发条件访问
	StrategyConditionConditionPage = 2

	WxCoupinStatuxSended = "SENDED"
)

// marketing_strategy_conditions 时间类型
const (
	// 秒
	TimeTypeSecond = 1 + iota
	// 分
	TimeTypeMin
	// 时
	TimeTypeHour
	// 天
	TimeTypeDay
)

// marketing_strategy 策略状态
const (
	// 禁用
	DisableStrategy = 1 + iota
	// 启用
	EnableStrategy
)

// marketing_strategy 策略是否完成
const (
	// 未完成
	NotFinished = 1 + iota
	// 已完成
	Finished
)

// marketing_strategy_conditions 指定类型
const (
	// 全部
	AssignTypeAll = 1 + iota
	AssignTypePart
)

// 是否为异常情况暂停的策略
const (
	// 异常情况暂停的策略
	AbnormalPauseYes = 1 + iota
	AbnormalPauseNo
)

// ListRequest 列表请求
type ListRequest struct {
	StrategyName           string `json:"strategy_name" desc:"策略名称" db:"strategy_name"`
	CustomerGroupID        int64  `json:"customer_group_id" desc:"客户分组ID" db:"customer_group_id"`
	EnterpriseID           int64  `json:"enterprise_id" desc:"企业ID" db:"enterprise_id"`
	AdminID                int64  `json:"admin_id" desc:"操作人ID" db:"admin_id"`
	StrategyType           int8   `json:"strategy_type" desc:"策略类型: 1=客户策略，2=内置策略" db:"strategy_type"`
	StrategyStatus         int8   `json:"strategy_status" desc:"状态: 1=未开始，2=投放中，3=暂停中，4=已结束" db:"strategy_status"`
	StrategyStatusMulti    []int8 `json:"strategy_status_multi" desc:"状态: 1=未开始，2=投放中，3=暂停中，4=已结束"`
	StrategyStatusMultiCnt int    `json:"strategy_status_multi_cnt"`
	StrategyMethod         int8   `json:"strategy_method" desc:"执行方式: 1=长期执行，2=立即执行，3=定时执行" db:"strategy_method"`
	StartTime              int64  `json:"start_time" desc:"执行开始时间" db:"start_time"`
	EndTime                int64  `json:"end_time" desc:"执行结束时间" db:"end_time"`
	StrategyPolicy         int8   `json:"strategy_policy" desc:"触达规则：1=通知一个渠道，2=通知所有渠道" db:"strategy_policy"`
	JobImplement           bool   `json:"job_implement"`
	Channels               int8   `json:"channels" db:"channels"`

	// 是否需要查询渠道数据
	IsChannels int `json:"is_channels"`
	// 执行时间
	CurrentTime int64 `json:"current_time"`
	Page        int   `json:"page"`
	Size        int   `json:"size"`
}

// MarketingStrategy 策略信息
type MarketingStrategy struct {
	ID              int64          `json:"id" desc:"策略ID" db:"id"`
	StrategyName    string         `json:"strategy_name" desc:"策略名称" db:"strategy_name"`
	CustomerGroupID int64          `json:"customer_group_id" desc:"客户分组ID" db:"customer_group_id"`
	EnterpriseID    int64          `json:"enterprise_id" desc:"企业ID" db:"enterprise_id"`
	AdminID         int64          `json:"admin_id" desc:"操作人ID" db:"admin_id"`
	StrategyType    int8           `json:"strategy_type" desc:"策略类型: 1=客户策略，2=内置策略" db:"strategy_type"`
	StrategyMethod  int8           `json:"strategy_method" desc:"执行方式: 1=长期执行，2=立即执行，3=定时执行" db:"strategy_method"`
	StartTime       int64          `json:"start_time" desc:"执行开始时间" db:"start_time"`
	EndTime         int64          `json:"end_time" desc:"执行结束时间" db:"end_time"`
	StrategyEndTime int64          `json:"strategy_end_time" desc:"执行结束时间"`
	StrategyPolicy  int8           `json:"strategy_policy" desc:"触达规则：1=通知一个渠道，2=通知所有渠道" db:"strategy_policy"`
	StrategyStatus  int8           `json:"strategy_status" desc:"状态: 1=未开始，2=投放中，3=暂停中，4=已结束" db:"strategy_status"`
	Status          int8           `json:"status" desc:"状态:1=禁用,2=启用" db:"status"`
	IsFinished      int8           `json:"is_finished" desc:"是否完成；1：未完成，2：已完成" db:"is_finished"`
	Mobile          string         `json:"mobile" desc:"联系方式" db:"mobile"`
	Condition       int8           `json:"condition" desc:"触发条件;1=不限，2=访问，3=加购，4=下单，5=支付" db:"condition"`
	InnerCondition  int8           `json:"inner_condition" desc:"触发条件;1=不设置条件，2=购买金额，3=购买商品" db:"inner_condition"`
	AssignType      int8           `json:"assign_type" desc:"指定类型；1=全部，2=指定" db:"assign_type"`
	RelationType    int8           `json:"relation_type" desc:"关联类型；1=首页，2=商品，3=直播，4=活动页，5=活动" db:"relation_type"`
	BelongType      int8           `json:"belong_type" desc:"所属类型" db:"belong_type"`
	StartPoint      int64          `json:"start_point" desc:"起始点" db:"start_point"`
	EndPoint        int64          `json:"end_point" desc:"结束点" db:"end_point"`
	TimeInterval    int64          `json:"time_interval" desc:"间隔时间（之后执行）" db:"time_interval"`
	TimeType        int8           `json:"time_type" desc:"时间类型1=秒，2=分，3=时，4=天" db:"time_type"`
	Time            int64          `json:"time" desc:"间隔时间/秒" db:"time"`
	HistoryStatus   int            `json:"history_status" db:"history_status" desc:"存量数据是否已经执行"`
	Channels        []ChannelInfo  `json:"channels" desc:"渠道信息"`
	Relations       []RelationInfo `json:"relations" desc:"触发关联信息"`
	SettingDay      int64          `json:"setting_day" desc:"创建时间" db:"setting_day"`
	CreatedAt       int64          `json:"created_at" desc:"创建时间" db:"created_at"`
	UpdatedAt       int64          `json:"updated_at" desc:"更新时间" db:"updated_at"`
	DeletedAt       int64          `json:"deleted_at" desc:"删除时间" db:"deleted_at"`
	Page            int            `json:"page"`
	Size            int            `json:"size"`

	StrategyMethodText string `json:"strategy_method_text"`
	StrategyPolicyText string `json:"strategy_policy_text"`
	StrategyStatusText string `json:"strategy_status_text"`
	StatusText         string `json:"status_text"`
	ChannelsText       string `json:"channels_text" db:"channels_text"`
	CustomerGroupName  string `json:"customer_group_name" db:"customer_group_name"`
	TriggerCount       int64  `json:"trigger_count" db:"trigger_count"`
	IsNotice           int64  `json:"is_notice" db:"is_notice"`
	AbnormalPause      int8   `json:"abnormal_pause" db:"abnormal_pause"`
}

type UpdateMarketingStrategyStatus struct {
	ID             int64 `json:"id" desc:"策略ID" db:"id"`
	StrategyStatus int64 `json:"strategy_status" desc:"状态: 1=未开始，2=投放中，3=暂停中，4=已结束" db:"strategy_status"`
}

// MarketingStrategies 策略信息
type MarketingStrategies struct {
	ID              int64  `json:"id" desc:"策略ID" db:"id"`
	StrategyName    string `json:"strategy_name" db:"strategy_name"`
	EnterpriseID    int64  `json:"enterprise_id" db:"enterprise_id"`
	Mobile          string `json:"mobile" db:"mobile"`
	EndTime         int64  `json:"end_time" db:"end_time"`
	StartTime       int64  `json:"start_time" db:"start_time"`
	IsNotice        int64  `json:"is_notice" db:"is_notice"`
	CustomerGroupID int64  `json:"customer_group_id" desc:"客户分组ID" db:"customer_group_id"`
	AdminID         int64  `json:"admin_id" desc:"操作人ID" db:"admin_id"`
	StrategyType    int8   `json:"strategy_type" desc:"策略类型: 1=客户策略，2=内置策略" db:"strategy_type"`
	StrategyMethod  int8   `json:"strategy_method" desc:"执行方式: 1=长期执行，2=立即执行，3=定时执行" db:"strategy_method"`
	StrategyPolicy  int8   `json:"strategy_policy" desc:"触达规则：1=通知一个渠道，2=通知所有渠道" db:"strategy_policy"`
	StrategyStatus  int8   `json:"strategy_status" desc:"状态: 1=未开始，2=投放中，3=暂停中，4=已结束" db:"strategy_status"`
	Status          int8   `json:"status" desc:"状态:1=禁用,2=启用" db:"status"`
	IsFinished      int8   `json:"is_finished" desc:"是否完成；1：未完成，2：已完成" db:"is_finished"`
	HistoryStatus   int    `json:"history_status" desc:"历史数据是否已经执行" db:"history_status"`
	TriggerCount    int64  `json:"trigger_count" db:"trigger_count"`
	CreatedAt       int64  `json:"created_at" desc:"创建时间" db:"created_at"`
	UpdatedAt       int64  `json:"updated_at" desc:"更新时间" db:"updated_at"`
	DeletedAt       int64  `json:"deleted_at" desc:"删除时间" db:"deleted_at"`
}

// MarketingStrategyConditions 策略条件
type MarketingStrategyCondition struct {
	ID             int64 `json:"id" desc:"策略ID" db:"id"`
	Condition      int8  `json:"condition" desc:"触发条件;1=不限，2=访问，3=加购，4=下单，5=支付" db:"condition"`
	AssignType     int8  `json:"assign_type" desc:"指定类型；1=全部，2=指定" db:"assign_type"`
	RelationType   int8  `json:"relation_type" desc:"关联类型；1=首页，2=商品，3=直播，4=活动页，5=活动" db:"relation_type"`
	BelongType     int8  `json:"belong_type" desc:"所属类型" db:"belong_type"`
	InnerCondition int8  `json:"inner_condition" desc:"触发条件;1=不设置条件，2=购买金额，3=购买商品" db:"inner_condition"`
	StartPoint     int64 `json:"start_point" desc:"起始点" db:"start_point"`
	EndPoint       int64 `json:"end_point" desc:"结束点" db:"end_point"`
	TimeInterval   int64 `json:"time_interval" desc:"间隔时间（之后执行）" db:"time_interval"`
	TimeType       int8  `json:"time_type" desc:"时间类型1=秒，2=分，3=时，4=天" db:"time_type"`
	Time           int64 `json:"time" desc:"间隔时间/秒" db:"time"`
	CreatedAt      int64 `json:"created_at" desc:"创建时间" db:"created_at"`
	UpdatedAt      int64 `json:"updated_at" desc:"更新时间" db:"updated_at"`
	DeletedAt      int64 `json:"deleted_at" desc:"删除时间" db:"deleted_at"`
}

// MarketingInfo 策略信息整理
func MarketingInfo(m *MarketingStrategy) *MarketingStrategy {
	arr := strings.Split(m.ChannelsText, ",")
	var build strings.Builder
	var channels []ChannelInfo
	for key, val := range arr {
		channelTypeInt, _ := strconv.Atoi(val)
		channelId, _ := strconv.ParseInt(val, 10, 64)
		channels = append(channels, ChannelInfo{ID: channelId})
		if key != 0 {
			build.WriteString(",")
		}
		build.WriteString(ChannelTypes[int8(channelTypeInt)])
	}
	channelsText := build.String()

	return &MarketingStrategy{
		ID:                 m.ID,
		StrategyName:       m.StrategyName,
		CustomerGroupID:    m.CustomerGroupID,
		EnterpriseID:       m.EnterpriseID,
		AdminID:            m.AdminID,
		StrategyType:       m.StrategyType,
		StrategyMethod:     m.StrategyMethod,
		StartTime:          m.StartTime,
		EndTime:            m.EndTime,
		StrategyPolicy:     m.StrategyPolicy,
		StrategyStatus:     m.StrategyStatus,
		StrategyMethodText: StrategyMethod[m.StrategyMethod],
		StrategyPolicyText: StrategyPolicies[m.StrategyPolicy],
		StrategyStatusText: StrategyStatus[m.StrategyStatus],
		ChannelsText:       channelsText,
		CustomerGroupName:  m.CustomerGroupName,
		TriggerCount:       m.TriggerCount,
		Status:             m.Status,
		StatusText:         StatusList[m.Status],
		IsFinished:         m.IsFinished,
		SettingDay:         m.SettingDay,
		Channels:           channels,
	}
}

// ChannelInfo 渠道信息
type ChannelInfo struct {
	ID            int64                  `json:"id" desc:"ID" db:"id"`
	StrategyID    int64                  `json:"strategy_id" desc:"策略ID" db:"strategy_id"`
	ChannelType   int8                   `json:"channel_type" desc:"渠道；1=短信，2=企微群发，3=自动发券，4=公众号" db:"channel_type"`
	TextContent   string                 `json:"text_content" desc:"数据表内容" db:"content"`
	JuniorID      int64                  `json:"junior_id" desc:"从属id" db:"junior_id"`
	Order         int64                  `json:"order" desc:"排序:默认1（此处正序排列）" db:"order"`
	CreatedAt     int64                  `json:"created_at" desc:"创建时间" db:"created_at"`
	UpdatedAt     int64                  `json:"updated_at" desc:"更新时间" db:"updated_at"`
	DeletedAt     int64                  `json:"deleted_at" desc:"删除时间" db:"deleted_at"`
	ChannelDetail ChannelDetail          `json:"channel_detail"`
	Content       map[string]interface{} `json:"content" desc:"接受和返回内容"`
}

// ChannelDetail 渠道详情
type ChannelDetail struct {
	Name       string `json:"name"`
	MessageStr string `json:"message_str"`
}

// RelationInfo 关联信息
type RelationInfo struct {
	ID             int64          `json:"id" desc:"ID" db:"id"`
	StrategyID     int64          `json:"strategy_id" desc:"策略ID" db:"strategy_id"`
	RelationID     int64          `json:"relation_id" desc:"关联id" db:"relation_id"`
	RelationDetail RelationDetail `json:"relation_detail"`
	CreatedAt      int64          `json:"created_at" desc:"创建时间" db:"created_at"`
	UpdatedAt      int64          `json:"updated_at" desc:"更新时间" db:"updated_at"`
	DeletedAt      int64          `json:"deleted_at" desc:"删除时间" db:"deleted_at"`
}

// RelationList 关联信息列表
type RelationDetail struct {
	ImgUrl string `json:"img_url"`
	Name   string `json:"name"`
}

// MarketingDetailRequest 营销策略详情
type MarketingDetailRequest struct {
	ID           int64 `json:"id" db:"id"`
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id"`
}

// DeleteRequest 策略删除请求
type DeleteRequest struct {
	ID           int64 `json:"id" db:"id"`
	EnterpriseID int64 `json:"enterprise_id"`
	DeletedAt    int64 `json:"deleted_at" db:"deleted_at"`
}

// UpdateStatusRequest 更新策略状态
type UpdateStatusRequest struct {
	ID             int64 `json:"id" db:"id"`
	EnterpriseID   int64 `json:"enterprise_id" db:"enterprise_id"`
	AdminID        int64 `json:"admin_id" db:"admin_id"`
	StrategyStatus int8  `json:"strategy_status" desc:"状态: 1=未开始，2=投放中，3=暂停中，4=已结束" db:"strategy_status"`
	Status         int8  `json:"status" desc:"状态:1=禁用,2=启用" db:"status"`
	IsFinished     int8  `json:"is_finished" desc:"是否完成；1：未完成，2：已完成" db:"is_finished"`
	AbnormalPause  int8  `json:"abnormal_pause" desc:"是否为异常情况暂停的策略，1:是,0:否" db:"abnormal_pause"`
	UpdatedAt      int64 `json:"updated_at" db:"updated_at"`
	SettingDay     int   `json:"setting_day" db:"setting_day"`
	StartTime      int64 `json:"start_time" db:"start_time"`
	EndTime        int64 `json:"end_time" db:"end_time"`
}

// MarketingStrategyListResponse 营销策略列表
type MarketingStrategyListResponse struct {
	List  []MarketingStrategy `json:"list"`
	Page  int                 `json:"page"`
	Size  int                 `json:"size"`
	Total int64               `json:"total"`
}

type MarketTrigRequest struct {
	EnterpriseID   int64    `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID不能为空"`
	RelationType   int      `json:"relation_type" db:"relation_type"`
	RelationID     int      `json:"relation_id" db:"relation_id"`
	UserUnionID    string   `json:"user_union_id" db:"user_union_id"`
	UserUnionIDs   []string `json:"user_union_ids" db:"user_union_ids"`
	Condition      int      `json:"condition" db:"condition"`
	OrderID        string   `json:"order_id" db:"order_id"`
	AddressName    string   `json:"address_name" db:"address_name"`
	Mobile         string   `json:"mobile" db:"mobile"`
	StrategyStatus int      `json:"strategy_status" db:"strategy_status"`
	GroupIds       []int64  `json:"group_ids" db:"group_ids"`
	UserName       string   `json:"user_name" db:"user_name"`
	VisitTemplate  string   `json:"visit_template"`
	VisitDataId    string   `json:"visit_data_id"`
	Count          int      `json:"count" db:"count"`
	Point          int64    `json:"point"`
	OrderGoodsSpu  []int64  `json:"order_goods_spu_id"`
}

type MarketTrigExecRequest struct {
	EnterpriseID   int64  `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID不能为空"`
	UserID         int64  `json:"user_id" db:"user_id" validate:"required" required:"UserID不能为空"`
	UnionID        string `json:"union_id" db:"union_id" validate:"required" required:"UnionID不能为空"`
	ExternalUserid string `json:"external_userid" db:"external_userid"`
	ChannelID      int64  `json:"channel_id" db:"channel_id"`
	ChannelType    int8   `json:"channel_type" db:"channel_type"`
	JuniorID       int64  `json:"junior_id" db:"junior_id"`
	Mobile         string `json:"mobile" db:"mobile"`
	UserName       string `json:"user_name" db:"user_name"`
	AddressName    string `json:"address_name" db:"address_name"`
	OrderID        string `json:"order_id" db:"order_id"`
	ID             int64  `json:"id" db:"id"`
	StrategyID     int64  `json:"strategy_id" db:"strategy_id"`
	MsgID          string `json:"msg_id" db:"msg_id"`
}

// RecordCronRequest 营销策略触达人群请求
type RecordCronRequest struct {
	CurrentTime int64 `json:"current_time" db:"current_time"`
	Status      int8  `json:"status" db:"status"`
	ID          int64 `json:"id" db:"id"`
	StrategyID  int64 `json:"strategy_id" db:"strategy_id"`
}

// MarketStrategyRecord 策略触达结果
type MarketStrategyRecord struct {
	ID              int64  `json:"id" db:"id"`
	UserId          int64  `json:"user_id" db:"user_id" desc:"enterprise_user的id"`
	UserUnionID     string `json:"user_union_id" db:"user_union_id"`
	ExternalUserid  string `json:"external_userid" db:"external_userid"`
	StrategyID      int64  `json:"strategy_id" db:"strategy_id"`
	ChannelID       int64  `json:"channel_id" db:"channel_id"`
	RelationID      int64  `json:"relation_id" db:"relation_id"`
	Status          int    `json:"status" db:"status"`
	RelationType    int8   `json:"relation_type" db:"relation_type"`
	ChannelType     int8   `json:"channel_type" db:"channel_type"`
	Mobile          string `json:"mobile" db:"mobile"`
	EnterpriseID    int64  `json:"enterprise_id" db:"enterprise_id"`
	OrderID         string `json:"order_id" db:"order_id"`
	AddressUserName string `json:"address_user_name" db:"address_user_name"`
	UserName        string `json:"user_name" db:"user_name"`
	MsgID           string `json:"msg_id" db:"msg_id"`
	CreatedAt       int64  `json:"created_at" db:"created_at"`
	UpdatedAt       int64  `json:"updated_at" db:"updated_at"`
	StartTime       int64  `json:"start_time" db:"start_time"`
	EndTime         int64  `json:"end_time" db:"end_time"`
	TouchedAt       int64  `json:"touched_at" db:"touched_at"`
}

type MarketTrigResponse []*struct {
	EnterpriseID             int64  `json:"enterprise_id" db:"enterprise_id"`
	StrategyStatus           int    `json:"strategy_status" db:"strategy_status"`
	StrategyID               int64  `json:"strategy_id" db:"strategy_id"`
	ChannelID                int64  `json:"channel_id" db:"channel_id"`
	RelationID               int64  `json:"relation_id" db:"relation_id"`
	StartTime                int64  `json:"start_time" json:"start_time"`
	EndTime                  int64  `json:"end_time" db:"end_time"`
	ChannelType              int8   `json:"channel_type" db:"channel_type"`
	StrategyMethod           int    `json:"strategy_method" db:"strategy_method"`
	StrategyPolicy           int    `json:"strategy_policy" db:"strategy_policy"`
	RelationIDs              string `json:"relation_ids" db:"relation_ids"`
	StrategyType             int64  `json:"strategy_type" db:"strategy_type"`
	ConditionsConditions     int64  `json:"condition" db:"condition"`
	ConditionsInnerCondition int64  `json:"inner_condition" db:"inner_condition"`
	ConditionsAssignType     int64  `json:"assign_type" db:"assign_type"`
	ConditionsRelationType   int64  `json:"relation_type" db:"relation_type"`
	ConditionsStartPoint     int64  `json:"start_point" db:"start_point"`
	ConditionsEndPoint       int64  `json:"end_point" db:"end_point"`
	ConditionsTime           int64  `json:"time" db:"time"`
}

type MarketTrigDetailResponse struct {
	EnterpriseID             int64  `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID不能为空"`
	StrategyStatus           int    `json:"strategy_status" db:"strategy_status"`
	StrategyID               int64  `json:"strategy_id" db:"strategy_id"`
	ChannelID                int64  `json:"channel_id" db:"channel_id"`
	RelationID               int64  `json:"relation_id" db:"relation_id"`
	StartTime                int64  `json:"start_time" json:"start_time"`
	EndTime                  int64  `json:"end_time" db:"end_time"`
	ChannelType              int8   `json:"channel_type" db:"channel_type"`
	StrategyMethod           int    `json:"strategy_method" db:"strategy_method"`
	StrategyPolicy           int    `json:"strategy_policy" db:"strategy_policy"`
	RelationIDs              string `json:"relation_ids" db:"relation_ids"`
	StrategyType             int64  `json:"strategy_type" db:"strategy_type"`
	ConditionsConditions     int64  `json:"condition" db:"condition"`
	ConditionsInnerCondition int64  `json:"inner_condition" db:"inner_condition"`
	ConditionsAssignType     int64  `json:"assign_type" db:"assign_type"`
	ConditionsRelationType   int64  `json:"relation_type" db:"relation_type"`
	ConditionsStartPoint     int64  `json:"start_point" db:"start_point"`
	ConditionsEndPoint       int64  `json:"end_point" db:"end_point"`
	ConditionsTime           int64  `json:"time" db:"time"`
}

type MarketTrigExEcRequest struct {
	StrategyID   int64 `json:"strategy_id"`
	EnterpriseID int64 `json:"enterprise_id"`
}

// MarketingAddResponse 策略添加结果
type MarketingAddResponse struct {
	ID int64 `json:"id" desc:"新建D"`
}

// ActivityTypes 活动类型
var ActivityTypes = map[int8]string{
	6:  "好友助力付邮",
	7:  "好友助力买赠",
	12: "限时购",
	13: "助力抽奖",
	14: "拼团",
	17: "先享卡",
	99: "全部类型",
}

// 好友助力付邮
const ActivityTypeHelpPostage = 6

// 好友助力买赠
const ActivityTypeHelpGift = 7

// 限时购
const ActivityTypeFlashSale = 12

// 助力抽奖
const ActivityTypeLuck = 13

// 拼团
const ActivityTypeGroup = 14

// 先享卡
const ActivityTypeDiscountCard = 17

// RelationTypes 关联类型
var RelationTypes = map[int8]string{
	1: "首页",
	2: "商品",
	3: "直播",
	4: "活动页",
	5: "活动",
}

// AssignTypes 指定类型
var AssignTypes = map[int8]string{
	1: "全部",
	2: "指定",
}

// TimeTypes 时间类型
var TimeTypes = map[int8]string{
	1: "秒",
	2: "分",
	3: "时",
	4: "天",
}

// TriggerConditions 触发条件
var TriggerConditions = map[int8]string{
	1: "不限",
	2: "访问",
	3: "加购",
	4: "下单",
	5: "支付",
}

// TriggerInnerConditions 内置触发条件
var TriggerInnerConditions = map[int8]string{
	1: "不设置条件",
	2: "购买金额",
	3: "购买商品",
}

// StrategyStatus 策略状态
var StrategyStatus = map[int8]string{
	1: "未开始",
	2: "投放中",
	3: "暂停中",
	4: "已结束",
}

// StatusList 状态
var StatusList = map[int8]string{
	1: "禁用",
	2: "启用",
}

// StrategyTypes 策略类型
var StrategyTypes = map[int8]string{
	1: "客户策略",
	2: "内置策略",
}

// StrategyMethod 执行方式
var StrategyMethod = map[int8]string{
	1: "长期执行",
	2: "立即执行",
	3: "定时执行",
}

// StrategyPolicies 触达规则
var StrategyPolicies = map[int8]string{
	1: "通知一个渠道",
	2: "通知所有渠道",
}

// ChannelTypes 渠道类型
var ChannelTypes = map[int8]string{
	1: "短信",
	2: "企微群发",
	3: "自动发券",
	4: "公众号",
}

// MarketingDataDetailRequest 营销策略数据详情
type MarketingDataDetailRequest struct {
	ID           int64   `json:"id" db:"id"`
	EnterpriseId int64   `json:"enterprise_id" binding:"required" db:"enterprise_id"`
	Day          int     `json:"day"  db:"day"`
	BeginTime    int     `json:"begin_time" db:"begin_time"`
	EndTime      int     `json:"end_time"  db:"end_time"`
	UserIDs      []int64 `json:"user_ids"  db:"user_ids"`
	SignType     int64   `json:"sign_type"  db:"sign_type"`
}

// 策略详情
type MarketingDataDetail struct {
	Type               int64 `json:"type" db:"type"`
	SendMemberNumber   int64 `json:"send_member_number" db:"send_member_number"`
	ReachMemberNumber  int64 `json:"reach_member_number" db:"reach_member_number"`
	FansNumber         int64 `json:"fans_number" db:"fans_number"`
	VisitorNumber      int64 `json:"visitor_number" db:"visitor_number"`
	AddOrderNumber     int64 `json:"add_order_number" db:"add_order_number"`
	PayOrder           int64 `json:"pay_order" db:"pay_order"`
	FinalAmount        int64 `json:"final_amount" db:"final_amount"`
	RefundAmount       int64 `json:"refund_amount" db:"refund_amount"`
	RefundAmountNumber int64 `json:"refund_amount_number" db:"refund_amount_number"`
}

// 短信日志表
type SmsSendLog struct {
	ID             int64 `json:"id" desc:"短信发送ID" db:"id"`
	EnterpriseID   int64 `json:"enterprise_id" desc:"企业ID" db:"enterprise_id"`
	ChannelValueId int64 `json:"channel_value_id" desc:"渠道id" db:"channel_value_id"`
	UserId         int64 `json:"user_id" desc:"用户id" db:"user_id"`
	SmsTemplateId  int64 `json:"sms_template_id" desc:"短信模板主键" db:"sms_template_id"`
	ChannelType    int64 `json:"channel_type" desc:"短信触发渠道（1：营销策略）" db:"channel_type"`
	Status         int64 `json:"status" desc:"短信发送状态（1：成功 2：失败）" db:"status"`
	CreatedAt      int64 `json:"created_at" desc:"创建时间" db:"created_at"`
	UpdatedAt      int64 `json:"updated_at" desc:"更新时间" db:"updated_at"`
	DeletedAt      int64 `json:"deleted_at" desc:"数据删除时间" db:"deleted_at"`
}

// UserGroup
type UserGroup struct {
	ID                         int64  `json:"id" db:"id"`
	EnterpriseID               int64  `json:"enterprise_id" db:"enterprise_id"`
	Name                       string `json:"name" db:"name"`
	Type                       int64  `json:"type" db:"type"`
	CreatedAt                  int64  `json:"created_at" db:"created_at"`
	UpdatedAt                  int64  `json:"updated_at" db:"updated_at"`
	DeletedAt                  int64  `json:"deleted_at" db:"deleted_at"`
	MemberCount                int64  `json:"member_count"`
	MobileUserCount            int64  `json:"mobile_user_count"`
	QwUserCount                int64  `json:"qw_user_count"`
	ShopUserCount              int64  `json:"shop_user_count"`
	WeChatOfficialAccountCount int64  `json:"chat_official_account_count"`
}

// UserCategory
type UserCategory struct {
	ID              int64   `json:"id" db:"id"`
	TagOrGroupID    int64   `json:"tag_or_group_id" db:"tag_or_group_id"`
	Type            string  `json:"type" db:"type"`
	EnterpriseID    int64   `json:"enterprise_id" db:"enterprise_id"`
	UnionID         string  `json:"union_id" db:"union_id"`
	CreatedAt       int64   `json:"created_at" db:"created_at"`
	UpdatedAt       int64   `json:"updated_at" db:"updated_at"`
	DeletedAt       int64   `json:"deleted_at" db:"deleted_at"`
	ChannelTypeList []int64 `json:"channels"`
	StrategyId      int64   `json:"strategy_id"`
}

type ConsumeMsg struct {
	User       UserInfo        `json:"user"`
	SourceItem string          `json:"sourceItem"`
	SourceUrl  string          `json:"sourceUrl"`
	Args       map[string]Args `json:"args"`
}

type UserInfo struct {
	ID             int64       `json:"id"`
	EnterpriseID   int64       `json:"enterprise_id"`
	EnterpriseHash string      `json:"enterprise_hash"`
	RealName       string      `json:"real_name"`
	Phone          string      `json:"phone"`
	Mobile         string      `json:"mobile"`
	OpenID         string      `json:"open_id"`
	UnionID        string      `json:"union_id" db:"union_id"`
	NickName       string      `json:"nick_name"`
	Address        AddressInfo `json:"address"`
}

type Args struct {
	Scene      string `json:"scene"`
	ID         string `json:"id"`
	LiveID     string `json:"live_id"`
	ActivityID string `json:"activity_id"`
}

type AddressInfo struct {
	DeliveryName string `json:"delivery_name"`
}

// WorkExternalUserRelation 用户外部联系人查询条件
type WorkExternalUserRelation struct {
	EnterpriseId int64  `json:"enterprise_id" db:"enterprise_id"`
	UnionId      string `json:"union_id" db:"union_id"`
}

// CustomersRequest 客户分组
type CustomersRequest struct {
	TagOrGroupId int64  `json:"tag_or_group_id" db:"tag_or_group_id"`
	EnterpriseID int64  `json:"enterprise_id"`
	Type         string `json:"type" db:"type"`
	Page         int    `json:"page"`
	Size         int    `json:"size"`
}

// CustomersResponse 客户分组返回结果
type CustomersResponse struct {
	Mobile         string `json:"mobile" db:"mobile"`
	UnionID        string `json:"union_id" db:"union_id"`
	ExternalUserid string `json:"external_userid" db:"external_userid"`
}

// UserStatisticEntity 用户营销策略实体结构
type UserStatisticEntity struct {
	ID           int64 `db:"id"`
	EnterpriseID int64 `db:"enterprise_id"`
	UserID       int64 `db:"user_id"`
	Count        int   `db:"count"`
	CreatedAt    int64 `db:"created_at"`
	UpdatedAt    int64 `db:"updated_at"`
	DeletedAt    int64 `db:"deleted_at"`
}

// UserEntity user表实体结构
type UserEntity struct {
	ID           int64  `json:"id" db:"id"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id"`
	UnionID      string `json:"union_id" db:"union_id"`
}

// MsgResultInfoRequest 群发消息发送结果
type MsgResultInfoRequest struct {
	EnterpriseID int64  `json:"enterprise_id"`
	MsgID        string `json:"msgid"`
}

// BatchUpdateRecordsRequest 批量修改记录表
type BatchUpdateRecordsRequest struct {
	IDs       []int64 `json:"ids"`
	Status    int     `json:"status" db:"status"`
	MsgID     string  `json:"msg_id" db:"msg_id"`
	UpdatedAt int64   `db:"updated_at"`
}

// QwMsgSendLog 企业微信群发记录表
type QyMsgSendLog struct {
	ID             int64  `json:"id" desc:"id" db:"id"`
	EnterpriseID   int64  `json:"enterprise_id" desc:"enterprise_id" db:"enterprise_id"`
	StrategyID     int64  `json:"strategy_id" desc:"策略id" db:"strategy_id"`
	ChannelID      int64  `json:"channel_id" desc:"渠道id" db:"channel_id"`
	ChatType       string `json:"chat_type" desc:"群发任务类型 single:客户 group:客户群" db:"chat_type"`
	ExternalUserid string `json:"external_userid" desc:"external_userid" db:"external_userid"`
	Sender         string `json:"sender" desc:"sender" db:"sender"`
	MsgID          string `json:"msg_id" desc:"msg_id" db:"msg_id"`
	CreatedAt      int64  `db:"created_at"`
	UpdatedAt      int64  `db:"updated_at"`
	DeletedAt      int64  `db:"deleted_at"`
}

// QyMsgSendResult 群发发送消息
type QyMsgSendResult struct {
	ID             int64  `json:"id" db:"id"`
	EnterpriseID   int64  `json:"enterprise_id" db:"enterprise_id"`
	SendLogID      int64  `json:"send_log_id" db:"send_log_id" desc:"发送记录ID"`
	ExternalUserID string `json:"external_userid" db:"external_userid" desc:"外部联系人userid"`
	ChatID         string `json:"chat_id" db:"chat_id" desc:"外部客户群id"`
	UserID         string `json:"userid" db:"userid" desc:"发送员工userid"`
	Status         int64  `json:"status" db:"status" desc:"发送状态（0-未发送 1-已发送 2-因客户不是好友导致发送失败 3-因客户已经收到其他群发消息导致发送失败）"`
	MsgID          string `json:"msg_id" db:"msg_id" desc:"企业群发消息的id"`
	SendTime       int64  `json:"send_time" db:"send_time" desc:"发送时间"`
	CreatedAt      int64  `json:"created_at" db:"created_at"`
	UpdatedAt      int64  `json:"updated_at" db:"updated_at"`
	DeletedAt      int64  `json:"deleted_at" db:"deleted_at"`
}

// ContentEntity
type ContentEntity struct {
	ContentType int           `json:"type"`
	Text        string        `json:"text"`
	Applet      ContentApplet `json:"applet"`
	Link        ContentLink   `json:"link"`
	Picture     ContentImage  `json:"picture"`
}

type ContentApplet struct {
	Title    string `json:"title"`
	AppID    string `json:"app_id"`
	Page     string `json:"path"`
	CoverUrl string `json:"cover_url"`
}

type ContentLink struct {
	Title    string `json:"title"`
	CoverUrl string `json:"cover_url"`
	Desc     string `json:"description"`
	LinkUrl  string `json:"link_url"`
}

type ContentImage struct {
	Url string `json:"url"`
}

// StrategyChannelListRequest 营销策略渠道检索
type StrategyChannelListRequest struct {
	EnterpriseID    int64   `json:"enterprise_id" db:"enterprise_id"`
	StrategyID      int64   `json:"strategy_id" db:"strategy_id"`
	StrategyIds     []int64 `json:"strategy_ids"`
	ChannelType     int     `json:"channel_type" db:"channel_type"`
	Page            int     `json:"page"`
	Size            int     `json:"size"`
	StrategyIdCount int
}

// StrategyConditionItem 策略配置条件
type StrategyConditionItem struct {
	StrategyID     int64 `json:"strategy_id" db:"strategy_id"`
	Condition      int8  `json:"condition" db:"condition"`
	InnerCondition int8  `json:"inner_condition" db:"inner_condition"`
}

// UserOrderInfo 用户的订单信息
type UserOrderInfo struct {
	DeliveryMobile string `json:"delivery_mobile" db:"mobile"`
}

// WxChannelInfo 场景值渠道
type WxChannelInfo struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Val  string `json:"val" db:"val"`
	Pid  int64  `json:"pid" db:"pid"`
	Type string `json:"type" db:"type"`
}

// 获取商品详情请求
type SpuDetailReq struct {
	PageUrl string `json:"page_url"`
}

type ShortUrlDataReq struct {
	Id int64 `json:"id" db:"id"`
}
type ShortUrlDelReq struct {
	Id int64 `json:"id" db:"id"`
}

type TwoLevelIds struct {
	TwoLevelId []ShortUrlDelReq
}

// 获取短链列表
type ShortUrlListReq struct {
	EnterpriseId int64  `json:"enterprise_id" db:"enterprise_id"`
	OneLevelId   int64  `json:"one_level_id" db:"one_level_id"`
	TwoLevelId   int64  `json:"two_level_id" db:"two_level_id"`
	PageName     string `json:"page_name" db:"page_name"`
	UrlType      int64  `json:"url_type" db:"url_type"`
	Operator     string `json:"operator" db:"operator"`
	CreateStart  int64  `json:"create_start" db:"create_start"`
	CreateEnd    int64  `json:"create_end" db:"create_end"`
	PageNum      int    `json:"page_num"`
	CurrentPage  int    `json:"current_page"`
}

type GetOriginIdReq struct {
	EnterpriseId int64 `json:"enterprise_id" db:"enterprise_id"`
}

// 短链以及配置信息添加
type ShortUrlAddReq struct {
	EnterpriseId int64  `json:"enterprise_id" db:"enterprise_id"`
	Id           int64  `json:"id" db:"id"`
	PageUrl      string `json:"page_url"` // 二级渠道page_url
	Operator     string `json:"operator" db:"operator"`
	OneLevelId   int64  `json:"one_level_id" db:"one_level_id"`
	TwoLevelId   int64  `json:"two_level_id" db:"two_level_id"`
	OneLevelName string `json:"one_level_name" db:"one_level_name"` // 一级渠道名称
	TwoLevelName string `json:"two_level_name" db:"two_level_name"` // 二级渠道名称
	ShortUrl     string `json:"short_url" db:"short_url"`           // 通用短链
	WxUrl        string `json:"wx_url" db:"weixin_url"`             // 生成的通用短链对应的微信快捷链接
	PageName     string `json:"page_name" db:"page_name"`
	UrlType      int64  `json:"url_type" db:"url_type"`
	MiddleData
	TmpMiddleData []byte `json:"-" db:"config_params"`
	CreatedAt     int64  `json:"created_at" db:"created_at"`
}
type MiddleData struct {
	MiddleConfig
	SpuData []*GoodsData `json:"spu_data"`
	ActiveData
}
type MiddleConfig struct {
	BrandImg          string `json:"brand_img"`           // 外链中间页品牌图
	CardTitle         string `json:"card_title"`          // 卡片标题
	CardDesc          string `json:"card_desc"`           // 卡片活动说明
	BackgroundImg     string `json:"background_img"`      // 外链中间页背景图
	IsHide            int64  `json:"is_hide"`             // 是否显示按钮
	ButtonName        string `json:"button_name"`         // 按钮名称
	ButtonColor       string `json:"button_color"`        // 按钮颜色
	ButtonFontColor   string `json:"button_font_color"`   // 按钮字体颜色
	LocationDesc      string `json:"location_desc"`       // 跳转说明
	LocationFontColor string `json:"location_font_color"` // 跳转说明文字颜色
	AdditionalType    int64  `json:"additional_type"`     // 1商品2活动
}
type GoodsData struct {
	SpuOneLevelId int64   `json:"spu_one_level_id"`
	SpuTwoLevelId int64   `json:"spu_two_level_id"`
	SpuId         int64   `json:"spu_id"`
	SpuImg        string  `json:"spu_img"`
	SpuName       string  `json:"spu_name"`
	SpuPageUrl    string  `json:"spu_page_url"`
	WxSpuPageUrl  string  `json:"wx_spu_page_url"` // 微信生成的可在手机跳转的链接
	SkuAttribute  string  `json:"sku_attribute"`
	SkuPrice      float32 `json:"sku_price"`
}

type ActiveData struct {
	ActiveImg        string `json:"active_img"`
	ActivePageUrl    string `json:"active_page_url"`
	WxActivePageUrl  string `json:"wx_active_page_url"` // 微信生成的可在手机跳转的链接
	ActiveOneLevelId int64  `json:"active_one_level_id"`
	ActiveTwoLevelId int64  `json:"active_two_level_id"`
}

type ShortUrlListRes struct {
	ShortCount int64               `json:"total"`
	ShortList  *[]*ShortUrlDataRes `json:"list"`
}

type OriginIdRes struct {
	OriginId string `json:"origin_id" db:"origin_id"`
}
type ShortUrlDataRes struct {
	Id           int64  `json:"id" db:"id"`
	OneLevelId   int64  `json:"one_level_id" db:"one_level_id"`
	OneLevelName string `json:"one_level_name" db:"one_level_name"`
	TwoLevelId   int64  `json:"two_level_id" db:"two_level_id"`
	TwoLevelName string `json:"two_level_name" db:"two_level_name"`
	PageName     string `json:"page_name" db:"page_name"`
	UrlType      int64  `json:"url_type" db:"url_type"`
	Operator     string `json:"operator" db:"operator"`
	ShortUrl     string `json:"short_url" db:"short_url"`
	WeixinUrl    string `json:"weixin_url" db:"weixin_url"`
	CreatedAt    string `json:"created_at" db:"created_at"`

	ConfigParams string      `json:"-" db:"config_params"`
	ConfigParam  interface{} `json:"config_params"`
}

type SpuDetailRes struct {
	SpuId        int64   `json:"spu_id" db:"spu_id"`
	SpuImg       string  `json:"spu_img" db:"picture"`
	SkuAttribute string  `json:"sku_attribute" db:"sku_attribute"`
	SkuPrice     float32 `json:"sku_price" db:"sell_price"`
	SkuName      string  `json:"spu_name" db:"spu_name"`
}

type SkuAttribute struct {
	AttributeName      string `json:"attribute_name"`
	AttributeItemValue string `json:"attribute_item_value"`
}

type (
	// 自动发券，检测券是否已使用中
	MarketingStrategyCoupons struct {
		Id         int64 `json:"id" db:"id"`
		StrategyId int64 `json:"strategy_id" db:"strategy_id"`
		WxCouponId int64 `json:"wx_coupon_id" db:"wx_coupon_id"`
		CreatedAt  int64 `json:"created_at" db:"created_at"`
		UpdatedAt  int64 `json:"updated_at" db:"updated_at"`
		DeletedAt  int64 `json:"deleted_at" db:"deleted_at"`
	}

	WxCouponDatail struct {
		CouponStockId string `json:"coupon_stock_id" db:"coupon_stock_id"`
		StartTime     int64  `json:"start_time" db:"start_time"`
		EndTime       int64  `json:"end_time" db:"end_time"`
		DeletedAt     int64  `json:"deleted_at" db:"deleted_at"`
	}

	MarketingCouponIsUsedRequest struct {
		WxCouponId int64 `json:"wx_coupon_id" db:"wx_coupon_id"`
	}
	MarketingCouponIsUsedResponse struct {
		WxCouponId int64 `json:"wx_coupon_id" db:"wx_coupon_id"`
		IsUsed     int   `json:"is_used"`
	}

	// 自动发券内容体
	AutoCouponContent struct {
		QueryCoupon          int                 `json:"query_coupon"`
		CustomerChannel      int                 `json:"customer_channel"`
		CustomerChannelValue map[string][]string `json:"customer_channel_value"`
		CouponId             []int64             `json:"coupon_id"`
		CouponStockId        []string            `json:"coupon_stock_id"`
		CouponList           []CouponDetail      `json:"coupon_list"`
		CouponDetailList     []CouponDetail      `json:"coupon_detail_list"`
	}

	// 优惠券匹配请求
	CouponMatchRequest struct {
		EnterpriseID     int64   `json:"enterprise_id" desc:"企业ID" db:"enterprise_id"`
		APIKey           string  `json:"apikey"`
		Scene            int64   `json:"scene"`
		Page             string  `json:"page"`
		UID              int64   `json:"uid"`
		AppID            string  `json:"appid"`
		OpenID           string  `json:"openid"`
		Options          string  `json:"options"`
		Source           string  `json:"source"`
		CustomerGroupIds []int64 `json:"customer_group_ids,omitempty"`
		Query            map[string]interface{}
		Now              int64  `json:"now"`
		UnionId          string `json:"unionid"`
	}

	CouponDetail struct {
		ID                int64  `json:"id" db:"id"`
		EnterpriseId      int64  `json:"enterprise_id" db:"enterprise_id"`
		CouponId          int64  `json:"coupon_id" db:"coupon_id"`
		CouponStockId     string `json:"coupon_stock_id" db:"coupon_stock_id"`
		CouponName        string `json:"coupon_name" db:"coupon_name"`
		CouponValue       int64  `json:"coupon_value" db:"coupon_value"`
		CouponMininumn    int64  `json:"coupon_mininumn" db:"coupon_mininumn"`
		CouponStatus      int8   `json:"coupon_status" db:"coupon_status"`
		CouponStockStatus int8   `json:"coupon_stock_status" db:"coupon_stock_status"`
		CouponTotal       int64  `json:"coupon_total" db:"coupon_total"`
		UseDayLimit       int64  `json:"use_day_limit" db:"use_day_limit"`
		MaxQuota          int64  `json:"max_quota" db:"max_quota"`
		IsSendNum         int64  `json:"is_send_num" db:"is_send_num"`
		IsAll             int64  `json:"is_all" db:"is_all"`
		CouponBudget      int64  `json:"coupon_budget" db:"coupon_budget"`
		Scope             string `json:"scope" db:"scope"`
		FavorType         string `json:"favor_type" db:"favor_type"`
		StartTime         int64  `json:"start_time" db:"start_time"`
		EndTime           int64  `json:"end_time" db:"end_time"`
		IsOverlay         int64  `json:"is_overlay" db:"is_overlay"`
		OverlayCoupons    string `json:"overlay_coupons" db:"overlay_coupons"`
		Channel           int64  `json:"channel" db:"channel"`
		Range             string `json:"range" db:"range"`
		RangeHide         string `json:"range_hide" db:"range_hide"`
		TagGoods          string `json:"tag_goods" db:"tag_goods"`
		TagName           string `json:"tag_name" db:"tag_name"`
		TagRuleType       int64  `json:"tag_rule_type" db:"tag_rule_type"`
		CreatedAt         int64  `json:"created_at" db:"created_at"`
		UpdatedAt         int64  `json:"updated_at" db:"updated_at"`
		DeletedAt         int64  `json:"deleted_at" db:"deleted_at"`
		StockUseRule      struct {
			MaxCoupons        int64 `json:"max_coupons"`
			MaxAmount         int64 `json:"max_amount"`
			MaxAmountByDay    int64 `json:"max_amount_by_day"`
			FixedNormalCoupon struct {
				CouponAmount       int64 `json:"coupon_amount"`
				TransactionMinimum int64 `json:"transaction_minimum"`
			} `json:"fixed_normal_coupon"`
			MaxCouponsPerUser int64  `json:"max_coupons_per_user"`
			CouponType        string `json:"coupon_type"`
			CombineUse        bool   `json:"combine_use"`
		} `json:"stock_use_rule"`
	}
	// 自动发券条件
	AutoSendCouponMatch struct {
		Now                 int64  `json:"now"`
		Uid                 int64  `json:"uid"`
		StartTime           int64  `json:"start_time" db:"start_time"`
		StrategyId          int64  `json:"strategy_id" db:"strategy_id"`
		CouponContent       string `json:"coupon_content" db:"content" desc:"发券内容体"`
		Acc                 AutoCouponContent
		TriggerCondition    int     `json:"trigger_condition" db:"condition" desc:"触点（1:不限,2:访问）"`
		TriggerRelationType int     `json:"trigger_relation_type" db:"relation_type" desc:"触点为2，访问什么页面"`
		TriggerAssignType   int     `json:"trigger_assign_type" db:"assign_type" desc:"触点为2,访问非首页时，全部还是访问具体页面；1=全部，2=指定"`
		TriggerRelationsIds []int64 `json:"trigger_relations_ids" desc:"指定页面的ID"`
		MatchVersion        int
	}

	// 用户已经获得的某个票的数量
	UserHasCouponNum struct {
		CouponID int64 `json:"coupon_id" db:"coupon_id"`
		C        int   `json:"c" db:"c"`
	}

	MarketingStrategiesCouponLog struct {
		Id                  int64  `json:"id" db:"id"`
		EnterpriseId        int64  `json:"enterprise_id" db:"enterprise_id"`
		ShopUserId          int64  `json:"shop_user_id" db:"shop_user_id"`
		Openid              string `json:"openid" db:"openid"`
		UnionId             string `json:"union_id" db:"union_id"`
		MarketingStrategyId int64  `json:"marketing_strategy_id" db:"marketing_strategy_id"`
		ScrmWxCouponId      int64  `json:"scrm_wx_coupon_id" db:"scrm_wx_coupon_id"`
		WxStockId           string `json:"wx_stock_id" db:"wx_stock_id"`
		WxCouponId          int64  `json:"wx_coupon_id" db:"wx_coupon_id"`
		Type                int    `json:"type" db:"type"`
		CreatedAt           int64  `json:"created_at" db:"created_at"`
		UpdatedAt           int64  `json:"updated_at" db:"updated_at"`
		DeletedAt           int64  `json:"deleted_at" db:"deleted_at"`
	}

	MarketingStrategyRecord struct {
		Id              int64  `json:"id" db:"id"`
		UserId          int64  `json:"user_id" db:"user_id"`
		EnterpriseId    int64  `json:"enterprise_id" db:"enterprise_id"`
		UnionId         string `json:"union_id" db:"union_id"`
		StrategyId      int64  `json:"strategy_id" db:"strategy_id"`
		ChannelType     int    `json:"channel_type" db:"channel_type"`
		RelationId      int64  `json:"relation_id" db:"relation_id"`
		ChannelId       int64  `json:"channel_id" db:"channel_id"`
		Status          int    `json:"status" db:"status"`
		UserName        string `json:"user_name" db:"user_name"`
		Mobile          string `json:"mobile" db:"mobile"`
		AddressUserName string `json:"address_user_name" db:"address_user_name"`
		OrderId         string `json:"order_id" db:"order_id"`
		StartTime       int64  `json:"start_time" db:"start_time"`
		EndTime         int64  `json:"end_time" db:"end_time"`
		CreatedAt       int64  `json:"created_at" db:"created_at"`
		UpdatedAt       int64  `json:"updated_at" db:"updated_at"`
		TouchedAt       int64  `json:"touched_at" db:"touched_at"`
		DeletedAt       int64  `json:"deleted_at" db:"deleted_at"`
		ExternalUserid  int64  `json:"external_userid" db:"external_userid"`
		MsgId           string `json:"msg_id" db:"msg_id"`
	}

	// CouponDataDetailResponse 自动发券策略数据详情
	CouponDataDetailResponse struct {
		Data  CouponData  `json:"list" desc:"数据明细"`
		Total CouponTotal `json:"total" desc:"数据合计"`
	}

	CouponData struct {
		ID          int64              `json:"id" desc:"策略ID"`
		ChannelName string             `json:"channel_name" desc:"渠道类型名称"`
		SendNum     int64              `json:"send_num" desc:"发券人数"`
		CouponList  []CouponDataDetail `json:"coupon_list" desc:"优惠券数据详情"`
	}

	CouponTotal struct {
		ChannelName string           `json:"channel_name" desc:"渠道类型名称"`
		SendNum     int64            `json:"send_num" desc:"发券人数"`
		CouponList  CouponDataDetail `json:"coupon_list" desc:"优惠券数据详情"`
	}

	CouponDataDetail struct {
		StrategyCoupon
		TouchNum          int64   `json:"touch_num" desc:"触达人数"`
		IsSendNum         int64   `json:"is_send_num" desc:"发送张数"`
		UseNum            int64   `json:"use_num"desc:"核销张数"`
		PayNum            int64   `json:"pay_num"  desc:"支付买家数"`
		CashFee           float64 `json:"cash_fee" desc:"实付金额"`
		RefundCustomerNum int64   `json:"refund_customer_num"  desc:"退款买家数"`
		RefundFee         float64 `json:"refund_fee" desc:"退款金额"`
	}

	StrategyCoupon struct {
		CouponName string `json:"coupon_name" db:"coupon_name" desc:"微信优惠券名称"`
		StockID    string `json:"stock_id" db:"stock_id" desc:"批次号"`
	}

	StrategyCouponStatistics struct {
		StockID string `json:"stock_id" db:"wx_stock_id" desc:"批次号"`
		SendNum int64  `json:"send_num" db:"send_num" desc:"发送张数"`
	}

	StrategyCouponTouchStatistics struct {
		StockID  string `json:"stock_id" db:"wx_stock_id" desc:"批次号"`
		TouchNum int64  `json:"touch_num" db:"touch_num" desc:"触达人数"`
	}

	UserCouPonTotalNum struct {
		CouponId int64 `json:"coupon_id" db:"coupon_id"`
		C        int   `json:"c" db:"c"`
	}
)

// CustomerLoadCarQueueData 客户浏览加购队列数据结构
type CustomerLoadCarQueueData struct {
	Job  string `json:"job"`
	Data *ConsumeMsg
}

// GetGoOnStrategyResp
type GetGoOnStrategyResp struct {
	StrategyID   int64  `json:"id" db:"id"`
	JuniorID     int64  `json:"junior_id" db:"junior_id"`
	Mobile       string `json:"mobile" desc:"联系方式" db:"mobile"`
	StrategyName string `json:"strategy_name" desc:"策略名称" db:"strategy_name"`
}

// BatchUpdateStrategy
type BatchUpdateStrategy struct {
	StrategyIDs    []int64 `json:"strategy_id" desc:"策略ID"`
	StrategyStatus int8    `json:"enterprise_id" desc:"企业ID" db:"enterprise_id"`
	UpdatedAt      int64   `json:"updated_at" desc:"更新时间" db:"updated_at"`
	AbnormalPause  int8    `json:"abnormal_pause" desc:"是否为异常情况暂停的策略" db:"abnormal_pause"`
}

type TagUserCount struct {
	MemberCount                int64 `json:"member_count"`
	MobileUserCount            int64 `json:"mobile_user_count"`
	QwUserCount                int64 `json:"qw_user_count"`
	ShopUserCount              int64 `json:"shop_user_count"`
	WeChatOfficialAccountCount int64 `json:"chat_official_account_count"`
	DeliveryUserCount          int64 `json:"delivery_user_count"`
}
