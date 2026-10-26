package scrm_errs

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"

type Errs struct {
	Code    int64
	Message string
}

var ErrorStatus = 401

var (
	DefaultError       = &Errs{Code: 200000, Message: "接口异常"}
	AccountErr         = &Errs{Code: 200001, Message: "账号或密码有误，请联系管理员"}
	MobileErr          = &Errs{Code: 200002, Message: "手机或验证码有误，请联系管理员"}
	BindMobileErr      = &Errs{Code: 200003, Message: "手机或验证码有误，请联系管理员"}
	BindParamsErr      = &Errs{Code: 200004, Message: "参数有误"}
	AccountNotExists   = errorx.NewErrorMsg(200005, "数据异常，请联系管理员")
	EnterpriseIdNull   = errorx.NewErrorMsg(200006, "账号密码错误或无管理后台登录权限")
	NotFoundFileUpload = errorx.NewErrorMsg(121001, "没有要上传的文件")
	FileUploadErr      = errorx.NewErrorMsg(121002, "文件上传失败")
	FileUploadExtErr   = errorx.NewErrorMsg(121003, "文件类型不允许上传")
	FileUploadSize10Mb = errorx.NewErrorMsg(121004, "文件大小需在10Mb以下")
	FileUploadSize2Mb  = errorx.NewErrorMsg(121005, "文件大小需在2Mb以下")
)
