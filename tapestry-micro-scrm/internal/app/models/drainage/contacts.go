package drainageContactModel

import (
	"micro-scrm/internal/app/models"
)

type DrainageContactsModel struct {
	models.MicroBaseModel
	Name                 string                       `gorm:"column:name;type:varchar(255);not null;index:name;comment:二维码名称"`
	RealUserIds          string                       `gorm:"column:real_user_ids;type:varchar(2000);not null;index:name;comment:实际的二维码用户"`
	ContactGroupId       int64                        `gorm:"column:contact_group_id;type:bigint(20);not null;default:0;comment:联系我活动id（coach)"`
	IsContactGroupSingle int                          `gorm:"column:is_contact_group_single;type:tinyint(1) unsigned;not null;default:0;comment:是否是联系我活动中的单人码 0 否 1是"`
	ChannelID            int64                        `gorm:"column:channel_id;type:bigint(20);not null;default:0;comment:渠道id,标识企业微信state"`
	AutoFriend           int                          `gorm:"column:auto_friend;type:int(1);not null;default:0;comment:是否自动成为好友 1 开启 0 不开启"`
	QrCodeURL            string                       `gorm:"column:qr_code_url;type:varchar(1000);not null;comment:二维码地址"`
	MemberCount          int                          `gorm:"column:member_count;type:int(10);not null;default:0;comment:创建员工数"`
	SumAddMemberNumber   int                          `gorm:"column:sum_add_member_number;type:int(10);not null;default:0;comment:累加添加人数"`
	ConfigID             string                       `gorm:"column:config_id;type:varchar(255);not null;comment:微信返回配置id"`
	Type                 uint32                       `gorm:"column:type;type:tinyint(1) unsigned;not null;default:1;comment:欢迎语类型：1-系统联系我，2-裂变活动默认二维码,3-裂变活动发起人二维码"`
	Members              []DrainageContactMemberModel `json:"members" gorm:"foreignKey:ContactId;references:ID"`
}

// TableName returns the table name of the DisplayTask model
func (DrainageContactsModel) TableName() string {
	return "drainage_contacts"
}

var DrainageContactGroupSingle = 1
