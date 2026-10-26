package appJob

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

const (
	cosPrivateLockTTL           = 60 * time.Second
	cosPrivateLockRenewInterval = 20 * time.Second
	// cosPrivateRedisNamespaceVersion identifies the persistent Redis layout.
	// v4 restarts progress after adding public-image routing; failure ZSET members remain zero-padded 19-digit IDs at score 0.
	// change this version whenever a persisted key representation becomes incompatible.
	cosPrivateRedisNamespaceVersion = "v4"
)

var (
	errCosPrivateInvalidLockTTL = errors.New("COS private migration lock TTL must be at least one millisecond")
	errCosPrivateInvalidToken   = errors.New("COS private migration lock token cannot be empty")
	errCosPrivateLockLost       = errors.New("COS private migration lock ownership lost")
)

const cosPrivateReleaseLockScript = `
if redis.call('get',KEYS[1]) == ARGV[1] then
    return redis.call('del',KEYS[1])
end
return 0
`

const cosPrivateRenewLockScript = `
if redis.call('get',KEYS[1]) == ARGV[1] then
    return redis.call('pexpire',KEYS[1],ARGV[2])
end
return 0
`

const cosPrivateMutationLuaHelpers = `
local max_int64 = '9223372036854775807'

local function key_type(key)
    local result = redis.call('type',key)
    if type(result) == 'table' then
        return result.ok
    end
    return result
end

local function require_key_type(key,expected)
    local actual = key_type(key)
    if actual ~= 'none' and actual ~= expected then
        error('invalid Redis key type for '..key..': '..actual)
    end
end

local function valid_uint64(value)
    if type(value) ~= 'string' or value == '' then
        return false
    end
    if value ~= '0' and string.match(value,'^[1-9]%d*$') == nil then
        return false
    end
    if #value > 19 or (#value == 19 and value > max_int64) then
        return false
    end
    return true
end

local function valid_positive_int64(value)
    return valid_uint64(value) and value ~= '0'
end

local function require_row_id(row_id,row_member)
    if not valid_positive_int64(row_id) then
        error('invalid row ID')
    end
    if type(row_member) ~= 'string' or #row_member ~= 19 or string.match(row_member,'^%d+$') == nil then
        error('invalid row member')
    end
    local member_id = string.gsub(row_member,'^0+','')
    if member_id == '' then
        member_id = '0'
    end
    if member_id ~= row_id then
        error('row member does not match row ID')
    end
end

local function checked_counter(key,field)
    local value = redis.call('hget',key,field) or '0'
    if not valid_uint64(value) then
        error('invalid counter '..field)
    end
    return value
end

local function require_incrementable(value,field)
    if value == max_int64 then
        error('counter overflow '..field)
    end
end

local function require_decrementable(value,field)
    if value == '0' then
        error('counter underflow '..field)
    end
end

local function increment_decimal(value)
    require_incrementable(value,'retry_count')
    local result = ''
    local carry = 1
    for index = #value,1,-1 do
        local digit = string.byte(value,index) - 48 + carry
        if digit == 10 then
            digit = 0
            carry = 1
        else
            carry = 0
        end
        result = string.char(48 + digit)..result
    end
    if carry == 1 then
        result = '1'..result
    end
    return result
end

local function id_greater(left,right)
    if #left ~= #right then
        return #left > #right
    end
    return left > right
end

local lock_key = KEYS[1]
local lock_token = ARGV[1]
if redis.call('get',lock_key) ~= lock_token then
    return 0
end
`

