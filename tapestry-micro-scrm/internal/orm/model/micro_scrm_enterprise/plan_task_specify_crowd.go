package micro_scrm_enterprise

import "micro-scrm/internal/orm/model"

const TaskCrowdTypeSa = 1
const TaskCrowdTypeDept = 2

type PlanTaskSpecifyCrowdModel struct {
	model.BranchModel

	PlanId   int64  `gorm:"column:plan_id;NOT NULL" json:"plan_id"`             // 计划编号
	TaskId   int64  `gorm:"column:task_id;NOT NULL" json:"task_id"`             // 任务编号
	UserId   string `gorm:"column:user_id;NOT NULL" json:"user_id"`             // 推送的用户编号
	UserName string `gorm:"column:user_name;NOT NULL" json:"user_name"`         // 推送的用户名称
	IsBa     int    `gorm:"column:is_ba;NOT NULL" json:"is_ba"`                 // 推送的用户是否是ba 1 否 2 是
	ThirdId  int    `gorm:"column:third_id;NOT NULL;default:0" json:"third_id"` // 企微部门id
	Type     int    `gorm:"column:type;NOT NULL;default:1" json:"type"`         // 范围类型 1 导购 2部门

}

func (PlanTaskSpecifyCrowdModel) TableName() string {
	return "plan_task_specify_crowd"
}

type PlanTaskSpecifyCrowdData struct {
	PlanId   int64  `json:"plan_id"`   // 计划编号
	TaskId   int64  `json:"task_id"`   // 任务编号
	UserId   string `json:"user_id"`   // 推送的用户编号
	UserName string `json:"user_name"` // 推送的用户名称
	IsBa     int    `json:"is_ba"`     // 推送的用户是否是ba 1 否 2 是
	ThirdId  int    `json:"third_id"`  // 企微部门id
	Type     int    `json:"type"`      // 范围类型 1 导购 2部门

}
