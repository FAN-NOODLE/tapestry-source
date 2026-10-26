package appJob

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	toolboxCtx "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/config"
	"micro-scrm/internal/pkg/db/gormx"
	pkgRedis "micro-scrm/internal/pkg/redis"
	"micro-scrm/internal/pkg/xjob"
)

const cosPrivateDefaultBatchSize = 500

var (
	errCosPrivateMigrationAlreadyRunning = errors.New("COS private DB migration is already running")
	errCosPrivateAuditFoundURLs          = errors.New("COS private audit found outstanding old private URLs")
	errCosPrivateOutstandingFailures     = errors.New("COS private migration has outstanding failures")
)

type CosPrivateDbMigrate struct {
	rows      cosPrivateRowRepository
	copier    *cosPrivateCopier
	progress  cosPrivateProgressStore
	rewriter  cosPrivateURLRewriter
	fields    func(*xjob.JobHandler) ([]cosPrivateField, error)
	batchSize int

	lockTTL       time.Duration
	renewInterval time.Duration
}

type cosPrivateRowMigrationError struct {
	err error
}

func (e *cosPrivateRowMigrationError) Error() string {
	return e.err.Error()
}

func (e *cosPrivateRowMigrationError) Unwrap() error {
	return e.err
}

type cosPrivateAuditError struct {
	count int
}

func (e *cosPrivateAuditError) Error() string {
	return fmt.Sprintf("%s: %d row(s)", errCosPrivateAuditFoundURLs, e.count)
}

func (e *cosPrivateAuditError) Unwrap() error {
	return errCosPrivateAuditFoundURLs
}

type cosPrivateCombinedError struct {
	errs []error
}

func (e *cosPrivateCombinedError) Error() string {
	parts := make([]string, 0, len(e.errs))
	for _, err := range e.errs {
		parts = append(parts, err.Error())
	}
	return strings.Join(parts, "; ")
}

func (e *cosPrivateCombinedError) Unwrap() error {
	if len(e.errs) == 0 {
		return nil
	}
	return e.errs[0]
}

