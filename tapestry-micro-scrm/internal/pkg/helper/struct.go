package helper

import (
	"reflect"
)

// CoverStruct 用 b的字段的值覆盖a中重复的字段的值
func CoverStruct(a interface{}, b interface{}) {
	// 获取结构体 a 和 b 的反射值
	valA := reflect.ValueOf(a).Elem()
	valB := reflect.ValueOf(b).Elem()

	// 遍历结构体 b 的字段
	for i := 0; i < valB.NumField(); i++ {
		// 获取字段名和字段值
		fieldB := valB.Type().Field(i)
		fieldBValue := valB.Field(i)

		// 尝试在结构体 a 中找到匹配的字段
		fieldA := valA.FieldByName(fieldB.Name)

		// 如果在结构体 a 中找到了同名字段并且字段可以被设置（可导出且为可写字段）
		if fieldA.IsValid() && fieldA.CanSet() {
			// 将结构体 b 中该字段的值赋给结构体 a 中相应的字段
			if fieldA.Type() == fieldBValue.Type() {
				fieldA.Set(fieldBValue)
			} else {
				// 类型不一样时，尝试进行转换
				if (fieldA.Kind() == reflect.Int64 || fieldA.Kind() == reflect.Int32 || fieldA.Kind() == reflect.Int) && (fieldBValue.Kind() == reflect.Int64 || fieldBValue.Kind() == reflect.Int32 || fieldBValue.Kind() == reflect.Int) {
					fieldA.SetInt(fieldBValue.Int())
				}
				// 也可以根据需要扩展更多类型转换，比如 float、string 等
			}
		}
	}
}
