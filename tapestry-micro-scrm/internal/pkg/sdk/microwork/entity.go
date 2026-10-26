package microwork

// SendMiniMessageRequest 发送智汇导购应用消息
type SendMiniMessageRequest struct {
	EnterpriseID int64               `json:"enterprise_id"`
	MemUserIds   []string            `json:"mem_user_ids"`
	Title        string              `json:"title"`
	Content      []map[string]string `json:"content"`
	Page         string              `json:"page"`
}

// SendMiniMessageResponse 根据id查询微信优惠券信息返回数据格式
type SendMiniMessageResponse struct {
	ResponseBase
}

type AddGuideRequest struct {
	Type         int64                   `json:"type"` // 1、门店导购，2、社群导购，3、企业导购  4.企业管理 5、分销合伙人
	EnterpriseID int64                   `json:"enterprise_id"`
	Member       []AddGuideMemberRequest `json:"member"`
	OperatorId   int                     `json:"operator_id"`
	OperatorName string                  `json:"operator_name"`
}
type AddGuideMemberRequest struct {
	MemUid   string `json:"mem_uid"`
	UserName string `json:"user_name"`
	Mobile   string `json:"mobile"`
	UserId   string `json:"user_id"`
	Avatar   string `json:"avatar"`
}

type AddGuideResponse struct {
	ResponseBase
}

type UserInfoRequest struct {
	EnterpriseID int64  `json:"enterprise_id"`
	UnionId      string `json:"union_id"`
	CustomerCode string `json:"customer_code"`
	Membership   string `json:"membership"`
}

type UserInfoResponse struct {
	ResponseBase
	Data MemberInfoResp `json:"data"`
}

type MemberInfoResp struct {
	Status                bool    `json:"status"`
	Id                    string  `json:"id"`
	CustomerCode          string  `json:"customer_code"`
	CustomerLevel         string  `json:"customer_level"`
	CustomerCard          string  `json:"customer_card"`
	CustomerCardStr       string  `json:"customer_card_str"`
	CustomerType          string  `json:"customer_type"`
	UnionId               string  `json:"union_id"`
	Mobile                string  `json:"mobile"`
	Name                  string  `json:"name"`
	BirthYear             int     `json:"birth_year"`
	BirthMonth            int     `json:"birth_month"`
	BirthDay              int     `json:"birth_day"`
	Birthday              string  `json:"birthday"`
	ExpiryDate            string  `json:"expiry_date"`
	UpgradeDateTime       string  `json:"upgrade_date_time"`
	UpgradeNeedPoints     float32 `json:"upgrade_need_points"`
	Last12mPurchaseAmount float32 `json:"last_12m_purchase_amount"`
	Last36mPurchaseAmount float32 `json:"last_36m_purchase_amount"`
	Membership            string  `json:"membership"`
	ClientType            string  `json:"client_type"`
}

type CheckMemberResponse struct {
	ResponseBase
	Data CheckMemberDataResponse `json:"data"`
}
type CheckMemberDataResponse struct {
	Status bool                          `json:"status"`
	Config CheckMemberDataConfigResponse `json:"config"`
}
type CheckMemberDataConfigResponse struct {
	AppId  string `json:"app_id"`
	Path   string `json:"path"`
	Source string `json:"source"`
}

type WelfareHavedRequest struct {
	EnterpriseID int64  `json:"enterprise_id"`
	UnionId      string `json:"union_id"`
	ExtUserid    string `json:"extUserid"`
	ShopCode     string `json:"shopCode"`
	ClerkCode    string `json:"clerkCode"`
	StoreType    int32  `json:"storeType"`
}

type WelfareHavedResponse struct {
	ResponseBase
	Data WelfareHavedDataResponse `json:"data"`
}
type WelfareHavedDataResponse struct {
	Haved bool `json:"haved"`
}

// CustomerPurchaseReq crm查询客户订单接口
type CustomerPurchaseReq struct {
	UnionId       string `json:"union_id"`
	EnterpriseId  int64  `json:"enterprise_id"`
	StartDateTime int    `json:"startDateTime"`
	EndDateTime   int64  `json:"endDateTime"`
}

type CustomerPurchaseResp struct {
	ResponseBase
	Data CustomerPurchaseDataResp `json:"data"`
}

