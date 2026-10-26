package entity

type BossEnterpriseAdminDetail struct {
	AdminId       int64                        `json:"id"`
	EnterpriseId  int64                        `json:"enterprise_id"`
	Username      string                       `json:"username"`
	Account       string                       `json:"account"`
	Password      string                       `json:"password"`
	Salt          string                       `json:"salt"`
	LastLoginTime int32                        `json:"last_login_time"`
	Role          string                       `json:"role"`
	Status        int8                         `json:"status"`
	IsSuper       int8                         `json:"is_super"`
	Enterprise    BossEnterpriseDetailResponse `json:"enterprise"`
}

type BossEnterpriseAdminDetailRequest struct {
	AdminId      string `json:"admin_id"`
	Account      string `json:"account"`
	MemberId     string `json:"member_id"`
	EnterpriseId string `json:"enterprise_id"`
	Mobile       string `json:"mobile"`
}

type BossEnterpriseAdminSaveRequest struct {
	EnterpriseId int64  `json:"enterprise_id"`
	AdminId      int64  `json:"admin_id"`
	UserName     string `json:"username"`
	Account      string `json:"account"`
	OperatorId   int64  `json:"operator_id"`
	RoleId       int64  `json:"role_id"`
	MemberId     string `json:"member_id"`
}