func (e *cosPrivateCombinedError) Is(target error) bool {
	for _, err := range e.errs {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

func combineCosPrivateErrors(errs ...error) error {
	combined := make([]error, 0, len(errs))
	for _, err := range errs {
		if err != nil {
			combined = append(combined, err)
		}
	}
	if len(combined) == 0 {
		return nil
	}
	if len(combined) == 1 {
		return combined[0]
	}
	return &cosPrivateCombinedError{errs: combined}
}

func NewCosPrivateDbMigrate() *CosPrivateDbMigrate {
	storage := config.CfgData.Storage.COS
	sourceBucket := storage.Bucket
	targetBucket := cosPrivateBucket()
	oldOriginHost := cosOriginHost(sourceBucket)
	privateHost := cosOriginHost(targetBucket)
	return &CosPrivateDbMigrate{
		rows: cosPrivateRowStore{
			db:            gormx.GetInstance(),
			oldOriginHost: oldOriginHost,
			oldCdnHost:    storage.CdnHost,
		},
		copier: &cosPrivateCopier{
			store:        newCosPrivateTencentStore(sourceBucket, targetBucket),
			sourceBucket: sourceBucket,
			targetBucket: targetBucket,
		},
		progress:      newCosPrivateRedisProgressStore(pkgRedis.Client(), config.CfgData.Env),
		rewriter:      newCosPrivateURLRewriter(oldOriginHost, storage.CdnHost, privateHost),
		fields:        cosPrivateMigrationFields,
		batchSize:     cosPrivateDefaultBatchSize,
		lockTTL:       cosPrivateLockTTL,
		renewInterval: cosPrivateLockRenewInterval,
	}
}

func (job *CosPrivateDbMigrate) Run(requestCtx toolboxCtx.Context, handler *xjob.JobHandler) error {
	token, err := newCosPrivateLockToken()
	if err != nil {
		return err
	}
	lockTTL := job.effectiveLockTTL()
	acquired, err := job.progress.Acquire(requestCtx, token, lockTTL)
	if err != nil {
		return fmt.Errorf("acquire COS private migration lock: %w", err)
	}
	if !acquired {
		return errCosPrivateMigrationAlreadyRunning
	}

	workContext, cancelWork := context.WithCancel(requestCtx)
	workCtx := toolboxCtx.Wrap(workContext)
	renewal := job.startLockRenewal(workCtx, cancelWork, token, lockTTL)

	status, runErr := job.runLocked(workCtx, handler)
	renewErr := renewal.stop()
	if renewErr != nil {
		runErr = combineCosPrivateErrors(renewErr, runErr)
	}

	var finishErr error
	if renewErr == nil {
		if err := job.progress.Finish(requestCtx, status); err != nil {
			finishErr = fmt.Errorf("finish COS private migration with status %s: %w", status, err)
		}
	}
	releaseErr := job.progress.Release(requestCtx, token)
	if releaseErr != nil {
		releaseErr = fmt.Errorf("release COS private migration lock: %w", releaseErr)
	}
	cancelWork()
	return combineCosPrivateErrors(runErr, finishErr, releaseErr)
}

func (job *CosPrivateDbMigrate) runLocked(requestCtx toolboxCtx.Context, handler *xjob.JobHandler) (string, error) {
	fields, err := job.fields(handler)
	if err != nil {
		return job.failureStatus(requestCtx, fmt.Errorf("load COS private migration fields: %w", err))
	}
	for _, field := range fields {
		if err := job.migrateField(requestCtx, field); err != nil {
			return job.failureStatus(requestCtx, err)
		}
		auditErr := job.auditField(requestCtx, field)
		if auditErr != nil && !errors.Is(auditErr, errCosPrivateAuditFoundURLs) {
			return job.failureStatus(requestCtx, auditErr)
		}
		progress, err := job.progress.Load(requestCtx, field)
		if err != nil {
			return job.failureStatus(requestCtx, fmt.Errorf("load audited progress for %s: %w", field.String(), err))
		}
		fieldStatus := "completed"
		if progress.OutstandingFailureCount > 0 || auditErr != nil {
			fieldStatus = "partial_failed"
		}
		if err := job.progress.SetCurrent(requestCtx, field, 0, fieldStatus); err != nil {
			return job.failureStatus(requestCtx, fmt.Errorf("finish migration field %s: %w", field.String(), err))
		}
	}

	outstanding, err := job.progress.OutstandingFailures(requestCtx)
	if err != nil {
		return "failed", fmt.Errorf("load COS private outstanding failures: %w", err)
	}
	if outstanding > 0 {
		return "partial_failed", fmt.Errorf("%w: %d", errCosPrivateOutstandingFailures, outstanding)
	}
	return "completed", nil
}

func (job *CosPrivateDbMigrate) failureStatus(requestCtx toolboxCtx.Context, runErr error) (string, error) {
	if requestCtx.Err() != nil {
		return "failed", runErr
	}
	outstanding, err := job.progress.OutstandingFailures(requestCtx)
	if err != nil {
		return "failed", combineCosPrivateErrors(runErr, fmt.Errorf("load COS private outstanding failures: %w", err))
	}
	if outstanding > 0 {
		return "partial_failed", combineCosPrivateErrors(
			runErr,
			fmt.Errorf("%w: %d", errCosPrivateOutstandingFailures, outstanding),
		)
	}
	return "failed", runErr
}

func (job *CosPrivateDbMigrate) migrateField(requestCtx toolboxCtx.Context, field cosPrivateField) error {
	highWater, err := job.rows.MaxID(requestCtx, field)
	if err != nil {
		return fmt.Errorf("load migration high-water for %s: %w", field.String(), err)
	}
	retryIDs, err := job.progress.RetryIDs(requestCtx, field)
	if err != nil {
		return fmt.Errorf("load retry IDs for %s: %w", field.String(), err)
	}
	for _, rowID := range retryIDs {
		if err := requestCtx.Err(); err != nil {
			return err
		}
		row, exists, err := job.rows.LoadByID(requestCtx, field, rowID)
		if err != nil {
			return fmt.Errorf("load retry row %s id=%d: %w", field.String(), rowID, err)
		}
		if !exists {
			if err := job.progress.ResolveFailure(requestCtx, field, rowID); err != nil {
				return fmt.Errorf("resolve missing retry row %s id=%d: %w", field.String(), rowID, err)
			}
			continue
		}
		if err := job.migrateRow(requestCtx, field, row, false); err != nil {
			var rowErr *cosPrivateRowMigrationError
			if errors.As(err, &rowErr) {
				continue
			}
			return err
		}
	}

	progress, err := job.progress.Load(requestCtx, field)
	if err != nil {
		return fmt.Errorf("load progress for %s: %w", field.String(), err)
	}
	afterID := progress.LastID
	batchSize := job.effectiveBatchSize()
	if afterID >= highWater {
		return nil
	}
	for {
		if err := requestCtx.Err(); err != nil {
			return err
		}
		rows, err := job.rows.ScanAfter(requestCtx, field, afterID, batchSize)
		if err != nil {
			return fmt.Errorf("scan rows for %s after id=%d: %w", field.String(), afterID, err)
		}
		reachedHighWater := false
		for _, row := range rows {
			if row.ID > highWater {
				reachedHighWater = true
				break
			}
			if err := requestCtx.Err(); err != nil {
				return err
			}
			err := job.migrateRow(requestCtx, field, row, true)
			if err != nil {
				var rowErr *cosPrivateRowMigrationError
				if !errors.As(err, &rowErr) {
					return err
				}
			}
			afterID = row.ID
		}
		if reachedHighWater || afterID >= highWater || len(rows) < batchSize {
			return nil
		}
	}
}

func (job *CosPrivateDbMigrate) migrateRow(requestCtx toolboxCtx.Context, field cosPrivateField, row cosPrivateRow, advanceCursor bool) error {
	if err := requestCtx.Err(); err != nil {
		return err
	}
	if err := job.progress.SetCurrent(requestCtx, field, row.ID, "running"); err != nil {
		return fmt.Errorf("set current migration row %s id=%d: %w", field.String(), row.ID, err)
	}
	items := job.rewriter.extract(row.Value)
	if len(items) == 0 {
		var err error
		if advanceCursor {
			err = job.progress.RecordSkip(requestCtx, field, row.ID)
		} else {
			err = job.progress.ResolveFailure(requestCtx, field, row.ID)
		}
		if err != nil {
			return fmt.Errorf("record skipped migration row %s id=%d: %w", field.String(), row.ID, err)
		}
		return nil
	}

	privateCopied := make(map[string]struct{}, len(items))
	publicCopied := make(map[string]struct{}, len(items))
	for _, item := range items {
		if err := requestCtx.Err(); err != nil {
			return err
		}
		var copyErr error
		if cosPrivateShouldMoveImagePublic(field, item.Key) {
			copyErr = job.copier.ensurePublicCopied(requestCtx, item.Key, cosPrivatePublicTargetKey(item.Key))
		} else {
			copyErr = job.copier.ensureCopied(requestCtx, item.Key)
		}
		if copyErr != nil {
			if requestCtx.Err() != nil {
				return requestCtx.Err()
			}
			if !isCosPrivateObjectRowError(copyErr) {
				return fmt.Errorf("copy migration infrastructure failure %s id=%d: %w", field.String(), row.ID, copyErr)
			}
			failure := newCosPrivateFailure(row.ID, row.Value, item.RawURL, copyErr, time.Now())
			if recordErr := job.progress.RecordFailure(requestCtx, field, failure, advanceCursor); recordErr != nil {
				return combineCosPrivateErrors(
					fmt.Errorf("copy migration row %s id=%d: %w", field.String(), row.ID, copyErr),
					fmt.Errorf("record migration row failure %s id=%d: %w", field.String(), row.ID, recordErr),
				)
			}
			return &cosPrivateRowMigrationError{err: fmt.Errorf("copy migration row %s id=%d: %w", field.String(), row.ID, copyErr)}
		}
		if cosPrivateShouldMoveImagePublic(field, item.Key) {
			publicCopied[item.Key] = struct{}{}
		} else {
			privateCopied[item.Key] = struct{}{}
		}
	}

	if err := requestCtx.Err(); err != nil {
		return err
	}
	newValue := job.rewriter.rewrite(row.Value, privateCopied)
	newValue = job.rewriter.rewritePublic(newValue, publicCopied)
	updated, err := job.rows.UpdateIfUnchanged(requestCtx, field, row, newValue)
	if err != nil {
		if requestCtx.Err() != nil {
			return requestCtx.Err()
		}
		failure := newCosPrivateFailure(row.ID, row.Value, "", err, time.Now())
		if recordErr := job.progress.RecordFailure(requestCtx, field, failure, advanceCursor); recordErr != nil {
			return combineCosPrivateErrors(
				fmt.Errorf("update migration row %s id=%d: %w", field.String(), row.ID, err),
				fmt.Errorf("record migration row failure %s id=%d: %w", field.String(), row.ID, recordErr),
			)
		}
		return fmt.Errorf("update migration row %s id=%d: %w", field.String(), row.ID, err)
	}
	if !updated {
		conflictErr := fmt.Errorf("concurrent update prevented migration for %s id=%d", field.String(), row.ID)
		failure := newCosPrivateFailure(row.ID, row.Value, "", conflictErr, time.Now())
		if err := job.progress.RecordFailure(requestCtx, field, failure, advanceCursor); err != nil {
			return combineCosPrivateErrors(conflictErr, fmt.Errorf("record migration conflict %s id=%d: %w", field.String(), row.ID, err))
		}
		return &cosPrivateRowMigrationError{err: conflictErr}
	}
	if err := requestCtx.Err(); err != nil {
		return err
	}
	var resultErr error
	if advanceCursor {
		resultErr = job.progress.RecordSuccess(requestCtx, field, row.ID)
	} else {
		resultErr = job.progress.ResolveFailure(requestCtx, field, row.ID)
	}
	if resultErr != nil {
		return fmt.Errorf("record successful migration row %s id=%d: %w", field.String(), row.ID, resultErr)
	}
	return nil
}

func (job *CosPrivateDbMigrate) auditField(requestCtx toolboxCtx.Context, field cosPrivateField) error {
	highWater, err := job.rows.MaxID(requestCtx, field)
	if err != nil {
		return fmt.Errorf("load audit high-water for %s: %w", field.String(), err)
	}
	retryIDs, err := job.progress.RetryIDs(requestCtx, field)
	if err != nil {
		return fmt.Errorf("load audit failures for %s: %w", field.String(), err)
	}
	existingFailures := make(map[int64]struct{}, len(retryIDs))
	for _, rowID := range retryIDs {
		existingFailures[rowID] = struct{}{}
	}
	afterID := int64(0)
	batchSize := job.effectiveBatchSize()
	found := 0
	for {
		if err := requestCtx.Err(); err != nil {
			return err
		}
		rows, err := job.rows.ScanAfter(requestCtx, field, afterID, batchSize)
		if err != nil {
			return fmt.Errorf("audit rows for %s after id=%d: %w", field.String(), afterID, err)
		}
		reachedHighWater := false
		for _, row := range rows {
			if row.ID > highWater {
				reachedHighWater = true
				break
			}
			if err := requestCtx.Err(); err != nil {
				return err
			}
			items := job.rewriter.extract(row.Value)
			if len(items) > 0 {
				found++
				if _, exists := existingFailures[row.ID]; exists {
					afterID = row.ID
					continue
				}
				if err := job.progress.SetCurrent(requestCtx, field, row.ID, "auditing"); err != nil {
					return fmt.Errorf("set current audit row %s id=%d: %w", field.String(), row.ID, err)
				}
				auditErr := fmt.Errorf("old private URL remains after migration audit")
				failure := newCosPrivateFailure(row.ID, row.Value, items[0].RawURL, auditErr, time.Now())
				if err := job.progress.RecordFailure(requestCtx, field, failure, false); err != nil {
					return fmt.Errorf("record audit failure %s id=%d: %w", field.String(), row.ID, err)
				}
				existingFailures[row.ID] = struct{}{}
			}
			afterID = row.ID
		}
		if reachedHighWater || afterID >= highWater || len(rows) < batchSize {
			break
		}
	}
	if found > 0 {
		return &cosPrivateAuditError{count: found}
	}
	return nil
}

type cosPrivateLockRenewal struct {
	stopCh chan struct{}
	doneCh chan struct{}
	errCh  chan error
}

func (r *cosPrivateLockRenewal) stop() error {
	close(r.stopCh)
	<-r.doneCh
	select {
	case err := <-r.errCh:
		return err
	default:
		return nil
	}
}

func (job *CosPrivateDbMigrate) startLockRenewal(requestCtx toolboxCtx.Context, cancel context.CancelFunc, token string, ttl time.Duration) *cosPrivateLockRenewal {
	renewal := &cosPrivateLockRenewal{
		stopCh: make(chan struct{}),
		doneCh: make(chan struct{}),
		errCh:  make(chan error, 1),
	}
	interval := job.effectiveRenewInterval()
	go func() {
		defer close(renewal.doneCh)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-renewal.stopCh:
				return
			case <-requestCtx.Done():
				return
			case <-ticker.C:
				if err := job.progress.Renew(requestCtx, token, ttl); err != nil {
					renewErr := fmt.Errorf("renew COS private migration lock: %w", err)
					renewal.errCh <- renewErr
					cancel()
					return
				}
			}
		}
	}()
	return renewal
}

func (job *CosPrivateDbMigrate) effectiveBatchSize() int {
	if job.batchSize > 0 {
		return job.batchSize
	}
	return cosPrivateDefaultBatchSize
}

func (job *CosPrivateDbMigrate) effectiveLockTTL() time.Duration {
	if job.lockTTL > 0 {
		return job.lockTTL
	}
	return cosPrivateLockTTL
}

func (job *CosPrivateDbMigrate) effectiveRenewInterval() time.Duration {
	if job.renewInterval > 0 {
		return job.renewInterval
	}
	return cosPrivateLockRenewInterval
}

func newCosPrivateLockToken() (string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("generate COS private migration lock token: %w", err)
	}
	return hex.EncodeToString(raw), nil
}

func (field cosPrivateField) String() string {
	return field.Database + "." + field.Table + "." + field.Field
}
