package membersJob

import (
	"errors"
	"testing"

	toolboxCtx "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/pkg/sdk/boss"
)

func TestCompleteMemberSyncDoesNotSyncStoreWhenSuccessRecordFails(t *testing.T) {
	recordErr := errors.New("record member sync success")
	syncCalled := false
	err := completeMemberSync(
		func() error { return recordErr },
		func() error {
			syncCalled = true
			return nil
		},
	)
	if !errors.Is(err, recordErr) {
		t.Fatalf("completeMemberSync() error = %v, want %v", err, recordErr)
	}
	if syncCalled {
		t.Fatal("store sync was called before member sync success was recorded")
	}
}

func TestSyncHkStoreAfterMemberSyncSkipsEnterpriseWithoutRootDepartmentConfig(t *testing.T) {
	syncCalled := false
	err := syncHkStoreAfterMemberSync(
		toolboxCtx.New(),
		123,
		func(_ toolboxCtx.Context, request interface{}) (*boss.EnterpriseConfigResponse, error) {
			assertHkStoreConfigRequest(t, request, 123)
			return &boss.EnterpriseConfigResponse{}, nil
		},
		func(_ toolboxCtx.Context, _ int64) error {
			syncCalled = true
			return nil
		},
	)
	if err != nil {
		t.Fatalf("syncHkStoreAfterMemberSync() error = %v", err)
	}
	if syncCalled {
		t.Fatal("store sync was called without hkmo_store_root_dept_id")
	}
}

func TestSyncHkStoreAfterMemberSyncCallsSyncForConfiguredEnterprise(t *testing.T) {
	var syncedEnterpriseID int64
	err := syncHkStoreAfterMemberSync(
		toolboxCtx.New(),
		123,
		func(_ toolboxCtx.Context, request interface{}) (*boss.EnterpriseConfigResponse, error) {
			assertHkStoreConfigRequest(t, request, 123)
			return &boss.EnterpriseConfigResponse{Data: []boss.EnterpriseConfigList{{
				Key:   "hkmo_store_root_dept_id",
				Value: 456,
			}}}, nil
		},
		func(_ toolboxCtx.Context, enterpriseID int64) error {
			syncedEnterpriseID = enterpriseID
			return nil
		},
	)
	if err != nil {
		t.Fatalf("syncHkStoreAfterMemberSync() error = %v", err)
	}
	if syncedEnterpriseID != 123 {
		t.Fatalf("synced enterprise ID = %d, want 123", syncedEnterpriseID)
	}
}

func TestSyncHkStoreAfterMemberSyncReturnsConfigAndSyncErrors(t *testing.T) {
	configErr := errors.New("get config")
	err := syncHkStoreAfterMemberSync(
		toolboxCtx.New(),
		123,
		func(_ toolboxCtx.Context, _ interface{}) (*boss.EnterpriseConfigResponse, error) {
			return nil, configErr
		},
		func(_ toolboxCtx.Context, _ int64) error {
			t.Fatal("store sync was called after config lookup failed")
			return nil
		},
	)
	if !errors.Is(err, configErr) {
		t.Fatalf("config error = %v, want %v", err, configErr)
	}

	syncErr := errors.New("sync store")
	err = syncHkStoreAfterMemberSync(
		toolboxCtx.New(),
		123,
		func(_ toolboxCtx.Context, _ interface{}) (*boss.EnterpriseConfigResponse, error) {
			return &boss.EnterpriseConfigResponse{Data: []boss.EnterpriseConfigList{{Value: 456}}}, nil
		},
		func(_ toolboxCtx.Context, _ int64) error { return syncErr },
	)
	if !errors.Is(err, syncErr) {
		t.Fatalf("sync error = %v, want %v", err, syncErr)
	}
}

func assertHkStoreConfigRequest(t *testing.T, request interface{}, enterpriseID int64) {
	t.Helper()
	params, ok := request.(map[string]interface{})
	if !ok {
		t.Fatalf("config request type = %T, want map[string]interface{}", request)
	}
	if params["enterprise_id"] != enterpriseID {
		t.Fatalf("config enterprise_id = %v, want %d", params["enterprise_id"], enterpriseID)
	}
	if params["key"] != "hkmo_store_root_dept_id" {
		t.Fatalf("config key = %v, want hkmo_store_root_dept_id", params["key"])
	}
}
