package mysql

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	groupmember "micro-scrm/internal/modules/group-member"
)

type GroupRepo struct {
	*repos.SQLRepo
}

// SelectGroupNum 查询群总数
func (group *GroupRepo) SelectGroupNum(c ctx.Context, param groupmember.GroupParam) (int64, error) {
	stmt, err := group.DB.Template("t_work_group_chat/select-group-num", group.Cond2Args(param))
	if err != nil {
		return 0, err
	}

	res := &groupmember.AsRes{}
	err = stmt.NamedGet(c, res, param)
	if err != nil {
		return 0, err
	}
	return res.RelationNum, nil
}

// SelectReGroupUv 查询留存
func (group *GroupRepo) SelectReGroupUv(c ctx.Context, param groupmember.GroupParam) (int64, error) {
	stmt, err := group.DB.Template("t_work_group_chat_member/select-re-group-uv", group.Cond2Args(param))
	if err != nil {
		return 0, err
	}

	res := &groupmember.AsRes{}
	err = stmt.NamedGet(c, res, param)
	if err != nil {
		return 0, err
	}
	return res.RelationNum, nil
}

// SelectAddGroupUv 查询新增群成员数
func (group *GroupRepo) SelectAddGroupUv(c ctx.Context, param groupmember.GroupParam) (int64, error) {
	stmt, err := group.DB.Template("t_work_group_chat_member/select-add-group-uv", group.Cond2Args(param))
	if err != nil {
		return 0, err
	}

	res := &groupmember.AsRes{}
	err = stmt.NamedGet(c, res, param)
	if err != nil {
		return 0, err
	}
	return res.RelationNum, nil
}

// SelectLossGroupUv 查询流失群成员数
func (group *GroupRepo) SelectLossGroupUv(c ctx.Context, param groupmember.GroupParam) (int64, error) {
	stmt, err := group.DB.Template("t_work_group_chat_member/select-loss-group-uv", group.Cond2Args(param))
	if err != nil {
		return 0, err
	}

	res := &groupmember.AsRes{}
	err = stmt.NamedGet(c, res, param)
	if err != nil {
		return 0, err
	}
	return res.RelationNum, nil
}
