package micro_scrm_enterprise

import (
	"fmt"

	"micro-scrm/internal/orm/model"
)

type OperationActivitySignInLog struct {
	model.BranchModel
	ActivityId     int64  `gorm:"column:activity_id" json:"activity_id"`           // 活动id
	UnionId        string `gorm:"column:union_id" json:"union_id"`                 // 用户的
	ActivityUserId int64  `gorm:"column:activity_user_id" json:"activity_user_id"` // acvivity_user表的id
	SignInContent  string `gorm:"column:sign_in_content" json:"sign_in_content"`   // 签到内容
	SignInPic      string `gorm:"column:sign_in_pic" json:"sign_in_pic"`           // 签到图片
}

func (o OperationActivitySignInLog) TableName() string {
	return "operation_activity_sign_in_log"
}

func (o OperationActivitySignInLog) AllTableName(enterpriseID int64) string {
	return fmt.Sprintf(" %s%d.%s ", o.DBPrefix(), enterpriseID, o.TableName())
}
