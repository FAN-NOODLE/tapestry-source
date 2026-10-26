package openSvc

import (
	"errors"
	"sort"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	sourceDao "micro-scrm/internal/app/dao/source"
	openReq "micro-scrm/internal/app/http/request/open"
	openResp "micro-scrm/internal/app/http/response/open"
	"micro-scrm/internal/app/models/source"
	sourceGroup "micro-scrm/internal/modules/sourceGroup"
	"micro-scrm/internal/orm/model/micro_scrm"
)

type SourceService struct {
	sourceDao *sourceDao.SourceGroupDao
}

func NewSourceService() *SourceService {
	return &SourceService{
		sourceDao: sourceDao.NewSourceGroupDao(),
	}
}

// getSourceGroupIds 根据ThirdIds获取可见的素材分组ID列表
func getSourceGroupIds(c app.Ctx, enterpriseId int64, thirdIds []int) []int64 {
	relationDao := sourceDao.NewGroupRelationDao()
	if len(thirdIds) > 0 {
		sourceGroupIds, _ := relationDao.GetManySourceGroupIds(c, enterpriseId, map[string]interface{}{
			"third_id": thirdIds,
		})
		return sourceGroupIds
	}
	return make([]int64, 0)
}

// getBrandGoodsMap 批量查询关联商品，返回 sourceId -> []BrandGoods 的映射
func getBrandGoodsMap(c app.Ctx, enterpriseId int64, sourceIds []int64) map[int64][]openResp.BrandGoods {
	sourceSkuMap := make(map[int64][]openResp.BrandGoods)
	if len(sourceIds) == 0 {
		return sourceSkuMap
	}
	goods, _ := sourceDao.NewRelationGoodsDao().GetMany(ctx.Wrap(c), "enterprise_id =? AND source_id in ? AND deleted_at =0", enterpriseId, sourceIds)
	for _, model := range goods {
		sourceSkuMap[model.SourceID] = append(sourceSkuMap[model.SourceID], openResp.BrandGoods{
			Name:  model.Name,
			Sku:   model.SKU,
			Upc:   model.UPC,
			Url:   model.Url,
			Color: model.Color,
		})
	}
	return sourceSkuMap
}

// LinkList 获取对外链接列表
func (s SourceService) LinkList(c app.Ctx, request openReq.SearchListRequest) ([]openResp.LinkListResponse, int64, error) {
	sourceGroupIds := getSourceGroupIds(c, request.EnterpriseID, request.ThirdIds)
	openList, total, err := sourceDao.NewGroupRelationDao().GetLinkOpenList(c, request, sourceGroupIds)
	if err != nil {
		return nil, 0, err
	}
	var sourceIds []int64
	for _, info := range openList {
		sourceIds = append(sourceIds, info.ID)
	}
	sourceSkuMap := getBrandGoodsMap(c, request.EnterpriseID, sourceIds)
	for i, info := range openList {
		openList[i].BrandGoods = sourceSkuMap[info.ID]
	}
	return openList, total, nil
}

// SourceTree 素材分组树
func (s SourceService) SourceTree(request openReq.SourceTreeRequest) ([]*micro_scrm.SourceTree, error) {
	return sourceGroup.SourceTree(request.EnterpriseID, 0)
}

