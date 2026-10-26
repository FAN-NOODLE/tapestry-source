package callback

import "micro-scrm/internal/pkg/sdk/work"

// QwExternalCallbackQueue 外部联系回调队列
const QwExternalCallbackQueue = "qw_external_callback_queue"

var (
	// AddExternalContact 添加企业客户事件
	AddExternalContact = "add_external_contact"
	// EditExternalContact 编辑企业客户事件
	EditExternalContact = "edit_external_contact"
	// AddHalfExternalContact 外部联系人免验证添加成员事件
	AddHalfExternalContact = "add_half_external_contact"
	// DelExternalContact 成员删除外部联系人
	DelExternalContact = "del_external_contact"
	// DelFollowUser 外部联系人员删除成员
	DelFollowUser = "del_follow_user"
	// TransferFail 客户接替失败事件
	TransferFail = "transfer_fail"
	// MigrateExternalContact 迁移企微历史数据
	MigrateExternalContact = "migrate_external_contact"
)

// 群状态
const (
	GroupUpdateStatus = "update"
)

const (
	WorkChatGroupMemberStatusJoin  = 1 // 企微群状态：加群
	WorkChatGroupMemberStatusLeave = 2 // 企微群状态：离群
	// 外部联系人回调地址
	CallbackUrl = "%s/micro-scrm/v1/callback/member-relation"
)

// ExternalCallbackEventOfChangeExternalContact 外部联系人通讯录变更事件
const ExternalCallbackEventOfChangeExternalContact = "change_external_contact"

// ExternalCallbackEventOfChangeExternalChat 外部联系人群变更通知
const ExternalCallbackEventOfChangeExternalChat = "change_external_chat"

// 群事件变更类型
const (
	Event_ChangeExternalChat_ChangeType_Create  = "create"
	Event_ChangeExternalChat_ChangeType_Update  = "update"
	Event_ChangeExternalChat_ChangeType_Dismiss = "dismiss"
)

// Callback 业务实体（结构体）
type Callback struct{}

// QwExternalCallback  回调参数
type QwExternalCallback struct {
	MsgSignature string `form:"msg_signature"`
	Nonce        string `form:"nonce"`
	Timestamp    string `form:"timestamp"`
	ToUserName   string `xml:"ToUserName"`
}

// ExternalCallbackMsg 企业微信回调信息获取
type ExternalCallbackMsg struct {
}

// ExternalGroup 群变更
type ExternalGroup struct {
	ChatID string `json:"chat_id" desc:"群ID"`
}

// Text 文本
type Text struct {
	Content string `json:"content"`
}

// Image 图片
type Image struct {
	MediaId string `json:"media"`
	PicUrl  string `json:"pic_url"`
}

// Link 链接
type Link struct {
	Title  string `json:"title"`
	PicUrl string `json:"pic_url"`
	Desc   string `json:"desc"`
	Url    string `json:"url"`
}

// Miniprogram 小程序
type Miniprogram struct {
	Title      string `json:"title"`
	PicMediaId string `json:"pic_media_id"`
	Appid      string `json:"appid"`
	Page       string `json:"page"`
	PicUrl     string `json:"pic_url"`
}

// WelcomePush 欢迎语
type WelcomePush struct {
	WelcomeCode string      `json:"welcome_code"`
	Text        Text        `json:"text"`
	Image       Image       `json:"image"`
	Link        Link        `json:"link"`
	Miniprogram Miniprogram `json:"miniprogram"`
}

//UserTagHandlerRequest 用户标签处理请求
type UserTagHandlerRequest struct {
	AccessToken    string   `json:"access_token" desc:"调用接口凭证"`
	Userid         string   `json:"userid" desc:"添加外部联系人的userid"`
	ExternalUserid string   `json:"external_userid" desc:"外部联系人userid"`
	AddTag         []string `json:"add_tag" desc:"要标记的标签列表"`
	RemoveTag      []string `json:"remove_tag" desc:"要移除的标签列表"`
}

// EnterpriseOfficialRequest 企业公众号绑定信息同步请求参数
type EnterpriseOfficialRequest struct {
	EnterpriseID int64  `json:"enterprise_id" binding:"required" desc:"企业id"`
	Type         int    `json:"type" binding:"required" desc:"绑定类型：1-绑定，2-解绑"`
	AppID        string `json:"appid" binding:"required" desc:"公众号appid"`
	Name         string `json:"name" binding:"required" desc:"公众号名称"`
}

// SyncMemberJob 同步任务消息
type ExternalCallbackJob struct {
	Job  string `json:"job"`
	Data *work.MemberExternalCallback
}

// GroupChatData 消息队列中群变化的消息
type GroupChatData struct {
	ChatID       int64 `json:"chat_id"`
	JoinExternal []struct {
		UserID  string `json:"external_user_id"`
		AddTime int64  `json:"add_time"`
	} `json:"add_external"`
	LeaveExternal []string
}
