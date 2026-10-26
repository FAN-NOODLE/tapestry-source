package externalSearchJob

import (
	externalSearchModel "micro-scrm/internal/app/models/external_search"
	"testing"
)

func TestBuildCustomerExtensionBackfillData(t *testing.T) {
	rows := []externalSearchModel.ExternalExtraModel{
		{UserUnionid: "", ValidReserveAt: 10},
		{UserUnionid: "empty"},
		{UserUnionid: "u1", ValidReserveAt: 10, TaskFinishAt: 20, CouponStatus: 1, CouponUpdatedAt: 100},
		{UserUnionid: "u1", ValidReserveAt: 30, TaskFinishAt: 15, CouponStatus: 0, CouponBy30daysStatus: 1, CouponUpdatedAt: 200},
	}
	data := buildCustomerExtensionBackfillData(rows, 300)
	if len(data) != 1 {
		t.Fatalf("len(data)=%d, want 1", len(data))
	}
	got := data[0]
	if got.Unionid != "u1" || got.ValidReserveAt != 30 || got.TaskFinishAt != 20 {
		t.Fatalf("unexpected merged timestamps: %+v", got)
	}
	if got.CouponStatus != 0 || got.CouponBy30daysStatus != 1 || got.CouponUpdatedAt != 200 {
		t.Fatalf("unexpected coupon data: %+v", got)
	}
}

func TestBuildCustomerExtensionBackfillDataClearedCoupons(t *testing.T) {
	rows := []externalSearchModel.ExternalExtraModel{
		{UserUnionid: "u2", CouponStatus: 1, CouponUpdatedAt: 100},
		{UserUnionid: "u2", CouponUpdatedAt: 200},
	}
	data := buildCustomerExtensionBackfillData(rows, 300)
	if len(data) != 1 {
		t.Fatalf("len(data)=%d, want 1", len(data))
	}
	if data[0].Unionid != "u2" || data[0].CouponStatus != 0 || data[0].CouponUpdatedAt != 200 {
		t.Fatalf("cleared coupons not retained: %+v", data[0])
	}
}

func TestBuildRelationshipExtensionBackfillData(t *testing.T) {
	rows := []externalSearchModel.ExternalRelationshipModel{
		{WorkExternalUserID: 1, InteractTime: 10},
		{},
		{Unionid: "u1", UserID: "m1", InteractTime: 10},
	}
	rows[0].ID = 1
	rows[1].ID = 2
	rows[2].ID = 1
	rows[2].ActiveTime = 20
	data := buildRelationshipExtensionBackfillData(rows, 300)
	if len(data) != 1 {
		t.Fatalf("len(data)=%d, want 1", len(data))
	}
	got := data[0]
	if got.RelationID != 1 || got.InteractTime != 10 || got.ActiveTime != 20 || got.Unionid != "u1" {
		t.Fatalf("unexpected relationship data: %+v", got)
	}
}
