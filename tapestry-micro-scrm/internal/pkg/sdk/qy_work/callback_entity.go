package qyWork

const (
	ChangeContactEvent         = "change_contact"          // 通讯录事件
	ChangeExternalContactEvent = "change_external_contact" // 企业客户事件
	ChangeExternalChatEvent    = "change_external_chat"    // 客户群创建事件
	ChangeExternalTagEvent     = "change_external_tag"     // 企业客户标签事件
)

const (
	CreateUser = "create_user" // 新增成员事件
	UpdateUser = "update_user" // 更新成员事件
	DeleteUser = "delete_user" // 删除成员事件
)

const (
	CreateParty = "create_party" //  新增部门事件
	UpdateParty = "update_party" //  更新成员事件
	DeleteParty = "delete_party" //  删除成员事件
)

const (
	UpdateTag = "update_tag" //  更新标签事件
)

const (
	AddExternalContact     = "add_external_contact"      // 添加企业客户事件
	EditExternalContact    = "edit_external_contact"     // 编辑企业客户事件
	AddHalfExternalContact = "add_half_external_contact" // 外部联系人免验证添加成员事件
	DelExternalContact     = "del_external_contact"      // 成员删除外部联系人
	DelFollowUser          = "del_follow_user"           // 外部联系人员删除成员
	TransferFail           = "transfer_fail"             // 客户接替失败事件
	MigrateExternalContact = "migrate_external_contact"  // 迁移企微历史数据
)

type ExternalContent struct {
	ToUserName     string `json:"ToUserName"`
	FromUserName   string `json:"FromUserName"`
	CreateTime     string `json:"CreateTime"`
	MsgType        string `json:"MsgType"`
	Event          string `json:"Event"`
	ChangeType     string `json:"ChangeType"`
	UserID         string `json:"UserID"`
	ExternalUserID string `json:"ExternalUserID"`
	WelcomeCode    string `json:"WelcomeCode"`
}

// Verify 接口验证
type Verify struct {
	CorpID       string `form:"corp_id"`
	MsgSignature string `form:"msg_signature"`
	Timestamp    string `form:"timestamp"`
	Nonce        string `form:"nonce"`
	EchoStr      string `form:"echostr"`
}

// Request 回调请求参数
type Request struct {
	MsgSignature string `json:"msg_signature"`
	Timestamp    string `json:"timestamp"`
	Nonce        string `json:"nonce"`
	ToUserName   string `xml:"ToUserName"`
	AgentID      string `xml:"AgentID"`
	Encrypt      string `xml:"Encrypt"`
}

// DecryptBaseData 解密基本数据
type DecryptBaseData struct {
	ToUserName string `xml:"ToUserName"`
	CreateTime int    `xml:"createTime"`
	MsgType    string `xml:"msgType"`
	Event      string `xml:"event"`
	ChangeType string `xml:"changeType"`
}
