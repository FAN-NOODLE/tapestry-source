package customerReq

type PageParam struct {
	Page int `json:"page" binding:"min=1" comment:"页码"`
	Size int `json:"size" binding:"min=10,max=99999999" comment:"分页数"`
}

type CustomerBaseReq struct {
	EnterpriseID int64 `json:"enterprise_id" binding:"required" comment:"企业ID"`
}

type CustomerListReq struct {
	CustomerBaseReq
	PageParam
	HandoverUserId       string `json:"handover_userid" comment:"离职成员的userid"`
	HandoverName         string `json:"handover_name" comment:"离职成员姓名"`
	HandoverDepartmentId int64  `json:"handover_department_id" comment:"离职成员所属门店"`
	Status               int    `json:"status"  comment:"状态: 0 未分配 1已分配"`
}

type ExternalListReq struct {
	SearchName string `json:"search_name"`
	HandOverBase
	TakeoverUserId string `json:"takeover_userid"`
	PageParam
}

type ExternalListRes struct {
	ExternalUserid   string `json:"external_userid"`
	ExternalNickname string `json:"external_nickname"`
	ExternalAvatar   string `json:"external_avatar"`
}

type ExternalListResult struct {
	Total int64             `json:"total"`
	List  []ExternalListRes `json:"list"`
}

type HandOverBase struct {
	HandoverUserId string `json:"handover_userid"  binding:"required"  comment:"离职成员的userid"`
	CustomerBaseReq
}
type GroupChatList struct {
	HandOverBase
	PageParam
}

type GroupChatInfo struct {
	ID     int64  `json:"id"`
	ChatID string `json:"chat_id"`
	Name   string `json:"name"`
}

type GroupChatRes struct {
	Total int64           `json:"total"`
	List  []GroupChatInfo `json:"list"`
}

type TakeOverInfo struct {
	TakeoverUserid         string `json:"takeover_userid"`
	TakeoverName           string `json:"takeover_name"`
	TakeoverDepartmentName string `json:"takeover_department_name"`
}

type TakeOverList struct {
	Total int64          `json:"total"`
	List  []TakeOverInfo `json:"list"`
}

type syncInsertData struct {
	ID     int64  `json:"id"`
	ChatID string `json:"chat_id"`
	Name   string `json:"name"`
}

type TakeOverExternalInfo struct {
	TakeOverTime       int    `json:"take_over_time"`
	ExternalStatusName string `json:"external_status_name"`
	ExternalListRes
}
type TakeOverExternalRes struct {
	Total int64                  `json:"total"`
	List  []TakeOverExternalInfo `json:"list"`
}
type TakeOverGroupInfo struct {
	ChatID string `json:"chat_id"`
	Name   string `json:"name"`
	TakeOverInfo
	TakeoverTime int64 `json:"takeover_time"`
	Status       int   `json:"status"`
}

type TakeOverGroupRes struct {
	Total int64               `json:"total"`
	List  []TakeOverGroupInfo `json:"list"`
}

type SyncCustomDataRes struct {
	CustomerBaseReq
	ExeSync bool `json:"exe_sync"`
}

type InsertChatNumList struct {
	HandoverUserid string `json:"handover_userid"`
	Num            int    `json:"num"`
}
