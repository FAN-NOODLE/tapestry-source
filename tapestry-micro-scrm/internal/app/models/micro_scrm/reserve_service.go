package microScrmModel

import "micro-scrm/internal/app/models"

type ReserveServiceModel struct {
	models.BaseModel
	ReserveName   string `json:"reserve_name" gorm:"column:reserve_name;type:varchar(255);not null;comment:服务类型名称"`
	ReserveImg    string `json:"reserve_img" gorm:"column:reserve_img;type:varchar(255);not null;comment:服务类型的图片"`
	ReserveTime   int32  `json:"reserve_time" gorm:"column:reserve_time;type:tinyint(1);not null;default:1;comment:服务时长,单位小时"`
	ReserveStatus int32  `json:"reserve_status" gorm:"column:reserve_status;type:tinyint(1);not null;default:0;comment:服务类型状态，上下架，0为下架，1为上架"`
	EndTime       int    `json:"end_time" gorm:"column:end_time;type:int(11);not null;default:0;comment:指定时间的截止时间"`
	ServiceType   int32  `json:"service_type" gorm:"column:service_type;type:tinyint(1);not null;default:0;comment:活动有效期，1长期有效，0指定时间"`
	StartTime     int    `json:"start_time" gorm:"column:start_time;type:int(11);not null;default:0;comment:指定时间的开始时间"`
	ServiceRole   int32  `json:"service_role" gorm:"column:service_role;type:tinyint(1);not null;default:0;comment:使用范围，0为全部，1为选择范围"`
	TimeSlot      string `json:"-" gorm:"column:time_slot;type:text;not null;default:'';comment:时间段"`
}

// TableName returns the table name of the ReserveService model
func (ReserveServiceModel) TableName() string {
	return "reserve_service"
}
