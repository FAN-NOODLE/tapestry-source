package socialDraw

import "micro-scrm/internal/modules/chat"

const (
	ChatMemberAll   = 0 // 全部
	ChatMemberJoin  = 1 // 新入群
	ChatMemberLeave = 2 // 新退群
)

// WorkGroupChat 客户群信息
type WorkGroupChat struct {
	chat.WorkGroupChatBase `json:"chat_info"`
	WorkGroupChatStatistic `json:"target"`
	NewGroup               []*Chart `json:"new"  desc:"新入群用户图表"`
	Valid                  []*Chart `json:"valid" desc:"留存用户图表"`
	UpdateTime             string   `json:"update_time" desc:"数据更新时间，固定为当天早上8：00"`
}

type WorkGroupChatInfo struct {
	chat.WorkGroupChatBase     `json:"chat_info"`
	WorkGroupChatStatisticInfo `json:"target"`
	NewGroup                   []*Chart `json:"new"  desc:"新入群用户图表"`
	Valid                      []*Chart `json:"valid" desc:"留存用户图表"`
	UpdateTime                 string   `json:"update_time" desc:"数据更新时间，固定为当天早上4：00"`
}

// WorkGroupChatStatistic 社群统计数据
type WorkGroupChatStatistic struct {
	ValidUv       int `json:"valid_uv" db:"lay_uv" desc:"留存群成员数"`
	AllLossUv     int `json:"all_loss_uv" db:"del_uv" desc:"累计退群人数"`
	AddUv         int `json:"add_uv" db:"add_uv" desc:"新入群人数"`
	NewLossUv     int `json:"new_loss_uv" db:"new_del" desc:"新退群人数"`
	MemberUv      int `json:"member_uv" db:"member_uv" desc:"会员人数"`
	GeneralUv     int `json:"general_uv" db:"general_uv" desc:"非会员人数"`
	WhiteGlodDown int `json:"glod_down" db:"glod_down" desc:"白金卡以下用户数"`
}

type WorkGroupChatStatisticInfo struct {
	WorkGroupChatStatistic
	MemberCardUv string `json:"member_card_uv"` // 会员率
	LossUvFee    string `json:"loss_uv_fee"`    // 退群率

}

type Chart struct {
	Name  string `json:"name" db:"name"`
	Value string `json:"value" db:"value"`
}

// WorkGroupChatMember 群成员
type WorkGroupChatMember struct {
	chat.WorkGroupChatMember
	WorkGroupChatMemberExt
}

type WorkGroupChatMemberList struct {
	UID       int64  `json:"-" db:"id"`
	UserID    string `json:"user_id" db:"user_id"`
	UnionID   string `json:"-" db:"unionid"`
	Name      string `json:"name" db:"name"`
	Nick      string `json:"nick" db:"nick"`
	Avatar    string `json:"avatar" db:"avatar"`
	Type      int    `json:"type" db:"type"` //1 - 企业成员,2 - 外部联系人
	Time      string `json:"time"`
	TmpTime   int64  `json:"-" db:"times"`
	TmpCkTime string `json:"-" db:"ck_time"`
	Friend    int64  `json:"friend" db:"friend"`
	MemberID  int64  `json:"-" db:"member_id"`
	JoinScene int    `json:"join_scene" db:"join_scene"`
	ActNum    int    `json:"act_num" db:"act_num"`
}
type ActInfo struct {
	ID         int64  `json:"-" db:"id"`
	Name       string `json:"name" db:"name"`
	Status     int    `json:"status" db:"status"`
	JoinTime   int64  `json:"join_time"`
	ActEndTime int64  `json:"end_time" db:"end_time"`
	DelTime    int64  `json:"deleted_at" db:"deleted_at"`
}

type ActInfoBase struct {
	Fission []Fission `json:"fission"`
	Market  []Market  `json:"market"`
}
type Fission struct {
	ActivityID int64 `json:"activity_id"`
	JoinAt     int64 `json:"join_at"`
}
type Market struct {
	ActivityID int64 `json:"activity_id"`
	JoinAt     int64 `json:"join_at"`
}

// WorkGroupChatMemberExt 群成员补充信息
type WorkGroupChatMemberExt struct {
	Name          string   `json:"name" db:"name" desc:"群成员名称"`
	GroupNickName string   `json:"group_nick_name" db:"group_nick_name" desc:"群成员昵称"`
	ChatID        string   `json:"chat_id" db:"chat_id" desc:"群ID"`
	VipStatus     int      `json:"vip_status" db:"-" desc:"是否会员：0-否，1-是"`
	VipCode       string   `json:"vip_code" db:"vip_code" desc:"会员卡号"`
	VipLevel      string   `json:"vip_level" db:"vip_level" desc:"会员卡等级"`
	Mobile        string   `json:"mobile" db:"mobile" desc:"手机号"`
	InviteUserID  string   `json:"invite_userid" db:"invite_userid" desc:"邀请人userid"`
	InviteUser    string   `json:"invite_username" db:"-" desc:"邀请人名称"`
	CustomerType  string   `json:"-" db:"customer_type" desc:"客户类型"`
	Guiders       []string `json:"guiders" db:"-" desc:"已添加导购"`
}

