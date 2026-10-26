package thirdModel

import (
	"micro-scrm/internal/app/models"
)

// 0 未绑定 1 已绑定 2 已解绑
const (
	BindStatusNotBind = iota
	BindStatusBind
	BindStatusLeave
)

type ReserveRecordModel struct {
	models.MicroBaseModel

	StoreCode       string `json:"store_code" gorm:"column:store_code;type:varchar(255);not null;comment:门店编码"`
	DepartmentID    int64  `json:"department_id" gorm:"column:department_id;type:bigint(20);not null;default:0;comment:部门id"`
	ReserveNumberID int64  `json:"reserve_number_id" gorm:"column:reserve_number_id;type:bigint(20);not null;uniqueIndex:unique_reserve_number_id_index;default:0;comment:预约单号"`
	ServiceName     string `json:"service_name" gorm:"column:service_name;type:varchar(50);not null;comment:服务名称"`
	Unionid         string `json:"unionid" gorm:"column:unionid;type:varchar(255);not null;index:unionid_index;comment:用户唯一id"`
	MemberID        int64  `json:"member_id" gorm:"column:member_id;type:bigint(20);not null;default:0;comment:绑定员工id"`
	ConfigID        string `json:"config_id" gorm:"column:config_id;type:varchar(255);not null;comment:新增联系方式的配置id"`
	ContactQrcode   string `json:"contact_qrcode" gorm:"column:contact_qrcode;type:varchar(255);not null;comment:联系我二维码链接"`
	BindStatus      int32  `json:"bind_status" gorm:"column:bind_status;type:tinyint(1);not null;default:0;comment:绑定状态 0 未绑定 1 已绑定 2 已解绑"`
}

// TableName 三方预约记录分配记录表
func (ReserveRecordModel) TableName() string {
	return "third_reserve_record"
}
