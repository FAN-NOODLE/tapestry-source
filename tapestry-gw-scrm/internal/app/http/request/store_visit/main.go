package storeVisitReq

import "gw-scrm/internal/app/http/request"

type ListReq struct {
	request.EnterpriseBaseReq
	request.PageReq
	StoreName      string `json:"store_name"`
	DeptName       string `json:"dept_name"`
	Status         int    `json:"status"`
	StartTime      int64  `json:"start_time"`
	EndTime        int64  `json:"end_time"`
	VisitStartTime int64  `json:"visit_start_time"`
	VisitEndTime   int64  `json:"visit_end_time"`
}

type StoreVisitBaseReq struct {
	request.EnterpriseBaseReq
	Id int64 `json:"id"`
}

type DashboardReq struct {
	request.EnterpriseBaseReq
	request.PageReq
	DeptName  string `json:"dept_name"`
	Type      int    `json:"type"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
	AreaId    int64  `json:"area_id"`
}
