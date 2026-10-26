package helper

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/helper"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/json"
	"math"
	"strconv"
	"strings"
)

func MarshalToString(requestData interface{}) string {
	if requestData == nil {
		return ""
	}
	return json.MarshalToString(requestData)
}

func Unmarshal(byteStr string, data interface{}) {
	if byteStr != "" {
		helper.Panic(json.Unmarshal([]byte(byteStr), &data))
	}
}

func StringToInt(param string) int {
	value, _ := strconv.Atoi(param)
	return value
}

func StringFloatToStringInt(str string) string {
	float, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return "0"
	}
	return strconv.FormatInt(int64(math.Round(float)), 10)
}

func Int64ToString(param int64) string {
	return strconv.FormatInt(param, 10)
}

// StringToIntSlice 将逗号分隔的字符串转换为 int 切片
func StringToIntSlice(str, sep string) []int {
	if str == "" {
		return []int{}
	}
	parts := strings.Split(str, sep)
	result := make([]int, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			val, err := strconv.Atoi(part)
			if err == nil {
				result = append(result, val)
			}
		}
	}
	return result
}

// StructToMap 结构体转map
func StructToMap(data interface{}) map[string]interface{} {
	strData, _ := json.Marshal(&data)
	m := make(map[string]interface{})
	_ = json.Unmarshal(strData, &m)
	return m
}

// MapStrToStruct map转结构体
func MapStrToStruct(m map[string]any, dest any) error {
	marshal, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = json.Unmarshal(marshal, &dest)
	if err != nil {
		return err
	}

	return nil
}
