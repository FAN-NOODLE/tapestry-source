package microScrmModel

import "micro-scrm/internal/app/models"

type PlanTaskImportDataModel struct {
	models.BaseModel
	FileID             int64  `json:"file_id" gorm:"column:file_id;type:int(10);not null;default:0"`
	TaskID             int64  `json:"task_id" gorm:"column:task_id;type:bigint(20);not null;default:0"`
	UnionId            string `json:"unionid" gorm:"column:unionid;type:varchar(255);not null"`
	Said               string `json:"said" gorm:"column:said;type:varchar(255);not null;comment:导购的SAID"`
	Name               string `json:"name" gorm:"column:name;type:varchar(255);not null;comment:导购姓名"`
	WorkExternalUserId string `json:"work_external_user_id" gorm:"column:work_external_user_id;type:varchar(255);not null;comment:用户的外部联系人id"`
	CustomerID         string `json:"customer_id" gorm:"column:customer_id;type:varchar(255);not null;comment:用户的customerid信息"`
	IsMatch            int    `json:"is_match" gorm:"column:is_match;type:tinyint(2);not null;default:1;comment:是否匹配，默认为1，如果customerid查不到用户手机号，则匹配失败，为2"`
	UserName           string `json:"user_name" gorm:"column:user_name;type:varchar(255);not null;comment:顾客姓名"`
}

// TableName returns the table name of the PlanTaskImportData model
func (PlanTaskImportDataModel) TableName() string {
	return "plan_task_import_data"
}
