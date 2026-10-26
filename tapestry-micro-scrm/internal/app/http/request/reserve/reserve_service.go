package reserveReq

import microScrmModel "micro-scrm/internal/app/models/micro_scrm"

type RoleRules struct {
	ID                 int64  `json:"id"` // is_dp为false时，id为member_id
	ThirdId            int    `json:"third_id"`
	Name               string `json:"name"` // is_dp为false时，name为member_name;is_dp为ture,name为store_name
	IsDP               bool   `json:"is_dp"`
	StoreType          int32  `json:"store_type"`
	MemSubscribeStatus int    `json:"mem_subscribe_status"`
	UserId             string `json:"user_id"`
}

type ReserveServiceCreateReq struct {
	EnterpriseBase
	ReserveName   string      `json:"reserve_name" comment:"服务类型名称"`
	ReserveImg    string      `json:"reserve_img" comment:"服务类型的图片"`
	ReserveTime   int32       `json:"reserve_time" comment:"服务时长"`
	EndTime       int         `json:"end_time" comment:"截止时间"`
	ReserveStatus int32       `json:"reserve_status" comment:"服务类型状态"`
	ServiceType   int32       `json:"service_type"`
	StartTime     int         `json:"start_time"`
	ServiceRole   int32       `json:"service_role"`
	RoleRules     []RoleRules `json:"role_rules"`
	TimeSlot      []TimeSlot  `json:"time_slot"`
}

type ReserveServiceUpdateReq struct {
	EnterpriseBase
	ID            int64       `json:"id" binding:"required" comment:"ID"`
	ReserveName   string      `json:"reserve_name" comment:"服务类型名称"`
	ReserveImg    string      `json:"reserve_img" comment:"服务类型的图片"`
	ReserveTime   int32       `json:"reserve_time" comment:"服务时长"`
	EndTime       int         `json:"end_time" comment:"截止时间"`
	ReserveStatus int32       `json:"reserve_status" comment:"服务类型状态"`
	ServiceType   int32       `json:"service_type"`
	StartTime     int         `json:"start_time"`
	ServiceRole   int32       `json:"service_role"`
	RoleRules     []RoleRules `json:"role_rules"`
	TimeSlot      []TimeSlot  `json:"time_slot"`
}

type ReserveServiceDeleteReq struct {
	EnterpriseBase
	ID int64 `json:"id" binding:"required" comment:"ID"`
}

type ReserveServiceListReq struct {
	EnterpriseBase
	PageReq
	ReserveName   string  `json:"reserve_name" comment:"服务类型名称"`
	ReserveIds    []int64 `json:"reserve_ids" comment:"服务类型ids"`
	ReserveStatus int     `json:"reserve_status" comment:"是否有效, -1全部，0为有效，1为失效"`
	StartTime     int64   `json:"start_time"`
	EndTime       int64   `json:"end_time"`
}

type ReserveServiceListRes struct {
	microScrmModel.ReserveServiceModel
	RoleRules []RoleRules `json:"role_rules"`
	TimeSlot  []TimeSlot  `json:"time_slot"`
}

type ReserveServiceListByMemberReq struct {
	EnterpriseBase
	Mobile string `json:"mobile"`
}

type TimeSlot struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type ReserveServiceListByMemberRes struct {
	microScrmModel.ReserveServiceModel
	TimeSlot []TimeSlot `json:"time_slot"`
}
