package member_level_http

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
)

func Endpoints() []*swagger.Endpoint {
	return []*swagger.Endpoint{
		endpoint.Post("/member-level/list", "会员等级列表",
			endpoint.Handler(List),
			endpoint.Tags("等级直达列表")),
		endpoint.Post("/member-level/set", "会员等级设置",
			endpoint.Handler(Set),
			endpoint.Tags("会员等级设置")),
		endpoint.Post("/member-level/detail", "会员等级详情",
			endpoint.Handler(Detail),
			endpoint.Tags("会员等级详情")),
		endpoint.Post("/member-level/update-status", "会员等级状态修改",
			endpoint.Handler(UpdateStatus),
			endpoint.Tags("会员等级状态修改")),
		endpoint.Post("/member-level/content-sort", "等级内容排序",
			endpoint.Handler(ContentSort),
			endpoint.Tags("等级内容排序")),
		endpoint.Post("/member/level-up-list", "等级跃迁列表",
			endpoint.Handler(LevelUpList),
			endpoint.Tags("等级跃迁")),
		endpoint.Post("/member/level-up-add", "新增等级跃迁",
			endpoint.Handler(AddLevelUp),
			endpoint.Tags("等级跃迁")),
		endpoint.Post("/member/level-up-edit", "新增等级跃迁",
			endpoint.Handler(EditLevelUp),
			endpoint.Tags("等级跃迁")),
		endpoint.Post("/member/level-up-delete", "新增等级跃迁",
			endpoint.Handler(DeleteLevelUp),
			endpoint.Tags("等级跃迁")),
		endpoint.Post("/member/level-up-info", "等级跃迁详情",
			endpoint.Handler(LevelUpInfo),
			endpoint.Tags("等级跃迁")),
		endpoint.Post("/member-level/level-list", "下拉框使用会员等级列表",
			endpoint.Handler(LevelList),
			endpoint.Tags("下拉框使用会员等级列表")),
	}
}
