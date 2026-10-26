package auth

import "micro-scrm/internal/orm/model/micro_scrm"

// GetOpenPlatformInfo 获取对接接口的secret
func (auth *AuthService) GetOpenPlatformInfo(clientID string) (data micro_scrm.OpenPlatformConfigModel, err error) {
	err = auth.DB.Where("client_id = ?", clientID).
		First(&data).Error

	return
}
