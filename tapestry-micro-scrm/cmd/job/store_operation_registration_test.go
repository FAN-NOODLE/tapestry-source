package main

import (
	"testing"
	"time"

	toolboxCtx "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/pkg/xjob"
)

type registrationRunner struct {
	jobs map[string]xjob.Job
}

func (runner *registrationRunner) AddFunc(string, func(toolboxCtx.Context, *xjob.JobHandler) error, ...xjob.JobHandlerOption) {
}
func (runner *registrationRunner) AddJob(name string, job xjob.Job, _ ...xjob.JobHandlerOption) {
	if runner.jobs == nil {
		runner.jobs = make(map[string]xjob.Job)
	}
	runner.jobs[name] = job
}
func (*registrationRunner) RunCmd(string) error                                          { return nil }
func (*registrationRunner) Specs() []string                                              { return nil }
func (*registrationRunner) HasSpec() bool                                                { return false }
func (*registrationRunner) Version() string                                              { return "test" }
func (*registrationRunner) App() string                                                  { return "micro-scrm" }
func (*registrationRunner) SpecEnable(string, string) bool                               { return true }
func (*registrationRunner) Lock(toolboxCtx.Context, string, string, time.Duration) error { return nil }
func (*registrationRunner) UnLock(toolboxCtx.Context, string, string)                    {}

func TestRegisterStoreOperationJobsAddsAllCommands(t *testing.T) {
	runner := &registrationRunner{}
	registerStoreOperationJobs(runner)

	for _, name := range []string{"syncCoachLMSStoreMaster", "syncCoachWorkdayEmployeeMaster", "syncCoachLMSDashboardArea"} {
		if runner.jobs[name] == nil {
			t.Fatalf("job %q was not registered", name)
		}
	}
	if len(runner.jobs) != 3 {
		t.Fatalf("registered jobs = %v", runner.jobs)
	}
}
