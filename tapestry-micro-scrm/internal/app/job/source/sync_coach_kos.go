package sourceJob

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"micro-scrm/internal/app"
	appJob "micro-scrm/internal/app/job"
	microBossModel "micro-scrm/internal/app/models/micro_boss"
	cacheSvc "micro-scrm/internal/app/service/cache"
	sourceSvc "micro-scrm/internal/app/service/source"
	"micro-scrm/internal/pkg/xjob"

	toolboxCtx "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

// SyncCoachKOS 同步 COACH 企业的 KOS 素材。
type SyncCoachKOS struct {
	getEnterpriseIDs    func(*xjob.JobHandler) ([]int64, error)
	getEnterpriseConfig func(app.Ctx, int64, string) (microBossModel.EnterpriseConfigModel, error)
	syncEnterprise      func(app.Ctx, int64, time.Time) error
	now                 func() time.Time
	logError            func(app.Ctx, string, ...any)
}

// NewSyncCoachKOS 创建 COACH KOS 素材同步任务。
func NewSyncCoachKOS() *SyncCoachKOS {
	service := sourceSvc.NewCoachKOSSyncService()
	return &SyncCoachKOS{
		getEnterpriseIDs:    appJob.GetEnterpriseIds,
		getEnterpriseConfig: cacheSvc.NewRedisData().GetEnterpriseConfig,
		syncEnterprise:      service.SyncEnterprise,
		now:                 time.Now,
		logError: func(c app.Ctx, format string, args ...any) {
			app.Logger().Error(c).Msgf(format, args...)
		},
	}
}

// Run 执行唯一 COACH 企业的 KOS 素材同步。
func (job *SyncCoachKOS) Run(c toolboxCtx.Context, handler *xjob.JobHandler) error {
	enterpriseIDs, err := job.getEnterpriseIDs(handler)
	if err != nil {
		return err
	}

	configErrors := make([]error, 0)
	for _, enterpriseID := range enterpriseIDs {
		config, configErr := job.getEnterpriseConfig(c, enterpriseID, "store_name")
		if configErr != nil {
			wrappedErr := fmt.Errorf("enterprise_id=%d: %w", enterpriseID, configErr)
			configErrors = append(configErrors, wrappedErr)
			job.logError(c, "读取企业 store_name 配置失败：%v", wrappedErr)
			continue
		}
		if config.Name != "coach" {
			continue
		}

		runSync := func() error {
			return job.syncEnterprise(c, enterpriseID, job.now().AddDate(0, 0, -7))
		}

		var syncErr error
		if handler == nil {
			syncErr = runSync()
		} else {
			syncErr = handler.SpecDoErr(c, strconv.FormatInt(enterpriseID, 10), runSync)
		}
		return aggregateCoachKOSErrors(configErrors, syncErr)
	}

	return aggregateCoachKOSErrors(configErrors, nil)
}

func aggregateCoachKOSErrors(configErrors []error, syncErr error) error {
	if len(configErrors) == 0 {
		return syncErr
	}

	messages := make([]string, 0, len(configErrors))
	for _, err := range configErrors {
		messages = append(messages, err.Error())
	}
	configMessage := strings.Join(messages, "; ")
	if syncErr != nil {
		return fmt.Errorf("%s; 同步失败: %w", configMessage, syncErr)
	}
	return errors.New(configMessage)
}
