package microScrmModel

import (
	"micro-scrm/internal/app/models"
)

type NotificationDepartmentModel struct {
	models.BaseModel
	Name   string `json:"name" gorm:"column:name;type:varchar(255);not null;comment:部门名称"`
	Status int32  `json:"status" gorm:"column:status;type:tinyint(1);not null;default:0;comment:部门的禁用状态，0禁用，1正常"`
}

// TableName returns the table name of the NotificationDepartment model
func (NotificationDepartmentModel) TableName() string {
	return "notification_department"
}
