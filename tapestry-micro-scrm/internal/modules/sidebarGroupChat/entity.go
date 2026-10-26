package sidebarGroupChat

import "micro-scrm/internal/pkg/common"

//const (
//	BUSINESS_USER_JOIN_FISSION_EVENT = "user_join_fission_event"
//)

// GroupChatInfoRequest 社群详情请求参数
type GroupChatInfoRequest struct {
	EnterpriseId int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业Id"`
	GroupChatId  string `json:"group_chat_id" db:"group_chat_id" desc:"社群Id"`
	Type         int64  `json:"type" desc:"排序方式（1：默认成员列表，2：活跃排行榜，3：成交排行榜）"`
	Page         int64  `json:"page"`
	Size         int64  `json:"size"`
}

// WorkGroupChat 社群表返回数据
type WorkGroupChat struct {
	ID            int64  `json:"id" db:"id"`
	EnterpriseID  int64  `json:"enterprise_id" db:"enterprise_id"`
	ChatID        string `json:"chat_id" db:"chat_id"`
	Name          string `json:"name" db:"name"`
	Type          int64  `json:"type" db:"type"`
	OwnerMemberID int64  `json:"owner_member_id" db:"owner_member_id"`
	Notice        string `json:"notice" db:"notice"`
	CreateTime    int64  `json:"create_time" db:"create_time"`
	Status        int64  `json:"status" db:"status"`
	DissolveTime  int64  `json:"dissolve_time" db:"dissolve_time"`
	CreatedAt     int64  `json:"created_at" db:"created_at"`
	UpdatedAt     int64  `json:"updated_at" db:"updated_at"`
	DeletedAt     int64  `json:"deleted_at" db:"deleted_at"`
}

// WorkGroupChat 社群表返回数据
type WorkGroupChatMember struct {
	ChatId                string `json:"chat_id" db:"chat_id"`
	GroupChatId           int64  `json:"group_chat_id" db:"group_chat_id"`
	Unionid               string `json:"unionid" db:"unionid"`
	ExternalUserId        string `json:"external_userid" db:"external_userid"`
	JoinTime              int64  `json:"join_time" db:"join_time"`
	LeaveTime             int64  `json:"leave_time" db:"leave_time"`
	JoinScene             int64  `json:"join_scene" db:"join_scene"`
	ShopUserId            int64  `json:"shop_user_id" db:"shop_user_id"`
	NickName              string `json:"nick_name" db:"nick_name"`
	Avatar                string `json:"avatar" db:"avatar"`
	Sex                   int64  `json:"sex" db:"sex"`
	Status                int64  `json:"status" db:"status"`
	MessageTotal          int64  `json:"message_total" db:"message_total"`
	TodayMessageTotal     int64  `json:"today_message_total" db:"today_message_total"`
	TodayDealOrderTotal   int64  `json:"today_deal_order_total" db:"today_deal_order_total"`
	TodayDealOrderAmount  int64  `json:"today_deal_order_total" db:"today_deal_order_total"`
	TodayJoinTotal        int64  `json:"today_join_total" db:"today_join_total"`
	TodayLeaveTotal       int64  `json:"today_leave_total" db:"today_leave_total"`
	TodayRefundOrderTotal int64  `json:"today_refund_order_total" db:"today_refund_order_total"`
	TodayRefundAmount     int64  `json:"today_refund_amount" db:"today_refund_amount"`
	DealOrderTotal        int64  `json:"deal_order_total" db:"deal_order_total"`
	DealOrderAmount       int64  `json:"deal_order_amount" db:"deal_order_amount"`
	RefundOrderTotal      int64  `json:"refund_order_total" db:"refund_order_total"`
	RefundAmount          int64  `json:"refund_amount" db:"refund_amount"`
	LastChatAt            int64  `json:"last_chat_at" db:"last_chat_at"`
	JoinDataTime          string `json:"join_data_time" db:"join_data_time"`
	LeaveDataTime         string `json:"leave_data_time" db:"leave_data_time"`
	PayDataAt             string `json:"pay_data_at" db:"pay_data_at"`
	PayAt                 int64  `json:"pay_at" db:"pay_at"`
	JoinNum               int64  `json:"join_num" db:"join_num"`
}

// WorkGroupChat 社群表返回数据
type WorkGroupCorpUser struct {
	Unionid   string `json:"unionid" db:"unionid"`
	UnionCode string `json:"union_code" db:"union_code"`
	Grade     string `json:"grade" db:"grade"`
}

// ChatMsgRequest 获取群聊天记录请求
type ChatMsgRequest struct {
	RoomId   string `json:"room_id" db:"room_id"`
	FromUser string `json:"from_user" db:"from_user"`
}

