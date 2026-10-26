package change_external_contact

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/events/entity"
)

// 删除跟进成员事件

// 参数说明：
//
// 参数	说明
// ToUserName	企业微信CorpID
// FromUserName	此事件该值固定为sys，表示该消息由系统生成
// CreateTime	消息创建时间 （整型）
// MsgType	消息的类型，此时固定为event
// Event	事件的类型，此时固定为change_external_contact
// ChangeType	此时固定为del_follow_user
// UserID	企业服务人员的UserID
// ExternalUserID	外部联系人的userid，注意不是企业成员的帐号

type delFollowContact struct {
}

func NewDelFollowContact() *delFollowContact {
	return &delFollowContact{}
}

func (listener *delFollowContact) Handle(ctx ctx.Context, data entity.ExternalUserCallbackEntity) {
	log.Info(ctx).Msgf("BA_BIND_Bind_ONE::BA删除跟进顾客:: data: %v", data)
	if !(data.Event == entity.EVENT_CHANGE_EXTERNAL_CONTACT && data.ChangeType == entity.CHANGE_TYPE_DEL_FOLLOW_USER) {
		return
	}

	// 任务计划更新状态
}
