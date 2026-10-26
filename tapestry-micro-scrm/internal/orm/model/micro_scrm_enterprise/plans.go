package micro_scrm_enterprise

import (
	"fmt"

	"micro-scrm/internal/orm/model"
)

type PlanModel struct {
	model.BranchModel

	Name            string `gorm:"column:name;NOT NULL" json:"name"`
	TaskNum         int    `gorm:"column:task_num;default:0;NOT NULL" json:"task_num"`
	OperationUserId int64  `gorm:"column:operation_user_id;default:0;NOT NULL" json:"operation_user_id"`
	StartTime       int    `gorm:"column:start_time;default:0;NOT NULL" json:"start_time"`
	EndTime         int    `gorm:"column:end_time;default:0;NOT NULL" json:"end_time"`
	IsOpen          int    `gorm:"column:is_open;default:0;NOT NULL" json:"is_open"`
	Status          int    `gorm:"column:status;default:0;NOT NULL" json:"status"`
	IsPush          int    `gorm:"column:is_push;default:0;NOT NULL" json:"is_push"`
	Remark          string `gorm:"column:remark;NOT NULL" json:"remark"`
	Type            int    `json:"type" gorm:"column:type;type:int(1);default:0" comment:"任务的类型,0为手动创建,1为模板创建(MA创建)"`
	Category        int32  `json:"category" gorm:"column:category;type:tinyint(1);not null;default:0;comment:任务标签 1 手动 2 模板"`
	PlanTemplateID  int64  `json:"plan_template_id" gorm:"column:plan_template_id;type:bigint(20);not null;default:0;comment:任务模板id"`
}

func (PlanModel) TableName() string {
	return "plans"
}

func (p PlanModel) AllTableName(enterpriseID int64) string {
	return fmt.Sprintf(" %s%d.%s ", p.DBPrefix(), enterpriseID, p.TableName())
}
