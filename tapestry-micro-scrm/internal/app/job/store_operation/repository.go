package storeOperationJob

import (
	"micro-scrm/internal/app"
	"micro-scrm/internal/orm/model"
	enterpriseModel "micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/pkg/db/gormx"

	"gorm.io/gorm"
)

type syncStats struct {
	Inserted  int
	Updated   int
	Unchanged int
}

type masterDataRepository interface {
	SyncStores(c app.Ctx, enterpriseID int64, records []storeSourceRecord, now int64) (syncStats, error)
	SyncEmployees(c app.Ctx, enterpriseID int64, records []employeeSourceRecord, now int64) (syncStats, error)
}

type gormMasterDataRepository struct {
	db    *gorm.DB
	newID func() int64
}

func newMasterDataRepository() *gormMasterDataRepository {
	return &gormMasterDataRepository{db: gormx.GetInstance(), newID: app.UID}
}

func (repository *gormMasterDataRepository) SyncStores(c app.Ctx, enterpriseID int64, records []storeSourceRecord, now int64) (stats syncStats, err error) {
	if len(records) == 0 {
		return stats, nil
	}
	err = repository.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		codes := make([]string, 0, len(records))
		for _, record := range records {
			codes = append(codes, record.StoreCode)
		}

		var existing []enterpriseModel.StoreMasterModel
		if queryErr := tx.Scopes(model.BranchTable(&enterpriseModel.StoreMasterModel{}, enterpriseID)).
			Where("store_code IN ?", codes).
			Find(&existing).Error; queryErr != nil {
			return queryErr
		}
		existingByCode := make(map[string]enterpriseModel.StoreMasterModel, len(existing))
		for _, item := range existing {
			existingByCode[item.StoreCode] = item
		}

		inserts := make([]enterpriseModel.StoreMasterModel, 0)
		for _, record := range records {
			item, found := existingByCode[record.StoreCode]
			if !found {
				inserts = append(inserts, newStoreModel(record, repository.newID(), now))
				stats.Inserted++
				continue
			}
			if item.DeletedAt == 0 && !storeSourceChanged(item, record) {
				stats.Unchanged++
				continue
			}
			if updateErr := tx.Session(&gorm.Session{SkipHooks: true}).
				Scopes(model.BranchTable(&enterpriseModel.StoreMasterModel{}, enterpriseID)).
				Where("id = ?", item.ID).
				Updates(storeUpdateValues(record, now)).Error; updateErr != nil {
				return updateErr
			}
			stats.Updated++
		}

		if len(inserts) > 0 {
			if createErr := tx.Session(&gorm.Session{SkipHooks: true}).
				Scopes(model.BranchTable(&enterpriseModel.StoreMasterModel{}, enterpriseID)).
				CreateInBatches(&inserts, 500).Error; createErr != nil {
				return createErr
			}
		}
		return nil
	})
	return stats, err
}

func (repository *gormMasterDataRepository) SyncEmployees(c app.Ctx, enterpriseID int64, records []employeeSourceRecord, now int64) (stats syncStats, err error) {
	if len(records) == 0 {
		return stats, nil
	}
	err = repository.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		storeCodes, queryErr := activeEmployeeStoreCodes(tx, enterpriseID, records)
		if queryErr != nil {
			return queryErr
		}
		applyEmployeeStoreDepartments(records, storeCodes)

		codes := make([]string, 0, len(records))
		for _, record := range records {
			codes = append(codes, record.EmployeeCode)
		}

		var existing []enterpriseModel.EmployeeMasterModel
		if queryErr := tx.Scopes(model.BranchTable(&enterpriseModel.EmployeeMasterModel{}, enterpriseID)).
			Where("employee_code IN ?", codes).
			Find(&existing).Error; queryErr != nil {
			return queryErr
		}
		existingByCode := make(map[string]enterpriseModel.EmployeeMasterModel, len(existing))
		for _, item := range existing {
			existingByCode[item.EmployeeCode] = item
		}

		inserts := make([]enterpriseModel.EmployeeMasterModel, 0)
		for _, record := range records {
			item, found := existingByCode[record.EmployeeCode]
			if !found {
				inserts = append(inserts, newEmployeeModel(record, repository.newID(), now))
				stats.Inserted++
				continue
			}
			if item.DeletedAt == 0 && !employeeSourceChanged(item, record) {
				stats.Unchanged++
				continue
			}
			if updateErr := tx.Session(&gorm.Session{SkipHooks: true}).
				Scopes(model.BranchTable(&enterpriseModel.EmployeeMasterModel{}, enterpriseID)).
				Where("id = ?", item.ID).
				Updates(employeeUpdateValues(record, now)).Error; updateErr != nil {
				return updateErr
			}
			stats.Updated++
		}

		if len(inserts) > 0 {
			if createErr := tx.Session(&gorm.Session{SkipHooks: true}).
				Scopes(model.BranchTable(&enterpriseModel.EmployeeMasterModel{}, enterpriseID)).
				CreateInBatches(&inserts, 500).Error; createErr != nil {
				return createErr
			}
		}
		return nil
	})
	return stats, err
}

