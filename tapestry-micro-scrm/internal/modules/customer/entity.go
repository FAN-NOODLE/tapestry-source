package customer

import (
	"encoding/json"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/typex"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/sdk/work"
)

const (
	ExternalCallbackStatusInit    = 0 // 外部联系人回调处理状态：初始化
	ExternalCallbackStatusSendMq  = 1 // 外部联系人回调处理状态：入队列成功
	ExternalCallbackStatusProcess = 2 // 外部联系人回调处理状态：队列消费
	ExternalCallbackStatusSuccess = 3 // 外部联系人回调处理状态：处理成功
	ExternalCallbackStatusFailed  = 4 // 外部联系人回调处理状态：处理失败
)

const (
	CustomerTypeShop           = 1 //商城用户
	CustomerTypeExternal       = 2 //企业微信用户
	CustomerTypeOpen           = 3 //公众号用户
	CustomerTypeUnion          = 7
	CustomerTypeMobile         = 10
	CustomerTypeDeliveryMobile = 11
)

var sortKeyList = common.ListSortKeySlice{
	"pay_order_num",
	"pay_total",
	"refund_order_num",
	"refund_total",
	"per_customer_transaction",
	"latest_marketing_time",
}

// Callback 业务实体（结构体）
type Callback struct {
}

// ShopUserID 外部联系人
type ShopUserID struct {
	ShopUserID int64 `json:"shop_user_id" db:"shop_user_id"`
}

// RequestCustomer 客户请求
type RequestCustomer struct {
	MemberID       int64 `json:"member_id"`
	ExternalUserID string
	UserID         int64  `json:"user_id" db:"user_id" validate:"required" required:"客户ID必填"`
	EnterpriseID   int64  `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID必填"`
	Page           int    `json:"page"`
	Size           int    `json:"size"`
	SortKey        string `json:"sort_key"`
	SortType       string `json:"sort_type"`
	ChannelType    int8   `json:"channel_type"`
	ConsumerID     string `json:"consumerID"`
	TransType      string `json:"transType"`
}

// RequestCustomerCoupons 客户卡券列表请求
type RequestCustomerCoupons struct {
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID必填"`
	ConsumerID   string `json:"consumerID"`
	TransType    string `json:"transType"`
}

// RequestCustomerGroup 客户群请求
type RequestCustomerGroup struct {
	RequestCustomer
	Status int `json:"status" db:"status" desc:"群状态，1=入群 2=退群"`
}

// OrderGoodsRequest 订单商品接口请求参数
type OrderGoodsRequest struct {
	OrderID int64 `json:"order_id" validate:"required" required:"订单ID必填"`
	Page    int   `json:"page"`
	Size    int   `json:"size"`
}

type OrderGoodsResponse struct {
	Total int       `json:"total"`
	Page  int       `json:"page"`
	Size  int       `json:"size"`
	List  GoodsList `json:"list"`
}

// OrderSortKeyList 订单列表排序字段
var OrderSortKeyList = common.ListSortKeySlice{
	"cash_fee",
	"total_refund_fee",
	"pay_at",
	"refund_at",
	"final_amount",
}

// ThirdUserID 第三方用户标识
type ThirdUserID struct {
	ID      int64  `json:"id"`       // 内部用户标识
	ThirdID string `json:"third_id"` // 第三方标识
}

// EnterpriseUserBindDB 客户关联信息(DB)
type EnterpriseUserBindDB struct {
	ID         int64  `db:"id"`
	WxUnion    []byte `db:"wx_union"`
	WxOpen     []byte `db:"wx_open"`
	WxShop     []byte `db:"wx_shop"`
	WxExternal []byte `db:"wx_external"`
	BindTo     int64  `db:"bind_to"`
}

// ToUser 转换结构
func (u *EnterpriseUserBindDB) ToUser() *EnterpriseUserBind {
	user := &EnterpriseUserBind{ID: u.ID}
	json.Unmarshal(u.WxExternal, &user.WxExternal)
	if user.WxExternal == nil {
		user.WxExternal = make(map[string]*ThirdUserID)
	}
	json.Unmarshal(u.WxUnion, &user.WxUnion)
	if user.WxUnion == nil {
		user.WxUnion = make(map[string]string)
	}
	json.Unmarshal(u.WxShop, &user.WxShop)
	if user.WxShop == nil {
		user.WxShop = make(map[string]*ThirdUserID)
	}
	json.Unmarshal(u.WxOpen, &user.WxOpen)
	if user.WxOpen == nil {
		user.WxOpen = make(map[string]*ThirdUserID)
	}
	return user
}

// EnterpriseUserBind 客户关联信息
type EnterpriseUserBind struct {
	ID         int64
	WxUnion    map[string]string
	WxShop     map[string]*ThirdUserID
	WxOpen     map[string]*ThirdUserID
	WxExternal map[string]*ThirdUserID
}

