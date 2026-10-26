package change_external_contact

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/events/entity"
)

// 外部联系人免验证添加成员事件

// 参数	说明
// ToUserName	企业微信CorpID
// FromUserName	此事件该值固定为sys，表示该消息由系统生成
// CreateTime	消息创建时间 （整型）
// MsgType	消息的类型，此时固定为event
// Event	事件的类型，此时固定为change_external_contact
// ChangeType	此时固定为add_half_external_contact
// UserID	企业服务人员的UserID
// ExternalUserID	外部联系人的userid，注意不是企业成员的帐号
// State	添加此用户的「联系我」方式配置的state参数，可用于识别添加此用户的渠道
// WelcomeCode	欢迎语code，可用于发送欢迎语

type addHalfExternalContact struct {
}

func NewAddHalfExternalContact() *addHalfExternalContact {
	return &addHalfExternalContact{}
}

func (listener *addHalfExternalContact) Handle(ctx ctx.Context, data entity.ExternalUserCallbackEntity) {
	if !(data.Event == entity.EVENT_CHANGE_EXTERNAL_CONTACT && data.ChangeType == entity.CHANGE_TYPE_ADD_HALF_EXTERNAL_CONTACT) {
		return
	}

	// 任务计划更新状态

}
