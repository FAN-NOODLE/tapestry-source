package http

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/common"
	"strconv"
)

var ClientObject = clients.NewHTTPClient()

// 调用内部服务接口Client
type SdkClient struct {
	Client *clients.HTTPClient
}

// 返回参数结构
type GwResponse struct {
	Code      int64       `json:"code"`
	Message   string      `json:"message"`
	RequestId string      `json:"request_id"`
	Data      interface{} `json:"data"`
}

func SdkErr(c *gin.Context, res *GwResponse) {
	strInt64 := strconv.FormatInt(res.Code, 10)
	codeInt, _ := strconv.Atoi(strInt64)
	c.JSON(codeInt, res)
}

// Post 发送 Post 请求
func (s *SdkClient) Post(c ctx.Context, url string, options ...clients.Option) (*clients.ResponseProxy, error) {
	return s.Client.Post(c, url, options...)
}

// Get 发送 Get 请求
func (s *SdkClient) Get(c ctx.Context, url string, options ...clients.Option) (*clients.ResponseProxy, error) {
	return s.Client.Get(c, url, options...)
}

func GetTokenInfo(c *gin.Context) *common.TokenInfo {
	var tokenInfo *common.TokenInfo
	if c.Value(common.ContextKeyEnum.TokenInfo) != nil {
		tokenInfo = c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	}
	return tokenInfo
}
