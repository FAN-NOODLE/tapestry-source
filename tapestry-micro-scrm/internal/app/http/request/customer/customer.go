package customerReq

// UpdateRemarkRequest 修改备注参数
type UpdateRemarkRequest struct {
	EnterpriseID   int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	MemberID       int64  `json:"member_id"`
	Mobile         string `json:"mobile" binding:"required"`
	ExternalUserID string `json:"external_userid" db:"external_userid"`
	Remark         string `json:"remark"`
}

type SimpleInfoRequest struct {
	MemberID       int64  `json:"member_id"`
	EnterpriseID   int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	ExternalUserID string `json:"external_userid" db:"external_userid"`
}

type SimpleInfoRes struct {
	UnionId        string `json:"union_id"`
	ExternalUserID string `json:"external_userid"`
	Nick           string `json:"nick"`
	Name           string `json:"name"`
	Avatar         string `json:"avatar"`
}

type OpenIdInfoReq struct {
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	UnionId      string `json:"union_id"`
	OpenID       string `json:"open_id"`
}

type MemberListsReq struct {
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	UnionId      string `json:"union_id"`
}

type UserInfoRequest struct {
	EnterpriseID   int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	UnionId        string `json:"union_id"`
	ExternalUserID string `json:"external_userid" `
}

type UserSaRequest struct {
	UserInfoRequest
	WorkExternalUserID int64 `json:"work_external_user_id" `
}

type UserDetailTipsRequest struct {
	EnterpriseID   int64  `json:"enterprise_id" desc:"企业ID"`
	ExternalUserID string `json:"external_userid" binding:"required" `
	MemberType     int    `json:"member_type" binding:"required"`
	Mobile         string `json:"mobile" binding:"required"`
	MemberID       int64  `json:"member_id"`
	ThirdIds       []int  `json:"third_ids" binding:"required_if=MemberType 4" required_if:"third_ids不能为空" comment:"当前SM管理的门店"`
}

type UserDetailTipsData struct {
	ReserveNum  int64 `json:"reserve_num"`
	SowTaskNum  int64 `json:"sow_task_num"`
	PlanTaskNum int64 `json:"plan_task_num"`
}

type UserBindRelationRequest struct {
	EnterpriseID   int64   `json:"enterprise_id" desc:"企业ID"`
	Mobile         string  `json:"mobile" binding:"required"`
	ExternalUserID string  `json:"external_userid" binding:"required" `
	MemberId       int64   `json:"member_id"  comment:"导购id"`
	MemberType     int     `json:"member_type"  comment:"当前顾问类型"`
	ThirdIds       []int   `json:"third_ids" comment:"当前SM管理的门店"`
	MemberList     []int64 `json:"member_list"  comment:"管理的导购ids"`
}

type UserUnionIdAndMemberId struct {
	Unionid  string `json:"unionid,omitempty"`
	MemberId int64  `json:"member_id"`
}

type UserList struct {
	EnterpriseID      int64    `json:"enterprise_id" comment:"企业ID"`
	ManageType        int      `json:"manage_type"`
	MemberList        []int64  `json:"member_list"`
	CardType          []string `json:"card_type"`
	ImpressionTagList []int64  `json:"impression_tag_list"`
	TagList           []int64  `json:"tag_list"`
	CltType           []string `json:"clt_type"`
	Mobile            string   `json:"mobile"`
	Name              string   `json:"name"`
	Remark            string   `json:"remark"`
	Page              int      `json:"page" comment:"页码"`
	Size              int      `json:"size" comment:"分页数"`
}

type SearchUserIdsReq struct {
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	Name         string `json:"name"`
	UnionCode    string `json:"union_code"`
}

// UpdateDescriptionRequest 修改描述参数
type UpdateDescriptionRequest struct {
	EnterpriseID   int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	MemberID       int64  `json:"member_id"`
	Mobile         string `json:"mobile" binding:"required"`
	ExternalUserID string `json:"external_userid" db:"external_userid"`
	Description    string `json:"description"`
}
