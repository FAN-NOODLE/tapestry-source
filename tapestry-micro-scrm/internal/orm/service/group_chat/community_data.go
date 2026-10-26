package group_chat

import (
	"fmt"
	"micro-scrm/internal/orm/model/micro_scrm"

	"micro-scrm/internal/orm/model/micro_scrm/group_chat"
)

// 根据ba查询ba所在的群
func (g *GroupChatService) FindGroupsBymemberIds(params GropuNameQueryconditions, BAUserIds []string) (data []group_chat.WorkGroupChatModel, total int64, err error) {
	query := g.Model(&group_chat.WorkGroupChatMemberModel{}).
		Joins("LEFT JOIN work_group_chat b ON work_group_chat_member.group_chat_id = b.id AND b.type = 1 AND b.status = 1 AND b.deleted_at=0").
		Where("work_group_chat_member.enterprise_id=?  AND work_group_chat_member.status=? AND work_group_chat_member.deleted_at=? AND work_group_chat_member.userid IN (?)", params.EnterpriseID, 1, 0, BAUserIds).
		Where("b.id IS NOT NULL AND b.enterprise_id=? AND b.create_time <= ? AND b.deleted_at=?", params.EnterpriseID, params.EndTime, 0).
		Group("b.id,b.name").Order("b.create_time DESC")
	if params.Page.PageNum != 0 && params.Page.PageSize != 0 {
		query.Count(&total).Limit(params.Page.PageSize).Offset((params.Page.PageNum - 1) * params.Page.PageSize)
	} else {
		query.Count(&total)
	}
	err = query.Scan(&data).Error
	return
}

/**
*查询群累计人数
 */
func (g *GroupChatService) FindGroupsMembersByChatIds(chatIds []int64, params GropuNameQueryconditions) (result []group_chat.ChatMembers, err error) {
	var (
		sqlParams = []interface{}{params.EnterpriseID, 0, 2, params.EndTime, micro_scrm.EXTERNAL_CONTACTS_STRING, chatIds}
	)
	sql := `
	SELECT 
		group_chat_id,SUM(number) number 
    FROM
	(SELECT 
	  work_group_chat_member_log.group_chat_id,
	  CASE WHEN work_group_chat_member_log.operation_status = 2 
		THEN 0 ELSE 1 END number 
	FROM work_group_chat_member_log 
	WHERE work_group_chat_member_log.enterprise_id=? AND work_group_chat_member_log.deleted_at = ? AND work_group_chat_member_log.member_type = ? 
	  AND work_group_chat_member_log.created_at <= ?
      AND work_group_chat_member_log.external_userid LIKE ?
	  AND work_group_chat_member_log.group_chat_id IN(?)
	GROUP BY work_group_chat_member_log.group_chat_id,
	  work_group_chat_member_log.union_id ) t GROUP BY group_chat_id 
	`
	err = g.Raw(sql, sqlParams...).Scan(&result).Error
	return
}

// 群新增人数
func (g *GroupChatService) FindGroupsAddMembersByChatIds(chatIds []int64, params GropuNameQueryconditions) (result []group_chat.ChatMembers, err error) {
	var (
		sqlParams = []interface{}{params.EnterpriseID, 0, 2, 1, params.StartTime, params.EndTime, micro_scrm.EXTERNAL_CONTACTS_STRING, chatIds}
	)
	sql := `
	SELECT 
	group_chat_id,
	SUM(1) number 
  FROM
	(SELECT 
		DISTINCT group_chat_id,union_id 
	FROM work_group_chat_member_log 
	WHERE enterprise_id=? AND deleted_at =? AND member_type = ? AND operation_status = ?
	  AND created_at BETWEEN ? AND ?
      AND work_group_chat_member_log.external_userid LIKE ?
	  AND group_chat_id IN(?) GROUP BY group_chat_id,union_id) t 
   GROUP BY group_chat_id 
	`
	err = g.Raw(sql, sqlParams...).Scan(&result).Error
	return
}

