package source

import (
	"errors"
	"fmt"
	"micro-scrm/internal/app"
	membersDao "micro-scrm/internal/app/dao/members"
	microBossDao "micro-scrm/internal/app/dao/micro_boss"
	sourceDao "micro-scrm/internal/app/dao/source"
	sourceRequest "micro-scrm/internal/app/http/request/source"
	sourceResp "micro-scrm/internal/app/http/response/source"
	"micro-scrm/internal/app/models"
	microBossModel "micro-scrm/internal/app/models/micro_boss"
	"micro-scrm/internal/app/models/source"
	"micro-scrm/internal/app/service"
	"micro-scrm/internal/config"
	"micro-scrm/internal/modules/enterprise"
	sourceModules "micro-scrm/internal/modules/source"
	"micro-scrm/internal/orm"
	bossModel "micro-scrm/internal/orm/model/boss"
	"micro-scrm/internal/pkg/assist"
	"micro-scrm/internal/pkg/errs"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sdk/boss"
	"micro-scrm/internal/pkg/sdk/work"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
)

type Service struct {
	service.Base
	categoryRelationDao *sourceDao.CategoryRelationDao
	sourceDao           *sourceDao.Dao
	sourceGroupDao      *sourceDao.SourceGroupDao
}

func NewService() Service {
	return Service{
		categoryRelationDao: sourceDao.NewCategoryRelationDao(),
		sourceDao:           sourceDao.NewSourceDao(),
		sourceGroupDao:      sourceDao.NewSourceGroupDao(),
	}
}

func (s Service) Update(c app.Ctx, request sourceRequest.Request) error {
	updateIds := make([]int64, 0)
	if request.SourceType == source.TypeImg || request.SourceType == source.TypeFriend {
		updateIds = request.Ids
	} else {
		updateIds = append(updateIds, int64(request.ID))
	}
	if request.StartTime > 0 || request.EndTime > 0 {
		// 更新开始结束时间
		query := map[string]interface{}{
			"id":            updateIds,
			"enterprise_id": request.EnterpriseID,
			"type":          request.SourceType,
		}
		column := map[string]interface{}{
			"start_time": request.StartTime,
			"end_time":   request.EndTime,
			"group_id":   request.GroupId,
		}
		err := s.sourceDao.Update(c, query, column)
		if err != nil {
			return err
		}
	}

	// 产品分类是否需要修改
	if len(request.CategoryIds) > 0 {
		tx := app.DB().WithContext(c).Begin()
		s.categoryRelationDao.WithTx(tx)
		// 删除产品分类和素材关系
		relationQuery := map[string]interface{}{
			"source_id":     updateIds,
			"enterprise_id": request.EnterpriseID,
			"type":          request.SourceType,
		}
		if err := s.categoryRelationDao.ForceDelete(c, relationQuery); err != nil {
			tx.Rollback()
			return err
		}

		var categoryRelationData []source.CategoryRelationModel

		for _, categoryId := range request.CategoryIds {
			for _, id := range updateIds {
				createData := source.CategoryRelationModel{
					MicroBaseModel: models.MicroBaseModel{
						BaseModel: models.BaseModel{
							ID: app.UID(),
						},
						EnterpriseID: int64(int(request.EnterpriseID)),
					},
					CategoryID: categoryId,
					SourceID:   uint64(id),
					Type:       request.SourceType,
				}
				categoryRelationData = append(categoryRelationData, createData)
			}

		}
		// 关联新的产品分类和素材关系
		if err := s.categoryRelationDao.BatchCreate(c, categoryRelationData); err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()
	}

	if len(request.CategoryIds) == 0 { // 如果分类标签给空数组，则意味着全部删除
		// 删除产品分类和素材关系
		relationQuery := map[string]interface{}{
			"source_id":     updateIds,
			"enterprise_id": request.EnterpriseID,
			"type":          request.SourceType,
		}
		if err := s.categoryRelationDao.ForceDelete(c, relationQuery); err != nil {
			return err
		}
	}
	return nil
}

func (s Service) UpdateStatus(c app.Ctx, request sourceRequest.UpdateStatusRequest) error {

	// 更新启用或禁用状态
	query := map[string]interface{}{
		"enterprise_id": request.EnterpriseID,
		"type":          request.SourceType,
	}
	nowTime := time.Now().Unix()
	column := map[string]interface{}{
		"status":     request.SourceStatus,
		"updated_at": uint(nowTime),
	}
	err := s.sourceDao.UpdateStatus(c, query, column, request.IDs)
	if err != nil {
		return err
	}
	return nil

}

func (s Service) GetChild(c app.Ctx, requestChild sourceRequest.RequestChildList) (error, []int64) {

	groupIDs := []int64{requestChild.ParentID}

	err, groupList := s.sourceGroupDao.GetChildIDs(c, requestChild.EnterpriseID, groupIDs)
	if err != nil {
		return err, nil
	}

	var groupIDList []int64
	groupIDList = append(groupIDList, s.getChildAllGroupID(c, requestChild.EnterpriseID, groupList, groupIDs)...)
	if requestChild.ParentID == 0 {
		groupIDList = append(groupIDList, -1)
	}
	if len(groupIDList) == 0 {
		return err, nil
	}

	return nil, groupIDList

}

func (s Service) getChildAllGroupID(c app.Ctx, enterpriseID int64, groupInfo []source.SourceGroupModel, groupIDLists []int64) (nodeIDLists []int64) {
	nodeIDLists = append(nodeIDLists, groupIDLists...)
	var newGroupIds []int64
	for _, item := range groupInfo {
		newGroupIds = append(newGroupIds, item.ID)
		nodeIDLists = append(nodeIDLists, item.ID)
	}
	_, groupList := s.sourceGroupDao.GetChildIDs(c, enterpriseID, newGroupIds)
	if len(groupList) == 0 {
		return
	}
	return append(nodeIDLists, s.getChildAllGroupID(c, enterpriseID, groupList, nodeIDLists)...)
}

