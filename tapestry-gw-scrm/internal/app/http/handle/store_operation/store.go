package store_operation

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gw-scrm/internal/app/http/handle"
	storeReq "gw-scrm/internal/app/http/request/store_operation"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type StoreHandle struct {
	handle.BaseHandle
}

type OrganizationHandle struct {
	handle.BaseHandle
}

func NewStoreHandle() StoreHandle {
	return StoreHandle{}
}

func NewOrganizationHandle() OrganizationHandle {
	return OrganizationHandle{}
}

func (h StoreHandle) List(c *gin.Context) {
	var req storeReq.StoreListReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseID, _ := common.GetEnterprise(c)
	req.EnterpriseID = enterpriseID
	result, err := microscrm.StoreOperationList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, result)
}

func (h StoreHandle) Detail(c *gin.Context) {
	var req storeReq.StoreDetailReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseID, _ := common.GetEnterprise(c)
	req.EnterpriseID = enterpriseID
	result, err := microscrm.StoreOperationDetail(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, result)
}

func (h StoreHandle) Update(c *gin.Context) {
	var req storeReq.StoreUpdateReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	tokenInfo, ok := c.Value(common.ContextKeyEnum.TokenInfo).(*common.TokenInfo)
	if !ok || tokenInfo == nil {
		h.ResError(c, fmt.Errorf("身份信息异常"))
		return
	}
	req.EnterpriseID = tokenInfo.EnterpriseId
	req.OperatorAdminID = tokenInfo.AdminId
	req.OperatorName = tokenInfo.Username
	result, err := microscrm.StoreOperationUpdate(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, result)
}

func (h OrganizationHandle) Tree(c *gin.Context) {
	var req storeReq.OrganizationTreeReq
	enterpriseID, _ := common.GetEnterprise(c)
	req.EnterpriseID = enterpriseID
	result, err := microscrm.StoreOperationOrganizationTree(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, result)
}

func (h OrganizationHandle) StoreList(c *gin.Context) {
	var req storeReq.OrganizationStoreListReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseID, _ := common.GetEnterprise(c)
	req.EnterpriseID = enterpriseID
	result, err := microscrm.StoreOperationOrganizationStoreList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, result)
}
