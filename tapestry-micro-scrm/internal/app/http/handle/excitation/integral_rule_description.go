package excitation

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	excitationSvc "micro-scrm/internal/app/service/excitation"
)

type IntegralRuleDescriptionHandle struct {
	handle.BaseHandle
	integralRuleDescriptionSvc *excitationSvc.IntegralRuleDescriptionService
}

func NewIntegralRuleDescriptionHandle() IntegralRuleDescriptionHandle {
	return IntegralRuleDescriptionHandle{
		integralRuleDescriptionSvc: excitationSvc.NewIntegralRuleDescriptionService(),
	}
}

func (h IntegralRuleDescriptionHandle) CreateOUpdate(c *gin.Context) {
	var request excitationReq.DescriptionCreateOrUpdateReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	err := h.integralRuleDescriptionSvc.CreateOUpdate(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, "")
	return
}

// Detail 规则说明详情
func (h IntegralRuleDescriptionHandle) Detail(c *gin.Context) {
	var request excitationReq.DescriptionReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	resp, err := h.integralRuleDescriptionSvc.Detail(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, resp)
	return

}