// TextConsoleList 查看文本素材列表
func (s Service) TextConsoleList(c app.Ctx, request sourceRequest.RequestConsoleList) (error, []source.ResponseConsoleTextList, int64) {
	// 七天前的起止时间戳
	sevenDaysStartTime, _ := helper.YesterdayPeriodOfTimestamp(7)
	topNum, err := s.sourceDao.GetTopNum(c, request.EnterpriseID, request.Type, request.GroupIDList)
	if err != nil {
		return err, nil, 0
	}
	err, result, total := s.sourceGroupDao.GetTextConsoleList(c, request.EnterpriseID, request.GroupIDList, request.CategoryIds, request.Keyword, request.SourceStatus, request.JustMe, request.AdminID, request.IsUse, request.CreatedStart, request.CreatedEnd, request.Order, request.Sort, request.Page, request.Size)

	if total > 0 {
		var adminIds []int64
		for k, v := range result {
			// 七天内上新素材
			if int64(v.CreatedAt) > sevenDaysStartTime {
				v.IsNew = 1
			}
			if v.AdminID > 0 {
				adminIds = append(adminIds, v.AdminID)
			}
			if topNum < 2 {
				result[k].ResponseConsoleList.CanTop = true
			}
		}

		adminList, err := microBossDao.NewEnterpriseAdmin().GetMany(c, map[string]any{"id": adminIds, "enterprise_id": request.EnterpriseID})
		if err != nil {
			return err, result, total
		}

		adminNameMap := make(map[int64]string)
		for _, admin := range adminList {
			adminNameMap[admin.ID] = admin.Username
		}

		for k, v := range result {
			if adminName, ok := adminNameMap[v.AdminID]; ok {
				result[k].AdminName = adminName
			}
		}
	}

	return err, result, total
}

func (s Service) ImageConsoleList(c app.Ctx, request sourceRequest.RequestConsoleList) (error, []sourceResp.ImageList, int64) {
	// 七天前的起止时间戳
	sevenDaysStartTime, _ := helper.YesterdayPeriodOfTimestamp(7)
	topNum, err := s.sourceDao.GetTopNum(c, request.EnterpriseID, request.Type, request.GroupIDList)
	if err != nil {
		return err, nil, 0
	}
	err, res, total := s.sourceGroupDao.GetImageConsoleList(c, request.EnterpriseID, request.GroupIDList, request.CategoryIds, request.Keyword, request.SourceStatus, request.JustMe, request.AdminID, request.IsUse, request.CreatedStart, request.CreatedEnd, request.Order, request.Sort, request.ProductName, request.SkuItem, request.SkuColor, request.GoodsRelation, request.Page, request.Size)

	resp := make([]sourceResp.ImageList, 0)
	for _, list := range res {
		resp = append(resp, sourceResp.ImageList{
			ResponseConsoleList: list.ResponseConsoleList,
			Name:                list.Name,
			URL:                 list.URL,
			MediaID:             list.MediaID,
			MediaExpireTime:     list.MediaExpireTime,
			MediaStatus:         list.MediaStatus,
			MergeId:             list.MergeId,
			MergeImage:          make([]sourceResp.MergeImage, 0),
		})
	}
	var goodsSourceIds []int64

	if total > 0 {
		var adminIds []int64
		var mergeIds []int64

		for _, v := range resp {
			// 七天内上新素材
			if int64(v.CreatedAt) > sevenDaysStartTime {
				v.IsNew = 1
			}
			if v.AdminID > 0 {
				adminIds = append(adminIds, v.AdminID)
			}
			if v.MergeId > 0 {
				mergeIds = append(mergeIds, v.MergeId)
			}
			goodsSourceIds = append(goodsSourceIds, v.ID)
		}

		adminNameMap := make(map[int64]string)
		// 查询管理员昵称
		if len(adminIds) > 0 {
			adminList, err := microBossDao.NewEnterpriseAdmin().GetMany(c, map[string]any{"id": adminIds, "enterprise_id": request.EnterpriseID})
			if err != nil {
				return err, resp, total
			}

			for _, admin := range adminList {
				adminNameMap[admin.ID] = admin.Username
			}
		}
		sourceMap := make(map[int64][]sourceResp.MergeImage)
		// 查询组合图
		if len(mergeIds) > 0 {
			sourceList, err := sourceDao.NewSourceDao().GetMany(c, map[string]any{"enterprise_id": request.EnterpriseID, "merge_id": mergeIds})
			if err != nil {
				return err, resp, total
			}

			var sourceIds []int64
			for _, val := range sourceList {
				sourceIds = append(sourceIds, val.ID)
			}

			sourceImageMap := make(map[int64]source.ImageModel)
			sourceImageList, err := sourceDao.NewImageDao().GetMany(c, request.EnterpriseID, map[string]any{"id": sourceIds})
			if err != nil {
				return err, resp, total
			}

			for _, image := range sourceImageList {
				sourceImageMap[image.ID] = image
			}

			for _, val := range sourceList {
				tmpImage := sourceResp.MergeImage{
					Id:   sourceImageMap[val.ID].ID,
					Name: sourceImageMap[val.ID].Name,
					URL:  sourceImageMap[val.ID].URL,
					Sort: val.Sort,
				}
				goodsSourceIds = append(goodsSourceIds, sourceImageMap[val.ID].ID)
				sourceMap[val.MergeID] = append(sourceMap[val.MergeID], tmpImage)
			}
		}
		//获取素材关联商品sku
		sourceSkuMap, _ := s.GetSourceRelationGoods(c, request.EnterpriseID, helper.ArrayUniq(goodsSourceIds))

		for k, v := range resp {
			if adminName, ok := adminNameMap[v.AdminID]; ok {
				resp[k].AdminName = adminName
			}

			// 组装合并图片
			if sourceImage, ok := sourceMap[v.MergeId]; ok {
				sort.Slice(sourceImage, func(i, j int) bool {
					return sourceImage[i].Sort < sourceImage[j].Sort
				})
				resp[k].MergeImage = append(resp[k].MergeImage, sourceImage...)
			}
			if topNum < 2 {
				resp[k].CanTop = true
			}
			//将素材关联商品sku写入结果
			var tmpSkuList = make([]string, 0)
			if _, ok := sourceSkuMap[v.ID]; ok {
				tmpSkuList = sourceSkuMap[v.ID]
			}
			for _, image := range resp[k].MergeImage {
				if _, ok := sourceSkuMap[image.Id]; ok {
					tmpSkuList = append(tmpSkuList, sourceSkuMap[image.Id]...)
				}
			}
			tmpSkuList = helper.ArrayUniq(tmpSkuList)
			resp[k].SkuList = tmpSkuList
		}
	}

	return err, resp, total
}

