package appJob

import (
	"errors"
	"net/http"
	"testing"

	"github.com/tencentyun/cos-go-sdk-v5"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/config"
)

type fakeCosPrivateTencentClient struct {
	headSize     int64
	headErr      error
	copyErr      error
	copyCalls    int
	copiedKey    string
	copiedSource string
}

func (c *fakeCosPrivateTencentClient) Head(_ ctx.Context, _ string) (int64, error) {
	return c.headSize, c.headErr
}

func (c *fakeCosPrivateTencentClient) Copy(_ ctx.Context, key, source string) error {
	c.copyCalls++
	c.copiedKey = key
	c.copiedSource = source
	return c.copyErr
}

type fakeCosStore struct {
	objects      map[string]int64
	copiedSize   int64
	copyCalls    int
	failKeys     map[string]error
	statErrors   map[string]error
	statCalls    map[string]int
	statErrorsAt map[string]map[int]error
}

func (s *fakeCosStore) Stat(_ ctx.Context, bucket, key string) (cosPrivateObjectMeta, bool, error) {
	objectName := bucket + "/" + key
	if s.statCalls == nil {
		s.statCalls = make(map[string]int)
	}
	s.statCalls[objectName]++
	if err := s.statErrors[objectName]; err != nil {
		return cosPrivateObjectMeta{}, false, err
	}
	if err := s.statErrorsAt[objectName][s.statCalls[objectName]]; err != nil {
		return cosPrivateObjectMeta{}, false, err
	}
	size, ok := s.objects[objectName]
	return cosPrivateObjectMeta{Size: size}, ok, nil
}

func (s *fakeCosStore) Copy(_ ctx.Context, _, target, sourceKey, targetKey string) error {
	if err := s.failKeys[sourceKey]; err != nil {
		return err
	}
	s.copyCalls++
	s.objects[target+"/"+targetKey] = s.copiedSize
	return nil
}

func TestEnsureCopiedSkipsMatchingTarget(t *testing.T) {
	store := &fakeCosStore{objects: map[string]int64{"old/a": 10, "private/a": 10}}
	err := (&cosPrivateCopier{store: store, sourceBucket: "old", targetBucket: "private"}).ensureCopied(ctx.New(), "a")
	if err != nil || store.copyCalls != 0 {
		t.Fatalf("err=%v calls=%d", err, store.copyCalls)
	}
}

func TestEnsureCopiedCopiesMissingTarget(t *testing.T) {
	store := &fakeCosStore{objects: map[string]int64{"old/a": 10}, copiedSize: 10}
	err := (&cosPrivateCopier{store: store, sourceBucket: "old", targetBucket: "private"}).ensureCopied(ctx.New(), "a")
	if err != nil || store.copyCalls != 1 {
		t.Fatalf("err=%v calls=%d", err, store.copyCalls)
	}
}

func TestEnsureCopiedFailsForMissingSource(t *testing.T) {
	store := &fakeCosStore{objects: map[string]int64{}}
	err := (&cosPrivateCopier{store: store, sourceBucket: "old", targetBucket: "private"}).ensureCopied(ctx.New(), "a")
	if err == nil || store.copyCalls != 0 {
		t.Fatalf("err=%v calls=%d", err, store.copyCalls)
	}
}

func TestEnsureCopiedClassifiesSource404AsRowError(t *testing.T) {
	notFound := &cos.ErrorResponse{Response: &http.Response{StatusCode: http.StatusNotFound}}
	store := cosPrivateTencentStore{
		sourceBucket: "old",
		targetBucket: "private",
		source:       &fakeCosPrivateTencentClient{headErr: notFound},
		target:       &fakeCosPrivateTencentClient{},
	}
	err := (&cosPrivateCopier{store: store, sourceBucket: "old", targetBucket: "private"}).ensureCopied(ctx.New(), "a")
	var rowErr *cosPrivateObjectRowError
	var sdkErr *cos.ErrorResponse
	if !isCosPrivateObjectRowError(err) || !errors.As(err, &rowErr) || !errors.As(err, &sdkErr) || !errors.Is(err, errCosPrivateSourceObjectNotFound) {
		t.Fatalf("err=%v", err)
	}
}

