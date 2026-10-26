package microScrmModel

import "micro-scrm/internal/app/models"

const (
	// StatusTypeLongTerm StatusTypeValid 有效期，0长期有效，1指定时间
	StatusTypeLongTerm = 0
	StatusTypeValid    = 1

	IsUp   = 1 // 上架
	IsDown = 2 // 下架

	UserIsShow = 1 // 显示顾客署名
	UserIsHide = 2 // 不显示顾客署名

	MemberIsShow = 1 // 显示导购署名
	MemberIsHide = 2 // 不显示导购署名
)

type InviteTemplateModel struct {
	models.BaseModel
	ReserveServiceID int64  `json:"reserve_service_id" gorm:"column:reserve_service_id;type:bigint(20);not null;index:reserve_service_id_index;default:0;comment:关联reserve_service表id"`
	TitleName        string `json:"title_name" gorm:"column:title_name;type:varchar(255);not null;comment:标题名称"`
	StatusType       int    `json:"status_type" gorm:"column:status_type;type:tinyint(1);not null;default:0;comment:有效期，0长期有效，1指定时间"`
	StartTime        int    `json:"start_time" gorm:"column:start_time;type:int(11);not null;default:0;comment:开始时间"`
	EndTime          int    `json:"end_time" gorm:"column:end_time;type:int(11);not null;default:0;comment:结束时间"`
	BackgroundURL    string `json:"background_url" gorm:"column:background_url;type:text;comment:背景图多个url"`
	Content          string `json:"content" gorm:"column:content;type:varchar(255);not null;comment:文案"`
	ShareTitle       string `json:"share_title" gorm:"column:share_title;type:varchar(255);not null;comment:分享卡片标题"`
	ShareURL         string `json:"share_url" gorm:"column:share_url;type:text;comment:分享图片url"`
	IsUp             int    `json:"is_up" gorm:"column:is_up;type:tinyint(1);not null;default:1;comment:是否上下架:1上架2下架"`
	UserIsShow       int    `json:"user_is_show" gorm:"column:user_is_show;type:tinyint(1);not null;default:1;comment:是否展示顾客称谓:1展示2不展示"`
	MemberIsShow     int    `json:"member_is_show" gorm:"column:member_is_show;type:tinyint(1);not null;default:1;comment:是否展示导购称谓:1展示2不展示"`
}

// TableName returns the table name of the InviteTemplate model
func (InviteTemplateModel) TableName() string {
	return "invite_template"
}

type InviteTemplateInfo struct {
	InviteTemplateModel
	ReserveServiceName      string `json:"reserve_service_name" `
	ThirdIds                string `json:"third_ids" `
	ThirdNames              string `json:"third_names" `
	StatusStr               string `json:"status_str" `
	StatusTimeStr           string `json:"status_time_str" `
	ReserveServiceTime      int    `json:"reserve_service_time" `
	ReserveServiceStartTime int    `json:"reserve_service_start_time" `
	ReserveServiceEndTime   int    `json:"reserve_service_end_time" `
	ReserveServiceRole      int    `json:"reserve_service_role" `
	InviteOrderNum          int    `json:"invite_order_num" `
	IsUpStr                 string `json:"is_up_str" `
}

type InviteTemplateInfoList struct {
	InviteTemplateInfo
	RoleRules []RoleRule `json:"role_rules" `
}

type RoleRule struct {
	ID                 int64  `json:"id"` //is_dp为false时，id为member_id
	ThirdId            int    `json:"third_id"`
	Name               string `json:"name"` //is_dp为false时，name为member_name;is_dp为ture,name为store_name
	IsDP               bool   `json:"is_dp"`
	StoreType          int32  `json:"store_type"`
	MemSubscribeStatus int    `json:"mem_subscribe_status"`
	UserId             string `json:"user_id"`
}

type InviteTemplateListInfo struct {
	InviteTemplateModel
	ReserveServiceName string `json:"reserve_service_name" `
	ThirdIds           string `json:"third_ids" `
	ThirdNames         string `json:"third_names" `
	ReserveServiceTime int    `json:"reserve_service_time" `
}

type InviteTemplateDetailInfo struct {
	InviteTemplateModel
	ReserveServiceName string `json:"reserve_service_name" `
	ReserveServiceTime int    `json:"reserve_service_time" `
}
