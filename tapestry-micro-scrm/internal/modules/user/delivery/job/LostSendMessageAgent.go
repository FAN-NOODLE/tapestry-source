package job

import (
	"encoding/json"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	messageSvc "micro-scrm/internal/app/service/message"
	"micro-scrm/internal/modules/enterprise"
	enterpriseMysql "micro-scrm/internal/modules/enterprise/respository/mysql"
	"micro-scrm/internal/orm/model/micro_scrm"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sdk/boss"
	"micro-scrm/internal/pkg/sdk/microwork"
	"micro-scrm/internal/pkg/xjob"
	"strconv"
	"time"
)

var corpId string

// LostCustomersSendMessage 流失客户定时推送到智汇导购小程序
func LostCustomersSendMessage(c ctx.Context, jobHandler *xjob.JobHandler) error {
	// 获取命令行参数配置的企业
	specsArgs := jobHandler.Runner().Specs()
	if len(specsArgs) <= 0 {
		// 获取所有企业
		enterpriseService := enterprise.New(enterpriseMysql.New())
		total, enterprises, err := enterpriseService.ListEnterprise(c, nil, 1, 0)
		if err != nil || total < 0 {
			return err
		}
		for _, e := range *enterprises {
			specsArgs = append(specsArgs, strconv.FormatInt(e.EnterpriseID, 10))
		}
	}

	for _, spec := range specsArgs {
		enterpriseID, _ := strconv.ParseInt(spec, 10, 64)
		// 通过企业id获取企业信息
		enterpriseDetail, err := boss.EnterpriseDetail(c, &boss.EnterpriseDetailRequest{EnterpriseID: enterpriseID})
		if err != nil {
			return err
		}
		if corpId == "" {
			corpId = enterpriseDetail.CorpID
		}

		// // 推送给企微流失配置
		// if err = handleCustomersLost(c, enterpriseID); err != nil {
		//	log.Error(c).Err(err)
		// }

		// 推送组织架构配置人员
		if err = handleDepartmentLost(c, enterpriseID); err != nil {
			log.Error(c).Err(err)
		}
	}

	return nil
}

// 处理流失列表配置推送
func handleCustomersLost(c ctx.Context, enterpriseID int64) error {
	// 查询删除客户消息配置表是否开启发送开关
	var customersDelMessages []micro_scrm.CustomersDelMessageModel
	if err := gormx.GetInstance().Where("enterprise_id", enterpriseID).Find(&customersDelMessages).Error; err != nil {
		log.Error(c).Msgf("查询删除客户消息配置表失败, enterpriseID:%d, err:%s", enterpriseID, err.Error())
		return err
	}

	// 从配置中获取配置的要发送的用户ID
	for _, message := range customersDelMessages {
		// 判断推送开关是否开启
		if message.Status != 1 {
			log.Error(c).Msgf("推送消息开关未开启, enterpriseID:%d, sendTYpe:%s", enterpriseID, message.SendType)
			continue
		}

		// 要发送的用户ID集合
		var memUserIds []string
		var customersDelMessagesStatTxt struct {
			Member []struct {
				Name   string `json:"name"`
				UserID string `json:"user_id"`
			} `json:"member"`
		}
		_ = json.Unmarshal([]byte(message.StatTxt), &customersDelMessagesStatTxt)
		for _, StatTxt := range customersDelMessagesStatTxt.Member {
			memUserIds = append(memUserIds, StatTxt.UserID)
		}

		// 查询当天删除好友关系记录
		var deleteFriendRelationRecords []micro_scrm.DeleteFriendRelationRecordModel
		startTimestamp, endTimestamp := helper.YesterdayPeriodOfTimestamp(1)
		// startTimestamp, endTimestamp := helper.PeriodOfTimestamp()
		err := gormx.GetInstance().Debug().
			Select("external_user_id,any_value(delete_type) as delete_type").
			Where("enterprise_id=? and created_at >= ? and created_at < ?", enterpriseID, startTimestamp, endTimestamp).
			Where("delete_type", message.SendType).
			Group("external_user_id").
			Find(&deleteFriendRelationRecords).Error
		if err != nil {
			log.Error(c).Msgf("查询当天删除好友关系记录失败, enterpriseID:%d, err:%s", enterpriseID, err.Error())
			continue
		}

		var numberOfActiveChurn, passiveChurn int
		// 主动流失数量,被动流失数量
		for _, record := range deleteFriendRelationRecords {
			if record.DeleteType == micro_scrm.DeleteFriendDeleteTypeExternalUser {
				passiveChurn++
				fmt.Println("用户删除员工+++++", passiveChurn)
			} else {
				numberOfActiveChurn++
				fmt.Println("员工删除用户+++++", numberOfActiveChurn)
			}
		}
		err = sendAgentMessage(c, enterpriseID, memUserIds, numberOfActiveChurn, passiveChurn, startTimestamp, "")
		log.Error(c).Err(err)
	}

	return nil
}

