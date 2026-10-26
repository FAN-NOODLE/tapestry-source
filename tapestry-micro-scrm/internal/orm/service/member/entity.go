package member

type DeptMemberList struct {
	UserId             string `json:"user_id" db:"user_id" desc:"members.userid"`
	UserName           string `json:"user_name" db:"user_name" desc:"members.name"`
	MemberId           int64  `json:"member_id" db:"member_id" desc:"员工ID"`
	MemSubscribeStatus int    `json:"mem_subscribe_status" db:"mem_subscribe_status" desc:"是不是子部门"`
	DeptName           string `json:"dept_name" db:"dept_name"`
	Mobile             string `json:"mobile" db:"mobile"`
}
