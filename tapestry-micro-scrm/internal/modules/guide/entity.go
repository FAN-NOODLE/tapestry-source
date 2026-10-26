package guide

type DDCRequestBase struct {
	Method   string `json:"method"`
	Platform string `json:"platform"`
}

// SyncUserBingGuideRequest 同步绑定关系接口请求参数
type SyncUserBingGuideRequest struct {
	DDCRequestBase
	DataUserBingGuide DataUserBingGuideList `json:"data"`
}

type DataUserBingGuideList struct {
	UserBingGuideList []BasicUserBingGuideInfoEntity `json:"list"`
}

type BasicUserBingGuideInfoEntity struct {
	// 主键ID
	ID           int64  `db:"id"`
	UnionID      string `json:"unionid"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id"`
	ThirdUserID  int64  `json:"third_user_id" db:"third_user_id"`
	MembersID    int64  `db:"members_id"`
	UserID       int64  `db:"user_id"`
	CurUserID    string `json:"cur_user_id" db:"members_user_id"`
	OriUserID    string `json:"ori_user_id"`
	BindTime     int64  `json:"bind_time" db:"bind_time"`
	CreatedAt    int64  `json:"created_at" db:"created_at"`
	UpdatedAt    int64  `json:"updated_at" db:"updated_at"`
	DeletedAt    int64  `json:"deleted_at" db:"deleted_at"`
}

type BatchUserBingGuideInfoEntity struct {
	ThirdUserIds []interface{}
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id"`
	DeletedAt    int64 `json:"deleted_at" db:"deleted_at"`
}

// 总数返回值
type CountResponse struct {
	Count int64 `json:"count" db:"count"`
}

// CustomerGuideRelation 返回信息
type CustomerGuideRelation struct {
	UserID         int64  `json:"user_id" db:"user_id" desc:"scrm 系统的user_id"`
	EnterpriseID   int64  `json:"enterprise_id" db:"enterprise_id"`
	MembersID      int64  `json:"members_id" db:"members_id" desc:"员工ID"`
	MemberUserid   string `json:"member_userid" db:"member_userid"`
	MemberUsername string `json:"member_username" db:"member_username"`
}
