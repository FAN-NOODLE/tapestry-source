package source

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	sourceDao "micro-scrm/internal/app/dao/source"
	sourceRequest "micro-scrm/internal/app/http/request/source"
	sourceResp "micro-scrm/internal/app/http/response/source"
	"micro-scrm/internal/app/models/source"
	"time"
)

type VideoService struct {
	videoDao  *sourceDao.VideoDao
	sourceDao *sourceDao.Dao
}

func NewVideoService() VideoService {
	return VideoService{
		videoDao:  sourceDao.NewVideoDao(),
		sourceDao: sourceDao.NewSourceDao(),
	}
}

func (s VideoService) Detail(c app.Ctx, request sourceRequest.DetailRequest) ([]sourceResp.VideoDetail, error) {
	res, err := s.sourceDao.GetSourceVideoDetail(c, request.SourceID, request.EnterpriseID)
	if err != nil {
		return nil, err
	}

	var resp []sourceResp.VideoDetail
	for _, detail := range res {
		resp = append(resp, sourceResp.VideoDetail{
			VideoCategoryRelationDetail: detail,
		})
	}

	// 查询关联商品
	var sourceIds []int64
	for _, detail := range resp {
		sourceIds = append(sourceIds, detail.SourceId)
	}

	goodsMap := make(map[int64][]sourceRequest.BrandGoods)
	if len(sourceIds) > 0 {
		goods, _ := sourceDao.NewRelationGoodsDao().GetMany(ctx.Wrap(c), "source_id IN ? AND type = ?", sourceIds, source.TypeVideo)
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
		if rg, ok := goodsMap[v.SourceId]; ok {
			resp[k].RelationGoods = rg
		}
	}

	return resp, nil
}

func (s VideoService) Update(ctx app.Ctx, request sourceRequest.VideoUpdateRequest) error {

	// 更新视频文件名称
	query := map[string]interface{}{
		"id":            request.ID,
		"enterprise_id": request.EnterpriseID,
	}
	column := map[string]interface{}{
		"name":       request.Name,
		"updated_at": time.Now().Unix(),
	}
	if err := s.videoDao.Update(ctx, query, column); err != nil {

		return err
	}

	request.SourceType = source.TypeVideo
	err := NewService().Update(ctx, request.Request)
	return err
}
