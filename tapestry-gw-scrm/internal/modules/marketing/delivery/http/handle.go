package marketinghttp

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gw-scrm/internal/modules/marketing"
	"gw-scrm/internal/pkg/errs"
	"gw-scrm/internal/pkg/http"
	"gw-scrm/internal/pkg/sdk/cdp"
	"gw-scrm/internal/pkg/sdk/gwshop"
	"gw-scrm/internal/pkg/sdk/microscrm"
	"gw-scrm/internal/pkg/sdk/microshop"
	"gw-scrm/internal/pkg/sdk/shorturl"
	"gw-scrm/internal/pkg/sdk/wxapi"
	"gw-scrm/internal/pkg/validate_params"
	"gw-scrm/internal/pkg/ws_file"
	"gw-scrm/pkg/microactivity"
	"io"
	"strconv"
	"strings"
	"time"
)

//userGroups 用户分组
func userGroups(c *gin.Context) {
	request := marketing.UserGroupRequest{}
	tokenInfo := http.GetTokenInfo(c)
	request.EnterpriseID = tokenInfo.EnterpriseId
	response, err := microscrm.UserGroupList(ctx.Wrap(c), nil)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	bytes, _ := json.Marshal(response)
	var data marketing.GroupList
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	arr := []marketing.UserGroup{}
	for _, aiplVal := range data.Aipl {
		arr = append(arr, aiplVal)
	}
	for _, customerVal := range data.Customer {
		arr = append(arr, customerVal)
	}
	for _, rfmVal := range data.Customer {
		arr = append(arr, rfmVal)
	}
	resp := marketing.UserGroupResponse{List: arr}

	ginplus.ResSuccess(c, resp)
}

