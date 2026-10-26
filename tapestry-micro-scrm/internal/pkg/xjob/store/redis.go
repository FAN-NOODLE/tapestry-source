package store

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/lock"
	"log"
	"micro-scrm/internal/pkg/timingwheel"
	"micro-scrm/internal/pkg/xjob"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	jobSpecRunnerKey     = "JOB_SPEC_RUNNER_KEY"
	jobSpecRunningKey    = "JOB_SPEC_RUNNING_KEY"
	jobSpecRunnerLOCKKey = "JOB_SPEC_RUNNER_LOCK_KEY"
)

// RedisClient 使用的 redis 客户端接口
type RedisClient interface {
	HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	HGet(ctx context.Context, key, field string) *redis.StringCmd
	HMGet(ctx context.Context, key string, fields ...string) *redis.SliceCmd
	HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd
}

type specRedisStore struct {
	client      RedisClient
	lf          lock.LockerFactory
	timingWheel *timingwheel.TimingWheel
	watcher     *sync.Map
}

func (s *specRedisStore) Lock(runner xjob.Runner, cmd string, spec string, timeout time.Duration) error {
	lockKey := s.key(runner, cmd, spec)
	c := ctx.New()

	runningKey := s.getXJobKey(jobSpecRunningKey, cmd)

	mux, err := s.lf.New(c,
		lock.Key(runningKey+"_"+lockKey),
		lock.Retry(lock.LimitRetry(lock.LinearBackoff(100*time.Millisecond), 110)),
		lock.TTL(10*time.Second))
	if err != nil {
		return err
	}
	err = mux.Lock()
	if err != nil {
		return err
	}

	defer func() {
		_ = mux.Unlock()
	}()

	v, err := s.client.HGet(c, runningKey, lockKey).Result()
	if v != "" {
		val, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		if int64(val) > time.Now().Unix() {
			return errors.New(fmt.Sprintf("xjob [%s] 其他实例正在运行", lockKey))
		}
		_, err = s.client.HDel(c, runningKey, lockKey).Result()
		if err != nil {
			return err
		}
	}

	_, err = s.client.HSet(c, runningKey, lockKey, time.Now().Add(timeout).Unix()).Result()
	if err != nil {
		return err
	}
	s.startWatchRunning(c, runningKey, lockKey, timeout)
	return nil
}

func (s *specRedisStore) startWatchRunning(c ctx.Context, runningKey, lockKey string, timeout time.Duration) {
	if watcher, ok := s.watcher.Load(lockKey); ok {
		close(watcher.(chan struct{}))
	}

	watcher := make(chan struct{})
	s.watcher.Store(lockKey, watcher)
	go func() {
		for {

			select {
			case <-watcher:
				return
			case <-s.timingWheel.After(timeout / 10):
				_, err := s.client.HSet(c, runningKey, lockKey, time.Now().Add(timeout).Unix()).Result()
				log.Println(fmt.Sprintf("xjob [%s] running status:", lockKey), err)
			}
		}
	}()

}

func (s *specRedisStore) stopWatchRunning(c ctx.Context, key string) {
	if watcher, ok := s.watcher.LoadAndDelete(key); ok {
		close(watcher.(chan struct{}))
	}
}

func (s *specRedisStore) UnLock(runner xjob.Runner, cmd string, spec string) {
	lockKey := s.key(runner, cmd, spec)
	c := ctx.New()

	runningKey := s.getXJobKey(jobSpecRunningKey, cmd)
	mux, err := s.lf.New(c,
		lock.Key(runningKey+"_"+lockKey),
		lock.Retry(lock.LimitRetry(lock.LinearBackoff(100*time.Millisecond), 110)),
		lock.TTL(10*time.Second))
	if err != nil {
		return
	}
	err = mux.Lock()
	if err != nil {
		return
	}

	defer func() {
		_ = mux.Unlock()
	}()

	s.stopWatchRunning(c, lockKey)
	_, err = s.client.HDel(c, runningKey, lockKey).Result()
	if err != nil {
		return
	}
}

func (s *specRedisStore) HasSpecRunner(runner xjob.Runner, cmd string, spec string) bool {
	lockKey := s.key(runner, cmd, spec)

	runningKey := s.getXJobKey(jobSpecRunningKey, cmd)

	c := ctx.New()
	v, err := s.client.HGet(c, runningKey, lockKey).Result()
	if err == nil && v != "" {
		return true
	}
	return false
}

func (s *specRedisStore) AddSpecCmdRunner(cmd string, runner xjob.Runner) error {
	specs := runner.Specs()
	version := runner.Version()

	runningKey := s.getXJobKey(jobSpecRunnerLOCKKey, cmd)

	c := ctx.New()

	mux, err := s.lf.New(c,
		lock.Key(runningKey),
		lock.Retry(lock.LimitRetry(lock.LinearBackoff(100*time.Millisecond), 110)),
		lock.TTL(10*time.Second))
	if err != nil {
		return err
	}
	err = mux.Lock()
	if err != nil {
		return err
	}

	defer func() {
		_ = mux.Unlock()
	}()

	var specRunnerKeys []string
	for _, spec := range specs {
		specRunnerKeys = append(specRunnerKeys, s.key(runner, cmd, spec))
	}

	runnersVersion, err := s.client.HMGet(c, jobSpecRunnerKey, specRunnerKeys...).Result()
	if err != nil && err != redis.Nil {
		return err
	}

	for i, v := range runnersVersion {
		if v != nil && v.(string) != version {
			return errors.New(fmt.Sprintf("%s job %s@%s exists %s", runner.App(), cmd, specs[i], v))
		}
	}

	hMsetStr := strings.Join(specRunnerKeys, ","+version+",") + "," + version
	hMsetKV := strings.Split(hMsetStr, ",")

	err = s.client.HSet(c, jobSpecRunnerKey, hMsetKV).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *specRedisStore) key(runner xjob.Runner, cmd string, spec string) string {
	return runner.App() + "_" + cmd + "_" + spec
}

// getXJobKey 获取 XJob Key
func (s *specRedisStore) getXJobKey(key string, executor string) string {
	return fmt.Sprintf("%s:%s", key, strings.ToUpper(executor))
}

// NewSpecRedisStore 新建 spec store redis 存储
func NewSpecRedisStore(c RedisClient, factory lock.LockerFactory) xjob.SpecStore {
	store := &specRedisStore{
		client: c,
		lf:     factory,
		// TODO 这种方式没有提供更主动关闭
		timingWheel: timingwheel.NewTimingWheel(1*time.Second, 10),
		watcher:     &sync.Map{},
	}
	return store
}
