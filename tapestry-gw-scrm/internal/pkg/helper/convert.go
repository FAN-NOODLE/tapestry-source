package helper

import (
	"encoding/json"
)

// StructToMap 结构体转map
func StructToMap(data interface{}) map[string]interface{} {
	strData, _ := json.Marshal(&data)
	m := make(map[string]interface{})
	_ = json.Unmarshal(strData, &m)
	return m
}
