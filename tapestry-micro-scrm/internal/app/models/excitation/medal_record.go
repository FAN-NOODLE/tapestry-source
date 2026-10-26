package excitationModels

import "micro-scrm/internal/app/models"

type ExcitationMedalRecordModel struct {
	models.MicroBaseModel

	MedalID       int64  `json:"medal_id" gorm:"column:medal_id;type:bigint(20);not null;index:enterprise_medal_id_index;default:0;comment:奖章表主键id"`
	MedalStyleID  int64  `json:"medal_style_id" gorm:"column:medal_style_id;type:bigint(20);not null;default:0;comment:奖章样式表主键id"`
	MemberID      int64  `json:"member_id" gorm:"column:member_id;type:bigint(20);not null;comment:员工id"`
	IsShow        int32  `json:"is_show" gorm:"column:is_show;type:tinyint(1);not null;default:0;comment:是否已查看 0 未查看 1 已查看"`
	IsReceive     int32  `json:"is_receive" gorm:"column:is_receive;type:tinyint(1);not null;default:0;comment:是否已领取 0 未领取 1 已领取"`
	IsPop         int32  `json:"is_pop" gorm:"column:is_pop;type:tinyint(1);not null;default:0;comment:是否已弹出 0 未弹出 1 已弹出"`
	CoinsRecordID int64  `json:"coins_record_id" gorm:"column:coins_record_id;type:bigint(10) unsigned;not null;default:0;comment:相关金币记录表id"`
	Ranking       string `json:"ranking" gorm:"column:ranking;type:varchar(50);not null;comment:排名文字"`
	RankingValue  string `json:"ranking_value" gorm:"column:ranking_value;type:varchar(255);not null;comment:具体业务数据说明"`
}

// TableName returns the table name of the WjMedalRecord model
func (ExcitationMedalRecordModel) TableName() string {
	return "excitation_medal_record"
}

// 是否已查看 0 未查看 1 已查看
const (
	MedalRecordIsShowFalse = 0
	MedalRecordIsShowTrue  = 1
)

// 是否已领取 0 未领取 1 已领取
const (
	MedalRecordIsReceiveFalse = 0
	MedalRecordIsReceiveTrue  = 1
)

// 是否已弹出 0 未弹出 1 已弹出
const (
	MedalRecordIsPopFalse = 0
	MedalRecordIsPopTrue  = 1
)

type ReceivedNumGroupStyleId struct {
	MedalStyleId int64 `json:"medal_style_id"`
	Total        int64 `json:"total"`
}
