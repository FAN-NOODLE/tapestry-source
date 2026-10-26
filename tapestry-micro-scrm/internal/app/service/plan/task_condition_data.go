package planService

import (
	"fmt"
	"micro-scrm/internal/app"
	membersDao "micro-scrm/internal/app/dao/members"
	planDao "micro-scrm/internal/app/dao/micro_scrm"
	microWorkDao "micro-scrm/internal/app/dao/micro_work"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	memberReq "micro-scrm/internal/app/http/request/member"
	planReq "micro-scrm/internal/app/http/request/plan"
	planModels "micro-scrm/internal/app/models/micro_scrm"
	excitationSvc "micro-scrm/internal/app/service/excitation"
	"micro-scrm/internal/pkg/common"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sdk/macallback"
	"micro-scrm/internal/pkg/sdk/microwork"
	"micro-scrm/internal/pkg/sdk/work"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
)

type PlanTaskConditionData struct {
	planTaskConditionDataDao *planDao.PlanTaskConditionDataDao
}

func NewPlanTaskConditionData() *PlanTaskConditionData {
	return &PlanTaskConditionData{
		planTaskConditionDataDao: planDao.NewPlanTaskConditionDataDao(),
	}
}

func (p PlanTaskConditionData) InsertConditionData(c app.Ctx, enterpriseID, planId int64, conditionModels []planModels.TaskConditionDataModel) error {
	return p.planTaskConditionDataDao.InsertConditionData(c, enterpriseID, planId, conditionModels)
}

// Deprecated:CreateConditionData 方法
func (p PlanTaskConditionData) CreateConditionData(c app.Ctx, enterpriseID, planId int64, conditionModels planModels.TaskConditionDataModel) error {
	return p.planTaskConditionDataDao.CreateConditionData(c, enterpriseID, planId, conditionModels)
}

// ConfirmCall 导购确认点击打电话按钮
func (p PlanTaskConditionData) ConfirmCall(c ctx.Context, request planReq.ConfirmReq, isCoach bool) (bool, error) {
	one, err := planDao.NewPlanTaskAllotLogDao().GetOne(c, request.EnterpriseId, map[string]interface{}{
		"id": request.TaskAllotLogId,
	})
	if err != nil {
		return false, err
	}
	if !isCoach {
		whereMap := map[string]interface{}{
			"task_allot_log_id": request.TaskAllotLogId,
			"user_id":           request.UserId}
		if request.UnionId == "" && request.CustomerId != "" {
			whereMap["customer_id"] = request.CustomerId
		} else {
			whereMap["union_id"] = request.UnionId
		}
		return false, p.planTaskConditionDataDao.UpdateData(c, request.EnterpriseId, one.PlanId, whereMap, map[string]interface{}{"is_call": request.Status})
	}
	return p.finishTaskConditionData(c, request, map[string]interface{}{
		"is_call":     request.Status,
		"is_finished": planModels.FinishedYes,
	}, isCoach)
}

// ConfirmContact 导购取人点击发消息按钮
func (p PlanTaskConditionData) ConfirmContact(c ctx.Context, request planReq.ConfirmReq, isCoach bool) (bool, error) {
	one, err := planDao.NewPlanTaskAllotLogDao().GetOne(c, request.EnterpriseId, map[string]interface{}{
		"id": request.TaskAllotLogId,
	})
	if err != nil {
		return false, err
	}
	if !isCoach {
		whereMap := map[string]interface{}{
			"task_allot_log_id": request.TaskAllotLogId,
			"user_id":           request.UserId}
		if request.UnionId == "" && request.CustomerId != "" {
			whereMap["customer_id"] = request.CustomerId
		} else {
			whereMap["union_id"] = request.UnionId
		}
		return false, p.planTaskConditionDataDao.UpdateData(c, request.EnterpriseId, one.PlanId, whereMap, map[string]interface{}{"is_contact": request.Status})
	}
	return p.finishTaskConditionData(c, request, map[string]interface{}{
		"is_contact":  request.Status,
		"is_finished": planModels.FinishedYes,
	}, isCoach)
}

