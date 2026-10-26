package storeOperationJob

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"micro-scrm/internal/app"
	thirdCrmSftpDao "micro-scrm/internal/app/dao/third_crm_sftp"
	"micro-scrm/internal/app/models"
	crmSftpFileModels "micro-scrm/internal/app/models/crm_sftp_file"
	planService "micro-scrm/internal/app/service/plan"
	"micro-scrm/internal/pkg/sftp"

	"gorm.io/gorm"
)

type remoteFile struct {
	Name    string
	ModTime time.Time
	IsDir   bool
}

type fileIdentity struct {
	FullPath string
	FileName string
}

type fileState int

const (
	fileStateNew fileState = iota
	fileStatePending
	fileStateComplete
)

type fileSelection struct {
	Matched   int
	Completed int
	Pending   []remoteFile
}

type remoteFS interface {
	ReadDir(directory string) ([]remoteFile, error)
	Open(fullPath string) (io.ReadCloser, error)
	Close() error
}

type sftpRemoteFS struct {
	client *sftp.SFTP
}

func connectRemoteFS(config sftp.SSHConfig) (remoteFS, error) {
	client, err := sftp.NewSFTPClient(config)
	if err != nil {
		return nil, err
	}
	return &sftpRemoteFS{client: client}, nil
}

func (fs *sftpRemoteFS) ReadDir(directory string) ([]remoteFile, error) {
	entries, err := fs.client.Client.ReadDir(directory)
	if err != nil {
		return nil, err
	}
	files := make([]remoteFile, 0, len(entries))
	for _, entry := range entries {
		files = append(files, remoteFile{Name: entry.Name(), ModTime: entry.ModTime(), IsDir: entry.IsDir()})
	}
	return files, nil
}

func (fs *sftpRemoteFS) Open(fullPath string) (io.ReadCloser, error) {
	return fs.client.Client.Open(fullPath)
}

func (fs *sftpRemoteFS) Close() error {
	return fs.client.Client.Close()
}

type fileLogStore interface {
	Find(c app.Ctx, enterpriseID int64, fileType int, identity fileIdentity) (crmSftpFileModels.CrmSftpDownloadedFileModel, error)
	CreatePending(c app.Ctx, enterpriseID int64, fileType int, identity fileIdentity) (crmSftpFileModels.CrmSftpDownloadedFileModel, error)
	MarkComplete(c app.Ctx, enterpriseID, id int64) error
}

type crmFileLogStore struct {
	dao *thirdCrmSftpDao.ThirdCrmSftpFile
}

func newCRMFileLogStore() *crmFileLogStore {
	return &crmFileLogStore{dao: thirdCrmSftpDao.NewThirdCrmSftpFile()}
}

func (store *crmFileLogStore) Find(c app.Ctx, enterpriseID int64, fileType int, identity fileIdentity) (crmSftpFileModels.CrmSftpDownloadedFileModel, error) {
	return store.dao.GetFileByIdentity(c, enterpriseID, fileType, identity.FullPath, identity.FileName)
}

func (store *crmFileLogStore) CreatePending(c app.Ctx, enterpriseID int64, fileType int, identity fileIdentity) (crmSftpFileModels.CrmSftpDownloadedFileModel, error) {
	fileID := app.UID()
	record := crmSftpFileModels.CrmSftpDownloadedFileModel{
		MicroBaseModel: models.CreateMicroBase(enterpriseID, fileID),
		FileName:       identity.FileName,
		FileFullPath:   identity.FullPath,
		FileType:       int32(fileType),
		FileURL:        "",
		IsExtract:      0,
	}
	if err := store.dao.CreateData(c, record); err != nil {
		return crmSftpFileModels.CrmSftpDownloadedFileModel{}, err
	}
	return record, nil
}

func (store *crmFileLogStore) MarkComplete(c app.Ctx, enterpriseID, id int64) error {
	return store.dao.UpdateFileExTract(c, enterpriseID, id, 2)
}

type importNotifier func(c app.Ctx, enterpriseName, taskName, message string)

func notifyImportResult(c app.Ctx, enterpriseName, taskName, message string) {
	if !strings.HasPrefix(message, "成功：") && !strings.HasPrefix(message, "失败：") {
		message = importFailureMessage(message)
	}
	planService.NewTaskNotice().SftpSendNotice(c, planService.SFTPMessageData{
		EnterpriseName: enterpriseName,
		SFTPName:       taskName,
		UploadRes:      message,
	})
}

func importFailureMessage(message string) string {
	return "失败：" + message
}

func importSuccessMessage(fileName string, stats syncStats) string {
	return fmt.Sprintf("成功：读取文件%s，新增%d，更新%d，未变化%d", fileName, stats.Inserted, stats.Updated, stats.Unchanged)
}

func ensurePendingFile(c app.Ctx, logs fileLogStore, enterpriseID int64, fileType int, identity fileIdentity) (crmSftpFileModels.CrmSftpDownloadedFileModel, error) {
	record, err := logs.Find(c, enterpriseID, fileType, identity)
	if err == nil {
		return record, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return crmSftpFileModels.CrmSftpDownloadedFileModel{}, err
	}
	return logs.CreatePending(c, enterpriseID, fileType, identity)
}

func warningSummary(missingKeyRows int, duplicateKeys []string) string {
	parts := make([]string, 0, 2)
	if missingKeyRows > 0 {
		parts = append(parts, fmt.Sprintf("缺少业务键%d行", missingKeyRows))
	}
	if len(duplicateKeys) > 0 {
		keys := duplicateKeys
		if len(keys) > 10 {
			keys = keys[:10]
		}
		parts = append(parts, fmt.Sprintf("重复业务键%d个（示例：%s）", len(duplicateKeys), strings.Join(keys, ",")))
	}
	return strings.Join(parts, "；")
}
