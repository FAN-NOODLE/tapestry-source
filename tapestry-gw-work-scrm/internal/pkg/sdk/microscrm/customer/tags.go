package microScrmCustomer

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	microScrm "gw-work-scrm/internal/pkg/sdk/microscrm"
)

type CustomerTagsApi struct {
	microScrm.BaseSDK
}

func NewCustomerTagsApi() CustomerTagsApi {
	return CustomerTagsApi{}
}

func (a CustomerTagsApi) Operate(c ctx.Context, request interface{}) (interface{}, error) {
	return a.Post(c, a.GetUrl("%s/micro-scrm/v1/customer/tags/impression/operate"), request)
}
