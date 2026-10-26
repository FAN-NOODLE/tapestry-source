package work

import (
	"fmt"
	"gw-work-scrm/internal/config"
	"gw-work-scrm/internal/pkg/sdk"
)

var urlEnum = struct {
	UserPoint            string
	UserOrderList        string
	UserInfo             string
	UserGrade            string
	UserOrderListGoods   string
	UserPointsDetailList string
	UserPointsExpire     string
	UserGradeCalculator  string
	AllGrade             string
	UserCouponList       string
	SwCrmUserCouponList  string
	CheckUserSource      string
}{
	//UserOrderList:        "%s/service/v1/data-lake/query-member-order-list",
	//UserInfo:             "%s/service/v1/data-lake/query-member-info",
	//UserPointsDetailList: "%s/service/v1/data-lake/query-member-points-detail-list",
	UserPoint:           "%s/service/v1/data-lake/query-member-point",
	UserGrade:           "%s/service/v1/data-lake/query-member-grade",
	UserOrderListGoods:  "%s/service/v1/data-lake/query-member-order-detail",
	UserPointsExpire:    "%s/service/v1/data-lake/query-member-points-expire",
	UserGradeCalculator: "%s/service/v1/data-lake/query-member-grade-calculator",
	AllGrade:            "%s/service/v1/data-lake/query-tier-list",

	UserInfo:             "%s/service/crm/get-member-info",
	UserOrderList:        "%s/service/crm/get-customer-purchase",
	UserPointsDetailList: "%s/service/crm/get-point-history",
	UserCouponList:       "%s/service/crm/get-customer-coupon",
	SwCrmUserCouponList:  "%s/service/crm/get-sw-customer-coupon",
	CheckUserSource:      "%s/service/v1/member/check-source",
}

// ResponseBase 服务接口返回数据结构体头
type ResponseBase struct {
	sdk.ResponseBase
}

// httpClient http请求代理
//var httpClient = sdk.NewHTTPClient(sdk.ParseJSON{})

// GetURL 获取shop请求链接
func GetURL(uri string) string {
	return fmt.Sprintf(uri, config.CfgData.Service.Work.Domain)
}
