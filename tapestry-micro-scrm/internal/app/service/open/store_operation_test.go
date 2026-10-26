package openSvc

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"micro-scrm/internal/app"
	openDao "micro-scrm/internal/app/dao/open"
	openReq "micro-scrm/internal/app/http/request/open"
	enterpriseModel "micro-scrm/internal/orm/model/micro_scrm_enterprise"
	pkgKMS "micro-scrm/internal/pkg/kms"
)

type storeMasterFakeRepository struct {
	stores          []enterpriseModel.StoreMasterModel
	total           int64
	departmentIDs   map[string]int64
	receivedRequest openReq.StoreMasterListReq
	receivedCodes   []string
}

type employeeMasterFakeRepository struct {
	employees       []openDao.EmployeeMasterRow
	total           int64
	wecom           map[string]openDao.EmployeeWeComData
	receivedRequest openReq.EmployeeMasterListReq
	receivedCodes   []string
}

func (f *employeeMasterFakeRepository) EmployeeMasterPage(_ app.Ctx, _ int64, req openReq.EmployeeMasterListReq) ([]openDao.EmployeeMasterRow, int64, error) {
	f.receivedRequest = req
	return f.employees, f.total, nil
}

func (f *employeeMasterFakeRepository) EmployeeWeComData(_ app.Ctx, _ int64, employeeCodes []string) (map[string]openDao.EmployeeWeComData, error) {
	f.receivedCodes = append([]string(nil), employeeCodes...)
	return f.wecom, nil
}

type storeMasterFakeKMS struct {
	dataKey      []byte
	encryptedKey string
	err          error
	calls        int
	decryptCalls int
}

func (f *storeMasterFakeKMS) GenerateDataKey(context.Context) ([]byte, string, error) {
	f.calls++
	if f.err != nil {
		return nil, "", f.err
	}
	return append([]byte(nil), f.dataKey...), f.encryptedKey, nil
}

func (f *storeMasterFakeKMS) DecryptDataKey(_ context.Context, encryptedKey string) ([]byte, error) {
	f.decryptCalls++
	if f.err != nil {
		return nil, f.err
	}
	if encryptedKey != f.encryptedKey {
		return nil, errors.New("unknown encrypted data key")
	}
	return append([]byte(nil), f.dataKey...), nil
}

func (f *storeMasterFakeRepository) StoreMasterPage(_ app.Ctx, _ int64, req openReq.StoreMasterListReq) ([]enterpriseModel.StoreMasterModel, int64, error) {
	f.receivedRequest = req
	return f.stores, f.total, nil
}

func (f *storeMasterFakeRepository) WechatDepartmentIDs(_ app.Ctx, _ int64, storeCodes []string) (map[string]int64, error) {
	f.receivedCodes = append([]string(nil), storeCodes...)
	return f.departmentIDs, nil
}

func TestStoreMasterDefaultsAndCountryValidation(t *testing.T) {
	req := openReq.StoreMasterListReq{}
	if err := normalizeStoreMasterRequest(&req); err != nil {
		t.Fatal(err)
	}
	if req.Page != 1 || req.Size != 100 || !reflect.DeepEqual(req.StoreCountryCodes, []string{"CN"}) {
		t.Fatalf("默认参数=%+v", req)
	}

	req = openReq.StoreMasterListReq{Page: 1, Size: 100, StoreCountryCodes: []string{"cn", " HK ", "CN", "MO", "TW"}}
	if err := normalizeStoreMasterRequest(&req); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(req.StoreCountryCodes, []string{"CN", "HK", "MO", "TW"}) {
		t.Fatalf("国家编码=%v", req.StoreCountryCodes)
	}

	for _, invalid := range []openReq.StoreMasterListReq{
		{Page: -1, Size: 100}, {Page: 1, Size: 1001}, {Page: 1, Size: 100, UpdatedAt: -1},
		{Page: 1, Size: 100, StoreCountryCodes: []string{"US"}},
	} {
		if err := normalizeStoreMasterRequest(&invalid); err == nil {
			t.Fatalf("非法参数应失败: %+v", invalid)
		}
	}
}

