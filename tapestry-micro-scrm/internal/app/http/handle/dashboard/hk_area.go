package dashboard

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	dashboardReq "micro-scrm/internal/app/http/request/dashboard"
	dashboardSvc "micro-scrm/internal/app/service/dashboard"
)

type HkAreaHandle struct {
	handle.BaseHandle
	areaSvc *dashboardSvc.HkAreaService
}

func NewHkAreaHandle() HkAreaHandle {
	return HkAreaHandle{
		areaSvc: dashboardSvc.NewHkAreaService(),
	}
}

func (h HkAreaHandle) List(c *gin.Context) {
	var req dashboardReq.AreaTree
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	list, err := h.areaSvc.GetAreaList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}
