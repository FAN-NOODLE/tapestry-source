package marketing

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/json"
	"gw-scrm/internal/pkg/sdk"
)

type UserGroupRequest struct {
	EnterpriseID int64 `json:"enterprise_id"`
}

//GroupList 用户标签列表
type GroupList struct {
	Aipl     []UserGroup `json:"aipl"`
	Customer []UserGroup `json:"custom"`
	Rfm      []UserGroup `json:"rfm"`
}

//UserGroup 用户标签
type UserGroup struct {
	GroupID   int64  `json:"group_id"`
	GroupName string `json:"group_name"`
}

//UserGroupResponse 用户标签返回
type UserGroupResponse struct {
	List []UserGroup `json:"list"`
}

//ActivitiesRequest 活动列表请求
type ActivitiesRequest struct {
	Type         int64 `json:"type"`
	Page         int64 `json:"page"`
	Size         int64 `json:"size"`
	EnterpriseID int64 `json:"enterprise_id"`
}

//LivesRequest 直播列表
type LivesRequest struct {
	EnterpriseID int64  `json:"enterprise_id"`
	Page         int64  `json:"page"`
	Size         int64  `json:"size"`
	SearchType   int8   `json:"search_type"`
	Search       string `json:"search"`
}

//Activities 活动列表
type Activities struct {
	ActivityType int64  `json:"activity_type"`
	Name         string `json:"name"`
	ID           int64  `json:"id"`
}

//ActivitiesResponse 活动列表返回数据
type ActivitiesResponse struct {
	List  []Activities `json:"list"`
	Page  int64        `json:"page"`
	Size  int64        `json:"size"`
	Total int64        `json:"total"`
}

//LivesResponse 直播列表返回
type LivesResponse struct {
	List  []Live `json:"list"`
	Page  int64  `json:"page"`
	Size  int64  `json:"size"`
	Total int64  `json:"total"`
}

//Live 直播信息
type Live struct {
	ID        int64  `json:"id"`
	LiveTitle string `json:"live_title"`
}

//ActivityPagesRequest 活动页请求
type ActivityPagesRequest struct {
	TemplateType string `json:"template_type"`
	Page         int64  `json:"page"`
	Size         int64  `json:"size"`
	EnterpriseID int64  `json:"enterprise_id"`
	SearchType   int8   `json:"search_type"`
	Search       string `json:"search"`
}

//ActivityPages 活动页
type ActivityPages struct {
	Name string `json:"name"`
	ID   int64  `json:"id"`
}

//ActivityPagesResponse 活动页返回
type ActivityPagesResponse struct {
	Data []ActivityPages `json:"data"`
}

//ActivityType 活动类型列表
type ActivityType struct {
	ActivityType int64  `json:"activity_type"`
	ActivityName string `json:"activity_name"`
}

//MarketingStrategy 策略信息
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
	StrategyPolicy  int8           `json:"strategy_policy" desc:"触达规则：1=通知一个渠道，2=通知所有渠道" db:"strategy_policy"`
	StrategyStatus  int8           `json:"strategy_status" desc:"状态: 1=未开始，2=投放中，3=暂停中，4=已结束"`
	Status          int8           `json:"status" desc:"状态:1=禁用,2=启用" db:"status"`
	IsFinished      int8           `json:"is_finished" desc:"是否完成；1：未完成，2：已完成"`
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
	Channels        []ChannelInfo  `json:"channels" desc:"渠道信息"`
	Relations       []RelationInfo `json:"relations" desc:"触发关联信息"`
	CreatedAt       int64          `json:"created_at" desc:"创建时间" db:"created_at"`
	UpdatedAt       int64          `json:"updated_at" desc:"更新时间" db:"updated_at"`
	DeletedAt       int64          `json:"deleted_at" desc:"删除时间" db:"deleted_at"`
	Page            int            `json:"page"`
	Size            int            `json:"size"`
	TriggerCount    int            `json:"trigger_count"`
}

