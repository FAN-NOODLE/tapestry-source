package microScrmModel

import "micro-scrm/internal/app/models"

const (
	ExcitationCoinsRecordAutoType_Points = iota + 1
	ExcitationCoinsRecordAutoType_Mall
	ExcitationCoinsRecordAutoType_Hand
	ExcitationCoinsRecordAutoType_Expire
)

type ExcitationCoinsRecordModel struct {
	models.BaseModel
	MemberID       int64  `json:"member_id" gorm:"column:member_id;type:bigint(10) unsigned;not null;index:member_id_type_deleted_at_index;default:0;comment:导购member表id"`
	RatioRemark    string `json:"ratio_remark" gorm:"column:ratio_remark;type:varchar(255);not null;comment:增加金币时的积分与金币兑换比例"`
	Coins          string `json:"coins" gorm:"column:coins;type:varchar(64);not null;comment:通过该记录获得的金币"`
	SourceType     int    `json:"source_type" gorm:"column:source_type;type:tinyint(1);not null;index:source_type_index;default:0;comment:来源类型 1 积分同步 2 兑礼中心 3 手动补发 4 积分过期"`
	PointsRecordID int64  `json:"points_record_id" gorm:"column:points_record_id;type:bigint(20);not null;default:0;comment:积分记录相关id"`
	MedalRecordID  int64  `json:"medal_record_id" gorm:"column:medal_record_id;type:bigint(20);not null;default:0;comment:奖章记录相关id"`
	CoinsRemark    string `json:"coins_remark" gorm:"column:coins_remark;type:varchar(255);not null;comment:变动原因"`
	OrderNum       string `json:"order_num" gorm:"column:order_num;type:varchar(255);not null;comment:三方订单号"`
}

// TableName returns the table name of the ExcitationCoinsRecord model
func (ExcitationCoinsRecordModel) TableName() string {
	return "excitation_coins_record"
}
