package storeOperationJob

import (
	"errors"
	"testing"
	"time"

	toolboxCtx "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	appModels "micro-scrm/internal/app/models"
	dashboardModel "micro-scrm/internal/app/models/dashboard"
	"micro-scrm/internal/orm/model/micro_scrm"
	enterpriseModel "micro-scrm/internal/orm/model/micro_scrm_enterprise"
)

func TestBuildDashboardAreaNodesBuildsWhitelistedDistrictFourLevelTree(t *testing.T) {
	now := time.Date(2026, 9, 2, 10, 0, 0, 0, time.Local)
	stores := []enterpriseModel.StoreMasterModel{
		newAreaTestStore("OC268", "OPEN", "", "RTL", "Greater China Full-Price", "CAFP_CN_ED", "East", "E1", "East 1"),
		newAreaTestStore("OC267", "PENDING", "08/15/2022", "RTL", "Greater China Full-Price", "CAFP_CN_ED", "East", "E1", "East 1"),
	}
	stores[0].SalesChannel = "Retail"
	stores[1].SalesChannel = "Retail"

	nodes, skipped, err := buildDashboardAreaNodes(stores, map[string]int{"OC267": 287}, now)
	if err != nil {
		t.Fatal(err)
	}

	if len(skipped) != 0 {
		t.Fatalf("skipped = %v, want empty", skipped)
	}
	if len(nodes) != 5 {
		t.Fatalf("node count = %d, want 5", len(nodes))
	}
	assertAreaNode(t, nodes[0], 1, "RTL", "Retail", "", 0, 0, 1)
	assertAreaNode(t, nodes[1], 2, "CAFP_CN_ED", "East", "RTL", 0, 0, 1)
	assertAreaNode(t, nodes[2], 3, "E1", "East 1", "CAFP_CN_ED", 0, 0, 1)
	assertAreaNode(t, nodes[3], 4, "OC267", "OC267", "E1", 287, 0, 1)
	assertAreaNode(t, nodes[4], 4, "OC268", "OC268", "E1", 0, 0, 2)
}

func TestBuildDashboardAreaNodesSkipsInvalidCodesAndFallsBackToCodeForEmptyDescriptions(t *testing.T) {
	now := time.Date(2026, 9, 2, 10, 0, 0, 0, time.Local)
	valid := newAreaTestStore("OC267", "OPEN", "", "RTL", "N/A", "CAFP_CN_ED", "", "E1", " ")
	invalid := newAreaTestStore("OC999", "OPEN", "", "N/A", "Retail", "CAFP_CN_ED", "East", "E1", "East 1")

	nodes, skipped, err := buildDashboardAreaNodes([]enterpriseModel.StoreMasterModel{invalid, valid}, nil, now)
	if err != nil {
		t.Fatal(err)
	}

	if len(skipped) != 1 || skipped[0] != "OC999" {
		t.Fatalf("skipped = %v, want [OC999]", skipped)
	}
	if len(nodes) != 4 {
		t.Fatalf("node count = %d, want 4", len(nodes))
	}
	if nodes[0].Name != "RTL" || nodes[1].Name != "CAFP_CN_ED" || nodes[2].Name != "E1" {
		t.Fatalf("fallback names = %q, %q, %q", nodes[0].Name, nodes[1].Name, nodes[2].Name)
	}
}

func TestDashboardAreaHierarchyUsesDistrictCodeWhitelistInsteadOfRegionDescription(t *testing.T) {
	now := time.Date(2026, 9, 24, 10, 0, 0, 0, time.Local)
	allowed := newAreaTestStore("OC267", "OPEN", "", "RTL", "Retail", "CAFP_CN_ED", "Retail-East", "EA1", "East 1")
	allowed.RegionDesc = "Hong Kong and Macao Region"
	rejected := newAreaTestStore("OC999", "OPEN", "", "RTL", "Retail", "CAFP_CN_OTHER", "Other", "OA1", "Other 1")
	rejected.RegionDesc = "China Region"

	hierarchyStores, _, _ := partitionDashboardAreaStores([]enterpriseModel.StoreMasterModel{allowed, rejected}, now)
	if len(hierarchyStores) != 1 || hierarchyStores[0].StoreCode != "OC267" {
		t.Fatalf("hierarchy stores = %+v, want only district CAFP_CN_ED", hierarchyStores)
	}

	nodes, skipped, err := buildDashboardAreaNodes([]enterpriseModel.StoreMasterModel{allowed, rejected}, nil, now)
	if err != nil {
		t.Fatal(err)
	}
	if len(skipped) != 0 {
		t.Fatalf("skipped = %v, want empty", skipped)
	}
	if len(nodes) != 4 || nodes[1].LMSID != "CAFP_CN_ED" || nodes[3].LMSID != "OC267" {
		t.Fatalf("nodes = %+v, want only the allowed district hierarchy", nodes)
	}
}

