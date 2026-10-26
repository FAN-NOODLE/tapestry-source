package mysql

import (
	"fmt"
	"strings"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/chat"
	"micro-scrm/internal/pkg/sdk/work"
)

type ChatMsgRepo struct {
	*repos.SQLRepo
}

// 记录是否存在
func (repo ChatMsgRepo) Exists(c ctx.Context, enterpriseId int64, msgId string) (bool, error) {
	params := struct {
		MsgId        string `json:"msg_id" db:"msg_id"`
		EnterpriseId int64  `json:"enterprise_id" db:"enterprise_id"`
	}{
		MsgId:        msgId,
		EnterpriseId: enterpriseId,
	}

	stmt, err := repo.DB.Template("t_chat_msg/exists")
	if err != nil {
		return false, err
	}

	res := struct {
		Id int64 `json:"id" db:"id"`
	}{}
	getErr := stmt.NamedGet(c, &res, params)
	if getErr != nil {
		return false, getErr
	}

	if res.Id == 0 {
		return false, nil
	}

	return true, nil
}

// LastRow 最后一条记录
func (repo ChatMsgRepo) LastRow(c ctx.Context, result interface{}, enterpriseId int64) error {
	stmt, err := repo.DB.Template("t_chat_msg/last-row")
	if err != nil {
		return err
	}

	params := struct {
		EnterpriseId int64 `json:"enterprise_id" db:"enterprise_id"`
	}{
		EnterpriseId: enterpriseId,
	}

	return stmt.NamedGet(c, result, params)
}

// ListByIDAsc 根据ID向下找附近的记录
func (repo ChatMsgRepo) ListByID(c ctx.Context, cond chat.MsgListByIDRequest) ([]chat.Msg, error) {
	option := &repos.ListOption{
		TplPrefix: "t_chat_msg/list-by-id",
		Fn: func(context ctx.Context, i int64) interface{} {
			return &[]chat.Msg{}
		},
	}
	res, err := repo.List(c, cond, option)

	return *res.(*[]chat.Msg), err
}

type ChatMemberRepo struct {
	*repos.SQLRepo
}

// FindMember 找员工记录
func (repo ChatMemberRepo) FindMember(c ctx.Context, cond interface{}, result interface{}) error {
	stmt, err := repo.DB.Template("t_chat_member/find")
	if err != nil {
		return err
	}

	return stmt.NamedGet(c, result, cond)
}

type ReceiveUserRepo struct {
	*repos.SQLRepo
}

type ChatObjectRepo struct {
	*repos.SQLRepo
}

type MemberToGroupRepo struct {
	*repos.SQLRepo
}

type MsgMediaRepo struct {
	*repos.SQLRepo
}

type GroupChatRepo struct {
	*repos.SQLRepo
}

type GroupChatMemberRepo struct {
	*repos.SQLRepo
}

// Find 找一条记录
func (repo ChatObjectRepo) Find(c ctx.Context, cond interface{}, result interface{}) error {
	stmt, err := repo.DB.Template("t_chat_object/find")
	if err != nil {
		return err
	}

	return stmt.NamedGet(c, result, cond)
}

// GetNotMediaList 获取为转换媒体文件的聊天记录
func (repo ChatMsgRepo) GetNotMediaList(ctx ctx.Context, req []string, enterpriseId int64, StartAt int64) (*[]chat.Msg, error) {
	mapping := make(map[string]interface{})
	mapping["type"] = req
	mapping["enterprise_id"] = enterpriseId
	mapping["StartAt"] = StartAt
	stmt, err := repo.DB.Template("t_chat_msg/get-not-media-list", &mapping)
	if err != nil {
		return nil, err
	}
	list := &[]chat.Msg{}
	stmtErr := stmt.NamedSelect(ctx, list, mapping)
	if stmtErr != nil {
		return nil, stmtErr
	}
	return list, nil
}

// GetMembersList 获取指定userids的用户信息
func (repo ChatMsgRepo) GetMembersList(c ctx.Context, req []string, enterpriseId int64) (*[]chat.UserInfo, error) {
	mapping := make(map[string]interface{})
	mapping["user_ids"] = req
	mapping["enterprise_id"] = enterpriseId
	stmt, err := repo.DB.Template("t_members/select-member-by-userids", &mapping)
	if err != nil {
		return nil, err
	}
	list := &[]chat.UserInfo{}
	stmtErr := stmt.NamedSelect(c, list, mapping)
	if stmtErr != nil {
		return nil, stmtErr
	}
	return list, nil
}

// GetExternalUsersList 获取指定userids的外部联系人信息
func (repo ChatMsgRepo) GetExternalUsersList(c ctx.Context, req []string, enterpriseId int64) (*[]chat.UserInfo, error) {
	mapping := make(map[string]interface{})
	mapping["user_ids"] = req
	mapping["enterprise_id"] = enterpriseId
	stmt, err := repo.DB.Template("t_work_external_user/select-user-by-userids", &mapping)
	if err != nil {
		return nil, err
	}
	list := &[]chat.UserInfo{}
	stmtErr := stmt.NamedSelect(c, list, mapping)
	if stmtErr != nil {
		return nil, stmtErr
	}
	return list, nil
}

