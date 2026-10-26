package chat

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/pkg/errorIgnore"
)

// 员工与群聊
type MemberToGroupService struct {
	repo MemberToGroupRepository
}

// 接收方列表
func (s *MemberToGroupService) List(c ctx.Context, cond MemberToGroupPageListRequest) (MemberToGroupPageListResponse, error) {
	var cmp MemberToGroupPageListResponse

	offset := (cond.Page - 1) * cond.Size
	pageVo := repos.PageVO{Offset: int(offset), Limit: int(cond.Size)}
	option := &repos.FindPageOption{
		TplPrefix: "t_chat_member_to_group/list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]MemberToGroupInfo{}
		},
	}
	total, list, err := s.repo.FindPage(c, pageVo, cond, option)
	if errorIgnore.Ignore(err) != nil {
		return cmp, err
	}
	resList := list.(*[]MemberToGroupInfo)

	cmp.Page = cond.Page
	cmp.Size = cond.Size
	cmp.List = *resList
	cmp.Total = total

	return cmp, nil
}
