package welcome

import "micro-scrm/internal/orm/model"

// 新老用户
const (
	CodeCustomTypeNewCustomer = "new_customer"
	CodeCustomTypeOldCustomer = "old_customer"
)

// 领券状态
const (
	ReceiveCouponStatusYes = "received"
	ReceiveCouponStatusNo  = "unreceived"
)

type WelcomesExternalTypeDataModel struct {
	model.MicroScrmBaseModel

	WelcomesExternalTypeID int64  `json:"welcomes_external_type_id" gorm:"column:welcomes_external_type_id;type:bigint(20);not null;default:0;comment:外部数据类型ID"` // 外部数据类型ID
	Name                   string `json:"name" gorm:"column:name;type:varchar(20);not null;comment:外部数据类型名称"`                                                    // 外部数据类型名称
	Code                   string `json:"code" gorm:"column:code;type:varchar(20);not null;comment:数据code"`                                                      // 数据code
	Sort                   int    `json:"sort" gorm:"column:sort;type:int(11);not null;default:0;comment:类型数据排序"`                                                // 类型数据排序
	Remark                 string `json:"remark" gorm:"column:remark;type:varchar(50);not null;comment:数据说明"`                                                    // 数据说明

	TypeInfo WelcomesExternalTypeModel `json:"type_info" gorm:"foreignKey:WelcomesExternalTypeID;references:ID"`
}

func (WelcomesExternalTypeDataModel) TableName() string {
	return "welcomes_external_type_data"
}
