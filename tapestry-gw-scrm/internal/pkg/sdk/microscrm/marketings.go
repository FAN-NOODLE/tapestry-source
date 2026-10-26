package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/modules/marketing"
	"gw-scrm/internal/pkg/sdk"
)

//GetSelectorText 获取选择项
func GetSelectorText(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.SelectorTextUrl), request)
}

//CreateMarketing 创建营销策略
func CreateMarketing(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.MarketingAddUrl), request)
}

//MarketingList 营销策略列表
func MarketingList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.MarketingListUrl), request)
}

//MarketingDetail 营销策略详情
func MarketingDetail(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.MarketingDetailUrl), request)
}

//UpdateMarketing 更新营销策略
func UpdateMarketing(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.MarketingUpdateUrl), request)
}

//DeleteMarketing 删除营销策略
func DeleteMarketing(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.MarketingDeleteUrl), request)
}

//MarketingStrategyUpdateStatus 营销策略更新状态
func MarketingStrategyUpdateStatus(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.MarketingUpdateStatusUrl), request)
}

//MarketingStrategyDataDetail 营销策略数据详情
func MarketingStrategyDataDetail(c ctx.Context, request interface{}) (*marketing.StrategyDataDetails, error) {
	//return sdk.SendPost(c, GetUrl(urlEnum.MarketingStrategyDataDetailUrl), request)
	var response marketing.StrategyDataDetails
	err := client.Post(c, GetUrl(urlEnum.MarketingStrategyDataDetailUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// 营销策略触发
func MarketingStrategyTrigger(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.MarketingStrategyTriggerUrl), request)
}
func UpdateStrategyDataTime(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.UpdateSettingDayUrl), request)
}

// 营销策略触发
func MarketingUserGroupList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.MarketingUserGroupListUrl), request)
}

// TagGroupUserCount 标签中用户总数
func TagGroupUserCount(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.TagGroupUserCountUrl), request)
}

// GetSpuDetail 获取商品详情
func GetSpuDetail(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.GetSpuDetail), request)
}

// GetTwoLevelIds 获取店铺的二级渠道id
func GetTwoLevelIds(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.GetTwoLevelIds), request)
}

// GetShortUrlData 获取短链详情
func GetShortUrlData(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.ShortUrlData), request)
}

// GetShortUrlList 获取短链列表
func GetShortUrlList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.GetShortUrlList), request)
}

// ShortUrlDel 删除短链
func ShortUrlDel(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.ShortUrlDel), request)
}

//短链添加
func ShortUrlAdd(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.ShortUrlAdd), request)
}

func GetOriginId(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.GetOriginId), request)
}

// 券是否使用中
func GouponIsUsed(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.GouponIsUsedUrl), request)
}

// 券匹配
func CouponMatch(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.CouponMatchUrl), request)
}

//MarketingDetailCoupon 营销策略自动发券详情
func MarketingDataCoupon(c ctx.Context, request interface{}) (*marketing.StrategyCouponDataDetails, error) {
	var response marketing.StrategyCouponDataDetails
	err := client.Post(c, GetUrl(urlEnum.MarketingDataCouponUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetShortUrlList 获取短链列表
func GetAllShortUrlList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.GetAllShortUrlListUrl), request)
}


func GetTemplateList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.GetTemplateList), request)
}

func CreateTemplateInfo(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.CreateTemplateInfo), request)
}

func DelTemplateInfo(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.DelTemplateInfo), request)
}
func CheckTemplateInfo(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.CheckTemplateInfo), request)
}
func GetTemplateInfo(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.GetTemplateInfo), request)
}
