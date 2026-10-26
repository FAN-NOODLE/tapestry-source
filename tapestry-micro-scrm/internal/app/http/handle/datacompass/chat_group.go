package dataCompass

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	dataCompassReq "micro-scrm/internal/app/http/request/datacompass"
	dataCompassSvc "micro-scrm/internal/app/service/datacompass"
)

type ChatGroup struct {
	handle.BaseHandle
}

func NewChatGroupHandle() ChatGroup {
	return ChatGroup{}
}

func (h ChatGroup) List(c *gin.Context) {
	var request dataCompassReq.ChatGroupReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	result, total, err := dataCompassSvc.NewChatGroup().List(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResList(c, total, request.PageNum, request.PageSize, result)
	return
}
