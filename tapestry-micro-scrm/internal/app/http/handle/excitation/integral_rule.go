package excitation

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	excitationSvc "micro-scrm/internal/app/service/excitation"
)

// 积分规则

type IntegralRuleHandle struct {
	handle.BaseHandle
	integralRuleSvc *excitationSvc.IntegralRuleService
}

func NewIntegralRuleHandle() IntegralRuleHandle {
	return IntegralRuleHandle{
		integralRuleSvc: excitationSvc.NewIntegralRuleService(),
	}

}

//List 别表
func (h IntegralRuleHandle) List(c *gin.Context) {
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

	resp, total, err := h.integralRuleSvc.GetList(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResList(c, total, request.Page, request.Size, resp)

	return
}

//ParameterConfiguration 规则配置信息
func (h IntegralRuleHandle) ParameterConfiguration(c *gin.Context) {
	var request excitationReq.ParameterConfiguration
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	resp := h.integralRuleSvc.ParameterConfiguration(request)
	h.ResSuccess(c, resp)
}

func (h IntegralRuleHandle) Create(c *gin.Context) {
	var request excitationReq.RuleCreate
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	err := h.integralRuleSvc.Create(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, "")
	return
}

// Detail 规则积分详情
func (h IntegralRuleHandle) Detail(c *gin.Context) {
	var request excitationReq.RuleDetail
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	resp, err := h.integralRuleSvc.Detail(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, resp)
	return

}

func (h IntegralRuleHandle) Update(c *gin.Context) {
	var request excitationReq.RuleUpdate
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	err := h.integralRuleSvc.Update(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, "")
	return
}

// Destroy 删除
func (h IntegralRuleHandle) Destroy(c *gin.Context) {
	var request excitationReq.RuleDestroy
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	err := h.integralRuleSvc.Destroy(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, "")
	return
}
