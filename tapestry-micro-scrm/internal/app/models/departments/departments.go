package departmentsModels

import "micro-scrm/internal/app/models"

type DepartmentsModel struct {
	models.MicroBaseModel

	ParentID         int64  `json:"parent_id" gorm:"column:parent_id;type:bigint(20);not null;index:parent_id_index;default:0;comment:业务中父部门id"`
	ThirdID          int    `json:"third_id" gorm:"primaryKey;column:third_id;type:int(10);default:0;comment:企业微信中的组织ID"`
	ThirdParentID    int    `json:"third_parent_id" gorm:"column:third_parent_id;type:int(10);not null;default:0;comment:企业微信中的组织父ID"`
	Name             string `json:"name" gorm:"column:name;type:varchar(128);not null;comment:部门名称"`
	NameEn           string `json:"name_en" gorm:"column:name_en;type:varchar(128);not null;comment:英文名称"`
	Order            int    `json:"order" gorm:"column:order;type:int(10);not null;default:0;comment:在父部门中的次序值"`
	ChildDepartments string `json:"child_departments" gorm:"column:child_departments;type:text;comment:子部门、子子部门"`
}

// TableName 企业微信组织结构
func (DepartmentsModel) TableName() string {
	return "departments"
}

type DepartmentsMembers struct {
	DepartmentId string `json:"department_id"`
	Userid       string `json:"userid"`
	MemberId     int64  `json:"member_id"`
	MemberName   string `json:"member_name"`
	Position     string `json:"position"`
}

type DepartmentsMembersInfo struct {
	DepartmentId   string `json:"department_id"`
	MemberId       int64  `json:"member_id"`
	DepartmentName string `json:"department_name"`
	ThirdId        int    `json:"third_id"`
}

type DepartmentsMembersData struct {
	MemberId   int64  `json:"member_id"`
	MemberName string `json:"member_name"`
	Mobile     string `json:"mobile"`
	Userid     string `json:"userid"`
	ThirdId    int    `json:"third_id"`
	ThirdName  string `json:"third_name"`
}