func TestEnsureCopiedLeavesInfrastructureFailuresNonRow(t *testing.T) {
	authErr := errors.New("authorization denied")
	networkErr := errors.New("network unavailable")
	serverErr := &cos.ErrorResponse{Response: &http.Response{StatusCode: http.StatusInternalServerError}}
	targetHeadErr := errors.New("target head failed")
	copyErr := errors.New("copy failed")
	verifyHeadErr := errors.New("target verify failed")

	tests := []struct {
		name  string
		store *fakeCosStore
		want  error
	}{
		{
			name:  "认证错误",
			store: &fakeCosStore{objects: map[string]int64{"old/a": 10}, statErrors: map[string]error{"old/a": authErr}},
			want:  authErr,
		},
		{
			name:  "网络错误",
			store: &fakeCosStore{objects: map[string]int64{"old/a": 10}, statErrors: map[string]error{"old/a": networkErr}},
			want:  networkErr,
		},
		{
			name:  "服务端错误",
			store: &fakeCosStore{objects: map[string]int64{"old/a": 10}, statErrors: map[string]error{"old/a": serverErr}},
			want:  serverErr,
		},
		{
			name:  "目标HEAD错误",
			store: &fakeCosStore{objects: map[string]int64{"old/a": 10}, statErrorsAt: map[string]map[int]error{"private/a": {1: targetHeadErr}}},
			want:  targetHeadErr,
		},
		{
			name:  "复制错误",
			store: &fakeCosStore{objects: map[string]int64{"old/a": 10}, failKeys: map[string]error{"a": copyErr}},
			want:  copyErr,
		},
		{
			name:  "复制后目标HEAD错误",
			store: &fakeCosStore{objects: map[string]int64{"old/a": 10}, copiedSize: 10, statErrorsAt: map[string]map[int]error{"private/a": {2: verifyHeadErr}}},
			want:  verifyHeadErr,
		},
		{
			name:  "最终大小校验错误",
			store: &fakeCosStore{objects: map[string]int64{"old/a": 10}, copiedSize: 9},
			want:  errCosPrivateTargetObjectVerification,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := (&cosPrivateCopier{store: tt.store, sourceBucket: "old", targetBucket: "private"}).ensureCopied(ctx.New(), "a")
			var rowErr *cosPrivateObjectRowError
			if err == nil || isCosPrivateObjectRowError(err) || errors.As(err, &rowErr) || !errors.Is(err, tt.want) {
				t.Fatalf("err=%v", err)
			}
		})
	}
}

func TestEnsureCopiedFailsForPostCopySizeMismatch(t *testing.T) {
	store := &fakeCosStore{objects: map[string]int64{"old/a": 10}, copiedSize: 9}
	err := (&cosPrivateCopier{store: store, sourceBucket: "old", targetBucket: "private"}).ensureCopied(ctx.New(), "a")
	if err == nil {
		t.Fatal("expected size mismatch")
	}
}

func TestEnsureCopiedReplacesMismatchedTarget(t *testing.T) {
	store := &fakeCosStore{objects: map[string]int64{"old/a": 10, "private/a": 9}, copiedSize: 10}
	err := (&cosPrivateCopier{store: store, sourceBucket: "old", targetBucket: "private"}).ensureCopied(ctx.New(), "a")
	if err != nil || store.copyCalls != 1 || store.objects["private/a"] != 10 {
		t.Fatalf("err=%v calls=%d size=%d", err, store.copyCalls, store.objects["private/a"])
	}
}

func TestEnsureCopiedPropagatesStatError(t *testing.T) {
	want := errors.New("stat failed")
	store := &fakeCosStore{
		objects:    map[string]int64{"old/a": 10},
		statErrors: map[string]error{"old/a": want},
	}
	err := (&cosPrivateCopier{store: store, sourceBucket: "old", targetBucket: "private"}).ensureCopied(ctx.New(), "a")
	if !errors.Is(err, want) || store.copyCalls != 0 {
		t.Fatalf("err=%v calls=%d", err, store.copyCalls)
	}
}

func TestEnsureCopiedFailsForPreCopyTargetStatError(t *testing.T) {
	want := errors.New("target stat failed")
	store := &fakeCosStore{
		objects: map[string]int64{"old/a": 10},
		statErrorsAt: map[string]map[int]error{
			"private/a": {1: want},
		},
	}
	err := (&cosPrivateCopier{store: store, sourceBucket: "old", targetBucket: "private"}).ensureCopied(ctx.New(), "a")
	if !errors.Is(err, want) || store.copyCalls != 0 {
		t.Fatalf("err=%v calls=%d", err, store.copyCalls)
	}
}

func TestEnsureCopiedFailsForPostCopyTargetStatError(t *testing.T) {
	want := errors.New("target verify failed")
	store := &fakeCosStore{
		objects:    map[string]int64{"old/a": 10},
		copiedSize: 10,
		statErrorsAt: map[string]map[int]error{
			"private/a": {2: want},
		},
	}
	err := (&cosPrivateCopier{store: store, sourceBucket: "old", targetBucket: "private"}).ensureCopied(ctx.New(), "a")
	if !errors.Is(err, want) || store.copyCalls != 1 {
		t.Fatalf("err=%v calls=%d", err, store.copyCalls)
	}
}

