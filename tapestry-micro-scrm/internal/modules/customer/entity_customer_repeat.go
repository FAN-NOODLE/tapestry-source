// Package customer 客户
package customer

const (
	UserCategoryTypeGroup = 1
	UserCategoryTypeTag   = 2
)

var UserCategoryType = map[string]int{
	"group": UserCategoryTypeGroup,
	"tag":   UserCategoryTypeTag,
}

// RepeatListRequest 重复客户请求
type RepeatListRequest struct {
	EnterpriseID int64        `json:"enterprise_id" db:"enterprise_id" desc:"企业id" validate:"required" required:"企业ID不能为空"`
	Page         int          `json:"page" validate:"min=1" min:"页面参数不能小于1"`
	Size         int          `json:"size" validate:"min=1" min:"分页参数不能小于1"`
	Search       ParamsSearch `json:"search"`
}

// ParamsSearch 参数请求
type ParamsSearch struct {
	CustomerName string `json:"customer_name"`
	DepartmentID int64  `json:"department_id"`
	MemberID     int64  `json:"member_id"`
}

// RepeatListParams 重复客户列表参数
type RepeatListParams struct {
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业id" validate:"required" required:"企业ID不能为空"`
	Page         int    `json:"page" validate:"min=1" min:"页面参数不能小于1"`
	Size         int    `json:"size" validate:"min=1" min:"分页参数不能小于1"`
	CustomerName string `json:"customer_name" db:"customer_name"`
	DepartmentID int64  `json:"department_id" db:"department_id"`
	MemberID     int64  `json:"member_id" db:"member_id"`
}

//重复客户定时计算
type (
	WorkExternalRepeatCustomer struct {
		Id                 int64 `json:"id" db:"id"`
		EnterpriseId       int64 `json:"enterprise_id" db:"enterprise_id"`
		WorkExternalUserId int64 `json:"work_external_user_id" db:"work_external_user_id"`
		MemberNum          int64 `json:"member_num" db:"member_num"`
		CreatedAt          int32 `json:"created_at" db:"created_at"`
		UpdateAt           int32 `json:"update_at" db:"update_at" `
	}
	WorkExternalRepeatCustomerResult struct {
		WorkExternalUserId int64 `json:"work_external_user_id" db:"work_external_user_id"`
		MemberNum          int   `json:"member_num" db:"member_num"`
	}
)

// RepeatListResponse 重复客户列表返回
type RepeatListResponse []struct {
	UserWorkID int64  `json:"user_work_id" db:"work_external_user_id"`
	MemberNum  int    `json:"member_num" db:"member_num"`
	CreateTime int64  `json:"create_time" db:"create_time"`
	Name       string `json:"name" db:"name"`
	Nick       string `json:"nick" db:"nick"`
	Position   string `json:"position" db:"position"`
	Avatar     string `json:"avatar" db:"avatar"`
	UserID     int64  `json:"user_id" db:"user_id"`
}

// WorkUserMemberRelationMemberIds 外部联系人列表
type WorkUserMemberRelationMemberIds []*struct {
	MemberID        int64  `json:"member_id" db:"member_id"`
	MemberUserid    string `json:"member_userid" db:"member_userid"`
	Name            string `json:"name" db:"name"`
	Avatar          string `json:"avatar" json:"avatar"`
	AssistantFlag   int    `json:"assistant_flag"`
	DepartmentName  string `json:"department_name"`
	Channel         string `json:"channel"`
	CreateTime      string `json:"create_time" db:"createtime"`
	LastConnectTime int64  `json:"last_connect_time"`
	ConnectNum      int    `json:"connect_num"`
	IsBindGuide     int    `json:"is_bind_guide" desc:"是否是绑定的导购1=是0=否"`
}

// ExternalRelationByUserIdRequest 外部联系人请求
type ExternalRelationByUserIdRequest struct {
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id" desc:"企业id" validate:"required" required:"企业ID不能为空"`
	UserWorkID   int64 `json:"user_work_id" db:"user_work_id" validate:"required" required:"UserWorkID不能为空"`
	UserID       int64 `json:"user_id" db:"user_id"`
	DepartmentID int64 `json:"department_id" db:"department_id"`
	MemberID     int64 `json:"member_id" db:"member_id"`
	Page         int   `json:"page" validate:"min=1" min:"页面参数不能小于1"`
	Size         int   `json:"size" validate:"min=1" min:"分页参数不能小于1"`
}

// DepartmentListByMemberRequest 组织列表请求
type DepartmentListByMemberRequest struct {
	MemberID int64 `json:"member_id" db:"member_id"`
}

// Department 组织
type Department struct {
	ID                int64  `json:"dp_id" db:"id" desc:"部门ID"`
	DpName            string `json:"dp_name" db:"name" desc:"部门名称"`
	ThirdID           int64  `json:"third_id" db:"third_id" desc:"企业微信中组织ID"`
	ThirdParentID     int64  `json:"third_parent_id" db:"third_parent_id" desc:"企业微信中的父ID"`
	Order             int32  `json:"dp_display_order" db:"order" desc:"排序"`
	IsChildDepartment int    `json:"is_child_department" db:"is_child_department" desc:"是不是子部门"`
}

type LossCustomerAddRelationTime struct {
	Id             int64  `json:"id" db:"id"`
	CreateTime     int64  `json:"create_time" db:"create_time"`
	UserId         string `json:"user_id" db:"user_id" `
	ExternalUserId string `json:"external_user_id" db:"external_user_id"`
}
