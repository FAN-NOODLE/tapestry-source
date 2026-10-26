package microScrmModel

import (
	"micro-scrm/internal/app/models"
)

const (
	VMFinishStatusIng         = iota //未提交
	VMFinishStatusCheck              //待审核
	VMFinishStatusFail               //未通过
	VMFinishStatusSuccess            //已通过
	VMFinishStatusDoubleCheck        //待复审
)

const (
	VMSpecialCheckNo   = iota // special未审核
	VMSpecialCheckWait        // special待审核
	VMSpecialCheckYes         // special已审核
)

type DisplayTaskAllotModel struct {
	models.BaseModel
	DisplayName   string                        `json:"display_name" gorm:"column:display_name;type:varchar(255);not null;comment:VM任务名称"`
	DisplayTaskID int64                         `json:"display_task_id" gorm:"column:display_task_id;type:bigint(20);not null;default:0;comment:VM任务表id"`
	StoreCode     string                        `json:"store_code" gorm:"column:store_code;type:varchar(255);not null;comment:门店code"`
	StoreID       int                           `json:"store_id" gorm:"column:store_id;type:int(10);not null;default:0;comment:门店id"`
	StoreName     string                        `json:"store_name" gorm:"column:store_name;type:varchar(255);not null;comment:门店名称"`
	StoreType     int32                         `json:"store_type" gorm:"column:store_type;type:tinyint(1);not null;default:0;comment:门店类型，1为奥莱，2为零售"`
	StoreCity     string                        `json:"store_city" gorm:"column:store_city;type:varchar(255);not null;comment:门店所属城市"`
	FinishStatus  int32                         `json:"finish_status" gorm:"column:finish_status;type:tinyint(1);not null;default:0;comment:任务完成状态,0为未提交，1为待审核，2为未通过，3为已通过"`
	Grade         uint32                        `json:"grade" gorm:"column:grade;type:tinyint(1) unsigned;not null;default:0;comment:任务评分"`
	OwnerMemberID int64                         `json:"owner_member_id" gorm:"column:owner_member_id;type:bigint(20);not null;default:0;comment:所属VM的member_id"`
	OwnerUserID   string                        `json:"owner_user_id" gorm:"column:owner_user_id;type:varchar(100);not null;comment:所属VM的工号"`
	OwnerName     string                        `json:"owner_name" gorm:"column:owner_name;type:varchar(255);not null;comment:所属VM的昵称"`
	ReviewRemark  string                        `json:"review_remark" gorm:"column:review_remark;type:text;not null;comment:整体评价留言"`
	SubmitTime    int64                         `json:"submit_time" gorm:"column:submit_time;type:int(10);not null;default:0;comment:提交时间"`
	ReviewTime    int64                         `json:"review_time" gorm:"column:review_time;type:int(10);not null;default:0;comment:审核时间"`
	SpecialUserId string                        `json:"special_user_id" gorm:"column:special_user_id;type:varchar(255);not null;comment:初审的员工user_id"`
	SpecialCheck  uint32                        `json:"special_check" gorm:"column:special_check;tinyint(1) unsigned;not null;default:0;comment:special员工是否审核，0否，1是"`
	SpecialRemark string                        `json:"special_remark" gorm:"column:special_remark;varchar(255);not null;comment:special的评价留言"`
	ImageList     []DisplayTaskAllotDataModel   `json:"image_list" gorm:"foreignKey:AllotID;references:ID"`
	Review        []DisplayTaskAllotReviewModel `json:"review" gorm:"foreignKey:AllotID;references:ID"`
	Submit        []DisplayTaskAllotSubmitModel `json:"submit" gorm:"foreignKey:AllotID;references:ID"`
}

// TableName returns the table name of the DisplayTaskAllot model
func (DisplayTaskAllotModel) TableName() string {
	return "display_task_allot"
}
