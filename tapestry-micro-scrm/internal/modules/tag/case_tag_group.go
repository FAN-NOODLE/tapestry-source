package tag

import (
	"errors"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/app"
	"micro-scrm/internal/modules/tag/entity"
	workSDK "micro-scrm/internal/pkg/sdk/work"
	"micro-scrm/internal/pkg/workwechat"
	"time"
)

type WorkExternalTagGroupService struct {
	Repo WorkTagGroupRepository
}

// AddTagGroup 新增标签组
// 新增标签组时，1、可以带标签信息，2、也可以不带标签;
// 不带标签时，生成标签组在本地
func (g *WorkExternalTagGroupService) AddTagGroup(c ctx.Context, params entity.CreateManualTagGroupParams) ([]string, error) {
	// 1、检测标签组是否存在
	isExist, existErr := g.Repo.GetTagGroupInfoByName(c, params)
	if existErr != nil && !isExist {
		return nil, existErr
	}
	if isExist && existErr == nil {
		return nil, errors.New("标签组名称已经存在")
	}

	// 只增加标签时，则只在数据库中增加组
	if 0 == len(params.TagNames) {
		return nil, errors.New("新建标签组，必须添加至少一个标签")
	}

	// 检测添加的标签名称是否已经被使用
	isExistTagName, isExistErr := g.Repo.IsExistTagName(c, params.EnterpriseID, params.TagNames)
	if isExistErr != nil {
		return nil, isExistErr
	}

	if isExistTagName {
		return nil, errors.New("添加的标签已经存在，请重新命名")
	}

	// 增加组和标签
	workWeChatObj, workErr := workwechat.NewWorkWeChat(c, params.EnterpriseID)
	if workErr != nil {
		return nil, workErr
	}

	data := workSDK.AddCorpTagRequest{
		GroupName: params.GroupName,
	}
	for _, tagName := range params.TagNames {
		data.TagList = append(data.TagList, workSDK.AddCorpTag{
			Name: tagName,
		})
	}
	response, resErr := workWeChatObj.AddTag(data)
	if resErr != nil {
		return nil, resErr
	}

	// 增加标签组和标签
	tagGroup := entity.WorkExternalTagGroupModel{
		ID:           app.UID(),
		EnterpriseID: params.EnterpriseID,
		GroupID:      response.TagGroup.GroupID,
		GroupName:    response.TagGroup.GroupName,
		Order:        response.TagGroup.Order,
		CreateTime:   response.TagGroup.CreateTime,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}
	_, insErr := g.Repo.AddTagGroup(c, tagGroup)
	if insErr != nil {
		log.Info(c).Msgf("新增手动标签组，同步标签数据失败，错误信息：%+v", insErr)
	}

	var workTagIds []string
	for _, tagInfo := range response.TagGroup.TagList {
		workTagIds = append(workTagIds, tagInfo.ID)
		tagErr := g.Repo.AddTag(c, entity.WorkExternalTagModel{
			ID:           app.UID(),
			EnterpriseID: params.EnterpriseID,
			GroupID:      tagGroup.ID,
			WorkTagID:    tagInfo.ID,
			Name:         tagInfo.Name,
			Type:         1,
			Order:        tagInfo.Order,
			CreateTime:   tagInfo.CreateTime,
			ExpiredTime:  params.ExpiredTime,
			CreatedAt:    time.Now().Unix(),
			UpdatedAt:    time.Now().Unix(),
		})
		if tagErr != nil {
			log.Info(c).Msgf("新增手动标签组，同步标签数据失败，错误信息：%+v", tagErr)
		}
	}

	return workTagIds, nil
}

// DeleteTagGroup 删除标签组
// 企业微信删除标签组的方式为：把标签组中的标签删除完后，该标签组就自动删除了
// 故，删除标签组，只要保证该标签组下没有标签，直接删除库中的标签即可
func (g *WorkExternalTagGroupService) DeleteTagGroup(c ctx.Context, params entity.DeleteManualTagGroupRequestParam) error {
	// 检测标签组是否有外部联系人
	tagGroupInfo, existErr := g.Repo.IsExistExternalUser(c, params)
	if existErr != nil {
		return existErr
	}

	if tagGroupInfo.UseNum > 0 {
		return errors.New("标签组中的某些标签已经标记了外部联系人，不能删除")
	}

	// 删除标签组
	// 1、删除企业微信标签组信息
	workWeChatObj, workErr := workwechat.NewWorkWeChat(c, params.EnterpriseID)
	if workErr != nil {
		return workErr
	}

	delCorpParam := workSDK.DelCorpTagParam{}
	delCorpParam.GroupID = append(delCorpParam.GroupID, tagGroupInfo.WorkWeChatGroupID)
	delErr := workWeChatObj.DeleteTag(delCorpParam)
	if delErr != nil {
		return delErr
	}
	// 2、删除数据库标签
	delParam := entity.DeleteTagGroup{
		EnterpriseID:  params.EnterpriseID,
		GroupID:       tagGroupInfo.GroupId,
		OperationType: entity.GroupOperation,
	}
	// 2.1 删除标签组
	_, groupErr := g.Repo.DeleteTagGroup(c, delParam)
	if groupErr != nil {
		log.Info(c).Msgf("删除手动标签组，同步标签数据失败，错误信息：%+v", groupErr)
	}
	// 2.2 删除标签
	delParam.OperationType = entity.TagOperation
	_, tagErr := g.Repo.DeleteTagGroup(c, delParam)
	if tagErr != nil {
		log.Info(c).Msgf("删除手动标签组，同步标签数据失败，错误信息：%+v", tagErr)
	}

	return nil
}
