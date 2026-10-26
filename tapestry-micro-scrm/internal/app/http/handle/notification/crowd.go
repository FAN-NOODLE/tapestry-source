package notification

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	request "micro-scrm/internal/app/http/request/notification"
	"micro-scrm/internal/app/service/notification"
)

type CrowdHandle struct {
	handle.BaseHandle
	crowdSvc *notificationSvc.CrowdSvc
}

func NewCrowdHandle() CrowdHandle {
	return CrowdHandle{
		crowdSvc: notificationSvc.NewCrowdSvc(),
	}
}

func (h CrowdHandle) Save(c *gin.Context) {
	var req request.CrowdReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.crowdSvc.Save(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h CrowdHandle) List(c *gin.Context) {
	var req request.CrowdListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	list, total, err := h.crowdSvc.List(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResList(c, total, req.Page, req.Size, list)
}

func (h CrowdHandle) Detail(c *gin.Context) {
	var req request.CrowdDetailReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	info, err := h.crowdSvc.Detail(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, info)
}

func (h CrowdHandle) DeleteCrowd(c *gin.Context) {
	var req request.CrowdDetailReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.crowdSvc.DeleteCrowd(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, "")
}
