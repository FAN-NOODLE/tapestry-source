package drainageContactModel

import (
	"micro-scrm/internal/app/models"
)

type DrainageContactGroupModel struct {
	models.BaseModel
	Name       string                          `gorm:"column:name;type:varchar(255);not null;index:name;comment:二维码名称"`
	MemberIds  string                          `gorm:"column:member_ids;type:varchar(255);not null;index:name;comment:选择的人员id,英文逗号分割"`
	ThirdIds   string                          `gorm:"column:third_ids;type:varchar(255);not null;index:name;comment:选择的部门id,英文逗号分割"`
	ChannelID  int64                           `gorm:"column:channel_id;type:bigint(20);not null;default:0;comment:渠道id,标识企业微信state"`
	AutoFriend int                             `gorm:"column:auto_friend;type:int(1);not null;default:0;comment:是否自动成为好友 1 开启 0 不开启"`
	Type       int                             `gorm:"column:type;type:tinyint(1) unsigned;not null;default:1;comment:欢迎语类型：1-系统联系我，2-裂变活动默认二维码,3-裂变活动发起人二维码"`
	SourceType int                             `gorm:"column:source_type;type:tinyint(1);not null;default:0;comment:来源类型：0-后台创建，1-门店导购码占位，2-服务预约码占位"`
	Contacts   []DrainageContactsModel         `json:"contacts,omitempty" gorm:"foreignKey:ContactGroupId;references:ID"`
	Welcomes   []DrainageContactsWelcomesModel `json:"welcomes,omitempty" gorm:"foreignKey:ContactGroupId;references:ID"`
}

// TableName returns the table name of the DisplayTask model
func (DrainageContactGroupModel) TableName() string {
	return "drainage_contact_group"
}

const DrainageContactGroupStore = 1 //门店码+员工码

const (
	DrainageContactGroupSourceConsole      = 0 // 后台创建
	DrainageContactGroupSourceStoreWorker  = 1 // 门店导购码占位
	DrainageContactGroupSourceReserveGuide = 2 // 服务预约码占位
)

//ALTER TABLE `micro_scrm`.`drainage_contacts`
//ADD COLUMN `contact_group_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '联系我活动id（coach）' AFTER `type`
