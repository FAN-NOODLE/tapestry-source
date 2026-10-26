package microwork

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"

type Crm struct {
}

func NewCrm() *Crm {
	return &Crm{}
}

// GetUserInfoByUnionId 查询用户信息
func (s Crm) GetUserInfoByUnionId(c ctx.Context, request UserInfoRequest) (UserInfoResponse, error) {
	var response UserInfoResponse
	err := httpClient.Post(c, GetURL("%s/service/crm/get-member-info"), request, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

// CheckMemberByUnionId 检测会员
func (s Crm) CheckMemberByUnionId(c ctx.Context, request UserInfoRequest) (CheckMemberResponse, error) {
	var response CheckMemberResponse
	err := httpClient.Post(c, GetURL("%s/service/crm/check-member"), request, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

// WelfareHaved 用户是否领券
func (s Crm) WelfareHaved(c ctx.Context, request WelfareHavedRequest) (WelfareHavedResponse, error) {
	var response WelfareHavedResponse
	err := httpClient.Post(c, GetURL("%s/service/crm/welfare-haved"), request, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

// GetCustomerPurchase crm查询客户订单接口
func (s Crm) GetCustomerPurchase(c ctx.Context, request CustomerPurchaseReq) (CustomerPurchaseResp, error) {
	var response CustomerPurchaseResp
	err := httpClient.Post(c, GetURL("%s/service/crm/get-customer-purchase"), request, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

// GetCustomerInfoByUnionId 根据unionid或者customerid获取用户信息
func (s Crm) GetCustomerInfoByUnionId(c ctx.Context, request UserInfoRequest) (UserInfoResponse, error) {
	var response UserInfoResponse
	err := httpClient.Post(c, GetURL("%s/service/crm/customer-info"), request, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}
