package chat

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/pkg/errorIgnore"
	"micro-scrm/internal/pkg/uid"
	"time"
)

// 员工服务
type MemberService struct {
	repo MemberRepository
}

// 员工分页列表
func (s *MemberService) MemberPageList(c ctx.Context, cond MemberPageListRequest) (MemberPageListResponse, error) {
	var cmp MemberPageListResponse

	offset := (cond.Page - 1) * cond.Size
	pageVo := repos.PageVO{Offset: int(offset), Limit: int(cond.Size)}
	option := &repos.FindPageOption{
		TplPrefix: "t_chat_member/list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]MemberInfo{}
		},
	}
	total, list, err := s.repo.FindPage(c, pageVo, cond, option)
	if errorIgnore.Ignore(err) != nil {
		return cmp, err
	}
	resList := list.(*[]MemberInfo)

	cmp.Page = cond.Page
	cmp.Size = cond.Size
	cmp.List = *resList
	cmp.Total = total

	return cmp, nil
}

// Add 新增员工聊天关系
func (s MemberService) Add(c ctx.Context, row MemberRelation) int64 {
	// 判断是否存在,存在则更新最后聊天时间
	memberRelation := s.FindMember(c, row.EnterpriseID, row.MemberID)
	var num int64
	if memberRelation.ID != 0 {
		return s.Update(c, memberRelation.ID)
	}

	row.ID = uid.ID(c)
	num, err := s.repo.Add(c, row)
	if err != nil {
		log.Info(c).Msgf("新增员工聊天关系错误:" + err.Error())
	}

	return num
}

func (s MemberService) FindMember(c ctx.Context, enterpriseID int64, memberID int64) MemberRelation {
	var memberRelation MemberRelation
	err := s.repo.FindMember(c, struct {
		MemberID     int64 `json:"member_id" db:"member_id"`
		EnterpriseID int64 `json:"enterprise_id" db:"enterprise_id"`
	}{
		MemberID:     memberID,
		EnterpriseID: enterpriseID,
	}, &memberRelation)
	if err != nil {
		log.Info(c).Msgf("查找员工聊天关系错误:" + err.Error())
	}

	return memberRelation
}

func (s MemberService) Update(c ctx.Context, id int64) int64 {
	t := time.Now().Unix()
	num, err := s.repo.Update(c, struct {
		ID         int64 `json:"id" db:"id"`
		LastChatAt int64 `json:"last_chat_at" db:"last_chat_at"`
		UpdateAt   int64 `json:"update_at" db:"update_at"`
	}{
		ID:         id,
		LastChatAt: t,
		UpdateAt:   t,
	})
	if err != nil {
		log.Info(c).Msgf("更新员工聊天关系错误:" + err.Error())
	}

	return num
}