// EnterpriseUserBase 客户基本信息
type EnterpriseUserBase struct {
	Name           string `db:"name" json:"name"`
	Nick           string `db:"nick" json:"nick"`
	Avatar         string `db:"avatar" json:"avatar"`
	Birthday       string `json:"birthday" db:"birthday"`
	UnionCode      string `json:"union_code" db:"union_code"`
	UnionId        string `json:"union_id" db:"union_id"`
	Point          string `json:"point" db:"point"`
	Grade          string `json:"grade" db:"grade"`
	Sex            int8   `db:"sex" json:"sex"`
	Mobile         string `db:"mobile" json:"mobile"`
	DeliveryMobile string `db:"delivery_mobile" json:"delivery_mobile"`
	StoreCode      string `json:"store_code" db:"store_code"`
	BACode         string `json:"ba_code" db:"ba_code"`
	//UserName       string `db:"user_name" json:"user_name"`
}

// EnterpriseUserLocation 客户地域信息
type EnterpriseUserLocation struct {
	Language string `json:"language" db:"language"`
	Country  string `db:"country" json:"country"`
	Province string `db:"province" json:"province"`
	City     string `db:"city" json:"city"`
}

// EnterpriseUserAttr 客户属性信息
type EnterpriseUserAttr struct {
	Tags           []int64 `json:"tags" db:"-"`
	Groups         []int64 `json:"groups" db:"-"`
	PayOrder       int     `json:"pay_order" db:"pay_order"`
	RefundOrder    int     `json:"refund_order" db:"refund_order"`
	TotalFee       int     `json:"total_fee" db:"total_fee"`
	RefundFee      int     `json:"refund_fee" db:"refund_fee"`
	LastMarketTime int64   `json:"latest_market_time" db:"latest_market_time"`
	MarketCount    int     `json:"market_count" db:"market_count"`
	LevelID        int64   `json:"level_id" db:"level_id"`
}

// AddTagOrGroup 增加分组标签
func (u *EnterpriseUserAttr) AddTagOrGroup(t string, tag int64) bool {
	if t != "group" {
		for _, v := range u.Tags {
			if v == tag {
				return false
			}
		}
		u.Tags = append(u.Tags, tag)
	} else {
		for _, v := range u.Groups {
			if v == tag {
				return false
			}
		}
		u.Groups = append(u.Groups, tag)
	}
	return true
}

// DelTagOrGroup 删除分组标签
func (u *EnterpriseUserAttr) DelTagOrGroup(t string, tag int64) bool {
	if t != "group" {
		for i, v := range u.Tags {
			if v == tag {
				u.Tags[i] = u.Tags[len(u.Tags)-1]
				u.Tags = u.Tags[0 : len(u.Tags)-1]
				return true
			}
		}
	} else {
		for i, v := range u.Groups {
			if v == tag {
				u.Groups[i] = u.Groups[len(u.Groups)-1]
				u.Groups = u.Groups[0 : len(u.Groups)-1]
				return true
			}
		}
	}
	return false
}

// EnterpriseUserSubscribe 用户的订阅信息
type EnterpriseUserSubscribe struct {
	SubscribeTime   int    `json:"subscribe_time" db:"subscribe_at"`
	SubscribeScene  string `json:"subscribe_scene" db:"subscribe_scene"`
	SubscribeStatus int    `json:"subscribe_status" db:"subscribe_status"`
}

// EnterpriseUserExternal 外部联系人信息
type EnterpriseUserExternal struct {
	Position        string `json:"position" db:"position" desc:"外部联系人职位"`
	CorpName        string `json:"corp_name" db:"corp_name" desc:"企业简称"`
	CorpFullName    string `json:"corp_full_name" db:"corp_full_name" desc:"企业全称"`
	Type            int8   `json:"type" db:"type" desc:"外部联系人的类型"`
	ExternalProfile string `json:"external_profile" db:"external_profile" desc:"外部联系人的自定义展示信息"`
}

// EnterpriseUserOther 用户的其他信息
type EnterpriseUserOther struct {
	QRScene    string `json:"qr_scene" db:"qr_scene"`
	QRSceneStr string `json:"qr_scene_str" db:"qr_scene_str"`
	Remark     string `json:"remark" db:"remark"`
}

// EnterpriseUser 企业用户
type EnterpriseUser struct {
	EnterpriseUserBind
	EnterpriseUserBase
	EnterpriseUserLocation
	EnterpriseUserAttr
	Remark   string
	DtUpdate string
}

type EnterpriseUserGroup struct {
	ID    int64   `db:"id"`
	Group []int64 `db:"groups"`
}

func (user *EnterpriseUser) GetOnlyShopUser() *ThirdUserID {
	for _, eu := range user.WxShop {
		return eu
	}
	return nil
}

func (user *EnterpriseUser) GetOnlyExternal() *ThirdUserID {
	for _, eu := range user.WxExternal {
		return eu
	}
	return nil
}

// GetShopUser 获取商城用户
func (user *EnterpriseUser) GetShopUser(appid string) *ThirdUserID {
	if len(user.WxShop) == 0 {
		return nil
	}

	return user.WxShop[appid]
}

