package drainageContactModel

import (
	"micro-scrm/internal/app/models"
)

type DrainageContactMemberModel struct {
	models.MicroBaseModel
	ContactId int64 `gorm:"column:contact_id;type:bigint(20);not null;default:0;comment:联系我表主键"`
	MemberId  int64 `gorm:"column:member_id;type:bigint(20);not null;default:0;comment:导购id"`
	ThirdId   int   `gorm:"column:third_id;type:int(10);not null;default:0;comment:部门id"`
}

// TableName returns the table name of the DisplayTask model
func (DrainageContactMemberModel) TableName() string {
	return "drainage_contact_members"
}