// SetRemark 联系说明
func (p PlanTaskConditionData) SetRemark(c ctx.Context, request planReq.SetRemark) error {
	one, err := planDao.NewPlanTaskAllotLogDao().GetOne(c, request.EnterpriseId, map[string]interface{}{
		"id": request.TaskAllotLogId,
	})
	if err != nil {
		return err
	}
	whereMap := map[string]interface{}{
		"task_allot_log_id": request.TaskAllotLogId,
		"user_id":           request.UserId}
	if request.UnionId == "" && request.CustomerId != "" {
		whereMap["customer_id"] = request.CustomerId
	} else {
		whereMap["union_id"] = request.UnionId
	}
	return p.planTaskConditionDataDao.UpdateData(c, request.EnterpriseId, one.PlanId, whereMap, map[string]interface{}{
		"remark": request.Remark,
	})
}

// ConfirmFinish 点击标记确认完成
func (p PlanTaskConditionData) ConfirmFinish(c ctx.Context, request planReq.ConfirmReq, isCoach bool) (bool, error) {
	return p.finishTaskConditionData(c, request, map[string]interface{}{
		"is_finished": planModels.FinishedYes,
	}, isCoach)
}

// 任务的标记联系如果最后一个标记为完成，则子任务自动完成
func (p PlanTaskConditionData) finishTaskConditionData(c ctx.Context, request planReq.ConfirmReq, updateMap map[string]interface{}, isCoach bool) (bool, error) {
	canAddPoints := false
	one, err := planDao.NewPlanTaskAllotLogDao().GetOne(c, request.EnterpriseId, map[string]interface{}{
		"id": request.TaskAllotLogId,
	})
	if err != nil {
		return false, err
	}
	enterpriseId := request.EnterpriseId
	id := request.TaskAllotLogId
	userId := request.UserId
	unionId := request.UnionId
	customerId := request.CustomerId
	whereMap := map[string]interface{}{
		"task_allot_log_id": id,
		"user_id":           userId,
	}
	if unionId == "" && customerId != "" {
		whereMap["customer_id"] = customerId
	} else {
		whereMap["union_id"] = unionId
	}
	updateMap["updated_at"] = time.Now().Unix()
	info, err := p.planTaskConditionDataDao.GetOne(c, enterpriseId, one.PlanId, whereMap)
	if err != nil {
		return false, err
	}
	if info.IsFinished == planModels.FinishedNo {
		canAddPoints = true
		updateMap["finished_at"] = time.Now().Unix()
	}
	err = p.planTaskConditionDataDao.UpdateData(c, enterpriseId, one.PlanId, whereMap, updateMap)
	taskInfo, err := planDao.NewPlanTaskDao().GetOne(c, enterpriseId, map[string]interface{}{
		"id": info.TaskId,
	})
	if err != nil {
		log.Info(c).Msgf("查询子任务信息报错 %+v", err)
		return canAddPoints, err
	}
	finishCount, err := planDao.NewPlanTaskConditionDataDao().GetGroupCount(c, request.EnterpriseId, one.PlanId, map[string]interface{}{
		"task_allot_log_id": request.TaskAllotLogId,
		"is_finished":       planModels.FinishedYes,
	})
	if err != nil {
		return canAddPoints, err
	}
	finishNum := decimal.NewFromInt(finishCount)
	customerNum := decimal.NewFromInt(int64(one.CustomerNum))
	divNum := finishNum.Div(customerNum).Round(2).Mul(decimal.NewFromInt(100)).IntPart()
	err = planDao.NewPlanTaskAllotLogDao().UpdateData(c, request.EnterpriseId, map[string]interface{}{
		"id": one.ID,
	}, map[string]interface{}{"finish_rate": divNum, "contacted_count": finishCount})
	if taskInfo.GroupSendType == planModels.PLANTASK1V1SEND { // 如果是企业沟通类型的任务,更新完成率和已触达人数
		if int(divNum) >= taskInfo.FinishRate && one.IsFinished == planModels.PLAN_FISISHING {
			// 任务完成增加积分
			excitationSvc.NewPointsProduceService().TaskFinished(ctx.Wrap(c), excitationReq.TaskFinishedReq{
				TaskAllotLogId: one.ID,
				EnterpriseId:   request.EnterpriseId,
				MemberId:       request.MemberId,
			})
			err := planDao.NewPlanTaskAllotLogDao().UpdateData(c, request.EnterpriseId, map[string]interface{}{
				"id": one.ID,
			}, map[string]interface{}{"is_finished": planModels.PLAN_FINISHED, "end_time": time.Now().Unix()})
			if err != nil {
				return canAddPoints, err
			}
		}
	} else {
		finishedCount, err := p.planTaskConditionDataDao.GetCount(c, enterpriseId, one.PlanId, map[string]interface{}{
			"user_id":           userId,
			"task_allot_log_id": id,
			"is_finished":       planModels.PLAN_FISISHING,
		})
		if err != nil {
			log.Info(c).Msgf("执行任务完成统计数报错 %+v", err)
			return canAddPoints, err
		}
		go func(finishedCount int64) {
			if finishedCount == 0 {
				if taskInfo.TaskCategory == planModels.PLANTASKCONTACT && taskInfo.GroupSendType < planModels.PLANTASK1V1SEND { // 沟通任务，全部标记为已联系则任务完成，朋友圈不是
					if one.IsFinished != 1 && finishedCount == 0 {
						excitationSvc.NewPointsProduceService().TaskFinished(ctx.Wrap(c), excitationReq.TaskFinishedReq{
							TaskAllotLogId: id,
							EnterpriseId:   request.EnterpriseId,
							MemberId:       request.MemberId,
						})
						err = planDao.NewPlanTaskAllotLogDao().UpdateData(c, enterpriseId, map[string]interface{}{
							"id":          id,
							"is_finished": planModels.PLAN_FISISHING,
						}, map[string]interface{}{
							"is_finished": planModels.PLAN_FINISHED,
							"end_time":    time.Now().Unix(),
						})
					}
					// 企业群发的沟通任务，如果任务已完成，群发消息未发送，则取消群发助手消息的发送
					if taskInfo.GroupSendType == planModels.PLANTASKQYSEND && isCoach {
						if one.IsGroupSend == planModels.QyGroupMsgSend {
							info, err := common.GetEnterpriseInfo(c, enterpriseId)
							if err != nil {
								return
							}
							if one.GroupMsgId != "" {
								qyWorkApi := new(work.QyWorkImpl).SetDomain(info.Domain)
								var request work.CancelGroupMsgSendRequest
								request.MsgId = one.GroupMsgId
								_, err = qyWorkApi.CancelGroupMsgSend(c, request)
								if err != nil {
									log.Info(c).Msgf("停止企业群发报错 err:%v\n", err)
									return
								}
							}
							err = planDao.NewPlanTaskAllotLogDao().UpdateData(c, enterpriseId, map[string]interface{}{
								"id": id,
							}, map[string]interface{}{
								"is_group_send": planModels.QyGroupMsgSendCancel,
							})
						}
					}
				}

			}
		}(finishedCount)
	}
	if info.IsFinished == planModels.FinishedNo && isCoach { // 顾客联系从未完成变成已完成，需要给MA推送完成的回执信息,目前只有coach企业对接了MA系统
		go func() {
			planInfo, err := planDao.NewPlanTaskDao().GetPlanDetail(c, planReq.PlanInfo{
				BaseReq: planReq.BaseReq{
					EnterpriseId: enterpriseId,
				},
				Id: info.PlanId,
			})
			if !helper.InArray(planInfo.Category, []int{planModels.PlanMaTypeYes, planModels.PlanMaTypeNewYes}) { // 非MA任务不推送
				return
			}
			if err != nil {
				log.Info(c).Msgf("标记联系MA回执查询主任务信息报错 %+v", err)
				return
			}
			var templateName string
			if taskInfo.PlanTemplateId > 0 {
				templateInfo, err := planDao.NewPlanTemplateDao().GetTemplateInfo(c, enterpriseId, taskInfo.PlanTemplateId)
				if err != nil {
					return
				}
				templateName = templateInfo.TemplateName
			}
			var events []microwork.MaEvents
			date := time.Now().UTC().Format("2006-01-02T15:04:05Z")
			memberInfo, err := membersDao.NewMembers().GetOne(c, map[string]interface{}{
				"enterprise_id": enterpriseId,
				"userid":        info.UserId,
				"deleted_at":    0,
			})
			if err != nil {
				log.Info(c).Msgf("查询导购信息身份失败 %+v", err)
				return
			}
			originType := taskInfo.OriginType
			event := microwork.MaEvents{
				IdentityType:      "wechat-unionid",
				IdentityValue:     info.UnionId,
				Date:              date,
				CResult:           "已完成",
				CSaName:           memberInfo.Name,
				CSaCode:           info.UserId,
				CTaskId:           strconv.FormatInt(info.TaskId, 10),
				CTaskName:         taskInfo.Name,
				CCrowdId:          taskInfo.CrowdID,
				CTaskRemark:       taskInfo.Remark,
				FlowId:            taskInfo.FlowId,
				Campaign:          taskInfo.Campaign,
				CPlanTemplateId:   strconv.FormatInt(taskInfo.PlanTemplateId, 10),
				CPlanTemplateName: templateName,
				CTaskClassifation: planInfo.Name,
			}
			events = append(events, event)
			params := microwork.SendMaFinishReq{
				EnterpriseID: enterpriseId,
				OriginType:   originType,
				Events:       events,
				CommonProps: microwork.MACommonProps{
					Event: "c_sa_completed_task_result",
				},
			}
			err = macallback.SendFinish(c, params)
			if err != nil {
				log.Info(c).Err(err).Interface("params", params).Msg("标记联系请求MA的任务发送结果回执")
			}
		}()

	}
	return canAddPoints, err
}

