package dashboard

import (
	"github.com/gin-gonic/gin"
	"gw-scrm/internal/app/http/handle"
	dashboardReq "gw-scrm/internal/app/http/request/dashboard"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type HkStoreHandle struct {
	handle.BaseHandle
	api microscrm.DashboardApi
}

func NewHkStoreHandle() HkStoreHandle {
	return HkStoreHandle{
		api: microscrm.NewDashboardApi(),
	}
}

func (h HkStoreHandle) Sync(c *gin.Context) {
	var req dashboardReq.HkStoreSync
	enterpriseID, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseID

	res, err := h.api.HkStoreSync(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h HkStoreHandle) LastSyncTime(c *gin.Context) {
	var req dashboardReq.HkStoreLastSyncTime
	enterpriseID, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseID

	res, err := h.api.HkStoreLastSyncTime(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h HkStoreHandle) List(c *gin.Context) {
	var req dashboardReq.HkStoreList
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

	res, err := h.api.HkStoreList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h HkStoreHandle) Detail(c *gin.Context) {
	var req dashboardReq.HkStoreDetail
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseID, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseID

	res, err := h.api.HkStoreDetail(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h HkStoreHandle) Save(c *gin.Context) {
	var req dashboardReq.HkStoreSave
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseID, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseID

	res, err := h.api.HkStoreSave(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h HkStoreHandle) Close(c *gin.Context) {
	var req dashboardReq.HkStoreClose
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseID, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseID

	res, err := h.api.HkStoreClose(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}
