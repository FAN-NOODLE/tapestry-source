package microwork

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"

type QRCode struct {
}

func NewQRCode() *QRCode {
	return &QRCode{}
}

// GetSunQRCode 获取太阳码
func (s QRCode) GetSunQRCode(c ctx.Context, request SunQRCodeReq) (SunQRCodeResp, error) {
	var response SunQRCodeResp
	err := httpClient.Post(c, GetURL("%s/service/v1/common/get-qrcode-and-link"), request, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}
