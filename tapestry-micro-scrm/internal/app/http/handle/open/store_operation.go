package open

import (
	"github.com/gin-gonic/gin"

	"micro-scrm/internal/app/http/handle"
	"micro-scrm/internal/app/http/middleware"
	openReq "micro-scrm/internal/app/http/request/open"
	openSvc "micro-scrm/internal/app/service/open"
)

type StoreOperationHandle struct {
	handle.BaseHandle
	service *openSvc.StoreOperationService
}

func NewStoreOperationHandle() *StoreOperationHandle {
	return &StoreOperationHandle{service: openSvc.NewStoreOperationService()}
}

// StoreMasterList 获取门店主数据列表。
func (h *StoreOperationHandle) StoreMasterList(c *gin.Context) {
	var request openReq.StoreMasterListReq
	if err := c.ShouldBindJSON(&request); err != nil {
		h.ResError(c, err)
		return
	}
	request.EnterpriseID = c.GetInt64(middleware.EnterpriseId)
	list, total, err := h.service.StoreMasterList(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResList(c, total, normalizedPage(request.Page), normalizedStorePageSize(request.Size), list)
}

// EmployeeMasterList 获取员工主数据列表。
func (h *StoreOperationHandle) EmployeeMasterList(c *gin.Context) {
	var request openReq.EmployeeMasterListReq
	if err := c.ShouldBindJSON(&request); err != nil {
		h.ResError(c, err)
		return
	}
	request.EnterpriseID = c.GetInt64(middleware.EnterpriseId)
	list, total, err := h.service.EmployeeMasterList(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResList(c, total, normalizedPage(request.Page), normalizedStorePageSize(request.Size), list)
}

func normalizedPage(page int) int {
	if page == 0 {
		return 1
	}
	return page
}

func normalizedStorePageSize(size int) int {
	if size == 0 {
		return 100
	}
	return size
}
