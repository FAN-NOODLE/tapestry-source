package cgema

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestSign(t *testing.T) {
	got := sign("secret", "app-key1717200000000abcdefghijklmnop{}")
	want := "c8f02e2108562cf26120df297011169a9143849073b0be775c1e25760c4091e4"
	if got != want {
		t.Fatalf("sign() = %s, want %s", got, want)
	}
}

func TestBuildIdempotencyKey(t *testing.T) {
	got := buildIdempotencyKey([]byte(`{"task_id":"1"}`), "1717200000000", "abcdefghijklmnop")
	want := "21e376d81446f9cf70dacf77bb05b34b"
	if got != want {
		t.Fatalf("buildIdempotencyKey() = %s, want %s", got, want)
	}
}

func TestBuildRequestHeaders(t *testing.T) {
	headers := buildRequestHeaders("app-key", "1717200000000", "nonce", "signature", "idempotency")
	want := map[string]string{
		"X-App-Key":         "REDACTED_FOR_SOURCE_SHARING",
		"X-Timestamp":       "1717200000000",
		"X-Nonce":           "nonce",
		"X-Signature":       "signature",
		"X-Idempotency-Key": "idempotency",
	}

	if len(headers) != len(want) {
		t.Fatalf("headers length = %d, want %d", len(headers), len(want))
	}
	for key, value := range want {
		if headers[key] != value {
			t.Fatalf("headers[%s] = %s, want %s", key, headers[key], value)
		}
	}
	if _, ok := headers["X-App-Secret"]; ok {
		t.Fatal("headers should not contain app secret")
	}
}

func TestNewNonce(t *testing.T) {
	nonce, err := newNonce()
	if err != nil {
		t.Fatal(err)
	}
	if len(nonce) != 32 {
		t.Fatalf("nonce length = %d, want 32", len(nonce))
	}
}

func TestIsSuccessCode(t *testing.T) {
	for _, code := range []int{0, 1006} {
		if !isSuccessCode(code) {
			t.Fatalf("code %d should be successful", code)
		}
	}
	if isSuccessCode(1007) {
		t.Fatal("code 1007 should fail")
	}
}

func TestSendResultRequestDoesNotContainLegacyFields(t *testing.T) {
	body, err := json.Marshal(SendResultRequest{
		TaskID:       "1",
		CrowdID:      "crowd",
		TaskCategory: 1,
		CrowdResult: []CrowdResultItem{{
			UnionID: "unionid",
		}},
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, field := range []string{"enterprise_id", "origin_type", "c_plan_template_id", "c_plan_template_name"} {
		if strings.Contains(string(body), field) {
			t.Fatalf("request contains legacy field %s: %s", field, body)
		}
	}
}

func TestRelationURLTrimsDcsPath(t *testing.T) {
	client := &Client{domain: "https://cgeapi-dev.coach.com.cn/dcs"}
	got := client.relationURL()
	want := "https://cgeapi-dev.coach.com.cn/expand/api/v1/projects/wecom/data-collect"
	if got != want {
		t.Fatalf("relationURL() = %s, want %s", got, want)
	}
}

func TestRelationEventRequestJSON(t *testing.T) {
	body, err := json.Marshal(RelationEventRequest{
		EventType: "bind",
		EventTime: "2026-06-17 10:00:00",
		Customer: RelationCustomer{
			UnionID:        "unionid",
			ExternalUserID: "wm_xxx",
		},
		Sa: RelationSa{
			UserID:    "sa001",
			StoreCode: "store001",
		},
		Relation: RelationInfo{
			FirstBindTime:  "2026-06-01 10:00:00",
			LatestBindTime: "2026-06-17 10:00:00",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	text := string(body)
	for _, field := range []string{"event_type", "event_time", "customer", "unionid", "external_userid", "sa", "userid", "store_code", "relation", "first_bind_time", "latest_bind_time"} {
		if !strings.Contains(text, field) {
			t.Fatalf("relation request missing field %s: %s", field, text)
		}
	}
	if strings.Contains(text, "unbind_reason") || strings.Contains(text, "unbind_time") {
		t.Fatalf("bind request contains unbind fields: %s", text)
	}
}
