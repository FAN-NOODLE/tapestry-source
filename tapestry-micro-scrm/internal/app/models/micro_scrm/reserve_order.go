package microScrmModel

import "micro-scrm/internal/app/models"

type ReserveOrderModel struct {
	models.BaseModel
	ReserveNumber      string                   `json:"reserve_number" gorm:"column:reserve_number;type:varchar(255);not null;default:0;comment:预约单的编号"`
	ReserveID          int64                    `json:"reserve_id" gorm:"column:reserve_id;type:bigint(20);not null;default:0;comment:预约服务的项目id"`
	ReserveTime        int                      `json:"reserve_time" gorm:"column:reserve_time;type:int(10);not null;default:0;comment:预约服务的时间"`
	StoreID            int64                    `json:"store_id" gorm:"column:store_id;type:bigint(20);not null;index:store_id_index;default:0;comment:门店id(部门id)"`
	StoreName          string                   `json:"store_name" gorm:"column:store_name;type:varchar(255);not null;comment:预约的门店名称"`
	StoreAddress       string                   `json:"store_address" gorm:"column:store_address;type:varchar(255);not null;comment:预约的门店地址"`
	ReserveType        int32                    `json:"reserve_type" gorm:"column:reserve_type;type:tinyint(1);default:0;comment:预约订单是否是分配单,0:默认预约单，1:分配的预约单"`
	Status             int32                    `json:"status" gorm:"column:status;type:tinyint(1);not null;default:0;comment:预约单的状态"`
	UserName           string                   `json:"user_name" gorm:"column:user_name;type:varchar(255);not null;comment:预约服务的用户名称"`
	UserUnionid        string                   `json:"user_unionid" gorm:"column:user_unionid;type:varchar(255);not null;index:unionid_index;comment:预约服务的用户unionid"`
	UserExternalID     string                   `json:"user_external_id" gorm:"column:user_external_id;type:varchar(255);not null;comment:预约服务的用户external_id"`
	MemberID           int64                    `json:"member_id" gorm:"column:member_id;type:bigint(20);not null;index:member_id_index;default:0;comment:预约单服务的导购id"`
	MemberName         string                   `json:"member_name" gorm:"column:member_name;type:varchar(255);not null;comment:预约单服务的导购名称"`
	ReserveRemark      string                   `json:"reserve_remark" gorm:"column:reserve_remark;type:varchar(255);not null;comment:预约服务的备注信息(导购)"`
	UserRemark         string                   `json:"user_remark" gorm:"column:user_remark;type:varchar(255);not null;comment:预约服务的备注信息(用户)"`
	IsOpen             int32                    `json:"is_open" gorm:"column:is_open;type:tinyint(1);not null;comment:用户是否打开预约单 0否，1是"`
	IsCommit           int32                    `json:"is_commit" gorm:"column:is_commit;type:tinyint(1);not null;comment:用户是否确认预约单 0否，1是"`
	InviteNum          int                      `json:"invite_num" gorm:"column:invite_num;type:int(10);not null;comment:邀约次数"`
	MemberUserID       string                   `json:"member_user_id" gorm:"column:member_user_id;type:varchar(255);not null;comment:预约单服务的导购user_id"`
	ReserveOrderSn     string                   `json:"reserve_order_sn" gorm:"column:reserve_order_sn;type:varchar(255);not null;comment:coach核销单号"`
	ReserveOrderPrice  int                      `json:"reserve_order_price" gorm:"column:reserve_order_price;type:int(10);not null;default:0;comment:coach核销导购的订单金额（元）"`
	IsSign             int32                    `json:"is_sign" gorm:"column:is_sign;type:tinyint(1);not null;default:0;comment:用户是否点击签到"`
	IsComment          int32                    `json:"is_comment" gorm:"column:is_comment;type:tinyint(1);not null;default:0;comment:用户是否评论，2为未评论,1为已评论"`
	CommentTime        int                      `json:"comment_time" gorm:"column:comment_time;type:int(10);not null;default:0;comment:评论时间"`
	ReserveOrderTime   int                      `json:"reserve_order_time" gorm:"column:reserve_order_time;type:int(10);not null;default:0;comment:coach核销单订单时间"`
	IsStore            int32                    `json:"is_store" gorm:"column:is_store;type:tinyint(1);not null;default:0;comment:是否是门店购买：1是，2否"`
	IsPay              int32                    `json:"is_pay" gorm:"column:is_pay;type:tinyint(1);not null;default:1;comment:是否成单：2否，1是"`
	IsConfirm          int32                    `json:"is_confirm" gorm:"column:is_confirm;type:tinyint(1);not null;default:1;comment:是否确认到店：2否，1是"`
	OrderOrigin        int32                    `json:"order_origin" gorm:"column:order_origin;type:tinyint(1);not null;default:1;comment:订单来源，1为到点邀约单，2为邀请函"`
	CustomerStatus     int32                    `json:"customer_status" gorm:"column:customer_status;type:tinyint(1);not null;default:0;comment:是否是新客，老客，1新客，2老客"`
	ReserveOrderStatus int32                    `json:"reserve_order_status" gorm:"column:reserve_order_status;type:tinyint(1);not null;default:0;comment:coach的cdp安客诚推送过来的订单匹配情况，1匹配成功,2取消匹配,3邀约到期"`
	IsVip              int32                    `json:"is_vip" gorm:"column:is_vip;type:tinyint(1);not null;default:0;comment:是否是会员，默认为0，1为是"`
	ReserveEndTime     int                      `json:"reserve_end_time" gorm:"column:reserve_end_time;type:int(10);not null;default:0;comment:邀约单到店时间的结束时间"`
	ReserveService     ReserveServiceModel      `json:"-" gorm:"foreignKey:ReserveID;references:ID"`
	ReserveOrderSales  []ReserveOrderSalesModel `json:"reserve_order_sales" gorm:"foreignKey:ReserveNumber;references:ReserveNumber"`
}

// TableName returns the table name of the ReserveOrder model
func (ReserveOrderModel) TableName() string {
	return "reserve_order"
}

const (
	RUNNING    = iota // 0 进行中
	FINISH            // 1 邀约成功：用户到店且有关联订单，若用户在邀约完成阶段未评价，SA在该阶段仍可发送邀请评价链接;
	WAITFINISH        // 2 邀约完成 (用户到店，且导购手动确认用户到店，该阶段可邀请客户评价)
	CANCEL            // 3 已取消
	OVERTIME          // 4 已过期，
)

const (
	ORDER_ORIGIN_RESERVE = iota + 1 // 订单来源，1为到点邀约单，2为邀请函
	ORDER_ORIGIN_INVITE
)

const (
	PAY_SUCCESS = iota + 1 // 是否成单：2否，1是
	PAY_FAIL
)

const (
	CONFIRM_YES = iota + 1 // 是否确认到店：2否，1是
	CONFIRM_NO
)

const (
	COMMENT_YES = iota + 1 // 用户是否评论，2为未评论,1为已评论'
	COMMENT_NO
)

const (
	STORE_NORMAL = iota
	STORE_YES    // 是否是门店购买：1是，2否
	STORE_NO
)
