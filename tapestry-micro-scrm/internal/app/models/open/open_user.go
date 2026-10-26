package openModel

import "micro-scrm/internal/app/models"

const (
	SuperUser = iota + 1 //超级管理
)

type OpenUserModel struct {
	models.MicroBaseModel
	SecretID  int    `json:"secret_id" gorm:"column:secret_id;type:int(11);not null;default:0;comment:账号"`
	SecretKey string `json:"secret_key" gorm:"column:secret_key;type:varchar(64);not null;comment:密码"`
	AuthType  uint32 `json:"auth_type" gorm:"column:auth_type;type:tinyint(3) unsigned;not null;default:1;comment:导购权限类型：1、全部权限 2部分权限"`
	ApiList   string `json:"api_list" gorm:"column:api_list;type:text;not null;comment:api权限列表 英文逗号分割"`
}

// TableName returns the table name of the OpenUser model
func (OpenUserModel) TableName() string {
	return "open_user"
}
