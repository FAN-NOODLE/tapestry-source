package event_delete_customer_group

import (
	"context"
	"micro-scrm/internal/events/base"
)

type DeleteCustomerGroupEvent struct {
	base.BaseEvent
}

func NewDeleteCustomerGroupEvent() *DeleteCustomerGroupEvent {
	event := &DeleteCustomerGroupEvent{
		BaseEvent: base.BaseEvent{
			Name: "Event::DeleteCustomerGroup",
			Ctx:  context.TODO(),
		},
	}

	// 注册监听者
	event.RegisterListener(
		NewSyncEnterpriseUserGroup(event.Ctx),
		NewSyncDeleteUserGroupDetail(event.Ctx),
	)

	return event
}
