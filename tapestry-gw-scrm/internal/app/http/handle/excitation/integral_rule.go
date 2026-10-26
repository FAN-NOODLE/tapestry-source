package excitation

import (
	"github.com/gin-gonic/gin"
	"gw-scrm/internal/app/http/handle"
	excitationReq "gw-scrm/internal/app/http/request/excitation"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type IntegralRuleHandle struct {
	handle.BaseHandle
	api microscrm.IntegralRuleApi
}

func NewIntegralRuleHandle() IntegralRuleHandle {
	return IntegralRuleHandle{
		api: microscrm.NewIntegralRuleApi(),
	}
}

//List 列表
func (h *IntegralRuleHandle) List(c *gin.Context) {
	var request excitationReq.RuleList
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	if request.Page == 0 {
		request.Page = 1
	}

	if request.Size == 0 {
		request.Size = 10
	}
	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseID
	list, err := h.api.List(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, list)
	return
}

//ParameterConfiguration 配置信息
func (h *IntegralRuleHandle) ParameterConfiguration(c *gin.Context) {
	var request excitationReq.ParameterConfiguration
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseID
	resp, err := h.api.ParameterConfiguration(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, resp)
	return
}

//Create 创建
func (h *IntegralRuleHandle) Create(c *gin.Context) {
	var request excitationReq.RuleCreate
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseID
	resp, err := h.api.Create(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, resp)
	return
}

//Detail 详情
func (h *IntegralRuleHandle) Detail(c *gin.Context) {
	var request excitationReq.RuleDetail
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseID
	resp, err := h.api.Detail(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, resp)
	return
}

//Destroy 删除
func (h *IntegralRuleHandle) Destroy(c *gin.Context) {
	var request excitationReq.RuleDestroy
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseID
	resp, err := h.api.Destroy(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, resp)
	return
}

//Update 更新
func (h *IntegralRuleHandle) Update(c *gin.Context) {
	var request excitationReq.RuleUpdate
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseID
	resp, err := h.api.Update(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, resp)
	return
}

func (h *IntegralRuleHandle) DescriptionSave(c *gin.Context) {
	var request excitationReq.DescriptionSave
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseID
	resp, err := h.api.DescriptionSave(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, resp)
	return
}

func (h *IntegralRuleHandle) DescriptionDetail(c *gin.Context) {
	var request excitationReq.DescriptionDetail

	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseID
	resp, err := h.api.DescriptionDetail(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, resp)
	return
}
