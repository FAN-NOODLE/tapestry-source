package cdp

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/modules/index"
	"gw-scrm/internal/modules/marketing"
	"strconv"
)

// FansTargetQuery 公众号指标查询
func FansTargetQuery(c ctx.Context, fans index.ShopFeeParam, body index.FansBody) (interface{}, error) {
	uri := fmt.Sprintf(urlEnum.FansTargetQuery, "%s", strconv.FormatInt(fans.EnterpriseId, 10), strconv.FormatInt(fans.Day, 10))
	response, err := SendPost(c, GetUrl(uri), body)
	return response, err
}

// 导购列表查询
func ShopperTargetQuery(c ctx.Context, shopper index.ShopperListParam, body index.ShopperListBody) (interface{}, error) {
	uri := fmt.Sprintf(urlEnum.ShopperListTarget, "%s", strconv.FormatInt(shopper.EnterpriseId, 10), strconv.FormatInt(shopper.Day, 10))
	response, err := SendPost(c, GetUrl(uri), body)
	return response, err
}

// 综合指标 - 商城变现、客户总转化、企微好友、群成员
func CompositeIndicator(c ctx.Context, shopFee index.ShopFeeParam, body index.ShopFeeBody) (interface{}, error) {
	uri := fmt.Sprintf(urlEnum.ShopFeeTarget, "%s", strconv.FormatInt(shopFee.EnterpriseId, 10), strconv.FormatInt(shopFee.Day, 10), strconv.FormatInt(shopFee.Type, 10))
	response, err := SendPost(c, GetUrl(uri), body)
	return response, err
}

// 获取自定义渠道列表"%s/apis/%s/chanList/%s"
func ChanListQuery(c ctx.Context, chanListP marketing.ChanListRequest, body marketing.ChanListRequestBody) (interface{}, error) {
	uri := fmt.Sprintf(urlEnum.ChanListQuery, "%s", strconv.FormatInt(chanListP.EnterpriseID, 10), strconv.FormatInt(chanListP.Pid, 10))
	response, err := SendPost(c, GetUrl(uri), body)
	return response, err
}
