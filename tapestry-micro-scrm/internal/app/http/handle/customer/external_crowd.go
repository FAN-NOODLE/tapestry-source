package customer

import (
	"github.com/gin-gonic/gin"
	externalSearchDao "micro-scrm/internal/app/dao/external_search"
	"micro-scrm/internal/app/http/handle"
	"micro-scrm/internal/app/http/request"
	customerSearchReq "micro-scrm/internal/app/http/request/customer_search"
	customerSearchResp "micro-scrm/internal/app/http/response/customer_search"
	customerSvc "micro-scrm/internal/app/service/customer"
	"micro-scrm/internal/pkg/helper"
)

type ExternalCrowdHandle struct {
	handle.BaseHandle
	svc *customerSvc.ExternalCrowd
}

func NewExternalCrowdHandle() ExternalCrowdHandle {
	return ExternalCrowdHandle{
		svc: customerSvc.NewExternalCrowd(),
	}
}

func (h ExternalCrowdHandle) Add(c *gin.Context) {
	var req customerSearchReq.AddCrowdReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	cc := h.WarpContext(c)
	var searchValue customerSearchReq.ExternalSearchListBaseReq
	helper.Unmarshal(req.Search.Value, &searchValue)
	searchValue.EnterpriseID = req.EnterpriseId
	searchValue.MemberId = req.MemberId
	searchValue.ThirdIds = req.ThirdIds
	searchValue.RelationshipType = customerSearchReq.IsUnionTure
	// 获取所有的external_userId
	ids, err := customerSvc.NewExternalSearch().GetExternalSearchCk(cc, searchValue, "id")
	if err != nil {
		h.ResError(c, err)
		return
	}
	var externalIds []int64
	for _, v := range ids {
		// 使用类型断言将 interface{} 转换为 int64
		if val, ok := v.(int64); ok {
			externalIds = append(externalIds, val)
		}
	}
	err = h.svc.Add(cc, req, externalIds)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)

}

func (h ExternalCrowdHandle) Info(c *gin.Context) {
	var req customerSearchReq.CrowdInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	cc := h.WarpContext(c)
	req.SearchList.MemberId = req.MemberId
	searchQueryResp := externalSearchDao.NewExternalExtraCk().SearchExternalOptionQuery(req.SearchList)
	crowdList, total, err := h.svc.PageListByExternalConn(cc, req, searchQueryResp)
	if err != nil {
		h.ResError(c, err)
		return
	}
	resList := make([]customerSearchResp.ExternalCrowdList, 0)
	if len(crowdList) == 0 {
		h.ResList(c, 0, req.Page, req.Size, resList)
		return
	}
	var externalIds []int64
	for _, row := range crowdList {
		externalIds = append(externalIds, row.WorkExternalUserId)
	}
	externalList, _, err := externalSearchDao.NewExternalExtraCk().ExternalSearchListCk(cc, customerSearchReq.ExternalSearchListBaseReq{
		GetAll: true,
		SearchList: customerSearchReq.SearchList{
			EnterpriseID:    req.EnterpriseID,
			ExternalUserIds: externalIds,
			MemberId:        req.MemberId,
		},
	})
	if err != nil {
		h.ResError(c, err)
		return
	}
	var externalMap = make(map[int64]customerSearchResp.ExternalUserList)
	for _, row := range externalList {
		externalMap[row.ID] = row
	}
	for _, crowdInfo := range crowdList {
		tmp := customerSearchResp.ExternalCrowdList{}
		if external, ok := externalMap[crowdInfo.WorkExternalUserId]; ok {
			tmp.ExternalUserList = external
		}
		tmp.IsContact = crowdInfo.IsContact
		resList = append(resList, tmp)
	}
	h.ResList(c, total, req.Page, req.Size, resList)
}

func (h ExternalCrowdHandle) GetOpenIds(c *gin.Context) {
	var req customerSearchReq.CrowdInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	cc := h.WarpContext(c)
	externalList, err := h.svc.ListByConn(cc, req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	if len(externalList) == 0 {
		h.ResSuccess(c, []string{})
		return
	}
	for _, row := range externalList {
		req.ExternalUserIds = append(req.ExternalUserIds, row.WorkExternalUserId)
	}
	if req.RelationshipType > 0 {
		req.MemberIds = append(req.MemberIds, req.MemberId)
	}
	openIds, err := externalSearchDao.NewExternalExtraCk().ExternalSearchPluckIds(cc, req.SearchList, "openid")
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, openIds)
}

func (h ExternalCrowdHandle) List(c *gin.Context) {
	var req customerSearchReq.CrowdListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	res, err := h.svc.List(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}

func (h ExternalCrowdHandle) Edit(c *gin.Context) {
	var req customerSearchReq.CrowdEditReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.svc.CheckAuth(h.WarpContext(c), customerSearchReq.CrowdCheckReq{
		EnterpriseBaseReq: request.EnterpriseBaseReq{EnterpriseId: req.EnterpriseId},
		Id:                req.CrowdId,
		MemberId:          req.MemberId,
	})
	if err != nil {
		h.ResError(c, err)
		return
	}
	err = h.svc.Edit(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h ExternalCrowdHandle) Del(c *gin.Context) {
	var req customerSearchReq.CrowdEditReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.svc.CheckAuth(h.WarpContext(c), customerSearchReq.CrowdCheckReq{
		EnterpriseBaseReq: request.EnterpriseBaseReq{EnterpriseId: req.EnterpriseId},
		Id:                req.CrowdId,
		MemberId:          req.MemberId,
	})
	if err != nil {
		h.ResError(c, err)
		return
	}
	err = h.svc.Del(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h ExternalCrowdHandle) UpdateSort(c *gin.Context) {
	var req customerSearchReq.CrowdUpdateSortReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.svc.UpdateSort(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}

func (h ExternalCrowdHandle) UpdateContact(c *gin.Context) {
	var req customerSearchReq.CrowdUpdateContactReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	err := h.svc.CheckAuth(h.WarpContext(c), customerSearchReq.CrowdCheckReq{
		EnterpriseBaseReq: request.EnterpriseBaseReq{EnterpriseId: req.EnterpriseId},
		Id:                req.CrowdId,
		MemberId:          req.MemberId,
	})
	if err != nil {
		h.ResError(c, err)
		return
	}
	err = h.svc.UpdateContact(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}
