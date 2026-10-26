package storeOperationJob

import (
	"errors"
	"testing"
	"time"

	toolboxCtx "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	crmSftpFileModels "micro-scrm/internal/app/models/crm_sftp_file"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/pkg/sftp"
)

func TestSelectWorkdayFilesFiltersCSVAndOrdersPendingOldestFirst(t *testing.T) {
	files := []remoteFile{
		{Name: "Coach_Staff_20260819_2.csv", ModTime: time.Unix(20, 0)},
		{Name: "notes.txt", ModTime: time.Unix(1, 0)},
		{Name: "Coach_Staff_20260818_1.CSV", ModTime: time.Unix(10, 0)},
		{Name: "Other_Staff_20260819.csv", ModTime: time.Unix(5, 0)},
	}
	completedIdentity := fileIdentity{FullPath: "/Tapestry/Coach_Staff_20260819_2.csv", FileName: "Coach_Staff_20260819_2.csv"}

	selection := selectWorkdayFiles(files, map[fileIdentity]fileState{completedIdentity: fileStateComplete})
	if selection.Matched != 2 || selection.Completed != 1 || len(selection.Pending) != 1 {
		t.Fatalf("selection = %+v", selection)
	}
	if selection.Pending[0].Name != "Coach_Staff_20260818_1.CSV" {
		t.Fatalf("pending = %+v", selection.Pending)
	}

	selection = selectWorkdayFiles(files, nil)
	if len(selection.Pending) != 2 || selection.Pending[0].Name != "Coach_Staff_20260818_1.CSV" || selection.Pending[1].Name != "Coach_Staff_20260819_2.csv" {
		t.Fatalf("ordered pending = %+v", selection.Pending)
	}
}

func TestHasExpectedWorkdayFileUsesPreviousDateInFileName(t *testing.T) {
	now := time.Date(2026, 8, 25, 6, 0, 0, 0, time.Local)
	if hasExpectedWorkdayFile([]remoteFile{{Name: "Coach_Staff_20260825_12345.csv"}}, now) {
		t.Fatal("当天日期文件不应满足昨日交付检查")
	}
	if !hasExpectedWorkdayFile([]remoteFile{{Name: "Coach_Staff_20260824_12345.csv"}}, now) {
		t.Fatal("未识别昨日日期的Workday文件")
	}
}

func TestExpectedWorkdayFileDateHandlesMonthAndYearBoundaries(t *testing.T) {
	tests := []struct {
		now      time.Time
		expected string
	}{
		{now: time.Date(2026, 9, 1, 6, 0, 0, 0, time.Local), expected: "20260831"},
		{now: time.Date(2027, 1, 1, 6, 0, 0, 0, time.Local), expected: "20261231"},
	}
	for _, test := range tests {
		if actual := expectedWorkdayFileDate(test.now).Format("20060102"); actual != test.expected {
			t.Fatalf("now=%s，期望文件日期=%s，实际=%s", test.now.Format("20060102"), test.expected, actual)
		}
	}
}

type workdayFileLogStore struct {
	records        map[fileIdentity]crmSftpFileModels.CrmSftpDownloadedFileModel
	identitiesByID map[int64]fileIdentity
	completed      []fileIdentity
	nextID         int64
}

func newWorkdayFileLogStore() *workdayFileLogStore {
	return &workdayFileLogStore{
		records:        make(map[fileIdentity]crmSftpFileModels.CrmSftpDownloadedFileModel),
		identitiesByID: make(map[int64]fileIdentity),
		nextID:         100,
	}
}

func (store *workdayFileLogStore) Find(_ app.Ctx, _ int64, _ int, identity fileIdentity) (crmSftpFileModels.CrmSftpDownloadedFileModel, error) {
	record, ok := store.records[identity]
	if !ok {
		return crmSftpFileModels.CrmSftpDownloadedFileModel{}, gorm.ErrRecordNotFound
	}
	return record, nil
}

func (store *workdayFileLogStore) CreatePending(_ app.Ctx, _ int64, _ int, identity fileIdentity) (crmSftpFileModels.CrmSftpDownloadedFileModel, error) {
	store.nextID++
	record := crmSftpFileModels.CrmSftpDownloadedFileModel{FileName: identity.FileName, FileFullPath: identity.FullPath}
	record.ID = store.nextID
	store.records[identity] = record
	store.identitiesByID[record.ID] = identity
	return record, nil
}

func (store *workdayFileLogStore) MarkComplete(_ app.Ctx, _ int64, id int64) error {
	identity := store.identitiesByID[id]
	record := store.records[identity]
	record.IsExtract = 2
	store.records[identity] = record
	store.completed = append(store.completed, identity)
	return nil
}

type orderedEmployeeRepository struct {
	codes    []string
	failCode string
}

func (*orderedEmployeeRepository) SyncStores(app.Ctx, int64, []storeSourceRecord, int64) (syncStats, error) {
	return syncStats{}, nil
}

