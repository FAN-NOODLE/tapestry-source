package storeOperationJob

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	toolboxCtx "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	thirdCrmSftpSvc "micro-scrm/internal/app/service/third_crm_sftp"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/pkg/sdk/boss"
	"micro-scrm/internal/pkg/xjob"
)

const lmsDashboardAreaTaskName = "Coach LMS组织架构同步"

type SyncCoachLMSDashboardArea struct {
	getConfigs             func(int) ([]micro_scrm.BrandConfigSftpDetail, error)
	getStoreRootDepartment func(app.Ctx, int64) (int, error)
	repository             dashboardAreaRepository
	now                    func() time.Time
	notify                 importNotifier
}

func NewSyncCoachLMSDashboardArea() *SyncCoachLMSDashboardArea {
	return &SyncCoachLMSDashboardArea{
		getConfigs:             thirdCrmSftpSvc.NewThirdCrmSftp().GetSFTPConfigList,
		getStoreRootDepartment: loadStoreRootDepartment,
		repository:             newDashboardAreaRepository(),
		now:                    time.Now,
		notify:                 notifyImportResult,
	}
}

func (job *SyncCoachLMSDashboardArea) Run(c toolboxCtx.Context, handler *xjob.JobHandler) error {
	configs, err := job.getConfigs(micro_scrm.SftpTypeDEL)
	if err != nil {
		job.notify(c, "coach", lmsDashboardAreaTaskName, "读取Coach企业配置失败："+err.Error())
		return err
	}
	coachConfigs := make([]micro_scrm.BrandConfigSftpDetail, 0, len(configs))
	for _, config := range configs {
		if strings.EqualFold(strings.TrimSpace(config.BrandName), "coach") {
			coachConfigs = append(coachConfigs, config)
		}
	}
	if len(coachConfigs) == 0 {
		err = errors.New("未找到Coach LMS SFTP配置")
		job.notify(c, "coach", lmsDashboardAreaTaskName, err.Error())
		return err
	}

	errorsFound := make([]string, 0)
	for _, config := range coachConfigs {
		config := config
		run := func() error { return job.syncEnterprise(c, config) }
		var syncErr error
		if handler == nil {
			syncErr = run()
		} else {
			syncErr = handler.SpecDoErr(c, strconv.FormatInt(config.EnterpriseID, 10), run)
		}
		if syncErr != nil {
			errorsFound = append(errorsFound, fmt.Sprintf("enterprise_id=%d: %v", config.EnterpriseID, syncErr))
		}
	}
	if len(errorsFound) > 0 {
		return errors.New(strings.Join(errorsFound, "; "))
	}
	return nil
}

func (job *SyncCoachLMSDashboardArea) syncEnterprise(c app.Ctx, config micro_scrm.BrandConfigSftpDetail) error {
	rootThirdID, err := job.getStoreRootDepartment(c, config.EnterpriseID)
	if err != nil {
		job.notify(c, config.BrandName, lmsDashboardAreaTaskName, "读取门店根部门配置失败："+err.Error())
		return err
	}
	stats, skipped, err := job.repository.SyncDashboardAreas(c, config.EnterpriseID, rootThirdID, job.now())
	if err != nil {
		job.notify(c, config.BrandName, lmsDashboardAreaTaskName, "同步组织架构失败："+err.Error())
		return err
	}
	job.notify(c, config.BrandName, lmsDashboardAreaTaskName, dashboardAreaSuccessMessage(stats, skipped))
	return nil
}

func loadStoreRootDepartment(c app.Ctx, enterpriseID int64) (int, error) {
	response, err := boss.EnterpriseConfig(c, map[string]interface{}{
		"enterprise_id": enterpriseID,
		"key":           "store_root_dept_id",
	})
	if err != nil {
		return 0, err
	}
	if response == nil || len(response.Data) == 0 || response.Data[0].Value <= 0 {
		return 0, errors.New("未配置store_root_dept_id")
	}
	return int(response.Data[0].Value), nil
}

func dashboardAreaSuccessMessage(stats dashboardAreaSyncStats, skipped []string) string {
	message := fmt.Sprintf("成功：同步LMS组织架构，新增%d，更新%d，未变化%d，失效%d，跳过%d", stats.Inserted, stats.Updated, stats.Unchanged, stats.Invalidated, len(skipped))
	if len(skipped) == 0 {
		return message
	}
	examples := skipped
	if len(examples) > 10 {
		examples = examples[:10]
	}
	return message + "（" + strings.Join(examples, ",") + "）"
}
