package micro_scrm

import "micro-scrm/internal/orm/model"

type MemberDepartmentsModel struct {
	model.MicroScrmBaseModel

	MemberID     int64 `json:"member_id" gorm:"column:member_id;type:bigint(20);not null;index:member_id_index;default:0;comment:成员ID"`             // 成员ID
	DepartmentID int64 `json:"department_id" gorm:"column:department_id;type:bigint(20);not null;index:department_id_index;default:0;comment:组织ID"` // 组织ID 	// 数据删除时间

	Department DepartmentsModel `gorm:"foreignKey:DepartmentID;references:ID"`
	Member     MemberModel      `gorm:"foreignKey:MemberID;references:ID"`
}

func (MemberDepartmentsModel) TableName() string {
	return "member_departments"
}
