package microwork

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"

type EcHub struct {
}

func NewEcHub() *EcHub {
	return &EcHub{}
}

// SaAssignment sa绑定
func (c EcHub) SaAssignment(ctx ctx.Context, request SaAssignmentEcHubReq) (SaAssignmentEcHubResp, error) {
	var response SaAssignmentEcHubResp
	err := httpClient.Post(ctx, GetURL("%s/service/ec-hub/sa-assignment"), request, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

// SaUnAssignment sa取消绑定
func (c EcHub) SaUnAssignment(ctx ctx.Context, request SaUnAssignmentEcHubReq) (SaAssignmentEcHubResp, error) {
	var response SaAssignmentEcHubResp
	err := httpClient.Post(ctx, GetURL("%s/service/ec-hub/sa-un-assignment"), request, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

// GetQrCodeAndLink  获取EC太阳码和link
func (c EcHub) GetQrCodeAndLink(ctx ctx.Context, request SunQRCodeReq) (SunQRCodeResp, error) {
	var response SunQRCodeResp
	err := httpClient.Post(ctx, GetURL("%s/service/ec-hub/get-qrcode-and-link"), request, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}
