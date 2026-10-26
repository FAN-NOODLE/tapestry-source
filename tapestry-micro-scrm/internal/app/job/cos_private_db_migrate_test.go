package appJob

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/pkg/xjob"
)

type migrationEventLog struct {
	mu     sync.Mutex
	events []string
}

func (l *migrationEventLog) add(event string) {
	l.mu.Lock()
	l.events = append(l.events, event)
	l.mu.Unlock()
}

func (l *migrationEventLog) snapshot() []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return append([]string(nil), l.events...)
}

type fakeRows struct {
	mu          sync.Mutex
	events      *migrationEventLog
	data        map[int64]cosPrivateRow
	updates     []string
	updatedIDs  []int64
	processed   map[int64]struct{}
	conflict    bool
	conflictIDs map[int64]bool
	scanErr     error
	loadErr     error
	updateErr   error
	scanWait    <-chan struct{}
}

func newFakeRows(events *migrationEventLog) *fakeRows {
	return &fakeRows{
		events:      events,
		data:        make(map[int64]cosPrivateRow),
		processed:   make(map[int64]struct{}),
		conflictIDs: make(map[int64]bool),
	}
}

func (r *fakeRows) ScanAfter(requestCtx ctx.Context, _ cosPrivateField, afterID int64, limit int) ([]cosPrivateRow, error) {
	r.events.add(fmt.Sprintf("scan:%d", afterID))
	if r.scanWait != nil {
		select {
		case <-r.scanWait:
		case <-requestCtx.Done():
			return nil, requestCtx.Err()
		}
	}
	if err := requestCtx.Err(); err != nil {
		return nil, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.scanErr != nil {
		return nil, r.scanErr
	}
	ids := make([]int64, 0, len(r.data))
	for id := range r.data {
		if id > afterID {
			ids = append(ids, id)
		}
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	if len(ids) > limit {
		ids = ids[:limit]
	}
	rows := make([]cosPrivateRow, 0, len(ids))
	for _, id := range ids {
		rows = append(rows, r.data[id])
	}
	return rows, nil
}

func (r *fakeRows) MaxID(requestCtx ctx.Context, _ cosPrivateField) (int64, error) {
	if err := requestCtx.Err(); err != nil {
		return 0, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	var maxID int64
	for id := range r.data {
		if id > maxID {
			maxID = id
		}
	}
	return maxID, nil
}

func (r *fakeRows) LoadByID(requestCtx ctx.Context, _ cosPrivateField, id int64) (cosPrivateRow, bool, error) {
	r.events.add(fmt.Sprintf("load:%d", id))
	if err := requestCtx.Err(); err != nil {
		return cosPrivateRow{}, false, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.loadErr != nil {
		return cosPrivateRow{}, false, r.loadErr
	}
	row, ok := r.data[id]
	return row, ok, nil
}

func (r *fakeRows) UpdateIfUnchanged(requestCtx ctx.Context, _ cosPrivateField, row cosPrivateRow, value string) (bool, error) {
	if err := requestCtx.Err(); err != nil {
		return false, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.processed[row.ID] = struct{}{}
	if r.updateErr != nil {
		return false, r.updateErr
	}
	if r.conflict || r.conflictIDs[row.ID] {
		return false, nil
	}
	current, exists := r.data[row.ID]
	if !exists || current.Value != row.Value {
		return false, nil
	}
	current.Value = value
	r.data[row.ID] = current
	r.updates = append(r.updates, value)
	r.updatedIDs = append(r.updatedIDs, row.ID)
	return true, nil
}

func (r *fakeRows) snapshotUpdates() ([]int64, []string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]int64(nil), r.updatedIDs...), append([]string(nil), r.updates...)
}

func (r *fakeRows) wasProcessed(id int64) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.processed[id]
	return ok
}

var errFakeProgressWrite = errors.New("progress write failed")

type fakeProgress struct {
	mu                  sync.Mutex
	events              *migrationEventLog
	acquired            bool
	acquireErr          error
	renewErr            error
	releaseErr          error
	loadErr             error
	retryErr            error
	finishErr           error
	outstandingErr      error
	setCurrentErr       error
	loaded              cosPrivateProgress
	retryIDs            []int64
	failures            []cosPrivateFailure
	failureAdvance      []bool
	failureIDs          map[int64]struct{}
	successIDs          []int64
	skippedIDs          []int64
	resolvedIDs         []int64
	lastSkippedID       int64
	failAfterRowID      int64
	finishes            []string
	currentStatuses     []string
	lastToken           string
	renewCount          int
	released            bool
	renewAfterRelease   bool
	renewSignal         chan struct{}
	renewSignalOnce     sync.Once
	outstanding         int64
	outstandingExplicit bool
}

func newFakeProgress(events *migrationEventLog) *fakeProgress {
	return &fakeProgress{
		events:     events,
		acquired:   true,
		failureIDs: make(map[int64]struct{}),
	}
}

func (p *fakeProgress) Acquire(_ ctx.Context, token string, _ time.Duration) (bool, error) {
	p.mu.Lock()
	p.lastToken = token
	acquired, err := p.acquired, p.acquireErr
	p.mu.Unlock()
	p.events.add("acquire")
	return acquired, err
}

func (p *fakeProgress) Renew(_ ctx.Context, token string, _ time.Duration) error {
	p.mu.Lock()
	if p.released {
		p.renewAfterRelease = true
	}
	p.renewCount++
	err := p.renewErr
	wantToken := p.lastToken
	p.mu.Unlock()
	p.events.add("renew")
	if p.renewSignal != nil {
		p.renewSignalOnce.Do(func() { close(p.renewSignal) })
	}
	if token != wantToken {
		return errors.New("renewed with wrong token")
	}
	return err
}

func (p *fakeProgress) Release(_ ctx.Context, token string) error {
	p.mu.Lock()
	p.released = true
	err := p.releaseErr
	wantToken := p.lastToken
	p.mu.Unlock()
	p.events.add("release")
	if token != wantToken {
		return errors.New("released with wrong token")
	}
	return err
}

func (p *fakeProgress) Load(_ ctx.Context, _ cosPrivateField) (cosPrivateProgress, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.loadErr != nil {
		return cosPrivateProgress{}, p.loadErr
	}
	result := p.loaded
	result.OutstandingFailureCount = int64(len(p.failureIDs))
	return result, nil
}

func (p *fakeProgress) RetryIDs(_ ctx.Context, _ cosPrivateField) ([]int64, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return append([]int64(nil), p.retryIDs...), p.retryErr
}

func (p *fakeProgress) RecordSuccess(_ ctx.Context, _ cosPrivateField, rowID int64) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if rowID == p.failAfterRowID {
		return errFakeProgressWrite
	}
	p.successIDs = append(p.successIDs, rowID)
	delete(p.failureIDs, rowID)
	return nil
}

func (p *fakeProgress) RecordSkip(_ ctx.Context, _ cosPrivateField, rowID int64) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if rowID == p.failAfterRowID {
		return errFakeProgressWrite
	}
	p.skippedIDs = append(p.skippedIDs, rowID)
	p.lastSkippedID = rowID
	delete(p.failureIDs, rowID)
	return nil
}

