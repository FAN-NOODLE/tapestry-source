package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

const (
	datacompassCommunityDataURL        = "%s/micro-scrm/v1/data-compass/chat-group/list"
	datacompassFriendDataExportURL     = "%s/micro-scrm/v1/datacompass/friend_data"
	datacompassCommunityDataSummaryURL = "%s/micro-scrm/v1/data-com-pass/get_compass_summary_ck"
)

func DatacompassCommunityData(ctx ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(ctx, GetUrl(datacompassCommunityDataURL), request)
	return response, err
}

func DatacompassFriendData(ctx ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(ctx, GetUrl(datacompassFriendDataExportURL), request)
	return response, err
}

func GetDatacompassCommunityDataSummary(ctx ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(ctx, GetUrl(datacompassCommunityDataSummaryURL), request)
	return response, err
}
