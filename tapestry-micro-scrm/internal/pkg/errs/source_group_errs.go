package errs

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"

var (
	SourceGroupNameRepeated = errorx.NewErrorMsg(170001, "素材分组名称已存在")
	SourceGroupIDEmpty      = errorx.NewErrorMsg(170002, "素材分组ID不能为空")
	SourceGroupDelErr       = errorx.NewErrorMsg(170003, "删除素材分组失败")
)
