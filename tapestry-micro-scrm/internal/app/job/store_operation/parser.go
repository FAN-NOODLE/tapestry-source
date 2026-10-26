package storeOperationJob

import (
	"encoding/csv"
	"fmt"
	"io"
	"sort"
	"strings"
)

type storeSourceRecord struct {
	StoreCode         string
	CountryCode       string
	SalesChannel      string
	OrganizationName  string
	StoreStatus       string
	OpenDate          string
	ClosedDate        string
	RegionDesc        string
	State             string
	City              string
	AddressLine1      string
	AddressLine2      string
	ZipCode           string
	Longitude         string
	Latitude          string
	StoreWorkingHours string
	BrandCode         string
	OrganizationCode  string
	CompanyCode       string
	StoreType         string
	OrgChannelCode    string
	OrgChannelName    string
	OrgRegionCode     string
	OrgRegionName     string
	OrgDistrictCode   string
	OrgDistrictName   string
	OrgAreaCode       string
	OrgAreaName       string
	AreaManager       string
	StoreDistrictMgr  string
}

type employeeSourceRecord struct {
	EmployeeID        string
	EmployeeCode      string
	MobileNumber      string
	Name              string
	AliasName         string
	NameZh            string
	JobLocation       string
	JobPosition       string
	Gender            string
	EmailAddress      string
	Status            string
	StoreCode         string
	Department        string
	HireDate          string
	TerminationDate   string
	JobFunction       string
	IsStoreManagement string
}

type storeParseResult struct {
	Records        []storeSourceRecord
	MissingKeyRows int
	DuplicateKeys  []string
}

type employeeParseResult struct {
	Records        []employeeSourceRecord
	MissingKeyRows int
	DuplicateKeys  []string
}

func deriveEmployeeCode(employeeID string) string {
	if strings.HasPrefix(employeeID, "000") {
		return strings.TrimPrefix(employeeID, "000")
	}
	return employeeID
}

func parseLMSStores(r io.Reader) (storeParseResult, error) {
	headers, rows, err := readCSV(r)
	if err != nil {
		return storeParseResult{}, err
	}
	index := newHeaderIndex(headers)
	if _, ok := index.find("retail_pos_number"); !ok {
		return storeParseResult{}, fmt.Errorf("LMS文件缺少必需字段 retail_pos_number")
	}

	records := make([]storeSourceRecord, 0, len(rows))
	keyCounts := make(map[string]int, len(rows))
	missingKeyRows := 0
	for _, row := range rows {
		storeCode := strings.TrimSpace(index.value(row, "retail_pos_number"))
		if storeCode == "" {
			missingKeyRows++
			continue
		}
		record := storeSourceRecord{
			StoreCode:         storeCode,
			CountryCode:       index.value(row, "country_code"),
			SalesChannel:      index.value(row, "channel"),
			OrganizationName:  index.value(row, "location_desc"),
			StoreStatus:       index.value(row, "status"),
			OpenDate:          index.value(row, "opn_date"),
			ClosedDate:        index.value(row, "close_date"),
			RegionDesc:        index.value(row, "region_desc"),
			State:             index.value(row, "state"),
			City:              index.value(row, "city"),
			AddressLine1:      index.value(row, "address"),
			AddressLine2:      index.value(row, "address2"),
			ZipCode:           index.value(row, "zip"),
			Longitude:         index.value(row, "longitude"),
			Latitude:          index.value(row, "latitude"),
			StoreWorkingHours: index.value(row, "str_hour"),
			BrandCode:         index.value(row, "brand_code"),
			OrganizationCode:  index.value(row, "sap_site_id"),
			CompanyCode:       index.value(row, "company_code"),
			StoreType:         index.value(row, "store_type"),
			OrgChannelCode:    index.value(row, "chnl_code"),
			OrgChannelName:    index.value(row, "chnl_desc"),
			OrgRegionCode:     index.value(row, "region_code"),
			OrgRegionName:     index.value(row, "region_name"),
			OrgDistrictCode:   index.value(row, "district_code"),
			OrgDistrictName:   index.value(row, "district_name"),
			OrgAreaCode:       index.value(row, "area_code"),
			OrgAreaName:       index.value(row, "area_name"),
			AreaManager:       index.value(row, "area_manager"),
			StoreDistrictMgr:  index.value(row, "district_manager"),
		}
		records = append(records, record)
		keyCounts[storeCode]++
	}

	valid, duplicates := filterDuplicateStores(records, keyCounts)
	return storeParseResult{Records: valid, MissingKeyRows: missingKeyRows, DuplicateKeys: duplicates}, nil
}