func (s Service) LinkConsoleList(c app.Ctx, request sourceRequest.RequestConsoleList) (error, []source.ResponseConsoleLinkList, int64) {
	// 七天前的起止时间戳
	sevenDaysStartTime, _ := helper.YesterdayPeriodOfTimestamp(7)
	topNum, err := s.sourceDao.GetTopNum(c, request.EnterpriseID, request.Type, request.GroupIDList)
	if err != nil {
		return err, nil, 0
	}
	err, result, total := s.sourceGroupDao.GetLinkConsoleList(c, request.EnterpriseID, request.GroupIDList, request.CategoryIds, request.Keyword, request.SourceStatus, request.JustMe, request.AdminID, request.IsUse, request.CreatedStart, request.CreatedEnd, request.Order, request.Sort, request.ProductName, request.SkuItem, request.SkuColor, request.GoodsRelation, request.SourceOrigin, request.Page, request.Size)

	if total > 0 {
		var adminIds []int64
		var sourceIds []int64
		for _, v := range result {
			// 七天内上新素材
			if int64(v.CreatedAt) > sevenDaysStartTime {
				v.IsNew = 1
			}
			if v.AdminID > 0 {
				adminIds = append(adminIds, v.AdminID)
			}
			sourceIds = append(sourceIds, v.ID)
		}

		adminList, err := microBossDao.NewEnterpriseAdmin().GetMany(c, map[string]any{"id": adminIds, "enterprise_id": request.EnterpriseID})
		if err != nil {
			return err, result, total
		}

		adminNameMap := make(map[int64]string)
		for _, admin := range adminList {
			adminNameMap[admin.ID] = admin.Username
		}
		var sourceSkuMap = make(map[int64][]string)
		sourceSkuMap, _ = s.GetSourceRelationGoods(c, request.EnterpriseID, sourceIds)

		corp, err := orm.BossService.Enterprise.GetInfo(bossModel.EnterpriseModel{ID: request.EnterpriseID})
		if err != nil {
			return err, result, total
		}

		corpHash := corp.Hash
		for k, v := range result {
			if adminName, ok := adminNameMap[v.AdminID]; ok {
				result[k].AdminName = adminName
			}

			if v.LinkType == sourceModules.LinkUrlTypeSelf {
				domain := config.CfgData.Service.WebH5.Domain
				linkUrl := fmt.Sprintf(sourceModules.SelfLinkUrl, domain, corpHash) + assist.StrVal(v.ID)
				v.LinkUrl = linkUrl
			}
			if topNum < 2 {
				result[k].CanTop = true
			}
			if _, ok := sourceSkuMap[v.ID]; ok {
				result[k].SkuList = sourceSkuMap[v.ID]
			}
			result[k].SkuList = helper.ArrayUniq(result[k].SkuList)
		}
	}

	return err, result, total
}
func (s Service) GetSourceRelationGoods(c app.Ctx, enterpriseId int64, sourceIds []int64) (map[int64][]string, error) {
	var sourceSkuMap = make(map[int64][]string)
	goods, err := sourceDao.NewRelationGoodsDao().GetMany(ctx.Wrap(c), "enterprise_id =? AND source_id in ?", enterpriseId, sourceIds)
	if len(goods) > 0 {
		for _, model := range goods {
			if _, ok := sourceSkuMap[model.SourceID]; !ok {
				sourceSkuMap[model.SourceID] = []string{model.SKU}
			} else {
				sourceSkuMap[model.SourceID] = append(sourceSkuMap[model.SourceID], model.SKU)
			}
			sourceSkuMap[model.SourceID] = helper.ArrayUniq(sourceSkuMap[model.SourceID])
		}
	}
	return sourceSkuMap, err
}

func (s Service) MiniProgramConsoleList(c app.Ctx, request sourceRequest.RequestConsoleList) (error, []source.ResponseConsoleMiniProgramList, int64) {
	// 七天前的起止时间戳
	sevenDaysStartTime, _ := helper.YesterdayPeriodOfTimestamp(7)
	topNum, err := s.sourceDao.GetTopNum(c, request.EnterpriseID, request.Type, request.GroupIDList)
	if err != nil {
		return err, nil, 0
	}
	err, result, total := s.sourceGroupDao.GetMiniProgramConsoleList(c, request.EnterpriseID, request.GroupIDList, request.CategoryIds, request.Keyword, request.SourceStatus, request.JustMe, request.AdminID, request.IsUse, request.CreatedStart, request.CreatedEnd, request.Order, request.Sort, request.Page, request.Size)

	if total > 0 {
		var adminIds []int64
		for _, v := range result {
			// 七天内上新素材
			if int64(v.CreatedAt) > sevenDaysStartTime {
				v.IsNew = 1
			}
			if v.AdminID > 0 {
				adminIds = append(adminIds, v.AdminID)
			}
		}

		adminList, err := microBossDao.NewEnterpriseAdmin().GetMany(c, map[string]any{"id": adminIds, "enterprise_id": request.EnterpriseID})
		if err != nil {
			return err, result, total
		}

		adminNameMap := make(map[int64]string)
		for _, admin := range adminList {
			adminNameMap[admin.ID] = admin.Username
		}

		for k, v := range result {
			if adminName, ok := adminNameMap[v.AdminID]; ok {
				result[k].AdminName = adminName
			}
			if topNum < 2 {
				result[k].CanTop = true
			}
		}
	}

	return err, result, total
}

func (s Service) VideoConsoleList(c app.Ctx, request sourceRequest.RequestConsoleList) (error, []source.ResponseConsoleVideoList, int64) {
	// 七天前的起止时间戳
	sevenDaysStartTime, _ := helper.YesterdayPeriodOfTimestamp(7)
	topNum, err := s.sourceDao.GetTopNum(c, request.EnterpriseID, request.Type, request.GroupIDList)
	if err != nil {
		return err, nil, 0
	}
	err, result, total := s.sourceGroupDao.GetVideoConsoleList(c, request.EnterpriseID, request.GroupIDList, request.CategoryIds, request.Keyword, request.SourceStatus, request.JustMe, request.AdminID, request.IsUse, request.CreatedStart, request.CreatedEnd, request.Order, request.Sort, request.ProductName, request.SkuItem, request.SkuColor, request.GoodsRelation, request.Page, request.Size)
	if total > 0 {
		var adminIds []int64
		var sourceIds []int64
		for _, v := range result {
			// 七天内上新素材
			if int64(v.CreatedAt) > sevenDaysStartTime {
				v.IsNew = 1
			}
			if v.AdminID > 0 {
				adminIds = append(adminIds, v.AdminID)
			}
			sourceIds = append(sourceIds, v.ID)
		}

		adminList, err := microBossDao.NewEnterpriseAdmin().GetMany(c, map[string]any{"id": adminIds, "enterprise_id": request.EnterpriseID})
		if err != nil {
			return err, result, total
		}

		adminNameMap := make(map[int64]string)
		for _, admin := range adminList {
			adminNameMap[admin.ID] = admin.Username
		}
		//获取素材关联商品sku
		sourceSkuMap, _ := s.GetSourceRelationGoods(c, request.EnterpriseID, sourceIds)
		for k, v := range result {
			if adminName, ok := adminNameMap[v.AdminID]; ok {
				result[k].AdminName = adminName
			}
			if topNum < 2 {
				result[k].CanTop = true
			}
			result[k].SkuList = make([]string, 0)
			if _, ok := sourceSkuMap[v.ID]; ok {
				result[k].SkuList = sourceSkuMap[v.ID]
			}
			result[k].SkuList = helper.ArrayUniq(result[k].SkuList)
		}
	}

	return err, result, total
}

