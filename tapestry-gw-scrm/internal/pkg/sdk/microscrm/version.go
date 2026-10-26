package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

const (
	URL_VERSION_GetCurrentVersionInfo = "%s/micro-scrm/v1/version/current-version-info"
)

func GetCurrentVersionInfo(ctx ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(ctx, GetUrl(URL_VERSION_GetCurrentVersionInfo), request)
	return response, err
}
