package request

type EnterpriseBaseReq struct {
	EnterpriseId int64 `json:"enterprise_id" binding:"required" comment:"企业ID"`
}

type PageReq struct {
	Page int `json:"page" binding:"min=1" comment:"页码"`
	Size int `json:"size" binding:"min=1,max=999999" comment:"分页数"`
}

type MemberInfoReq struct {
	MemberName   string `json:"member_name"`
	MemberUserid string `json:"member_userid"`
	MemberID     int    `json:"member_third_id"`
}

type MemberSearchReq struct {
	UserId       string  `json:"user_id"`
	UserName     string  `json:"user_name"`
	ThirdId      int64   `json:"third_id"`
	MemberIdArr  []int64 `json:"member_id_arr"`
	EnterpriseId int64   `json:"enterprise_id"`
}

// SearchList 小程序端顾客列表，邀约任务，筛选信息的结构体数据
type SearchList struct {
	ManageType             int64    `json:"manage_type"`
	TagList                []int64  `json:"tag_list,omitempty" comment:"客户标签"`
	ImpressionTagList      []int64  `json:"impression_tag_list"  comment:"印象标签"`
	CltType                []string `json:"clt_type"  comment:"clt类型"`
	CardType               []string `json:"card_type"  comment:"卡等"`
	LastPurchaseStartTime  int      `json:"last_purchase_start_time" comment:"最近购买起始时间"`
	LastPurchaseEndTime    int      `json:"last_purchase_end_time" comment:"最近购买结束时间"`
	BirthdayStart          int      `json:"birthday_start" comment:"生日时间"`
	BirthdayEnd            int      `json:"birthday_end" comment:"生日时间"`
	CustomerType           []int    `json:"customer_type"  comment:"sw会员类型"`
	UnionCodes             []string `json:"union_codes"  comment:"会员code"`
	ReserveType            int      `json:"reserve_type" comment:"1 有进行中的邀约 2 无邀约 3 有成功的邀约" `
	TaskType               int      `json:"task_type" comment:"1 未触达 2 已触达" `
	TaskHave               int      `json:"task_have" comment:"1 5天内无任务 2 5天内有任务"`
	CouponType             int      `json:"coupon_type" comment:"优惠券筛选,1 有30天内临期券,2 有待使用券"`
	ConsumeType            []int    `json:"consume_type" comment:"近一年消费金额筛选,两个消费金额之间是-分割"`
	RecentConsumeStartTime int      `json:"recent_consume_start_time"  comment:"近期消费时间筛选：开始时间"`
	RecentConsumeEndTime   int      `json:"recent_consume_end_time"  comment:"近期消费时间筛选：结束时间"`
	IsSearch               int      `json:"is_search" comment:"是否筛选,0为否,1为是"`
	AgeTag                 int      `json:"age_tag" comment:"顾客列表查询，1为Gen Z(Z世代)，2为YM (千禧一代)"`
	IsUnion                int      `json:"is_union" comment:"是否好友关系，1好友，2非好友(任务顾客列表筛选)"`
	ChatTag                []int64  `json:"chat_tag,omitempty" comment:"顾客沟通意愿标签"`
}

type SendCardMessageToH5 struct {
	EnterpriseBaseReq
	Title    string              `json:"title"`
	Uri      string              `json:"uri"`
	Party    []string            `json:"party"`
	UserIds  []string            `json:"user_ids"`
	Content  []map[string]string `json:"content"`
	AppSign  string              `json:"app_sign"`
	Language string              `json:"language"`
}