// SetForward 点击转发给顾客的按钮上报
func (p PlanTaskConditionData) SetForward(c ctx.Context, request planReq.SetForward, coach bool) (err error) {
	taskAllotInfo, err := planDao.NewPlanTaskAllotLogDao().GetInfoById(c, request.EnterpriseId, request.Id)
	if err != nil {
		return err
	}
	if taskAllotInfo.IsForward == 1 || taskAllotInfo.IsFinished == planModels.FinishedYes {
		return nil
	}
	taskInfo, err := planDao.NewPlanTaskDao().GetOne(c, request.EnterpriseId, map[string]interface{}{
		"id": taskAllotInfo.TaskId,
	})
	if err != nil {
		return err
	}
	if taskInfo.Status == planModels.PlanEnd {
		return nil
	}
	planId := taskAllotInfo.PlanId
	nowTime := time.Now().Unix()
	updateMap := map[string]interface{}{
		"is_forward": 1,
	}
	memberInfo, err := membersDao.NewMembers().GetOne(c, map[string]interface{}{
		"mobile":        request.Mobile,
		"enterprise_id": request.EnterpriseId,
		"deleted_at":    0,
	})
	conditionUpdateMap := make(map[string]interface{}, 0)
	if taskInfo.TaskCategory == planModels.PLANTASKMOMENT {
		conditionUpdateMap["is_touch"] = planModels.TouchStatusYes
		updateMap["is_finished"] = planModels.PLAN_FINISHED
		updateMap["end_time"] = nowTime
	} else {
		conditionUpdateMap = map[string]interface{}{
			"is_finished": planModels.PLAN_FINISHED,
			"finished_at": nowTime,
		}
	}
	if err != nil {
		return err
	}
	if coach {
		conditionsData, err := p.planTaskConditionDataDao.GetMany(c, request.EnterpriseId, planId, map[string]interface{}{
			"task_allot_log_id": request.Id,
			"is_finished":       planModels.PLAN_FISISHING,
			"deleted_at":        0,
		})
		if err != nil {
			return err
		}
		planInfo, err := planDao.NewPlanTaskDao().GetPlanDetail(c, planReq.PlanInfo{
			BaseReq: planReq.BaseReq{
				EnterpriseId: request.EnterpriseId,
			},
			Id: planId,
		})
		if err != nil {
			log.Info(c).Msgf("标记联系MA回执查询主任务信息报错 %+v", err)
		}
		isMaTask := false
		if helper.InArray(planInfo.Category, []int{planModels.PlanMaTypeYes, planModels.PlanMaTypeNewYes}) {
			isMaTask = true
		}
		var templateName string
		if taskInfo.PlanTemplateId > 0 {
			templateInfo, err := planDao.NewPlanTemplateDao().GetTemplateInfo(c, request.EnterpriseId, taskInfo.PlanTemplateId)
			if err != nil {
				return err
			}
			templateName = templateInfo.TemplateName
		}
		go func(enterpriseId int64, isMaTask bool) {
			var events []microwork.MaEvents
			var ids []int64
			var allNum, updateNum int
			for _, v := range conditionsData {
				info := v
				allNum++
				if v.IsUnion != 1 && v.GroupChatId == 0 {
					continue
				}
				if isMaTask {
					date := time.Now().UTC().Format("2006-01-02T15:04:05Z")
					event := microwork.MaEvents{
						IdentityType:      "wechat-unionid",
						IdentityValue:     info.UnionId,
						Date:              date,
						CResult:           "已完成",
						CSaName:           info.UserName,
						CSaCode:           info.UserId,
						CTaskId:           strconv.FormatInt(info.TaskId, 10),
						CTaskName:         taskInfo.Name,
						CCrowdId:          taskInfo.CrowdID,
						CTaskRemark:       taskInfo.Remark,
						FlowId:            taskInfo.FlowId,
						Campaign:          taskInfo.Campaign,
						CPlanTemplateId:   strconv.FormatInt(taskInfo.PlanTemplateId, 10),
						CPlanTemplateName: templateName,
						CTaskClassifation: planInfo.Name,
					}
					events = append(events, event)
				}
				ids = append(ids, info.ID)
			}
			updateNum = len(ids)
			originType := taskInfo.OriginType
			if isMaTask {
				params := microwork.SendMaFinishReq{
					EnterpriseID: enterpriseId,
					Events:       events,
					OriginType:   originType,
					CommonProps: microwork.MACommonProps{
						Event: "c_sa_completed_task_result",
					},
				}
				err = macallback.SendFinish(c, params)
				if err != nil {
					log.Info(c).Err(err).Interface("params", params).Msg("标记联系请求MA的任务发送结果回执")
				}
			}
			if allNum > updateNum { // 如果外部联系人不都是unionid情况，则只标记unionid的为已联系
				err = p.planTaskConditionDataDao.UpdateData(c, request.EnterpriseId, planId, map[string]interface{}{
					"id": ids,
				}, conditionUpdateMap)
				if taskInfo.TaskCategory == planModels.PLANTASKMOMENT {
					excitationSvc.NewPointsProduceService().TaskFinished(ctx.Wrap(c), excitationReq.TaskFinishedReq{
						TaskAllotLogId: request.Id,
						EnterpriseId:   request.EnterpriseId,
						MemberId:       memberInfo.ID,
					})
					err = planDao.NewPlanTaskAllotLogDao().UpdateData(c, request.EnterpriseId, map[string]interface{}{
						"id": request.Id,
					}, updateMap)
				}
			}
			if allNum == updateNum {
				excitationSvc.NewPointsProduceService().TaskFinished(ctx.Wrap(c), excitationReq.TaskFinishedReq{
					TaskAllotLogId: request.Id,
					EnterpriseId:   request.EnterpriseId,
					MemberId:       memberInfo.ID,
				})
				updateMap["end_time"] = nowTime
				updateMap["is_finished"] = planModels.PLAN_FINISHED
				err = planDao.NewPlanTaskAllotLogDao().UpdateData(c, request.EnterpriseId, map[string]interface{}{
					"id": request.Id,
				}, updateMap)
				err = p.planTaskConditionDataDao.UpdateData(c, request.EnterpriseId, planId, map[string]interface{}{
					"task_allot_log_id": request.Id,
					"is_union":          1,
				}, conditionUpdateMap)
			}
			if taskInfo.TaskCategory == planModels.PLANTASKMOMENT { // 目前转发给顾客的按钮只有朋友圈有该逻辑
				momentCount, _ := planDao.NewPlanTaskConditionDataDao().GetGroupCount(c, enterpriseId, planId, map[string]interface{}{
					"task_allot_log_id": taskAllotInfo.ID,
					"is_touch":          planModels.TouchStatusYes,
				})
				momentLookCount := decimal.NewFromInt(momentCount)
				customerNum := decimal.NewFromInt(int64(taskAllotInfo.CustomerNum))
				divNum := momentLookCount.Div(customerNum).Round(2).Mul(decimal.NewFromInt(100)).IntPart()
				err = planDao.NewPlanTaskAllotLogDao().UpdateData(c, enterpriseId, map[string]interface{}{
					"id": taskAllotInfo.ID,
				}, map[string]interface{}{"finish_rate": divNum, "moment_look_count": momentCount})
			}

		}(request.EnterpriseId, isMaTask)
	} else {
		excitationSvc.NewPointsProduceService().TaskFinished(ctx.Wrap(c), excitationReq.TaskFinishedReq{
			TaskAllotLogId: request.Id,
			EnterpriseId:   request.EnterpriseId,
			MemberId:       memberInfo.ID,
		})
		updateMap["end_time"] = nowTime
		updateMap["is_finished"] = planModels.PLAN_FINISHED
		err = planDao.NewPlanTaskAllotLogDao().UpdateData(c, request.EnterpriseId, map[string]interface{}{
			"id": request.Id,
		}, updateMap)
		if err != nil {
			return err
		}
		err = p.planTaskConditionDataDao.UpdateData(c, request.EnterpriseId, planId, map[string]interface{}{
			"task_allot_log_id": request.Id,
			"is_union":          1,
		}, conditionUpdateMap)
	}

	return err
}