func TestEnsureCopiedPropagatesCopyError(t *testing.T) {
	want := errors.New("copy failed")
	store := &fakeCosStore{
		objects:  map[string]int64{"old/a": 10},
		failKeys: map[string]error{"a": want},
	}
	err := (&cosPrivateCopier{store: store, sourceBucket: "old", targetBucket: "private"}).ensureCopied(ctx.New(), "a")
	if !errors.Is(err, want) {
		t.Fatalf("err=%v", err)
	}
}

func TestCosPrivateTencentStoreStatReturnsExistingObject(t *testing.T) {
	source := &fakeCosPrivateTencentClient{headSize: 10}
	store := cosPrivateTencentStore{sourceBucket: "old-123", targetBucket: "private-123", source: source, target: &fakeCosPrivateTencentClient{}}
	meta, exists, err := store.Stat(ctx.New(), "old-123", "a")
	if err != nil || !exists || meta.Size != 10 {
		t.Fatalf("meta=%#v exists=%t err=%v", meta, exists, err)
	}
}

func TestCosPrivateTencentStoreStatUsesTargetClient(t *testing.T) {
	store := cosPrivateTencentStore{
		sourceBucket: "old-123",
		targetBucket: "private-123",
		source:       &fakeCosPrivateTencentClient{headSize: 10},
		target:       &fakeCosPrivateTencentClient{headSize: 20},
	}
	meta, exists, err := store.Stat(ctx.New(), "private-123", "a")
	if err != nil || !exists || meta.Size != 20 {
		t.Fatalf("meta=%#v exists=%t err=%v", meta, exists, err)
	}
}

func TestCosPrivateTencentStoreStatTreatsOnlyNotFoundAsMissing(t *testing.T) {
	notFound := &cos.ErrorResponse{Response: &http.Response{StatusCode: http.StatusNotFound}}
	store := cosPrivateTencentStore{
		sourceBucket: "old-123",
		targetBucket: "private-123",
		source:       &fakeCosPrivateTencentClient{headErr: notFound},
		target:       &fakeCosPrivateTencentClient{},
	}
	_, exists, err := store.Stat(ctx.New(), "old-123", "a")
	if err != nil || exists {
		t.Fatalf("exists=%t err=%v", exists, err)
	}
}

func TestCosPrivateTencentStoreStatPropagatesSDKError(t *testing.T) {
	want := errors.New("access denied")
	store := cosPrivateTencentStore{
		sourceBucket: "old-123",
		targetBucket: "private-123",
		source:       &fakeCosPrivateTencentClient{headErr: want},
		target:       &fakeCosPrivateTencentClient{},
	}
	_, exists, err := store.Stat(ctx.New(), "old-123", "a")
	if !errors.Is(err, want) || exists {
		t.Fatalf("exists=%t err=%v", exists, err)
	}
}

func TestCosPrivateTencentStoreCopyUsesTargetClientAndSDKCopySource(t *testing.T) {
	oldOrigin := config.CfgData.Storage.COS.Origin
	config.CfgData.Storage.COS.Origin = "ap-shanghai"
	t.Cleanup(func() { config.CfgData.Storage.COS.Origin = oldOrigin })

	target := &fakeCosPrivateTencentClient{}
	store := cosPrivateTencentStore{
		sourceBucket: "old-123",
		targetBucket: "private-123",
		source:       &fakeCosPrivateTencentClient{},
		target:       target,
	}
	key := "folder/a b+%.txt"
	err := store.Copy(ctx.New(), "old-123", "private-123", key, key)
	if err != nil || target.copyCalls != 1 || target.copiedKey != key {
		t.Fatalf("err=%v calls=%d key=%q", err, target.copyCalls, target.copiedKey)
	}
	wantSource := "old-123.cos.ap-shanghai.myqcloud.com/folder/a b+%.txt"
	if target.copiedSource != wantSource {
		t.Fatalf("copy source=%q want=%q", target.copiedSource, wantSource)
	}
}

func TestCosPrivateTencentStoreCopyPropagatesSDKError(t *testing.T) {
	want := errors.New("copy denied")
	store := cosPrivateTencentStore{
		sourceBucket: "old-123",
		targetBucket: "private-123",
		source:       &fakeCosPrivateTencentClient{},
		target:       &fakeCosPrivateTencentClient{copyErr: want},
	}
	err := store.Copy(ctx.New(), "old-123", "private-123", "a", "a")
	if !errors.Is(err, want) {
		t.Fatalf("err=%v", err)
	}
}
