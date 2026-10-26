package storeVisitModel

import (
	"micro-scrm/internal/app/models"
)

type StoreVisitMemberRelationModel struct {
	models.BaseModel
	VisitID      int64  `json:"visit_id" gorm:"column:visit_id;type:bigint(20);not null;index:visit_id_index;comment:巡店表id"`
	MemberID     int64  `json:"member_id" gorm:"column:member_id;type:bigint(20);not null;index:member_id_index;comment:导购id"`
	Type         int    `json:"type" gorm:"column:type;type:tinyint(1);not null;default:0;comment:1 巡店发起人 2 巡店参与人 3 巡店店长 4 巡店店员 5 巡店跟进人员 6被@的人"`
	IsMM         int    `json:"is_mm" gorm:"column:is_mm;type:tinyint(1);not null;default:0;comment:是否是MM 0 否 1是"`
	Email        string `json:"email" gorm:"column:email;type:varchar(32);not null;default:'';comment:邮箱"`
	ReportIsShow int8   `json:"report_is_show" gorm:"column:report_is_show;type:tinyint(1);not null;default:0;comment:巡店报告是否可见 0 否 1是"`
}

func (StoreVisitMemberRelationModel) TableName() string {
	return "store_visit_member_relation"
}

const VisitCreator = 1      //创建者
const VisitVisitor = 2      //巡店参与人
const VisitStoreManager = 3 //店长
const VisitCopyMember = 4   //被同步者

const VisitReportIsShow = 1 //可查看报告

const IsMM = 1 //是MM管理
