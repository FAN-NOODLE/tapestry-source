package work

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/json"
	"gw-work-scrm/internal/pkg/sdk"
	"time"
)

type ResultData struct {
	ID      string `json:"id"`
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
	Data    GoData `json:"data"`
}

type GoPoints struct {
	UsedPoint    int    `json:"usedPoint"`
	EvalidPoints int64  `json:"evalidPoints"`
	PointType    string `json:"pointType"`
	SumPoint     int    `json:"sumPoint"`
}

type GoData struct {
	UnionCode string     `json:"union_code"`
	Points    []GoPoints `json:"points"`
}

// Enterprise Boss里企业信息返回结果
type UserPoint struct {
	UsedPoint    int64 `json:"used_point"`
	EvalidPoints int64 `json:"evalid_points"`
	PointType    int64 `json:"point_type"`
	SumPoint     int64 `json:"sum_point"`
}

// UserRequest 用户信息请求
type UserRequest struct {
	EnterpriseID       int64  `json:"enterprise_id"`
	Type               string `json:"type"`
	Value              string `json:"value"`
	TradeNo            string `json:"trade_no"`
	Tradeno            string `json:"tradeNo"`
	PointTypeGroup     string `json:"pointTypeGroup"`
	Grade              string `json:"grade"`
	UpGrade            string `json:"upGrade"`
	Startdate          string `json:"Startdate"`
	Enddate            string `json:"Enddate"`
	PageNo             string `json:"pageNo"`
	Limit              string `json:"limit"`
	UnionId            string `json:"union_id"`
	EuId               int64  `json:"eu_id"`
	WorkExternalUserID int64  `json:"work_external_user_id" db:"work_external_user_id"`
}

// CrmUserRequest CRM 用户信息请求
type CrmUserRequest struct {
	EnterpriseID   int64  `json:"enterprise_id"`
	PointTypeGroup string `json:"pointTypeGroup"`
	Grade          string `json:"grade"`
	UpGrade        string `json:"upGrade"`
	Startdate      string `json:"Startdate"`
	Enddate        string `json:"Enddate"`
	PageNum        int    `json:"page_num"`
	PageSize       int    `json:"page_size"`
	UnionId        string `json:"union_id"`
	StartDateTime  int    `json:"startDateTime"`
	EndDateTime    int    `json:"endDateTime"`
}

type PointsListRequest struct {
	EnterpriseID   int64  `json:"enterprise_id"`
	Type           string `json:"type"`
	Value          string `json:"value"`
	PointTypeGroup string `json:"pointTypeGroup"`
	Startdate      string `json:"startDate"`
	Enddate        string `json:"endDate"`
	PageNo         string `json:"pageNo"`
	Limit          string `json:"limit"`
}

type CrmPointsListRequest struct {
	EnterpriseID   int64  `json:"enterprise_id"`
	UnionId        string `json:"union_id"`
	PointTypeGroup string `json:"pointTypeGroup"`
	StartTime      string `json:"startTime"`
	EndTime        string `json:"endTime"`
	PageNum        int    `json:"pageNum"`
	PageSize       int    `json:"pageSize"`
}

type CustomerCrmCouponsRequest struct {
	EnterpriseID int64  `json:"enterprise_id"`
	UnionId      string `json:"union_id"`
	ChannelType  string `json:"channel_type"`
	Status       int    `json:"status"`
	Source       int    `json:"source"`
	PageNum      int    `json:"pageNum"`
	PageSize     int    `json:"pageSize"`
}

// UserPoint 查询用户积分信息
func GetUserPoint(ctx ctx.Context, req *UserRequest) ([]GoPoints, error) {
	client := sdk.NewSdkClient(clients.Timeout(time.Second * 15))
	rs, err := client.Post(ctx, GetURL(urlEnum.UserPoint), clients.JSONBody(req))
	if err != nil {
		return nil, err
	}
	resStr, rErr := rs.Read(clients.Text)
	if rErr != nil {
		return nil, rErr
	}
	var response = ResultData{}
	jErr := json.Unmarshal([]byte(resStr.(string)), &response)
	if jErr != nil {
		return nil, jErr
	}
	if response.Status != 200 {
		return nil, errorx.NewErrorMsg(response.Status, response.Message)
	}
	return response.Data.Points, nil
}

