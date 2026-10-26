package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

const (
	URL_ACTIVITY_SIGN_IN_LIST   = "%s/micro-scrm/v1/activities/sign/list"
	URL_ACTIVITY_SIGN_IN_DETAIL = "%s/micro-scrm/v1/activities/sign/detail"
)

func GetActivitySignInList(ctx ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(ctx, GetUrl(URL_ACTIVITY_SIGN_IN_LIST), request)
	return response, err
}

func GetActivitySignInDetail(ctx ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(ctx, GetUrl(URL_ACTIVITY_SIGN_IN_DETAIL), request)
	return response, err
}
