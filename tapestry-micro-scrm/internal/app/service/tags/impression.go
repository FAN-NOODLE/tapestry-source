package tagsSvc

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	tagsDao "micro-scrm/internal/app/dao/tags"
	tagsReq "micro-scrm/internal/app/http/request/tags"
	"micro-scrm/internal/app/models"
	tagModels "micro-scrm/internal/app/models/tag"
	"micro-scrm/internal/app/service"
	"sort"
)

type Impression struct {
	impressionDao         *tagsDao.Impression
	impressionRelationDao *tagsDao.ImpressionRelation
	service.Base
}

func NewImpression() *Impression {
	return &Impression{
		impressionDao:         tagsDao.NewImpression(),
		impressionRelationDao: tagsDao.NewImpressionRelation(),
	}
}

func (s Impression) GroupSearchList(c app.Ctx, request tagsReq.GroupSearchList) ([]tagModels.ImpressionTagGroupModel, int64, error) {
	return tagsDao.NewImpressionGroup().GetList(c, request)
}

func (s Impression) TagSearchList(c app.Ctx, request tagsReq.TagList) ([]tagModels.ImpressionTagModel, int64, error) {
	tagInfoList, total, err := s.impressionDao.GetList(c, request)
	if err != nil {
		return nil, total, err
	}

	var tagIds []int64
	for _, tagInfo := range tagInfoList {
		tagIds = append(tagIds, tagInfo.ID)
	}

	// 取出标签对应的人数
	tagCount, err := tagsDao.NewImpressionRelation().GetTagCount(ctx.Wrap(c), request.EnterpriseID, tagIds)
	if err != nil {
		return nil, total, err
	}

	// 循环查询每个标签对应的人数
	for k, v := range tagInfoList {
		if countVal, ok := tagCount[v.ID]; ok {
			v.UseNum = countVal
		}
		tagInfoList[k] = v
	}

	return tagInfoList, total, nil
}

func (s Impression) CreateTag(c app.Ctx, data tagsReq.CreateTag) error {
	// 判断标签是否存在，存在则更新,不存在则新增标签数据
	tagInfo, err := s.impressionDao.GetInfo(c, tagModels.ImpressionTagModel{
		ThirdTagID: data.ThirdTagId,
		MicroBaseModel: models.MicroBaseModel{
			EnterpriseID: data.EnterpriseID,
		},
	})
	if err != nil {
		return err
	}

	if tagInfo.ID > 0 {
		if tagInfo.Name != data.Name {
			tagInfo.Name = data.Name
			tagInfo.Description = data.Description
			err = s.impressionDao.Update(c, tagInfo)
		}
		return err
	}

	impressionData := tagModels.ImpressionTagModel{
		ThirdTagID:     data.ThirdTagId,
		Name:           data.Name,
		Description:    data.Description,
		GroupID:        data.GroupId,
		MicroBaseModel: models.CreateMicroBase(data.EnterpriseID, app.UID()),
	}
	err = s.impressionDao.Create(c, impressionData)

	return err
}

// CreateGroupAndTag 创建标签组合标签
func (s Impression) CreateGroupAndTag(c app.Ctx, data tagsReq.CreateGroupAndTag) error {
	groupId, err := NewImpressionGroup().CreateGroup(c, tagsReq.CreateGroup{
		EnterpriseID: data.EnterpriseID,
		ThirdGroupId: data.ThirdGroupId,
		GroupName:    data.GroupName,
	})
	if err != nil {
		return err
	}

	err = s.CreateTag(c, tagsReq.CreateTag{
		EnterpriseID: data.EnterpriseID,
		ThirdTagId:   data.ThirdTagId,
		Name:         data.Name,
		Description:  data.Description,
		GroupId:      groupId,
	})

	return err
}

// TagsAndTagGroupList 客户印象标签分组和标签
func (s Impression) TagsAndTagGroupList(c app.Ctx, req tagsReq.UserTagsReq) (info []tagModels.UserTagsGroupList, err error) {
	tagsAndGroupTagListRes, _, err := s.impressionDao.TagsAndTagGroupList(c, req.EnterpriseID)
	if err != nil {
		return nil, err
	}

	if len(req.Userids) == 0 {
		req.Userids = append(req.Userids, req.Userid)
	}
	// 查询给用户已打的印象标签
	userTagsList, errs := s.impressionRelationDao.GetListByExternalUserId(c, req.EnterpriseID, req.WorkExternalUserId, req.Userids)
	if errs != nil {
		return nil, errs
	}
	userTagsListMap := make(map[int64]tagModels.ImpressionTagRelationModel, 0)
	for _, item := range userTagsList {
		userTagsListMap[item.TagID] = item
	}
	var tagsGroupList []tagModels.UserTagsGroupList
	userGroupInfoMap := make(map[int64]tagModels.GroupInfoRes)
	userTagsMap := make(map[int64][]tagModels.ImpressionTagsInfo, 0)
	for _, v := range tagsAndGroupTagListRes {
		if _, ok := userGroupInfoMap[v.GroupID]; !ok {
			userGroupInfoMap[v.GroupID] = tagModels.GroupInfoRes{
				GroupBaseInfo: tagModels.GroupBaseInfo{
					GroupId:   v.GroupID,
					GroupName: v.GroupName,
				},
			}
		}
		// 是否有此标签
		var isHave = 0
		if _, ok := userTagsListMap[v.TagID]; ok {
			isHave = 1
		}
		userTagsMap[v.GroupID] = append(userTagsMap[v.GroupID], tagModels.ImpressionTagsInfo{
			TagId:  v.TagID,
			Name:   v.Name,
			IsHave: isHave,
		})
	}

	for k := range userGroupInfoMap {
		temp := tagModels.UserTagsGroupList{
			GroupInfoRes:       userGroupInfoMap[k],
			ImpressionTagsList: userTagsMap[k],
		}
		tagsGroupList = append(tagsGroupList, temp)
	}

	// 标签排序
	sort.Slice(tagsGroupList, func(i, j int) bool {
		return tagsGroupList[i].GroupId > tagsGroupList[j].GroupId
	})
	for _, v := range tagsGroupList {
		sort.Slice(v.ImpressionTagsList, func(i, j int) bool {
			return v.ImpressionTagsList[i].TagId > v.ImpressionTagsList[j].TagId
		})
	}

	return tagsGroupList, nil
}

func (s Impression) GetList(c app.Ctx, request tagsReq.ReqList) ([]tagModels.ImpressionTagModel, error) {
	return s.impressionDao.GetTagList(c, request)
}