func activeEmployeeStoreCodes(tx *gorm.DB, enterpriseID int64, records []employeeSourceRecord) (map[string]struct{}, error) {
	storeCodes := make([]string, 0, len(records))
	seen := make(map[string]struct{}, len(records))
	for _, record := range records {
		if record.StoreCode == "" {
			continue
		}
		if _, exists := seen[record.StoreCode]; exists {
			continue
		}
		seen[record.StoreCode] = struct{}{}
		storeCodes = append(storeCodes, record.StoreCode)
	}
	activeCodes := make(map[string]struct{}, len(storeCodes))
	if len(storeCodes) == 0 {
		return activeCodes, nil
	}
	var matchedCodes []string
	err := tx.Scopes(model.BranchTable(&enterpriseModel.StoreMasterModel{}, enterpriseID)).
		Where("deleted_at = ?", 0).
		Where("store_code IN ?", storeCodes).
		Pluck("store_code", &matchedCodes).Error
	if err != nil {
		return nil, err
	}
	for _, storeCode := range matchedCodes {
		activeCodes[storeCode] = struct{}{}
	}
	return activeCodes, nil
}

func applyEmployeeStoreDepartments(records []employeeSourceRecord, activeStoreCodes map[string]struct{}) {
	for index := range records {
		if _, matched := activeStoreCodes[records[index].StoreCode]; matched {
			records[index].Department = records[index].StoreCode
		}
	}
}

func storeSourceChanged(existing enterpriseModel.StoreMasterModel, incoming storeSourceRecord) bool {
	return storeSourceFromModel(existing) != incoming
}

func employeeSourceChanged(existing enterpriseModel.EmployeeMasterModel, incoming employeeSourceRecord) bool {
	return employeeSourceFromModel(existing) != incoming
}

func storeSourceFromModel(item enterpriseModel.StoreMasterModel) storeSourceRecord {
	return storeSourceRecord{
		StoreCode:         item.StoreCode,
		CountryCode:       item.CountryCode,
		SalesChannel:      item.SalesChannel,
		OrganizationName:  item.OrganizationName,
		StoreStatus:       item.StoreStatus,
		OpenDate:          item.OpenDate,
		ClosedDate:        item.ClosedDate,
		RegionDesc:        item.RegionDesc,
		State:             item.State,
		City:              item.City,
		AddressLine1:      item.AddressLine1,
		AddressLine2:      item.AddressLine2,
		ZipCode:           item.ZipCode,
		Longitude:         item.Longitude,
		Latitude:          item.Latitude,
		StoreWorkingHours: item.StoreWorkingHours,
		BrandCode:         item.BrandCode,
		OrganizationCode:  item.OrganizationCode,
		CompanyCode:       item.CompanyCode,
		StoreType:         item.StoreType,
		OrgChannelCode:    item.OrgChannelCode,
		OrgChannelName:    item.OrgChannelName,
		OrgRegionCode:     item.OrgRegionCode,
		OrgRegionName:     item.OrgRegionName,
		OrgDistrictCode:   item.OrgDistrictCode,
		OrgDistrictName:   item.OrgDistrictName,
		OrgAreaCode:       item.OrgAreaCode,
		OrgAreaName:       item.OrgAreaName,
		AreaManager:       item.AreaManager,
		StoreDistrictMgr:  item.StoreDistrictMgr,
	}
}

func employeeSourceFromModel(item enterpriseModel.EmployeeMasterModel) employeeSourceRecord {
	return employeeSourceRecord{
		EmployeeID:        item.EmployeeID,
		EmployeeCode:      item.EmployeeCode,
		MobileNumber:      item.MobileNumber,
		Name:              item.Name,
		AliasName:         item.AliasName,
		NameZh:            item.NameZh,
		JobLocation:       item.JobLocation,
		JobPosition:       item.JobPosition,
		Gender:            item.Gender,
		EmailAddress:      item.EmailAddress,
		Status:            item.Status,
		StoreCode:         item.StoreCode,
		Department:        item.Department,
		HireDate:          item.HireDate,
		TerminationDate:   item.TerminationDate,
		JobFunction:       item.JobFunction,
		IsStoreManagement: item.IsStoreManagement,
	}
}