func parseWorkdayEmployees(r io.Reader) (employeeParseResult, error) {
	headers, rows, err := readCSV(r)
	if err != nil {
		return employeeParseResult{}, err
	}
	index := newHeaderIndex(headers)
	if _, ok := index.find("EmployeeId"); !ok {
		return employeeParseResult{}, fmt.Errorf("Workday文件缺少必需字段 EmployeeId")
	}

	records := make([]employeeSourceRecord, 0, len(rows))
	keyCounts := make(map[string]int, len(rows))
	missingKeyRows := 0
	for _, row := range rows {
		employeeID := strings.TrimSpace(index.value(row, "EmployeeId"))
		employeeCode := deriveEmployeeCode(employeeID)
		if employeeCode == "" {
			missingKeyRows++
			continue
		}
		record := employeeSourceRecord{
			EmployeeID:        employeeID,
			EmployeeCode:      employeeCode,
			MobileNumber:      index.value(row, "MobileNumber"),
			Name:              index.value(row, "Name"),
			AliasName:         index.value(row, "AliasName"),
			NameZh:            index.value(row, "NameZh", "ChineseName", "name_zh"),
			JobLocation:       index.value(row, "JobLocation"),
			JobPosition:       index.value(row, "JobPosition"),
			Gender:            index.value(row, "Gender"),
			EmailAddress:      index.value(row, "EmailAddress"),
			Status:            index.value(row, "Status"),
			StoreCode:         index.value(row, "StoreCode"),
			Department:        index.value(row, "Department"),
			HireDate:          index.value(row, "HireDate", "Hire Date", "hire_date"),
			TerminationDate:   index.value(row, "TerminationDate", "Termination Date", "termination_date"),
			JobFunction:       index.value(row, "JobFunction", "Job Function", "job_function"),
			IsStoreManagement: index.value(row, "IsStoreManagement", "Is Store Management", "is_store_management"),
		}
		records = append(records, record)
		keyCounts[employeeCode]++
	}

	valid, duplicates := filterDuplicateEmployees(records, keyCounts)
	return employeeParseResult{Records: valid, MissingKeyRows: missingKeyRows, DuplicateKeys: duplicates}, nil
}

func readCSV(r io.Reader) ([]string, [][]string, error) {
	reader := csv.NewReader(r)
	reader.FieldsPerRecord = -1
	headers, err := reader.Read()
	if err != nil {
		if err == io.EOF {
			return nil, nil, fmt.Errorf("CSV文件为空")
		}
		return nil, nil, fmt.Errorf("读取CSV表头失败: %w", err)
	}

	rows := make([][]string, 0)
	for rowNumber := 2; ; rowNumber++ {
		row, readErr := reader.Read()
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return nil, nil, fmt.Errorf("读取CSV第%d行失败: %w", rowNumber, readErr)
		}
		rows = append(rows, row)
	}
	return headers, rows, nil
}

type headerIndex map[string]int

func newHeaderIndex(headers []string) headerIndex {
	index := make(headerIndex, len(headers))
	for position, header := range headers {
		normalized := normalizeHeader(header)
		if _, exists := index[normalized]; !exists {
			index[normalized] = position
		}
	}
	return index
}

func normalizeHeader(header string) string {
	header = strings.TrimPrefix(strings.TrimSpace(header), "\ufeff")
	header = strings.ToLower(header)
	return strings.Map(func(r rune) rune {
		switch r {
		case ' ', '\t', '\r', '\n', '_', '-':
			return -1
		default:
			return r
		}
	}, header)
}

func (index headerIndex) find(aliases ...string) (int, bool) {
	for _, alias := range aliases {
		position, ok := index[normalizeHeader(alias)]
		if ok {
			return position, true
		}
	}
	return 0, false
}

func (index headerIndex) value(row []string, aliases ...string) string {
	position, ok := index.find(aliases...)
	if !ok || position >= len(row) {
		return ""
	}
	return row[position]
}

func filterDuplicateStores(records []storeSourceRecord, counts map[string]int) ([]storeSourceRecord, []string) {
	valid := make([]storeSourceRecord, 0, len(records))
	duplicates := make([]string, 0)
	for key, count := range counts {
		if count > 1 {
			duplicates = append(duplicates, key)
		}
	}
	sort.Strings(duplicates)
	for _, record := range records {
		if counts[record.StoreCode] == 1 {
			valid = append(valid, record)
		}
	}
	return valid, duplicates
}

func filterDuplicateEmployees(records []employeeSourceRecord, counts map[string]int) ([]employeeSourceRecord, []string) {
	valid := make([]employeeSourceRecord, 0, len(records))
	duplicates := make([]string, 0)
	for key, count := range counts {
		if count > 1 {
			duplicates = append(duplicates, key)
		}
	}
	sort.Strings(duplicates)
	for _, record := range records {
		if counts[record.EmployeeCode] == 1 {
			valid = append(valid, record)
		}
	}
	return valid, duplicates
}