// WorkGroupChatType 群所属类型
type WorkGroupChatType struct {
	ID           int64    `json:"id" db:"id" desc:"主键"`
	Name         string   `json:"name" db:"name" desc:"群类型"`
	Names        []string `json:"-" db:"-" desc:"支持批量查询"`
	EnterpriseID int64    `json:"enterprise_id" db:"enterprise_id" desc:"企业ID" binding:"required"`
	Opt          int      `json:"opt" db:"-" desc:"类型管理：1-添加，2-编辑，3-删除"`
	GroupChatNum int      `json:"g_num" db:"g_num" desc:"设置了改类型的社群数"`
}

// WorkGroupListRequest 群列表请求参数
type WorkGroupListRequest struct {
	EnterpriseID int64    `json:"enterprise_id" db:"enterprise_id" desc:"企业ID" binding:"required"`
	ChatID       string   `json:"chat_id" db:"chat_id" desc:"群ID"`
	ChatIds      []string `json:"chat_ids" db:"chat_ids" desc:"群ID主键列表，支持批量查询"`
	ChatName     string   `json:"chat_name" db:"chat_name" desc:"群名"`
	ChannelID    int64    `json:"channel_id" db:"channel_id" desc:"群所属渠道"`
	Owner        int64    `json:"owner" db:"owner" desc:"群主member_id"`
	Operator     string   `json:"operator" db:"operator" desc:"运营者昵称"`
	DepartmentID int64    `json:"department_id" db:"department_id" desc:"运营者所属部门ID"`
	OperatorIds  []string `json:"-" db:"-" desc:"运营者+部门转换成运营者userid筛选"`
	ChatTypeID   int64    `json:"chat_type_id" db:"chat_type_id" desc:"群类型ID"`
	StartAt      int64    `json:"sts" db:"sts" desc:"筛选开始时间"`
	EndAt        int64    `json:"ets" db:"ets" desc:"筛选结束时间"`
}

// WorkGroupChatResponse 群列表请求数据
type WorkGroupChatResponse struct {
	chat.WorkGroupChatBase
	WorkGroupChatStatistic
}

// WorkGroupChatSettingRequest defines work group chat request
type WorkGroupChatSettingRequest struct {
	EnterpriseID  int64    `json:"enterprise_id" db:"enterprise_id" desc:"企业ID" validate:"required"`
	ChatID        string   `json:"chat_id" db:"chat_id" desc:"群ID"`
	ChatName      string   `json:"name" db:"name" desc:"群名"`
	ChannelID     int64    `json:"channel_id" db:"channel_id" desc:"群所属渠道"`
	ChannelName   string   `json:"channel_name" db:"channel_name" desc:"渠道名称"`
	ChatTypeID    int64    `json:"chat_type_id" db:"chat_type_id" desc:"群类型ID"`
	OperatorID    string   `json:"operator_id" db:"operator_id" desc:"运营者user_id"`
	Mobile        string   `json:"mobile" db:"mobile" desc:"运营者手机号，导入时使用"`
	Operator      string   `json:"-" db:"operator" desc:"运营者昵称"`
	DepartmentID  []int64  `json:"-" desc:"运营者所属部门ID"`
	DepartmentIds string   `json:"-" db:"department_ids" desc:""`
	Departments   []string `json:"-" db:"departments" desc:"部门列表"`
	DepartmentStr string   `json:"department_str" db:"department_str"`
	ChatTarget    string   `json:"chat_target" db:"chat_target" desc:"群目标"`
}

// WorkGroupChatImportSetting 导入群设置数据
type WorkGroupChatImportSetting struct {
	WorkGroupChatSettingRequest
	LineNum  int      `json:"line_num" db:"-" desc:"行数"`
	ChatType string   `json:"chat_type" db:"chat_type" desc:"群类型"`
	Errors   []string `json:"errors" db:"-" desc:"错误信息"`
}

// BatchImportResponse 批量导入响应
type BatchImportResponse struct {
	Success    int    `json:"success" desc:"成功数"`
	Repeated   int    `json:"repeated" desc:"重复数"`
	Errors     int    `json:"errors" desc:"错误数"`
	SettingKey string `json:"setting_key" desc:"设置存储的key"`

	SuccessItems interface{} `json:"-" desc:"成功列表"`
	ErrorItems   interface{} `json:"error_items" desc:"失败列表"`
}

// WorkGroupChatBatchSettingRequest 批量设置请求
type WorkGroupChatBatchSettingRequest struct {
	SettingKey string `json:"setting_key" desc:"前端提交的缓存key" validate:"required"`
}