//ChannelInfo 渠道信息
type ChannelInfo struct {
	ID            int64                  `json:"id" desc:"ID" db:"id"`
	StrategyID    int64                  `json:"strategy_id" desc:"策略ID" db:"strategy_id"`
	ChannelType   int8                   `json:"channel_type" desc:"渠道；1=短信，2=企微群发，3=自动发券，4=公众号" db:"channel_type"`
	Content       map[string]interface{} `json:"content" desc:"接受和返回内容"`
	TextContent   string                 `json:"text_content" desc:"数据表内容" db:"content"`
	ChannelDetail ChannelDetail          `json:"channel_detail"`
	JuniorID      int64                  `json:"junior_id" desc:"从属id" db:"junior_id"`
	Order         int64                  `json:"order" desc:"排序:默认1（此处正序排列）" db:"order"`
	CreatedAt     int64                  `json:"created_at" desc:"创建时间" db:"created_at"`
	UpdatedAt     int64                  `json:"updated_at" desc:"更新时间" db:"updated_at"`
	DeletedAt     int64                  `json:"deleted_at" desc:"删除时间" db:"deleted_at"`
}

// ChannelDetail 渠道详情
type ChannelDetail struct {
	Name string `json:"name"`
}

//RelationInfo 关联信息
type RelationInfo struct {
	ID             int64          `json:"id" desc:"ID" db:"id"`
	StrategyID     int64          `json:"strategy_id" desc:"策略ID" db:"strategy_id"`
	RelationID     int64          `json:"relation_id" desc:"关联id" db:"relation_id"`
	RelationDetail RelationDetail `json:"relation_detail"`
	CreatedAt      int64          `json:"created_at" desc:"创建时间" db:"created_at"`
	UpdatedAt      int64          `json:"updated_at" desc:"更新时间" db:"updated_at"`
	DeletedAt      int64          `json:"deleted_at" desc:"删除时间" db:"deleted_at"`
}

//RelationList 关联信息列表
type RelationDetail struct {
	ImgUrl string `json:"img_url"`
	Name   string `json:"name"`
}

//UpdateStatusRequest 更新策略状态
type UpdateStatusRequest struct {
	ID             int64 `json:"id"`
	EnterpriseID   int64 `json:"enterprise_id"`
	AdminID        int64 `json:"admin_id"`
	StrategyStatus int8  `json:"strategy_status" desc:"状态: 1=未开始，2=投放中，3=暂停中，4=已结束"`
	Status         int8  `json:"status" desc:"状态:1=禁用,2=启用"`
	UpdatedAt      int64 `json:"updated_at"`
}

//MarketingDetailRequest 营销策略详情
type MarketingDetailRequest struct {
	ID           int64 `json:"id"`
	Day          int   `json:"day"  db:"day"`
	BeginTime    int64 `json:"begin_time"  db:"begin_time"`
	EndTime      int64 `json:"end_time"  db:"end_time"`
	EnterpriseID int64 `json:"enterprise_id"`
}

//ListRequest 列表请求
type ListRequest struct {
	StrategyName    string `json:"strategy_name" desc:"策略名称" db:"strategy_name"`
	CustomerGroupID int64  `json:"customer_group_id" desc:"客户分组ID" db:"customer_group_id"`
	EnterpriseID    int64  `json:"enterprise_id" desc:"企业ID" db:"enterprise_id"`
	AdminID         int64  `json:"admin_id" desc:"操作人ID" db:"admin_id"`
	StrategyType    int8   `json:"strategy_type" desc:"策略类型: 1=客户策略，2=内置策略" db:"strategy_type"`
	StrategyStatus  int8   `json:"strategy_status" desc:"状态: 1=未开始，2=投放中，3=暂停中，4=已结束" db:"strategy_status"`
	StrategyMethod  int8   `json:"strategy_method" desc:"执行方式: 1=长期执行，2=立即执行，3=定时执行" db:"strategy_method"`
	StartTime       int64  `json:"start_time" desc:"执行开始时间" db:"start_time"`
	EndTime         int64  `json:"end_time" desc:"执行结束时间" db:"end_time"`
	StrategyPolicy  int8   `json:"strategy_policy" desc:"触达规则：1=通知一个渠道，2=通知所有渠道" db:"strategy_policy"`
	Channels        int8   `json:"channels" db:"channels"`
	Page            int    `json:"page"`
	Size            int    `json:"size"`
}

//DeleteRequest 策略删除请求
type DeleteRequest struct {
	ID           int64 `json:"id" db:"id"`
	EnterpriseID int64 `json:"enterprise_id"`
	DeletedAt    int64 `json:"deleted_at" db:"deleted_at"`
}

