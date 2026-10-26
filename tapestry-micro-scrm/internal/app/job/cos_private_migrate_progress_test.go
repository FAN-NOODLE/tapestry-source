package appJob

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

type recordingCosPrivateRedis struct {
	mu              sync.Mutex
	setNXKey        string
	setNXValue      interface{}
	setNXExpiration time.Duration
	setNXResult     bool
	setNXErr        error
	setNXCalls      int
	evalScript      string
	evalKeys        []string
	evalArgs        []interface{}
	evalResult      interface{}
	evalErr         error
	evalCalls       int
	hgetAllResult   map[string]string
	hgetAllErr      error
	zrangeResult    []string
	zrangeErr       error
	hgetResult      string
	hgetErr         error
	hgetKey         string
	hgetField       string
}

func (r *recordingCosPrivateRedis) SetNX(_ context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.setNXCalls++
	r.setNXKey = key
	r.setNXValue = value
	r.setNXExpiration = expiration
	return redis.NewBoolResult(r.setNXResult, r.setNXErr)
}

func (r *recordingCosPrivateRedis) Eval(_ context.Context, script string, keys []string, args ...interface{}) *redis.Cmd {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.evalCalls++
	r.evalScript = script
	r.evalKeys = append([]string(nil), keys...)
	r.evalArgs = append([]interface{}(nil), args...)
	return redis.NewCmdResult(r.evalResult, r.evalErr)
}

func (r *recordingCosPrivateRedis) HGetAll(context.Context, string) *redis.StringStringMapCmd {
	return redis.NewStringStringMapResult(r.hgetAllResult, r.hgetAllErr)
}

func (r *recordingCosPrivateRedis) ZRange(context.Context, string, int64, int64) *redis.StringSliceCmd {
	return redis.NewStringSliceResult(r.zrangeResult, r.zrangeErr)
}

func (r *recordingCosPrivateRedis) HGet(_ context.Context, key string, field string) *redis.StringCmd {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.hgetKey = key
	r.hgetField = field
	return redis.NewStringResult(r.hgetResult, r.hgetErr)
}

func acquireCosPrivateProgressStore(t *testing.T, client *recordingCosPrivateRedis) *cosPrivateRedisProgressStore {
	t.Helper()
	client.setNXResult = true
	store := newCosPrivateRedisProgressStore(client, "uat")
	acquired, err := store.Acquire(ctx.New(), "run-token", time.Minute)
	if err != nil || !acquired {
		t.Fatalf("acquire=%v err=%v", acquired, err)
	}
	return store
}

func TestProgressKeysAreFieldScoped(t *testing.T) {
	keys := cosPrivateProgressKeys("uat", cosPrivateField{"micro_scrm", "download_resource", "source_data"})
	want := "cos:private:migrate:v4:uat:micro_scrm:download_resource:source_data:progress"
	if keys.Progress != want {
		t.Fatalf("got %q", keys.Progress)
	}
}

func TestProgressKeyComponentsCannotCollide(t *testing.T) {
	left := cosPrivateProgressKeys("uat:a", cosPrivateField{"b", "c", "d"})
	right := cosPrivateProgressKeys("uat", cosPrivateField{"a:b", "c", "d"})
	if left.Progress == right.Progress || left.Task == right.Task {
		t.Fatalf("colliding keys: left=%+v right=%+v", left, right)
	}
	percent := cosPrivateProgressKeys("uat", cosPrivateField{"a%3Ab", "c", "d"})
	if percent.Progress == right.Progress {
		t.Fatalf("encoded input collided: percent=%q colon=%q", percent.Progress, right.Progress)
	}
}

func TestFailureDetailHashesRawValue(t *testing.T) {
	f := newCosPrivateFailure(7, "secret raw value", "https://old/x", errors.New("copy failed"), time.Unix(10, 0))
	body, err := json.Marshal(f)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(body, []byte("secret raw value")) || f.ValueHash == "" {
		t.Fatalf("failure=%s", body)
	}
}

func TestFailureDetailRedactsRawValueFromReason(t *testing.T) {
	rawValue := `{"private":"secret value"}`
	failure := newCosPrivateFailure(7, rawValue, "https://old/x", fmt.Errorf("copy failed for %s", rawValue), time.Unix(10, 0))
	body, err := json.Marshal(failure)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(failure.Reason, rawValue) || bytes.Contains(body, []byte(rawValue)) {
		t.Fatalf("failure leaked raw value: %s", body)
	}
	if !strings.Contains(failure.Reason, "copy failed") {
		t.Fatalf("reason lost useful context: %q", failure.Reason)
	}
}

