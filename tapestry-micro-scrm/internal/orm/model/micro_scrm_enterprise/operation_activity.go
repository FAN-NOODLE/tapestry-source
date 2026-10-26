package micro_scrm_enterprise

import (
	"fmt"

	"micro-scrm/internal/orm/model"
)

var STATUSTEXT = map[int64]string{
	1: "未开始",
	2: "进行中",
	3: "已结束",
}

type OperationActivityModel struct {
	model.BranchModel
	Name          string `gorm:"name" json:"name"`
	Type          int64  `gorm:"type" json:"type"`
	ParentId      int64  `gorm:"column:parent_id;default:0" json:"parent_id"` // 依赖的活动编号
	Nums          int64  `gorm:"nums" json:"nums"`
	Status        int64  `gorm:"status" json:"status"`
	StatusText    string `json:"status_text"`
	StartTime     int64  `gorm:"start_time" json:"start_time"`
	StartTimeText string `json:"start_time_text"`
	EndTime       int64  `gorm:"end_time" json:"end_time"`
	EndTimeText   string `json:"end_time_text"`
}

func (o OperationActivityModel) TableName() string {
	return "operation_activities"
}
func (o OperationActivityModel) AllTableName(enterpriseID int64) string {
	return fmt.Sprintf(" %s%d.%s ", o.DBPrefix(), enterpriseID, o.TableName())
}

var USERSTATUSTEXT = map[int64]string{
	1: "待审核",
	2: "审核通过",
}

type UserModel struct {
	model.BranchModel
	WxName          string `gorm:"wx_name" json:"wx_name"`
	MemberCode      string `gorm:"member_code" json:"member_code"`
	ActivityId      int64  `gorm:"activity_id" json:"activity_id"`
	UnionId         string `gorm:"union_id" json:"union_id"`
	WxPhone         string `gorm:"wx_phone" json:"wx_phone"`
	Content         string `gorm:"content" json:"content"`
	PictureUrl      string `gorm:"picture_url" json:"picture_url"`
	ReceiverName    string `gorm:"receiver_name" json:"receiver_name"`
	ReceiverPhone   string `gorm:"receiver_phone" json:"receiver_phone"`
	ReceiverAddress string `gorm:"receiver_address" json:"receiver_address"`
	Status          int64  `gorm:"status" json:"status"`
}
type OperationActivityUserListModel struct {
	model.BranchModel
	WxName          string `gorm:"wx_name" json:"wx_name"`
	ActivityId      int64  `gorm:"activity_id" json:"activity_id"`
	ReceiverName    string `gorm:"receiver_name" json:"receiver_name"`
	ReceiverPhone   string `gorm:"receiver_phone" json:"receiver_phone"`
	ReceiverAddress string `gorm:"receiver_address" json:"receiver_address"`
	Status          int64  `gorm:"status" json:"status"`
	StatusText      string `json:"status_text"`
	CreatedTimeText string `json:"created_time_text"`
}

type OperationActivityUserModel struct {
	UserModel
	StatusText      string `json:"status_text"`
	CreatedTimeText string `json:"created_time_text"`
}

func (o OperationActivityUserModel) TableName() string {
	return "operation_activity_users"
}
func (o OperationActivityUserModel) AllTableName(enterpriseID int64) string {
	return fmt.Sprintf(" %s%d.%s ", o.DBPrefix(), enterpriseID, o.TableName())
}
