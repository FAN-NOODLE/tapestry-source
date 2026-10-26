package entity

// 请求标签信息的数据集合
type TagParam struct {
	EnterpriseID int64 `db:"enterprise_id"`
	TagGroupInfoRequest
}

// 搜索查询标签信息参数集合
type SearchTagParam struct {
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id"`
	Name         string `json:"name" db:"tag_group_name"`
	SearchType   int8
	Type         int
}

// 后台管理，手动删除标签的operation参数
type TagOperationParam struct {
	EnterpriseID int64 `json:"enterprise_id"`
	TagId        int64 `json:"tag_id"`
	TagGroupID   int64 `json:"tag_group_id"`
}

// 更新手动标签名称参数
type UpdateManualTagParam struct {
	EnterpriseID int64  `db:"enterprise_id"`
	TagID        int64  `db:"id"`
	NewTagName   string `db:"tag_name"`
}

// 新增标签参数
type CreateManualTagParam struct {
	EnterpriseID int64 `db:"enterprise_id"`
	TagGroupID   int64 `db:"tag_group_id"`
	TagName      string
	ExpiredTime  int64 `json:"expired_time" db:"end_time"`
}

// 创建标签组
type CreateManualTagGroupParams struct {
	EnterpriseID int64 `db:"enterprise_id"`
	CreateGroupRequest
}

// 查询客户分组信息参数
type QueryUserGroupParam struct {
	EnterpriseID int64 `db:"enterprise_id"`
	UserGroupID  int64 `db:"user_group_id"`
	IncludeType  int8  `db:"include_type"`
}
