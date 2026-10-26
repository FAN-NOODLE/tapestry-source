package customerSvc

import (
	"encoding/json"
	"errors"
	"fmt"
	"micro-scrm/internal/app"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	tagsDao "micro-scrm/internal/app/dao/tags"
	customerReq "micro-scrm/internal/app/http/request/customer"
	"micro-scrm/internal/app/models"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	tagModels "micro-scrm/internal/app/models/tag"
	"micro-scrm/internal/app/service"
	"micro-scrm/internal/pkg/helper"
)

type ImpressionTag struct {
	service.Base
	tagRelationDao *tagsDao.ImpressionRelation
}

func NewImpressionTag() *ImpressionTag {
	return &ImpressionTag{
		tagRelationDao: tagsDao.NewImpressionRelation(),
	}
}

func (s ImpressionTag) Update(c app.Ctx, request customerReq.CreateTags) error {
	if len(request.AddTagIds) > 0 {
		// 添加标签关系
		if err := s.createUserTag(c, request); err != nil {
			s.Log().Error(c).Msgf("添加标签关系,错误%v", err)
			return err
		}
	}

	if len(request.DelTagIds) > 0 {
		// 删除标签关系
		if err := s.deleteUserTag(c, request); err != nil {
			s.Log().Error(c).Msgf("删除标签关系,错误%v", err)
			return err
		}
	}

	return nil
}

// 创建用户印象标签
func (s ImpressionTag) createUserTag(c app.Ctx, data customerReq.CreateTags) error {
	existedTagIds, err := s.tagRelationDao.GetExistedTagIdsByUserIdAndTagId(c, data.EnterpriseID, data.WorkExternalUserId, data.AddTagIds, data.OperateUserId)
	if err != nil {
		return err
	}

	// 过滤已存在关系的标签
	addTagIds := helper.DelItems[int64](data.AddTagIds, existedTagIds)
	if len(addTagIds) == 0 {
		return errors.New("标签关系已存在")
	}

	var createRelationData []tagModels.ImpressionTagRelationModel
	for _, tag := range addTagIds {
		// 添加用户标签关联关系
		relationData := tagModels.ImpressionTagRelationModel{
			ExternalUserID: data.WorkExternalUserId,
			TagID:          tag,
			UserID:         data.OperateUserId,
			MicroBaseModel: models.CreateMicroBase(data.EnterpriseID, app.UIDContext(c)),
		}
		createRelationData = append(createRelationData, relationData)
	}

	if err = s.tagRelationDao.BatchCreate(c, createRelationData); err != nil {
		return err
	}

	// err = s.syncUserTags(c, data.WorkExternalUserId, data.EnterpriseID, addTagIds, 1)
	// if err != nil {
	//	s.Log().Error(c).Msgf("同步用户标签,错误:%v", err)
	//	return err
	// }

	// 查询关联的标签信息
	tagListInfo, err := tagsDao.NewImpression().GetListByTagId(c, addTagIds)
	if err != nil {
		s.Log().Error(c).Msgf("查询关联的标签信息,错误:%v", err)
		return err
	}

	var thirdTagIds []string
	for _, tag := range tagListInfo {
		thirdTagIds = append(thirdTagIds, tag.ThirdTagID)
	}

	req := customerReq.CreateTagOperateLogReq{
		EnterpriseID:       data.EnterpriseID,
		WorkExternalUserId: data.WorkExternalUserId,
		ThirdTagIds:        thirdTagIds,
		OperateUserId:      data.OperateUserId,
		OperateName:        data.OperateName,
		OperateMobile:      data.OperateMobile,
		ActionType:         tagModels.ImpressionTagActionTypeAdd,
	}
	return s.createTagOperateLog(c, req)
}

