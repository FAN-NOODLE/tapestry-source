package customer

// 企微员工和外部联系人

// 关系查询参数
type QueryRelationParams struct {
	EnterpriseID   int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	ExternalUserID string `json:"external_userid" db:"work_external_user.external_userid" desc:"微信外部联系人ID"`
	MemberID       int64  `json:"member_id" db:"member_id" desc:"企微员工ID"`
}

// ExternalUserRelationEntity 关系信息实体
type ExternalUserRelationEntity struct {
	RelationID     int64  `json:"relation_id" db:"relation_id" desc:"work_external_user_relation.id"`
	EnterpriseID   int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	ExternalUserID string `json:"external_user_id" db:"external_user_id" desc:"外部联系人ID"`
	MemberID       int64  `json:"member_id" db:"member_id" desc:"企业员工ID"`
}
