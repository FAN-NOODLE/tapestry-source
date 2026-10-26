package event_change_user_tag

import (
	"micro-scrm/internal/events/base"
	"micro-scrm/internal/events/entity"
	"micro-scrm/internal/logic"
)

// 同步用户标签
// 把用户的标签同步到分库表enterprise_user.groups

var _ base.IListener = &syncUserGroupListener{}

type syncUserGroupListener struct {
	base.BaseListener
}

func NewSyncUserGroupListener() *syncUserGroupListener {
	return &syncUserGroupListener{
		BaseListener: base.BaseListener{
			Name: "Listener::同步用户标签客户分组",
		},
	}
}

// 并发
func (l syncUserGroupListener) Handle(data entity.Data) {
	// 同步外部联系人的客户分组信息
	go logic.NewUserGroupLogic(data.Ctx, data.EnterpriseID).SetRelationID(data.RelationID).SyncEnterpriseUserGroups()
}

// 同步
func (l syncUserGroupListener) SyncHandle(data entity.Data) {
	logic.NewUserGroupLogic(data.Ctx, data.EnterpriseID).SetRelationID(data.RelationID).SyncEnterpriseUserGroups()
}
