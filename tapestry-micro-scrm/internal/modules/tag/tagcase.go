package tag

import (
	"database/sql"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/modules/tag/api/cdp"
	"micro-scrm/internal/modules/tag/entity"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/uid"
)

// Service 应用服务
type Service struct {
	Repo   Repository
	Market MarketingInterface
}

// MarketingInterface 营销策略接口
type MarketingInterface interface {
	IsUserGroupInMarketing(c ctx.Context, groupID int64) (bool, error)
}

// CreateManualTag 创建手动标签逻辑实现
func (s *Service) CreateManualTag(c ctx.Context, request entity.CreateTagRequest) (*entity.CreateTagResponse, error) {
	// 调用cdp接口创建标签
	tagID, err := cdp.CreateManualTag(c, request.TagGroupID, request.TagName)
	if err != nil {
		return nil, err
	}

	var response = entity.CreateTagResponse{
		TagID: tagID,
	}
	return &response, nil
}

// CreateUserGroup 创建用户分组
func (s *Service) CreateUserGroup(c ctx.Context, param entity.CreateUserGroupBody) (*entity.CreateUserGroupResponse, error) {
	createGroupUser, err := cdp.CreateUserGroup(c, param)
	if err != nil {
		return nil, err
	}
	return createGroupUser, nil
}

// EditUserGroup 编辑用户分组
func (s *Service) EditUserGroup(c ctx.Context, param entity.CreateUserGroupBody) (*entity.CreateUserGroupResponse, error) {
	createGroupUser, err := cdp.EditUserGroup(c, param)
	if err != nil {
		return nil, err
	}
	return createGroupUser, nil
}

// DeleteUserGroup 删除用户分组
func (s *Service) DeleteUserGroup(c ctx.Context, param entity.CreateUserGroupBody) (*entity.CreateUserGroupResponse, error) {
	createGroupUser, err := cdp.DeleteUserGroup(c, param)
	if err != nil {
		return nil, err
	}
	return createGroupUser, nil
}

// DeleteUserGroup 删除用户分组
func (s *Service) UserGroupDetail(c ctx.Context, ID int64) (*entity.UserGroupDetail, error) {
	detail, err := cdp.UserGroupDetail(c, ID)
	if err != nil {
		return nil, err
	}
	inMarket, err := s.Market.IsUserGroupInMarketing(c, detail.ID)
	if err != nil {
		return detail, err
	}
	detail.CanEdit = !inMarket
	return detail, nil
}

// UserGroupList 获取用户分组列表
func (s *Service) UserGroupList(c ctx.Context) (*entity.UserGroupList, error) {
	userGroupList, err := cdp.GroupUserList(c)
	if err != nil {
		return nil, err
	}
	//
	//// 取各个模型下的最小时间，作为当前用户分组的更新时间
	//var updatedAt = time.Now().Unix()
	//for _, v := range userGroupList.Rfm {
	//	if v.ValueUpdateAt < updatedAt {
	//		updatedAt = v.ValueUpdateAt
	//	}
	//}
	//
	//for _, v := range userGroupList.Aipl {
	//	if v.ValueUpdateAt < updatedAt {
	//		updatedAt = v.ValueUpdateAt
	//	}
	//}
	//
	//for _, v := range userGroupList.Custom {
	//	if v.ValueUpdateAt < updatedAt {
	//		updatedAt = v.ValueUpdateAt
	//	}
	//}
	//userGroupList.ValueUpdateAt = updatedAt

	return userGroupList, nil
}

// QueryListByIDList 通过id查询标签列表
func (s *Service) QueryListByIDList(c ctx.Context, IDList []int64) ([]*entity.TagEntity, error) {
	list, err := s.Repo.QueryTags(c, IDList)
	if err == sql.ErrNoRows {
		return make([]*entity.TagEntity, 0), nil
	}
	if err != nil {
		return nil, err
	}
	return list, nil
}

