package microScrmModel

import "micro-scrm/internal/app/models"

type ReserveOrderDataAcxiomSftpLogModel struct {
	models.MicroBaseModel
	FileType  int32  `json:"file_type" gorm:"column:file_type;type:tinyint(1);not null;default:0;comment:文件的行为类型,0为推送文件给acxiom，1为去acxiom拉取数据"`
	LocalFile string `json:"local_file" gorm:"column:local_file;type:varchar(255);not null;comment:本地的文件名称"`
	FileURL   string `json:"file_url" gorm:"column:file_url;type:varchar(255);not null;comment:文件的cos地址"`
	Status    int32  `json:"status" gorm:"column:status;type:tinyint(1);not null;default:0;comment:文件的状态,0为未完成，1为已完成"`
	FileTime  int    `json:"file_time" gorm:"column:file_time;type:int(10);not null;default:0;comment:文件的时间，每天的零点时间"`
	PushType  string `json:"push_type" gorm:"column:push_type;type:varchar(255);not null;default:push_day_order;comment:推送类型,push_day_order,successful_appointment,STORE_APPOINTMENTQTY,Report_Data"`
	Year      int64  `json:"year" gorm:"column:year;type:int(10);not null;default:0;comment:财年"`
	Period    int64  `json:"period" gorm:"column:period;type:int(10);not null;default:0;comment:财月"`
}

// TableName returns the table name of the ReserveOrderDataAcxiomSftpLog model
func (ReserveOrderDataAcxiomSftpLogModel) TableName() string {
	return "reserve_order_data_acxiom_sftp_log"
}
