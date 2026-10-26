package sourceJob

import (
	"errors"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	toolboxCtx "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	microBossModel "micro-scrm/internal/app/models/micro_boss"
	"micro-scrm/internal/pkg/xjob"
)

func TestSyncCoachKOSSelectsOnlyCoachAndUsesSevenDayWindow(t *testing.T) {
	fixedNow := time.Date(2026, 7, 14, 10, 30, 0, 0, time.Local)
	var configIDs []int64
	var syncedEnterpriseID int64
	var syncedUpdatedAt time.Time

	job := &SyncCoachKOS{
		getEnterpriseIDs: func(*xjob.JobHandler) ([]int64, error) {
			return []int64{101, 102, 103}, nil
		},
		getEnterpriseConfig: func(_ app.Ctx, enterpriseID int64, key string) (microBossModel.EnterpriseConfigModel, error) {
			if key != "store_name" {
				t.Fatalf("配置键 = %q，期望 store_name", key)
			}
			configIDs = append(configIDs, enterpriseID)
			name := "sw"
			if enterpriseID == 102 {
				name = "coach"
			}
			return microBossModel.EnterpriseConfigModel{Name: name}, nil
		},
		syncEnterprise: func(_ app.Ctx, enterpriseID int64, updatedAt time.Time) error {
			syncedEnterpriseID = enterpriseID
			syncedUpdatedAt = updatedAt
			return nil
		},
		now:      func() time.Time { return fixedNow },
		logError: func(app.Ctx, string, ...any) {},
	}

	if err := job.Run(toolboxCtx.New(), nil); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if syncedEnterpriseID != 102 {
		t.Fatalf("同步企业 ID = %d，期望 102", syncedEnterpriseID)
	}
	if !syncedUpdatedAt.Equal(fixedNow.AddDate(0, 0, -7)) {
		t.Fatalf("updatedAt = %v，期望 %v", syncedUpdatedAt, fixedNow.AddDate(0, 0, -7))
	}
	if len(configIDs) != 2 || configIDs[0] != 101 || configIDs[1] != 102 {
		t.Fatalf("配置查询企业 = %v，期望匹配 COACH 后停止", configIDs)
	}
}

func TestSyncCoachKOSContinuesAfterConfigErrorAndAggregatesIt(t *testing.T) {
	configErr := errors.New("读取企业配置失败")
	var syncedEnterpriseID int64

	job := &SyncCoachKOS{
		getEnterpriseIDs: func(*xjob.JobHandler) ([]int64, error) {
			return []int64{201, 202, 203}, nil
		},
		getEnterpriseConfig: func(_ app.Ctx, enterpriseID int64, _ string) (microBossModel.EnterpriseConfigModel, error) {
			switch enterpriseID {
			case 201:
				return microBossModel.EnterpriseConfigModel{}, configErr
			case 202:
				return microBossModel.EnterpriseConfigModel{Name: "ks"}, nil
			default:
				return microBossModel.EnterpriseConfigModel{Name: "coach"}, nil
			}
		},
		syncEnterprise: func(_ app.Ctx, enterpriseID int64, _ time.Time) error {
			syncedEnterpriseID = enterpriseID
			return nil
		},
		now:      time.Now,
		logError: func(app.Ctx, string, ...any) {},
	}

	err := job.Run(toolboxCtx.New(), nil)
	if err == nil || !strings.Contains(err.Error(), configErr.Error()) {
		t.Fatalf("Run() error = %v，期望聚合配置错误", err)
	}
	if syncedEnterpriseID != 203 {
		t.Fatalf("同步企业 ID = %d，期望配置错误后继续找到 203", syncedEnterpriseID)
	}
}

func TestSyncCoachKOSReturnsSafelyWhenCoachIsAbsent(t *testing.T) {
	var syncCalls int32
	job := &SyncCoachKOS{
		getEnterpriseIDs: func(*xjob.JobHandler) ([]int64, error) {
			return []int64{301, 302}, nil
		},
		getEnterpriseConfig: func(_ app.Ctx, enterpriseID int64, _ string) (microBossModel.EnterpriseConfigModel, error) {
			return microBossModel.EnterpriseConfigModel{Name: "sw"}, nil
		},
		syncEnterprise: func(app.Ctx, int64, time.Time) error {
			atomic.AddInt32(&syncCalls, 1)
			return nil
		},
		now:      time.Now,
		logError: func(app.Ctx, string, ...any) {},
	}

	if err := job.Run(toolboxCtx.New(), nil); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if calls := atomic.LoadInt32(&syncCalls); calls != 0 {
		t.Fatalf("同步调用次数 = %d，期望 0", calls)
	}
}

type coachKOSSpecStore struct {
	mu    sync.Mutex
	locks []string
}

func (s *coachKOSSpecStore) Lock(_ xjob.Runner, cmd, spec string, _ time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.locks = append(s.locks, cmd+":"+spec)
	return nil
}
func (*coachKOSSpecStore) UnLock(xjob.Runner, string, string)             {}
func (*coachKOSSpecStore) HasSpecRunner(xjob.Runner, string, string) bool { return false }
func (*coachKOSSpecStore) AddSpecCmdRunner(string, xjob.Runner) error     { return nil }

func TestSyncCoachKOSSpecAndNonSpecBothUseEnterpriseLockAndReturnSyncError(t *testing.T) {
	for _, test := range []struct {
		name  string
		specs []string
	}{
		{name: "non-spec"},
		{name: "spec", specs: []string{"501"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			store := &coachKOSSpecStore{}
			options := make([]xjob.SimpleRunnerOption, 0, 1)
			if len(test.specs) > 0 {
				options = append(options, xjob.Specs(test.specs...))
			}
			runner := xjob.NewSimpleRunner("micro-scrm", "test", store, options...)
			syncErr := errors.New("同步失败")
			var syncCalls int32
			job := &SyncCoachKOS{
				getEnterpriseIDs: func(handler *xjob.JobHandler) ([]int64, error) {
					if specs := handler.Runner().Specs(); len(specs) > 0 {
						return []int64{501}, nil
					}
					return []int64{501}, nil
				},
				getEnterpriseConfig: func(app.Ctx, int64, string) (microBossModel.EnterpriseConfigModel, error) {
					return microBossModel.EnterpriseConfigModel{Name: "coach"}, nil
				},
				syncEnterprise: func(app.Ctx, int64, time.Time) error {
					atomic.AddInt32(&syncCalls, 1)
					return syncErr
				},
				now:      time.Now,
				logError: func(app.Ctx, string, ...any) {},
			}
			runner.AddJob("syncCoachKOSSourceLink", job, xjob.Timeout(time.Hour))

			if err := runner.RunCmd("syncCoachKOSSourceLink"); !errors.Is(err, syncErr) {
				t.Fatalf("RunCmd() error = %v，期望 %v", err, syncErr)
			}
			if calls := atomic.LoadInt32(&syncCalls); calls != 1 {
				t.Fatalf("主同步执行次数 = %d，期望 1", calls)
			}
			store.mu.Lock()
			locks := append([]string(nil), store.locks...)
			store.mu.Unlock()
			found := false
			for _, lock := range locks {
				if lock == "syncCoachKOSSourceLink:501" {
					found = true
				}
			}
			if !found {
				t.Fatalf("未通过企业 spec 锁执行，实际锁：%v", locks)
			}
		})
	}
}