// ChatMsg 聊天数据返回结构体
type ChatMsgResponse struct {
	ID           int64  `json:"id" db:"id"`
	MsgID        string `json:"msg_id" db:"msg_id"`
	Action       string `json:"action" db:"action"`
	FromUser     string `json:"from_user" db:"from_user"`
	FromUserType int64  `json:"from_user_type" db:"from_user_type"`
	ToList       string `json:"to_list" db:"to_list"`
	RoomID       string `json:"room_id" db:"room_id"`
	MsgAt        int64  `json:"msg_at" db:"msg_at"`
	MsgTime      int64  `json:"msg_time" db:"msg_time"`
	Type         string `json:"type" db:"type"`
	Content      string `json:"content" db:"content"`
	Seq          int64  `json:"seq" db:"seq"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id"`
	CreatedAt    int64  `json:"created_at" db:"created_at"`
	UpdatedAt    int64  `json:"updated_at" db:"updated_at"`
	DeletedAt    int64  `json:"deleted_at" db:"deleted_at"`
	SendTime     string `json:"send_time" db:"send_time"`
	Time         string `json:"time" db:"time"`
}

// GroupChatInfoResponse 社群详情返回参数
type GroupChatInfoResponse struct {
	Name                 string         `json:"name" desc:"群名称"`
	CreateTime           int64          `json:"create_time" desc:"建群时间"`
	GroupLeaderName      string         `json:"group_leader_name" desc:"群主名称"`
	GroupImageUrl        string         `json:"group_image_url" desc:"群头像"`
	MaleTotal            int64          `json:"male_total" desc:"男性成员总数"`
	GenderFemaleTotal    int64          `json:"gender_female_total" desc:"女性成员数量"`
	UnknownGenderTotal   int64          `json:"unknown_gender_total" desc:"未知性别成功数量"`
	TodayActiveUser      int64          `json:"today_active_user" desc:"今日活跃用户总数"`
	TodayDealUser        int64          `json:"today_deal_user" desc:"今日成交用户总数"`
	TodayDealOrder       int64          `json:"today_deal_order" desc:"今日成交订单总数	"`
	TodayDealOrderAmount int64          `json:"today_deal_order_amount" desc:"今日实收金额（分）"`
	TodayJoinUser        int64          `json:"today_join_user" desc:"今日入群人数"`
	TodayLeaveUser       int64          `json:"today_leave_user" desc:"今日退群人数"`
	TodayRefundOrder     int64          `json:"today_refund_order" desc:"今日退款订单"`
	TodayRefundAmount    int64          `json:"today_refund_amount" desc:"今日退款金额"`
	GroupMemberNum       int64          `json:"group_member_num" desc:"群成员数量"`
	GroupVipMemberNum    int            `json:"group_vip_member_num" desc:"群会员数量"`
	GroupVipGroup        map[string]int `json:"group_vip_group" desc:"群会员数量"`
}

// GroupChatInfoResponse 社群详情返回参数
type GroupChatInfo struct {
	Unionid string `json:"unionid"`
}

// MsgPageListResponse 分页消息响应
type GroupChatUserListResponse struct {
	common.PageType
	List []WorkGroupChatMember `json:"list" desc:"列表数据"`
}

// GroupUsersAddEmployeesRequest 用户添加员工请求结构体
type GroupUsersAddEmployeesRequest struct {
	ExternalUserId string `json:"external_userid" db:"external_userid"`
	EnterpriseId   int64  `json:"enterprise_id" db:"enterprise_id"`
	Page           int64  `json:"page"`
	Size           int64  `json:"size"`
}

type GroupUsersJoinActivityRequest struct {
	UnionID      string `json:"union_id" db:"union_id"`
	EnterpriseId int64  `json:"enterprise_id" db:"enterprise_id"`
	Page         int64  `json:"page"`
	Size         int64  `json:"size"`
}

// MsgPageListResponse 分页消息响应
type GroupUsersAddEmployeesResponse struct {
	common.PageType
	List []GroupUsersAddEmployees `json:"list" desc:"列表数据"`
}

// MsgPageListResponse 分页消息响应
type GroupUsersJoinActivityResponse struct {
	common.PageType
	List []GroupUsersJoinActivity `json:"list" desc:"列表数据"`
}

// GroupUsersAddEmployeesResponse 用户添加员工返回结果结构体
type GroupUsersAddEmployees struct {
	Avatar                  string `json:"avatar" db:"avatar"`
	Name                    string `json:"name" db:"name"`
	Createtime              int64  `json:"createtime" db:"createtime"`
	CustomerRelationGuideId int64  `json:"customer_relation_guide_id" db:"customer_relation_guide_id"`
	IsBindingMember         int64  `json:"is_binding_member" db:"is_binding_member"`
}

// GroupUsersJoinActivity
type GroupUsersJoinActivity struct {
	Name      string `json:"name" db:"name"`
	CreatedAt int64  `json:"created_at" db:"created_at"`
	Status    string `json:"status" db:"status"`
}
