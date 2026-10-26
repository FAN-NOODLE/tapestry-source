package errs

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"

// ErrorStatus 错误状态码
var ErrorStatus = 401

// DefaultErrorCode 默认码
var DefaultErrorCode = 110000

// DefaultErrorMsg 默认错误信息
var DefaultErrorMsg = "参数错误"
var (
	// ParamsErr 参数错误
	ParamsErr            = errorx.NewErrorMsg(DefaultErrorCode, DefaultErrorMsg)
	SourceTypeErr        = errorx.NewErrorMsg(1100001, "素材类型不合法")
	TokenEnterpriseIdErr = errorx.NewErrorMsg(1100002, "解析企业ID失败")
)