func (p PlanTaskConditionData) TaskDetailExportBase(c ctx.Context, request planReq.TaskDetailExportReq) (planReq.TaskAllotDetailExportBase, error) {
	var err error
	var result planReq.TaskAllotDetailExportBase
	enterpriseId := request.EnterpriseId
	var userIds []string
	if request.BaName != "" && request.DepartmentId > 0 {
		memberLists, err := membersDao.NewMembers().GetMemberInfoAndDepartmentsByNameAndDepartmentId(c, enterpriseId, memberReq.MemberSearch{
			Name: request.BaName,
		}, []int64{request.DepartmentId})
		if err != nil {
			return result, err
		}
		for _, v := range memberLists {
			info := v
			if !helper.InArray[string](info.MembersUserid, userIds) {
				userIds = append(userIds, info.MembersUserid)
			}
		}
	} else {
		if request.BaName != "" {
			userIds, err = membersDao.NewMembers().GetUserIdsByName(c, enterpriseId, request.BaName)
			if err != nil {
				return result, err
			}
		}
		if request.DepartmentId > 0 {
			memberLists, err := membersDao.NewMembers().GetMemberByDepartmentId(c, enterpriseId, []int64{request.DepartmentId})
			if err != nil {
				return result, err
			}
			for _, v := range memberLists {
				info := v
				if !helper.InArray[string](info.MembersUserid, userIds) {
					userIds = append(userIds, info.MembersUserid)
				}
			}
		}
	}

	taskAllotLogInfos, err := planDao.NewPlanTaskAllotLogDao().GetTaskAllotLogs(c, enterpriseId, request, userIds)
	if err != nil || len(taskAllotLogInfos) == 0 {
		return result, err
	}
	var taskAllotUserIds []string
	taskForward := make(map[int64]int)
	taskFinish := make(map[int64]planReq.TaskAllotFinish)
	taskFinishRate := make(map[int64]int)
	for _, v := range taskAllotLogInfos {
		info := v
		if info.UserId != "" {
			if !helper.InArray[string](info.UserId, taskAllotUserIds) {
				taskAllotUserIds = append(taskAllotUserIds, info.UserId)
			}
		}
		if info.ID > 0 {
			if !helper.InArray[int64](info.ID, result.TaskAllotLogIds) {
				result.TaskAllotLogIds = append(result.TaskAllotLogIds, info.ID)
				taskFinish[info.ID] = planReq.TaskAllotFinish{
					IsFinished: info.IsFinished,
					EndTime:    info.EndTime,
				}
			}
		}
		taskForward[info.ID] = info.IsForward
		taskFinishRate[info.ID] = info.FinishRate
	}
	result.TaskForwards = taskForward
	result.TaskFinish = taskFinish
	result.TaskFinishRate = taskFinishRate
	planId := taskAllotLogInfos[0].PlanId
	taskId := taskAllotLogInfos[0].TaskId
	taskInfo, err := planDao.NewPlanTaskDao().GetOne(c, enterpriseId, map[string]interface{}{
		"id": taskId,
	})
	if err != nil {
		return result, err
	}
	result.TaskName = taskInfo.Name
	result.TaskCategory = taskInfo.TaskCategory
	result.GroupSendType = taskInfo.GroupSendType
	planInfo, err := planDao.NewPlanTaskDao().GetPlanDetail(c, planReq.PlanInfo{
		BaseReq: planReq.BaseReq{
			EnterpriseId: enterpriseId,
		},
		Id: planId,
	})
	if err != nil {
		return result, err
	}
	result.PlanId = planId
	result.PlanName = planInfo.Name
	result.Category = planInfo.Category
	memberLists, err := membersDao.NewMembers().GetUnscopedMemberInfoAndDepartmentsByUserId(c, enterpriseId, taskAllotUserIds)
	if err != nil {
		return result, err
	}
	var outLets, retail []string
	outLets, err = microWorkDao.NewStoreDepartment().GetDpIdsByType(c, enterpriseId, microWorkDao.STORE_TYPE_OUTLETS)
	if err != nil {
		return result, err
	}
	retail, err = microWorkDao.NewStoreDepartment().GetDpIdsByType(c, enterpriseId, microWorkDao.STORE_TYPE_RETAIL)
	if err != nil {
		return result, err
	}
	for _, v := range memberLists {
		info := v
		storeType := ""
		thirdId := strconv.FormatInt(info.DepartmentThirdId, 10)
		if helper.InArray(thirdId, outLets) {
			storeType = "奥莱"
		}
		if helper.InArray(thirdId, retail) {
			storeType = "零售"
		}
		item := planReq.TaskAllotMemberLists{
			SaName:    info.MembersName,
			SaCode:    info.MembersUserid,
			StoreName: info.DepartmentName,
			StoreType: storeType,
		}
		result.MemberLists = append(result.MemberLists, item)
	}
	return result, err
}

