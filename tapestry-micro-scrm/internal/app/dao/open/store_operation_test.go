package openDao

import (
	"strings"
	"testing"

	openReq "micro-scrm/internal/app/http/request/open"
	enterpriseModel "micro-scrm/internal/orm/model/micro_scrm_enterprise"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestStoreMasterQueryUsesBranchTableAndRequiredFilters(t *testing.T) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "gorm:gorm@tcp(localhost:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local",
		SkipInitializeWithVersion: true,
	}), &gorm.Config{DryRun: true, DisableAutomaticPing: true})
	if err != nil {
		t.Fatal(err)
	}
	req := openReq.StoreMasterListReq{StoreCountryCodes: []string{"CN", "HK"}, UpdatedAt: 1780000000}
	var result []enterpriseModel.StoreMasterModel
	statement := buildStoreMasterQuery(db, 3360610327843968, req).
		Order("updated_at ASC, id ASC").Offset(100).Limit(100).Find(&result).Statement
	sql := statement.SQL.String()
	for _, fragment := range []string{
		"`micro_scrm_enterprise_3360610327843968`.`store_master`",
		"deleted_at = ?",
		"country_code IN (?,?)",
		"updated_at >= ?",
		"ORDER BY updated_at ASC, id ASC",
		"LIMIT 100 OFFSET 100",
	} {
		if !strings.Contains(sql, fragment) {
			t.Fatalf("SQL缺少%q: %s", fragment, sql)
		}
	}
	if len(statement.Vars) != 4 {
		t.Fatalf("SQL参数=%#v", statement.Vars)
	}
}

func TestStoreMasterQueryOmitsIncrementFilterForFullSync(t *testing.T) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "gorm:gorm@tcp(localhost:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local",
		SkipInitializeWithVersion: true,
	}), &gorm.Config{DryRun: true, DisableAutomaticPing: true})
	if err != nil {
		t.Fatal(err)
	}
	var result []enterpriseModel.StoreMasterModel
	sql := buildStoreMasterQuery(db, 1, openReq.StoreMasterListReq{StoreCountryCodes: []string{"CN"}}).
		Find(&result).Statement.SQL.String()
	if strings.Contains(sql, "updated_at >=") {
		t.Fatalf("全量查询不应包含增量条件: %s", sql)
	}
}

func TestEmployeeMasterQueryAppliesScopesAndIncrement(t *testing.T) {
	db := newStoreOperationDryRunDB(t)
	req := openReq.EmployeeMasterListReq{
		StoreScope:        "open_store",
		EmploymentScope:   "active",
		StoreCountryCodes: []string{"CN", "HK"},
		UpdatedAt:         1780000000,
	}
	var result []EmployeeMasterRow
	statement := buildEmployeeMasterQuery(db, 3360610327843968, req).
		Select("em.*, sm.organization_code AS organization_code").Find(&result).Statement
	sql := statement.SQL.String()
	for _, fragment := range []string{
		"micro_scrm_enterprise_3360610327843968.employee_master AS em",
		"CONVERT(sm.store_code USING utf8mb4) COLLATE utf8mb4_unicode_ci = CONVERT(em.store_code USING utf8mb4) COLLATE utf8mb4_unicode_ci",
		"em.deleted_at = ?",
		"em.status = ?",
		"sm.id IS NOT NULL AND sm.store_status = ?",
		"sm.country_code IN (?,?)",
		"em.updated_at >= ?",
	} {
		if !strings.Contains(sql, fragment) {
			t.Fatalf("SQL缺少%q: %s", fragment, sql)
		}
	}
}

func TestEmployeeMasterAllScopeKeepsEmployeesWithoutStore(t *testing.T) {
	db := newStoreOperationDryRunDB(t)
	var result []EmployeeMasterRow
	sql := buildEmployeeMasterQuery(db, 1, openReq.EmployeeMasterListReq{StoreScope: "all", EmploymentScope: "all"}).
		Find(&result).Statement.SQL.String()
	for _, forbidden := range []string{"sm.id IS NOT NULL", "em.status =", "em.status <>"} {
		if strings.Contains(sql, forbidden) {
			t.Fatalf("all范围不应包含%q: %s", forbidden, sql)
		}
	}
}

func TestEmployeeWeComMatchPrefersUseridOverDirtyOpenUserid(t *testing.T) {
	members := []employeeMemberRow{
		{ID: 1, Userid: "other", OpenUserid: "931251"},
		{ID: 2, Userid: "931251", OpenUserid: "dirty"},
		{ID: 3, Userid: "fallback", OpenUserid: "882211"},
	}
	matched := matchEmployeeMembers([]string{"931251", "882211", "missing"}, members)
	if matched["931251"].ID != 2 {
		t.Fatalf("userid精确匹配应优先，实际=%+v", matched["931251"])
	}
	if matched["882211"].ID != 3 {
		t.Fatalf("未匹配userid时应使用open_userid，实际=%+v", matched["882211"])
	}
	if _, exists := matched["missing"]; exists {
		t.Fatal("不存在的企微员工不应产生匹配")
	}
}

func TestEmployeeWeComDeletedMemberUsesLeaveStatusAndTime(t *testing.T) {
	data := assembleEmployeeWeComData(map[string]employeeMemberRow{
		"931251": {ID: 2, Userid: "931251", Status: 1, CreatedAt: 1772182945, DeletedAt: 1773000000},
	}, nil)["931251"]
	if data.WecomUserid != "931251" || data.WecomJoinedAt != 1772182945 || data.VerificationStatus != 5 || data.WecomLeftAt != 1773000000 {
		t.Fatalf("离开企微员工状态=%+v", data)
	}
	if data.DepartmentChanges == nil || len(data.DepartmentChanges) != 0 {
		t.Fatalf("无部门变动应返回空数组，实际=%#v", data.DepartmentChanges)
	}
}

func TestDepartmentChangesBuildJoinRejoinAndLeaveInOrder(t *testing.T) {
	changes := buildDepartmentChanges([]employeeDepartmentRow{
		{MemberID: 1, CreatedAt: 200, UpdatedAt: 300, DeletedAt: 300, WechatDepartmentID: 443, DepartmentName: "OC339"},
		{MemberID: 1, CreatedAt: 100, UpdatedAt: 100, WechatDepartmentID: 287, DepartmentName: "OC267"},
	})[1]
	if len(changes) != 4 {
		t.Fatalf("部门变动数量=%d: %+v", len(changes), changes)
	}
	wantTypes := []string{"join", "join", "join", "leave"}
	wantTimes := []int64{100, 200, 300, 300}
	for index := range changes {
		if changes[index].EventType != wantTypes[index] || changes[index].EventTime != wantTimes[index] {
			t.Fatalf("第%d条变动=%+v", index, changes[index])
		}
	}
}

func newStoreOperationDryRunDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "gorm:gorm@tcp(localhost:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local",
		SkipInitializeWithVersion: true,
	}), &gorm.Config{DryRun: true, DisableAutomaticPing: true})
	if err != nil {
		t.Fatal(err)
	}
	return db
}
