package fission

import "micro-scrm/internal/orm/model/micro_scrm/fission"

// InsertFissionMpUser
func (f *FissionService) InsertFissionMpUser(model fission.FissionMpUserModel) error {
	return f.Create(&model).Error
}
