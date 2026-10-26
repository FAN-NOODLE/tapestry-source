package microshop

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

func AuthorizationDetail(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.AuthorizationDetail), request)
	return response, err
}

func AuthorizationSave(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.AuthorizationSave), request)
	return response, err
}