//营销策略数据参数
type GetMarketingDataDetailRequest struct {
	ID           int64 `json:"id"  binding:"required" db:"id"`
	Day          int   `json:"day"  db:"day"`
	BeginTime    int64 `json:"begin_time"  db:"begin_time"`
	EndTime      int64 `json:"end_time"  db:"end_time"`
	EnterpriseID int   `json:"enterprise_id"  db:"enterprise_id"`
}
type MarketingDataDetailRequest struct {
	ID           int64 `json:"id"  binding:"required" db:"id"`
	Day          int   `json:"day"  db:"day"`
	BeginTime    int64 `json:"begin_time"  db:"begin_time"`
	EndTime      int64 `json:"end_time"  db:"end_time"`
	EnterpriseID int64 `json:"enterprise_id"  db:"enterprise_id"`
}

type MarketingDataDetailNewRequest struct {
	ID           int64 `json:"id"  binding:"required" db:"id"`
	Day          int   `json:"day"  db:"day"`
	BeginTime    int64 `json:"begin_time"  db:"begin_time"`
	EndTime      int64 `json:"end_time"  db:"end_time"`
	EnterpriseID int64 `json:"enterprise_id"  db:"enterprise_id"`
}

type MarketingDataDetail struct {
	Type               int `json:"type" db:"type"`
	SendMemberNumber   int `json:"send_member_number" db:"send_member_number"`
	ReachMemberNumber  int `json:"reach_member_number" db:"reach_member_number"`
	FansNumber         int `json:"fans_number" db:"fans_number"`
	VisitorNumber      int `json:"visitor_number" db:"visitor_number"`
	AddOrderNumber     int `json:"add_order_number" db:"add_order_number"`
	PayOrder           int `json:"pay_order" db:"pay_order"`
	FinalAmount        int `json:"final_amount" db:"final_amount"`
	RefundAmount       int `json:"refund_amount" db:"refund_amount"`
	RefundAmountNumber int `json:"refund_amount_number" db:"refund_amount_number"`
}

type StrategyDataDetails struct {
	sdk.ResponseBase
	Data *StrategyData `json:"data"`
}

type StrategyFinlaDataDetails struct {
	Data  []*StrategyItem `json:"data"`
	Total *StrategyItem   `json:"total"`
}

type StrategyData struct {
	Data  []*StrategyItem `json:"data"`
	Total *StrategyItem   `json:"total"`
}

type StrategyItem struct {
	SendMemberNumber   int     `json:"send_member_number"`
	ReachMemberNumber  int     `json:"reach_member_number"`
	FansNumber         int     `json:"fans_number"`
	VisitorNumber      int     `json:"visitor_number"`
	AddOrderNumber     int     `json:"add_order_number"`
	PayOrder           int     `json:"pay_order"`
	FinalAmount        float64 `json:"final_amount"`
	RefundAmountNumber int     `json:"refund_amount_number"`
	RefundAmount       float64 `json:"refund_amount"`
	Types              int     `json:"type"`
}

// 内置营销触发数据参数
type MarketTrigRequest struct {
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" `
	RelationType int    `json:"relation_type" db:"relation_type" validate:"required" required:"触发类型不能为空"`
	RelationID   int    `json:"relation_id" db:"relation_id" validate:"required" required:"触发关联ID不能为空"`
	UserID       int64  `json:"user_id" validate:"required" required:"用户ID不能为空"`
	UserUnionID  string `json:"user_union_id" validate:"required" required:"用户UnionID不能为空"`
}

// 内置营销触发数据参数
type UserGroupList struct {
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id" `
}

// 内置营销触发数据参数
type TagGroupUserCountRequest struct {
	EnterpriseID    int64   `json:"enterprise_id" db:"enterprise_id" `
	TagOrGroupId    int64   `json:"tag_or_group_id" db:"tag_or_group_id" `
	ChannelTypeList []int64 `json:"channels"`
	StrategyId      int64   `json:"strategy_id" `
}

type RedisKeyTime struct {
	Day       int   `json:"day"`
	BeginTime int64 `json:"begin_time"`
	EndTime   int64 `json:"end_time"`
}

func (r RedisKeyTime) MarshalBinary() (data []byte, err error) {
	return json.Marshal(r)
}

