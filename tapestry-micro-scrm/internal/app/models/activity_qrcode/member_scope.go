package activityQRCodeModel

import (
	"micro-scrm/internal/app/models"
)

type ActivityQRCodeMemberScopeModel struct {
	models.BaseModel
	ActivityID int64 `json:"activity_id" gorm:"column:activity_id;type:bigint(20);not null;default:0;index:activity_id_index;comment:活动id"`
	MemberId   int64 `json:"member_id" gorm:"column:member_id;type:bigint(20);not null;default:0;comment:member_id"`
	AreaId     int   `json:"area_id" gorm:"column:area_id;type:bigint(20);not null;default:0;comment:area表id"`
}

func (ActivityQRCodeMemberScopeModel) TableName() string {
	return "activity_qrcode_member_scope"
}
