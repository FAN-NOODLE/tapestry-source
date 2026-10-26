package work

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/json"
	"gw-work-scrm/internal/app/http/request/source"
	"gw-work-scrm/internal/pkg/sdk"
	"time"
)

func CrmGetUserOrderList(ctx ctx.Context, req *CrmUserRequest) (*CrmOrderData, error) {
	client := sdk.NewSdkClient(clients.Timeout(time.Second * 15))
	rs, err := client.Post(ctx, GetURL(urlEnum.UserOrderList), clients.JSONBody(req))
	if err != nil {
		return nil, err
	}
	resStr, rErr := rs.Read(clients.Text)
	if rErr != nil {
		return nil, rErr
	}
	var response = CrmOrderResultData{}
	json.Unmarshal([]byte(resStr.(string)), &response)

	if response.Code != 200 {
		return nil, errorx.NewErrorMsg(response.Code, response.Message)
	}
	return &response.Data, nil
}

type CrmOrderResultData struct {
	RequestId string       `json:"request_id"`
	Code      int          `json:"code"`
	Message   string       `json:"message"`
	Error     string       `json:"error"`
	Data      CrmOrderData `json:"data"`
}

type CrmOrderData struct {
	Records     []Records `json:"records"`
	PageNum     int       `json:"pageNum"`
	PageSize    int       `json:"pageSize"`
	IsFirstPage bool      `json:"isFirstPage"`
	IsLastPage  bool      `json:"isLastPage"`
	HasNextPage bool      `json:"hasNextPage"`
	Total       int       `json:"total"`
	Pages       int       `json:"pages"`
}
type Records struct {
	CustomerID   string         `json:"customer_id"`
	PurchaseDate int64          `json:"purchase_date"` // 下单时间
	StoreCode    string         `json:"store_code"`    // store_code
	StoreName    string         `json:"store_name"`    // store_name
	PurchaseID   string         `json:"purchase_id"`
	OrderNo      string         `json:"order_no"` // 订单编号
	Amount       string         `json:"amount"`   // 实际支付总额
	PurchaseSkus []PurchaseSkus `json:"purchase_skus"`

	//Ssociate string `json:"ssociate"`
	//SsociateName string `json:"ssociate_name"`
	//Style string `json:"style"`
	//ItemDescription string `json:"item_description"`
	//Color string `json:"color"`
	//Size string `json:"size"`
	//Quantity string `json:"quantity"`
	//Amount string `json:"amount"`
	//Currency string `json:"currency"`
	//GiftType string `json:"gift_type"`
	//Gift string `json:"gift"`
	//DiscountRate string `json:"discount_rate"`
	//Discount interface{} `json:"discount"`
	//Total interface{} `json:"total"`
	//PayType interface{} `json:"pay_type"`
	//Promotions interface{} `json:"promotions"`
	//Mobile interface{} `json:"mobile"`
	//Membership string `json:"membership"`
	//CounterCode interface{} `json:"counter_code"`
	//OrderType interface{} `json:"order_type"`
	//IdentityKey interface{} `json:"identity_key"`
	//IdentityValue interface{} `json:"identity_value"`
	//CustomerCode interface{} `json:"customer_code"`
}
type PurchaseSkus struct {
	ID       string `json:"id"`
	OrderNo  string `json:"order_no"` // 订单编号
	SkuName  string `json:"sku_name"` // sku_name
	Quantity string `json:"quantity"` // 数量
	Amount   string `json:"amount"`   // 实际支付总额

	//SubOrderNumber interface{} `json:"sub_order_number"`
	//SkuCode string `json:"sku_code"`
	//ReturnFlag string `json:"return_flag"`
	//Currency string `json:"currency"`
	//ReferenceAmount string `json:"reference_amount"`
	//Cost string `json:"cost"`
	//RetailPrice string `json:"retail_price"`
	//DiscountFee interface{} `json:"discount_fee"`
	//OriginalPrice string `json:"original_price"`
	//Remark interface{} `json:"remark"`
	//ProductVoid interface{} `json:"product_void"`
	//ProductVoidRef interface{} `json:"product_void_ref"`
	//DiscountCode string `json:"discount_code"`
	//SkuEvent interface{} `json:"sku_event"`
	//Code string `json:"code"`
	//SubCode interface{} `json:"sub_code"`
	//PurchaseDateTime string `json:"purchase_date_time"`
	//ChannelCode interface{} `json:"channel_code"`
	//CounterCode interface{} `json:"counter_code"`
	//AdvisorCode interface{} `json:"advisor_code"`
	//Color string `json:"color"`
	//Size string `json:"size"`
	//Style string `json:"style"`
	//ItemDescription interface{} `json:"item_description"`
	//Gift interface{} `json:"gift"`
	//Discount string `json:"discount"`
	//StoreCode string `json:"store_code"`
	//StoreName string `json:"store_name"`
	//Ssociate string `json:"ssociate"`
	//SsociateName string `json:"ssociate_name"`
}

