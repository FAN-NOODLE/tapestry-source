package micro_scrm

import "micro-scrm/internal/orm/model"

type DepartmentsModel struct {
	model.MicroScrmBaseModel

	ParentID      int64  `json:"parent_id" gorm:"column:parent_id;type:bigint(20);not null;index:parent_id_index;default:0;comment:业务中父部门id"`
	ThirdID       int    `json:"third_id" gorm:"column:third_id;type:int(10);not null;default:0;comment:企业微信中的组织ID"`
	ThirdParentID int    `json:"third_parent_id" gorm:"column:third_parent_id;type:int(10);not null;default:0;comment:企业微信中的组织父ID"`
	Name          string `json:"name" gorm:"column:name;type:varchar(128);not null;comment:部门名称"`
	NameEn        string `json:"name_en" gorm:"column:name_en;type:varchar(128);not null;comment:英文名称"`
	Order         int    `json:"order" gorm:"column:order;type:int(10);not null;default:0;comment:在父部门中的次序值"`
}

func (DepartmentsModel) TableName() string {
	return "departments"
}
