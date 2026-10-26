package source

import (
	"reflect"
	"testing"

	"micro-scrm/internal/app/models"
)

func TestLinkKOSModelContract(t *testing.T) {
	model := LinkKOSModel{}
	if got := model.TableName(); got != "source_link_kos" {
		t.Fatalf("表名错误：%q", got)
	}

	typ := reflect.TypeOf(model)
	base, ok := typ.FieldByName("MicroBaseModel")
	if !ok || base.Type != reflect.TypeOf(models.MicroBaseModel{}) {
		t.Fatal("LinkKOSModel 未嵌入 models.MicroBaseModel")
	}

	wantTypes := map[string]reflect.Kind{
		"SourceID": reflect.Int64, "SourceLinkID": reflect.Int64,
		"PostID": reflect.String, "SourceURL": reflect.String, "DataDate": reflect.String,
		"PublishedAt": reflect.String, "APIUpdatedAt": reflect.String,
		"SourceUpdatedAt": reflect.String,
	}
	for name, wantKind := range wantTypes {
		field, ok := typ.FieldByName(name)
		if !ok {
			t.Fatalf("LinkKOSModel 缺少 %s 字段", name)
		}
		if field.Type.Kind() != wantKind {
			t.Fatalf("%s 类型错误：%v", name, field.Type)
		}
	}
}
