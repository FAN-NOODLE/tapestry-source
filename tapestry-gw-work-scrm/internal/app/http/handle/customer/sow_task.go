package customer

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-work-scrm/internal/app/http/handle"
	customerReq "gw-work-scrm/internal/app/http/request/customer"
	microScrmCustomer "gw-work-scrm/internal/pkg/sdk/microscrm/customer"
)

type SowTaskHandle struct {
	handle.BaseHandle
	SowTaskSdk *microScrmCustomer.SowTaskApi
}

func NewSowTaskHandle() SowTaskHandle {
	return SowTaskHandle{
		BaseHandle: handle.BaseHandle{},
		SowTaskSdk: microScrmCustomer.NewSowTaskApi(),
	}
}

// List 顾客的播种列表
func (h SowTaskHandle) List(c *gin.Context) {
	var request customerReq.ListReq

	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	request.EnterpriseID = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	request.MemberMobile = memberInfo.Mobile

	list, err := h.SowTaskSdk.Lists(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
	return

}

// Create 创建播种
func (h SowTaskHandle) Create(c *gin.Context) {
	var request customerReq.SowTaskCreate

	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	request.EnterpriseID = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	request.MemberMobile = memberInfo.Mobile
	res, err := h.SowTaskSdk.Create(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return

}

// Update 修改自设任务播种
func (h SowTaskHandle) Update(c *gin.Context) {
	var request customerReq.SowTaskUpdate

	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	request.EnterpriseID = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	request.MemberMobile = memberInfo.Mobile
	res, err := h.SowTaskSdk.Update(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return

}

// Finished 完成自设任务播种
func (h SowTaskHandle) Finished(c *gin.Context) {
	var request customerReq.SowTaskFinished

	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	request.EnterpriseID = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	request.MemberMobile = memberInfo.Mobile
	res, err := h.SowTaskSdk.Finished(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return

}

// Delete 删除自设任务播种
func (h SowTaskHandle) Delete(c *gin.Context) {
	var request customerReq.SowTaskFinished

	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	request.EnterpriseID = h.EnterpriseID(c)
	memberInfo := h.MemberInfo(c)
	request.MemberMobile = memberInfo.Mobile
	res, err := h.SowTaskSdk.Delete(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return

}
