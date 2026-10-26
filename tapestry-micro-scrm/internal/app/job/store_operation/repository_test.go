package storeOperationJob

import (
	"testing"

	enterpriseModel "micro-scrm/internal/orm/model/micro_scrm_enterprise"
)

func TestStoreSourceChangedIgnoresLocalFieldsAndTimestamps(t *testing.T) {
	existing := enterpriseModel.StoreMasterModel{}
	existing.StoreCode = "OC267"
	existing.StoreStatus = "OPEN"
	existing.StoreNameZh = "青岛海信广场"
	existing.AddressZh = "上海市"
	existing.StoreArea = "4928"
	existing.UpdatedAt = 100
	incoming := storeSourceRecord{StoreCode: "OC267", StoreStatus: "OPEN"}

	if storeSourceChanged(existing, incoming) {
		t.Fatal("local fields and timestamps must not trigger an LMS update")
	}
}

func TestStoreSourceChangedDetectsUpstreamChange(t *testing.T) {
	existing := enterpriseModel.StoreMasterModel{}
	existing.StoreCode = "OC267"
	existing.StoreStatus = "OPEN"
	incoming := storeSourceRecord{StoreCode: "OC267", StoreStatus: "CLOSED"}

	if !storeSourceChanged(existing, incoming) {
		t.Fatal("changed LMS status must trigger an update")
	}
}

func TestEmployeeSourceChangedDetectsEmptyOverwrite(t *testing.T) {
	existing := enterpriseModel.EmployeeMasterModel{}
	existing.EmployeeID = "000931251"
	existing.EmployeeCode = "931251"
	existing.EmailAddress = "old@coach.com"
	incoming := employeeSourceRecord{EmployeeID: "000931251", EmployeeCode: "931251", EmailAddress: ""}

	if !employeeSourceChanged(existing, incoming) {
		t.Fatal("empty Workday value must overwrite an old source value")
	}
}

func TestStoreUpdateValuesExcludeLocalFields(t *testing.T) {
	record := storeSourceRecord{StoreCode: "OC267", StoreStatus: "OPEN", OrganizationName: "SHANGHAI MIXC"}
	values := storeUpdateValues(record, 1787100000)

	for _, forbidden := range []string{"store_area", "store_name_zh", "address_zh", "day_phone", "email_id", "payment_methods", "created_at"} {
		if _, ok := values[forbidden]; ok {
			t.Fatalf("store update unexpectedly contains local field %q", forbidden)
		}
	}
	if values["store_status"] != "OPEN" || values["updated_at"] != int64(1787100000) {
		t.Fatalf("store update values = %#v", values)
	}
}

func TestEmployeeUpdateValuesIncludeEmptySourceFields(t *testing.T) {
	record := employeeSourceRecord{EmployeeID: "000931251", EmployeeCode: "931251", EmailAddress: ""}
	values := employeeUpdateValues(record, 1787100000)

	value, ok := values["email_address"]
	if !ok || value != "" {
		t.Fatalf("employee update email_address = %#v, exists=%v", value, ok)
	}
	if _, ok := values["created_at"]; ok {
		t.Fatal("employee update must preserve created_at")
	}
}

func TestApplyEmployeeStoreDepartmentsUsesStoreCodeForStoreEmployee(t *testing.T) {
	records := []employeeSourceRecord{{
		EmployeeCode: "931251",
		StoreCode:    "OC267",
		Department:   "Sales - Direct",
	}}

	applyEmployeeStoreDepartments(records, map[string]struct{}{"OC267": {}})

	if records[0].Department != "OC267" {
		t.Fatalf("store employee department = %q, want OC267", records[0].Department)
	}
}

func TestApplyEmployeeStoreDepartmentsKeepsWorkdayDepartmentWithoutStore(t *testing.T) {
	records := []employeeSourceRecord{
		{EmployeeCode: "931251", StoreCode: "CO0003", Department: "Marketing"},
		{EmployeeCode: "931252", StoreCode: "", Department: "Finance"},
	}

	applyEmployeeStoreDepartments(records, map[string]struct{}{"OC267": {}})

	if records[0].Department != "Marketing" {
		t.Fatalf("non-store employee department = %q, want Marketing", records[0].Department)
	}
	if records[1].Department != "Finance" {
		t.Fatalf("employee without store code department = %q, want Finance", records[1].Department)
	}
}

func TestStoreEmployeeDepartmentDoesNotTriggerRepeatedUpdate(t *testing.T) {
	existing := enterpriseModel.EmployeeMasterModel{}
	existing.EmployeeID = "000931251"
	existing.EmployeeCode = "931251"
	existing.StoreCode = "OC267"
	existing.Department = "OC267"
	incoming := []employeeSourceRecord{{
		EmployeeID:   "000931251",
		EmployeeCode: "931251",
		StoreCode:    "OC267",
		Department:   "Sales - Direct",
	}}

	applyEmployeeStoreDepartments(incoming, map[string]struct{}{"OC267": {}})

	if employeeSourceChanged(existing, incoming[0]) {
		t.Fatal("unchanged store employee must not be updated repeatedly")
	}
}
