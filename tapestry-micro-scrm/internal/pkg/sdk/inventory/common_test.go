package inventory

import (
	"errors"
	"testing"
)

func TestDecodeAndValidateResponseAcceptsEmptyDetailList(t *testing.T) {
	var resp Response
	err := decodeAndValidateResponse(false, ParamsReq{UpcCode: []string{"UPC-1"}}, `{
		"header":{"resultCode":"SUCCESS","resultMessage":"ok"},
		"detailList":[]
	}`, &resp)
	if err != nil {
		t.Fatalf("expected legal empty detailList, got %v", err)
	}
	if resp.Header == nil || resp.Header.ResultCode != "SUCCESS" {
		t.Fatalf("unexpected header: %#v", resp.Header)
	}
	if len(resp.DetailList) != 0 {
		t.Fatalf("expected empty detailList, got %#v", resp.DetailList)
	}
}

func TestDecodeAndValidateResponseRejectsMalformedJSONAsContractError(t *testing.T) {
	var resp Response
	err := decodeAndValidateResponse(false, ParamsReq{}, `{"header":`, &resp)
	if err == nil || !IsContractError(err) {
		t.Fatalf("expected contract error, got %v", err)
	}
}

func TestDecodeAndValidateResponseRejectsMissingResponseStructure(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{name: "missing header", body: `{"detailList":[]}`},
		{name: "missing result code", body: `{"header":{},"detailList":[]}`},
		{name: "missing detail list", body: `{"header":{"resultCode":"SUCCESS"}}`},
		{name: "null detail list", body: `{"header":{"resultCode":"SUCCESS"},"detailList":null}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp Response
			err := decodeAndValidateResponse(false, ParamsReq{}, tt.body, &resp)
			if err == nil || !IsContractError(err) {
				t.Fatalf("expected contract error, got %v", err)
			}
		})
	}
}

func TestDecodeAndValidateResponseRejectsMissingOrMistypedHHTFields(t *testing.T) {
	tests := []struct {
		name   string
		detail string
	}{
		{name: "missing upcCode", detail: `{"locCode":"S001","onHand":0,"inTransit":0}`},
		{name: "null onHand", detail: `{"upcCode":"UPC-1","locCode":"S001","onHand":null,"inTransit":0}`},
		{name: "string inTransit", detail: `{"upcCode":"UPC-1","locCode":"S001","onHand":0,"inTransit":"0"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp Response
			body := `{"header":{"resultCode":"SUCCESS"},"detailList":[` + tt.detail + `]}`
			err := decodeAndValidateResponse(false, ParamsReq{UpcCode: []string{"UPC-1"}}, body, &resp)
			if err == nil || !IsContractError(err) {
				t.Fatalf("expected contract error, got %v", err)
			}
		})
	}
}

