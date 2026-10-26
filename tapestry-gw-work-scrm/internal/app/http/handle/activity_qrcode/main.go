package activityQrCodeHandle

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-work-scrm/internal/app/http/handle"
	activityQRCodeReq "gw-work-scrm/internal/app/http/request/activity_qrcode"
	"gw-work-scrm/internal/pkg/sdk/boss"
	activityQrCodeSdk "gw-work-scrm/internal/pkg/sdk/microscrm/activity_qrcode"
)

type ActivityQrCodeHandle struct {
	handle.BaseHandle
	activityQrCodeSdk *activityQrCodeSdk.ActivityQrCode
}

func NewActivityQrCodeHandle() *ActivityQrCodeHandle {
	return &ActivityQrCodeHandle{
		activityQrCodeSdk: activityQrCodeSdk.NewActivityQrCode(),
	}
}

func (h ActivityQrCodeHandle) GetQrCodeByUnionid(c *gin.Context) {
	var req activityQRCodeReq.GetQrCodeReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseHash := c.Request.Header.Get("enterprise-hash")
	enterpriseReq := boss.EnterpriseDetailRequest{
		Hash: enterpriseHash,
	}
	enterpriseInfo, err := boss.EnterpriseDetail(ctx.Wrap(c), &enterpriseReq)
	req.EnterpriseId = enterpriseInfo.ID
	create, err := h.activityQrCodeSdk.GetQrCodeByUnionid(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, create)
}

func (h ActivityQrCodeHandle) GetAllStoreList(c *gin.Context) {
	var req activityQRCodeReq.GetQrCodeReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseHash := c.Request.Header.Get("enterprise-hash")
	enterpriseReq := boss.EnterpriseDetailRequest{
		Hash: enterpriseHash,
	}
	enterpriseInfo, err := boss.EnterpriseDetail(ctx.Wrap(c), &enterpriseReq)
	req.EnterpriseId = enterpriseInfo.ID
	create, err := h.activityQrCodeSdk.GetAllStoreList(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, create)
}

func (h ActivityQrCodeHandle) ActivityDetail(c *gin.Context) {
	var req activityQRCodeReq.ActivityInfoReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseHash := c.Request.Header.Get("enterprise-hash")
	enterpriseReq := boss.EnterpriseDetailRequest{
		Hash: enterpriseHash,
	}
	enterpriseInfo, err := boss.EnterpriseDetail(ctx.Wrap(c), &enterpriseReq)
	req.EnterpriseId = enterpriseInfo.ID
	create, err := h.activityQrCodeSdk.ActivityDetail(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, create)
}
