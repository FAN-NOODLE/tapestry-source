package helper

import (
	"fmt"
	"github.com/samber/lo"
	"golang.org/x/exp/constraints"
	"math/rand"
	"micro-scrm/internal/app"
	"strconv"
	"strings"
	"time"
)

// RuleThree 三元运算
func RuleThree[T app.Types](ok bool, resOk, resFail T) T {
	return lo.Ternary(ok, resOk, resFail)
}

// IntSliceToString 数字类型的切片转换成逗号分割的字符串
func IntSliceToString[T app.Int](arr []T, sep string) string {
	strSlice := make([]string, len(arr))
	for i, v := range arr {
		strSlice[i] = strconv.FormatFloat(float64(v), 'f', -1, 64)
	}
	return strings.Join(strSlice, sep)
}

// RuleThreeForMap  T 表示 map 中值的类型
func RuleThreeForMap[T app.Types, K comparable, V any](m map[K]V, key K, trueVal, falseVal T) T {
	// 如果 key 存在于 map 中，返回 trueVal，否则返回 falseVal
	if _, exists := m[key]; exists {
		return trueVal
	}
	return falseVal
}

// ArraySearch 返回数组中第一个出现的坐标和值
func ArraySearch[T app.Types](needle T, haystack []T) (int, T) {
	val, key, _ := lo.FindIndexOf(haystack, func(item T) bool {
		return item == needle
	})

	return key, val
}

// InArray 判断key是否在数组中
func InArray[T app.Types](checkData T, data []T) bool {
	return lo.Contains(data, checkData)
}

// ArrayDifference 返回两个集合的差异。
// 第一个值是不在 list2 中的元素集合。
// 第二个值是 list1 中没有元素的集合。
func ArrayDifference[T app.Types](list1, list2 []T) ([]T, []T) {
	return lo.Difference(list1, list2)
}

// ArrayIntersect 返回两个集合的交叉点
func ArrayIntersect[T app.Types](list1 []T, list2 []T) []T {
	return lo.Intersect(list1, list2)
}

// ArrayMerge 合并数组
func ArrayMerge[T app.Types](arrays ...[]T) []T {
	return lo.Interleave(arrays...)
}

// ArrayRand 返回数组中随机一个值
func ArrayRand[T app.Types](data []T) T {
	index := Rand(len(data))
	return data[index]
}

func ArrayUniq[T app.Types](arr []T) []T {
	return lo.Uniq[T](arr)
}

// DelItems 从切片中删除指定的值
func DelItems[T app.Types](source []T, items []T) []T {
	if len(source) <= 0 || len(items) == 0 {
		return source
	}
	var newSlice []T
	for _, old := range source {
		if !InArray[T](old, items) {
			newSlice = append(newSlice, old)
		}
	}
	return newSlice
}

// ArrayDiff 计算数组的差集
func ArrayDiff(array []int64, othersParams ...[]int64) []int64 {
	if len(array) == 0 {
		return array
	}
	if len(array) > 0 && len(othersParams) == 0 {
		return array
	}
	var tmp = make(map[int64]int, len(array))
	for _, v := range array {
		tmp[v] = 1
	}
	for _, param := range othersParams {
		for _, arg := range param {
			if tmp[arg] != 0 {
				tmp[arg]++
			}
		}
	}
	var res = make([]int64, 0, len(tmp))
	for k, v := range tmp {
		if v == 1 {
			res = append(res, k)
		}
	}
	return res
}

// GenericArrayDiff  计算数组的差集(泛型）
func GenericArrayDiff[T app.Types](array []T, othersParams ...[]T) []T {
	if len(array) == 0 {
		return array
	}
	if len(array) > 0 && len(othersParams) == 0 {
		return array
	}
	var tmp = make(map[T]int, len(array))
	for _, v := range array {
		tmp[v] = 1
	}
	for _, param := range othersParams {
		for _, arg := range param {
			if tmp[arg] != 0 {
				tmp[arg]++
			}
		}
	}
	var res = make([]T, 0, len(tmp))
	for k, v := range tmp {
		if v == 1 {
			res = append(res, k)
		}
	}
	return res
}

// AppendIfNotExists 添加不重复的元素
func AppendIfNotExists[T app.Types](slice []T, item T) []T {
	for _, v := range slice {
		if v == item {
			return slice
		}
	}
	return append(slice, item)
}

func SliceMax[T constraints.Ordered](slice []T) T {
	if len(slice) == 0 {
		var zero T
		return zero
	}

	max := slice[0]
	for _, value := range slice[1:] {
		if value > max {
			max = value
		}
	}
	return max
}

// GetRandomValue 泛型函数，适用于任何类型切片
func GetRandomValue[T any](arr []T) T {
	if len(arr) == 0 {
		var zero T
		return zero
	}

	rand.Seed(time.Now().UnixNano())
	return arr[rand.Intn(len(arr))]
}

// SlicesEqualUnorderedUnique 比较两个切片数据是否一样
func SlicesEqualUnorderedUnique[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	// 创建map存储a中的元素
	set := make(map[T]bool, len(a))
	for _, v := range a {
		set[v] = true
	}

	// 检查b中的每个元素是否都在set中
	for _, v := range b {
		if !set[v] {
			return false
		}
	}

	// 由于长度相同且b的所有元素都在a中，所以a的所有元素也都在b中
	return true
}

func SliceToString[T app.Int](slice []T) string {
	if len(slice) == 0 {
		return ""
	}

	var result string
	for i, v := range slice {
		if i > 0 {
			result += ","
		}
		// 使用 fmt.Sprintf 将任意类型转换为字符串
		result += fmt.Sprintf("%v", v)
	}
	return result
}