// GetAllGrade 所有会员等级
func GetAllGrade(ctx ctx.Context, req *UserRequest) ([]AllGrade, error) {
	client := sdk.NewSdkClient(clients.Timeout(time.Second * 15))
	rs, err := client.Post(ctx, GetURL(urlEnum.AllGrade), clients.JSONBody(req))
	if err != nil {
		return nil, err
	}
	resStr, rErr := rs.Read(clients.Text)
	if rErr != nil {
		return nil, rErr
	}
	var response = RepAllGradeData{}
	jErr := json.Unmarshal([]byte(resStr.(string)), &response)
	if jErr != nil {
		return nil, jErr
	}
	if response.Status != 200 {
		return nil, errorx.NewErrorMsg(response.Status, response.Message)
	}
	return response.Data, nil
}

type PointListResultData struct {
	ID      string        `json:"id"`
	Status  int           `json:"status"`
	Message string        `json:"message"`
	Error   string        `json:"error"`
	Data    PointListData `json:"data"`
}
type PointListData struct {
	Details []PointList `json:"details"`
}

type PointList struct {
	TradeNo        string `json:"tradeNo"`
	PointType      string `json:"pointType"`
	IsAdd          int    `json:"is_add"`
	Points         int64  `json:"points"`
	SourceCode     string `json:"sourceCode"`
	CreateDateTime string `json:"createDateTime"`
	SortDateTime   int64  `json:"sortDateTime"`
}

type OrderResultData struct {
	ID      string    `json:"id"`
	Status  int       `json:"status"`
	Message string    `json:"message"`
	Error   string    `json:"error"`
	Data    OrderData `json:"data"`
}
type Orders struct {
	OrderType        int    `json:"orderType"`
	Consignee        string `json:"consignee"`
	TradeNo          string `json:"tradeNo"`
	City             string `json:"city"`
	EmployeeNo       string `json:"employee_no"`
	Remark           string `json:"remark"`
	Source           int    `json:"source"`
	Type             int    `json:"type"`
	LogisticsInfo    string `json:"logisticsInfo"`
	TradeTime        string `json:"tradeTime"`
	StoreNo          string `json:"storeNo"`
	DeliverMethod    string `json:"deliverMethod"`
	Province         string `json:"province"`
	Phone            string `json:"phone"`
	LogisticsCompany string `json:"logisticsCompany"`
	StoreName        string `json:"store_name"`
	LogisticsNum     string `json:"logisticsNum"`
	Addr             string `json:"addr"`
	TradePrice       int64  `json:"tradePrice"`
	//Goods            []OrdersGoods `json:"goods"`
	//Goods []interface{} `json:"goods"`
}
type OrderData struct {
	Total  int      `json:"total"`
	Orders []Orders `json:"orders"`
}

func GetUserOrderList(ctx ctx.Context, req *UserRequest) ([]Orders, error) {
	client := sdk.NewSdkClient(clients.Timeout(time.Second * 15))
	rs, err := client.Post(ctx, GetURL(urlEnum.UserOrderList), clients.JSONBody(req))
	if err != nil {
		return nil, err
	}
	resStr, rErr := rs.Read(clients.Text)
	if rErr != nil {
		return nil, rErr
	}
	var response = OrderResultData{}
	json.Unmarshal([]byte(resStr.(string)), &response)
	if response.Status != 200 {
		return nil, errorx.NewErrorMsg(response.Status, response.Message)
	}
	return response.Data.Orders, nil
}