const cosPrivateRecordFailureScript = cosPrivateMutationLuaHelpers + `
local progress_key = KEYS[2]
local failure_key = KEYS[3]
local detail_key = KEYS[4]
local task_key = KEYS[5]
local row_id = ARGV[2]
local row_member = ARGV[3]
local detail_prefix = ARGV[4]
local detail_suffix = ARGV[5]
local initial_retry_count = ARGV[6]
local updated_at = ARGV[7]
local advance_cursor = ARGV[8]

require_key_type(progress_key,'hash')
require_key_type(failure_key,'zset')
require_key_type(detail_key,'hash')
require_key_type(task_key,'hash')
require_row_id(row_id,row_member)
if not valid_positive_int64(initial_retry_count) then
    error('invalid initial retry count')
end
if advance_cursor ~= '0' and advance_cursor ~= '1' then
    error('invalid cursor flag')
end

local last_id = redis.call('hget',progress_key,'last_id') or '0'
if not valid_uint64(last_id) then
    error('invalid last_id')
end
local failure_total = checked_counter(progress_key,'failure_total')
require_incrementable(failure_total,'failure_total')
local existing = redis.call('zscore',failure_key,row_member)
local retry_count = initial_retry_count
if existing then
    local stored_retry_count = redis.call('hget',detail_key,'retry:'..row_id) or initial_retry_count
    if not valid_positive_int64(stored_retry_count) then
        error('invalid retry_count')
    end
    retry_count = increment_decimal(stored_retry_count)
else
    local progress_outstanding = checked_counter(progress_key,'outstanding_failure_count')
    local task_outstanding = checked_counter(task_key,'outstanding_failure_count')
    require_incrementable(progress_outstanding,'outstanding_failure_count')
    require_incrementable(task_outstanding,'task outstanding_failure_count')
end
if advance_cursor == '1' then
    local scanned_count = checked_counter(progress_key,'scanned_count')
    require_incrementable(scanned_count,'scanned_count')
end

if not existing then
    redis.call('hincrby',progress_key,'outstanding_failure_count',1)
    redis.call('hincrby',task_key,'outstanding_failure_count',1)
end
redis.call('zadd',failure_key,0,row_member)
redis.call('hset',detail_key,row_id,detail_prefix..retry_count..detail_suffix,'retry:'..row_id,retry_count)
redis.call('hincrby',progress_key,'failure_total',1)
if advance_cursor == '1' then
    if id_greater(row_id,last_id) then
        redis.call('hset',progress_key,'last_id',row_id)
    end
    redis.call('hincrby',progress_key,'scanned_count',1)
end
redis.call('hset',progress_key,'status','running','updated_at',updated_at)
redis.call('hset',task_key,'updated_at',updated_at)
return 1
`

const cosPrivateRecordSuccessScript = cosPrivateMutationLuaHelpers + `
local progress_key = KEYS[2]
local failure_key = KEYS[3]
local detail_key = KEYS[4]
local task_key = KEYS[5]
local row_id = ARGV[2]
local row_member = ARGV[3]
local updated_at = ARGV[4]

require_key_type(progress_key,'hash')
require_key_type(failure_key,'zset')
require_key_type(detail_key,'hash')
require_key_type(task_key,'hash')
require_row_id(row_id,row_member)
local last_id = redis.call('hget',progress_key,'last_id') or '0'
if not valid_uint64(last_id) then
    error('invalid last_id')
end
local success_count = checked_counter(progress_key,'success_count')
require_incrementable(success_count,'success_count')
local existing = redis.call('zscore',failure_key,row_member)
if existing then
    local progress_outstanding = checked_counter(progress_key,'outstanding_failure_count')
    local task_outstanding = checked_counter(task_key,'outstanding_failure_count')
    require_decrementable(progress_outstanding,'outstanding_failure_count')
    require_decrementable(task_outstanding,'task outstanding_failure_count')
else
    local scanned_count = checked_counter(progress_key,'scanned_count')
    require_incrementable(scanned_count,'scanned_count')
end

redis.call('zrem',failure_key,row_member)
redis.call('hdel',detail_key,row_id,'retry:'..row_id)
if existing then
    redis.call('hincrby',progress_key,'outstanding_failure_count',-1)
    redis.call('hincrby',task_key,'outstanding_failure_count',-1)
else
    redis.call('hincrby',progress_key,'scanned_count',1)
end
if id_greater(row_id,last_id) then
    redis.call('hset',progress_key,'last_id',row_id)
end
redis.call('hincrby',progress_key,'success_count',1)
redis.call('hset',progress_key,'status','running','updated_at',updated_at)
redis.call('hset',task_key,'updated_at',updated_at)
return 1
`

