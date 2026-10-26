package invite

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-work-scrm/internal/app/http/handle"
	"gw-work-scrm/internal/app/http/request/invite"
	invite2 "gw-work-scrm/internal/pkg/sdk/microscrm/invite"
)

type Handle struct {
	handle.BaseHandle
	inviteSdk *invite2.Api
}

func NewInviteHandle() Handle {
	return Handle{
		BaseHandle: handle.BaseHandle{},
		inviteSdk:  invite2.NewInviteApi(),
	}
}

// ConfigList 邀请函配置页列表
func (h Handle) ConfigList(c *gin.Context) {
	var request invite.TemplateConfigListReq
	request.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	request.Mobile = memberInfo.Mobile
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	list, err := h.inviteSdk.ConfigList(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}

// ConfigDetail 邀请函模板配置详情
func (h Handle) ConfigDetail(c *gin.Context) {
	var request invite.TemplateConfigDetailReq
	request.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	request.Mobile = memberInfo.Mobile
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	info, err := h.inviteSdk.ConfigDetail(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

// CreateRecord 发送邀约单记录
func (h Handle) CreateRecord(c *gin.Context) {
	var request invite.RecordInfoReq
	request.EnterpriseId = h.EnterpriseID(c)
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	memberInfo := h.MemberInfo(c)
	request.Mobile = memberInfo.Mobile
	info, err := h.inviteSdk.CreateRecord(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

// RecordInfo 用户查看邀请函信息
func (h Handle) RecordInfo(c *gin.Context) {
	var request invite.GetRecordInfoReq
	request.EnterpriseId = h.EnterpriseID(c)
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	info, err := h.inviteSdk.RecordInfo(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h Handle) RecordDelete(c *gin.Context) {
	var request invite.GetRecordInfoReq
	request.EnterpriseId = h.EnterpriseID(c)
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	info, err := h.inviteSdk.RecordDelete(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}
