package customerReq

type UserMemberRelationReq struct {
	EnterpriseID        int64   `json:"enterprise_id" desc:"企业ID"`
	MemberIds           []int64 `json:"member_ids" binding:"required" required:"员工Id集合"`
	WorkExternalUserIds []int64 `json:"work_external_user_ids" binding:"required" required:"外部联系人Id集合"`
}

type UserMemberRelationResp struct {
	MemberId           int64 `json:"member_id"`
	WorkExternalUserId int64 `json:"work_external_user_id"`
	RelationStatus     int   `json:"relation_status"`
}
