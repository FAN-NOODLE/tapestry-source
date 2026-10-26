package xjob

import (
	"fmt"
	"log"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/fx"
)

// Job 任务接口
type Job interface {
	Run(c ctx.Context, job *JobHandler) error
}

// JobFunc Job func 类型
type JobFunc func(c ctx.Context, job *JobHandler) error

// Run 执行接口
func (jf JobFunc) Run(c ctx.Context, job *JobHandler) error {
	return jf(c, job)
}

type JobHandlerOption func(*JobHandler)

// JobHandler job 执行代理
type JobHandler struct {
	job     Job
	cmd     string
	runner  Runner
	timeout time.Duration
}

// gogogo 运行 job ，供 runner 调用
func (jh *JobHandler) gogogo(c ctx.Context) error {
	return jh.job.Run(c, jh)
}

// SpecDo 检查 spec 并运行，供业务端调用，进行业务切分
func (jh *JobHandler) SpecDo(c ctx.Context, spec string, f func()) {

	if jh.runner.SpecEnable(jh.cmd, spec) {
		// 锁定当前 spec running，然后执行
		err := jh.runner.Lock(c, jh.cmd, spec, jh.timeout)
		if err != nil {
			log.Println("[WARN] SpecDo 没有运行：", spec, err)
			return
		}
		_, _ = fx.RunSafe(func() error {
			log.Println("[INFO] SpecDo 开始运行：", spec)
			f()
			log.Println("[INFO] SpecDo 完成运行：", spec)
			return nil
		})
		jh.runner.UnLock(c, jh.cmd, spec)
	}
}

// SpecDoErr 检查 spec、加锁并执行任务，同时把业务错误返回给调用方。
// 未命中当前 runner 的 spec 时不执行任务并返回 nil。
func (jh *JobHandler) SpecDoErr(c ctx.Context, spec string, f func() error) error {
	if !jh.runner.SpecEnable(jh.cmd, spec) {
		return nil
	}
	if err := jh.runner.Lock(c, jh.cmd, spec, jh.timeout); err != nil {
		log.Println("[WARN] SpecDoErr 没有运行：", spec, err)
		return fmt.Errorf("xjob SpecDoErr 获取 spec=%s 执行锁失败: %w", spec, err)
	}
	defer jh.runner.UnLock(c, jh.cmd, spec)

	log.Println("[INFO] SpecDoErr 开始运行：", spec)
	err := f()
	log.Println("[INFO] SpecDoErr 完成运行：", spec)
	return err
}

// Runner 返回当前 runner
func (jh *JobHandler) Runner() Runner {
	return jh.runner
}

func Timeout(timeout time.Duration) JobHandlerOption {
	return func(handler *JobHandler) {
		handler.timeout = timeout
	}
}
