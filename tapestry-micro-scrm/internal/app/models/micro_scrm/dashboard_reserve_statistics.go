package microScrmModel

import (
	"micro-scrm/internal/app/models"
)

type DashboardReserveStatisticsModel struct {
	ID                uint             `json:"id" gorm:"primaryKey;column:id;type:int(10) unsigned auto_increment;comment:自增主键"`
	MemberID          uint64           `json:"member_id" gorm:"column:member_id;type:bigint(20) unsigned;not null;index:member_id_index;default:0;comment:员工id"`
	UserID            string           `json:"user_id" gorm:"column:user_id;type:varchar(50);not null;comment:导购id"`
	SourceID          uint             `json:"source_id" gorm:"column:source_id;type:int(10) unsigned;not null;default:0;comment:渠道id"`
	RegionID          uint             `json:"region_id" gorm:"column:region_id;type:int(10) unsigned;not null;default:0;comment:大区id"`
	AreaID            uint             `json:"area_id" gorm:"column:area_id;type:int(10) unsigned;not null;default:0;comment:区域id"`
	ThirdID           uint             `json:"third_id" gorm:"column:third_id;type:int(10) unsigned;not null;default:0;comment:企微部门id"`
	ServiceID         uint             `json:"service_id" gorm:"column:service_id;type:int(20) unsigned;not null;default:0;comment:邀约服务id"`
	SendNumber        uint32           `json:"send_number" gorm:"column:send_number;type:smallint(5) unsigned;not null;default:0;comment:发送数量"`
	StoreNumber       uint32           `json:"store_number" gorm:"column:store_number;type:smallint(5) unsigned;not null;default:0;comment:到店数量"`
	SuccessNumber     uint32           `json:"success_number" gorm:"column:success_number;type:smallint(5) unsigned;not null;default:0;comment:成功数量"`
	OnlineNumber      uint32           `json:"online_number" gorm:"column:online_number;type:smallint(5) unsigned;not null;default:0;comment:线上数量"`
	OfflineNumber     uint32           `json:"offline_number" gorm:"column:offline_number;type:smallint(5) unsigned;not null;default:0;comment:线下数量"`
	SuccessMoney      uint             `json:"success_money" gorm:"column:success_money;type:int(10) unsigned;not null;default:0;comment:成功金额（单位:分）"`
	OnlineProportion  uint32           `json:"online_proportion" gorm:"column:online_proportion;type:smallint(5) unsigned;not null;default:0;comment:线上比例"`
	OfflineProportion uint32           `json:"offline_proportion" gorm:"column:offline_proportion;type:smallint(5) unsigned;not null;default:0;comment:线下比例"`
	CreatedAt         uint             `json:"created_at" gorm:"column:created_at;type:int(11) unsigned;not null;default:0;comment:创建时间"`
	UpdatedAt         uint             `json:"updated_at" gorm:"column:updated_at;type:int(11) unsigned;not null;default:0;comment:更新时间"`
	DeletedAt         models.DeletedAt `gorm:"softDelete;column:deleted_at;type:int(10) unsigned;not null;default:0" json:"deleted_at"`
}

// TableName dashboard预约统计
func (DashboardReserveStatisticsModel) TableName() string {
	return "dashboard_reserve_statistics"
}