const cosPrivateRecordSkipScript = cosPrivateMutationLuaHelpers + `
local progress_key = KEYS[2]
local failure_key = KEYS[3]
local detail_key = KEYS[4]
local task_key = KEYS[5]
local row_id = ARGV[2]
local row_member = ARGV[3]
local updated_at = ARGV[4]

require_key_type(progress_key,'hash')
require_key_type(failure_key,'zset')
require_key_type(detail_key,'hash')
require_key_type(task_key,'hash')
require_row_id(row_id,row_member)
local last_id = redis.call('hget',progress_key,'last_id') or '0'
if not valid_uint64(last_id) then
    error('invalid last_id')
end
local skip_count = checked_counter(progress_key,'skip_count')
require_incrementable(skip_count,'skip_count')
local existing = redis.call('zscore',failure_key,row_member)
if existing then
    local progress_outstanding = checked_counter(progress_key,'outstanding_failure_count')
    local task_outstanding = checked_counter(task_key,'outstanding_failure_count')
    require_decrementable(progress_outstanding,'outstanding_failure_count')
    require_decrementable(task_outstanding,'task outstanding_failure_count')
else
    local scanned_count = checked_counter(progress_key,'scanned_count')
    require_incrementable(scanned_count,'scanned_count')
end

redis.call('zrem',failure_key,row_member)
redis.call('hdel',detail_key,row_id,'retry:'..row_id)
if existing then
    redis.call('hincrby',progress_key,'outstanding_failure_count',-1)
    redis.call('hincrby',task_key,'outstanding_failure_count',-1)
else
    redis.call('hincrby',progress_key,'scanned_count',1)
end
if id_greater(row_id,last_id) then
    redis.call('hset',progress_key,'last_id',row_id)
end
redis.call('hincrby',progress_key,'skip_count',1)
redis.call('hset',progress_key,'status','running','updated_at',updated_at)
redis.call('hset',task_key,'updated_at',updated_at)
return 1
`

const cosPrivateResolveFailureScript = cosPrivateMutationLuaHelpers + `
local progress_key = KEYS[2]
local failure_key = KEYS[3]
local detail_key = KEYS[4]
local task_key = KEYS[5]
local row_id = ARGV[2]
local row_member = ARGV[3]
local updated_at = ARGV[4]

require_key_type(progress_key,'hash')
require_key_type(failure_key,'zset')
require_key_type(detail_key,'hash')
require_key_type(task_key,'hash')
require_row_id(row_id,row_member)
local existing = redis.call('zscore',failure_key,row_member)
if existing then
    local progress_outstanding = checked_counter(progress_key,'outstanding_failure_count')
    local task_outstanding = checked_counter(task_key,'outstanding_failure_count')
    require_decrementable(progress_outstanding,'outstanding_failure_count')
    require_decrementable(task_outstanding,'task outstanding_failure_count')
end

redis.call('zrem',failure_key,row_member)
redis.call('hdel',detail_key,row_id,'retry:'..row_id)
if existing then
    redis.call('hincrby',progress_key,'outstanding_failure_count',-1)
    redis.call('hincrby',task_key,'outstanding_failure_count',-1)
end
redis.call('hset',progress_key,'updated_at',updated_at)
redis.call('hset',task_key,'updated_at',updated_at)
return 1
`

const cosPrivateSetCurrentScript = cosPrivateMutationLuaHelpers + `
local progress_key = KEYS[2]
local task_key = KEYS[3]
local database = ARGV[2]
local table_name = ARGV[3]
local field = ARGV[4]
local row_id = ARGV[5]
local status = ARGV[6]
local updated_at = ARGV[7]

require_key_type(progress_key,'hash')
require_key_type(task_key,'hash')
if not valid_uint64(row_id) then
    error('invalid current row ID')
end

if redis.call('hexists',progress_key,'started_at') == 0 then
    redis.call('hset',progress_key,'started_at',updated_at)
end
if redis.call('hexists',task_key,'started_at') == 0 then
    redis.call('hset',task_key,'started_at',updated_at)
end
redis.call('hset',progress_key,'status',status,'updated_at',updated_at)
redis.call('hset',task_key,
    'current_database',database,
    'current_table',table_name,
    'current_field',field,
    'current_id',row_id,
    'status',status,
    'updated_at',updated_at)
if status == 'completed' or status == 'partial_failed' or status == 'failed' then
    redis.call('hset',progress_key,'completed_at',updated_at)
end
return 1
`

