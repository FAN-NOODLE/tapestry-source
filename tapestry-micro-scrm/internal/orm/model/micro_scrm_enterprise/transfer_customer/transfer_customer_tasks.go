package transfer_customer

import "micro-scrm/internal/orm/model"

const (
	TCT_STATUS_ZERO = iota
	TCT_STATUS_NOT
	TCT_STATUS_RUNNING
	TCT_STATUS_FINISH
	TCT_STATUS_CANCEL
)

// 任务数据拉取的状态
const (
	TCT_PullDataStatus_No  = 0 // 未拉取数据
	TCT_PullDataStatus_Yes = 1 // 已经拉取数据
)

const (
	TCT_Type_Only   = 1 // 统一到一个管理BA上
	TCT_Type_Random = 2 // 所以选择继承到管理BA上
)

type TransferCustomerTasksModel struct {
	model.BranchModel

	TransferID     string `json:"transfer_id" gorm:"column:transfer_id;type:char(16);not null;default:0;comment:任务自定义ID"`     // 任务自定义ID
	Name           string `json:"name" gorm:"column:name;type:varchar(20);not null;default:0;comment:任务名称"`                   // 任务名称
	CustomerNum    int    `json:"customer_num" gorm:"column:customer_num;type:int(10);not null;default:0;comment:符合条件的客户数"`   // 符合条件的客户数
	StartTime      int64  `json:"start_time" gorm:"column:start_time;type:int(10);not null;default:0;comment:任务开始时间"`         // 任务开始时间
	EndTime        int64  `json:"end_time" gorm:"column:end_time;type:int(10);not null;default:0;comment:任务完成时间"`             // 任务完成时间
	Status         uint8  `json:"status" gorm:"column:status;type:tinyint(2);not null;default:0;comment:任务状态；1未开始，1进行中，2已完成"` // 任务状态；1未开始，1进行中，2已完成
	Type           uint8  `json:"type" gorm:"column:type;type:tinyint(2);not null;default:1;comment:1唯一继承，2随机继承"`             // 1唯一继承，2随机继承
	PullDataStatus uint8  `json:"pull_data_status" gorm:"column:pull_data_status;type:tinyint(2);not null;default:0;comment:任务数据拉取状态：0 未拉取完成，1 拉取完成"`
}

func (TransferCustomerTasksModel) TableName() string {
	return "transfer_customer_tasks"
}