func (p *fakeProgress) RecordFailure(_ ctx.Context, _ cosPrivateField, failure cosPrivateFailure, advance bool) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if failure.RowID == p.failAfterRowID {
		return errFakeProgressWrite
	}
	p.failures = append(p.failures, failure)
	p.failureAdvance = append(p.failureAdvance, advance)
	p.failureIDs[failure.RowID] = struct{}{}
	return nil
}

func (p *fakeProgress) ResolveFailure(_ ctx.Context, _ cosPrivateField, rowID int64) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if rowID == p.failAfterRowID {
		return errFakeProgressWrite
	}
	p.resolvedIDs = append(p.resolvedIDs, rowID)
	delete(p.failureIDs, rowID)
	return nil
}

func (p *fakeProgress) SetCurrent(_ ctx.Context, _ cosPrivateField, rowID int64, status string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.setCurrentErr != nil {
		return p.setCurrentErr
	}
	p.currentStatuses = append(p.currentStatuses, status)
	return nil
}

func (p *fakeProgress) Finish(_ ctx.Context, status string) error {
	p.mu.Lock()
	p.finishes = append(p.finishes, status)
	err := p.finishErr
	p.mu.Unlock()
	p.events.add("finish:" + status)
	return err
}

func (p *fakeProgress) OutstandingFailures(_ ctx.Context) (int64, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.outstandingErr != nil {
		return 0, p.outstandingErr
	}
	if p.outstandingExplicit {
		return p.outstanding, nil
	}
	return int64(len(p.failureIDs)), nil
}

