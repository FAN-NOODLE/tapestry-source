package fission

import "micro-scrm/internal/orm/model"

type FissionUserModel struct {
	model.MicroScrmBaseModel

	ActivityID        int64 `gorm:"column:activity_id;type:bigint(20);not null;index:enterprise_activity_user_idx;default:0;comment:活动id"`
	UserID            int64 `gorm:"column:user_id;type:bigint(20);not null;index:enterprise_activity_user_idx;default:0;comment:用户id"`
	TotalAssistantNum int   `gorm:"column:total_assistant_num;type:int(10);not null;default:0;comment:累计助力人数"`
	CurAssistantNum   int   `gorm:"column:cur_assistant_num;type:int(10);not null;default:0;comment:当前剩余助力人数"`
	ContactID         int64 `gorm:"column:contact_id;type:bigint(20);not null;index:contact_id_idx;default:0;comment:联系我二维码id，drainage_contacts表id"`
	FissionNum        uint  `gorm:"column:fission_num;type:int(10) unsigned;not null;default:0;comment:好友裂变数"`
	RunAwayNum        uint  `gorm:"column:run_away_num;type:int(10) unsigned;not null;default:0;comment:流失人数"`
	JoinGroupNum      uint  `gorm:"column:join_group_num;type:int(10) unsigned;not null;default:0;comment:群成员裂变数"`
	LeaveGroupNum     uint  `gorm:"column:leave_group_num;type:int(10) unsigned;not null;default:0;comment:群成员退群数"`
}

func (FissionUserModel) TableName() string {
	return "fission_user"
}
