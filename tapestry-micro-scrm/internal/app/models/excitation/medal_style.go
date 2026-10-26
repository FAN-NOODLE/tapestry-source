package excitationModels

import "micro-scrm/internal/app/models"

type ExcitationMedalStyleModel struct {
	models.MicroBaseModel

	MedalID   int64  `json:"medal_id" gorm:"column:medal_id;type:bigint(20);not null;comment:奖章表主键id"`
	TimeValue int32  `json:"time_value" gorm:"column:time_value;type:smallint(1);not null;default:0;comment:时间值"`
	TimeName  string `json:"time_name" gorm:"column:time_name;type:varchar(20);not null;comment:时间名称"`
	MedalName string `json:"medal_name" gorm:"column:medal_name;type:varchar(100);not null;comment:奖章名称"`
	LightIcon string `json:"light_icon" gorm:"column:light_icon;type:varchar(255);not null;comment:点亮图标"`
	UnlitIcon string `json:"unlit_icon" gorm:"column:unlit_icon;type:varchar(255);not null;comment:未点亮图标"`
}

// TableName returns the table name of the WjMedalStyle model
func (ExcitationMedalStyleModel) TableName() string {
	return "excitation_medal_style"
}