func storeUpdateValues(record storeSourceRecord, now int64) map[string]interface{} {
	return map[string]interface{}{
		"country_code":        record.CountryCode,
		"sales_channel":       record.SalesChannel,
		"organization_name":   record.OrganizationName,
		"store_status":        record.StoreStatus,
		"open_date":           record.OpenDate,
		"closed_date":         record.ClosedDate,
		"region_desc":         record.RegionDesc,
		"state":               record.State,
		"city":                record.City,
		"address_line1":       record.AddressLine1,
		"address_line2":       record.AddressLine2,
		"zip_code":            record.ZipCode,
		"longitude":           record.Longitude,
		"latitude":            record.Latitude,
		"store_working_hours": record.StoreWorkingHours,
		"brand_code":          record.BrandCode,
		"organization_code":   record.OrganizationCode,
		"company_code":        record.CompanyCode,
		"store_type":          record.StoreType,
		"org_channel_code":    record.OrgChannelCode,
		"org_channel_name":    record.OrgChannelName,
		"org_region_code":     record.OrgRegionCode,
		"org_region_name":     record.OrgRegionName,
		"org_district_code":   record.OrgDistrictCode,
		"org_district_name":   record.OrgDistrictName,
		"org_area_code":       record.OrgAreaCode,
		"org_area_name":       record.OrgAreaName,
		"area_manager":        record.AreaManager,
		"store_district_mgr":  record.StoreDistrictMgr,
		"deleted_at":          0,
		"updated_at":          now,
	}
}

func employeeUpdateValues(record employeeSourceRecord, now int64) map[string]interface{} {
	return map[string]interface{}{
		"employee_id":         record.EmployeeID,
		"mobile_number":       record.MobileNumber,
		"name":                record.Name,
		"alias_name":          record.AliasName,
		"name_zh":             record.NameZh,
		"job_location":        record.JobLocation,
		"job_position":        record.JobPosition,
		"gender":              record.Gender,
		"email_address":       record.EmailAddress,
		"status":              record.Status,
		"store_code":          record.StoreCode,
		"department":          record.Department,
		"hire_date":           record.HireDate,
		"termination_date":    record.TerminationDate,
		"job_function":        record.JobFunction,
		"is_store_management": record.IsStoreManagement,
		"deleted_at":          0,
		"updated_at":          now,
	}
}

func newStoreModel(record storeSourceRecord, id, now int64) enterpriseModel.StoreMasterModel {
	item := enterpriseModel.StoreMasterModel{}
	item.ID = id
	item.CreatedAt = now
	item.UpdatedAt = now
	item.StoreCode = record.StoreCode
	applyStoreSource(&item, record)
	return item
}

func newEmployeeModel(record employeeSourceRecord, id, now int64) enterpriseModel.EmployeeMasterModel {
	item := enterpriseModel.EmployeeMasterModel{}
	item.ID = id
	item.CreatedAt = now
	item.UpdatedAt = now
	item.EmployeeCode = record.EmployeeCode
	applyEmployeeSource(&item, record)
	return item
}

func applyStoreSource(item *enterpriseModel.StoreMasterModel, record storeSourceRecord) {
	item.CountryCode = record.CountryCode
	item.SalesChannel = record.SalesChannel
	item.OrganizationName = record.OrganizationName
	item.StoreStatus = record.StoreStatus
	item.OpenDate = record.OpenDate
	item.ClosedDate = record.ClosedDate
	item.RegionDesc = record.RegionDesc
	item.State = record.State
	item.City = record.City
	item.AddressLine1 = record.AddressLine1
	item.AddressLine2 = record.AddressLine2
	item.ZipCode = record.ZipCode
	item.Longitude = record.Longitude
	item.Latitude = record.Latitude
	item.StoreWorkingHours = record.StoreWorkingHours
	item.BrandCode = record.BrandCode
	item.OrganizationCode = record.OrganizationCode
	item.CompanyCode = record.CompanyCode
	item.StoreType = record.StoreType
	item.OrgChannelCode = record.OrgChannelCode
	item.OrgChannelName = record.OrgChannelName
	item.OrgRegionCode = record.OrgRegionCode
	item.OrgRegionName = record.OrgRegionName
	item.OrgDistrictCode = record.OrgDistrictCode
	item.OrgDistrictName = record.OrgDistrictName
	item.OrgAreaCode = record.OrgAreaCode
	item.OrgAreaName = record.OrgAreaName
	item.AreaManager = record.AreaManager
	item.StoreDistrictMgr = record.StoreDistrictMgr
}

func applyEmployeeSource(item *enterpriseModel.EmployeeMasterModel, record employeeSourceRecord) {
	item.EmployeeID = record.EmployeeID
	item.MobileNumber = record.MobileNumber
	item.Name = record.Name
	item.AliasName = record.AliasName
	item.NameZh = record.NameZh
	item.JobLocation = record.JobLocation
	item.JobPosition = record.JobPosition
	item.Gender = record.Gender
	item.EmailAddress = record.EmailAddress
	item.Status = record.Status
	item.StoreCode = record.StoreCode
	item.Department = record.Department
	item.HireDate = record.HireDate
	item.TerminationDate = record.TerminationDate
	item.JobFunction = record.JobFunction
	item.IsStoreManagement = record.IsStoreManagement
}
