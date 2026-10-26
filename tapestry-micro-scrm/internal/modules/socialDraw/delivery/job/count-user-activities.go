package job

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/pkg/xjob"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/pkg/db"
	"micro-scrm/internal/pkg/redis"
)

// UserActivity 参与活动统计数据
// 分类存储数据 {"fission": [{"activity_id": 0, "join_at": 0}], "market": [{"activity_id": 0, "join_at": 0}]}
type UserActivity struct {
	EnterpriseID      int64         `db:"enterprise_id"`
	ActivityNum       int           `db:"activity_num"`
	UnionID           string        `db:"unionid" desc:"fission_mp_user的unionid"`
	ChatID            int64         `json:"-" db:"chat_id" desc:"社群主键"`
	FissionActivities []*ActiveTime `json:"-" desc:"参与裂变活动列表"`
	MarketActivities  []*ActiveTime `json:"-" desc:"参与营销活动列表"`
	TmpActivities     []byte        `json:"-" db:"act_ids" desc:"参与的所有活动"`
}

type ActiveTime struct {
	UnionID    string `json:"-" db:"unionid"`
	ActivityID int64  `json:"activity_id" db:"activity_id"`
	ChatID     int64  `json:"-" db:"chat_id"`
	JoinAt     int64  `json:"join_at" db:"join_at" desc:"参与时间"`
}

// getLastCountTime 获取上一次查询时间
func getOrSetLastCountTime(c ctx.Context, enterpriseId int64, sets int64) int64 {
	cacheKey := fmt.Sprintf("work_goup_chat_member_activity_last_count_time:%+v", enterpriseId)
	client := redis.Client()
	timestamp, _ := client.Get(c, cacheKey).Int64()
	client.Set(c, cacheKey, sets, 0)
	return timestamp
}

// CountUserActivities 统计用户参与活动数
// 只统计裂变活动-主库、营销活动（问卷、签到、报名）-分库、外部联系人-分库
// 初始化时，全量统计；以后只统计到昨天的数据
func CountUserActivities(c ctx.Context, jobHandler *xjob.JobHandler) error {
	log.Debug(c).Msgf("==================== start 同步社群用户参与活动数 =======================")
	var (
		mdb           = db.DB()
		enterpriseIds []int64 // 所有企业
	)
	if err := mdb.Select(c, &enterpriseIds, `select enterprise_id from enterprise_corp`); err != nil {
		return err
	}
	for _, eid := range enterpriseIds {
		var (
			totalMap       = make(map[string]*UserActivity)
			countEndAt     = time.Now().Unix()                         // 结束时间
			countStartAt   = getOrSetLastCountTime(c, eid, countEndAt) // 开始时间
			memberActivity []*UserActivity                             // 参与活动群成员
		)

		// 裂变活动
		fissionActivityUsers := countUserFissionActivities(c, eid, countStartAt, countEndAt)
		for _, v := range fissionActivityUsers {
			mk := fmt.Sprintf("%s%d", v.UnionID, v.ChatID)
			totalMap[mk] = v
		}
		// 营销活动
		marketActivityUsers := countUserMarketActivities(c, eid, countStartAt, countEndAt)
		for _, v := range marketActivityUsers {
			mk := fmt.Sprintf("%s%d", v.UnionID, v.ChatID)
			if _, ok := totalMap[mk]; !ok {
				totalMap[mk] = v
			} else {
				totalMap[mk].ActivityNum += v.ActivityNum
				totalMap[mk].MarketActivities = v.MarketActivities
			}
		}
		if countStartAt > 0 { // 全覆盖更新
			// 查询历史参与活动数
			_ = mdb.Select(c, &memberActivity, `select unionid,group_chat_id chat_id,act_num activity_num,act_ids from work_group_chat_member where enterprise_id=?`, eid)
			for _, v := range memberActivity {
				mk := fmt.Sprintf("%s%d", v.UnionID, v.ChatID)
				if _, ok := totalMap[mk]; ok {
					var tmpMemberActMap = make(map[string][]*ActiveTime)
					_ = json.Unmarshal(v.TmpActivities, &tmpMemberActMap)
					log.Debug(c).Msgf("------------------------------------------------- %+v", tmpMemberActMap)
					totalMap[mk].FissionActivities = append(totalMap[mk].FissionActivities, tmpMemberActMap["fission"]...)
					totalMap[mk].MarketActivities = append(totalMap[mk].MarketActivities, tmpMemberActMap["market"]...)

					totalMap[mk].ActivityNum += v.ActivityNum
				}
			}
		}

		for _, activity := range totalMap {
			userActivity := make(map[string][]*ActiveTime)
			if len(activity.FissionActivities) > 0 {
				userActivity["fission"] = activity.FissionActivities
			}
			if len(activity.MarketActivities) > 0 {
				userActivity["market"] = activity.MarketActivities
			}

			userActs, _ := json.Marshal(userActivity)
			query := `update work_group_chat_member set act_num=?,act_ids=? where enterprise_id=? and unionid=? and group_chat_id=? and status=1`
			if _, err := mdb.DB.Exec(c, query, activity.ActivityNum, string(userActs), eid, activity.UnionID, activity.ChatID); err != nil {
				log.Error(c).Msgf("exec sql【%+v】error %+v", query, err)
				continue
			}
		}
	}
	return nil
}