type GroupInfoReq struct {
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID" validate:"required"`
	ChatID       string `json:"chat_id" db:"chat_id" desc:"群ID"`
	CID          int64  `json:"-" db:"cid" desc:"群聊表ID"`
	STS          int64  `json:"sts" db:"sts" desc:"开始时间"`
	MonthTS      int64  `json:"month_ts" db:"month_ts" desc:"一月前开始时间"`
	ETS          int64  `json:"ets" db:"ets" desc:"开始时间"`
	CardOn       bool   `json:"card_on" db:"-" desc:"是否展示会员卡数据"`
}

type ActivityReq struct {
	ChatID       string `json:"chat_id" db:"chat_id" desc:"群聊表id"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID" validate:"required"`
	UserID       string `json:"user_id" db:"user_id" desc:"参与活动用户ID"`
	Cid          int64  `json:"-" db:"cid" desc:"群聊表ID"`
}

// WorkGroupFilterRequest 群列表筛选项请求参数
type WorkGroupFilterRequest struct {
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID" validate:"required"`
	FiledName    string `json:"filter_name" db:"-" desc:"自定义筛选字段，非固定"`
}

// WorkGroupFilterResponse 群列表筛选项响应
type WorkGroupFilterResponse struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

// WorkGroupChatMembersRequest 成员明细请求参数
type WorkGroupChatMembersRequest struct {
	CID          int64  `json:"cid" db:"-" desc:"群ID"`
	ListType     int    `json:"list_type" db:"list_type" desc:"列表类型 1白金卡以下未入群 2群内会员 3 群内非会员"`
	Friend       int    `json:"friend" db:"friend" desc:"好友类型 0 全部 1 好友 2 非好友"`
	ChatID       string `json:"chat_id" db:"chat_id" desc:"社群ID"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id"`
	MemberID     int64  `json:"member_id" db:"member_id"`
	UserId       string `json:"user_id" db:"user_id"`
	MarsID       string `json:"-" db:"mars_id"`
	MemberType   int    `json:"member_type" db:"member_type" desc:"群成员类型：0-全部，1-新入群，2-新退群"`
	StartAt      int64  `json:"sts" db:"sts"`
	EndAt        int64  `json:"ets" db:"ets"`
	Pagination
}

// Pagination 分页请求
type Pagination struct {
	Page   int `json:"page" db:"-"`
	Offset int `json:"-" db:"offset"`
	Size   int `json:"size" db:"size"`
}

// HandleOffset 初始化页码，计算偏移量
func (p *Pagination) HandleOffset() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Size <= 0 {
		p.Size = 0
	}
	p.Offset = (p.Page - 1) * p.Size
}

// PaginationResponse 翻页响应
type PaginationResponse struct {
	Items      interface{} `json:"items"`
	TotalCount int         `json:"total_count"`
}

// WorkGroupChatStaticsRequest 企微社群统计请求参数
type WorkGroupChatStaticsRequest struct {
	GroupChatID []int64 `json:"-"`
	StartAt     int64   `json:"sts" db:"sts"`
	EndAt       int64   `json:"ets" db:"ets"`
}

// WxUserCorpExternal 微信外部联系人信息
type WxUserCorpExternal struct {
	ID              int64  `json:"-" db:"id" desc:"主键ID"`
	AppID           string `json:"-" db:"app_id" desc:"corp_id"`
	OpenID          string `json:"-" db:"openid" desc:"对应的是userid"`
	UnionID         string `json:"union_id" db:"unionid" desc:"微信unionid"`
	Nick            string `json:"nick" db:"nick" desc:"昵称"`
	Name            string `json:"name" db:"name" desc:"姓名"`
	Position        string `json:"position" db:"position" desc:"用户职称"`
	Avatar          string `json:"avatar" db:"avatar" desc:"头像"`
	UnionCode       string `json:"union_code" db:"union_code" desc:"会员卡号"`
	Point           string `json:"point" db:"point" desc:"会员积分"`
	Grade           string `json:"grade" db:"grade" desc:"会员等级"`
	StoreCode       string `json:"store_code" db:"store_code" desc:"门店编号"`
	BaCode          string `json:"ba_code" db:"ba_code" desc:"BA编号"`
	CorpName        string `json:"-" db:"corp_name" desc:"企业简称"`
	CorpFullName    string `json:"-" db:"corp_full_name" desc:"企业全陈"`
	UserType        int    `json:"user_type" db:"type" desc:"类型：1-微信用户，2-企业微信用户"`
	Sex             int    `json:"sex" db:"sex" desc:"性别：0-未知，1-男性，2-女性"`
	ExternalProfile string `json:"-" db:"external_profile" desc:"自定义展示信息"`
	EUid            int64  `json:"euid" db:"euid" desc:"enterprise_use.id"`
	Mobile          string `json:"mobile" db:"mobile" desc:"会员手机号"`
}

// WorkGroupChatDataGridResponse 统计数据视图结构
type WorkGroupChatDataGridResponse struct {
	Target *WorkGroupChatStatistic `json:"target" desc:"数据指标"`
	Chart  *DataGridChart          `json:"chart" desc:"饼图"`
}

// DataGridChart 饼图结构
type DataGridChart struct {
	New   []*Chart `json:"new"`
	Valid []*Chart `json:"valid"`
}
