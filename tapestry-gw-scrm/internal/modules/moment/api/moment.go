package api

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/modules/moment"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type MomentApi struct {
}

// MemberTopic 一个员工的朋友圈列表
func (m MomentApi) MemberTopic(c ctx.Context, req *moment.MemberTopicRequest) (interface{}, error) {
	res, err := microscrm.MemberTopic(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// TopicEdit 朋友圈编辑
func (m MomentApi) TopicEdit(c ctx.Context, req *moment.TopicEditRequest) (interface{}, error) {
	res, err := microscrm.TopicEdit(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// TopicDelete 朋友圈删除
func (m MomentApi) TopicDelete(c ctx.Context, req *moment.TopicDeleteRequest) (interface{}, error) {
	res, err := microscrm.TopicDelete(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// TopicCreate 分布朋友圈
func (m MomentApi) TopicCreate(c ctx.Context, req *moment.TopicCreateRequest) (interface{}, error) {
	res, err := microscrm.TopicCreate(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// TopicDisplay 分布朋展示
func (m MomentApi) TopicDisplay(c ctx.Context, req *moment.TopicDisplayRequest) (interface{}, error) {
	res, err := microscrm.TopicDisplay(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// ContentList 朋友圈按内容列表展示
func (m MomentApi) ContentList(c ctx.Context, req *moment.ContentListRequest) (interface{}, error) {
	res, err := microscrm.ContentList(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// MemberList 朋友圈列表按员工
func (m MomentApi) MemberList(c ctx.Context, req *moment.MemberListRequest) (interface{}, error) {
	res, err := microscrm.MemberList(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// MemberList 获取朋友圈形象
func (m MomentApi) AppearanceGet(c ctx.Context, req *moment.MomentMemberGetRequest) (interface{}, error) {
	res, err := microscrm.AppearanceGet(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// AppearanceEdit 编辑朋友圈形象
func (m MomentApi) AppearanceEdit(c ctx.Context, req *moment.MomentMemberEditRequest) (interface{}, error) {
	res, err := microscrm.AppearanceEdit(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// JsConfig JS配置
func (m MomentApi) JsConfig(c ctx.Context, req *moment.JsConfigRequest) (interface{}, error) {
	res, err := microscrm.JsConfig(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// TopicMemberDelete 删除员工朋友圈
func (m MomentApi) TopicMemberDelete(c ctx.Context, req *moment.TopicMemberDeleteRequest) (interface{}, error) {
	res, err := microscrm.TopicMemberDelete(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// TopicMemberEdit 编辑员工朋友圈
func (m MomentApi) TopicMemberEdit(c ctx.Context, req *moment.TopicMemberEditRequest) (interface{}, error) {
	res, err := microscrm.TopicMemberEdit(c, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