func TestEmployeeMasterDefaultsAndScopeValidation(t *testing.T) {
	req := openReq.EmployeeMasterListReq{}
	if err := normalizeEmployeeMasterRequest(&req); err != nil {
		t.Fatal(err)
	}
	if req.Page != 1 || req.Size != 100 || req.StoreScope != "all" || req.EmploymentScope != "active" || len(req.StoreCountryCodes) != 0 {
		t.Fatalf("默认参数=%+v", req)
	}

	req = openReq.EmployeeMasterListReq{Page: 1, Size: 100, StoreScope: "store", EmploymentScope: "inactive", StoreCountryCodes: []string{"cn", " HK ", "CN"}}
	if err := normalizeEmployeeMasterRequest(&req); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(req.StoreCountryCodes, []string{"CN", "HK"}) {
		t.Fatalf("国家编码=%v", req.StoreCountryCodes)
	}

	for _, invalid := range []openReq.EmployeeMasterListReq{
		{Page: 1, Size: 1001}, {Page: 1, Size: 100, StoreScope: "closed_store"},
		{Page: 1, Size: 100, EmploymentScope: "employed"},
		{Page: 1, Size: 100, StoreCountryCodes: []string{"US"}},
	} {
		if err := normalizeEmployeeMasterRequest(&invalid); err == nil {
			t.Fatalf("非法参数应失败: %+v", invalid)
		}
	}
}

func TestEmployeeMasterMappingIncludesWeComAndUnmatchedDefaults(t *testing.T) {
	first := employeeMasterTestRow()
	second := employeeMasterTestRow()
	second.EmployeeCode = "882211"
	second.EmployeeID = "000882211"
	second.StoreCode = ""
	repository := &employeeMasterFakeRepository{
		employees: []openDao.EmployeeMasterRow{first, second},
		total:     2,
		wecom: map[string]openDao.EmployeeWeComData{
			"931251": {
				WecomUserid:        "931251",
				WecomJoinedAt:      1772182945,
				VerificationStatus: 5,
				WecomLeftAt:        1773000000,
				DepartmentChanges: []openDao.EmployeeDepartmentChange{{
					EventType: "join", EventTime: 1772182945, WechatDepartmentID: 443, DepartmentName: "OC339",
				}},
			},
		},
	}
	service := newStoreOperationServiceWithEmployeeRepository(repository)
	result, total, err := service.loadEmployeeMasterList(nil, openReq.EmployeeMasterListReq{EnterpriseID: 3360610327843968})
	if err != nil {
		t.Fatal(err)
	}
	if total != 2 || len(result) != 2 {
		t.Fatalf("total=%d result=%d", total, len(result))
	}
	if result[0].WecomUserid != "931251" || result[0].Response.VerificationStatus != 5 || result[0].Response.WecomLeftAt != 1773000000 {
		t.Fatalf("企微字段=%+v", result[0])
	}
	if len(result[0].Response.DepartmentChanges) != 1 || result[0].Response.DepartmentChanges[0].DepartmentName != "OC339" {
		t.Fatalf("部门变动=%+v", result[0].Response.DepartmentChanges)
	}
	if result[1].WecomUserid != "" || result[1].Response.VerificationStatus != 0 || result[1].Response.WecomJoinedAt != 0 || result[1].Response.DepartmentChanges == nil {
		t.Fatalf("未匹配员工附带字段=%+v", result[1])
	}
	if !reflect.DeepEqual(repository.receivedCodes, []string{"931251", "882211"}) {
		t.Fatalf("企微批量查询工号=%v", repository.receivedCodes)
	}
	if repository.receivedRequest.EmploymentScope != "active" || repository.receivedRequest.StoreScope != "all" {
		t.Fatalf("标准化查询参数=%+v", repository.receivedRequest)
	}
}

