package excitationModels

import "micro-scrm/internal/app/models"

const (
	TaskTypePoint = iota + 1
	TaskTypeMedal
	TaskTypeIncentive
)

type ExcitationTaskRecordModel struct {
	models.MicroBaseModel
	UniqueStr string `json:"unique_str" gorm:"column:unique_str;type:varchar(32);not null;index:unique_str_index;comment:本次任务的唯一字符串，可以是hash也可以是拼接"`
	Type      int64  `json:"type" gorm:"column:type;type:tinyint(1);not null;default:0;comment:执行类型 0无,1积分相关，2奖章相关，3能力图谱相关"`
}

// TableName returns the table name of the ExcitationTaskRecord model
func (ExcitationTaskRecordModel) TableName() string {
	return "excitation_task_record"
}
