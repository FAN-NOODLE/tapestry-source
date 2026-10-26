package scrm_errs

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"

var (
	MemberIdIsEmpty     = errorx.NewErrorMsg(120001, "成员ID不能为空")
	EnterpriseIdIsEmpty = errorx.NewErrorMsg(120002, "企业id不能为空")
	ParamsErr           = errorx.NewErrorMsg(120003, "请求参数有误")
	UploadErr           = errorx.NewErrorMsg(120004, "文件上传出错")
	InputEmpty          = errorx.NewErrorMsg(120005, "保存数据不能为空")
)
