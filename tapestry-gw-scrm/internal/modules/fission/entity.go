package fission

import "gw-scrm/internal/pkg/sdk"

// ExchangePrizeRecordsEntity 奖品兑换记录
type ExchangePrizeRecordsEntity struct {
	ID               int64  `db:"id" json:"id"`
	EnterpriseID     int64  `db:"enterprise_id" json:"enterprise_id"`     // EnterpriseID 企业id
	ActivityID       int64  `db:"activity_id" json:"activity_id"`         // ActivityID 活动id
	PrizeID          int64  `db:"prize_id" json:"prize_id"`               // PrizeID 奖品id
	FissionUserID    int64  `db:"fission_user_id" json:"fission_user_id"` // FissionUserID 裂变用户关系表id
	UserID           int64  `db:"user_id"`
	CostAssistantNum int    `db:"cost_assistant_num" json:"cost_assistant_num"` // CostAssistantNum 消耗的助力人数
	DeliverName      string `db:"deliver_name" json:"deliver_name"`             // DeliverName 收件人姓名
	DeliverMobile    string `db:"deliver_mobile" json:"deliver_mobile"`         // DeliverMobile 收件人手机号
	DeliverAddress   string `db:"deliver_address" json:"deliver_address"`       // DeliverAddress 收件人地址
	DeliverArea      string `db:"deliver_area" json:"deliver_area"`
	DeliverAreaCode  string `db:"deliver_area_code" json:"deliver_area_code"`
	Num              int    `db:"num"`
	CreatedAt        int64  `db:"created_at" json:"created_at"` // CreatedAt 创建时间
	UpdatedAt        int64  `db:"updated_at" json:"updated_at"` // UpdatedAt 更新时间
	DeletedAt        int64  `db:"deleted_at" json:"deleted_at"` // DeletedAt 删除时间
}

// ExchangeInfo 奖品兑换信息
type ExchangeInfo struct {
	ID               int64  `json:"id" `
	ExchangeTime     int64  `json:"exchange_time" db:"exchange_time"`
	PrizePic         string `json:"prize_pic" db:"prize_pic"`
	PrizeName        string `json:"prize_name" db:"prize_name"`
	CostAssistantNum int    `json:"cost_assistant_num"`
	Num              int    `json:"num" db:"num"`
}

type ExchangeRecordInfo struct {
	ExchangePrizeRecordsEntity
	ExchangeInfo
	Nickname     string `json:"nick_name" db:"nick_name" desc:"用户昵称"`
	ImgUrl       string `json:"img_url" db:"img_url" desc:"头像"`
	ActivityName string `json:"activity_name" db:"activity_name" desc:"活动名称"`
}

type ExchangeDetailRequest struct {
	EnterpriseID  int64  `db:"fission_user_id" json:"enterprise_id"`
	ActivityID    int64  `db:"fission_user_id" json:"activity_id"` //活动id
	FissionUserID int64  `db:"user_id" json:"user_id"`             //裂变用户id
	Page          int64  `json:"page" desc:"页码"`
	Size          int64  `json:"size" desc:"数量"`
	Nickname      string `json:"nick_name" desc:"用户昵称"`
	StartTime     int64  `json:"start_time" desc:"开始时间"`
	EndTime       int64  `json:"end_time" desc:"结束时间"`
	DeliverName   string `db:"deliver_name" json:"deliver_name"`     // DeliverName 收件人姓名
	DeliverMobile string `db:"deliver_mobile" json:"deliver_mobile"` // DeliverMobile 收件人手机号

}
type ExchangeDetailResponse struct {
	sdk.ResponseBase
	List  []*ExchangeRecordInfo `json:"list"`
	Page  int64                 `json:"page"`
	Size  int64                 `json:"size"`
	Total int64                 `json:"total"`
}
