package workexternal

import "encoding/json"

type SelectKeep struct {
	Day          int64 `json:"day" desc:"开始时间" validate:"gte=0" required:"day 参数不正确"`
	EndAt        int64 `json:"end_at" desc:"结束时间" db:"end_at"`
	StartAt      int64 `json:"start_at" desc:"开始时间" db:"start_at"`
	EnterpriseId int64 `json:"enterprise_id" desc:"企业ID" db:"enterprise_id"`
}

type SelectUvRes struct {
	RelationNum int64 `json:"relation_num" db:"relation_num"`
}

type QwFriend struct {
	ReUv       int64 `json:"re_uv"`
	ReUvThan   int64 `json:"re_uv_than"`
	AddUv      int64 `json:"add_uv"`
	AddUvThan  int64 `json:"add_uv_than"`
	LossUv     int64 `json:"loss_uv"`
	LossUvThan int64 `json:"loss_uv_than"`
}

// DelRelationChannelTrace 外部联系人好友删员工关系
type DelRelationChannelTrace struct {
	ID             int64  `json:"id" db:"id"`
	ExternalUserID string `json:"external_user_id" db:"external_user_id"` // ExternalUserID 外部联系人的userid
	UserID         string `json:"user_id" db:"user_id"`                   // UserID 企业服务人员的UserID
	EnterpriseID   int64  `json:"enterprise_id" db:"enterprise_id"`       // EnterpriseID 企业ID
	DeleteTime     int    `json:"delete_time" db:"delete_time"`           // DeleteTime 表示好友删员工的时间
	DeleteChannel  int64  `json:"delete_channel" db:"delete_channel"`     // DeleteChannel 删除时存在的渠道
	AddTime        int    `json:"add_time" db:"add_time"`                 // AddTime 表示好友添加员工的时间
	AddChannel     int64  `json:"add_channel" db:"add_channel"`           // AddChannel 添加时存在的渠道
	Change         uint8  `json:"change" db:"change"`                     // Change 1=渠道值发送了变化，2=未发生变化
}

// WorkExternalTrace 外部联系人状态
type ExternalTrace struct {
	ID             int64  `json:"id" db:"id"`
	ExternalUserID string `json:"external_user_id" db:"external_user_id"` // ExternalUserID 外部联系人的userid
	EnterpriseID   int64  `json:"enterprise_id" db:"enterprise_id"`       // EnterpriseID 企业ID
	JoinTime       int    `json:"join_time" db:"join_time"`               // JoinTime 表示进入系统
	LeaveTime      int    `json:"leave_time" db:"leave_time"`             // LeaveTime 表示不存在系统
	LastLeaveTime  int    `json:"last_leave_time" db:"last_leave_time"`   // LastLeaveTime 上次离开时间
	NextJoinTime   int    `json:"next_join_time" db:"next_join_time"`     // NextJoinTime 下一次加入时间
}

// ExternalChannelTrace 渠道外部联系人状态
type ExternalChannelTrace struct {
	ID             int64  `json:"id" db:"id"`
	ExternalUserID string `json:"external_user_id" db:"external_user_id"` // ExternalUserID 外部联系人的userid
	EnterpriseID   int64  `json:"enterprise_id" db:"enterprise_id"`       // EnterpriseID 企业ID
	JoinTime       int    `json:"join_time" db:"join_time"`               // JoinTime 表示加入渠道的时间
	AddChannel     int64  `json:"add_channel" db:"add_channel"`           // AddChannel 添加时存在的渠道
	LeaveTime      int    `json:"leave_time" db:"leave_time"`             // LeaveTime 表示不存在系统
	LastLeaveTime  int    `json:"last_leave_time" db:"last_leave_time"`   // LastLeaveTime 上次离开时间
	NextJoinTime   int    `json:"next_join_time" db:"next_join_time"`     // NextJoinTime 下一次加入时间
}

// WorkRelationAddRecord 添加记录
type WorkRelationAddRecord struct {
	CorpID         string `db:"corp_id" json:"corp_id"`
	AddTime        int64  `db:"add_time" json:"add_time"`
	UserId         string `db:"user_id" json:"user_id"`
	ExternalUserId string `db:"external_user_id" json:"external_user_id"`
	AddChannel     int64  `db:"add_channel" json:"add_channel"`
	HasCK          int8   `db:"has_ck" json:"-"`
}

// 判断该记录关联的删除记录是否以及进入CK
func (r *WorkRelationAddRecord) InCK() bool {
	if r.HasCK == 1 {
		return true
	}

	return false
}

type KafkaMsg struct {
	*WorkRelationAddRecord
	DelTime int64 `db:"del_time" json:"del_time"`
}

// ToJSONString 输出 json 字符串
func (r *KafkaMsg) ToJSONString() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// FromJSONString 从 json 字符串读取
func (r *KafkaMsg) FromJSONString(s string) error {
	return json.Unmarshal([]byte(s), &r)
}

// UC 历史数据结构体
type HistoryUcRelationLog struct {
	ID             int64  `db:"id" json:"id"`
	ToUserName     string `db:"to_user_name" json:"to_user_name"`
	CreateTime     int64  `db:"create_time" json:"create_time"`
	ChangeType     string `db:"change_type" json:"change_type"`
	UserId         string `db:"user_id" json:"user_id"`
	ExternalUserId string `db:"external_user_id" json:"external_user_id"`
	State          string `db:"state" json:"state"`
	Status         int    `db:"status" json:"status"`
}
