package online_warehouse

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gw-scrm/internal/app/http/handle"
	onlineWarehouseReq "gw-scrm/internal/app/http/request/online_warehouse"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type OrderHandle struct {
	handle.BaseHandle
	onlineWarehouse microscrm.OnlineWarehouseApi
}

func NewOrderHandle() OrderHandle {
	return OrderHandle{
		onlineWarehouse: microscrm.NewOnlineWarehouseApi(),
	}
}

func (h OrderHandle) List(c *gin.Context) {
	var req onlineWarehouseReq.OrderReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, adminId := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	req.AdminID = adminId
	list, err := h.onlineWarehouse.List(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}

func (h OrderHandle) Submit(c *gin.Context) {
	var req onlineWarehouseReq.SubmitOrderReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, adminId := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	req.AdminID = adminId
	res, err := h.onlineWarehouse.Submit(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h OrderHandle) Info(c *gin.Context) {
	var req onlineWarehouseReq.InfoReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, adminId := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	req.AdminID = adminId
	res, err := h.onlineWarehouse.Info(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h OrderHandle) StoreList(c *gin.Context) {
	var req onlineWarehouseReq.StoreListReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, adminId := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	req.AdminID = adminId
	res, err := h.onlineWarehouse.StoreList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h OrderHandle) Edit(c *gin.Context) {
	var req onlineWarehouseReq.SubmitOrderReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, adminId := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	req.AdminID = adminId
	res, err := h.onlineWarehouse.Edit(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h OrderHandle) StoreOrderList(c *gin.Context) {
	var req onlineWarehouseReq.StoreOrderListReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, adminId := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	req.AdminID = adminId
	res, err := h.onlineWarehouse.StoreOrderList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h OrderHandle) Export(c *gin.Context) {
	var req onlineWarehouseReq.OrderReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, adminId := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	req.AdminID = adminId
	result, err := sdk.SendPostNoLog(h.WarpContext(c), microscrm.GetUrl("%s/micro-scrm/v1/online_warehouse/export"), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, result)
}

func (h OrderHandle) ImportOrder(c *gin.Context) {
	var req onlineWarehouseReq.ImportOrderReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	// 获取enterpriseID
	enterpriseId, adminId := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	req.AdminID = adminId
	// 将上传的文件url地址给到micro_scrm处理解析文件
	res, err := h.onlineWarehouse.ImportOrder(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h OrderHandle) BatchTransfer(c *gin.Context) {
	var req onlineWarehouseReq.BatchTransferReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, adminId := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	req.AdminID = adminId
	res, err := h.onlineWarehouse.BatchTransfer(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h OrderHandle) Refund(c *gin.Context) {
	var req onlineWarehouseReq.InfoReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, adminId := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	req.AdminID = adminId
	res, err := h.onlineWarehouse.Refund(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h OrderHandle) DepList(c *gin.Context) {
	var req onlineWarehouseReq.BaseReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, adminId := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	req.AdminID = adminId
	res, err := h.onlineWarehouse.DepList(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h OrderHandle) TransferExport(c *gin.Context) {
	var req onlineWarehouseReq.TransferExportReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, adminId := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	req.AdminID = adminId
	if req.ThirdId == 0 {
		h.ResError(c, errors.New("请选择门店"))
		return
	}
	res, err := h.onlineWarehouse.TransferExport(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h OrderHandle) TransferImport(c *gin.Context) {
	var req onlineWarehouseReq.TransferImportReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseId, adminId := common.GetEnterprise(c)
	req.EnterpriseId = enterpriseId
	req.AdminID = adminId
	if req.ThirdId == 0 {
		h.ResError(c, errors.New("请选择门店"))
		return
	}
	res, err := h.onlineWarehouse.TransferImport(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}
