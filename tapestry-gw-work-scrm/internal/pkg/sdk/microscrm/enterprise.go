package microScrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
)

// EnterpriseApi 企业SDK
type EnterpriseApi struct {
	BaseSDK
}

// NewEnterprise 返回API实例
func NewEnterprise(host string) *EnterpriseApi {
	urlHost = host
	return &EnterpriseApi{BaseSDK{c: &clients.HTTPClient{}}}
}

// JsAgentConfig 应用签名
func (e EnterpriseApi) JsAgentConfig(c ctx.Context, request JsConfigRequest) (*AgentConfigResponse, error) {
	path := "/micro-scrm/v1/js/agent-config"
	url := urlHost + path
	res, err := e.Post(c, url, request)
	if err != nil {
		return nil, err
	}

	result := &Response{
		Data: &AgentConfigResponse{},
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	if result.Data == nil {
		return nil, err
	}

	return result.Data.(*AgentConfigResponse), err
}

// JsConfig 企业签名
func (e EnterpriseApi) JsConfig(c ctx.Context, request JsConfigRequest) (*JsConfigResponse, error) {
	path := "/micro-scrm/v1/js/config"
	url := urlHost + path
	res, err := e.Post(c, url, request)
	if err != nil {
		return nil, err
	}

	result := &Response{
		Data: &JsConfigResponse{},
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	if result.Data == nil {
		return nil, err
	}

	return result.Data.(*JsConfigResponse), err
}
