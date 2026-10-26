package externalSearchModel

import (
	"micro-scrm/internal/app/models"
)

type ExternalTagModel struct {
	Unionid   string           `gorm:"column:unionid;type:varchar(255);not null;comment:外部联系人的unionid"`
	TagType   string           `gorm:"column:tag_type;type:varchar(255);not null;comment:标签类型"`
	TagKey    string           `gorm:"column:tag_key;type:varchar(255);not null;comment:标签key"`
	TagValue  string           `gorm:"column:tag_value;type:varchar(255);not null;comment:成员设置的外部联系人备注"`
	CreatedAt int64            `gorm:"column:created_at;type:int(10);not null;default:0" json:"created_at"`
	UpdatedAt int64            `gorm:"column:updated_at;type:int(10);not null;default:0" json:"updated_at"`
	DeletedAt models.DeletedAt `gorm:"softDelete;column:deleted_at;type:int(10) unsigned;not null;default:0" json:"deleted_at"`
}

// TableName returns the table name of the WorkExternalUserRelation model
func (ExternalTagModel) TableName() string {
	return "external_tag"
}

const SWAtAmount = "at_amount" //季度消费
