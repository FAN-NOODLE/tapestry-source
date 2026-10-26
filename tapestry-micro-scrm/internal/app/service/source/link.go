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

type LinkService struct {
	sourceDao *sourceDao.Dao
	linkDao   *sourceDao.LinkDao
}

func NewLinkService() LinkService {
	return LinkService{
		sourceDao: sourceDao.NewSourceDao(),
		linkDao:   sourceDao.NewLinkDao(),
	}
}
func (s LinkService) Detail(c app.Ctx, request sourceRequest.DetailRequest) ([]sourceResp.SourceDetail, error) {
	detailList, err := s.sourceDao.GetSourceLinkDetail(c, request.SourceID, request.EnterpriseID)
	if err != nil || len(detailList) == 0 {
		return nil, err
	}
	resDetailList := make([]sourceResp.SourceDetail, len(detailList))

	var sourceSkuMap = make(map[int64][]sourceRequest.BrandGoods)
	goods, _ := sourceDao.NewRelationGoodsDao().GetMany(ctx.Wrap(c), "source_id = ?", detailList[0].SourceId)
	if len(goods) > 0 {
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
			sourceSkuMap[model.SourceID] = append(sourceSkuMap[model.SourceID], tmp)
		}
	}

	for i, relationDetail := range detailList {
		// 如果使用第一种方式（已分配长度）：
		resDetailList[i] = sourceResp.SourceDetail{
			LinkCategoryRelationDetail: relationDetail,
		}

		if goods, ok := sourceSkuMap[relationDetail.SourceId]; ok {
			resDetailList[i].RelationGoods = goods
		}
	}

	return resDetailList, nil
}

func (s LinkService) Update(c app.Ctx, request sourceRequest.LinkUpdateRequest) error {
	// 更新链接标题，描述
	query := map[string]interface{}{
		"id":            request.ID,
		"enterprise_id": request.EnterpriseID,
	}
	column := map[string]interface{}{
		"title":       request.Title,
		"description": request.Description,
		"cover_url":   request.CoverUrl,
		"updated_at":  time.Now().Unix(),
	}
	if err := s.linkDao.Update(c, query, column); err != nil {
		return err
	}
	request.SourceType = source.TypeLink
	err := NewService().Update(c, request.Request)
	return err
}
