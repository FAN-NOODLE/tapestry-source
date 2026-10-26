package plan_task

import (
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/orm/service"
)

// sftp配置类型:1 计划任务 2 ba业绩看板 3 社群营销数据
const (
	SFTP_TYPE_PLAN_TASK      = 1 + iota //计划任务
	SFTP_TYPE_BA_KANBAN                 //ba业绩看板
	SFTP_TYPE_COMMUNITY_DATE            //社群营销数据
)

// GetBrandSFTPConfig 获取配置了计划任务功能的SFTP信息
func (p *PlanTaskService) GetBrandSFTPConfig(sftpType int) (data []micro_scrm.BrandConfigSftpModel, err error) {
	err = p.Scopes(service.ScopeCond_EQDelete).Where("sftp_type=?", sftpType).
		Find(&data).Error

	return
}
