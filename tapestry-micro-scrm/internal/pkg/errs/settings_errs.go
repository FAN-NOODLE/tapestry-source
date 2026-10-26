package errs

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"

var (
	// SettingsInfoErr 获取配置失败
	SettingsInfoErr = errorx.NewErrorMsg(170006, "获取配置失败")
	// GenKeyErr 生成密钥对失败
	GenKeyErr = errorx.NewErrorMsg(170006, "密钥对生成失败")
)
