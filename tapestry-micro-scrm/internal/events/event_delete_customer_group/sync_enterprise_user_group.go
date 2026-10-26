package event_delete_customer_group

import (
	"context"
	"encoding/json"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/events/base"
	"micro-scrm/internal/events/entity"
	"micro-scrm/internal/orm"
)

var _ base.IListener = &syncEnterpriseUserGroup{}

type syncEnterpriseUserGroup struct {
	base.BaseListener
}

func NewSyncEnterpriseUserGroup(ctx context.Context) *syncEnterpriseUserGroup {
	return &syncEnterpriseUserGroup{
		BaseListener: base.BaseListener{
			Name: "Listener::同步外部联系人的客户分组;enterprise_user.groups",
			Ctx:  ctx,
		},
	}
}

func (l syncEnterpriseUserGroup) Handle(data entity.Data) {
	go l.syncCustomerGroup(data.Ctx, data.EnterpriseID, data.CustomerGroupID)
}

func (l syncEnterpriseUserGroup) SyncHandle(data entity.Data) {
	l.syncCustomerGroup(data.Ctx, data.EnterpriseID, data.CustomerGroupID)
}

func (l syncEnterpriseUserGroup) syncCustomerGroup(ctx ctx.Context, enterpriseID, customerGroupID int64) {
	// 获取包含该标签的用户集合
	infoList, err := orm.Service.WorkExternalUser.GetExternalUserByCustomerGroupID(enterpriseID, customerGroupID)
	if err != nil {
		log.Info(ctx).Msgf("%s::查询聚合用户信息错误: %+v", l.Name, err)
	}

	for _, item := range infoList {
		groupIDList := []int64{}
		json.Unmarshal([]byte(item.Groups), &groupIDList)
		newGroupIDList := []int64{}
		for _, id := range groupIDList {
			if id == customerGroupID {
				continue
			}
			newGroupIDList = append(newGroupIDList, id)
		}

		// 更新聚合用户信息
		if updateErr := orm.Service.WorkExternalUser.UpdateExternalUserGroupsOne(enterpriseID, item.ID, newGroupIDList); updateErr != nil {
			log.Info(ctx).Msgf("%s::查询聚合用户信息错误: %+v", l.Name, updateErr)
		}
	}
}
