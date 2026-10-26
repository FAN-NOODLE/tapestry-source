package dashboardModel

import (
	"micro-scrm/internal/app/models"
)

type DashboardLogModel struct {
	models.MicroBaseModel
	FileType string `json:"file_type" gorm:"column:file_type;type:varchar(255);not null;comment:文件类型"`
	FileName string `json:"file_name" gorm:"column:file_name;type:varchar(255);not null;comment:文件名称"`
	FileURL  string `json:"file_url" gorm:"column:file_url;type:varchar(255);not null;comment:文件cos地址"`
}

// TableName returns the table name of the DashboardLog model
func (DashboardLogModel) TableName() string {
	return "dashboard_log"
}
