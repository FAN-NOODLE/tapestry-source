package microScrm

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
