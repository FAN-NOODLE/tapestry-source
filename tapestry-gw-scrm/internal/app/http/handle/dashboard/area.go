package dashboard

import (
	"github.com/gin-gonic/gin"
	"gw-scrm/internal/app/http/handle"
	"gw-scrm/internal/app/http/request"
	dashboardReq "gw-scrm/internal/app/http/request/dashboard"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type AreaHandle struct {
	handle.BaseHandle
	api microscrm.DashboardApi
}

func NewAreaHandle() AreaHandle {
	return AreaHandle{
		api: microscrm.NewDashboardApi(),
	}
}

// NotRelationStoreList  待绑定门店列表
func (h AreaHandle) NotRelationStoreList(c *gin.Context) {
	var params request.EnterpriseBaseReq
	enterpriseID, _ := common.GetEnterprise(c)
	params.EnterpriseId = enterpriseID
	res, err := h.api.NotRelationStoreList(h.WarpContext(c), params)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, res)

	return
}

//// StoreList  已绑定门店列表
//func (h AreaHandle) StoreList(c *gin.Context) {
//	var params dashboardReq.StoreList
//	if err := c.ShouldBind(&params); err != nil {
//		h.ResError(c, err)
//		return
//	}
//	enterpriseID, _ := common.GetEnterprise(c)
//	params.EnterpriseId = enterpriseID
//	res, err := h.api.StoreList(h.WarpContext(c), params)
//	if err != nil {
//		h.ResError(c, err)
//		return
//	}
//
//	h.ResSuccess(c, res)
//
//	return
//}

// AreaTree dashBoard组织架构
func (h AreaHandle) AreaTree(c *gin.Context) {
	var params dashboardReq.AreaTree
	if err := c.ShouldBind(&params); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseID, _ := common.GetEnterprise(c)
	params.EnterpriseId = enterpriseID
	res, err := h.api.AreaTree(h.WarpContext(c), params)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, res)

	return
}

// Save  插入/更新节点
func (h AreaHandle) Save(c *gin.Context) {
	var params dashboardReq.AreaSave
	if err := c.ShouldBind(&params); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseID, _ := common.GetEnterprise(c)
	params.EnterpriseId = enterpriseID
	res, err := h.api.AreaSave(h.WarpContext(c), params)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, res)

	return
}

// Remove 删除节点
func (h AreaHandle) Remove(c *gin.Context) {
	var params dashboardReq.AreaRemove
	if err := c.ShouldBind(&params); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseID, _ := common.GetEnterprise(c)
	params.EnterpriseId = enterpriseID
	res, err := h.api.AreaRemove(h.WarpContext(c), params)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, res)

	return
}

func (h AreaHandle) AreaList(c *gin.Context) {
	var params request.EnterpriseBaseReq
	if err := c.ShouldBind(&params); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseID, _ := common.GetEnterprise(c)
	params.EnterpriseId = enterpriseID
	res, err := h.api.AreaList(h.WarpContext(c), params)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, res)

	return
}