const cosPrivateFinishScript = cosPrivateMutationLuaHelpers + `
local task_key = KEYS[2]
local status = ARGV[2]
local updated_at = ARGV[3]
require_key_type(task_key,'hash')
redis.call('hset',task_key,'status',status,'updated_at',updated_at,'completed_at',updated_at)
return 1
`

type cosPrivateProgress struct {
	LastID                  int64
	ScannedCount            int64
	SuccessCount            int64
	SkipCount               int64
	FailureTotal            int64
	OutstandingFailureCount int64
	Status                  string
}

type cosPrivateFailure struct {
	RowID      int64  `json:"row_id"`
	ValueHash  string `json:"value_hash"`
	FailedURL  string `json:"failed_url"`
	Reason     string `json:"reason"`
	RetryCount int64  `json:"retry_count"`
	FailedAt   int64  `json:"failed_at"`
}

type cosPrivateProgressStore interface {
	Acquire(ctx.Context, string, time.Duration) (bool, error)
	Renew(ctx.Context, string, time.Duration) error
	Release(ctx.Context, string) error
	Load(ctx.Context, cosPrivateField) (cosPrivateProgress, error)
	RetryIDs(ctx.Context, cosPrivateField) ([]int64, error)
	RecordSuccess(ctx.Context, cosPrivateField, int64) error
	RecordSkip(ctx.Context, cosPrivateField, int64) error
	RecordFailure(ctx.Context, cosPrivateField, cosPrivateFailure, bool) error
	ResolveFailure(ctx.Context, cosPrivateField, int64) error
	SetCurrent(ctx.Context, cosPrivateField, int64, string) error
	Finish(ctx.Context, string) error
	OutstandingFailures(ctx.Context) (int64, error)
}

type cosPrivateRedisClient interface {
	SetNX(context.Context, string, interface{}, time.Duration) *redis.BoolCmd
	Eval(context.Context, string, []string, ...interface{}) *redis.Cmd
	HGetAll(context.Context, string) *redis.StringStringMapCmd
	ZRange(context.Context, string, int64, int64) *redis.StringSliceCmd
	HGet(context.Context, string, string) *redis.StringCmd
}

type cosPrivateRedisProgressStore struct {
	client cosPrivateRedisClient
	env    string
	now    func() time.Time
	tokenM sync.RWMutex
	token  string
}

type cosPrivateRedisKeys struct {
	Lock           string
	Task           string
	Progress       string
	Failures       string
	FailureDetails string
}

func cosPrivateProgressKeys(env string, field cosPrivateField) cosPrivateRedisKeys {
	root := fmt.Sprintf("cos:private:migrate:%s:%s", cosPrivateRedisNamespaceVersion, cosPrivateRedisKeyComponent(env))
	fieldRoot := fmt.Sprintf(
		"%s:%s:%s:%s",
		root,
		cosPrivateRedisKeyComponent(field.Database),
		cosPrivateRedisKeyComponent(field.Table),
		cosPrivateRedisKeyComponent(field.Field),
	)
	return cosPrivateRedisKeys{
		Lock:           root + ":lock",
		Task:           root + ":task",
		Progress:       fieldRoot + ":progress",
		Failures:       fieldRoot + ":failures",
		FailureDetails: fieldRoot + ":failure_details",
	}
}

