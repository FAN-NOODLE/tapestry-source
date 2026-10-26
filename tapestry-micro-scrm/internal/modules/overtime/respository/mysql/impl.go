package mysql

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/overtime"
	"micro-scrm/internal/pkg/errorIgnore"
)

type OvertimeRepo struct {
	*repos.SQLRepo
}

type ReceiveUserRepo struct {
	*repos.SQLRepo
}

// SelectList 获取超时记录列表
func (s *OvertimeRepo) SelectRecordsList(c ctx.Context, request overtime.ListRequest) (int64, *[]overtime.Records, error) {
	offset := (request.Page - 1) * request.Size
	pageVo := repos.PageVO{Offset: offset, Limit: request.Size}

	option := &repos.FindPageOption{
		TplPrefix: "t_chat_overtime_records/select-list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]overtime.Records{}
		},
	}
	total, list, err := s.FindPage(c, pageVo, request, option)
	if errorIgnore.Ignore(err) != nil {
		return 0, nil, err
	}

	resList := list.(*[]overtime.Records)
	return total, resList, err
}

// SelectChatMsgList 获取超时对应聊天记录
func (s *OvertimeRepo) SelectChatMsgList(ctx ctx.Context, ids []int64) (*[]overtime.ChatMsg, error) {
	mapping := make(map[string]interface{})
	mapping["ids"] = ids
	stmt, err := s.DB.Template("t_chat_msg/select-chat-msg-list", &mapping)
	if err != nil {
		return nil, err
	}
	list := &[]overtime.ChatMsg{}
	stmtErr := stmt.NamedSelect(ctx, list, mapping)
	if stmtErr != nil {
		return nil, stmtErr
	}
	return list, nil
}

// SelectUserRelationList 获取聊天记录对应好友关系
func (s *OvertimeRepo) SelectUserRelation(ctx ctx.Context, request overtime.ChatMsgMembers) (*overtime.ChatMsgMembers, error) {
	mapping := make(map[string]interface{})
	mapping["cond"] = request
	stmt, err := s.DB.Template("t_work_external_user/work-external-user-by-relation", &mapping)
	if err != nil {
		return nil, err
	}
	workUser := &overtime.ChatMsgMembers{}
	errStmt := stmt.NamedGet(ctx, workUser, request)
	if errStmt != nil {
		return nil, errStmt
	}
	return workUser, nil
}

// MemberMultiDepartment 根据memberID查询所在部门信息
func (s *OvertimeRepo) MemberMultiDepartments(cc ctx.Context, memberIDs []int64) (map[int64]overtime.MemberDepartsList, error) {
	memberIDsCond := map[string]interface{}{
		"memberIDs": memberIDs,
	}

	memberDepartsList := &overtime.MemberDepartsList{}
	memberOption := &repos.ListOption{
		TplPrefix: "t_departments/select-dept-by-memberids",
		Fn: func(context ctx.Context, i int64) interface{} {
			return memberDepartsList
		},
	}

	departmentListById := make(map[int64]overtime.MemberDepartsList)
	_, err := s.List(cc, memberIDsCond, memberOption)
	if err != nil {
		return departmentListById, err
	}

	for _, department := range *memberDepartsList {
		departmentListById[department.MemberID] = append(departmentListById[department.MemberID], department)
	}

	return departmentListById, err
}

// 查询外部联系人会话同意记录
func (s *OvertimeRepo) SelectReplyMsg(ctx ctx.Context, request overtime.ChatMsg) (*[]overtime.ChatMsg, error) {
	mapping := make(map[string]interface{})
	mapping["cond"] = request
	stmt, err := s.DB.Template("t_chat_msg/select-reply-msg", &mapping)
	if err != nil {
		return nil, err
	}
	list := &[]overtime.ChatMsg{}
	stmtErr := stmt.NamedSelect(ctx, list, mapping)
	if stmtErr != nil {
		return nil, stmtErr
	}
	return list, nil
}

// 记录是否存在
func (s *OvertimeRepo) FindOvertime(ctx ctx.Context, request overtime.Records) (*overtime.Records, error) {
	mapping := make(map[string]interface{})
	mapping["cond"] = request
	stmt, err := s.DB.Template("t_chat_overtime_records/exists", &mapping)
	if err != nil {
		return nil, err
	}
	detail := &overtime.Records{}
	errStmt := stmt.NamedGet(ctx, detail, request)
	if errStmt != nil {
		return nil, errStmt
	}
	return detail, nil
}

// SelectMemberByWhere 根据员工查询条件，返回员工userid
func (s *OvertimeRepo) SelectMemberByWhere(ctx ctx.Context, request overtime.ListRequest) (*[]overtime.SearchMemberResponse, error) {
	mapping := make(map[string]interface{})
	mapping["cond"] = request
	stmt, err := s.DB.Template("t_members/select-member-by-where", &mapping)
	if err != nil {
		return nil, err
	}
	list := &[]overtime.SearchMemberResponse{}
	stmtErr := stmt.NamedSelect(ctx, list, mapping)
	if stmtErr != nil {
		return nil, stmtErr
	}
	return list, nil
}

// GetMembersList 获取指定userids的用户信息
func (s *OvertimeRepo) SelectMembersList(c ctx.Context, req []string, enterpriseId int64) (*[]overtime.MemberInfo, error) {
	mapping := make(map[string]interface{})
	mapping["user_ids"] = req
	mapping["enterprise_id"] = enterpriseId
	stmt, err := s.DB.Template("t_members/select-member-by-userids", &mapping)
	if err != nil {
		return nil, err
	}
	list := &[]overtime.MemberInfo{}
	stmtErr := stmt.NamedSelect(c, list, mapping)
	if stmtErr != nil {
		return nil, stmtErr
	}
	return list, nil
}
