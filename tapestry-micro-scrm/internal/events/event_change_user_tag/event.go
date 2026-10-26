package event_change_user_tag

import (
	"micro-scrm/internal/events/base"
)

type ChangeUserTagEvent struct {
	base.BaseEvent
}

func NewChangeUserTagEvent() *ChangeUserTagEvent {
	event := &ChangeUserTagEvent{
		BaseEvent: base.BaseEvent{
			Name: "Event::ChangeUserTag",
		},
	}

	// 注册监听者
	event.RegisterListener(
		NewSyncUserTagListener(),
		//NewSyncUserGroupListener(),
	)

	return event
}
