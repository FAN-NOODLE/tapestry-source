package storeOperationJob

import (
	"errors"
	"fmt"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	toolboxCtx "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	crmSftpFileModels "micro-scrm/internal/app/models/crm_sftp_file"
	thirdCrmSftpSvc "micro-scrm/internal/app/service/third_crm_sftp"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/pkg/sftp"
	"micro-scrm/internal/pkg/xjob"
)

const (
	workdayRemoteDirectory = "/Tapestry/"
	workdayImportTaskName  = "Coach Workday员工主数据同步"
)

type SyncCoachWorkdayEmployeeMaster struct {
	getConfigs func(int) ([]micro_scrm.BrandConfigSftpDetail, error)
	connect    func(sftp.SSHConfig) (remoteFS, error)
	logs       fileLogStore
	repository masterDataRepository
	now        func() time.Time
	notify     importNotifier
}

func NewSyncCoachWorkdayEmployeeMaster() *SyncCoachWorkdayEmployeeMaster {
	return &SyncCoachWorkdayEmployeeMaster{
		getConfigs: thirdCrmSftpSvc.NewThirdCrmSftp().GetSFTPConfigList,
		connect:    connectRemoteFS,
		logs:       newCRMFileLogStore(),
		repository: newMasterDataRepository(),
		now:        time.Now,
		notify:     notifyImportResult,
	}
}

