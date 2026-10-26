package drainage

import "micro-scrm/internal/orm/model"

// receive_coupon_status 领券状态
const (
	UnReceived = iota // 未领取
	Received          // 已领取
)

//  reserve_status 预约回复
const (
	ReserveStatusNo  = "no_reserve"  // 否
	ReserveStatusYes = "yes_reserve" // 是
)

type DrainageWelcomesModel struct {
	model.MicroScrmBaseModel

	Content   string `gorm:"column:content;type:text;not null;comment:内容"`
	IsDefault int    `gorm:"column:is_default;type:int(1);not null;default:0;comment:是否是默认欢迎语 1是 0不是"`
	ChannelID int64  `gorm:"column:channel_id;type:bigint(20);not null;default:0;comment:渠道id,标识企业微信state"`
	Type      uint32 `gorm:"column:type;type:tinyint(1) unsigned;not null;default:1;comment:欢迎语类型：1-系统欢迎语，2-裂变活动欢迎语"`
}

func (DrainageWelcomesModel) TableName() string {
	return "drainage_welcomes"
}
