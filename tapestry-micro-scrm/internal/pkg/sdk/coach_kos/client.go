package coachkos

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

const (
	listPath        = "/api/biz/coach_mark_post_results"
	defaultPageSize = 1000
	maxPageSize     = 1000
)

// Client 是 COACH KOS HTTP 客户端。
type Client interface {
	List(ctx.Context, ListRequest) (PageData, error)
}

type client struct {
	baseURL  string
	token    string
	pageSize int
	http     *http.Client
}

// NewClient 创建 COACH KOS HTTP 客户端。
func NewClient(baseURL, token string, pageSize int) Client {
	return &client{
		baseURL:  strings.TrimRight(baseURL, "/"),
		token:    token,
		pageSize: pageSize,
		http:     &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *client) List(context ctx.Context, request ListRequest) (PageData, error) {
	if c.baseURL == "" {
		return PageData{}, fmt.Errorf("COACH KOS配置错误：baseURL不能为空")
	}
	if c.token == "" {
		return PageData{}, fmt.Errorf("COACH KOS配置错误：token不能为空")
	}

	query := url.Values{}
	query.Set("updated_at", request.UpdatedAt)
	query.Set("page", strconv.Itoa(request.Page))
	query.Set("size", strconv.Itoa(c.requestPageSize(request.Size)))
	httpRequest, err := http.NewRequestWithContext(context, http.MethodGet, c.baseURL+listPath+"?"+query.Encode(), nil)
	if err != nil {
		return PageData{}, fmt.Errorf("创建COACH KOS请求失败：%w", err)
	}
	httpRequest.Header.Set("X-Coach-Mark-Token", c.token)

	httpResponse, err := c.http.Do(httpRequest)
	if err != nil {
		return PageData{}, fmt.Errorf("请求COACH KOS接口失败：%w", err)
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode != http.StatusOK {
		return PageData{}, fmt.Errorf("COACH KOS接口HTTP状态异常，status=%d", httpResponse.StatusCode)
	}

	var response listResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return PageData{}, fmt.Errorf("解析COACH KOS接口JSON失败：%w", err)
	}
	if response.Code != 0 {
		return PageData{}, fmt.Errorf("COACH KOS接口返回失败，code=%d，message=%s", response.Code, response.Message)
	}
	return response.Data, nil
}

func (c *client) requestPageSize(requested int) int {
	size := requested
	if size <= 0 {
		size = c.pageSize
	}
	if size <= 0 {
		size = defaultPageSize
	}
	if size > maxPageSize {
		return maxPageSize
	}
	return size
}
