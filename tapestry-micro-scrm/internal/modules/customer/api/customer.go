package api

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/pkg/sdk/microwork"
)

type CustomerApi struct {
}

func NewCustomerApi() *CustomerApi {
	return &CustomerApi{}
}

// GetUserInfoByUnionId 查询会员
func (m CustomerApi) GetUserInfoByUnionId(c ctx.Context, req microwork.UserInfoRequest) *microwork.UserInfoResponse {
	res, err := microwork.NewCrm().GetUserInfoByUnionId(c, req)
	fmt.Println(res)
	if err != nil {
		return nil
	}
	return &res
}




