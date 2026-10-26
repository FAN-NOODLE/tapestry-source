package source

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	sourceRequest "micro-scrm/internal/app/http/request/source"
	sourceService "micro-scrm/internal/app/service/source"
)

type ImageHandle struct {
	handle.BaseHandle
	imageService sourceService.ImageService
}

func NewImageHandle() ImageHandle {
	return ImageHandle{
		imageService: sourceService.NewImageService(),
	}
}

func (h ImageHandle) Detail(c *gin.Context) {
	var request sourceRequest.DetailRequest
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	res, err := h.imageService.Detail(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return
}

func (h ImageHandle) Update(c *gin.Context) {
	var request sourceRequest.ImageUpdateRequest
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	if err := h.imageService.Update(h.WarpContext(c), request); err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
	return
}
