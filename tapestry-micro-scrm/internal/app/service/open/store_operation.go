package openSvc

import (
	"fmt"
	"strings"

	"micro-scrm/internal/app"
	openDao "micro-scrm/internal/app/dao/open"
	openReq "micro-scrm/internal/app/http/request/open"
	openResp "micro-scrm/internal/app/http/response/open"
	"micro-scrm/internal/config"
	enterpriseModel "micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/pkg/errs"
	pkgKMS "micro-scrm/internal/pkg/kms"
)

const (
	defaultStoreMasterPageSize = 100
	maxStoreMasterPageSize     = 1000
)

var validStoreCountryCodes = map[string]struct{}{
	"CN": {}, "HK": {}, "MO": {}, "TW": {},
}

var validEmployeeStoreScopes = map[string]struct{}{
	"all": {}, "store": {}, "open_store": {},
}

var validEmploymentScopes = map[string]struct{}{
	"active": {}, "inactive": {}, "all": {},
}

type storeMasterRepository interface {
	StoreMasterPage(app.Ctx, int64, openReq.StoreMasterListReq) ([]enterpriseModel.StoreMasterModel, int64, error)
	WechatDepartmentIDs(app.Ctx, int64, []string) (map[string]int64, error)
}

type employeeMasterRepository interface {
	EmployeeMasterPage(app.Ctx, int64, openReq.EmployeeMasterListReq) ([]openDao.EmployeeMasterRow, int64, error)
	EmployeeWeComData(app.Ctx, int64, []string) (map[string]openDao.EmployeeWeComData, error)
}

type StoreOperationService struct {
	repository         storeMasterRepository
	employeeRepository employeeMasterRepository
	kmsClient          pkgKMS.DataKeyClient
	kmsInitErr         error
}

type storeMasterPlainItem struct {
	Response             openResp.StoreMasterResp
	AreaManager          string
	StoreDistrictManager string
}

type employeeMasterPlainItem struct {
	Response     openResp.EmployeeMasterResp
	Name         string
	NameZh       string
	EmployeeID   string
	EmployeeCode string
	WecomUserid  string
	MobileNumber string
	EmailAddress string
}

func NewStoreOperationService() *StoreOperationService {
	client, err := pkgKMS.NewTencentClient(config.CfgData.KMS)
	return &StoreOperationService{
		repository:         openDao.NewStoreOperationDao(),
		employeeRepository: openDao.NewStoreOperationDao(),
		kmsClient:          client,
		kmsInitErr:         err,
	}
}

func newStoreOperationServiceWithRepository(repository storeMasterRepository) *StoreOperationService {
	return &StoreOperationService{repository: repository}
}

func newStoreOperationServiceWithDependencies(repository storeMasterRepository, client pkgKMS.DataKeyClient) *StoreOperationService {
	return &StoreOperationService{repository: repository, kmsClient: client}
}

func newStoreOperationServiceWithEmployeeRepository(repository employeeMasterRepository) *StoreOperationService {
	return &StoreOperationService{employeeRepository: repository}
}

func newStoreOperationServiceWithEmployeeDependencies(repository employeeMasterRepository, client pkgKMS.DataKeyClient) *StoreOperationService {
	return &StoreOperationService{employeeRepository: repository, kmsClient: client}
}

func (s *StoreOperationService) StoreMasterList(c app.Ctx, req openReq.StoreMasterListReq) ([]openResp.StoreMasterResp, int64, error) {
	if s.kmsInitErr != nil {
		return nil, 0, fmt.Errorf("初始化KMS客户端失败: %w", s.kmsInitErr)
	}
	plainItems, total, err := s.loadStoreMasterList(c, req)
	if err != nil {
		return nil, total, err
	}
	if len(plainItems) == 0 {
		return make([]openResp.StoreMasterResp, 0), total, nil
	}
	pageCipher, err := pkgKMS.NewPageCipher(c, s.kmsClient)
	if err != nil {
		return nil, 0, err
	}
	defer pageCipher.Close()

	result := make([]openResp.StoreMasterResp, 0, len(plainItems))
	for _, plainItem := range plainItems {
		item := plainItem.Response
		item.AreaManager, err = pageCipher.EncryptString(plainItem.AreaManager)
		if err != nil {
			return nil, 0, fmt.Errorf("加密门店PI字段失败: %w", err)
		}
		item.StoreDistrictManager, err = pageCipher.EncryptString(plainItem.StoreDistrictManager)
		if err != nil {
			return nil, 0, fmt.Errorf("加密门店PI字段失败: %w", err)
		}
		result = append(result, item)
	}
	return result, total, nil
}

