package microScrmModel

import "micro-scrm/internal/app/models"

type ReserveServiceScoreModel struct {
	models.BaseModel
	ReserveNumber  string `json:"reserve_number" gorm:"column:reserve_number;type:varchar(100);not null;default:0;comment:预约单编号"`
	ReserveOrderID int64  `json:"reserve_order_id" gorm:"column:reserve_order_id;type:bigint(20);not null;comment:预约单表id"`
	UserName       string `json:"user_name" gorm:"column:user_name;type:varchar(255);not null;comment:评论的用户昵称"`
	UserUnionid    string `json:"user_unionid" gorm:"column:user_unionid;type:varchar(255);not null;comment:评论的用户unionid"`
	UserExternalID string `json:"user_external_id" gorm:"column:user_external_id;type:varchar(255);not null;comment:评论的用户external_id"`
	MemberID       int64  `json:"member_id" gorm:"column:member_id;type:bigint(20);not null;default:0;comment:服务的导购id"`
	MemberName     string `json:"member_name" gorm:"column:member_name;type:varchar(255);not null;comment:服务的导购昵称"`
	MemberUserID   string `json:"member_user_id" gorm:"column:member_user_id;type:varchar(255);not null;comment:服务导购的user_id"`
	StoreID        int64  `json:"store_id" gorm:"column:store_id;type:bigint(20);not null;default:0;comment:服务的门店id"`
	StoreName      string `json:"store_name" gorm:"column:store_name;type:varchar(255);not null;comment:服务的门店名称"`
	ServiceScore   int32  `json:"service_score" gorm:"column:service_score;type:tinyint(1);not null;default:0;comment:默认服务分数为5分（KS端店员服务,sw端活动安排）"`
	StoreScore     int32  `json:"store_score" gorm:"column:store_score;type:tinyint(1);not null;default:0;comment:门店的评分，默认5分（KS端店铺环境,sw端整体服务）"`
	Remark         string `json:"remark" gorm:"column:remark;type:varchar(255);not null;comment:用户的评论信息"`
	FlowScore      int32  `json:"flow_score" gorm:"column:flow_score;type:tinyint(1);not null;default:0;comment:流程的评分（ks端服务流程）"`
	ReserveTime    int    `json:"reserve_time" gorm:"column:reserve_time;type:int(10);not null;default:0;comment:预约时间(预约单的到店时间)"`
}

// TableName returns the table name of the ReserveServiceScore model
func (ReserveServiceScoreModel) TableName() string {
	return "reserve_service_score"
}
