package microScrmModel

import "micro-scrm/internal/app/models"

type InviteRecordModel struct {
	models.BaseModel
	InviteTemplateID int64  `json:"invite_template_id" gorm:"column:invite_template_id;type:bigint(20);not null;index:invite_template_id_idx;default:0;comment:邀约模板id"`
	ReserveOrderID   int64  `json:"reserve_order_id" gorm:"column:reserve_order_id;type:bigint(20);not null;index:reserve_order_id_idx;default:0;comment:邀约订单id"`
	MemberID         int64  `json:"member_id" gorm:"column:member_id;type:bigint(20);not null;default:0;comment:邀约单导购id"`
	MemberUserid     string `json:"member_userid" gorm:"column:member_userid;type:varchar(255);not null;comment:邀约单导购userid"`
	UserName         string `json:"user_name" gorm:"column:user_name;type:varchar(255);not null;comment:邀约单用户名称"`
	UserUnionid      string `json:"user_unionid" gorm:"column:user_unionid;type:varchar(255);not null;index:user_unionid_idx;comment:邀约单用户unionid"`
	UserOpenid       string `json:"user_openid" gorm:"column:user_openid;type:varchar(255);not null;comment:邀约单用户openid"`
	StyleType        int    `json:"style_type" gorm:"column:style_type;type:int(11);not null;comment:选择的样式"`
	BackgroundUrl    string `json:"background_url" gorm:"column:background_url;type:varchar(255);not null;comment:背景url"`
	SendContent      string `json:"send_content" gorm:"column:send_content;type:varchar(255);not null;comment:编辑后发送的内容"`
}

// TableName returns the table name of the InviteRecord model
func (InviteRecordModel) TableName() string {
	return "invite_record"
}
