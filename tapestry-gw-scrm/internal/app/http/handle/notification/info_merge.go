package notification

import (
	"github.com/gin-gonic/gin"
	"gw-scrm/internal/app/http/handle"
	notificationReq "gw-scrm/internal/app/http/request/notification"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type InfoMergeHandle struct {
	handle.BaseHandle
	notificationSdk microscrm.NotificationApi
}

func NewInfoMergeHandle() InfoMergeHandle {
	return InfoMergeHandle{
		notificationSdk: microscrm.NewNotificationApi(),
	}
}

func (h InfoMergeHandle) Save(c *gin.Context) {
	var req notificationReq.InfoMergeReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, adminId, mobile := common.GetEnterpriseAndMobile(c)
	req.EnterpriseId = enterpriseId
	req.AdminId = adminId
	req.Mobile = mobile
	res, err := h.notificationSdk.InfoMergeSave(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h InfoMergeHandle) Detail(c *gin.Context) {
	var req notificationReq.InfoMergeDetailReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, adminId, _ := common.GetEnterpriseAndMobile(c)
	req.EnterpriseId = enterpriseId
	req.AdminID = adminId

	res, err := h.notificationSdk.InfoMergeDetail(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h InfoMergeHandle) List(c *gin.Context) {
	var req notificationReq.InfoMergeListReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, adminId, _ := common.GetEnterpriseAndMobile(c)
	req.EnterpriseId = enterpriseId
	req.AdminID = adminId
	if req.Page == 0 {
		req.Page = 1
	}

	if req.Size == 0 {
		req.Size = 10
	}
	res, err := h.notificationSdk.InfoMergeList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h InfoMergeHandle) Release(c *gin.Context) {
	var req notificationReq.InfoMergeDetailReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, adminId, mobile := common.GetEnterpriseAndMobile(c)
	req.EnterpriseId = enterpriseId
	req.AdminID = adminId
	req.Mobile = mobile
	res, err := h.notificationSdk.InfoMergeRelease(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h InfoMergeHandle) SavePdf(c *gin.Context) {
	var req notificationReq.SavePdfReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	res, err := h.notificationSdk.SavePdf(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}