// SourceList  素材列表
func (s SourceService) SourceList(c app.Ctx, request openReq.SourceOpenReq) ([]openResp.SourceListResponse, int64, error) {
	// 确保分页参数有效
	if request.Page < 1 {
		request.Page = 1
	}
	if request.Size < 1 {
		request.Size = 10
	}

	var (
		resList []openResp.SourceListResponse
		total   int64
		err     error
	)
	switch request.Type {
	case source.TypeText:
		var list []source.TextModel
		list, total, err = sourceDao.NewTextDao().GetListByOpen(c, request)
		for _, model := range list {
			resList = append(resList, openResp.SourceListResponse{
				ID:        model.ID,
				GroupId:   model.Source.GroupID,
				Content:   model.Content,
				Status:    model.Source.Status,
				StartTime: model.Source.StartTime,
				EndTime:   model.Source.EndTime,
				CreatedAt: model.CreatedAt,
				UpdatedAt: model.UpdatedAt,
				DeleteAt:  model.DeletedAt,
			})
		}
	case source.TypeImg:
		var list []source.ImageModel
		list, total, err = sourceDao.NewImageDao().GetListByOpen(c, request)
		for _, model := range list {
			resList = append(resList, openResp.SourceListResponse{
				ID:        model.ID,
				GroupId:   model.Source.GroupID,
				MergeId:   model.Source.MergeID,
				Name:      model.Name,
				URL:       model.URL,
				Ext:       model.Ext,
				Status:    model.Source.Status,
				StartTime: model.Source.StartTime,
				EndTime:   model.Source.EndTime,
				CreatedAt: model.CreatedAt,
				UpdatedAt: model.UpdatedAt,
				DeleteAt:  model.DeletedAt,
			})
		}
	case source.TypeLink:
		var list []source.LinkModel
		list, total, err = sourceDao.NewLinkDao().GetListByOpen(c, request)
		for _, model := range list {
			resList = append(resList, openResp.SourceListResponse{
				ID:        model.ID,
				GroupId:   model.Source.GroupID,
				Name:      model.Title,
				Content:   model.Description,
				CoverURL:  model.CoverURL,
				URL:       model.LinkURL,
				Status:    model.Source.Status,
				StartTime: model.Source.StartTime,
				EndTime:   model.Source.EndTime,
				CreatedAt: model.CreatedAt,
				UpdatedAt: model.UpdatedAt,
				DeleteAt:  model.DeletedAt,
			})
		}
	case source.TypeVideo:
		var list []source.VideoModel
		list, total, err = sourceDao.NewVideoDao().GetListByOpen(c, request)
		for _, model := range list {
			resList = append(resList, openResp.SourceListResponse{
				ID:        model.ID,
				GroupId:   model.Source.GroupID,
				Name:      model.Name,
				CoverURL:  model.CoverURL,
				URL:       model.URL,
				Ext:       model.Ext,
				Status:    model.Source.Status,
				StartTime: model.Source.StartTime,
				EndTime:   model.Source.EndTime,
				CreatedAt: model.CreatedAt,
				UpdatedAt: model.UpdatedAt,
				DeleteAt:  model.DeletedAt,
			})
		}
	case source.TypeFriend:
		var list []source.FriendModel
		list, total, err = sourceDao.NewFriendDao().GetListByOpen(c, request)
		for _, model := range list {
			resList = append(resList, openResp.SourceListResponse{
				ID:        model.ID,
				GroupId:   model.Source.GroupID,
				MergeId:   model.Source.MergeID,
				Name:      model.Name,
				Content:   model.TextContent,
				CoverURL:  model.VideoCoverURL,
				URL:       model.OriginalURL,
				Ext:       model.Ext,
				Status:    model.Source.Status,
				StartTime: model.Source.StartTime,
				EndTime:   model.Source.EndTime,
				CreatedAt: model.CreatedAt,
				UpdatedAt: model.UpdatedAt,
				DeleteAt:  model.DeletedAt,
			})
		}
	default:
		return nil, 0, errors.New("素材类型参数不正确")
	}

	// 查询关联商品
	if err == nil && len(resList) > 0 {
		var sourceIds []int64
		for _, item := range resList {
			sourceIds = append(sourceIds, item.ID)
		}
		sourceSkuMap := getBrandGoodsMap(c, request.EnterpriseID, sourceIds)
		for i := range resList {
			resList[i].BrandGoods = sourceSkuMap[resList[i].ID]
		}
		if request.Type == source.TypeLink {
			sourceURLMap, err := sourceDao.NewLinkKOSDao().GetSourceURLMap(c, request.EnterpriseID, sourceIds)
			if err != nil {
				return nil, 0, err
			}
			for i := range resList {
				resList[i].SourceURL = sourceURLMap[resList[i].ID]
			}
		}
		if request.IncludeCategories {
			categoryMap, err := sourceDao.NewCategoryRelationDao().GetOpenCategoriesBySourceIds(c, request.EnterpriseID, sourceIds)
			if err != nil {
				return nil, 0, err
			}
			for i := range resList {
				categories := categoryMap[resList[i].ID]
				if categories == nil {
					categories = make([]openResp.SourceCategory, 0)
				}
				resList[i].Categories = &categories
			}
		}
	}

	return resList, total, err
}

