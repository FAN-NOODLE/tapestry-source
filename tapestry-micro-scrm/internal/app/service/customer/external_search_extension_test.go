package customerSvc

import (
	"reflect"
	"testing"

	openReq "micro-scrm/internal/app/http/request/open"
	"micro-scrm/internal/app/models"
	externalSearchModel "micro-scrm/internal/app/models/external_search"
)

func TestBuildCouponExtensionGroupsKeepsPartialNoCouponUpdate(t *testing.T) {
	groups := buildCouponExtensionGroups(map[string]openReq.UpdateCouponParams{
		"union-1": {
			UnionId:           "union-1",
			HasCoupon:         2,
			HasWelcomeCoupon:  1,
			HasUpgradeCoupon:  0,
			HasBirthdayCoupon: 0,
			UpdatedAt:         100,
		},
	}, 200)

	mask := couponStatusMask | welcomeCouponStatusMask
	group, ok := groups[mask]
	if !ok {
		t.Fatalf("expected coupon group with mask %d", mask)
	}
	wantColumns := []string{"coupon_status", "welcome_coupon_status"}
	if !reflect.DeepEqual(group.columns, wantColumns) {
		t.Fatalf("unexpected columns: got %v want %v", group.columns, wantColumns)
	}
	if len(group.data) != 1 {
		t.Fatalf("unexpected data size: %d", len(group.data))
	}
	got := group.data[0]
	if got.CouponStatus != externalSearchModel.NoCoupon {
		t.Fatalf("expected no-coupon status, got %d", got.CouponStatus)
	}
	if got.WelcomeCouponStatus != externalSearchModel.HasCoupon {
		t.Fatalf("expected welcome coupon status, got %d", got.WelcomeCouponStatus)
	}
	if got.CouponUpdatedAt != 100 || got.CreatedAt != 200 || got.UpdatedAt != 200 {
		t.Fatalf("unexpected timestamps: %+v", got)
	}
}

func TestBuildCouponExtensionGroupsSkipsEmptyUpdates(t *testing.T) {
	groups := buildCouponExtensionGroups(map[string]openReq.UpdateCouponParams{
		"union-1": {UnionId: "union-1", UpdatedAt: 100},
		"":        {HasCoupon: 1, UpdatedAt: 100},
	}, 200)
	if len(groups) != 0 {
		t.Fatalf("expected no groups, got %d", len(groups))
	}
}

func TestRelationshipIdentityDoesNotConcatenateAmbiguously(t *testing.T) {
	if relationshipIdentity("ab", "c") == relationshipIdentity("a", "bc") {
		t.Fatal("relationship identity must separate unionid and user id")
	}
}

func TestBuildRelationshipExtensionDataKeepsEveryMatchedRelation(t *testing.T) {
	data := buildRelationshipExtensionData(
		[]openReq.UpdateActiveParams{{
			UnionId:      "union-1",
			UserId:       "user-1",
			ActiveTime:   100,
			InteractTime: 200,
		}},
		[]externalSearchModel.ExternalRelationshipModel{
			{BaseModel: models.CreateBase(11), Unionid: "union-1", UserID: "user-1"},
			{BaseModel: models.CreateBase(12), Unionid: "union-1", UserID: "user-1"},
		},
		300,
	)
	if len(data) != 2 {
		t.Fatalf("expected two relationship rows, got %d", len(data))
	}
	if data[0].RelationID != 11 || data[1].RelationID != 12 {
		t.Fatalf("unexpected relation ids: %+v", data)
	}
	for _, item := range data {
		if item.ActiveTime != 100 || item.InteractTime != 200 || item.CreatedAt != 300 || item.UpdatedAt != 300 {
			t.Fatalf("unexpected extension row: %+v", item)
		}
	}
}
