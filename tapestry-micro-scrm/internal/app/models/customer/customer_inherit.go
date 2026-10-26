package customerModels

import "micro-scrm/internal/app/models"

const (
	InheritStatusUnassigned = iota // 0 未分配
	InheritStatusAssigned          // 1 已分配
)

type CustomerInheritModel struct {
	models.MicroBaseModel
	HandoverUserid         string `json:"handover_userid" gorm:"column:handover_userid;type:varchar(255);not null;index:handover_userid_index;comment:离职成员的userid"`
	HandoverName           string `json:"handover_name" gorm:"column:handover_name;type:varchar(255);not null;comment:离职成员姓名"`
	HandoverMobile         string `json:"handover_mobile" gorm:"column:handover_mobile;type:varchar(15);not null;comment:离职成员企业微信手机号"`
	HandoverDepartmentId   string `json:"handover_department_id" gorm:"column:handover_department_id;type:int(10);not null;default:0;comment:离职成员部门id"`
	HandoverDepartmentName string `json:"handover_department_name" gorm:"column:handover_department_name;type:varchar(255);not null;comment:离职成员部门名"`
	ExternalUserid         string `json:"external_userid" gorm:"column:external_userid;type:varchar(255);not null;index:external_userid_index;comment:外部联系人userid(openid)"`
	ExternalNickname       string `json:"external_nickname" gorm:"column:external_nickname;type:varchar(255);not null;comment:外部联系人昵称"`
	ExternalAvatar         string `json:"external_avatar" gorm:"column:external_avatar;type:varchar(255);not null;comment:外部联系人头像"`
	TakeoverUserid         string `json:"takeover_userid" gorm:"column:takeover_userid;type:varchar(255);not null;comment:接替人userid"`
	TakeoverName           string `json:"takeover_name" gorm:"column:takeover_name;type:varchar(255);not null;comment:接替人姓名"`
	TakeoverDepartmentName string `json:"takeover_department_name" gorm:"column:takeover_department_name;type:varchar(255);not null;comment:接替人部门名"`
	DimissionTime          int    `json:"dimission_time" gorm:"column:dimission_time;type:int(10);not null;default:0;comment:离职时间"`
	Status                 int32  `json:"status" gorm:"column:status;type:tinyint(1);not null;default:0;comment:状态: 0 未分配 1已分配"`
}

// ResponseCustomerInheritList 返回结构
type ResponseCustomerInheritList struct {
	ID                     int64  `json:"id"`
	HandoverUserid         string `json:"handover_userid"`
	HandoverName           string `json:"handover_name"`
	HandoverMobile         string `json:"handover_mobile"`
	HandoverDepartmentId   int64  `json:"handover_department_id"`
	HandoverDepartmentName string `json:"handover_department_name"`
	DimissionTime          int    `json:"dimission_time"`
	UnassignedUserNum      int    `json:"unassigned_user_num"`
	AssignedUserNum        int    `json:"assigned_user_num"`
	UnassignedGroupNum     int    `json:"unassigned_group_num"`
	AssignedGroupNum       int    `json:"assigned_group_num"`
}

// TableName returns the table name of the CustomerInherit model
func (CustomerInheritModel) TableName() string {
	return "customer_inherit"
}