func (s ImpressionTag) deleteUserTag(c app.Ctx, data customerReq.CreateTags) error {
	// 删除用户标签关系
	if err := s.tagRelationDao.RemoveByExternalUserIdTagIds(c, data.EnterpriseID, data.WorkExternalUserId, data.DelTagIds); err != nil {
		return err
	}

	// err := s.syncUserTags(c, data.WorkExternalUserId, data.EnterpriseID, data.DelTagIds, 2)
	// if err != nil {
	//	s.Log().Error(c).Msgf("同步用户标签,错误:%v", err)
	//	return err
	// }

	// 查询关联的标签信息
	tagListInfo, err := tagsDao.NewImpression().GetListByTagId(c, data.DelTagIds)
	if err != nil {
		s.Log().Error(c).Msgf("查询关联的标签信息,错误:%v", err)
		return err
	}

	var thirdTagIds []string
	for _, tag := range tagListInfo {
		thirdTagIds = append(thirdTagIds, tag.ThirdTagID)
	}

	req := customerReq.CreateTagOperateLogReq{
		EnterpriseID:       data.EnterpriseID,
		WorkExternalUserId: data.WorkExternalUserId,
		ThirdTagIds:        thirdTagIds,
		OperateUserId:      data.OperateUserId,
		OperateName:        data.OperateName,
		OperateMobile:      data.OperateMobile,
		ActionType:         tagModels.ImpressionTagActionTypeDel,
	}
	// 添加操作记录
	return s.createTagOperateLog(c, req)
}

// CreateTagOperateLog 创建印象标签操作记录
func (s ImpressionTag) createTagOperateLog(c app.Ctx, data customerReq.CreateTagOperateLogReq) error {
	// 获取用户集合数据
	condParam := fmt.Sprintf("id = %d", data.WorkExternalUserId)
	wxUserCorpInfo, err := microScrmDao.NewWxUserCorpExternal().GetOne(c, condParam, data.EnterpriseID)
	if err != nil {
		s.Log().Error(c).Msgf("获取用户数据,错误:%v", err)
		return err
	}

	// 组装操作记录插入数据
	var createTagOperateLogData []tagModels.ImpressionTagOperateLogModel
	for _, thirdTagId := range data.ThirdTagIds {
		tagOperateLogData := tagModels.ImpressionTagOperateLogModel{
			ExternalUsername: wxUserCorpInfo.Name,
			ExternalUnionId:  wxUserCorpInfo.Unionid,
			ThirdTagID:       thirdTagId,
			OperateUserId:    data.OperateUserId,
			OperateName:      data.OperateName,
			ActionType:       data.ActionType,
			MicroBaseModel:   models.CreateMicroBase(data.EnterpriseID, app.UIDContext(c)),
		}
		createTagOperateLogData = append(createTagOperateLogData, tagOperateLogData)
		// 删除同一个BA、客户、标签的历史标签操作记录
		err = tagsDao.NewTagsOperateLog().DeleteBeforeRecord(c, data.EnterpriseID, tagOperateLogData)
	}

	// 添加标签操作记录
	err = tagsDao.NewTagsOperateLog().BatchCreate(c, createTagOperateLogData)
	if err != nil {
		s.Log().Error(c).Msgf("添加标签操作记录,错误:%v", err)
		return err
	}

	return nil
}

// 同步用户标签
func (s ImpressionTag) syncUserTags(c app.Ctx, enterpriseUserID, enterpriseId int64, tagIds []int64, operateType int) error {
	// 获取用户集合数据
	condParam := microScrmModel.EnterpriseUserModel{ID: enterpriseUserID}
	enterpriseUserInfo, err := microScrmDao.NewEnterpriseUser().GetInfo(c, condParam, enterpriseId)
	if err != nil {
		s.Log().Error(c).Msgf("获取用户集合数据,错误:%v", err)
		return err
	}

	// operateType 1 添加 2 删除
	var updateEnterpriseUserData microScrmModel.EnterpriseUserModel
	if operateType == 1 {
		updateEnterpriseUserData.Tags = helper.MergeTags(enterpriseUserInfo.Tags, tagIds)
	}

	if operateType == 2 {
		var oldTagIds []int64
		_ = json.Unmarshal([]byte(enterpriseUserInfo.Tags), &oldTagIds)
		newTagIds := helper.DelItems[int64](oldTagIds, tagIds)
		newTagIdStr, _ := json.Marshal(newTagIds)
		updateEnterpriseUserData.Tags = string(newTagIdStr)
	}

	// 同步标签id到聚合表中
	err = microScrmDao.NewEnterpriseUser().UpdateById(c, updateEnterpriseUserData, enterpriseId, enterpriseUserInfo.ID)
	if err != nil {
		s.Log().Error(c).Msgf("同步标签id到聚合表中,错误:%v", err)
		return err
	}

	return nil
}