// countUserFissionActivities 分企业统计裂变活动用户参与活动数
// unionid -> 活动 -> 更新到chat_member中
func countUserFissionActivities(c ctx.Context, enterpriseId int64, startAt, endAt int64) []*UserActivity {
	var (
		fissionUsers []*UserActivity
		users        []string
		mdb          = db.DB()
	)
	query := `select count(1) activity_num,fmu.union_id unionid,fcg.chat_id from fission_user fu inner join fission_mp_user fmu on fu.user_id = fmu.id
		inner join fission_chat_group fcg on fu.activity_id = fcg.activity_id where fu.deleted_at = 0 and fmu.union_id != '' and 
		fu.enterprise_id=? and fu.created_at>=? and fu.created_at<? group by fcg.chat_id,unionid`
	if err := mdb.Select(c, &fissionUsers, query, enterpriseId, startAt, endAt); err != nil {
		log.Error(c).Msgf("exec sql【%+v】error %+v", query, err)
		return nil
	}
	for _, v := range fissionUsers {
		v.EnterpriseID = enterpriseId
		users = append(users, v.UnionID)
	}

	var (
		userActivityTimes []*ActiveTime
		timeMap           = make(map[string][]*ActiveTime)
	)
	query, args, err := sqlx.In(`select fu.activity_id, fu.created_at join_at, fmu.union_id unionid,fcg.chat_id from fission_user fu
         inner join fission_mp_user fmu on fu.user_id = fmu.id inner join fission_chat_group fcg on fu.activity_id = fcg.activity_id
		where fmu.union_id in (?) and fu.enterprise_id = ? and fu.created_at>=? and fu.created_at<?`,
		users, enterpriseId, startAt, endAt)
	if err != nil {
		return nil
	}
	if err := mdb.Select(c, &userActivityTimes, query, args...); err != nil {
		return nil
	}
	for _, v := range userActivityTimes {
		mk := fmt.Sprintf("%s%d", v.UnionID, v.ChatID)
		timeMap[mk] = append(timeMap[mk], v)
	}
	for _, v := range fissionUsers {
		mk := fmt.Sprintf("%s%d", v.UnionID, v.ChatID)
		if _, ok := timeMap[mk]; ok {
			v.FissionActivities = timeMap[mk]
		}
	}
	return fissionUsers
}

// countUserMarketActivities 统计指定企业营销活动参与者的参与数
func countUserMarketActivities(c ctx.Context, enterpriseId int64, startAt, endAt int64) []*UserActivity {
	var (
		marketUsers []*UserActivity
		users       []string
		mdb         = db.DB()
	)
	query := fmt.Sprintf(`select count(1) activity_num,union_id unionid,wgc.id chat_id from micro_scrm_enterprise_%+v.activity_user au 
		inner join work_group_chat wgc on au.chat_id = convert(wgc.chat_id using 'utf8mb4')
		where union_id != '' and au.created_at >=? and au.created_at<? group by wgc.chat_id, unionid`, enterpriseId)
	if err := mdb.Select(c, &marketUsers, query, startAt, endAt); err != nil {
		log.Error(c).Msgf("exec sql 【%+v】 error %+v", query, err)
		return nil
	}
	for _, v := range marketUsers {
		v.EnterpriseID = enterpriseId
		users = append(users, v.UnionID)
	}

	var (
		userActivityTimes []*ActiveTime
		timeMap           = make(map[string][]*ActiveTime)
	)
	query, args, err := sqlx.In(fmt.Sprintf(`select union_id unionid,activity_id,au.created_at join_at,wgc.id chat_id 
		from micro_scrm_enterprise_%+v.activity_user au
		inner join work_group_chat wgc on au.chat_id = convert(wgc.chat_id using 'utf8mb4')
		where union_id in (?) and au.created_at >=? and au.created_at<?`, enterpriseId), users, startAt, endAt)
	if err != nil {
		return nil
	}
	if err := mdb.Select(c, &userActivityTimes, query, args...); err != nil {
		return nil
	}
	for _, v := range userActivityTimes {
		mk := fmt.Sprintf("%s%d", v.UnionID, v.ChatID)
		timeMap[mk] = append(timeMap[mk], v)
	}
	for _, v := range marketUsers {
		mk := fmt.Sprintf("%s%d", v.UnionID, v.ChatID)
		if _, ok := timeMap[mk]; ok {
			v.MarketActivities = timeMap[mk]
		}
	}
	return marketUsers
}
