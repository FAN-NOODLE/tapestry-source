package tag

import (
	"errors"
	"micro-scrm/internal/app"
	tagsDao "micro-scrm/internal/app/dao/tags"
	workDao "micro-scrm/internal/app/dao/work"
	"micro-scrm/internal/pkg/db/gormx"
	"time"

	"gorm.io/gorm"

	"micro-scrm/internal/modules/tag/entity"
	workSDK "micro-scrm/internal/pkg/sdk/work"
	"micro-scrm/internal/pkg/workwechat"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
)

type WorkExternalTagService struct {
	Repo WorkRepository
}

func (work *WorkExternalTagService) isDelete(isDel bool) int {
	if isDel {
		return 0
	}
	return 1
}

// CreateTag 新增标签
func (work *WorkExternalTagService) CreateTag(c ctx.Context, params entity.CreateManualTagParam) ([]string, error) {
	// 1、检测标签组是否存在
	tagGroupInfo, tagGroupErr := work.Repo.GetTagGroupInfoByID(c, params)
	if tagGroupErr != nil {
		return nil, tagGroupErr
	}

	// 2、检测新增的标签名称是否已经使用
	tagIsExist, existErr := work.Repo.IsExistTagByName(c, entity.UpdateManualTagParam{
		EnterpriseID: params.EnterpriseID,
		NewTagName:   params.TagName,
	})
	if existErr != nil {
		return nil, existErr
	}

	if tagIsExist {
		return nil, errors.New("标签名称已经存在")
	}

	// 3、在企业微信中新增标签
	workWeChatObj, workErr := workwechat.NewWorkWeChat(c, params.EnterpriseID)
	if workErr != nil {
		return nil, workErr
	}

	request := workSDK.AddCorpTagRequest{
		GroupID: tagGroupInfo.GroupID,
	}
	request.TagList = append(request.TagList, workSDK.AddCorpTag{
		Name: params.TagName,
	})
	response, resErr := workWeChatObj.AddTag(request)

	if resErr != nil {
		return nil, resErr
	}

	var workTagIds []string
	// 4、同步数据库标签work_external_tag
	for _, tagInfo := range response.TagGroup.TagList {
		workTagIds = append(workTagIds, tagInfo.ID)
		addNum, addErr := work.Repo.AddTag(c, entity.WorkExternalTagModel{
			ID:           app.UID(),
			EnterpriseID: params.EnterpriseID,
			WorkTagID:    tagInfo.ID,
			Name:         params.TagName,
			Type:         1,
			Order:        tagInfo.Order,
			CreateTime:   tagInfo.CreateTime,
			GroupID:      tagGroupInfo.ID,
			ExpiredTime:  params.ExpiredTime,
			CreatedAt:    time.Now().Unix(),
			UpdatedAt:    time.Now().Unix(),
		})
		if addErr != nil || 0 == addNum {
			log.Info(c).Msgf("新增手动标签，同步标签数据失败，错误信息：%+v", addErr)
		}
	}

	return workTagIds, nil
}

