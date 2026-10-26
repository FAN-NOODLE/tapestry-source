package customer

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	"micro-scrm/internal/app/http/request/customer_search"
	customerSvc "micro-scrm/internal/app/service/customer"
)

type ExternalSearchSwHandle struct {
	handle.BaseHandle
}

func NewExternalSearchSwHandle() ExternalSearchSwHandle {
	return ExternalSearchSwHandle{}
}

// SearchOption 筛选外部联系人列表的条件
func (h ExternalSearchSwHandle) SearchOption(c *gin.Context) {

	var request customerSearchReq.ReqChoice
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	// SM管理的门店导购及自己
	choiceOptionList, err := customerSvc.NewExternalSearchSw().GetSearchOption(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, choiceOptionList)

	return
}

// List 根据条件筛选导购和外部联系人列表
func (h ExternalSearchSwHandle) List(c *gin.Context) {
	var request customerSearchReq.ExternalSearchListBaseSwReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	request.RelationshipType = customerSearchReq.IsUnionTure
	result, total, err := customerSvc.NewExternalSearchSw().GetExternalSearchListCk(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResList(c, total, request.Page, request.Size, result)

	return
}

// GetOpenIds   全选获取所有openid集合
func (h ExternalSearchSwHandle) GetOpenIds(c *gin.Context) {
	var request customerSearchReq.ExternalSearchListBaseSwReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	request.RelationshipType = customerSearchReq.IsUnionTure
	openIds, err := customerSvc.NewExternalSearchSw().GetExternalSearchCk(h.WarpContext(c), request, "openid")
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, openIds)

	return
}
