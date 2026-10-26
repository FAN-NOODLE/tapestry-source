package open

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	"micro-scrm/internal/app/http/middleware"
	openReq "micro-scrm/internal/app/http/request/open"
	openSvc "micro-scrm/internal/app/service/open"
)

type HkStoreWorkerHandle struct {
	handle.BaseHandle

	openSvc *openSvc.HkStoreWorkerService
}

func NewHkStoreWorkerHandle() HkStoreWorkerHandle {
	return HkStoreWorkerHandle{
		openSvc: openSvc.NewHkStoreWorkerService(),
	}
}

// StoreList 获取香港门店列表
func (h HkStoreWorkerHandle) StoreList(c *gin.Context) {
	var request openReq.HkStoreListReq
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
}

// GetStoreWorkList 获取香港门店导购信息列表
func (h HkStoreWorkerHandle) GetStoreWorkList(c *gin.Context) {
	var request openReq.HkStoreWorkListReq
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
}

// AreaManager 香港门店管理者
func (h HkStoreWorkerHandle) AreaManager(c *gin.Context) {
	var req openReq.HkAreaManagerReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	req.EnterpriseId = c.GetInt64(middleware.EnterpriseId)
	list, err := h.openSvc.AreaManager(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, list)
}