func normalizeStoreMasterRequest(req *openReq.StoreMasterListReq) error {
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Size == 0 {
		req.Size = defaultStoreMasterPageSize
	}
	if req.Page < 1 || req.Size < 1 || req.Size > maxStoreMasterPageSize || req.UpdatedAt < 0 {
		return errs.OpenAPIParamsErr
	}
	if len(req.StoreCountryCodes) == 0 {
		req.StoreCountryCodes = []string{"CN"}
		return nil
	}
	seen := make(map[string]struct{}, len(req.StoreCountryCodes))
	countryCodes := make([]string, 0, len(req.StoreCountryCodes))
	for _, value := range req.StoreCountryCodes {
		countryCode := strings.ToUpper(strings.TrimSpace(value))
		if _, valid := validStoreCountryCodes[countryCode]; !valid {
			return errs.OpenAPIParamsErr
		}
		if _, exists := seen[countryCode]; exists {
			continue
		}
		seen[countryCode] = struct{}{}
		countryCodes = append(countryCodes, countryCode)
	}
	req.StoreCountryCodes = countryCodes
	return nil
}

func normalizeEmployeeMasterRequest(req *openReq.EmployeeMasterListReq) error {
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Size == 0 {
		req.Size = defaultStoreMasterPageSize
	}
	if req.StoreScope == "" {
		req.StoreScope = "all"
	}
	if req.EmploymentScope == "" {
		req.EmploymentScope = "active"
	}
	if req.Page < 1 || req.Size < 1 || req.Size > maxStoreMasterPageSize || req.UpdatedAt < 0 {
		return errs.OpenAPIParamsErr
	}
	if _, valid := validEmployeeStoreScopes[req.StoreScope]; !valid {
		return errs.OpenAPIParamsErr
	}
	if _, valid := validEmploymentScopes[req.EmploymentScope]; !valid {
		return errs.OpenAPIParamsErr
	}
	seen := make(map[string]struct{}, len(req.StoreCountryCodes))
	countryCodes := make([]string, 0, len(req.StoreCountryCodes))
	for _, value := range req.StoreCountryCodes {
		countryCode := strings.ToUpper(strings.TrimSpace(value))
		if _, valid := validStoreCountryCodes[countryCode]; !valid {
			return errs.OpenAPIParamsErr
		}
		if _, exists := seen[countryCode]; exists {
			continue
		}
		seen[countryCode] = struct{}{}
		countryCodes = append(countryCodes, countryCode)
	}
	req.StoreCountryCodes = countryCodes
	return nil
}

func (s *StoreOperationService) loadStoreMasterList(c app.Ctx, req openReq.StoreMasterListReq) ([]storeMasterPlainItem, int64, error) {
	if err := normalizeStoreMasterRequest(&req); err != nil {
		return nil, 0, err
	}
	stores, total, err := s.repository.StoreMasterPage(c, req.EnterpriseID, req)
	if err != nil || len(stores) == 0 {
		return nil, total, err
	}
	storeCodes := make([]string, 0, len(stores))
	for _, store := range stores {
		storeCodes = append(storeCodes, store.StoreCode)
	}
	departmentIDs, err := s.repository.WechatDepartmentIDs(c, req.EnterpriseID, storeCodes)
	if err != nil {
		return nil, 0, err
	}
	result := make([]storeMasterPlainItem, 0, len(stores))
	for _, store := range stores {
		result = append(result, mapStoreMaster(store, departmentIDs[store.StoreCode]))
	}
	return result, total, nil
}

func mapStoreMaster(store enterpriseModel.StoreMasterModel, wechatDepartmentID int64) storeMasterPlainItem {
	return storeMasterPlainItem{
		Response: openResp.StoreMasterResp{
			StoreCode:               store.StoreCode,
			CountryCode:             store.CountryCode,
			Channel:                 store.SalesChannel,
			StoreNameZh:             store.StoreNameZh,
			StoreNameEn:             store.OrganizationName,
			Status:                  store.StoreStatus,
			OpenDate:                store.OpenDate,
			CloseDate:               store.ClosedDate,
			RegionDesc:              store.RegionDesc,
			State:                   store.State,
			City:                    store.City,
			AddressZh:               store.AddressZh,
			Address:                 store.AddressLine1,
			Address2:                store.AddressLine2,
			ZipCode:                 store.ZipCode,
			Phone:                   store.DayPhone,
			Email:                   store.EmailID,
			Longitude:               store.Longitude,
			Latitude:                store.Latitude,
			StoreWorkingHours:       store.StoreWorkingHours,
			SquareFootage:           store.StoreArea,
			AppointmentServiceTypes: store.AppointmentServiceTypes,
			BrandCode:               store.BrandCode,
			OrganizationCode:        store.OrganizationCode,
			PaymentMethods:          store.PaymentMethods,
			CompanyCode:             store.CompanyCode,
			StoreClass:              store.StoreClass,
			WechatDepartmentID:      wechatDepartmentID,
			BusinessOrganization: openResp.StoreBusinessOrganization{
				RegionCode:   store.OrgRegionCode,
				Region:       store.OrgRegionName,
				DistrictCode: store.OrgDistrictCode,
				District:     store.OrgDistrictName,
				AreaCode:     store.OrgAreaCode,
				Area:         store.OrgAreaName,
			},
			StoreDesignConcept: store.StoreDesignConcept,
			DPAP:               store.DPAP,
			UpdatedAt:          store.UpdatedAt,
		},
		AreaManager:          store.AreaManager,
		StoreDistrictManager: store.StoreDistrictMgr,
	}
}

