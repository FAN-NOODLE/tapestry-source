package customer

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	"micro-scrm/internal/app/http/request/customer_search"
	customerSvc "micro-scrm/internal/app/service/customer"
)

type ExternalSearchHandle struct {
	handle.BaseHandle
}

func NewExternalSearchHandle() ExternalSearchHandle {
	return ExternalSearchHandle{}
}

// SearchOption 筛选外部联系人列表的条件
func (h ExternalSearchHandle) SearchOption(c *gin.Context) {

	var request customerSearchReq.ReqChoice
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	// SM管理的门店导购及自己
	choiceOptionList, err := customerSvc.NewExternalSearch().GetSearchOption(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, choiceOptionList)

	return
}

// List 根据条件筛选导购和外部联系人列表
func (h ExternalSearchHandle) List(c *gin.Context) {
	var request customerSearchReq.ExternalSearchListBaseReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	request.RelationshipType = customerSearchReq.IsUnionTure
	result, total, err := customerSvc.NewExternalSearch().GetExternalSearchListCk(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResList(c, total, request.Page, request.Size, result)

	return
}

// Detail  顾客详情
func (h ExternalSearchHandle) Detail(c *gin.Context) {
	var request customerSearchReq.ExternalDetailReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	detail, err := customerSvc.NewExternalSearch().GetExternalDetail(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, detail)
	return
}

// GetOpenIds   全选获取所有openid集合
func (h ExternalSearchHandle) GetOpenIds(c *gin.Context) {
	var request customerSearchReq.ExternalSearchListBaseReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	request.RelationshipType = customerSearchReq.IsUnionTure
	openIds, err := customerSvc.NewExternalSearch().GetExternalSearchCk(h.WarpContext(c), request, "openid")
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, openIds)

	return
}
