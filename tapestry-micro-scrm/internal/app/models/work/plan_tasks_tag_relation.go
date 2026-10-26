package workModels

type PlanTasksTagRelationModel struct {
	ID           int64 `json:"id" gorm:"primaryKey;column:id;type:bigint(20)"`
	EnterpriseID int64 `json:"enterprise_id" gorm:"column:enterprise_id;type:bigint(20);not null;default:0;comment:企业id"`
	PlanID       int64 `json:"plan_id" gorm:"column:plan_id;type:bigint(20);not null;default:0;comment:子任务对应的计划的ID"`
	TagID        int64 `json:"tag_id" gorm:"column:tag_id;type:bigint(20);not null;default:0;comment:标签表ID"`
	StartTime    int   `json:"start_time" gorm:"column:start_time;type:int(10);not null;default:0"`
	EndTime      int   `json:"end_time" gorm:"column:end_time;type:int(10);not null;default:0"`
}

// TableName returns the table name of the PlanTasksTagRelation model
func (PlanTasksTagRelationModel) TableName() string {
	return "plan_tasks_tag_relation"
}
