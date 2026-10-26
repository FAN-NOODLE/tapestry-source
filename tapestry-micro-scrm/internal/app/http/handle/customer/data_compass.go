package customer

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	customerReq "micro-scrm/internal/app/http/request/customer"
	customerSvc "micro-scrm/internal/app/service/customer"
)

type DataCompassHandle struct {
	handle.BaseHandle
	CompassSvc *customerSvc.DataCompassSvc
}

var _ DataCompassHandle

func NewDataCompassHandle() *DataCompassHandle {
	return &DataCompassHandle{
		BaseHandle: handle.BaseHandle{},
		CompassSvc: customerSvc.NewDataCompassSvc(),
	}
}

// ChatGroupSummary 获取小程序数据罗盘的社群数据统计
func (h *DataCompassHandle) ChatGroupSummary(c *gin.Context) {
	var request customerReq.CompassReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	result, err := h.CompassSvc.ChatGroupCompass(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, result)
}

// FriendsSummary 获取小程序数据罗盘,好友数据的统计
func (h *DataCompassHandle) FriendsSummary(c *gin.Context) {
	var request customerReq.CompassReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	result, err := h.CompassSvc.FriendsCompass(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, result)
}
