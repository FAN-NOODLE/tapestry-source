package excitation

import (
	"github.com/gin-gonic/gin"
	"gw-scrm/internal/app/http/handle"
	excitationReq "gw-scrm/internal/app/http/request/excitation"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type MedalRuleHandle struct {
	handle.BaseHandle
	api microscrm.MedalRuleApi
}

func NewMedalRuleHandle() MedalRuleHandle {
	return MedalRuleHandle{
		api: microscrm.NewMedalRuleApi(),
	}
}

func (h MedalRuleHandle) List(c *gin.Context) {
	var request excitationReq.MedalRuleList
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
	data, err := h.api.List(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, data)

	return
}