func (s Service) FileConsoleList(c app.Ctx, request sourceRequest.RequestConsoleList) (error, []source.ResponseConsoleFileList, int64) {
	// 七天前的起止时间戳
	sevenDaysStartTime, _ := helper.YesterdayPeriodOfTimestamp(7)

	topNum, err := s.sourceDao.GetTopNum(c, request.EnterpriseID, request.Type, request.GroupIDList)
	if err != nil {
		return err, nil, 0
	}
	err, result, total := s.sourceGroupDao.GetFileConsoleList(c, request.EnterpriseID, request.GroupIDList, request.CategoryIds, request.Keyword, request.SourceStatus, request.JustMe, request.AdminID, request.IsUse, request.CreatedStart, request.CreatedEnd, request.Order, request.Sort, request.Page, request.Size)

	if total > 0 {
		var adminIds []int64
		for _, v := range result {
			// 七天内上新素材
			if int64(v.CreatedAt) > sevenDaysStartTime {
				v.IsNew = 1
			}
			if v.AdminID > 0 {
				adminIds = append(adminIds, v.AdminID)
			}
		}

		adminList, err := microBossDao.NewEnterpriseAdmin().GetMany(c, map[string]any{"id": adminIds, "enterprise_id": request.EnterpriseID})
		if err != nil {
			return err, result, total
		}

		adminNameMap := make(map[int64]string)
		for _, admin := range adminList {
			adminNameMap[admin.ID] = admin.Username
		}

		for k, v := range result {
			if adminName, ok := adminNameMap[v.AdminID]; ok {
				result[k].AdminName = adminName
			}
			if topNum < 2 {
				result[k].CanTop = true
			}
		}
	}

	return err, result, total
}

func (s Service) FriendConsoleList(c app.Ctx, request sourceRequest.RequestConsoleList) (error, []sourceResp.FriendList, int64) {
	// 七天前的起止时间戳
	sevenDaysStartTime, _ := helper.YesterdayPeriodOfTimestamp(7)

	topNum, err := s.sourceDao.GetTopNum(c, request.EnterpriseID, request.Type, request.GroupIDList)
	if err != nil {
		return err, nil, 0
	}
	err, res, total := s.sourceGroupDao.GetFriendConsoleList(c, request.EnterpriseID, request.GroupIDList, request.CategoryIds, request.Keyword, request.SourceStatus, request.JustMe, request.AdminID, request.IsUse, request.CreatedStart, request.CreatedEnd, request.Order, request.Sort, request.ProductName, request.SkuItem, request.SkuColor, request.GoodsRelation, request.Page, request.Size)

	resp := make([]sourceResp.FriendList, 0)
	for _, list := range res {
		resp = append(resp, sourceResp.FriendList{
			ResponseConsoleList: list.ResponseConsoleList,
			TextContent:         list.TextContent,
			Name:                list.Name,
			URL:                 list.OriginalURL,
			MergeId:             list.MergeId,
			MergeImage:          make([]sourceResp.MergeImage, 0),
		})
	}
	var goodsSourceIds []int64
	if total > 0 {
		var adminIds []int64
		var mergeIds []int64

		for _, v := range resp {
			// 七天内上新素材
			if int64(v.CreatedAt) > sevenDaysStartTime {
				v.IsNew = 1
			}
			if v.AdminID > 0 {
				adminIds = append(adminIds, v.AdminID)
			}
			if v.MergeId > 0 {
				mergeIds = append(mergeIds, v.MergeId)
			}
			goodsSourceIds = append(goodsSourceIds, v.ID)

		}

		adminNameMap := make(map[int64]string)
		// 查询管理员昵称
		if len(adminIds) > 0 {
			adminList, err := microBossDao.NewEnterpriseAdmin().GetMany(c, map[string]any{"id": adminIds, "enterprise_id": request.EnterpriseID})
			if err != nil {
				return err, resp, total
			}

			for _, admin := range adminList {
				adminNameMap[admin.ID] = admin.Username
			}
		}

		sourceMap := make(map[int64][]sourceResp.MergeImage)
		// 查询组合图
		if len(mergeIds) > 0 {
			sourceList, err := sourceDao.NewSourceDao().GetMany(c, map[string]any{"enterprise_id": request.EnterpriseID, "merge_id": mergeIds})
			if err != nil {
				return err, resp, total
			}

			var sourceIds []int64
			for _, val := range sourceList {
				sourceIds = append(sourceIds, val.ID)
			}
			goodsSourceIds = append(goodsSourceIds, sourceIds...)
			sourceFriendMap := make(map[int64]source.FriendModel)
			sourceFriendList, err := sourceDao.NewFriendDao().GetMany(c, request.EnterpriseID, map[string]any{"id": sourceIds})
			if err != nil {
				return err, resp, total
			}

			for _, friend := range sourceFriendList {
				if friend.OriginalURL != "" && friend.Type == source.FriendTypeImg {
					sourceFriendMap[friend.ID] = friend
				}
			}

			for _, val := range sourceList {
				if _, ok := sourceFriendMap[val.ID]; !ok {
					continue
				}
				tmpImage := sourceResp.MergeImage{
					Id:   sourceFriendMap[val.ID].ID,
					Name: sourceFriendMap[val.ID].Name,
					URL:  sourceFriendMap[val.ID].OriginalURL,
					Sort: val.Sort,
				}
				sourceMap[val.MergeID] = append(sourceMap[val.MergeID], tmpImage)
			}
		}
		//获取素材关联商品sku
		sourceSkuMap, _ := s.GetSourceRelationGoods(c, request.EnterpriseID, helper.ArrayUniq(goodsSourceIds))

		for k, v := range resp {
			if adminName, ok := adminNameMap[v.AdminID]; ok {
				resp[k].AdminName = adminName
			}

			// 组装合并图片
			if sourceImage, ok := sourceMap[v.MergeId]; ok {
				resp[k].MergeImage = append(resp[k].MergeImage, sourceImage...)
			}
			if topNum < 2 {
				resp[k].CanTop = true
			}
			//将素材关联商品sku写入结果
			var tmpSkuList = make([]string, 0)
			if _, ok := sourceSkuMap[v.ID]; ok {
				tmpSkuList = sourceSkuMap[v.ID]
			}
			for _, image := range resp[k].MergeImage {
				if _, ok := sourceSkuMap[image.Id]; ok {
					tmpSkuList = append(tmpSkuList, sourceSkuMap[image.Id]...)
				}
			}
			tmpSkuList = helper.ArrayUniq(tmpSkuList)
			resp[k].SkuList = tmpSkuList
		}
	}

	return err, resp, total
}