// 废弃
func (s *Service) SyncUserGroup(c ctx.Context) error {
	enterpriseID := common.EnterpriseID(c)
	userGroupList, err := cdp.GroupUserList(c)
	if err != nil {
		return err
	}
	var cdpUserGroupMap = make(map[int64]*entity.UserGroupEntity)
	for _, v := range userGroupList.Aipl {
		cdpUserGroupMap[v.GroupID] = &entity.UserGroupEntity{
			UserGroupID: v.GroupID,
			Name:        v.GroupName,
			Type:        v.GroupType,
		}
	}
	for _, v := range userGroupList.Rfm {
		cdpUserGroupMap[v.GroupID] = &entity.UserGroupEntity{
			UserGroupID: v.GroupID,
			Name:        v.GroupName,
			Type:        v.GroupType,
		}
	}
	for _, v := range userGroupList.Custom {
		cdpUserGroupMap[v.GroupID] = &entity.UserGroupEntity{
			UserGroupID: v.GroupID,
			Name:        v.GroupName,
		}
	}
	req := entity.GetTagListRequest{
		EnterpriseID: enterpriseID,
		Page:         1,
		Size:         9999,
	}
	_, dbUserGroupList, err := s.Repo.GetUserGroupListPage(c, &req)

	var deleteUserGroupList []int64
	var updateUserGroupList []*entity.UserGroupEntity
	for _, v := range dbUserGroupList {
		cdpUserGroup, ok := cdpUserGroupMap[v.UserGroupID]
		if !ok {
			deleteUserGroupList = append(deleteUserGroupList, v.UserGroupID)
			continue
		}
		if cdpUserGroup.Name == v.Name {
			delete(cdpUserGroupMap, v.UserGroupID)
			continue
		} else {
			v.Name = cdpUserGroup.Name
			updateUserGroupList = append(updateUserGroupList, v)
			delete(cdpUserGroupMap, v.UserGroupID)
			continue
		}
	}
	var newUserGroupList []*entity.UserGroupEntity
	for _, v := range cdpUserGroupMap {
		v.EnterpriseID = enterpriseID
		v.ID = uid.ID(c)
		v.CreatedAt = int(time.Now().Unix())
		v.UpdatedAt = int(time.Now().Unix())
		newUserGroupList = append(newUserGroupList, v)
	}
	_, err = s.Repo.AddUserGroup(c, newUserGroupList)
	_ = s.DeleteUserGroupFromDB(c, enterpriseID, deleteUserGroupList)
	_ = s.UpdateUserGroupToDB(c, updateUserGroupList)
	return err
}

// UpdateUserGroupToDB 更新数据库中用户分组数据
func (s *Service) UpdateUserGroupToDB(c ctx.Context, userGroupList []*entity.UserGroupEntity) error {
	if len(userGroupList) == 0 {
		return nil
	}
	for _, v := range userGroupList {
		v.UpdatedAt = int(time.Now().Unix())
		_, err := s.Repo.UpdateUserGroup(c, v)
		if err != nil {
			log.Info(c).Msgf("更新用户分组失败：%s, %+v", err.Error(), v)
		}
	}
	return nil
}

// DeleteUserGroupFromDB 从数据库中删除用户分组
func (s *Service) DeleteUserGroupFromDB(c ctx.Context, enterpriseID int64, groupIDList []int64) error {
	if len(groupIDList) == 0 {
		return nil
	}
	_, err := s.Repo.DeleteUserGroup(c, enterpriseID, groupIDList)
	if err != nil {
		log.Info(c).Msgf("删除用户分组失败：%s, %d, %v", err.Error(), enterpriseID, groupIDList)
		return err
	}
	return nil
}

