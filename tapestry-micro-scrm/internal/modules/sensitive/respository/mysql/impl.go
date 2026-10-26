package mysql

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/sensitive"
	"micro-scrm/internal/pkg/errorIgnore"
)

type SensitiveRepo struct {
	*repos.SQLRepo
}

// SelectList 获取敏感词列表
func (s *SensitiveRepo) SelectList(c ctx.Context, request sensitive.ListRequest) (int64, *[]sensitive.Detail, error) {
	offset := (request.Page - 1) * request.Size
	pageVo := repos.PageVO{Offset: offset, Limit: request.Size}

	option := &repos.FindPageOption{
		TplPrefix: "t_chat_sensitive/select-list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]sensitive.Detail{}
		},
	}
	total, list, err := s.FindPage(c, pageVo, request, option)
	if errorIgnore.Ignore(err) != nil {
		return 0, nil, err
	}

	resList := list.(*[]sensitive.Detail)
	return total, resList, err
}

// CountBySensitive 获取敏感词每天触发次数
func (s *SensitiveRepo) CountBySensitive(ctx ctx.Context, ids []int64) (*[]sensitive.CountBySensitive, error) {
	mapping := make(map[string]interface{})
	mapping["ids"] = ids
	stmt, err := s.DB.Template("t_chat_sensitive_records/count-by-sensitive", &mapping)
	if err != nil {
		return nil, err
	}
	list := &[]sensitive.CountBySensitive{}
	stmtErr := stmt.NamedSelect(ctx, list, mapping)
	if stmtErr != nil {
		return nil, stmtErr
	}
	return list, nil
}

// FindSensitive 查询敏感词信息
func (s *SensitiveRepo) FindSensitive(ctx ctx.Context, request sensitive.SaveRequest) (*sensitive.Detail, error) {
	mapping := make(map[string]interface{})
	mapping["cond"] = request
	stmt, err := s.DB.Template("t_chat_sensitive/find-sensitive", &mapping)
	if err != nil {
		return nil, err
	}
	detail := &sensitive.Detail{}
	errStmt := stmt.NamedGet(ctx, detail, request)
	if errStmt != nil {
		return nil, errStmt
	}
	return detail, nil
}

// DeleteSensitive 删除敏感词
func (s *SensitiveRepo) DeleteSensitive(c ctx.Context, request sensitive.DeleteRequest) (int64, error) {
	stmt, stmtErr := s.DB.Template("t_chat_sensitive/delete-sensitive")
	if stmtErr != nil {
		return 0, stmtErr
	}
	result, execError := stmt.NamedExec(c, request)
	if execError != nil {
		return 0, execError
	}
	return result.RowsAffected()
}

// DeleteSensitiveRecords 删除敏感词聊天记录
func (s *SensitiveRepo) DeleteSensitiveRecords(c ctx.Context, request sensitive.DeleteRequest) (int64, error) {
	stmt, stmtErr := s.DB.Template("t_chat_sensitive_records/delete-sensitive-records")
	if stmtErr != nil {
		return 0, stmtErr
	}
	result, execError := stmt.NamedExec(c, request)
	if execError != nil {
		return 0, execError
	}
	return result.RowsAffected()
}

// SelectRecordsList 获取敏感词关联聊天记录列表
func (s *SensitiveRepo) SelectRecordsList(c ctx.Context, request sensitive.DetailListRequest) (int64, *[]sensitive.Records, error) {
	offset := (request.Page - 1) * request.Size
	pageVo := repos.PageVO{Offset: offset, Limit: request.Size}

	option := &repos.FindPageOption{
		TplPrefix: "t_chat_sensitive_records/select-list",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]sensitive.Records{}
		},
	}
	total, list, err := s.FindPage(c, pageVo, request, option)
	if errorIgnore.Ignore(err) != nil {
		return 0, nil, err
	}

	resList := list.(*[]sensitive.Records)
	return total, resList, err
}

// SelectChatMsgList 获取敏感词对应聊天记录
func (s *SensitiveRepo) SelectChatMsgList(ctx ctx.Context, ids []int64) (*[]sensitive.ChatMsg, error) {
	mapping := make(map[string]interface{})
	mapping["ids"] = ids
	stmt, err := s.DB.Template("t_chat_msg/select-chat-msg-list", &mapping)
	if err != nil {
		return nil, err
	}
	list := &[]sensitive.ChatMsg{}
	stmtErr := stmt.NamedSelect(ctx, list, mapping)
	if stmtErr != nil {
		return nil, stmtErr
	}
	return list, nil
}

// SelectUserRelationList 获取聊天记录对应好友关系
func (s *SensitiveRepo) SelectUserRelation(ctx ctx.Context, request sensitive.ChatMsgMembers) (*sensitive.DetailListResponse, error) {
	mapping := make(map[string]interface{})
	mapping["cond"] = request
	stmt, err := s.DB.Template("t_work_external_user/work-external-user-by-relation", &mapping)
	if err != nil {
		return nil, err
	}
	workUser := &sensitive.DetailListResponse{}
	errStmt := stmt.NamedGet(ctx, workUser, request)
	if errStmt != nil {
		return nil, errStmt
	}
	return workUser, nil
}

// InsertSensitiveRecords 插入敏感词聊天记录
func (s *SensitiveRepo) InsertSensitiveRecords(c ctx.Context, request sensitive.Records) (int64, error) {
	stmt, stmtErr := s.DB.Template("t_chat_sensitive_records/insert")
	if stmtErr != nil {
		return 0, stmtErr
	}
	result, execError := stmt.NamedExec(c, request)
	if execError != nil {
		return 0, execError
	}
	return result.RowsAffected()
}

// GetMembersList 获取指定userids的用户信息
func (s *SensitiveRepo) GetMembersList(c ctx.Context, req []string, enterpriseId int64) (*[]sensitive.UserInfo, error) {
	mapping := make(map[string]interface{})
	mapping["user_ids"] = req
	mapping["enterprise_id"] = enterpriseId
	stmt, err := s.DB.Template("t_members/select-member-by-userids", &mapping)
	if err != nil {
		return nil, err
	}
	list := &[]sensitive.UserInfo{}
	stmtErr := stmt.NamedSelect(c, list, mapping)
	if stmtErr != nil {
		return nil, stmtErr
	}
	return list, nil
}
