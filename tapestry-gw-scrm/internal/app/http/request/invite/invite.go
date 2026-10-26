package invite

type EnterpriseBase struct {
	EnterpriseId int64 `json:"enterprise_id" comment:"企业ID"`
}

type PageReq struct {
	Page int `json:"page" comment:"页码"`
	Size int `json:"size" comment:"分页数"`
}

type TemplateCreateReq struct {
	EnterpriseBase
	ReserveServiceId int64    `json:"reserve_service_id" binding:"required" comment:"服务类型id"`
	TitleName        string   `json:"title_name" binding:"required" comment:"标题名称"`
	StartTime        int      `json:"start_time"  comment:"有效开始时间"`
	EndTime          int      `json:"end_time"  comment:"有效截止时间"`
	BackgroundUrl    []string `json:"background_url"  binding:"required" comment:"背景图多个url"`
	Content          string   `json:"content"  binding:"required" comment:"文案"`
	ShareTitle       string   `json:"share_title"  binding:"required" comment:"分享卡片标题"`
	ShareUrl         []string `json:"share_url"  binding:"required" comment:"分享图片url"`
	ThirdIds         []int    `json:"third_ids"  binding:"required" comment:"部门id"`
}

type TemplateSaveReq struct {
	EnterpriseBase
	ID               int64    `json:"id"`
	ReserveServiceId int64    `json:"reserve_service_id" binding:"required" comment:"服务类型id"`
	TitleName        string   `json:"title_name" binding:"required" comment:"标题名称"`
	StatusType       int      `json:"status_type" binding:"min=0,max=1" comment:"有效状态"`
	StartTime        int      `json:"start_time"  comment:"有效开始时间"`
	EndTime          int      `json:"end_time"  comment:"有效截止时间"`
	BackgroundUrl    []string `json:"background_url"  binding:"required" comment:"背景图多个url"`
	Content          string   `json:"content"  binding:"required" comment:"文案"`
	ShareTitle       string   `json:"share_title"  binding:"required" comment:"分享卡片标题"`
	ShareUrl         []string `json:"share_url"  binding:"required" comment:"分享图片url"`
	ThirdIds         []int    `json:"third_ids"   comment:"部门id"`
	UserIsShow       int      `json:"user_is_show" comment:"是否展示顾客称谓:1展示2不展示"`
	MemberIsShow     int      `json:"member_is_show" comment:"是否展示导购称谓:1展示2不展示"`
}

type TemplateListReq struct {
	EnterpriseBase
	PageReq
	Keyword          string `json:"keyword" comment:"搜索关键字"`
	ReserveServiceId int64  `json:"reserve_service_id" comment:"服务类型id"`
	IsUp             int    `json:"is_up" comment:"是否上下架:1上架2下架"`
}

type TemplateInfoReq struct {
	EnterpriseBase
	Id int64 `json:"id"`
}

type TemplateDeleteReq struct {
	EnterpriseBase
	ID int64 `json:"id"`
}

type RecordOrderListReq struct {
	EnterpriseBase
	PageReq
	TitleName               string `json:"title_name" comment:"搜索模板关键字"`
	InviteTemplateId        int64  `json:"invite_template_id"`
	ReserveServiceID        int64  `json:"reserve_service_id" `
	MemberName              string `json:"member_name" comment:"搜索导购姓名"`
	UserName                string `json:"user_name" comment:"搜索顾客姓名"`
	ThirdIds                []int  `json:"third_ids" comment:"门店ids"`
	ReserveServiceStartTime int    `json:"reserve_service_start_time" `
	ReserveServiceEndTime   int    `json:"reserve_service_end_time" `
	OrderCreateStartTime    int    `json:"order_create_start_time" `
	OrderCreateEndTime      int    `json:"order_create_end_time" `
	Total                   int64  `json:"total"`
}

type RecordDetailInfo struct {
	Id                 int64  `json:"id" `
	MemberId           int64  `json:"member_id" `
	MemberUserid       string `json:"member_userid" `
	UserName           string `json:"user_name" `
	UserUnionid        string `json:"user_unionid" `
	UserOpenid         string `json:"user_openid" `
	TitleName          string `json:"title_name" `
	BackgroundURL      string `json:"background_url" `
	Content            string `json:"content" `
	ShareTitle         string `json:"share_title" `
	ShareURL           string `json:"share_url" `
	ReserveServiceID   int64  `json:"reserve_service_id" `
	ReserveServiceName string `json:"reserve_service_name" `
	ReserveServiceTime int    `json:"reserve_service_time" `
}

type RecordOrder struct {
	RecordDetailInfo
	MemberName     string `json:"member_name" `
	StoreId        int    `json:"store_id" `
	StoreName      string `json:"store_name" `
	UserName       string `json:"user_name" `
	ReserveTime    int    `json:"reserve_time" `
	ReserveTimeStr string `json:"reserve_time_str" `
	CreatedAt      int    `json:"created_at" `
	CreateTimeStr  string `json:"create_time_str" `
}

type RecordOrderList struct {
	Total int64 `json:"total"`
	List  []RecordOrder
}

type TemplateUpdateStatusReq struct {
	EnterpriseBase
	ID   int64 `json:"id"`
	IsUp int   `json:"is_up"`
}