func (p PlanTaskConditionData) TaskDetailExportList(c ctx.Context, request planReq.TaskDetailExportListReq) ([]planReq.TaskDetailExportRes, error) {
	var result []planReq.TaskDetailExportRes
	lists, err := planDao.NewPlanTaskConditionDataDao().GetListsByPage(c, request.EnterpriseId, request.PlanId, request.Page, request.Size, request.TaskAllotLogId)
	if err != nil {
		return result, err
	}
	for _, v := range lists {
		info := v
		item := planReq.TaskDetailExportRes{
			Id:             info.TaskAllotLogId,
			UserId:         info.UserId,
			IsFinished:     int(info.IsFinished),
			FinishedAt:     info.FinishedAt,
			UnionID:        info.UnionId,
			CustomerID:     info.CustomerID,
			IsMatch:        info.IsMatch,
			IsContact:      info.IsContact,
			IsCall:         info.IsCall,
			Remark:         info.Remark,
			IsAdd:          0,
			CustomerStatus: 0,
			IsTouch:        info.IsTouch,
		}
		result = append(result, item)
	}
	return result, err
}

func (p PlanTaskConditionData) ConfirmFriend(c ctx.Context, request planReq.ConfirmReq) error {
	one, err := planDao.NewPlanTaskAllotLogDao().GetOne(c, request.EnterpriseId, map[string]interface{}{
		"id": request.TaskAllotLogId,
	})
	if err != nil {
		return err
	}
	return p.planTaskConditionDataDao.UpdateData(c, request.EnterpriseId, one.PlanId, map[string]interface{}{
		"id":      request.TaskAllotLogId,
		"user_id": request.UserId,
	}, map[string]interface{}{
		"is_add": request.Status,
	})
}

