package tagModels

import "micro-scrm/internal/app/models"

const (
	ImpressionTagActionTypeAdd = iota + 1 // 1 新增客户标签
	ImpressionTagActionTypeDel            // 2 删除客户标签
)

type ImpressionTagOperateLogModel struct {
	OperateUserId    string `json:"operate_user_id" gorm:"column:operate_user_id;type:varchar(50);not null;comment:操作人Id(企业微信userid)"`
	OperateName      string `json:"operate_name" gorm:"column:operate_name;type:varchar(20);not null;comment:操作人-姓名"`
	ThirdTagGroupID  int    `json:"third_tag_group_id" gorm:"column:third_tag_group_id;type:int(10);not null;default:0;comment:三方标签组ID"`
	ThirdTagID       string `json:"third_tag_id" gorm:"column:third_tag_id;type:varchar(50);not null;index:index_tag_id;comment:三方标签ID"`
	ExternalUsername string `json:"external_username" gorm:"column:external_username;type:varchar(50);not null;comment:外部联系人姓名"`
	ExternalUnionId  string `json:"external_unionid" gorm:"column:external_unionid;type:varchar(50);not null;comment:外部联系人unionid"`
	ActionType       int32  `json:"action_type" gorm:"column:action_type;type:tinyint(1);not null;default:0;comment:操作类型：1 新增客户标签 2删除客户标签"`
	models.MicroBaseModel
}

// TableName 外部联系人标签操作记录表
func (ImpressionTagOperateLogModel) TableName() string {
	return "impression_tag_operate_log"
}
