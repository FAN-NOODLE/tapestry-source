package display

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-work-scrm/internal/app/http/handle"
	"gw-work-scrm/internal/app/http/request/display"
	display2 "gw-work-scrm/internal/pkg/sdk/microscrm/display"
)

type DashboardHandle struct {
	handle.BaseHandle
	dashboardSdk *display2.Dashboard
}

func NewDashboardHandle() *DashboardHandle {
	return &DashboardHandle{
		dashboardSdk: display2.NewDashboard(),
	}
}

func (h DashboardHandle) Tree(c *gin.Context) {
	var request display.TreeReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	request.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	request.MemberId = memberInfo.ID
	request.Language = c.Request.Header.Get("Accept-Language")
	info, err := h.dashboardSdk.Tree(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h DashboardHandle) SearchStore(c *gin.Context) {
	var request display.SearchStore
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	request.EnterpriseId = h.EnterpriseID(c)
	request.Language = c.Request.Header.Get("Accept-Language")
	info, err := h.dashboardSdk.SearchStore(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h DashboardHandle) MemberAuth(c *gin.Context) {
	var request display.MemberAuthDetail
	request.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	request.MemberId = memberInfo.ID
	info, err := h.dashboardSdk.MemberAuth(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}
