package store

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	toolboxCtx "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/pkg/timingwheel"
)

type recordingRedisClient struct {
	mu     sync.Mutex
	writes [][2]string
}

func (c *recordingRedisClient) HSet(_ context.Context, key string, values ...interface{}) *redis.IntCmd {
	c.mu.Lock()
	if len(values) >= 1 {
		c.writes = append(c.writes, [2]string{key, values[0].(string)})
	}
	c.mu.Unlock()
	return redis.NewIntResult(1, nil)
}
func (*recordingRedisClient) HGet(context.Context, string, string) *redis.StringCmd {
	return redis.NewStringResult("", redis.Nil)
}
func (*recordingRedisClient) HMGet(context.Context, string, ...string) *redis.SliceCmd {
	return redis.NewSliceResult(nil, nil)
}
func (*recordingRedisClient) HDel(context.Context, string, ...string) *redis.IntCmd {
	return redis.NewIntResult(1, nil)
}

func TestStartWatchRunningRenewsOriginalCommandHashAndLockField(t *testing.T) {
	client := &recordingRedisClient{}
	store := &specRedisStore{
		client:      client,
		timingWheel: timingwheel.NewTimingWheel(time.Millisecond, 20),
		watcher:     &sync.Map{},
	}
	runningKey := "JOB_SPEC_RUNNING_KEY:SYNCCOACHKOSSOURCELINK"
	lockKey := "micro-scrm_syncCoachKOSSourceLink_501"
	store.startWatchRunning(toolboxCtx.New(), runningKey, lockKey, 20*time.Millisecond)
	t.Cleanup(func() { store.stopWatchRunning(toolboxCtx.New(), lockKey) })

	deadline := time.After(time.Second)
	for {
		client.mu.Lock()
		writes := append([][2]string(nil), client.writes...)
		client.mu.Unlock()
		if len(writes) > 0 {
			if writes[0][0] != runningKey || writes[0][1] != lockKey {
				t.Fatalf("续租位置 = %v，期望 hash=%q field=%q", writes[0], runningKey, lockKey)
			}
			return
		}
		select {
		case <-deadline:
			t.Fatal("等待续租写入超时")
		case <-time.After(5 * time.Millisecond):
		}
	}
}
