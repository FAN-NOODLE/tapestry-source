package externalSearchModel

type ExternalRelationshipExtensionModel struct {
	RelationID   int64  `json:"relation_id" gorm:"primaryKey;column:relation_id;type:bigint unsigned;not null"`
	Unionid      string `json:"unionid" gorm:"column:unionid;type:varbinary(64);not null;default:''"`
	UserID       string `json:"user_id" gorm:"column:user_id;type:varbinary(255);not null;default:''"`
	InteractTime int64  `json:"interact_time" gorm:"column:interact_time;type:int unsigned;not null;default:0"`
	ActiveTime   int64  `json:"active_time" gorm:"column:active_time;type:int unsigned;not null;default:0"`
	CreatedAt    int64  `json:"created_at" gorm:"column:created_at;type:int unsigned;not null;default:0"`
	UpdatedAt    int64  `json:"updated_at" gorm:"column:updated_at;type:int unsigned;not null;default:0"`
}

func (ExternalRelationshipExtensionModel) TableName() string {
	return "external_relationship_extension"
}
