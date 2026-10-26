package customer

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	customerReq "micro-scrm/internal/app/http/request/customer"
	customerSvc "micro-scrm/internal/app/service/customer"
)

type OnInheritHandle struct {
	handle.BaseHandle
	onInheritService *customerSvc.OnInherit
}

func NewOnInheritHandle() OnInheritHandle {
	return OnInheritHandle{
		onInheritService: customerSvc.NewOnInherit(),
	}
}

// GetTakeoverMemberList 获取可分配的导购列表
func (h OnInheritHandle) GetTakeoverMemberList(c *gin.Context) {
	var request customerReq.OnInheritReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	cc := h.WarpContext(c)
	result, err := h.onInheritService.GetTakeoverMemberList(cc, request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, result)
}

// GetSelectAllNum 在职分配的全选人员统计
func (h OnInheritHandle) GetSelectAllNum(c *gin.Context) {
	var request customerReq.OnInheritReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	cc := h.WarpContext(c)
	num, err := h.onInheritService.GetSelectAllNum(cc, request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, num)
}

// TransferCount  在职分配的全选人员统计(新版）
func (h OnInheritHandle) TransferCount(c *gin.Context) {
	var request customerReq.OnInheritSearchReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	cc := h.WarpContext(c)
	memberIds, err := customerSvc.NewExternalSearch().GetMembersByManageType(cc, request.EnterpriseID, request.MemberId, request.ThirdIds, request.MemberIds)
	request.MemberIds = memberIds
	num, err := h.onInheritService.TransferCount(cc, request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, num)
}

func (h OnInheritHandle) GetSelectList(c *gin.Context) {
	var request customerReq.OnInheritReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	cc := h.WarpContext(c)
	num, err := h.onInheritService.GetSelectAllList(cc, request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, num)
}

// TransferList 可转移列表
func (h OnInheritHandle) TransferList(c *gin.Context) {
	var request customerReq.OnInheritSearchReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	cc := h.WarpContext(c)
	memberIds, err := customerSvc.NewExternalSearch().GetMembersByManageType(cc, request.EnterpriseID, request.MemberId, request.ThirdIds, request.MemberIds)
	request.MemberIds = memberIds
	num, err := h.onInheritService.TransferList(cc, request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, num)
}
