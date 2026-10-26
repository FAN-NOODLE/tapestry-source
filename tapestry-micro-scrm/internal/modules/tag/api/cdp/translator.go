package cdp

import (
	"micro-scrm/internal/modules/tag/entity"
	"micro-scrm/internal/pkg/cdp"
)

// GroupUserListResponseTranslator GroupUserList接口返回数据转换
func GroupUserListResponseTranslator(sdkResponse *cdp.UserGroupListResponse) *entity.UserGroupList {
	if sdkResponse == nil {
		return nil
	}
	var userGroup entity.UserGroupList

	userGroup.ValueUpdateAt = sdkResponse.Data.ValueUpdateAt
	userGroup.Aipl = make([]entity.UserGroupModel, 0)
	userGroup.Rfm = make([]entity.UserGroupModel, 0)
	userGroup.Custom = make([]entity.CustomUserGroupModel, 0)
	for _, v := range sdkResponse.Data.Aipl {
		userGroup.Aipl = append(userGroup.Aipl, entity.UserGroupModel{
			GroupID:       v.GroupID,
			GroupName:     v.GroupName,
			GroupType:     v.GroupType,
			Value:         v.Value,
			ValueUpdateAt: v.ValueUpdateAt,
		})
	}

	for _, v := range sdkResponse.Data.Rfm {
		userGroup.Rfm = append(userGroup.Rfm, entity.UserGroupModel{
			GroupID:       v.GroupID,
			GroupName:     v.GroupName,
			GroupType:     v.GroupType,
			Value:         v.Value,
			ValueUpdateAt: v.ValueUpdateAt,
		})
	}

	for _, v := range sdkResponse.Data.Custom {

		customerModel := entity.CustomUserGroupModel{
			GroupID:       v.GroupID,
			GroupName:     v.GroupName,
			Value:         v.Value,
			ValueUpdateAt: v.ValueUpdateAt,
			NewAggeResult: v.NewAggeResult,
			Conditions:    []entity.UserGroupCondition{},
		}
		for _, vCondition := range v.Conditions {
			condition := entity.UserGroupCondition{
				IncludeType: vCondition.IncludeType,
				TagGroups:   []entity.UserGroupTagInfo{},
			}
			for _, tagGroup := range vCondition.TagGroups {
				if tagGroup.TagIds == nil {
					tagGroup.TagIds = []int64{}
				}
				condition.TagGroups = append(condition.TagGroups, entity.UserGroupTagInfo{
					TagGroupID: tagGroup.TagGroupID,
					TagIds:     tagGroup.TagIds,
				})
			}
			customerModel.Conditions = append(customerModel.Conditions, condition)
		}
		userGroup.Custom = append(userGroup.Custom, customerModel)
	}

	return &userGroup
}

// TagGroupInfoResponseTranslator 标签组详情返回值转换
func TagGroupInfoResponseTranslator(sdkResponse *cdp.TagGroupInfoResponse) *entity.TagGroup {
	sdkData := sdkResponse.Data
	var tagGroup = entity.TagGroup{
		TagGroupBase: entity.TagGroupBase{
			GroupID:           sdkData.GroupID,
			GroupName:         sdkData.GroupName,
			ValueUpdateAt:     sdkData.ValueUpdateAt,
			SyncQW:            sdkData.SyncQW,
			TagGroupFirstType: sdkData.TagGroupFirstType,
			TagGroupType:      sdkData.TagGroupType,
			Description:       sdkData.Description,
			PageCount:         sdkData.PageCount,
		},
		TagList: []*entity.TagInfo{},
	}
	for _, v := range sdkData.TagList {
		var condition map[string]interface{}
		switch v.Condition.(type) {
		case map[string]interface{}:
			condition = v.Condition.(map[string]interface{})
		}
		tagGroup.TagList = append(tagGroup.TagList, &entity.TagInfo{
			ID:            v.ID,
			Name:          v.Name,
			Value:         v.Value,
			NewAggeResult: v.NewAggeResult,
			Condition:     condition,
		})
	}
	return &tagGroup
}

// UserGroupDetailTranslator 用户分组接口返回值转换
func UserGroupDetailTranslator(sdkResponse *cdp.UserGroupDetailResponse) *entity.UserGroupDetail {
	detail := sdkResponse.Data
	var response = entity.UserGroupDetail{
		Name:  detail.Name,
		ID:    detail.ID,
		Value: detail.Value,
	}
	for _, v := range detail.Conditions {
		var cond = entity.UserGroupDetailCondition{
			IncludeType: v.IncludeType,
		}
		cond.TagGroups = make([]*entity.UserGroupDetailTagGroup, 0)
		for _, vv := range v.TagGroups {
			var group = entity.UserGroupDetailTagGroup{
				TagGroupID:   vv.TagGroupID,
				TagGroupName: vv.TagGroupName,
			}
			group.Tags = make([]*entity.UserGroupDetailTag, 0)
			for _, vvv := range vv.Tags {
				group.Tags = append(group.Tags, &entity.UserGroupDetailTag{
					TagID:   vvv.TagID,
					TagName: vvv.TagName,
				})
			}
			cond.TagGroups = append(cond.TagGroups, &group)
		}
		response.Conditions = append(response.Conditions, &cond)
	}
	return &response
}