// FinishALL  sw和ks标记完成，所有顾客标记为已联系
func (p PlanTaskConditionData) FinishALL(c ctx.Context, enterpriseId, taskAllotLogId int64, userId string) error {
	one, err := planDao.NewPlanTaskAllotLogDao().GetOne(c, enterpriseId, map[string]interface{}{
		"id": taskAllotLogId,
	})
	if err != nil {
		return err
	}
	return p.planTaskConditionDataDao.UpdateData(c, enterpriseId, one.PlanId, map[string]interface{}{
		"task_allot_log_id": taskAllotLogId,
		"user_id":           userId,
	}, map[string]interface{}{
		"is_finished": planModels.PLAN_FINISHED,
		"finished_at": time.Now().Unix(),
	})
}

func (p PlanTaskConditionData) GetNoTouchLists(c ctx.Context, enterpriseId, planId, taskId int64) ([]planModels.TaskConditionDataModel, error) {
	return p.planTaskConditionDataDao.GetNoTouchLists(c, enterpriseId, planId, taskId)
}

// CreateShardTable 根据主任务id创建task_condition_data分表
func (p PlanTaskConditionData) CreateShardTable(c ctx.Context, enterpriseId, planId int64) error {
	return p.planTaskConditionDataDao.CreateShardTable(c, enterpriseId, planId)
}