func TestLockScriptsRequireMatchingToken(t *testing.T) {
	if !strings.Contains(cosPrivateReleaseLockScript, "redis.call('get',KEYS[1]) == ARGV[1]") {
		t.Fatal("missing token check")
	}
	if !strings.Contains(cosPrivateRenewLockScript, "redis.call('get',KEYS[1]) == ARGV[1]") {
		t.Fatal("missing renew token check")
	}
}

func TestProgressAcquireUsesEnvironmentLockAndTTL(t *testing.T) {
	client := &recordingCosPrivateRedis{setNXResult: true}
	store := newCosPrivateRedisProgressStore(client, "uat")
	acquired, err := store.Acquire(ctx.New(), "run-token", time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if !acquired || client.setNXKey != "cos:private:migrate:v4:uat:lock" || client.setNXValue != "run-token" || client.setNXExpiration != time.Minute {
		t.Fatalf("acquired=%v key=%q value=%v ttl=%v", acquired, client.setNXKey, client.setNXValue, client.setNXExpiration)
	}
}

func TestProgressAcquireRejectsInvalidTTLWithoutRedisCall(t *testing.T) {
	for _, ttl := range []time.Duration{0, -time.Second, time.Nanosecond, 999 * time.Microsecond} {
		t.Run(ttl.String(), func(t *testing.T) {
			client := &recordingCosPrivateRedis{setNXResult: true}
			store := newCosPrivateRedisProgressStore(client, "uat")
			if acquired, err := store.Acquire(ctx.New(), "run-token", ttl); err == nil || acquired {
				t.Fatalf("ttl=%v acquired=%v err=%v", ttl, acquired, err)
			}
			if client.setNXCalls != 0 {
				t.Fatalf("ttl=%v redis calls=%d", ttl, client.setNXCalls)
			}
		})
	}
}

func TestProgressAcquireRejectsEmptyTokenWithoutRedisCall(t *testing.T) {
	client := &recordingCosPrivateRedis{setNXResult: true}
	store := newCosPrivateRedisProgressStore(client, "uat")
	if acquired, err := store.Acquire(ctx.New(), "", time.Minute); err == nil || acquired {
		t.Fatalf("acquired=%v err=%v", acquired, err)
	}
	if client.setNXCalls != 0 {
		t.Fatalf("redis calls=%d", client.setNXCalls)
	}
}

func TestProgressRenewRejectsInvalidTTLWithoutRedisCall(t *testing.T) {
	for _, ttl := range []time.Duration{0, -time.Second, time.Nanosecond, 999 * time.Microsecond} {
		t.Run(ttl.String(), func(t *testing.T) {
			client := &recordingCosPrivateRedis{setNXResult: true, evalResult: int64(1)}
			store := newCosPrivateRedisProgressStore(client, "uat")
			if acquired, err := store.Acquire(ctx.New(), "run-token", time.Minute); err != nil || !acquired {
				t.Fatalf("acquire=%v err=%v", acquired, err)
			}
			if err := store.Renew(ctx.New(), "run-token", ttl); err == nil {
				t.Fatalf("ttl=%v should fail", ttl)
			}
			if client.evalCalls != 0 {
				t.Fatalf("ttl=%v redis calls=%d", ttl, client.evalCalls)
			}
		})
	}
}

func TestProgressLockOwnershipLossReturnsError(t *testing.T) {
	for _, operation := range []struct {
		name string
		call func(*cosPrivateRedisProgressStore) error
	}{
		{"renew", func(store *cosPrivateRedisProgressStore) error {
			return store.Renew(ctx.New(), "run-token", time.Minute)
		}},
		{"release", func(store *cosPrivateRedisProgressStore) error { return store.Release(ctx.New(), "run-token") }},
	} {
		t.Run(operation.name, func(t *testing.T) {
			client := &recordingCosPrivateRedis{evalResult: int64(0)}
			store := acquireCosPrivateProgressStore(t, client)
			if err := operation.call(store); err == nil {
				t.Fatal("operation should fail after ownership loss")
			}
		})
	}
}

func TestProgressLoadParsesPersistentHash(t *testing.T) {
	client := &recordingCosPrivateRedis{hgetAllResult: map[string]string{
		"last_id":                   "17",
		"scanned_count":             "20",
		"success_count":             "11",
		"skip_count":                "5",
		"failure_total":             "4",
		"outstanding_failure_count": "2",
		"status":                    "running",
	}}
	store := newCosPrivateRedisProgressStore(client, "uat")
	got, err := store.Load(ctx.New(), cosPrivateField{"micro_scrm", "sample", "url"})
	if err != nil {
		t.Fatal(err)
	}
	want := cosPrivateProgress{17, 20, 11, 5, 4, 2, "running"}
	if got != want {
		t.Fatalf("got=%+v want=%+v", got, want)
	}
}

func TestProgressLoadParsesMaxInt64Cursor(t *testing.T) {
	client := &recordingCosPrivateRedis{hgetAllResult: map[string]string{"last_id": "9223372036854775807"}}
	store := newCosPrivateRedisProgressStore(client, "uat")
	got, err := store.Load(ctx.New(), cosPrivateField{"micro_scrm", "sample", "url"})
	if err != nil {
		t.Fatal(err)
	}
	if got.LastID != math.MaxInt64 {
		t.Fatalf("last_id=%d", got.LastID)
	}
}

func TestProgressRetryIDsPreserveSortedSetOrder(t *testing.T) {
	client := &recordingCosPrivateRedis{zrangeResult: []string{"2", "7", "19"}}
	store := newCosPrivateRedisProgressStore(client, "uat")
	got, err := store.RetryIDs(ctx.New(), cosPrivateField{"micro_scrm", "sample", "url"})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, []int64{2, 7, 19}) {
		t.Fatalf("ids=%v", got)
	}
}

