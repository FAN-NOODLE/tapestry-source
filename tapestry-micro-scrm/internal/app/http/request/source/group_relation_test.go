package source

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestRequestConsoleList包含素材来源(t *testing.T) {
	payload, err := json.Marshal(RequestConsoleList{SourceOrigin: SourceOriginKOSSync})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(payload), `"source_origin":2`) {
		t.Fatalf("console-list 请求缺少素材来源：%s", payload)
	}
}
