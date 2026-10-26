package source

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	sourceRequest "micro-scrm/internal/app/http/request/source"
	sourceService "micro-scrm/internal/app/service/source"
)

type MiniProgramHandle struct {
	handle.BaseHandle
	miniProgramService sourceService.MiniProgramService
}

func NewMiniProgramHandle() MiniProgramHandle {
	return MiniProgramHandle{
		miniProgramService: sourceService.NewMiniProgramService(),
	}
}

func (h MiniProgramHandle) Detail(c *gin.Context) {
	var request sourceRequest.DetailRequest
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	res, err := h.miniProgramService.Detail(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return
}

func (h MiniProgramHandle) Update(c *gin.Context) {
	var request sourceRequest.MiniProgramUpdateRequest
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	if err := h.miniProgramService.Update(h.WarpContext(c), request); err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
	return
}
