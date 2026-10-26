package microScrmModel

import (
	"micro-scrm/internal/app/models"
)

type PlanTaskSpecifyCrowdModel struct {
	models.BaseModel
	ID       uint64 `json:"id" gorm:"primaryKey;column:id;type:bigint(20) unsigned auto_increment;comment:主键"`
	PlanID   uint64 `json:"plan_id" gorm:"column:plan_id;type:bigint(20) unsigned;not null;comment:计划编号"`
	TaskID   uint64 `json:"task_id" gorm:"column:task_id;type:bigint(20) unsigned;not null;index:index_task_id_is_ba;comment:任务编号"`
	UserID   string `json:"user_id" gorm:"column:user_id;type:varchar(64);not null;comment:推送的用户编号"`
	UserName string `json:"user_name" gorm:"column:user_name;type:varchar(128);not null;comment:推送的用户名称"`
	IsBa     int32  `json:"is_ba" gorm:"column:is_ba;type:tinyint(1);not null;index:index_task_id_is_ba;comment:推送的用户是否是ba 1 否 2 是"`
}

// TableName returns the table name of the PlanTaskSpecifyCrowd model
func (PlanTaskSpecifyCrowdModel) TableName() string {
	return "plan_task_specify_crowd"
}