// TextList open API文本素材列表
func (s SourceService) TextList(c app.Ctx, request openReq.SearchListRequest) ([]openResp.TextListResponse, int64, error) {
	sourceGroupIds := getSourceGroupIds(c, request.EnterpriseID, request.ThirdIds)
	openList, total, err := sourceDao.NewGroupRelationDao().GetTextOpenList(c, request, sourceGroupIds)
	if err != nil {
		return nil, 0, err
	}
	var sourceIds []int64
	for _, info := range openList {
		sourceIds = append(sourceIds, info.ID)
	}
	sourceSkuMap := getBrandGoodsMap(c, request.EnterpriseID, sourceIds)
	for i := range openList {
		openList[i].BrandGoods = sourceSkuMap[openList[i].ID]
	}
	return openList, total, nil
}

// ImageList open API图片素材列表（含组合图分组）
func (s SourceService) ImageList(c app.Ctx, request openReq.SearchListRequest) ([]openResp.ImageListResponse, int64, error) {
	sourceGroupIds := getSourceGroupIds(c, request.EnterpriseID, request.ThirdIds)
	openList, total, err := sourceDao.NewGroupRelationDao().GetImageOpenList(c, request, sourceGroupIds)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return openList, total, nil
	}

	// 收集所有 sourceId 和 mergeId
	var allSourceIds []int64
	var mergeIds []int64
	for i, v := range openList {
		allSourceIds = append(allSourceIds, v.ID)
		if v.MergeId > 0 {
			// 同一 mergeId 只查一次
			mergeIds = append(mergeIds, v.MergeId)
			// 移除重复的 mergeId 对应的记录（只保留第一条作为封面）
			for j := i + 1; j < len(openList); {
				if openList[j].MergeId == v.MergeId {
					// 将被移除的记录ID也加入 allSourceIds（用于查商品关联）
					allSourceIds = append(allSourceIds, openList[j].ID)
					openList = append(openList[:j], openList[j+1:]...)
					total--
				} else {
					j++
				}
			}
		}
	}

	// 查询组合图子图
	sourceMap := make(map[int64][]openResp.OpenMergeImage)
	if len(mergeIds) > 0 {
		sourceList, err := sourceDao.NewSourceDao().GetMany(c, map[string]any{"enterprise_id": request.EnterpriseID, "merge_id": mergeIds})
		if err != nil {
			return nil, 0, err
		}
		var subSourceIds []int64
		for _, val := range sourceList {
			subSourceIds = append(subSourceIds, val.ID)
			allSourceIds = append(allSourceIds, val.ID)
		}
		sourceImageMap := make(map[int64]source.ImageModel)
		sourceImageList, err := sourceDao.NewImageDao().GetMany(c, request.EnterpriseID, map[string]any{"id": subSourceIds})
		if err != nil {
			return nil, 0, err
		}
		for _, image := range sourceImageList {
			sourceImageMap[image.ID] = image
		}
		for _, val := range sourceList {
			if img, ok := sourceImageMap[val.ID]; ok {
				sourceMap[val.MergeID] = append(sourceMap[val.MergeID], openResp.OpenMergeImage{
					Id:   img.ID,
					Name: img.Name,
					URL:  img.URL,
					Sort: val.Sort,
				})
			}
		}
	}

	// 批量查所有关联商品
	sourceSkuMap := getBrandGoodsMap(c, request.EnterpriseID, allSourceIds)

	for k, v := range openList {
		// 组装合并图片
		if subImages, ok := sourceMap[v.MergeId]; ok {
			sort.Slice(subImages, func(i, j int) bool {
				return subImages[i].Sort < subImages[j].Sort
			})
			// 子图的 BrandGoods
			for si := range subImages {
				subImages[si].BrandGoods = sourceSkuMap[subImages[si].Id]
			}
			openList[k].MergeImage = subImages
			if len(subImages) > 0 {
				openList[k].URL = subImages[0].URL
			}
		}
		openList[k].BrandGoods = sourceSkuMap[v.ID]
	}

	return openList, total, nil
}

