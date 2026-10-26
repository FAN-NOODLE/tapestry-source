package events

import (
	"micro-scrm/internal/events/event_change_user_tag"
	"micro-scrm/internal/events/event_delete_customer_group"
	"micro-scrm/internal/events/event_external_user_callback"
)

var ScrmEvent *scrmEvent

type scrmEvent struct {
	ChangeUserTag        *event_change_user_tag.ChangeUserTagEvent
	DeleteCustomerGroup  *event_delete_customer_group.DeleteCustomerGroupEvent
	ExternalUserCallback *event_external_user_callback.ExternalUserCallbackEvent
}

func EventInit() {
	ScrmEvent = &scrmEvent{
		ChangeUserTag:        event_change_user_tag.NewChangeUserTagEvent(),
		DeleteCustomerGroup:  event_delete_customer_group.NewDeleteCustomerGroupEvent(),
		ExternalUserCallback: event_external_user_callback.NewExternalUserCallbackEvent(),
	}
}
