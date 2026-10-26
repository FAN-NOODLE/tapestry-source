package entity

// 有关work标签信息实体

// 标签信息实体实体
type WorkExternalTagGroupEntity struct {
	WorkWeChatGroupID string `db:"work_group_id" desc:"企业微信中标签组ID"`
	GroupName         string `db:"group_name" desc:"企业微信中标签组名称"`
	GroupId           int64  `db:"group_id" desc:"work_external_tag_group.id，对应work_external_tag.group_id"`
	UseNum            int    `db:"use_num" desc:"标签组中标签使用的人数"`
	ValueUpdateAt     int    `db:"value_update_at"`
}

// 标签信息实体
type WorkExternalTagEntity struct {
	WorkWeChatTagID string `json:"work_tag_id" db:"work_tag_id" desc:"work_external_tag.tag_id别名，企业微信标签ID" gorm:"column:work_tag_id"`
	TagId           int64  `json:"tag_id" db:"tag_id" desc:"work_external_tag.id别名"`
	TagGroupId      int64  `json:"tag_group_id" db:"tag_group_id"`
	TagName         string `json:"tag_name" db:"tag_name" desc:"work_external_tag.name标签名称"`
	UseNum          int    `json:"use_num" db:"use_num" desc:"标签使用的人数"`
	IsDelete        int    `json:"is_delete"`
	CreateTime      int64  `json:"create_time" db:"create_time"`
	ExpiredTime     int64  `json:"expired_time" db:"expired_time"`
}

// 标签集合信息
type WorkTagInfoEntity struct {
	WorkExternalTagGroupEntity
	WorkExternalTagEntity
}

// 删除标签组
const (
	GroupOperation = 1
	TagOperation   = 2
)

type DeleteTagGroup struct {
	OperationType int8  // 1删除标签组，2删除标签
	EnterpriseID  int64 `db:"enterprise_id"`
	GroupID       int64 `db:"tag_group_id"`
	DeletedAt     int64 `db:"deleted_at"`
}
