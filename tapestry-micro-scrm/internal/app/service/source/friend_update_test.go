package source

import (
	"reflect"
	"testing"

	toolboxCtx "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	sourceRequest "micro-scrm/internal/app/http/request/source"
	"micro-scrm/internal/app/models"
	sourceModel "micro-scrm/internal/app/models/source"
)

func TestFriendUpdateTextOnlyContinuesToUpdateMainSourceGroup(t *testing.T) {
	const sourceID uint64 = 101
	friendUpdated := false
	var mainUpdate sourceRequest.Request
	service := FriendService{
		getSource: func(app.Ctx, map[string]any) (sourceModel.SourceModel, error) {
			return sourceModel.SourceModel{
				MicroBaseModel: models.MicroBaseModel{
					BaseModel: models.BaseModel{ID: int64(sourceID)},
				},
			}, nil
		},
		updateFriend: func(app.Ctx, map[string]any, map[string]interface{}) error {
			friendUpdated = true
			return nil
		},
		updateSource: func(_ app.Ctx, request sourceRequest.Request) error {
			mainUpdate = request
			return nil
		},
	}

	err := service.Update(app.Ctx(toolboxCtx.New()), sourceRequest.FriendUpdateReq{
		ID: sourceID,
		FriendCreateReq: sourceRequest.FriendCreateReq{
			EnterpriseID: 1,
			StartTime:    100,
			EndTime:      200,
			SourceType:   sourceModel.TypeFriend,
			GroupId:      99,
			TextContent:  "text only",
		},
	})
	if err != nil {
		t.Fatalf("update returned error: %v", err)
	}
	if !friendUpdated {
		t.Fatal("text content was not updated")
	}
	if mainUpdate.GroupId != 99 || !reflect.DeepEqual(mainUpdate.Ids, []int64{int64(sourceID)}) {
		t.Fatalf("main source update did not carry group and source ID: %#v", mainUpdate)
	}
}
