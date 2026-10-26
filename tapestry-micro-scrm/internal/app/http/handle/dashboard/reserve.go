package dashboard

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	dashboardReq "micro-scrm/internal/app/http/request/dashboard"
	dashboardSvc "micro-scrm/internal/app/service/dashboard"
)

type ReserveHandle struct {
	handle.BaseHandle
}

func NewReserveHandle() ReserveHandle {
	return ReserveHandle{}
}

func (s ReserveHandle) DataKanban(c *gin.Context) {
	var request dashboardReq.DataKanban
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}

	result, err := dashboardSvc.NewReserveService().DataKanban()

	if err != nil {
		s.ResError(c, err)
		return
	}

	s.ResSuccess(c, result)

	return
}

func (s ReserveHandle) DataList(c *gin.Context) {
	var request dashboardReq.DataList
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}

	result, total, err := dashboardSvc.NewReserveService().DataList()
	if err != nil {
		s.ResError(c, err)
		return
	}

	s.ResList(c, total, request.Page, request.Size, result)
	return
}