// GetMemberByGroup 获取群聊里的员工列表
func (repo ChatObjectRepo) GetMemberByGroup(c ctx.Context, req string, enterpriseId int64) (*[]chat.GroupMemberInfo, error) {
	mapping := make(map[string]interface{})
	mapping["chat_id"] = req
	mapping["enterprise_id"] = enterpriseId
	stmt, err := repo.DB.Template("t_work_group_chat/get-member-by-group", &mapping)
	if err != nil {
		return nil, err
	}
	list := &[]chat.GroupMemberInfo{}
	stmtErr := stmt.NamedSelect(c, list, mapping)
	if stmtErr != nil {
		return nil, stmtErr
	}
	return list, nil
}

// ListByFileId 根据文件id列表找记录
func (repo MsgMediaRepo) ListByFileId(c ctx.Context, enterpriseID int64, sdkFileidList []string) ([]chat.MsgMedia, error) {
	var list []chat.MsgMedia
	cond := map[string]interface{}{
		"enterprise_id":    enterpriseID,
		"sdk_file_id_list": sdkFileidList,
	}
	stmt, err := repo.DB.Template("t_chat_msg_media/list-by-fileid", &cond)
	if err != nil {
		return nil, err
	}
	errStmt := stmt.NamedSelect(c, &list, cond)
	if errStmt != nil {
		return nil, err
	}

	return list, nil
}

// selectReceiveUserByMemberUserIDs 根据员工userid，返回员工消息IDs
func (repo ReceiveUserRepo) SelectReceiveUserByMemberUserIDs(c ctx.Context, memberUserIDs []string) ([]chat.SearchReceiveResponse, error) {
	mapping := make(map[string]interface{})
	mapping["memberUserIDs"] = memberUserIDs
	stmt, err := repo.DB.Template("t_chat_receive_user/select-list-by-touser", &mapping)
	if err != nil {
		return nil, err
	}
	list := &[]chat.SearchReceiveResponse{}
	stmtErr := stmt.NamedSelect(c, list, mapping)
	if stmtErr != nil {
		return nil, stmtErr
	}
	return *list, nil
}

// FindWorkGroupChatDetail 群聊详情
func (repo GroupChatRepo) FindWorkGroupChatDetail(ctx ctx.Context, req chat.WorkGroupChat) (*chat.WorkGroupChat, error) {
	mapping := make(map[string]interface{})
	mapping["cond"] = req
	stmt, err := repo.DB.Template("t_work_group_chat/detail", &mapping)
	if err != nil {
		return nil, err
	}
	chat := &chat.WorkGroupChat{}
	sErr := stmt.NamedGet(ctx, chat, req)
	if sErr != nil {
		return nil, sErr
	}
	chat.EnterpriseID = req.EnterpriseID
	return chat, nil
}

// BatchUpdateWorkGroupChatMembers 批量更新群聊表
func (repo GroupChatMemberRepo) BatchUpdateWorkGroupChatMembers(ctx ctx.Context, where map[string]interface{}, update chat.WorkGroupChatMember) (int64, error) {
	stmt, err := repo.DB.Template("t_work_group_chat_member/batch-update", &where)
	if err != nil {
		return 0, err
	}
	res, sErr := stmt.NamedExec(ctx, update)
	if sErr != nil {
		return 0, sErr
	}
	return res.RowsAffected()
}

// BatchUpdateExternalWorkGroupChatMembersInfo 批量更新外部联系人群成员信息
func (repo GroupChatMemberRepo) BatchUpdateExternalWorkGroupChatMembersInfo(c ctx.Context, EnterpriseId int64,
	groupChatId int64, chatMembers []work.ExternalGroupMemberList) (int64, error) {

	var (
		nickNameThen, nameThen string
		joinThen, inviteThen   string
		updateStrs             []string
	)
	for _, v := range chatMembers {
		if v.GroupNickname != "" {
			nickNameThen += fmt.Sprintf(" when '%+v' then '%+v' ", v.UserID, v.GroupNickname)
		}
		if v.JoinScene == 1 {
			joinThen += fmt.Sprintf(" when '%+v' then '%+v' ", v.UserID, v.JoinScene)
		}
		if v.Invitor.UserID != "" {
			inviteThen += fmt.Sprintf(" when '%+v' then '%+v' ", v.UserID, v.Invitor.UserID)
		}
		nameThen += fmt.Sprintf(" when '%+v' then '%+v' ", v.UserID, v.Name)
	}
	if nickNameThen != "" {
		updateStrs = append(updateStrs, fmt.Sprintf(`nickname= case userid %s end`, nickNameThen+fmt.Sprintf(" else '' ")))
	}
	if nameThen != "" {
		updateStrs = append(updateStrs, fmt.Sprintf(`name= case userid %s end`, nameThen+fmt.Sprintf(" else name ")))
	}
	if joinThen != "" {
		updateStrs = append(updateStrs, fmt.Sprintf(`join_scene= case userid %s end`, joinThen+fmt.Sprintf(" else join_scene ")))
	}
	if inviteThen != "" {
		updateStrs = append(updateStrs, fmt.Sprintf(`invite_userid= case userid %s end`, inviteThen+fmt.Sprintf(" else '' ")))
	}
	if len(updateStrs) == 0 {
		return 0, nil
	}
	rows, err := repo.DB.DB.Exec(c, fmt.Sprintf(`update work_group_chat_member set %s where status=1 and enterprise_id=? and group_chat_id = ?`,
		strings.Join(updateStrs, ",")), EnterpriseId, groupChatId)
	if err != nil {
		return 0, err
	}
	return rows.RowsAffected()
}