// VideoList open API视频素材列表
func (s SourceService) VideoList(c app.Ctx, request openReq.SearchListRequest) ([]openResp.VideoListResponse, int64, error) {
	sourceGroupIds := getSourceGroupIds(c, request.EnterpriseID, request.ThirdIds)
	openList, total, err := sourceDao.NewGroupRelationDao().GetVideoOpenList(c, request, sourceGroupIds)
	if err != nil {
		return nil, 0, err
	}
	var sourceIds []int64
	for _, info := range openList {
		sourceIds = append(sourceIds, info.ID)
	}
	sourceSkuMap := getBrandGoodsMap(c, request.EnterpriseID, sourceIds)
	for i := range openList {
		openList[i].BrandGoods = sourceSkuMap[openList[i].ID]
	}
	return openList, total, nil
}

// FriendList open API朋友圈素材列表（含组合图分组）
func (s SourceService) FriendList(c app.Ctx, request openReq.SearchListRequest) ([]openResp.FriendListResponse, int64, error) {
	sourceGroupIds := getSourceGroupIds(c, request.EnterpriseID, request.ThirdIds)
	openList, total, err := sourceDao.NewGroupRelationDao().GetFriendOpenList(c, request, sourceGroupIds)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return openList, total, nil
	}

	// 收集所有 sourceId、mergeId
	var allSourceIds []int64
	var mergeIds []int64
	for i, v := range openList {
		allSourceIds = append(allSourceIds, v.ID)
		if v.MergeId > 0 {
			mergeIds = append(mergeIds, v.MergeId)
			for j := i + 1; j < len(openList); {
				if openList[j].MergeId == v.MergeId {
					allSourceIds = append(allSourceIds, openList[j].ID)
					openList = append(openList[:j], openList[j+1:]...)
					total--
				} else {
					j++
				}
			}
		}
	}

	// 查询组合图子图（只取 FriendTypeImg 类型）
	sourceMap := make(map[int64][]openResp.OpenMergeImage)
	if len(mergeIds) > 0 {
		sourceList, err := sourceDao.NewSourceDao().GetMany(c, map[string]any{"enterprise_id": request.EnterpriseID, "merge_id": mergeIds})
		if err != nil {
			return nil, 0, err
		}
		var subSourceIds []int64
		for _, val := range sourceList {
			subSourceIds = append(subSourceIds, val.ID)
			allSourceIds = append(allSourceIds, val.ID)
		}
		sourceFriendMap := make(map[int64]source.FriendModel)
		sourceFriendList, err := sourceDao.NewFriendDao().GetMany(c, request.EnterpriseID, map[string]any{"id": subSourceIds})
		if err != nil {
			return nil, 0, err
		}
		for _, friend := range sourceFriendList {
			if friend.OriginalURL != "" && friend.Type == source.FriendTypeImg {
				sourceFriendMap[friend.ID] = friend
			}
		}
		for _, val := range sourceList {
			if friend, ok := sourceFriendMap[val.ID]; ok {
				sourceMap[val.MergeID] = append(sourceMap[val.MergeID], openResp.OpenMergeImage{
					Id:   friend.ID,
					Name: friend.Name,
					URL:  friend.OriginalURL,
					Sort: val.Sort,
				})
			}
		}
	}

	// 批量查所有关联商品
	sourceSkuMap := getBrandGoodsMap(c, request.EnterpriseID, allSourceIds)

	for k, v := range openList {
		// 组装合并图片
		if subImages, ok := sourceMap[v.MergeId]; ok {
			sort.Slice(subImages, func(i, j int) bool {
				return subImages[i].Sort < subImages[j].Sort
			})
			for si := range subImages {
				subImages[si].BrandGoods = sourceSkuMap[subImages[si].Id]
			}
			openList[k].MergeImage = subImages
		}
		openList[k].BrandGoods = sourceSkuMap[v.ID]
	}

	return openList, total, nil
}

