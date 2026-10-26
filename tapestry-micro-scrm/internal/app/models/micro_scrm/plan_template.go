package microScrmModel

import "micro-scrm/internal/app/models"

const (
	OpenBackYes = iota + 1 // 开启备份模板
)

type PlanTaskTemplateModel struct {
	models.BaseModel
	TemplateName        string `json:"template_name" gorm:"column:template_name;type:varchar(255);not null;comment:任务模板名称"`
	TemplateDesc        string `json:"template_desc" gorm:"column:template_desc;type:varchar(255);not null;comment:任务模板的描述"`
	TemplateContent     string `json:"template_content" gorm:"column:template_content;type:tinytext;comment:任务模板的内容"`
	TemplateContentType int32  `json:"template_content_type" gorm:"column:template_content_type;type:tinyint;comment:消息类型；1文本；2图片；3语音；4视频；5位置；6链接,7 小程序"`
	Status              int32  `json:"status" gorm:"column:status;type:tinyint(2);not null;default:0;comment:任务模板的状态，0禁用,1启用"`
	AdminID             int64  `json:"admin_id" gorm:"column:admin_id;type:bigint(20);not null;default:0;comment:后台创建的admin_id"`
	GroupSendType       int    `json:"group_send_type" gorm:"column:group_send_type;type:tinyint(2);not null;default:0;comment:任务群发类型，标签任务群发/企业任务群发"`
	FinishRate          int    `json:"finish_rate" gorm:"column:finish_rate;type:tinyint(1);not null;default:100;comment:任务完成率"`
	ALLowSelect         int    `json:"allow_select" gorm:"column:allow_select;type:tinyint(1);not null;default:0;comment:是否允许重新选择发送客户，0否，1是"`
	OpenBackUp          int    `json:"open_backup" gorm:"column:open_backup;type:tinyint(1);not null;default:0;comment:是否存在备份模版 0 否 1是"`
	TaskCategory        int    `json:"task_category" gorm:"column:task_category;type:tinyint(1);not null;default:0;comment:模版类型 0 沟通任务 1朋友圈任务"`
	AdminName           string `json:"admin_name" gorm:"-"`
}

// TableName returns the table name of the PlanTaskTemplate model
func (PlanTaskTemplateModel) TableName() string {
	return "plan_task_template"
}

// MaTemplateLists MA任务三方接口获取的任务模板列表
type MaTemplateLists struct {
	Id            int64  `json:"template_id"`
	TemplateName  string `json:"template_name"`
	TemplateDesc  string `json:"template_desc"`
	GroupSendType int    `json:"group_send_type"`
	TaskCategory  int    `json:"task_category"`
}

const TaskCategoryMoment = 1 //朋友圈

const GroupSendTypeGroupMsg = 1 //群发助手
