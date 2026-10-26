package customerReq

import (
	"micro-scrm/internal/app/http/request"
	customerSearchReq "micro-scrm/internal/app/http/request/customer_search"
	tagsReq "micro-scrm/internal/app/http/request/tags"
	tagModels "micro-scrm/internal/app/models/tag"
)

type PageReq struct {
	Page int `json:"page" binding:"min=1" comment:"页码"`
	Size int `json:"size" binding:"min=10,max=99999999" comment:"分页数"`
}

type ExternalUserListsReq struct {
	PageReq
	ExternalUserListBaseReq
}

type ExternalUserListBaseReq struct {
	EnterpriseID         int64    `json:"enterprise_id" comment:"企业ID"`
	IsGroup              bool     `json:"is_group" comment:"是否分组"`
	Mobile               string   `json:"mobile"  comment:"当前导购手机号"`
	MemberType           int      `json:"member_type"  comment:"当前顾问类型"`
	ThirdIds             []int    `json:"third_ids" binding:"required_if=MemberType 4" required_if:"third_ids不能为空" comment:"当前SM管理的门店"`
	Name                 string   `json:"name"  comment:"顾客姓名"`
	MemberId             int64    `json:"member_id"  comment:"导购id"`
	Userid               string   `json:"userid"  comment:"导购userid"`
	Userids              []string `json:"userids"  comment:"导购userids"`
	MemberList           []int64  `json:"member_list"  comment:"选择的导购ids"`
	MemberLists          []int64  `json:"member_lists"  comment:"导购权限下的ids"`
	SearchType           int      `json:"search_type"`
	IsSelf               int      `json:"is_self" comment:"是否只看自己的好友关系，0:否，1:是"`
	FilterExternalUserId []int64  `json:"filter_external_user_id" comment:"需要过滤的外部联系人id"`
	FilterRelationId     []int64  `json:"filter_relation_id" comment:"需要过滤的relation_id"`
	request.SearchList
}

const (
	HasValidReserve   = iota + 1 // 有进行中的邀约
	NoReserve                    // 无邀约
	HasSuccessReserve            // 有成功的邀约
)

const (
	NoFinishTask  = iota + 1 // 任务未触达
	HasFinishTask            // 任务已触达
)

const (
	TaskHave   = iota + 1 // 无任务
	TaskNoHave            // 有任务
)

const (
	Has30DayCoupon = iota + 1 // 有30天待使用券
	HasCoupon                 // 有待使用券
)

type CustomerListsReq struct {
	EnterpriseID int64    `json:"enterprise_id" comment:"企业ID"`
	Mobile       string   `json:"mobile"  comment:"当前导购手机号"`
	MemberId     int64    `json:"member_id"  comment:"导购id"`
	Userid       string   `json:"userid"  comment:"导购userid"`
	UnionCodes   []string `json:"union_codes"  comment:"会员code"`
	MemberType   int      `json:"member_type"  comment:"当前顾问类型"`
	ThirdIds     []int    `json:"third_ids" binding:"required_if=MemberType 4" required_if:"third_ids不能为空" comment:"当前SM管理的门店"`
	MemberList   []int64  `json:"member_list"  comment:"选择的导购ids"`
}

type ExternalUseridReq struct {
	EnterpriseID int64    `json:"enterprise_id" comment:"企业ID"`
	Mobile       string   `json:"mobile"  comment:"当前导购手机号"`
	MemberType   int      `json:"member_type"  comment:"当前顾问类型"`
	ThirdIds     []int    `json:"third_ids" binding:"required_if=MemberType 4" required_if:"third_ids不能为空" comment:"当前SM管理的门店"`
	MemberId     int64    `json:"member_id"  comment:"导购id"`
	Userid       string   `json:"userid"  comment:"导购userid"`
	Userids      []string `json:"userids"  comment:"导购userids"`
	ManageType   int      `json:"manage_type"  comment:"管理类型"`
}

const (
	HIGH_GRADE = "10" // 优质客
	POTENTIAL  = "20" // 潜在客

	PROSPECT = "PROSPECT" // 注册会员
)

type ReqChoice struct {
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	Mobile       string `json:"mobile"  comment:"当前导购手机号"`
	MemberType   int    `json:"member_type"  comment:"当前顾问类型"`
	ThirdIds     []int  `json:"third_ids" comment:"当前SM管理的门店"`
}

