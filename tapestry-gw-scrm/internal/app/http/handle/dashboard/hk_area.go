package dashboard

import (
	"github.com/gin-gonic/gin"
	"gw-scrm/internal/app/http/handle"
	dashboardReq "gw-scrm/internal/app/http/request/dashboard"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type HkAreaHandle struct {
	handle.BaseHandle
	api microscrm.DashboardApi
}

func NewHkAreaHandle() HkAreaHandle {
	return HkAreaHandle{
		api: microscrm.NewDashboardApi(),
	}
}

func (h HkAreaHandle) List(c *gin.Context) {
	var req dashboardReq.AreaTree
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseID, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseID

	res, err := h.api.HkAreaList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}
