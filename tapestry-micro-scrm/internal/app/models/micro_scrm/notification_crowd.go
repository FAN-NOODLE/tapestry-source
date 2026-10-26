package microScrmModel

import (
	"micro-scrm/internal/app/models"
)

type NotificationCrowdModel struct {
	models.BaseModel
	Name      string `json:"name" gorm:"column:name;type:varchar(255);not null;comment:人群包名称"`
	FileUrl   string `json:"file_url" gorm:"column:file_url;type:varchar(255);not null;comment:人群包文件"`
	Type      int    `json:"type" gorm:"column:type;type:tinyint(1);not null;default:0;comment:人群包类型"`
	MemberNum int    `json:"member_num" gorm:"column:member_num;type:int(10);not null;default:0;comment:人群包人数"`
	AdminID   int64  `json:"admin_id" gorm:"column:admin_id;type:bigint(20);not null;default:0;comment:人群包创建人"`
	AdminName string `json:"admin_name" gorm:"column:admin_name;type:varchar(255);not null;comment:创建人名称"`
	//关联导购列表
	DataList []NotificationCrowdDataModel `json:"data_list" gorm:"foreignKey:CrowdId;references:ID"`
}

// TableName returns the table name of the NotificationCorwd model
func (NotificationCrowdModel) TableName() string {
	return "notification_crowd"
}