func (s Service) SelectEndListAndUpdate(c app.Ctx, enterpriseID int64) error {

	// 查询过期的id
	nowTime := time.Now().Unix()
	sourceIds, err := s.sourceDao.SelectEndList(c, enterpriseID, nowTime)

	if err != nil {
		return err
	}
	if len(sourceIds) == 0 {
		return nil
	}
	// 更新禁用状态
	query := map[string]interface{}{
		"enterprise_id": enterpriseID,
	}
	column := map[string]interface{}{
		"status":     1,
		"updated_at": uint(nowTime),
	}
	err = s.sourceDao.UpdateStatus(c, query, column, sourceIds)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) GetNewSourceCountForMember(c app.Ctx, enterpriseId, startTime, endTime int64) (res []source.NewSourceList, err error) {
	ignoreGroupIds, err := s.sourceGroupDao.GetIgnoreIndexGroupIds(c, enterpriseId)
	if err != nil {
		return nil, err
	}
	return s.sourceDao.GetNewSourceCountForMember(c, enterpriseId, startTime, endTime, ignoreGroupIds)
}

func (s Service) GetNewSourceCount(c app.Ctx, enterpriseId, startTime, endTime int64, thirdIds []int) (res []source.NewSourceCountNum, err error) {
	return s.sourceDao.GetNewSourceCount(c, enterpriseId, startTime, endTime, thirdIds)
}

func (s Service) SendTextMessage(c app.Ctx, data sourceRequest.SendRequest, enterpriseData *enterprise.Corp) (bool, error) {

	req := boss.WorkPlugin{
		EnterpriseID: data.EnterpriseID,
	}
	// 获取自建应用信息
	pluginInfo, errInfo := boss.WorkPluginInfo(c, req)
	if errInfo != nil {
		return false, errs.WorkPluginErr
	}

	if pluginInfo.Data.AgentID == 0 {
		return false, errs.PluginInfoRrr
	}

	// 昨天的起止时间戳
	startTime, endTime := helper.YesterdayPeriodOfTimestamp(1)
	// 查询需要推送消息的部门 部门下导购可查看素材权限的新素材条数
	sendNewSourceList, err1 := s.GetNewSourceCountForMember(c, data.EnterpriseID, startTime, endTime)
	if err1 != nil {
		return false, err1
	}

	// 处理相同数量的门店
	sendMessageThirdGroupMap := make(map[int64]sourceRequest.SendMessageThirdGroup)
	for _, v := range sendNewSourceList {
		sendMessageThirdGroupMap[v.IdCount] = sourceRequest.SendMessageThirdGroup{
			ThirdIds:  append(sendMessageThirdGroupMap[v.IdCount].ThirdIds, v.ThirdId),
			StartTime: v.StartTime,
			IdCount:   v.IdCount,
		}
	}

	// 企业H5的域名地址
	enterpriseConfig, err := boss.EnterpriseConfig(c, map[string]interface{}{
		"enterprise_id": data.EnterpriseID,
		"key":           "web_h5_domain",
	})
	if err != nil {
		return false, err
	}
	if enterpriseConfig.Data[0].Name == "" {
		log.Info(c).Msgf("企业%d未配置H5的域名地址", data.EnterpriseID)
		return false, err
	}
	urlH5 := enterpriseConfig.Data[0].Name

	workImpl := new(work.QyWorkImpl).SetDomain(enterpriseData.Domain)
	for _, v := range sendMessageThirdGroupMap {
		if len(v.ThirdIds) == 0 {
			log.Info(c).Msgf("企业%d未配置文本消息发送人群", data.EnterpriseID)
			continue
		}
		var messMem = sourceRequest.MessageMember{}
		// 根据门店获取要发送的员工
		memberList, _ := membersDao.NewMembers().GetMemberByThirdIds(c, data.EnterpriseID, v.ThirdIds, []int64{})
		for _, member := range memberList {
			messMem.Member = append(messMem.Member, member.MembersUserid)
		}

		fmt.Println(messMem.Member)
		messUser := strings.Join(messMem.Member, "|")

		startTimeStr := time.Unix(v.StartTime, 0).Format("2006-01-02")
		memberUserIdsTotal := len(messMem.Member)
		if memberUserIdsTotal > 100 {
			maxPage := len(messMem.Member) / 100
			maxPage += 1
			for i := 0; i < maxPage; i++ {
				start := i * 100
				end := (i + 1) * 100
				var memberUserIds []string
				if end >= memberUserIdsTotal {
					memberUserIds = messMem.Member[start:]
				} else {
					memberUserIds = messMem.Member[start:end]
				}
				messUser = strings.Join(memberUserIds, "|")
				textCardContent := work.MessageTypeTextCardContent{
					Title:       "素材库更新提醒",
					Description: fmt.Sprintf("<div> </div> <div class=\"gray\">素材数量：%d</div> <div class=\"gray\">生效时间：%s</div>", v.IdCount, startTimeStr),
					URL:         fmt.Sprintf("%s/work/%s/material-library?is_applets=1", urlH5, enterpriseData.Hash),
				}

				deptReq := work.SendMessageRequest{
					Touser:                 messUser,
					Msgtype:                "textcard",
					TextCard:               &textCardContent,
					Agentid:                pluginInfo.Data.AgentID,
					DuplicateCheckInterval: 1800,
					QyGwRequest: work.QyGwRequest{
						ReqPL: pluginInfo.Data.Identifier,
					},
				}
				res, errSend := workImpl.SendMessage(c, deptReq)
				if errSend != nil {
					log.Info(c).Msgf("企业%d下的发送信息应用[%s]的成员出错, 错误信息：%+v", data.EnterpriseID, pluginInfo.Data.Identifier, errSend)
					return false, errSend
				}
				if res.Errcode == 0 && res.Errmsg == "ok" {
					log.Info(c).Msgf("企业%d下的发送信息应用[%s]的成员成功, 返回信息：%+v", data.EnterpriseID, pluginInfo.Data.Identifier, res)
				}
			}
		} else {
			messUser = strings.Join(messMem.Member, "|")
			textCardContent := work.MessageTypeTextCardContent{
				Title:       "素材库更新提醒",
				Description: fmt.Sprintf("<div> </div> <div class=\"gray\">素材数量：%d</div> <div class=\"gray\">生效时间：%s</div>", v.IdCount, startTimeStr),
				URL:         fmt.Sprintf("%s/work/%s/material-library?is_applets=1", urlH5, enterpriseData.Hash),
			}

			deptReq := work.SendMessageRequest{
				Touser:                 messUser,
				Msgtype:                "textcard",
				TextCard:               &textCardContent,
				Agentid:                pluginInfo.Data.AgentID,
				DuplicateCheckInterval: 1800,
				QyGwRequest: work.QyGwRequest{
					ReqPL: pluginInfo.Data.Identifier,
				},
			}

			res, errSend := workImpl.SendMessage(c, deptReq)
			if errSend != nil {
				log.Info(c).Msgf("企业%d下的发送信息应用[%s]的成员出错, 错误信息：%+v", data.EnterpriseID, pluginInfo.Data.Identifier, errSend)
				return false, errSend
			}
			if res.Errcode == 0 && res.Errmsg == "ok" {
				log.Info(c).Msgf("企业%d下的发送信息应用[%s]的成员成功, 返回信息：%+v", data.EnterpriseID, pluginInfo.Data.Identifier, res)
			}
		}

	}

	return true, nil

}