func (p *fakeProgress) snapshot() (failures []cosPrivateFailure, advances []bool, finishes []string, renewCount int, renewAfterRelease bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return append([]cosPrivateFailure(nil), p.failures...), append([]bool(nil), p.failureAdvance...), append([]string(nil), p.finishes...), p.renewCount, p.renewAfterRelease
}

type migrationFixture struct {
	job      *CosPrivateDbMigrate
	field    cosPrivateField
	rows     *fakeRows
	progress *fakeProgress
	objects  *fakeCosStore
	events   *migrationEventLog
}

func newMigrationFixture() *migrationFixture {
	events := &migrationEventLog{}
	rows := newFakeRows(events)
	progress := newFakeProgress(events)
	objects := &fakeCosStore{objects: map[string]int64{}, copiedSize: 10, failKeys: map[string]error{}}
	job := &CosPrivateDbMigrate{
		rows:      rows,
		copier:    &cosPrivateCopier{store: objects, sourceBucket: "old", targetBucket: "private"},
		progress:  progress,
		rewriter:  newCosPrivateURLRewriter("old.example", "cdn.example", "private.example"),
		batchSize: 500,
	}
	return &migrationFixture{
		job: job, field: cosPrivateField{"micro_scrm", "sample", "url"}, rows: rows,
		progress: progress, objects: objects, events: events,
	}
}

func (f *migrationFixture) addRow(id int64, value string) {
	f.rows.mu.Lock()
	f.rows.data[id] = cosPrivateRow{ID: id, Value: value}
	f.rows.mu.Unlock()
	for _, item := range f.job.rewriter.extract(value) {
		f.objects.objects["old/"+item.Key] = 10
	}
}

func (f *migrationFixture) failCopy(key string) {
	f.objects.failKeys[key] = &cosPrivateObjectRowError{err: fmt.Errorf("%w: %s", errCosPrivateSourceObjectNotFound, key)}
}

func (f *migrationFixture) runWithField() error {
	f.job.fields = func(*xjob.JobHandler) ([]cosPrivateField, error) {
		return []cosPrivateField{f.field}, nil
	}
	return f.job.Run(ctx.New(), nil)
}

func TestMigrateRowUpdatesOnlyAfterEveryCopySucceeds(t *testing.T) {
	f := newMigrationFixture()
	f.addRow(1, "https://old.example/a https://old.example/b")
	f.failCopy("b")

	err := f.job.migrateRow(ctx.New(), f.field, f.rows.data[1], true)
	_, updates := f.rows.snapshotUpdates()
	failures, advances, _, _, _ := f.progress.snapshot()
	if err == nil || len(updates) != 0 || len(failures) != 1 || !reflect.DeepEqual(advances, []bool{true}) {
		t.Fatalf("err=%v updates=%v failures=%#v advances=%v", err, updates, failures, advances)
	}
}

func TestMigrateRowPreservesPublicAndUpdatesPrivate(t *testing.T) {
	f := newMigrationFixture()
	f.addRow(1, "https://cdn.example/uploads/a https://old.example/b")

	if err := f.job.migrateRow(ctx.New(), f.field, f.rows.data[1], true); err != nil {
		t.Fatal(err)
	}
	_, updates := f.rows.snapshotUpdates()
	want := "https://cdn.example/uploads/a https://private.example/b"
	if !reflect.DeepEqual(updates, []string{want}) {
		t.Fatalf("updates=%q want=%q", updates, want)
	}
}

func TestMigrateRowMovesSelectedPNGAndJPGFieldsIntoUploads(t *testing.T) {
	tests := []struct {
		name  string
		field cosPrivateField
		url   string
		want  string
		key   string
	}{
		{
			name:  "member poster PNG",
			field: cosPrivateField{Database: "micro_scrm", Table: "members", Field: "qr_poster_code"},
			url:   "https://old.example/2024/08/member.PNG?version=1",
			want:  "https://cdn.example/uploads/2024/08/member.PNG?version=1",
			key:   "old/uploads/2024/08/member.PNG",
		},
		{
			name:  "activity point JPG",
			field: cosPrivateField{Database: "micro_scrm_enterprise_1", Table: "activity_qrcode_point", Field: "code_url"},
			url:   "https://cdn.example/qrcode/point.jpg",
			want:  "https://cdn.example/uploads/qrcode/point.jpg",
			key:   "old/uploads/qrcode/point.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newMigrationFixture()
			f.field = tt.field
			f.addRow(1, tt.url)

			if err := f.job.migrateRow(ctx.New(), f.field, f.rows.data[1], true); err != nil {
				t.Fatal(err)
			}
			_, updates := f.rows.snapshotUpdates()
			if !reflect.DeepEqual(updates, []string{tt.want}) {
				t.Fatalf("updates=%q want=%q", updates, tt.want)
			}
			if _, exists := f.objects.objects[tt.key]; !exists {
				t.Fatalf("public target %q was not copied: %#v", tt.key, f.objects.objects)
			}
		})
	}
}

