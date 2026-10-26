package errs

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"

var (
	// 外部联系人错误
	FindWorkUserParamErr = errorx.NewErrorMsg(1300001, "外部联系人查询external_user_id、union_id、id 不能同时为空")
	// NotFundRelationFromQwErr 未从企业微信查询到好友关系
	NotFundRelationFromQwErr = errorx.NewErrorMsg(1300002, "未从企业微信查询到好友关系")
	// CallbackQueryErr 企业微信回调信息未正确获取到query参数
	CallbackQueryErr = errorx.NewErrorMsg(1300003, "企业微信回调信息未正确获取到query参数")
	// CallbackXmlErr 企业微信回调信息未正确获取xml参数
	CallbackXmlErr = errorx.NewErrorMsg(1300004, "企业微信回调信息未正确获取xml参数")
	// ExternalEditRemarksErrCode 编辑外部联系人备注错误码
	ExternalEditRemarksErrCode  = 1300005
	ExternalRemarksLengthErr    = errorx.NewErrorMsg(1300006, "企微外部联系人描述长度超过限制")
	ExternalRelationNotExistErr = errorx.NewErrorMsg(1300007, "企微外部联系人数据不存在")
	// GetCorpTagListErr 获取企业标签列表错误码
	GetCorpTagListErr = 1300008
	// 删除企业微信标签错误code
	DelCorpTagErrCode = 1300009
	ParseResultErr    = errorx.NewErrorMsg(1300010, "解析结果错误")
)
