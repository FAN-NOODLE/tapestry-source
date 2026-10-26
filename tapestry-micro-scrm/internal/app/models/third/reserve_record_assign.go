package thirdModel

import (
	"micro-scrm/internal/app/models"
)

type ReserveRecordAssignModel struct {
	models.MicroBaseModel

	DepartmentID int64 `json:"department_id" gorm:"column:department_id;type:bigint(20);not null;index:department_member_id_index;default:0;comment:部门id"`
	MemberID     int64 `json:"member_id" gorm:"column:member_id;type:bigint(20);not null;index:department_member_id_index;default:0;comment:员工id"`
	AssignNumber int   `json:"assign_number" gorm:"column:assign_number;type:int(10);not null;default:0;comment:分配次数"`
}

// TableName 三方预约记录分配记录表
func (ReserveRecordAssignModel) TableName() string {
	return "third_reserve_record_assign"
}
