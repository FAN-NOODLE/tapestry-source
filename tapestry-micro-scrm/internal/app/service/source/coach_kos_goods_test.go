package source

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"micro-scrm/internal/app"
	"micro-scrm/internal/pkg/sdk/baozun"
)

func TestCoachKOSGoodsProductBrief(t *testing.T) {
	const productNo = "  CV758_B4NQ4  "
	var response baozun.DataResponse
	if err := json.Unmarshal([]byte(`{
		"data": [{
			"extData": {"resource": {
				"other": {"name": "CV758_other.jpg", "url": "https://img.example/other.jpg"},
				"main": {"name": "CV758_a0.jpg", "url": "https://img.example/CV758_a0.jpg?width=1200&token=secret"}
			}},
			"master": {
				"productCode": "CV758_B4NQ4",
				"properties": {"productName": "经典手袋", "productNo": "CV758"},
				"fragments": [{"properties": {"productName_en": "Classic Bag", "color": "B4NQ4"}}]
			},
			"variants": [
				{"variantCode": "OTHER", "fragments": [{"properties": {"Upc": "wrong-upc"}}]},
				{"variantCode": "  CV758_B4NQ4  ", "fragments": [{"properties": {"Upc": "right-upc"}}]}
			]
		}]
	}`), &response); err != nil {
		t.Fatalf("构造宝尊响应失败: %v", err)
	}

	brief, ok := coachKOSProductBrief(productNo, response)
	if !ok {
		t.Fatal("期望成功解析商品")
	}
	if brief.Name != "经典手袋" {
		t.Fatalf("商品名 = %q, 期望 %q", brief.Name, "经典手袋")
	}
	if brief.Upc != "right-upc" {
		t.Fatalf("UPC = %q, 期望匹配 variant 的 UPC %q", brief.Upc, "right-upc")
	}
	if brief.Url != "https://img.example/CV758_a0.jpg" {
		t.Fatalf("URL = %q, 期望移除查询参数", brief.Url)
	}
	if brief.Color != "B4NQ4" {
		t.Fatalf("颜色 = %q, 期望主商品 fragment 颜色", brief.Color)
	}
	briefJSON, err := json.Marshal(brief)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(briefJSON), `"product_no":"CV758"`) {
		t.Fatalf("商品简要信息未保存 PIM productNo：%s", briefJSON)
	}
}

func TestCoachKOSGoodsProductBriefFallsBackToEnglishName(t *testing.T) {
	var response baozun.DataResponse
	if err := json.Unmarshal([]byte(`{
		"data": [{
			"master": {
				"productCode": "CV758",
				"properties": {"productName": "", "productNo": "CV758"},
				"fragments": [{"properties": {"productName_en": "Classic Bag", "color": "BLACK"}}]
			},
			"variants": [{"variantCode": "CV758", "fragments": [{"properties": {"Upc": "upc-1"}}]}]
		}]
	}`), &response); err != nil {
		t.Fatalf("构造宝尊响应失败: %v", err)
	}

	brief, ok := coachKOSProductBrief("CV758", response)
	if !ok {
		t.Fatal("期望成功解析商品")
	}
	if brief.Name != "Classic Bag" {
		t.Fatalf("商品名 = %q, 期望英文回退名", brief.Name)
	}
}

func TestCoachKOSGoodsProductBriefFallsBackToFirstNonEmptyVariantUPC(t *testing.T) {
	var response baozun.DataResponse
	if err := json.Unmarshal([]byte(`{
		"data": [{
			"master": {"productCode": "MATCH", "properties": {"productNo": "MATCH"}},
			"variants": [
				{"variantCode": "MATCH", "fragments": [{"properties": {"Upc": ""}}]},
				{"variantCode": "SECOND", "fragments": [{"properties": {"Upc": "fallback-upc"}}]}
			]
		}]
	}`), &response); err != nil {
		t.Fatalf("构造宝尊响应失败: %v", err)
	}

	brief, ok := coachKOSProductBrief("MATCH", response)
	if !ok {
		t.Fatal("期望成功解析商品")
	}
	if brief.Upc != "fallback-upc" {
		t.Fatalf("UPC = %q, 期望无精确匹配时回退首个 variant UPC", brief.Upc)
	}
}

