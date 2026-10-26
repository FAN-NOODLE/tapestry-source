package customer

// PortraitRequest 客户画像详情请求
type PortraitRequest struct {
	ExternalUserID string `json:"external_userid"`
	MemberID       int64  `json:"member_id"`
	MemberUserID   string `json:"member_userid"`
	EnterpriseID   int64  `json:"enterprise_id"`
}

// PortraitResponse 客户画像详情返回结果
type PortraitResponse struct {
	UserID             int64  `json:"user_id" db:"user_id"`
	UnionID            string `json:"union_id" db:"union_id"`
	ExternalUserID     string `json:"external_userid" db:"external_userid"`
	WorkExternalUserID int64  `json:"work_external_user_id" db:"work_external_user_id"`
	NickName           string `json:"nick_name" db:"nick_name"`
	Avatar             string `json:"avatar" db:"avatar"`
	Mobile             string `json:"mobile" db:"mobile"`
	Gender             int    `json:"gender"`
	RelationID         int64  `json:"relation_id" db:"relation_id"`
	Remark             string `json:"remark" db:"remark"`
	//Birthday           int64  `json:"birthday" db:"birthday"`
	//CreateTime         int64  `json:"createtime" db:"createtime"`
	//BindGuideName         string              `json:"bind_guide_name" db:"bind_guide_name"`
	//ChannelName           string              `json:"channel_name"`
	//GroupName             []string            `json:"group_name"`
	//GroupNumber           int64               `json:"group_number"`
	//TagNames              []string            `json:"tag_names"`
	//BindGuideMemberID     int64               `json:"bind_guide_member_id" db:"bind_member_id"`
	//BindGuideMemberUserid string              `json:"bind_guide_member_userid"`
	//Age                   int                 `json:"age" db:"age"`
	//QQ                    string              `json:"qq" db:"qq"`
	//Email                 string              `json:"email" db:"email"`
	//Describe string `json:"describe" db:"describe"`

	//GenderIndex           int                 `json:"gender_index" db:"gender"`
	//UserOrder UserOrderStatistics `json:"user_order"`
	//Sex                   int8                `json:"sex"`
	//IsRelation            int                 `json:"is_relation"`
}

// UserOrderStatistics 订单统计的
type UserOrderStatistics struct {
	NumberOfOrders   int `json:"number_of_orders" desc:"订单成交数"`
	NumberOfRefunded int `json:"number_of_refunded" desc:"退款订单数"`
	OrdersAmount     int `json:"orders_amount" desc:"实收金额"`
	RefundAmount     int `json:"refund_amount" desc:"退款金额"`
}

// PortraitShoppingAddress 客户画像收货地址
type PortraitShoppingAddress struct {
	CustomerName    string `json:"customer_name" db:"customer_name"`
	CustomerMobile  string `json:"customer_mobile" db:"customer_mobile"`
	CustomerAddress string `json:"customer_address" db:"customer_address"`
}

// RemarkRequest 客户备注信息
type RemarkInfo struct {
	UserID       int64  `json:"user_id" db:"user_id" desc:"客户userID" validate:"required" required:"客户ID不能为空"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业ID" validate:"required" required:"企业ID不能为空"`
	Age          int    `json:"age" db:"age" desc:"年龄"`
	Birthday     int64  `json:"birthday" db:"birthday" desc:"生日时间戳"`
	BirthdayDate int    `json:"birthday_date" db:"birthday_date" desc:"生日格式:20200908"`
	Gender       int    `json:"gender" db:"gender" desc:"性别：0-未知 1-男性 2-女性"`
	Mobile       string `json:"mobile" db:"mobile" desc:"电话"`
	Email        string `json:"email" db:"email" desc:"邮箱"`
	Qq           string `json:"qq" db:"qq" desc:"QQ号"`
	Describe     string `json:"describe" db:"describe" desc:"描述信息"`
	CreatedAt    int64  `json:"created_at" db:"created_at" desc:"添加时间"`
	UpdatedAt    int64  `json:"updated_at" db:"updated_at" desc:"更新时间"`
	DeletedAt    int64  `json:"deleted_at" db:"deleted_at" desc:"删除时间"`
}

// ReturnUID 更新客户备注信息返回数据
type ReturnUID struct {
	UserID int64 `json:"user_id" db:"user_id"`
}

// PortraitChannel 渠道信息
type PortraitChannel struct {
	QrCodeID    int64  `json:"qr_code_id"`
	ChannelID   int64  `json:"channel_id"`
	ChannelName string `json:"channel_name"`
}

// PortraitBindGuide 绑定的导购信息
type PortraitBindGuide struct {
	BindMemberID   int64  `json:"bind_member_id" db:"bind_member_id"`
	BindMemberName string `json:"bind_member_name" db:"bind_member_name"`
	UserID         int64  `json:"user_id"`
}
