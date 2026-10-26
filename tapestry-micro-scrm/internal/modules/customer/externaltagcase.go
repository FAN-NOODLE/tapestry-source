package customer

import (
	"database/sql"
	"errors"
	"micro-scrm/internal/pkg/helper"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/pkg/sdk/work"
	"micro-scrm/internal/pkg/uid"
)

type ExternalTagService struct {
	repo Repository
}

// var workApi api.Work

// GetExternalTagOfficial 根据标签id查询标签列表
func (e *ExternalTagService) GetExternalTagOfficial(c ctx.Context, enterpriseID int64, tagIDList []string) ([]*WorkExternalTag, error) {
	if len(tagIDList) == 0 {
		return nil, nil
	}
	list, err := e.repo.GetExternalTagsOfficial(c, enterpriseID, tagIDList)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return list, nil
}

// AddExternalTag 添加标签
func (e *ExternalTagService) AddExternalTag(c ctx.Context, data []*WorkExternalTag) (int64, error) {
	if len(data) == 0 {
		return 0, nil
	}
	return e.repo.AddExternalTag(c, data)
}

// GetAndCreateTagListFromWork 根据标签id从企业微信获取标签信息，然后加入到数据库中
func (e *ExternalTagService) GetAndCreateTagListFromWork(c ctx.Context, enterpriseID int64, domain string, idList []string) ([]*WorkExternalTag, error) {
	if len(idList) == 0 {
		return nil, nil
	}
	// 调用企业微信接口查询标签列表信息
	corpTagList, err := e.GetCorpTagList(c, domain, idList)
	if err != nil {
		return nil, err
	}
	// 从企业微信结果中获取标签组列表
	var tagGroupList []*WorkExternalTagGroup
	for _, v := range corpTagList {
		tagGroupList = append(tagGroupList, v.TagGroupInfo)
	}

	// 从数据库中查询标签组信息，没有的话创建新的数据
	dbTagGroupList, err := e.GetAndCreateTagGroup(c, enterpriseID, tagGroupList)
	if err != nil {
		return nil, err
	}
	var dbTagGroupMap = make(map[string]*WorkExternalTagGroup)
	for _, v := range dbTagGroupList {
		dbTagGroupMap[v.GroupID] = v
	}
	// 从企微返回的数据中组装
	var tagList []*WorkExternalTag
	for _, v := range corpTagList {
		dbTagGroupInfo, ok := dbTagGroupMap[v.TagGroupInfo.GroupID]
		if !ok {
			log.Info(c).Msgf("数据库中未查询到标签组信息: groupID:%s, groupMap:%+v", v.TagGroupInfo.GroupID, dbTagGroupMap)
		}
		for _, vv := range v.Tags {
			vv.GroupID = dbTagGroupInfo.ID
			vv.CreatedAt = time.Now().Unix()
			vv.UpdatedAt = time.Now().Unix()
			vv.EnterpriseID = enterpriseID
			vv.ID = uid.ID(c)
			tagList = append(tagList, vv)
		}
	}
	_, err = e.repo.AddExternalTag(c, tagList)
	if err != nil {
		return nil, err
	}
	return tagList, nil
}

func (e *ExternalTagService) GetAndCreateTagGroup(c ctx.Context, enterpriseID int64, tagGroupList []*WorkExternalTagGroup) ([]*WorkExternalTagGroup, error) {
	var groupIDList []string
	for _, v := range tagGroupList {
		groupIDList = append(groupIDList, v.GroupID)
	}
	// 从数据库查询已存在的标签组
	dbTagGroupList, err := e.GetExternalTagGroup(c, enterpriseID, groupIDList)
	if err != nil {
		return nil, err
	}
	var newTagGroupList []*WorkExternalTagGroup
	var dbTagGroupMap = make(map[string]*WorkExternalTagGroup)
	for _, v := range dbTagGroupList {
		dbTagGroupMap[v.GroupID] = v
	}
	// 生成数据库中不存在的标签组实体列表
	for _, v := range tagGroupList {
		if _, ok := dbTagGroupMap[v.GroupID]; !ok {
			v.ID = uid.ID(c)
			v.CreatedAt = time.Now().Unix()
			v.UpdatedAt = time.Now().Unix()
			v.EnterpriseID = enterpriseID
			newTagGroupList = append(newTagGroupList, v)
		}
	}
	if len(newTagGroupList) > 0 {
		_, err = e.repo.AddExternalTagGroup(c, newTagGroupList)
		if err != nil {
			log.Info(c).Msgf("添加标签组失败, err:%s, groupList:%+v", err.Error(), newTagGroupList)
			return nil, err
		}
	}

	ret := append(dbTagGroupList, newTagGroupList...)
	return ret, nil
}

