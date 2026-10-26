package event_change_user_tag

import (
	"fmt"

	"micro-scrm/internal/pkg/db"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/sqlplus"
	"micro-scrm/internal/events/base"
	"micro-scrm/internal/events/entity"
	"micro-scrm/internal/orm"
)

// 同步用户标签
// 把用户的标签同步到分库表enterprise_user.tag

var dbClient *sqlplus.ClusterDB

var _ base.IListener = &syncUserTagListener{}

type syncUserTagListener struct {
	base.BaseListener
}

func NewSyncUserTagListener() *syncUserTagListener {
	dbClient = db.DB()
	return &syncUserTagListener{
		BaseListener: base.BaseListener{
			Name: "Listener::同步用户标签",
		},
	}
}

func (l syncUserTagListener) Handle(data entity.Data) {
	go l.syncUserTag(data.Ctx, data.EnterpriseID, data.RelationID)
}

func (l syncUserTagListener) SyncHandle(data entity.Data) {
	l.syncUserTag(data.Ctx, data.EnterpriseID, data.RelationID)
}

// 根据work_external_user_relation.id 获取用户的标签集合，更新enterprise_user.tags
func (l syncUserTagListener) syncUserTag(ctx ctx.Context, enterpriseID, relationID int64) {
	//defer dbClient.Close()
	externalUserID, userErr := orm.Service.WorkExternalUser.GetExternalUserIDByRelationID(enterpriseID, relationID)
	if userErr != nil {
		log.Info(ctx).Msgf("%s, 获取用户uid错误:%+v", l.Name, userErr)
		return
	}
	var tagIDs []int64
	sqlStr := fmt.Sprintf("SELECT DISTINCT(external_tag_id) FROM work_external_user_tag WHERE relation_id in (SELECT id FROM work_external_user_relation WHERE enterprise_id = ? AND work_external_user_id =? AND relation_status = 0 AND deleted_at=0) AND deleted_at = 0")
	err := dbClient.Select(ctx, &tagIDs, sqlStr, enterpriseID, externalUserID)
	if err != nil {
		log.Info(ctx).Msgf("externalUserID:%+v, 获取用户标签集合错误:%+v", externalUserID, err)
		return
	}
	// 需要对用户的tag进行去重
	tagM := make(map[int64]bool)
	var updateTagIDs []int64
	for _, v := range tagIDs {
		if _, ok := tagM[v]; ok {
			continue
		}
		tagM[v] = true
		updateTagIDs = append(updateTagIDs, v)
	}
	// 同步tag
	updateErr := orm.Service.WorkExternalUser.UpdateExternalUserTag(enterpriseID, externalUserID, updateTagIDs)
	if updateErr != nil {
		log.Info(ctx).Msgf("%s, %+v", l.Name, updateErr)
	}
}