type UpdateSettingDayRequest struct {
	ID           int64 `json:"id" binding:"required" db:"id"`
	SettingDay   int64 `json:"setting_day" binding:"required" db:"setting_day"`
	EnterpriseID int   `json:"enterprise_id"  db:"enterprise_id"`
}

// 获取渠道列表请求参数
type ChanListRequest struct {
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id" `
	Pid          int64 `json:"pid" db:"pid" `
}

// 获取短链请求参数
type ShortUrlRequest struct {
	ShortType int64  `json:"short_type" desc:"是否scrm生成短链1是" db:"short_type"`
	PageUrl   string `json:"longURL"`
}

type ShortUrlRes struct {
	ShortUrl string `json:"shortURL" desc:"短链" db:"short_url"`
}

type WxUrlRes struct {
	ErrCode  int64  `json:"errcode"`
	ErrMsg   string `json:"errmsg"`
	OpenLink string `json:"openlink"`
}

//获取小程序原始id
type AuthDetailRequest struct {
	EnterpriseID int64 `json:"enterprise_id"`
	AuthType     int   `json:"auth_type"`
}

//获取短链详情
type ShortUrlDataReq struct {
	Id int64 `json:"id" db:"id"`
}

type ShortUrlDelReq struct {
	Id int64 `json:"id" db:"id"`
}

//获取短链列表
type ShortUrlListReq struct {
	EnterpriseId int64  `json:"enterprise_id"`
	OneLevelId   int64  `json:"one_level_id"`
	TwoLevelId   int64  `json:"two_level_id"`
	PageName     string `json:"page_name"`
	UrlType      int64  `json:"url_type"`
	Operator     string `json:"operator"`
	CreateStart  int64  `json:"create_start"`
	CreateEnd    int64  `json:"create_end"`
	PageNum      int    `json:"page_num"`
	CurrentPage  int    `json:"current_page"`
}

//短链以及配置信息添加
type ShortUrlAddReq struct {
	EnterpriseId int64  `json:"enterprise_id"`
	Id           int64  `json:"id" db:"id"`
	PageUrl      string `json:"page_url"` //二级渠道page_url
	OneLevelId   int64  `json:"one_level_id"`
	Operator     string `json:"operator"`
	TwoLevelId   int64  `json:"two_level_id"`
	OneLevelName string `json:"one_level_name"` //一级渠道名称
	TwoLevelName string `json:"two_level_name"` //二级渠道名称
	ShortUrl     string `json:"short_url"`      //通用短链
	WxUrl        string `json:"wx_url"`         //生成的通用短链对应的微信快捷链接
	PageName     string `json:"page_name"`
	UrlType      int64  `json:"url_type"`
	MiddleData
}

type MiddleConfig struct {
	BrandImg          string `json:"brand_img"`           //外链中间页品牌图
	CardTitle         string `json:"card_title"`          //卡片标题
	CardDesc          string `json:"card_desc"`           //卡片活动说明
	BackgroundImg     string `json:"background_img"`      //外链中间页背景图
	IsHide            int64  `json:"is_hide"`             //是否显示按钮
	ButtonName        string `json:"button_name"`         //按钮名称
	ButtonColor       string `json:"button_color"`        //按钮颜色
	ButtonFontColor   string `json:"button_font_color"`   //按钮字体颜色
	LocationDesc      string `json:"location_desc"`       //跳转说明
	LocationFontColor string `json:"location_font_color"` //跳转说明文字颜色
	AdditionalType    int64  `json:"additional_type"`     //1商品2活动
}
type GoodsData struct {
	SpuOneLevelId int64   `json:"spu_one_level_id"`
	SpuTwoLevelId int64   `json:"spu_two_level_id"`
	SpuId         int64   `json:"spu_id"`
	SkuAttribute  string  `json:"sku_attribute"`
	SkuPrice      float32 `json:"sku_price"`
	SpuImg        string  `json:"spu_img"`
	SpuName       string  `json:"spu_name"`
	SpuPageUrl    string  `json:"spu_page_url"`
	WxSpuPageUrl  string  `json:"wx_spu_page_url"` //微信生成的可在手机跳转的链接

}

