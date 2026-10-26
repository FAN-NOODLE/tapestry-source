package enterprise

import (
	"errors"
	"gorm.io/gorm"
	enterpriseConfigModel "micro-scrm/internal/orm/model/micro_scrm/enterprise"
)

// 根据企业ID，获取企业config信息
func (e *EnterpriseService) GetConfigListOneEnterprise(enterpriseID int64) (data []enterpriseConfigModel.EnterpriseConfigModel, err error) {
	result := e.Where("enterprise_id = ?", enterpriseID).
		Find(&data)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return data, nil
		}
	}

	return data, result.Error
}
