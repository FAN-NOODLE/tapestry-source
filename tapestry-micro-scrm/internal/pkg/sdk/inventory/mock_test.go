package inventory

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestFindInventoryMockUsesExactHHTUPC(t *testing.T) {
	want := inventoryMock{StatusCode: 200, Body: `{"header":{"resultCode":"SUCCESS"},"detailList":[]}`, Delay: 150 * time.Millisecond}
	got, found, err := findInventoryMock("uat", false, ParamsReq{UpcCode: []string{"UPC-1"}}, func(key string) (map[string]string, error) {
		if key != "inventory:mock:uat:hht:UPC-1" {
			t.Fatalf("unexpected key: %s", key)
		}
		return map[string]string{"status_code": "200", "body": want.Body, "delay_ms": "150"}, nil
	})
	if err != nil || !found || !reflect.DeepEqual(got, want) {
		t.Fatalf("got mock=%#v found=%v err=%v", got, found, err)
	}
}

func TestFindInventoryMockUsesMatchingUPCDuringDCLBatch(t *testing.T) {
	var keys []string
	got, found, err := findInventoryMock("uat", true, ParamsReq{UpcCode: []string{"UPC-1", "UPC-2"}}, func(key string) (map[string]string, error) {
		keys = append(keys, key)
		if key == "inventory:mock:uat:dcl:UPC-2" {
			return map[string]string{"status_code": "500"}, nil
		}
		return map[string]string{}, nil
	})
	if err != nil || !found || got.StatusCode != 500 {
		t.Fatalf("got mock=%#v found=%v err=%v", got, found, err)
	}
	wantKeys := []string{"inventory:mock:uat:dcl:UPC-1", "inventory:mock:uat:dcl:UPC-2"}
	if !reflect.DeepEqual(keys, wantKeys) {
		t.Fatalf("got keys=%v want=%v", keys, wantKeys)
	}
}

func TestFindInventoryMockFallsBackOutsideUAT(t *testing.T) {
	called := false
	_, found, err := findInventoryMock("prod", false, ParamsReq{UpcCode: []string{"UPC-1"}}, func(string) (map[string]string, error) {
		called = true
		return map[string]string{"status_code": "500"}, nil
	})
	if err != nil || found || called {
		t.Fatalf("found=%v called=%v err=%v", found, called, err)
	}
}

func TestFindInventoryMockFallsBackWhenKeyDoesNotExist(t *testing.T) {
	_, found, err := findInventoryMock("uat", false, ParamsReq{UpcCode: []string{"UPC-1"}}, func(string) (map[string]string, error) {
		return map[string]string{}, nil
	})
	if err != nil || found {
		t.Fatalf("found=%v err=%v", found, err)
	}
}

func TestApplyInventoryMockRunsContractValidation(t *testing.T) {
	var resp Response
	err := applyInventoryMock(false, ParamsReq{UpcCode: []string{"UPC-1"}}, inventoryMock{
		StatusCode: 200,
		Body:       `{"header":{"resultCode":"SUCCESS"},"detailList":[{"upcCode":"UPC-1","onHand":1,"inTransit":0},{"upcCode":"UPC-1","onHand":0,"inTransit":0}]}`,
	}, &resp)
	if err == nil || !IsContractError(err) {
		t.Fatalf("expected contract error, got %v", err)
	}
}

func TestApplyInventoryMockSimulatesHTTPFailureWithoutContractAlert(t *testing.T) {
	err := applyInventoryMock(false, ParamsReq{}, inventoryMock{StatusCode: 500}, &Response{})
	if err == nil || IsContractError(err) {
		t.Fatalf("expected non-contract HTTP error, got %v", err)
	}
}

func TestFindInventoryMockReturnsRedisErrors(t *testing.T) {
	want := errors.New("redis unavailable")
	_, found, err := findInventoryMock("uat", false, ParamsReq{UpcCode: []string{"UPC-1"}}, func(string) (map[string]string, error) {
		return nil, want
	})
	if found || !errors.Is(err, want) {
		t.Fatalf("found=%v err=%v", found, err)
	}
}
