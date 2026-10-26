package dashboardSvc

import (
	"encoding/json"
	"testing"

	appModels "micro-scrm/internal/app/models"
	dashboardModel "micro-scrm/internal/app/models/dashboard"
)

func TestRecursionAreaInfoForOpenReturnsLMSIDForEveryLevel(t *testing.T) {
	area := dashboardModel.DashboardAreaModel{
		BaseModel: appModels.BaseModel{ID: 1},
		Type:      dashboardModel.AreaSource,
		Name:      "Retail",
		LMSID:     "CHANNEL_RETAIL",
		Children: []dashboardModel.DashboardAreaModel{{
			BaseModel: appModels.BaseModel{ID: 2},
			Type:      dashboardModel.AreaRegion,
			Name:      "East",
			LMSID:     "DISTRICT_EAST",
		}},
	}

	service := &AreaService{}
	result := service.recursionAreaInfoForOpen(area, nil, nil, nil, nil)

	if result.LMSID != "CHANNEL_RETAIL" {
		t.Fatalf("root lms_id = %q, want CHANNEL_RETAIL", result.LMSID)
	}
	if len(result.Children) != 1 || result.Children[0].LMSID != "DISTRICT_EAST" {
		t.Fatalf("children = %+v, want child lms_id DISTRICT_EAST", result.Children)
	}
	encoded, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}
	var payload map[string]interface{}
	if err = json.Unmarshal(encoded, &payload); err != nil {
		t.Fatal(err)
	}
	if payload["lms_id"] != "CHANNEL_RETAIL" {
		t.Fatalf("JSON lms_id = %#v, want CHANNEL_RETAIL", payload["lms_id"])
	}
}
