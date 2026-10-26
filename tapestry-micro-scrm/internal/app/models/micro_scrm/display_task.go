package microScrmModel

import (
	"micro-scrm/internal/app/models"
)

const (
	DisplayStatus = iota // 任务状态
	DisplayStatusIng
	DisplayStatusEnd
)

const (
	ShowStatusDown = iota // 上下架
	ShowStatusUp
)

const (
	VMStoreTypeAll = iota // 指定门店类型
	VMStoreTypeLimit
)

const (
	VMTaskReviewTypeNormal = iota
	VMTaskReviewTypeRand
	VMTaskReviewTypeFile
)

type DisplayTaskModel struct {
	models.BaseModel
	DisplayName   string                  `json:"display_name" gorm:"column:display_name;type:varchar(255);not null;comment:VM陈列任务标题"`
	StoreType     int32                   `json:"store_type" gorm:"column:store_type;type:tinyint(1);not null;default:0;comment:通知门店类型，0为全部门店,1为指定门店"`
	Content       string                  `json:"content" gorm:"column:content;type:text;not null;comment:陈列任务的内容"`
	ClassURL      string                  `json:"class_url" gorm:"column:class_url;type:text;not null;comment:关联课程的链接"`
	TemplateInfo  string                  `json:"template_info" gorm:"column:template_info;type:text;not null;comment:上传的模板字段"`
	ShowType      int32                   `json:"show_type" gorm:"column:show_type;type:tinyint(1);not null;default:0;comment:默认上下架状态，1上架，0下架"`
	Status        int32                   `json:"status" gorm:"column:status;type:tinyint(1);not null;default:0;comment:任务的状态，0为未开始，1为进行中，2为已结束"`
	StartTime     int                     `json:"start_time" gorm:"column:start_time;type:int(10);not null;default:0;comment:任务的开始时间"`
	EndTime       int                     `json:"end_time" gorm:"column:end_time;type:int(10);not null;default:0;comment:任务的结束时间"`
	ReviewType    int                     `json:"review_type" gorm:"column:review_type;type:tinyint(1);not null;default:0;comment:审批类型 0 默认VM直接审批 1 随机分配SL审批 2名单上传SL审批"`
	ReviewFile    string                  `json:"review_file" gorm:"column:review_file;type:varchar(255);not null;comment:定向分配名单文件"`
	TaskAllotInfo []DisplayTaskAllotModel `json:"task_allot_info" gorm:"foreignKey:DisplayTaskID;references:ID"`
}

// TableName returns the table name of the DisplayTask model
func (DisplayTaskModel) TableName() string {
	return "display_task"
}
