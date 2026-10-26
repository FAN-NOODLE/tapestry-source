package baiqiuSdk

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/config"
	"micro-scrm/internal/pkg/sdk"
	"sort"
	"strings"
	"time"
)

// BatchProductClient 批量产品查询客户端
type BatchProductClient struct {
	baseURL    string
	customerId string
	secret     string
	timeout    time.Duration
}

// NewBatchProductClient 创建批量产品查询客户端
func NewBatchProductClient() *BatchProductClient {
	conf := config.CfgData.BaiQiu
	return &BatchProductClient{
		baseURL:    conf.BaseURL,
		customerId: conf.CustomerId,
		secret:     conf.Secret,
		timeout:    time.Duration(60) * time.Second,
	}
}

// BatchSearch 批量产品搜索
func (c *BatchProductClient) BatchSearch(ctx ctx.Context, req BatchSearchRequest) (BatchSearchResponse, error) {
	var resp BatchSearchResponse

	// 验证请求参数
	if err := c.validateRequest(req); err != nil {
		return resp, err
	}

	// 构建查询字符串
	queryStr := c.buildQueryString(req.Products)

	// 构建请求参数
	params := map[string]string{
		"query":       queryStr,
		"customer_id": c.customerId,
		"timestamp":   fmt.Sprintf("%d", time.Now().Unix()),
	}
	params["signature"] = c.generateSignature(params)

	// 构建请求体
	requestBody := params

	// 发送请求
	url := c.baseURL + "/api/product/batch_search"
	err := c.doPostRequest(ctx, url, requestBody, &resp)
	if err != nil || resp.Code != 200 {
		go sdk.SendMessage(ctx, "Sw三方商品图片接口请求失败", url, requestBody, resp, err)
	}

	return resp, nil
}

// validateRequest 验证请求参数
func (c *BatchProductClient) validateRequest(req BatchSearchRequest) error {
	if len(req.Products) == 0 {
		return errorx.NewErrorMsg(400, "产品列表不能为空")
	}

	if len(req.Products) > 20 {
		return errorx.NewErrorMsg(400, fmt.Sprintf("批量查询最多支持20个产品，当前数量：%d", len(req.Products)))
	}

	for i, product := range req.Products {
		if strings.TrimSpace(product.ProductCode) == "" {
			return errorx.NewErrorMsg(400, fmt.Sprintf("第%d个产品的货号不能为空", i+1))
		}
	}

	return nil
}

// buildQueryString 构建查询字符串
func (c *BatchProductClient) buildQueryString(products []ProductQuery) string {
	var queries []string
	for _, product := range products {
		query := strings.ToUpper(strings.TrimSpace(product.ProductCode))
		if strings.TrimSpace(product.Color) != "" {
			query += "*" + strings.ToUpper(strings.TrimSpace(product.Color))
		}
		queries = append(queries, query)
	}
	return strings.Join(queries, ",")
}

// generateSignature 生成签名
func (c *BatchProductClient) generateSignature(params map[string]string) string {
	// 创建参数副本，避免修改原参数
	signParams := make(map[string]string)
	for k, v := range params {
		signParams[k] = v
	}

	// 按键名排序
	keys := make([]string, 0, len(signParams))
	for k := range signParams {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 构建参数字符串
	var paramParts []string
	for _, k := range keys {
		paramParts = append(paramParts, fmt.Sprintf("%s=%s", k, signParams[k]))
	}
	paramStr := strings.Join(paramParts, "&")

	// 追加密钥
	signStr := paramStr + "&key=" + c.secret

	// 计算MD5
	hash := md5.Sum([]byte(signStr))
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}

// doPostRequest 执行POST请求
func (c *BatchProductClient) doPostRequest(ctx ctx.Context, url string, body interface{}, resp interface{}) error {
	var opts []clients.Option

	// 设置请求头
	opts = append(opts, clients.Header("Content-Type", "application/json"))

	// 设置请求体
	if body != nil {
		opts = append(opts, clients.JSONBody(body))
	}

	// 直接执行请求，不移除重试
	client := clients.NewHTTPClient()
	if c.timeout > 0 {
		client = clients.NewHTTPClient(clients.Timeout(120 * time.Second))
	}

	res, err := client.Post(ctx, url, opts...)
	if err != nil {
		log.Error(ctx).Msgf("HTTP请求失败: url=%s, error=%v", url, err)
		return errors.New("HTTP请求失败" + err.Error())
	}

	// 检查HTTP状态码
	if res.Response().StatusCode >= 500 {
		log.Error(ctx).Msgf("HTTP状态码异常: url=%s, status=%d", url, res.Response().StatusCode)
		return errors.New(fmt.Sprintf("HTTP状态码异常: url=%s, status=%d", url, res.Response().StatusCode))
	}

	// 读取响应
	rt, rErr := res.Read(clients.Text)
	if rErr != nil {
		log.Error(ctx).Msgf("读取响应失败: url=%s, error=%v", url, rErr)
		return errors.New(fmt.Sprintf("读取响应失败: url=%s, error=%v", url, rErr))
	}

	respStr := rt.(string)

	// 解析JSON响应
	if err := json.Unmarshal([]byte(respStr), resp); err != nil {
		log.Error(ctx).Msgf("JSON解析失败: url=%s, response=%s, error=%v", url, respStr, err)
		return errors.New(fmt.Sprintf("JSON解析失败: url=%s, response=%s, error=%v", url, respStr, err))
	}

	log.Info(ctx).Msgf("批量产品查询成功: url=%s", url)
	return nil
}

// GetFoundProducts 获取查询成功的产品列表
func (r *BatchSearchResponse) GetFoundProducts() []Product {
	if r.Data == nil {
		return nil
	}
	return r.Data.Products
}

// GetNotFoundProducts 获取未找到的产品列表
func (r *BatchSearchResponse) GetNotFoundProducts() []NotFound {
	if r.Data == nil {
		return nil
	}
	return r.Data.NotFound
}

// IsSuccess 判断请求是否成功
func (r *BatchSearchResponse) IsSuccess() bool {
	return r.Code == 200
}

//
//// SimpleBatchSearch 简化版批量查询（直接传入查询字符串）
//func (c *BatchProductClient) SimpleBatchSearch(ctx ctx.Context, queryStr string) (BatchSearchResponse, error) {
//	var resp BatchSearchResponse
//
//	// 验证查询字符串
//	if strings.TrimSpace(queryStr) == "" {
//		return resp, errorx.NewErrorMsg(400, "查询字符串不能为空")
//	}
//
//	// 检查查询数量
//	queries := strings.Split(queryStr, ",")
//	if len(queries) > 20 {
//		return resp, errorx.NewErrorMsg(400, fmt.Sprintf("批量查询最多支持20个产品，当前数量：%d", len(queries)))
//	}
//
//	// 构建请求参数
//	params := map[string]string{
//		"query": queryStr,
//	}
//
//	// 添加时间戳和签名
//	timestamp := fmt.Sprintf("%d", time.Now().Unix())
//	params["timestamp"] = timestamp
//	params["signature"] = c.generateSignature(params)
//
//	// 构建请求体
//	requestBody := map[string]string{
//		"query":     queryStr,
//		"timestamp": timestamp,
//		"signature": params["signature"],
//	}
//
//	// 发送请求
//	url := c.baseURL + "/api/product/batch_search"
//	err := c.doPostRequest(ctx, url, requestBody, &resp)
//	if err != nil {
//		return resp, err
//	}
//
//	return resp, nil
//}
