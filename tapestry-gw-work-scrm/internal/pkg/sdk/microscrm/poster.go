package microScrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
)

// PosterApi 海报 SDK
type PosterApi struct {
	BaseSDK
}

// NewPoster 返回API实例
func NewPoster(host string) *PosterApi {
	urlHost = host
	return &PosterApi{BaseSDK{c: &clients.HTTPClient{}}}
}

// Poster 下载海报
func (s PosterApi) Poster(c ctx.Context, request PosterRequest) (interface{}, error) {
	path := "/micro-scrm/v1/marketing-promotion/poster"
	url := urlHost + path
	res, err := s.Post(c, url, request)
	if err != nil {
		return nil, err
	}

	result := &Response{}

	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}
