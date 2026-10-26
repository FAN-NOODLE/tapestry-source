package hkStoreModel

import "micro-scrm/internal/app/models"

type StoreManagerModel struct {
	models.BaseModel

	MemberId      int64  `json:"member_id" gorm:"column:member_id;type:bigint(20);not null;default:0;comment:manager member id"`
	AreaIds       string `json:"area_ids" gorm:"column:area_ids;type:varchar(255);not null;default:'';comment:managed area ids"`
	Position      string `json:"position" gorm:"column:position;type:varchar(255);not null;default:'';comment:position"`
	PositionLevel int64  `json:"position_level" gorm:"column:position_level;type:tinyint(1);not null;default:0;comment:position level"`
}

func (StoreManagerModel) TableName() string {
	return "hkmo_store_manager"
}