// 处理组织架构和智慧导购单人消息推送
func handleDepartmentLost(c ctx.Context, enterpriseID int64) error {
	var userIds []string
	// 获取已配的的人员集合
	gormx.GetInstance().Model(micro_scrm.CustomersDelUserMessageModel{}).
		Where("enterprise_id", enterpriseID).
		Group("user_id").
		Pluck("user_id", &userIds)

	for _, userId := range userIds {
		var customersDelUserMessage []micro_scrm.CustomersDelUserMessageModel
		// 查询用户对应的组织架构配置信息
		gormx.GetInstance().Select("dp_id,send_type").Where("enterprise_id", enterpriseID).Where("user_id", userId).Find(&customersDelUserMessage)

		var sendTypeDpIDs1, sendTypeDpIDs2 []int64
		for _, customer := range customersDelUserMessage {
			if customer.SendType == 0 || customer.SendType == 1 {
				sendTypeDpIDs1 = append(sendTypeDpIDs1, customer.DpID)
			}
			if customer.SendType == 0 || customer.SendType == 2 {
				sendTypeDpIDs2 = append(sendTypeDpIDs2, customer.DpID)
			}
		}

		startTimestamp, endTimestamp := helper.PeriodOfTimestamp()
		startTimestamp = startTimestamp - 86400
		endTimestamp = endTimestamp - 86400
		// 查询导购主动流失好友统计
		var deleteFriendRelationRecords2 []micro_scrm.DeleteFriendRelationRecordModel
		gormx.GetInstance().Select("external_user_id,any_value(delete_type) as delete_type").
			Where("enterprise_id=? and created_at >= ? and created_at < ?", enterpriseID, startTimestamp, endTimestamp).
			Where("delete_type", micro_scrm.DeleteFriendDeleteTypeUser).
			Where("dp_id in (?)", sendTypeDpIDs2).
			Group("external_user_id").
			Find(&deleteFriendRelationRecords2)

		// 查询导购被动流失好友统计
		var deleteFriendRelationRecords1 []micro_scrm.DeleteFriendRelationRecordModel
		gormx.GetInstance().Select("external_user_id,any_value(delete_type) as delete_type").
			Where("enterprise_id=? and created_at >= ? and created_at < ?", enterpriseID, startTimestamp, endTimestamp).
			Where("delete_type", micro_scrm.DeleteFriendDeleteTypeExternalUser).
			Where("dp_id in (?)", sendTypeDpIDs1).
			Group("external_user_id").
			Find(&deleteFriendRelationRecords1)

		// 主动流失数量,被动流失数量
		var numberOfActiveChurn, passiveChurn = len(deleteFriendRelationRecords2), len(deleteFriendRelationRecords1)

		var sendUserIds []string
		sendUserIds = append(sendUserIds, userId)
		err := sendAgentMessage(c, enterpriseID, sendUserIds, numberOfActiveChurn, passiveChurn, startTimestamp, "")
		log.Error(c).Err(err)
	}

	return nil
}

// sendAgentMessage 发送小程序应用消息
func sendAgentMessage(c ctx.Context, enterpriseID int64, memUserIds []string, numberOfActiveChurn, passiveChurn int, startTimestamp int64, page string) error {
	year := time.Now().Format("2006")
	month := time.Now().Format("01")
	day := time.Now().AddDate(0, 0, -1).Format("02")

	content := make([]map[string]string, 0)
	content = append(content, map[string]string{"key": "客户流失时间", "value": fmt.Sprintf("%s年%s月%s日", year, month, day)})
	content = append(content, map[string]string{"key": "被动流失数量", "value": fmt.Sprintf("%d人", passiveChurn)})
	content = append(content, map[string]string{"key": "主动流失数量", "value": fmt.Sprintf("%d人", numberOfActiveChurn)})
	content = append(content, map[string]string{"key": "说明", "value": "你管辖的BA导购有客户流失情况,点击查看详情"})

	params := microwork.SendMiniMessageRequest{
		EnterpriseID: enterpriseID,
		MemUserIds:   memUserIds,
		Title:        "客户流失提醒",
		Content:      content,
		Page:         fmt.Sprintf("pages/customer-lost/customer-lost?time=%d&corp_id=%s"+page, startTimestamp, corpId),
	}
	_, err := messageSvc.NewMessageService().SendMiniMessage(c, params)
	return err
}
