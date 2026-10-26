package welcome

import "micro-scrm/internal/orm/model"

const (
	MemberLevel         = "member_level"          // 会员等级
	MemberLevelFY26     = "member_level_26"       // FY26会员等级
	StoreType           = "store_type"            // 门店类型
	CustomType          = "custom_type"           // 顾客类型
	ReceiveCouponStatus = "receive_coupon_status" // 领券状态
	ReserveStatus       = "reserve_status"        // 预约状态
)

type WelcomesExternalTypeModel struct {
	model.MicroScrmBaseModel

	Name   string `json:"name" gorm:"column:name;type:varchar(50);not null;comment:外部数据类型名称"` // 外部数据类型名称
	Code   string `json:"code" gorm:"column:code;type:varchar(20);not null;comment:外部数据类型标志"` // 外部数据类型标志
	Remark string `json:"remark" gorm:"column:remark;type:varchar(50);not null;comment:数据说明"` // 数据说明

	TypeDataList []WelcomesExternalTypeDataModel `json:"type_data_list" gorm:"foreignKey:WelcomesExternalTypeID;references:ID"`
}

func (WelcomesExternalTypeModel) TableName() string {
	return "welcomes_external_type"
}
