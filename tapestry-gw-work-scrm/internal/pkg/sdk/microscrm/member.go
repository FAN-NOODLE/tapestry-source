package microScrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
)

// MemberApi 员工 SDK
type MemberApi struct {
	BaseSDK
}

// NewMember 返回API实例
func NewMember(host string) *MemberApi {
	urlHost = host
	return &MemberApi{BaseSDK{c: &clients.HTTPClient{}}}
}

// AuthInfo 获取授权信息
func (m MemberApi) AuthInfo(c ctx.Context, request AuthInfoRequest) (*AuthInfoResponse, error) {
	path := "/micro-scrm/v1/members/auth-info"
	url := urlHost + path
	res, err := m.Post(c, url, request)
	if err != nil {
		return nil, err
	}

	result := &Response{
		Data: &AuthInfoResponse{},
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	if result.Data == nil {
		return nil, err
	}

	return result.Data.(*AuthInfoResponse), err
}

// DetailByUserid 通过用户userid 查询详情
func (m MemberApi) DetailByUserid(c ctx.Context, request DetailUseridRequest) (*Member, error) {
	path := "/micro-scrm/v1/members/detail-userid"
	url := urlHost + path
	res, err := m.Post(c, url, request)
	if err != nil {
		return nil, err
	}

	result := &Response{
		Data: &Member{},
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	if result.Data == nil {
		return nil, err
	}

	return result.Data.(*Member), err
}
