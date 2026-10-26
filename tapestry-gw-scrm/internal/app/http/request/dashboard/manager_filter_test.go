package dashboardReq

import (
	"reflect"
	"testing"
)

func TestManagerAuthListRequestsExposePositionFilter(t *testing.T) {
	for _, request := range []interface{}{ManagerAuthList{}, HkManagerAuthList{}} {
		field, ok := reflect.TypeOf(request).FieldByName("Position")
		if !ok {
			t.Fatalf("%T must expose a Position filter", request)
		}
		if field.Tag.Get("json") != "position" {
			t.Fatalf("%T Position json tag = %q", request, field.Tag.Get("json"))
		}
	}
}
