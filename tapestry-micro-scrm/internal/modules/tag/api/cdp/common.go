package cdp

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"micro-scrm/internal/config"
)

// ResponseBase cdp接口返回结构体头
type ResponseBase struct {
	App       string `json:"app"`
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	RequestID string `json:"request_id"`
}

// GetURL 获取cdp请求链接
func GetURL(uri string) string {
	return fmt.Sprintf(uri, config.CfgData.Service.Cdp.Domain)
}

// CheckResponse 检查cdp接口返回结果
func CheckResponse(response ResponseBase) error {
	if response.Code == 0 {
		return nil
	}
	return errorx.NewErrorMsg(response.Code, response.Msg)
}
