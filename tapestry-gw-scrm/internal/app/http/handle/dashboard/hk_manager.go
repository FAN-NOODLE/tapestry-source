package dashboard

import (
	"github.com/gin-gonic/gin"
	"gw-scrm/internal/app/http/handle"
	"gw-scrm/internal/app/http/request"
	dashboardReq "gw-scrm/internal/app/http/request/dashboard"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type HkManagerHandle struct {
	handle.BaseHandle
	api microscrm.DashboardApi
}

func NewHkManagerHandle() HkManagerHandle {
	return HkManagerHandle{
		api: microscrm.NewDashboardApi(),
	}
}

func (h HkManagerHandle) StoreList(c *gin.Context) {
	var req dashboardReq.HkManagerStoreList
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseID, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseID

	res, err := h.api.HkManagerStoreList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h HkManagerHandle) AuthList(c *gin.Context) {
	var req dashboardReq.HkManagerAuthList
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Size == 0 {
		req.Size = 10
	}
	enterpriseID, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseID

	res, err := h.api.HkManagerAuthList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h HkManagerHandle) AuthSave(c *gin.Context) {
	var req dashboardReq.HkManagerAuthSave
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseID, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseID

	res, err := h.api.HkManagerAuthSave(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h HkManagerHandle) AuthDetail(c *gin.Context) {
	var req dashboardReq.HkManagerAuthDetail
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseID, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseID

	res, err := h.api.HkManagerAuthDetail(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h HkManagerHandle) AuthDelete(c *gin.Context) {
	var req dashboardReq.HkManagerRemove
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseID, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseID

	res, err := h.api.HkManagerAuthDelete(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h HkManagerHandle) AllMemberIds(c *gin.Context) {
	var req request.EnterpriseBaseReq
	enterpriseID, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseID

	res, err := h.api.HkManagerAllMemberIds(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}