// UploadMedia 小程序端导购的素材上传
func (s Service) UploadMedia(c app.Ctx, request sourceRequest.UploadMediaRequest) ([]int64, []sourceRequest.SourcesMediaInfo, int64, error) {
	adminInfo, err := microBossDao.NewEnterpriseAdmin().GetInfo(c, microBossModel.EnterpriseAdminModel{
		Mobile:         request.Mobile,
		MicroBaseModel: models.MicroBaseModel{EnterpriseID: int64(request.EnterpriseID)},
	})
	if err != nil {
		return nil, nil, 0, err
	}

	// 原资源集合
	var originMediaList []sourceRequest.SourcesMediaInfo
	// 压缩资源集合（仅视频，图片改用COS参数）
	var compressMediaList []sourceRequest.SourcesMediaInfo
	// 图片扩展名列表
	imageExts := []string{"PNG", "JPG", "JPEG", "HEIF"}
	videoExts := []string{"MP4", "MOV"}

	// 兼容只写文案不上传素材情况
	if request.Content != "" && len(request.Sources) == 0 {
		originMediaList = append(originMediaList, sourceRequest.SourcesMediaInfo{
			SourcesInfo: sourceRequest.SourcesInfo{SourceId: app.UID(), SourceType: source.TypeFriend},
			MediaInfo:   sourceRequest.MediaInfo{IsOrigin: 1},
		})
		// 不合并素材组
		request.ShowType = 0
	}

	var sourceIds []int64
	for _, sourcePtr := range request.Sources {
		sourcesMediaInfo := sourceRequest.SourcesMediaInfo{
			SourcesInfo: sourceRequest.SourcesInfo{SourceId: app.UID(), SourceType: request.SourceType},
			MediaInfo:   sourcePtr,
		}
		sourceIds = append(sourceIds, sourcesMediaInfo.SourceId)

		isImage := helper.InArray(strings.ToUpper(sourcePtr.Ext), imageExts)
		isVideo := helper.InArray(strings.ToUpper(sourcePtr.Ext), videoExts)

		// 图片：不再异步压缩，统一走originMediaList（朋友圈场景在DAO层拼COS压缩URL）
		if isImage {
			originMediaList = append(originMediaList, sourcesMediaInfo)
			continue
		}

		// 视频：小于10m不压缩
		if isVideo && sourcePtr.Size < 10*(1024*1024) {
			originMediaList = append(originMediaList, sourcesMediaInfo)
			continue
		}

		// 原图且不转码
		if sourcePtr.IsOrigin == 1 && sourcePtr.IsTranscode != 1 {
			originMediaList = append(originMediaList, sourcesMediaInfo)
			continue
		}

		// 压缩资源集合（仅视频）
		compressMediaList = append(compressMediaList, sourcesMediaInfo)
	}

	compressId := app.UID()
	adminId := adminInfo.ID

	// 图片不再走CreateCompressImage，统一走CreateImage
	switch request.SourceType {
	case sourceModules.TypeText:
		sourceInfo := sourceRequest.SourcesInfo{SourceId: app.UID(), SourceType: request.SourceType}
		sourceIds = append(sourceIds, sourceInfo.SourceId)
		err = sourceDao.NewTextDao().CreateText(c, sourceInfo, request, adminId)
	case sourceModules.TypeImg:
		// 图片全部走originMediaList，不再分流
		if len(originMediaList) > 0 {
			err = sourceDao.NewImageDao().CreateImage(c, originMediaList, request, adminId)
		}
	case sourceModules.TypeVideo:
		if len(originMediaList) > 0 {
			err = sourceDao.NewVideoDao().CreateVideo(c, originMediaList, request, adminId)
		}
		if len(compressMediaList) > 0 {
			err = sourceDao.NewVideoDao().CreateCompressVideo(c, compressMediaList, request, compressId, adminId)
		}
	case sourceModules.TypeCollect:
		// 判断是否为组合展示,组合展示需要新增合并id

		mergeId := app.UID()

		if len(originMediaList) > 0 {
			err = sourceDao.NewFriendDao().CreateFriend(c, originMediaList, request, 0, adminId, mergeId)
		}
		if len(compressMediaList) > 0 {
			err = sourceDao.NewFriendDao().CreateFriend(c, compressMediaList, request, compressId, adminId, mergeId)
		}
	}

	// 保存关联商品
	for i, sourcePtr := range request.Sources {
		if len(sourcePtr.RelationGoods) == 0 {
			continue
		}
		sourceId := sourceIds[i]
		switch request.SourceType {
		case sourceModules.TypeImg:
			_ = sourceDao.NewRelationGoodsDao().DeleteData(c, "source_id = ? AND type = ?", sourceId, sourceModules.TypeImg)
			s.saveRelationGoods(c, sourceId, int(request.EnterpriseID), sourceModules.TypeImg, sourcePtr.RelationGoods)
		case sourceModules.TypeVideo:
			_ = sourceDao.NewRelationGoodsDao().DeleteData(c, "source_id = ? AND type = ?", sourceId, sourceModules.TypeVideo)
			s.saveRelationGoods(c, sourceId, int(request.EnterpriseID), sourceModules.TypeVideo, sourcePtr.RelationGoods)
		case sourceModules.TypeCollect:
			_ = sourceDao.NewRelationGoodsDao().DeleteData(c, "source_id = ? AND type = ?", sourceId, sourceModules.TypeCollect)
			s.saveRelationGoods(c, sourceId, int(request.EnterpriseID), sourceModules.TypeCollect, sourcePtr.RelationGoods)
		}
	}

	return sourceIds, compressMediaList, compressId, err
}

