package notification

import (
	"errors"
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	notificationReq "micro-scrm/internal/app/http/request/notification"
	"micro-scrm/internal/app/service/notification"
)

type InfoMergeHandle struct {
	handle.BaseHandle
	mergeInfoSvc *notificationSvc.InfoMergeSvc
}

func NewInfoMergeHandle() InfoMergeHandle {
	return InfoMergeHandle{
		mergeInfoSvc: notificationSvc.NewInfoMergeSvc(),
	}
}

func (h InfoMergeHandle) Save(c *gin.Context) {
	var req notificationReq.InfoMergeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.mergeInfoSvc.Save(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h InfoMergeHandle) List(c *gin.Context) {
	var req notificationReq.InfoMergeListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	ctx := h.WarpContext(c)
	info, err := notificationSvc.NewInfoSvc().IsRoRole(ctx, req.AdminManager)
	if !info {
		h.ResError(c, errors.New("该账号无RO权限"))
		return
	}

	list, total, err := h.mergeInfoSvc.List(ctx, req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResList(c, total, req.Page, req.Size, list)
}

func (h InfoMergeHandle) Detail(c *gin.Context) {
	var req notificationReq.InfoMergeDetailReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	info, err := h.mergeInfoSvc.Detail(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, info)
}

func (h InfoMergeHandle) Release(c *gin.Context) {
	var req notificationReq.InfoMergeDetailReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.mergeInfoSvc.Release(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResOk(c)
}

func (h InfoMergeHandle) ListShow(c *gin.Context) {
	var req notificationReq.InfoMergeListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	list, total, err := h.mergeInfoSvc.ListShow(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResList(c, total, req.Page, req.Size, list)
}

func (h InfoMergeHandle) SavePdf(c *gin.Context) {
	var req notificationReq.SavePdfReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.mergeInfoSvc.SavePdf(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}
