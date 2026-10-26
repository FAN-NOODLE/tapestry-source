package departmentsModels

import "micro-scrm/internal/app/models"

type MemberDepartmentsModel struct {
	models.MicroBaseModel

	MemberID     int64 `json:"member_id" gorm:"primaryKey;column:member_id;type:bigint(20);index:member_id_index;default:0;comment:成员ID"`
	DepartmentID int64 `json:"department_id" gorm:"primaryKey;column:department_id;type:bigint(20);index:department_id_index;default:0;comment:组织ID"`
}

// TableName returns the table name of the MemberDepartments model
func (MemberDepartmentsModel) TableName() string {
	return "member_departments"
}
