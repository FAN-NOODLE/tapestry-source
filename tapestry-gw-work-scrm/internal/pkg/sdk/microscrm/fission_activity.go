package microScrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
)

// MemberApi 员工 SDK
type FissionActivityApi struct {
	BaseSDK
}

// NewMember 返回API实例
func NewFissionActivity(host string) *FissionActivityApi {
	urlHost = host
	return &FissionActivityApi{BaseSDK{c: &clients.HTTPClient{}}}
}

// FissionActivityList 活动裂变列表
func (s FissionActivityApi) FissionActivityList(c ctx.Context, request FissionActivityListRequest) (*FissionActivityListResponse, error) {
	path := "/micro-scrm/v1/fission/qw-list"
	url := urlHost + path
	res, err := s.Post(c, url, request)
	if err != nil {
		return nil, err
	}

	result := &Response{
		Data: &FissionActivityListResponse{},
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	if result.Data == nil {
		return nil, err
	}

	return result.Data.(*FissionActivityListResponse), err
}
