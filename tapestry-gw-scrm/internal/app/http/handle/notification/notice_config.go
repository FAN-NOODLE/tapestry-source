package notification

import (
	"github.com/gin-gonic/gin"
	"gw-scrm/internal/app/http/handle"
	notificationReq "gw-scrm/internal/app/http/request/notification"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type NoticeConfigHandle struct {
	handle.BaseHandle
	notificationSdk microscrm.NotificationApi
}

func NewNoticeConfigHandle() NoticeConfigHandle {
	return NoticeConfigHandle{
		notificationSdk: microscrm.NewNotificationApi(),
	}
}

func (h NoticeConfigHandle) Save(c *gin.Context) {
	var req notificationReq.NoticeConfigSaveReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	res, err := h.notificationSdk.NoticeConfigSave(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h NoticeConfigHandle) Detail(c *gin.Context) {
	var req notificationReq.NoticeConfigDetailReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId

	res, err := h.notificationSdk.NoticeConfigDetail(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}
