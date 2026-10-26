package customerSearchReq

type ChoiceOptionSwList struct {
	CardType           []ChoiceOption `json:"card_type"`
	CustomerType       []ChoiceList   `json:"customer_type"`
	BirthMonth         []ChoiceOption `json:"birth_month"`
	ManageType         []ChoiceList   `json:"manage_type"`
	MemberList         []ChoiceList   `json:"member_list"`
	ReserveType        []ChoiceList   `json:"reserve_type"`
	TaskType           []ChoiceList   `json:"task_type"`
	TaskHave           []ChoiceList   `json:"task_have"`
	RelCreateTime      []ChoiceList   `json:"rel_create_time"`
	RelInteractTime    []ChoiceList   `json:"rel_interact_time"`
	CouponType         []ChoiceList   `json:"coupon_type"`
	NotPurchasedDate   []ChoiceList   `json:"not_purchased_date"`
	CategoryPreference []ChoiceOption `json:"category_preference"`
	PriceTag           []ChoiceOption `json:"price_tag"`
	AgeTag             []ChoiceOption `json:"age_tag"`
	AtPurchasedYear    []ChoiceOption `json:"at_purchased_year"`
	AtPurchasedQuarter []ChoiceOption `json:"at_purchased_quarter"`
	AtPurchasedAmount  []ChoiceList   `json:"at_purchased_amount"`
	TagList            []ChoiceList   `json:"tag_list"`
	ImpressionTagList  []ChoiceList   `json:"impression_tag_list"`
	ChatTag            []ChoiceList   `json:"chat_tag" comment:"Coach的沟通意愿"`
}

// SearchListSw  小程序端顾客列表，邀约任务，筛选信息的结构体数据
type SearchListSw struct {
	EnterpriseID         int64        `json:"enterprise_id"`
	CardType             []string     `json:"card_type"  comment:"卡等"`
	CustomerType         []int        `json:"customer_type"`
	BirthMonth           []string     `json:"birth_month" comment:"生日月份"`
	ManageType           int64        `json:"manage_type"`
	MemberIds            []int64      `json:"member_ids"  comment:"管理导购筛选的导购或者管理导购权限下的导购ids"`
	MemberId             int64        `json:"member_id"  comment:"当前导购id"`
	UserId               string       `json:"user_id" comment:"当前导购企微id"`
	Name                 string       `json:"name" comment:"顾客名称"`
	ChatTag              int64        `json:"chat_tag" comment:"顾客沟通意愿标签"`
	ReserveType          int          `json:"reserve_type" comment:"1 有进行中的邀约 2 无邀约 3 有成功的邀约" `
	TaskType             int          `json:"task_type" comment:"1 未触达 2 已触达" `
	TaskHave             int          `json:"task_have" comment:"1 5天内无任务 2 5天内有任务"`
	RelCreateTime        []int        `json:"rel_create_time" comment:"添加好友时间"`
	RelInteractTime      []int        `json:"rel_interact_time" comment:"上次触达时间"`
	UnionIds             []string     `json:"union_ids" comment:"使用unionid集合作为参数"`
	FilterExternalUserId []int64      `json:"filter_external_user_id" comment:"需要过滤的外部联系人id"`
	RelationId           []int64      `json:"relation_id" comment:"选中的relationId"`
	FilterRelationId     []int64      `json:"filter_relation_id" comment:"过滤的relationId"`
	ExternalUserIds      []int64      `json:"external_user_ids" comment:"使用external_user_id集合作为参数"`
	UnionCodes           []string     `json:"union_codes"  comment:"会员code"`
	PurchasedDate        SearchTime   `json:"purchase_date" comment:"消费时间"`
	CouponType           []int        `json:"coupon_type" comment:"优惠券筛选,1 有30天内临期券,2 有待使用券,3 有生日券 4 新人 5 回购 6活动"`
	NotPurchasedDate     int          `json:"not_purchased_date" comment:"近期未消费时间"`
	R12PurchaseAmount    SearchNumber `json:"r12_purchase_amount" comment:"年消费金额"`
	CategoryPreference   []string     `json:"category_preference" comment:"商品偏好"`
	PriceTag             []string     `json:"price_tag"`
	AgeTag               []string     `json:"age_tag"`
	AtPurchasedYear      string       `json:"at_purchased_year"`
	AtPurchasedQuarter   string       `json:"at_purchased_quarter"`
	AtPurchasedAmount    []int        `json:"at_purchased_amount"`
	TagList              []int64      `json:"tag_list" comment:"客户标签"`
	ImpressionTagList    []int64      `json:"impression_tag_list"  comment:"印象标签"`
	RelationshipType     int          `json:"relationship_type" comment:"是否好友关系，1好友，2非好友"`
}

type ExternalSearchListBaseSwReq struct {
	OnlyCount      bool  `json:"only_count" comment:"是否只获取总数"`
	ThirdIds       []int `json:"third_ids" binding:"required_if=MemberType 4" required_if:"third_ids不能为空" comment:"当前SM管理的门店"`
	IsRelationship bool  `json:"is_relationship" comment:"检索类型,是否是查顾客关系,默认false只查顾客"`
	GetAll         bool  `json:"get_all"`
	IsInherit      bool  `json:"is_inherit"`
	Page           int   `json:"page"  comment:"页码"`
	Size           int   `json:"size"  comment:"分页数"`
	SearchListSw
}

type CrowdInfoSwReq struct {
	CrowdId         int64 `json:"crowd_id" binding:"required" required:"人群包id必传"`
	ContactedStatus int   `json:"contacted_status"`
	OnlyCount       bool  `json:"only_count" comment:"是否只获取总数"`
	Page            int   `json:"page"  comment:"页码"`
	Size            int   `json:"size"  comment:"分页数"`
	SearchListSw
}

const (
	SwCategoryDX  = "短靴"
	SwCategoryCX  = "长靴"
	SwCategoryLX  = "凉鞋"
	SwCategoryGGX = "高跟鞋"
	SwCategoryYDX = "运动鞋"
	SwCategoryPDX = "平底鞋"
	//SwCategoryPGHDX = "坡跟厚底鞋"
)

const (
	SwHas30DayCoupon    = iota + 1 // 有30天待使用券
	SwHasCoupon                    // 有待使用券
	SwBirthCoupon                  //生日礼券
	SwNewMemberCoupon              //新人礼券
	SwRebateCoupon                 //回购礼券
	SwPromotionalCoupon            //活动礼券
)

const (
	SwAtPurchasedAmountGt8000     = iota + 1 // >8000
	SwAtPurchasedAmount6000To8000            // >=6000 && < 8000
	SwAtPurchasedAmount4000To6000            // >=4000 && < 6000
	SwAtPurchasedAmountLt4000                // <4000
)

const (
	SwPriceTag1 = "仅折扣"
	SwPriceTag2 = "仅正价"
	SwPriceTag3 = "折扣+正价"
)
