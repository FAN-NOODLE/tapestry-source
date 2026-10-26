package overtimeJob

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/modules/overtime"
	overtimeMysql "micro-scrm/internal/modules/overtime/respository/mysql"
	"micro-scrm/internal/modules/settings"
	settingRepository "micro-scrm/internal/modules/settings/respository/mysql"
	"micro-scrm/internal/pkg/xjob"
)

// OvertimeJob 会话存档超时检测
func OvertimeJob(c ctx.Context, jobHandler *xjob.JobHandler) error {
	overtimeService := overtime.New(overtimeMysql.New())

	// 遍历chat_settings表的企业开始处理
	settingService := settings.New(settingRepository.New())
	for _, enterpriseSetting := range settingService.EnterpriseList(c) {

		jobHandler.SpecDo(c, fmt.Sprintf("%d", enterpriseSetting.EnterpriseId), func() {
			overtimeService.GetOvertimeOne(c, enterpriseSetting)
		})
	}

	return nil
}
