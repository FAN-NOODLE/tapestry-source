package version

import "micro-scrm/internal/orm/model/micro_scrm/version"

// RecordVersionPrompt 记录版本信息展示记录
func (s *VersionService) RecordVersionPrompt(data version.VersionPromptLogModel) error {
	return s.Create(&data).Error
}