// SyncTagData 同步cdp标签数据
func (s *Service) SyncTagData(c ctx.Context) error {
	enterpriseID := common.EnterpriseID(c)
	groupList, err := cdp.GetTagGroupList(c)
	if err != nil {
		log.Info(c).Msgf("获取cdp标签分组列表")
		return err
	}
	var tagGroupIDList []int64
	var tagGroupMap = make(map[int64]*entity.TagGroup)
	var tagList []*entity.TagInfo
	for i, v := range groupList.AutoList {
		tagGroupIDList = append(tagGroupIDList, v.GroupID)
		tagGroupMap[v.GroupID] = &groupList.AutoList[i]
		tagGroupMap[v.GroupID].Type = entity.GroupTypeAuto
		for _, vv := range v.TagList {
			vv.GroupID = v.GroupID
		}
		if v.TagGroupType != int(entity.TagGroupTypeEnum.Region) {
			tagList = append(tagList, v.TagList...)
		} else {
			tagGroupInfo, err := cdp.TagGroupInfo(c, &entity.TagGroupInfoRequest{
				TagGroupID: v.GroupID,
				Order:      1,
				Page:       1,
				Size:       9999,
			})
			if err != nil {
				log.Info(c).Msgf("获取地区标签数据失败, err:%s", err.Error())
			} else {
				for _, vv := range tagGroupInfo.TagList {
					vv.GroupID = v.GroupID
				}
				tagList = append(tagList, tagGroupInfo.TagList...)
			}
		}
	}
	for i, v := range groupList.ManualList {
		tagGroupIDList = append(tagGroupIDList, v.GroupID)
		tagGroupMap[v.GroupID] = &groupList.ManualList[i]
		tagGroupMap[v.GroupID].Type = entity.GroupTypeManual
		tagList = append(tagList, v.TagList...)
	}

	_ = s.SyncTagGroup(c, enterpriseID, tagGroupMap)
	_ = s.SyncTag(c, enterpriseID, tagList)

	return nil
}

