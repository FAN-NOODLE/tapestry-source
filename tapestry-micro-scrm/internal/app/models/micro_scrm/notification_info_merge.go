package microScrmModel

import (
	"micro-scrm/internal/app/models"
)

type NotificationInfoMergeModel struct {
	models.BaseModel
	Type         int                                 `json:"type" gorm:"column:type;type:tinyint(1);not null;default:0;comment:传讯类型，1周传讯，2转货传讯,3紧急传讯"`
	WeekNum      int                                 `json:"week_num" gorm:"column:week_num;type:int(10);not null;default:0;comment:第几周"`
	Title        string                              `json:"title" gorm:"column:title;type:varchar(255);not null;comment:标题"`
	Cover        string                              `json:"cover" gorm:"column:cover;type:varchar(255);not null;comment:封面图"`
	CrowdType    int                                 `json:"crowd_type" gorm:"column:crowd_type;type:int(10);not null;default:1;comment:通知人群 1 白名单id 2企微人员id"`
	CrowdIds     string                              `json:"crowd_ids" gorm:"column:crowd_ids;type:text;not null;comment:通知的人群id集合，英文逗号分割"`
	SendType     int                                 `json:"send_type" gorm:"column:send_type;type:tinyint(1);not null;default:0;comment:发布类型，0定时发布，1保存后发布，2暂不发布"`
	SendRoName   string                              `json:"send_ro_name" gorm:"column:send_ro_name;type:varchar(255);not null;comment:发布人姓名"`
	SendTime     int64                               `json:"send_time" gorm:"column:send_time;type:int(10);not null;default:0;comment:发布的时间,(定时发布的时间，已发布的时间)"`
	Status       int                                 `json:"status" gorm:"column:status;type:tinyint(1);not null;default:0;comment:发布的状态，0未发布，1已发布"`
	AdminID      int64                               `json:"admin_id" gorm:"column:admin_id;type:bigint(20);not null;default:0;comment:管理员id"`
	MemberID     int64                               `json:"member_id" gorm:"column:member_id;type:bigint(20);not null;default:0;comment:导购的member_id"`
	MemberName   string                              `json:"member_name" gorm:"column:member_name;type:varchar(255);not null;comment:创建人姓名"`
	PdfUrl       string                              `json:"pdf_url" gorm:"column:pdf_url;type:varchar(255);not null;comment:pdf文件cos地址"`
	EmailLoopUrl string                              `json:"email_loop_url" gorm:"column:email_loop_url;type:varchar(255);not null;comment:邮件发送的周传讯loop网址"`
	ZipUrl       string                              `json:"zip_url" gorm:"column:zip_url;type:varchar(255);not null;comment:打包的zip url"`
	ContentList  []NotificationInfoMergeContentModel `json:"content_list" gorm:"foreignKey:MergeId;references:ID"`
}

// TableName returns the table name of the NotificationInfoMerge model
func (NotificationInfoMergeModel) TableName() string {
	return "notification_info_merge"
}

const (
	NotificationCrowdWhite = iota + 1 // 白名单
	NotificationMember                // 导购
)

const (
	NotificationMergeNotSend = iota // 未发送
	NotificationMergeSend           // 已发送
)

const (
	NotificationMergeTimeSend = iota // 0定时发布
	NotificationMergeNowSend         // 1保存后发布
)

const (
	NotificationTypeWeek = 1 // 周传讯
)