func TestProgressRowMemberOrdersFullInt64(t *testing.T) {
	if got := cosPrivateRedisRowMember(2); got != "0000000000000000002" {
		t.Fatalf("small member=%q", got)
	}
	if got := cosPrivateRedisRowMember(math.MaxInt64); got != "9223372036854775807" {
		t.Fatalf("max member=%q", got)
	}
	if !(cosPrivateRedisRowMember(1<<53+1) < cosPrivateRedisRowMember(math.MaxInt64-1)) {
		t.Fatal("fixed-width members must preserve numeric order")
	}
}

func TestFailureDetailKeepsMaxInt64WithoutLuaJSONRoundTrip(t *testing.T) {
	failure := newCosPrivateFailure(math.MaxInt64, "raw", "https://old/x", errors.New("copy failed"), time.Unix(10, 0))
	prefix, suffix := cosPrivateFailureDetailParts(failure)
	body := prefix + "1" + suffix
	decoder := json.NewDecoder(strings.NewReader(body))
	decoder.UseNumber()
	var detail map[string]interface{}
	if err := decoder.Decode(&detail); err != nil {
		t.Fatal(err)
	}
	if got := detail["row_id"].(json.Number).String(); got != "9223372036854775807" {
		t.Fatalf("row_id=%q body=%s", got, body)
	}
}

func TestProgressRetryIDsParseFullInt64Members(t *testing.T) {
	client := &recordingCosPrivateRedis{zrangeResult: []string{
		"0009007199254740993",
		"9223372036854775806",
	}}
	store := newCosPrivateRedisProgressStore(client, "uat")
	got, err := store.RetryIDs(ctx.New(), cosPrivateField{"micro_scrm", "sample", "url"})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, []int64{1<<53 + 1, math.MaxInt64 - 1}) {
		t.Fatalf("ids=%v", got)
	}
}

func TestProgressFailureScriptAvoidsFloatingPointRowIDs(t *testing.T) {
	if strings.Contains(cosPrivateRecordFailureScript, "cjson") || strings.Contains(cosPrivateRecordFailureScript, "tonumber(row_id)") {
		t.Fatal("failure script must keep row IDs as opaque decimal strings")
	}
	if !strings.Contains(cosPrivateRecordFailureScript, "redis.call('zadd',failure_key,0,row_member)") {
		t.Fatal("failure IDs must use fixed-width members at a constant score")
	}
}