func TestMigrateRowKeepsOtherSelectedFieldExtensionsPrivate(t *testing.T) {
	f := newMigrationFixture()
	f.field = cosPrivateField{Database: "micro_scrm", Table: "members", Field: "qr_poster_code"}
	f.addRow(1, "https://old.example/2024/08/member.jpeg")

	if err := f.job.migrateRow(ctx.New(), f.field, f.rows.data[1], true); err != nil {
		t.Fatal(err)
	}
	_, updates := f.rows.snapshotUpdates()
	want := []string{"https://private.example/2024/08/member.jpeg"}
	if !reflect.DeepEqual(updates, want) {
		t.Fatalf("updates=%q want=%q", updates, want)
	}
}

func TestMigrateRowConflictCreatesRetryFailure(t *testing.T) {
	f := newMigrationFixture()
	f.addRow(1, "https://old.example/a")
	f.rows.conflict = true

	if err := f.job.migrateRow(ctx.New(), f.field, f.rows.data[1], true); err == nil {
		t.Fatal("expected conflict")
	}
	failures, _, _, _, _ := f.progress.snapshot()
	if len(failures) != 1 || !strings.Contains(failures[0].Reason, "concurrent") {
		t.Fatalf("failures=%#v", failures)
	}
}

func TestRunRetriesFailuresBeforeCheckpointScan(t *testing.T) {
	f := newMigrationFixture()
	f.progress.retryIDs = []int64{2}
	f.addRow(2, "https://old.example/a")
	f.addRow(10, "https://old.example/b")

	if err := f.job.migrateField(ctx.New(), f.field); err != nil {
		t.Fatal(err)
	}
	events := f.events.snapshot()
	if len(events) < 2 || !reflect.DeepEqual(events[:2], []string{"load:2", "scan:0"}) {
		t.Fatalf("events=%v", events)
	}
}

func TestRunAdvancesCheckpointForPublicOnlyRow(t *testing.T) {
	f := newMigrationFixture()
	f.addRow(3, "https://cdn.example/uploads/a")

	if err := f.job.migrateRow(ctx.New(), f.field, f.rows.data[3], true); err != nil {
		t.Fatal(err)
	}
	if f.progress.lastSkippedID != 3 {
		t.Fatalf("skip id=%d", f.progress.lastSkippedID)
	}
}

func TestRunAuditReportsLowIDOldPrivateURLWithoutMutation(t *testing.T) {
	f := newMigrationFixture()
	f.addRow(1, "https://old.example/a")

	if err := f.job.auditField(ctx.New(), f.field); err == nil {
		t.Fatal("expected audit failure")
	}
	ids, updates := f.rows.snapshotUpdates()
	failures, advances, _, _, _ := f.progress.snapshot()
	if len(failures) != 1 || len(ids) != 0 || len(updates) != 0 || f.objects.copyCalls != 0 || !reflect.DeepEqual(advances, []bool{false}) {
		t.Fatalf("failures=%#v ids=%v updates=%v copies=%d advances=%v", failures, ids, updates, f.objects.copyCalls, advances)
	}
	events := f.events.snapshot()
	if len(events) == 0 || events[0] != "scan:0" {
		t.Fatalf("events=%v", events)
	}
}

func TestRunStopsOnProgressWriteError(t *testing.T) {
	f := newMigrationFixture()
	f.addRow(1, "https://old.example/a")
	f.addRow(2, "https://old.example/b")
	f.progress.failAfterRowID = 1

	if err := f.job.migrateField(ctx.New(), f.field); !errors.Is(err, errFakeProgressWrite) {
		t.Fatalf("error=%v", err)
	}
	if f.rows.wasProcessed(2) {
		t.Fatal("row 2 processed after progress error")
	}
}

