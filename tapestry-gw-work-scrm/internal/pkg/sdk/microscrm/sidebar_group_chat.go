package microScrm

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
)

// MemberApi 员工 SDK
type SidebarGroupChatApi struct {
	BaseSDK
}

// NewMember 返回API实例
func NewSidebarGroupChat(host string) *SidebarGroupChatApi {
	urlHost = host
	return &SidebarGroupChatApi{BaseSDK{c: &clients.HTTPClient{}}}
}

// GroupChatInfo 获取社群侧边栏基础信息
func (s SidebarGroupChatApi) GroupChatInfo(c ctx.Context, request GroupChatInfoRequest) (*GroupChatInfoResponse, error) {
	path := "/micro-scrm/v1/sidebar/group-chat/info"
	url := urlHost + path
	res, err := s.Post(c, url, request)
	if err != nil {
		return nil, err
	}

	result := &Response{
		Data: &GroupChatInfoResponse{},
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	if result.Data == nil {
		return nil, err
	}

	return result.Data.(*GroupChatInfoResponse), err
}

// GroupChatMember 侧边栏社群成员列表
func (s SidebarGroupChatApi) GroupChatMember(c ctx.Context, request GroupChatInfoRequest) (*GroupChatUserListResponse, error) {
	path := "/micro-scrm/v1/sidebar/group-chat/group-chat-member"
	url := urlHost + path
	res, err := s.Post(c, url, request)
	if err != nil {
		return nil, err
	}

	result := &Response{
		Data: &GroupChatUserListResponse{},
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	if result.Data == nil {
		return nil, err
	}

	return result.Data.(*GroupChatUserListResponse), err
}

// GroupUsersAddEmployeesRequest GroupChatMember 侧边栏社群成员列表
func (s SidebarGroupChatApi) GroupUsersAddEmployeesRequest(c ctx.Context, request GroupUsersAddEmployeesRequest) (*GroupUsersAddEmployeesResponse, error) {
	path := "/micro-scrm/v1/sidebar/group-chat/group-users-add-employees"
	url := urlHost + path
	res, err := s.Post(c, url, request)
	if err != nil {
		return nil, err
	}

	result := &Response{
		Data: &GroupUsersAddEmployeesResponse{},
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	if result.Data == nil {
		return nil, err
	}

	return result.Data.(*GroupUsersAddEmployeesResponse), err
}

// GroupChatMember 侧边栏社群成员列表
func (s SidebarGroupChatApi) GroupUsersJoinActivityRequest(c ctx.Context, request GroupUsersJoinActivityRequest) (*GroupUsersJoinActivityResponse, error) {
	path := "/micro-scrm/v1/sidebar/group-chat/group-users-join-activity"
	url := urlHost + path
	res, err := s.Post(c, url, request)
	if err != nil {
		return nil, err
	}

	result := &Response{
		Data: &GroupUsersJoinActivityResponse{},
	}
	_, err = res.Read(clients.JSON(result))
	if result.Code != StatusOK {
		return nil, errorx.NewErrorMsg(result.Code, result.Message)
	}
	if result.Data == nil {
		return nil, err
	}

	return result.Data.(*GroupUsersJoinActivityResponse), err
}