func (s *StoreOperationService) loadEmployeeMasterList(c app.Ctx, req openReq.EmployeeMasterListReq) ([]employeeMasterPlainItem, int64, error) {
	if err := normalizeEmployeeMasterRequest(&req); err != nil {
		return nil, 0, err
	}
	employees, total, err := s.employeeRepository.EmployeeMasterPage(c, req.EnterpriseID, req)
	if err != nil || len(employees) == 0 {
		return nil, total, err
	}
	employeeCodes := make([]string, 0, len(employees))
	for _, employee := range employees {
		employeeCodes = append(employeeCodes, employee.EmployeeCode)
	}
	wecomData, err := s.employeeRepository.EmployeeWeComData(c, req.EnterpriseID, employeeCodes)
	if err != nil {
		return nil, 0, err
	}
	result := make([]employeeMasterPlainItem, 0, len(employees))
	for _, employee := range employees {
		result = append(result, mapEmployeeMaster(employee, wecomData[employee.EmployeeCode]))
	}
	return result, total, nil
}

func (s *StoreOperationService) EmployeeMasterList(c app.Ctx, req openReq.EmployeeMasterListReq) ([]openResp.EmployeeMasterResp, int64, error) {
	if s.kmsInitErr != nil {
		return nil, 0, fmt.Errorf("初始化KMS客户端失败: %w", s.kmsInitErr)
	}
	plainItems, total, err := s.loadEmployeeMasterList(c, req)
	if err != nil {
		return nil, total, err
	}
	if len(plainItems) == 0 {
		return make([]openResp.EmployeeMasterResp, 0), total, nil
	}
	pageCipher, err := pkgKMS.NewPageCipher(c, s.kmsClient)
	if err != nil {
		return nil, 0, err
	}
	defer pageCipher.Close()

	result := make([]openResp.EmployeeMasterResp, 0, len(plainItems))
	for _, plainItem := range plainItems {
		item := plainItem.Response
		if item.Name, err = pageCipher.EncryptString(plainItem.Name); err != nil {
			return nil, 0, employeeEncryptionError(err)
		}
		if item.NameZh, err = pageCipher.EncryptString(plainItem.NameZh); err != nil {
			return nil, 0, employeeEncryptionError(err)
		}
		if item.EmployeeID, err = pageCipher.EncryptString(plainItem.EmployeeID); err != nil {
			return nil, 0, employeeEncryptionError(err)
		}
		if item.EmployeeCode, err = pageCipher.EncryptString(plainItem.EmployeeCode); err != nil {
			return nil, 0, employeeEncryptionError(err)
		}
		if item.WecomUserid, err = pageCipher.EncryptString(plainItem.WecomUserid); err != nil {
			return nil, 0, employeeEncryptionError(err)
		}
		if item.MobileNumber, err = pageCipher.EncryptString(plainItem.MobileNumber); err != nil {
			return nil, 0, employeeEncryptionError(err)
		}
		if item.EmailAddress, err = pageCipher.EncryptString(plainItem.EmailAddress); err != nil {
			return nil, 0, employeeEncryptionError(err)
		}
		result = append(result, item)
	}
	return result, total, nil
}

func employeeEncryptionError(err error) error {
	return fmt.Errorf("加密员工PI字段失败: %w", err)
}

func mapEmployeeMaster(employee openDao.EmployeeMasterRow, wecom openDao.EmployeeWeComData) employeeMasterPlainItem {
	changes := make([]openResp.DepartmentChange, 0, len(wecom.DepartmentChanges))
	for _, change := range wecom.DepartmentChanges {
		changes = append(changes, openResp.DepartmentChange{
			EventType:          change.EventType,
			EventTime:          change.EventTime,
			WechatDepartmentID: change.WechatDepartmentID,
			DepartmentName:     change.DepartmentName,
		})
	}
	return employeeMasterPlainItem{
		Response: openResp.EmployeeMasterResp{
			Gender:             employee.Gender,
			StoreCode:          employee.StoreCode,
			Department:         employee.Department,
			OrganizationCode:   employee.OrganizationCode,
			Status:             employee.Status,
			JobPosition:        employee.JobPosition,
			HireDate:           employee.HireDate,
			TerminationDate:    employee.TerminationDate,
			WecomJoinedAt:      wecom.WecomJoinedAt,
			VerificationStatus: wecom.VerificationStatus,
			WecomLeftAt:        wecom.WecomLeftAt,
			DepartmentChanges:  changes,
			JobFunction:        employee.JobFunction,
			IsStoreManagement:  employee.IsStoreManagement,
			UpdatedAt:          employee.UpdatedAt,
		},
		Name:         employee.Name,
		NameZh:       employee.NameZh,
		EmployeeID:   employee.EmployeeID,
		EmployeeCode: employee.EmployeeCode,
		WecomUserid:  wecom.WecomUserid,
		MobileNumber: employee.MobileNumber,
		EmailAddress: employee.EmailAddress,
	}
}
