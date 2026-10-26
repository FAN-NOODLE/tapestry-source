package storeVisitModel

import (
	"micro-scrm/internal/app/models"
)

type StoreVisitMainModel struct {
	models.BaseModel
	Date           string `json:"date" gorm:"column:date;type:varchar(255);not null;default:'';comment:巡店日期"`
	TimeSpan       string `json:"time_span" gorm:"column:time_span;type:varchar(255);not null;default:'';comment:巡店时段"`
	UnixTime       int64  `json:"unix_time" gorm:"column:unix_time;type:int(10);not null;default:0;comment:预约到店开始时间"`
	ThirdId        int    `json:"third_id" gorm:"column:third_id;type:int(10);not null;default:0;index:third_id_index;comment:门店企微id"`
	SendMsgToStore int    `json:"send_msg_to_store" gorm:"column:send_msg_to_store;type:tinyint(1);not null;default:0;comment:是否给门店发消息通知"`
	VisitorIsShow  int    `json:"visitor_is_show" gorm:"column:visitor_is_show;type:tinyint(1);not null;default:0;comment:巡店人员可见:0否 1是"`
	StoreIsShow    int    `json:"store_is_show" gorm:"column:store_is_show;type:tinyint(1);not null;default:0;comment:门店人员可见:0否 1是"`
	Remark         string `json:"remark" gorm:"column:remark;type:text;comment:备注"`
	ToStoreTime    int    `json:"to_store_time" gorm:"column:to_store_time;type:int(10);not null;default:0;comment:到店打卡时间"`
	Status         int    `json:"status" gorm:"column:status;type:tinyint(1);not null;default:0;index:status_index;comment:状态 1 已预约 2 巡店中 3已发送报告"`
	Overview       string `json:"overview" gorm:"column:overview;type:text;not null;comment:概览"`
	Content        string `json:"content" gorm:"column:content;type:mediumtext;not null;comment:整体表现"`
	PdfUrl         string `json:"pdf_url" gorm:"column:pdf_url;type:varchar(1024);not null;default:'';comment:pdf文件Url"`
	IsReport       int    `json:"is_report" gorm:"column:is_report;type:int(1);not null;default:0;comment:是否已经通知巡店完成 0否 1是"`

	MemberList []StoreVisitMemberRelationModel `json:"member_list" gorm:"foreignKey:VisitID;references:ID"`
}

func (StoreVisitMainModel) TableName() string {
	return "store_visit_main"
}

const VisitStatusPlan = 1 //预约后
const VisitVisiting = 2   //巡店中
const VisitComplete = 3   //巡店完成
const VisitClose = 4      //关闭

const VisitSendReport = 1 //完成通知

const VisitorDefaultShow = 1 //巡店人员默认可见

const StoreDefaultShow = 1 //门店人员默认可见