//selectorText 选择项
func selectorText(c *gin.Context) {
	var request interface{}
	response, err := microscrm.GetSelectorText(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

//activities 活动列表
func activities(c *gin.Context) {
	var request marketing.ActivitiesRequest
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	tokenInfo := http.GetTokenInfo(c)
	request.EnterpriseID = tokenInfo.EnterpriseId
	response, err := microactivity.GetActivityLists(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	bytes, _ := json.Marshal(response)
	var data marketing.ActivitiesResponse
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, data)
}

//activityPages 活动页列表
func activityPages(c *gin.Context) {
	var request marketing.ActivityPagesRequest
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	tokenInfo := http.GetTokenInfo(c)
	req := microshop.ActivityPagesRequest{
		EnterpriseID: tokenInfo.EnterpriseId,
		Page:         request.Page,
		Size:         request.Size,
		TemplateType: request.TemplateType,
		Keywords:     request.Search,
	}
	if request.SearchType == 1 {
		req.KeyType = "id"
	} else {
		req.KeyType = "template_name"
	}
	response, err := microshop.GetActivitiesPages(ctx.Wrap(c), req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response.Data)
}

//lives 直播列表
func lives(c *gin.Context) {
	var request marketing.LivesRequest
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	tokenInfo := http.GetTokenInfo(c)
	req := microshop.LivesRequest{
		Size:         request.Size,
		Page:         request.Page,
		EnterpriseID: tokenInfo.EnterpriseId,
	}
	if request.SearchType == 1 {
		req.Search.ID = request.Search
	} else {
		req.Search.LiveTitle = request.Search
	}

	response, err := microshop.LiveList(ctx.Wrap(c), req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	bytes, _ := json.Marshal(response)
	var data marketing.LivesResponse
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, data)
}

//marketingStrategy 添加营销策略
func marketingStrategy(c *gin.Context) {
	var request marketing.MarketingStrategy
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	tokenInfo := http.GetTokenInfo(c)
	request.EnterpriseID = tokenInfo.EnterpriseId
	response, err := microscrm.CreateMarketing(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

//marketingStrategies 营销策略列表
func marketingStrategies(c *gin.Context) {
	var request marketing.ListRequest
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	tokenInfo := http.GetTokenInfo(c)
	request.EnterpriseID = tokenInfo.EnterpriseId
	response, err := microscrm.MarketingList(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

//marketingStrategyDetail 营销策略详情
func marketingStrategyDetail(c *gin.Context) {
	var request marketing.MarketingDetailRequest
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	tokenInfo := http.GetTokenInfo(c)
	request.EnterpriseID = tokenInfo.EnterpriseId
	response, err := microscrm.MarketingDetail(ctx.Wrap(c), request)
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
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	tokenInfo := http.GetTokenInfo(c)
	request.EnterpriseID = tokenInfo.EnterpriseId
	response, err := microscrm.UpdateMarketing(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

//marketingStrategyDelete 营销策略删除
func marketingStrategyDelete(c *gin.Context) {
	var request marketing.DeleteRequest
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	tokenInfo := http.GetTokenInfo(c)
	request.EnterpriseID = tokenInfo.EnterpriseId
	response, err := microscrm.DeleteMarketing(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

//策略数据导出
func marketingStrategyDetailExport(c *gin.Context) {
	var request marketing.MarketingDataDetailRequest
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	tokenInfo := http.GetTokenInfo(c)
	request.EnterpriseID = tokenInfo.EnterpriseId
	log.Info(ctx.New()).Msgf("策略数据详情请求参数%+v", request)
	response, err := microscrm.MarketingStrategyDataDetail(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	log.Info(ctx.New()).Msgf("策略数据详情导出数据%+v", response)
	var typeName string = " "
	excelFile := excelize.NewFile()
	sheetName := "数据详情跟踪"
	// 创建sheet
	index := excelFile.NewSheet(sheetName)
	// 删除默认的sheet Sheet1
	excelFile.DeleteSheet(ws_file.DefaultCreatSheetName)

	// 写入表头
	_ = excelFile.SetCellStr(sheetName, "A1", "渠道类型")
	_ = excelFile.SetCellStr(sheetName, "B1", "发送人数")
	_ = excelFile.SetCellStr(sheetName, "C1", "触达人数")
	rowStart := 2
	if response.Data.Data != nil {
		for k, detailStrategData := range response.Data.Data {
			// 起始行不增长了
			if k != 0 {
				rowStart = rowStart + 1
			}
			rowStartStr := strconv.Itoa(rowStart)
			switch detailStrategData.Types {
			case 1:
				typeName = "短信"
			case 2:
				typeName = "企业群发"
			case 3:
				typeName = "自动发券"
			case 4:
				typeName = "公众号"
			default:
				typeName = "短信"
			}
			excelFile.SetCellStr(sheetName, "A"+rowStartStr, typeName)
			excelFile.SetCellInt(sheetName, "B"+rowStartStr, detailStrategData.SendMemberNumber)
			excelFile.SetCellInt(sheetName, "C"+rowStartStr, detailStrategData.ReachMemberNumber)
		}
	}
	//最后一条总计
	if response.Data.Total != nil {
		excelFile.SetCellStr(sheetName, "A"+strconv.Itoa(rowStart+1), "总计")
		excelFile.SetCellInt(sheetName, "B"+strconv.Itoa(rowStart+1), response.Data.Total.SendMemberNumber)
		excelFile.SetCellInt(sheetName, "C"+strconv.Itoa(rowStart+1), response.Data.Total.ReachMemberNumber)
	}

	excelFile.SetActiveSheet(index)
	fileName := "策略数据" + time.Now().Format("200601021504") + ".xlsx"
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	des, _ := excelFile.WriteToBuffer()
	_, _ = c.Writer.Write(des.Bytes())
	return
}

//获取导出地址
/*
func getStrategyDataDetailUrl(c *gin.Context) {
	var request marketing.MarketingDataDetailRequest
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	requestTLS := c.Request.TLS
	scheme := "http://"
	if requestTLS != nil {
		scheme = "https://"
	}
	tokenInfo := http.GetTokenInfo(c)
	enterpriseId := tokenInfo.EnterpriseId //int64(133)
	strategyDataDownloadUrl := scheme + c.Request.Host +
		"/gw-scrm/marketing-strategies/detail/export?id=" + strconv.FormatInt(request.ID, 10) +
		"&day=" + strconv.Itoa(request.Day) +
		"&begin_time=" + strconv.FormatInt(request.BeginTime, 10) +
		"&end_time=" + strconv.FormatInt(request.EndTime, 10) +
		"&enterprise_id=" + strconv.FormatInt(enterpriseId, 10)
	ginplus.ResSuccess(c, map[string]string{
		"url": strategyDataDownloadUrl,
	})
}

*/
func marketingStrategyDataDetail(c *gin.Context) {
	var request marketing.MarketingDataDetailRequest
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	tokenInfo := http.GetTokenInfo(c)
	request.EnterpriseID = tokenInfo.EnterpriseId //133
	response, err := microscrm.MarketingStrategyDataDetail(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	resp := marketing.StrategyFinlaDataDetails{
		Data:  response.Data.Data,
		Total: response.Data.Total,
	}
	ginplus.ResSuccess(c, resp)
}

func updateStrategyDataTime(c *gin.Context) {
	var request marketing.UpdateSettingDayRequest
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	tokenInfo := http.GetTokenInfo(c)
	request.EnterpriseID = int(tokenInfo.EnterpriseId) //
	response, err := microscrm.UpdateStrategyDataTime(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

//marketingStrategyUpdateStatus 营销策略更新状态
func marketingStrategyUpdateStatus(c *gin.Context) {
	var request marketing.UpdateStatusRequest
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	tokenInfo := http.GetTokenInfo(c)
	request.EnterpriseID = tokenInfo.EnterpriseId
	response, err := microscrm.MarketingStrategyUpdateStatus(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// 营销策略触发
func MarketingStrategyTrigger(c *gin.Context) {

	var request marketing.MarketTrigRequest
	ok := validate_params.ValidateParams(c, &request)
	if !ok {
		return
	}

	tokenInfo := http.GetTokenInfo(c)
	request.EnterpriseID = tokenInfo.EnterpriseId
	response, err := microscrm.MarketingStrategyTrigger(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// paramsList 短信参数列表
func UserGroupList(c *gin.Context) {
	var request marketing.UserGroupList

	tokenInfo := http.GetTokenInfo(c)
	request.EnterpriseID = tokenInfo.EnterpriseId

	response, err := microscrm.MarketingUserGroupList(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

// TagGroupUserCount 获取指定标签中用户总数
func TagGroupUserCount(c *gin.Context) {
	var request marketing.TagGroupUserCountRequest
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}

	tokenInfo := http.GetTokenInfo(c)
	request.EnterpriseID = tokenInfo.EnterpriseId

	response, err := microscrm.TagGroupUserCount(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

func GetShortUrl(c *gin.Context) {
	request := marketing.ShortUrlRequest{}
	resp := marketing.ShortUrlRes{}

	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	log.Info(ctx.New()).Msgf("GetShortUrl获取短链req is %v\n", request)
	response, err := shorturl.GetShortUrlQuery(ctx.Wrap(c), request)
	if err != nil {
		log.Info(ctx.New()).Msgf("GetShortUrl获取短链失败 err is %v\n", err)
		ginplus.ResError(c, err)
		return
	}
	json.Unmarshal([]byte(response.(string)), &resp)
	ginplus.ResSuccess(c, resp)

}

func ShortUrlDel(c *gin.Context) {
	request := marketing.ShortUrlDelReq{}
	bindErr := c.BindJSON(&request)

	if bindErr != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	response, err := microscrm.ShortUrlDel(ctx.Wrap(c), request)
	if err != nil {
		log.Info(ctx.New()).Msgf("ShortUrlDel删除短链失败 err is %v\n", err)
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

//获取短链列表
func ShortUrlList(c *gin.Context) {
	request := marketing.ShortUrlListReq{}
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	log.Info(ctx.New()).Msgf("ShortUrlList获取短链列表req is %v\n", request)

	response, err := microscrm.GetShortUrlList(ctx.Wrap(c), request)
	if err != nil {
		log.Info(ctx.New()).Msgf("ShortUrlList获取短链列表失败 err is %v\n", err)
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

//获取短链详情
func ShortUrlData(c *gin.Context) {
	request := marketing.ShortUrlDataReq{}
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	response, err := microscrm.GetShortUrlData(ctx.Wrap(c), request)
	if err != nil {
		log.Info(ctx.New()).Msgf("ShortUrlData接口获取短链详情失败 err is %v\n", err)
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

//添加短链以及配置
func ShortUrlAdd(c *gin.Context) {
	request := marketing.ShortUrlAddReq{}
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	log.Info(ctx.New()).Msgf("ShortUrlAdd添加短链req is %v\n", request)

	tokenInfo := http.GetTokenInfo(c)
	w := md5.New()
	io.WriteString(w, "a3f24f4a11fd0135627ddd8ab9f40cbe")

	params := make(map[string]interface{})
	params["enterprise_id"] = tokenInfo.EnterpriseId
	params["Auth-Sign"] = fmt.Sprintf("%x", w.Sum(nil))
	//获取AccessToken
	response, err := gwshop.GetAccessTokenQuery(ctx.Wrap(c), params)
	if err != nil {
		log.Info(ctx.New()).Msgf("ShortUrlAdd接口获取AccessToken失败 err is %v\n", err)
		ginplus.ResError(c, err)
		return
	}
	accessTo := response.(map[string]interface{})
	pageUrl := strings.Split(request.PageUrl, "?")
	var wxUrl string
	var code int64
	code, wxUrl, _ = wxapi.GetWxUrlQuery(ctx.Wrap(c), pageUrl, accessTo["access_token"])
	if code == 40001 {
		i := 1
		for i <= 10 {
			response, err := gwshop.GetAccessTokenQuery(ctx.Wrap(c), params)
			if err != nil {
				log.Info(ctx.New()).Msgf("ShortUrlAdd接口获取AccessToken失败 次数%v err is %v\n", i, err)
				ginplus.ResError(c, err)
				continue
			}
			accessTo := response.(map[string]interface{})
			code, wxUrl, _ = wxapi.GetWxUrlQuery(ctx.Wrap(c), pageUrl, accessTo["access_token"])
			if code == 40001 {
				log.Info(ctx.New()).Msgf("ShortUrlAdd接口获取wx微信openLink失败 次数%v err is %v\n", i, err)
				i = i + 1
				continue
			} else {
				break
			}
		}
	}

	request.WxUrl = wxUrl
	var wxActiveUrl string
	//获取活动渠道对应的微信openlink
	if request.ActivePageUrl != "" {
		pageUrl := strings.Split(request.ActivePageUrl, "?")
		code, wxActiveUrl, _ = wxapi.GetWxUrlQuery(ctx.Wrap(c), pageUrl, accessTo["access_token"])
		if code == 40001 {
			i := 1
			for i <= 10 {
				response, err := gwshop.GetAccessTokenQuery(ctx.Wrap(c), params)
				if err != nil {
					log.Info(ctx.New()).Msgf("ShortUrlAdd接口获取AccessToken失败 err is %v\n", err)
					ginplus.ResError(c, err)
					continue
				}
				accessTo := response.(map[string]interface{})
				code, wxActiveUrl, _ = wxapi.GetWxUrlQuery(ctx.Wrap(c), pageUrl, accessTo["access_token"])
				if code == 40001 {
					log.Info(ctx.New()).Msgf("ShortUrlAdd接口获取active微信openLink失败 次数 %v err is %v\n", i, err)
					i = i + 1
					continue
				} else {
					break
				}
			}
		}
		request.WxActivePageUrl = wxActiveUrl
	}
	if len(request.SpuData) > 0 {

		for _, v := range request.SpuData {
			var code int64
			var wxSpuUrl string
			pageUrl := strings.Split(v.SpuPageUrl, "?")
			code, wxSpuUrl, err = wxapi.GetWxUrlQuery(ctx.Wrap(c), pageUrl, accessTo["access_token"])
			if code == 40001 {
				i := 1
				for i <= 10 {
					response, err := gwshop.GetAccessTokenQuery(ctx.Wrap(c), params)
					if err != nil {
						log.Info(ctx.New()).Msgf("ShortUrlAdd接口获取AccessToken失败 err is %v\n", err)
						ginplus.ResError(c, err)
						continue
					}
					accessTo := response.(map[string]interface{})
					code, wxSpuUrl, _ = wxapi.GetWxUrlQuery(ctx.Wrap(c), pageUrl, accessTo["access_token"])
					if code == 40001 {
						log.Info(ctx.New()).Msgf("ShortUrlAdd接口获取微信openLink失败 次数%v err is %v\n", i, err)
						i = i + 1
						continue
					} else {
						break
					}
				}
			}
			v.WxSpuPageUrl = wxSpuUrl
		}
	}
	spuData, err := json.Marshal(request.SpuData)

	response, err = microscrm.ShortUrlAdd(ctx.Wrap(c), request)
	if err != nil {
		log.Info(ctx.New()).Msgf("ShortUrlAdd接口添加短链失败 err is %v\n", err)
		ginplus.ResError(c, err)
		return
	}
	ShortUrlSlice := strings.Split(request.ShortUrl, "/")
	shortJson := fmt.Sprintf(shorturl.ShortHtml, request.CardTitle, request.CardDesc, request.BrandImg, request.BrandImg, request.BackgroundImg, request.ButtonColor, request.ButtonName, request.ButtonFontColor, request.LocationDesc, request.LocationFontColor, request.IsHide, request.ActiveImg, request.ActivePageUrl, request.WxActivePageUrl, string(spuData), request.WxUrl)
	requestAuth := marketing.AuthDetailRequest{EnterpriseID: request.EnterpriseId, AuthType: 1}

	response, err = microscrm.GetOriginId(ctx.Wrap(c), requestAuth)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	fmt.Printf("getoriginid is %v err is %v\n", response, err)
	ShortAddReq := marketing.ShortAddReq{}
	ShortAddReq.ShortUrl = ShortUrlSlice[3]
	ShortAddReq.PageUrl = request.PageUrl

	ShortAddReq.OriginId = response.(string)
	ShortAddReq.LongUrl = "https://crs.vchangyi.com"
	ShortAddReq.MiddleHtml = shortJson
	ShortAddReq.ShortUrlType = 2
	fmt.Printf("MiddleData is %v\n", request.MiddleData)
	ShortAddReq.MiddleData = request.MiddleData
	ShortAddReq.EnterpriseId = request.EnterpriseId
	if request.Id != 0 {
		ShortAddReq.ReqType = 1
	}
	fmt.Printf("getoriginidis is %v err is %v\n", ShortAddReq.PageUrl, err)

	response, err = shorturl.AddShortUrlQuery(ctx.Wrap(c), ShortAddReq)
	if err != nil {
		log.Info(ctx.New()).Msgf("ShortUrlAdd接口调用shortUrl接口失败 err is %v\n", err)
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nil)
}

//获取商品详情
func GetSpuDetail(c *gin.Context) {
	request := marketing.SpuDetailReq{}

	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	response, err := microscrm.GetSpuDetail(ctx.Wrap(c), request)
	fmt.Sprintf("get-spu-detail err %v\n", err)
	if err != nil {
		log.Info(ctx.New()).Msgf("GetSpuDetail接口获取商品详情失败 err is %v\n", err)
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)

}

// GetChanList 获取渠道列表
func GetChanList(c *gin.Context) {
	var request marketing.ChanListRequest
	var resp marketing.ChanListCDPResp
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}

	tokenInfo := http.GetTokenInfo(c)
	request.EnterpriseID = tokenInfo.EnterpriseId

	cdpSign := marketing.ChanListRequestBody{}
	cdpSign.SignTime = time.Now().Unix()
	sign, err := cdp.GetSign(map[string]interface{}{
		"sign_time": cdpSign.SignTime,
	})

	cdpSign.Sign = sign
	response, err := cdp.ChanListQuery(ctx.Wrap(c), request, cdpSign)
	if response == nil {
		resp.List = make([]*marketing.ChanInfo, 0)
		ginplus.ResSuccess(c, resp)
		return
	}

	chanListRes := marketing.ChanListCDPRes{}
	var m map[string]interface{}
	_ = json.Unmarshal([]byte(response.(string)), &m)
	if ok := m["code"]; ok == nil {
		_ = json.Unmarshal([]byte(response.(string)), &chanListRes)
	}
	resp.List = chanListRes.Data
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	response, err = microscrm.GetTwoLevelIds(ctx.Wrap(c), request)
	if err != nil {
		log.Info(ctx.New()).Msgf("GetChanList接口获取二级渠道id失败 err is %v\n", err)
		ginplus.ResError(c, err)
		return
	}
	if response != nil {
		re := response.([]interface{})
		for _, j := range resp.List {
			for _, v := range re {
				mp := v.(map[string]interface{})
				if float64(j.ID) == mp["id"].(float64) {
					j.IsHide = 1
					break
				}
			}
		}
	}
	ginplus.ResSuccess(c, resp)
}

//查询券是否在使用中
func GouponIsUsed(c *gin.Context) {
	var request marketing.CouponIsUsedRequest
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}

	resp, err := microscrm.GouponIsUsed(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, resp)
}

//券匹配
func CouponMatch(c *gin.Context) {
	var request marketing.CouponMatchRequest
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}

	resp, err := microscrm.CouponMatch(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, resp)
}

//策略自动发券数据导出--新
func marketingStrategyDetailExportNew(c *gin.Context) {
	var request marketing.MarketingDataDetailNewRequest
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	//request.EnterpriseID = 133//调试完成删掉
	tokenInfo := http.GetTokenInfo(c)
	request.EnterpriseID = tokenInfo.EnterpriseId //int64(133)
	response, err := microscrm.MarketingDataCoupon(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	log.Info(ctx.New()).Msgf("策略数据详情导出数据%+v", response)
	excelFile := excelize.NewFile()
	sheetName := "数据详情跟踪"
	excelFile.SetSheetName("Sheet1", sheetName)
	rowStart := 1

	_ = excelFile.SetCellStr(sheetName, "A1", "渠道类型")
	_ = excelFile.SetCellStr(sheetName, "B1", "发送人数")
	_ = excelFile.SetCellStr(sheetName, "C1", "微信优惠券")
	_ = excelFile.SetCellStr(sheetName, "D1", "批次号")
	_ = excelFile.SetCellStr(sheetName, "E1", "触达人数")
	_ = excelFile.SetCellStr(sheetName, "F1", "发放张数")
	_ = excelFile.SetCellStr(sheetName, "G1", "核销张数")
	_ = excelFile.SetCellStr(sheetName, "H1", "核销率(%)")
	_ = excelFile.SetCellStr(sheetName, "I1", "支付买家数")
	_ = excelFile.SetCellStr(sheetName, "J1", "触达支付转化率(%)")
	_ = excelFile.SetCellStr(sheetName, "K1", "实付金额")
	_ = excelFile.SetCellStr(sheetName, "L1", "退款买家数")
	_ = excelFile.SetCellStr(sheetName, "M1", "退款金额")
	total := 0
	for _, v := range response.Data.Data.CouponList {
		total += 1
		rowStart = rowStart + 1
		rowStartStr := strconv.Itoa(rowStart)
		excelFile.SetCellStr(sheetName, "A"+rowStartStr, response.Data.Data.ChannelName)
		excelFile.SetCellInt(sheetName, "B"+rowStartStr, response.Data.Data.SendNum)
		excelFile.SetCellStr(sheetName, "C"+rowStartStr, v.CouponName)
		excelFile.SetCellStr(sheetName, "D"+rowStartStr, v.StockID)
		excelFile.SetCellInt(sheetName, "E"+rowStartStr, v.TouchNum)
		excelFile.SetCellInt(sheetName, "F"+rowStartStr, v.IsSendNum)
		excelFile.SetCellInt(sheetName, "G"+rowStartStr, v.UseNum)
		if v.IsSendNum != 0 {
			thanInt := (float64(v.UseNum) / float64(v.IsSendNum)) * 100
			useNumThan, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", thanInt), 64)
			excelFile.SetCellFloat(sheetName, "H"+rowStartStr, useNumThan, -1, 64)
		} else {
			excelFile.SetCellStr(sheetName, "H"+rowStartStr, "0")
		}

		excelFile.SetCellInt(sheetName, "I"+rowStartStr, v.PayNum)

		if v.TouchNum != 0 {
			payInt := (float64(v.PayNum) / float64(v.TouchNum)) * 100
			payNumThan, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", payInt), 64)
			excelFile.SetCellFloat(sheetName, "J"+rowStartStr, payNumThan, -1, 64)
		} else {
			excelFile.SetCellStr(sheetName, "J"+rowStartStr, "0")
		}

		excelFile.SetCellValue(sheetName, "K"+rowStartStr, v.CashFee)
		excelFile.SetCellInt(sheetName, "L"+rowStartStr, v.RefundCustomerNum)
		excelFile.SetCellValue(sheetName, "M"+rowStartStr, v.RefundFee)
	}

	//最后一条总计
	if response.Data.Total != nil {
		excelFile.SetCellStr(sheetName, "A"+strconv.Itoa(total+2), response.Data.Total.ChannelName)
		excelFile.SetCellInt(sheetName, "B"+strconv.Itoa(total+2), response.Data.Total.SendNum)
		excelFile.SetCellStr(sheetName, "C"+strconv.Itoa(total+2), response.Data.Total.CouponList.CouponName)
		excelFile.SetCellStr(sheetName, "D"+strconv.Itoa(total+2), response.Data.Total.CouponList.StockID)
		excelFile.SetCellInt(sheetName, "E"+strconv.Itoa(total+2), response.Data.Total.CouponList.TouchNum)
		excelFile.SetCellInt(sheetName, "F"+strconv.Itoa(total+2), response.Data.Total.CouponList.IsSendNum)
		excelFile.SetCellInt(sheetName, "G"+strconv.Itoa(total+2), response.Data.Total.CouponList.UseNum)
		if response.Data.Total.CouponList.IsSendNum != 0 {
			thanIntH := (float64(response.Data.Total.CouponList.UseNum) / float64(response.Data.Total.CouponList.IsSendNum)) * 100
			useNumThanH, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", thanIntH), 64)
			excelFile.SetCellFloat(sheetName, "H"+strconv.Itoa(total+2), useNumThanH, -1, 64)
		} else {
			excelFile.SetCellStr(sheetName, "H"+strconv.Itoa(total+2), "0")
		}
		excelFile.SetCellInt(sheetName, "I"+strconv.Itoa(total+2), response.Data.Total.CouponList.PayNum)
		if response.Data.Total.CouponList.TouchNum != 0 {
			payIntH := (float64(response.Data.Total.CouponList.PayNum) / float64(response.Data.Total.CouponList.TouchNum)) * 100
			payNumThanH, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", payIntH), 64)
			excelFile.SetCellFloat(sheetName, "J"+strconv.Itoa(total+2), payNumThanH, -1, 64)
		} else {
			excelFile.SetCellStr(sheetName, "J"+strconv.Itoa(total+2), "0")
		}
		excelFile.SetCellValue(sheetName, "K"+strconv.Itoa(total+2), response.Data.Total.CouponList.CashFee)
		excelFile.SetCellInt(sheetName, "L"+strconv.Itoa(total+2), response.Data.Total.CouponList.RefundCustomerNum)
		excelFile.SetCellValue(sheetName, "M"+strconv.Itoa(total+2), response.Data.Total.CouponList.RefundFee)
	}

	excelFile.MergeCell(sheetName, "A2", "A"+strconv.Itoa(total+1))
	excelFile.MergeCell(sheetName, "B2", "B"+strconv.Itoa(total+1))

	style, _ := excelFile.NewStyle(`{"alignment":{"horizontal":"center","vertical":"center"}}`)
	_ = excelFile.SetCellStyle(sheetName, "A2", "A"+strconv.Itoa(total+1), style)

	fileName := "策略数据" + time.Now().Format("200601021504") + ".xlsx"
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	des, _ := excelFile.WriteToBuffer()
	_, _ = c.Writer.Write(des.Bytes())
	return
}

//marketingStrategyDetailCoupon 营销策略自动发券详情
func marketingStrategyDetailCoupon(c *gin.Context) {
	var request marketing.MarketingDetailRequest
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	tokenInfo := http.GetTokenInfo(c)
	request.EnterpriseID = tokenInfo.EnterpriseId
	response, err := microscrm.MarketingDataCoupon(ctx.Wrap(c), request)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

//获取短链列表
func GetAllShortUrlList(c *gin.Context) {
	request := marketing.ShortUrlListReq{}
	bindErr := c.BindJSON(&request)
	if bindErr != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	response, err := microscrm.GetAllShortUrlList(ctx.Wrap(c), request)
	if err != nil {
		log.Info(ctx.New()).Msgf("ShortUrlList获取短链列表失败 err is %v\n", err)
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

func GetTemplateList(c *gin.Context) {
	var request marketing.TemplateListReq
	bindErr := c.ShouldBind(&request)
	if bindErr != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	tokenInfo := http.GetTokenInfo(c)
	request.EnterpriseId = tokenInfo.EnterpriseId
	response, err := microscrm.GetTemplateList(ctx.Wrap(c), request)
	if err != nil {
		log.Info(ctx.New()).Msgf("获取贺卡模板列表失败 err is %v\n", err)
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

func CreateTemplateInfo(c *gin.Context) {
	var request marketing.TemplateCreate
	bindErr := c.ShouldBind(&request)
	if bindErr != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	tokenInfo := http.GetTokenInfo(c)
	request.EnterpriseId = tokenInfo.EnterpriseId
	response, err := microscrm.CreateTemplateInfo(ctx.Wrap(c), request)
	if err != nil {
		log.Info(ctx.New()).Msgf("创建贺卡模板列表失败 err is %v\n", err)
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

func DelTemplateInfo(c *gin.Context) {
	var request marketing.TemplateDel
	bindErr := c.ShouldBind(&request)
	if bindErr != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	tokenInfo := http.GetTokenInfo(c)
	request.EnterpriseId = tokenInfo.EnterpriseId
	response, err := microscrm.DelTemplateInfo(ctx.Wrap(c), request)
	if err != nil {
		log.Info(ctx.New()).Msgf("删除贺卡模板列表失败 err is %v\n", err)
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}

func CheckTemplateInfo(c *gin.Context) {
	var request marketing.TemplateDel
	bindErr := c.ShouldBind(&request)
	if bindErr != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	tokenInfo := http.GetTokenInfo(c)
	request.EnterpriseId = tokenInfo.EnterpriseId
	response, err := microscrm.CheckTemplateInfo(ctx.Wrap(c), request)
	if err != nil {
		log.Info(ctx.New()).Msgf("检验贺卡模板列表失败 err is %v\n", err)
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}


func GetTemplateInfo(c *gin.Context){
	var request marketing.TemplateDel
	bindErr := c.ShouldBind(&request)
	if bindErr != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	tokenInfo := http.GetTokenInfo(c)
	request.EnterpriseId = tokenInfo.EnterpriseId
	response, err := microscrm.GetTemplateInfo(ctx.Wrap(c), request)
	if err != nil {
		log.Info(ctx.New()).Msgf("删除贺卡模板列表失败 err is %v\n", err)
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, response)
}
