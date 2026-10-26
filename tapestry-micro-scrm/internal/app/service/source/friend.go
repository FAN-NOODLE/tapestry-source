package source

import (
	"errors"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	sourceDao "micro-scrm/internal/app/dao/source"
	sourceRequest "micro-scrm/internal/app/http/request/source"
	sourceResp "micro-scrm/internal/app/http/response/source"
	"micro-scrm/internal/app/models"
	"micro-scrm/internal/app/models/source"
	"micro-scrm/internal/app/service"
	"micro-scrm/internal/pkg/helper"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type FriendService struct {
	service.Base
	getSource    func(app.Ctx, map[string]any) (source.SourceModel, error)
	updateFriend func(app.Ctx, map[string]any, map[string]interface{}) error
	updateSource func(app.Ctx, sourceRequest.Request) error
}

func NewFriendService() FriendService {
	return FriendService{
		getSource: func(c app.Ctx, query map[string]any) (source.SourceModel, error) {
			return sourceDao.NewSourceDao().GetOne(c, query)
		},
		updateFriend: func(c app.Ctx, query map[string]any, column map[string]interface{}) error {
			return sourceDao.NewFriendDao().Update(c, query, column)
		},
		updateSource: func(c app.Ctx, request sourceRequest.Request) error {
			return NewService().Update(c, request)
		},
	}
}

func (s FriendService) Create(c app.Ctx, request sourceRequest.FriendCreateReq) error {
	request.TextContent = helper.CleanUserInput(request.TextContent)

	// 资源集合
	var originMediaList []sourceRequest.SourcesMediaInfo

	// 兼容只写文案不上传素材情况
	if request.TextContent != "" && len(request.Sources) == 0 {
		originMediaList = append(originMediaList, sourceRequest.SourcesMediaInfo{
			SourcesInfo: sourceRequest.SourcesInfo{SourceId: app.UID(), SourceType: source.TypeFriend},
			MediaInfo:   sourceRequest.MediaInfo{IsOrigin: 1},
		})
		// 不合并素材组
		request.ShowType = 0
	}

	for key, sourcePtr := range request.Sources {
		if sourcePtr.Url == "" {
			return errors.New("素材不能为空")
		}

		sourcesMediaInfo := sourceRequest.SourcesMediaInfo{
			SourcesInfo: sourceRequest.SourcesInfo{SourceId: app.UID(), SourceType: source.TypeFriend},
			MediaInfo: sourceRequest.MediaInfo{
				Url:      sourcePtr.Url,
				Name:     sourcePtr.Name,
				Ext:      strings.TrimLeft(filepath.Ext(sourcePtr.Name), "."),
				Size:     sourcePtr.Size,
				Sort:     int32(key),
				CoverUrl: sourcePtr.VideoCoverUrl,
				Type:     request.FileType,
				IsOrigin: 1, // 保持原图
			},
		}
		originMediaList = append(originMediaList, sourcesMediaInfo)
	}

	createFriendReq := sourceRequest.UploadMediaRequest{
		EnterpriseID: request.EnterpriseID,
		StartTime:    request.StartTime,
		EndTime:      request.EndTime,
		CategoryIds:  request.CategoryIds,
		SourceType:   source.TypeFriend,
		GroupId:      request.GroupId,
		Content:      request.TextContent,
		ShowType:     request.ShowType,
		FileType:     request.FileType,
	}
	// 判断是否为组合展示,组合展示需要新增合并id

	mergeId := app.UID()
	err := sourceDao.NewFriendDao().CreateFriend(c, originMediaList, createFriendReq, 0, request.AdminId, mergeId)
	if err != nil {
		return err
	}

	// 保存关联商品
	for i, sourcePtr := range request.Sources {
		if len(sourcePtr.RelationGoods) > 0 {
			_ = s.UpdateRelationGoods(c, uint64(originMediaList[i].SourceId), request.EnterpriseID, sourcePtr.RelationGoods)
		}
	}

	return nil
}

func (s FriendService) Detail(c app.Ctx, request sourceRequest.FriendDetailReq) ([]sourceResp.FriendDetail, error) {
	res, err := sourceDao.NewSourceDao().GetSourceFriendDetail(c, request.SourceID, request.EnterpriseID)
	if err != nil {
		return nil, err
	}

	var resp []sourceResp.FriendDetail
	for _, list := range res {
		resp = append(resp, sourceResp.FriendDetail{
			FriendDetail: list,
			MergeImage:   make([]sourceResp.MergeImage, 0),
		})
	}

	sourceMap := make(map[int64][]sourceResp.MergeImage)

	for _, detail := range resp {
		if detail.MergeId == 0 {
			continue
		}

		sourceList, err := sourceDao.NewSourceDao().GetMany(c, map[string]any{"enterprise_id": request.EnterpriseID, "merge_id": detail.MergeId})
		if err != nil {
			return nil, err
		}

		var sourceIds []int64
		for _, val := range sourceList {
			sourceIds = append(sourceIds, val.ID)
		}

		sourceImageMap := make(map[int64]source.FriendModel)
		sourceImageList, err := sourceDao.NewFriendDao().GetMany(c, int64(request.EnterpriseID), map[string]any{"id": sourceIds})
		if err != nil {
			return nil, err
		}

		for _, image := range sourceImageList {
			if image.OriginalURL != "" && image.Type == source.FriendTypeImg {
				sourceImageMap[image.ID] = image
			}
		}

		for _, val := range sourceList {
			if _, ok := sourceImageMap[val.ID]; !ok {
				continue
			}
			tmpImage := sourceResp.MergeImage{
				Id:   sourceImageMap[val.ID].ID,
				Name: sourceImageMap[val.ID].Name,
				URL:  sourceImageMap[val.ID].OriginalURL,
				Sort: val.Sort,
			}
			sourceMap[val.MergeID] = append(sourceMap[val.MergeID], tmpImage)
		}
	}

	// 查询关联商品
	var allSourceIds []int64
	for _, detail := range resp {
		allSourceIds = append(allSourceIds, detail.SourceId)
	}
	for _, images := range sourceMap {
		for _, img := range images {
			allSourceIds = append(allSourceIds, img.Id)
		}
	}

	goodsMap := make(map[int64][]sourceRequest.BrandGoods)
	if len(allSourceIds) > 0 {
		goods, _ := sourceDao.NewRelationGoodsDao().GetMany(ctx.Wrap(c), "source_id IN ? AND type = ?", allSourceIds, source.TypeFriend)
		for _, model := range goods {
			tmp := sourceRequest.BrandGoods{
				Name:    model.Name,
				Sku:     model.SKU,
				Upc:     model.UPC,
				Url:     model.Url,
				Color:   model.Color,
				AddType: model.AddType,
				Score:   model.Score,
			}
			goodsMap[model.SourceID] = append(goodsMap[model.SourceID], tmp)
		}
	}

	for k, v := range resp {
		// 组装合并图片
		if sourceImage, ok := sourceMap[v.MergeId]; ok {
			sort.Slice(sourceImage, func(i, j int) bool {
				return sourceImage[i].Sort < sourceImage[j].Sort
			})
			// 为每个图片填充关联商品
			for idx := range sourceImage {
				if rg, ok := goodsMap[sourceImage[idx].Id]; ok {
					sourceImage[idx].RelationGoods = rg
				}
			}
			resp[k].MergeImage = append(resp[k].MergeImage, sourceImage...)
		}
		// 如果是单个图片（非合并组），也填充关联商品
		if v.SourceId > 0 {
			if rg, ok := goodsMap[v.SourceId]; ok {
				resp[k].RelationGoods = append(resp[k].RelationGoods, rg...)
			}
		}
	}

	return resp, nil
}

func (s FriendService) Update(c app.Ctx, request sourceRequest.FriendUpdateReq) error {
	request.TextContent = helper.CleanUserInput(request.TextContent)

	// 查询当前素材的merge_id
	sourceInfo, err := s.getSource(c, map[string]any{"id": request.ID, "enterprise_id": request.EnterpriseID})
	if err != nil {
		return err
	}

	if sourceInfo.ID == 0 {
		return errors.New("数据记录不存在")
	}
	updateInfo := sourceRequest.Request{
		ID:           request.ID,
		EnterpriseID: request.EnterpriseID,
		StartTime:    request.StartTime,
		EndTime:      request.EndTime,
		CategoryIds:  request.CategoryIds,
		SourceType:   request.SourceType,
		GroupId:      request.GroupId,
	}
	// 文本不为空且素材集合为空,则更新所有素材
	if request.TextContent != "" && len(request.Sources) == 0 {
		sourceIds := []int64{sourceInfo.ID}
		if sourceInfo.MergeID > 0 {
			// 获取merge_id对应的所有素材id
			sourceList, err := sourceDao.NewSourceDao().GetMany(c, map[string]any{"id": request.ID, "merge_id": sourceInfo.MergeID})
			if err != nil {
				return err
			}
			sourceIds = []int64{}
			for _, model := range sourceList {
				sourceIds = append(sourceIds, model.ID)
			}
		}

		query := map[string]any{
			"id":            sourceIds,
			"enterprise_id": request.EnterpriseID,
		}
		column := map[string]interface{}{
			"text_content": request.TextContent,
			"updated_at":   time.Now().Unix(),
			"original_url": "",
		}

		// 更新所有merge_id对应资源id的文本内容
		if err = s.updateFriend(c, query, column); err != nil {
			return err
		}
		updateInfo.Ids = sourceIds
	}

	// 更新素材
	for key, val := range request.Sources {
		query := map[string]any{
			"id":            val.Id,
			"enterprise_id": request.EnterpriseID,
		}
		column := map[string]interface{}{
			"original_url": val.Url,
			"text_content": request.TextContent,
			"updated_at":   time.Now().Unix(),
		}

		// 更新主表排序
		err = sourceDao.NewSourceDao().Update(c, query, map[string]any{"sort": key})
		if err != nil {
			return err
		}

		// 更新子表数据
		err = sourceDao.NewFriendDao().Update(c, query, column)
		if err != nil {
			return err
		}

		// 更新关联商品
		_ = s.UpdateRelationGoods(c, uint64(val.Id), request.EnterpriseID, val.RelationGoods)

		updateInfo.Ids = append(updateInfo.Ids, val.Id)
	}

	// 更新主表 分组和标签
	request.SourceType = source.TypeFriend
	err = s.updateSource(c, updateInfo)

	return err
}

// UpdateRelationGoods 更新朋友圈的关联商品
func (s FriendService) UpdateRelationGoods(c app.Ctx, sourceId, enterpriseId uint64, goodsList []sourceRequest.BrandGoods) error {
	relationGoodsDao := sourceDao.NewRelationGoodsDao()

	// 删除旧的关联商品
	_ = relationGoodsDao.DeleteData(c, "source_id = ? AND type = ?", sourceId, source.TypeFriend)

	// 批量插入新的关联商品
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
			SourceID: int64(sourceId),
			Type:     source.TypeFriend,
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
		_ = relationGoodsDao.CreateInBatchesData(c, goodsModels)
	}
	return nil
}
