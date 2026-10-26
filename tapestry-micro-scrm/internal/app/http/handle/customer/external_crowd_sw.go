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

type ExternalCrowdSwHandle struct {
	handle.BaseHandle
	svc *customerSvc.ExternalCrowd
}

func NewExternalCrowdSwHandle() ExternalCrowdSwHandle {
	return ExternalCrowdSwHandle{
		svc: customerSvc.NewExternalCrowd(),
	}
}

func (h ExternalCrowdSwHandle) Add(c *gin.Context) {
	var req customerSearchReq.AddCrowdReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	cc := h.WarpContext(c)
	var searchValue customerSearchReq.ExternalSearchListBaseSwReq
	helper.Unmarshal(req.Search.Value, &searchValue)
	searchValue.EnterpriseID = req.EnterpriseId
	searchValue.MemberId = req.MemberId
	searchValue.ThirdIds = req.ThirdIds
	// 获取所有的external_userId
	searchValue.RelationshipType = customerSearchReq.IsUnionTure
	ids, err := customerSvc.NewExternalSearchSw().GetExternalSearchCk(cc, searchValue, "id")
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

func (h ExternalCrowdSwHandle) Info(c *gin.Context) {
	var req customerSearchReq.CrowdInfoSwReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	cc := h.WarpContext(c)
	req.SearchListSw.MemberId = req.MemberId
	searchQuerySwResp := externalSearchDao.NewExternalExtraCkSw().SearchExternalOptionQuery(req.SearchListSw)
	crowdList, total, err := h.svc.PageListByExternalConnSw(cc, req, searchQuerySwResp)
	if err != nil {
		h.ResError(c, err)
		return
	}
	resList := make([]customerSearchResp.ExternalCrowdSwList, 0)
	if len(crowdList) == 0 {
		h.ResList(c, 0, req.Page, req.Size, resList)
		return
	}
	var externalIds []int64
	for _, row := range crowdList {
		externalIds = append(externalIds, row.WorkExternalUserId)
	}
	externalList, _, err := externalSearchDao.NewExternalExtraCkSw().ExternalSearchListCk(cc, customerSearchReq.ExternalSearchListBaseSwReq{
		GetAll: true,
		SearchListSw: customerSearchReq.SearchListSw{
			EnterpriseID:    req.EnterpriseID,
			ExternalUserIds: externalIds,
			MemberId:        req.MemberId,
		},
	})
	if err != nil {
		h.ResError(c, err)
		return
	}
	var externalMap = make(map[int64]customerSearchResp.ExternalUserListSw)
	for _, row := range externalList {
		externalMap[row.ID] = row
	}
	for _, crowdInfo := range crowdList {
		tmp := customerSearchResp.ExternalCrowdSwList{}
		if external, ok := externalMap[crowdInfo.WorkExternalUserId]; ok {
			tmp.ExternalUserListSw = external
		}
		tmp.IsContact = crowdInfo.IsContact
		resList = append(resList, tmp)
	}
	h.ResList(c, total, req.Page, req.Size, resList)
}

func (h ExternalCrowdSwHandle) GetOpenIds(c *gin.Context) {
	var req customerSearchReq.CrowdInfoSwReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ResError(c, err)
		return
	}
	cc := h.WarpContext(c)
	externalList, err := h.svc.ListByConnSw(cc, req)
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
	openIds, err := externalSearchDao.NewExternalExtraCkSw().ExternalSearchPluckIds(cc, req.SearchListSw, "openid")
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, openIds)
}

func (h ExternalCrowdSwHandle) List(c *gin.Context) {
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

func (h ExternalCrowdSwHandle) Edit(c *gin.Context) {
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

func (h ExternalCrowdSwHandle) Del(c *gin.Context) {
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

func (h ExternalCrowdSwHandle) UpdateSort(c *gin.Context) {
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

func (h ExternalCrowdSwHandle) UpdateContact(c *gin.Context) {
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

func (h ExternalCrowdSwHandle) BatchUpdateContact(c *gin.Context) {
	var req customerSearchReq.BatchCrowdUpdateContactReq
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
	ids, err := externalSearchDao.NewExternalExtraCkSw().GetExternalIdsByOpenIds(h.WarpContext(c), req.EnterpriseId, req.OpenIds)
	if err != nil {
		return
	}
	req.WorkExternalUserIds = ids
	err = h.svc.BatchUpdateContact(h.WarpContext(c), req)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResOk(c)
}
