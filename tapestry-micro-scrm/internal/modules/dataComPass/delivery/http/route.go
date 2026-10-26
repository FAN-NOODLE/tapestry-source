package http

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
)

func Endpoints() []*swagger.Endpoint {
	friendHandler := NewGoodFriendsHandler()
	return []*swagger.Endpoint{
		endpoint.Post("/datacompass/community_data", "社群数据明细",
			endpoint.Handler(friendHandler.datacompassCommunityData),
			endpoint.Tags("数据罗盘-社群数据")),
		endpoint.Post("/datacompass/friend_data", "好友数据明细",
			endpoint.Handler(friendHandler.datacompassFriendData),
			endpoint.Tags("数据罗盘-好友数据")),
		endpoint.Post("/data-com-pass/get_compass_summary", "社群，好友，活动社群数据汇总",
			endpoint.Handler(friendHandler.getGroupNumber)),
		endpoint.Post("/data-com-pass/get_ba_department", "获取导购的一级部门",
			endpoint.Handler(friendHandler.getBaFirstDepartment)),

		endpoint.Post("/datacompass/rank_data", "好友数据明细",
			endpoint.Handler(friendHandler.rankData),
			endpoint.Tags("顾客总计-顾客数据")),

		endpoint.Post("/data-com-pass/get_compass_summary_ck", "社群，好友，活动社群数据汇总(Ck)",
			endpoint.Handler(friendHandler.getGroupNumberCk)),
	}
}
