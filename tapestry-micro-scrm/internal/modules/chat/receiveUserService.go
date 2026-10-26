package chat

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/typex"
	"micro-scrm/internal/pkg/uid"
)

// 消息接收人服务
type ReceiveUserService struct {
	repo ReceiveUserRepository
}

// 消息接收人新增
func (s ReceiveUserService) Add(c ctx.Context, row ReceiveUser) typex.ID {
	if row.ID == 0 {
		row.ID = uid.ID(c)
	}
	num, err := s.repo.Add(c, row)
	if err != nil {
		log.Error(ctx.Wrap(c)).Msgf("新增消息接收人记录错误:" + err.Error())

		return 0
	}

	return typex.ID(num)
}

// SelectMsgIDByUserIds 根据员工ID查询消息ID
func (s ReceiveUserService) SelectMsgIDByUserIds(c ctx.Context, memberUserIDs []string) ([]SearchReceiveResponse, error) {
	list, err := s.repo.SelectReceiveUserByMemberUserIDs(c, memberUserIDs)
	if err != nil {
		log.Error(ctx.Wrap(c)).Msgf("根据员工ID查询消息ID记录错误:" + err.Error())
		return nil, err
	}
	return list, nil
}
