package customerSearchReq

import "micro-scrm/internal/app/http/request"

const (
	IsUnionTure  = iota + 1 // 筛选是否是好友，1为好友
	IsUnionFalse            // 筛选是否是好友，2为非好友
)

const (
	Contacted    = iota + 1 // 已联系
	NotContacted            // 未联系
)

const (
	IsKeySA = iota + 1
	NotKeySA
)
const (
	AgeTagOne   = 1 // Z世代
	AgeTagTwo   = 2 // 千禧一代
	AgeTagThree = 3 // OTHERS
)
const ValidDays = 30         // 检索的有效期
const TaskHaveValidDays = 30 // 检索是否有任务的有效期

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
	Self       = iota + 1 // 本人
	Enterprise            // 企业
)

const (
	TaskNoHave = 1 // 无任务
	TaskHave   = 2 // 有任务
)

const (
	Has30DayCoupon = iota + 1 // 有30天待使用券
	HasCoupon                 // 有待使用券
	BirthCoupon
	WelcomeCoupon
	RenewCoupon
	UpgradeCoupon
)

const (
	TaskTypeKey = "task_type"
	TaskHaveKey = "task_have"
)

const (
	TimeLt7D    = iota + 1 // 7天内
	Time8To15D             // 6-15
	Time16To30D            // 16-30
	TimeLt30D              // 30天内
	TimeGt30D              // 大于30天
	Time31To60D            // 31-60
	Time61To90D            //61-90
	TimeGt90D              //大于90
	TimeGt180D             //大于180
	TimeLt1D               // 1天内
	Time1To30D             // 1-30天
	TimeYesterday
	Time2To7D
	TimeGt15D // 大于15天
	TimeToday //今天
)

const (
	CategoryBD  = "包袋"  // 包袋
	CategoryXPJ = "小皮具" // 小皮具
	CategoryFS  = "服饰"  //服饰
	CategoryXL  = "鞋履"  //鞋履
	CategorySP  = "饰品"  //饰品
)

var CdpCategoryMap = map[string][]string{
	CategoryBD:  {"WHB", "MHB"},
	CategoryXPJ: {"WLG", "MLG"},
	CategoryFS:  {"WRTW", "MRTW"},
	CategoryXL:  {"WFTW", "MFTW"},
	CategorySP:  {"ACC"},
}

var AgeMap = map[int]string{
	AgeTagOne:   "Z世代",
	AgeTagTwo:   "千禧一代",
	AgeTagThree: "", // OTHERS
}

type ReqChoice struct {
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	MemberId     int64 `json:"memberId"  comment:"当前导购ID"`
	ThirdIds     []int `json:"third_ids" comment:"当前SM管理的门店"`
}

type ChoiceOptionList struct {
	CardType           []ChoiceOption `json:"card_type"`
	BirthMonth         []ChoiceOption `json:"birth_month"`
	ManageType         []ChoiceList   `json:"manage_type"`
	MemberList         []ChoiceList   `json:"member_list"`
	ReserveType        []ChoiceList   `json:"reserve_type"`
	TaskTypeValue      int            `json:"task_type_value"`
	TaskType           []ChoiceList   `json:"task_type"`
	TaskHaveValue      int            `json:"task_have_value"`
	TaskHave           []ChoiceList   `json:"task_have"`
	CouponType         []ChoiceList   `json:"coupon_type"`
	RelCreateTime      []ChoiceList   `json:"rel_create_time"`
	NotPurchasedDate   []ChoiceList   `json:"not_purchased_date"`
	RelInteractTime    []ChoiceList   `json:"rel_interact_time"`
	RelActiveTime      []ChoiceList   `json:"rel_active_time"`
	ChatTag            []ChoiceList   `json:"chat_tag"`
	HighWilling        []ChoiceOption `json:"high_willing"`
	PriceSensitivity   []ChoiceOption `json:"price_sensitivity"`
	CategoryPreference []ChoiceOption `json:"category_preference"`
	AgeTag             []ChoiceList   `json:"age_tag"`
}

type ChoiceList struct {
	Key   int64  `json:"key"`
	Value string `json:"value"`
}
type ChoiceOption struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ExternalSearchListBaseReq struct {
	OnlyCount      bool  `json:"only_count" comment:"是否只获取总数"`
	ThirdIds       []int `json:"third_ids" binding:"required_if=MemberType 4" required_if:"third_ids不能为空" comment:"当前SM管理的门店"`
	IsRelationship bool  `json:"is_relationship" comment:"检索类型,是否是查顾客关系,默认false只查顾客"`
	GetAll         bool  `json:"get_all"`
	Page           int   `json:"page"  comment:"页码"`
	Size           int   `json:"size"  comment:"分页数"`
	SearchList
}

