package dashboardModel

import "micro-scrm/internal/app/models"

type DashboardAreaModel struct {
	models.BaseModel
	Type       int    `json:"type" gorm:"column:type;type:int(10);not null;default:0;comment:1 渠道 2大区 3分区 4门店"`
	Name       string `json:"name" gorm:"column:name;type:varchar(100);not null;comment:名称"`
	ParentID   int    `json:"parent_id" gorm:"column:parent_id;type:int(10);not null;default:0;comment:父级id"`
	Path       string `json:"path" gorm:"column:path;type:varchar(255);not null;comment:层级路径，英文逗号分割"`
	ThirdID    int    `json:"third_id" gorm:"column:third_id;type:int(10);not null;default:0;comment:企微部门id"`
	IsVirtual  int    `json:"is_virtual" gorm:"column:is_virtual;type:int(10);not null;default:0;comment:是否虚拟单元"`
	LMSID      string `json:"lms_id" gorm:"column:lms_id;type:varchar(128);not null;default:'';comment:LMS区域或门店唯一编码"`
	FiscalYear string `json:"fiscal_year" gorm:"column:fiscal_year;type:varchar(255);not null;comment:渠道财年标识"`
	ClosedAt   int64  `json:"closed_at" gorm:"column:closed_at;type:int(10);not null;default:0;comment:闭店时间" `
	Sort       int64  `json:"sort" gorm:"column:sort;type:int(10);not null;default:0;comment:同级别中的排序，默认升序" `
	//获取子区域
	Children []DashboardAreaModel `json:"children"  gorm:"foreignKey:ParentID"`
}

// TableName returns the table name of the DashboardArea model
func (DashboardAreaModel) TableName() string {
	return "dashboard_area"
}

func (DashboardAreaModel) DBPrefix() string {
	return "micro_scrm_enterprise_"
}

const (
	VirtualArea = 1 //虚拟区域单元
)

const (
	AreaSource = iota + 1 //渠道
	AreaRegion            //大区
	AreaArea              //分区
	AreaStore             //门店
)