// GetExternalTagGroup 查询标签组列表
func (e *ExternalTagService) GetExternalTagGroup(c ctx.Context, enterpriseID int64, groupIDList []string) ([]*WorkExternalTagGroup, error) {
	list, err := e.repo.GetExternalTagGroups(c, enterpriseID, groupIDList)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return list, nil
}

// GetCorpTagList 调用企业微信sdk获取标签信息
func (e *ExternalTagService) GetCorpTagList(ctx ctx.Context, domain string, idList []string) ([]*WorkTagListItem, error) {
	var workImpl work.QyWorkImpl
	workImpl.SetDomain(domain)
	response, err := workImpl.GetCorpTagList(ctx, idList)
	if err != nil {
		return nil, err
	}
	var result []*WorkTagListItem
	for _, v := range response.TagGroup {
		var item WorkTagListItem
		item.TagGroupInfo = &WorkExternalTagGroup{
			GroupID:    v.GroupID,
			GroupName:  v.GroupName,
			Order:      int64(v.Order),
			CreateTime: int64(v.CreateTime),
		}
		for _, vv := range v.Tag {
			item.Tags = append(item.Tags, &WorkExternalTag{
				TagID:      vv.ID,
				Name:       vv.Name,
				Order:      int64(vv.Order),
				CreateTime: int64(vv.CreateTime),
				Type:       ExternalTagTypeCompany,
			})
		}
		result = append(result, &item)
	}
	return result, nil
}

func (e *ExternalTagService) AddExternalUserTag(c ctx.Context, enterpriseID int64, relationID int64, tagIDList []int64) (int64, error) {
	if len(tagIDList) == 0 {
		return 0, nil
	}

	externalUserTagIds, err := e.repo.GetExternalUserTagsByTagId(c, enterpriseID, relationID, tagIDList)
	if err != nil {
		return 0, err
	}

	var userTags []*WorkExternalUserTag
	for _, v := range tagIDList {
		if !helper.InArray[int64](v, externalUserTagIds) {
			userTags = append(userTags, &WorkExternalUserTag{
				ID:            uid.ID(c),
				EnterpriseID:  enterpriseID,
				RelationID:    relationID,
				ExternalTagId: v,
				CreatedAt:     time.Now().Unix(),
				UpdatedAt:     time.Now().Unix(),
			})
		}
	}
	if len(userTags) == 0 {
		return 0, errors.New("存在相同关系,去重后没有要添加的数据")
	}
	return e.repo.AddExternalUserTags(c, userTags)
}

// GetExternalUserTags 获取外部联系人关系标签数据
func (e *ExternalTagService) GetExternalUserTags(c ctx.Context, relationID int64) ([]*WorkExternalUserTag, error) {
	return e.repo.QueryExternalUserTags(c, relationID)
}

// GetExternalUserRelation 获取外部联系人关系数据
func (e *ExternalTagService) GetExternalUserRelation(c ctx.Context, enterpriseID int64, workExternalUserID int64, memberUserID string) (*WorkUserMemberRelation, error) {
	return e.repo.GetExternalUserRelation(c, enterpriseID, workExternalUserID, memberUserID)
}

// GetExternalTagInfo 根据id获取标签信息
func (e *ExternalTagService) GetExternalTagInfo(c ctx.Context, tagIDList []int64) ([]*WorkExternalTag, error) {
	return e.repo.GetExternalTagsByID(c, tagIDList)
}

// DeleteExternalUserTags 删除外部联系人标签数据
func (e *ExternalTagService) DeleteExternalUserTags(c ctx.Context, idList []int64) (int64, error) {
	return e.repo.DeleteExternalUserTags(c, idList)
}

// GetCustomerTags 查询外部联系人自定义标签
func (e *ExternalTagService) GetCustomerTags(c ctx.Context, enterpriseID int64, memberID int64, tags []string) ([]*WorkExternalTag, error) {
	return e.repo.QueryCustomerTags(c, enterpriseID, memberID, tags)
}

func (e *ExternalTagService) GetExternalUserUsedTag(c ctx.Context, enterpriseID int64, tagID int64) ([]WorkExternalRelationTag, error) {
	return e.repo.GetExternalUserUsedTag(c, enterpriseID, tagID)
}

// 根据外部联系人用户ID,删除不需要的关系
func (e *ExternalTagService) DeleteUserTagsByUserID(c ctx.Context, enterpriseID int64, externalUserID int64) {
	e.repo.DeleteUserTagsByUserID(c, enterpriseID, externalUserID)
}
