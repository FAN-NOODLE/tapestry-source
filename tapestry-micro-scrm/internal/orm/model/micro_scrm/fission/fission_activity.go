package fission

import (
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm/drainage"
)

const (
	FA_FissionType_AddEmp = 1 + iota
	FA_FissionType_AddGroup
	FA_FissionType_VIP
)

type FissionActivityModel struct {
	model.MicroScrmBaseModel

	Name               string                         `gorm:"column:name;type:varchar(120);not null;comment:活动名称"`
	Desc               string                         `gorm:"column:desc;type:text;not null;comment:活动描述"`
	Status             int32                          `gorm:"column:status;type:tinyint(1);not null;default:1;comment:活动状态：1-上架中，2-已下架"`
	ChannelID          int64                          `gorm:"column:channel_id;type:bigint(20);not null;index:channel_id_index;default:0;comment:渠道id"`
	StartTime          int                            `gorm:"column:start_time;type:int(10);not null;default:0;comment:活动开始时间"`
	EndTime            int                            `gorm:"column:end_time;type:int(10);not null;default:0;comment:活动结束时间"`
	ExchangeEndTime    int                            `gorm:"column:exchange_end_time;type:int(10);not null;default:0;comment:兑换奖品截至时间"`
	FissionType        int32                          `gorm:"column:fission_type;type:tinyint(1);not null;default:1;comment:裂变类型：1-添加员工，2-添加员工且入群"`
	CondType           int32                          `gorm:"column:cond_type;type:tinyint(1);not null;default:1;comment:参与条件：1-添加任意员工，2-添加裂变员工"`
	MutualType         int32                          `gorm:"column:mutual_type;type:tinyint(1);not null;default:2;comment:互斥类型：1-互斥，2-不互斥"`
	DrainageWelcomesID int64                          `gorm:"column:drainage_welcomes_id;type:bigint(20);not null;default:0;comment:欢迎语id"`
	DrainageWelcomes   drainage.DrainageWelcomesModel `gorm:"foreignKey:DrainageWelcomesID;references:ID"`
	ContactID          int64                          `gorm:"column:contact_id;type:bigint(20);not null;default:0;comment:活动默认联系我二维码id，drainage_contacts表id"`
	DrainageContacts   drainage.DrainageContactsModel `gorm:"foreignKey:ContactID;references:ID"`
	UserNum            uint                           `gorm:"column:user_num;type:int(10) unsigned;not null;default:0;comment:发起人数"`
	FissionNum         uint                           `gorm:"column:fission_num;type:int(10) unsigned;not null;default:0;comment:好友裂变数"`
	RunAwayNum         uint                           `gorm:"column:run_away_num;type:int(10) unsigned;not null;default:0;comment:流失人数"`
	JoinGroupNum       uint                           `gorm:"column:join_group_num;type:int(10) unsigned;not null;default:0;comment:群成员裂变数"`
	LeaveGroupNum      uint                           `gorm:"column:leave_group_num;type:int(10) unsigned;not null;default:0;comment:群成员退群数"`

	FissionChatGroupList []FissionChatGroupModel `gorm:"foreignKey:activity_id;reference:ID"`
}

func (FissionActivityModel) TableName() string {
	return "fission_activity"
}
