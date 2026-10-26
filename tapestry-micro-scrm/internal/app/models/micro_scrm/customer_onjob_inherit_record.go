package microScrmModel

import (
	"micro-scrm/internal/app/models"
)

const (
	TransferSuccess = iota + 1
	TransferWait
	TransferRefund
	TransferLimit
	TransferIng
)

type CustomerOnjobInheritRecordModel struct {
	models.MicroBaseModel
	InheritID           int64  `json:"inherit_id" gorm:"column:inherit_id;type:bigint(20);not null;index:enterprise_id_inherit_id_index,inherit_id_external_userid_index;default:0;comment:分配记录id"`
	HandoverMemberID    int64  `json:"handover_member_id" gorm:"column:handover_member_id;type:bigint(20);not null;default:0;comment:原跟进用户id"`
	HandoverUserid      string `json:"handover_userid" gorm:"column:handover_userid;type:varchar(255);not null;index:handover_userid_index;comment:原跟进成员的userid"`
	HandoverName        string `json:"handover_name" gorm:"column:handover_name;type:varchar(255);not null;comment:原跟成员姓名"`
	HandoverMobile      string `json:"handover_mobile" gorm:"column:handover_mobile;type:varchar(15);not null;comment:企业微信手机号"`
	TakeoverMemberID    int64  `json:"takeover_member_id" gorm:"column:takeover_member_id;type:bigint(20);not null;default:0;comment:接替人用户id"`
	TakeoverUserid      string `json:"takeover_userid" gorm:"column:takeover_userid;type:varchar(255);not null;index:enterprise_id_external_id_index;comment:接替人userid"`
	TakeoverName        string `json:"takeover_name" gorm:"column:takeover_name;type:varchar(255);not null;comment:接替人姓名"`
	TakeoverMobile      string `json:"takeover_mobile" gorm:"column:takeover_mobile;type:varchar(15);comment:接替人电话"`
	ExternalID          int64  `json:"external_id" gorm:"column:external_id;type:bigint(20);not null;default:0;comment:外部联系人id"`
	ExternalUserid      string `json:"external_userid" gorm:"column:external_userid;type:varchar(255);not null;index:inherit_id_external_userid_index;comment:外部联系人userid(openid)"`
	ExternalNickname    string `json:"external_nickname" gorm:"column:external_nickname;type:varchar(255);comment:外部联系人昵称"`
	ExternalName        string `json:"external_name" gorm:"column:external_name;type:varchar(255);comment:外部联系人名称"`
	ExternalAvatar      string `json:"external_avatar" gorm:"column:external_avatar;type:varchar(255);comment:外部联系人头像"`
	Status              int32  `json:"status" gorm:"column:status;type:tinyint(4);not null;default:2;comment:状态: 1-接替完毕 2-等待接替 3-客户拒绝 4-接替成员客户达到上限"`
	Errmsg              string `json:"errmsg" gorm:"column:errmsg;type:varchar(255);comment:分配结果"`
	TakeoverTime        uint   `json:"takeover_time" gorm:"column:takeover_time;type:int(11) unsigned;default:0;comment:分配接收时间"`
	RequestCallbackCode int    `json:"request_callback_code" gorm:"column:request_callback_code;type:int(11);default:1;comment:发起分配返回状态 1是默认值"`
}

// TableName returns the table name of the CustomerOnjobInheritRecord model
func (CustomerOnjobInheritRecordModel) TableName() string {
	return "customer_onjob_inherit_record"
}

type OnInheritExternalIds struct {
	ExternalID int64 `json:"external_id,omitempty"`
}
