package members

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/typex"

// 业务实体（结构体）
type Member struct {
	ID             typex.ID `json:"id" desc:"员工ID"`
	Userid         string   `json:"userid" desc:"企业微信userid"`
	Name           string   `json:"name" desc:"员工姓名"`
	Mobile         string   `json:"mobile" desc:"员工手机号"`
	Position       string   `json:"position" desc:"部门"`
	Gender         int8     `json:"gender" desc:""`
	Email          string   `json:"email" desc:"员工邮箱"`
	IsLeaderInDept string   `json:"is_leader_in_dept" desc:""`
	Avatar         string   `json:"avatar" desc:"头像"`
	ThumbAvatar    string   `json:"thumb_avatar" desc:""`
	Telephone      string   `json:"telephone" desc:"电话"`
	Alias          string   `json:"alias" desc:"别名"`
	Status         int8     `json:"status" desc:"状态"`
	QrCode         string   `json:"qr_code" desc:"二维码"`
	Address        string   `json:"address" desc:"地址"`
	OpenUserid     string   `json:"open_userid" desc:""`
	EnterpriseID   int64    `json:"enterprise_id" desc:"企业ID"`
	CreatedAt      int64    `json:"created_at" desc:"添加时间"`
	UpdatedAt      int64    `json:"updated_at" desc:"更新时间"`
}

// 列表请求参数
type MemberListRequest struct {
	ID           int64  `json:"id" desc:"员工ID"`
	Userid       string `json:"userid" desc:"企业微信userid"`
	EnterpriseID int64  `json:"enterprise_id" desc:"企业ID"`
	DepartmentID int64  `json:"department_id" desc:"部门id"`
	Status       int8   `json:"status" desc:"激活状态1=激活 2=禁用 4=未关注"`
	Name         string `json:"name" desc:"员工姓名"`
	Page         int32  `json:"page" desc:"页码"`
	Size         int32  `json:"size" desc:"分页"`
}

// 员工详情请求参数
type MemberDetailRequest struct {
	ID           int64 `json:"id" desc:"员工ID"`
	EnterpriseID int64 `json:"enterprise_id" desc:"企业ID"`
}

// MomentListRequest 朋友圈创建选择员工请求
type MomentListRequest struct {
	GetNoDp      int64 `json:"get_no_dp" desc:"没有部门的员工"`
	Page         int64 `json:"page" desc:"第几页"`
	Size         int64 `json:"size" desc:"每页条数"`
	EnterpriseId int64 `json:"enterprise_id" desc:"企业ID"`
	Search       struct {
		DpId        int64  `json:"dp_id" desc:"部门ID"`
		MemUsername string `json:"mem_username" desc:"名称"`
		Userid      string `json:"userid"`
		Mobile      string `json:"mobile"`
	} `json:"search" desc:"搜索条件"`
}

// MomentList 朋友圈创建选择员工返回
type MomentListResponse struct {
	List     []MomentList `json:"list"`
	NoDpList []MomentList `json:"no_dp_list"`
	Total    int64        `json:"total" desc:"total"`
	Page     int64        `json:"page" desc:"page"`
	Size     int64        `json:"size" desc:"size"`
}

// MomentList 朋友圈创建选择员工返回子结构体
type MomentList struct {
	UserId             string `json:"user_id" db:"user_id" desc:"members.userid"`
	Name               string `json:"user_name" db:"user_name" desc:"members.name"`
	MemberId           int64  `json:"member_id" db:"member_id" desc:"员工ID"`
	MemSubscribeStatus int    `json:"mem_subscribe_status" db:"mem_subscribe_status" desc:"是不是子部门"`
}

type MemberListByThirdIdsReq struct {
	EnterpriseId int64 `json:"enterprise_id"`
	ThirdIds     []int `json:"third_ids"`
}
