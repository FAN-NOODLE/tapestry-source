package errs

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"

var (
	SmsSignNotReviewed     = errorx.NewErrorMsg(150001, "签名未审核通过")
	SmsTemplateNameRepeat  = errorx.NewErrorMsg(150002, "模板名称已被使用")
	SmsTemplateNonExistent = errorx.NewErrorMsg(150003, "短信模板不存在")
	TemplateSignFormat     = errorx.NewErrorMsg(150004, "短信签名格式不正确")
	TemplateContentFormat  = errorx.NewErrorMsg(150004, "短信模板不能为空")
	TemplateStatusNomarl   = errorx.NewErrorMsg(150005, "短信模板未赋值，不能开启")
	NotSmsNumber           = errorx.NewErrorMsg(150005, "发送短信失败-短信条数不足")
	CannotEditInReview     = errorx.NewErrorMsg(150006, "短信模板审核中不可编辑")
	CannotEditNot          = errorx.NewErrorMsg(150006, "短信模板不可编辑")
)
