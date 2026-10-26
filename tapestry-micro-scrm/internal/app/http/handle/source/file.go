package source

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	sourceRequest "micro-scrm/internal/app/http/request/source"
	sourceService "micro-scrm/internal/app/service/source"
)

type FileHandle struct {
	handle.BaseHandle
	fileService sourceService.FileService
}

func NewFileHandle() FileHandle {
	return FileHandle{
		fileService: sourceService.NewFileService(),
	}
}

func (h FileHandle) Detail(c *gin.Context) {
	var request sourceRequest.DetailRequest
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	res, err := h.fileService.Detail(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return
}

func (h FileHandle) Update(c *gin.Context) {
	var request sourceRequest.FileUpdateRequest
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	if err := h.fileService.Update(h.WarpContext(c), request); err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
	return
}