// SearchList 小程序端顾客列表，邀约任务，筛选信息的结构体数据
type SearchList struct {
	IsTaskList           bool         `json:"is_task_list"`
	IsSearch             int          `json:"is_search" comment:"是否筛选,0为否,1为是"`
	EnterpriseID         int64        `json:"enterprise_id"`
	IsSelf               bool         `json:"is_self" comment:"是否只看自己的好友关系，0:否，1:是"`
	CardType             []string     `json:"card_type"  comment:"卡等"`
	BirthMonth           []string     `json:"birth_month" comment:"生日月份"`
	ManageType           int64        `json:"manage_type"`
	MemberIds            []int64      `json:"member_ids"  comment:"管理导购筛选的导购或者管理导购权限下的导购ids"`
	MemberId             int64        `json:"member_id"  comment:"当前导购id"`
	Name                 string       `json:"name" comment:"顾客名称"`
	ChatTag              int64        `json:"chat_tag" comment:"顾客沟通意愿标签"`
	ReserveType          int          `json:"reserve_type" comment:"1 有进行中的邀约 2 无邀约 3 有成功的邀约" `
	TaskType             int          `json:"task_type" comment:"1 未触达 2 已触达" `
	TaskHave             int          `json:"task_have" comment:"1 5天内无任务 2 5天内有任务"`
	CouponType           []int        `json:"coupon_type" comment:"优惠券筛选,1 有30天内临期券,2 有待使用券,3 有生日券 4 续会 5 欢迎 升级"`
	PurchasedDate        SearchTime   `json:"purchased_date" comment:"历史消费时间"`
	NotPurchasedDate     int          `json:"not_purchased_date" comment:"近期未消费时间"`
	HistoryAmountTotal   SearchNumber `json:"history_amount_total" comment:"历史消费总金额"`
	AmountTotal          SearchNumber `json:"amount_total" comment:"近一年消费金额"`
	RelCreateTime        []int        `json:"rel_create_time" comment:"添加好友时间"`
	RelCreateTimeType    int          `json:"rel_create_time_type" comment:"添加好友时间类型 1本人 2企业"`
	RelInteractTime      []int        `json:"rel_interact_time" comment:"上次触达时间"`
	RelInteractTimeType  int          `json:"rel_interact_time_type" comment:"上次触达时间类型 1本人 2企业"`
	RelActiveTime        []int        `json:"rel_active_time" comment:"上次活跃沟通过时间"`
	RelActiveTimeType    int          `json:"rel_active_time_type" comment:"上次活跃沟通过时间类型 1本人 2企业"`
	HighWilling          string       `json:"high_willing" comment:"高意愿客户 Y是 N否"`
	PriceSensitivity     string       `json:"price_sensitivity" comment:"价格敏感度 H高 L低"`
	CategoryPreference   []string     `json:"category_preference" comment:"商品偏好"`
	AgeTag               int          `json:"age_tag" comment:"顾客列表查询，1为Gen Z(Z世代)，2为YM (千禧一代) 3为 Other"`
	RelationshipType     int          `json:"relationship_type" comment:"是否好友关系，1好友，2非好友(任务顾客列表筛选)"`
	UnionIds             []string     `json:"union_ids" comment:"使用unionid集合作为参数"`
	FilterExternalUserId []int64      `json:"filter_external_user_id" comment:"需要过滤的外部联系人id"`
	RelationId           []int64      `json:"relation_id" comment:"选中的relationId"`
	FilterRelationId     []int64      `json:"filter_relation_id" comment:"过滤的relationId"`
	ExternalUserIds      []int64      `json:"external_user_ids" comment:"使用external_user_id集合作为参数"`
	UnionCodes           []string     `json:"union_codes"  comment:"会员code"`
}

type SearchTime struct {
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}
type SearchNumber struct {
	StartNumber int64 `json:"start_number"`
	EndNumber   int64 `json:"end_number"`
}

type AddCrowdReq struct {
	CrowdId int64 `json:"crowd_id"`
	request.EnterpriseBaseReq
	MemberId int64         `json:"member_id"`
	Name     string        `json:"name"`
	ThirdIds []int         `json:"third_ids" binding:"required_if=MemberType 4" required_if:"third_ids不能为空" comment:"当前SM管理的门店"`
	Search   AddCrowSearch `json:"search" comment:"筛选条件"`
	Sort     int           `json:"sort"`
}

type CrowdListReq struct {
	request.EnterpriseBaseReq
	MemberId int64 `json:"member_id"`
}

type CrowdInfoReq struct {
	CrowdId         int64 `json:"crowd_id" binding:"required" required:"人群包id必传"`
	ContactedStatus int   `json:"contacted_status"`
	OnlyCount       bool  `json:"only_count" comment:"是否只获取总数"`
	Page            int   `json:"page"  comment:"页码"`
	Size            int   `json:"size"  comment:"分页数"`
	SearchList
}

type ExternalDetailReq struct {
	request.EnterpriseBaseReq
	MemberID int64  `json:"member_id"`
	Openid   string `json:"openid"`
}

type AddCrowSearch struct {
	Map   string `json:"map" comment:"前端使用的值"`
	Value string `json:"value"`
}

type CrowdEditReq struct {
	request.EnterpriseBaseReq
	CrowdId  int64  `json:"crowd_id" binding:"required" required:"人群包id必传"`
	MemberId int64  `json:"member_id"`
	Name     string `json:"name"`
}

type CrowdCheckReq struct {
	request.EnterpriseBaseReq
	Id       int64 `json:"id"`
	MemberId int64 `json:"member_id"`
}

type UpdateSortInfo struct {
	CrowdId int64 `json:"crowd_id"`
	Sort    int   `json:"sort"`
}

type CrowdUpdateSortReq struct {
	request.EnterpriseBaseReq
	SortList []UpdateSortInfo `json:"sort_list"`
}

type CrowdUpdateContactReq struct {
	request.EnterpriseBaseReq
	CrowdId            int64 `json:"crowd_id"`
	MemberId           int64 `json:"member_id"`
	WorkExternalUserId int64 `json:"work_external_user_id"`
}

type BatchCrowdUpdateContactReq struct {
	request.EnterpriseBaseReq
	CrowdId             int64    `json:"crowd_id"`
	MemberId            int64    `json:"member_id"`
	OpenIds             []string `json:"open_ids"`
	WorkExternalUserIds []int64  `json:"work_external_user_ids"`
}
