package inventorySvc

import (
	"errors"
	"reflect"
	"testing"

	"micro-scrm/internal/app"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
)

type fakeSkuUpcBatchReader struct {
	batch microScrmModel.InventorySkuUpcImportBatchModel
	err   error
}

func (f fakeSkuUpcBatchReader) GetEnabledActiveBatch(app.Ctx, int64) (microScrmModel.InventorySkuUpcImportBatchModel, error) {
	return f.batch, f.err
}

type fakeSkuUpcMappingReader struct {
	upcs []string
	err  error
}

func (f fakeSkuUpcMappingReader) GetUpcCodesBySkuItem(app.Ctx, int64, string, string) ([]string, error) {
	return f.upcs, f.err
}

func TestSkuUpcMappingLookupUsesEnabledActiveBatch(t *testing.T) {
	resolver := skuUpcMappingResolver{
		batches: fakeSkuUpcBatchReader{batch: microScrmModel.InventorySkuUpcImportBatchModel{
			BatchNo:        "2026-07-001",
			Status:         microScrmModel.InventorySkuUpcBatchStatusSuccess,
			IsActive:       microScrmModel.InventorySkuUpcBatchActive,
			MappingEnabled: microScrmModel.InventorySkuUpcMappingEnabled,
		}},
		mappings: fakeSkuUpcMappingReader{upcs: []string{"198685202526", "198685202540"}},
	}

	upcs, enabled, err := resolver.Resolve(nil, 123, "ceu87")
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}
	if !enabled {
		t.Fatalf("expected DB mapping to be enabled")
	}
	if !reflect.DeepEqual(upcs, []string{"198685202526", "198685202540"}) {
		t.Fatalf("unexpected upcs: %#v", upcs)
	}
}

func TestSkuUpcMappingLookupDisabledWhenNoEnabledBatch(t *testing.T) {
	resolver := skuUpcMappingResolver{
		batches:  fakeSkuUpcBatchReader{},
		mappings: fakeSkuUpcMappingReader{upcs: []string{"should-not-read"}},
	}

	upcs, enabled, err := resolver.Resolve(nil, 123, "CEU87")
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}
	if enabled {
		t.Fatalf("expected DB mapping disabled")
	}
	if len(upcs) != 0 {
		t.Fatalf("expected no DB upcs when disabled, got %#v", upcs)
	}
}

func TestSkuUpcMappingLookupPropagatesMappingError(t *testing.T) {
	resolver := skuUpcMappingResolver{
		batches: fakeSkuUpcBatchReader{batch: microScrmModel.InventorySkuUpcImportBatchModel{
			BatchNo:        "2026-07-001",
			MappingEnabled: microScrmModel.InventorySkuUpcMappingEnabled,
		}},
		mappings: fakeSkuUpcMappingReader{err: errors.New("db failed")},
	}

	_, _, err := resolver.Resolve(nil, 123, "CEU87")
	if err == nil {
		t.Fatalf("expected mapping error")
	}
}

func TestResolveSkuUpcListFallsBackToPimWhenDbDisabled(t *testing.T) {
	resolver := skuUpcMappingResolver{batches: fakeSkuUpcBatchReader{}, mappings: fakeSkuUpcMappingReader{}}
	pimCalled := false
	upcs, source, dbMissFallback, err := resolveSkuUpcList(nil, resolver, 123, "CEU87", func(app.Ctx, int64, string) ([]string, error) {
		pimCalled = true
		return []string{"PIM-UPC"}, nil
	})
	if err != nil {
		t.Fatalf("resolveSkuUpcList returned error: %v", err)
	}
	if source != skuUpcSourcePIM || !pimCalled || !reflect.DeepEqual(upcs, []string{"PIM-UPC"}) {
		t.Fatalf("unexpected fallback result source=%s pimCalled=%v upcs=%#v", source, pimCalled, upcs)
	}
	if dbMissFallback {
		t.Fatalf("expected dbMissFallback false when DB mapping is disabled")
	}
}

func TestResolveSkuUpcListUsesDbWhenHit(t *testing.T) {
	resolver := skuUpcMappingResolver{
		batches: fakeSkuUpcBatchReader{batch: microScrmModel.InventorySkuUpcImportBatchModel{
			BatchNo:        "2026-07-001",
			MappingEnabled: microScrmModel.InventorySkuUpcMappingEnabled,
		}},
		mappings: fakeSkuUpcMappingReader{upcs: []string{"DB-UPC"}},
	}
	pimCalled := false
	upcs, source, dbMissFallback, err := resolveSkuUpcList(nil, resolver, 123, "CEU87", func(app.Ctx, int64, string) ([]string, error) {
		pimCalled = true
		return []string{"PIM-UPC"}, nil
	})
	if err != nil {
		t.Fatalf("resolveSkuUpcList returned error: %v", err)
	}
	if source != skuUpcSourceDB || pimCalled || !reflect.DeepEqual(upcs, []string{"DB-UPC"}) {
		t.Fatalf("unexpected DB result source=%s pimCalled=%v upcs=%#v", source, pimCalled, upcs)
	}
	if dbMissFallback {
		t.Fatalf("expected dbMissFallback false when DB mapping hits")
	}
}

func TestResolveSkuUpcListFallsBackToPimWhenDbMiss(t *testing.T) {
	resolver := skuUpcMappingResolver{
		batches: fakeSkuUpcBatchReader{batch: microScrmModel.InventorySkuUpcImportBatchModel{
			BatchNo:        "2026-07-001",
			MappingEnabled: microScrmModel.InventorySkuUpcMappingEnabled,
		}},
		mappings: fakeSkuUpcMappingReader{upcs: []string{}},
	}
	pimCalled := false
	upcs, source, dbMissFallback, err := resolveSkuUpcList(nil, resolver, 123, "CEU87", func(app.Ctx, int64, string) ([]string, error) {
		pimCalled = true
		return []string{"PIM-UPC"}, nil
	})
	if err != nil {
		t.Fatalf("resolveSkuUpcList returned error: %v", err)
	}
	if source != skuUpcSourcePIM || !pimCalled || !reflect.DeepEqual(upcs, []string{"PIM-UPC"}) {
		t.Fatalf("unexpected DB miss fallback source=%s pimCalled=%v upcs=%#v", source, pimCalled, upcs)
	}
	if !dbMissFallback {
		t.Fatalf("expected dbMissFallback true when enabled DB mapping misses")
	}
}