type CustomerPurchaseDataResp struct {
	Total    int                               `json:"total"`
	PageNum  int                               `json:"pageNum"`
	PageSize int                               `json:"pageSize"`
	Records  []CustomerPurchaseDataRecordsResp `json:"records"`
}

type CustomerPurchaseDataRecordsResp struct {
	StoreName    string                                        `json:"store_name"`
	OrderNo      string                                        `json:"order_no"`
	PurchaseDate int64                                         `json:"purchase_date"`
	PurchaseSkus []CustomerPurchaseDataRecordsPurchaseSkusResp `json:"purchase_skus"`
}

type CustomerPurchaseDataRecordsPurchaseSkusResp struct {
	StoreName string `json:"store_name"`
	OrderNo   string `json:"order_no"`
	SkuName   string `json:"sku_name"`
	Quantity  string `json:"quantity"`
	Amount    string `json:"amount"`
}

type SubscribeMessageSendReq struct {
	EnterpriseId int64                       `json:"enterprise_id"`
	OpenId       string                      `json:"open_id"`
	Type         int                         `json:"type" desc:"业务类型 1到店服务提醒 2服务变更提醒"`
	Data         SubscribeMessageSendDataReq `json:"data"`
}

type SubscribeMessageSendDataReq struct {
	ReserveOrderId int64  `json:"reserve_order_id" desc:"预约表主键id"`
	ServiceName    string `json:"service_name" desc:"服务名称"`
	StoreName      string `json:"store_name" desc:"门店名称"`
	ReserveTime    int    `json:"reserve_time" desc:"到店时间"`
	ReserveEndTime int    `json:"reserve_end_time" desc:"到店时间"`
}

type SubscribeMessageSendResp struct {
	ResponseBase
}

type SaAssignmentReq struct {
	EnterpriseID int64                 `json:"enterprise_id"`
	Params       SaAssignmentParamsReq `json:"params"`
}

type SaAssignmentParamsReq struct {
	SaAccount      string `json:"saAccount"`
	FirstBindDate  string `json:"firstbinddate"`
	LatestBindDate string `json:"latestbinddate"`
	UnionId        string `json:"unionId"`
	StoreCode      string `json:"storecode"`
	Source         string `json:"source" comment:"SW的接口请求增加"`
	BindTime       string `json:"bindTime" comment:"SW的接口请求增加"`
}

type SaUnAssignmentReq struct {
	EnterpriseID int64                   `json:"enterprise_id"`
	Params       SaUnAssignmentParamsReq `json:"params"`
}
type SaUnAssignmentParamsReq struct {
	SaAccount  string `json:"saAccount"`
	UnionId    string `json:"unionId"`
	UnbindTime string `json:"unbindTime"`
}

type ThirdCdpResp struct {
	ResponseCode string `json:"responseCode"`
	ResponseDesc string `json:"responseDesc"`
}

type SaAssignmentResp struct {
	ResponseBase
	Data ThirdCdpResp `json:"data"`
}

type SendMaResultReq struct {
	EnterpriseID int64              `json:"enterprise_id"`
	TaskId       int64              `json:"task_id"`
	CrowdId      string             `json:"crowd_id"`
	TaskCategory int                `json:"task_category"`
	CrowdResult  []SendMaResultData `json:"crowd_result"`
}

type SendMaResultData struct {
	Date              string `json:"date"`
	SaName            string `json:"sa_name"`
	SaCode            string `json:"sa_code"`
	FlowId            string `json:"flowId"`
	Campaign          string `json:"campaign"`
	CPlanTemplateId   string `json:"c_plan_template_id"`
	CPlanTemplateName string `json:"c_plan_template_name"`
	CTaskRemark       string `json:"c_task_remark"`
	CTaskClassifation string `json:"c_task_classification"`
	Unionid           string `json:"unionid"`
	Result            string `json:"result"`
	ResultMessage     string `json:"result_message"`
	TaskRemark        string `json:"task_remark"`
}

type SendMaMomentResultReq struct {
	EnterpriseID        int64    `json:"enterprise_id"`
	SaCode              string   `json:"sa_code"`
	TaskId              int64    `json:"task_id"`
	CrowdId             string   `json:"crowd_id"`
	Campaign            string   `json:"campaign"`
	CustomerUnionidList []string `json:"customer_unionid_list"`
}