func TestProgressMutationsRequireAcquiredLock(t *testing.T) {
	client := &recordingCosPrivateRedis{evalResult: int64(1)}
	store := newCosPrivateRedisProgressStore(client, "uat")
	err := store.RecordSuccess(ctx.New(), cosPrivateField{"micro_scrm", "sample", "url"}, 7)
	if err == nil || client.evalCalls != 0 {
		t.Fatalf("err=%v eval_calls=%d", err, client.evalCalls)
	}
}

func TestProgressMutationsFenceExpiredOwner(t *testing.T) {
	field := cosPrivateField{"micro_scrm", "sample", "url"}
	failure := newCosPrivateFailure(7, "raw", "https://old/x", errors.New("copy failed"), time.Unix(10, 0))
	tests := []struct {
		name string
		call func(*cosPrivateRedisProgressStore) error
	}{
		{"failure", func(store *cosPrivateRedisProgressStore) error {
			return store.RecordFailure(ctx.New(), field, failure, true)
		}},
		{"success", func(store *cosPrivateRedisProgressStore) error { return store.RecordSuccess(ctx.New(), field, 7) }},
		{"skip", func(store *cosPrivateRedisProgressStore) error { return store.RecordSkip(ctx.New(), field, 7) }},
		{"resolve", func(store *cosPrivateRedisProgressStore) error { return store.ResolveFailure(ctx.New(), field, 7) }},
		{"current", func(store *cosPrivateRedisProgressStore) error {
			return store.SetCurrent(ctx.New(), field, 7, "running")
		}},
		{"finish", func(store *cosPrivateRedisProgressStore) error { return store.Finish(ctx.New(), "failed") }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := &recordingCosPrivateRedis{setNXResult: true, evalResult: int64(0)}
			store := newCosPrivateRedisProgressStore(client, "uat")
			acquired, err := store.Acquire(ctx.New(), "old-owner", time.Minute)
			if err != nil || !acquired {
				t.Fatalf("acquire=%v err=%v", acquired, err)
			}
			if err := test.call(store); err == nil {
				t.Fatal("stale owner mutation should fail")
			}
			if len(client.evalKeys) == 0 || client.evalKeys[0] != "cos:private:migrate:v4:uat:lock" {
				t.Fatalf("keys=%v", client.evalKeys)
			}
			if len(client.evalArgs) == 0 || client.evalArgs[0] != "old-owner" {
				t.Fatalf("args=%#v", client.evalArgs)
			}
		})
	}
}

func TestProgressMutationScriptsCheckLockBeforeWrites(t *testing.T) {
	scripts := []string{
		cosPrivateRecordFailureScript,
		cosPrivateRecordSuccessScript,
		cosPrivateRecordSkipScript,
		cosPrivateResolveFailureScript,
		cosPrivateSetCurrentScript,
		cosPrivateFinishScript,
	}
	for _, script := range scripts {
		lockCheck := strings.Index(script, "redis.call('get',lock_key) ~= lock_token")
		firstWrite := len(script)
		for _, command := range []string{"redis.call('hset'", "redis.call('hincrby'", "redis.call('zadd'", "redis.call('zrem'", "redis.call('hdel'"} {
			if index := strings.Index(script, command); index >= 0 && index < firstWrite {
				firstWrite = index
			}
		}
		if lockCheck < 0 || firstWrite == len(script) || lockCheck > firstWrite {
			t.Fatalf("lock check must precede writes: %q", script)
		}
	}
}

func TestProgressMutationScriptsPreflightBeforeWrites(t *testing.T) {
	tests := []struct {
		script       string
		preflight    string
		firstCommand string
	}{
		{cosPrivateRecordFailureScript, "local failure_total = checked_counter", "redis.call('zadd'"},
		{cosPrivateRecordSuccessScript, "local success_count = checked_counter", "redis.call('zrem'"},
		{cosPrivateRecordSkipScript, "local skip_count = checked_counter", "redis.call('zrem'"},
		{cosPrivateResolveFailureScript, "require_key_type(task_key,'hash')", "redis.call('zrem'"},
		{cosPrivateSetCurrentScript, "require_key_type(task_key,'hash')", "redis.call('hset'"},
		{cosPrivateFinishScript, "require_key_type(task_key,'hash')", "redis.call('hset'"},
	}
	for _, test := range tests {
		if strings.Contains(test.script, "cjson") || strings.Contains(test.script, "tonumber") {
			t.Fatal("mutation script must not parse exact integers through Lua numbers")
		}
		preflight := strings.Index(test.script, test.preflight)
		firstWrite := strings.Index(test.script, test.firstCommand)
		if preflight < 0 || firstWrite < 0 || preflight > firstWrite {
			t.Fatalf("preflight must precede writes: %q", test.script)
		}
	}
}

