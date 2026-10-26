package drainage

import "micro-scrm/internal/orm/model"

type DrainageContactsModel struct {
	model.MicroScrmBaseModel

	Name               string `gorm:"column:name;type:varchar(255);not null;index:name;comment:二维码名称"`
	ChannelID          int64  `gorm:"column:channel_id;type:bigint(20);not null;default:0;comment:渠道id,标识企业微信state"`
	AutoFriend         int    `gorm:"column:auto_friend;type:int(1);not null;default:0;comment:是否自动成为好友 1 开启 0 不开启"`
	QrCodeURL          string `gorm:"column:qr_code_url;type:varchar(1000);not null;comment:二维码地址"`
	MemberCount        int    `gorm:"column:member_count;type:int(10);not null;default:0;comment:创建员工数"`
	SumAddMemberNumber int    `gorm:"column:sum_add_member_number;type:int(10);not null;default:0;comment:累加添加人数"`
	ConfigID           string `gorm:"column:config_id;type:varchar(255);not null;comment:微信返回配置id"`
	Type               uint32 `gorm:"column:type;type:tinyint(1) unsigned;not null;default:1;comment:欢迎语类型：1-系统联系我，2-裂变活动默认二维码,3-裂变活动发起人二维码"`
}

func (DrainageContactsModel) TableName() string {
	return "drainage_contacts"
}
