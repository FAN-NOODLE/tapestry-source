package tag

import (
	"micro-scrm/internal/pkg/sdk/boss"
	"sort"
	"strconv"

	"micro-scrm/internal/modules/tag/api/cdp"
	"micro-scrm/internal/modules/tag/api/microgoods"
	"micro-scrm/internal/modules/tag/api/shop"
	"micro-scrm/internal/modules/tag/entity"
	cdpSdk "micro-scrm/internal/pkg/cdp"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/workwechat"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

// GroupService 应用服务
type GroupService struct {
	Repo Repository
}

// TagGroupList 获取标签组列表
func (g *GroupService) TagGroupList(c ctx.Context, enterpriseID int64, tagType int) (*entity.TagGroupList, error) {
	// 查询企业微信中的标签信息
	workWeChatObj, workWeChatError := workwechat.NewWorkWeChat(c, enterpriseID)
	if workWeChatError != nil {
		return nil, workWeChatError
	}
	tagGroupInfoList, tagGroupErr := workWeChatObj.TagList()
	if tagGroupErr != nil {
		return nil, tagGroupErr
	}

	isThirdTage := 0
	// 获取enterprise配置
	configData, err := boss.EnterpriseConfig(c, map[string]interface{}{
		"enterprise_id": enterpriseID,
		"key":           "only_third_tag",
	})
	if err != nil {
		return nil, tagGroupErr
	}
	if len(configData.Data) > 0 {
		isThirdTage = int(configData.Data[0].Value)
	}
	// 过滤，并对数据添加数据库的信息数据
	// 1、获取企业对应的标签对应组的信息集合
	workTagInfoList, workTagInfoErr := g.Repo.GetTagInfoList(c, enterpriseID, tagType, isThirdTage)
	if workTagInfoErr != nil {
		return nil, workTagInfoErr
	}
	mapping := make(map[string]entity.WorkTagInfoEntity, 0)
	mappingGroup := make(map[string]entity.WorkTagInfoEntity, 0)
	for _, item := range workTagInfoList {
		mapping[item.WorkWeChatTagID] = item
		mappingGroup[item.WorkWeChatGroupID] = item
	}

	// 还原老旧的数据格式 entity.TagGroupList
	var groupList *entity.TagGroupList = &entity.TagGroupList{}
	for _, itemTagGroup := range tagGroupInfoList {
		if _, isExistGroup := mappingGroup[itemTagGroup.GroupID]; !isExistGroup {
			continue
		}

		if itemTagGroup.Delete == true {
			continue
		}
		// 标签组信息赋值
		tagGroup := entity.TagGroup{}
		isSetGroup := false
		// 标签组内标签赋值
		var tagList entity.TagList = make([]*entity.TagInfo, 0)
		for _, itemTag := range itemTagGroup.Tag {
			if _, isExistTag := mapping[itemTag.ID]; !isExistTag {
				continue
			}

			if 0 == tagGroup.IsDelete {
				if mapping[itemTag.ID].WorkExternalTagEntity.UseNum > 0 {
					tagGroup.IsDelete = 1
				}
			}

			if !isSetGroup {
				tagGroup.WxTagGroupID = itemTagGroup.GroupID
				tagGroup.GroupName = itemTagGroup.GroupName
				tagGroup.GroupID = mapping[itemTag.ID].GroupId
				tagGroup.ValueUpdateAt = mapping[itemTag.ID].ValueUpdateAt

				isSetGroup = true
			}

			tagList = append(tagList, &entity.TagInfo{
				WxTagID: itemTag.ID,
				Name:    itemTag.Name,
				ID:      mapping[itemTag.ID].TagId,
				Value:   strconv.Itoa(mapping[itemTag.ID].WorkExternalTagEntity.UseNum),
			})
		}
		tagGroup.TagList = tagList
		groupList.ManualList = append(groupList.ManualList, tagGroup)
	}

	return groupList, nil
}

// CreateGroup 创建标签组
func (g *GroupService) CreateGroup(c ctx.Context, request entity.CreateGroupRequest) (int64, error) {
	groupID, err := cdp.CreateTagGroup(c, request)
	if err != nil {
		return groupID, err
	}
	return groupID, nil
}

// UpdateTagGroup 更新标签组设置
func (g *GroupService) UpdateTagGroup(c ctx.Context, request entity.AutoTagUpdateRequest) error {
	err := cdp.UpdateAutoTag(c, request)
	if err != nil {
		return err
	}
	return nil
}

// UpdateTagGroupSync 更新标签组同步到企业微信配置
func (g *GroupService) UpdateTagGroupSync(c ctx.Context, request entity.UpdateTagGroupSyncRequest) error {
	_, err := cdpSdk.UpdateTagGroupSync(c, request)
	if err != nil {
		return err
	}
	return nil
}

// UserGroupPrecalculateCount 客户分组预计算人数
func (g *GroupService) UserGroupPrecalculateCount(c ctx.Context, request *entity.UserGroupPrecalculateRequest) (*entity.UserGroupPrecalculateResponse, error) {
	//sdkResponse, err := cdpSdk.UserGroupPrecalculateCount(c, request)
	//if err != nil || sdkResponse == nil {
	//	return nil, err
	//}

	var response = entity.UserGroupPrecalculateResponse{
		UserCount: 0,
	}

	return &response, nil
}

// TagGroupInfo 标签组详情
func (g *GroupService) TagGroupInfo(c ctx.Context, request *entity.TagGroupInfoRequest) (*entity.TagGroup, error) {
	tagGroupInfo, err := cdp.TagGroupInfo(c, request)
	if err != nil || tagGroupInfo == nil {
		return nil, err
	}
	// TODO 分类型返回数据类型
	switch entity.TagGroupType(tagGroupInfo.TagGroupType) {
	case entity.TagGroupTypeEnum.AddGoodToCart,
		entity.TagGroupTypeEnum.PayForGood,
		entity.TagGroupTypeEnum.VisitGood:
		err := g.setGoodInfoForTagGroup(c, tagGroupInfo)
		if err != nil {
			return nil, err
		}
	case entity.TagGroupTypeEnum.UseCoupon,
		entity.TagGroupTypeEnum.ReceiveCoupon:
		err := g.setCouponForTagGroup(c, tagGroupInfo)
		if err != nil {
			return nil, err
		}
	case entity.TagGroupTypeEnum.VisitDepth,
		entity.TagGroupTypeEnum.VisitActivation,
		entity.TagGroupTypeEnum.VisitFrequency,
		entity.TagGroupTypeEnum.PayFrequency,
		entity.TagGroupTypeEnum.PerCustomerTransaction,
		entity.TagGroupTypeEnum.GMV,
		entity.TagGroupTypeEnum.RePurchaseCycle,
		entity.TagGroupTypeEnum.VisitDuration:
		if request.Order == 0 {
			sort.Sort(tagGroupInfo.TagList)
		}
	case entity.TagGroupTypeEnum.PurchaseFrequency,
		entity.TagGroupTypeEnum.RefundRate:
		if request.Order == 0 {
			sort.Sort(sort.Reverse(tagGroupInfo.TagList))
		}
	}

	return tagGroupInfo, nil
}

func (g *GroupService) setGoodInfoForTagGroup(c ctx.Context, tagGroup *entity.TagGroup) error {
	var spuIDList []int64
	for i := range tagGroup.TagList {
		spuID, ok := tagGroup.TagList[i].Condition["spu_id"]
		if ok {
			spuIDList = append(spuIDList, int64(spuID.(float64)))
		}
	}
	if len(spuIDList) == 0 {
		return nil
	}
	goodsInfo, err := microgoods.GetSpuListBySpuIds(c, spuIDList, common.EnterpriseID(c))
	if err != nil {
		return nil
	}

	for i := range tagGroup.TagList {
		spuID, ok := tagGroup.TagList[i].Condition["spu_id"]
		if ok {
			tagGroup.TagList[i].GoodData = goodsInfo[int64(spuID.(float64))]
		} else {
			tagGroup.TagList[i].GoodData = &entity.TagGoodInfo{}
		}
	}
	return nil
}

// setCouponForTagGroup 设置标签组中优惠券信息
func (g *GroupService) setCouponForTagGroup(c ctx.Context, tagGroup *entity.TagGroup) error {
	var couponIDList []int64
	for i := range tagGroup.TagList {
		couponID, ok := tagGroup.TagList[i].Condition["coupon_id"]
		if ok {
			couponIDInt, _ := strconv.ParseInt(couponID.(string), 10, 64)
			couponIDList = append(couponIDList, couponIDInt)
		}
	}
	if len(couponIDList) == 0 {
		return nil
	}
	couponInfo, err := shop.GetWxCouponListByIds(c, couponIDList)
	if err != nil {
		return err
	}
	for i := range tagGroup.TagList {
		couponID, ok := tagGroup.TagList[i].Condition["coupon_id"]
		if ok {
			couponIDInt, _ := strconv.ParseInt(couponID.(string), 10, 64)
			tagGroup.TagList[i].CouponData = couponInfo[couponIDInt]
		} else {
			tagGroup.TagList[i].CouponData = &entity.TagCouponInfo{}
		}
	}
	return nil
}

// SortTagGroup 标签组排序
func (g *GroupService) SortTagGroup(c ctx.Context, request *entity.SortTagGroupRequest) error {
	_, err := cdpSdk.SortTagGroup(c, request)
	if err != nil {
		return err
	}
	return nil
}

// SearchTag 标签搜索
func (g *GroupService) SearchTag(c ctx.Context, request *entity.SearchTagRequest) (*entity.TagGroupList, error) {
	tagGroupList, err := cdp.SearchTag(c, request)
	if err != nil || tagGroupList == nil {
		return nil, err
	}
	return tagGroupList, nil
}

// DeleteManualTagGroup 删除手动标签组
func (g *GroupService) DeleteManualTagGroup(c ctx.Context, request *entity.DeleteManualTagGroupRequest) error {
	_, err := cdpSdk.DeleteManualTagGroup(c, request)
	if err != nil {
		return err
	}
	return nil
}

// UpdateManualTag 更新手动标签信息
func (g *GroupService) UpdateManualTag(c ctx.Context, request *entity.UpdateManualTagRequest) error {
	_, err := cdpSdk.UpdateManualTag(c, request)
	if err != nil {
		return err
	}
	return nil
}

// DeleteManualTag 删除手动标签
func (g *GroupService) DeleteManualTag(c ctx.Context, request *entity.DeleteManualTagRequest) error {
	_, err := cdpSdk.DeleteManualTag(c, request)
	if err != nil {
		return err
	}
	return nil
}

// TagEnterpriseInit 初始化企业
func (g *GroupService) TagEnterpriseInit(c ctx.Context, request *entity.TagInitRequest) error {
	_, err := cdpSdk.TagEnterpriseInit(c, request)
	if err != nil {
		return err
	}
	return nil
}
