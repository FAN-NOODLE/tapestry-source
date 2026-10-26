package micro_scrm

type RelationCallbackRecordsModel struct {
	ID              int64  `json:"id" gorm:"column:id;type:bigint(20);not null"`
	ToUserName      string `json:"to_user_name" gorm:"column:to_user_name;type:varchar(32);not null;index:idx_corp_change_relation,to_user_name,to_user_name_2,to_user_name_index;comment:企业ID"`
	FromUserName    string `json:"from_user_name" gorm:"column:from_user_name;type:varchar(32);not null;comment:此事件该值固定为sys，表示该消息由系统生成"`
	CreateTime      uint   `json:"create_time" gorm:"column:create_time;type:int(10) unsigned;not null;index:to_user_name_2;default:0;comment:消息创建时间 （整型）"`
	MsgType         string `json:"msg_type" gorm:"column:msg_type;type:varchar(32);not null;comment:消息的类型，此时固定为event"`
	Event           string `json:"event" gorm:"column:event;type:varchar(32);not null;comment:事件的类型"`
	ChangeType      string `json:"change_type" gorm:"column:change_type;type:varchar(32);not null;index:idx_corp_change_relation,to_user_name,to_user_name_2;comment:此时固定为add_external_contact"`
	UserID          string `json:"user_id" gorm:"column:user_id;type:varchar(32);not null;index:idx_corp_change_relation;comment:企业服务人员的UserID"`
	ExternalUserID  string `json:"external_user_id" gorm:"column:external_user_id;type:varchar(32);not null;index:idx_corp_change_relation;comment:外部联系人的userid，注意不是企业成员的帐号"`
	State           string `json:"state" gorm:"column:state;type:varchar(32);not null;comment:添加此用户的「联系我」方式配置的state参数，可用于识别添加此用户的渠道"`
	WelcomeCode     string `json:"welcome_code" gorm:"column:welcome_code;type:varchar(128);not null;comment:欢迎语welcome code"`
	FailReason      string `json:"fail_reason" gorm:"column:fail_reason;type:varchar(32);not null;comment:失败理由"`
	ChatID          string `json:"chat_id" gorm:"column:chat_id;type:varchar(32);not null;comment:群ID"`
	Status          uint32 `json:"status" gorm:"column:status;type:tinyint(1) unsigned;not null;default:0;comment:状态：0-未处理，1-入队列成功，2-队列消费，3-处理成功，4-处理失败"`
	ChannelID       int64  `json:"channel_id" gorm:"column:channel_id;type:bigint(20);not null;default:0;comment:渠道id"`
	StatusTrace     uint32 `json:"status_trace" gorm:"column:status_trace;type:tinyint(1) unsigned;not null;default:0"`
	AddRelationTime int    `json:"add_relation_time" gorm:"column:add_relation_time;type:int(10);not null;index:to_user_name;default:0;comment:建立关系的时间"`
}

// TableName returns the table name of the RelationCallbackRecords model
func (RelationCallbackRecordsModel) TableName() string {
	return "relation_callback_records"
}