func TestDashboardAreaHierarchyAcceptsAllConfiguredDistrictCodes(t *testing.T) {
	now := time.Date(2026, 9, 24, 10, 0, 0, 0, time.Local)
	districtCodes := []string{
		"CAFP_CN_CD",
		"CAFP_CN_ED",
		"CAFP_CN_HD",
		"CAFP_CN_ND",
		"CAFP_CN_SWD",
		"CAFS_CN_CD",
		"CAFS_CN_ESD",
		"CAFS_CN_HD",
		"CAFS_CN_ND",
		"CAFS_CN_WD",
		"CAFP_CN_CSD",
	}
	stores := make([]enterpriseModel.StoreMasterModel, 0, len(districtCodes))
	for index, districtCode := range districtCodes {
		store := newAreaTestStore(
			"OC"+districtCode,
			"OPEN",
			"",
			"RTL",
			"Retail",
			districtCode,
			districtCode,
			"AREA"+districtCode,
			"Area "+districtCode,
		)
		store.RegionDesc = ""
		if index == 0 {
			store.OrgDistrictCode = " cafp_cn_cd "
		}
		stores = append(stores, store)
	}

	hierarchyStores, _, _ := partitionDashboardAreaStores(stores, now)
	if len(hierarchyStores) != len(districtCodes) {
		t.Fatalf("hierarchy store count = %d, want %d", len(hierarchyStores), len(districtCodes))
	}
}

func TestDashboardAreaClosedAtUsesRunTimeWhenCloseDateCannotBeParsed(t *testing.T) {
	now := time.Date(2026, 9, 2, 10, 0, 0, 0, time.Local)
	if got := dashboardAreaClosedAt("OPEN", "08/15/2022", now); got != 0 {
		t.Fatalf("OPEN closed_at = %d, want 0", got)
	}
	if got := dashboardAreaClosedAt("PENDING", "N/A", now); got != 0 {
		t.Fatalf("PENDING closed_at = %d, want 0", got)
	}
	if got := dashboardAreaClosedAt("CLOSED", "2099/12/31", now); got != now.Unix() {
		t.Fatalf("out-of-range CLOSED closed_at = %d, want %d", got, now.Unix())
	}
}

func TestPartitionDashboardAreaStoresUsesDistrictWhitelistAndClosedHierarchy(t *testing.T) {
	now := time.Date(2026, 9, 8, 10, 0, 0, 0, time.Local)
	stores := []enterpriseModel.StoreMasterModel{
		newAreaTestStore("OC267", "OPEN", "", "RTL", "Retail", "CAFP_CN_ED", "East", "E1", "East 1"),
		newAreaTestStore("OC268", "PENDING", "", "RTL", "Retail", "CAFP_CN_ED", "East", "E1", "East 1"),
		newAreaTestStore("OC269", "CLOSED", "08/15/2022", "N/A", "N/A", "N/A", "N/A", "E1", "N/A"),
		newAreaTestStore("OC270", "Location Status N/A", "", "N/A", "N/A", "N/A", "N/A", "N/A", "N/A"),
		newAreaTestStore("HK001", "OPEN", "", "RTL", "Retail", "HK", "Hong Kong", "HK1", "Hong Kong 1"),
		newAreaTestStore("OC271", "CLOSED", "08/16/2022", "RTL", "Retail", "CAFP_CN_ED", "East", "E1", "East 1"),
		newAreaTestStore("COFFEE01", "OPEN", "", "COFFEE", "Coffee Shop", "CAFS_CN_CD", "Coffee District", "COFFEE-A", "Coffee Area"),
	}
	stores[2].RegionDesc = ""
	stores[3].RegionDesc = "Hong Kong and Macao Region"
	stores[4].RegionDesc = "Hong Kong and Macao Region"
	stores[6].RegionDesc = "China Retail Region - Coffee Shop"

	hierarchyStores, closedAtByStore, protectedStoreCodes := partitionDashboardAreaStores(stores, now)

	if len(hierarchyStores) != 4 || hierarchyStores[0].StoreCode != "OC267" || hierarchyStores[1].StoreCode != "OC268" ||
		hierarchyStores[2].StoreCode != "OC271" || hierarchyStores[3].StoreCode != "COFFEE01" {
		t.Fatalf("hierarchy stores = %+v, want supported OPEN, PENDING and CLOSED stores", hierarchyStores)
	}
	wantLegacyClosedAt := time.Date(2022, 8, 15, 0, 0, 0, 0, time.Local).Unix()
	wantHierarchyClosedAt := time.Date(2022, 8, 16, 0, 0, 0, 0, time.Local).Unix()
	if len(closedAtByStore) != 2 || closedAtByStore["OC269"] != wantLegacyClosedAt || closedAtByStore["OC271"] != wantHierarchyClosedAt {
		t.Fatalf("closed stores = %+v", closedAtByStore)
	}
	if len(protectedStoreCodes) != 3 {
		t.Fatalf("protected stores = %+v, want CLOSED and Location Status N/A stores", protectedStoreCodes)
	}
	if _, exists := protectedStoreCodes["OC269"]; !exists {
		t.Fatal("CLOSED store must be protected from hierarchy invalidation")
	}
	if _, exists := protectedStoreCodes["OC270"]; !exists {
		t.Fatal("Location Status N/A store must be protected from hierarchy invalidation")
	}
	if _, exists := protectedStoreCodes["OC271"]; !exists {
		t.Fatal("CLOSED hierarchy store must be protected from hierarchy invalidation")
	}

	nodes, skipped, err := buildDashboardAreaNodes(hierarchyStores, nil, now)
	if err != nil {
		t.Fatal(err)
	}
	if len(skipped) != 0 {
		t.Fatalf("skipped = %v, want empty", skipped)
	}
	var closedNodeFound bool
	for _, node := range nodes {
		if node.Type == 4 && node.LMSID == "OC271" {
			closedNodeFound = true
			if node.ClosedAt != wantHierarchyClosedAt {
				t.Fatalf("closed hierarchy node = %+v", node)
			}
		}
	}
	if !closedNodeFound {
		t.Fatal("CLOSED store in an allowed district must participate in hierarchy")
	}
}

