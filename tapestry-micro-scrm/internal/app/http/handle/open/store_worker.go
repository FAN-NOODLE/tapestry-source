package open

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	"micro-scrm/internal/app/http/middleware"
	openReq "micro-scrm/internal/app/http/request/open"
	dashboardSvc "micro-scrm/internal/app/service/dashboard"
	openSvc "micro-scrm/internal/app/service/open"
)

type StoreWorkerHandle struct {
	handle.BaseHandle

	openSvc *openSvc.StoreWorkerService
}

func NewStoreWorkerHandle() StoreWorkerHandle {
	return StoreWorkerHandle{
		openSvc: openSvc.NewStoreWorkerService(),
	}
}

// SaListByCustomer 顾客的导购
func (h StoreWorkerHandle) SaListByCustomer(c *gin.Context) {
	var request openReq.StoreWorkerSaListReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	request.EnterpriseId = c.GetInt64(middleware.EnterpriseId)

	resp, err := h.openSvc.SaListByCustomer(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, resp)
	return
}

// StoreList 获取门店列表
func (h StoreWorkerHandle) StoreList(c *gin.Context) {
	var request openReq.StoreListReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	request.EnterpriseId = c.GetInt64(middleware.EnterpriseId)
	if request.Page == 0 {
		request.Page = 1
	}

	if request.Size == 0 {
		request.Size = 10
	}

	resp, total, err := h.openSvc.GetStoreList(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResList(c, total, request.Page, request.Size, resp)
	return
}

// SaList 获取导购列表
func (h StoreWorkerHandle) SaList(c *gin.Context) {
	var request openReq.SaListReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	request.EnterpriseId = c.GetInt64(middleware.EnterpriseId)

	resp, total, err := h.openSvc.GetSaList(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResList(c, total, request.Page, request.Size, resp)
	return
}

// GetStoreWorkList 获取导购信息列表
func (h StoreWorkerHandle) GetStoreWorkList(c *gin.Context) {
	var request openReq.StoreWorkListReq
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

	request.EnterpriseId = c.GetInt64(middleware.EnterpriseId)
	resp, total, err := h.openSvc.GetStoreWorkList(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResList(c, total, request.Page, request.Size, resp)
	return
}

// GetCustomerStoreWorkList 获取客户添加的导购列表
func (h StoreWorkerHandle) GetCustomerStoreWorkList(c *gin.Context) {
	var request openReq.CustomerStoreWorkListReq
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
	request.EnterpriseId = c.GetInt64(middleware.EnterpriseId)
	resp, total, err := h.openSvc.GetCustomerStoreWorkList(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResList(c, total, request.Page, request.Size, resp)
	return
}

// GenQrCode 生成导购联系我二维码
func (h StoreWorkerHandle) GenQrCode(c *gin.Context) {
	var req openReq.GenQrCodeReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}

	req.EnterpriseId = c.GetInt64(middleware.EnterpriseId)
	resp, err := h.openSvc.GenQrCode(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, resp)
	return
}

// AreaTree 组织架构树
func (h StoreWorkerHandle) AreaTree(c *gin.Context) {
	var req openReq.AreaTreeReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.EnterpriseId = c.GetInt64(middleware.EnterpriseId)
	list, err := dashboardSvc.NewAreaService().AreaListForOpen(h.WarpContext(c), req)

	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}

// AreaManager 组织架构管理者
func (h StoreWorkerHandle) AreaManager(c *gin.Context) {
	var req openReq.AreaManagerReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.EnterpriseId = c.GetInt64(middleware.EnterpriseId)
	list, err := dashboardSvc.NewAreaService().AreaManager(h.WarpContext(c), req)

	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}
