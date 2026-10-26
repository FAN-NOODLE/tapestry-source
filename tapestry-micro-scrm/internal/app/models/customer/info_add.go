package customerModels

import "micro-scrm/internal/app/models"

type CustomerInfoAddModel struct {
	OperateUserId   string `json:"operate_user_id" gorm:"column:operate_user_id;type:varchar(50);not null;comment:操作人Id(企业微信userid)"`
	OperateName     string `json:"operate_name" gorm:"column:operate_name;type:varchar(20);not null;comment:操作人-姓名"`
	Address         string `json:"address" gorm:"column:address;type:varchar(50);not null;comment:地址"`
	ShoeSize        string `json:"shoe_size" gorm:"column:shoe_size;type:varchar(50);not null;comment:鞋码"`
	Month           int    `json:"month" gorm:"column:month;type:int(10);not null;default:0;comment:生日（月）"`
	Day             int    `json:"day" gorm:"column:day;type:int(10);not null;default:0;comment:生日（日）"`
	ExternalUnionId string `json:"external_unionid" gorm:"column:external_unionid;type:varchar(50);not null;comment:外部联系人unionid"`
	models.MicroBaseModel
}

// TableName 外部联系人标签操作记录表
func (CustomerInfoAddModel) TableName() string {
	return "customer_info_add"
}