// 退群人数
func (g *GroupChatService) FindGroupsCancleMembersByChatIds(chatIds []int64, params GropuNameQueryconditions) (result []group_chat.ChatMembers, err error) {
	var (
		sqlParams = []interface{}{params.EnterpriseID, 0, 2, 2, params.StartTime, params.EndTime, micro_scrm.EXTERNAL_CONTACTS_STRING}
		whereChat = ""
	)
	for k, v := range chatIds {
		if k == 0 {
			whereChat = "?"
		} else {
			whereChat += ",?"
		}
		sqlParams = append(sqlParams, v)
	}
	sql := `
	SELECT 
	group_chat_id,
	SUM(1) number 
  FROM
	(SELECT 
		DISTINCT group_chat_id,union_id 
	FROM work_group_chat_member_log 
	WHERE enterprise_id=? AND deleted_at =? AND member_type = ? AND operation_status = ?
	  AND created_at BETWEEN ? AND ?
      AND work_group_chat_member_log.external_userid LIKE ?
	  AND group_chat_id IN (` + whereChat + `) GROUP BY group_chat_id,union_id) t 
   GROUP BY group_chat_id 
	`
	err = g.Raw(sql, sqlParams...).Scan(&result).Error
	return
}

// 群成员加好友数
func (g *GroupChatService) FindGroupsMembersAddhaoyouByChatIds(chatIds []int64, params GropuNameQueryconditions, members []int64) (result []group_chat.ChatMembers, err error) {
	userexternalTable := " micro_scrm_enterprise_" + fmt.Sprintf("%d", params.EnterpriseID) + ".wx_user_corp_external "

	var (
		sqlParams = []interface{}{params.EnterpriseID, params.EnterpriseID, 2, 1, params.EndTime, params.EndTime, chatIds, members, chatIds}
	)
	sql := `
	
	 	SELECT group_chat_id,SUM(1) number FROM (
			SELECT a.group_chat_id,c.work_external_user_id 
		   FROM work_group_chat_member_log a  
	  		LEFT JOIN ` + userexternalTable + ` b  ON a.union_id = b.unionid AND b.unionid != '' AND b.deleted_at = 0 
	  		LEFT JOIN work_external_user_relation c ON b.id = c.work_external_user_id AND c.enterprise_id = ? AND c.relation_status = 0 AND c.deleted_at = 0 
			WHERE a.enterprise_id = ? AND a.member_type = ? AND a.operation_status = ? 
	  		AND a.deleted_at = 0  
			AND a.created_at <= ? 
	  		AND c.createtime <= ? 
	  		AND a.group_chat_id IN (?)
	  		AND c.member_id IN (?)
	  		AND c.id IS NOT NULL
     		GROUP BY a.group_chat_id,c.work_external_user_id
		 ) table1 GROUP BY group_chat_id
	`
	err = g.Raw(sql, sqlParams...).Scan(&result).Error
	return
}

// 微信会员绑定人数
func (g *GroupChatService) FindGroupsMembershuiyuanByChatIds(chatIds []int64, params GropuNameQueryconditions) (result []group_chat.ChatMembers, err error) {
	var (
		sqlParams = []interface{}{params.EnterpriseID, 0, 2, 1, params.EndTime}
		whereChat = ""
	)
	userexternalTable := " micro_scrm_enterprise_" + fmt.Sprintf("%d", params.EnterpriseID) + ".wx_user_corp_external "

	for k, v := range chatIds {
		if k == 0 {
			whereChat = "?"
		} else {
			whereChat += ",?"
		}
		sqlParams = append(sqlParams, v)
	}
	sql := `
	SELECT 
		group_chat_id,SUM(1) number 
  	FROM
	(SELECT 
	 DISTINCT group_chat_id,union_id 
	FROM
	  work_group_chat_member_log a
	  LEFT JOIN ` + userexternalTable + ` b ON a.union_id = b.unionid AND b.unionid != '' AND b.deleted_at = 0  AND b.union_code!=""
	WHERE a.enterprise_id = ? AND a.deleted_at = ? AND a.member_type = ? 
	  AND a.operation_status = ? AND a.created_at <= ? AND a.group_chat_id IN (` + whereChat + `) 
	  AND b.id IS NOT NULL  AND b.union_code!="" GROUP BY a.group_chat_id,a.union_id) t 
     GROUP BY group_chat_id 
	`
	err = g.Raw(sql, sqlParams...).Scan(&result).Error
	return
}