// saveRelationGoods 批量保存关联商品
func (s Service) saveRelationGoods(c app.Ctx, sourceId int64, enterpriseId int, goodsType int8, goodsList []sourceRequest.BrandGoods) {
	var goodsModels []source.RelationGoodsModel
	for _, item := range goodsList {
		goods := source.RelationGoodsModel{
			MicroBaseModel: models.MicroBaseModel{
				BaseModel: models.BaseModel{
					ID:        app.UID(),
					CreatedAt: time.Now().Unix(),
					UpdatedAt: time.Now().Unix(),
				},
				EnterpriseID: int64(enterpriseId),
			},
			SourceID: sourceId,
			Type:     goodsType,
			UPC:      item.Upc,
			SKU:      item.Sku,
			Name:     item.Name,
			Color:    item.Color,
			SkuCode:  fmt.Sprintf("%s %s", item.Sku, item.Color),
			Url:      item.Url,
			AddType:  item.AddType,
			Score:    item.Score,
		}
		goodsModels = append(goodsModels, goods)
	}
	if len(goodsModels) > 0 {
		_ = sourceDao.NewRelationGoodsDao().CreateInBatchesData(c, goodsModels)
	}
}

// MediaCompress 资源信息压缩（仅处理视频，图片已改为COS参数）
func (s Service) MediaCompress(c ctx.Context, request sourceRequest.UploadMediaRequest, medias []sourceRequest.SourcesMediaInfo, Corp *enterprise.Corp, compressId int64) {
	var successNum int
	for _, media := range medias {
		whereQuery := map[string]any{"id": media.SourceId}

		// 获取文件名和文件后缀
		fileName := filepath.Base(media.Url)
		fileExt := strings.TrimLeft(filepath.Ext(fileName), ".")

		// 图片不再走ffmpeg压缩，跳过
		if helper.InArray(strings.ToUpper(fileExt), []string{"PNG", "JPG", "JPEG", "HEIF"}) {
			continue
		}

		// 构建压缩对象和参数（仅视频）
		ffmpegClient := helper.NewFfmpeg(media.Url)
		if helper.InArray(strings.ToUpper(fileExt), []string{"MP4", "MOV"}) {
			ffmpegClient.VideoBuildKwArgs()
		}
		// 执行压缩
		_, compressFileSize, err := ffmpegClient.Compress()
		// 上传cos
		filePath, err := ffmpegClient.UploadCos()

		if err != nil {
			s.Log().Error(c).Msgf("压缩文件失败：%+v", err)
			err = sourceDao.NewSourceDao().Update(c, whereQuery, map[string]any{"is_deal": source.DealFail, "deleted_at": time.Now().Unix()})
			if err != nil {
				s.Log().Error(c).Msgf("更新Source表异常：%+v", err)
			}
			continue
		}

		successNum++
		switch request.SourceType {
		case int32(sourceModules.TypeVideo):
			err = sourceDao.NewVideoDao().Update(c, whereQuery, map[string]any{"size": compressFileSize, "ext": fileExt, "url": filePath})
		case int32(sourceModules.TypeCollect):
			err = sourceDao.NewFriendDao().Update(c, whereQuery, map[string]any{"size": compressFileSize, "compression_url": filePath})
		}

		if err != nil {
			s.Log().Error(c).Msgf("更新SourceImage/SourceVideo异常：%+v", err)
		}

		err = sourceDao.NewSourceDao().Update(c, whereQuery, map[string]interface{}{"is_deal": source.DealSuccess})
		if err != nil {
			s.Log().Error(c).Msgf("更新Source表异常：%+v", err)
		}
	}

	err := s.compressResultSendTextMessage(c, request, len(medias), successNum, Corp, compressId)
	if err != nil {
		s.Log().Error(c).Msgf("压缩结果消息推送异常：%+v", err)
		return
	}
}

// compressResultSendTextMessage 压缩结果消息推送
func (s Service) compressResultSendTextMessage(c ctx.Context, request sourceRequest.UploadMediaRequest, allNum, successNum int, Corp *enterprise.Corp, compressId int64) error {
	req := boss.WorkPlugin{
		EnterpriseID: int64(request.EnterpriseID),
	}
	// 获取自建应用信息
	pluginInfo, errInfo := boss.WorkPluginInfo(c, req)
	if errInfo != nil {
		return errs.WorkPluginErr
	}

	if pluginInfo.Data.AgentID == 0 {
		return errs.PluginInfoRrr
	}
	enterpriseConfig, _ := boss.EnterpriseConfig(c, map[string]interface{}{
		"enterprise_id": request.EnterpriseID,
		"key":           "web_h5_domain",
	})
	domain := enterpriseConfig.Data[0].Name
	workImpl := new(work.QyWorkImpl).SetDomain(Corp.Domain)
	textCardContent := work.MessageTypeTextCardContent{
		Title:       "素材库上传压缩结果提醒",
		Description: fmt.Sprintf("<div> </div> <div class=\"gray\">总压缩数量：%d</div> <div class=\"gray\">成功数量：%d</div>", allNum, successNum),
		URL:         fmt.Sprintf("%s/work/%s/material-result?id=%d&type=%d", domain, Corp.Hash, compressId, request.SourceType),
	}

	deptReq := work.SendMessageRequest{
		Touser:                 request.UserId,
		Msgtype:                "textcard",
		TextCard:               &textCardContent,
		Agentid:                pluginInfo.Data.AgentID,
		DuplicateCheckInterval: 1800,
		QyGwRequest: work.QyGwRequest{
			ReqPL: pluginInfo.Data.Identifier,
		},
	}

	res, errSend := workImpl.SendMessage(c, deptReq)
	if errSend != nil {
		log.Info(c).Msgf("企业%d下的发送信息应用[%s]的成员出错, 错误信息：%+v", request.EnterpriseID, pluginInfo.Data.Identifier, errSend)
		return errSend
	}
	if res.Errcode == 0 && res.Errmsg == "ok" {
		log.Info(c).Msgf("企业%d下的发送信息应用[%s]的成员成功, 返回信息：%+v", request.EnterpriseID, pluginInfo.Data.Identifier, res)
	}
	return nil
}

// GetCompress 获取压缩结果
func (s Service) GetCompress(c ctx.Context, request sourceRequest.CompressReq) (result []sourceRequest.CompressRes, err error) {
	switch request.SourceType {
	case sourceModules.TypeImg:
		result, err = sourceDao.NewImageDao().GetCompress(c, request)
	case sourceModules.TypeVideo:
		result, err = sourceDao.NewVideoDao().GetCompress(c, request)
	}
	for k, res := range result {
		info := res
		if info.IsDeal == source.DealFail {
			info.Reason = "尺寸过大，请修改后重新上传"
		}
		result[k] = info
	}
	return result, err

}

