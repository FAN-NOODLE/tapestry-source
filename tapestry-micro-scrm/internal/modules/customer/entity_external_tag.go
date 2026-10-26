package customer

const (
	ExternalTagTypeCompany  = 1 // 外部联系人标签类型：企业标签
	ExternalTagTypeCustomer = 2 // 外部联系人标签类型：自定义标签
)

// WorkExternalTagGroup 企业微信标签组库
type WorkExternalTagGroup struct {
	ID           int64  `json:"id" db:"id" desc:"主键"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	GroupID      string `json:"group_id" db:"group_id" desc:"企业微信标签组ID"`
	GroupName    string `json:"group_name" db:"group_name" desc:"企业微信标签组名称"`
	Order        int64  `json:"order" db:"order" desc:"次序"`
	CreateTime   int64  `json:"create_time" db:"create_time" desc:"企业微信标签组创建时间"`
	CreatedAt    int64  `json:"created_at" db:"created_at"`
	UpdatedAt    int64  `json:"updated_at" db:"updated_at"`
	DeletedAt    int64  `json:"deleted_at" db:"deleted_at"`
}

// WorkExternalUserTag 外部联系人的标签
type WorkExternalUserTag struct {
	ID            int64  `json:"id" db:"id"`
	EnterpriseID  int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	RelationID    int64  `json:"relation_id" db:"relation_id" desc:"work_external_user_relation表id"`
	ExternalTagId int64  `json:"external_tag_id" db:"external_tag_id" desc:"标签表标签id"`
	CreatedAt     int64  `json:"created_at" db:"created_at"`
	UpdatedAt     int64  `json:"updated_at" db:"updated_at"`
	DeletedAt     int64  `json:"deleted_at" db:"deleted_at"`
	GroupName     string `json:"group_name"`
	TagID         string `json:"tag_id"`
	TagName       string `json:"tag_name"`
	Type          int8   `json:"type"`
}

// WorkExternalTag 企业微信标签库
type WorkExternalTag struct {
	ID           int64  `json:"id" db:"id"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	GroupID      int64  `json:"group_id" db:"group_id" desc:"企业微信标签组ID"`
	Type         int8   `json:"type" db:"type" desc:"标签类型"`
	TagID        string `json:"tag_id" db:"tag_id" desc:"标签id"`
	Name         string `json:"name" db:"name" desc:"标签名称"`
	Order        int64  `json:"order" db:"order" desc:"次序"`
	MemberID     int64  `json:"member_id" db:"member_id"`
	CreateTime   int64  `json:"create_time" db:"create_time" desc:"创建时间"`
	ExpiredTime  int64  `json:"expired_time" db:"expired_time" desc:"过期时间"`
	CreatedAt    int64  `json:"created_at" db:"created_at"`
	UpdatedAt    int64  `json:"updated_at" db:"updated_at"`
	DeletedAt    int64  `json:"deleted_at" db:"deleted_at"`
}

// WorkTagListItem 企业微信返回标签数据格式
type WorkTagListItem struct {
	TagGroupInfo *WorkExternalTagGroup
	Tags         []*WorkExternalTag
}

// SetExternalTags 设置客户企业标签请求参数
type SetExternalTags struct {
	EnterpriseID int64    `json:"enterprise_id" db:"enterprise_id" desc:"企业ID"`
	RelationID   int64    `json:"relation_id" db:"relation_id" desc:"work_external_user_relation表id"`
	AddTag       []string `json:"add_tag" db:"external_tag_id" desc:"新增的标签IDS"`
	RemoveTag    []string `json:"remove_tag" db:"external_tag_id" desc:"删除的标签IDS"`
	CreatedAt    int64    `json:"created_at" db:"created_at"`
	UpdatedAt    int64    `json:"updated_at" db:"updated_at"`
	DeletedAt    int64    `json:"deleted_at" db:"deleted_at"`
}

type WorkExternalRelationTag struct {
	EnterpriseID       int64 `json:"enterprise_id" db:"enterprise_id"`
	RelationID         int64 `json:"relation_id" db:"relation_id"`
	ExternalTagID      int64 `json:"external_tag_id" db:"external_tag_id"`
	MemberID           int64 `json:"member_id" db:"member_id"`
	WorkExternalUserId int64 `json:"work_external_user_id" db:"work_external_user_id"`
}