func TestDecodeAndValidateResponseAcceptsSingleHHTDetailWithoutLocation(t *testing.T) {
	tests := []struct {
		name   string
		detail string
	}{
		{name: "missing locCode", detail: `{"upcCode":"UPC-1","onHand":-1,"inTransit":0}`},
		{name: "empty locCode", detail: `{"upcCode":"UPC-1","locCode":"","onHand":0,"inTransit":0}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp Response
			body := `{"header":{"resultCode":"SUCCESS"},"detailList":[` + tt.detail + `]}`
			err := decodeAndValidateResponse(false, ParamsReq{UpcCode: []string{"UPC-1"}}, body, &resp)
			if err != nil {
				t.Fatalf("single HHT detail without usable locCode should be ignored, got %v", err)
			}
			if len(resp.DetailList) != 0 {
				t.Fatalf("expected unusable detail to be filtered, got %#v", resp.DetailList)
			}
		})
	}
}

func TestDecodeAndValidateResponseFiltersHHTDetailWithoutLocation(t *testing.T) {
	var resp Response
	body := `{
		"header":{"resultCode":"SUCCESS"},
		"detailList":[
			{"upcCode":"UPC-1","locCode":"S001","onHand":2,"inTransit":0},
			{"upcCode":"UPC-1","onHand":-1,"inTransit":0}
		]
	}`
	err := decodeAndValidateResponse(false, ParamsReq{UpcCode: []string{"UPC-1"}}, body, &resp)
	if err != nil {
		t.Fatalf("HHT detail without locCode should be ignored when usable rows exist, got %v", err)
	}
	if len(resp.DetailList) != 1 || resp.DetailList[0].LocCode != "S001" {
		t.Fatalf("unexpected filtered detailList: %#v", resp.DetailList)
	}
}

func TestDecodeAndValidateResponseRejectsHHTWhenAllLocationsAreMissing(t *testing.T) {
	var resp Response
	body := `{
		"header":{"resultCode":"SUCCESS"},
		"detailList":[
			{"upcCode":"UPC-1","onHand":-1,"inTransit":0},
			{"upcCode":"UPC-1","locCode":"","onHand":0,"inTransit":0}
		]
	}`
	err := decodeAndValidateResponse(false, ParamsReq{UpcCode: []string{"UPC-1"}}, body, &resp)
	if err == nil || !IsContractError(err) {
		t.Fatalf("expected contract error, got %v", err)
	}
}

func TestDecodeAndValidateResponseAcceptsUnusedDCLSegment(t *testing.T) {
	var resp Response
	body := `{"header":{"resultCode":"SUCCESS"},"detailList":[{"upcCode":"UPC-1","locCode":"CN10","segment":"ECM","onHand":2}]}`
	err := decodeAndValidateResponse(true, ParamsReq{
		UpcCode:         []string{"UPC-1"},
		ExpectedSegment: "RTL",
	}, body, &resp)
	if err != nil {
		t.Fatalf("unused DCL segment must be accepted and filtered by the business layer, got %v", err)
	}
	if len(resp.DetailList) != 1 || resp.DetailList[0].Segment != "ECM" {
		t.Fatalf("unexpected decoded detailList: %#v", resp.DetailList)
	}
}

func TestDecodeAndValidateResponseAcceptsUnusedDCLLocation(t *testing.T) {
	var resp Response
	body := `{"header":{"resultCode":"SUCCESS"},"detailList":[{"upcCode":"UPC-1","locCode":"CN01","segment":"RTL","onHand":9},{"upcCode":"UPC-1","locCode":"CN10","segment":"RTL","onHand":2}]}`
	err := decodeAndValidateResponse(true, ParamsReq{UpcCode: []string{"UPC-1"}}, body, &resp)
	if err != nil {
		t.Fatalf("unused DCL locations must be accepted and filtered by the business layer, got %v", err)
	}
	if len(resp.DetailList) != 2 || resp.DetailList[0].LocCode != "CN01" || resp.DetailList[1].LocCode != "CN10" {
		t.Fatalf("unexpected decoded detailList: %#v", resp.DetailList)
	}
}

func TestDecodeAndValidateResponseRejectsDCLContractViolations(t *testing.T) {
	tests := []struct {
		name   string
		detail string
	}{
		{name: "missing segment", detail: `{"upcCode":"UPC-1","locCode":"CN10","onHand":2}`},
		{name: "missing locCode", detail: `{"upcCode":"UPC-1","segment":"RTL","onHand":2}`},
		{name: "unknown upc", detail: `{"upcCode":"UPC-X","locCode":"CN10","segment":"RTL","onHand":2}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp Response
			body := `{"header":{"resultCode":"SUCCESS"},"detailList":[` + tt.detail + `]}`
			err := decodeAndValidateResponse(true, ParamsReq{UpcCode: []string{"UPC-1"}}, body, &resp)
			if err == nil || !IsContractError(err) {
				t.Fatalf("expected contract error, got %v", err)
			}
		})
	}
}

func TestDecodeAndValidateResponseRejectsUnexpectedRequestedLocation(t *testing.T) {
	var resp Response
	body := `{"header":{"resultCode":"SUCCESS"},"detailList":[{"upcCode":"UPC-1","locCode":"S002","onHand":1,"inTransit":0}]}`
	err := decodeAndValidateResponse(false, ParamsReq{UpcCode: []string{"UPC-1"}, LocCode: []string{"S001"}}, body, &resp)
	if err == nil || !IsContractError(err) {
		t.Fatalf("expected contract error, got %v", err)
	}
}

func TestIsContractErrorDoesNotClassifyTransportErrors(t *testing.T) {
	if IsContractError(errors.New("http 500")) {
		t.Fatal("transport errors must not trigger contract alerts")
	}
}
