package baozun

import (
	"encoding/json"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/config"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/redis"
	"micro-scrm/internal/pkg/sdk"
	"strings"
	"time"
)

const (
	TokenRedisKey       = "baozun:resources-token:%d"
	ProductInfoRedisKey = "baozun:product-info:%d:%s" // enterprise_id:product_no
)

func getToken(ctx ctx.Context, enterpriseId int64) (string, error) {
	var resp TokenResponse
	var opts []clients.Option
	// 获取url
	conf := config.CfgData.Baozun
	url := conf.Domain + "/auth/realms/open-gateway/protocol/openid-connect/token"
	// 请求参数
	var reqMap = map[string]string{
		"client_id":     conf.ClientId,
		"client_secret": conf.ClientSecret,
		"grant_type":    "client_credentials",
	}
	body := clients.Form(reqMap)
	opts = append(opts, body)
	// 返回值
	var respStr string
	reqStr := helper.MarshalToString(reqMap)
	err := helper.WithRetry(1, time.Second, func() error {
		opts = append(opts, clients.Header("Content-Type", "application/x-www-form-urlencoded"))
		client := clients.NewHTTPClient(clients.Timeout(120 * time.Second))
		res, err := client.Post(ctx, url, opts...)
		if err != nil {
			return errorx.NewErrorMsg(500, fmt.Sprintf("三方商品图片接口-token获取失败：%+v，res:%v", err, res))
		}

		if res.Response().StatusCode != 200 {
			return errorx.NewErrorMsg(res.Response().StatusCode, fmt.Sprintf("三方商品图片接口-获取token状态码异常：%s", res.Response().Status))
		}

		rt, rErr := res.Read(clients.Text)
		if rErr != nil {
			return errorx.NewErrorMsg(res.Response().StatusCode, fmt.Sprintf("三方商品图片接口-token数据获取失败：%+v", rErr))
		}

		log.Info(ctx).Msgf("三方商品图片接口获取token, 地址：%s, 参数：%+v, 返回结果：%s", url, reqStr, rt.(string))

		respStr = rt.(string)
		uIntErr := json.Unmarshal([]byte(respStr), &resp)
		if uIntErr != nil || resp.Error != "" {
			return errorx.NewErrorMsg(res.Response().StatusCode, fmt.Sprintf("三方商品图片接口-token返回数据解析失败：%+v", respStr))
		}
		return nil
	})

	if err != nil || resp.AccessToken == "" {
		resErr := err.(*errorx.Error)
		log.Error(ctx).Msgf("%s：%+v", resErr.Error(), resErr)
		go sdk.SendMessage(ctx, "三方商品图片接口-获取token异常", url, reqStr, respStr, resErr)
		return "", err
	}
	redis.Client().Set(ctx, fmt.Sprintf(TokenRedisKey, enterpriseId), resp.AccessToken, 1500*time.Second)

	return resp.AccessToken, nil
}

