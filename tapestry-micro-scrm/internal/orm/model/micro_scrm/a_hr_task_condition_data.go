package micro_scrm

type ARrTaskConditionDataModel struct {
	ID                 int64  `json:"id" gorm:"primaryKey;column:id;type:bigint(20) auto_increment"`
	TaskAllotLogID     int    `json:"task_allot_log_id" gorm:"column:task_allot_log_id;type:int(10);not null"`
	TagName            string `json:"tag_name" gorm:"column:tag_name;type:varchar(64);not null;comment:CampaignTool标签名称"`                        // CampaignTool标签名称
	WwTagID            string `json:"ww_tag_id" gorm:"column:ww_tag_id;type:varchar(255);not null;comment:企微标签id"`                               // 企微标签id
	UnionID            string `json:"union_id" gorm:"column:union_id;type:char(64);not null;comment:顾客的UnionID"`                                 // 顾客的UnionID
	WorkExternalUserID string `json:"work_external_user_id" gorm:"column:work_external_user_id;type:char(64);not null;comment:外部联系人user_id,可能空"` // 外部联系人user_id,可能空
	IsUnion            int32  `json:"is_union" gorm:"column:is_union;type:tinyint(2);not null;default:0;comment:关系状态；0没有好友关系、1有好友关系"`            // 关系状态；0没有好友关系、1有好友关系
	Mobile             string `json:"mobile" gorm:"column:mobile;type:varchar(20);comment:手机"`                                                   // 手机
	MarsID             string `json:"mars_id" gorm:"column:mars_id;type:varchar(20);not null;comment:BA的MarsCode"`                               // BA的MarsCode
	UserID             string `json:"user_id" gorm:"column:user_id;type:varchar(64);not null;comment:员工的企微user_id"`                              // 员工的企微user_id
	ConditionID        int64  `json:"condition_id" gorm:"column:condition_id;type:bigint(20);not null;default:0;comment:条件对应的子任务的ID"`            // 条件对应的子任务的ID
	TaskID             int64  `json:"task_id" gorm:"column:task_id;type:bigint(20);not null;default:0;comment:条件对应的子任务的ID"`                      // 条件对应的子任务的ID
	PlanID             int64  `json:"plan_id" gorm:"column:plan_id;type:bigint(20);not null;default:0;comment:条件对应的子任务的ID"`                      // 条件对应的子任务的ID
	CreatedAt          int    `json:"created_at" gorm:"column:created_at;type:int(10);not null;default:0"`
	UpdatedAt          int    `json:"updated_at" gorm:"column:updated_at;type:int(10);not null;default:0"`
	DeletedAt          int    `json:"deleted_at" gorm:"column:deleted_at;type:int(10);not null;default:0"`
	IsDelete           uint32 `json:"is_delete" gorm:"column:is_delete;type:tinyint(1) unsigned;not null;default:0"`
}

// TableName returns the table name of the ARrTaskConditionDataModel model
func (ARrTaskConditionDataModel) TableName() string {
	return "a_hr_task_condition_data"
}
