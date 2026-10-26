package microScrmModel

import "micro-scrm/internal/app/models"

const (
	PlanMaTypeNo = iota + 1
	PlanMaTypeYes
	PlanMaTypeNewYes
)
const (
	PlanNotStart = iota // 任务的状态，0未开始,1进行中,2已结束
	PlanIng
	PlanEnd
)

const (
	FINISHPUSHMA = iota
	FINISHPUSHMAING
	FINISHPUSHMAEND
)

type PlansModel struct {
	models.BaseModel
	Name            string `json:"name" gorm:"column:name;type:varchar(64);not null;comment:计划名称"`
	TaskNum         int    `json:"task_num" gorm:"column:task_num;type:int(10);not null;default:0;comment:该计划下关联任务的数量"`
	OperationUserID int64  `json:"operation_user_id" gorm:"column:operation_user_id;type:bigint(11);not null;default:0;comment:成员表id"`
	StartTime       int    `json:"start_time" gorm:"column:start_time;type:int(10);not null;default:0"`
	EndTime         int    `json:"end_time" gorm:"column:end_time;type:int(10);not null;default:0"`
	IsOpen          int32  `json:"is_open" gorm:"column:is_open;type:tinyint(2);not null;default:0;comment:任务在小程序中是否可见。0不可见，1可见；默认可见"`
	Status          int32  `json:"status" gorm:"column:status;type:tinyint(2);not null;default:0;comment:任务状态；0未开始、1进行中、2已结束"`
	Remark          string `json:"remark" gorm:"column:remark;type:varchar(200);not null;comment:计划描述"`
	IsPush          int    `json:"is_push" gorm:"column:is_push;type:int(1);default:0;comment:任务是否推送 1已推送 0 未推送"`
	Type            int    `json:"type" gorm:"column:type;type:int(1);default:0"`
	Category        int    `json:"category" gorm:"column:category;type:tinyint(1);default:0;comment:任务的类型,1为手动创建,2为模板创建(MA创建)"`
	PlanTemplateID  int64  `json:"plan_template_id" gorm:"column:plan_template_id;type:bigint(20);default:0;comment:任务模板id"`
}

// TableName returns the table name of the Plans model
func (PlansModel) TableName() string {
	return "plans"
}

type PlanTopList struct {
	MemberUserId string `json:"member_user_id"`
	Total        int64  `json:"total"`
	Ranking      string `json:"ranking"`
	RankingValue string `json:"ranking_value"`
}
