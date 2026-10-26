package common

import "encoding/json"

var ShareSource = "scrm"

// 系统执行默认的管理员ID
const SystemAdminID = -1

// 系统管理员默认名称
const SystemAdminName = "scrm-system"

// ContextKeyEnum ctx.Context中存储的数据key枚举
var ContextKeyEnum = struct {
	TokenInfo      string
	EnterpriseHash string
	EnterpriseID   string
}{
	TokenInfo:      "token_info",
	EnterpriseHash: "enterprise_hash",
	EnterpriseID:   "enterprise_id",
}

// HeaderKeyEnum 请求Header中存储的key枚举
var HeaderKeyEnum = struct {
	AuthToken      string
	EnterpriseHash string
	EnterpriseID   string
	XRequestID     string
}{
	AuthToken:      "Auth-Token",
	EnterpriseHash: "Enterprise-Hash",
	EnterpriseID:   "Enterprise-ID",
	XRequestID:     "X-Request-ID",
}

// EnterpriseInfo 企业信息结构体
type EnterpriseInfo struct {
	ID             int64  `json:"id"`
	EnterpriseID   int64  `json:"enterprise_id"`
	Hash           string `json:"hash"`
	Name           string `json:"name"`
	Logo           string `json:"logo"`
	IndustryID     int64  `json:"industry_id"`
	CorpID         string `json:"corp_id"`
	Domain         string `json:"domain"`
	Status         int    `json:"status"`
	BusinessSource int    `json:"business_source"`
	Version        int    `json:"version"`
	VersionTag     string `json:"version_tag"`
	IsSplitFlow    int    `json:"is_split_flow"`
	Identifier     string `json:"identifier"`
}

// MarshalBinary 将结构体转码为字符串
func (e *EnterpriseInfo) MarshalBinary() (data []byte, err error) {
	return json.Marshal(e)
}

// UnmarshalBinary 解码将字符串转为结构体
func (e *EnterpriseInfo) UnmarshalBinary(data []byte) error {
	err := json.Unmarshal(data, e)
	if err != nil {
		return err
	}
	return nil
}

// Event 事件类型
type Event string

// String 字符串格式化
func (e Event) String() string {
	return string(e)
}

// EventEnum 事件枚举
var EventEnum = struct {
	FirstVisitShop             Event
	AddExternalUser            Event
	DeleteExternalUser         Event
	DeleteExternalUserByMember Event
	JoinWorkChatGroup          Event
	LeaveWorkChatGroup         Event
	PayOrder                   Event
	UserJoinActivity           Event
	SubscribeOfficial          Event
	UnsubscribeOfficial        Event
	ViewMoment                 Event
	UserJoinFissionActivity    Event
}{
	FirstVisitShop:             "first_visit_shop",
	AddExternalUser:            "add_external_user",
	DeleteExternalUser:         "delete_external_user",
	DeleteExternalUserByMember: "delete_external_user_by_member",
	JoinWorkChatGroup:          "join_work_chat_group",
	LeaveWorkChatGroup:         "leave_work_chat_group",
	PayOrder:                   "pay_order",
	UserJoinActivity:           "shop_user_join_activity",
	SubscribeOfficial:          "subscribe_official",
	UnsubscribeOfficial:        "unsubscribe_official",
	ViewMoment:                 "view_moment",
	UserJoinFissionActivity:    "user_join_fission_activity",
}

const DDCMethodHeader = "method"

type DdcEvent string

// String 格式化
func (d DdcEvent) String() string {
	return string(d)
}

