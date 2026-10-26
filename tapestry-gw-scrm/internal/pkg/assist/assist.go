package assist

import (
	"time"
)

func Timestamp() int64 {
	now := time.Now().Unix()
	return now
}
