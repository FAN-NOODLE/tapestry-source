package customer

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/app/http/handle"
	customerReq "gw-scrm/internal/app/http/request/customer"
	"gw-scrm/internal/pkg/common"
	scrm_errs "gw-scrm/internal/pkg/errs"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type InheritHandle struct {
	handle.BaseHandle
}

func NewInheritHandle() InheritHandle {
	return InheritHandle{}
}

// List 离职成员列表
func (h InheritHandle) List(c *gin.Context) {
	var request customerReq.CustomerListReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseId

	result, err := microscrm.GetList(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, result)
	return
}

func (h InheritHandle) TransferCustomer(c *gin.Context) {
	var request customerReq.ResignedTransferCustomerReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseId

	result, err := microscrm.TransferCustomer(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, result)
	return
}

func (h InheritHandle) TransferCustomerGroup(c *gin.Context) {
	var request customerReq.ResignedTransferCustomerGroupReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseId
	result, err := microscrm.TransferCustomerGroup(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, result)

	return
}

func (h InheritHandle) GetExternalAllotList(c *gin.Context) {
	var request customerReq.RequestCustomerExternalList
	err := c.ShouldBind(&request)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseId

	res, err := microscrm.GetExternalAllotList(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
}
func (h InheritHandle) GetGroupChatAllotList(c *gin.Context) {
	var request customerReq.RequestCustomerGroupList
	err := c.ShouldBind(&request)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseId

	res, err := microscrm.GetGroupChatAllotList(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
}

func (h InheritHandle) GetTakeOverList(c *gin.Context) {
	var request customerReq.RequestCustomerGroupList
	err := c.ShouldBind(&request)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseId

	res, err := microscrm.GetTakeOverList(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
}

func (h InheritHandle) GetTakeOverExternalList(c *gin.Context) {
	var request customerReq.RequestTakeOverExternalList
	err := c.ShouldBind(&request)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseId

	res, err := microscrm.GetTakeOverExternalList(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
}

func (h InheritHandle) GetTakeOverGroupList(c *gin.Context) {
	var request customerReq.RequestCustomerGroupList
	err := c.ShouldBind(&request)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseId

	res, err := microscrm.GetTakeOverGroupList(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
}

func (h InheritHandle) SyncCustomerData(c *gin.Context) {
	var request customerReq.SyncCustomDataRes
	err := c.ShouldBind(&request)
	if err != nil {
		ginplus.ResError(c, scrm_errs.ParamsErr)
		return
	}
	enterpriseId, _ := common.GetEnterprise(c)
	request.EnterpriseID = enterpriseId

	res, err := microscrm.SyncCustomerData(ctx.Wrap(c), request)

	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, res)
}