// SourceUserBase 基础用户
type SourceUserBase struct {
	ID           int64  `db:"id"`
	AppID        string `db:"appid"`
	OpenID       string `db:"openid"`
	UnionID      string `db:"unionid"`
	EuID         int64  `db:"euid"`
	CustomerType int
}

// TagUserCount 标签人数
type TagUserCount struct {
	MemberCount                int64 `json:"member_count"`
	MobileUserCount            int64 `json:"mobile_user_count"`
	QwUserCount                int64 `json:"qw_user_count"`
	ShopUserCount              int64 `json:"shop_user_count"`
	WeChatOfficialAccountCount int64 `json:"chat_official_account_count"`
}

type TagUserCountAggs struct {
	ChatOfficialAccountCount struct {
		DocCount int64 `json:"doc_count"`
	} `json:"chat_official_account_count"`
	MemberCount struct {
		Value int64 `json:"value"`
	} `json:"member_count"`
	MobileUserCount struct {
		DocCount int64 `json:"doc_count"`
	} `json:"mobile_user_count"`
	QwUserCount struct {
		DocCount int64 `json:"doc_count"`
	} `json:"qw_user_count"`
	ShopUserCount struct {
		DocCount int64 `json:"doc_count"`
	} `json:"shop_user_count"`
}

// ResponseFootprintList 足迹列表返回结果
type ResponseFootprintList []struct {
	ID        int64  `json:"id" db:"id"`
	Event     string `json:"event" db:"event"`
	Object    string `json:"object" db:"object"`
	EventTxt  string `json:"event_txt" db:"event_txt"`
	EventDesc string `json:"event_desc" db:"event_desc"`
	CreatedAt int64  `json:"created_at" db:"created_at"`
}

// ResponseMemberList 成员列表返回结果
type ResponseMemberList []struct {
	MemberId               int64                 `json:"member_id" db:"member_id"`
	MemberUserid           string                `json:"member_userid" db:"member_userid"`
	MemberName             string                `json:"member_name" db:"member_name"`
	DepartmentList         MemberDepartsList     `json:"department_list"`
	CreateTime             int64                 `json:"create_time" db:"create_time" usage:"关系建立时间"`
	State                  string                `json:"state" json:"state" usage:"来源"`
	BelongToGroupChatCount int64                 `json:"belong_to_group_chat_count" usage:"当前客户占员工群总数的个数"`
	BelongToGroupChatList  BelongToGroupChatList `json:"belong_to_group_chat_list" usage:"当前客户占员工群列表"`
	Remark                 string                `json:"remark" db:"description" usage:"备注"`
	RemarkList             []string              `json:"remark_list" usage:"备注"`
}

// BelongToGroupChat 群聊
type BelongToGroupChat struct {
	GroupChatId   int64                 `json:"group_chat_id" db:"group_chat_id"`
	GroupChatName string                `json:"group_chat_name" db:"group_chat_name"`
	MemberId      int64                 `json:"member_id" db:"member_id"`
	MemberName    string                `json:"member_name" db:"member_name"`
	UserId        string                `json:"user_id" db:"user_id"`
	JoinTime      int64                 `json:"join_time" db:"join_time"`
	GroupMember   []WorkGroupChatMember `json:"group_member"`
}

// BelongToGroupChatList 群列表
type BelongToGroupChatList []BelongToGroupChat

// RequestOrderDetail 订单详情请求数据
type RequestOrderDetail struct {
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID必填"`
	OrderNo      int64 `json:"order_no" db:"order_no"`
}

// ResponseOrderDetail 订单详情返回结果
type ResponseOrderDetail struct {
	ID             int64          `json:"id" db:"id"`
	OrderNo        int64          `json:"order_no" db:"order_no"`
	Status         int            `json:"status" db:"status" usage:"订单状态"`
	PaymentChannel int            `json:"payment_channel" db:"payment_channel" usage:"支付渠道"`
	PayAt          int64          `json:"pay_at" db:"pay_at"`
	ThirdPaymentNo string         `json:"third_payment_no" db:"third_payment_no" usage:"第三方支付流水号"`
	GoodsPayPrice  int            `json:"goods_pay_price" db:"goods_pay_price"`
	LogisticsPrice int            `json:"logistics_price" db:"logistics_price" usage:"物流邮费"`
	DiscountFee    int            `json:"discount_fee" db:"discount_fee" usage:"优惠总金额"`
	CouponFee      int            `json:"coupon_fee" db:"coupon_fee" usage:"微信优惠金额"`
	OtherFee       int            `json:"other_fee" db:"other_fee" usage:"微信优惠金额"`
	TotalAmount    int            `json:"total_amount" db:"total_amount" usage:"原价"`
	CashFee        int            `json:"cash_fee" db:"cash_fee" usage:"实际支付金额"`
	TotalRefundFee int64          `json:"total_refund_fee" db:"total_refund_fee" usage:"退款金额"`
	OrderAt        int64          `json:"order_at" db:"order_at"`
	GoodsList      []GoodsEntity  `json:"goods_list" usage:"商品列表"`
	DeliveryInfo   DeliveryEntity `json:"delivery_info"`
}

