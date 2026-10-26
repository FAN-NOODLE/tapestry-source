package errs

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"

var (
	SourceNotFund   = errorx.NewErrorMsg(170004, "素材信息不存在")
	SourceMediaFail = errorx.NewErrorMsg(170005, "生成临时素材失败")
	SourceTypeErr   = errorx.NewErrorMsg(170006, "临时素材类型错误")
)