func CrmGetUserPointList(ctx ctx.Context, req *CrmPointsListRequest) (*CrmPointData, error) {
	client := sdk.NewSdkClient(clients.Timeout(time.Second * 15))
	rs, err := client.Post(ctx, GetURL(urlEnum.UserPointsDetailList), clients.JSONBody(req))
	if err != nil {
		return nil, err
	}
	resStr, rErr := rs.Read(clients.Text)
	if rErr != nil {
		return nil, rErr
	}
	var response = CrmPointResultData{}
	json.Unmarshal([]byte(resStr.(string)), &response)
	if response.Code != 200 {
		return nil, errorx.NewErrorMsg(response.Code, response.Message)
	}
	return &response.Data, nil
}

type CrmPointResultData struct {
	Code      int          `json:"code"`
	Message   string       `json:"message"`
	RequestID string       `json:"request_id"`
	Data      CrmPointData `json:"data"`
}
type CrmPointData struct {
	Records     []PointRecords `json:"records"`
	PageNum     int            `json:"pageNum"`
	PageSize    int            `json:"pageSize"`
	IsFirstPage bool           `json:"isFirstPage"`
	IsLastPage  bool           `json:"isLastPage"`
	HasNextPage bool           `json:"hasNextPage"`
	Total       int            `json:"total"`
	Pages       int            `json:"pages"`
}
type PointRecords struct {
	ID               string      `json:"id"`
	CustomerID       string      `json:"customer_id"`
	Type             int         `json:"type"`
	Points           int         `json:"points"`
	ChannelID        string      `json:"channel_id"`
	Description      interface{} `json:"description"`
	PointDateTime    string      `json:"point_date_time"`
	PurchaseCode     string      `json:"purchase_code"`
	CounterCode      string      `json:"counter_code"`
	CampaignCode     interface{} `json:"campaign_code"`
	CreationDateTime string      `json:"creation_date_time"`
	ExpiryStartDate  string      `json:"expiry_start_date"`
	ExpiryEndDate    string      `json:"expiry_end_date"`
	CreateUID        interface{} `json:"create_uid"`
}

func CrmGetUserCouponList(ctx ctx.Context, req *CustomerCrmCouponsRequest) (*CrmCouponData, error) {
	client := sdk.NewSdkClient(clients.Timeout(time.Second * 15))
	rs, err := client.Post(ctx, GetURL(urlEnum.UserCouponList), clients.JSONBody(req))
	if err != nil {
		return nil, err
	}
	resStr, rErr := rs.Read(clients.Text)
	if rErr != nil {
		return nil, rErr
	}
	var response = CrmCouponResultData{}
	json.Unmarshal([]byte(resStr.(string)), &response)
	if response.Code != 200 {
		return nil, errorx.NewErrorMsg(response.Code, response.Message)
	}
	return &response.Data, nil
}

type CrmCouponResultData struct {
	Code      int           `json:"code"`
	Message   string        `json:"message"`
	RequestID string        `json:"request_id"`
	Data      CrmCouponData `json:"data"`
}
type CrmCouponData struct {
	Records     []CouponRecords `json:"records"`
	PageNum     int             `json:"pageNum"`
	PageSize    int             `json:"pageSize"`
	IsFirstPage bool            `json:"isFirstPage"`
	IsLastPage  bool            `json:"isLastPage"`
	HasNextPage bool            `json:"hasNextPage"`
	Total       int             `json:"total"`
	Pages       int             `json:"pages"`
}
type CouponRecords struct {
	CouponSerialCode string      `json:"coupon_serial_code"`
	CouponCode       string      `json:"coupon_code"`
	Name             string      `json:"name"`
	Description      string      `json:"description"`
	ChannelType      int         `json:"channel_type"`
	ChannelTypeStr   string      `json:"channel_type_str"`
	Type             int         `json:"type"`
	TypeStr          string      `json:"type_str"`
	StartDate        string      `json:"startDate"`
	EndDate          string      `json:"endDate"`
	Discount         interface{} `json:"discount"`
	FullAmout        interface{} `json:"fullAmout"`
	ReductionAmout   interface{} `json:"reductionAmout"`
	Status           int         `json:"status"`
	StatusStr        string      `json:"status_str"`
}

func GetCrmUserName(ctx ctx.Context, req *UserRequest) (*CrmUserData, error) {
	client := sdk.NewSdkClient(clients.Timeout(time.Second * 15))
	rs, err := client.Post(ctx, GetURL(urlEnum.UserInfo), clients.JSONBody(req))
	if err != nil {
		return nil, err
	}
	var response = CrmUserInfoResultData{}
	_, rErr := rs.Read(clients.JSON(&response))
	if rErr != nil {
		return nil, rErr
	}
	if response.Code != 200 {
		return nil, errorx.NewErrorMsg(response.Code, response.Message)
	}

	return &response.Data, nil
}