func GetUserPointsDetailList(ctx ctx.Context, req PointsListRequest) (*[]PointList, error) {
	client := sdk.NewSdkClient(clients.Timeout(time.Second * 15))
	rs, err := client.Post(ctx, GetURL(urlEnum.UserPointsDetailList), clients.JSONBody(req))
	if err != nil {
		return nil, err
	}
	resStr, rErr := rs.Read(clients.Text)
	if rErr != nil {
		return nil, rErr
	}
	var response = PointListResultData{}
	json.Unmarshal([]byte(resStr.(string)), &response)
	if response.Status != 200 {
		return nil, errorx.NewErrorMsg(response.Status, response.Message)
	}
	return &response.Data.Details, nil
}

type OrderGoodsResultData struct {
	ID      string `json:"id"`
	Code    string `json:"code"`
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
	Data    Data   `json:"data"`
}
type OrdersGoods struct {
	TradeNo      string `json:"tradeNo"`
	TradePrice   int64  `json:"tradePrice"`
	Amt          int    `json:"amt"`
	ProductCount int    `json:"productCount"`
	SkuCode      string `json:"skuCode"`
	ProductName  string `json:"productName"`
}
type Data struct {
	OrderDetails []OrdersGoods `json:"orderDetails"`
}

func GetUserOrderListGoods(ctx ctx.Context, req *UserRequest) ([]OrdersGoods, error) {
	//rs, err := http.ClientObject.Post(ctx, GetURL(urlEnum.UserOrderListGoods), clients.JSONBody(req))
	client := sdk.NewSdkClient(clients.Timeout(time.Second * 15))
	rs, err := client.Post(ctx, GetURL(urlEnum.UserOrderListGoods), clients.JSONBody(req))
	if err != nil {
		return nil, err
	}
	resStr, rErr := rs.Read(clients.Text)
	if rErr != nil {
		return nil, rErr
	}
	var response = OrderGoodsResultData{}
	json.Unmarshal([]byte(resStr.(string)), &response)
	if response.Status != 200 {
		return nil, errorx.NewErrorMsg(response.Status, response.Message)
	}
	return response.Data.OrderDetails, nil
}

type GradeResultData struct {
	ID string `json:"id"`
	//Code    string    `json:"code"`
	Status  int       `json:"status"`
	Message string    `json:"message"`
	Error   string    `json:"error"`
	Data    GradeData `json:"data"`
}
type GradeData struct {
	IsUpdateGrade  string `json:"isUpdateGrade"`
	Grade          string `json:"grade"`
	GradeEndDate   string `json:"gradeEndDate"`
	GradeCode      string `json:"gradeCode"`
	GradeStartDate string `json:"gradeStartDate"`
}

func GetUserGrade(ctx ctx.Context, req *UserRequest) (*GradeData, error) {
	client := sdk.NewSdkClient(clients.Timeout(time.Second * 15))
	rs, err := client.Post(ctx, GetURL(urlEnum.UserGrade), clients.JSONBody(req))
	if err != nil {
		return nil, err
	}
	resStr, rErr := rs.Read(clients.Text)
	if rErr != nil {
		return nil, rErr
	}
	var response = GradeResultData{}
	jErr := json.Unmarshal([]byte(resStr.(string)), &response)
	if jErr != nil {
		return nil, jErr
	}
	if response.Status != 200 {
		return nil, errorx.NewErrorMsg(response.Status, response.Message)
	}
	return &response.Data, nil
}

//过期积分
func GetUserPointsExpire(ctx ctx.Context, req *UserRequest) (*UserPointExpire, error) {
	client := sdk.NewSdkClient(clients.Timeout(time.Second * 15))
	rs, err := client.Post(ctx, GetURL(urlEnum.UserPointsExpire), clients.JSONBody(req))
	if err != nil {
		return nil, err
	}
	resStr, rErr := rs.Read(clients.Text)
	if rErr != nil {
		return nil, rErr
	}
	var response = RepUserPointExpireData{}
	json.Unmarshal([]byte(resStr.(string)), &response)
	//if jErr != nil {
	//	return nil, jErr
	//}
	if response.Status != 200 {
		return nil, errorx.NewErrorMsg(response.Status, response.Message)
	}
	return &response.Data, nil
}