/**
*活动参与人数
 */

func (g *GroupChatService) FindGroupszhutaizuliByChatIds(chatIds []int64, params GropuNameQueryconditions) (result []group_chat.ChatMembers, err error) {
	var (
		sqlParams = []interface{}{params.EnterpriseID, params.EndTime, params.EnterpriseID, params.EndTime}
		whereChat = ""
	)

	for k, v := range chatIds {
		if k == 0 {
			whereChat = "?"
		} else {
			whereChat += ",?"
		}
		sqlParams = append(sqlParams, v)
	}
	sql := `
	SELECT 
		group_chat_id,SUM(1) number 
  	FROM
	(
		SELECT DISTINCT a.chat_id group_chat_id,b.user_id 
		FROM fission_chat_group a 
		LEFT JOIN fission_user b ON a.activity_id=b.activity_id AND b.deleted_at=0
		LEFT JOIN fission_mp_user c ON  b.user_id=c.id AND c.deleted_at=0 AND c.union_id!=""
		LEFT JOIN work_group_chat_member_log d ON d.enterprise_id=? AND d.created_at <= ?  AND d.group_chat_id=a.chat_id  AND c.union_id=d.union_id AND d.union_id!="" AND d.deleted_at=0 AND d.operation_status=1
		WHERE a.enterprise_id=?  AND b.created_at<=? AND a.chat_id IN (` + whereChat + `) AND c.id IS NOT NULL  AND c.union_id!="" AND d.union_id!=""
		GROUP BY a.chat_id,b.user_id) t GROUP BY group_chat_id 
	`
	err = g.Raw(sql, sqlParams...).Scan(&result).Error

	return
}

/**
*客太助力
 */

func (g *GroupChatService) FindGroupsMembersfissionassistantByChatIds(chatIds []int64, params GropuNameQueryconditions) (result []group_chat.ChatMembers, err error) {
	var (
		sqlParams = []interface{}{params.EnterpriseID, params.EndTime}
		whereChat = ""
	)

	for k, v := range chatIds {
		if k == 0 {
			whereChat = "?"
		} else {
			whereChat += ",?"
		}
		sqlParams = append(sqlParams, v)
	}
	sql := `
	SELECT 
		group_chat_id,SUM(1) number 
  	FROM
	(
		SELECT DISTINCT a.chat_id group_chat_id,b.assistant_user_id 
		FROM fission_chat_group a 
		LEFT JOIN fission_assistant_records b ON a.activity_id=b.activity_id AND b.deleted_at=0
		WHERE a.enterprise_id=?  AND b.created_at<=? AND a.chat_id IN (` + whereChat + `)
		GROUP BY a.chat_id,b.assistant_user_id) t GROUP BY group_chat_id 
	`
	err = g.Raw(sql, sqlParams...).Scan(&result).Error
	return
}

// ba数量

func (g *GroupChatService) FindGroupsBaByChatIds(chatIds []int64, params GropuNameQueryconditions, BAUserIds []string) (result []group_chat.ChatMembers, err error) {
	var (
		sqlParams = []interface{}{params.EnterpriseID, 0, params.EndTime}
		whereChat = ""
		userWHERE = ""
	)
	for k, v := range chatIds {
		if k == 0 {
			whereChat = "?"
		} else {
			whereChat += ",?"
		}
		sqlParams = append(sqlParams, v)
	}
	for k, v := range BAUserIds {
		if k == 0 {
			userWHERE = "?"
		} else {
			userWHERE += ",?"
		}
		sqlParams = append(sqlParams, v)
	}
	sql := `
	SELECT 
	group_chat_id,
	SUM(1) number 
  FROM
	(SELECT 
		DISTINCT group_chat_id,userid 
	FROM work_group_chat_member 
	WHERE enterprise_id=? AND deleted_at =? 
	  AND created_at <= ? 
	  AND group_chat_id IN(` + whereChat + `) 
	  AND userid IN (` + userWHERE + `) 
	  GROUP BY group_chat_id,userid) t 
   GROUP BY group_chat_id 
	`
	err = g.Raw(sql, sqlParams...).Scan(&result).Error
	return
}