func newCosPrivateFailure(rowID int64, rawValue, failedURL string, failureErr error, failedAt time.Time) cosPrivateFailure {
	digest := sha256.Sum256([]byte(rawValue))
	valueHash := hex.EncodeToString(digest[:])
	reason := ""
	if failureErr != nil {
		reason = failureErr.Error()
		if rawValue != "" {
			reason = strings.ReplaceAll(reason, rawValue, "[redacted:"+valueHash[:12]+"]")
		}
	}
	return cosPrivateFailure{
		RowID:      rowID,
		ValueHash:  valueHash,
		FailedURL:  failedURL,
		Reason:     reason,
		RetryCount: 1,
		FailedAt:   failedAt.Unix(),
	}
}

func newCosPrivateRedisProgressStore(client cosPrivateRedisClient, env string) *cosPrivateRedisProgressStore {
	return &cosPrivateRedisProgressStore{client: client, env: env, now: time.Now}
}

func cosPrivateRedisKeyComponent(value string) string {
	const hexDigits = "0123456789ABCDEF"
	var encoded strings.Builder
	encoded.Grow(len(value))
	for index := 0; index < len(value); index++ {
		current := value[index]
		if current >= 'a' && current <= 'z' || current >= 'A' && current <= 'Z' || current >= '0' && current <= '9' || current == '_' || current == '-' {
			encoded.WriteByte(current)
			continue
		}
		encoded.WriteByte('%')
		encoded.WriteByte(hexDigits[current>>4])
		encoded.WriteByte(hexDigits[current&0x0f])
	}
	return encoded.String()
}

func cosPrivateRedisRowMember(rowID int64) string {
	return fmt.Sprintf("%019d", rowID)
}

func cosPrivateFailureDetailParts(failure cosPrivateFailure) (string, string) {
	valueHash, _ := json.Marshal(failure.ValueHash)
	failedURL, _ := json.Marshal(failure.FailedURL)
	reason, _ := json.Marshal(failure.Reason)
	prefix := fmt.Sprintf(
		`{"row_id":%d,"value_hash":%s,"failed_url":%s,"reason":%s,"retry_count":`,
		failure.RowID,
		valueHash,
		failedURL,
		reason,
	)
	suffix := fmt.Sprintf(`,"failed_at":%d}`, failure.FailedAt)
	return prefix, suffix
}

func (s *cosPrivateRedisProgressStore) Acquire(requestCtx ctx.Context, token string, ttl time.Duration) (bool, error) {
	if token == "" {
		return false, errCosPrivateInvalidToken
	}
	if err := validateCosPrivateLockTTL(ttl); err != nil {
		return false, err
	}
	acquired, err := s.client.SetNX(requestCtx, cosPrivateProgressKeys(s.env, cosPrivateField{}).Lock, token, ttl).Result()
	if err != nil || !acquired {
		return acquired, err
	}
	s.tokenM.Lock()
	s.token = token
	s.tokenM.Unlock()
	return true, nil
}

func (s *cosPrivateRedisProgressStore) Renew(requestCtx ctx.Context, token string, ttl time.Duration) error {
	if err := validateCosPrivateLockTTL(ttl); err != nil {
		return err
	}
	ownedToken, err := s.ownedToken(token)
	if err != nil {
		return err
	}
	result, err := s.client.Eval(
		requestCtx,
		cosPrivateRenewLockScript,
		[]string{cosPrivateProgressKeys(s.env, cosPrivateField{}).Lock},
		ownedToken,
		ttl.Milliseconds(),
	).Int64()
	if err != nil {
		return err
	}
	if result != 1 {
		s.clearToken(ownedToken)
		return errCosPrivateLockLost
	}
	return nil
}

func (s *cosPrivateRedisProgressStore) Release(requestCtx ctx.Context, token string) error {
	ownedToken, err := s.ownedToken(token)
	if err != nil {
		return err
	}
	result, err := s.client.Eval(
		requestCtx,
		cosPrivateReleaseLockScript,
		[]string{cosPrivateProgressKeys(s.env, cosPrivateField{}).Lock},
		ownedToken,
	).Int64()
	if err != nil {
		return err
	}
	s.clearToken(ownedToken)
	if result != 1 {
		return errCosPrivateLockLost
	}
	return nil
}

