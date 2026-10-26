package inventory

import (
	"strings"
	"testing"
)

func TestParseSkuUpcCsvRowsUsesFirstTwoColumns(t *testing.T) {
	input := strings.NewReader("198685202526,ceu87,EUD,36\n198685202540,CEU87,EUD,38\n")
	rows, summary, err := parseSkuUpcCsvRows(input)
	if err != nil {
		t.Fatalf("parseSkuUpcCsvRows returned error: %v", err)
	}
	if summary.TotalRows != 2 || summary.SuccessRows != 2 || summary.FailedRows != 0 {
		t.Fatalf("unexpected summary: %#v", summary)
	}
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}
	if rows[0].SkuItem != "CEU87" || rows[0].UpcCode != "198685202526" {
		t.Fatalf("unexpected first row: %#v", rows[0])
	}
}

func TestParseSkuUpcCsvRowsDeduplicatesSkuUpc(t *testing.T) {
	input := strings.NewReader("198685202526,CEU87\n198685202526,ceu87\n")
	rows, summary, err := parseSkuUpcCsvRows(input)
	if err != nil {
		t.Fatalf("parseSkuUpcCsvRows returned error: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 deduped row, got %d", len(rows))
	}
	if summary.DuplicateRows != 1 {
		t.Fatalf("expected 1 duplicate, got %#v", summary)
	}
}

func TestParseSkuUpcCsvRowsRejectsMissingCoreFields(t *testing.T) {
	input := strings.NewReader(",CEU87\n198685202526,\n")
	rows, summary, err := parseSkuUpcCsvRows(input)
	if err != nil {
		t.Fatalf("parseSkuUpcCsvRows returned error: %v", err)
	}
	if len(rows) != 0 {
		t.Fatalf("expected no valid rows, got %#v", rows)
	}
	if summary.FailedRows != 2 {
		t.Fatalf("expected 2 failed rows, got %#v", summary)
	}
}

func TestNewInventorySkuUpcMappingModelsLeaveIDForAutoIncrement(t *testing.T) {
	rows := []skuUpcCsvRow{
		{SkuItem: "CEU87", UpcCode: "198685202526"},
	}

	mappings := newInventorySkuUpcMappingModels("batch-001", "2026-07", rows)
	if len(mappings) != 1 {
		t.Fatalf("expected 1 mapping, got %d", len(mappings))
	}
	if mappings[0].ID != 0 {
		t.Fatalf("expected zero ID for DB auto increment, got %d", mappings[0].ID)
	}
	if mappings[0].CreatedAt == 0 || mappings[0].UpdatedAt == 0 {
		t.Fatalf("expected timestamps to be set, got %#v", mappings[0])
	}
}
