package logic

// 标签相关的业务逻辑

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm/work_external_user"
	"micro-scrm/internal/pkg/uid"
)

type workExternalUserTagLogic struct {
	Name         string
	enterpriseID int64
	relationID   int64
	ctx          ctx.Context
}

func NewTagLogic(ctx ctx.Context, enterpriseID, relationID int64) *workExternalUserTagLogic {
	return &workExternalUserTagLogic{
		ctx:          ctx,
		Name:         "标签相关逻辑",
		enterpriseID: enterpriseID,
		relationID:   relationID,
	}
}

// SetWorkExternalUserTag 通过外部联系人的企微标签同步本地数据库标签
func (t workExternalUserTagLogic) SetWorkExternalUserTag(wwTagIDList []string) (err error) {
	wwLen := len(wwTagIDList)
	if 0 == wwLen {
		return
	}

	tagIDList := []int64{}
	tagIDList = orm.Service.Tag.GetTagIDListByQYID(t.enterpriseID, wwTagIDList)
	if len(tagIDList) != wwLen {
		// TODO 同步企微标签
		if err = NewSyncTagLogic(t.ctx, t.enterpriseID).SyncHandle(); err != nil {
			return
		}
		tagIDList = orm.Service.Tag.GetTagIDListByQYID(t.enterpriseID, wwTagIDList)
	}

	// 给外部联系人打本地数据库标签
	userTagList := []work_external_user.WorkExternalUserTagModel{}
	for _, tagID := range tagIDList {
		userTagList = append(userTagList, work_external_user.WorkExternalUserTagModel{
			MicroScrmBaseModel: model.MicroScrmBaseModel{
				BaseModel: model.BaseModel{
					ID: uid.ID(t.ctx),
				},
				EnterpriseID: t.enterpriseID,
			},
			RelationID:    t.relationID,
			ExternalTagID: tagID,
		})
	}
	err = orm.Service.WorkExternalUser.InsertUserTag(userTagList)
	if err != nil {
		log.Info(t.ctx).Msgf("%s::批量新增用户标签错误: %+v", t.Name, err)
		return
	}

	return
}
