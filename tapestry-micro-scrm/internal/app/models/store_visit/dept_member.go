package storeVisitModel

import (
	"micro-scrm/internal/app/models"
	membersModel "micro-scrm/internal/app/models/members"
	"time"
)

type StoreVisitDeptMemberModel struct {
	models.BaseModel
	MemberID            int64                    `json:"member_id" gorm:"column:member_id;type:bigint(20);not null;default:0;index:member_id_index;comment:成员ID"`
	Name                string                   `json:"name" gorm:"column:name;type:varchar(32);not null;default:'';comment:姓名"`
	Userid              string                   `json:"user_id" gorm:"column:user_id;type:varchar(32);not null;default:'';comment:企微id"`
	Email               string                   `json:"email" gorm:"column:email;type:varchar(32);not null;default:'';comment:邮箱"`
	DeptName            string                   `json:"dept_name" gorm:"column:dept_name;type:varchar(32);not null;default:'';comment:部门名"`
	IsStoreVisitManager int                      `json:"is_store_visit_manager" gorm:"column:is_store_visit_manager;type:tinyint(1);not null;default:0;comment:是否是巡店管理员 0否 1是"`
	UpdatedTime         time.Time                `json:"updated_time" gorm:"column:updated_time;type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:是否是巡店管理员 0否 1是"`
	Members             membersModel.MemberModel `gorm:"foreignKey:ID;references:MemberID"`
}

func (StoreVisitDeptMemberModel) TableName() string {
	return "store_visit_dept_member"
}