type SendMaMomentCommentsReq struct {
	EnterpriseID int64                      `json:"enterprise_id"`
	SaCode       string                     `json:"sa_code"`
	TaskId       int64                      `json:"task_id"`
	CrowdId      string                     `json:"crowd_id"`
	Campaign     string                     `json:"campaign"`
	CommentList  []SendMaMomentCommentsInfo `json:"comment_list"`
	LikeList     []SendMaMomentCommentsInfo `json:"like_list"`
}

type SendMaMomentCommentsInfo struct {
	UnionId   string `json:"unionid"`
	CreateTme int64  `json:"create_time"`
}

type SendMaFinishReq struct {
	EnterpriseID int64         `json:"enterprise_id"`
	OriginType   int           `json:"origin_type"`
	Events       []MaEvents    `json:"events"`
	CommonProps  MACommonProps `json:"commonProps"`
}
type MaEvents struct {
	IdentityType      string `json:"identityType"`
	IdentityValue     string `json:"identityValue"`
	Date              string `json:"date"`
	CResult           string `json:"c_result"`
	CSaName           string `json:"c_sa_name"`
	CSaCode           string `json:"c_sa_code"`
	CTaskId           string `json:"c_task_id"`
	CTaskName         string `json:"c_task_name"`
	CCrowdId          string `json:"c_crowd_id"`
	FlowId            string `json:"flowId"`
	Campaign          string `json:"campaign"`
	CPlanTemplateId   string `json:"c_plan_template_id"`
	CPlanTemplateName string `json:"c_plan_template_name"`
	CTaskRemark       string `json:"c_task_remark"`
	CTaskClassifation string `json:"c_task_classification"`
}
type MACommonProps struct {
	Event string `json:"event" comment:"c_sa_completed_task_result"`
}

type ItemSetReq struct {
	EnterpriseId int64  `json:"enterprise_id"`
	ItemType     string `json:"item_type"`
	ItemId       string `json:"item_id"`
	Properties   any    `json:"properties"`
}

type ProfileSetReq struct {
	EnterpriseId int64  `json:"enterprise_id"`
	DistinctId   string `json:"distinct_id"`
	IsLoginId    bool   `json:"is_login_id"`
	Profiles     any    `json:"profiles"`
}

type TrackReq struct {
	EnterpriseId int64  `json:"enterprise_id"`
	DistinctId   string `json:"distinct_id"`
	EventName    string `json:"event_name"`
	IsLoginId    bool   `json:"is_login_id"`
	Profiles     any    `json:"profiles"`
}

type SunQRCodeReq struct {
	EnterpriseID int64  `json:"enterprise_id"`
	Path         string `json:"path"`
	Scene        string `json:"scene"`
	EnvVersion   string `json:"env_version"`
}

type SunQRCodeResp struct {
	ResponseBase
	CodeUrl string `json:"code_url"`
	Link    string `json:"link"`
}

type SaAssignmentEcHubReq struct {
	EnterpriseID int64                      `json:"enterprise_id"`
	Params       SaAssignmentParamsEcHubReq `json:"params"`
}

// SaAssignmentParamsEcHubReq echub的sa绑定
type SaAssignmentParamsEcHubReq struct {
	SaAccount      string `json:"saAccount"`
	FirstBindDate  string `json:"firstbinddate"`
	LatestBindDate string `json:"latestbinddate"`
	UnionId        string `json:"unionId"`
	Nickname       string `json:"nickname"`
	AvatarUrl      string `json:"avatarUrl"`
	Channel        string `json:"channel"`
}

type SaUnAssignmentEcHubReq struct {
	EnterpriseID int64                        `json:"enterprise_id"`
	Params       SaUnAssignmentParamsEcHubReq `json:"params"`
}

type SaUnAssignmentParamsEcHubReq struct {
	SaAccount  string `json:"saAccount"`
	UnionId    string `json:"unionId"`
	UnbindDate string `json:"unbinddate"`
}
type EcHubResp struct {
	Msg        string `json:"msg"`
	StatusCode string `json:"statusCode"`
}
type SaAssignmentEcHubResp struct {
	ResponseBase
	Data EcHubResp `json:"data"`
}
