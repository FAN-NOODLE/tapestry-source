package microScrmModel

import (
	"micro-scrm/internal/app/models"
)

type NotificationInfoMergeContentModel struct {
	models.BaseModel
	GroupId        int64  `json:"group_id" gorm:"column:group_id;type:bigint(20);not null;default:0;comment:合并传讯详情中的每一个部门组的id"`
	MergeId        int64  `json:"merge_id" gorm:"column:merge_id;type:bigint(20);not null;default:0;comment:合并传讯id"`
	DepartmentName string `json:"department_name" gorm:"column:department_name;type:varchar(255);not null;comment:部门名"`
	InfoId         int64  `json:"info_id" gorm:"column:info_id;type:bigint(20);not null;default:0;comment:传讯id"`
}

// TableName returns the table name of the NotificationInfoMerge model
func (NotificationInfoMergeContentModel) TableName() string {
	return "notification_info_merge_content"
}