func TestProgressTokenStateSupportsConcurrentAccess(t *testing.T) {
	client := &recordingCosPrivateRedis{setNXResult: true, evalResult: int64(1)}
	store := newCosPrivateRedisProgressStore(client, "uat")
	field := cosPrivateField{"micro_scrm", "sample", "url"}
	var wait sync.WaitGroup
	for index := 1; index <= 100; index++ {
		token := fmt.Sprintf("run-token-%d", index)
		wait.Add(2)
		go func() {
			defer wait.Done()
			_, _ = store.Acquire(ctx.New(), token, time.Minute)
		}()
		go func(rowID int64) {
			defer wait.Done()
			_ = store.RecordSuccess(ctx.New(), field, rowID)
		}(int64(index))
	}
	wait.Wait()
}

func TestProgressRecordFailureUsesAtomicFieldTransition(t *testing.T) {
	client := &recordingCosPrivateRedis{evalResult: int64(1)}
	store := acquireCosPrivateProgressStore(t, client)
	field := cosPrivateField{"micro_scrm", "sample", "url"}
	failure := newCosPrivateFailure(7, "secret raw value", "https://old/x", errors.New("copy failed"), time.Unix(10, 0))
	if err := store.RecordFailure(ctx.New(), field, failure, true); err != nil {
		t.Fatal(err)
	}
	wantKeys := []string{
		"cos:private:migrate:v4:uat:lock",
		"cos:private:migrate:v4:uat:micro_scrm:sample:url:progress",
		"cos:private:migrate:v4:uat:micro_scrm:sample:url:failures",
		"cos:private:migrate:v4:uat:micro_scrm:sample:url:failure_details",
		"cos:private:migrate:v4:uat:task",
	}
	if client.evalScript != cosPrivateRecordFailureScript || !reflect.DeepEqual(client.evalKeys, wantKeys) {
		t.Fatalf("script=%q keys=%v", client.evalScript, client.evalKeys)
	}
	prefix, prefixOK := client.evalArgs[3].(string)
	suffix, suffixOK := client.evalArgs[4].(string)
	body := prefix + "1" + suffix
	if !prefixOK || !suffixOK || strings.Contains(body, "secret raw value") || !strings.Contains(body, `"row_id":7`) {
		t.Fatalf("failure detail leaked raw value: %#v", client.evalArgs)
	}
}

func TestProgressRecordFailurePassesMaxInt64AsStrings(t *testing.T) {
	client := &recordingCosPrivateRedis{evalResult: int64(1)}
	store := acquireCosPrivateProgressStore(t, client)
	failure := newCosPrivateFailure(math.MaxInt64, "raw", "https://old/x", errors.New("copy failed"), time.Unix(10, 0))
	if err := store.RecordFailure(ctx.New(), cosPrivateField{"micro_scrm", "sample", "url"}, failure, true); err != nil {
		t.Fatal(err)
	}
	if client.evalArgs[1] != "9223372036854775807" || client.evalArgs[2] != "9223372036854775807" {
		t.Fatalf("args=%#v", client.evalArgs)
	}
	prefix := client.evalArgs[3].(string)
	suffix := client.evalArgs[4].(string)
	if !strings.Contains(prefix+"1"+suffix, `"row_id":9223372036854775807`) {
		t.Fatalf("detail=%s", prefix+"1"+suffix)
	}
}

func TestProgressRecordSuccessUsesAtomicResolveTransition(t *testing.T) {
	client := &recordingCosPrivateRedis{evalResult: int64(1)}
	store := acquireCosPrivateProgressStore(t, client)
	field := cosPrivateField{"micro_scrm", "sample", "url"}
	if err := store.RecordSuccess(ctx.New(), field, 7); err != nil {
		t.Fatal(err)
	}
	keys := cosPrivateProgressKeys("uat", field)
	wantKeys := []string{keys.Lock, keys.Progress, keys.Failures, keys.FailureDetails, keys.Task}
	if client.evalScript != cosPrivateRecordSuccessScript || !reflect.DeepEqual(client.evalKeys, wantKeys) {
		t.Fatalf("script=%q keys=%v", client.evalScript, client.evalKeys)
	}
}

