package xjob

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"time"
)

// Runner job 执行器
type Runner interface {
	// AddFunc 添加可执行函数
	AddFunc(subCmd string, f func(c ctx.Context, jobHandler *JobHandler) error, options ...JobHandlerOption)
	// AddJob  添加 Job
	AddJob(subCmd string, j Job, options ...JobHandlerOption)
	// RunCmd 运行 Job
	RunCmd(subCmd string) error
	// Specs 得到执行器特定标识，在 SpecDo 时传入的标识必须在这里面
	Specs() []string
	// HasSpec 是否配置 spec
	HasSpec() bool
	Version() string
	App() string
	SpecEnable(cmd string, spec string) bool
	Lock(c ctx.Context, cmd string, spec string, timeout time.Duration) error
	UnLock(c ctx.Context, cmd string, spec string)
}

// SpecStore runner spec store
type SpecStore interface {
	Lock(runner Runner, cmd string, spec string, timeout time.Duration) error
	UnLock(runner Runner, cmd string, spec string)
	HasSpecRunner(runner Runner, cmd string, spec string) bool
	AddSpecCmdRunner(cmd string, runner Runner) error
}

// ----------------------------------------------------------------

type simpleRunner struct {
	jobs      map[string]*JobHandler
	app       string
	version   string
	specs     []string
	specStore SpecStore
}

func (s *simpleRunner) Lock(c ctx.Context, cmd string, spec string, timeout time.Duration) error {
	return s.specStore.Lock(s, cmd, spec, timeout)
}

func (s *simpleRunner) UnLock(c ctx.Context, cmd string, spec string) {
	s.specStore.UnLock(s, cmd, spec)
}

func (s *simpleRunner) Version() string {
	return s.version
}

func (s *simpleRunner) App() string {
	return s.app
}

func (s *simpleRunner) AddFunc(subCmd string, f func(c ctx.Context, jobHandler *JobHandler) error, options ...JobHandlerOption) {
	h := &JobHandler{job: JobFunc(f), cmd: subCmd, runner: s, timeout: 1 * time.Minute}
	for _, option := range options {
		option(h)
	}
	s.jobs[subCmd] = h
}

func (s *simpleRunner) AddJob(subCmd string, j Job, options ...JobHandlerOption) {
	h := &JobHandler{job: j, cmd: subCmd, runner: s, timeout: 1 * time.Minute}
	for _, option := range options {
		option(h)
	}
	s.jobs[subCmd] = h
}

func (s *simpleRunner) RunCmd(subCmd string) error {
	if job, ok := s.jobs[subCmd]; ok {
		specs := s.Specs()

		if len(specs) == 0 {
			if err := s.specStore.Lock(s, subCmd, subCmd, 1*time.Minute); err != nil {
				panic(fmt.Sprintf("xjob cmdLock %s error: %s", subCmd, err.Error()))
			}
			defer s.specStore.UnLock(s, subCmd, subCmd)
		} else {
			if err := s.specStore.AddSpecCmdRunner(subCmd, s); err != nil {
				panic(fmt.Sprintf("xjob AddSpecCmdRunner error: %s", err.Error()))
			}
		}

		return job.gogogo(ctx.New())
	}
	return errorx.New(fmt.Sprintf("%s job not found", subCmd))
}

func (s *simpleRunner) Specs() []string {
	return s.specs
}

func (s *simpleRunner) HasSpec() bool {
	return len(s.specs) > 0
}

func (s *simpleRunner) SpecEnable(cmd string, spec string) bool {
	if s.HasSpec() {
		// 指定 spec 时，判断当前执行器是否包含 spec
		for _, sp := range s.specs {
			if sp == spec {
				return true
			}
		}
		return false
	}
	// 没有指定 spec 时，查看 spec 是否存在其他执行器
	// 如果存在不执行
	if s.specStore.HasSpecRunner(s, cmd, spec) {
		return false
	}
	return true
}

// NewSimpleRunner 新建一个简单 runner
func NewSimpleRunner(app string, version string, specStore SpecStore, options ...SimpleRunnerOption) Runner {
	runner := &simpleRunner{
		app:       app,
		version:   version,
		specStore: specStore,
		jobs:      make(map[string]*JobHandler),
	}
	for _, option := range options {
		option(runner)
	}
	return runner
}

// ----------------------------------------------------------------

// SimpleRunnerOption 配置参数
type SimpleRunnerOption func(runner *simpleRunner)

// Specs 设置 runner Specs
func Specs(specs ...string) SimpleRunnerOption {
	return func(runner *simpleRunner) {
		runner.specs = specs
	}
}