func TestMigrateFieldContinuesAfterRecordedRowFailure(t *testing.T) {
	f := newMigrationFixture()
	f.addRow(1, "https://old.example/a")
	f.addRow(2, "https://old.example/b")
	f.failCopy("a")

	if err := f.job.migrateField(ctx.New(), f.field); err != nil {
		t.Fatal(err)
	}
	ids, _ := f.rows.snapshotUpdates()
	failures, advances, _, _, _ := f.progress.snapshot()
	if !reflect.DeepEqual(ids, []int64{2}) || len(failures) != 1 || !reflect.DeepEqual(advances, []bool{true}) {
		t.Fatalf("updated=%v failures=%#v advances=%v", ids, failures, advances)
	}
}

func TestMigrateFieldRetryFailureDoesNotAdvanceCursor(t *testing.T) {
	f := newMigrationFixture()
	f.progress.retryIDs = []int64{2}
	f.progress.loaded.LastID = 7
	f.addRow(2, "https://old.example/a")
	f.failCopy("a")

	if err := f.job.migrateField(ctx.New(), f.field); err != nil {
		t.Fatal(err)
	}
	failures, advances, _, _, _ := f.progress.snapshot()
	if len(failures) != 1 || !reflect.DeepEqual(advances, []bool{false}) {
		t.Fatalf("failures=%#v advances=%v", failures, advances)
	}
	events := f.events.snapshot()
	if !reflect.DeepEqual(events, []string{"load:2"}) {
		t.Fatalf("events=%v", events)
	}
}

func TestMigrateFieldResumesFromPersistentCheckpoint(t *testing.T) {
	f := newMigrationFixture()
	f.progress.loaded.LastID = 7
	f.addRow(3, "https://old.example/a")
	f.addRow(10, "https://old.example/b")

	if err := f.job.migrateField(ctx.New(), f.field); err != nil {
		t.Fatal(err)
	}
	ids, _ := f.rows.snapshotUpdates()
	if !reflect.DeepEqual(ids, []int64{10}) {
		t.Fatalf("updated=%v", ids)
	}
	events := f.events.snapshot()
	if len(events) == 0 || events[0] != "scan:7" {
		t.Fatalf("events=%v", events)
	}
}

func TestMigrateFieldResolvesDeletedRetryRow(t *testing.T) {
	f := newMigrationFixture()
	f.progress.retryIDs = []int64{2}

	if err := f.job.migrateField(ctx.New(), f.field); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(f.progress.resolvedIDs, []int64{2}) {
		t.Fatalf("resolved=%v", f.progress.resolvedIDs)
	}
}

func TestRunRejectsAnotherLockOwnerWithoutBusinessMutation(t *testing.T) {
	f := newMigrationFixture()
	f.progress.acquired = false
	f.job.fields = func(*xjob.JobHandler) ([]cosPrivateField, error) {
		t.Fatal("fields called without lock")
		return nil, nil
	}

	err := f.job.Run(ctx.New(), nil)
	_, _, finishes, _, _ := f.progress.snapshot()
	if err == nil || !strings.Contains(err.Error(), "already running") || len(finishes) != 0 || f.progress.released {
		t.Fatalf("err=%v finishes=%v released=%t", err, finishes, f.progress.released)
	}
}

func TestRunReturnsPartialFailedWhenFailuresRemain(t *testing.T) {
	f := newMigrationFixture()
	f.addRow(1, "https://old.example/a")
	f.failCopy("a")

	err := f.runWithField()
	_, _, finishes, _, _ := f.progress.snapshot()
	if err == nil || !strings.Contains(err.Error(), "outstanding") || !reflect.DeepEqual(finishes, []string{"partial_failed"}) {
		t.Fatalf("err=%v finishes=%v", err, finishes)
	}
}

func TestRunCompletesAfterScanAndAuditHaveNoFailures(t *testing.T) {
	f := newMigrationFixture()
	f.addRow(1, "https://old.example/a")

	if err := f.runWithField(); err != nil {
		t.Fatal(err)
	}
	_, _, finishes, _, _ := f.progress.snapshot()
	if !reflect.DeepEqual(finishes, []string{"completed"}) || !f.progress.released {
		t.Fatalf("finishes=%v released=%t", finishes, f.progress.released)
	}
}