func TestPlanDashboardAreaSyncReusesRestoresAndInvalidatesNodes(t *testing.T) {
	existing := []dashboardModel.DashboardAreaModel{
		{BaseModel: appModels.BaseModel{ID: 10, CreatedAt: 100, UpdatedAt: 100}, Type: 1, Name: "Old Retail", LMSID: "RTL", Sort: 9},
		{BaseModel: appModels.BaseModel{ID: 11, CreatedAt: 100, UpdatedAt: 100}, Type: 2, Name: "East", LMSID: "EAST", ParentID: 10, Path: "10", Sort: 1},
		{BaseModel: appModels.BaseModel{ID: 12, CreatedAt: 100, UpdatedAt: 100, DeletedAt: appModels.DeletedAt(200)}, Type: 3, Name: "East 1", LMSID: "E1", ParentID: 11, Path: "10,11", Sort: 1},
		{BaseModel: appModels.BaseModel{ID: 13, CreatedAt: 100, UpdatedAt: 100}, Type: 4, Name: "OC999", LMSID: "OC999", ParentID: 12, Path: "10,11,12", Sort: 1},
	}
	desired := []dashboardAreaNode{
		{Type: 1, LMSID: "RTL", Name: "Retail", Sort: 1},
		{Type: 2, LMSID: "EAST", Name: "East", ParentLMSID: "RTL", Sort: 1},
		{Type: 3, LMSID: "E1", Name: "East 1", ParentLMSID: "EAST", Sort: 1},
		{Type: 4, LMSID: "OC267", Name: "OC267", ParentLMSID: "E1", ThirdID: 287, Sort: 1},
	}

	plan, err := planDashboardAreaSync(existing, desired, nil, nil, 1, 300)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Creates) != 1 || plan.Creates[0].ID != 14 || plan.Creates[0].ParentID != 12 || plan.Creates[0].Path != "10,11,12" {
		t.Fatalf("creates = %+v", plan.Creates)
	}
	if len(plan.Updates) != 2 {
		t.Fatalf("updates = %+v, want channel update and area restore", plan.Updates)
	}
	if plan.Updates[10]["name"] != "Retail" || plan.Updates[10]["sort"] != int64(1) {
		t.Fatalf("channel update = %+v", plan.Updates[10])
	}
	if plan.Updates[12]["deleted_at"] != int64(0) {
		t.Fatalf("area restore update = %+v", plan.Updates[12])
	}
	if len(plan.Invalidates) != 1 || plan.Invalidates[0] != 13 {
		t.Fatalf("invalidates = %v, want [13]", plan.Invalidates)
	}
	if plan.Stats.Inserted != 1 || plan.Stats.Updated != 2 || plan.Stats.Unchanged != 1 || plan.Stats.Invalidated != 1 {
		t.Fatalf("stats = %+v", plan.Stats)
	}
}

