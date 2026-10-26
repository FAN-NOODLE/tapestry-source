package microWorkModel

import "micro-scrm/internal/app/models"

type StoreEmployeeModel struct {
	models.MicroBaseModel

	Name           string `json:"name" gorm:"column:name;type:varchar(32);not null;comment:姓名"`
	AliasName      string `json:"alias_name" gorm:"column:alias_name;type:varchar(32);not null;comment:别名"`
	EmployeeCode   string `json:"employee_code" gorm:"column:employee_code;type:varchar(32);not null;uniqueIndex:unique_employee_code_index;comment:员工编号"`
	JobLocation    string `json:"job_location" gorm:"column:job_location;type:varchar(255);not null;comment:工作地点"`
	JobPosition    string `json:"job_position" gorm:"column:job_position;type:varchar(255);not null;comment:职位"`
	Gender         uint32 `json:"gender" gorm:"column:gender;type:tinyint(1) unsigned;not null;default:0;comment:性别 1男 2女"`
	Email          string `json:"email" gorm:"column:email;type:varchar(255);not null;comment:邮箱"`
	Status         uint32 `json:"status" gorm:"column:status;type:tinyint(1) unsigned;not null;default:0;comment:状态 1正常 2异常"`
	SyncStatus     int32  `json:"sync_status" gorm:"column:sync_status;type:tinyint(1);not null;default:0;comment:企微同步状态 0未同步 1已同步"`
	SyncAt         int    `json:"sync_at" gorm:"column:sync_at;type:int(10);not null;default:0;comment:同步时间"`
	ThirdStoreCode string `json:"third_store_code" gorm:"column:third_store_code;type:varchar(32);not null;index:enterprise_store_code_index;comment:三方门店code"`
	WechatUserID   string `json:"wechat_user_id" gorm:"column:wechat_user_id;type:varchar(64);not null;index:unique_user_id_index;comment:企微user_id"`
	Mobile         string `json:"mobile" gorm:"column:mobile;type:varchar(50);not null;uniqueIndex:unique_mobile_index;comment:电话"`
}

// TableName 三方员工表（业务架构）
func (StoreEmployeeModel) TableName() string {
	return "store_employee"
}
