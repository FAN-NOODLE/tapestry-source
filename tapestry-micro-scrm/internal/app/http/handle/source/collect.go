package source

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	sourceRequest "micro-scrm/internal/app/http/request/source"
	sourceService "micro-scrm/internal/app/service/source"
)

type CollectHandle struct {
	handle.BaseHandle
	CollectSvc *sourceService.CollectService
}

func NewCollectHandle() CollectHandle {
	return CollectHandle{
		CollectSvc: sourceService.NewCollectService(),
	}
}

func (h *CollectHandle) Edit(c *gin.Context) {
	var req sourceRequest.CollectReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	res, err := h.CollectSvc.CollectSource(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h *CollectHandle) List(c *gin.Context) {
	var req sourceRequest.CollectListReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	result, total, err := h.CollectSvc.List(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResList(c, int64(total), req.Page, req.Size, result)
}
