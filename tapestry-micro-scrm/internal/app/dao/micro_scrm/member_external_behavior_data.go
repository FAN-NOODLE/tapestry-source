package microScrmDao

import (
	"fmt"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/dao"
	memberReq "micro-scrm/internal/app/http/request/member"
	reserveReq "micro-scrm/internal/app/http/request/reserve"
	"micro-scrm/internal/app/models"
	excitationModels "micro-scrm/internal/app/models/excitation"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
	"strings"
)

type MemberExternalBehaviorData struct {
	dao   dao.Dao
	model *microScrmModel.MemberExternalBehaviorDataModel
}

func NewMemberExternalBehaviorData() *MemberExternalBehaviorData {
	return &MemberExternalBehaviorData{
		dao:   dao.Dao{DB: gormx.GetInstance()},
		model: new(microScrmModel.MemberExternalBehaviorDataModel),
	}
}

// WithTx 切换db连接,用于事物
func (d *MemberExternalBehaviorData) WithTx(tx *gorm.DB) *MemberExternalBehaviorData {
	d.dao.DB = tx
	return d
}

func (d *MemberExternalBehaviorData) GetOne(c app.Ctx, query interface{}, args ...interface{}) (microScrmModel.MemberExternalBehaviorDataModel, error) {
	var result microScrmModel.MemberExternalBehaviorDataModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

func (d *MemberExternalBehaviorData) GetMany(c app.Ctx, enterpriseId, query interface{}, args ...interface{}) ([]microScrmModel.MemberExternalBehaviorDataModel, error) {
	var result []microScrmModel.MemberExternalBehaviorDataModel
	err := d.dao.WithContext(c).Model(d.model).Where(query, args...).Find(&result).Error
	return result, err
}

// GetContactCountTopList 获取新增人员排行榜
func (d *MemberExternalBehaviorData) GetContactCountTopList(c app.Ctx, enterpriseId int64, inviteTop int, startTime, endTime int64) ([]models.TopList, error) {
	var result []models.TopList
	var lastTopTotal int64
	orm := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.model.TableName(), enterpriseId)).
		Select("SUM(new_contact_cnt) AS total").
		Where("created_at BETWEEN ? AND ?", startTime, endTime).
		Group("member_id").
		Order("total desc")
	// 获取top以及和top最后一位数量相同的其他导购
	newOrm := orm.Session(&gorm.Session{NewDB: true}).WithContext(c)
	err := orm.Offset(inviteTop-1).Limit(1).Pluck("total", &lastTopTotal).Error
	err = newOrm.Select("member_id,SUM(new_contact_cnt) AS total").
		Having("total > 0 AND total >= ?", lastTopTotal).Scan(&result).Error
	return result, err
}

// GetChatCountTopList 获取聊天排行榜
func (d *MemberExternalBehaviorData) GetChatCountTopList(c app.Ctx, enterpriseId int64, inviteTop int, startTime, endTime int64) ([]models.TopList, error) {
	var result []models.TopList
	var lastTopTotal int64
	orm := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.model.TableName(), enterpriseId)).
		Select("SUM(chat_cnt) AS total").
		Where("created_at BETWEEN ? AND ?", startTime, endTime).
		Group("member_id").
		Order("total desc")
	// 获取top以及和top最后一位数量相同的其他导购
	newOrm := orm.Session(&gorm.Session{NewDB: true}).WithContext(c)
	err := orm.Offset(inviteTop-1).Limit(1).Pluck("total", &lastTopTotal).Error
	err = newOrm.Select("member_id,SUM(chat_cnt) AS total").
		Having("total >= ?", lastTopTotal).Scan(&result).Error
	return result, err
}

// GetAverageReplyTimeTopList 获取平均回复时长榜单
func (d *MemberExternalBehaviorData) GetAverageReplyTimeTopList(c app.Ctx, enterpriseId int64, inviteTop int, startTime, endTime int64) ([]microScrmModel.TopList, error) {
	var result []microScrmModel.TopList
	var lastTopTotal int64
	orm := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.model.TableName(), enterpriseId)).
		Select("AVG(avg_reply_time) AS total").
		Where("reply_percentage >= 0").
		Where("created_at BETWEEN ? AND ?", startTime, endTime).
		Group("member_id").
		Order("total asc")
	// 获取top以及和top最后一位数量相同的其他导购
	newOrm := orm.Session(&gorm.Session{NewDB: true}).WithContext(c)
	err := orm.Offset(inviteTop-1).Limit(1).Pluck("total", &lastTopTotal).Error
	err = newOrm.Select("member_id,AVG(avg_reply_time) AS total").
		Having("total <= ?", lastTopTotal).Scan(&result).Error
	return result, err
}

// GetChatCount 获取所有导购聊天总数
func (d *MemberExternalBehaviorData) GetChatCount(c app.Ctx, enterpriseId int64, startTime int64, endTime int64) (list []models.BehaviorList, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.model.TableName(), enterpriseId)).
		Where("stat_time BETWEEN ? AND ?  ", startTime, endTime).
		Select("member_id,SUM(chat_cnt) AS total").
		Group("member_id").
		Having("total > ?", 0).
		Scan(&list).Error
	return list, err
}

