package micro_scrm

import "micro-scrm/internal/orm/model"

type UserGroupModel struct {
	model.MicroScrmBaseModel

	Name string `gorm:"column:name;type:varchar(255);not null;comment:客户分组名称"`
	Type int32  `gorm:"column:type;type:tinyint(3) unsigned;not null;default:0;comment:客户分组类型：1-忠诚度模型AIPL,2-RFM模型，3-自定义"`

	UserGroupDetails []UserGroupDetailModel `gorm:"foreignKey:UserGroupID;references:ID"`
}

// TableName returns the table name of the struct_name model
func (UserGroupModel) TableName() string {
	return "user_group"
}
