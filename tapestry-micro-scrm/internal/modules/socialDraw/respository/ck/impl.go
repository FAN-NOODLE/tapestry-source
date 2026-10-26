package ck

import (
	"fmt"
	"sync"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"

	"github.com/jmoiron/sqlx"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/modules/socialDraw"
)

// GroupUserUvStatistics 群基本用户统计
func (repo GroupChatRepo) GroupUserUvStatistics(ctx ctx.Context, param socialDraw.GroupInfoReq) (*socialDraw.WorkGroupChatStatistic, error) {
	res := &socialDraw.WorkGroupChatStatistic{}
	//留存-新增-退群
	noLayUvQuery := fmt.Sprintf(`select countIf(1, operatime >= %+v and status = 1) as add_uv,countIf(1, operatime <= %+v and status = 2) as del_uv,
		countIf(1,operatime>= %+v and status=2) as new_del
		from (select user_id,toDate(opera_time) day, argMax(status, opera_time) as status, max(opera_time) as operatime
      from micro_scrm.group_chat_member_log where opera_time <= %+v and group_chat_id = :cid  group by user_id,day)`, param.MonthTS, param.ETS, param.STS, param.ETS)
	if err := repo.DB.NamedGet(ctx, res, noLayUvQuery, param); err != nil {
		log.Err(ctx, err).Msgf("exec sql [%+v] error", noLayUvQuery)
		return nil, err
	}

	onlyLayUvQuery := fmt.Sprintf(`select countIf(1, status = 1) as lay_uv from (select user_id, argMax(status, opera_time) as status, max(opera_time) as operatime
		from micro_scrm.group_chat_member_log where opera_time <= %+v and group_chat_id = :cid group by user_id)`, param.ETS)
	if err := repo.DB.NamedGet(ctx, res, onlyLayUvQuery, param); err != nil {
		log.Err(ctx, err).Msgf("exec sql [%+v] error", onlyLayUvQuery)
	}
	return res, nil
}

// GroupMemberStatistics 会员统计
func (repo GroupChatRepo) GroupMemberStatistics(ctx ctx.Context, param socialDraw.GroupInfoReq) (*socialDraw.WorkGroupChatStatistic, error) {
	res := &socialDraw.WorkGroupChatStatistic{}
	sql := `select uniqExactIf(user_id,grade!='' ) as member_uv,uniqExactIf(user_id,grade='') as general_uv,uniqExactIf(user_id,(grade='白金卡' or grade='金卡' or grade='名媛' or grade='GC')) as glod_down from (select l.*,u.grade
		from (select user_id
			  from (select argMax(status, opera_time) as sta, user_id
					from micro_scrm.group_chat_member_log
					where opera_time <= ?
					  and group_chat_id = ?
					  and type = 2
					group by user_id) as a
			  where a.sta = 1) as l
				 left join micro_scrm_enterprise_%d.wx_user_corp_external as u on l.user_id = u.openid)`
	if err := repo.DB.DB.Get(ctx, res, fmt.Sprintf(sql, param.EnterpriseID), param.ETS, param.CID); err != nil {
		return nil, err
	}
	return res, nil
}

func (repo GroupChatRepo) GroupMemberStatisticsForCrm(ctx ctx.Context, param socialDraw.GroupInfoReq) (*socialDraw.WorkGroupChatStatistic, error) {
	res := &socialDraw.WorkGroupChatStatistic{}
	sql := `select uniqIf(user_id,union_code  !='') as member_uv,
uniqIf(user_id,union_code = '') as general_uv 
from (select l.*,u.union_code
		from (select user_id
			  from (select argMax(status, opera_time) as sta, user_id
					from micro_scrm.group_chat_member_log
					where opera_time <= ?
					  and group_chat_id = ?
					  and type = 2
					group by user_id) as a
			  where a.sta = 1) as l
				 left join micro_scrm_enterprise_%d.wx_user_corp_external as u on l.user_id = u.openid)`
	if err := repo.DB.DB.Get(ctx, res, fmt.Sprintf(sql, param.EnterpriseID), param.ETS, param.CID); err != nil {
		return nil, err
	}
	return res, nil
}

