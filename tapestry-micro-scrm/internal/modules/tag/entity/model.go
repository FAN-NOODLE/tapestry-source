package entity

// 标签信息实体（为表work_external_tag中字段）
type WorkExternalTagModel struct {
	ID           int64  `db:"id" desc:"标签ID"`
	EnterpriseID int64  `db:"enterprise_id" desc:"企业ID"`
	GroupID      int64  `db:"group_id" desc:"标签组ID，work_external_tag_group.id"`
	WorkTagID    string `db:"tag_id" desc:"企业微信中的标签ID"`
	Name         string `db:"name" desc:"标签名称，和企业微信中的标签名称一致"`
	Type         int8   `db:"type" desc:"1-企业设置, 2-用户自定义"`
	Order        int    `db:"order" desc:"在企业微信中标签的排序分值"`
	MemberID     int64  `db:"member_id" desc:"自定义标签的员工对应的员工ID，member表"`
	CreateTime   int64  `db:"create_time" desc:""`
	ExpiredTime  int64  `db:"expired_time" desc:""`
	CreatedAt    int64  `db:"created_at" desc:""`
	UpdatedAt    int64  `db:"updated_at" desc:""`
	DeletedAt    int64  `db:"deleted_at" desc:""`
}

// 标签组信息实体（为表work_external_tag_group中字段）
type WorkExternalTagGroupModel struct {
	ID           int64  `db:"id"`
	EnterpriseID int64  `db:"enterprise_id"`
	GroupID      string `db:"group_id" desc:"企业微信标签组ID"`
	GroupName    string `db:"group_name" desc:"标签组名称"`
	Order        int    `db:"order" desc:"标签组次序值"`
	CreateTime   int64  `db:"create_time" desc:""`
	CreatedAt    int64  `db:"created_at" desc:""`
	UpdatedAt    int64  `db:"updated_at" desc:"标签组数据更新时间"`
	DeletedAt    int64  `db:"deleted_at" desc:"标签组数据删除时间"`
}

// 用户分组的类型
const (
	UserGroupTypeAIPL = 1
	UserGroupTypeRFM  = 2
	UserGroupTypeSelf = 3

	// TagGroupBasic 标签组类型
	TagGroupBasic = 1 // 基础标签组
	TagGroupThird = 2 // 第三方标签组

	// TagThirdPart 标签类型
	TagThirdPart = 3 // 第三方标签

	// TaskTag 异步标签删除任务
	TaskTag   = 2
	TaskGroup = 1
)

// 客户分组模型
type UserGroupModelEntity struct {
	ID           int64  `db:"id"`
	EnterpriseID int64  `db:"enterprise_id"`
	Name         string `db:"name"`
	Type         int8   `db:"type"`
	CreatedAt    int64  `db:"created_at"`
	UpdatedAt    int64  `db:"updated_at"`
	DeletedAt    int64  `db:"deleted_at"`
}

// 客户分组详情模型
type UserGroupDetailModelEntity struct {
	ID            int64 `db:"id"`
	EnterpriseID  int64 `db:"enterprise_id"`
	UserGroupID   int64 `db:"user_group_id"`
	TagGroupID    int64 `db:"tag_group_id"`
	TagID         int64 `db:"tag_id"`
	IncludeType   int8  `db:"include_type"`
	CreatedAt     int64 `db:"created_at"`
	UpdatedAt     int64 `db:"updated_at"`
	DeletedAt     int64 `db:"deleted_at"`
	UserGroupType int   `db:"user_group_type"`
	FileID        int64 `db:"file_id"`
	DepartmentID  int64 `db:"department_id"`
	IsExtract     int   `db:"is_extract"`
}

type UnionIDData struct {
	UnionID string `json:"unionid"`
}

type ExportData struct {
	UnionID string `json:"union_id"`
	Result  string `json:"result"`
}
