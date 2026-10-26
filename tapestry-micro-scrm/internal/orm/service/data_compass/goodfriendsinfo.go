package data_compass

import (
	"fmt"
	"micro-scrm/internal/orm/model/micro_scrm/work_external_user"
)

func (d *DataCompassService) GetNewGoodFriendsAllBymembers(enterpriseID int64, memberIds []int64, startTime, endTime int64) (result []work_external_user.Bastatistical, err error) {

	var (
		params = []interface{}{enterpriseID, startTime, endTime, memberIds}
	)

	sql := `
	SELECT member_id,SUM(1) number  FROM (SELECT 
		DISTINCT member_id,work_external_user_id 
		  FROM
		    work_external_user_relation 
		  WHERE enterprise_id=? AND  deleted_at = 0 AND relation_status=0 AND unionid!=''
		  AND createtime BETWEEN ? AND ?  AND member_id IN (?)) t GROUP BY member_id 
	`
	err = d.Raw(sql, params...).Scan(&result).Error
	return
}
func (d *DataCompassService) GetNewGoodFriendsAllHuiYuANBymembers(enterpriseID int64, memberIds []int64, startTime, endTime int64) (result []work_external_user.Bastatistical, err error) {
	userexternalTable := " micro_scrm_enterprise_" + fmt.Sprintf("%d", enterpriseID) + ".wx_user_corp_external "

	var (
		params = []interface{}{enterpriseID, startTime, endTime, memberIds}
	)

	sql := `
	SELECT member_id,SUM(1) number  FROM (SELECT 
		DISTINCT a.member_id,a.work_external_user_id 
		  FROM
			work_external_user_relation a 
			LEFT JOIN ` + userexternalTable + ` b
		  ON a.work_external_user_id = b.id  AND b.deleted_at=0 AND b.union_code!=""
		  WHERE a.enterprise_id=? AND  a.deleted_at = 0 AND a.relation_status=0 AND b.type=1  AND b.union_code!=""
		  AND a.createtime BETWEEN ? AND ?  AND a.member_id IN (?) ) t GROUP BY member_id 
	`
	err = d.Raw(sql, params...).Scan(&result).Error
	return
}

func (d *DataCompassService) GetNewGoodFriendsAllRuQunANBymembers(enterpriseID int64, memberIds, groupIds []int64, startTime, endTime int64) (result []work_external_user.Bastatistical, err error) {
	var (
		params = []interface{}{enterpriseID, startTime, endTime, memberIds, groupIds}
	)

	sql := `
	SELECT member_id,SUM(1) number  FROM (SELECT 
		DISTINCT weur.member_id,weur.work_external_user_id 
		  FROM
		    work_external_user_relation weur
		  LEFT JOIN work_group_chat_member wgcm ON wgcm.unionid=weur.unionid AND wgcm.status=1
		  WHERE weur.enterprise_id=? AND weur.deleted_at = 0  AND weur.relation_status=0 AND weur.unionid!='' AND wgcm.id IS NOT NULL
		  AND weur.createtime BETWEEN ? AND ?  AND weur.member_id IN (?) AND wgcm.group_chat_id IN (?)) t GROUP BY member_id 
	`
	err = d.Raw(sql, params...).Scan(&result).Error
	return
}

func (d *DataCompassService) GetDrainFriendsBymembers(enterpriseID int64, memberIds []int64, startTime, endTime int64) (result []work_external_user.Bastatistical, err error) {
	var (
		params = []interface{}{enterpriseID, startTime, endTime, startTime, endTime, memberIds}
	)

	sql := `
	SELECT member_id,SUM(1) number  FROM (SELECT 
		DISTINCT a.member_id,a.work_external_user_id 
		  FROM
		  	work_external_user_relation as a
		  WHERE a.enterprise_id=? AND  a.deleted_at = 0  AND a.unionid!='' 
		  AND ((a.member_deleted_time>=? AND a.member_deleted_time  <=?) or(a.external_deleted_time>=? AND  a.external_deleted_time<=? ))  AND member_id IN (?) ) t GROUP BY member_id 
	`
	err = d.Raw(sql, params...).Scan(&result).Error
	return
}

func (d *DataCompassService) GetLeijiFriendsAllBymembers(enterpriseID int64, memberIds []int64, endTime int64) (result []work_external_user.Bastatistical, err error) {
	var (
		params = []interface{}{enterpriseID, endTime, memberIds}
	)
	sql := `
	SELECT member_id,SUM(number) number  FROM (SELECT 
		member_id,work_external_user_id, 1 number 
		  FROM
			work_external_user_relation
		  WHERE enterprise_id=? AND  deleted_at = 0  AND unionid!='' AND relation_status=0
		  AND createtime <= ?  AND member_id IN (?) GROUP BY member_id,work_external_user_id) t GROUP BY member_id 
	`
	err = d.Raw(sql, params...).Scan(&result).Error
	return
}

func (d *DataCompassService) GetLeijiFriendsAllRUqunBymembers(enterpriseID int64, memberIds, groupIds []int64, endTime int64) (result []work_external_user.Bastatistical, err error) {

	var (
		params = []interface{}{enterpriseID, endTime, memberIds, groupIds}
	)
	sql := `
	SELECT member_id,SUM(number) number  FROM (SELECT 
		a.member_id,a.work_external_user_id,1 number
		  FROM
			work_external_user_relation a
		  LEFT JOIN work_group_chat_member c ON c.unionid=a.unionid AND c.status=1
		  WHERE a.enterprise_id=? AND  a.deleted_at = 0 AND a.relation_status=0 AND a.unionid!=''  AND c.id IS NOT NULL
		  AND a.createtime <= ?  AND a.member_id IN (?) AND c.group_chat_id IN (?) GROUP BY a.member_id,a.work_external_user_id) t GROUP BY member_id 
	`
	err = d.Raw(sql, params...).Scan(&result).Error
	return
}

// GetGoodFriendsAllBymembers 导购添加好友数量
func (d *DataCompassService) GetGoodFriendsAllBymembers(enterpriseID int64, memberIds []int64, startTime, endTime int64, page int, size int, ptype int) (result []work_external_user.Bastatistical, err error) {

	if ptype == 2 {

		var (
			params = []interface{}{enterpriseID, startTime, endTime, memberIds}
		)
		sql := `
	SELECT member_id,SUM(1) number  FROM (SELECT 
		DISTINCT a.member_id,a.work_external_user_id 
		  FROM
		    work_external_user_relation a 
		  WHERE a.enterprise_id=? AND  a.deleted_at = 0 AND a.relation_status=0 AND a.unionid!='' 
				AND a.createtime BETWEEN ? AND ?  AND a.member_id IN (?) ) t GROUP BY member_id ORDER BY number desc
	`
		err = d.Debug().Raw(sql, params...).Scan(&result).Error
	} else {
		var (
			params = []interface{}{enterpriseID, startTime, endTime, memberIds}
		)
		sql := `
	SELECT member_id,SUM(1) number  FROM (SELECT 
		DISTINCT a.member_id,a.work_external_user_id 
		  FROM
		    work_external_user_relation a
		  WHERE a.enterprise_id=? AND  a.deleted_at = 0 AND a.relation_status=0 AND a.unionid!='' 
			  AND a.createtime BETWEEN ? AND ?  AND a.member_id IN (?) ) t GROUP BY member_id ORDER BY number desc
				limit 50
	`
		err = d.Debug().Raw(sql, params...).Scan(&result).Error
	}

	return
}