func TestPlanDashboardAreaSyncRejectsEmptySourceTree(t *testing.T) {
	existing := []dashboardModel.DashboardAreaModel{{
		BaseModel: appModels.BaseModel{ID: 10}, Type: 1, Name: "Retail", LMSID: "RTL",
	}}

	plan, err := planDashboardAreaSync(existing, nil, nil, nil, 0, 300)
	if err == nil {
		t.Fatal("empty source tree must fail instead of invalidating all existing nodes")
	}
	if len(plan.Invalidates) != 0 {
		t.Fatalf("invalidates = %v, want empty", plan.Invalidates)
	}
}

func TestBuildDashboardAreaNodesRejectsConflictingLMSHierarchy(t *testing.T) {
	now := time.Date(2026, 9, 2, 10, 0, 0, 0, time.Local)
	stores := []enterpriseModel.StoreMasterModel{
		newAreaTestStore("OC267", "OPEN", "", "RTL", "Retail", "CAFP_CN_ED", "East", "E1", "East 1"),
		newAreaTestStore("OC268", "OPEN", "", "OUTLET", "Outlet", "CAFP_CN_ED", "East", "E2", "East 2"),
	}

	if _, _, err := buildDashboardAreaNodes(stores, nil, now); err == nil {
		t.Fatal("same district code under different channels must fail")
	}
	stores[0], stores[1] = stores[1], stores[0]
	if _, _, err := buildDashboardAreaNodes(stores, nil, now); err == nil {
		t.Fatal("conflict detection must not depend on input order")
	}
}

func TestPlanDashboardAreaSyncProtectsSkippedStoreAndAncestorsFromInvalidation(t *testing.T) {
	existing := []dashboardModel.DashboardAreaModel{
		{BaseModel: appModels.BaseModel{ID: 10}, Type: 1, Name: "Retail", LMSID: "RTL"},
		{BaseModel: appModels.BaseModel{ID: 11}, Type: 2, Name: "East", LMSID: "EAST", ParentID: 10, Path: "10"},
		{BaseModel: appModels.BaseModel{ID: 12}, Type: 3, Name: "East 1", LMSID: "E1", ParentID: 11, Path: "10,11"},
		{BaseModel: appModels.BaseModel{ID: 13}, Type: 4, Name: "OC999", LMSID: "OC999", ParentID: 12, Path: "10,11,12"},
		{BaseModel: appModels.BaseModel{ID: 20}, Type: 1, Name: "Stale", LMSID: "STALE"},
	}
	desired := []dashboardAreaNode{{Type: 1, LMSID: "NEW", Name: "New", Sort: 1}}

	plan, err := planDashboardAreaSync(existing, desired, map[string]struct{}{"OC999": {}}, nil, 1, 300)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Invalidates) != 1 || plan.Invalidates[0] != 20 {
		t.Fatalf("invalidates = %v, want only unrelated stale node 20", plan.Invalidates)
	}
}

func TestPlanDashboardAreaSyncUpdatesOnlyClosedAtForClosedStore(t *testing.T) {
	existing := []dashboardModel.DashboardAreaModel{
		{BaseModel: appModels.BaseModel{ID: 10}, Type: 1, Name: "Retail", LMSID: "RTL"},
		{BaseModel: appModels.BaseModel{ID: 11}, Type: 2, Name: "East", LMSID: "EAST", ParentID: 10, Path: "10"},
		{BaseModel: appModels.BaseModel{ID: 12}, Type: 3, Name: "East 1", LMSID: "E1", ParentID: 11, Path: "10,11"},
		{BaseModel: appModels.BaseModel{ID: 13}, Type: 4, Name: "OC267", LMSID: "OC267", ParentID: 12, Path: "10,11,12", ThirdID: 287, Sort: 3},
		{BaseModel: appModels.BaseModel{ID: 20}, Type: 1, Name: "Stale", LMSID: "STALE"},
	}
	closedAt := time.Date(2022, 8, 15, 0, 0, 0, 0, time.Local).Unix()

	plan, err := planDashboardAreaSync(
		existing,
		nil,
		map[string]struct{}{"OC267": {}},
		map[string]int64{"OC267": closedAt},
		0,
		300,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Creates) != 0 || len(plan.Invalidates) != 1 || plan.Invalidates[0] != 20 {
		t.Fatalf("creates = %+v, invalidates = %v; closed store must retain historical hierarchy and invalidate unrelated nodes", plan.Creates, plan.Invalidates)
	}
	updates := plan.Updates[13]
	if len(updates) != 2 || updates["closed_at"] != closedAt || updates["updated_at"] != int64(300) {
		t.Fatalf("closed store updates = %+v, want only closed_at and updated_at", updates)
	}
	for _, field := range []string{"name", "parent_id", "path", "third_id", "sort", "deleted_at"} {
		if _, exists := updates[field]; exists {
			t.Fatalf("closed store must not update %s: %+v", field, updates)
		}
	}
}

