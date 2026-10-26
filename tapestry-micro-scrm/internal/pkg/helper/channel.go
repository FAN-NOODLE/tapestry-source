package helper

import (
	"time"
)

// chanSize 信道容量
func chanSize(total int) int {
	l := total / 5
	if l > 100 {
		l = 100
	} else if l < 5 {
		l = 5
	}
	return l
}

// ChannelRateLimit 信道速率限制
func ChannelRateLimit(src chan any, window time.Duration, speed int) chan interface{} {
	ch := make(chan interface{}, chanSize(speed))
	go func() {
		defer close(ch)
		var count = 0         // 已处理的元素数量计数器
		barrier := time.Now() // 当前时间窗口的开始时间
		for d := range src {
			if time.Now().Sub(barrier) > window {
				// 如果当前窗口开始时间已经超过了窗口期则重置计数器和窗口开始时间
				count = 0
				barrier = time.Now()
			}
			// 处理元素，计数器++
			ch <- d
			count++
			if count >= speed {
				// 如果计数器超过了限频要求
				cc := window - time.Now().Sub(barrier)
				if cc > 0 {
					// 计算窗口期剩余时间并等待到下一个窗口期
					<-time.After(cc)
				}
				// 重置计数器和窗口开始时间
				count = 0
				barrier = time.Now()
			}
		}
	}()
	return ch
}
