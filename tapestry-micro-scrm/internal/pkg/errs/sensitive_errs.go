package errs

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"

var (
	// SensitiveTextEmptyErr 敏感词不能为空
	SensitiveTextEmptyErr = errorx.NewErrorMsg(170020, "敏感词不能为空")
	// SensitiveTextRepeatErr 敏感词不能重复
	SensitiveTextRepeatErr = errorx.NewErrorMsg(170021, "敏感词不能重复")
	// SensitiveCartIsExistsErr 敏感词名片类型已存在
	SensitiveCartIsExistsErr = errorx.NewErrorMsg(170022, "名片类型只能添加一次")
)
