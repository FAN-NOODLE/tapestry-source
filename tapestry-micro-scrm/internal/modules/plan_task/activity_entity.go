package plan_task

import "micro-scrm/internal/orm/model/micro_scrm_enterprise"

type BaseReq struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type UserListReq struct {
	ExportUsreReq
	BaseReq
}
type ExportUsreReq struct {
	ActivityId int64 `json:"activity_id"`
	StartTime  int64 `json:"start_time"`
	EndTime    int64 `json:"end_time"`
	Status     int64 `json:"status"`
}

type GetActivityIdReq struct {
	EnterpriseID int64 `json:"enterprise_id"` //EnterpriseID 企业ID
}

type CheckStatusReq struct {
	EnterpriseID int64  `json:"enterprise_id"` //EnterpriseID 企业ID
	ActivityId   int64  `json:"activity_id"`
	UnionId      string `json:"union_id"`
}

type InsertUserReq struct {
	EnterpriseID    int64  `json:"enterprise_id"` //EnterpriseID 企业ID
	UnionId         string `json:"union_id"`
	WxName          string `json:"wx_name"`
	ActivityId      int64  `json:"activity_id"`
	Content         string `gorm:"content" json:"content"`
	PictureUrl      string `gorm:"picture_url" json:"picture_url"`
	ReceiverName    string `gorm:"receiver_name" json:"receiver_name"`
	ReceiverPhone   string `gorm:"receiver_phone" json:"receiver_phone"`
	ReceiverAddress string `gorm:"receiver_address" json:"receiver_address"`
}

type InsertUserResp struct {
	IsJoin bool                                             `json:"is_join"`
	Data   micro_scrm_enterprise.OperationActivityUserModel `json:"data"`
}

type UserInfoReq struct {
	ActivityUserId int64 `json:"activity_user_id"`
}