// CustomerListRequest 客户列表接口请求参数结构
type ListRequest struct {
	ListCond
	Page     int    `json:"page"`
	Size     int    `json:"size"`
	SortKey  string `json:"sort_key"`
	SortType string `json:"sort_type"`
}

// ListResponse 客户列表接口返回数据结构
type ListResponse struct {
	List     []ListItem `json:"list"`
	Page     int        `json:"page"`
	Size     int        `json:"size"`
	Total    int64      `json:"total"`
	DtUpdate string     `json:"update"`
}

// getLevelIDList 获取列表数据中的等级id
func (l *ListResponse) getLevelIDList() []int64 {
	var result []int64
	for _, v := range l.List {
		if v.LevelID != 0 {
			result = append(result, v.LevelID)
		}
	}
	return result
}

// ListItem 客户列表客户数据结构
type ListItem struct {
	UserID                 int64                     `json:"user_id" db:"user_id"`
	Name                   string                    `json:"name" db:"name"`
	Nick                   string                    `json:"nick" db:"nick"`
	Avatar                 string                    `json:"avatar" db:"avatar"`
	Mobile                 string                    `json:"mobile" db:"mobile"`
	Tags                   []*TagItem                `json:"tags,omitempty"`
	CustomerGroups         []*TagItem                `json:"customer_groups,omitempty"`
	PayOrderNum            int                       `json:"pay_order_num" db:"pay_order_num"`
	PayTotal               int                       `json:"pay_total" db:"pay_total"`
	RefundOrderNum         int                       `json:"refund_order_num" db:"refund_order_num"`
	RefundTotal            int                       `json:"refund_total" db:"refund_total"`
	PerCustomerTransaction int                       `json:"per_customer_transaction" db:"per_customer_transaction"`
	Remark                 string                    `db:"remark"`
	RemarkList             []string                  `json:"remark_list"`
	LatestMarketingTime    int64                     `json:"latest_marketing_time" db:"latest_marketing_time"`
	Relations              []*WorkExternalUserRemark `json:"relations"`
	LevelName              string                    `json:"level_name" db:"level_name"`
	LevelID                int64                     `json:"level_id" db:"level_id"`
	Birthday               string                    `json:"birthday" db:"birthday"`
	UnionCode              string                    `json:"union_code" db:"union_code"`
	Point                  string                    `json:"point" db:"point"`
	Grade                  string                    `json:"grade" db:"grade"`
	UnionId                string                    `json:"union_id" db:"union_id"`
	//UserName               string                    `json:"user_name" db:"user_name"`
	UserSaInfo
}

// UserSaInfo 外部联系人SA信息
type UserSaInfo struct {
	MemberId       int64  `json:"member_id"`
	Description    string `json:"description"`
	DepartmentId   int64  `json:"department_id"`
	DepartmentName string `json:"department_name"`
	IsKeySa        int    `json:"is_key_sa" db:"is_key_sa"`   // 是否专属导购：0否，1是
	SaUserId       string `json:"sa_user_id" db:"sa_user_id"` // 是否专属导购：0否，1是
	SaUserName     string `json:"sa_user_name" db:"sa_user_name"`
}

type TagItem struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type WorkExternalUserRemark struct {
	WorkExternalUserID int64 `db:"work_external_user_id"`
	RelationID         int64 `db:"relation_id" json:"relation_id"`
	WorkExternalRelationForEdit
}

// WorkExternalRelationForEdit 外部联系人编辑结构体
type WorkExternalRelationForEdit struct {
	Remark         string   `db:"description"`
	MemberUserID   string   `json:"member_userid" db:"member_userid"`
	ExternalUserID string   `json:"external_userid" db:"external_userid"`
	MemberName     string   `db:"member_name" json:"member_name"`
	RemarkList     []string `json:"remarks"`
}

// ListCond 列表筛选条件
type ListCond struct {
	EnterpriseID    int64  `json:"enterprise_id"`
	NickName        string `json:"nick_name"`
	Nick            string `json:"nick"`
	Name            string `json:"name"`
	Mobile          string `json:"mobile"`
	MarketTimeStart int    `json:"market_time_start"`
	MarketTimeEnd   int    `json:"market_time_end"`
	UnionCode       int    `json:"union_code"`
	TagCond         struct {
		Include   bool    `json:"include"`
		TagIDList []int64 `json:"tag_id_list"`
	} `json:"tag_cond"`
	CustomerGroupCond struct {
		Include     bool    `json:"include"`
		GroupIDList []int64 `json:"tag_id_list"`
	} `json:"customer_group_cond"`
	MemberLevelCond struct {
		Include     bool     `json:"include"`
		LevelIDList []string `json:"level_id_list"`
	} `json:"member_level_cond"`
	Remark struct {
		Include bool   `json:"include"`
		Remark  string `json:"remark"`
	} `json:"remark"`
	CustomerType   int `json:"customer_type"`
	CustomerSource []int
	Eids           []int64 `json:"eids"`
}

