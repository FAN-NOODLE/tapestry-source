package storeOperationJob

import (
	"strings"
	"time"

	"micro-scrm/internal/app"
	dashboardModel "micro-scrm/internal/app/models/dashboard"
	departmentsModel "micro-scrm/internal/app/models/departments"
	"micro-scrm/internal/orm/model"
	enterpriseModel "micro-scrm/internal/orm/model/micro_scrm_enterprise"
	"micro-scrm/internal/pkg/db/gormx"

	"gorm.io/gorm"
)

type dashboardAreaRepository interface {
	SyncDashboardAreas(c app.Ctx, enterpriseID int64, rootThirdID int, now time.Time) (dashboardAreaSyncStats, []string, error)
}

type gormDashboardAreaRepository struct {
	db *gorm.DB
}

func newDashboardAreaRepository() *gormDashboardAreaRepository {
	return &gormDashboardAreaRepository{db: gormx.GetInstance()}
}

func (repository *gormDashboardAreaRepository) SyncDashboardAreas(c app.Ctx, enterpriseID int64, rootThirdID int, now time.Time) (stats dashboardAreaSyncStats, skipped []string, err error) {
	var stores []enterpriseModel.StoreMasterModel
	err = repository.db.WithContext(c).
		Scopes(model.BranchTable(&enterpriseModel.StoreMasterModel{}, enterpriseID)).
		Where("deleted_at = ?", 0).
		Find(&stores).Error
	if err != nil {
		return stats, nil, err
	}

	var departments []departmentsModel.DepartmentsModel
	err = repository.db.WithContext(c).
		Model(&departmentsModel.DepartmentsModel{}).
		Where("enterprise_id = ?", enterpriseID).
		Where("third_parent_id = ?", rootThirdID).
		Where("deleted_at = ?", 0).
		Find(&departments).Error
	if err != nil {
		return stats, nil, err
	}
	storeThirdIDs := make(map[string]int, len(departments))
	for _, department := range departments {
		name := strings.TrimSpace(department.Name)
		if name != "" {
			storeThirdIDs[name] = department.ThirdID
		}
	}

	hierarchyStores, closedAtByStore, protectedStoreCodes := partitionDashboardAreaStores(stores, now)
	desired, skipped, buildErr := buildDashboardAreaNodes(hierarchyStores, storeThirdIDs, now)
	if buildErr != nil {
		return stats, skipped, buildErr
	}
	for _, storeCode := range skipped {
		if storeCode != "" {
			protectedStoreCodes[storeCode] = struct{}{}
		}
	}
	err = repository.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		var existing []dashboardModel.DashboardAreaModel
		if queryErr := tx.
			Scopes(model.BranchTable(&dashboardModel.DashboardAreaModel{}, enterpriseID)).
			Unscoped().
			Find(&existing).Error; queryErr != nil {
			return queryErr
		}

		plan, planErr := planDashboardAreaSync(existing, desired, protectedStoreCodes, closedAtByStore, len(hierarchyStores), now.Unix())
		if planErr != nil {
			return planErr
		}
		areaTable := func(db *gorm.DB) *gorm.DB {
			return db.Scopes(model.BranchTable(&dashboardModel.DashboardAreaModel{}, enterpriseID)).Unscoped()
		}
		if len(plan.Creates) > 0 {
			if createErr := areaTable(tx.Session(&gorm.Session{SkipHooks: true})).CreateInBatches(&plan.Creates, 200).Error; createErr != nil {
				return createErr
			}
		}
		for id, values := range plan.Updates {
			if updateErr := areaTable(tx).Where("id = ?", id).UpdateColumns(values).Error; updateErr != nil {
				return updateErr
			}
		}
		if len(plan.Invalidates) > 0 {
			if updateErr := areaTable(tx).
				Where("id IN ?", plan.Invalidates).
				UpdateColumns(map[string]interface{}{"deleted_at": now.Unix(), "updated_at": now.Unix()}).Error; updateErr != nil {
				return updateErr
			}
		}
		stats = plan.Stats
		stats.Skipped = len(skipped)
		return nil
	})
	return stats, skipped, err
}
