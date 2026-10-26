package wechat

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"io"
	netHttp "net/http"
)

// HTTPGet get 请求
func HTTPGet(ctx ctx.Context, uri string) ([]byte, error) {
	response, err := clients.NewHTTPClient().Get(ctx, uri)
	if err != nil {
		return nil, err
	}
	if response.Response().StatusCode != netHttp.StatusOK {
		return nil, fmt.Errorf("http get error : uri=%v , statusCode=%v", uri, response.Response().StatusCode)
	}

	return io.ReadAll(response.Response().Body)
}
