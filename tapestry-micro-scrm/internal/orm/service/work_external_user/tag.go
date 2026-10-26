package work_external_user

import (
	"errors"
	"gorm.io/gorm"
	"micro-scrm/internal/orm/model/micro_scrm/work_external_user"
	"micro-scrm/internal/orm/service"
)

func (w *WorkExternalUserService) InsertUserTag(dataList []work_external_user.WorkExternalUserTagModel) error {
	result := w.Create(&dataList)

	return result.Error
}

func (w *WorkExternalUserService) IsExistUserTag(param work_external_user.WorkExternalUserTagModel) (data work_external_user.WorkExternalUserTagModel, isExist bool, err error) {
	err = w.Scopes(service.ScopeCond_EQDelete).Where(&param).First(&data).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
		return
	}

	if err != nil {
		return
	}

	isExist = true
	return
}
