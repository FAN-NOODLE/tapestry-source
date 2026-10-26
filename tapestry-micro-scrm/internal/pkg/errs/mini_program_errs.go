package errs

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"

var (
	// AuthorizeMiniProgramCannotOperate 授权小程序不允许操作
	MiniProgramEmpty                  = errorx.NewErrorMsg(170001, "小程序不存在")
	AuthorizeMiniProgramCannotOperate = errorx.NewErrorMsg(170002, "授权小程序不允许操作")
)
