package helper

import (
	"encoding/json"
	"time"
)

// WithRetry 定时重试
func WithRetry(count int, interval time.Duration, f func() error) error {
	var err error
	for i := 0; i < count; i++ {
		if err = f(); err == nil {
			return nil
		}
		if i+1 < count {
			time.Sleep(interval)
		}
	}
	return err
}

// MergeTags 合并新旧标签
func MergeTags(oldTags string, newTags []int64) string {
	if len(newTags) <= 0 {
		return oldTags
	}
	var tagIds []int64
	_ = json.Unmarshal([]byte(oldTags), &tagIds)
	tagIds = append(tagIds, newTags...)
	removeRepeatedTagIds := ArrayUniq[int64](tagIds)
	tagIdsStr, _ := json.Marshal(removeRepeatedTagIds)
	return string(tagIdsStr)
}
