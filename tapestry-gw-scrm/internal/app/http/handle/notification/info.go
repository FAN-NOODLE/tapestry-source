package notification

import (
	"github.com/gin-gonic/gin"
	"gw-scrm/internal/app/http/handle"
	"gw-scrm/internal/app/http/request"
	"gw-scrm/internal/app/http/request/notification"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type InfoHandle struct {
	handle.BaseHandle
	notificationSdk microscrm.NotificationApi
}

func NewInfoHandle() InfoHandle {
	return InfoHandle{
		notificationSdk: microscrm.NewNotificationApi(),
	}
}

func (h InfoHandle) Save(c *gin.Context) {
	var req notificationReq.InfoReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, adminId, mobile := common.GetEnterpriseAndMobile(c)
	req.EnterpriseId = enterpriseId
	req.AdminID = adminId
	req.Mobile = mobile
	create, err := h.notificationSdk.Save(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, create)
}

func (h InfoHandle) List(c *gin.Context) {
	var req notificationReq.SearchListReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, adminId := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	req.AdminID = adminId
	create, err := h.notificationSdk.List(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, create)
}

func (h InfoHandle) Info(c *gin.Context) {
	var req notificationReq.SearchId
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, adminId := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	req.AdminID = adminId
	info, err := h.notificationSdk.Info(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h InfoHandle) Departments(c *gin.Context) {
	var req request.EnterpriseBaseReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	info, err := h.notificationSdk.Departments(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h InfoHandle) Copy(c *gin.Context) {
	var req notificationReq.SearchId
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, adminId := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	req.AdminID = adminId
	info, err := h.notificationSdk.Copy(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h InfoHandle) Recall(c *gin.Context) {
	var req notificationReq.SearchId
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, adminId := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	req.AdminID = adminId
	info, err := h.notificationSdk.Recall(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h InfoHandle) Check(c *gin.Context) {
	var req notificationReq.InfoReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, adminId, mobile := common.GetEnterpriseAndMobile(c)
	req.EnterpriseId = enterpriseId
	req.AdminID = adminId
	req.Mobile = mobile
	req.MemberID = 0
	info, err := h.notificationSdk.Check(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h InfoHandle) Manager(c *gin.Context) {
	var req notificationReq.AdminManager
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, adminId := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	req.AdminID = adminId
	info, err := h.notificationSdk.Manager(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h InfoHandle) Cancel(c *gin.Context) {
	var req notificationReq.SearchId
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, adminId := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	req.AdminID = adminId
	info, err := h.notificationSdk.Cancel(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}
