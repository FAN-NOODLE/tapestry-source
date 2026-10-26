package excitation

import (
	"github.com/gin-gonic/gin"
	excitationDao "micro-scrm/internal/app/dao/excitation"
	"micro-scrm/internal/app/http/handle"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	excitationResp "micro-scrm/internal/app/http/response/excitation"
)

type MedalRuleHandle struct {
	handle.BaseHandle
	medalRuleDao *excitationDao.MedalRule
}

func NewMedalRuleHandle() MedalRuleHandle {
	return MedalRuleHandle{
		medalRuleDao: excitationDao.NewMedalRule(),
	}
}

func (h MedalRuleHandle) List(c *gin.Context) {
	var request excitationReq.MedalRuleList
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	data, total, err := h.medalRuleDao.GetPagination(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	var resp []excitationResp.MedalRuleList
	for _, val := range data {
		resp = append(resp, excitationResp.MedalRuleList{
			Id:          val.ID,
			Pid:         val.Pid,
			RuleName:    val.RuleName,
			TimeType:    val.TimeType,
			Description: val.Description,
		})
	}

	h.ResList(c, total, request.Page, request.Size, resp)

	return
}
