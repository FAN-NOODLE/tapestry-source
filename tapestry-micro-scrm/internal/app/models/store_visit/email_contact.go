package storeVisitModel

import (
	"micro-scrm/internal/app/models"
)

type StoreVisitEmailContactModel struct {
	models.BaseModel
	MemberID int64  `json:"member_id" gorm:"column:member_id;type:bigint(20);not null;default:0;index:member_id_index;comment:成员ID"`
	Email    string `json:"email" gorm:"column:email;type:varchar(32);not null;default:'';index:email_index;comment:邮箱"`
}

func (StoreVisitEmailContactModel) TableName() string {
	return "store_visit_email_contact"
}
