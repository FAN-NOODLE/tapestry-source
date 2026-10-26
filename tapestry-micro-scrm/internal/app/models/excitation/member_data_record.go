package excitationModels

import "micro-scrm/internal/app/models"

type ExcitationMemberDataRecordModel struct {
	models.MicroBaseModel
	MemberID      int64  `json:"member_id" gorm:"column:member_id;type:bigint(20);not null;index:member_id_index;default:0;comment:员工id"`
	Level         string `json:"level" gorm:"column:level;type:varchar(50);not null;comment:等级名称"`
	InviteValue   string `json:"invite_value" gorm:"column:invite_value;type:varchar(50);not null;comment:邀约能力指标"`
	TaskValue     string `json:"task_value" gorm:"column:task_value;type:varchar(50);not null;comment:任务能力指标"`
	InteractValue string `json:"interact_value" gorm:"column:interact_value;type:varchar(50);not null;comment:互动能力指标"`
	TrainValue    string `json:"train_value" gorm:"column:train_value;type:varchar(50);not null;comment:培训能力指标"`
	OnlineValue   string `json:"online_value" gorm:"column:online_value;type:varchar(50);not null;comment:线上能力指标"`
}

// TableName returns the table name of the ExcitationMemberDataRecord model
func (ExcitationMemberDataRecordModel) TableName() string {
	return "excitation_member_data_record"
}
