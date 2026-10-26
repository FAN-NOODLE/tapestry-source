package dashboardModel

import "micro-scrm/internal/app/models"

type DashboardManagerModel struct {
	models.BaseModel
	MemberId      int64  `json:"member_id" gorm:"column:member_id;type:bigint(20);not null;default:0;comment:管理者member_id"`
	AreaIds       string `json:"area_ids" gorm:"column:area_ids;type:varchar(250);not null;comment:权限范围area_id集合，英文逗号分割"`
	Position      string `json:"position" gorm:"column:position;type:varchar(100);not null;comment:职位名称"`
	PositionLevel int64  `json:"position_level" gorm:"column:position_level;type:tinyint(1);not null;default:0;comment:职位级别 越小越高"`
}

// TableName returns the table name of the DashboardArea model
func (DashboardManagerModel) TableName() string {
	return "dashboard_manager"
}
