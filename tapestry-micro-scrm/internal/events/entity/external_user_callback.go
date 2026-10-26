package entity

import "micro-scrm/internal/pkg/sdk/work"

// 外部联系人变更事件
const (
	EVENT_CHANGE_EXTERNAL_CONTACT         = "change_external_contact"
	CHANGE_TYPE_ADD_EXTERNAL_CONTACT      = "add_external_contact"
	CHANGE_TYPE_EDIT_EXTERNAL_CONTACT     = "edit_external_contact"
	CHANGE_TYPE_ADD_HALF_EXTERNAL_CONTACT = "add_half_external_contact"
	CHANGE_TYPE_DEL_EXTERNAL_CONTACT      = "del_external_contact"
	CHANGE_TYPE_DEL_FOLLOW_USER           = "del_follow_user"

	CHANGEEXTERNALCHAT         = "change_external_chat"
	CHANGEEXTERNALCHAT_DISMISS = "dismiss"

	UpdateDetailAddMember = "add_member"
	UpdateDetailDelMember = "del_member"
)

// ExternalUserCallbackEntity
// 企业微信客户回调事件数据实体
type ExternalUserCallbackEntity struct {
	EnterpriseID int64
	work.MemberExternalCallback
}
type ExternalUserInfoEntity struct {
	EnterpriseID int64
	GroupChatID  int64
	ToUserName   string
}