func TestLegacyCosPrivateJobsFailClosed(t *testing.T) {
	jobs := []xjob.Job{NewCosPrivateCopy(), NewCosPrivateDbURLRefresh(), NewCosPrivateDelete()}
	for _, job := range jobs {
		err := job.Run(ctx.New(), nil)
		if !errors.Is(err, errCosPrivateLegacyJobDisabled) || !strings.Contains(err.Error(), "cos_private_db_migrate") {
			t.Fatalf("legacy job error=%v", err)
		}
	}
}

func TestRunReportsFinishAndReleaseErrors(t *testing.T) {
	f := newMigrationFixture()
	finishErr := errors.New("finish failed")
	releaseErr := errors.New("release failed")
	f.progress.finishErr = finishErr
	f.progress.releaseErr = releaseErr
	f.job.fields = func(*xjob.JobHandler) ([]cosPrivateField, error) { return nil, nil }

	err := f.job.Run(ctx.New(), nil)
	if !errors.Is(err, finishErr) || !errors.Is(err, releaseErr) || !strings.Contains(err.Error(), finishErr.Error()) || !strings.Contains(err.Error(), releaseErr.Error()) {
		t.Fatalf("err=%v", err)
	}
}

type cancelAwareCosStore struct {
	copyStarted chan struct{}
	startOnce   sync.Once
}

func (s *cancelAwareCosStore) Stat(_ ctx.Context, bucket, _ string) (cosPrivateObjectMeta, bool, error) {
	if bucket == "old" {
		return cosPrivateObjectMeta{Size: 10}, true, nil
	}
	return cosPrivateObjectMeta{}, false, nil
}

func (s *cancelAwareCosStore) Copy(requestCtx ctx.Context, _, _, _, _ string) error {
	s.startOnce.Do(func() { close(s.copyStarted) })
	<-requestCtx.Done()
	return requestCtx.Err()
}

func TestRunRenewalFailureCancelsWorkBeforeDbOrProgressWrites(t *testing.T) {
	f := newMigrationFixture()
	f.addRow(1, "https://old.example/a")
	renewErr := errors.New("renew failed")
	f.progress.renewErr = renewErr
	f.progress.renewSignal = make(chan struct{})
	store := &cancelAwareCosStore{copyStarted: make(chan struct{})}
	f.job.copier = &cosPrivateCopier{store: store, sourceBucket: "old", targetBucket: "private"}
	f.job.renewInterval = time.Millisecond
	f.job.lockTTL = 50 * time.Millisecond
	f.job.fields = func(*xjob.JobHandler) ([]cosPrivateField, error) { return []cosPrivateField{f.field}, nil }

	err := f.job.Run(ctx.New(), nil)
	ids, updates := f.rows.snapshotUpdates()
	failures, _, finishes, renewCount, _ := f.progress.snapshot()
	if !errors.Is(err, renewErr) || len(ids) != 0 || len(updates) != 0 || len(failures) != 0 || len(finishes) != 0 || renewCount != 1 || !f.progress.released {
		t.Fatalf("err=%v ids=%v updates=%v failures=%#v finishes=%v renews=%d released=%t", err, ids, updates, failures, finishes, renewCount, f.progress.released)
	}
}

func TestRunStopsRenewalBeforeReleaseWithoutGoroutineLeak(t *testing.T) {
	f := newMigrationFixture()
	f.progress.renewSignal = make(chan struct{})
	f.rows.scanWait = f.progress.renewSignal
	f.job.renewInterval = time.Millisecond
	f.job.lockTTL = 50 * time.Millisecond
	f.job.fields = func(*xjob.JobHandler) ([]cosPrivateField, error) { return []cosPrivateField{f.field}, nil }

	if err := f.job.Run(ctx.New(), nil); err != nil {
		t.Fatal(err)
	}
	_, _, _, renewCount, renewAfterRelease := f.progress.snapshot()
	time.Sleep(5 * time.Millisecond)
	_, _, _, laterRenewCount, laterRenewAfterRelease := f.progress.snapshot()
	if renewCount == 0 || laterRenewCount != renewCount || renewAfterRelease || laterRenewAfterRelease {
		t.Fatalf("renews=%d later=%d after_release=%t/%t", renewCount, laterRenewCount, renewAfterRelease, laterRenewAfterRelease)
	}
}