type CondUnionID struct {
	EnterpriseID int64    `json:"enterprise_id"`
	UnionId      []string `json:"union_id"`
	EUid         []int64  `json:"_id"`
}

type TagsInfo struct {
	ID      int64  `json:"external_tag_id" db:"external_tag_id"`
	TagName string `json:"name" db:"name"`
}

// RequestDelivery 收货人信息请求参数
type RequestDelivery struct {
	ShopOrderID int64 `json:"shop_order_id" db:"shop_order_id"`
}

// DeliveryEntity 收货人信息返回参数
type DeliveryEntity struct {
	Name     string `json:"name" db:"name"`
	Mobile   string `json:"mobile" db:"mobile"`
	Province string `json:"province" db:"province"`
	City     string `json:"city" db:"city"`
	Area     string `json:"area" db:"area"`
	Address  string `json:"address" db:"address"`
}

// RequestOrderList 订单列表请求数据
type RequestOrderList struct {
	UserID       int64  `json:"user_id" db:"user_id" validate:"required" required:"客户ID必填"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" validate:"required" required:"企业ID必填"`
	OrderNo      int    `json:"order_no" db:"order_no"`
	Status       int    `json:"status" db:"status"`
	ShowCancel   int    `json:"show_cancel" usage:"是否展示未付款的订单 0=不展示（默认） 1=展示"`
	Page         int    `json:"page"`
	Size         int    `json:"size"`
	SortKey      string `json:"sort_key"`
	SortType     string `json:"sort_type"`
}

// ResponseOrderList 订单列表返回结果
type ResponseOrderList []struct {
	ID             int64                  `json:"id" db:"id"`
	OrderNo        string                 `json:"order_no" db:"order_no"`
	Status         int                    `json:"status" db:"status" usage:"订单状态"`
	FinalAmount    int64                  `json:"final_amount" db:"final_amount" usage:"含券支付金额"`
	CashFee        int64                  `json:"cash_fee" db:"cash_fee" usage:"实际支付金额"`
	TotalRefundFee int64                  `json:"total_refund_fee" db:"total_refund_fee" usage:"退款金额"`
	SellPrice      int                    `json:"sell_price" db:"sell_price"`
	GoodsPayPrice  int                    `json:"goods_pay_price" db:"goods_pay_price"`
	PayAt          int64                  `json:"pay_at" db:"pay_at"`
	RefundAt       int64                  `json:"refund_at" db:"refund_at"`
	GoodsList      GoodsList              `json:"goods_list" usage:"商品列表"`
	CouponList     CouponList             `json:"coupon_list" usage:"优惠券列表"`
	UserInfo       *OrderListCustomerInfo `json:"user_info"`
}

// OrderListCustomerInfo 订单列表用户信息
type OrderListCustomerInfo struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

// GoodsList 商品列表
type GoodsList []GoodsEntity

// GoodsEntity 商品实体
type GoodsEntity struct {
	ID               int64                    `json:"id" db:"id"`
	ShopOrderID      int64                    `json:"shop_order_id" db:"shop_order_id"`
	SpuName          string                   `json:"spu_name" db:"spu_name"`
	SkuPicture       string                   `json:"sku_picture" db:"sku_picture"`
	Quantity         int                      `json:"quantity" db:"quantity"`
	SellPrice        int                      `json:"sell_price" db:"sell_price"`
	SubPrice         int                      `json:"sub_price" db:"sub_price"`
	GoodsType        int                      `json:"goods_type" db:"goods_type"`
	SkuId            int64                    `json:"sku_id" db:"sku_id"`
	SpuId            int64                    `json:"spu_id" db:"spu_id"`
	SkuAttribute     string                   `json:"sku_attribute" db:"sku_attribute"`
	SkuAttributeList []map[string]interface{} `json:"sku_attribute_list"`
}

