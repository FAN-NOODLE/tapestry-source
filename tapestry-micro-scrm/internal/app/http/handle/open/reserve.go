package open

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	"micro-scrm/internal/app/http/middleware"
	openReq "micro-scrm/internal/app/http/request/open"
	openSvc "micro-scrm/internal/app/service/open"
)

type ReserveHandle struct {
	handle.BaseHandle
	reserveSvc *openSvc.ReserveService
}

func NewReserveHandle() *ReserveHandle {
	return &ReserveHandle{
		reserveSvc: openSvc.NewReserve(),
	}
}

// AssignReserveGuide 分配预约导购
func (h ReserveHandle) AssignReserveGuide(c *gin.Context) {
	var params openReq.GenerateReserveRequest
	if err := c.ShouldBind(&params); err != nil {
		h.ResError(c, err)
		return
	}
	params.EnterpriseId = c.GetInt64(middleware.EnterpriseId)
	reserve, err := h.reserveSvc.GenerateReserve(h.WarpContext(c), params)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, reserve)
	return
}

// AssignGuide 分配导购
func (h ReserveHandle) AssignGuide(c *gin.Context) {
	var params openReq.AssignGuideRequest
	if err := c.ShouldBind(&params); err != nil {
		h.ResError(c, err)
		return
	}
	params.EnterpriseId = c.GetInt64(middleware.EnterpriseId)
	reserve, err := h.reserveSvc.AssignGuide(h.WarpContext(c), params)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, reserve)
	return
}

// GetReserve 获取预约信息
func (h ReserveHandle) GetReserve(c *gin.Context) {
	var params openReq.GetRelationRequest
	if err := c.ShouldBind(&params); err != nil {
		h.ResError(c, err)
		return
	}
	params.EnterpriseId = c.GetInt64(middleware.EnterpriseId)
	reserve, err := h.reserveSvc.GetRelation(h.WarpContext(c), params)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, reserve)
	return
}

func (h ReserveHandle) UserLists(c *gin.Context) {
	var params openReq.UserReserveListReq

	if err := c.ShouldBind(&params); err != nil {
		h.ResError(c, err)
		return
	}
	params.EnterpriseId = c.GetInt64(middleware.EnterpriseId)

	lists, total, err := h.reserveSvc.UserReserveList(h.WarpContext(c), params)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResList(c, total, params.Page, params.Size, lists)
	return
}
