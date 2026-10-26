package storeOperationJob

import (
	"io"
	"strings"
	"testing"
	"time"

	toolboxCtx "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	crmSftpFileModels "micro-scrm/internal/app/models/crm_sftp_file"
)

func TestLMSPathUsesRunDate(t *testing.T) {
	now := time.Date(2026, 8, 19, 6, 0, 0, 0, time.Local)
	if got := lmsRemoteDirectory(now); got != "outbound/mstr_location/20260819/" {
		t.Fatalf("directory = %q", got)
	}
	if got := lmsFileName(now); got != "mstr_location_20260819.csv" {
		t.Fatalf("file = %q", got)
	}
}

func TestSelectLMSFilesOnlyReturnsExactPendingFile(t *testing.T) {
	now := time.Date(2026, 8, 19, 6, 0, 0, 0, time.Local)
	files := []remoteFile{
		{Name: "mstr_location_20260818.csv", ModTime: time.Unix(10, 0)},
		{Name: "mstr_location_20260819.csv", ModTime: time.Unix(20, 0)},
		{Name: "mstr_location_20260819.csv.bak", ModTime: time.Unix(30, 0)},
	}
	identity := fileIdentity{
		FullPath: "outbound/mstr_location/20260819/mstr_location_20260819.csv",
		FileName: "mstr_location_20260819.csv",
	}

	selection := selectLMSFiles(files, now, map[fileIdentity]fileState{identity: fileStatePending})
	if selection.Matched != 1 || len(selection.Pending) != 1 || selection.Pending[0].Name != "mstr_location_20260819.csv" {
		t.Fatalf("selection = %+v", selection)
	}

	selection = selectLMSFiles(files, now, map[fileIdentity]fileState{identity: fileStateComplete})
	if selection.Matched != 1 || len(selection.Pending) != 0 || selection.Completed != 1 {
		t.Fatalf("completed selection = %+v", selection)
	}
}

type fakeRemoteFS struct {
	files   map[string]string
	list    []remoteFile
	readErr error
	openErr error
	closed  bool
}

func (fs *fakeRemoteFS) ReadDir(string) ([]remoteFile, error) { return fs.list, fs.readErr }
func (fs *fakeRemoteFS) Open(name string) (io.ReadCloser, error) {
	if fs.openErr != nil {
		return nil, fs.openErr
	}
	return io.NopCloser(strings.NewReader(fs.files[name])), nil
}
func (fs *fakeRemoteFS) Close() error { fs.closed = true; return nil }

type fakeFileLogStore struct {
	record      crmSftpFileModels.CrmSftpDownloadedFileModel
	findErr     error
	createErr   error
	completeErr error
	created     bool
	completedID int64
}

func (store *fakeFileLogStore) Find(app.Ctx, int64, int, fileIdentity) (crmSftpFileModels.CrmSftpDownloadedFileModel, error) {
	return store.record, store.findErr
}
func (store *fakeFileLogStore) CreatePending(_ app.Ctx, _ int64, _ int, identity fileIdentity) (crmSftpFileModels.CrmSftpDownloadedFileModel, error) {
	store.created = true
	if store.createErr != nil {
		return crmSftpFileModels.CrmSftpDownloadedFileModel{}, store.createErr
	}
	store.record.ID = 91
	store.record.FileName = identity.FileName
	store.record.FileFullPath = identity.FullPath
	return store.record, nil
}
func (store *fakeFileLogStore) MarkComplete(_ app.Ctx, _ int64, id int64) error {
	store.completedID = id
	return store.completeErr
}

type fakeMasterDataRepository struct {
	stores        []storeSourceRecord
	storeStats    syncStats
	storeErr      error
	employees     []employeeSourceRecord
	employeeStats syncStats
	employeeErr   error
}

func (repository *fakeMasterDataRepository) SyncStores(_ app.Ctx, _ int64, records []storeSourceRecord, _ int64) (syncStats, error) {
	repository.stores = append([]storeSourceRecord(nil), records...)
	return repository.storeStats, repository.storeErr
}
func (repository *fakeMasterDataRepository) SyncEmployees(_ app.Ctx, _ int64, records []employeeSourceRecord, _ int64) (syncStats, error) {
	repository.employees = append([]employeeSourceRecord(nil), records...)
	return repository.employeeStats, repository.employeeErr
}

func TestProcessLMSFileMarksCompleteOnlyAfterStoreSync(t *testing.T) {
	directory := "outbound/mstr_location/20260819/"
	file := remoteFile{Name: "mstr_location_20260819.csv"}
	fullPath := directory + file.Name
	remote := &fakeRemoteFS{files: map[string]string{fullPath: "retail_pos_number,status\nOC267,OPEN\n"}}
	logs := &fakeFileLogStore{findErr: gorm.ErrRecordNotFound}
	repository := &fakeMasterDataRepository{storeStats: syncStats{Inserted: 1, Updated: 2, Unchanged: 3}}
	var notifications []string

	err := processLMSFile(toolboxCtx.New(), 3360610327843968, "coach", remote, directory, file, logs, repository, time.Unix(1787100000, 0), func(_ app.Ctx, _, _, message string) {
		notifications = append(notifications, message)
	})
	if err != nil {
		t.Fatal(err)
	}
	if !logs.created || logs.completedID != 91 {
		t.Fatalf("created=%v completedID=%d", logs.created, logs.completedID)
	}
	if len(repository.stores) != 1 || repository.stores[0].StoreCode != "OC267" {
		t.Fatalf("stores=%+v", repository.stores)
	}
	if len(notifications) != 1 || notifications[0] != "成功：读取文件mstr_location_20260819.csv，新增1，更新2，未变化3" {
		t.Fatalf("notifications=%v", notifications)
	}
}

func TestProcessLMSFileLeavesPendingWhenCSVIsInvalid(t *testing.T) {
	directory := "outbound/mstr_location/20260819/"
	file := remoteFile{Name: "mstr_location_20260819.csv"}
	fullPath := directory + file.Name
	remote := &fakeRemoteFS{files: map[string]string{fullPath: "retail_pos_number,status\n\"broken\n"}}
	logs := &fakeFileLogStore{findErr: gorm.ErrRecordNotFound}
	repository := &fakeMasterDataRepository{}

	var notifications []string
	err := processLMSFile(toolboxCtx.New(), 3360610327843968, "coach", remote, directory, file, logs, repository, time.Unix(1787100000, 0), func(_ app.Ctx, _, _, message string) {
		notifications = append(notifications, message)
	})
	if err == nil {
		t.Fatal("invalid CSV must fail")
	}
	if !logs.created || logs.completedID != 0 || len(repository.stores) != 0 {
		t.Fatalf("created=%v completedID=%d stores=%+v", logs.created, logs.completedID, repository.stores)
	}
	if len(notifications) != 1 || !strings.HasPrefix(notifications[0], "失败：解析文件失败：") {
		t.Fatalf("notifications=%v", notifications)
	}
}
