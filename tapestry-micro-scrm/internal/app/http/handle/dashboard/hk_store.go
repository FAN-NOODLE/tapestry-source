package dashboard

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	dashboardReq "micro-scrm/internal/app/http/request/dashboard"
	dashboardSvc "micro-scrm/internal/app/service/dashboard"
)

type HkStoreHandle struct {
	handle.BaseHandle
	storeSvc *dashboardSvc.HkStoreService
}

func NewHkStoreHandle() HkStoreHandle {
	return HkStoreHandle{
		storeSvc: dashboardSvc.NewHkStoreService(),
	}
}

// Sync 同步企微香港门店
func (h HkStoreHandle) Sync(c *gin.Context) {
	var req dashboardReq.HkStoreSync
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	result, err := h.storeSvc.Sync(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, result)
}

// LastSyncTime 获取香港门店最后同步时间
func (h HkStoreHandle) LastSyncTime(c *gin.Context) {
	var req dashboardReq.HkStoreLastSyncTime
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	result, err := h.storeSvc.LastSyncTime(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, result)
}

// List 获取香港门店列表
func (h HkStoreHandle) List(c *gin.Context) {
	var req dashboardReq.HkStoreList
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Size == 0 {
		req.Size = 10
	}
	list, total, err := h.storeSvc.List(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResList(c, total, req.Page, req.Size, list)
}

// Detail 获取香港门店详情
func (h HkStoreHandle) Detail(c *gin.Context) {
	var req dashboardReq.HkStoreDetail
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	detail, err := h.storeSvc.Detail(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, detail)
}

// Save 更新香港门店
func (h HkStoreHandle) Save(c *gin.Context) {
	var req dashboardReq.HkStoreSave
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	if err := h.storeSvc.Save(h.WarpContext(c), req); err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, "")
}

// Close 关闭香港门店
func (h HkStoreHandle) Close(c *gin.Context) {
	var req dashboardReq.HkStoreClose
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}
	if err := h.storeSvc.Close(h.WarpContext(c), req); err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, "")
}
