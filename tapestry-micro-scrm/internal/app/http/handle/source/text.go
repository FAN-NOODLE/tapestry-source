package source

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	sourceRequest "micro-scrm/internal/app/http/request/source"
	sourceService "micro-scrm/internal/app/service/source"
)

type TextHandle struct {
	handle.BaseHandle
	textService sourceService.TextService
}

func NewTextHandle() TextHandle {
	return TextHandle{
		textService: sourceService.NewTextService(),
	}
}

func (h TextHandle) Detail(c *gin.Context) {
	var request sourceRequest.DetailRequest
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	res, err := h.textService.Detail(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return
}

func (h TextHandle) Update(c *gin.Context) {
	var request sourceRequest.TextUpdateRequest
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	if err := h.textService.Update(h.WarpContext(c), request); err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
	return
}

func (h TextHandle) GetList(c *gin.Context) {
	var request sourceRequest.TextUpdateRequest
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	if err := h.textService.Update(h.WarpContext(c), request); err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
	return
}