// DdcEventEnum ddc事件枚举
var DdcEventEnum = struct {
	UserAdd                   DdcEvent
	UserUpdate                DdcEvent
	UserAuth                  DdcEvent
	UserAddressAdd            DdcEvent
	UserAddressUpdate         DdcEvent
	UserAddressDelete         DdcEvent
	UserJoinActivity          DdcEvent
	OrderCreate               DdcEvent
	OrderPaySuccess           DdcEvent
	OrderRefund               DdcEvent
	OrderRefundPartial        DdcEvent
	OrderRefundPartialOverdue DdcEvent
	GZHUserSubscribe          DdcEvent
	GZHUserUnsubscribe        DdcEvent
	GZHUserSubscribeScene     DdcEvent
	GZHUserScene              DdcEvent
	ActivityAdd               DdcEvent
	ActivityUpdate            DdcEvent
	ActivityDelete            DdcEvent
	WxCouponAdd               DdcEvent
	WxCouponUpdate            DdcEvent
	WxCouponDelete            DdcEvent
	CouponUserAdd             DdcEvent
	CouponUserUpdate          DdcEvent
	UserBindingGuideAdd       DdcEvent
	UserBindingGuideUpdate    DdcEvent
	UserBindingGuideDelete    DdcEvent
	GoodsBase                 DdcEvent
	GoodsChange               DdcEvent
	GoodsDelete               DdcEvent
	ShopMemberAdd             DdcEvent
	ShopMemberUpdate          DdcEvent
	ShopMemberLevelAdd        DdcEvent
	ShopMemberLevelUpdate     DdcEvent
}{
	UserAdd:                   "user.shop_user_add",
	UserUpdate:                "user.shop_user_update",
	UserAuth:                  "user.shop_user_auth",
	UserAddressAdd:            "user.user_address_add",
	UserAddressUpdate:         "user.user_address_update",
	UserAddressDelete:         "user.user_address_delete",
	UserJoinActivity:          "user.shop_user_join_activity",
	OrderCreate:               "order.create_order",
	OrderPaySuccess:           "order.wx_pay_success_no_confirm",
	OrderRefund:               "order.all_refund_success",
	OrderRefundPartial:        "order.partial_refund_success",
	OrderRefundPartialOverdue: "order.partial_refund_success_overdue_after_sale",
	GZHUserSubscribe:          "gzh.subscribe",
	GZHUserUnsubscribe:        "gzh.unsubscribe",
	GZHUserSubscribeScene:     "gzh.subscribe_scene",
	GZHUserScene:              "gzh.scene",
	ActivityAdd:               "activity.shop_activity_add",
	ActivityUpdate:            "activity.shop_activity_update",
	ActivityDelete:            "activity.shop_activity_delete",
	WxCouponAdd:               "coupon.wx_coupon_add",
	WxCouponUpdate:            "coupon.wx_coupon_update",
	WxCouponDelete:            "coupon.wx_coupon_delete",
	CouponUserAdd:             "coupon.user_coupon_add",
	CouponUserUpdate:          "coupon.user_coupon_update",
	UserBindingGuideAdd:       "guide.user_binding_add",
	UserBindingGuideUpdate:    "guide.user_binding_update",
	UserBindingGuideDelete:    "guide.user_binding_delete",
	GoodsBase:                 "goods.base",
	GoodsChange:               "goods.change",
	GoodsDelete:               "goods.delete",
	ShopMemberAdd:             "member.shop_member_add",
	ShopMemberUpdate:          "member.shop_member_update",
	ShopMemberLevelAdd:        "member.shop_member_level_add",
	ShopMemberLevelUpdate:     "member.shop_member_level_update",
}

// PageType 分页列表通用
type PageType struct {
	Page  int64 `json:"page" desc:"页数"`
	Size  int64 `json:"size" desc:"每页记录数"`
	Total int64 `json:"total" desc:"总记录数"`
}

// ListSortKeySlice 客户列表排序字段数组
type ListSortKeySlice []string

func (l ListSortKeySlice) Contains(key string) bool {
	for _, v := range l {
		if key == v {
			return true
		}
	}
	return false
}

type Params map[string]string

type Inner struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type requestParams struct {
	Args []Inner `json:"args"`
}

type ResponseParams struct {
	ID int64 `json:"id"`
}

// SearchListCount 列表查询总数实体
type SearchListCount struct {
	Total int64 `json:"total"`
}

// Text 文本
type Text struct {
	Content string `json:"content"`
}

// Image 图片
type Image struct {
	MediaId string `json:"media_id"`
	PicUrl  string `json:"pic_url"`
}

// Link 链接
type Link struct {
	Title    string `json:"title"`
	WxPicUrl string `json:"picurl"`
	Desc     string `json:"desc"`
	Url      string `json:"url"`
	PicUrl   string `json:"pic_url"`
}

// MiniProgram 小程序
type MiniProgram struct {
	Title      string `json:"title"`
	PicMediaId string `json:"pic_media_id"`
	Appid      string `json:"appid"`
	Page       string `json:"page"`
	CoverUrl   string `json:"cover_url"`
}

type RequestPage struct {
	Page int64 `json:"page"`
	Size int64 `json:"size"`
}

type RequestID struct {
	ID int64 `json:"id"`
}