func SearchResources(ctx ctx.Context, searchParams SearchParams, resp *DataResponse) (err error) {
	accessToken := redis.Client().Get(ctx, fmt.Sprintf(TokenRedisKey, searchParams.EnterpriseId)).Val()
	if accessToken == "" {
		accessToken, err = getToken(ctx, searchParams.EnterpriseId)
		if accessToken == "" {
			return err
		}
	}
	var opts []clients.Option
	var reqStr string
	//处理like
	likeStr := helper.RuleThree(searchParams.IsLike, "like:", "")
	if searchParams.Upc != "" || searchParams.ProductNo != "" {
		//upcCode查询
		var tmpRequest ParamsUpcReq
		tmpRequest.Sort.SortType = "desc"
		tmpRequest.Sort.Column = "id"
		if searchParams.Upc != "" {
			tmpRequest.Filter.SkfParam.Properties.Upc = searchParams.Upc
		} else if searchParams.ProductNo != "" {
			tmpRequest.Filter.SkfParam.VariantCode = likeStr + searchParams.ProductNo
		}
		tmpRequest.Filter.IncludeNoneImageProduct = searchParams.IncludeNoneImageProduct
		if !searchParams.IncludeNoneImageProduct {
			tmpRequest.Filter.ResourceParam.Name = "!:suffix:swatch.jpg"
		}
		if searchParams.StoreCode != "" {
			tmpRequest.ShopCode = searchParams.StoreCode
		}
		tmpRequest.Pages.PageSize = searchParams.Size
		body := clients.JSONBody(tmpRequest)
		opts = append(opts, body)
		reqStr = helper.MarshalToString(tmpRequest)
		tmpRequest.Pages.PageSize = searchParams.Size
	} else if len(searchParams.UpcSlice) > 0 {
		//upcCode查询
		var tmpRequest ParamsUpcReq
		tmpRequest.Sort.SortType = "desc"
		tmpRequest.Sort.Column = "id"
		tmpRequest.Filter.SkfParam.Properties.Upc = "in:" + strings.Join(searchParams.UpcSlice, ",")
		tmpRequest.Filter.IncludeNoneImageProduct = searchParams.IncludeNoneImageProduct
		if !searchParams.IncludeNoneImageProduct {
			tmpRequest.Filter.ResourceParam.Name = "!:suffix:swatch.jpg"
		}
		if searchParams.StoreCode != "" {
			tmpRequest.ShopCode = searchParams.StoreCode
		}
		tmpRequest.Pages.PageSize = searchParams.Size
		body := clients.JSONBody(tmpRequest)
		opts = append(opts, body)
		reqStr = helper.MarshalToString(tmpRequest)
	} else if searchParams.ProductName != "" || searchParams.SkuColor != "" || searchParams.SkuItem != "" {
		//商品名称检索
		var tmpRequest ParamsSpuReq
		tmpRequest.Sort.SortType = "desc"
		tmpRequest.Sort.Column = "id"
		tmpRequest.Filter.IncludeNoneImageProduct = searchParams.IncludeNoneImageProduct
		if searchParams.ProductName != "" {
			tmpRequest.Filter.SpuParam.Properties.ProductName = likeStr + searchParams.ProductName
		} else if searchParams.SkuItem != "" {
			tmpRequest.Filter.SpuParam.Properties.ProductNo = searchParams.SkuItem
		} else {
			tmpRequest.Filter.SpuParam.ProductCode = likeStr + searchParams.SkuColor
		}
		if !searchParams.IncludeNoneImageProduct {
			tmpRequest.Filter.ResourceParam.Name = "!:suffix:swatch.jpg"
		}
		if searchParams.StoreCode != "" {
			tmpRequest.ShopCode = searchParams.StoreCode
		}
		tmpRequest.Pages.PageSize = searchParams.Size
		tmpRequest.Pages.LastCode = searchParams.LastCode
		body := clients.JSONBody(tmpRequest)
		opts = append(opts, body)
		reqStr = helper.MarshalToString(tmpRequest)
	} else if len(searchParams.ProductCode) > 0 {
		//商品名称检索
		var tmpRequest ParamsSpuReq
		tmpRequest.Sort.SortType = "desc"
		tmpRequest.Sort.Column = "id"
		tmpRequest.Filter.IncludeNoneImageProduct = searchParams.IncludeNoneImageProduct
		if !searchParams.IncludeNoneImageProduct {
			tmpRequest.Filter.ResourceParam.Name = "!:suffix:swatch.jpg"
		}
		tmpRequest.Filter.SpuParam.ProductCode = "in:" + strings.Join(searchParams.ProductCode, ",")
		if searchParams.StoreCode != "" {
			tmpRequest.ShopCode = searchParams.StoreCode
		}
		tmpRequest.Pages.PageSize = searchParams.Size
		tmpRequest.Pages.LastCode = searchParams.LastCode
		body := clients.JSONBody(tmpRequest)
		opts = append(opts, body)
		reqStr = helper.MarshalToString(tmpRequest)
	} else {
		return errorx.NewErrorMsg(500, "参数错误")
	}

	// 获取url
	conf := config.CfgData.Baozun
	url := conf.Domain + "/rest/product/product/queryDetailRoll"
	// 请求参数

	// 返回值
	var respStr string
	err = helper.WithRetry(1, time.Second, func() error {
		opts = append(opts, clients.Header("Auth_Authorization", accessToken))
		opts = append(opts, clients.Header("Auth_ClientId", "TapestryCoach"))
		opts = append(opts, clients.Header("catalog", "COACH"))
		opts = append(opts, clients.Header("version", "1.0"))

		client := clients.NewHTTPClient(clients.Timeout(120 * time.Second))
		res, err := client.Post(ctx, url, opts...)
		if err != nil {
			return errorx.NewErrorMsg(500, fmt.Sprintf("三方商品图片接口-失败：%+v，res:%v", err, res))
		}

		if res.Response().StatusCode != 200 {
			return errorx.NewErrorMsg(res.Response().StatusCode, fmt.Sprintf("三方商品图片接口-状态码异常：%s", res.Response().Status))
		}

		rt, rErr := res.Read(clients.Text)
		if rErr != nil {
			return errorx.NewErrorMsg(res.Response().StatusCode, fmt.Sprintf("三方商品图片接口-数据获取失败：%+v", rErr))
		}

		log.Info(ctx).Msgf("三方商品图片接口, 地址：%s, 参数：%+v, 返回结果：%s", url, reqStr, rt.(string))

		respStr = rt.(string)
		uIntErr := json.Unmarshal([]byte(respStr), resp)
		if uIntErr != nil || resp.Code != "0" {
			return errorx.NewErrorMsg(res.Response().StatusCode, fmt.Sprintf("三方商品图片接口-返回数据解析失败：%+v", uIntErr))
		}
		return nil
	})

	if err != nil {
		resErr := err.(*errorx.Error)
		log.Error(ctx).Msgf("%s：%+v", resErr.Error(), resErr)
		go sdk.SendMessage(ctx, "三方商品图片接口-异常", url, reqStr, respStr, resErr)
		return err
	}
	return nil
}
