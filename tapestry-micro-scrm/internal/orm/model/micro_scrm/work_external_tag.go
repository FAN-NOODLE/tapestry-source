package micro_scrm

import (
	"fmt"
	"micro-scrm/internal/orm/model"
)

const (
	WET_Type_Enterprise = 1 + iota
	WET_Type_Custom
	WET_Type_CT
)

type WorkExternalTagModel struct {
	model.MicroScrmBaseModel

	GroupID    int64  `gorm:"column:group_id;type:bigint(20);not null;index:group_id_index;comment:group表标签组id"`
	TagID      string `gorm:"column:tag_id;type:varchar(255);not null;index:tag_id_index;comment:标签id"`
	Name       string `gorm:"column:name;type:varchar(255);not null;comment:标签名称"`
	Type       int32  `gorm:"column:type;type:tinyint(1);not null;default:0;comment:1-企业设置, 2-用户自定义"`
	Order      int    `gorm:"column:order;type:int(10);not null;default:0;comment:次序值"`
	MemberID   int64  `gorm:"column:member_id;type:bigint(20) unsigned;not null;default:0;comment:自定义标签对应员工id"`
	CreateTime int    `gorm:"column:create_time;type:int(10);not null;default:0;comment:创建时间"`
}

// TableName returns the table name of the struct_name model
func (WorkExternalTagModel) TableName() string {
	return "work_external_tag"
}

func (WorkExternalTagModel) DBPrefix() string {
	return "micro_scrm"
}

// 查询后的值替换
//func (w WorkExternalTagModel) AfterFind(tx *gorm.DB) (err error) {
//	return
//}

func (w WorkExternalTagModel) AllTableName(enterpriseID int64) string {
	return fmt.Sprintf(" %s%d.%s ", w.DBPrefix(), enterpriseID, w.TableName())
}
