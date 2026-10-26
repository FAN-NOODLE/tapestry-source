package change_external_chat

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/events/entity"
	"micro-scrm/internal/orm"
)

// 群解散事件触发

type dismiss struct {
	name string
}

func NewDismiss() *dismiss {
	return &dismiss{
		name: "群解散事件触发",
	}
}

func (listener dismiss) Handle(ctx ctx.Context, data entity.ExternalUserCallbackEntity) {
	if !(data.Event == entity.CHANGEEXTERNALCHAT && data.ChangeType == entity.CHANGEEXTERNALCHAT_DISMISS) {
		return
	}

	// 处理解散群事件
	if err := orm.Service.GroupChat.UpdateWorkGroupChatDismiss(data.EnterpriseID, data.ChatID); err != nil {
		log.Info(ctx).Msgf("%s::操作失败:%v", listener.name, err)
	}
}
