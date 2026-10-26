package membersModel

import (
	"micro-scrm/internal/app/models"
	departmentsModels "micro-scrm/internal/app/models/departments"
)

type MemberDepartmentsModel struct {
	models.MicroBaseModel

	MemberID     int64                              `json:"member_id" gorm:"primaryKey;column:member_id;type:bigint(20);default:0;comment:成员ID"`
	DepartmentID int64                              `json:"department_id" gorm:"primaryKey;column:department_id;type:bigint(20);index:department_id_index;default:0;comment:组织ID"`
	Member       MemberModel                        `json:"member" gorm:"foreignKey:ID;references:MemberID"`
	Department   departmentsModels.DepartmentsModel `json:"department" gorm:"foreignKey:ID;references:DepartmentID"`
}

// TableName 企业微信成员,组织关系
func (MemberDepartmentsModel) TableName() string {
	return "member_departments"
}