func TestCoachKOSGoodsProductBriefScansAllFragmentsForExactVariant(t *testing.T) {
	var response baozun.DataResponse
	payload := `{"data":[{"master":{"productCode":"MATCH","properties":{"productName":"手袋","productNo":"MATCH"}},"variants":[{"variantCode":"MATCH","fragments":[{"properties":{"Upc":""}},{"properties":{"Upc":"exact-upc"}}]},{"variantCode":"OTHER","fragments":[{"properties":{"Upc":"fallback-upc"}}]}]}]}`
	if err := json.Unmarshal([]byte(payload), &response); err != nil {
		t.Fatal(err)
	}
	brief, ok := coachKOSProductBrief("MATCH", response)
	if !ok || brief.Upc != "exact-upc" {
		t.Fatalf("UPC = %q, ok=%v，期望扫描精确 variant 的全部 fragments", brief.Upc, ok)
	}
}

func TestCoachKOSGoodsCacheWriteFailureStillUpdatesRelation(t *testing.T) {
	original := coachKOSGoodsDeps
	t.Cleanup(func() { coachKOSGoodsDeps = original })
	var updated bool
	var logs []string
	coachKOSGoodsDeps = coachKOSGoodsDependencies{
		cacheGet: func(app.Ctx, string) (baozun.ProductBriefCache, bool, error) {
			return baozun.ProductBriefCache{}, false, nil
		},
		cacheSet: func(app.Ctx, string, baozun.ProductBriefCache, time.Duration) error {
			return errors.New("redis unavailable")
		},
		search: func(_ app.Ctx, _ baozun.SearchParams, response *baozun.DataResponse) error {
			*response = coachKOSGoodsResponse(t, "CV758", "手袋", "", "BLACK", "upc-1")
			return nil
		},
		updateRelation: func(app.Ctx, map[string]any, map[string]any) error {
			updated = true
			return nil
		},
		now: func() int64 { return 123 },
		logError: func(_ app.Ctx, format string, args ...any) {
			logs = append(logs, fmt.Sprintf(format, args...))
		},
	}

	EnrichCoachKOSRelationGoods(nil, []CoachKOSGoodsTask{{EnterpriseID: 1, SourceID: 2, ProductCodes: []string{"CV758"}}})
	if !updated {
		t.Fatal("缓存写入失败后仍应更新商品关系")
	}
	if len(logs) != 1 || !strings.Contains(logs[0], "缓存写入失败") {
		t.Fatalf("缓存错误日志 = %v", logs)
	}
}

func TestCoachKOSGoodsProductBriefRejectsEmptyResponse(t *testing.T) {
	if _, ok := coachKOSProductBrief("CV758", baozun.DataResponse{}); ok {
		t.Fatal("空响应不应解析为有效商品")
	}
}

func TestCoachKOSGoods查无商品时删除占位关系(t *testing.T) {
	original := coachKOSGoodsDeps
	t.Cleanup(func() { coachKOSGoodsDeps = original })

	var deletedQuery map[string]any
	coachKOSGoodsDeps = coachKOSGoodsDependencies{
		cacheGet: func(app.Ctx, string) (baozun.ProductBriefCache, bool, error) {
			return baozun.ProductBriefCache{}, false, nil
		},
		cacheSet: func(app.Ctx, string, baozun.ProductBriefCache, time.Duration) error {
			t.Fatal("查无商品时不应写入缓存")
			return nil
		},
		search: func(_ app.Ctx, _ baozun.SearchParams, response *baozun.DataResponse) error {
			*response = baozun.DataResponse{}
			return nil
		},
		updateRelation: func(app.Ctx, map[string]any, map[string]any) error {
			t.Fatal("查无商品时不应更新占位关系")
			return nil
		},
		deleteRelation: func(_ app.Ctx, query map[string]any) error {
			deletedQuery = query
			return nil
		},
		now: func() int64 { return 123 },
		logError: func(_ app.Ctx, format string, args ...any) {
			t.Fatalf("查无商品并成功删除占位关系不应记录错误："+format, args...)
		},
	}

	EnrichCoachKOSRelationGoods(nil, []CoachKOSGoodsTask{{
		EnterpriseID: 1001,
		SourceID:     2002,
		ProductCodes: []string{"CV758_B4NQ4"},
	}})

	assertCoachKOSGoodsMapValue(t, deletedQuery, "enterprise_id", int64(1001))
	assertCoachKOSGoodsMapValue(t, deletedQuery, "source_id", int64(2002))
	assertCoachKOSGoodsMapValue(t, deletedQuery, "type", int8(3))
	assertCoachKOSGoodsMapValue(t, deletedQuery, "sku_code", "CV758 B4NQ4")
	assertCoachKOSGoodsMapValue(t, deletedQuery, "deleted_at", 0)
}

