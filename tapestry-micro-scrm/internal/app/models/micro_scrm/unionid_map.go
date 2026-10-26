package microScrmModel

import (
	"micro-scrm/internal/app/models"
)

type UnionidMapModel struct {
	models.BaseModel
	Nick       string `json:"nick" gorm:"column:nick;type:char(64);not null"`
	Openid     string `json:"openid" gorm:"column:openid;type:char(64);not null"`
	OldUnionid string `json:"old_unionid" gorm:"column:old_unionid;type:char(64);not null"`
	NewUnionid string `json:"new_unionid" gorm:"column:new_unionid;type:char(64);not null;"`
}

// TableName returns the table name of the struct_name model
func (UnionidMapModel) TableName() string {
	return "unionid_map"
}