func (s Service) CompressFile(c ctx.Context, filePath, fileType string) (string, error) {
	// 获取文件名和文件后缀
	fileName := filepath.Base(filePath)
	fileExt := filepath.Ext(fileName)
	if fileType == "image" {
		if !helper.InArray(strings.ToUpper(fileExt), []string{".PNG", ".JPG", ".JPEG", ".HEIF"}) {
			return "", errors.New(fmt.Sprintf("不支持[%s]类型", fileExt))
		}
		md5Name := helper.Md5V(fileName) + helper.Int64ToString(time.Now().UnixNano())
		// 压缩文件地址
		fileLocalPath := fmt.Sprintf("storage/upload/%s%s", md5Name, fileExt)
		err := helper.DownloadFileToLocal(filePath, fileLocalPath)
		defer func() { _ = os.Remove(filePath) }()
		if err != nil {
			return "", err
		}
		width, height, err := helper.GetImageSize(fileLocalPath)
		if err != nil {
			return "", err
		}
		// 计算图片尺寸
		newWidth, _ := helper.CalculateNewDimensions(width, height, 1440, 1080)
		imageCosCompressUrl := fmt.Sprintf("%s?imageMogr2/thumbnail/%dx", filePath, newWidth)
		return imageCosCompressUrl, nil
	}
	ffmpegClient := helper.NewFfmpeg(filePath)
	if fileType == "video" {
		if !helper.InArray(strings.ToUpper(fileExt), []string{".MP4", ".MOV"}) {
			return "", errors.New(fmt.Sprintf("不支持[%s]类型", fileExt))
		}
		ffmpegClient.VideoBuildKwArgs()
	}
	// 执行压缩
	compressFilePath, _, err := ffmpegClient.Compress()
	defer func() { _ = os.Remove(compressFilePath) }()
	if err != nil {
		return "", err
	}
	// 上传cos
	cosFileAddr, err := ffmpegClient.UploadCos()
	if err != nil {
		return "", err
	}

	return cosFileAddr, nil
}

// EditVideoTime 修改视频创建时间
func (s Service) EditVideoTime(filePath string) (string, error) {
	// 获取文件名和文件后缀
	fileName := filepath.Base(filePath)
	fileExt := filepath.Ext(fileName)
	ffmpegClient := helper.NewFfmpeg(filePath)

	if !helper.InArray(strings.ToUpper(fileExt), []string{".MP4", ".MOV"}) {
		return "", errors.New(fmt.Sprintf("不支持[%s]类型", fileExt))
	}
	ffmpegClient.VideoBuildKwArgs()
	// 执行修改
	compressFilePath, err := ffmpegClient.EditVideoCreateTime()
	defer func() { _ = os.Remove(compressFilePath) }()
	if err != nil {
		return "", err
	}
	// 上传cos
	cosFileAddr, err := ffmpegClient.UploadCos()
	if err != nil {
		return "", err
	}
	return cosFileAddr, nil
}

func (s Service) ReportTimes(c ctx.Context, req sourceRequest.ReportTimesReq) (any, error) {
	var incrementColumn string
	switch req.Type {
	case 1:
		incrementColumn = "download_times"
	case 2:
		incrementColumn = "forwarding_times"
	case 3:
		incrementColumn = "copy_times"
	default:
		return nil, errors.New("上报类型错误")
	}
	reportTimes := make(map[string]map[string]int)

	if len(req.ResourceKeys) > 0 {
		err := sourceDao.NewBaozunCacheDao().Increment(c, map[string]any{"resource_key": req.ResourceKeys, "enterprise_id": req.EnterpriseID}, incrementColumn, 1)
		if err != nil {
			return nil, err
		}
		// 获取素材对应的下载,复制,转发次数
		sourceList, err := sourceDao.NewBaozunCacheDao().GetMany(c, req.EnterpriseID, map[string]any{"resource_key": req.ResourceKeys})
		if err != nil {
			return nil, err
		}
		for _, val := range sourceList {
			reportTimesVal := make(map[string]int)
			reportTimesVal["copy_times"] = val.CopyTimes
			reportTimesVal["download_times"] = val.DownloadTimes
			reportTimesVal["forwarding_times"] = val.ForwardingTimes
			reportTimes[val.ResourceKey] = reportTimesVal
		}
	} else if len(req.Ids) > 0 {
		err := sourceDao.NewSourceDao().Increment(c, map[string]any{"id": req.Ids, "enterprise_id": req.EnterpriseID}, incrementColumn, 1)
		if err != nil {
			return nil, err
		}

		// 获取素材对应的下载,复制,转发次数
		sourceList, err := sourceDao.NewSourceDao().GetMany(c, map[string]any{"enterprise_id": req.EnterpriseID, "id": req.Ids})
		if err != nil {
			return nil, err
		}
		for _, val := range sourceList {
			reportTimesVal := make(map[string]int)
			reportTimesVal["copy_times"] = val.CopyTimes
			reportTimesVal["download_times"] = val.DownloadTimes
			reportTimesVal["forwarding_times"] = val.ForwardingTimes
			reportTimes[helper.Int64ToString(val.ID)] = reportTimesVal
		}
	}

	return reportTimes, nil
}

func (s Service) MakeCover(c ctx.Context, filePath string) (string, string, error) {
	// 获取文件名和文件后缀

	ffmpegClient := helper.NewFfmpeg(filePath)
	ffmpegClient.SetFileExt("mp4")
	ffmpegClient = ffmpegClient.BuildKwArgs(map[string]any{
		"c":        "copy",
		"metadata": "creation_time=",
	})

	// 执行压缩
	compressFilePath, imageUrl, err := ffmpegClient.MakeCover()
	defer func() { _ = os.Remove(compressFilePath) }()
	if err != nil {
		return "", "", err
	}
	// 上传cos
	cosFileAddr, err := ffmpegClient.UploadCos()

	imageCosUrl, err := helper.UploadFileToCos(c, imageUrl, false)
	if err != nil {
		return "", "", err
	}
	return cosFileAddr, imageCosUrl, nil
}

func (s Service) SetTop(c ctx.Context, req sourceRequest.SourceSetTop) error {
	var err error
	whereQuery := map[string]any{
		"enterprise_id": req.EnterpriseId,
		"type":          req.SourceType,
		"is_top":        source.TopYes,
	}
	if req.IsTop == source.TopYes { //要置顶的素材
		sourceList, err := sourceDao.NewSourceDao().GetMany(c, whereQuery)
		if err != nil {
			return err
		}
		if len(sourceList) >= 2 {
			return errors.New("最多只能置顶两条素材")
		}
		err = sourceDao.NewSourceDao().Update(c, map[string]interface{}{
			"id": req.SourceId,
		}, map[string]interface{}{"is_top": source.TopYes})
	} else {
		err = sourceDao.NewSourceDao().Update(c, map[string]interface{}{
			"id": req.SourceId,
		}, map[string]interface{}{"is_top": source.TopNo})
	}
	return err
}
