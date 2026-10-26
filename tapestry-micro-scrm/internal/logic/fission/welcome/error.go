package welcome

// 缺少活动信息
type welcomeParamError struct{}

func (welcomeParamError) Error() string {
	return "查询活动欢迎语参数不足"
}

type welcomeObjError struct{}

func (welcomeObjError) Error() string {
	return "构建欢迎语对象异常"
}