// SearchAll 合并检索：通过商品条件一次性查出所有素材，按类型分类返回（不分页）
func (s SourceService) SearchAll(c app.Ctx, request openReq.SearchAllRequest) (openResp.SearchAllResponse, error) {
	var result openResp.SearchAllResponse

	// 1. 按商品条件一次性查出所有 source_id + type
	sourceIdTypes, err := sourceDao.NewGroupRelationDao().GetSourceIdsByProduct(c, request.EnterpriseID, request.ProductSku, request.ProductColor, request.ProductName)
	if err != nil {
		return result, err
	}
	if len(sourceIdTypes) == 0 {
		return result, nil
	}

	// 2. 权限过滤：获取可见的素材分组ID
	sourceGroupIds := getSourceGroupIds(c, request.EnterpriseID, request.ThirdIds)

	// 3. 按 type 分组 source_id
	sourceIdsByType := make(map[int8][]int64)
	for _, item := range sourceIdTypes {
		sourceIdsByType[item.Type] = append(sourceIdsByType[item.Type], item.SourceID)
	}

	// 4. merge_id 扩展：图片(2)和朋友圈(7)需要通过 merge_id 让封面图也能被检索到
	expandMergeIds(c, request.EnterpriseID, sourceIdsByType)

	relationDao := sourceDao.NewGroupRelationDao()

	// 5. 按类型批量查询各素材详情（不分页）
	textIds := sourceIdsByType[source.TypeText]
	imageIds := sourceIdsByType[source.TypeImg]
	linkIds := sourceIdsByType[source.TypeLink]
	videoIds := sourceIdsByType[source.TypeVideo]
	friendIds := sourceIdsByType[source.TypeFriend]

	textList, _ := relationDao.GetTextOpenAll(c, request.EnterpriseID, textIds, sourceGroupIds)
	imageList, _ := relationDao.GetImageOpenAll(c, request.EnterpriseID, imageIds, sourceGroupIds)
	linkList, _ := relationDao.GetLinkOpenAll(c, request.EnterpriseID, linkIds, sourceGroupIds)
	videoList, _ := relationDao.GetVideoOpenAll(c, request.EnterpriseID, videoIds, sourceGroupIds)
	friendList, _ := relationDao.GetFriendOpenAll(c, request.EnterpriseID, friendIds, sourceGroupIds)

	// 6. 组合图子图查询 + BrandGoods 填充
	allSourceIds := collectAllSourceIds(textList, imageList, linkList, videoList, friendList)

	// 图片组合图子图
	imageList = assembleImageMerge(c, request.EnterpriseID, imageList, &allSourceIds)
	// 朋友圈组合图子图
	friendList = assembleFriendMerge(c, request.EnterpriseID, friendList, &allSourceIds)

	// 批量查询所有关联商品
	sourceSkuMap := getBrandGoodsMap(c, request.EnterpriseID, allSourceIds)

	// 填充 BrandGoods
	for i := range textList {
		textList[i].BrandGoods = sourceSkuMap[textList[i].ID]
	}
	for i := range linkList {
		linkList[i].BrandGoods = sourceSkuMap[linkList[i].ID]
	}
	for i := range videoList {
		videoList[i].BrandGoods = sourceSkuMap[videoList[i].ID]
	}
	for i := range imageList {
		imageList[i].BrandGoods = sourceSkuMap[imageList[i].ID]
		for si := range imageList[i].MergeImage {
			imageList[i].MergeImage[si].BrandGoods = sourceSkuMap[imageList[i].MergeImage[si].Id]
		}
	}
	for i := range friendList {
		friendList[i].BrandGoods = sourceSkuMap[friendList[i].ID]
		for si := range friendList[i].MergeImage {
			friendList[i].MergeImage[si].BrandGoods = sourceSkuMap[friendList[i].MergeImage[si].Id]
		}
	}

	result = openResp.SearchAllResponse{
		TextList:   textList,
		ImageList:  imageList,
		LinkList:   linkList,
		VideoList:  videoList,
		FriendList: friendList,
	}
	return result, nil
}

