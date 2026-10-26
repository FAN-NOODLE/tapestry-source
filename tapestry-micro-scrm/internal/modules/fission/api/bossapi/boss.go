package bossapi

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/modules/fission"
	"micro-scrm/internal/pkg/sdk/boss"
)

func NewApi() *ApiImpl {
	return &ApiImpl{}
}

type ApiImpl struct {
}

// GetShopAuthInfo 获取商城授权信息
func (a *ApiImpl) GetShopAuthInfo(c ctx.Context, enterpriseID int64, authType int) (*fission.ShopAuthInfo, error) {
	var authRequest = boss.AuthDetailRequest{
		EnterpriseID: enterpriseID,
		AuthType:     authType,
	}

	response, err := boss.AuthDetail(c, &authRequest)
	if err != nil {
		return nil, err
	}
	var authInfo = fission.ShopAuthInfo{
		Name: response.Data.Name,
		Logo: response.Data.Logo,
	}
	return &authInfo, nil
}
