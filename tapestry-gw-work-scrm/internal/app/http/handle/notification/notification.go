package notification

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-work-scrm/internal/app/http/handle"
	notificationReq "gw-work-scrm/internal/app/http/request/notification"
	notification2 "gw-work-scrm/internal/pkg/sdk/microscrm/notification"
)

type InfoHandle struct {
	handle.BaseHandle
	notificationSdk *notification2.Info
}

func NewInfoHandle() *InfoHandle {
	return &InfoHandle{
		notificationSdk: notification2.NewInfo(),
	}
}

func (h InfoHandle) InfoList(c *gin.Context) {
	var req notificationReq.SearchListReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	create, err := h.notificationSdk.InfoList(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, create)
}

func (h InfoHandle) MergeList(c *gin.Context) {
	var req notificationReq.InfoMergeListReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	req.IsH5 = true
	create, err := h.notificationSdk.MergeList(ctx.Wrap(c), req)
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
	enterpriseId := h.EnterpriseID(c)
	req.EnterpriseId = enterpriseId
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	info, err := h.notificationSdk.Info(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h InfoHandle) MergeDetail(c *gin.Context) {
	var req notificationReq.InfoMergeDetailReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId := h.EnterpriseID(c)
	req.EnterpriseId = enterpriseId
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	req.IsH5 = true
	info, err := h.notificationSdk.MergeDetail(ctx.Wrap(c), req)
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
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	info, err := h.notificationSdk.Recall(ctx.Wrap(c), req)
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
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	req.IsH5 = true
	info, err := h.notificationSdk.Check(ctx.Wrap(c), req)
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
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	info, err := h.notificationSdk.Manager(ctx.Wrap(c), req)
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
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	info, err := h.notificationSdk.Cancel(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h InfoHandle) SendFileByEmail(c *gin.Context) {
	var req notificationReq.SendFileByEmailReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.EnterpriseId = h.EnterpriseID(c)
	info, err := h.notificationSdk.SendFileByEmail(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h InfoHandle) ListReadNum(c *gin.Context) {
	var req notificationReq.AdminManager
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.MemberId = memberInfo.ID
	info, err := h.notificationSdk.ListReadNum(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h InfoHandle) SavePdf(c *gin.Context) {
	var req notificationReq.SavePdfReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.EnterpriseId = h.EnterpriseID(c)
	info, err := h.notificationSdk.SavePdf(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)

}
