package externalSearchModel

import "micro-scrm/internal/app/models"

type ExternalCrowdModel struct {
	models.BaseModel
	MemberID       int64  `json:"member_id" gorm:"column:member_id;type:bigint(10) unsigned;not null;default:0;comment:导购member表id"`
	Name           string `json:"name" gorm:"column:name;type:varchar(255);comment:名称"`
	SearchConnJson string `json:"search_conn_json" gorm:"column:search_conn_json;type:text;comment:检索条件"`
	Sort           int    `json:"sort" gorm:"column:sort;type:tinyint(1) unsigned;comment:排序，越小排在前面"`
}

// TableName returns the table name of the ExcitationAchievement model
func (ExternalCrowdModel) TableName() string {
	return "external_crowd"
}

const (
	ExternalCrowdNum = 8
)
