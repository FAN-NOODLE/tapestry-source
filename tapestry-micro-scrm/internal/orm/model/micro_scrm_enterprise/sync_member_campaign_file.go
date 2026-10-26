package micro_scrm_enterprise

import "micro-scrm/internal/orm/model"

type SyncMemberCampaignFileModel struct {
	model.BranchModel

	FileName  string `gorm:"column:file_name;NOT NULL" json:"file_name"`    // 文件名
	IsExtract int    `gorm:"column:is_extract;default:0" json:"is_extract"` // 数据文件是否已经提取完成：1: 文件数据提取完成；2、预处理数据完成
}

func (SyncMemberCampaignFileModel) TableName() string {
	return "sync_member_campaign_file"
}