// TaskDetailExportBaseNew 子任务数据明细的数据导出---下载中心
func (p PlanTaskConditionData) TaskDetailExportBaseNew(c ctx.Context, request planReq.TaskDetailExportListNewReq) ([]planReq.TaskDetailExportNewRes, int64, error) {
	var result []planReq.TaskDetailExportNewRes
	var total int64
	lists, total, err := planDao.NewPlanTaskConditionDataDao().GetListsByPageNew(c, request.EnterpriseId, request.PlanId, request.Page, request.Size, request.TaskAllotLogId)
	if err != nil {
		return result, total, err
	}
	memberLists := request.MemberLists
	memberListInfo := make(map[string]planReq.TaskAllotMemberLists)
	for _, v := range memberLists {
		info := v
		if _, ok := memberListInfo[info.SaCode]; !ok {
			memberListInfo[info.SaCode] = info
		}
	}
	planName := request.PlanName
	taskName := request.TaskName
	categoryName := "--"
	switch request.Category {
	case 1:
		categoryName = "手动创建"
	case 2:
		categoryName = "任务模板创建"
	}
	taskCategoryName := "--"
	switch request.TaskCategory {
	case 0:
		taskCategoryName = "沟通任务"
	case 1:
		taskCategoryName = "朋友圈任务"
	}

	for _, v := range lists {
		info := v
		touchName := "--"
		if request.GroupSendType == 1 {
			if request.TaskCategory == 0 {
				switch info.IsTouch {
				case 0:
					touchName = "未发送"
				case 1:
					touchName = "成功"
				case 2:
					touchName = "好友已删除 "
				case 3:
					touchName = "接收次数已达上限 "
				}
			} else {
				touchName = "未发送"
				if info.IsTouch == 1 {
					touchName = "可展示"
				}
				if request.TaskFinish[info.TaskAllotLogId].IsFinished == 1 && info.IsTouch == 0 {
					touchName = "不可展示"
				}
			}
		}
		finishedTimeStr := "--"
		if request.TaskFinish[info.TaskAllotLogId].EndTime > 0 {
			finishedTimeStr = helper.GetFormatDate(int64(request.TaskFinish[info.TaskAllotLogId].EndTime), 4)
		}
		isForward := "否"
		if request.GroupSendType < planModels.PLANTASK1V1SEND {
			isForward = threeBoolValue(request.TaskForwards[info.TaskAllotLogId])
		} else {
			isForward = threeBoolValue(info.Forward)
		}
		finishRate := "0%"
		if value, ok := request.TaskFinishRate[info.TaskAllotLogId]; ok {
			finishRate = fmt.Sprintf("%d%%", value)
		}
		touchType := "--"
		if info.Forward == planModels.AllForward {
			touchType = "群发触达"
		}
		if info.Forward == planModels.SingleForward && (info.IsCall > 0 || info.IsContact > 0) {
			touchType = "单独沟通"
		}
		customerFinishedTimeStr := "--"
		if info.FinishedAt > 0 {
			customerFinishedTimeStr = helper.GetFormatDateByParams(int64(info.FinishedAt))
		}
		item := planReq.TaskDetailExportNewRes{
			PlanName:                planName,
			TaskName:                taskName,
			CategoryName:            categoryName,
			TaskCategoryName:        taskCategoryName,
			SaName:                  memberListInfo[info.UserId].SaName,
			SaCode:                  memberListInfo[info.UserId].SaCode,
			StoreName:               memberListInfo[info.UserId].StoreName,
			StoreType:               memberListInfo[info.UserId].StoreType,
			IsTaskFinished:          threeBoolValue(request.TaskFinish[info.TaskAllotLogId].IsFinished),
			FinishedTimeStr:         finishedTimeStr,
			UnionID:                 info.UnionId,
			CustomerID:              info.CustomerID,
			IsForward:               isForward,
			IsMatch:                 threeBoolValue(info.IsMatch),
			IsContact:               threeBoolValue(info.IsContact),
			IsCall:                  threeBoolValue(info.IsCall),
			Remark:                  info.Remark,
			IsAdd:                   threeBoolValue(info.IsAdd),
			CustomerStatus:          threeContactValue(int(info.IsFinished)),
			IsTouch:                 info.IsTouch,
			TouchName:               touchName,
			FinishRate:              finishRate,
			TouchType:               touchType,
			CustomerFinishedTimeStr: customerFinishedTimeStr,
		}
		result = append(result, item)
	}
	return result, total, err
}

func threeBoolValue(value int) string {
	str := "否"
	if value == 1 {
		str = "是"
	}
	return str
}
func threeContactValue(value int) string {
	str := "待联系"
	if value == 1 {
		str = "已联系"
	}
	return str
}
