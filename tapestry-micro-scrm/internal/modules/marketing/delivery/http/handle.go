package marketinghttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/modules/customer"
	"micro-scrm/internal/modules/enterprise"
	"micro-scrm/internal/modules/marketing"
	"micro-scrm/internal/modules/sms"
	"micro-scrm/internal/modules/user"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/errs"
	"micro-scrm/internal/pkg/sdk/boss"
)

var service *marketing.Service
var userService *user.Service
var smsService *sms.Service
var customerService *customer.Service
var enterpriseService *enterprise.Service

// marketingStrategy 添加营销策略
func marketingStrategy(c *gin.Context) {
	var request marketing.MarketingStrategy
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	response, err := service.Add(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	resp := marketing.MarketingAddResponse{ID: response}
	ginplus.ResSuccess(c, resp)
}

//marketingStrategyInner 创建内置营销策略
func marketingStrategyInner(c *gin.Context) {
	var request marketing.MarketingStrategy
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}

	//内置营销策略初始值
	request.StrategyName = "加粉短信"
	request.StrategyMethod = 1
	request.InnerCondition = 1
	request.CustomerGroupID = 1
	request.Channels = []marketing.ChannelInfo{}
	// 查询短信模板
	templateID, err := smsService.CheckDefaultSmsTemplateByEnterpriseID(ctx.Wrap(c), request.EnterpriseID)
	if err != nil {
		ginplus.ResSuccess(c, err)
		return
	}
	channel := marketing.ChannelInfo{
		ChannelType: 1,
		JuniorID:    templateID,
		Order:       1,
	}
	request.Channels = append(request.Channels, channel)
	request.StrategyPolicy = 1
	request.StrategyType = 2
	request.Status = 1
	request.IsFinished = 1
	response, err := service.Add(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	resp := marketing.MarketingAddResponse{ID: response}
	ginplus.ResSuccess(c, resp)
}

//marketingStrategies 营销策略列表
func marketingStrategies(c *gin.Context) {
	var request marketing.ListRequest
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	service.SmsService = smsService
	total, resp := service.List(ctx.Wrap(c), request)
	res := marketing.MarketingStrategyListResponse{}
	res.List = make([]marketing.MarketingStrategy, 0)
	if resp != nil {
		res.List = *resp
	}
	res.Total = total
	res.Page = request.Page
	res.Size = request.Size
	ginplus.ResSuccess(c, res)
}

// marketingStrategyDetail 营销策略详情
func marketingStrategyDetail(c *gin.Context) {
	var request marketing.MarketingDetailRequest
	service.SmsService = smsService
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	response, err := service.Detail(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

//marketingStrategyUpdate 营销策略更新
func marketingStrategyUpdate(c *gin.Context) {
	var request marketing.MarketingStrategy
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	response, err := service.Update(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// marketingStrategyDataDetail 策略数据详情
func marketingStrategyDataDetail(c *gin.Context) {
	var request marketing.MarketingDataDetailRequest
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	cc := ctx.Wrap(c)
	_, _, hash, _ := enterpriseService.GetInfo(cc, request.EnterpriseId)
	if hash == "" {
		ginplus.ResError(c, errs.EnterpriseNotFund)
		return
	}
	cc = ctx.Wrap(context.WithValue(cc, common.ContextKeyEnum.EnterpriseHash, hash))
	response, err := service.GetMarketingStrategyDetails(cc, request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, response)
}

// marketingStrategyDelete 营销策略删除
func marketingStrategyDelete(c *gin.Context) {
	var request marketing.DeleteRequest
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	if request.ID == 0 {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	err := service.Delete(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nil)
}

// marketingStrategyUpdateStatus 营销策略开始
func marketingStrategyUpdateStatus(c *gin.Context) {
	var request marketing.UpdateStatusRequest
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}

	resp, err := service.UpdateStatus(ctx.Wrap(c), &request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, resp)
}

func SelectorText(c *gin.Context) {
	var selectorText = map[string]map[int8]string{
		"ActivityTypes":          marketing.ActivityTypes,
		"RelationTypes":          marketing.RelationTypes,
		"AssignTypes":            marketing.AssignTypes,
		"TimeTypes":              marketing.TimeTypes,
		"TriggerConditions":      marketing.TriggerConditions,
		"TriggerInnerConditions": marketing.TriggerInnerConditions,
		"StrategyStatus":         marketing.StrategyStatus,
		"StrategyTypes":          marketing.StrategyTypes,
		"StrategyMethod":         marketing.StrategyMethod,
		"StrategyPolicies":       marketing.StrategyPolicies,
		"ChannelTypes":           marketing.ChannelTypes,
	}

	ginplus.ResSuccess(c, selectorText)
}

func MarketTrigExec(c *gin.Context) {
	params := &marketing.MarketTrigExecRequest{}
	if err := c.ShouldBind(params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	cc := ctx.Wrap(c)
	service.SmsService = smsService
	res, err := service.MarketTrigChannelExec(cc, params)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, res)

}

func UpdateStrategyTime(c *gin.Context) {
	var request marketing.UpdateStatusRequest
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	resp, err := service.UpdateStrategyTime(ctx.Wrap(c), &request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, resp)
}

// tagGroupList 标签列表
func tagGroupList(c *gin.Context) {
	var request marketing.UserGroup
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	res, err := service.TagGroupList(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
}

func stest(c *gin.Context) {
	service.SmsService = smsService
	service.CheckGoOnStrategyTemplate(ctx.Wrap(c), 89175222843648, service.SmsService)
}

// getGroupMsgResult 获取发送结果接口
func getGroupMsgResult(c *gin.Context) {
	err := service.GetGroupMsgResult(ctx.Wrap(c))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, struct {
	}{})
}

// 该接口暂时下线，因SQL对CPU营销较大。
func tagGroupUserCount(c *gin.Context) {
	var (
		request marketing.UserCategory
		err     error
		res     *customer.TagUserCount
	)
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	cond := &customer.ListCond{EnterpriseID: request.EnterpriseID}
	if request.TagOrGroupID > 0 {
		cond.CustomerGroupCond.Include = true
		cond.CustomerGroupCond.GroupIDList = []int64{request.TagOrGroupID}
	}

	for _, channel := range request.ChannelTypeList {
		switch channel {
		case marketing.StrategyChannelOfficial:
			cond.CustomerSource = append(cond.CustomerSource, customer.CustomerTypeOpen)
		case marketing.StrategyChannelQW:
			cond.CustomerSource = append(cond.CustomerSource, customer.CustomerTypeExternal)
		case marketing.StrategyChannelShop:
			cond.CustomerSource = append(cond.CustomerSource, customer.CustomerTypeShop)
		case marketing.StrategyChannelMobile:
			cond.CustomerSource = append(cond.CustomerSource, customer.CustomerTypeMobile)
			cond.CustomerSource = append(cond.CustomerSource, customer.CustomerTypeDeliveryMobile)
		}
	}

	request.Type = "group"
	res, err = userService.QueryUserCount(ctx.Wrap(c), cond)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
}

func MarketingTriggerUpdate(c *gin.Context) {
	return
}

//wxChannelList 场景值列表
func wxChannelList(c *gin.Context) {
	resp, _ := service.GetWxChannelList(ctx.Wrap(c))
	ginplus.ResSuccess(c, resp)
}

func GetOriginId(c *gin.Context) {
	request := boss.AuthDetailRequest{}
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	authDetailRes, err := boss.AuthDetail(ctx.Wrap(c), &request)
	if err != nil {
		log.Info(ctx.Wrap(c)).Msgf("GetOriginId获取小程序原始id失败:%v", err)
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, authDetailRes.Data.OriginId)
}

func TwoLevelIds(c *gin.Context) {
	request := marketing.ShortUrlAddReq{}
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	Ids, err := service.TwoLevelIds(ctx.Wrap(c), request)
	if err != nil {
		log.Info(ctx.Wrap(c)).Msgf("TwoLevelIds接口获取二级渠道失败:%v", err)

		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, Ids)
}

func ShortUrlDel(c *gin.Context) {
	request := marketing.ShortUrlDelReq{}

	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}

	err := service.ShortUrlDel(ctx.Wrap(c), request)
	if err != nil {
		log.Info(ctx.Wrap(c)).Msgf("ShortUrlDel接口删除失败:%v", err)

		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, err)
}
func ShortUrlList(c *gin.Context) {
	request := marketing.ShortUrlListReq{}
	res := marketing.ShortUrlListRes{}

	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	log.Info(ctx.Wrap(c)).Msgf("ShortUrlList接口获取短链列表请求参数:request is %v", request)

	shortUrlTotal, resp, err := service.ShortUrlList(ctx.Wrap(c), request)
	if err != nil {
		log.Info(ctx.Wrap(c)).Msgf("ShortUrlList接口获取短链列表失败:shortUrlTotal is %v resp is %v err is %v", shortUrlTotal, resp, err)

	}
	res.ShortCount = shortUrlTotal
	res.ShortList = resp

	ginplus.ResSuccess(c, res)

}

//获取短链详情
func ShortUrlData(c *gin.Context) {
	request := marketing.ShortUrlDataReq{}
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	shortUrlData, err := service.ShortUrlData(ctx.Wrap(c), request)
	if err != nil {
		log.Info(ctx.Wrap(c)).Msgf("ShortUrlData接口获取短链详情失败:%v", err)
		return
	}
	var MiddleConfig marketing.MiddleData
	err = json.Unmarshal([]byte(shortUrlData.ConfigParams), &MiddleConfig)
	shortUrlData.ConfigParam = MiddleConfig
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, shortUrlData)
}

//添加短链
func ShortUrlAdd(c *gin.Context) {
	request := marketing.ShortUrlAddReq{}
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	log.Info(ctx.Wrap(c)).Msgf("ShortUrlAdd接口添加请求参数:%v", request)
	err := service.ShortUrlAdd(ctx.Wrap(c), request)
	if err != nil {
		log.Info(ctx.Wrap(c)).Msgf("ShortUrlAdd接口添加失败，err:%v", err)
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, err)

}

//获取商品详情
func GetSpuDetail(c *gin.Context) {
	request := marketing.SpuDetailReq{}
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}
	urlData, err := url.ParseQuery(request.PageUrl)
	if err != nil {
		ginplus.ResError(c, err)
	}
	log.Info(ctx.Wrap(c)).Msgf("urlData商品pageUrl is:%v", urlData)
	spuId, _ := strconv.ParseInt(urlData["id"][0], 10, 64)
	resp, err := service.GetSpuDetail(ctx.Wrap(c), spuId)
	if err != nil {
		log.Info(ctx.Wrap(c)).Msgf("GetSpuDetail接口获取商品详情失败，err:%v", err)
		return
	}
	var SkuAttribute []marketing.SkuAttribute
	json.Unmarshal([]byte(resp.SkuAttribute), &SkuAttribute)
	resp.SkuAttribute = ""
	for _, v := range SkuAttribute {
		tmpAttr := fmt.Sprintf("%v:%v", v.AttributeName, v.AttributeItemValue)
		if resp.SkuAttribute == "" {
			resp.SkuAttribute += tmpAttr
		} else {
			resp.SkuAttribute += "," + tmpAttr
		}
	}
	resp.SkuPrice = resp.SkuPrice / 100
	ginplus.ResSuccess(c, resp)

}

//获取当前企业的所有短外链
func GetAllShortUrlList(c *gin.Context) {
	request := marketing.ShortUrlListReq{}
	res := marketing.ShortUrlListRes{}

	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, errs.ParamsErr)
		return
	}

	resp, err := service.AllShortUrlList(ctx.Wrap(c), request)
	if err != nil {
		log.Info(ctx.Wrap(c)).Msgf("ShortUrlList接口获取短链列表失败: resp is %v err is %v", resp, err)

	}
	res.ShortCount = int64(len(*resp))
	res.ShortList = resp

	ginplus.ResSuccess(c, res)
}
