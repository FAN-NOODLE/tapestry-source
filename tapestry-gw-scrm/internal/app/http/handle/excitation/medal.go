package excitation

import (
	"github.com/gin-gonic/gin"
	"gw-scrm/internal/app/http/handle"
	excitationReq "gw-scrm/internal/app/http/request/excitation"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type MedalHandle struct {
	handle.BaseHandle
	api microscrm.MedalApi
}

func NewMedalHandle() MedalHandle {
	return MedalHandle{
		api: microscrm.NewMedalApi(),
	}
}

func (h MedalHandle) List(c *gin.Context) {
	var request excitationReq.MedalList
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

func (h MedalHandle) Create(c *gin.Context) {
	var request excitationReq.MedalCreate
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseID
	list, err := h.api.Create(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, list)
	return
}

func (h MedalHandle) Detail(c *gin.Context) {
	var request excitationReq.MedalDetail
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseID
	list, err := h.api.Detail(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, list)
	return
}

func (h MedalHandle) Update(c *gin.Context) {
	var request excitationReq.MedalUpdate
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseID
	list, err := h.api.Update(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, list)
	return
}

func (h MedalHandle) Remove(c *gin.Context) {
	var request excitationReq.MedalDetail
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseID
	list, err := h.api.Remove(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, list)
	return
}

func (h MedalHandle) Status(c *gin.Context) {
	var request excitationReq.MedalStatus
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseID
	list, err := h.api.Status(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, list)
	return
}