type ActiveData struct {
	ActiveImg        string `json:"active_img"`
	ActivePageUrl    string `json:"active_page_url"`
	WxActivePageUrl  string `json:"wx_active_page_url"` //微信生成的可在手机跳转的链接
	ActiveOneLevelId int64  `json:"active_one_level_id"`
	ActiveTwoLevelId int64  `json:"active_two_level_id"`
}

type ShortAddReq struct {
	ShortUrl   string `json:"short_url"`
	LongUrl    string `json:"long_url"`
	MiddleHtml string `json:"middle_html"`
	MiddleData
	ShortUrlType int64  `json:"short_url_type"`
	ReqType      int64  `json:"req_type"`
	EnterpriseId int64  `json:"enterprise_id"`
	OriginId     string `json:"origin_id"`
	PageUrl      string `json:"page_url"`
}

type MiddleData struct {
	MiddleConfig
	SpuData []*GoodsData `json:"spu_data"`
	ActiveData
}

//获取商品详情请求
type SpuDetailReq struct {
	PageUrl string `json:"page_url"`
}

type SpuDetailRes struct {
	SpuImg  string `json:"spu_img"`
	SpuName string `json:"spu_name"`
}
type AccessTokenRes struct {
	AccessToken string `json:"access_token"`
}
type ChanListRequestBody struct {
	Sign     string `json:"sign"`
	SignTime int64  `json:"sign_time"`
}

type ChanListCDPRes struct {
	Data []*ChanInfo `json:"data"`
}

type ChanListCDPResp struct {
	List []*ChanInfo `json:"list"`
}

type ChanInfo struct {
	ID         int64  `json:"id"`
	CustomerID int64  `json:"customer_id"`
	Name       string `json:"name"`
	Indice     string `json:"indice"`
	PageUrl    string `json:"page_url"`
	UrlType    int64  `json:"url_type"`
	IsHide     int64  `json:"is_hide"`
}

type (
	//自动发券，检测券是否已使用中
	CouponIsUsedRequest struct {
		WxCouponId int64 `json:"wx_coupon_id" db:"wx_coupon_id"`
	}
	CouponIsUsedResponse struct {
		WxCouponId int64 `json:"wx_coupon_id" db:"wx_coupon_id"`
		IsUsed     int   `json:"is_used"`
	}

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
		UnionId          string `json:"union_id"`
	}
)

//优惠券数据详情
type StrategyCouponDataDetails struct {
	sdk.ResponseBase
	Data *StrategyCouponData `json:"data"`
}

type StrategyCouponData struct {
	Data  *StrategyCouponDetail      `json:"list"`
	Total *StrategyCouponTotalDetail `json:"total"`
}

type StrategyCouponDetail struct {
	ChannelName string               `json:"channel_name"`
	SendNum     int                  `json:"send_num"`
	CouponList  []StrategyCouponItem `json:"coupon_list"`
}

type StrategyCouponTotalDetail struct {
	ChannelName string             `json:"channel_name"`
	SendNum     int                `json:"send_num"`
	CouponList  StrategyCouponItem `json:"coupon_list"`
}

type StrategyCouponItem struct {
	CouponName        string  `json:"coupon_name"`
	StockID           string  `json:"stock_id"`
	TouchNum          int     `json:"touch_num"`
	IsSendNum         int     `json:"is_send_num"`
	UseNum            int     `json:"use_num"`
	PayNum            int     `json:"pay_num"`
	CashFee           float64 `json:"cash_fee"`
	RefundCustomerNum int     `json:"refund_customer_num"`
	RefundFee         float64 `json:"refund_fee"`
}

type TemplateListReq struct {
	EnterpriseId int64  `json:"enterprise_id"`
	TemplateName string `json:"template_name"`
	Page         int    `json:"page" binding:"min=1" comment:"页码"`
	Size         int    `json:"size" binding:"min=10,max=99999999" comment:"分页数"`
}

type TemplateCreate struct {
	ID                 int      `json:"id"`
	EnterpriseId       int64    `json:"enterprise_id"`
	TemplateName       string   `json:"template_name"`
	TemplatePictureURL []string `json:"template_picture_url"`
	TemplateContent    string   `json:"template_content"`
}

type TemplateDel struct {
	ID                 int      `json:"id"`
	EnterpriseId       int64    `json:"enterprise_id"`
}
