package microScrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
)

// TagApi 员工 SDK
type TagApi struct {
	BaseSDK
}

// NewTag 返回API实例
func NewTag(host string) *TagApi {
	urlHost = host
	return &TagApi{BaseSDK{c: &clients.HTTPClient{}}}
}


// TagList 足迹列表
func (c TagApi)TagList(ctx ctx.Context, req interface{}) (interface{}, error) {
	path := "/micro-scrm/v1/tag/list"
	url := urlHost + path
	res, err := c.Post(ctx, url, req)
	if err != nil {
		return nil, err
	}

	result := &Response{
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

// CustomerManualList 足迹列表
func (c TagApi)CustomerManualList(ctx ctx.Context, req interface{}) (interface{}, error) {
	path := "/micro-scrm/v1/tag/customer-manual-list"
	url := urlHost + path
	res, err := c.Post(ctx, url, req)
	if err != nil {
		return nil, err
	}

	result := &Response{
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}