func (repository *orderedEmployeeRepository) SyncEmployees(_ app.Ctx, _ int64, records []employeeSourceRecord, _ int64) (syncStats, error) {
	if len(records) != 1 {
		return syncStats{}, errors.New("expected one employee")
	}
	code := records[0].EmployeeCode
	repository.codes = append(repository.codes, code)
	if code == repository.failCode {
		return syncStats{}, errors.New("employee sync failed")
	}
	return syncStats{Inserted: 1}, nil
}

func TestWorkdaySyncStopsAfterOlderFileFails(t *testing.T) {
	older := remoteFile{Name: "Coach_Staff_20260818_1.csv", ModTime: time.Unix(10, 0)}
	newer := remoteFile{Name: "Coach_Staff_20260819_1.csv", ModTime: time.Unix(20, 0)}
	remote := &fakeRemoteFS{
		list: []remoteFile{newer, older},
		files: map[string]string{
			"/Tapestry/Coach_Staff_20260818_1.csv": "EmployeeId,Name\n000111,Older\n",
			"/Tapestry/Coach_Staff_20260819_1.csv": "EmployeeId,Name\n000222,Newer\n",
		},
	}
	logs := newWorkdayFileLogStore()
	repository := &orderedEmployeeRepository{failCode: "111"}
	job := &SyncCoachWorkdayEmployeeMaster{
		connect:    func(sftp.SSHConfig) (remoteFS, error) { return remote, nil },
		logs:       logs,
		repository: repository,
		now:        func() time.Time { return time.Unix(1787100000, 0) },
		notify:     func(app.Ctx, string, string, string) {},
	}
	config := micro_scrm.BrandConfigSftpDetail{}
	config.EnterpriseID = 3360610327843968
	config.BrandName = "coach"

	err := job.syncConfig(toolboxCtx.New(), config)
	if err == nil {
		t.Fatal("older file failure must be returned")
	}
	if len(repository.codes) != 1 || repository.codes[0] != "111" {
		t.Fatalf("processed codes = %v", repository.codes)
	}
	if len(logs.completed) != 0 {
		t.Fatalf("completed files = %+v", logs.completed)
	}
}

func TestWorkdaySyncProcessesPendingFilesOldestFirst(t *testing.T) {
	older := remoteFile{Name: "Coach_Staff_20260818_1.csv", ModTime: time.Unix(10, 0)}
	newer := remoteFile{Name: "Coach_Staff_20260819_1.csv", ModTime: time.Unix(20, 0)}
	remote := &fakeRemoteFS{
		list: []remoteFile{newer, older},
		files: map[string]string{
			"/Tapestry/Coach_Staff_20260818_1.csv": "EmployeeId,Name\n000111,Older\n",
			"/Tapestry/Coach_Staff_20260819_1.csv": "EmployeeId,Name\n000222,Newer\n",
		},
	}
	logs := newWorkdayFileLogStore()
	repository := &orderedEmployeeRepository{}
	job := &SyncCoachWorkdayEmployeeMaster{
		connect:    func(sftp.SSHConfig) (remoteFS, error) { return remote, nil },
		logs:       logs,
		repository: repository,
		now:        func() time.Time { return time.Unix(1787100000, 0) },
		notify:     func(app.Ctx, string, string, string) {},
	}
	config := micro_scrm.BrandConfigSftpDetail{}
	config.EnterpriseID = 3360610327843968
	config.BrandName = "coach"

	if err := job.syncConfig(toolboxCtx.New(), config); err != nil {
		t.Fatal(err)
	}
	if len(repository.codes) != 2 || repository.codes[0] != "111" || repository.codes[1] != "222" {
		t.Fatalf("processed codes = %v", repository.codes)
	}
	if len(logs.completed) != 2 || logs.completed[0].FileName != older.Name || logs.completed[1].FileName != newer.Name {
		t.Fatalf("completed files = %+v", logs.completed)
	}
}

func TestProcessWorkdayFileNotifiesAfterSuccessfulImport(t *testing.T) {
	file := remoteFile{Name: "Coach_Staff_20260831_000002.csv"}
	fullPath := "/Tapestry/" + file.Name
	remote := &fakeRemoteFS{files: map[string]string{
		fullPath: "EmployeeId,StoreCode,JobPosition\n000931251,OC267,Sales Associate\n",
	}}
	logs := &fakeFileLogStore{findErr: gorm.ErrRecordNotFound}
	repository := &fakeMasterDataRepository{employeeStats: syncStats{Inserted: 4, Updated: 5, Unchanged: 6}}
	var notifications []string

	err := processWorkdayFile(toolboxCtx.New(), 3360610327843968, "coach", remote, file, logs, repository, time.Unix(1787100000, 0), func(_ app.Ctx, _, _, message string) {
		notifications = append(notifications, message)
	})
	if err != nil {
		t.Fatal(err)
	}
	if logs.completedID != 91 {
		t.Fatalf("completedID=%d", logs.completedID)
	}
	if len(notifications) != 1 || notifications[0] != "成功：读取文件Coach_Staff_20260831_000002.csv，新增4，更新5，未变化6" {
		t.Fatalf("notifications=%v", notifications)
	}
}
