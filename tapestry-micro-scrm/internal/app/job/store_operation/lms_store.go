package storeOperationJob

import (
	"errors"
	"fmt"
	"path"
	"strconv"
	"strings"
	"time"

	"micro-scrm/internal/app"
	crmSftpFileModels "micro-scrm/internal/app/models/crm_sftp_file"
	thirdCrmSftpSvc "micro-scrm/internal/app/service/third_crm_sftp"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/pkg/sftp"
	"micro-scrm/internal/pkg/xjob"

	toolboxCtx "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
)

const (
	lmsDirectoryRoot  = "outbound/mstr_location"
	lmsImportTaskName = "Coach LMS门店主数据同步"
)

type SyncCoachLMSStoreMaster struct {
	getConfigs func(int) ([]micro_scrm.BrandConfigSftpDetail, error)
	connect    func(sftp.SSHConfig) (remoteFS, error)
	logs       fileLogStore
	repository masterDataRepository
	now        func() time.Time
	notify     importNotifier
}

func NewSyncCoachLMSStoreMaster() *SyncCoachLMSStoreMaster {
	return &SyncCoachLMSStoreMaster{
		getConfigs: thirdCrmSftpSvc.NewThirdCrmSftp().GetSFTPConfigList,
		connect:    connectRemoteFS,
		logs:       newCRMFileLogStore(),
		repository: newMasterDataRepository(),
		now:        time.Now,
		notify:     notifyImportResult,
	}
}

func (job *SyncCoachLMSStoreMaster) Run(c toolboxCtx.Context, handler *xjob.JobHandler) error {
	configs, err := job.getConfigs(micro_scrm.SftpTypeDEL)
	if err != nil {
		job.notify(c, "coach", lmsImportTaskName, "读取SFTP配置失败："+err.Error())
		return err
	}

	coachConfigs := make([]micro_scrm.BrandConfigSftpDetail, 0, len(configs))
	for _, config := range configs {
		if strings.EqualFold(config.BrandName, "coach") {
			coachConfigs = append(coachConfigs, config)
		}
	}
	if len(coachConfigs) == 0 {
		err = errors.New("未找到Coach LMS SFTP配置")
		job.notify(c, "coach", lmsImportTaskName, err.Error())
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

func (job *SyncCoachLMSStoreMaster) syncConfig(c app.Ctx, config micro_scrm.BrandConfigSftpDetail) error {
	remote, err := job.connect(config.ConfigDetail)
	if err != nil {
		job.notify(c, config.BrandName, lmsImportTaskName, "连接SFTP失败："+err.Error())
		return err
	}
	defer remote.Close()

	now := job.now()
	directory := lmsRemoteDirectory(now)
	files, err := remote.ReadDir(directory)
	if err != nil {
		job.notify(c, config.BrandName, lmsImportTaskName, "读取目录失败："+err.Error())
		return err
	}

	states := make(map[fileIdentity]fileState)
	identity := fileIdentity{FullPath: path.Join(directory, lmsFileName(now)), FileName: lmsFileName(now)}
	record, findErr := job.logs.Find(c, config.EnterpriseID, crmSftpFileModels.Store_Master_LMS, identity)
	if findErr == nil {
		if record.IsExtract == 2 {
			states[identity] = fileStateComplete
		} else {
			states[identity] = fileStatePending
		}
	} else if !errors.Is(findErr, gorm.ErrRecordNotFound) {
		job.notify(c, config.BrandName, lmsImportTaskName, "读取文件处理记录失败："+findErr.Error())
		return findErr
	}

	selection := selectLMSFiles(files, now, states)
	if selection.Matched == 0 {
		err = fmt.Errorf("未找到当天文件%s", lmsFileName(now))
		job.notify(c, config.BrandName, lmsImportTaskName, err.Error())
		return err
	}
	if len(selection.Pending) == 0 {
		return nil
	}
	return processLMSFile(c, config.EnterpriseID, config.BrandName, remote, directory, selection.Pending[0], job.logs, job.repository, now, job.notify)
}

func processLMSFile(c app.Ctx, enterpriseID int64, brandName string, remote remoteFS, directory string, file remoteFile, logs fileLogStore, repository masterDataRepository, now time.Time, notify importNotifier) error {
	identity := fileIdentity{FullPath: path.Join(directory, file.Name), FileName: file.Name}
	fileRecord, err := ensurePendingFile(c, logs, enterpriseID, crmSftpFileModels.Store_Master_LMS, identity)
	if err != nil {
		notify(c, brandName, lmsImportTaskName, importFailureMessage("创建文件处理记录失败："+err.Error()))
		return err
	}
	if fileRecord.IsExtract == 2 {
		return nil
	}

	reader, err := remote.Open(identity.FullPath)
	if err != nil {
		notify(c, brandName, lmsImportTaskName, importFailureMessage("打开文件失败："+err.Error()))
		return err
	}
	result, parseErr := parseLMSStores(reader)
	closeErr := reader.Close()
	if parseErr != nil {
		notify(c, brandName, lmsImportTaskName, importFailureMessage("解析文件失败："+parseErr.Error()))
		return parseErr
	}
	if closeErr != nil {
		notify(c, brandName, lmsImportTaskName, importFailureMessage("关闭文件失败："+closeErr.Error()))
		return closeErr
	}
	if warning := warningSummary(result.MissingKeyRows, result.DuplicateKeys); warning != "" {
		notify(c, brandName, lmsImportTaskName, importFailureMessage(warning))
	}
	stats, err := repository.SyncStores(c, enterpriseID, result.Records, now.Unix())
	if err != nil {
		notify(c, brandName, lmsImportTaskName, importFailureMessage("写入门店主数据失败："+err.Error()))
		return err
	}
	if err = logs.MarkComplete(c, enterpriseID, fileRecord.ID); err != nil {
		notify(c, brandName, lmsImportTaskName, importFailureMessage("更新文件处理状态失败："+err.Error()))
		return err
	}
	notify(c, brandName, lmsImportTaskName, importSuccessMessage(file.Name, stats))
	return nil
}

func lmsRemoteDirectory(now time.Time) string {
	return path.Join(lmsDirectoryRoot, now.Format("20060102")) + "/"
}

func lmsFileName(now time.Time) string {
	return "mstr_location_" + now.Format("20060102") + ".csv"
}

func selectLMSFiles(files []remoteFile, now time.Time, states map[fileIdentity]fileState) fileSelection {
	selection := fileSelection{}
	directory := lmsRemoteDirectory(now)
	targetName := lmsFileName(now)
	for _, file := range files {
		if file.IsDir || file.Name != targetName {
			continue
		}
		selection.Matched++
		identity := fileIdentity{FullPath: path.Join(directory, file.Name), FileName: file.Name}
		if states[identity] == fileStateComplete {
			selection.Completed++
			continue
		}
		selection.Pending = append(selection.Pending, file)
	}
	return selection
}
