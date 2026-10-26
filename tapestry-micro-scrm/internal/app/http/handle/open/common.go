package open

import (
	"github.com/gin-gonic/gin"

	"micro-scrm/internal/app/http/handle"
	openReq "micro-scrm/internal/app/http/request/open"
	openSvc "micro-scrm/internal/app/service/open"
)

type CommonHandle struct {
	handle.BaseHandle
	decryptService *openSvc.DecryptService
}

func NewCommonHandle() *CommonHandle {
	return &CommonHandle{decryptService: openSvc.NewDecryptService()}
}

// Decrypt 批量解密 OpenAPI 返回的 PI 字段。
func (h *CommonHandle) Decrypt(c *gin.Context) {
	var request openReq.DecryptReq
	if err := c.ShouldBindJSON(&request); err != nil {
		h.ResError(c, err)
		return
	}
	response, err := h.decryptService.Decrypt(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, response)
}