// SyncTag 同步标签数据
func (s *Service) SyncTag(c ctx.Context, enterpriseID int64, tagList []*entity.TagInfo) error {
	var idList []int64
	var tagListMap = make(map[int64]*entity.TagInfo)
	for _, v := range tagList {
		idList = append(idList, v.ID)
		tagListMap[v.ID] = v
	}
	dbTagList, err := s.Repo.QueryEnterpriseTags(c, enterpriseID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err != nil {
		return err
	}
	var deleteTagList []int64
	var updateTagList []*entity.TagEntity
	for _, v := range dbTagList {
		tag, ok := tagListMap[v.TagID]
		if !ok {
			deleteTagList = append(deleteTagList, v.TagID)
			continue
		}
		if v.TagName == tag.Name {
			delete(tagListMap, v.TagID)
			continue
		} else {
			v.TagName = tag.Name
			updateTagList = append(updateTagList, v)
			delete(tagListMap, v.TagID)
			continue
		}
	}
	var newTagList []*entity.TagEntity
	for _, v := range tagListMap {
		newTagList = append(newTagList, &entity.TagEntity{
			ID:           uid.ID(c),
			EnterpriseID: enterpriseID,
			TagID:        v.ID,
			TagName:      v.Name,
			TagGroupID:   v.GroupID,
			CreatedAt:    int(time.Now().Unix()),
			UpdatedAt:    int(time.Now().Unix()),
		})
	}
	if len(newTagList) != 0 {
		_, err := s.Repo.AddTag(c, newTagList)
		if err != nil {
			return err
		}
	}
	_ = s.DeleteTagFromDB(c, enterpriseID, deleteTagList)
	_ = s.UpdateTagToDB(c, updateTagList)
	return nil
}

// UpdateTagToDB 更新数据库中标签信息
func (s *Service) UpdateTagToDB(c ctx.Context, tagList []*entity.TagEntity) error {
	if len(tagList) == 0 {
		return nil
	}
	for _, v := range tagList {
		v.UpdatedAt = int(time.Now().Unix())
		_, err := s.Repo.UpdateTag(c, v)
		if err != nil {
			log.Info(c).Msgf("更新标签数据失败：%s, tag:%+v", err.Error(), v)
		}
	}
	return nil
}

// DeleteTagFromDB 从数据库中删除标签
func (s *Service) DeleteTagFromDB(c ctx.Context, enterpriseID int64, tagIDList []int64) error {
	if len(tagIDList) == 0 {
		return nil
	}
	_, err := s.Repo.DeleteTag(c, enterpriseID, tagIDList)
	if err != nil {
		log.Info(c).Msgf("删除标签失败:%s, enterpriseID:%d, tagIDList:%+v", err.Error(), enterpriseID, tagIDList)
		return err
	}
	return nil
}

// SyncTagGroup 同步标签分组
func (s *Service) SyncTagGroup(c ctx.Context, enterpriseID int64, tagGroupMap map[int64]*entity.TagGroup) error {
	tagGroupList, err := s.Repo.QueryEnterpriseTagGroups(c, enterpriseID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	var deleteTagGroupList []int64
	var updateTagGroupList []*entity.TagGroupEntity
	for _, v := range tagGroupList {
		cdpTagGroup, ok := tagGroupMap[v.GroupID]
		// 数据库中有，cdp返回的数据中没有，从库中删除
		if !ok {
			deleteTagGroupList = append(deleteTagGroupList, v.GroupID)
			continue
		}
		// 数据库中和cdp返回的数据一致，不做处理
		if cdpTagGroup.GroupName == v.GroupName {
			delete(tagGroupMap, v.GroupID)
			continue
		} else {
			v.GroupName = cdpTagGroup.GroupName
			updateTagGroupList = append(updateTagGroupList, v)
			delete(tagGroupMap, v.GroupID)
			continue
		}
	}
	_, _ = s.AddTagGroup(c, enterpriseID, tagGroupMap)
	_ = s.DeleteTagGroupsFromDB(c, enterpriseID, deleteTagGroupList)
	_ = s.UpdateTagGroupToDB(c, updateTagGroupList)
	return nil
}

// UpdateTagGroupToDB 更新数据库中标签组信息
func (s *Service) UpdateTagGroupToDB(c ctx.Context, tagGroupList []*entity.TagGroupEntity) error {
	if len(tagGroupList) == 0 {
		return nil
	}
	for _, v := range tagGroupList {
		v.UpdatedAt = int(time.Now().Unix())
		_, err := s.Repo.UpdateTagGroup(c, v)
		if err != nil {
			log.Info(c).Msgf("更新标签组失败:err:%s, tagGroup:%+v", err.Error(), v)
		}
	}
	return nil
}

// DeleteTagGroupsFromDB 删除标签组
func (s *Service) DeleteTagGroupsFromDB(c ctx.Context, enterpriseID int64, groupIDList []int64) error {
	if len(groupIDList) == 0 {
		return nil
	}
	_, err := s.Repo.DeleteTagGroup(c, enterpriseID, groupIDList)
	if err != nil {
		log.Info(c).Msgf("删除标签组失败:err:%s, enterpriseID:%d, groupIDs:%v", err.Error(), enterpriseID, groupIDList)
		return err
	}
	return nil
}

// AddTagGroup 添加标签分组
func (s *Service) AddTagGroup(c ctx.Context, enterpriseID int64, mapList map[int64]*entity.TagGroup) (int64, error) {
	var newTagGroupList []*entity.TagGroupEntity
	for _, v := range mapList {
		newTagGroupList = append(newTagGroupList, &entity.TagGroupEntity{
			ID:           uid.ID(c),
			EnterpriseID: enterpriseID,
			GroupID:      v.GroupID,
			GroupName:    v.GroupName,
			Type:         v.Type,
			CreatedAt:    int(time.Now().Unix()),
			UpdatedAt:    int(time.Now().Unix()),
		})
	}
	if len(newTagGroupList) == 0 {
		return 0, nil
	}
	count, err := s.Repo.AddTagGroup(c, newTagGroupList)
	if err != nil {
		log.Info(c).Msgf("添加标签分组失败:%s, tagGroup:%+v", err.Error(), newTagGroupList)
		return 0, err
	}
	return count, nil
}

// GetTagList 获取标签列表
func (s *Service) GetTagList(c ctx.Context, req *entity.GetTagListRequest) (int64, []*entity.GetTagListItem, error) {
	req.EnterpriseID = common.EnterpriseID(c)
	total, list, err := s.Repo.GetTagListPage(c, req)
	if err != nil {
		return 0, nil, err
	}
	var data = make([]*entity.GetTagListItem, 0)
	for _, v := range list {
		data = append(data, &entity.GetTagListItem{
			ID:   v.TagID,
			Name: v.TagName,
		})
	}
	return total, data, nil
}

// GetUserGroupList 获取用户分组列表
func (s *Service) GetUserGroupList(c ctx.Context, req *entity.GetTagListRequest) (int64, []*entity.GetTagListItem, error) {
	req.EnterpriseID = common.EnterpriseID(c)
	total, list, err := s.Repo.GetUserGroupListPage(c, req)
	if err != nil {
		return 0, nil, err
	}
	var data = make([]*entity.GetTagListItem, 0)
	for _, v := range list {
		data = append(data, &entity.GetTagListItem{
			ID:   v.UserGroupID,
			Name: v.Name,
		})
	}
	return total, data, nil
}

// QueryTagGroups 通过groupIds 查询标签组
func (s Service) QueryTagGroups(ctx ctx.Context, tagGroupIds []int64) ([]*entity.TagGroupEntity, error) {
	return s.Repo.QueryTagGroups(ctx, tagGroupIds)
}

// QueryTags 通过标签ID 查询标签
func (s Service) QueryTags(ctx ctx.Context, tagIds []int64) ([]*entity.TagEntity, error) {
	return s.Repo.QueryTags(ctx, tagIds)
}

// QueryTags 通过标签ID 查询标签
func (s Service) QueryWorkExternalTags(ctx ctx.Context, tagIds []int64) ([]entity.WorkExternalTag, error) {
	return s.Repo.QueryWorkExternalTagList(ctx, tagIds, 0)
}

// CustomerTags 客户标签
func (s Service) CustomerTags(ctx ctx.Context, enterpriseID int64, relationID int64, tagGroupList *entity.TagGroupList) (result entity.CustomerTagResponse) {
	if len(tagGroupList.ManualList) > 0 {
		// 用户标签
		userTagList, _ := s.Repo.WorkExternalTagList(ctx, enterpriseID, relationID)
		userTagListMap := make(map[string]entity.WorkExternalTag)
		if userTagList != nil && len(*userTagList) > 0 {
			for _, vItem := range *userTagList {
				if vItem.TagID != "" {
					userTagListMap[vItem.TagID] = vItem
				}
			}
		}
		for _, vItem := range tagGroupList.ManualList {
			temp := entity.CustomerTagGroup{
				TagGroupBase: entity.TagGroupBase{
					GroupID:           vItem.GroupID,
					GroupName:         vItem.GroupName,
					WxTagGroupID:      vItem.WxTagGroupID,
					ValueUpdateAt:     vItem.ValueUpdateAt,
					SyncQW:            vItem.SyncQW,
					TagGroupFirstType: vItem.TagGroupFirstType,
					TagGroupType:      vItem.TagGroupType,
					Description:       vItem.Description,
					PageCount:         vItem.PageCount,
					IsDelete:          vItem.IsDelete,
					Type:              vItem.Type,
				},
			}
			for _, subItem := range vItem.TagList {
				tagTemp := &entity.CustomerTag{
					TagInfo:  *subItem,
					IsChoice: 0,
				}
				if _, ok := userTagListMap[subItem.WxTagID]; ok {
					tagTemp.IsChoice = 1
				}
				temp.TagList = append(temp.TagList, tagTemp)
			}
			if len(temp.TagList) == 0 {
				temp.TagList = make([]*entity.CustomerTag, 0)
			}
			result.ManualList = append(result.ManualList, temp)
		}
	}
	return result
}

// GetAllUsedTags
func (s Service) GetAllUsedExternalTags(ctx ctx.Context, enterpriseID int64) ([]int64, error) {
	return s.Repo.GetAllUsedExternalTags(ctx, enterpriseID)
}
