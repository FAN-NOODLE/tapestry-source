package marketing

import "micro-scrm/internal/orm/model"

const (
	MSR_ChannelType_SMS = 1 + iota
	MSR_ChannelType_WxMsg
	MSR_ChannelType_AutoCoupon
	MSR_ChannelType_Official
)

const (
	MSR_Status_No = 1 + iota
	MSR_Status_Success
	MSR_Status_Fail
)

type MarketingStrategyRecordsModel struct {
	model.BaseModel

	UserID          int64  `gorm:"column:user_id;type:bigint(20);not null;default:0;comment:用户id"`
	UnionID         string `gorm:"column:union_id;type:varchar(255);not null;index:union_id_index;uniqueIndex:uin;comment:union_id"`
	StrategyID      int64  `gorm:"column:strategy_id;type:bigint(20);not null;uniqueIndex:uin;default:0;comment:策略id"`
	ChannelType     int    `gorm:"column:channel_type;type:bigint(20);not null;uniqueIndex:uin;default:0;comment:渠道类型,1：短信；2：企微群发；3：自动发券；4：公众号；"`
	RelationID      int64  `gorm:"column:relation_id;type:bigint(20);not null;index:relation_id_index;default:0;comment:触发模板id"`
	ChannelID       int64  `gorm:"column:channel_id;type:bigint(20);not null;default:0;comment:渠道id"`
	Status          int32  `gorm:"column:status;type:tinyint(1);not null;default:1;comment:执行状态；1=未执行，2=执行成功，3=执行失败"`
	UserName        string `gorm:"column:user_name;type:varchar(50);not null;comment:用户名称"`
	Mobile          string `gorm:"column:mobile;type:varchar(20);not null;comment:用户电话"`
	AddressUserName string `gorm:"column:address_user_name;type:varchar(50);not null;comment:联系人名称"`
	OrderID         string `gorm:"column:order_id;type:varchar(255);not null;comment:订单号"`
	StartTime       int    `gorm:"column:start_time;type:int(11);not null;default:0;comment:开始时间"`
	EndTime         int    `gorm:"column:end_time;type:int(11);not null;default:0;comment:结束时间"`
	TouchedAt       int    `gorm:"column:touched_at;type:int(11);not null;default:0;comment:触达时间"`
	ExternalUserID  string `gorm:"column:external_userid;type:varchar(64);not null;index:exuid_msg_idx;comment:外部联系人userid"`
	MsgID           string `gorm:"column:msg_id;type:varchar(64);not null;index:exuid_msg_idx;comment:群发返回msgid"`
}

func (MarketingStrategyRecordsModel) TableName() string {
	return "marketing_strategy_records_v2"
}