// SearchTag 标签搜索
func (work *WorkExternalTagService) SearchTag(c ctx.Context, request entity.SearchTagRequest, enterpriseID int64) (*entity.TagGroupList, error) {
	// isThirdTage := 0
	// // 获取enterprise配置
	// configData, err := boss.EnterpriseConfig(c, map[string]interface{}{
	// 	"enterprise_id": enterpriseID,
	// 	"key":           "only_third_tag",
	// })
	// if err != nil {
	// 	return nil, err
	// }
	// if len(configData.Data) > 0 {
	// 	isThirdTage = int(configData.Data[0].Value)
	// }

	// 查询数据
	// dataList, dataErr := work.Repo.MatchTagGroupByGroupName(c, entity.SearchTagParam{
	//	EnterpriseID: enterpriseID,
	//	Name:         request.Keyword,
	//	SearchType:   int8(request.SearchType),
	//	//IsThirdTage:  isThirdTage,
	// })
	// if dataErr != nil {
	//	return nil, dataErr
	// }

	var dataList []entity.WorkExternalTagGroupEntity
	query := gormx.GetInstance().WithContext(c).Table("work_external_tag_group AS TG").
		Joins("LEFT JOIN work_external_tag T ON T.group_id = TG.id AND TG.deleted_at = 0").
		Select("TG.group_id AS work_group_id,TG.group_name,TG.id AS group_id,TG.updated_at AS value_update_at").
		Where("TG.enterprise_id", enterpriseID).
		Where("T.deleted_at = 0").
		Where("T.type = 1").
		Scopes(func(db *gorm.DB) *gorm.DB {
			if request.SearchType == 1 {
				db.Where("TG.group_name like ?", "%"+request.Keyword+"%")
			}
			if request.SearchType == 2 {
				db.Where("T.name like ?", "%"+request.Keyword+"%")
			}
			// if isThirdTage == 1 {
			// 	db.Where("TG.is_third_tag = ? OR TG.group_name = ?", isThirdTage, "任务标签")
			// }
			return db
		})

	query.Group("TG.id").Order("TG.order desc,TG.id desc").Scan(&dataList)

	// 封装数据为entity.TagGroupList
	res := entity.TagGroupList{
		AutoList:   []entity.TagGroup{},
		ManualList: []entity.TagGroup{},
	}

	for _, data := range dataList {
		var tagGroup entity.TagGroup
		tagGroup.GroupID = data.GroupId
		tagGroup.GroupName = data.GroupName
		tagGroup.WxTagGroupID = data.WorkWeChatGroupID
		tagGroup.IsDelete = work.isDelete(data.UseNum <= 0)
		tagGroup.ValueUpdateAt = data.ValueUpdateAt
		res.ManualList = append(res.ManualList, tagGroup)
	}

	return &res, nil
}

// QueryTagListByGroupId 根据标签组ID，获取组中标签
func (work *WorkExternalTagService) QueryTagListByGroupId(c ctx.Context, params entity.TagParam) ([]entity.WorkExternalTagEntity, int64, error) {
	// tagList, err := work.Repo.QueryTagListByGroupId(c, params)
	// if err != nil {
	//	return nil, err
	// }
	// return tagList, err

	var (
		total int64
		resp  []entity.WorkExternalTagEntity
		err   error
	)

	tagGroupIds := []int64{params.TagGroupID}

	if len(params.TagGroupIds) > 0 {
		tagGroupIds = params.TagGroupIds
	}

	query := gormx.GetInstance().WithContext(c).Table("work_external_tag").
		Select("name as tag_name,id as tag_id,tag_id as work_tag_id,create_time,expired_time,group_id as tag_group_id").
		Where("enterprise_id", params.EnterpriseID).
		Where("group_id in (?)", tagGroupIds).
		Where("deleted_at = 0").
		Where("type = 1").
		Scopes(func(db *gorm.DB) *gorm.DB {
			if params.TagKeywords != "" {
				db.Where("name like ?", "%"+params.TagKeywords+"%")
			}
			return db
		})
	query.Order("create_time DESC")
	query.Count(&total)

	offset := (params.Page - 1) * params.Size
	err = query.Limit(params.Size).Offset(offset).Scan(&resp).Error

	return resp, total, err
}

// GetTagGroupListByEnterpriseId 根据企业ID，获取标签组信息。在本地
func (work *WorkExternalTagService) GetTagGroupListByEnterpriseId(c ctx.Context, enterpriseId int64) ([]entity.WorkExternalTagGroupEntity, error) {
	// 删除企业信息中的企业标签

	// 删数据表 work_external_tag_group, work_external_tag
	tagGroupList, err := work.Repo.GetTagGroupListByEnterpriseId(c, enterpriseId)
	if err != nil {
		return nil, err
	}

	return tagGroupList, nil
}

