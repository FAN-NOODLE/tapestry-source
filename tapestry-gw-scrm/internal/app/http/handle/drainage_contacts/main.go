package drainage_contacts

import (
	"github.com/gin-gonic/gin"
	"gw-scrm/internal/app/http/handle"
	"gw-scrm/internal/app/http/request/drainage_contacts"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type DrainageContactsCoachHandle struct {
	handle.BaseHandle
	DrainageContactCoachApi microscrm.DrainageContactCoachApi
}

func NewDrainageContactsCoachHandle() DrainageContactsCoachHandle {
	return DrainageContactsCoachHandle{
		DrainageContactCoachApi: microscrm.NewDrainageContactCoachApi(),
	}
}

func (h DrainageContactsCoachHandle) Save(c *gin.Context) {
	var req drainageContactsReq.GroupSaveReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	res, err := h.DrainageContactCoachApi.Save(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return
}

func (h DrainageContactsCoachHandle) ZipUrl(c *gin.Context) {
	var req drainageContactsReq.ZipUrlReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	res, err := h.DrainageContactCoachApi.ZipUrl(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return
}

func (h DrainageContactsCoachHandle) Detail(c *gin.Context) {
	var req drainageContactsReq.InfoReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	res, err := h.DrainageContactCoachApi.Detail(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return
}

func (h DrainageContactsCoachHandle) Delete(c *gin.Context) {
	var req drainageContactsReq.InfoReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	res, err := h.DrainageContactCoachApi.Delete(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return
}

func (h DrainageContactsCoachHandle) List(c *gin.Context) {
	var req drainageContactsReq.GroupListReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	res, err := h.DrainageContactCoachApi.List(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return
}

func (h DrainageContactsCoachHandle) DetailList(c *gin.Context) {
	var req drainageContactsReq.DetailListReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	res, err := h.DrainageContactCoachApi.DetailList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
	return
}
