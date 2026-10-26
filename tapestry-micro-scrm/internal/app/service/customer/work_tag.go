package customerSvc

import (
	"errors"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	tagsDao "micro-scrm/internal/app/dao/tags"
	customerReq "micro-scrm/internal/app/http/request/customer"
	"micro-scrm/internal/app/models"
	tagModels "micro-scrm/internal/app/models/tag"
	"micro-scrm/internal/app/service"
	qyWork "micro-scrm/internal/pkg/sdk/qy_work"
	"time"
)

type WorkTag struct {
	service.Base
	userTag *tagsDao.WorkExternalUserTag
}

func NewWorkTag() *WorkTag {
	return &WorkTag{
		userTag: tagsDao.NewWorkExternalUserTag(),
	}
}

func (s WorkTag) Operate(c ctx.Context, req customerReq.OperateWorkTags) error {
	external, err := microScrmDao.NewWxUserCorpExternal().GetOne(c, map[string]interface{}{
		"id": req.WorkExternalUserId,
	}, req.EnterpriseID)
	if err != nil {
		return err
	}
	relation, err := microScrmDao.NewUserRelation().GetOne(c, map[string]interface{}{
		"enterprise_id":         req.EnterpriseID,
		"member_id":             req.MemberId,
		"work_external_user_id": req.WorkExternalUserId,
	})
	if err != nil {
		return err
	}
	relationID := relation.ID
	var addTags []string
	var removeTags []string
	if len(req.AddTagIds) > 0 {
		many, err := tagsDao.NewWorkTagsDao().GetMany(c, map[string]interface{}{
			"enterprise_id": req.EnterpriseID,
			"id":            req.AddTagIds,
		})
		if err != nil {
			return err
		}
		for _, v := range many {
			addTags = append(addTags, v.TagID)
		}

	}
	if len(req.DelTagIds) > 0 {
		many, err := tagsDao.NewWorkTagsDao().GetMany(c, map[string]interface{}{
			"enterprise_id": req.EnterpriseID,
			"id":            req.DelTagIds,
		})
		if err != nil {
			return err
		}
		for _, v := range many {
			removeTags = append(removeTags, v.TagID)
		}
	}
	if len(addTags) == 0 && len(removeTags) == 0 {
		return nil
	}
	// 调用企微的接口
	qyNewWorkApi := qyWork.NewQyWork(req.EnterpriseID)
	res, err := qyNewWorkApi.MarkTag(c, qyWork.MarkTagReq{
		Userid:         req.OperateUserId,
		ExternalUserid: external.Openid,
		AddTag:         addTags,
		RemoveTag:      removeTags,
	})
	if err != nil {
		return err
	}
	if res.ErrCode > 0 {
		return errors.New(res.ErrMsg)
	}
	if len(req.DelTagIds) > 0 {
		// 删除企微标签的关联关系
		err = s.userTag.DeleteData(c, map[string]interface{}{
			"enterprise_id":   req.EnterpriseID,
			"relation_id":     relationID,
			"external_tag_id": req.DelTagIds,
		})
		if err != nil {
			return err
		}
	}
	if len(req.AddTagIds) > 0 {
		// 增加企微标签的关联关系
		var data []tagModels.WorkExternalUserTagModel
		for _, id := range req.AddTagIds {
			info := tagModels.WorkExternalUserTagModel{
				MicroBaseModel: models.MicroBaseModel{
					BaseModel: models.BaseModel{
						ID:        app.UID(),
						CreatedAt: time.Now().Unix(),
						UpdatedAt: time.Now().Unix(),
					},
					EnterpriseID: req.EnterpriseID,
				},
				RelationID:    relationID,
				ExternalTagID: id,
			}
			data = append(data, info)
		}
		return s.userTag.CreateData(c, data)
	}
	return nil
}