// NewGroupUserChat 新入群用户图表
func (repo GroupChatRepo) NewGroupUserChat(ctx ctx.Context, param socialDraw.GroupInfoReq) ([]*socialDraw.Chart, error) {
	var res []*socialDraw.Chart
	sql := `select count(1) as value,name from (
		select l.*,if(grade='','非会员',u.grade) as name
		from (select user_id
			  from (select argMax(status, opera_time) as sta, user_id
					from micro_scrm.group_chat_member_log
					where opera_time >= ?
					  and opera_time <= ?
					  and group_chat_id = ?
					  and type = 2
					group by user_id) as a
			  where a.sta = 1) as l
				 left join micro_scrm_enterprise_%d.wx_user_corp_external as u on l.user_id = u.openid) group by name`
	if err := repo.DB.DB.Select(ctx, &res, fmt.Sprintf(sql, param.EnterpriseID), param.STS, param.ETS, param.CID); err != nil {
		return nil, err
	}
	return res, nil
}

func (repo GroupChatRepo) NewGroupUserChatForCrm(ctx ctx.Context, param socialDraw.GroupInfoReq) ([]*socialDraw.Chart, error) {
	var res []*socialDraw.Chart
	sql := `select count(distinct user_id ) as value,name from (
		select l.*,if(union_code  = '','非会员','会员') as name
		from (select user_id
			  from (select argMax(status, opera_time) as sta, user_id
					from micro_scrm.group_chat_member_log
					where opera_time >= ?
					  and opera_time <= ?
					  and group_chat_id = ?
					  and type = 2
					group by user_id) as a
			  where a.sta = 1) as l
				 left join micro_scrm_enterprise_%d.wx_user_corp_external as u on l.user_id = u.openid) group by name`
	if err := repo.DB.DB.Select(ctx, &res, fmt.Sprintf(sql, param.EnterpriseID), param.STS, param.ETS, param.CID); err != nil {
		return nil, err
	}
	return res, nil
}

// ValidGroupUserChat 留存群用户图表
func (repo GroupChatRepo) ValidGroupUserChat(ctx ctx.Context, param socialDraw.GroupInfoReq) ([]*socialDraw.Chart, error) {
	var res []*socialDraw.Chart
	sql := `select count(1) as value,name from (
		select l.*,if(grade='','非会员',u.grade) as name
		from (select user_id
			  from (select argMax(status, opera_time) as sta, user_id
					from micro_scrm.group_chat_member_log
					where 
					  opera_time <= ?
					  and group_chat_id = ?
					  and type = 2
					group by user_id) as a
			  where a.sta = 1) as l
				 left join micro_scrm_enterprise_%d.wx_user_corp_external as u on l.user_id = u.openid) group by name`
	if err := repo.DB.DB.Select(ctx, &res, fmt.Sprintf(sql, param.EnterpriseID), param.ETS, param.CID); err != nil {
		return nil, err
	}
	return res, nil
}

func (repo GroupChatRepo) ValidGroupUserChatForCrm(ctx ctx.Context, param socialDraw.GroupInfoReq) ([]*socialDraw.Chart, error) {
	var res []*socialDraw.Chart
	sql := `select count(distinct user_id ) as value,name from (
		select l.*,if(union_code  = '' ,'非会员','会员')  as name
		from (select user_id
			  from (select argMax(status, opera_time) as sta, user_id
					from micro_scrm.group_chat_member_log
					where 
					  opera_time <= ?
					  and group_chat_id = ?
					  and type = 2
					group by user_id) as a
			  where a.sta = 1) as l
				 left join micro_scrm_enterprise_%d.wx_user_corp_external as u on l.user_id = u.openid) group by name`
	if err := repo.DB.DB.Select(ctx, &res, fmt.Sprintf(sql, param.EnterpriseID), param.ETS, param.CID); err != nil {
		return nil, err
	}
	return res, nil
}

