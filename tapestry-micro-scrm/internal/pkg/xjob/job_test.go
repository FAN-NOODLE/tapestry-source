package xjob

import (
	"errors"
	"strings"
	"testing"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

type specDoErrRunner struct {
	Runner
	enabled     bool
	lockErr     error
	unlockCalls int
}

func (r *specDoErrRunner) SpecEnable(string, string) bool { return r.enabled }
func (r *specDoErrRunner) Lock(ctx.Context, string, string, time.Duration) error {
	return r.lockErr
}
func (r *specDoErrRunner) UnLock(ctx.Context, string, string) { r.unlockCalls++ }

func TestSpecDoErrReturnsLockErrorWithoutExecutingBusiness(t *testing.T) {
	lockErr := errors.New("redis lock unavailable")
	runner := &specDoErrRunner{enabled: true, lockErr: lockErr}
	handler := &JobHandler{cmd: "sync", runner: runner, timeout: time.Minute}
	executed := false

	err := handler.SpecDoErr(ctx.New(), "501", func() error {
		executed = true
		return nil
	})

	if !errors.Is(err, lockErr) || !strings.Contains(err.Error(), "spec=501") {
		t.Fatalf("SpecDoErr() error = %v，期望包装锁错误并包含 spec", err)
	}
	if executed {
		t.Fatal("获取锁失败时不应执行业务函数")
	}
	if runner.unlockCalls != 0 {
		t.Fatalf("获取锁失败时解锁次数 = %d，期望 0", runner.unlockCalls)
	}
}

func TestSpecDoErrDisabledSpecIsSkippedWithoutError(t *testing.T) {
	runner := &specDoErrRunner{enabled: false, lockErr: errors.New("不应加锁")}
	handler := &JobHandler{cmd: "sync", runner: runner, timeout: time.Minute}
	executed := false

	if err := handler.SpecDoErr(ctx.New(), "501", func() error { executed = true; return nil }); err != nil {
		t.Fatalf("未分配给当前 runner 的 spec 应跳过且不报错：%v", err)
	}
	if executed {
		t.Fatal("未分配给当前 runner 的 spec 不应执行业务函数")
	}
}