// GetReplyPercentageCount 获取所有导购聊天回复占比
func (d *MemberExternalBehaviorData) GetReplyPercentageCount(c app.Ctx, enterpriseId int64, startTime int64, endTime int64) (list []models.BehaviorList, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.model.TableName(), enterpriseId)).
		Where("stat_time BETWEEN ? AND ?  ", startTime, endTime).
		Where("reply_percentage >= 0"). // 如果无聊天则为-1
		Select("member_id,round(avg(reply_percentage)) AS total").
		Group("member_id").
		Find(&list).Error
	return list, err
}

func (d *MemberExternalBehaviorData) ActiveCustomers(c app.Ctx, request memberReq.PerformanceListReq, membersID []int64) ([]memberReq.PerformanceListActiveRes, int64, error) {

	var (
		total  int64
		result []memberReq.PerformanceListActiveRes
		err    error
	)
	offset := (request.Page - 1) * request.Size
	startTime := request.FinancialTime[0]
	endTime := request.FinancialTime[1]

	// var result []memberReq.PerformanceListActiveRes
	query := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.model.TableName(), request.EnterpriseId)).
		Select("member_id,sum(chat_cnt) as chat_cnt,sum(message_cnt) as message_cnt,COALESCE(AVG(IF(reply_percentage>=0,reply_percentage,null)),0) as reply_percentage,sum(new_contact_cnt) as new_contact_cnt").
		Where("stat_time BETWEEN ? AND ?", startTime, endTime)

	switch request.ActiveType {
	case 1:
		// 活跃榜
		// 活跃榜排序按聊天总数正序排列，聊天总数相同时，按已回复聊天占比正序排列，前两个字段均相同时，按发送消息数正序排列，如全部相同，按员工名称首字母排序；
		query.Where("member_id in (?)", membersID).
			Order("chat_cnt desc").
			Order("reply_percentage desc").
			Order("message_cnt desc").
			Order("member_id asc")
	case 2:
		// 顾客榜
		query.Where("member_id in (?)", membersID).
			Order("new_contact_cnt desc").
			Order("member_id asc")
	case 3:
		// 明星榜-活跃榜（全部门店）
		query.Order("chat_cnt desc").
			Order("reply_percentage desc").
			Order("message_cnt desc").
			Order("member_id asc")
	case 4:
		// 明星榜-顾客榜（全部门店）
		query.Order("new_contact_cnt desc").
			Order("member_id asc")
	default:
		// 默认数据按照活跃榜
		return result, total, nil
	}

	err = query.Offset(offset).Limit(request.Size).Group("member_id").Find(&result).Error

	// 如果是明星榜单就不传递导购id
	if len(membersID) > 0 {
		sqlMembersID := strings.Trim(strings.Replace(fmt.Sprint(membersID), " ", ",", -1), "[]")
		sql := fmt.Sprintf("SELECT count(*) as total FROM (select member_id FROM micro_scrm_enterprise_%d.member_external_behavior_data WHERE ( stat_time BETWEEN %d AND %d and member_id in (%s)) GROUP By member_id) a", request.EnterpriseId, startTime, endTime, sqlMembersID)
		d.dao.WithContext(c).Raw(sql).Scan(&total)
	}

	return result, total, err
}

// GetAverage 获取平均值
func (d *MemberExternalBehaviorData) GetAverage(c app.Ctx, params reserveReq.AchievementParams) (int64, error) {
	var res int64
	memberCount := int64(len(params.MemberIds))
	orm := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.model.TableName(), params.EnterpriseId)).
		Where("member_id in ?", params.MemberIds).
		Where("stat_time BETWEEN ? AND ?  ", params.StartTime, params.EndTime)
	switch params.Type {
	case excitationModels.RuleCategoryId12:
		err := orm.Select("sum(chat_cnt) as sum_chat_cnt").
			Pluck("sum_chat_cnt", &res).Error
		if err != nil {
			return 0, err
		}
		return res / memberCount, nil
	case excitationModels.RuleCategoryId13:
		err := orm.Where("reply_percentage >= 0").
			Select("round(AVG(reply_percentage)) as avg_reply_percentage").
			Pluck("avg_reply_percentage", &res).Error
		if err != nil {
			return 0, err
		}
		return res, nil
	}

	return 0, nil
}

// GetChatCntCounts 导购的互动数
func (d *MemberExternalBehaviorData) GetChatCntCounts(c app.Ctx, enterpriseId, startTime, endTime int64, membersID []int64) (int64, error) {

	type Result struct {
		ChatCnt int64 `json:"chat_cnt"`
	}
	var resultInfo Result
	err := d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.model.TableName(), enterpriseId)).
		Select("sum(chat_cnt) as chat_cnt ").Where("member_id in (?)", membersID).
		Where("stat_time BETWEEN ? AND ?", startTime, endTime).Scan(&resultInfo).Error

	return resultInfo.ChatCnt, err
}

// GetChatCountAllTop 获取所有导购时间段内聊天互动排行榜
func (d *MemberExternalBehaviorData) GetChatCountAllTop(c app.Ctx, enterpriseId int64, startTime int64, endTime int64) (list []models.MemberTopList, err error) {
	err = d.dao.WithContext(c).Scopes(d.dao.BranchTable(d.model.TableName(), enterpriseId)).
		Where("stat_time BETWEEN ? AND ?  ", startTime, endTime).
		Select("member_id,COUNT(chat_cnt) AS total").
		Group("member_id").
		Order("total desc").
		Scan(&list).Error
	return list, err
}
