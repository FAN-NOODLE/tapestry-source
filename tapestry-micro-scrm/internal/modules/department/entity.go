package department

import "micro-scrm/internal/pkg/sdk/work"

// SyncContactRecordStatusSuccess 同步通讯录记录同步状态：成功
const SyncContactRecordStatusSuccess = 2

// SyncContactRecordStatusFail 同步通讯录记录同步状态：失败
const SyncContactRecordStatusFail = 1

// SyncContactRecordStatusNull 同步通讯录记录同步状态:未开始
const SyncContactRecordStatusNull = 0

// SyncContactRecordStatusCustom 同步通讯录记录同步状态:推送成功
const SyncContactRecordStatusPush = 3

// SyncContactRecordStatusCustom 同步通讯录记录同步状态:异步处理
const SyncContactRecordStatusCustom = 4

// Department 组织业务实体（结构体）
type Department struct {
	ID            int64  `json:"id" db:"id" desc:"组织ID"`
	Name          string `json:"name" db:"name" desc:"组织名称"`
	NameEn        string `json:"name_en" db:"name_en" desc:"组织英文名称"`
	Order         int32  `json:"order" db:"order" desc:"排序"`
	ThirdID       int64  `json:"third_id" db:"third_id" desc:"企业微信中组织ID"`
	ThirdParentID int64  `json:"third_parent_id" db:"third_parent_id" desc:"企业微信中的父ID"`
	ParentID      int64  `json:"parent_id" db:"parent_id" desc:"父级ID"`
	EnterpriseID  int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	CreatedAt     int64  `json:"created_at" db:"created_at" desc:"创建时间"`
	UpdatedAt     int64  `json:"updated_at" db:"updated_at" desc:"更新时间"`
	DeletedAt     int64  `json:"deleted_at" db:"deleted_at"`
	ChildrenCount int64  `json:"children_count" db:"children_count" desc:"子类数量"`
	Identifier    string
}

// ChildrenForDepartment 部门，子部门
type ChildrenForDepartment struct {
	DepartmentID int64   `json:"dp_id" db:"id" desc:"父部门ID"`
	TmpIds       []byte  `json:"-" db:"child_departments" desc:"-"`
	Children     []int64 `json:"children" db:"-" desc:"子部门ID列表"`
}

// DepartmentDetailRequest 组织详情
type DepartmentDetailRequest struct {
	EnterpriseID int64 `json:"enterprise_id"`
}

// DepartmentListRequest 组织列表
type DepartmentListRequest struct {
	Pid          int64 `json:"pid" db:"third_parent_id"`
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id"`
}

// DepartmentUpdateRecord 更新记录
type DepartmentUpdateRecord struct {
	ID           int64  `json:"id" db:"id" desc:"记录ID"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	RequestID    string `json:"request_id" db:"request_id" desc:"请求ID"`
	MsgID        string `json:"msg_id" db:"msg_id" desc:"消息ID"`
	AdminID      int64  `json:"admin_id" db:"admin_id" desc:"操作人ID"`
	Type         int8   `json:"type" db:"type" desc:"更新类型"`
	Status       int8   `json:"status" db:"status" desc:"更新状态"`
	CreatedAt    int64  `json:"created_at" db:"created_at" desc:"创建时间"`
	UpdatedAt    int64  `json:"updated_at" db:"updated_at" desc:"更新时间"`
}

// TranslateDepartment 企业微信组织转成业务组织
func TranslateDepartment(wd work.WorkDepartment) *Department {
	return &Department{
		ThirdID:       wd.ID,
		Name:          wd.Name,
		NameEn:        wd.NameEn,
		ThirdParentID: wd.Parentid,
		Order:         wd.Order,
	}
}

// BatchDeleteDepartment 批量删除组织
type BatchDeleteDepartment struct {
	DeletedAt int64 `db:"deleted_at"`
	IDs       []interface{}
}

// BatchGetDepartment 批量获取部门
type BatchGetDepartment struct {
	DeletedAt int64 `db:"deleted_at"`
	ThirdIDs  []interface{}
}

// MomentListRequest 朋友圈获取部门请求
type MomentListRequest struct {
	ThirdId      int64  `json:"third_id" db:"third_id" desc:"部门ID"`
	Name         string `json:"dp_name" db:"name" desc:"部门名称"`
	IsTop        int    `json:"is_top" db:"is_top" desc:"是否顶级"`
	EnterpriseId int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
}

// MomentListResponse 朋友圈获取部门返回
type MomentListResponse struct {
	Data []MomentDepartment
}

// MomentDepartment 朋友圈获取部门返回子结构体
type MomentDepartment struct {
	ID                int64  `json:"dp_id" db:"id" desc:"部门ID"`
	DpName            string `json:"dp_name" db:"name" desc:"部门名称"`
	ThirdID           int64  `json:"third_id" db:"third_id" desc:"企业微信中组织ID"`
	ThirdParentID     int64  `json:"third_parent_id" db:"third_parent_id" desc:"企业微信中的父ID"`
	Order             int32  `json:"dp_display_order" db:"order" desc:"排序"`
	IsChildDepartment int    `json:"is_child_department" db:"is_child_department" desc:"是不是子部门"`
	StoreType         int    `json:"store_type" db:"store_type" desc:"部门类型,1为零售,2为奥莱,0为默认"`
}

// BatchGetDpnameByMembers 批量获取部门
type BatchGetDpnameByMembers struct {
	EnterpriseId int64 `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	DeletedAt    int64 `db:"deleted_at"`
	MemberIds    []interface{}
}

// TopicMember 朋友圈列表按员工搜索子结构体
type TopicMember struct {
	MemberId       int64  `json:"member_id" db:"member_id" desc:"员工id"`
	Mark           string `json:"mark" db:"mark" binding:"required" desc:"链接加密mark"`
	UserName       string `json:"user_name" db:"user_name" desc:"员工姓名"`
	DepartmentId   int64  `json:"department_id" db:"department_id" desc:"部门IDS"`
	DepartmentName string `json:"department_name" db:"department_name" desc:"部门名称"`
	DepartmentStr  string `json:"department_str" db:"department_str" desc:"部门名称"`
	TopicCount     int64  `json:"topic_count" db:"topic_count" desc:"朋友圈数"`
	UpdatedAt      int64  `json:"updated_at" db:"updated_at" desc:"更新时间"`
}
