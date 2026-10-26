package inventorySvc

import (
	"errors"
	"reflect"
	"testing"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	"micro-scrm/internal/pkg/sdk/inventory"
)

type inventoryProviderStub struct {
	hht    map[string]inventory.Response
	err    map[string]error
	dcl    inventory.Response
	dclErr error
}

func (s inventoryProviderStub) post(_ ctx.Context, isDcl bool, request inventory.ParamsReq, resp *inventory.Response, _ ...clients.Option) error {
	if !isDcl {
		return errors.New("unexpected batch HHT request")
	}
	if s.dclErr != nil {
		return s.dclErr
	}
	*resp = s.dcl
	return nil
}

func (s inventoryProviderStub) postOnce(_ ctx.Context, isDcl bool, request inventory.ParamsReq, resp *inventory.Response, _ ...clients.Option) error {
	if isDcl {
		return errors.New("unexpected single DCL request")
	}
	upc := request.UpcCode[0]
	if err := s.err[upc]; err != nil {
		return err
	}
	*resp = s.hht[upc]
	return nil
}

func newInventoryServiceWithProvider(stub inventoryProviderStub) *InventoryService {
	return &InventoryService{
		postInventory:     stub.post,
		postInventoryOnce: stub.postOnce,
	}
}

func TestFetchConcurrentlyKeepsSuccessfulUPCsAndReportsFailedUPCs(t *testing.T) {
	service := newInventoryServiceWithProvider(inventoryProviderStub{
		hht: map[string]inventory.Response{
			"UPC-OK": {DetailList: []inventory.RespDetail{{
				UpcCode: "UPC-OK", LocCode: "STORE-1", OnHand: 2, InTransit: 1, SkuAttr1: "BLACK",
			}}},
		},
		err: map[string]error{"UPC-FAIL": errors.New("http 500")},
		dcl: inventory.Response{DetailList: []inventory.RespDetail{{
			UpcCode: "UPC-OK", LocCode: dclLocCode, Segment: "RTL", OnHand: 7,
		}}},
	})

	result, err := service.FetchConcurrentlyBySkuItem(ctx.New(), "sku-1", []string{"UPC-OK", "UPC-FAIL"}, "STORE-1", microWorkDao.STORE_TYPE_RETAIL)
	if err != nil {
		t.Fatalf("FetchConcurrentlyBySkuItem returned error: %v", err)
	}
	if result.InventoryStatus != InventoryStatusPartial {
		t.Fatalf("expected partial status, got %q", result.InventoryStatus)
	}
	if !reflect.DeepEqual(result.FailedUPCs, []string{"UPC-FAIL"}) {
		t.Fatalf("unexpected failed UPCs: %#v", result.FailedUPCs)
	}
	if result.DCLStatus != DCLStatusSuccess {
		t.Fatalf("expected successful DCL status, got %q", result.DCLStatus)
	}
	if len(result.List) != 1 || result.List[0].UpcCode != "UPC-OK" || result.List[0].DclOnHand != 7 {
		t.Fatalf("unexpected successful inventory rows: %#v", result.List)
	}
}

func TestFetchConcurrentlyKeepsHHTRowsWhenDCLFails(t *testing.T) {
	service := newInventoryServiceWithProvider(inventoryProviderStub{
		hht: map[string]inventory.Response{
			"UPC-1": {DetailList: []inventory.RespDetail{{
				UpcCode: "UPC-1", LocCode: "STORE-1", OnHand: 2, InTransit: 1, SkuAttr1: "BLACK",
			}}},
		},
		dclErr: errors.New("timeout"),
	})

	result, err := service.FetchConcurrentlyBySkuItem(ctx.New(), "sku-1", []string{"UPC-1"}, "STORE-1", microWorkDao.STORE_TYPE_RETAIL)
	if err != nil {
		t.Fatalf("FetchConcurrentlyBySkuItem returned error: %v", err)
	}
	if result.InventoryStatus != InventoryStatusSuccess {
		t.Fatalf("expected successful HHT status, got %q", result.InventoryStatus)
	}
	if result.DCLStatus != DCLStatusFailed {
		t.Fatalf("expected failed DCL status, got %q", result.DCLStatus)
	}
	if len(result.List) != 1 || result.List[0].OnHand != 2 {
		t.Fatalf("expected HHT row to be retained, got %#v", result.List)
	}
}

func TestFetchConcurrentlyRejectsHHTContractErrorsInsteadOfReturningPartialData(t *testing.T) {
	service := newInventoryServiceWithProvider(inventoryProviderStub{
		hht: map[string]inventory.Response{
			"UPC-OK": {DetailList: []inventory.RespDetail{{
				UpcCode: "UPC-OK", LocCode: "STORE-1", OnHand: 2, InTransit: 1,
			}}},
		},
		err: map[string]error{
			"UPC-BAD": &inventory.ContractError{Reason: "locCode 缺失"},
		},
		dcl: inventory.Response{},
	})

	result, err := service.FetchConcurrentlyBySkuItem(ctx.New(), "sku-1", []string{"UPC-OK", "UPC-BAD"}, "STORE-1", microWorkDao.STORE_TYPE_RETAIL)
	if err == nil || !inventory.IsContractError(err) {
		t.Fatalf("expected HHT contract error, got %v", err)
	}
	if len(result.List) != 0 {
		t.Fatalf("contract-invalid HHT data must not be returned, got %#v", result.List)
	}
}
