package micro_scrm_enterprise

import (
	"fmt"
	"micro-scrm/internal/orm/model"
)

type OperationActivityUsers struct {
	model.BranchModel
	ActivityId      int64  `gorm:"column:activity_id;NOT NULL" json:"activity_id"`  // 活动id, operation_activities.id
	UnionId         string `gorm:"column:union_id;NOT NULL" json:"union_id"`        // 用户unionid
	MemberCode      string `gorm:"column:member_code;NOT NULL" json:"member_code"`  // 用户的会员号
	WxName          string `gorm:"column:wx_name;NOT NULL" json:"wx_name"`          // 用户微信昵称
	WxPhone         string `gorm:"column:wx_phone;NOT NULL" json:"wx_phone"`        // 用户微信联系电话
	Content         string `gorm:"column:content;NOT NULL" json:"content"`          // 留言内容
	PictureUrl      string `gorm:"column:picture_url;NOT NULL" json:"picture_url"`  // 用户上传图片地址
	ReceiverName    string `gorm:"column:receiver_name" json:"receiver_name"`       // 收货人姓名
	ReceiverPhone   string `gorm:"column:receiver_phone" json:"receiver_phone"`     // 收货人手机号
	ReceiverAddress string `gorm:"column:receiver_address" json:"receiver_address"` // 收货人地址
	Status          int    `gorm:"column:status;default:1;NOT NULL" json:"status"`  // 审核状态（默认1、未审核，2、审核通过）
}

func (o OperationActivityUsers) TableName() string {
	return "operation_activity_users"
}

func (o OperationActivityUsers) AllTableName(enterpriseID int64) string {
	return fmt.Sprintf(" %s%d.%s ", o.DBPrefix(), enterpriseID, o.TableName())
}
