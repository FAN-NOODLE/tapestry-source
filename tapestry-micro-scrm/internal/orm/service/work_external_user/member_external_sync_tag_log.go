package work_external_user

import (
	"micro-scrm/internal/orm/model/micro_scrm/work_external_user"

)

// CreateSyncLogData 根据用户的关系ID，获取用户的uid
func (w *WorkExternalUserService) CreateSyncLogData(data work_external_user.MemberExternalSyncTagLogModel) (work_external_user.MemberExternalSyncTagLogModel, int64, error) {

	result := w.Create(&data)

	return data, result.RowsAffected, result.Error
}
