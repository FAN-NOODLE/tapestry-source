package mysql

import (
	"fmt"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/sqlplus"
	"micro-scrm/internal/modules/sidebarGroupChat"
)

type SidebarGroupChatRepo struct {
	*repos.SQLRepo
}

// GroupChatInfo 获取社群详情
func (s *SidebarGroupChatRepo) GroupChatInfo(c ctx.Context, groupChatInfoRequest sidebarGroupChat.GroupChatInfoRequest) (*sidebarGroupChat.WorkGroupChat, error) {
	stmt, stmtErr := s.DB.Template("t_sidebar_group_chat/find_group_chat", s.Cond2Args(groupChatInfoRequest))
	if stmtErr != nil {
		return nil, stmtErr
	}
	GroupChatInfoResponse := &sidebarGroupChat.WorkGroupChat{}
	getErr := stmt.NamedGet(c, GroupChatInfoResponse, groupChatInfoRequest)
	if getErr != nil {
		return nil, getErr
	}
	return GroupChatInfoResponse, nil
}

// GroupChatLeaderInfo 获取群主信息
func (s *SidebarGroupChatRepo) GroupChatLeaderInfo(c ctx.Context, groupChatInfoRequest sidebarGroupChat.GroupChatInfoRequest) (*sidebarGroupChat.WorkGroupChat, error) {
	stmt, stmtErr := s.DB.Template("t_sidebar_group_chat/find_group_chat_leader", s.Cond2Args(groupChatInfoRequest))
	if stmtErr != nil {
		return nil, stmtErr
	}
	GroupChatInfoResponse := &sidebarGroupChat.WorkGroupChat{}
	getErr := stmt.NamedGet(c, GroupChatInfoResponse, groupChatInfoRequest)
	if getErr != nil {
		return nil, getErr
	}
	return GroupChatInfoResponse, nil
}

// GroupChatMember 获取社群用户基础信息
func (s *SidebarGroupChatRepo) GroupChatMember(c ctx.Context, groupChatInfoRequest sidebarGroupChat.GroupChatInfoRequest) ([]*sidebarGroupChat.WorkGroupChatMember, error) {
	stmt, err := s.DB.Template("t_sidebar_group_chat/group_chat_member_list")
	if err != nil {
		return nil, err
	}
	ChatMember := []*sidebarGroupChat.WorkGroupChatMember{}
	res, errStmt := stmt.NamedQuery(c, groupChatInfoRequest)
	defer res.Close()
	sqlplusErr := sqlplus.ScanAll(res, &ChatMember)
	fmt.Println(sqlplusErr)
	if errStmt != nil {
		return nil, err
	}
	return ChatMember, nil
}

// GroupChatMember 获取社群用户基础信息
func (s *SidebarGroupChatRepo) GroupChatMemberUnionid(c ctx.Context, groupChatInfoRequest sidebarGroupChat.GroupChatInfoRequest) (*[]sidebarGroupChat.WorkGroupChatMember, error) {
	ChatMember := &[]sidebarGroupChat.WorkGroupChatMember{}
	err := s.DB.Select(c, ChatMember, "SELECT mem.unionid FROM micro_scrm.work_group_chat_member as mem WHERE mem.group_chat_id IN (SELECT micro_scrm.work_group_chat.id FROM micro_scrm.work_group_chat WHERE micro_scrm.work_group_chat.chat_id =?) and mem.`status`=1 and mem.unionid<>''", groupChatInfoRequest.GroupChatId)
	if err != nil {
		return nil, err
	}
	//defer s.DB.Close()
	return ChatMember, nil
}

// GroupChatCorpMember 获取社群外部联系人DL会员等级
func (s *SidebarGroupChatRepo) GroupChatCorpUser(c ctx.Context, uid string, enterpriseId int64) (*[]sidebarGroupChat.WorkGroupCorpUser, error) {
	//var Db *sqlx.DB
	//database, err := sqlx.Open("mysql", "root:ga1k6vcazuagaqm8@tcp(10.99.94.98:30055)/micro_scrm_enterprise_3306120225798016")
	//if err != nil {
	//	fmt.Println("open mysql failed,", err)
	//	return nil, err
	//}
	//
	//Db = database
	//defer s.DB.Close()
	ChatMember := &[]sidebarGroupChat.WorkGroupCorpUser{}
	sql := "SELECT unionid,union_code,grade from micro_scrm_enterprise_%d.wx_user_corp_external as crop where crop.unionid in('%s') and crop.deleted_at=0"
	sql = fmt.Sprintf(sql, enterpriseId, uid)
	dbErr := s.DB.Select(c, ChatMember, sql)
	if dbErr != nil {
		return nil, dbErr
	}
	return ChatMember, nil
}

func (s *SidebarGroupChatRepo) GetCorpUserUnionid(c ctx.Context, openId string, enterpriseId int64) (sidebarGroupChat.WorkGroupCorpUser, error) {
	ChatMember := sidebarGroupChat.WorkGroupCorpUser{}
	sql := "SELECT unionid from micro_scrm_enterprise_%d.wx_user_corp_external  where openid='%s'"
	sql = fmt.Sprintf(sql, enterpriseId, openId)
	dbErr := s.DB.Select(c, ChatMember, sql)
	if dbErr != nil {
		return ChatMember, dbErr
	}
	return ChatMember, nil
}

// GroupChatMemberMsg 群聊天消息
func (s *SidebarGroupChatRepo) GroupChatMemberMsg(c ctx.Context, chatMsgRequest sidebarGroupChat.ChatMsgRequest) (*[]sidebarGroupChat.ChatMsgResponse, error) {
	stmt, err := s.DB.Template("t_sidebar_group_chat/group_chat_msg_list")
	if err != nil {
		return nil, err
	}
	ChatMember := &[]sidebarGroupChat.ChatMsgResponse{}
	res, errStmt := stmt.NamedQuery(c, chatMsgRequest)
	defer res.Close()
	sqlplusErr := sqlplus.ScanAll(res, ChatMember)
	fmt.Println(sqlplusErr)
	if errStmt != nil {
		return nil, err
	}
	return ChatMember, nil
}

// UsersAddEmployees 获取用户添加员工列表
func (s *SidebarGroupChatRepo) UsersAddEmployees(c ctx.Context, enterpriseID int64, workExternalUserID int64, euid int64) (*[]sidebarGroupChat.GroupUsersAddEmployees, error) {
	stmt, err := s.DB.Template("t_sidebar_group_chat/user_add_employees_list")
	if err != nil {
		return nil, err
	}
	ChatMember := &[]sidebarGroupChat.GroupUsersAddEmployees{}
	res, errStmt := stmt.Query(c, euid, workExternalUserID, enterpriseID)
	defer res.Close()
	sqlplusErr := sqlplus.ScanAll(res, ChatMember)
	fmt.Println(sqlplusErr)
	if errStmt != nil {
		return nil, err
	}
	return ChatMember, nil
}
