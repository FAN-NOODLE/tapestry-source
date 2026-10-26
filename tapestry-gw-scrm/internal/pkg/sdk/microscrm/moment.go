package microscrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

// MemberTopic 一个员工的朋友圈列表
func MemberTopic(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.MemberTopicUrl), request)
}

// TopicCreate 发布朋友圈
func TopicCreate(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.TopicCreateUrl), request)
}

// TopicEdit 编辑朋友圈
func TopicEdit(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.TopicEditUrl), request)
}

// TopicDelete 删除朋友圈
func TopicDelete(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.TopicDeleteUrl), request)
}

// TopicDisplay 朋友圈展示
func TopicDisplay(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.TopicDisplayUrl), request)
}

// ContentList 朋友圈按内容列表展示
func ContentList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.ContentListUrl), request)
}

// MemberList 朋友圈列表按员工搜索
func MemberList(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.MomentMemberListUrl), request)
}

// AppearanceGet 获取朋友圈形象
func AppearanceGet(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.AppearanceGetUrl), request)
}

// AppearanceEdit 编辑朋友圈形象
func AppearanceEdit(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.AppearanceEditUrl), request)
}

// JsConfig 获取jsConfig
func JsConfig(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.JsConfigUrl), request)
}

// TopicMemberDelete 删除员工朋友圈
func TopicMemberDelete(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.TopicMemberDeleteUrl), request)
}

// TopicMemberEdit 编辑员工朋友圈
func TopicMemberEdit(c ctx.Context, request interface{}) (interface{}, error) {
	return sdk.SendPost(c, GetUrl(urlEnum.TopicMemberEditUrl), request)
}
