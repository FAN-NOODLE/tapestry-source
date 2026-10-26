package brand_config

import "micro-scrm/internal/orm/model/micro_scrm/brand_config"

// GetCouponCenterConfig 获取品牌CouponCenter的uuid
func (service *BrandConfigService) GetCouponCenterConfig(enterpriseID int64) (info brand_config.BrandConfigCouponCenterModel, err error) {
	err = service.Where("enterprise_id = ?", enterpriseID).Take(&info).Error
	return
}
