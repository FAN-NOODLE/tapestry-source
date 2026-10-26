package store_visit

import (
	"github.com/gin-gonic/gin"
	"gw-scrm/internal/app/http/handle"
	"gw-scrm/internal/app/http/request"
	storeVisitReq "gw-scrm/internal/app/http/request/store_visit"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type StoreVisitMainHandle struct {
	handle.BaseHandle
	storeVisitApi microscrm.StoreVisitApi
}

func NewStoreVisitMainHandle() StoreVisitMainHandle {
	return StoreVisitMainHandle{
		storeVisitApi: microscrm.NewStoreVisitApi(),
	}
}

func (h StoreVisitMainHandle) List(c *gin.Context) {
	var req storeVisitReq.ListReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	res, err := h.storeVisitApi.List(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return

}

func (h StoreVisitMainHandle) Close(c *gin.Context) {
	var req storeVisitReq.StoreVisitBaseReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	res, err := h.storeVisitApi.Close(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return
}

func (h StoreVisitMainHandle) DeptName(c *gin.Context) {
	var req request.EnterpriseBaseReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	res, err := h.storeVisitApi.DeptNames(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return
}

func (h StoreVisitMainHandle) Export(c *gin.Context) {
	var req storeVisitReq.ListReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	res, err := h.storeVisitApi.Export(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h StoreVisitMainHandle) Dashboard(c *gin.Context) {
	var req storeVisitReq.DashboardReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	res, err := h.storeVisitApi.Dashboard(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return
}

func (h StoreVisitMainHandle) DashboardExport(c *gin.Context) {
	var req storeVisitReq.DashboardReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	res, err := h.storeVisitApi.DashboardExport(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return
}
