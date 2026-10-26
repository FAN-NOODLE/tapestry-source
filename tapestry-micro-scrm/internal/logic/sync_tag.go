package logic

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/pkg/sdk/work"
	"micro-scrm/internal/pkg/uid"
	"micro-scrm/internal/pkg/workwechat"
	"sync"
)

type syncTagLogic struct {
	ctx          ctx.Context
	enterpriseID int64
	Name         string
}

func NewSyncTagLogic(ctx ctx.Context, enterpriseID int64) *syncTagLogic {
	return &syncTagLogic{
		ctx:          ctx,
		enterpriseID: enterpriseID,
		Name:         "同步标签逻辑",
	}
}

func (s syncTagLogic) SyncHandle() error {
	tagList, tagErr := s.getTagList()
	if tagErr != nil {
		log.Info(s.ctx).Msgf("%s::获取企微所有标签信息错误: %+v", s.Name, tagErr)
		return tagErr
	}

	// 检测，并更新标签信息
	for _, tagGroup := range tagList {
		s.updateTag(*tagGroup)
	}

	return nil
}

func (s syncTagLogic) Handle() error {
	tagList, tagErr := s.getTagList()
	if tagErr != nil {
		log.Info(s.ctx).Msgf("%s::获取企微所有标签信息错误: %+v", s.Name, tagErr)
		return tagErr
	}

	// 检测，并更新标签信息
	waitGroup := sync.WaitGroup{}
	for _, tagGroup := range tagList {
		waitGroup.Add(1)
		go func(tagGroup work.CorpTagGroupEntity) {
			defer waitGroup.Done()
			s.updateTag(tagGroup)
		}(*tagGroup)
	}
	waitGroup.Done()

	return nil
}

// 获取企微所有标签
func (s syncTagLogic) getTagList() (tagList []*work.CorpTagGroupEntity, err error) {
	// 获取企微所有的标签
	wwObj, objErr := workwechat.NewWorkWeChat(s.ctx, s.enterpriseID)
	if objErr != nil {
		log.Info(s.ctx).Msgf("%s::获取企微对象错误: %+v", s.Name, objErr)
		err = objErr
		return
	}
	tagList, err = wwObj.TagList()
	if err != nil {
		log.Info(s.ctx).Msgf("%s::获取企微所有标签信息错误: %+v", s.Name, err)
		return
	}

	return
}

func (s syncTagLogic) updateTag(wwTagGroup work.CorpTagGroupEntity) (err error) {
	if true == wwTagGroup.Delete {
		return
	}

	// 判断客户分组是否存在。如果不存在则新建分组，以及分组中的标签
	isExist, tagGroupInfo, groupErr := orm.Service.Tag.GetTagGroupInfoByWwTagGroupID(s.enterpriseID, wwTagGroup.GroupID)
	if groupErr != nil {
		log.Info(s.ctx).Msgf("%s::获取企业标签组详情错误: %+v", s.Name, groupErr)
		err = groupErr
		return
	}

	tagList := []micro_scrm.WorkExternalTagModel{}
	if isExist {
		tagList, err = orm.Service.Tag.GetTagListByWwTagGroupID(s.enterpriseID, wwTagGroup.GroupID)
		if err != nil {
			log.Info(s.ctx).Msgf("%s::获取企业标签集合错误: %+v", s.Name, err)
			return
		}
	} else {
		tagGroupInfo.ID = uid.ID(s.ctx)
		tagGroupInfo.GroupID = wwTagGroup.GroupID
		tagGroupInfo.GroupName = wwTagGroup.GroupName
		tagGroupInfo.EnterpriseID = s.enterpriseID
		tagGroupInfo.Order = wwTagGroup.Order
		tagGroupInfo.CreateTime = wwTagGroup.CreateTime

		createGroupErr := orm.Service.Tag.CreateTagGroup(tagGroupInfo)
		if createGroupErr != nil {
			log.Info(s.ctx).Msgf("%s::新建标签组失败: %+v", s.Name, createGroupErr)
			err = createGroupErr
			return
		}
	}

	// 企微中没有删除的标签
	noDelWxTagList := func(wwTagList []*work.CorpTagEntity) (tagList []work.CorpTagEntity) {
		for _, tag := range wwTagList {
			if false == tag.Delete {
				tagList = append(tagList, work.CorpTagEntity{
					ID:         tag.ID,
					Name:       tag.Name,
					CreateTime: tag.CreateTime,
					Order:      tag.Order,
					Delete:     tag.Delete,
				})
			}
		}
		return
	}(wwTagGroup.Tag)

	if len(noDelWxTagList) == len(tagList) {
		return
	}

	mapping := map[string]bool{}
	for _, tag := range tagList {
		mapping[tag.TagID] = true
	}

	tagInfoList := []micro_scrm.WorkExternalTagModel{}
	for _, noDelWxTag := range noDelWxTagList {
		if _, ok := mapping[noDelWxTag.ID]; !ok {
			tagInfoList = append(tagInfoList, micro_scrm.WorkExternalTagModel{
				MicroScrmBaseModel: model.MicroScrmBaseModel{
					BaseModel: model.BaseModel{
						ID: uid.ID(s.ctx),
					},
					EnterpriseID: s.enterpriseID,
				},
				GroupID:    tagGroupInfo.ID,
				TagID:      noDelWxTag.ID,
				Name:       noDelWxTag.Name,
				Order:      noDelWxTag.Order,
				CreateTime: noDelWxTag.CreateTime,
				Type:       micro_scrm.WET_Type_Enterprise,
			})
		}
	}

	createTagErr := orm.Service.Tag.BatchesCreate(tagInfoList)
	if createTagErr != nil {
		log.Info(s.ctx).Msgf("%s::新建标签组失败: %+v", s.Name, createTagErr)
		err = createTagErr
		return
	}

	return
}