type ChoiceList struct {
	Key   int64  `json:"key"`
	Value string `json:"value"`
}
type ChoiceOption struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type ReserveOption struct {
	Key   int    `json:"key"`
	Value string `json:"value"`
}
type CouponOption struct {
	Key   int    `json:"key"`
	Value string `json:"value"`
}
type ConsumeOption struct {
	Key   int    `json:"key"`
	Value string `json:"value"`
}

// ConsumeOptionMap 条件对应的实际查询金额
var ConsumeOptionMap = map[int][]int{
	1: {1, 2999},
	2: {3000, 4999},
	3: {5000, 7999},
	4: {8000, 9999},
	5: {10000, 19999},
	6: {20000, 29999},
	7: {30000, 9999999},
}

type TaskOption struct {
	Key   int    `json:"key"`
	Value string `json:"value"`
}
type ChoiceOptionList struct {
	ReserveType       []ReserveOption `json:"reserve_type"`
	TaskType          []TaskOption    `json:"task_type"`
	CouponType        []CouponOption  `json:"coupon_type"`
	ConsumeType       []ConsumeOption `json:"consume_type"`
	ManageType        []ChoiceList    `json:"manage_type"`
	MemberList        []ChoiceList    `json:"member_list"`
	CardType          []ChoiceOption  `json:"card_type"`
	CltType           []ChoiceOption  `json:"clt_type"`
	TagList           []ChoiceList    `json:"tag_list"`
	ImpressionTagList []ChoiceList    `json:"impression_tag_list"`
	CustomerType      []ChoiceList    `json:"customer_type"  comment:"sw"`
	AgeTag            []ChoiceList    `json:"age_tag" comment:"Coach的顾客类型，Gen Z(Z世代),YM(千禧一代)"`
	TaskHave          []TaskOption    `json:"task_have"`
	IsUnion           []TaskOption    `json:"is_union"`
	ChatTag           []ChoiceList    `json:"chat_tag" comment:"Coach的沟通意愿"`
}

type ExternalUserDetailReq struct {
	EnterpriseId   int64  `json:"enterprise_id" binding:"required" comment:"企业ID"`
	Mobile         string `json:"mobile" comment:"手机号"`
	ExternalUserid string `json:"external_userid"`
	MemberId       int64  `json:"member_id"`
	MemberUserid   string `json:"member_userid"`
	MemberType     int    `json:"member_type"  comment:"当前顾问类型"`
	ThirdIds       []int  `json:"third_ids" binding:"required_if=MemberType 4" required_if:"third_ids不能为空" comment:"当前SM管理的门店"`
	IsWork         bool   `json:"is_work"  comment:"是否对话框打开 true 是 false 否"`
	CustomerCode   string `json:"customer_code"`
}

type ExternalUserTagsInfoRes struct {
	ID                 int64  `json:"id"`
	Avatar             string `json:"avatar"`
	Nick               string `json:"nick"`
	Name               string `json:"name"`
	WorkExternalUserId string `json:"work_external_user_id"`
	Unionid            string `json:"unionid"`
	Status             int    `json:"status"`
}

type ExternalUserTagsInfoListsRes struct {
	ExternalUserTagsInfoRes
	TgName []string `json:"tg_name"`
}

type UserDetailInfo struct {
	UserID             int64  `json:"user_id"`
	UnionID            string `json:"union_id"`
	ExternalUserID     string `json:"external_userid"`
	WorkExternalUserID int64  `json:"work_external_user_id"`
	NickName           string `json:"nick_name"`
	Avatar             string `json:"avatar"`
	Mobile             string `json:"mobile"`
	RelationID         int64  `json:"relation_id"`
	Remark             string `json:"remark"`
	IsRelation         int    `json:"is_relation"`
	UnionCode          string `json:"union_code"`

	UserTags               []tagsReq.UserTagsGroupList   `json:"user_tags"`
	UserImpressionTags     []tagModels.UserTagsGroupList `json:"user_impression_tags"`
	UserImpressionChatTags []tagModels.UserTagsGroupList `json:"user_impression_chat_tags"`
	ConfigStoreName        string                        `json:"config_store_name"`
	UserChatTags           tagModels.UserTagsGroupList   `json:"user_chat_tags"`
}

