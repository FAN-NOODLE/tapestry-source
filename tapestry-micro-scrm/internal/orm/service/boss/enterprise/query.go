package enterprise

import (
	"fmt"
	"micro-scrm/internal/orm/model/boss"
	"micro-scrm/internal/orm/service"
)

func (boss *EnterpriseBossService) GetInfo(param boss.EnterpriseModel) (data boss.EnterpriseModel, err error) {
	err = boss.DB.Where(&param).First(&data).Error

	return
}

func (boss *EnterpriseBossService) GetEnterpriseAdminList(param boss.EnterpriseAdminModel) (data []boss.EnterpriseAdminModel,
	err error) {

	err = boss.DB.Where("username like ?", fmt.Sprintf("%s%s%s", "%", param.Username, "%")).
		Find(&data).Error

	return
}

func (boss *EnterpriseBossService) GetEnterpriseAdminListByOperationUserId(operationUserId []int64) (
	data []boss.
		EnterpriseAdminModel,
	err error) {

	err = boss.DB.Where("id in ?", operationUserId).
		Find(&data).Error

	return
}

func (boss *EnterpriseBossService) GetAllInfo() (data []boss.EnterpriseModel, err error) {
	err = boss.DB.Scopes(service.ScopeCond_EQDelete).
		Find(&data).Error

	return
}