// expandMergeIds 对图片(2)和朋友圈(7)类型的source_id做merge_id扩展
// 子图关联了商品 → 通过 merge_id 找到同组封面图(sort=0)的 source_id 也加入集合
// 封面图关联了商品 → 其 source_id 已经在集合中，DAO查询时通过 OR merge_id IN 自动覆盖同组子图
func expandMergeIds(c app.Ctx, enterpriseId int64, sourceIdsByType map[int8][]int64) {
	for _, t := range []int8{source.TypeImg, source.TypeFriend} {
		ids, ok := sourceIdsByType[t]
		if !ok || len(ids) == 0 {
			continue
		}
		// 通过 GetMany 查出这些 source_id 对应的 source 记录，提取 merge_id
		sourceList, err := sourceDao.NewSourceDao().GetMany(c, map[string]any{"id": ids, "enterprise_id": enterpriseId})
		if err != nil || len(sourceList) == 0 {
			continue
		}

		var mergeIds []int64
		mergeIdSet := make(map[int64]bool)
		for _, val := range sourceList {
			if val.MergeID > 0 && !mergeIdSet[val.MergeID] {
				mergeIdSet[val.MergeID] = true
				mergeIds = append(mergeIds, val.MergeID)
			}
		}
		if len(mergeIds) == 0 {
			continue
		}

		// 查出这些 merge_id 对应的封面图(sort=0)的 source_id
		coverList, err := sourceDao.NewSourceDao().GetMany(c, map[string]any{"merge_id": mergeIds, "sort": 0, "enterprise_id": enterpriseId})
		if err != nil {
			continue
		}
		var coverIds []int64
		for _, cover := range coverList {
			coverIds = append(coverIds, cover.ID)
		}

		// 将封面图 id 追加到该类型的集合
		sourceIdsByType[t] = append(ids, coverIds...)
	}
}

// collectAllSourceIds 收集所有素材的 source_id（用于批量查 BrandGoods）
func collectAllSourceIds(
	textList []openResp.TextListResponse,
	imageList []openResp.ImageListResponse,
	linkList []openResp.LinkListResponse,
	videoList []openResp.VideoListResponse,
	friendList []openResp.FriendListResponse,
) []int64 {
	var ids []int64
	for _, v := range textList {
		ids = append(ids, v.ID)
	}
	for _, v := range linkList {
		ids = append(ids, v.ID)
	}
	for _, v := range videoList {
		ids = append(ids, v.ID)
	}
	for _, v := range imageList {
		ids = append(ids, v.ID)
	}
	for _, v := range friendList {
		ids = append(ids, v.ID)
	}
	return ids
}

