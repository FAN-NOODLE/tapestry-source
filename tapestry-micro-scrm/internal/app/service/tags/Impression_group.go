package tagsSvc

import (
	"micro-scrm/internal/app"
	tagsDao "micro-scrm/internal/app/dao/tags"
	tagsReq "micro-scrm/internal/app/http/request/tags"
	"micro-scrm/internal/app/models"
	tagModels "micro-scrm/internal/app/models/tag"
)

type ImpressionGroup struct {
	impressionGroupDao *tagsDao.ImpressionGroup
}

func NewImpressionGroup() *ImpressionGroup {
	return &ImpressionGroup{
		impressionGroupDao: tagsDao.NewImpressionGroup(),
	}
}

func (s ImpressionGroup) CreateGroup(c app.Ctx, data tagsReq.CreateGroup) (groupId int64, err error) {
	// 判读标签组是否存在，不存在则新增标签组
	tagGroupInfo, err := s.impressionGroupDao.GetInfo(c, tagModels.ImpressionTagGroupModel{
		GroupName: data.GroupName,
		MicroBaseModel: models.MicroBaseModel{
			EnterpriseID: data.EnterpriseID,
		},
	})
	if err != nil {
		return 0, err
	}
	// 添加标签组
	if tagGroupInfo.ID != 0 {
		return tagGroupInfo.ID, nil
	}
	impressionGroupData := tagModels.ImpressionTagGroupModel{
		GroupName:      data.GroupName,
		MicroBaseModel: models.CreateMicroBase(data.EnterpriseID, app.UID()),
	}
	groupId, err = s.impressionGroupDao.Create(c, impressionGroupData)
	if err != nil {
		return 0, err
	}
	return groupId, nil
}
