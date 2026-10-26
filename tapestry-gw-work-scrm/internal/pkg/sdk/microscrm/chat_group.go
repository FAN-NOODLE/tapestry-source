package microScrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
)

type GroupChatApi struct {
	BaseSDK
}

// NewCustomer 返回API实例
func NewGroupChat(host string) *GroupChatApi {
	urlHost = host
	return &GroupChatApi{BaseSDK{c: &clients.HTTPClient{}}}
}

// Index 首页数据统计
func (c GroupChatApi) Index(ctx ctx.Context, req GroupChatStaticsRequest) (interface{}, error) {
	path := "/micro-scrm/v1/chat-group/index_list"
	url := urlHost + path
	res, err := c.Post(ctx, url, req)
	if err != nil {
		return nil, err
	}

	result := &GroupChatBaseResponse{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

// MembersGuide 会员入群列表 非会员未入群列表 非会员列表
func (c GroupChatApi) MembersGuide(ctx ctx.Context, req GroupChatStaticsRequest) (interface{}, error) {
	path := "/micro-scrm/v1/chat-group/members-guide"
	url := urlHost + path
	res, err := c.Post(ctx, url, req)
	if err != nil {
		return nil, err
	}

	result := &GroupChatBaseResponse{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

// GroupMembers 群成员明细
func (c GroupChatApi) Members(ctx ctx.Context, req GroupChatStaticsRequest) (interface{}, error) {
	path := "/micro-scrm/v1/chat-group/members_info"
	url := urlHost + path
	res, err := c.Post(ctx, url, req)
	if err != nil {
		return nil, err
	}

	result := &GroupChatBaseResponse{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

// GetDataUpdateTime 数据更时间
func (c GroupChatApi) GetDataUpdateTime(ctx ctx.Context, req RequestBase) (interface{}, error) {
	path := "/micro-scrm/v1/chat-group/data-update-time"
	url := urlHost + path
	res, err := c.Post(ctx, url, req)
	if err != nil {
		return nil, err
	}

	result := &GroupChatBaseResponse{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}

// ActivityInfo 活动详情
func (c GroupChatApi) ActivityInfo(ctx ctx.Context, req ActivityInfo) (interface{}, error) {
	path := "/micro-scrm/v1/chat-group/activity-info"
	url := urlHost + path
	res, err := c.Post(ctx, url, req)
	if err != nil {
		return nil, err
	}

	result := &GroupChatBaseResponse{}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}

	return result.Data, err
}
