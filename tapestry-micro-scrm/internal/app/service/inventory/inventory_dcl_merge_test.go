package inventorySvc

import (
	"testing"

	inventoryResp "micro-scrm/internal/app/http/response/inventory"
	"micro-scrm/internal/pkg/sdk/inventory"
)

func TestAppendDclOnlySkuRows(t *testing.T) {
	result := make(map[string]inventoryResp.SkuInfo)
	dcl := map[string]inventory.RespDetail{
		"UPC-STOCK":  {UpcCode: "UPC-STOCK", SkuAttr1: "B4/MS", SkuAttr2: "ONE", OnHand: 464},
		"UPC-ZERO":   {UpcCode: "UPC-ZERO", OnHand: 0},
		"UPC-FAILED": {UpcCode: "UPC-FAILED", OnHand: 10},
	}
	hhtSuccess := map[string]struct{}{
		"UPC-STOCK": {},
		"UPC-ZERO":  {},
	}

	appendDclOnlySkuRows(result, dcl, hhtSuccess, "CDZ34")

	if len(result) != 1 {
		t.Fatalf("expected one DCL-only row, got %#v", result)
	}
	row, ok := result["UPC-STOCK"]
	if !ok || row.DclOnHand != 464 || row.Color != "B4/MS" || row.Style != "ONE" {
		t.Fatalf("unexpected DCL-only row: %#v", row)
	}
}
