package microScrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

// OvertimeApi 超时记录 SDK
type OvertimeApi struct {
	BaseSDK
}

func NewOvertime(host string) *OvertimeApi {
	urlHost = host
	return &OvertimeApi{BaseSDK{c: &clients.HTTPClient{}}}
}

// list 获取超时记录
func (api *OvertimeApi) ConfigInfo(c ctx.Context, overtime *ListRequest) (*PageListResponse, error) {
	path := "/micro-scrm/v1/overtime/list"
	url := urlHost + path
	res, err := api.Post(c, url, overtime)
	if err != nil {
		return nil, err
	}
	result := &Response{
		Data: &PageListResponse{},
	}
	_, err = res.Read(clients.JSON(result))
	return result.Data.(*PageListResponse), err
}
