package tag

import "micro-scrm/internal/orm/model/micro_scrm"

func (t *TagService) GetSyncData(enterpriseID int64, limit int) (data []*micro_scrm.DLOneBindLogModel, err error) {
	err = t.Where("status = ?", 0).Where("enterprise_id = ?", enterpriseID).Limit(limit).Find(&data).Error

	return
}

func (t *TagService) SyncStatus(data micro_scrm.DLOneBindLogModel) error {
	return t.Save(&data).Error
}
