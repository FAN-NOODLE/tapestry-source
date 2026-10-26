package plan

import (
	"github.com/gin-gonic/gin"
	"gw-scrm/internal/app/http/handle"
	planReq "gw-scrm/internal/app/http/request/plan"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type TemplateHandle struct {
	handle.BaseHandle
	api microscrm.PlanTemplateApi
}

func NewTemplateHandle() TemplateHandle {
	return TemplateHandle{
		BaseHandle: handle.BaseHandle{},
		api:        microscrm.NewPlanTemplateApi(),
	}
}

func (h TemplateHandle) List(c *gin.Context) {
	var request planReq.TemplateListsReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseID

	res, err := h.api.List(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, res)
}

// Create 任务模板的创建
func (h TemplateHandle) Create(c *gin.Context) {
	var request planReq.TemplateAddReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	enterpriseID, adminId := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseID
	request.AdminId = adminId

	res, err := h.api.Create(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, res)
}

// Update 任务模板的创建
func (h TemplateHandle) Update(c *gin.Context) {
	var request planReq.TemplateUpdateReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	enterpriseID, adminId := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseID
	request.AdminId = adminId

	res, err := h.api.Update(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h TemplateHandle) Status(c *gin.Context) {
	var request struct {
		Id           int64 `json:"id"`
		EnterpriseID int64 `json:"enterprise_id"`
		Status       int   `json:"status"`
	}
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseID

	res, err := h.api.Status(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h TemplateHandle) Details(c *gin.Context) {
	var request planReq.TemplateDetails
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseID

	res, err := h.api.Detail(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}