func TestCoachKOSGoodsMalformedCacheIsMiss(t *testing.T) {
	if _, found := coachKOSCachedProductBrief([]byte("{")); found {
		t.Fatal("损坏的商品缓存应视为未命中")
	}
}

func TestCoachKOSGoodsNullCacheRefreshesFromBaozun(t *testing.T) {
	original := coachKOSGoodsDeps
	t.Cleanup(func() { coachKOSGoodsDeps = original })

	const productNo = "CV758"
	var (
		searched      bool
		refreshed     bool
		updatedValues map[string]any
	)
	coachKOSGoodsDeps = coachKOSGoodsDependencies{
		cacheGet: func(_ app.Ctx, _ string) (baozun.ProductBriefCache, bool, error) {
			brief, found := coachKOSCachedProductBrief([]byte("null"))
			return brief, found, nil
		},
		cacheSet: func(_ app.Ctx, _ string, brief baozun.ProductBriefCache, ttl time.Duration) error {
			refreshed = true
			if ttl != 24*time.Hour {
				t.Fatalf("缓存 TTL = %v, 期望 24h", ttl)
			}
			if brief.Name != "宝尊商品" || brief.Upc != "fresh-upc" {
				t.Fatalf("刷新缓存商品信息不正确: %+v", brief)
			}
			return nil
		},
		search: func(_ app.Ctx, params baozun.SearchParams, response *baozun.DataResponse) error {
			searched = true
			if len(params.ProductCode) != 1 || params.ProductCode[0] != productNo {
				t.Fatalf("宝尊查询 productCode = %#v, 期望 [%q]", params.ProductCode, productNo)
			}
			*response = coachKOSGoodsResponse(t, productNo, "宝尊商品", "", "BLACK", "fresh-upc")
			return nil
		},
		updateRelation: func(_ app.Ctx, _ map[string]any, updates map[string]any) error {
			updatedValues = updates
			return nil
		},
		now: func() int64 { return 123 },
		logError: func(_ app.Ctx, format string, args ...any) {
			t.Fatalf("null 缓存刷新成功不应记录错误日志: "+format, args...)
		},
	}

	EnrichCoachKOSRelationGoods(nil, []CoachKOSGoodsTask{{
		EnterpriseID: 1,
		SourceID:     2,
		ProductCodes: []string{productNo},
	}})

	if !searched {
		t.Fatal("null 缓存应视为未命中并查询宝尊")
	}
	if !refreshed {
		t.Fatal("null 缓存查询宝尊后应刷新缓存")
	}
	assertCoachKOSGoodsMapValue(t, updatedValues, "name", "宝尊商品")
	assertCoachKOSGoodsMapValue(t, updatedValues, "upc", "fresh-upc")
}

