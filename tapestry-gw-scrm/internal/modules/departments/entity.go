package departments

// Department 组织结构
type Department struct {
	ID       int64  `json:"id" desc:"组织ID"`
	Name     string `json:"name" desc:"组织名称"`
	Order    int32  `json:"order" desc:"排序"`
	ParentID int64  `json:"parent_id" desc:"父级ID"`
}

// DepartmentListRequest 组织列表请求
type DepartmentListRequest struct {
	Page         int32 `json:"page" desc:"页码"`
	Size         int32 `json:"size" desc:"分页"`
	Pid          int64 `json:"pid" desc:"PID"`
	EnterpriseId int64 `json:"enterprise_id" desc:"企业ID"`
}

// DepartmentDetailRequest 组织详情请求
type DepartmentDetailRequest struct {
	ID       int64  `json:"id" desc:"组织ID"`
	Name     string `json:"name" desc:"组织名称"`
	Order    int32  `json:"order" desc:"排序"`
	ParentID int64  `json:"parent_id" desc:"父级ID"`
}

// DepartmentUpdateRecord 更新记录结果
type DepartmentUpdateRecord struct {
	ID           int64 `json:"id" db:"id" desc:"记录ID"`
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	AdminID      int64 `json:"admin_id" db:"admin_id" desc:"操作人ID"`
	Type         int8  `json:"type" db:"type" desc:"更新类型"`
	Status       int8  `json:"status" db:"status" desc:"更新状态"`
	CreatedAt    int64 `json:"created_at" db:"created_at" desc:"创建时间"`
	UpdatedAt    int64 `json:"updated_at" db:"updated_at" desc:"更新时间"`
}

// UpdateRecordRequest 更新记录请求
type UpdateRecordRequest struct {
	EnterpriseId int64 `json:"enterprise_id" desc:"企业ID"`
}

// MomentListRequest 朋友圈发布选择部门请求
type MomentListRequest struct {
	ThirdId      int64  `json:"third_id" db:"third_id" desc:"部门ID"`
	Name         string `json:"dp_name" db:"name" desc:"部门名称"`
	IsTop        int    `json:"is_top" db:"is_top" desc:"是否顶级"`
	EnterpriseId int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
}

// MomentListResponse 朋友圈发布选择部门返回
type MomentListResponse struct {
	Data []MomentDepartment
}

// MomentDepartment 朋友圈发布选择部门子结构体
type MomentDepartment struct {
	ID                int64  `json:"dp_id" db:"id" desc:"部门ID"`
	DpName            string `json:"dp_name" db:"name" desc:"部门名称"`
	ThirdID           int64  `json:"third_id" db:"third_id" desc:"企业微信中组织ID"`
	ThirdParentID     int64  `json:"third_parent_id" db:"third_parent_id" desc:"企业微信中的父ID"`
	Order             int32  `json:"dp_display_order" db:"order" desc:"排序"`
	IsChildDepartment int    `json:"is_child_department" db:"is_child_department" desc:"是不是子部门"`
}

type LostInfoRequest struct {
	DpId         int64 `json:"dp_id"`
	EnterpriseId int64 `json:"enterprise_id"`
}

type LostConfigRequest struct {
	EnterpriseId int64 `json:"enterprise_id" `
	DpId         int64 `json:"dp_id"`
	UserList     []struct {
		UserId   string `json:"user_id"`
		SendType uint32 `json:"send_type"`
	} `json:"user_ids" binding:"required"`
}
