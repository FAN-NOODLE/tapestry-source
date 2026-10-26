package microScrmModel

import (
	"micro-scrm/internal/app/models"
)

type PlanTipsRelationModel struct {
	models.BaseModel
	PlanID int64  `json:"plan_id" gorm:"column:plan_id;type:bigint(20);not null;index:plan_id_index;default:0;comment:主任务id"`
	UserID string `json:"user_id" gorm:"column:user_id;type:varchar(255);not null;comment:导购的工号"`
}

// TableName returns the table name of the PlanTipsRelation model
func (PlanTipsRelationModel) TableName() string {
	return "plan_tips_relation"
}