func TestProgressRecordSkipUsesAtomicResolveTransition(t *testing.T) {
	client := &recordingCosPrivateRedis{evalResult: int64(1)}
	store := acquireCosPrivateProgressStore(t, client)
	field := cosPrivateField{"micro_scrm", "sample", "url"}
	if err := store.RecordSkip(ctx.New(), field, 8); err != nil {
		t.Fatal(err)
	}
	if client.evalScript != cosPrivateRecordSkipScript {
		t.Fatalf("script=%q", client.evalScript)
	}
}

func TestProgressResolveFailureUsesAtomicCounts(t *testing.T) {
	client := &recordingCosPrivateRedis{evalResult: int64(1)}
	store := acquireCosPrivateProgressStore(t, client)
	field := cosPrivateField{"micro_scrm", "sample", "url"}
	if err := store.ResolveFailure(ctx.New(), field, 9); err != nil {
		t.Fatal(err)
	}
	keys := cosPrivateProgressKeys("uat", field)
	wantKeys := []string{keys.Lock, keys.Progress, keys.Failures, keys.FailureDetails, keys.Task}
	if client.evalScript != cosPrivateResolveFailureScript || !reflect.DeepEqual(client.evalKeys, wantKeys) {
		t.Fatalf("script=%q keys=%v", client.evalScript, client.evalKeys)
	}
}

func TestProgressSetCurrentUpdatesFieldAndTask(t *testing.T) {
	client := &recordingCosPrivateRedis{evalResult: int64(1)}
	store := acquireCosPrivateProgressStore(t, client)
	field := cosPrivateField{"micro_scrm", "sample", "url"}
	if err := store.SetCurrent(ctx.New(), field, 12, "running"); err != nil {
		t.Fatal(err)
	}
	keys := cosPrivateProgressKeys("uat", field)
	if client.evalScript != cosPrivateSetCurrentScript || !reflect.DeepEqual(client.evalKeys, []string{keys.Lock, keys.Progress, keys.Task}) {
		t.Fatalf("script=%q keys=%v", client.evalScript, client.evalKeys)
	}
	wantArgs := []interface{}{"run-token", "micro_scrm", "sample", "url", "12", "running"}
	if len(client.evalArgs) < len(wantArgs) || !reflect.DeepEqual(client.evalArgs[:len(wantArgs)], wantArgs) {
		t.Fatalf("args=%#v", client.evalArgs)
	}
}

func TestProgressFinishUpdatesTaskStatus(t *testing.T) {
	client := &recordingCosPrivateRedis{evalResult: int64(1)}
	store := acquireCosPrivateProgressStore(t, client)
	if err := store.Finish(ctx.New(), "partial_failed"); err != nil {
		t.Fatal(err)
	}
	if client.evalScript != cosPrivateFinishScript || !reflect.DeepEqual(client.evalKeys, []string{"cos:private:migrate:v4:uat:lock", "cos:private:migrate:v4:uat:task"}) {
		t.Fatalf("script=%q keys=%v", client.evalScript, client.evalKeys)
	}
}

func TestProgressOutstandingFailuresReadsTaskCounter(t *testing.T) {
	client := &recordingCosPrivateRedis{hgetResult: "3"}
	store := newCosPrivateRedisProgressStore(client, "uat")
	got, err := store.OutstandingFailures(ctx.New())
	if err != nil {
		t.Fatal(err)
	}
	if got != 3 || client.hgetKey != "cos:private:migrate:v4:uat:task" || client.hgetField != "outstanding_failure_count" {
		t.Fatalf("got=%d key=%q field=%q", got, client.hgetKey, client.hgetField)
	}
}

func TestProgressRedisWriteErrorStopsTransition(t *testing.T) {
	want := errors.New("redis unavailable")
	client := &recordingCosPrivateRedis{evalErr: want}
	store := acquireCosPrivateProgressStore(t, client)
	err := store.RecordSuccess(ctx.New(), cosPrivateField{"micro_scrm", "sample", "url"}, 7)
	if !errors.Is(err, want) {
		t.Fatalf("err=%v", err)
	}
}

func TestProgressStoreImplementsContract(t *testing.T) {
	var _ cosPrivateRedisClient = (*redis.Client)(nil)
	var _ cosPrivateProgressStore = newCosPrivateRedisProgressStore(&recordingCosPrivateRedis{}, "uat")
}
