package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

const (
	PLAN_TASK_DATA_STAT_URL            = "%s/micro-scrm/v1/plan/data/stat"
	PLAN_TASK_DATA_STAT_DETAIL_URL     = "%s/micro-scrm/v1/plan/data/stat/detail"
	PLAN_TASK_CONSUMER_LIST_URL        = "%s/micro-scrm/v1/plan/consumer/list"
	PLAN_TASK_NEW_DATA_STAT_URL        = "%s/micro-scrm/v1/plan/task/data/new-stat"
	PLAN_TASK_NEW_DATA_STAT_DETAIL_URL = "%s/micro-scrm/v1/plan/task/data/stat/new-detail"
)

func TaskDataStat(ctx ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(ctx, GetUrl(PLAN_TASK_DATA_STAT_URL), request)
	return response, err
}

func TaskDataStatDetail(ctx ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(ctx, GetUrl(PLAN_TASK_DATA_STAT_DETAIL_URL), request)
	return response, err
}

func TaskConsumeList(ctx ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(ctx, GetUrl(PLAN_TASK_CONSUMER_LIST_URL), request)
	return response, err
}

func TaskNewDataStat(ctx ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(ctx, GetUrl(PLAN_TASK_NEW_DATA_STAT_URL), request)
	return response, err
}

func TaskNewDataStatDetail(ctx ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(ctx, GetUrl(PLAN_TASK_NEW_DATA_STAT_DETAIL_URL), request)
	return response, err
}