func TestEmployeeMasterEncryptionCoversSevenPIFields(t *testing.T) {
	employee := employeeMasterTestRow()
	employee.NameZh = ""
	repository := &employeeMasterFakeRepository{
		employees: []openDao.EmployeeMasterRow{employee},
		total:     1,
		wecom: map[string]openDao.EmployeeWeComData{
			"931251": {WecomUserid: "931251"},
		},
	}
	client := &storeMasterFakeKMS{
		dataKey:      []byte("0123456789abcdef0123456789abcdef"),
		encryptedKey: "employee-page-data-key",
	}
	service := newStoreOperationServiceWithEmployeeDependencies(repository, client)
	result, total, err := service.EmployeeMasterList(nil, openReq.EmployeeMasterListReq{})
	if err != nil {
		t.Fatal(err)
	}
	if total != 1 || len(result) != 1 || client.calls != 1 {
		t.Fatalf("total=%d result=%d KMS调用=%d", total, len(result), client.calls)
	}
	item := result[0]
	values := []pkgKMS.EncryptedValue{
		item.Name, item.NameZh, item.EmployeeID, item.EmployeeCode,
		item.WecomUserid, item.MobileNumber, item.EmailAddress,
	}
	if values[1] != pkgKMS.EmptyEncryptedValue() {
		t.Fatalf("空中文名加密对象=%+v", values[1])
	}
	ivSet := make(map[string]struct{})
	for index, value := range values {
		if index == 1 {
			continue
		}
		if value.EncryptedDataKey != "employee-page-data-key" || value.IV == "" {
			t.Fatalf("第%d个PI加密对象=%+v", index, value)
		}
		ivSet[value.IV] = struct{}{}
	}
	if len(ivSet) != 6 {
		t.Fatalf("六个非空PI字段应使用独立IV，实际=%v", ivSet)
	}
}

func TestEmployeeMasterEncryptionFailureFailsWholeList(t *testing.T) {
	repository := &employeeMasterFakeRepository{employees: []openDao.EmployeeMasterRow{employeeMasterTestRow()}, total: 1}
	client := &storeMasterFakeKMS{err: errors.New("kms unavailable")}
	service := newStoreOperationServiceWithEmployeeDependencies(repository, client)
	result, total, err := service.EmployeeMasterList(nil, openReq.EmployeeMasterListReq{})
	if err == nil || result != nil || total != 0 {
		t.Fatalf("KMS失败应使员工列表整体失败，result=%v total=%d err=%v", result, total, err)
	}
}

func TestStoreOperationEncryptedValuesCanBeDecryptedEndToEnd(t *testing.T) {
	client := &storeMasterFakeKMS{
		dataKey:      []byte("0123456789abcdef0123456789abcdef"),
		encryptedKey: "end-to-end-page-key",
	}
	storeRepository := &storeMasterFakeRepository{
		stores: []enterpriseModel.StoreMasterModel{storeMasterTestModel()}, total: 1,
		departmentIDs: map[string]int64{"OC267": 287},
	}
	storeService := newStoreOperationServiceWithDependencies(storeRepository, client)
	stores, _, err := storeService.StoreMasterList(nil, openReq.StoreMasterListReq{})
	if err != nil {
		t.Fatal(err)
	}
	employeeRepository := &employeeMasterFakeRepository{
		employees: []openDao.EmployeeMasterRow{employeeMasterTestRow()}, total: 1,
		wecom: map[string]openDao.EmployeeWeComData{"931251": {WecomUserid: "931251"}},
	}
	employeeService := newStoreOperationServiceWithEmployeeDependencies(employeeRepository, client)
	employees, _, err := employeeService.EmployeeMasterList(nil, openReq.EmployeeMasterListReq{})
	if err != nil {
		t.Fatal(err)
	}
	if client.calls != 2 {
		t.Fatalf("门店页和员工页应各生成一次数据密钥，实际=%d", client.calls)
	}

	decryptService := NewDecryptServiceWithClient(client)
	response, err := decryptService.Decrypt(context.Background(), openReq.DecryptReq{Items: []openReq.DecryptItem{
		decryptRequestItem("area_manager", stores[0].AreaManager),
		decryptRequestItem("employee_name", employees[0].Name),
		decryptRequestItem("employee_code", employees[0].EmployeeCode),
	}})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"Jean Yuan", "Susu Wang", "931251"}
	for index, item := range response.Items {
		if !item.Success || item.Plaintext != want[index] {
			t.Fatalf("第%d项解密结果=%+v", index, item)
		}
	}
}

