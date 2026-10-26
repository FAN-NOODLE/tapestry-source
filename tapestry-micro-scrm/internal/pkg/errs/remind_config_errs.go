package errs

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"

var (
	// ConfigInfoErr 获取提醒配置
	ConfigInfoErr = errorx.NewErrorMsg(170001, "提醒配置获取失败")
	// LostTimeErr  每日接待时间不能为空
	LostTimeErr = errorx.NewErrorMsg(170002, "每日接待时间不能为空")
	// OvertimeErr  每日接待时间不能为空
	OvertimeErr = errorx.NewErrorMsg(170003, "聊天超时分钟设置有误")
)
