package micro_scrm

import "micro-scrm/internal/orm/model"

const (
	DeleteFriendDeleteTypeExternalUser = 1 // 用户删员工
	DeleteFriendDeleteTypeUser         = 2 // 员工删用户
)

type DeleteFriendRelationRecordData struct {
	DpName string `json:"dp_name"`
	DpID   string `json:"dp_id"`
	UserID string `json:"user_id"`
	Sex    int32  `json:"sex"`
	Nick   string `json:"nick"`
	Name   string `json:"name"`
	Mobile string `json:"mobile"`
	Avatar string `json:"avatar"`
}

type DeleteFriendRelationRecordModel struct {
	model.MicroScrmBaseModel

	ExternalUserID uint64 `json:"external_user_id" gorm:"column:external_user_id;type:bigint(20) unsigned;not null;index:external_user_id;default:0;comment:外部联系人的userid"`
	UserID         string `json:"user_id" gorm:"column:user_id;type:varchar(32);not null;index:user_id;comment:企业服务人员的UserID"`
	DpID           uint64 `json:"dp_id" gorm:"column:dp_id;type:bigint(20) unsigned;not null;index:dp_id;default:0;comment:部门Id"`
	DpName         string `json:"dp_name" gorm:"column:dp_name;type:varchar(50);not null;comment:部门名称"`
	DeleteType     uint32 `json:"delete_type" gorm:"column:delete_type;type:tinyint(1) unsigned;not null;default:0;comment:删除类型：1=用户删员工 2=员工删用户"`
}

// TableName 好友删除关系记录表
func (DeleteFriendRelationRecordModel) TableName() string {
	return "delete_friend_relation_record"
}
