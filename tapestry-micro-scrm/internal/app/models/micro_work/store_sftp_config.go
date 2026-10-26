package microWorkModel

import "micro-scrm/internal/app/models"

type StoreSftpConfigModel struct {
	models.BaseModel
	Host        string `json:"host" gorm:"column:host;type:varchar(30);not null;comment:地址"`
	Port        string `json:"port" gorm:"column:port;type:varchar(30);not null;comment:端口"`
	User        string `json:"user" gorm:"column:user;type:varchar(30);not null;comment:用户"`
	Password    string `json:"password" gorm:"column:password;type:varchar(255);not null;comment:密码"`
	RootPath    string `json:"root_path" gorm:"column:root_path;type:varchar(255);not null;comment:下载文件目录"`
	RecyclePath string `json:"recycle_path" gorm:"column:recycle_path;type:varchar(255);not null;comment:归档文件目录"`
	Brand       string `json:"brand" gorm:"column:brand;type:varchar(255);not null;comment:品牌文件标识"`
}

// TableName returns the table name of the StoreSftpConfig model
func (StoreSftpConfigModel) TableName() string {
	return "store_sftp_config"
}
