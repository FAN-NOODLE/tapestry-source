package externalSearchModel

import (
	"micro-scrm/internal/app/models"
)

type ExternalSearchConfigModel struct {
	models.BaseModel
	Key   string `json:"key" gorm:"column:key;type:varchar(64);not null;default:'';comment:key"`
	Value string `json:"value" gorm:"column:value;type:text;comment:value"`
}

func (ExternalSearchConfigModel) TableName() string {
	return "external_search_config"
}
