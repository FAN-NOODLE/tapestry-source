package fission

import (
	"errors"
	"gorm.io/gorm"
	"micro-scrm/internal/orm/model/micro_scrm/fission"
)

func (e *FissionService) GetFissionMpUser(enterpriseID int64, openID string) (isExist bool, data fission.FissionMpUserModel, err error) {
	result := e.Where(map[string]interface{}{
		"enterprise_id": enterpriseID,
		"open_id":       openID,
	}).First(&data)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			isExist = false
			return
		}

		err = result.Error
		return
	}

	isExist = true
	return
}

func (e *FissionService) GetFissionMpUserById(enterpriseID int64, ID string) (isExist bool, data fission.FissionMpUserModel, err error) {
	result := e.Where(map[string]interface{}{
		"enterprise_id": enterpriseID,
		"id":            ID,
	}).First(&data)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			isExist = false
			return
		}

		err = result.Error
		return
	}

	isExist = true
	return
}
