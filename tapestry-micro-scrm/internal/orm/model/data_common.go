package model

type GroupUserInfo struct {
	GroupChatID int64  `json:"group_chat_id" db:"group_chat_id" comment:"对应的是clickhouse里的字段"`
	UserID      string `json:"user_id" db:"user_id" comment:"对应的是clickhouse里的字段"`
}
type FissionChatGroup struct {
	ChatID     int64 `json:"chat_id"`
	ActivityID int64 `json:"activity_id"`
}
type ActivityFissionInfo struct {
	ExternalUserID string `json:"external_user_id"`
	ISVip          int    `json:"is_vip"`
}

type FissionUserInfo struct {
	AssistantUserID int64  `json:"assistant_user_id"`
	ExternalUserID  string `json:"external_user_id"`
}
type UnionInfo struct {
	UnionID         string `json:"union_id" db:"unionid" comment:"对应的是clickhouse里的字段"`
	OperationStatus int    `json:"operation_status" db:"status" comment:"对应的是clickhouse里的字段"`
}