func (s *cosPrivateRedisProgressStore) Load(requestCtx ctx.Context, field cosPrivateField) (cosPrivateProgress, error) {
	values, err := s.client.HGetAll(requestCtx, cosPrivateProgressKeys(s.env, field).Progress).Result()
	if err != nil {
		return cosPrivateProgress{}, err
	}
	progress := cosPrivateProgress{Status: "pending"}
	for name, target := range map[string]*int64{
		"last_id":                   &progress.LastID,
		"scanned_count":             &progress.ScannedCount,
		"success_count":             &progress.SuccessCount,
		"skip_count":                &progress.SkipCount,
		"failure_total":             &progress.FailureTotal,
		"outstanding_failure_count": &progress.OutstandingFailureCount,
	} {
		value := values[name]
		if value == "" {
			continue
		}
		parsed, parseErr := strconv.ParseInt(value, 10, 64)
		if parseErr != nil {
			return cosPrivateProgress{}, fmt.Errorf("parse COS private progress %s=%q: %w", name, value, parseErr)
		}
		*target = parsed
	}
	if values["status"] != "" {
		progress.Status = values["status"]
	}
	return progress, nil
}

func (s *cosPrivateRedisProgressStore) RetryIDs(requestCtx ctx.Context, field cosPrivateField) ([]int64, error) {
	values, err := s.client.ZRange(requestCtx, cosPrivateProgressKeys(s.env, field).Failures, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	ids := make([]int64, 0, len(values))
	for _, value := range values {
		id, parseErr := strconv.ParseInt(value, 10, 64)
		if parseErr != nil {
			return nil, fmt.Errorf("parse COS private failure row ID %q: %w", value, parseErr)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (s *cosPrivateRedisProgressStore) RecordFailure(requestCtx ctx.Context, field cosPrivateField, failure cosPrivateFailure, advanceCursor bool) error {
	if failure.RowID <= 0 {
		return errors.New("COS private migration failure row ID must be positive")
	}
	if failure.RetryCount <= 0 {
		return errors.New("COS private migration retry count must be positive")
	}
	keys := cosPrivateProgressKeys(s.env, field)
	advance := "0"
	if advanceCursor {
		advance = "1"
	}
	detailPrefix, detailSuffix := cosPrivateFailureDetailParts(failure)
	return s.evalMutation(
		requestCtx,
		cosPrivateRecordFailureScript,
		[]string{keys.Progress, keys.Failures, keys.FailureDetails, keys.Task},
		strconv.FormatInt(failure.RowID, 10),
		cosPrivateRedisRowMember(failure.RowID),
		detailPrefix,
		detailSuffix,
		strconv.FormatInt(failure.RetryCount, 10),
		strconv.FormatInt(s.now().Unix(), 10),
		advance,
	)
}

func (s *cosPrivateRedisProgressStore) RecordSuccess(requestCtx ctx.Context, field cosPrivateField, rowID int64) error {
	return s.recordResult(requestCtx, field, rowID, cosPrivateRecordSuccessScript)
}

func (s *cosPrivateRedisProgressStore) RecordSkip(requestCtx ctx.Context, field cosPrivateField, rowID int64) error {
	return s.recordResult(requestCtx, field, rowID, cosPrivateRecordSkipScript)
}

func (s *cosPrivateRedisProgressStore) recordResult(requestCtx ctx.Context, field cosPrivateField, rowID int64, script string) error {
	if rowID <= 0 {
		return errors.New("COS private migration result row ID must be positive")
	}
	keys := cosPrivateProgressKeys(s.env, field)
	return s.evalMutation(
		requestCtx,
		script,
		[]string{keys.Progress, keys.Failures, keys.FailureDetails, keys.Task},
		strconv.FormatInt(rowID, 10),
		cosPrivateRedisRowMember(rowID),
		strconv.FormatInt(s.now().Unix(), 10),
	)
}

func (s *cosPrivateRedisProgressStore) ResolveFailure(requestCtx ctx.Context, field cosPrivateField, rowID int64) error {
	if rowID <= 0 {
		return errors.New("COS private migration failure row ID must be positive")
	}
	keys := cosPrivateProgressKeys(s.env, field)
	return s.evalMutation(
		requestCtx,
		cosPrivateResolveFailureScript,
		[]string{keys.Progress, keys.Failures, keys.FailureDetails, keys.Task},
		strconv.FormatInt(rowID, 10),
		cosPrivateRedisRowMember(rowID),
		strconv.FormatInt(s.now().Unix(), 10),
	)
}

func (s *cosPrivateRedisProgressStore) SetCurrent(requestCtx ctx.Context, field cosPrivateField, rowID int64, status string) error {
	if rowID < 0 {
		return errors.New("COS private migration current row ID cannot be negative")
	}
	keys := cosPrivateProgressKeys(s.env, field)
	return s.evalMutation(
		requestCtx,
		cosPrivateSetCurrentScript,
		[]string{keys.Progress, keys.Task},
		field.Database,
		field.Table,
		field.Field,
		strconv.FormatInt(rowID, 10),
		status,
		strconv.FormatInt(s.now().Unix(), 10),
	)
}

func (s *cosPrivateRedisProgressStore) Finish(requestCtx ctx.Context, status string) error {
	return s.evalMutation(
		requestCtx,
		cosPrivateFinishScript,
		[]string{cosPrivateProgressKeys(s.env, cosPrivateField{}).Task},
		status,
		strconv.FormatInt(s.now().Unix(), 10),
	)
}

func validateCosPrivateLockTTL(ttl time.Duration) error {
	if ttl <= 0 || ttl.Milliseconds() <= 0 {
		return errCosPrivateInvalidLockTTL
	}
	return nil
}

func (s *cosPrivateRedisProgressStore) ownedToken(provided string) (string, error) {
	s.tokenM.RLock()
	owned := s.token
	s.tokenM.RUnlock()
	if owned == "" || provided != owned {
		return "", errCosPrivateLockLost
	}
	return owned, nil
}

func (s *cosPrivateRedisProgressStore) currentToken() (string, error) {
	s.tokenM.RLock()
	owned := s.token
	s.tokenM.RUnlock()
	if owned == "" {
		return "", errCosPrivateLockLost
	}
	return owned, nil
}

func (s *cosPrivateRedisProgressStore) clearToken(token string) {
	s.tokenM.Lock()
	if s.token == token {
		s.token = ""
	}
	s.tokenM.Unlock()
}

func (s *cosPrivateRedisProgressStore) evalMutation(requestCtx ctx.Context, script string, keys []string, args ...interface{}) error {
	token, err := s.currentToken()
	if err != nil {
		return err
	}
	lockKey := cosPrivateProgressKeys(s.env, cosPrivateField{}).Lock
	mutationKeys := make([]string, 0, len(keys)+1)
	mutationKeys = append(mutationKeys, lockKey)
	mutationKeys = append(mutationKeys, keys...)
	mutationArgs := make([]interface{}, 0, len(args)+1)
	mutationArgs = append(mutationArgs, token)
	mutationArgs = append(mutationArgs, args...)
	result, err := s.client.Eval(requestCtx, script, mutationKeys, mutationArgs...).Int64()
	if err != nil {
		return err
	}
	if result != 1 {
		s.clearToken(token)
		return errCosPrivateLockLost
	}
	return nil
}

func (s *cosPrivateRedisProgressStore) OutstandingFailures(requestCtx ctx.Context) (int64, error) {
	value, err := s.client.HGet(
		requestCtx,
		cosPrivateProgressKeys(s.env, cosPrivateField{}).Task,
		"outstanding_failure_count",
	).Result()
	if errors.Is(err, redis.Nil) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	count, parseErr := strconv.ParseInt(value, 10, 64)
	if parseErr != nil {
		return 0, fmt.Errorf("parse COS private outstanding failure count %q: %w", value, parseErr)
	}
	return count, nil
}
