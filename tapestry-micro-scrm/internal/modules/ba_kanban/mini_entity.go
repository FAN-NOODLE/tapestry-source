package ba_kanban

import "micro-scrm/internal/orm/model/micro_scrm_enterprise"

type MiniBaKanBanDataRequest struct {
	EnterpriseID int64   `json:"enterprise_id"`
	StartTime    int64   `json:"start_time"`
	EndTime      int64   `json:"end_time"`
	StoresID     int64   `json:"stores_id"`
	StoresIDs    []int64 `json:"stores_ids"`
	UserID       int64   `json:"user_id"`

	PageSize int `json:"page_size"`
	Page     int `json:"page"`
}
type BaKanBanDataRequest struct {
	Month              int64  `json:"month" binding:"required" desc:"当前月份"`
	Default_mobile     string `json:"default_mobile" binding:"required" desc:"默认查询电话"`
	Dept_leader_mobile string `json:"dept_leader_mobile" desc:"部门电话"`
	Ba_mobile          string `json:"ba_mobile" desc:"BA电话"`
	EnterpriseId       int64  `json:"enterprise_id" binding:"required" desc:"公司编号"`
	PageSize           int    `json:"page_size"`
	Page               int    `json:"page"`
}
type MiNiBaKanBanDataRequest struct {
	Month              int64  `json:"month" binding:"required" desc:"当前月份"`
	Default_mobile     string `json:"default_mobile" binding:"required" desc:"默认查询电话"`
	Dept_leader_mobile string `json:"dept_leader_mobile" desc:"部门电话"`
	Ba_mobile          string `json:"ba_mobile" desc:"BA电话"`
	EnterpriseId       int64  `json:"enterprise_id" binding:"required" desc:"公司编号"`
}

type BaKanDateUploadRequst struct {
	EnterpriseId int64                                        `json:"enterprise_id" binding:"required" desc:"企业编号"`
	Data         []micro_scrm_enterprise.BaMonthlyKanbanModel `json:"data"`
}
