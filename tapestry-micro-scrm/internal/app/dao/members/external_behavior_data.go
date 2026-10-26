package membersDao

import (
	"fmt"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	memberReq "micro-scrm/internal/app/http/request/member"
	membersModel "micro-scrm/internal/app/models/members"
	"micro-scrm/internal/pkg/db/gormx"
	"strings"
)

type MemberExternalBehaviorData struct {
	dao       dao.Dao
	tableName string
}

func NewMemberExternalBehaviorData() *MemberExternalBehaviorData {
	return &MemberExternalBehaviorData{
		dao:       dao.Dao{DB: gormx.GetInstance()},
		tableName: membersModel.MemberExternalBehaviorDataModel{}.TableName(),
	}
}

func (d MemberExternalBehaviorData) ActiveCustomers(c app.Ctx, request memberReq.PerformanceListReq, membersID []int64) ([]memberReq.PerformanceListActiveRes, int64, error) {

	var (
		total  int64
		result []memberReq.PerformanceListActiveRes
		err    error
	)
	offset := (request.Page - 1) * request.Size
	startTime := request.FinancialTime[0]
	endTime := request.FinancialTime[1]

	//var result []memberReq.PerformanceListActiveRes
	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, request.EnterpriseId)).
		Select("member_id,sum(chat_cnt) as chat_cnt,sum(message_cnt) as message_cnt,COALESCE(AVG(IF(reply_percentage>=0,reply_percentage,null)),0) as reply_percentage,sum(new_contact_cnt) as new_contact_cnt").
		Where("stat_time BETWEEN ? AND ?", startTime, endTime)

	switch request.ActiveType {
	case 1:
		//活跃榜
		//活跃榜排序按聊天总数正序排列，聊天总数相同时，按已回复聊天占比正序排列，前两个字段均相同时，按发送消息数正序排列，如全部相同，按员工名称首字母排序；
		query.Where("member_id in (?)", membersID).
			Order("chat_cnt desc").
			Order("reply_percentage desc").
			Order("message_cnt desc").
			Order("member_id asc")
	case 2:
		//顾客榜
		query.Where("member_id in (?)", membersID).
			Order("new_contact_cnt desc").
			Order("member_id asc")
	case 3:
		//明星榜-活跃榜（全部门店）
		query.Order("chat_cnt desc").
			Order("reply_percentage desc").
			Order("message_cnt desc").
			Order("member_id asc")
	case 4:
		//明星榜-顾客榜（全部门店）
		query.Order("new_contact_cnt desc").
			Order("member_id asc")
	default:
		//默认数据按照活跃榜
		return result, total, nil
	}

	err = query.Offset(offset).Limit(request.Size).Group("member_id").Find(&result).Error

	//如果是明星榜单就不传递导购id
	if len(membersID) > 0 {
		sqlMembersID := strings.Trim(strings.Replace(fmt.Sprint(membersID), " ", ",", -1), "[]")
		sql := fmt.Sprintf("SELECT count(*) as total FROM (select member_id FROM micro_scrm_enterprise_%d.member_external_behavior_data WHERE ( stat_time BETWEEN %d AND %d and member_id in (%s)) GROUP By member_id) a", request.EnterpriseId, startTime, endTime, sqlMembersID)
		d.dao.WithContext(c).Raw(sql).Scan(&total)
	}

	return result, total, err
}

// GetChatCntCounts 导购的互动数
func (d MemberExternalBehaviorData) GetChatCntCounts(c app.Ctx, enterpriseId, startTime, endTime int64, membersID []int64) (int64, error) {

	type Result struct {
		ChatCnt int64 `json:"chat_cnt"`
	}
	var resultInfo Result
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, enterpriseId)).
		Select("sum(chat_cnt) as chat_cnt ").Where("member_id in (?)", membersID).
		Where("stat_time BETWEEN ? AND ?", startTime, endTime).Scan(&resultInfo).Error

	return resultInfo.ChatCnt, err
}

//SwActive sw 活跃榜单
func (d MemberExternalBehaviorData) SwActive(c app.Ctx, request memberReq.SwActiveListReqReq, membersID []int64) (map[int64]memberReq.PerformanceListActiveRes, error) {

	var (
		resMap = make(map[int64]memberReq.PerformanceListActiveRes)
		result []memberReq.PerformanceListActiveRes
		err    error
	)
	startTime := request.FinancialTime[0]
	endTime := request.FinancialTime[1]

	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.tableName, request.EnterpriseId)).
		Select("member_id,sum(chat_cnt) as chat_cnt,sum(message_cnt) as message_cnt,COALESCE(AVG(IF(reply_percentage>=0,reply_percentage,null)),0) as reply_percentage,sum(new_contact_cnt) as new_contact_cnt").
		Where("stat_time BETWEEN ? AND ?", startTime, endTime)
	if len(membersID) > 0 {
		query.Where("member_id in (?)", membersID)
	}
	err = query.Group("member_id").Find(&result).Error
	if err != nil {
		return nil, err
	}

	for _, res := range result {
		resMap[res.MemberID] = res
	}
	return resMap, nil
}
