package storeVisitModel

import (
	"micro-scrm/internal/app/models"
)

type StoreVisitShowRecordModel struct {
	models.BaseModel
	VisitID   int64                  `json:"visit_id" gorm:"column:visit_id;type:bigint(20);not null;index:visit_id_index;comment:巡店表id"`
	MemberID  int64                  `json:"member_id" gorm:"column:member_id;type:bigint(20);not null;index:member_id_index;comment:导购id"`
	TypeName  string                 `json:"type_name" gorm:"column:type_name;type:varchar(128);not null;default:'';comment:类型名称"`
	Content   string                 `json:"content" gorm:"column:content;type:mediumtext;comment:表现内容"`
	MediaList []StoreVisitMediaModel `json:"media_list" gorm:"foreignKey:RelID;references:ID"`
}

func (StoreVisitShowRecordModel) TableName() string {
	return "store_visit_show_record"
}
