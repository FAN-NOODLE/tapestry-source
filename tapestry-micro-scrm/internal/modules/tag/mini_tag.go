package tag

import (
	"micro-scrm/internal/modules/tag/entity"
	"micro-scrm/internal/orm"
)

// 小程序标签逻辑层

type MiniTag struct {
}

// 获取企业标签组详情
func (MiniTag) GetTagGroupDetail(param entity.MiniBaseParam) ([]entity.MiniTagGroupDetail, error) {
	result := []entity.MiniTagGroupDetail{}

	dataList, err := orm.Service.Tag.GetTagGroupDetail(param.EnterpriseID)
	if err != nil {
		return result, err
	}

	// 过滤数据
	for _, tagGroup := range dataList {
		tagGroupInfo := entity.MiniTagGroupDetail{
			ID:        tagGroup.ID,
			GroupID:   tagGroup.GroupID,
			GroupName: tagGroup.GroupName,
			Order:     tagGroup.Order,
			TagList:   []entity.MiniTagDetail{},
		}
		for _, tag := range tagGroup.TagList {
			tagGroupInfo.TagList = append(tagGroupInfo.TagList, entity.MiniTagDetail{
				ID:    tag.ID,
				TagID: tag.TagID,
				Name:  tag.Name,
				Order: tag.Order,
			})
		}

		result = append(result, tagGroupInfo)
	}

	return result, nil
}
