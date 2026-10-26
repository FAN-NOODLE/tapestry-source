package helper

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

// MathRounded 四舍五入保留n位
func MathRounded(params float64, n int) float64 {
	if params == 0 {
		return params
	}
	n10 := math.Pow10(n)
	if params < 0 {
		return math.Trunc((params-0.5/n10)*n10) / n10
	}
	return math.Trunc((params+0.5/n10)*n10) / n10
}

// MathRoundedPer 处理百分比
func MathRoundedPer(params1, params2 float64) string {
	if params1 == 0 {
		return "0%"
	}
	params3 := MathRounded((params1/params2)*100, 2)
	if params3 == float64(int64(params3)) {
		return fmt.Sprintf("%.0f", params3) + "%"
	}
	return fmt.Sprintf("%.2f", params3) + "%"
}

var randGen *rand.Rand
var once sync.Once

func Rand(num int) int {
	once.Do(func() {
		randGen = rand.New(rand.NewSource(time.Now().UnixNano()))
	})

	return randGen.Intn(num)
}

// Random 获取去区间数的随机数
func Random(min, max int) int {
	once.Do(func() {
		randGen = rand.New(rand.NewSource(time.Now().UnixNano()))
	})
	num := randGen.Intn(max-min+1) + min
	return num
}
