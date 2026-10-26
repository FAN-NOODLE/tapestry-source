package index

import (
	"fmt"
	"math"
	"strconv"
)

// 对比
func ThreeDirection(cur int64, than int64) int64 {
	if cur > than {
		return 1
	} else if cur < than {
		return 2
	} else {
		return 0
	}
}

// 环比 百分比
func Percent(cur int64, than int64) string {
	if than == 0 {
		return ""
	}
	thanInt := math.Abs(float64(cur-than)) / float64(than) * 100

	if (cur - than) < 0 {
		return "-" + fmt.Sprintf("%.2f", thanInt) + "%"
	} else {
		return fmt.Sprintf("%.2f", thanInt) + "%"
	}
}

// 环比百分比
func ConversionPercent(cur int64, than int64) string {
	if than == 0 {
		return ""
	}
	thanInt := (float64(cur) / float64(than)) * 100
	return fmt.Sprintf("%.2f", thanInt) + "%"
}

// 排序用百分比
func ConversionPercentFloat(cur float64, than float64) float64 {
	if than == 0 {
		return 0
	}
	thanInt := (cur / than) * 100
	floatThan, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", thanInt), 64)
	return floatThan
}

// 字符串百分比
func ConversionPercentStr(cur float64, than float64) string {
	if than == 0 {
		return ""
	}
	thanInt := (cur / than) * 100
	return fmt.Sprintf("%.2f", thanInt) + "%"
}

// 单价
func UnitPrice(cur float64, than float64) float64 {
	if than == 0 {
		return 0
	}
	thanInt := cur / than
	floatThan, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", thanInt), 64)
	return floatThan
}
