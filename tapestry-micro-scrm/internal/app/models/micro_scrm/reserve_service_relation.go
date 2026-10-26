package microScrmModel

import (
	"micro-scrm/internal/app/models"
)

type ReserveServiceRelationModel struct {
	models.BaseModel
	ServiceID    int64  `json:"service_id" gorm:"column:service_id;type:bigint(20);not null;index:service_id_index;default:0;comment:预约服务的id"`
	StoreID      int    `json:"store_id" gorm:"column:store_id;type:int(10);not null;default:0;comment:部门id，departments表中的third_id，为企微的部门id"`
	StoreName    string `json:"store_name" gorm:"column:store_name;type:varchar(255);not null;comment:部门名称"`
	MemberID     int64  `json:"member_id" gorm:"column:member_id;type:bigint(20);not null;default:0;comment:导购的id，members表中的id"`
	MemberName   string `json:"member_name" gorm:"column:member_name;type:varchar(255);not null;comment:导购名称"`
	IsStore      int32  `json:"is_store" gorm:"column:is_store;type:tinyint(1);not null;default:0;comment:是否是门店,0为否,1为是"`
	StoreType    int32  `json:"store_type" gorm:"column:store_type;type:tinyint(1);not null;default:0;comment:门店类型,0为默认,1为零售,2为奥莱"`
	DepartmentId int64  `json:"department_id" gorm:"column:department_id;type:bigint(20);not null;default:0;comment:部门id,departments表中的id"`
	UserId       string `json:"user_id" gorm:"column:user_id;type:varchar(100);not null;comment:导购的user_id"`
}

// TableName returns the table name of the ReserveServiceRelation model
func (ReserveServiceRelationModel) TableName() string {
	return "reserve_service_relation"
}