func (job *SyncCoachWorkdayEmployeeMaster) Run(c toolboxCtx.Context, handler *xjob.JobHandler) error {
	configs, err := job.getConfigs(micro_scrm.SftpTypeWorkday)
	if err != nil {
		job.notify(c, "coach", workdayImportTaskName, "读取SFTP配置失败："+err.Error())
		return err
	}

	coachConfigs := make([]micro_scrm.BrandConfigSftpDetail, 0, len(configs))
	for _, config := range configs {
		if strings.EqualFold(config.BrandName, "coach") {
			coachConfigs = append(coachConfigs, config)
		}
	}
	if len(coachConfigs) == 0 {
		err = errors.New("未找到Coach Workday SFTP配置")
		job.notify(c, "coach", workdayImportTaskName, err.Error())
		return err
	}

	errorsFound := make([]string, 0)
	for _, config := range coachConfigs {
		config := config
		run := func() error { return job.syncConfig(c, config) }
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

func (job *SyncCoachWorkdayEmployeeMaster) syncConfig(c app.Ctx, config micro_scrm.BrandConfigSftpDetail) error {
	remote, err := job.connect(config.ConfigDetail)
	if err != nil {
		job.notify(c, config.BrandName, workdayImportTaskName, "连接SFTP失败："+err.Error())
		return err
	}
	defer remote.Close()

	files, err := remote.ReadDir(workdayRemoteDirectory)
	if err != nil {
		job.notify(c, config.BrandName, workdayImportTaskName, "读取目录失败："+err.Error())
		return err
	}
	runNow := job.now()
	expectedFileExists := hasExpectedWorkdayFile(files, runNow)
	var missingExpectedFileErr error
	if !expectedFileExists {
		expectedDate := expectedWorkdayFileDate(runNow)
		missingExpectedFileErr = fmt.Errorf("未找到昨日Workday文件，期望文件名前缀Coach_Staff_%s", expectedDate.Format("20060102"))
		job.notify(c, config.BrandName, workdayImportTaskName, missingExpectedFileErr.Error())
	}

	states := make(map[fileIdentity]fileState)
	for _, file := range files {
		if !isWorkdayFile(file) {
			continue
		}
		identity := fileIdentity{FullPath: path.Join(workdayRemoteDirectory, file.Name), FileName: file.Name}
		record, findErr := job.logs.Find(c, config.EnterpriseID, crmSftpFileModels.Employee_Master_Workday, identity)
		if findErr == nil {
			if record.IsExtract == 2 {
				states[identity] = fileStateComplete
			} else {
				states[identity] = fileStatePending
			}
		} else if !errors.Is(findErr, gorm.ErrRecordNotFound) {
			job.notify(c, config.BrandName, workdayImportTaskName, "读取文件处理记录失败："+findErr.Error())
			return findErr
		}
	}

	selection := selectWorkdayFiles(files, states)
	if selection.Matched == 0 {
		if missingExpectedFileErr != nil {
			return missingExpectedFileErr
		}
		err = errors.New("未找到Coach_Staff_*.csv文件")
		job.notify(c, config.BrandName, workdayImportTaskName, err.Error())
		return err
	}
	if len(selection.Pending) == 0 {
		return missingExpectedFileErr
	}
	for _, file := range selection.Pending {
		if err = processWorkdayFile(c, config.EnterpriseID, config.BrandName, remote, file, job.logs, job.repository, runNow, job.notify); err != nil {
			return err
		}
	}
	return missingExpectedFileErr
}

func processWorkdayFile(c app.Ctx, enterpriseID int64, brandName string, remote remoteFS, file remoteFile, logs fileLogStore, repository masterDataRepository, now time.Time, notify importNotifier) error {
	identity := fileIdentity{FullPath: path.Join(workdayRemoteDirectory, file.Name), FileName: file.Name}
	fileRecord, err := ensurePendingFile(c, logs, enterpriseID, crmSftpFileModels.Employee_Master_Workday, identity)
	if err != nil {
		notify(c, brandName, workdayImportTaskName, importFailureMessage("创建文件处理记录失败："+err.Error()))
		return err
	}
	if fileRecord.IsExtract == 2 {
		return nil
	}

	reader, err := remote.Open(identity.FullPath)
	if err != nil {
		notify(c, brandName, workdayImportTaskName, importFailureMessage("打开文件失败："+err.Error()))
		return err
	}
	result, parseErr := parseWorkdayEmployees(reader)
	closeErr := reader.Close()
	if parseErr != nil {
		notify(c, brandName, workdayImportTaskName, importFailureMessage("解析文件失败："+parseErr.Error()))
		return parseErr
	}
	if closeErr != nil {
		notify(c, brandName, workdayImportTaskName, importFailureMessage("关闭文件失败："+closeErr.Error()))
		return closeErr
	}
	if warning := warningSummary(result.MissingKeyRows, result.DuplicateKeys); warning != "" {
		notify(c, brandName, workdayImportTaskName, importFailureMessage(warning))
	}
	stats, err := repository.SyncEmployees(c, enterpriseID, result.Records, now.Unix())
	if err != nil {
		notify(c, brandName, workdayImportTaskName, importFailureMessage("写入员工主数据失败："+err.Error()))
		return err
	}
	if err = logs.MarkComplete(c, enterpriseID, fileRecord.ID); err != nil {
		notify(c, brandName, workdayImportTaskName, importFailureMessage("更新文件处理状态失败："+err.Error()))
		return err
	}
	notify(c, brandName, workdayImportTaskName, importSuccessMessage(file.Name, stats))
	return nil
}

func isWorkdayFile(file remoteFile) bool {
	return !file.IsDir && strings.HasPrefix(file.Name, "Coach_Staff_") && strings.EqualFold(path.Ext(file.Name), ".csv")
}

func expectedWorkdayFileDate(now time.Time) time.Time {
	return now.AddDate(0, 0, -1)
}

func hasExpectedWorkdayFile(files []remoteFile, now time.Time) bool {
	prefix := "Coach_Staff_" + expectedWorkdayFileDate(now).Format("20060102")
	for _, file := range files {
		if isWorkdayFile(file) && strings.HasPrefix(file.Name, prefix) {
			return true
		}
	}
	return false
}

func selectWorkdayFiles(files []remoteFile, states map[fileIdentity]fileState) fileSelection {
	selection := fileSelection{}
	for _, file := range files {
		if !isWorkdayFile(file) {
			continue
		}
		selection.Matched++
		identity := fileIdentity{FullPath: path.Join(workdayRemoteDirectory, file.Name), FileName: file.Name}
		if states[identity] == fileStateComplete {
			selection.Completed++
			continue
		}
		selection.Pending = append(selection.Pending, file)
	}
	sort.Slice(selection.Pending, func(i, j int) bool {
		if selection.Pending[i].ModTime.Equal(selection.Pending[j].ModTime) {
			return selection.Pending[i].Name < selection.Pending[j].Name
		}
		return selection.Pending[i].ModTime.Before(selection.Pending[j].ModTime)
	})
	return selection
}
