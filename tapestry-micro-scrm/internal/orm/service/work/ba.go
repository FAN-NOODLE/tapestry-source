package work

import (
	workDB "micro-scrm/internal/orm/model/work"
)

// 智慧导购

// GetAllBAInfoList 获取品牌所有的导购信息
func (work *WorkDBService) GetAllBAInfoList(enterpriseID int64, baType []uint8) (userInfoList []*workDB.StoreWorkerModel, err error) {
	model := work.Model(workDB.StoreWorkerModel{}).Where(map[string]interface{}{
		"enterprise_id": enterpriseID,
		//"worker_status": workDB.SW_WORKER_STATUS_ON,
		"deleted_at":    0,
	})
	//if len(baType) != 0 {
	//	model.Where("type IN (?)", baType)
	//}

	model.Where("type IN (?,?)", workDB.SW_TYPE_ENTERPRISE, workDB.SW_TYPE_ENTERPRISE_MANAGE)
	model.Where("worker_status IN (?,?)", workDB.SW_WORKER_STATUS_ON, workDB.SW_WORKER_STATUS_RESIGN)

	model.Find(&userInfoList).Order("created_at DESC")
	err = model.Error
	return
}

// GetBAInfo 根据BA的ID，获取导购的信息
func (work *WorkDBService) GetBAInfo(enterpriseID int64, memUserID string) (data workDB.StoreWorkerModel, err error) {
	err = work.Model(workDB.StoreWorkerModel{}).Where(map[string]interface{}{
		"enterprise_id": enterpriseID,
		"worker_status": workDB.SW_WORKER_STATUS_ON,
		"deleted_at":    0,
		"mem_uid":       memUserID,
	}).Take(&data).Error

	return
}
func (work *WorkDBService) GetBAInfoListPart(enterpriseID int64, baType []uint8, memUserIDList []string) (userInfoList []*workDB.StoreWorkerModel, err error) {
	model := work.Model(workDB.StoreWorkerModel{}).Where(map[string]interface{}{
		"enterprise_id": enterpriseID,
		//"worker_status": workDB.SW_WORKER_STATUS_ON,
		"deleted_at":    0,
		"mem_uid":       memUserIDList,
	})
	//if len(baType) != 0 {
	//	model.Where("type IN (?)", baType)
	//}

	model.Where("type IN (?,?)", workDB.SW_TYPE_ENTERPRISE, workDB.SW_TYPE_ENTERPRISE_MANAGE)
	model.Where("worker_status IN (?,?)", workDB.SW_WORKER_STATUS_ON, workDB.SW_WORKER_STATUS_RESIGN)

	model.Find(&userInfoList)
	err = model.Error
	return
}

func (work *WorkDBService) GetBAInfoListWithUserIdPart(enterpriseID int64, memUserIDList []string) (userInfoList []*workDB.StoreWorkerModel, err error) {
	model := work.Model(workDB.StoreWorkerModel{}).Where(map[string]interface{}{
		"enterprise_id": enterpriseID,
		"worker_status": workDB.SW_WORKER_STATUS_ON,
		"deleted_at":    0,
		"mem_userid":    memUserIDList,
	}).Where("mem_userid <> '' ").
		Find(&userInfoList)
	err = model.Error
	return
}

// GetPagingAllBAInfoList 分页获取品牌所有的导购信息
func (work *WorkDBService) GetPagingAllBAInfoList(enterpriseID int64, baType []uint8, memUserIDList []string, page, size int) (userInfoList []*workDB.StoreWorkerModel, total int64, err error) {

	model := work.Model(workDB.StoreWorkerModel{}).Where(map[string]interface{}{
		"enterprise_id": enterpriseID,
		//"worker_status": workDB.SW_WORKER_STATUS_ON,
		"deleted_at":    0,
	})
	if len(baType) != 0 {
		model.Where("type IN (?)", baType)
	}
	if len(memUserIDList) > 0 {
		model.Where("mem_uid IN (?)", memUserIDList)
	}
	model.Where("worker_status IN (?,?)", workDB.SW_WORKER_STATUS_ON, workDB.SW_WORKER_STATUS_RESIGN)
	model.Where("mem_userid !=''")
	model.Count(&total)
	if page != 0 && size != 0 {
		model.Select("worker_no,name,mem_userid")
		model.Limit(size).Offset((page - 1) * size)
	}
	err = model.Find(&userInfoList).Order("created_at DESC").Error
	return
}

// GetTypeOFEnterprise 获取品牌的编号
func (work *WorkDBService) GetTypeOFEnterprise(enterpriseID int64) (result string, err error) {
	err = work.Model(workDB.DataLakeConfigModel{}).
		Select("brand_code").
		Where("enterprise_id = ?", enterpriseID).
		Where("deleted_at = ?", workDB.DLC_DELETED).
		Find(&result).Error
	return
}


// GetAllBAList 导购信息
func (work *WorkDBService) GetAllBAList(enterpriseID int64, baType []uint8, memUserIDList []string) (userInfoList []*workDB.StoreWorkerModel, total int64, err error) {

	model := work.Model(workDB.StoreWorkerModel{}).Where(map[string]interface{}{
		"enterprise_id": enterpriseID,
		"worker_status": workDB.SW_WORKER_STATUS_ON,
		"deleted_at":    0,
	})
	if len(baType) != 0 {
		model.Where("type IN (?)", baType)
	}
	if len(memUserIDList) > 0 {
		model.Where("mem_uid IN (?)", memUserIDList)
	}
	model.Where("worker_status IN (?,?)", workDB.SW_WORKER_STATUS_ON, workDB.SW_WORKER_STATUS_RESIGN)
	model.Where("mem_userid !=''")
	model.Count(&total)
	err = model.Find(&userInfoList).Order("created_at DESC").Error
	return
}