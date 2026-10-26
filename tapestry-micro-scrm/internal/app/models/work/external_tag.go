package workModels

import (
	"micro-scrm/internal/app/models"
)

type WorkExternalTagModel struct {
	models.MicroBaseModel

	GroupID             int64  `json:"group_id" gorm:"column:group_id;type:bigint(20);not null;index:group_id_index;comment:group表标签组id"`
	TagID               string `json:"tag_id" gorm:"column:tag_id;type:varchar(255);not null;index:tag_id_index;comment:标签id"`
	Name                string `json:"name" gorm:"column:name;type:varchar(255);not null;comment:标签名称"`
	Type                int32  `json:"type" gorm:"column:type;type:tinyint(1);not null;default:0;comment:1-企业设置, 2-用户自定义"`
	Order               int    `json:"order" gorm:"column:order;type:int(10);not null;default:0;comment:次序值"`
	MemberID            uint64 `json:"member_id" gorm:"column:member_id;type:bigint(20) unsigned;not null;default:0;comment:自定义标签对应员工id"`
	ExpiredTime         int    `json:"expired_time" gorm:"column:expired_time;type:int(10);not null;default:0;comment:过期时间"`
	CreateTime          int    `json:"create_time" gorm:"column:create_time;type:int(10);not null;default:0;comment:创建时间"`
	ThirdTagID          string `json:"third_tag_id" gorm:"column:third_tag_id;type:varchar(100);not null;comment:tapestry同步第三方的标签id"`
	ThirdTagDescription string `json:"third_tag_description" gorm:"column:third_tag_description;type:varchar(255);not null;comment:tapestry同步三方标签的描述"`
}

// TableName returns the table name of the WorkExternalTag model
func (WorkExternalTagModel) TableName() string {
	return "work_external_tag"
}

const (
	WorkExternalTagTypeQy   = iota + 1 // 1 企业设置
	WorkExternalTagTypeUser            // 2 用户自定义
	_
	MomentTaskTempTag // 企业群发类型朋友圈任务的临时标签
)

const (
	ChatTagThirdType = "2" // coach的企微沟通意愿标签
)
