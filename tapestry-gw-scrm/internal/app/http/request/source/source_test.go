package source

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestRequestConsoleList透传素材来源(t *testing.T) {
	payload, err := json.Marshal(RequestConsoleList{SourceOrigin: 2})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(payload), `"source_origin":2`) {
		t.Fatalf("console-list 请求缺少素材来源：%s", payload)
	}
}

func TestRequestConsoleList透传上传人和门店名称(t *testing.T) {
	payload, err := json.Marshal(RequestConsoleList{UploaderName: "张三", DepartmentName: "上海门店"})
	if err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{`"uploader_name":"张三"`, `"department_name":"上海门店"`} {
		if !strings.Contains(string(payload), field) {
			t.Fatalf("console-list 请求缺少筛选字段 %s：%s", field, payload)
		}
	}
}