// GetWorkGroupStatics 企微社群聚合
func (repo GroupChatRepo) GetWorkGroupStatics(ctx ctx.Context, req socialDraw.WorkGroupChatStaticsRequest) (map[int64]*socialDraw.WorkGroupChatStatistic, error) {
	var (
		wg               sync.WaitGroup
		staticsNoLayUv   []*socialDraw.WorkGroupChatResponse
		staticsOnlyLayUv []*socialDraw.WorkGroupChatResponse
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		noLayUvQuery, args, err := sqlx.In(fmt.Sprintf(`select countIf(1, operatime >= %+v and status = 1) as add_uv,countIf(1, operatime <= %+v and status = 2) as del_uv, group_chat_id id 
		from (select user_id,toDate(opera_time) day, argMax(status, opera_time) as status, max(opera_time) as operatime,group_chat_id
      	from micro_scrm.group_chat_member_log where opera_time <= %+v and group_chat_id in (?)  group by group_chat_id,user_id,day) group by group_chat_id`,
			req.StartAt, req.EndAt, req.EndAt), req.GroupChatID)
		if err != nil {
			return
		}
		if err := repo.DB.Select(ctx, &staticsNoLayUv, noLayUvQuery, args...); err != nil {
			log.Err(ctx, err).Msgf("exec sql [%+v] error", noLayUvQuery)
			return
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		onlyLayUvQuery, args, err := sqlx.In(fmt.Sprintf(`select countIf(1, status = 1) as lay_uv,group_chat_id id from (select user_id, argMax(status, opera_time) as status, max(opera_time) as operatime,group_chat_id 
			from micro_scrm.group_chat_member_log where opera_time <= %+v and group_chat_id in (?) group by group_chat_id, user_id) group by group_chat_id`,
			req.EndAt), req.GroupChatID)
		if err != nil {
			return
		}
		if err := repo.DB.Select(ctx, &staticsOnlyLayUv, onlyLayUvQuery, args...); err != nil {
			log.Err(ctx, err).Msgf("exec sql [%+v] error", onlyLayUvQuery)

			return
		}
	}()

	wg.Wait()

	var data = make(map[int64]*socialDraw.WorkGroupChatStatistic)
	for _, v := range staticsNoLayUv {
		data[v.ID] = &v.WorkGroupChatStatistic
	}
	for _, v := range staticsOnlyLayUv {
		if _, ok := data[v.ID]; !ok {
			data[v.ID] = &v.WorkGroupChatStatistic
		}
		data[v.ID].ValidUv = v.ValidUv
	}

	return data, nil
}

// GetWorkGroupMembersList 获取未入群会员列表 会员列表 非会员列表
func (repo GroupChatRepo) GetWorkGroupMembersList(ctx ctx.Context, req socialDraw.WorkGroupChatMembersRequest) ([]*socialDraw.WorkGroupChatMemberList, int, error) {

	var data []*socialDraw.WorkGroupChatMemberList
	var argsList []interface{}
	var argCount []interface{}
	sql := ""
	ctsql := ""
	if req.ListType == 0 {
		fd := " 1=1 "
		if req.Friend == 1 {
			fd = " friend>0 "
		}
		if req.Friend == 2 {
			fd = " friend=0 "
		}
		sql = "select %s from (select a.*,times from (select openid as user_id,nick,name,avatar,unionid,ifnull(r.work_external_user_id,0) as friend from micro_scrm_enterprise_%d.wx_user_corp_external as u any left join (select  work_external_user_id from micro_scrm.work_external_user_relation where enterprise_id=? and relation_status=0 and member_id=?) as r on r.work_external_user_id=u.id where ba_code = ? AND ((grade = '白金卡') OR (grade = '金卡') OR (grade = '名媛') OR (grade = 'GC')) and openid NOT IN(SELECT user_id FROM(SELECT argMax(status, opera_time) AS sta,user_id FROM micro_scrm.group_chat_member_log WHERE (group_chat_id = ?) AND (type = 2) AND (enterprise_id = ?) GROUP BY user_id) AS a WHERE a.sta = 1) ) a any left join (select argMax(arrayReduce('max',arr),created_at) as times,union_id from (select array(pin_pai_buy_last_time,online_buy_last_time,offline_buy_time,tian_mao_buy_last_time,jing_dong_buy_last_time,mini_buy_last_time) as arr,union_id,created_at from micro_scrm_enterprise_%d.campaign_tool_cron_data) group by union_id) as r on a.unionid = r.union_id) where %s "
		ctsql = fmt.Sprintf(sql, "count(1)", req.EnterpriseID, req.EnterpriseID, fd)
		sql = fmt.Sprintf(sql+"order by times desc limit ?,?", "user_id,nick,name,avatar,unionid,friend,times as ck_time ", req.EnterpriseID, req.EnterpriseID, fd)
		argsList = append(argsList, req.EnterpriseID, req.MemberID, req.MarsID, req.CID, req.EnterpriseID, req.Offset, req.Size)
		argCount = append(argCount, req.EnterpriseID, req.MemberID, req.MarsID, req.CID, req.EnterpriseID)
	} else if req.ListType == 1 || req.ListType == 2 {
		sql = `select %s
		from (select argMax(status, opera_time) as sta,
					 user_id,
					 max(opera_time)            as times,
					 max(join_scene)            as join_scene
			  from micro_scrm.group_chat_member_log
			  where group_chat_id = ?
				and opera_time <= ?
				and enterprise_id = ?
				and type = 2
			  group by user_id) as a
				any left join (select id,openid, nick, name, avatar, unionid, grade from micro_scrm_enterprise_%d.wx_user_corp_external) as u on a.user_id = u.openid
		where a.sta = 1 `

		if req.ListType == 1 {
			sql += " and grade !='' "
		} else {
			sql += " and grade ='' "
		}
		if req.Friend == 1 {
			sql += fmt.Sprintf(" and u.id in (select work_external_user_id from micro_scrm.work_external_user_relation where enterprise_id = %d and relation_status = 0 and member_id = %d)", req.EnterpriseID, req.MemberID)
		}
		if req.Friend == 2 {
			sql += fmt.Sprintf(" and u.id not in (select work_external_user_id from micro_scrm.work_external_user_relation where enterprise_id = %d and relation_status = 0 and member_id = %d)", req.EnterpriseID, req.MemberID)
		}
		ctsql = fmt.Sprintf(sql, "count(1)", req.EnterpriseID)
		sql = fmt.Sprintf(sql+"limit ?,?", "user_id,times,join_scene,ifnull(nick,'') as nick, ifnull(name,'') as name, ifnull(avatar,'') as avatar, ifnull(unionid,'') as unionid,ifnull(u.id,0) as id", req.EnterpriseID)
		argsList = append(argsList, req.CID, req.EndAt, req.EnterpriseID, req.Offset, req.Size)
		argCount = append(argCount, req.CID, req.EndAt, req.EnterpriseID)
	} else {
		return nil, 0, nil
	}

	if err := repo.DB.DB.Select(ctx, &data, sql, argsList...); err != nil {
		return nil, 0, err
	}
	var exid []int64
	for _, item := range data {
		if req.ListType > 0 {
			item.Friend = 0
			if req.Friend == 0 && item.UID > 0 {
				exid = append(exid, item.UID)
			}
			if req.Friend == 1 {
				item.Friend = 1
			}
		}
		item.Type = 2
		if req.ListType == 0 {
			if item.TmpCkTime != "" {
				if s, err := time.Parse("20060102", item.TmpCkTime); err == nil && s.Unix() > 0 {
					item.Time = time.Unix(s.Unix(), 0).Format("2006-01-02")
				}
			}
		} else {
			item.Time = time.Unix(item.TmpTime, 0).Format("2006-01-02 15:04")
		}
	}
	if len(exid) > 0 {
		var uid []int64
		fdmp := make(map[int64]bool)
		ex := fmt.Sprintf("select work_external_user_id from micro_scrm.work_external_user_relation where enterprise_id = %d and relation_status = 0 and member_id = %d", req.EnterpriseID, req.MemberID)
		if err := repo.DB.DB.Select(ctx, &uid, ex); err != nil {
			return nil, 0, err
		}
		for _, u := range uid {
			fdmp[u] = true
		}
		if len(fdmp) > 0 {
			for _, item := range data {
				if _, ok := fdmp[item.UID]; ok {
					item.Friend = 1
				}
			}
		}
	}
	var count int
	if err := repo.DB.DB.Get(ctx, &count, ctsql, argCount...); err != nil {
		return nil, 0, err
	}

	return data, count, nil
}
