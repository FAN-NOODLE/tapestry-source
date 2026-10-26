package invite

type EnterpriseBase struct {
	EnterpriseId int64 `json:"enterprise_id" binding:"required" comment:"企业ID"`
}

type PageReq struct {
	Page int `json:"page" binding:"min=1" comment:"页码"`
	Size int `json:"size" binding:"min=10,max=1000" comment:"分页数"`
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
	ThirdIds         []int    `json:"third_ids" comment:"部门id"`
	UserIsShow       int      `json:"user_is_show" comment:"是否展示顾客称谓:1展示2不展示"`
	MemberIsShow     int      `json:"member_is_show" comment:"是否展示导购称谓:1展示2不展示"`
}

type TemplateSaveRes struct {
	ID int64 `json:"id"`
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
	ID int64 `json:"id"`
}

type TemplateDeleteReq struct {
	EnterpriseBase
	ID int64 `json:"id"`
}

type OrderCountNum struct {
	InviteTemplateId int64 `json:"invite_template_id"`
	Num              int   `json:"num"`
}

type TemplateConfigListReq struct {
	EnterpriseBase
	Mobile   string `json:"mobile" comment:"手机号"`
	ThirdId  int64  `json:"third_ids"`
	MemberId int64  `json:"member_id"`
}

type TemplateConfigDetailReq struct {
	EnterpriseBase
	Mobile string `json:"mobile" comment:"手机号"`
	ID     int64  `json:"id"`
}

type RecordInfoReq struct {
	EnterpriseBase
	RecordInfo
}

type RecordInfo struct {
	InviteTemplateId int64  `json:"invite_template_id" binding:"required"`
	MemberId         int64  `json:"member_id" `
	MemberUserid     string `json:"member_userid" `
	UserUnionid      string `json:"user_unionid" `
	UserOpenid       string `json:"user_openid" binding:"required"`
	Mobile           string `json:"mobile"`
	StyleType        int    `json:"style_type"  binding:"required"`
	BackgroundUrl    string `json:"background_url"  binding:"required"`
	SendContent      string `json:"send_content"  binding:"required"`
}

type GetRecordInfoReq struct {
	EnterpriseBase
	InviteRecordId int64  `json:"invite_record_id" binding:"required"`
	UnionId        string `json:"union_id"`
}

type RecordDetailInfo struct {
	Id                 int64  `json:"id" `
	MemberId           int64  `json:"member_id" `
	MemberName         string `json:"member_name" `
	MemberUserid       string `json:"member_userid" `
	UserName           string `json:"user_name" `
	UserUnionid        string `json:"user_unionid" `
	UserOpenid         string `json:"user_openid" `
	TitleName          string `json:"title_name" `
	Content            string `json:"content" `
	ShareTitle         string `json:"share_title" `
	ShareURL           string `json:"share_url" `
	ReserveServiceID   int64  `json:"reserve_service_id" `
	ReserveServiceName string `json:"reserve_service_name" `
	ReserveServiceTime int    `json:"reserve_service_time" `
	ReserveOrderId     int64  `json:"reserve_order_id" `
	StyleType          int    `json:"style_type" `
	BackgroundURL      string `json:"background_url" `
	SendContent        string `json:"send_content"`
	ReserveTime        int    `json:"reserve_time"`
	IsUp               int    `json:"is_up"`
	UserIsShow         int    `json:"user_is_show" comment:"是否展示顾客称谓:1展示2不展示"`
	MemberIsShow       int    `json:"member_is_show" comment:"是否展示导购称谓:1展示2不展示"`
	StoreId            int64  `json:"store_id"`
	StoreCode          string `json:"store_code"`
	StoreName          string `json:"store_name"`
	TimeSlotStr        string `json:"-" db:"time_slot_str"`
}

type RecordResult struct {
	InviteRecordId int64 `json:"invite_record_id" `
}

type CreateOrderReq struct {
	EnterpriseBase
	InviteRecordId int64  `json:"invite_record_id" binding:"required"`
	ReserveTime    int    `json:"reserve_time" binding:"required"`
	UnionId        string `json:"union_id" binding:"required"`
	ReserveOrderId int64  `json:"reserve_order_id" `
	ReserveEndTime int64  `json:"reserve_end_time"`
}

type CreateOrderRrs struct {
	Id            int64  `json:"id" `
	ReserveNumber string `json:"reserve_number" `
}

type RecordOrderListReq struct {
	EnterpriseBase
	PageReq
	TitleName               string `json:"title_name" comment:"搜索模板关键字"`
	InviteTemplateId        int64  `json:"invite_template_id"`
	ReserveServiceID        int64  `json:"reserve_service_id"`
	MemberName              string `json:"member_name" comment:"搜索导购姓名"`
	UserName                string `json:"user_name" comment:"搜索顾客姓名"`
	ThirdIds                []int  `json:"third_ids" comment:"门店ids"`
	ReserveServiceStartTime int    `json:"reserve_service_start_time" `
	ReserveServiceEndTime   int    `json:"reserve_service_end_time" `
	OrderCreateStartTime    int    `json:"order_create_start_time" `
	OrderCreateEndTime      int    `json:"order_create_end_time" `
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

type TemplateUpdateStatusReq struct {
	EnterpriseBase
	ID   int64 `json:"id"`
	IsUp int   `json:"is_up"`
}
type RecordDetailInfoRes struct {
	RecordDetailInfo
	TimeSlot       []TimeSlot `json:"time_slot"`
	ReserveEndTime int        `json:"reserve_end_time"`
}

type TimeSlot struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}
