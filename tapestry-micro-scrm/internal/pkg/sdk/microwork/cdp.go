package microwork

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

type Cdp struct {
}

func NewCdp() *Cdp {
	return &Cdp{}
}

// SaAssignment sa绑定
func (c Cdp) SaAssignment(ctx ctx.Context, request SaAssignmentReq) (SaAssignmentResp, error) {
	var response SaAssignmentResp
	err := httpClient.Post(ctx, GetURL("%s/service/crm/sa-assignment"), request, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

// SaUnAssignment sa取消绑定
func (c Cdp) SaUnAssignment(ctx ctx.Context, request SaUnAssignmentReq) (SaAssignmentResp, error) {
	var response SaAssignmentResp
	err := httpClient.Post(ctx, GetURL("%s/service/crm/sa-un-assignment"), request, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}
