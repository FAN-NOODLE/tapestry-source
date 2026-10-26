package http

/**
*数据罗盘
 */
import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
)

func Endpoints() []*swagger.Endpoint {
	return []*swagger.Endpoint{
		endpoint.Post("/datacompass/community_data", "社群数据明细",
			endpoint.Handler(datacompassCommunityData),
			endpoint.Tags("数据罗盘-社群数据")),
		endpoint.Post("/datacompass/community_data/export", "社群数据明细导出",
			endpoint.Handler(datacompassCommunityDataExport),
			endpoint.Tags("数据罗盘-社群数据")),
		endpoint.Post("/datacompass/friend_data", "好友数据明细",
			endpoint.Handler(datacompassFriendData),
			endpoint.Tags("数据罗盘-好友数据")),
		endpoint.Post("/datacompass/friend_data/export", "好友数据明细导出",
			endpoint.Handler(datacompassFriendDataExport),
			endpoint.Tags("数据罗盘-好友数据")),
		endpoint.Post("/datacompass/community_data_summary", "社群，好友，活动社群数据汇总",
			endpoint.Handler(dataCompassDataSummary),
			endpoint.Tags("数据罗盘-社群，好友，活动社群数据汇总")),
	}
}
