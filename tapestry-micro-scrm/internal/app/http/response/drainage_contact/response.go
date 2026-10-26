package drainageContactResp

type GroupBaseResp struct {
	ID                int64    `json:"id"`
	Name              string   `json:"name"`
	Type              int      `json:"type"`
	CodeNumber        int      `json:"code_number"`
	ChannelName       string   `json:"channel_name"`
	Tags              []string `json:"tags"`
	AddCustomerNumber int64    `json:"add_customer_number"`
	HasWelcome        int      `json:"has_welcome"`
	CreatedAt         int64    `json:"created_at"`
}

type FriendsNumberResp struct {
	State string `json:"state"`
	Total int64  `json:"total"`
}

type GroupInfoListResp struct {
	GroupBaseResp
	ContactList []ContactInfo `json:"contact_list"`
	Total       int64         `json:"total"`
}

type GroupDetailResp struct {
	ID             int64            `json:"id" db:"id" desc:"id"`
	Name           string           `json:"name" binding:"required" `
	ChannelId      int64            `json:"channel_id" binding:"required" desc:"渠道id"`
	AutoFriend     int              `json:"auto_friend"   db:"auto_friend" desc:"是否自动开启"`
	Type           int              `json:"type"  db:"type" desc:"类型 0 多人共享码（原本的） 1 员工码+门店码"`
	HasWelcome     int              `json:"has_welcome"`
	MemberData     []MemberData     `json:"member_data" db:"member_data" desc:"员工数据"`
	DepartmentData []DepartmentData `json:"department_data" db:"department_data" desc:"部门数据"`
	WelcomeList    []WelcomeInfo    `json:"welcome_list"`
}

type WelcomeInfo struct {
	Type    int                      `json:"type" commit:"0 通用 1 会员 2非会员"`
	Content []map[string]interface{} `json:"content" binding:"required" desc:"内容"`
}

// MemberData 联系我员工信息
type MemberData struct {
	MemberId   int64  `json:"member_id" db:"member_id" desc:"员工ID"`
	MemberName string `json:"member_name" db:"member_name" desc:"员工名称"`
}

// DepartmentData 联系我部门信息
type DepartmentData struct {
	ThirdId        int    `json:"third_id" db:"third_id" desc:"部门id"`
	DepartmentName string `json:"department_name" db:"department_name" desc:"部门名称"`
}

// RealMember  联系我员工信息
type RealMember struct {
	MemberId   int64  `json:"member_id" db:"member_id" desc:"员工ID"`
	MemberName string `json:"member_name" db:"member_name" desc:"员工名称"`
	ThirdId    int    `json:"third_id"`
}

// UnverifiedUser 未实名/未激活员工信息
type UnverifiedUser struct {
	Userid string `json:"userid"` // 员工userid
	Name   string `json:"name"`   // 员工姓名
}

// GroupSaveResp 活动保存响应
type GroupSaveResp struct {
	ID              int64            `json:"id"`
	UnverifiedUsers []UnverifiedUser `json:"unverified_users,omitempty"` // 未实名认证的员工列表
}

type ContactInfo struct {
	ID                int64            `json:"id"`
	Name              string           `json:"name" `
	CreatedAt         int64            `json:"created_at" desc:"创建时间"`
	Members           []string         `json:"members"  desc:"员工数据"`
	Departments       []DepartmentData `json:"departments"  desc:"部门数据"`
	RealMembers       []RealMember     `json:"real_members"`
	AddCustomerNumber int64            `json:"add_customer_number"`
	QrCodeUrl         string           `json:"qr_code_url"  db:"qr_code_url" desc:"二维码地址"`
}