type UserDetailInfoCoach struct {
	UserDetailInfo
	Birthday        string `json:"birthday"`
	CltType         string `json:"clt_type"`
	CustomerCode    string `json:"customer_code" comment:"会员customer_code"`
	CustomerCodeStr string `json:"customer_code_str" comment:"会员customer_code名称"`
	CardLevel       string `json:"card_level" comment:"会员卡等"`
	Membership      string `json:"membership"  comment:"会员卡号"`
	AmountTotal     int64  `json:"amount_total" comment:"12个月消费总金额"`
}

type UserDetailInfoSW struct {
	UserDetailInfo
	Birthday        string      `json:"birthday" `
	CltType         string      `json:"clt_type"`
	CustomerCode    string      `json:"customer_code" comment:"会员customer_code"`
	CustomerCodeStr string      `json:"customer_code_str" comment:"会员customer_code名称"`
	CardLevel       string      `json:"card_level" comment:"会员卡等"`
	Membership      string      `json:"membership"  comment:"会员卡号"`
	AmountTotal     int64       `json:"amount_total"  comment:"12个月消费总金额"`
	UserInfoAdd     UserInfoAdd `json:"user_info_add"`

	LastPurchaseDate      int      `json:"last_purchase_date"  comment:"最近购买时间"`
	LastPurchaseDateStr   string   `json:"last_purchase_date_str" comment:"最近购买时间"`
	LastPurchaseStoreCode string   `json:"last_purchase_store_code" comment:"最近购买门店code"`
	LastPurchaseStoreName string   `json:"last_purchase_store_name" comment:"最近购买门店名称"`
	UpgradeNeedAmount     int      `json:"upgrade_need_amount" comment:"升级所需金额"`
	CustomerType          int      `json:"customer_type" comment:"顾客类型"`
	CrmTag                []string `json:"crm_tag" comment:"CRM标签"`
}

type UserDetailInfoKS struct {
	UserDetailInfo
	Birthday        string `json:"birthday"`
	CustomerCode    string `json:"customer_code" comment:"会员customer_code"`
	CustomerCodeStr string `json:"customer_code_str" comment:"会员customer_code名称"`
	CardLevel       string `json:"card_level" comment:"会员卡等"`
	Membership      string `json:"membership" comment:"会员卡号"`
	AmountTotal     int64  `json:"amount_total"   comment:"ks36个月消费总金额"`
	City            string `json:"city"`

	LastPurchaseDate      int    `json:"last_purchase_date" comment:"最近购买时间"`
	LastPurchaseDateStr   string `json:"last_purchase_date_str" comment:"最近购买时间"`
	LastPurchaseStoreCode string `json:"last_purchase_store_code" comment:"最近购买门店code"`
	LastPurchaseStoreName string `json:"last_purchase_store_name" comment:"最近购买门店名称"`
	UpgradeNeedAmount     int    `json:"upgrade_need_amount" comment:"升级所需金额"`
	CustomerType          int    `json:"customer_type" comment:"顾客类型"`
}

// OnInheritReq 在职分配全选的条件
type OnInheritReq struct {
	EnterpriseID     int64   `json:"enterprise_id" comment:"企业ID"`
	ThirdIds         []int   `json:"third_ids" comment:"当前SM管理的门店"`
	Name             string  `json:"name"  comment:"顾客姓名"`
	MemberList       []int64 `json:"member_list"  comment:"选择的导购ids"`
	MemberId         int64   `json:"member_id"  comment:"导购的id"`
	RelationId       []int64 `json:"relation_id" comment:"选中的relationId"`
	FilterRelationId []int64 `json:"filter_relation_id" comment:"过滤的relationId"`
	request.SearchList
}

// OnInheritSearchReq  在职分配全选的条件
type OnInheritSearchReq struct {
	EnterpriseID int64 `json:"enterprise_id" comment:"企业ID"`
	ThirdIds     []int `json:"third_ids" comment:"当前SM管理的门店"`
	customerSearchReq.SearchList
}

// OnInheritExternalInfoRes 在职分配的顾客信息
type OnInheritExternalInfoRes struct {
	EnterpriseID       int64  `json:"enterprise_id" comment:"企业ID"`
	MemberId           int64  `json:"member_id"  comment:"导购的id"`
	WorkExternalUserId int64  `json:"work_external_user_id"`
	UserId             string `json:"user_id"`
}
