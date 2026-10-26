package errs

import "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"

// ErrorStatus 错误状态码
var ErrorStatus = 401

// DefaultErrorCode 默认码
var DefaultErrorCode = 100000

// DefaultErrorMsg 默认错误信息
var DefaultErrorMsg = "参数错误"
var (
	// ParamsErr 参数错误
	ParamsErr = errorx.NewErrorMsg(DefaultErrorCode, DefaultErrorMsg)
	// EnterpriseIdEmpty 企业ID不存在
	EnterpriseIdEmpty = errorx.NewErrorMsg(100001, "企业ID不存在")
	// EnterpriseNotFund 企业信息不存在
	EnterpriseNotFund = errorx.NewErrorMsg(100002, "企业信息不存在")
	// STokenFail sToken获取失败
	STokenFail = errorx.NewErrorMsg(100003, "sToken获取失败")
	// RedisSetFail 设置缓存失败
	RedisSetFail = errorx.NewErrorMsg(100004, "设置缓存失败")
	// DataTransformErr 数据转换失败
	DataTransformErr = errorx.NewErrorMsg(100005, "数据转换失败")
	// WorkPluginErr 应用列表查询失败
	WorkPluginErr = errorx.NewErrorMsg(100006, "应用列表查询失败")
	// 获取自建应用信息
	PluginInfoRrr = errorx.NewErrorMsg(100007, "自建应用信息未配置")
	// ResponseNotJson
	ResponseNotJson = errorx.NewErrorMsg(100008, "返回格式非json")
	// QwResponseParseErr
	QwResponseParseErr = errorx.NewErrorMsg(100009, "企业微信获取结果失败")
	// MethodHeaderErr
	MethodHeaderErr = errorx.NewErrorMsg(100010, "Header头method设置错误")
	// ActionErr
	ActionErr = errorx.NewErrorMsg(100011, "操作失败")
	//
	DBErr = errorx.NewErrorMsg(100012, "数据库错误")

	FailedTaAddTask = errorx.NewErrorMsg(100013, "新增任务计划失败")

	DataRepeated    = errorx.NewErrorMsg(100014, "数据重复")
	RejectedService = errorx.NewErrorMsg(100015, "服务暂不可用")
	FileIsEmpty     = errorx.NewErrorMsg(100016, "上传文件不能为空")

	DataExist = errorx.NewErrorMsg(100017, "数据已存在")

	// OrderFeeErr
	OrderFeeErr    = errorx.NewErrorMsg(200010, "订单实付金额与商品实付总额不一致")
	OrderAmountErr = errorx.NewErrorMsg(200011, "订单总额与商品总额不一致")

	DataAlreadyExists = errorx.NewErrorMsg(300000, "数据已存在")
)

/**
*业务操作失败返回信息
 */
func OperationsBackError(msg string) *errorx.Error {
	return errorx.NewErrorMsg(100011, msg)
}
