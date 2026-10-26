package notification

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	request "micro-scrm/internal/app/http/request/notification"
	"micro-scrm/internal/app/service/notification"
)

type NoticeConfigHandle struct {
	handle.BaseHandle
	crowdSvc *notificationSvc.NoticeConfigSvc
}

func NewNoticeConfigHandle() NoticeConfigHandle {
	return NoticeConfigHandle{
		crowdSvc: notificationSvc.NewNoticeConfigSvc(),
	}
}

func (h NoticeConfigHandle) Save(c *gin.Context) {
	var req request.NoticeConfigSaveReq
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

func (h NoticeConfigHandle) Detail(c *gin.Context) {
	var req request.NoticeConfigDetailReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	info, err := h.crowdSvc.Detail(h.WarpContext(c), req.EnterpriseId)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, info)
}
