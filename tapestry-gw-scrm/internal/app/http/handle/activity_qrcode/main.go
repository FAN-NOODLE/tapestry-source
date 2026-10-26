package activity_qrcode

import (
	"github.com/gin-gonic/gin"
	"gw-scrm/internal/app/http/handle"
	ActivityQRCodeReq "gw-scrm/internal/app/http/request/activity_qrcode"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type ActivityQRCodeMainHandle struct {
	handle.BaseHandle
	ActivityQRCodeApi microscrm.ActivityQRCodeApi
}

func NewActivityQRCodeMainHandle() ActivityQRCodeMainHandle {
	return ActivityQRCodeMainHandle{
		ActivityQRCodeApi: microscrm.NewActivityQRCodeApi(),
	}
}

func (h ActivityQRCodeMainHandle) ActivitySave(c *gin.Context) {
	var req ActivityQRCodeReq.ActivitySaveReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	res, err := h.ActivityQRCodeApi.ActivitySave(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return
}

func (h ActivityQRCodeMainHandle) ActivityDetail(c *gin.Context) {
	var req ActivityQRCodeReq.InfoReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	res, err := h.ActivityQRCodeApi.ActivityDetail(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return
}

func (h ActivityQRCodeMainHandle) ActivityDelete(c *gin.Context) {
	var req ActivityQRCodeReq.InfoReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	res, err := h.ActivityQRCodeApi.ActivityDelete(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return
}

func (h ActivityQRCodeMainHandle) ActivityList(c *gin.Context) {
	var req ActivityQRCodeReq.ActivityListReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	res, err := h.ActivityQRCodeApi.ActivityList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return
}

func (h ActivityQRCodeMainHandle) ActivityCategoryList(c *gin.Context) {
	var req ActivityQRCodeReq.ActivityCategoryListReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	res, err := h.ActivityQRCodeApi.ActivityCategoryList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return
}

func (h ActivityQRCodeMainHandle) WelcomeSave(c *gin.Context) {
	var req ActivityQRCodeReq.WelcomeInfoReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	res, err := h.ActivityQRCodeApi.WelcomeSave(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return
}

func (h ActivityQRCodeMainHandle) WelcomeDetail(c *gin.Context) {
	var req ActivityQRCodeReq.InfoReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	res, err := h.ActivityQRCodeApi.WelcomeDetail(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return
}

func (h ActivityQRCodeMainHandle) WelcomeDelete(c *gin.Context) {
	var req ActivityQRCodeReq.InfoReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	res, err := h.ActivityQRCodeApi.WelcomeDelete(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return
}

func (h ActivityQRCodeMainHandle) WelcomeList(c *gin.Context) {
	var req ActivityQRCodeReq.WelcomeListReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	res, err := h.ActivityQRCodeApi.WelcomeList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return
}

func (h ActivityQRCodeMainHandle) Statistics(c *gin.Context) {
	var req ActivityQRCodeReq.StatisticsReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	res, err := h.ActivityQRCodeApi.Statistics(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return
}
