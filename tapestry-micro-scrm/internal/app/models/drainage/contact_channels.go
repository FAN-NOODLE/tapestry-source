package drainageContactModel

import (
	"micro-scrm/internal/app/models"
)

type DrainageContactChannelsModel struct {
	models.MicroBaseModel
	Name   string                              `gorm:"column:name;type:varchar(255);not null;default:'';index:name;comment:渠道名称"`
	Labels []DrainageContactChannelLabelsModel `json:"labels" gorm:"foreignKey:ChannelID;references:ID"`
}

// TableName returns the table name of the DisplayTask model
func (DrainageContactChannelsModel) TableName() string {
	return "drainage_contact_channels"
}
