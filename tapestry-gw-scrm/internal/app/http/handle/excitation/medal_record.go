package excitation

import (
	"github.com/gin-gonic/gin"
	"gw-scrm/internal/app/http/handle"
	excitationReq "gw-scrm/internal/app/http/request/excitation"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type MedalRecordHandle struct {
	handle.BaseHandle
	api microscrm.MedalRecordApi
}

func NewMedalRecordHandle() MedalRecordHandle {
	return MedalRecordHandle{
		api: microscrm.NewMedalRecordApi(),
	}
}

func (h MedalRecordHandle) List(c *gin.Context) {
	var request excitationReq.MedalRecordList
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
	resp, err := h.api.List(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, resp)
	return
}

// Reissue 补发
func (h MedalRecordHandle) Reissue(c *gin.Context) {
	var request excitationReq.MedalRecordReissue
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseID
	resp, err := h.api.Reissue(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, resp)
	return
}

func (h MedalRecordHandle) Remove(c *gin.Context) {
	var request struct {
		Id           int64 `json:"id" binding:"required"`
		EnterpriseId int64 `json:"enterprise_id"`
	}
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseID
	resp, err := h.api.Remove(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, resp)
	return
}
