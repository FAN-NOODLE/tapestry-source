package source

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"micro-scrm/internal/app"
	sourcedao "micro-scrm/internal/app/dao/source"
	"micro-scrm/internal/pkg/sdk/baozun"
)

type coachKOSGoodsDependencies struct {
	cacheGet       func(app.Ctx, string) (baozun.ProductBriefCache, bool, error)
	cacheSet       func(app.Ctx, string, baozun.ProductBriefCache, time.Duration) error
	search         func(app.Ctx, baozun.SearchParams, *baozun.DataResponse) error
	updateRelation func(app.Ctx, map[string]any, map[string]any) error
	deleteRelation func(app.Ctx, map[string]any) error
	now            func() int64
	logError       func(app.Ctx, string, ...any)
}

var coachKOSGoodsDeps = coachKOSGoodsDependencies{
	cacheGet: func(c app.Ctx, key string) (baozun.ProductBriefCache, bool, error) {
		payload, err := app.Redis().Get(c, key).Bytes()
		if errors.Is(err, redis.Nil) {
			return baozun.ProductBriefCache{}, false, nil
		}
		if err != nil {
			return baozun.ProductBriefCache{}, false, err
		}
		brief, found := coachKOSCachedProductBrief(payload)
		return brief, found, nil
	},
	cacheSet: func(c app.Ctx, key string, brief baozun.ProductBriefCache, ttl time.Duration) error {
		payload, err := json.Marshal(brief)
		if err != nil {
			return err
		}
		return app.Redis().Set(c, key, payload, ttl).Err()
	},
	search: func(c app.Ctx, params baozun.SearchParams, response *baozun.DataResponse) error {
		return baozun.SearchResources(c, params, response)
	},
	updateRelation: func(c app.Ctx, query, updates map[string]any) error {
		return sourcedao.NewRelationGoodsDao().UpdateData(c, query, updates)
	},
	deleteRelation: func(c app.Ctx, query map[string]any) error {
		return sourcedao.NewRelationGoodsDao().ForceDelete(c, query)
	},
	now: func() int64 { return time.Now().Unix() },
	logError: func(c app.Ctx, format string, args ...any) {
		app.Logger().Error(c).Msgf(format, args...)
	},
}

func coachKOSCachedProductBrief(payload []byte) (baozun.ProductBriefCache, bool) {
	var brief *baozun.ProductBriefCache
	if err := json.Unmarshal(payload, &brief); err != nil || brief == nil {
		return baozun.ProductBriefCache{}, false
	}
	return *brief, true
}

// EnrichCoachKOSRelationGoods 使用宝尊商品信息补全已存在的 KOS 链接商品关系。
func EnrichCoachKOSRelationGoods(c app.Ctx, tasks []CoachKOSGoodsTask) {
	for _, task := range tasks {
		for _, productCode := range task.ProductCodes {
			if err := enrichCoachKOSRelationGoods(c, task.EnterpriseID, task.SourceID, productCode); err != nil {
				coachKOSGoodsDeps.logError(c,
					"KOS 商品补全失败：企业ID=%d，素材ID=%d，productCode=%s，错误=%v",
					task.EnterpriseID, task.SourceID, productCode, err,
				)
			}
		}
	}
}

func enrichCoachKOSRelationGoods(c app.Ctx, enterpriseID, sourceID int64, productCode string) error {
	query := map[string]any{
		"enterprise_id": enterpriseID,
		"source_id":     sourceID,
		"type":          int8(3),
		"sku_code":      coachKOSSKUCode(productCode),
		"deleted_at":    0,
	}
	cacheKey := fmt.Sprintf(baozun.ProductInfoRedisKey, enterpriseID, productCode)
	brief, found, err := coachKOSGoodsDeps.cacheGet(c, cacheKey)
	if err != nil {
		return fmt.Errorf("读取商品缓存失败: %w", err)
	}
	if found && brief.ProductNo == "" {
		found = false
	}

	if !found {
		var response baozun.DataResponse
		err = coachKOSGoodsDeps.search(c, baozun.SearchParams{
			EnterpriseId:            enterpriseID,
			ProductCode:             []string{productCode},
			IsLike:                  false,
			IncludeNoneImageProduct: true,
			Size:                    1,
		}, &response)
		if err != nil {
			return fmt.Errorf("查询宝尊商品失败: %w", err)
		}

		var ok bool
		brief, ok = coachKOSProductBrief(productCode, response)
		if !ok {
			if err = coachKOSGoodsDeps.deleteRelation(c, query); err != nil {
				return fmt.Errorf("删除未补全的商品关系失败: %w", err)
			}
			return nil
		}
		if err = coachKOSGoodsDeps.cacheSet(c, cacheKey, brief, 24*time.Hour); err != nil {
			coachKOSGoodsDeps.logError(c,
				"KOS 商品缓存写入失败，继续补全商品关系：企业ID=%d，素材ID=%d，productCode=%s，错误=%v",
				enterpriseID, sourceID, productCode, err,
			)
		}
	}

	updates := map[string]any{
		"upc":        brief.Upc,
		"sku":        brief.ProductNo,
		"name":       brief.Name,
		"url":        brief.Url,
		"color":      brief.Color,
		"updated_at": coachKOSGoodsDeps.now(),
	}
	if err = coachKOSGoodsDeps.updateRelation(c, query, updates); err != nil {
		return fmt.Errorf("更新商品关系失败: %w", err)
	}
	return nil
}

func coachKOSProductBrief(productCode string, data baozun.DataResponse) (baozun.ProductBriefCache, bool) {
	if len(data.Data) == 0 {
		return baozun.ProductBriefCache{}, false
	}

	var productIndex = -1
	for index := range data.Data {
		if strings.TrimSpace(data.Data[index].Master.ProductCode) == strings.TrimSpace(productCode) {
			productIndex = index
			break
		}
	}
	if productIndex < 0 {
		return baozun.ProductBriefCache{}, false
	}

	product := data.Data[productIndex]
	brief := baozun.ProductBriefCache{
		ProductNo: product.Master.Properties.ProductNo,
		Name:      product.Master.Properties.ProductName,
	}
	if len(product.Master.Fragments) > 0 {
		brief.Color = product.Master.Fragments[0].Properties.Color
		if brief.Name == "" {
			brief.Name = product.Master.Fragments[0].Properties.ProductNameEn
		}
	}

	for _, variant := range product.Variants {
		if strings.TrimSpace(variant.VariantCode) != strings.TrimSpace(productCode) {
			continue
		}
		for _, fragment := range variant.Fragments {
			if fragment.Properties.Upc != "" {
				brief.Upc = fragment.Properties.Upc
				break
			}
		}
		if brief.Upc != "" {
			break
		}
	}
	if brief.Upc == "" {
		for _, variant := range product.Variants {
			for _, fragment := range variant.Fragments {
				if fragment.Properties.Upc != "" {
					brief.Upc = fragment.Properties.Upc
					break
				}
			}
			if brief.Upc != "" {
				break
			}
		}
	}

	for _, resource := range product.ExtData.Resource {
		if strings.Contains(resource.Name, "_a0") {
			brief.Url = strings.SplitN(resource.Url, "?", 2)[0]
			break
		}
	}

	return brief, true
}