// CouponList 优惠券列表
type CouponList []struct {
	ShopOrderID int64  `json:"shop_order_id" db:"shop_order_id"`
	ID          string `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	StockId     string `json:"stock_id" db:"stock_id"`
}

// UserCategory 用户分类
type UserCategory struct {
	ID           typex.ID `json:"id" db:"id"`
	TagOrGroupId int64    `json:"tag_or_group_id" db:"tag_or_group_id"`
	EnterpriseId int64    `json:"enterprise_id" db:"enterprise_id"`
	Type         string   `json:"type" db:"type"`
	UnionId      string   `json:"union_id" db:"union_id"`
	CreatedAt    int64    `json:"created_at" db:"created_at"`
	UpdatedAt    int64    `json:"updated_at" db:"updated_at"`
	DeletedAt    int64    `json:"deleted_at" db:"deleted_at"`
}

// UserCategoryRequest 用户分类请求
type UserCategoryRequest struct {
	EnterpriseId int64  `json:"enterprise_id" db:"enterprise_id"`
	UnionId      string `json:"union_id" db:"union_id"`
	Type         string `json:"type" json:"type"`
}

// ResponseAddressList 收货人地址列表
type ResponseAddressList []struct {
	ID               int64  `json:"id" db:"id"`
	ConsigneeName    string `json:"consignee_name" db:"consignee_name" usage:"收货人姓名"`
	ConsigneeMobile  string `json:"consignee_mobile" db:"consignee_mobile" usage:"收货人手机号"`
	ConsigneeAddress string `json:"consignee_address" usage:"收货人地址"`
	Province         string `json:"province" json:"province"`
	City             string `json:"city" db:"city"`
	Area             string `json:"area" db:"area"`
	Address          string `json:"address" db:"address"`
}

// ResponseCustomerDetail 客户详情
type ResponseCustomerDetail struct {
	UserID       int64                    `json:"user_id" db:"user_id"`
	UnionID      string                   `json:"union_id" db:"union_id"`
	UserName     string                   `json:"user_name" db:"name"`
	Nickname     string                   `json:"nick_name" db:"nick"`
	Avatar       string                   `json:"avatar" db:"avatar"`
	Gender       int                      `json:"gender" db:"gender"`
	Country      string                   `json:"country" db:"country"`
	Province     string                   `json:"province" db:"province"`
	City         string                   `json:"city" db:"city"`
	Remark       string                   `json:"remark" db:"remark"`
	RemarkList   []string                 `json:"remark_list"`
	TagList      []string                 `json:"tag_list"`
	GroupList    []string                 `json:"group_list"`
	IsWorkUser   bool                     `json:"is_work_user"`
	OfficialList []*DetailOfficialAccount `json:"official_list"`
	ShowOfficial bool                     `json:"show_official"`
	LevelName    string                   `json:"level_name"`
	Birthday     string                   `json:"birthday" db:"birthday"`
	UnionCode    string                   `json:"union_code" db:"union_code"`
	Point        string                   `json:"point" db:"point"`
	Grade        string                   `json:"grade" db:"grade"`
	StoreCode    string                   `json:"store_code" db:"store_code"`
	BACode       string                   `json:"ba_code" db:"ba_code"`
	Mobile       string                   `json:"mobile" db:"mobile"`
}

// DetailOfficialAccount 用户公众号信息
type DetailOfficialAccount struct {
	UserID        int64  `json:"user_id" db:"user_id"`
	Name          string `json:"name" db:"name"`
	SubscribeTime int64  `json:"subscribe_time" db:"subscribe_time"`
}

// LossCustomersListRequest 流失客户列表接口请求参数结构
type LossCustomersListRequest struct {
	Type         int      `json:"type" binding:"required"`
	EnterpriseID int64    `json:"enterprise_id" db:"enterprise_id"`
	CorpId       string   `json:"corp_id" db:"corp_id"`
	CustomerName string   `json:"customer_name"`
	MemberID     int64    `json:"member_id"`
	DepartmentID int64    `json:"department_id"`
	ChannelId    int64    `json:"channel_id"`
	TimeSlotType int      `json:"time_slot_type"`
	TimeStart    int      `json:"time_start"`
	TimeEnd      int      `json:"time_end"`
	Page         int      `json:"page"`
	Size         int      `json:"size"`
	ChangeType   string   `json:"change_type"`
	UserId       []string `json:"user_id"`
}

// LossCustomersListResponse 流失客户列表接口返回数据结构
type LossCustomersListResponse struct {
	List  *[]LossCustomersListItem `json:"list"`
	Page  int                      `json:"page"`
	Size  int                      `json:"size"`
	Total int                      `json:"total"`
}

// MemberDeparts 成员组织
type MemberDeparts struct {
	MemberID       int64  `json:"member_id" db:"member_id"`
	DepartmentID   int64  `json:"department_id" db:"department_id"`
	DepartmentName string `json:"department_name" db:"department_name"`
}

// MemberDepartsList 成员组织列表
type MemberDepartsList []MemberDeparts

// LossCustomersListItem 流失客户列表客户数据结构
type LossCustomersListItem struct {
	ID                int64             `json:"id" db:"id"`
	EnterpriseID      int64             `json:"enterprise_id" db:"enterprise_id"`
	CustomerName      string            `json:"customer_name" db:"customer_name"`
	UserID            int64             `json:"user_id" db:"user_id"`
	WorkUserID        string            `json:"work_user_id" db:"work_user_id"`
	Remark            string            `json:"remark" db:"remark"`
	CustomerAvatar    string            `json:"customer_avatar" db:"customer_avatar"`
	MemberID          int64             `json:"member_id" db:"member_id"`
	MemberName        string            `json:"member_name" db:"member_name"`
	MemberAvatar      string            `json:"member_avatar" db:"member_avatar"`
	MemberDepartsList MemberDepartsList `json:"member_departs" usage:"部门信息"`
	CreatedAt         int64             `json:"created_at" db:"created_at"`
	RelieveAt         int64             `json:"relieve_at" db:"relieve_at"`
	CreateTime        int64             `json:"create_time" db:"createtime"`
}

// ExternalsRequest 外部联系人列表请求参数
type ExternalsRequest struct {
	Ids             []int64  `json:"ids"`
	UnionIds        []string `json:"union_ids"`
	ExternalUserIds []string `json:"external_user_ids"`
	Gender          int      `json:"gender"`
	Name            string   `json:"name"`
	Page            int      `json:"page"`
	Size            int      `json:"size"`
}

// WorkExternalUser 外部联系人信息
type WorkExternalUser struct {
	ID              int64  `json:"id" db:"id" desc:"外部联系人主键"`
	UserID          int64  `json:"user_id" db:"user_id" desc:"客户主表ID"`
	EnterpriseID    int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	UnionID         string `json:"union_id" db:"union_id" desc:"开放平台union_id"`
	ExternalUserID  string `json:"external_user_id" db:"external_userid" desc:"外部联系人UserID"`
	Name            string `json:"name" db:"name" desc:"外部联系人姓名"`
	Position        string `json:"position" db:"position" desc:"外部联系人职位"`
	Avatar          string `json:"avatar" db:"avatar" desc:"外部联系人头像"`
	CorpName        string `json:"corp_name" db:"corp_name" desc:"企业简称"`
	CorpFullName    string `json:"corp_full_name" db:"corp_full_name" desc:"企业全称"`
	Type            int8   `json:"type" db:"type" desc:"外部联系人的类型"`
	Gender          int8   `json:"gender" db:"gender" desc:"外部联系人性别"`
	ExternalProfile string `json:"external_profile" db:"external_profile" desc:"外部联系人的自定义展示信息"`
	CreatedAt       int64  `json:"created_at" db:"created_at" desc:"添加时间"`
	UpdatedAt       int64  `json:"updated_at" db:"updated_at" desc:"更新时间"`
}

// ReceiveCDPUserCategory cdp 用户分类
type ReceiveCDPUserCategory struct {
	ID           int64    `json:"id" db:"tag_or_group_id"`
	Type         string   `json:"type" db:"type"`
	EnterpriseID int64    `json:"enterprise_id" db:"enterprise_id"`
	Add          []string `json:"add"`
	Del          []string `json:"del"`
}

// TransferExternalUser 转换字段
func TransferExternalUser(response work.ExternalContactBase) *WorkExternalUser {
	return &WorkExternalUser{
		ExternalUserID:  response.ExternalUserID,
		UnionID:         response.UnionID,
		Name:            response.Name,
		Position:        response.Position,
		Avatar:          response.Avatar,
		CorpName:        response.CorpName,
		CorpFullName:    response.CorpFullName,
		Type:            response.Type,
		ExternalProfile: string(response.ExternalProfile),
		Gender:          response.Gender,
		CreatedAt:       time.Now().Unix(),
		UpdatedAt:       time.Now().Unix(),
	}
}

// ExternalUserRelation 获取联系人
type ExternalUserRelation struct {
	ID                 int64 `json:"id" db:"id"`
	WorkExternalUserId int64 `json:"work_external_user_id" db:"work_external_user_id"`
	State              int64 `json:"state" db:"state"`
	EnterpriseId       int64 `json:"enterprise_id" db:"enterprise_id"`
}

// ListChannelByQrCodeIdsResponse 根据二维码ID 查询渠道信息返回参数
type ListChannelByQrCodeIdsResponse struct {
	ChannelID    int64  `json:"channel_id"`
	QrCodeID     int64  `json:"qr_code_id"`
	ChannelName  string `json:"channel_name"`
	EnterpriseID int64  `json:"enterprise_id"`
}

// ExistsRelationByTimeReq 判断在该时间之前是否存在关系
type ExistsRelationByTimeReq struct {
	CustomerID   int64   `json:"customer_id" db:"customer_id"`
	EnterpriseID int64   `json:"enterprise_id" db:"enterprise_id"`
	MembersID    []int64 `json:"members_id" db:"members_id"`
	StartTime    int64   `json:"start_time" db:"start_time"`
}

type EnterpriseUserEntity struct {
	ID                     int64  `json:"id" db:"id"`
	Name                   string `db:"name" json:"name"`
	Nick                   string `db:"nick" json:"nick"`
	Avatar                 string `db:"avatar" json:"avatar"`
	Birthday               string `json:"birthday" db:"birthday"`
	Sex                    int8   `db:"sex" json:"sex"`
	Mobile                 string `db:"mobile" json:"mobile"`
	DeliveryMobile         string `json:"delivery_mobile" db:"delivery_mobile"`
	Language               string `json:"language" db:"language"`
	Country                string `db:"country" json:"country"`
	Province               string `db:"province" json:"province"`
	City                   string `db:"city" json:"city"`
	Remark                 []byte `json:"remark" db:"remark"`
	WxUnion                []byte `json:"wx_union" db:"wx_union"`
	WxShop                 []byte `json:"wx_shop" db:"wx_shop"`
	WxOpen                 []byte `json:"wx_open" db:"wx_open"`
	WxExternal             []byte `json:"wx_external" db:"wx_external"`
	Tags                   []byte `json:"tags" db:"tags"`
	Groups                 []byte `json:"groups" db:"groups"`
	Version                int    `json:"version" db:"version"`
	PayOrder               int    `json:"pay_order" db:"pay_order"`
	TotalFee               int    `json:"total_fee" db:"total_fee"`
	RefundOrder            int    `json:"refund_order" db:"refund_order"`
	RefundFee              int    `json:"refund_fee" db:"refund_fee"`
	PerCustomerTransaction int    `json:"per_customer_transaction" db:"per_customer_transaction"`
	LatestMarketTime       int    `json:"latest_market_time" db:"latest_market_time"`
	MarketCount            int    `json:"market_count" db:"market_count"`
	BindTo                 int64  `json:"bind_to" db:"bind_to"`
	LevelId                int64  `json:"level_id" db:"level_id"`
	CreatedAt              int    `json:"created_at" db:"created_at"`
	UpdatedAt              int    `json:"updated_at" db:"updated_at"`
}

type SyncEterpriseUserToEsRequest struct {
	EnterpriseId int64  `json:"enterprise_id"`
	BeginTime    string `json:"begin_time"`
	EndTime      string `json:"end_time"`
}

type EsQueryRes struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total    int64   `json:"total"`
		MaxScore float64 `json:"max_score"`
		Hits     []struct {
			Index  string  `json:"_index"`
			Type   string  `json:"_type"`
			ID     string  `json:"_id"`
			Score  float64 `json:"_score"`
			Source struct {
				Avatar                 string        `json:"avatar"`
				BindTo                 int           `json:"bind_to"`
				Birthday               string        `json:"birthday"`
				UnionCode              string        `json:"union_code"`
				Point                  string        `json:"point"`
				Grade                  string        `json:"grade"`
				City                   string        `json:"city"`
				Country                string        `json:"country"`
				CreatedAt              int           `json:"created_at"`
				Groups                 []interface{} `json:"groups"`
				ID                     int64         `json:"id"`
				Language               string        `json:"language"`
				LatestMarketTime       int           `json:"latest_market_time"`
				MarketCount            int           `json:"market_count"`
				Mobile                 string        `json:"mobile"`
				DeliveryMobile         string        `json:"delivery_mobile"`
				Name                   string        `json:"name"`
				Nick                   string        `json:"nick"`
				PayOrder               int           `json:"pay_order"`
				PerCustomerTransaction int           `json:"per_customer_transaction"`
				Province               string        `json:"province"`
				RefundFee              int           `json:"refund_fee"`
				RefundOrder            int           `json:"refund_order"`
				Remark                 string        `json:"remark"`
				Sex                    int8          `json:"sex"`
				Tags                   []interface{} `json:"tags"`
				TotalFee               int           `json:"total_fee"`
				UpdatedAt              int           `json:"updated_at"`
				Version                int           `json:"version"`
				LevelId                int64         `json:"level_id"`
				DtUpdate               string        `json:"dt_update"`
				WxExternal             []struct {
					AppID     string       `json:"app_id"`
					ThirdUser *ThirdUserID `json:"third_user"`
				} `json:"wx_external"`
				WxOpen []struct {
					AppID     string       `json:"app_id"`
					ThirdUser *ThirdUserID `json:"third_user"`
				} `json:"wx_open"`
				WxShop []struct {
					AppID     string       `json:"app_id"`
					ThirdUser *ThirdUserID `json:"third_user"`
				} `json:"wx_shop"`
				WxUnion []struct {
					Pid     string `json:"pid"`
					UnionID string `json:"union_id"`
				} `json:"wx_union"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

// LostRecordRequest 丢失客户记录请求参数
type LostRecordRequest struct {
	EnterpriseId int64  `json:"enterprise_id" binding:"required" comment:"企业Id"`
	Mobile       string `json:"mobile" binding:"required" comment:"员工手机号"`
	UserId       string `json:"user_id" comment:"员工Id"`
	DpId         int64  `json:"dp_id" comment:"部门Id"`
	Page         int    `json:"page" binding:"min=1" comment:"页码"`
	Size         int    `json:"size" binding:"min=1,max=100" comment:"分页数"`
	Time         int64  `json:"time"  binding:"required" comment:"时间戳"`
	DeleteType   uint   `json:"delete_type" binding:"min=1,max=2" comment:"类型"`
}

type LossMember struct {
	ID     int64  `json:"id" db:"id" desc:"员工ID"`
	Name   string `json:"name" db:"name" desc:"员工姓名"`
	Avatar string `json:"avatar" db:"avatar" desc:"头像"`
}

type EUidCount struct {
	Count int64 `json:"count" db:"count"`
}
type EUidList struct {
	EUid int64 `json:"euid" db:"euid"`
}

type UserInfoRequest struct {
	UnionId string `json:"union_id"`
}
