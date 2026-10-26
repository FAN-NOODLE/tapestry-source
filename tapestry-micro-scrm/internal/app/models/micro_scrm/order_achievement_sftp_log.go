package microScrmModel

import "micro-scrm/internal/app/models"

type OrderAchievementSftpLogModel struct {
	models.MicroBaseModel
	ID                int64  `json:"id" gorm:"primaryKey;column:id;type:bigint(20)"`
	FileName          string `json:"file_name" gorm:"column:file_name;type:varchar(255);comment:文件名称"`
	FilePath          string `json:"file_path" gorm:"column:file_path;type:varchar(255);comment:sftp文件路径"`
	FileURL           string `json:"file_url" gorm:"column:file_url;type:varchar(255);comment:cos文件路径"`
	Status            int64  `json:"status" gorm:"column:status;type:tinyint(4);default:0;comment:0.未完成 1.已完成"`
	FileTime          int64  `json:"file_time" gorm:"column:file_time;type:int(11);default:0;comment:同步文件时间"`
	FileNumber        int64  `json:"file_number" gorm:"column:file_number;type:int(11);default:0;comment:文件行数"`
	SyncSuccessNumber int64  `json:"sync_success_number" gorm:"column:sync_success_number;type:int(11);default:0;comment:同步成功行数"`
}

func (OrderAchievementSftpLogModel) TableName() string {
	return "order_achievement_sftp_log"
}
