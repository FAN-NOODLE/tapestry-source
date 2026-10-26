package microscrm

import (
	"fmt"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gw-scrm/internal/config"
	"gw-scrm/internal/pkg/sdk"
)

var urlEnum = struct {
	TagGroupDetail                 string
	TagList                        string
	TagSearch                      string
	CreateGroupUrl                 string
	CreateTagUrl                   string
	UpdateAutoTagUrl               string
	CreateUserGroupUrl             string
	GroupUserListUrl               string
	GroupUserUserCount             string
	TagUpdateUrl                   string
	TagInfoUrl                     string
	TagRandUrl                     string
	TagGroupDeleteUrl              string
	TagUpdateTimeUrl               string
	TagDeleteUrl                   string
	TagSyncQw                      string
	DepartmentListUrl              string
	DepartmentDetailUrl            string
	MemberListUrl                  string
	MemberDetailUrl                string
	UpdateRecordUrl                string
	SyncContactUrl                 string
	WorkExternalRelation           string
	GroupMemberRelation            string
	UpdateMessageUrl               string
	MessageListUrl                 string
	RepeatCustomerUrl              string
	RepeatCustomerMemberUrl        string
	MarketingAddUrl                string
	MarketingListUrl               string
	MarketingDetailUrl             string
	MarketingUpdateUrl             string
	MarketingDeleteUrl             string
	LossCustomersListUrl           string
	TemplateListUrl                string
	TemplateSignListUrl            string
	CreateSignUrl                  string
	TemplateCreateUrl              string
	TemplateArraignmentUrl         string
	ParamsListUrl                  string
	SignListUrl                    string
	TemplateDeleteUrl              string
	TemplateEditUrl                string
	TemplateInfoUrl                string
	TemplateSignDeleteUrl          string
	TemplateSendUrl                string
	TemplateTotalUrl               string
	MomentListUrl                  string
	MomentBaListUrl                string
	DepartmentMomentListUrl        string
	UserLostConfig                 string
	UserLostInfo                   string
	MemberTopicUrl                 string
	TopicCreateUrl                 string
	TopicEditUrl                   string
	TopicDeleteUrl                 string
	TopicDisplayUrl                string
	ContentListUrl                 string
	MomentMemberListUrl            string
	AppearanceGetUrl               string
	AppearanceEditUrl              string
	DepartmentMomentList           string
	SelectorTextUrl                string
	JsConfigUrl                    string
	CustomerListUrl                string
	TopicMemberDeleteUrl           string
	TopicMemberEditUrl             string
	MarketingStrategyDataDetailUrl string
	EditCustomerRemarkUrl          string
	MarketingUpdateStatusUrl       string
	MarketingStrategyTriggerUrl    string
	UpdateSettingDayUrl            string
	MarketingUserGroupListUrl      string
	TagGroupUserCountUrl           string
	CustomerTagListUrl             string
	SensitiveListUrl               string
	SensitiveSaveUrl               string
	SensitiveDeleteUrl             string
	SensitiveDetailUrl             string
	SensitiveDetailListUrl         string
	SettingsUpdateUrl              string
	ChatMsgTypeList                string `desc:"聊天-消息类型列表"`
	ChatMsgList                    string `desc:"聊天-消息列表"`
	ChatGroupMsgList               string `desc:"聊天-群消息列表"`
	ChatMsgListByID                string `desc:"根据消息ID找附近的消息"`
	ChatAgreeInfo                  string `desc:"聊天-会话同意信息"`
	ChatMemberList                 string `desc:"聊天-员工列表"`
	ChatObjectToMemberList         string `desc:"聊天-员工的员工列表"`
	ChatObjectToExternalList       string `desc:"聊天-员工的外部联系人列表"`
	ChatObjectToGroupList          string `desc:"聊天-员工的社群列表"`
	GetUserInfoUrl                 string
	ConfigInfoUrl                  string `desc:"提醒设置-获取"`
	ConfigSaveUrl                  string `desc:"提醒设置-保存"`
	SettingsInfoUrl                string `desc:"会话存档配置-获取"`
	OvertimeListUrl                string `desc:"聊天超时记录-列表"`
	CustomerListConfigUrl          string
	ChannelCreateUrl               string
	ChannelListUrl                 string
	ChannelDetailUrl               string
	ChannelEditUrl                 string
	ContactListUrl                 string
	ContactCreateUrl               string
	ContactUpdateUrl               string
	ContactDeleteUrl               string
	ContactDetailsUrl              string
	DownloadQrCodeUrl              string
	WelcomeCreateUrl               string
	WelcomeDetailUrl               string
	WelcomeDefaultUrl              string
	WelcomeEditUrl                 string
	WelcomeDeleteUrl               string
	WelcomeListUrl                 string
	WelcomeValidateUrl             string
	ChannelsDataUrl                string
	OrderGoodsUrl                  string
	BelongToChatsUrl               string
	CacheMediaUrl                  string
	ExternalWelcomeTypeUrl         string `desc:"所有外部条件"`
	UserGroupDetailUrl             string
	UserGroupExportDataUrl         string
	UserGroupEditUrl               string
	UserGroupDeleteUrl             string
	StatisticsAllByDayUrl          string
	StatisticsAllBySumUrl          string
	StatisticsAllByChannelUrl      string
	StatisticsItemByDayUrl         string
	MemberLevelListURL             string
	GroupMsgListUrl                string
	GroupSendAddUrl                string
	GroupSendEditUrl               string
	TaskInfoUrl                    string
	TaskDetailUrl                  string
	TaskMemberListUrl              string
	WxCouponAvailableUrl           string
	WxChannelUrl                   string
	GouponIsUsedUrl                string `desc:"券是否在使用中"`
	MarketingDataCouponUrl         string `desc:"获取自动发券策略数据详情"`
	CouponMatchUrl                 string `desc:"券匹配"`
	ExportExchangeRecordUrl        string
	FissionActivityListUrl         string `desc:"裂变活动-列表"`
	FissionActivityAddUrl          string `desc:"裂变活动-新增"`
	FissionActivityUpdateUrl       string `desc:"裂变活动-编辑"`
	FissionActivityDeleteUrl       string `desc:"裂变活动-删除"`
	FissionUpdateStatusUrl         string `desc:"裂变活动-更新状态"`
	FissionActivityInfoUrl         string `desc:"裂变活动-详情"`
	ChatGroupListUrl               string `desc:"选择社群列表"`
	MutualActivityListUrl          string `desc:"选择互斥活动列表"`
	ActivityDetailUrl              string `desc:"裂变活动-活动详情"`
	ActivityDetailListUrl          string `desc:"裂变活动-详情列表"`
	InviteDetailUrl                string `desc:"发起人邀请明细"`
	ExchangeRecordUrl              string `desc:"发起人兑奖明细"`
	OpenExchangeWarnUrl            string `desc:"奖品兑换提醒"`
	ExchangeTotalUrl               string `desc:"奖品兑换记录总数"`
	ActivityListSimpleUrl          string `desc:"活动列表 全部有效活动"`
	ExchangeConfigMemberUrl        string `desc:"获取兑换提醒的额外员工"`
	GetSpuDetail                   string `desc:"获取商品详情"`
	ShortUrlAdd                    string `desc:"短链添加"`
	ShortUrlData                   string
	GetShortUrlList                string
	ShortUrlDel                    string
	GetTwoLevelIds                 string
	GetOriginId                    string
	SyncEterpriseUserToESUrl       string `desc:"全量更新ES索引"`
	GetAllShortUrlListUrl          string
	GetTemplateList                string
	CreateTemplateInfo             string
	DelTemplateInfo                string
	CheckTemplateInfo              string
	GetTemplateInfo                string
	GetStyleConfigUrl              string `desc:"获取配色列表"`
	GetCheckStyleConfigUrl         string `desc:"获取默认配色"`
	SetStyleConfigUrl              string `desc:"添加配色"`
	AddStyleConfigUrl              string `desc:"设置默认配色"`
	DelStyleConfigUrl              string `desc:"删除默认配色"`
	GetAllWelcomeConfigUrl         string `desc:"获取所有欢迎语"`
	GetAllKDConfigUrl              string `desc:"获取所有欢迎语"`
	DelWelcomeConfigUrl            string `desc:"删除欢迎语"`
	WelcomeNewEditConfigUrl        string `desc:"删除欢迎语"`
	WelcomeDetailNewConfigUrl      string `desc:"欢迎语详情"`
	SaveOpenOrderUrl               string `desc:"聊天超时记录-列表"`

	WorkGroupChatListUrl           string `desc:"获取群列表"`
	WorkGroupChatInfoUrl           string `desc:"获取群列表"`
	WorkGroupChatSettingUrl        string `desc:"社群设置"`
	WorkGroupChatFiltersUrl        string `desc:"社群自定义筛选下拉"`
	WorkGroupChatTypesUrl          string `desc:"群类型列表"`
	WorkGroupChatTypeManageUrl     string `desc:"社群类型管理"`
	WorkGroupChatMembersUrl        string `desc:"社群成员明细"`
	WorkGroupChatDataGridUrl       string `desc:"数据统计"`
	WorkGroupChatUploadSettingUrl  string `desc:"批量上传群设置"`
	WorkGroupChatBatchSettingUrl   string `desc:"批量设置社群信息"`
	WorkGroupChatDataUpdateTimeUrl string `desc:"社群数据最新更新时间"`
	WorkGroupChatBatchSyncUrl      string `desc:"批量同步社群基本信息"`

	DrainageSimpleChannelListUrl string `desc:"渠道筛选列表"`

	SimpleMembersOfDepartmentsUrl string `desc:"多组织架构员工信息"`
	TaskMemberDetailUrl           string
	TaskMemberSendListUrl         string
}{
	TagGroupDetail: "%s/micro-scrm/v1/tag/mini/tag-group-detail-list",
	// 获取标签组列表
	TagList: "%s/micro-scrm/v1/tag/list",
	// 标签/标签分组搜索
	TagSearch: "%s/micro-scrm/v1/tag/search",
	// 创建标签组
	CreateGroupUrl: "%s/micro-scrm/v1/tag/create-group",
	// 新增手动标签
	CreateTagUrl: "%s/micro-scrm/v1/tag/create-tag",
	// 更新自动标签设置
	UpdateAutoTagUrl: "%s/micro-scrm/v1/tag/update",
	// 创建客户分组
	CreateUserGroupUrl: "%s/micro-scrm/v1/user/create-group",
	// 客户分组列表
	GroupUserListUrl: "%s/micro-scrm/v1/user/group-list",
	// 客户分组预计算人数
	GroupUserUserCount: "%s/micro-scrm/v1/user/group-precalculate",
	// 修改子标签
	TagUpdateUrl: "%s/micro-scrm/v1/tag/manual/update",
	// 标签组详情
	TagInfoUrl: "%s/micro-scrm/v1/tag/tag-info",
	// 标签组排序
	TagRandUrl: "%s/micro-scrm/v1/tag/group-sort",
	// 标签组删除接口
	TagGroupDeleteUrl: "%s/micro-scrm/v1/tag/tag-group-delete",
	// 刷新标签更新时间
	TagUpdateTimeUrl: "%s/micro-scrm/v1/tag/update-time",
	// 标签开启同步企业微信、
	TagSyncQw: "%s/micro-scrm/v1/tag/update-sync",
	// 子标签删除接口
	TagDeleteUrl: "%s/micro-scrm/v1/tag/manual/delete",
	// 获取企业列表
	DepartmentListUrl: "%s/micro-scrm/v1/departments/list",
	// 获取企业详情
	DepartmentDetailUrl: "%s/micro-scrm/v1/departments/detail",
	// 获取成员列表
	MemberListUrl: "%s/micro-scrm/v1/members/list",
	// 获取成员详情
	MemberDetailUrl: "%s/micro-scrm/v1/members/detail",
	// 获取更新记录
	UpdateRecordUrl: "%s/micro-scrm/v1/departments/update-record",
	// 同步通讯录
	SyncContactUrl: "%s/micro-scrm/v1/members/sync",
	// 外部联系人数据变更计算
	WorkExternalRelation: "%s/micro-scrm/v1/work-external/calc",
	// 群成员数据变更计算
	GroupMemberRelation: "%s/micro-scrm/v1/group-member/calc",
	// 添加营销策略
	MarketingAddUrl: "%s/micro-scrm/v1/marketing-strategies/store",
	// 营销策略列表
	MarketingListUrl: "%s/micro-scrm/v1/marketing-strategies/list",
	// 营销策略详情
	MarketingDetailUrl: "%s/micro-scrm/v1/marketing-strategies/detail",
	// 营销策略更新
	MarketingUpdateUrl: "%s/micro-scrm/v1/marketing-strategies/update",
	// 营销策略删除
	MarketingDeleteUrl: "%s/micro-scrm/v1/marketing-strategies/delete",
	// 更新删除消息配置
	UpdateMessageUrl: "%s/micro-scrm/v1/customer/update-message",
	// 删除消息列表
	MessageListUrl: "%s/micro-scrm/v1/customer/message-list",
	// 重复客户
	RepeatCustomerUrl: "%s/micro-scrm/v1/customer/repeat-customer",
	// 重复客户-所属员工
	RepeatCustomerMemberUrl: "%s/micro-scrm/v1/customer/repeat-find-member",
	// 流失客户列表
	LossCustomersListUrl: "%s/micro-scrm/v1/customer/loss-customers-list",

	// 短信模板列表
	TemplateListUrl: "%s/micro-scrm/v1/sms-template/template-list",
	// 短信签名列表
	TemplateSignListUrl: "%s/micro-scrm/v1/customer/loss-customers-list",
	// 创建短信签名
	CreateSignUrl: "%s/micro-scrm/v1/sms-template/create-sign",
	// 创建短信模板
	TemplateCreateUrl: "%s/micro-scrm/v1/sms-template/create",
	// 短信模板提审
	TemplateArraignmentUrl: "%s/micro-scrm/v1/sms-template/arraignment",
	// 短信参数列表
	ParamsListUrl: "%s/micro-scrm/v1/sms-template/params-list",
	// 短信签名列表
	SignListUrl: "%s/micro-scrm/v1/sms-template/sign-list",
	// 删除短信模板
	TemplateDeleteUrl: "%s/micro-scrm/v1/sms-template/delete",
	// 编辑短信模板
	TemplateEditUrl: "%s/micro-scrm/v1/sms-template/edit",
	// 短信模板详情
	TemplateInfoUrl: "%s/micro-scrm/v1/sms-template/info",
	// 删除短信签名
	TemplateSignDeleteUrl: "%s/micro-scrm/v1/sms-template/sign-delete",
	// 短信发送
	TemplateSendUrl: "%s/micro-scrm/v1/sms-template/send",
	// 企业短信统计
	TemplateTotalUrl: "%s/micro-scrm/v1/sms-template/total",

	// 朋友圈选择员工
	MomentListUrl: "%s/micro-scrm/v1/members/moment-list",

	// 朋友圈选择BA员工
	MomentBaListUrl: "%s/micro-scrm/v1/members/ba-list",
	// 朋友圈选择部门
	DepartmentMomentListUrl: "%s/micro-scrm/v1/departments/moment-list",
	// 导购丢失客户推送消息配置
	UserLostConfig: "%s/micro-scrm/v1/departments/user-lost-config",
	// 导购丢失客户推送消息配置
	UserLostInfo: "%s/micro-scrm/v1/departments/user-lost-info",
	// 一个员工的朋友圈列表
	MemberTopicUrl: "%s/micro-scrm/v1/moment/member-topic",
	// 发布朋友圈
	TopicCreateUrl: "%s/micro-scrm/v1/moment/topic-create",
	// 编辑朋友圈
	TopicEditUrl: "%s/micro-scrm/v1/moment/topic-edit",
	// 删除朋友圈
	TopicDeleteUrl: "%s/micro-scrm/v1/moment/topic-delete",
	// 朋友圈展示
	TopicDisplayUrl: "%s/micro-scrm/v1/moment/topic-display",
	// 朋友圈按内容列表展示
	ContentListUrl: "%s/micro-scrm/v1/moment/content-list",
	// 朋友圈列表按员工搜索
	MomentMemberListUrl: "%s/micro-scrm/v1/moment/member-list",
	// 获取朋友圈形象
	AppearanceGetUrl: "%s/micro-scrm/v1/moment/appearance-get",
	// 编辑朋友圈形象
	AppearanceEditUrl: "%s/micro-scrm/v1/moment/appearance-edit",
	// 朋友圈选择部门
	DepartmentMomentList: "%s/micro-scrm/v1/departments/moment-list",
	// 删除员工朋友圈
	TopicMemberDeleteUrl: "%s/micro-scrm/v1/moment/topic-member-delete",
	// 编辑员工朋友圈
	TopicMemberEditUrl: "%s/micro-scrm/v1/moment/topic-member-edit",

	// 选择项
	SelectorTextUrl:                "%s/micro-scrm/v1/marketing-strategies/selector-text",
	JsConfigUrl:                    "%s/micro-scrm/v1/moment/js-config",
	CustomerListUrl:                "%s/micro-scrm/v1/customer/user-list",
	EditCustomerRemarkUrl:          "%s/micro-scrm/v1/user/edit-remark",
	MarketingStrategyDataDetailUrl: "%s/micro-scrm/v1/marketing-strategies/data/detail",
	// 营销策略更新状态
	MarketingUpdateStatusUrl: "%s/micro-scrm/v1/marketing-strategies/update-status",
	// 营销策略触发
	MarketingStrategyTriggerUrl: "%s/micro-scrm/v1/marketing-strategies/market-trig",
	UpdateSettingDayUrl:         "%s/micro-scrm/v1/marketing-strategies/update-strategy-time",
	// 顾客列表标签下拉框列表
	CustomerTagListUrl:        "%s/micro-scrm/v1/tag/tag-list",
	MarketingUserGroupListUrl: "%s/micro-scrm/v1/marketing-strategies/tag-group-list",
	// 敏感词列表
	SensitiveListUrl: "%s/micro-scrm/v1/sensitive/list",
	// 敏感词保存
	SensitiveSaveUrl: "%s/micro-scrm/v1/sensitive/save",
	// 敏感词删除
	SensitiveDeleteUrl: "%s/micro-scrm/v1/sensitive/delete",
	// 敏感词详情
	SensitiveDetailUrl: "%s/micro-scrm/v1/sensitive/detail",
	// 敏感词明细列表
	SensitiveDetailListUrl: "%s/micro-scrm/v1/sensitive/detail-list",
	// 消息风控secret保存
	SettingsUpdateUrl:        "%s/micro-scrm/v1/settings/update",
	ChatMsgTypeList:          "%s/micro-scrm/v1/chat/msg-type-list",
	ChatAgreeInfo:            "%s/micro-scrm/v1/chat/agree-info",
	ChatMsgList:              "%s/micro-scrm/v1/chat/msg-list",
	ChatGroupMsgList:         "%s/micro-scrm/v1/chat/group-msg-list",
	ChatMsgListByID:          "%s/micro-scrm/v1/chat/msg-list-by-id",
	ChatMemberList:           "%s/micro-scrm/v1/chat/member-list",
	ChatObjectToMemberList:   "%s/micro-scrm/v1/chat/object/to-member-list",
	ChatObjectToExternalList: "%s/micro-scrm/v1/chat/object/to-external-list",
	ChatObjectToGroupList:    "%s/micro-scrm/v1/chat/object/to-group-list",
	GetUserInfoUrl:           "%s/micro-scrm/v1/chat/get-user-info",

	ConfigInfoUrl:   "%s/micro-scrm/v1/remind-config/config-info",
	ConfigSaveUrl:   "%s/micro-scrm/v1/remind-config/config-save",
	SettingsInfoUrl: "%s/micro-scrm/v1/settings/info",
	// 聊天超时记录-列表
	OvertimeListUrl:       "%s/micro-scrm/v1/overtime/list",
	CustomerListConfigUrl: "%s/micro-scrm/v1/user/list-config",
	TagGroupUserCountUrl:  "%s/micro-scrm/v1/marketing-strategies/tag-group-user-count",

	// 联系我渠道模块
	ChannelCreateUrl:  "%s/micro-scrm/v1/drainage/channels-create",          // 添加渠道
	ChannelListUrl:    "%s/micro-scrm/v1/drainage/channels-list",            // 渠道列表
	ChannelDetailUrl:  "%s/micro-scrm/v1/drainage/channels-detail",          // 渠道详情
	ChannelEditUrl:    "%s/micro-scrm/v1/drainage/channels-edit",            // 渠道编辑
	ContactListUrl:    "%s/micro-scrm/v1/drainage/contact/list",             // 联系我列表
	ContactCreateUrl:  "%s/micro-scrm/v1/drainage/contact/create",           // 联系我新增
	ContactUpdateUrl:  "%s/micro-scrm/v1/drainage/contact/update",           // 联系我修改
	ContactDeleteUrl:  "%s/micro-scrm/v1/drainage/contact/delete",           // 联系我删除
	ContactDetailsUrl: "%s/micro-scrm/v1/drainage/contact/details",          // 联系我详情
	DownloadQrCodeUrl: "%s/micro-scrm/v1/drainage/contact/download-qr-code", // 下载联系我二维码

	// 欢迎语模块
	// 新增欢迎语
	WelcomeCreateUrl: "%s/micro-scrm/v1/drainage/welcome-create",
	// 欢迎语详情
	WelcomeDetailUrl: "%s/micro-scrm/v1/drainage/welcome-detail",
	// 欢迎语默认详情
	WelcomeDefaultUrl: "%s/micro-scrm/v1/drainage/welcome-default",
	// 欢迎语编辑
	WelcomeEditUrl: "%s/micro-scrm/v1/drainage/welcome-edit",
	// 欢迎语删除
	WelcomeDeleteUrl: "%s/micro-scrm/v1/drainage/welcome-delete",
	// 欢迎语列表
	WelcomeListUrl: "%s/micro-scrm/v1/drainage/welcome-list",
	// 欢迎语渠道列表
	WelcomeValidateUrl: "%s/micro-scrm/v1/drainage/welcome-channel-list",
	// 缓存图片素材
	CacheMediaUrl: "%s/micro-scrm/v1/drainage/cache-media",
	// 获取所有外部条件
	ExternalWelcomeTypeUrl: "%s/micro-scrm/v1/drainage/all-type",

	// 渠道数据统计
	ChannelsDataUrl:        "%s/micro-scrm/v1/drainage/channels-data",
	OrderGoodsUrl:          "%s/micro-scrm/v1/customer/order-goods",
	BelongToChatsUrl:       "%s/micro-scrm/v1/customer/belong-to-chats",
	UserGroupDetailUrl:     "%s/micro-scrm/v1/user/group-detail",
	UserGroupExportDataUrl: "%s/micro-scrm/v1/user/group-export",
	UserGroupEditUrl:       "%s/micro-scrm/v1/user/group-edit",
	MemberLevelListURL:     "%s/micro-scrm/v1/shop-member/level-list",
	UserGroupDeleteUrl:     "%s/micro-scrm/v1/user/group-delete",

	// 渠道统计按天展示
	StatisticsAllByDayUrl: "%s/micro-scrm/v1/statistics/all-day",
	// 渠道统计总览
	StatisticsAllBySumUrl: "%s/micro-scrm/v1/statistics/all-sum",
	// 渠道统计渠道列表总览
	StatisticsAllByChannelUrl: "%s/micro-scrm/v1/statistics/channels-sum",
	// 渠道统计渠道按天展示
	StatisticsItemByDayUrl: "%s/micro-scrm/v1/statistics/channels-day",
	// 群发消息任务列表
	GroupMsgListUrl: "%s/micro-scrm/v1/group-msg/list",
	//新增群发消息
	GroupSendAddUrl: "%s/micro-scrm/v1/group-msg/group-send-add",
	//编辑群发消息
	GroupSendEditUrl: "%s/micro-scrm/v1/group-msg/group-send-edit",
	//任务明细-执行群主下拉列表
	TaskMemberListUrl: "%s/micro-scrm/v1/group-msg/member-list",
	// 群发任务明细
	TaskDetailUrl: "%s/micro-scrm/v1/group-msg/task-detail",
	// 群发任务详情
	TaskInfoUrl: "%s/micro-scrm/v1/group-msg/task-info",
	//微信优惠券模块
	// 可用的微信优惠券列表
	WxCouponAvailableUrl: "%s/micro-scrm/v1/coupon/wx-coupon-available-list",
	//微信渠道列表
	WxChannelUrl: "%s/micro-scrm/v1/marketing-strategies/get-wx-channel",
	//券是否在使用中
	GouponIsUsedUrl: "%s/micro-scrm/v1/marketing-strategies/coupon-is-used",

	MarketingDataCouponUrl: "%s/micro-scrm/v1/marketing-strategies/coupon-data/detail",

	CouponMatchUrl: "%s/micro-scrm/v1/marketing-strategies/coupon-match",
	//裂变活动
	ExportExchangeRecordUrl: "%s/micro-scrm/v1/fission/prize-exchange-record",
	// 裂变活动列表
	FissionActivityListUrl: "%s/micro-scrm/v1/fission/activity-list",
	// 裂变活动新增
	FissionActivityAddUrl: "%s/micro-scrm/v1/fission/activity-add",
	// 裂变活动编辑
	FissionActivityUpdateUrl: "%s/micro-scrm/v1/fission/activity-update",
	// 裂变活动删除
	FissionActivityDeleteUrl: "%s/micro-scrm/v1/fission/activity-delete",
	// 更新活动状态
	FissionUpdateStatusUrl: "%s/micro-scrm/v1/fission/update-status",
	// 裂变活动详情
	FissionActivityInfoUrl: "%s/micro-scrm/v1/fission/activity-info",
	// 选择社群列表
	ChatGroupListUrl: "%s/micro-scrm/v1/fission/chat-group-list",
	// 选择互斥活动列表
	MutualActivityListUrl: "%s/micro-scrm/v1/fission/mutual-activity-list",
	//裂变活动明细
	ActivityDetailUrl: "%s/micro-scrm/v1/fission/activity-detail",
	//裂变活动明细-列表
	ActivityDetailListUrl: "%s/micro-scrm/v1/fission/activity-detail-list",
	//发起人邀请明细
	InviteDetailUrl: "%s/micro-scrm/v1/fission/invite-detail",
	//发起人兑奖明细
	ExchangeRecordUrl: "%s/micro-scrm/v1/fission/exchange-record",
	//奖品兑换提醒
	OpenExchangeWarnUrl: "%s/micro-scrm/v1/fission/open-exchange-warn",
	//奖品兑换记录总数
	ExchangeTotalUrl: "%s/micro-scrm/v1/fission/exchange-total",
	//活动列表-全部有效活动
	ActivityListSimpleUrl: "%s/micro-scrm/v1/fission/activity-list-simple",
	//获取兑换提醒的额外员工
	ExchangeConfigMemberUrl: "%s/micro-scrm/v1/fission/get-exchange-config-member",
	GetSpuDetail:            "%s/micro-scrm/v1/marketing-strategies/get-spu-detail",
	ShortUrlAdd:             "%s/micro-scrm/v1/marketing-strategies/short-url-add",
	ShortUrlData:            "%s/micro-scrm/v1/marketing-strategies/short-url-data",
	GetShortUrlList:         "%s/micro-scrm/v1/marketing-strategies/short-url-list",
	ShortUrlDel:             "%s/micro-scrm/v1/marketing-strategies/short-url-del",
	GetTwoLevelIds:          "%s/micro-scrm/v1/marketing-strategies/two-level-ids",
	GetOriginId:             "%s/micro-scrm/v1/marketing-strategies/get-origin-id",

	SyncEterpriseUserToESUrl: "%s/micro-scrm/v1/customer/all-eterpriseuser-fores",
	GetAllShortUrlListUrl:    "%s/micro-scrm/v1/marketing-strategies/all-short-url-list",

	GetTemplateList:    "%s/micro-scrm/v1/marketing-strategies/template/list",
	CreateTemplateInfo: "%s/micro-scrm/v1/marketing-strategies/template/save",
	DelTemplateInfo:    "%s/micro-scrm/v1/marketing-strategies/template/delete",
	CheckTemplateInfo:  "%s/micro-scrm/v1/marketing-strategies/template/check",
	GetTemplateInfo:    "%s/micro-scrm/v1/marketing-strategies/template/detail",

	GetStyleConfigUrl:         "%s/micro-scrm/v1/fission/get-style-config",
	AddStyleConfigUrl:         "%s/micro-scrm/v1/fission/add-config",
	SetStyleConfigUrl:         "%s/micro-scrm/v1/fission/set-style-config",
	DelStyleConfigUrl:         "%s/micro-scrm/v1/fission/del-config",
	GetAllWelcomeConfigUrl:    "%s/micro-scrm/v1/drainage/welcome-all",
	GetAllKDConfigUrl:         "%s/micro-scrm/v1/drainage/get-detail",
	DelWelcomeConfigUrl:       "%s/micro-scrm/v1/drainage/welcome-del",
	WelcomeNewEditConfigUrl:   "%s/micro-scrm/v1/drainage/welcome-new-edit",
	WelcomeDetailNewConfigUrl: "%s/micro-scrm/v1/drainage/welcome-detailnew",
	SaveOpenOrderUrl:          "%s/micro-scrm/v1/open/save-order",

	WorkGroupChatListUrl:           "%s/micro-scrm/v1/chat-group/list",
	WorkGroupChatInfoUrl:           "%s/micro-scrm/v1/chat-group/info",
	WorkGroupChatTypesUrl:          "%s/micro-scrm/v1/chat-group/types",
	WorkGroupChatTypeManageUrl:     "%s/micro-scrm/v1/chat-group/type-manage",
	WorkGroupChatSettingUrl:        "%s/micro-scrm/v1/chat-group/setting",
	WorkGroupChatFiltersUrl:        "%s/micro-scrm/v1/chat-group/filters",
	WorkGroupChatMembersUrl:        "%s/micro-scrm/v1/chat-group/members",
	WorkGroupChatDataGridUrl:       "%s/micro-scrm/v1/chat-group/data-grid",
	WorkGroupChatBatchSettingUrl:   "%s/micro-scrm/v1/chat-group/batch-setting",
	WorkGroupChatUploadSettingUrl:  "%s/micro-scrm/v1/chat-group/upload-settings",
	WorkGroupChatDataUpdateTimeUrl: "%s/micro-scrm/v1/chat-group/data-update-time",
	WorkGroupChatBatchSyncUrl:      "%s/micro-scrm/v1/chat-group/sync",
	SimpleMembersOfDepartmentsUrl:  "%s/micro-scrm/v1/members/dp-members",

	DrainageSimpleChannelListUrl: "%s/micro-scrm/v1/drainage/channels-simple-list",

	TaskMemberDetailUrl:   "%s/micro-scrm/v1/group-msg/task-member-detail",
	TaskMemberSendListUrl: "%s/micro-scrm/v1/group-msg/task-member-send-list",
}

// 获取cdp请求链接
func GetUrl(uri string) string {
	return fmt.Sprintf(uri, config.CfgData.Service.MicroScrm.Domain)
}

var httpClient = sdk.NewHttpClient(sdk.ParseJson{})

type ScrmResponse struct {
	App       string      `json:"app"`
	Code      int         `json:"code"`
	Msg       string      `json:"message"`
	RequestId string      `json:"request_id"`
	Data      interface{} `json:"data"`
}

func (r *ScrmResponse) CheckCode() error {
	if r.Code == 200 {
		return nil
	}
	return errorx.NewErrorMsg(r.Code, r.Msg)

}

// CommonResponse 公共默认返回值
type CommonResponse struct {
	ScrmResponse
	Data interface{} `json:"data"`
}

var scrmClient = sdk.NewHttpClient(sdk.ParseJson{})