func TestCoachKOSGoodsEnrichesCacheMissAndUpdatesCurrentLinkRelation(t *testing.T) {
	original := coachKOSGoodsDeps
	t.Cleanup(func() { coachKOSGoodsDeps = original })

	const (
		enterpriseID = int64(1001)
		sourceID     = int64(2002)
		productNo    = "CV758_B4NQ4"
		now          = int64(1783987200)
	)
	cacheKey := fmt.Sprintf(baozun.ProductInfoRedisKey, enterpriseID, productNo)
	var (
		searched baozun.SearchParams
		cached   baozun.ProductBriefCache
		cacheTTL time.Duration
		query    map[string]any
		updates  map[string]any
	)
	coachKOSGoodsDeps = coachKOSGoodsDependencies{
		cacheGet: func(_ app.Ctx, key string) (baozun.ProductBriefCache, bool, error) {
			if key != cacheKey {
				t.Fatalf("缓存 key = %q, 期望 %q", key, cacheKey)
			}
			return baozun.ProductBriefCache{}, false, nil
		},
		cacheSet: func(_ app.Ctx, key string, brief baozun.ProductBriefCache, ttl time.Duration) error {
			if key != cacheKey {
				t.Fatalf("写缓存 key = %q, 期望 %q", key, cacheKey)
			}
			cached, cacheTTL = brief, ttl
			return nil
		},
		search: func(_ app.Ctx, params baozun.SearchParams, response *baozun.DataResponse) error {
			searched = params
			*response = coachKOSGoodsResponse(t, productNo, "经典手袋", "Classic Bag", "B4NQ4", "upc-1")
			return nil
		},
		updateRelation: func(_ app.Ctx, gotQuery, gotUpdates map[string]any) error {
			query, updates = gotQuery, gotUpdates
			return nil
		},
		now: func() int64 { return now },
		logError: func(_ app.Ctx, format string, args ...any) {
			t.Fatalf("成功流程不应记录错误日志: "+format, args...)
		},
	}

	EnrichCoachKOSRelationGoods(nil, []CoachKOSGoodsTask{{
		EnterpriseID: enterpriseID,
		SourceID:     sourceID,
		ProductCodes: []string{productNo},
	}})

	if searched.EnterpriseId != enterpriseID || len(searched.ProductCode) != 1 || searched.ProductCode[0] != productNo ||
		searched.ProductNo != "" || searched.IsLike ||
		!searched.IncludeNoneImageProduct || searched.Size != 1 {
		t.Fatalf("宝尊查询参数不正确: %+v", searched)
	}
	if cacheTTL != 24*time.Hour {
		t.Fatalf("缓存 TTL = %v, 期望 24h", cacheTTL)
	}
	if cached.Name != "经典手袋" || cached.Upc != "upc-1" || cached.Color != "B4NQ4" {
		t.Fatalf("缓存商品信息不正确: %+v", cached)
	}
	assertCoachKOSGoodsMapValue(t, query, "enterprise_id", enterpriseID)
	assertCoachKOSGoodsMapValue(t, query, "source_id", sourceID)
	assertCoachKOSGoodsMapValue(t, query, "type", int8(3))
	assertCoachKOSGoodsMapValue(t, query, "sku_code", "CV758 B4NQ4")
	assertCoachKOSGoodsMapValue(t, query, "deleted_at", 0)
	assertCoachKOSGoodsMapValue(t, updates, "upc", "upc-1")
	assertCoachKOSGoodsMapValue(t, updates, "sku", "CV758")
	assertCoachKOSGoodsMapValue(t, updates, "name", "经典手袋")
	assertCoachKOSGoodsMapValue(t, updates, "url", "https://img.example/CV758_a0.jpg")
	assertCoachKOSGoodsMapValue(t, updates, "color", "B4NQ4")
	assertCoachKOSGoodsMapValue(t, updates, "updated_at", now)
}

func TestCoachKOSGoodsUsesCacheWithoutSearchingBaozun(t *testing.T) {
	original := coachKOSGoodsDeps
	t.Cleanup(func() { coachKOSGoodsDeps = original })

	var updates map[string]any
	coachKOSGoodsDeps = coachKOSGoodsDependencies{
		cacheGet: func(_ app.Ctx, _ string) (baozun.ProductBriefCache, bool, error) {
			payload := []byte(`{"product_no":"CV758","name":"缓存商品","color":"BLACK","upc":"cached-upc"}`)
			brief, found := coachKOSCachedProductBrief(payload)
			return brief, found, nil
		},
		cacheSet: func(_ app.Ctx, _ string, _ baozun.ProductBriefCache, _ time.Duration) error {
			t.Fatal("缓存命中不应重新写缓存")
			return nil
		},
		search: func(_ app.Ctx, _ baozun.SearchParams, _ *baozun.DataResponse) error {
			t.Fatal("缓存命中不应查询宝尊")
			return nil
		},
		updateRelation: func(_ app.Ctx, _ map[string]any, gotUpdates map[string]any) error {
			updates = gotUpdates
			return nil
		},
		now:      func() int64 { return 123 },
		logError: func(_ app.Ctx, _ string, _ ...any) {},
	}

	EnrichCoachKOSRelationGoods(nil, []CoachKOSGoodsTask{{EnterpriseID: 1, SourceID: 2, ProductCodes: []string{"CV758"}}})

	assertCoachKOSGoodsMapValue(t, updates, "sku", "CV758")
}

