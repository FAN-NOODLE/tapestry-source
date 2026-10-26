package source

type CategoryListRequest struct {
	Name         string `json:"name"`
	EnterpriseID int64  `json:"enterprise_id"`
	AdminId      int64  `json:"admin_id"`
	Page         int    `json:"page"`
	Size         int    `json:"size"`
	MemberId     int64  `json:"member_id"`
	IsCollect    int    `json:"is_collect"`
	OriginH5     bool   `json:"origin_h5"`
}
