package customer

// DelMsgByCustomerStr 客户删除员工
const DelMsgByCustomerStr = "客户%s将员工%s从好友列表中删除"

// DelMsgByMemberStr 员工删除
const DelMsgByMemberStr = "员工%s将客户%s从好友列表中删除"

// SendTypeFromCustomer 客户删除员工
const SendTypeFromCustomer = 1

// SendTypeFromMember 员工删除客户
const SendTypeFromMember = 2

// 发送消息类型
var MessageType = map[int]string{
	1: "text",
	2: "textcard",
}

// CrsDomain 后台域名
var CrsDomain = "https://www.vchangyi.com"

// MessageTypeTextCard 消息模块
const MessageTypeTextCard = 2

// DelMessage 删除配置消息
type DelMessage struct {
	ID           int64  `json:"id" db:"id" desc:"ID"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业id"`
	Status       int    `json:"status" db:"status" desc:"开启开关"`
	SendType     int    `json:"send_type" db:"send_type" desc:"类别"`
	StatTxt      string `json:"stat_txt" db:"stat_txt" desc:"设置详情"`
	QyStatTxt    string `json:"qy_stat_txt" db:"qy_stat_txt" desc:"设置详情"`
	CreatedAt    int64  `json:"created_at" db:"created_at" desc:"添加时间"`
	UpdatedAt    int64  `json:"updated_at" db:"updated_at" desc:"更新时间"`
	DeletedAt    int64  `json:"deleted_at" db:"deleted_at" desc:"删除时间"`
}

// MessageDetailRequest message详情请求
type MessageDetailRequest struct {
	ID int64 `json:"id" db:"id" desc:"id" validate:"required" required:"企业ID不能为空"`
}

// MessageListRequest message list请求
type MessageListRequest struct {
	EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id" desc:"企业id" validate:"required" required:"企业ID不能为空"`
	SendType     int   `json:"send_type" db:"send_type" desc:"类别" validate:"required" required:"Type不能为空"`
}

// MessageListRequest message list返回
type MessageListResponse []*struct {
	ID           int64             `json:"id" db:"id" desc:"ID"`
	EnterpriseID int64             `json:"enterprise_id" db:"enterprise_id" desc:"企业id"`
	Status       int               `json:"status" db:"status" desc:"开启开关"`
	SendType     int               `json:"send_type" db:"send_type" desc:"类别"`
	StatTxt      string            `json:"stat_txt" db:"stat_txt" desc:"设置详情"`
	QyStatTxt    string            `json:"qy_stat_txt" db:"qy_stat_txt" desc:"设置详情"`
	StatTxtObj   *MessageMemberArr `json:"stat_txt_obj"`
	CreatedAt    int64             `json:"created_at" db:"created_at" desc:"添加时间"`
	UpdatedAt    int64             `json:"updated_at" db:"updated_at" desc:"更新时间"`
}

// StatResponse 返回结果
type StatResponse struct {
	Name   string `json:"name"`
	UserId string `json:"user_id"`
}

// UpdateMessageRequest 更新message请求
type UpdateMessageRequest struct {
	ID           int64  `json:"id" db:"id" desc:"ID"`
	EnterpriseID int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业id" validate:"required" required:"企业ID不能为空"`
	Status       int    `json:"status" db:"status" desc:"开启开关" validate:"required" required:"提醒设置开关必传"`
	SendType     int    `json:"send_type" db:"send_type" desc:"类别" validate:"required" required:"type不能为空"`
	StatTxt      string `json:"stat_txt" db:"stat_txt" desc:"设置详情" validate:"required" required:"设置不能为空"`
	QyStatTxt    string `json:"qy_stat_txt" db:"qy_stat_txt" desc:"设置详情" validate:"required" required:"设置不能为空"`
}

// AddMessageRequest 添加message 请求
type AddMessageRequest struct {
	ID           int64  `json:"id" db:"id" desc:"ID"`
	EnterpriseID int64  `json:"enterprise_id" desc:"企业id" validate:"required" required:"企业ID不能为空"`
	Status       int    `json:"status" desc:"开启开关" validate:"required" required:"提醒设置开关必传"`
	SendType     int    `json:"send_type" desc:"类别" validate:"required" required:"type不能为空"`
	StatTxt      string `json:"stat_txt" desc:"设置详情" validate:"required" required:"员工不能为空"`
	QyStatTxt    string `json:"qy_stat_txt" db:"qy_stat_txt" desc:"设置详情"`
}

// AddMessageResponse 添加message返回值
type AddMessageResponse struct {
	ID int64 `json:"id" db:"id" desc:"ID"`
}

// TranAddParam 转换请求
func TranAddParam(param *AddMessageRequest) *DelMessage {
	return &DelMessage{
		ID:           param.ID,
		EnterpriseID: param.EnterpriseID,
		Status:       param.Status,
		SendType:     param.SendType,
		StatTxt:      param.StatTxt,
	}
}

// TranUpdateParam 转换请求
func TranUpdateParam(param *UpdateMessageRequest) *DelMessage {
	return &DelMessage{
		ID:           param.ID,
		EnterpriseID: param.EnterpriseID,
		Status:       param.Status,
		SendType:     param.SendType,
		StatTxt:      param.StatTxt,
		QyStatTxt:    param.QyStatTxt,
	}
}

// SendMessageRequest 发送请求
type SendMessageRequest struct {
	EnterpriseID        int64  `json:"enterprise_id" db:"enterprise_id" desc:"企业id" validate:"required" required:"企业ID不能为空"`
	Type                int    `json:"type" validate:"required" required:"Type不能为空"`
	UserName            string `json:"user_name" validate:"required" required:"UserName不能为空"`
	MemberName          string `json:"member_name" validate:"required" required:"MemberName不能为空"`
	DeletedMemberUserID string `json:"deleted_member_userid" desc:"被删除/删除的员工企微userid"`
}

// MessageMember 发送消息成员
type MessageMember struct {
	Department []string `json:"department"`
	Member     []string `json:"member"`
}

// MessageMemberArr 发送成员
type MessageMemberArr struct {
	Department []StatResponse `json:"department"`
	Member     []StatResponse `json:"member"`
}
