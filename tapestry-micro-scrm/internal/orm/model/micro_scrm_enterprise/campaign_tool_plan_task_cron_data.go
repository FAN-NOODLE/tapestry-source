package micro_scrm_enterprise

import "micro-scrm/internal/orm/model"

type CampaignToolPlanTaskCronDataModel struct {
	model.BranchModel
	MemberCode   string `gorm:"column:member_code;NOT NULL" json:"member_code"` // 会员卡号
	UnionID      string `gorm:"column:union_id;NOT NULL" json:"union_id"`
	GradeCode    string `gorm:"column:grade_code;NOT NULL" json:"grade_code"`
	AdminBA      string `gorm:"column:admin_ba;NOT NULL" json:"admin_ba"`
	TagName      string `gorm:"column:tag_name;NOT NULL" json:"tag_name"`             // 标签
	CronFileName string `gorm:"column:cron_file_name;NOT NULL" json:"cron_file_name"` // sftp文件名
	CreateTime   string `gorm:"column:create_time;NOT NULL" json:"create_time"`       // campaign-tool推送时间
	FileID       int64  `gorm:"column:file_id;NOT NULL" json:"file_id"`
}

func (CampaignToolPlanTaskCronDataModel) TableName() string {
	return "campaign_tool_plan_task_cron_data"
}
