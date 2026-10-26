package customerModels

import "micro-scrm/internal/app/models"

type SowTemplateModel struct {
	models.MicroBaseModel
	Name       string                   `json:"name" gorm:"column:name;type:varchar(255);not null;comment:播种模版名称"`
	Theme      string                   `json:"theme" gorm:"column:theme;type:varchar(255);not null;comment:主题说明"`
	Status     int                      `json:"status" gorm:"column:status;type:tinyint(1);not null;default:2;comment:是否上架  1 是 2 否"`
	AdminName  string                   `json:"admin_name" gorm:"column:admin_name;type:varchar(255);not null;comment:创建人名称"`
	UseNum     int                      `json:"use_num" gorm:"column:use_num;type:int(10);not null;default:0;comment:引用次数"`
	NoticeList []SowTemplateNoticeModel `json:"notice_list"  gorm:"foreignKey:TemplateID;references:ID"`
}

// TableName 播种任务
func (SowTemplateModel) TableName() string {
	return "sow_template"
}

const (
	SowTemplateIsUp   = 1 // 1 已上架
	SowTemplateIsDown = 2 // 2 未上架
)