func TestCoachKOSGoodsLogsFailureAndContinues(t *testing.T) {
	original := coachKOSGoodsDeps
	t.Cleanup(func() { coachKOSGoodsDeps = original })

	var (
		updatedSKUs []string
		logs        []string
	)
	coachKOSGoodsDeps = coachKOSGoodsDependencies{
		cacheGet: func(_ app.Ctx, _ string) (baozun.ProductBriefCache, bool, error) {
			return baozun.ProductBriefCache{}, false, nil
		},
		cacheSet: func(_ app.Ctx, _ string, _ baozun.ProductBriefCache, _ time.Duration) error { return nil },
		search: func(_ app.Ctx, params baozun.SearchParams, response *baozun.DataResponse) error {
			if len(params.ProductCode) == 1 && params.ProductCode[0] == "FAIL" {
				return errors.New("baozun unavailable")
			}
			*response = coachKOSGoodsResponse(t, params.ProductCode[0], "成功商品", "", "BLUE", "upc-ok")
			return nil
		},
		updateRelation: func(_ app.Ctx, query, _ map[string]any) error {
			updatedSKUs = append(updatedSKUs, query["sku_code"].(string))
			return nil
		},
		now: func() int64 { return 123 },
		logError: func(_ app.Ctx, format string, args ...any) {
			logs = append(logs, fmt.Sprintf(format, args...))
		},
	}

	EnrichCoachKOSRelationGoods(nil, []CoachKOSGoodsTask{{
		EnterpriseID: 1001,
		SourceID:     2002,
		ProductCodes: []string{"FAIL", "OK"},
	}})

	if len(updatedSKUs) != 1 || updatedSKUs[0] != "OK" {
		t.Fatalf("失败后应继续更新 OK 商品，实际更新: %v", updatedSKUs)
	}
	if len(logs) != 1 {
		t.Fatalf("错误日志数 = %d, 期望 1: %v", len(logs), logs)
	}
	for _, want := range []string{"商品补全失败", "企业ID=1001", "素材ID=2002", "productCode=FAIL"} {
		if !strings.Contains(logs[0], want) {
			t.Fatalf("错误日志 %q 未包含 %q", logs[0], want)
		}
	}
}

func coachKOSGoodsResponse(t *testing.T, productNo, name, fallbackName, color, upc string) baozun.DataResponse {
	t.Helper()
	var response baozun.DataResponse
	payload := fmt.Sprintf(`{
		"data": [{
			"extData": {"resource": {"main": {"name": "CV758_a0.jpg", "url": "https://img.example/CV758_a0.jpg?token=secret"}}},
			"master": {
				"productCode": %q,
				"properties": {"productName": %q, "productNo": %q},
				"fragments": [{"properties": {"productName_en": %q, "color": %q}}]
			},
			"variants": [{"variantCode": %q, "fragments": [{"properties": {"Upc": %q}}]}]
		}]
	}`, productNo, name, strings.SplitN(productNo, "_", 2)[0], fallbackName, color, productNo, upc)
	if err := json.Unmarshal([]byte(payload), &response); err != nil {
		t.Fatalf("构造宝尊响应失败: %v", err)
	}
	return response
}

func assertCoachKOSGoodsMapValue(t *testing.T, values map[string]any, key string, want any) {
	t.Helper()
	if got := values[key]; got != want {
		t.Fatalf("%s = %#v, 期望 %#v；完整数据: %#v", key, got, want, values)
	}
}