func TestStoreMasterMappingAndBatchDepartmentLookup(t *testing.T) {
	store := storeMasterTestModel()

	repository := &storeMasterFakeRepository{
		stores:        []enterpriseModel.StoreMasterModel{store},
		total:         1,
		departmentIDs: map[string]int64{"OC267": 287},
	}
	service := newStoreOperationServiceWithRepository(repository)
	result, total, err := service.loadStoreMasterList(nil, openReq.StoreMasterListReq{EnterpriseID: 3360610327843968, UpdatedAt: 1780000000})
	if err != nil {
		t.Fatal(err)
	}
	if total != 1 || len(result) != 1 {
		t.Fatalf("total=%d result=%d", total, len(result))
	}
	got := result[0]
	if got.Response.StoreCode != "OC267" || got.Response.Channel != "RETAIL" || got.Response.StoreNameEn != "SHANGHAI MIXC" || got.Response.Status != "OPEN" {
		t.Fatalf("基础字段映射=%+v", got.Response)
	}
	if got.Response.WechatDepartmentID != 287 || got.Response.BusinessOrganization.AreaCode != "Area_Code_CAFP_CN_EA1" {
		t.Fatalf("关联字段映射=%+v", got.Response)
	}
	if got.AreaManager != "Jean Yuan" || got.StoreDistrictManager != "Sandy Tang" {
		t.Fatalf("PI源字段映射=%+v", got)
	}
	if !reflect.DeepEqual(repository.receivedCodes, []string{"OC267"}) {
		t.Fatalf("应按当前页批量查询门店编码，实际=%v", repository.receivedCodes)
	}
	if repository.receivedRequest.UpdatedAt != 1780000000 || !reflect.DeepEqual(repository.receivedRequest.StoreCountryCodes, []string{"CN"}) {
		t.Fatalf("查询参数=%+v", repository.receivedRequest)
	}
}

func TestStoreMasterEncryptionReusesPageKeyWithDistinctIVs(t *testing.T) {
	first := storeMasterTestModel()
	second := storeMasterTestModel()
	second.StoreCode = "OC268"
	second.AreaManager = "Alice Chen"
	second.StoreDistrictMgr = "Bob Li"
	repository := &storeMasterFakeRepository{
		stores:        []enterpriseModel.StoreMasterModel{first, second},
		total:         2,
		departmentIDs: map[string]int64{},
	}
	client := &storeMasterFakeKMS{
		dataKey:      []byte("0123456789abcdef0123456789abcdef"),
		encryptedKey: "page-encrypted-data-key",
	}
	service := newStoreOperationServiceWithDependencies(repository, client)
	result, total, err := service.StoreMasterList(nil, openReq.StoreMasterListReq{})
	if err != nil {
		t.Fatal(err)
	}
	if total != 2 || len(result) != 2 || client.calls != 1 {
		t.Fatalf("total=%d result=%d KMS调用=%d", total, len(result), client.calls)
	}
	values := []pkgKMS.EncryptedValue{
		result[0].AreaManager, result[0].StoreDistrictManager,
		result[1].AreaManager, result[1].StoreDistrictManager,
	}
	ivSet := make(map[string]struct{}, len(values))
	for _, value := range values {
		if value.EncryptedDataKey != "page-encrypted-data-key" {
			t.Fatalf("页面内密文密钥不一致: %+v", value)
		}
		if value.IV == "" {
			t.Fatal("非空PI字段必须生成IV")
		}
		ivSet[value.IV] = struct{}{}
	}
	if len(ivSet) != len(values) {
		t.Fatalf("每个PI字段必须使用独立IV，实际=%v", ivSet)
	}
}

func TestStoreMasterEncryptionFailureFailsWholeList(t *testing.T) {
	repository := &storeMasterFakeRepository{stores: []enterpriseModel.StoreMasterModel{storeMasterTestModel()}, total: 1}
	client := &storeMasterFakeKMS{err: errors.New("kms unavailable")}
	service := newStoreOperationServiceWithDependencies(repository, client)
	result, total, err := service.StoreMasterList(nil, openReq.StoreMasterListReq{})
	if err == nil || result != nil || total != 0 {
		t.Fatalf("KMS失败应使列表整体失败，result=%v total=%d err=%v", result, total, err)
	}
}