type CrmUserInfoResultData struct {
	ID      string      `json:"id"`
	Status  int         `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Error   string      `json:"error"`
	Data    CrmUserData `json:"data"`
}

// CrmUserData crm返回客户数据结构
type CrmUserData struct {
	Status                bool         `json:"status"`
	ID                    string       `json:"user_id"`
	CustomerCode          string       `json:"customer_code"`
	CustomerType          string       `json:"customer_type"`
	CustomerCard          string       `json:"customer_card"`     // 会员卡等类型
	CustomerCardStr       string       `json:"customer_card_str"` // 会员卡等类型
	UnionId               string       `json:"union_id"`
	Mobile                string       `json:"mobile"`
	Name                  string       `json:"name"`
	BirthYear             int          `json:"birth_year"`
	BirthMonth            int          `json:"birth_month"`
	BirthDay              int          `json:"birth_day"`
	Birthday              string       `json:"birthday"`
	ExpiryDate            string       `json:"expiry_date"`
	UpgradeDateTime       string       `json:"upgrade_date_time"`
	UpgradeNeedPoints     float32      `json:"upgrade_need_points"`
	Last12mPurchaseAmount float32      `json:"last_12m_purchase_amount"`
	Last36mPurchaseAmount float32      `json:"last_36m_purchase_amount"`
	TotalValidPoints      float32      `json:"total_valid_points"` // 积分
	ClientType            string       `json:"client_type"`
	Membership            string       `json:"membership"`
	Address               string       `json:"address"`
	Msg                   string       `json:"msg"`
	VipStatus             CrmVipStatus `json:"VipStatus"`
}

type CrmVipStatus struct {
	IssueDate    string `json:"issueDate"`
	IssueShop    string `json:"issueShop"`
	PtsToUpgrade string `json:"ptsToUpgrade"`
	ExpiryDate   string `json:"expiryDate"`
}

type CheckSourceResult struct {
	Code      int             `json:"code"`
	Message   string          `json:"message"`
	RequestId string          `json:"request_id"`
	Data      CheckSourceData `json:"data"`
}
type CheckSourceData struct {
	Auth int `json:"auth"`
}

func CheckSource(ctx ctx.Context, req *source.CheckAuth) (*CheckSourceData, error) {
	client := sdk.NewSdkClient(clients.Timeout(time.Second * 15))
	rs, err := client.Post(ctx, GetURL(urlEnum.CheckUserSource), clients.JSONBody(req))
	if err != nil {
		return nil, err
	}
	response := &CheckSourceResult{}
	_, rErr := rs.Read(clients.JSON(&response))
	if rErr != nil {
		return nil, rErr
	}
	if response.Code != 200 {
		return nil, errorx.NewErrorMsg(response.Code, response.Message)
	}

	return &response.Data, nil
}

func GetGuideList(ctx ctx.Context, req interface{}) (res interface{}, err error) {
	client := sdk.NewSdkClient(clients.Timeout(time.Second * 15))
	rs, err := client.Post(ctx, GetURL("%s/service/v1/store-worker/guide-list"), clients.JSONBody(req))
	if err != nil {
		return res, err
	}
	response := make([]interface{}, 0)
	_, rErr := rs.Read(clients.JSON(&response))
	if rErr != nil {
		return res, rErr
	}

	return response, nil
}

func GetStoreOrDepartmentList(ctx ctx.Context, req interface{}) (res interface{}, err error) {
	client := sdk.NewSdkClient(clients.Timeout(time.Second * 15))
	rs, err := client.Post(ctx, GetURL("%s/service/v1/department/list-by-scrm"), clients.JSONBody(req))
	if err != nil {
		return res, err
	}
	response := make([]interface{}, 0)
	_, rErr := rs.Read(clients.JSON(&response))
	if rErr != nil {
		return res, rErr
	}

	return response, nil
}

func GetFiscalMonth(ctx ctx.Context, req interface{}) (res interface{}, err error) {
	client := sdk.NewSdkClient(clients.Timeout(time.Second * 15))
	rs, err := client.Post(ctx, GetURL("%s/service/v1/data-compass/customer-statistic/fiscal-month"), clients.JSONBody(req))
	if err != nil {
		return res, err
	}
	response := make([]interface{}, 0)
	_, rErr := rs.Read(clients.JSON(&response))
	if rErr != nil {
		return res, rErr
	}

	return response, nil
}

type SwCrmCouponsRequest struct {
	EnterpriseID int64  `json:"enterprise_id"`
	UnionId      string `json:"union_id"`
}
type SwCrmCouponResultData struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	RequestID string      `json:"request_id"`
	Data      interface{} `json:"data"`
}

func SwUserCrmGetUserCouponList(ctx ctx.Context, req *SwCrmCouponsRequest) (res interface{}, err error) {
	client := sdk.NewSdkClient(clients.Timeout(time.Second * 15))
	rs, err := client.Post(ctx, GetURL(urlEnum.SwCrmUserCouponList), clients.JSONBody(req))
	if err != nil {
		return nil, err
	}
	var resp SwCrmCouponResultData

	_, rErr := rs.Read(clients.JSON(&resp))
	if rErr != nil {
		return res, rErr
	}
	return resp.Data, nil
}
