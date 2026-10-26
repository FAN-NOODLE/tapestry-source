package version

import (
	"errors"
	"gorm.io/gorm"
	"micro-scrm/internal/orm/model/micro_scrm/version"
	"micro-scrm/internal/orm/service"
)

// GetCurrentVersionInfo 获取最新版本信息
func (s *VersionService) GetCurrentVersionInfo() (result version.VersionRecordModel, err error) {
	err = s.Scopes(service.ScopeCond_EQDelete).
		Order("created_at desc").
		Take(&result).Error

	return
}

// IsHasLatestPromptRecord 判断用户是否已经有的版本信息提示记录
func (s *VersionService) IsHasLatestPromptRecord(cond version.VersionPromptLogModel) (bool, error) {
	var result version.VersionPromptLogModel
	err := s.Where(&cond).First(&result).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
