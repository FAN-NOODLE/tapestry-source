package microScrmModel

import (
	"micro-scrm/internal/app/models"
)

const (
	NotificationWeek     = iota + 1 // 周传讯
	NotificationTransfer            // 转货传讯
	NotificationQuick               // 紧急传讯
)

const (
	NotificationStatusWait     = iota // 待提交
	NotificationStatusChecking        // 待审核
	NotificationStatusRecall          // 已撤回
	NotificationStatusFail            // 未通过
	NotificationStatusSuccess         // 已通过
	NotificationStatusCancel          // 作废
)

const (
	NotificationNoDraft = iota //是否保存草稿，0为否，1为是
)

type NotificationInfoModel struct {
	models.BaseModel
	ThemeTitle     string `json:"theme_title" gorm:"column:theme_title;type:varchar(255);not null;comment:主题（紧急传讯需要）"`
	Title          string `json:"title" gorm:"column:title;type:varchar(255);not null;comment:标题"`
	UntilName      string `json:"until_name" gorm:"column:until_name;type:varchar(255);not null;comment:至"`
	Content        string `json:"content" gorm:"column:content;type:text;not null;comment:内容"`
	Perform        string `json:"perform" gorm:"column:perform;type:varchar(255);not null;comment:执行事务"`
	DateTime       string `json:"date_time" gorm:"column:date_time;type:varchar(255);not null;comment:时间"`
	ContactInfo    string `json:"contact_info" gorm:"column:contact_info;type:text;not null;comment:联系人信息"`
	Attachment     string `json:"attachment" gorm:"column:attachment;type:text;not null;comment:附件信息"`
	Type           int32  `json:"type" gorm:"column:type;type:tinyint(1);not null;index:type_index;default:0;comment:传讯类型，1周传讯，2转货传讯,3紧急传讯"`
	DepartmentName string `json:"department_name" gorm:"column:department_name;type:varchar(255);not null;comment:所属部门，部门名称"`
	AdminID        int64  `json:"admin_id" gorm:"column:admin_id;type:bigint(20);not null;default:0;comment:管理员id"`
	Status         int32  `json:"status" gorm:"column:status;type:tinyint(1);not null;index:status_index;default:0;comment:传讯的审核状态"`
	MergeRelease   int32  `json:"merge_release" gorm:"column:merge_release;type:tinyint(1);not null;index:merge_release_index;default:0;comment:是否已经发布 0 否 1是"`
	ReviewTime     int    `json:"review_time" gorm:"column:review_time;type:int(10);not null;default:0;comment:审核时间"`
	SubmitTime     int    `json:"submit_time" gorm:"column:submit_time;type:int(10);not null;default:0;comment:提交时间"`
	ReviewAdminID  int64  `json:"review_admin_id" gorm:"column:review_admin_id;type:bigint(20);not null;comment:审核的管理人员"`
	ReviewMemberID int64  `json:"review_member_id" gorm:"column:review_member_id;type:bigint(20);not null;index:review_member_id_index;comment:审核的member_id"`
	ReviewName     string `json:"review_name" gorm:"column:review_name;type:varchar(255);not null;comment:审核人员名称"`
	ReadMemberIds  string `json:"read_member_ids" gorm:"column:read_member_ids;type:varchar(255);not null;comment:状态流转时已查看的人（列表红点用)"`
	MemberID       int64  `json:"member_id" gorm:"column:member_id;type:bigint(20);not null;index:member_id_index;default:0;comment:导购的member_id"`
	MemberName     string `json:"member_name" gorm:"column:member_name;type:varchar(255);not null;comment:导购名称"`
	Comment        string `json:"comment" gorm:"column:comment;type:varchar(255);not null;comment:留言"`
	PageURL        string `json:"page_url" gorm:"column:page_url;type:varchar(255);not null;comment:页脚设置"`
	IsDraft        int32  `json:"is_draft" gorm:"column:is_draft;type:tinyint(1);not null;default:0;comment:是否保存草稿，0为否，1为是"`
	PdfUrl         string `json:"pdf_url" gorm:"column:pdf_url;type:varchar(255);not null;comment:pdf文件cos地址"`
	ZipUrl         string `json:"zip_url" gorm:"column:zip_url;type:varchar(255);not null;comment:打包的zip url"`
}

// TableName returns the table name of the NotificationInfo model
func (NotificationInfoModel) TableName() string {
	return "notification_info"
}

func (NotificationInfoModel) NotificationTypeStr(index int32) string {
	typeMap := map[int32]string{
		NotificationWeek:     "周传讯",
		NotificationTransfer: "转货传讯",
		NotificationQuick:    "紧急传讯",
	}
	return typeMap[index]
}

func (NotificationInfoModel) NotificationStatusStr(index int32) string {
	statusMap := map[int32]string{
		NotificationStatusWait:     "待提交",
		NotificationStatusChecking: "待审核",
		NotificationStatusRecall:   "已撤回",
		NotificationStatusFail:     "未通过",
		NotificationStatusSuccess:  "已通过",
		NotificationStatusCancel:   "已作废",
	}
	return statusMap[index]
}
