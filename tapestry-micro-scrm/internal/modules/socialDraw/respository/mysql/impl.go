package mysql

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/repos"
	"micro-scrm/internal/modules/socialDraw"
)

type GroupChatRepo struct {
	*repos.SQLRepo
}

// FindWorkGroupChatList 社群列表
func (repo GroupChatRepo) FindWorkGroupChatList(ctx ctx.Context, req socialDraw.WorkGroupListRequest) ([]*socialDraw.WorkGroupChatResponse, error) {
	stmt, err := repo.DB.Template("t_work_group_chat/list", repo.Cond2Args(req))
	if err != nil {
		return nil, err
	}

	var chats []*socialDraw.WorkGroupChatResponse
	if err := stmt.NamedSelect(ctx, &chats, req); err != nil {
		return nil, err
	}
	return chats, nil
}

// WorkGroupChatSetting 社群类型管理
func (repo GroupChatRepo) WorkGroupChatSetting(ctx ctx.Context, req socialDraw.WorkGroupChatSettingRequest) error {
	stmt, err := repo.DB.Template("t_work_group_chat/setting", repo.Cond2Args(req))
	if err != nil {
		return err
	}

	if _, err := stmt.NamedExec(ctx, req); err != nil {
		return err
	}
	return nil
}

// FindWorkGroupTypeList 社群类型列表
func (repo GroupChatRepo) FindWorkGroupTypeList(ctx ctx.Context, req socialDraw.WorkGroupChatType) ([]*socialDraw.WorkGroupChatType, error) {
	stmt, err := repo.DB.Template("t_work_group_chat_type/type_list", repo.Cond2Args(req))
	if err != nil {
		return nil, err
	}

	var data []*socialDraw.WorkGroupChatType
	if err := stmt.NamedSelect(ctx, &data, req); err != nil {
		return nil, err
	}
	return data, nil
}

func (repo GroupChatRepo) CountWorkGroupByType(ctx ctx.Context, req socialDraw.WorkGroupChatType) (map[int64]int, error) {
	// 统计社群类型使用数量
	var (
		countData []*socialDraw.WorkGroupChatType
		countMap  = make(map[int64]int)
	)
	stmt, err := repo.DB.Template("t_work_group_chat_type/count-group-by-type")
	if err != nil {
		return countMap, nil
	}
	if err := stmt.Select(ctx, &countData, req.EnterpriseID); err != nil {
		return countMap, nil
	}
	for _, v := range countData {
		countMap[v.ID] = v.GroupChatNum
	}
	return countMap, nil
}

// ManageWorkGroupType 社群类型管理
func (repo GroupChatRepo) ManageWorkGroupType(ctx ctx.Context, req socialDraw.WorkGroupChatType) (*socialDraw.WorkGroupChatType, error) {
	if req.Opt == 1 {
		stmt, err := repo.DB.Template("t_work_group_chat_type/insert")
		if err != nil {
			return nil, err
		}
		result, err := stmt.NamedExec(ctx, req)
		if err != nil {
			return nil, err
		}
		req.ID, _ = result.LastInsertId()
		return &req, err
	}
	if req.Opt == 2 {
		stmt, err := repo.DB.Template("t_work_group_chat_type/update")
		if err != nil {
			return nil, err
		}
		_, err = stmt.NamedExec(ctx, req)
		if err != nil {
			return nil, err
		}
		return &req, nil
	}
	if req.Opt == 3 {
		stmt, err := repo.DB.Template("t_work_group_chat_type/delete")
		if err != nil {
			return nil, err
		}
		_, err = stmt.Exec(ctx, req.ID, req.EnterpriseID)
		return nil, err
	}
	return nil, nil
}

// GetWorkGroupTypeInfo get work group type info
func (repo GroupChatRepo) GetWorkGroupTypeInfo(ctx ctx.Context, req socialDraw.WorkGroupChatType) (*socialDraw.WorkGroupChatType, error) {
	stmt, err := repo.DB.Template("t_work_group_chat_type/type_info", repo.Cond2Args(req))
	if err != nil {
		return nil, err
	}

	var data = &socialDraw.WorkGroupChatType{}
	if err := stmt.NamedGet(ctx, &data, req); err != nil {
		return nil, err
	}
	return data, nil
}

// WorkGroupFilters 社群列表筛选数据
func (repo GroupChatRepo) WorkGroupFilters(ctx ctx.Context, req socialDraw.WorkGroupFilterRequest) ([]*socialDraw.WorkGroupFilterResponse, error) {
	var data []*socialDraw.WorkGroupFilterResponse
	switch req.FiledName {
	case "chat_owner":
		if err := repo.DB.Select(ctx, &data, `select wgc.owner_member_id id, m.name from work_group_chat wgc inner join members m on wgc.owner_member_id = m.id 
			where wgc.dissolve_time = 0 and wgc.status = 1 and wgc.enterprise_id = ? group by owner_member_id`, req.EnterpriseID); err != nil {
			return nil, err
		}
	case "channels":
		fallthrough
	default:
		stmt, err := repo.DB.Template("t_work_group_chat/filter-channel-list")
		if err != nil {
			return nil, err
		}
		if err := stmt.NamedSelect(ctx, &data, req); err != nil {
			return nil, err
		}
	}

	return data, nil
}

