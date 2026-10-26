package excitationModels

import (
	"micro-scrm/internal/app/models"
)

type ExcitationMedalTaskRecordModel struct {
	models.MicroBaseModel

	MedalID  int64 `json:"medal_id" gorm:"column:medal_id;type:bigint(20);not null;comment:奖章表id"`
	TimeType int32 `json:"time_type" gorm:"column:time_type;type:tinyint(1);not null;default:0;comment:时间类型 0无,1月,2季度,3年"`
	Status   int32 `json:"status" gorm:"column:status;type:tinyint(0);not null;default:2;comment:执行状态 1成功,2执行中,3失败"`
}

// TableName returns the table name of the ExcitationMedalTaskRecord model
func (ExcitationMedalTaskRecordModel) TableName() string {
	return "excitation_medal_task_record"
}

// 执行状态 1成功,2执行中,3失败
const (
	MedalTaskRecordStatusDefault = 0
	MedalTaskRecordStatusOk      = 1
	MedalTaskRecordStatusRunning = 2
	MedalTaskRecordStatusFail    = 3
)