// DeleteTag 根据标签ID，删除标签信息
func (work *WorkExternalTagService) DeleteTag(c ctx.Context, params entity.TagOperationParam) error {
	// 判断当前标签是否过期
	workExternalTag, err := tagsDao.NewWorkExternalTag().GetOneById(c, params.EnterpriseID, params.TagId)
	if err != nil {
		return err
	}

	if workExternalTag.ID == 0 {
		return errors.New("标签信息未找到")
	}

	taskTagRelationRes, err := workDao.NewPlanTasksTagRelation().GetOneByTagId(c, params.EnterpriseID, params.TagId)
	if err != nil {
		return err
	}

	if taskTagRelationRes.EndTime > int(time.Now().Unix()) {
		return errors.New("该标签在任务关联中,任务未到期无法删除")
	}

	// 删除企业微信标签
	workWeChatObj, workErr := workwechat.NewWorkWeChat(c, params.EnterpriseID)
	if workErr != nil {
		return workErr
	}
	err = workWeChatObj.DeleteTag(workSDK.DelCorpTagParam{TagID: []string{workExternalTag.TagID}})
	if err != nil {
		return err
	}

	// 删除数据库标签
	_, err = work.Repo.DeleteTag(c, params)
	if err != nil {
		log.Error(c).Msgf("删除手动标签，同步标签数据失败，错误信息：%+v", err)
		return err
	}

	// 删除外部联系人标签关系
	err = tagsDao.NewWorkExternalUserTag().DeleteData(c, map[string]any{
		"enterprise_id":   params.EnterpriseID,
		"external_tag_id": params.TagId,
	})
	if err != nil {
		log.Error(c).Msgf("删除外部联系人标签关系失败，错误信息：%+v", err)
		return err
	}

	// 获取当前标签的标签组的标签
	groupTags, err := tagsDao.NewWorkExternalTag().GetMany(c, "enterprise_id = ? and group_id = ?", params.EnterpriseID, workExternalTag.GroupID)
	if err != nil {
		log.Error(c).Msgf("获取当前标签的标签组的标签失败，错误信息：%+v", err)
		return err
	}

	// 如果标签组下不存在标签则删除标签组(企微:标签组下所有的标签均被删除，则标签组会被自动删除)
	if len(groupTags) == 0 {
		err = tagsDao.NewWorkExternalTagGroup().DeleteData(c, map[string]any{
			"enterprise_id": params.EnterpriseID,
			"id":            workExternalTag.GroupID,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// UpdateManualTag 更新手动标签
func (work *WorkExternalTagService) UpdateManualTag(c ctx.Context, param entity.UpdateManualTagParam) error {
	// 1、根据标签ID，获取标签信息;并检测标签是否存在
	tagInfo, tagErr := work.Repo.GetTagInfoByID(c, param)
	if tagErr != nil {
		return tagErr
	}
	if tagInfo == nil {
		return errors.New("标签不存在")
	}

	// 2、检测新的标签名称是否被使用
	isExist, existErr := work.Repo.IsExistTagByName(c, param)
	if existErr != nil {
		return existErr
	}
	if isExist {
		return errors.New("标签名称已经被使用")
	}

	// 3、更新企业微信标签
	workWeChatObj, workErr := workwechat.NewWorkWeChat(c, param.EnterpriseID)
	if workErr != nil {
		return workErr
	}
	updateErr := workWeChatObj.UpdateTagName(workSDK.UpdateTagRequest{
		ID:   tagInfo.WorkTagID,
		Name: param.NewTagName,
	})
	if updateErr != nil {
		return updateErr
	}

	// 4、同步数据库标签
	go func() {
		updateErr := work.Repo.UpdateTagNameByID(c, param)
		if updateErr != nil {
			log.Info(c).Msgf("更新手动标签名称，同步标签数据失败，错误信息：%+v", updateErr)
		}
	}()

	return nil
}