// GetWorkGroupMembers 按类型获取企微群成员明细
func (repo GroupChatRepo) GetWorkGroupMembers(ctx ctx.Context, req socialDraw.WorkGroupChatMembersRequest) ([]*socialDraw.WorkGroupChatMember, int, error) {
	reqData := repo.Cond2Args(req)
	stmt, err := repo.DB.Template("t_work_group_chat_member/members-by-type", reqData)
	if err != nil {
		return nil, 0, err
	}

	var data []*socialDraw.WorkGroupChatMember
	if err := stmt.NamedSelect(ctx, &data, req); err != nil {
		return nil, 0, err
	}

	var count struct {
		TotalCount int `db:"total_count"`
	}
	countStmt, _ := repo.DB.Template("t_work_group_chat_member/members-count-by-type", reqData)
	_ = countStmt.NamedGet(ctx, &count, req)

	return data, count.TotalCount, nil
}

// CheckWorkGroupChatTypeUsed 检测群类型是否被使用
func (repo GroupChatRepo) CheckWorkGroupChatTypeUsed(ctx ctx.Context, enterpriseId int64, chatTypeId int64) (bool, error) {
	stmt, err := repo.DB.Template("t_work_group_chat_type/check-type-used")
	if err != nil {
		return false, err
	}

	var count int
	if err := stmt.Get(ctx, &count, enterpriseId, chatTypeId); err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetWorkGroupExternalUsers 获取社群外部联系人信息
func (repo GroupChatRepo) GetWorkGroupExternalUsers(ctx ctx.Context, enterpriseId int64, userIds []string) (map[string]*socialDraw.WxUserCorpExternal, error) {
	req := make(map[string]interface{})
	req["EnterpriseID"] = enterpriseId
	req["UserIds"] = userIds
	stmt, err := repo.DB.Template("t_work_group_chat_member/select-user-from-corp-external", repo.Cond2Args(req))
	if err != nil {
		return nil, err
	}

	var result []*socialDraw.WxUserCorpExternal
	if err := stmt.Select(ctx, &result); err != nil {
		return nil, err
	}

	var data = make(map[string]*socialDraw.WxUserCorpExternal)
	for _, v := range result {
		data[v.OpenID] = v
	}

	return data, nil
}

// GetWorkGroupChatListFilters 社群筛选列表查询
func (repo GroupChatRepo) GetWorkGroupChatListFilters(ctx ctx.Context, req socialDraw.WorkGroupListRequest) []string {
	stmt, err := repo.DB.Template("t_work_group_chat/filter-operator-department", repo.Cond2Args(req))
	if err != nil {
		return nil
	}
	var memberUserIds []string
	if err := stmt.Select(ctx, &memberUserIds, req.EnterpriseID); err != nil {
		return nil
	}
	return memberUserIds
}

func (repo GroupChatRepo) GetMarsID(ctx ctx.Context, enterpriseId int64, uid string) (string, error) {

	var marId string
	if err := repo.DB.DB.Get(ctx, &marId, fmt.Sprintf("select mars_id from micro_scrm_enterprise_%d.members where userid='%s'", enterpriseId, uid)); err != nil {
		return "", err
	}
	return marId, nil
}

func (repo GroupChatRepo) GetWorkGroupMembersInfo(ctx ctx.Context, req socialDraw.WorkGroupChatMembersRequest) ([]*socialDraw.WorkGroupChatMemberList, int, error) {

	var data []*socialDraw.WorkGroupChatMemberList
	var argList []interface{}
	var argCount []interface{}
	sql := "select %s from micro_scrm.work_group_chat_member as m left join micro_scrm_enterprise_%d.wx_user_corp_external as u on m.userid=u.openid where m.leave_time = 0 and m.status = 1 and m.deleted_at = 0 and m.group_chat_id=? "
	if req.ListType == 1 {
		sql = sql + " and m.act_num>0"
	}
	argList = append(argList, req.CID)
	argCount = append(argCount, req.CID)
	if req.Friend > 0 {
		fd := " and (u.id in (select work_external_user_id  from micro_scrm.work_external_user_relation where enterprise_id=? and member_id=? and relation_status=0) or member_id>0) "
		if req.Friend == 2 {
			fd = " and (u.id not in (select work_external_user_id  from micro_scrm.work_external_user_relation where enterprise_id=? and member_id=? and relation_status=0) or (isnull(u.id) and member_id=0 )) "
		}
		sql = sql + fd
		argList = append(argList, req.EnterpriseID, req.MemberID)
		argCount = append(argCount, req.EnterpriseID, req.MemberID)
	}
	order := " order by join_time desc "
	if req.ListType > 0 {
		order = " order by act_num desc "
	}

	list := fmt.Sprintf(sql+order+" limit ?,?", "userid as user_id,ifnull(u.id, 0) as id ,m.type,join_time as times,nickname as nick,m.name,member_id,join_scene,ifnull(u.avatar,'') as avatar,act_num", req.EnterpriseID)

	argList = append(argList, req.Offset, req.Size)
	if err := repo.DB.DB.Select(ctx, &data, list, argList...); err != nil {
		return nil, 0, err
	}
	var externalUids []interface{}
	for _, item := range data {

		if req.Friend == 0 {
			if item.UID == 0 {
				item.Friend = 0
			} else {
				externalUids = append(externalUids, item.UID)
			}

			if item.Type == 1 {
				item.Friend = 1
			}
		}

		if req.Friend == 1 {
			item.Friend = 1
		}
		if req.Friend == 2 {
			item.Friend = 0
		}
	}

	if len(externalUids) > 0 {
		var uids []int64
		externalUids = append(externalUids, req.EnterpriseID, req.MemberID)
		if err := repo.DB.DB.Select(ctx, &uids, "select work_external_user_id  from micro_scrm.work_external_user_relation where work_external_user_id in (?"+strings.Repeat(",?", len(externalUids)-3)+") and enterprise_id=? and member_id=? and relation_status=0", externalUids...); err != nil {
			return nil, 0, err
		}
		if len(uids) > 0 {
			umap := make(map[int64]bool)
			for _, r := range uids {
				umap[r] = true
			}
			for _, item := range data {
				if _, ok := umap[item.UID]; ok {
					item.Friend = 1
				}
			}
		}
	}

	var count int
	countQuery := fmt.Sprintf(sql, "count(1)", req.EnterpriseID)
	if err := repo.DB.DB.Get(ctx, &count, countQuery, argCount...); err != nil {
		return nil, 0, err
	}

	return data, count, nil
}

func (repo GroupChatRepo) GetWorkGroupUserActivityInfo(ctx ctx.Context, enterpriseId int64, UserID string, cid int64) ([]*socialDraw.ActInfo, error) {
	var acts []*socialDraw.ActInfo
	var data []byte
	var tmpData socialDraw.ActInfoBase
	if err := repo.DB.DB.Get(ctx, &data, "select act_ids from micro_scrm.work_group_chat_member where group_chat_id=? and userid=? and enterprise_id=?", cid, UserID, enterpriseId); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &tmpData); err != nil {
		return nil, nil
	}
	if tmpData.Fission != nil {
		var fids []interface{}
		fismp := make(map[int64]*socialDraw.ActInfo)
		for _, f := range tmpData.Fission {
			fismp[f.ActivityID] = &socialDraw.ActInfo{JoinTime: f.JoinAt}
			fids = append(fids, f.ActivityID)
		}
		if len(fids) > 0 {
			var f []*socialDraw.ActInfo
			if err := repo.DB.DB.Select(ctx, &f, "select id,name,status,deleted_at,end_time from micro_scrm.fission_activity where id in(?"+strings.Repeat(",?", len(fids)-1)+")", fids...); err != nil {
				return nil, err
			}
			for _, f1 := range f {
				if (f1.DelTime > 0 || f1.ActEndTime < time.Now().Unix()) && f1.Status == 1 {
					f1.Status = 2
				}
				if _, ok := fismp[f1.ID]; ok {
					fismp[f1.ID].Name = f1.Name
					fismp[f1.ID].Status = f1.Status
				}
			}
			for _, fs := range fismp {
				acts = append(acts, fs)
			}
		}
	}
	if tmpData.Market != nil {
		var mids []interface{}
		mkmp := make(map[int64]*socialDraw.ActInfo)
		for _, f := range tmpData.Market {
			mkmp[f.ActivityID] = &socialDraw.ActInfo{JoinTime: f.JoinAt}
			mids = append(mids, f.ActivityID)
		}
		if len(mids) > 0 {
			var m []*socialDraw.ActInfo
			if err := repo.DB.DB.Select(ctx, &m, fmt.Sprintf("select id,title as name ,online as status,deleted_at,end_time from micro_scrm_enterprise_%d.activity where id in(?"+strings.Repeat(",?", len(mids)-1)+")", enterpriseId), mids...); err != nil {
				return nil, err
			}
			for _, m1 := range m {
				if (m1.DelTime > 0 || m1.ActEndTime < time.Now().Unix()) && m1.Status == 1 {
					m1.Status = 2
				}
				if _, ok := mkmp[m1.ID]; ok {
					mkmp[m1.ID].Name = m1.Name
					mkmp[m1.ID].Status = m1.Status
				}
			}
			for _, ms := range mkmp {
				acts = append(acts, ms)
			}
		}
	}
	return acts, nil
}
