package event_delete_customer_group

import (
	"context"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/events/base"
	"micro-scrm/internal/events/entity"
	"micro-scrm/internal/orm"
)

var _ base.IListener = &syncDeleteUserGroupDetail{}

type syncDeleteUserGroupDetail struct {
	base.BaseListener
}

func NewSyncDeleteUserGroupDetail(ctx context.Context) *syncDeleteUserGroupDetail {
	return &syncDeleteUserGroupDetail{
		BaseListener: base.BaseListener{
			Name: "Listener::同步删除客户分组详情",
			Ctx:  ctx,
		},
	}
}

func (l syncDeleteUserGroupDetail) Handle(data entity.Data) {
	go l.syncDeleteCustomerGroupDetail(data.Ctx, data.EnterpriseID, data.CustomerGroupID)
}

func (l syncDeleteUserGroupDetail) SyncHandle(data entity.Data) {
	l.syncDeleteCustomerGroupDetail(data.Ctx, data.EnterpriseID, data.CustomerGroupID)
}

func (l syncDeleteUserGroupDetail) syncDeleteCustomerGroupDetail(ctx ctx.Context, enterpriseID, customerGroupID int64) {
	if err := orm.Service.UserGroup.DelCustomerGroupDetailByCustomerGroupID(enterpriseID, customerGroupID); err != nil {
		log.Info(ctx).Msgf("%s::%+v", l.Name, err)
	}
}
