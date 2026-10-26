package storeVisitModel

import (
	"micro-scrm/internal/app/models"
)

type StoreVisitConfigModel struct {
	models.BaseModel
	Key   string `json:"key" gorm:"column:key;type:varchar(64);not null;default:'';comment:key"`
	Value string `json:"value" gorm:"column:value;type:text;comment:value"`
}

func (StoreVisitConfigModel) TableName() string {
	return "store_visit_config"
}
