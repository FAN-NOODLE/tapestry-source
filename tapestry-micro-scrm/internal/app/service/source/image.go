package source

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	sourceDao "micro-scrm/internal/app/dao/source"
	sourceRequest "micro-scrm/internal/app/http/request/source"
	sourceResp "micro-scrm/internal/app/http/response/source"
	"micro-scrm/internal/app/models"
	"micro-scrm/internal/app/models/source"
	"sort"
	"time"
)

type ImageService struct {
	sourceDao *sourceDao.Dao
}

func NewImageService() ImageService {
	return ImageService{
		sourceDao: sourceDao.NewSourceDao(),
	}
}

func (s ImageService) Detail(c app.Ctx, request sourceRequest.DetailRequest) ([]sourceResp.ImageDetail, error) {
	res, err := s.sourceDao.GetSourceImageDetail(c, request.SourceID, request.EnterpriseID)
	if err != nil {
		return nil, err
	}

	var resp []sourceResp.ImageDetail
	for _, list := range res {
		resp = append(resp, sourceResp.ImageDetail{
			ImageCategoryRelationDetail: list,
			MergeImage:                  make([]sourceResp.MergeImage, 0),
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

		sourceImageMap := make(map[int64]source.ImageModel)
		sourceImageList, err := sourceDao.NewImageDao().GetMany(c, int64(request.EnterpriseID), map[string]any{"id": sourceIds})
		if err != nil {
			return nil, err
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

			sourceMap[val.MergeID] = append(sourceMap[val.MergeID], tmpImage)
		}
	}

	// 查询关联商品
	var allSourceIds []int64
	for _, detail := range resp {
		if detail.MergeId == 0 {
			allSourceIds = append(allSourceIds, detail.SourceId)
		}
	}
	for _, images := range sourceMap {
		for _, img := range images {
			allSourceIds = append(allSourceIds, img.Id)
		}
	}

	goodsMap := make(map[int64][]sourceRequest.BrandGoods)
	if len(allSourceIds) > 0 {
		goods, _ := sourceDao.NewRelationGoodsDao().GetMany(ctx.Wrap(c), "source_id IN ? AND type = ?", allSourceIds, source.TypeImg)
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
		// 如果是单个图片则把关联商品放外面
		if v.MergeId == 0 && v.SourceId > 0 {
			if rg, ok := goodsMap[v.SourceId]; ok {
				resp[k].RelationGoods = rg
			}
		}
	}

	return resp, nil
}

func (s ImageService) Update(c app.Ctx, request sourceRequest.ImageUpdateRequest) error {
	request.SourceType = source.TypeImg

	// 更新图片排序
	for k, val := range request.Sources {
		query := map[string]any{
			"id":            val.Id,
			"enterprise_id": request.EnterpriseID,
		}
		column := map[string]interface{}{
			"sort":       k,
			"updated_at": time.Now().Unix(),
		}
		request.Ids = append(request.Ids, val.Id)
		if err := sourceDao.NewSourceDao().Update(c, query, column); err != nil {
			return err
		}

		// 更新关联商品
		_ = s.UpdateRelationGoods(c, uint64(val.Id), request.EnterpriseID, val.RelationGoods, source.TypeImg)
	}

	err := NewService().Update(c, request.Request)
	return err
}

// UpdateRelationGoods 更新图片/朋友圈/视频的关联商品
func (s ImageService) UpdateRelationGoods(c app.Ctx, sourceId, enterpriseId uint64, goodsList []sourceRequest.BrandGoods, sourceType int8) error {
	relationGoodsDao := sourceDao.NewRelationGoodsDao()

	// 删除旧的关联商品
	_ = relationGoodsDao.DeleteData(c, "source_id = ? AND type = ?", sourceId, sourceType)

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
			Type:     sourceType,
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
