package storeOperationJob

import (
	"strings"
	"testing"
)

func TestDeriveEmployeeCodeRemovesExactlyThreeLeadingZeros(t *testing.T) {
	cases := map[string]string{
		"000931251": "931251",
		"001234":    "001234",
		"000000123": "000123",
		"xiaopeng":  "xiaopeng",
	}
	for input, want := range cases {
		if got := deriveEmployeeCode(input); got != want {
			t.Fatalf("deriveEmployeeCode(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestParseLMSStoresSupportsBOMAndMapsSourceValues(t *testing.T) {
	csvText := "\ufeffretail_pos_number,country_code,channel,location_desc,status,opn_date,close_date,region_desc,state,city,address,address2,zip,longitude,latitude,str_hour,sqr_ftage,brand_code,sap_site_id,company_code,store_type,chnl_code,chnl_desc,region_code,region_name,district_code,district_name,area_code,area_name,area_manager,district_manager\n" +
		"OC267,CN,RETAIL,SHANGHAI MIXC,OPEN,2021-07-01 00:00:00,2099-12-31 00:00:00,China Region,N/A,Shanghai,Line 1,Line 2,200000,-94.495340,36.596890,Mon-Sun 10:00-22:00,180 sqm,COACH,D369,3000,Physical Store,CN,China,RC_CAFP_CN,China Region,CAFP_CN_ED,Retail-East,Area_Code_CAFP_CN_EA1,East 1,Jean Yuan,Sandy Tang\n"

	result, err := parseLMSStores(strings.NewReader(csvText))
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Records) != 1 {
		t.Fatalf("record count = %d, want 1", len(result.Records))
	}
	got := result.Records[0]
	if got.StoreCode != "OC267" || got.OrganizationCode != "D369" || got.RegionDesc != "China Region" || got.AreaManager != "Jean Yuan" || got.StoreDistrictMgr != "Sandy Tang" {
		t.Fatalf("record = %+v", got)
	}
}

func TestParseWorkdayEmployeesMapsKnownAndOptionalHeaders(t *testing.T) {
	csvText := "EmployeeId,MobileNumber,Name,AliasName,ChineseName,JobLocation,JobPosition,Gender,EmailAddress,Status,StoreCode,Department,Hire Date,TerminationDate,Job Function,IsStoreManagement\n" +
		"000931251,18700001251,Susu Wang,Susu,王散,Shanghai,Sr Sales Associate,F,Lily@coach.com,1,OC267,Sales - Direct,2026/7/2,,VM,Y\n"

	result, err := parseWorkdayEmployees(strings.NewReader(csvText))
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Records) != 1 {
		t.Fatalf("record count = %d, want 1", len(result.Records))
	}
	got := result.Records[0]
	if got.EmployeeID != "000931251" || got.EmployeeCode != "931251" || got.NameZh != "王散" || got.HireDate != "2026/7/2" || got.JobFunction != "VM" || got.IsStoreManagement != "Y" {
		t.Fatalf("record = %+v", got)
	}
}

func TestParseWorkdayEmployeesSkipsAllRowsForDuplicateDerivedCode(t *testing.T) {
	csvText := "EmployeeId,Name\n000931251,First\n931251,Second\n,Missing\n000763514,Valid\n"

	result, err := parseWorkdayEmployees(strings.NewReader(csvText))
	if err != nil {
		t.Fatal(err)
	}
	if result.MissingKeyRows != 1 {
		t.Fatalf("missing key rows = %d, want 1", result.MissingKeyRows)
	}
	if len(result.DuplicateKeys) != 1 || result.DuplicateKeys[0] != "931251" {
		t.Fatalf("duplicate keys = %v, want [931251]", result.DuplicateKeys)
	}
	if len(result.Records) != 1 || result.Records[0].EmployeeCode != "763514" {
		t.Fatalf("records = %+v", result.Records)
	}
}

func TestParseLMSStoresRejectsMissingBusinessKeyHeader(t *testing.T) {
	_, err := parseLMSStores(strings.NewReader("country_code,status\nCN,OPEN\n"))
	if err == nil || !strings.Contains(err.Error(), "retail_pos_number") {
		t.Fatalf("error = %v, want missing retail_pos_number", err)
	}
}
