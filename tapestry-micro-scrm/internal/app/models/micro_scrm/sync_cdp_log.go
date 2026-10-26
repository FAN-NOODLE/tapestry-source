package microScrmModel

import (
	"micro-scrm/internal/app/models"
)

type SyncCdpLogModel struct {
	models.MicroBaseModel
	RelationID int64  `json:"relation_id" gorm:"column:relation_id;type:bigint(20);not null"`
	Result     string `json:"result" gorm:"column:result;type:varchar(255);not null"`
}

// TableName returns the table name of the SyncCdpLog model
func (SyncCdpLogModel) TableName() string {
	return "sync_cdp_log"
}
