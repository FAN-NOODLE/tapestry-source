package microScrmModel

import "micro-scrm/internal/app/models"

type PlanTaskCrowdDataLogModel struct {
	models.BaseModel
	CrowdID    string `json:"crowd_id" gorm:"column:crowd_id;type:varchar(255);not null;default:0"`
	TaskID     int64  `json:"task_id" gorm:"column:task_id;type:bigint(10);not null;default:0;comment:子任务id"`
	Unionid    string `json:"unionid" gorm:"column:unionid;type:varchar(255);not null"`
	SaCode     string `json:"sa_code" gorm:"column:sa_code;type:varchar(255);not null;comment:导购的sacode"`
	Name       string `json:"name" gorm:"column:name;type:varchar(255);not null;comment:导购姓名"`
	ImportRes  string `json:"import_res" gorm:"column:import_res;type:varchar(255);not null;comment:导入结果"`
	CustomerID string `json:"customer_id" gorm:"column:customer_id;type:varchar(255);not null;comment:用户的customerid信息"`
}

// TableName returns the table name of the PlanTaskCrowdDataLog model
func (PlanTaskCrowdDataLogModel) TableName() string {
	return "plan_task_crowd_data_log"
}
