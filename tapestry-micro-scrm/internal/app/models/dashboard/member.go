package dashboardModel

import "micro-scrm/internal/app/models"

type DashboardMemberModel struct {
	models.BaseModel
	AuthType     int    `json:"auth_type" gorm:"column:auth_type;type:int(10);not null;default:0;comment:1 渠道经理 2大区经理 3分区经理 4门店经理"`
	MemberID     int64  `json:"member_id" gorm:"column:member_id;type:bigint(20);not null;index:member_id_index;default:0;comment:员工id"`
	SourceID     int    `json:"source_id" gorm:"column:source_id;type:int(10);not null;default:0;comment:渠道id"`
	RegionID     int    `json:"region_id" gorm:"column:region_id;type:int(10);not null;default:0;comment:大区id"`
	AreaID       int    `json:"area_id" gorm:"column:area_id;type:int(10);not null;default:0;comment:区域id"`
	ThirdID      int    `json:"third_id" gorm:"column:third_id;type:int(10);not null;default:0;comment:企微部门id"`
	AuthThirdIds string `json:"auth_third_ids" gorm:"column:auth_third_ids;type:text;comment:门店管理范围，英文逗号分割"`
}

// TableName returns the table name of the DashboardMember model
func (DashboardMemberModel) TableName() string {
	return "dashboard_member"
}

const (
	AuthSource = iota + 1
	AuthRegion
	AuthArea
	AuthStore
)
