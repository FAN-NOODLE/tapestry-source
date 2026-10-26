package right_rule

type AdminDetailRequest struct {
	AdminId int64 `json:"admin_id"`
}

type AdminDetailListRequest struct {
	EnterpriseId int64   `json:"enterprise_id"`
	Role         []int   `json:"role"`
	AdminIds     []int64 `json:"admin_ids"`
	Size         int     `json:"size"`
}
