package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

// 群成员关系计算 - 首页用
func GroupRelationCalc(c ctx.Context, request interface{}) (interface{}, error) {
	response, err := sdk.SendPost(c, GetUrl(urlEnum.GroupMemberRelation), request)
	return response, err
}