func TestPlanDashboardAreaSyncRejectsAllInvalidOpenAndPendingHierarchy(t *testing.T) {
	plan, err := planDashboardAreaSync(
		nil,
		nil,
		map[string]struct{}{"OC999": {}},
		nil,
		1,
		300,
	)
	if err == nil {
		t.Fatal("OPEN or PENDING source rows with no valid hierarchy must fail")
	}
	if len(plan.Invalidates) != 0 {
		t.Fatalf("invalidates = %v, want empty", plan.Invalidates)
	}
}

type fakeDashboardAreaRepository struct {
	enterpriseID int64
	rootThirdID  int
	stats        dashboardAreaSyncStats
	skipped      []string
	err          error
}

func (repository *fakeDashboardAreaRepository) SyncDashboardAreas(_ app.Ctx, enterpriseID int64, rootThirdID int, _ time.Time) (dashboardAreaSyncStats, []string, error) {
	repository.enterpriseID = enterpriseID
	repository.rootThirdID = rootThirdID
	return repository.stats, repository.skipped, repository.err
}

func TestSyncCoachLMSDashboardAreaUsesCoachLMSEnterpriseAndStoreRootDepartment(t *testing.T) {
	repository := &fakeDashboardAreaRepository{stats: dashboardAreaSyncStats{Inserted: 4, Updated: 2, Unchanged: 3, Invalidated: 1}, skipped: []string{"OC999"}}
	var requestedSFTPType int
	var notifications []string
	job := &SyncCoachLMSDashboardArea{
		getStoreRootDepartment: func(_ app.Ctx, enterpriseID int64) (int, error) {
			if enterpriseID != 3360610327843968 {
				return 0, errors.New("unexpected enterprise")
			}
			return 340, nil
		},
		repository: repository,
		now:        func() time.Time { return time.Unix(300, 0) },
		notify: func(_ app.Ctx, _, _, message string) {
			notifications = append(notifications, message)
		},
	}
	jobConfig := micro_scrm.BrandConfigSftpDetail{BrandConfigSftpModel: micro_scrm.BrandConfigSftpModel{BrandName: "Coach"}}
	jobConfig.EnterpriseID = 3360610327843968
	job.getConfigs = func(sftpType int) ([]micro_scrm.BrandConfigSftpDetail, error) {
		requestedSFTPType = sftpType
		return []micro_scrm.BrandConfigSftpDetail{{BrandConfigSftpModel: micro_scrm.BrandConfigSftpModel{BrandName: "other"}}, jobConfig}, nil
	}

	if err := job.Run(toolboxCtx.New(), nil); err != nil {
		t.Fatal(err)
	}
	if requestedSFTPType != micro_scrm.SftpTypeDEL || repository.enterpriseID != 3360610327843968 || repository.rootThirdID != 340 {
		t.Fatalf("sftpType=%d enterpriseID=%d rootThirdID=%d", requestedSFTPType, repository.enterpriseID, repository.rootThirdID)
	}
	if len(notifications) != 1 || notifications[0] != "成功：同步LMS组织架构，新增4，更新2，未变化3，失效1，跳过1（OC999）" {
		t.Fatalf("notifications = %v", notifications)
	}
}

func newAreaTestStore(storeCode, status, closeDate, channelCode, channelName, districtCode, districtName, areaCode, areaName string) enterpriseModel.StoreMasterModel {
	return enterpriseModel.StoreMasterModel{
		StoreCode:       storeCode,
		StoreStatus:     status,
		ClosedDate:      closeDate,
		RegionDesc:      "China Region",
		OrgChannelCode:  channelCode,
		OrgChannelName:  channelName,
		OrgDistrictCode: districtCode,
		OrgDistrictName: districtName,
		OrgAreaCode:     areaCode,
		OrgAreaName:     areaName,
	}
}

func assertAreaNode(t *testing.T, node dashboardAreaNode, nodeType int, lmsID, name, parentLMSID string, thirdID int, closedAt int64, sort int64) {
	t.Helper()
	if node.Type != nodeType || node.LMSID != lmsID || node.Name != name || node.ParentLMSID != parentLMSID || node.ThirdID != thirdID || node.ClosedAt != closedAt || node.Sort != sort {
		t.Fatalf("node = %+v", node)
	}
}
