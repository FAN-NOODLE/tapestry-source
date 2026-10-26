package source

import (
	"testing"

	sourceResp "micro-scrm/internal/app/http/response/source"
	sourceModel "micro-scrm/internal/app/models/source"
)

func TestBuildBrandCacheRespUsesFirstMainImage(t *testing.T) {
	products := []sourceResp.BrandProductInfo{
		{
			ProductName: "Test product",
			SkuItem:     "SKU-1",
			UpcCode:     "UPC-1",
			Color:       "Black",
			BrandSourceList: []sourceResp.BrandSourceInfo{
				{ResourceKey: "first", URL: "https://example.com/first.jpg", IsMain: sourceModel.BrandIsMainImage},
				{ResourceKey: "second", URL: "https://example.com/second.jpg", IsMain: sourceModel.BrandIsMainImage},
			},
		},
	}

	got := buildBrandCacheResp(products)

	if got.ResourceKey != "first" || got.URL != "https://example.com/first.jpg" {
		t.Fatalf("expected first main image, got resource_key=%q url=%q", got.ResourceKey, got.URL)
	}
	if got.SourceNum != 2 {
		t.Fatalf("expected source_num=2, got %d", got.SourceNum)
	}
}

func TestBuildBrandCacheRespFallsBackToFirstSource(t *testing.T) {
	products := []sourceResp.BrandProductInfo{
		{
			BrandSourceList: []sourceResp.BrandSourceInfo{
				{ResourceKey: "first", URL: "https://example.com/first.jpg"},
				{ResourceKey: "second", URL: "https://example.com/second.jpg"},
			},
		},
	}

	got := buildBrandCacheResp(products)

	if got.ResourceKey != "first" || got.URL != "https://example.com/first.jpg" {
		t.Fatalf("expected first source fallback, got resource_key=%q url=%q", got.ResourceKey, got.URL)
	}
}

func TestBuildBrandCacheRespReturnsEmptyForNoSources(t *testing.T) {
	got := buildBrandCacheResp([]sourceResp.BrandProductInfo{{}})

	if got != (sourceResp.BrandCacheResp{}) {
		t.Fatalf("expected empty response, got %#v", got)
	}
}