// assembleImageMerge 组装图片组合图子图
func assembleImageMerge(c app.Ctx, enterpriseId int64, imageList []openResp.ImageListResponse, allSourceIds *[]int64) []openResp.ImageListResponse {
	if len(imageList) == 0 {
		return imageList
	}

	var mergeIds []int64
	for _, v := range imageList {
		if v.MergeId > 0 {
			mergeIds = append(mergeIds, v.MergeId)
		}
	}
	if len(mergeIds) == 0 {
		return imageList
	}

	// 查询组合图子图
	sourceMap := make(map[int64][]openResp.OpenMergeImage)
	sourceList, err := sourceDao.NewSourceDao().GetMany(c, map[string]any{"enterprise_id": enterpriseId, "merge_id": mergeIds})
	if err != nil {
		return imageList
	}

	var subSourceIds []int64
	for _, val := range sourceList {
		subSourceIds = append(subSourceIds, val.ID)
	}

	sourceImageMap := make(map[int64]source.ImageModel)
	sourceImageList, err := sourceDao.NewImageDao().GetMany(c, enterpriseId, map[string]any{"id": subSourceIds})
	if err != nil {
		return imageList
	}
	for _, image := range sourceImageList {
		sourceImageMap[image.ID] = image
	}

	for _, val := range sourceList {
		if img, ok := sourceImageMap[val.ID]; ok {
			sourceMap[val.MergeID] = append(sourceMap[val.MergeID], openResp.OpenMergeImage{
				Id:   img.ID,
				Name: img.Name,
				URL:  img.URL,
				Sort: val.Sort,
			})
		}
	}

	// 填充子图 source_id 到 allSourceIds（后续查 BrandGoods 用）
	*allSourceIds = append(*allSourceIds, subSourceIds...)

	for k, v := range imageList {
		if subImages, ok := sourceMap[v.MergeId]; ok {
			sort.Slice(subImages, func(i, j int) bool {
				return subImages[i].Sort < subImages[j].Sort
			})
			imageList[k].MergeImage = subImages
			if len(subImages) > 0 {
				imageList[k].URL = subImages[0].URL
			}
		}
	}

	return imageList
}

// assembleFriendMerge 组装朋友圈组合图子图（只取 FriendTypeImg 类型）
func assembleFriendMerge(c app.Ctx, enterpriseId int64, friendList []openResp.FriendListResponse, allSourceIds *[]int64) []openResp.FriendListResponse {
	if len(friendList) == 0 {
		return friendList
	}

	var mergeIds []int64
	for _, v := range friendList {
		if v.MergeId > 0 {
			mergeIds = append(mergeIds, v.MergeId)
		}
	}
	if len(mergeIds) == 0 {
		return friendList
	}

	// 查询组合图子图
	sourceMap := make(map[int64][]openResp.OpenMergeImage)
	sourceList, err := sourceDao.NewSourceDao().GetMany(c, map[string]any{"enterprise_id": enterpriseId, "merge_id": mergeIds})
	if err != nil {
		return friendList
	}

	var subSourceIds []int64
	for _, val := range sourceList {
		subSourceIds = append(subSourceIds, val.ID)
	}

	sourceFriendMap := make(map[int64]source.FriendModel)
	sourceFriendList, err := sourceDao.NewFriendDao().GetMany(c, enterpriseId, map[string]any{"id": subSourceIds})
	if err != nil {
		return friendList
	}
	for _, friend := range sourceFriendList {
		if friend.OriginalURL != "" && friend.Type == source.FriendTypeImg {
			sourceFriendMap[friend.ID] = friend
		}
	}

	for _, val := range sourceList {
		if friend, ok := sourceFriendMap[val.ID]; ok {
			sourceMap[val.MergeID] = append(sourceMap[val.MergeID], openResp.OpenMergeImage{
				Id:   friend.ID,
				Name: friend.Name,
				URL:  friend.OriginalURL,
				Sort: val.Sort,
			})
		}
	}

	// 填充子图 source_id 到 allSourceIds
	*allSourceIds = append(*allSourceIds, subSourceIds...)

	for k, v := range friendList {
		if subImages, ok := sourceMap[v.MergeId]; ok {
			sort.Slice(subImages, func(i, j int) bool {
				return subImages[i].Sort < subImages[j].Sort
			})
			friendList[k].MergeImage = subImages
		}
	}

	return friendList
}
