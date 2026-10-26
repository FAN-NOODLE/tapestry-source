package moment_http

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/docs/swagger/endpoint"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gw-scrm/internal/modules/moment"
	"net/http"
)

func Endpoints() []*swagger.Endpoint {
	return []*swagger.Endpoint{

		endpoint.Post("/moment/topic-create", "朋友圈手动发布",
			endpoint.Handler(topicCreate),
			endpoint.Body(moment.TopicCreateRequest{}, "", true),
			endpoint.Tags("朋友圈"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: moment.TopicCreateResponse{},
			}, "topicCreate", "接口调用成功")),

		endpoint.Post("/moment/topic-edit", "朋友圈编辑",
			endpoint.Handler(topicEdit),
			endpoint.Body(moment.TopicEditRequest{}, "", true),
			endpoint.Tags("朋友圈"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: moment.TopicEditResponse{},
			}, "topicEdit", "接口调用成功")),

		endpoint.Post("/moment/topic-delete", "朋友圈删除",
			endpoint.Handler(topicDelete),
			endpoint.Body(moment.TopicDeleteRequest{}, "", true),
			endpoint.Tags("朋友圈"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: moment.TopicDeleteResponse{},
			}, "topicDelete", "接口调用成功")),

		endpoint.Post("/moment/topic-display", "朋友圈展示",
			endpoint.Handler(topicDisplay),
			endpoint.Body(moment.TopicDisplayRequest{}, "", true),
			endpoint.Tags("朋友圈"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: moment.TopicDisplayResponse{},
			}, "topicCreate", "接口调用成功")),

		endpoint.Post("/moment/member-topic", "朋友圈一个员工的朋友圈列表",
			endpoint.Handler(memberTopic),
			endpoint.Body(moment.MemberTopicRequest{}, "", true),
			endpoint.Tags("朋友圈"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: moment.MemberTopicResponse{},
			}, "memberTopic", "接口调用成功")),

		endpoint.Post("/moment/content-list", "朋友圈按内容列表展示",
			endpoint.Handler(contentList),
			endpoint.Body(moment.ContentListRequest{}, "", true),
			endpoint.Tags("朋友圈"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: moment.ContentListResponse{},
			}, "contentList", "接口调用成功")),

		endpoint.Post("/moment/member-list", "朋友圈按员工列表展示",
			endpoint.Handler(memberList),
			endpoint.Body(moment.MemberListRequest{}, "", true),
			endpoint.Tags("朋友圈"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: moment.MemberListResponse{},
			}, "memberList", "接口调用成功")),

		endpoint.Post("/moment/appearance-get", "获取朋友圈形象",
			endpoint.Handler(appearanceGet),
			endpoint.Body(moment.MomentMemberGetRequest{}, "", true),
			endpoint.Tags("朋友圈"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: moment.MomentMemberGetResponse{},
			}, "appearanceGet", "接口调用成功")),

		endpoint.Post("/moment/appearance-get-admin", "管理后台调用获取朋友圈形象",
			endpoint.Handler(appearanceGet),
			endpoint.Body(moment.MomentMemberGetRequest{}, "", true),
			endpoint.Tags("朋友圈"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: moment.MomentMemberGetResponse{},
			}, "appearanceGet", "接口调用成功")),

		endpoint.Post("/moment/appearance-edit", "编辑朋友圈形象",
			endpoint.Handler(appearanceEdit),
			endpoint.Body(moment.MomentMemberEditRequest{}, "", true),
			endpoint.Tags("朋友圈"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: moment.MomentMemberEditResponse{},
			}, "appearanceEdit", "接口调用成功")),

		endpoint.Post("/moment/js-config", "前端jsconfig获取",
			endpoint.Handler(jsConfig),
			endpoint.Tags("朋友圈"),
			endpoint.Response(http.StatusOK, "", "jsconfig", "接口调用成功")),

		endpoint.Post("/moment/topic-member-delete", "删除员工朋友圈",
			endpoint.Handler(topicMemberDelete),
			endpoint.Body(moment.TopicMemberDeleteRequest{}, "", true),
			endpoint.Tags("朋友圈"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: moment.TopicMemberDeleteResponse{},
			}, "topicMemberDelete", "接口调用成功")),

		endpoint.Post("/moment/topic-member-edit", "编辑员工朋友圈",
			endpoint.Handler(topicMemberEdit),
			endpoint.Body(moment.TopicMemberEditRequest{}, "", true),
			endpoint.Tags("朋友圈"),
			endpoint.Response(http.StatusOK, ginplus.APIResponse{
				Data: moment.TopicMemberEditResponse{},
			}, "topicMemberEdit", "接口调用成功")),
	}
}