func TestStoreOperationEmptyPagesReturnEmptyListsWithoutKMS(t *testing.T) {
	client := &storeMasterFakeKMS{
		dataKey:      []byte("0123456789abcdef0123456789abcdef"),
		encryptedKey: "unused",
	}
	storeService := newStoreOperationServiceWithDependencies(&storeMasterFakeRepository{}, client)
	stores, total, err := storeService.StoreMasterList(nil, openReq.StoreMasterListReq{})
	if err != nil || total != 0 || stores == nil || len(stores) != 0 {
		t.Fatalf("空门店页 stores=%#v total=%d err=%v", stores, total, err)
	}
	employeeService := newStoreOperationServiceWithEmployeeDependencies(&employeeMasterFakeRepository{}, client)
	employees, total, err := employeeService.EmployeeMasterList(nil, openReq.EmployeeMasterListReq{})
	if err != nil || total != 0 || employees == nil || len(employees) != 0 {
		t.Fatalf("空员工页 employees=%#v total=%d err=%v", employees, total, err)
	}
	if client.calls != 0 {
		t.Fatalf("空分页不应调用KMS，实际=%d", client.calls)
	}
}

func storeMasterTestModel() enterpriseModel.StoreMasterModel {
	store := enterpriseModel.StoreMasterModel{}
	store.StoreCode = "OC267"
	store.CountryCode = "CN"
	store.SalesChannel = "RETAIL"
	store.StoreNameZh = "青岛海信广场"
	store.OrganizationName = "SHANGHAI MIXC"
	store.StoreStatus = "OPEN"
	store.OpenDate = "2021-07-01 00:00:00"
	store.ClosedDate = "2099-12-31 00:00:00"
	store.RegionDesc = "China Region"
	store.State = "N/A"
	store.City = "Shanghai"
	store.AddressZh = "陕西省西安市雁塔区长安中路123号"
	store.AddressLine1 = "Shop L134, L1 Floor, Shanghai MixC"
	store.AddressLine2 = "1599 Wuzhong Road, Minhang District"
	store.ZipCode = "200000"
	store.DayPhone = "041184323929"
	store.EmailID = "CNS360@coach.com"
	store.Longitude = "-94.495340"
	store.Latitude = "36.596890"
	store.StoreWorkingHours = "Mon-Sun 10:00-22:00"
	store.StoreArea = "180㎡"
	store.AppointmentServiceTypes = "皮牌压印、皮具护理"
	store.BrandCode = "COACH"
	store.OrganizationCode = "D369"
	store.PaymentMethods = "自收银"
	store.CompanyCode = "3000"
	store.StoreClass = "S"
	store.OrgRegionCode = "RC_CAFP_CN"
	store.OrgRegionName = "China Region"
	store.OrgDistrictCode = "CAFP_CN_ED"
	store.OrgDistrictName = "Retail-East"
	store.OrgAreaCode = "Area_Code_CAFP_CN_EA1"
	store.OrgAreaName = "East 1 - SHANGHAI + WUXI"
	store.AreaManager = "Jean Yuan"
	store.StoreDistrictMgr = "Sandy Tang"
	store.StoreDesignConcept = "Modern Luxury"
	store.UpdatedAt = 1780904216

	return store
}

func employeeMasterTestRow() openDao.EmployeeMasterRow {
	row := openDao.EmployeeMasterRow{OrganizationCode: "D440"}
	row.EmployeeID = "000931251"
	row.EmployeeCode = "931251"
	row.MobileNumber = "18700001251"
	row.Name = "Susu Wang"
	row.NameZh = "王小散"
	row.JobPosition = "Sales Associate"
	row.Gender = "F"
	row.EmailAddress = "Lily@coach.com"
	row.Status = "1"
	row.StoreCode = "OC287"
	row.Department = "Sales - Direct"
	row.HireDate = "2026/7/2 0:00:00"
	row.TerminationDate = ""
	row.JobFunction = "VM"
	row.IsStoreManagement = "否"
	row.UpdatedAt = 1780904216
	return row
}
