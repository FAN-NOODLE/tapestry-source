package customerReq

type CustomerInfoAddReq struct {
	EnterpriseID     int64  `json:"enterprise_id"`
	EnterpriseUserID int64  `json:"enterprise_user_id" binding:"required" comment:"客户ID"`
	Address          string `json:"address"`
	ShoeSize         string `json:"shoe_size"`
	Month            int    `json:"month"`
	Day              int    `json:"day"`
	OperateUserId    string `json:"operate_user_id"`
	OperateName      string `json:"operate_name"`
}

type GetCustomerInfoAddReq struct {
	EnterpriseID int64  `json:"enterprise_id"`
	UnionID      string `json:"union_id" db:"union_id"`
}

type UserInfoAdd struct {
	Address  string `json:"address"`
	ShoeSize string `json:"shoe_size"`
	Month    int    `json:"month"`
	Day      int    `json:"day"`
}

type MyCenterReq struct {
	EnterpriseID int64  `json:"enterprise_id"`
	UnionID      string `json:"union_id" db:"union_id"`
}
