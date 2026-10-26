package micro_scrm_enterprise

import "micro-scrm/internal/orm/model"

type SyncExternalUserSftpModel struct {
	model.BranchModel

	FileName string `gorm:"column:file_name;NOT NULL" json:"file_name"`  // 同步的文件名称
	FileTime int64  `gorm:"column:file_time;default:0" json:"file_time"` // 文件同步完成时间
}

func (SyncExternalUserSftpModel) TableName() string {
	return "sync_external_user_sftp"
}
