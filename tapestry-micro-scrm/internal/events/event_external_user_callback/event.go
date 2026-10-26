package event_external_user_callback

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/events/entity"
	"micro-scrm/internal/events/event_external_user_callback/change_external_chat"
	"micro-scrm/internal/events/event_external_user_callback/change_external_contact"
)

type IExternalListener interface {
	Handle(ctx.Context, entity.ExternalUserCallbackEntity)
}

type ExternalUserCallbackEvent struct {
	Name     string
	Listener []IExternalListener
}

func NewExternalUserCallbackEvent() *ExternalUserCallbackEvent {
	event := &ExternalUserCallbackEvent{
		Name: "Event::外部联系人事件回调",
	}

	event.registerListener(
		change_external_contact.NewAddExternalContact(),
		change_external_contact.NewAddHalfExternalContact(),
		change_external_contact.NewDelExternalContact(),
		change_external_contact.NewDelFollowContact(),
		// change_external_chat.NewChangeExternalChatAddMember(),
		// change_external_chat.NewChangeExternalChatDeleteMember(),
		change_external_chat.NewDismiss(),
	)

	return event
}

// RegisterListener 重构注册
func (event *ExternalUserCallbackEvent) registerListener(listenerList ...IExternalListener) {
	event.Listener = append(event.Listener, listenerList...)
}

func (event *ExternalUserCallbackEvent) Dispatch(ctx ctx.Context, data entity.ExternalUserCallbackEntity) {
	for _, item := range event.Listener {
		go item.Handle(ctx, data)
	}
}
