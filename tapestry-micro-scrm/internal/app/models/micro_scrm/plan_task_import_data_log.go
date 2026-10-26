package microScrmModel

import "micro-scrm/internal/app/models"

type PlanTaskImportDataLogModel struct {
	models.BaseModel
	FileID     int64  `json:"file_id" gorm:"column:file_id;type:int(10);not null;default:0"`
	TaskID     int64  `json:"task_id" gorm:"column:task_id;type:bigint(20);not null;default:0"`
	UnionId    string `json:"unionid" gorm:"column:unionid;type:varchar(255);not null"`
	Said       string `json:"said" gorm:"column:said;type:varchar(255);not null;comment:导购的SAID"`
	Name       string `json:"name" gorm:"column:name;type:varchar(255);not null;comment:导购姓名"`
	ImportRes  string `json:"import_res" gorm:"column:import_res;type:varchar(255);not null;comment:导入结果"`
	CustomerID string `json:"customer_id" gorm:"column:customer_id;type:varchar(255);not null;comment:用户的customerid信息"`
}

// TableName returns the table name of the PlanTaskImportDataLog model
func (PlanTaskImportDataLogModel) TableName() string {
	return "plan_task_import_data_log"
}