//用户升级信息
func GetUserGradeCalculator(ctx ctx.Context, req *UserRequest) (*UserGradeCalculator, error) {
	client := sdk.NewSdkClient(clients.Timeout(time.Second * 15))
	rs, err := client.Post(ctx, GetURL(urlEnum.UserGradeCalculator), clients.JSONBody(req))
	if err != nil {
		return nil, err
	}
	resStr, rErr := rs.Read(clients.Text)
	if rErr != nil {
		return nil, rErr
	}
	var response = RepUserGradeCalculatorData{}
	json.Unmarshal([]byte(resStr.(string)), &response)
	if response.Status != 200 {
		return nil, errorx.NewErrorMsg(response.Status, response.Message)
	}
	return &response.Data, nil
}

type UserInfoResultData struct {
	ID string `json:"id"`
	//Code    string   `json:"code"`
	Status  int      `json:"status"`
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Error   string   `json:"error"`
	Data    UserData `json:"data"`
}

//type UserData struct {
//	Name     string `json:"name"`
//	Mobile   string `json:"mobile"`
//	Birthday string `json:"birthday"`
//	Status   string `json:"status"`
//	Cardno   string `json:"cardno"`
//	Address  []UserAddress
//}

// UserData crm返回客户数据结构
type UserData struct {
	Status                bool    `json:"status"`
	ID                    string  `json:"user_id" db:"user_id"`
	CustomerCode          string  `json:"customer_code" db:"customer_code"`
	CustomerType          string  `json:"customer_type" db:"customer_type"`
	UnionId               string  `json:"union_id" db:"union_id"`
	Mobile                string  `json:"mobile" db:"mobile"`
	Name                  string  `json:"name" db:"name"`
	BirthYear             int     `json:"birth_year" db:"birth_year"`
	BirthMonth            int     `json:"birth_month" db:"birth_month"`
	BirthDay              int     `json:"birth_day" db:"birth_day"`
	Birthday              string  `json:"birthday" db:"birthday"`
	ExpiryDate            string  `json:"expiry_date" db:"expiry_date"`
	UpgradeDateTime       string  `json:"upgrade_date_time" db:"upgrade_date_time"`
	UpgradeNeedPoints     int     `json:"upgrade_need_points" db:"upgrade_need_points"`
	Last36mPurchaseAmount float32 `json:"last_36m_purchase_amount" db:"last_36m_purchase_amount"`
	Membership            string  `json:"membership" db:"membership"`
	Msg                   string  `json:"msg" db:"msg"`
}

type UserAddress struct {
	IsDefault string `json:"isDefault"`
	Location  string `json:"location"`
}

func GetUserName(ctx ctx.Context, req *UserRequest) (*UserData, error) {
	//rs, err := http.ClientObject.Post(ctx, GetURL(urlEnum.UserInfo), clients.JSONBody(req))
	client := sdk.NewSdkClient(clients.Timeout(time.Second * 15))
	rs, err := client.Post(ctx, GetURL(urlEnum.UserInfo), clients.JSONBody(req))
	if err != nil {
		return nil, err
	}
	resStr, rErr := rs.Read(clients.Text)
	if rErr != nil {
		return nil, rErr
	}
	var response = UserInfoResultData{}
	jErr := json.Unmarshal([]byte(resStr.(string)), &response)
	fmt.Println(response)
	fmt.Println(response.Code)
	if jErr != nil {
		return nil, jErr
	}
	if response.Code != 200 {
		fmt.Println(313131313)
		return nil, errorx.NewErrorMsg(response.Code, response.Message)
	}
	fmt.Println(response)

	return &response.Data, nil
}
