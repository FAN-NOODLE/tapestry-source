package member

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-work-scrm/internal/app/http/common"
	"gw-work-scrm/internal/app/http/handle"
	"gw-work-scrm/internal/app/http/request"
	"gw-work-scrm/internal/app/http/request/member"
	microScrmCustomer "gw-work-scrm/internal/pkg/sdk/microscrm/customer"
	"gw-work-scrm/internal/pkg/sdk/work"
)

type MemberHandle struct {
	handle.BaseHandle
	customerApiSdk microScrmCustomer.CustomerApi
}

func NewMemberHandle() MemberHandle {
	return MemberHandle{
		customerApiSdk: microScrmCustomer.NewCustomerApi(),
	}
}

// GuidePosition 导购身份和权限类型判定
func (h MemberHandle) GuidePosition(c *gin.Context) {

	var req member.GuidePositionReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.Mobile = memberInfo.Mobile
	req.MemberType = memberInfo.Type
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}

	var role int
	// 企业导购
	if req.MemberType == 3 {
		role = 3
	}
	// 企业管理
	if req.MemberType == 4 {
		role = 5
	}

	var res member.GuidePositionRes

	res.Role = role

	h.ResSuccess(c, res)
	return

}

func (h MemberHandle) BaseInfo(c *gin.Context) {
	memberInfo := h.MemberInfo(c)
	h.ResSuccess(c, gin.H{
		"member_id":   memberInfo.ID,
		"member_name": memberInfo.Name,
		"avatar":      memberInfo.Avatar,
		"user_id":     memberInfo.Userid,
	})
}

func (h MemberHandle) FirstOpen(c *gin.Context) {
	var req member.FirstOpenReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.Mobile = memberInfo.Mobile
	info, err := h.customerApiSdk.FirstOpen(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, info)
}

func (h MemberHandle) GuideList(c *gin.Context) {
	var req member.GuideListReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.Mobile = memberInfo.Mobile
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	cc := ctx.Wrap(c)
	if memberInfo.Type != 4 {
		h.ResError(c, errors.New("身份无权限"))
		return
	}
	if len(req.DpIds) == 0 {
		thirdIds, err := common.NewBaseHandle().GetMemberThirdsByType(cc, req.EnterpriseId, req.Mobile)
		if err != nil {
			h.ResError(c, err)
			return
		}
		req.DpIds = thirdIds
	}
	req.Domain = "domain"
	list, err := work.GetGuideList(cc, req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}

func (h MemberHandle) StoreOrDepartmentList(c *gin.Context) {
	var req member.StoreOrDepartmentsReq
	req.EnterpriseId = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	req.Mobile = memberInfo.Mobile
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	cc := ctx.Wrap(c)
	if memberInfo.Type != 4 {
		h.ResError(c, errors.New("身份无权限"))
		return
	}
	if len(req.DpIds) == 0 {
		thirdIds, err := common.NewBaseHandle().GetMemberThirdsByType(cc, req.EnterpriseId, req.Mobile)
		if err != nil {
			h.ResError(c, err)
			return
		}
		req.DpIds = thirdIds
	}
	list, err := work.GetStoreOrDepartmentList(cc, req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}

func (h MemberHandle) GetFiscalMonth(c *gin.Context) {
	var req request.EnterpriseBaseReq
	req.EnterpriseId = h.EnterpriseID(c)
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	list, err := work.GetFiscalMonth(ctx.Wrap(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}
