package excitationSvc

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	excitationDao "micro-scrm/internal/app/dao/excitation"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	excitationResp "micro-scrm/internal/app/http/response/excitation"
	"micro-scrm/internal/app/models"
	excitationModels "micro-scrm/internal/app/models/excitation"
	"micro-scrm/internal/app/service"
)

type LevelRuleService struct {
	service.Base
	levelRuleDao *excitationDao.LevelRule
}

func NewLevelRuleService() *LevelRuleService {
	return &LevelRuleService{
		levelRuleDao: excitationDao.NewLevelRule(),
	}
}

// CreateOrUpdate 创建或修改等级规则
func (s *LevelRuleService) CreateOrUpdate(c ctx.Context, params excitationReq.LevelRuleSaveArrReq) error {

	err := s.Base.DB().Transaction(func(tx *gorm.DB) error {

		var levelRules []excitationModels.ExcitationLevelRuleModel

		var levelRuleDeletes []int64
		for _, levelRuleDelete := range params.LevelRules {
			if levelRuleDelete.IsDelete == 1 {
				levelRuleDeletes = append(levelRuleDeletes, levelRuleDelete.ID)
			}
		}
		// 删除
		if len(levelRuleDeletes) > 0 {
			var deleteRuleWhere = map[string]interface{}{
				"enterprise_id": params.EnterpriseId,
				"id":            levelRuleDeletes,
			}
			err := excitationDao.NewLevelRule().WithTx(tx).DeleteData(c, deleteRuleWhere)
			if err != nil {
				return err
			}
		}

		for _, levelRule := range params.LevelRules {

			var ruleInviteID, ruleTaskID, ruleInteractID, ruleTrainID, ruleOnlineID int64
			var ruleConditionInvite, ruleConditionTask, ruleConditionInteract, ruleConditionTrain, ruleConditionOnline int

			if len(levelRule.LevelRuleAttachArr) > 0 {
				for _, levelRuleArr := range levelRule.LevelRuleAttachArr {
					if levelRuleArr.IncentiveConfigType == excitationModels.IncentiveConfigTypeInvite {
						ruleInviteID = levelRuleArr.IncentiveRelationRuleId
						ruleConditionInvite = levelRuleArr.RuleCondition
					}
					if levelRuleArr.IncentiveConfigType == excitationModels.IncentiveConfigTypeTask {
						ruleTaskID = levelRuleArr.IncentiveRelationRuleId
						ruleConditionTask = levelRuleArr.RuleCondition
					}
					if levelRuleArr.IncentiveConfigType == excitationModels.IncentiveConfigTypeInteract {
						ruleInteractID = levelRuleArr.IncentiveRelationRuleId
						ruleConditionInteract = levelRuleArr.RuleCondition
					}
					if levelRuleArr.IncentiveConfigType == excitationModels.IncentiveConfigTypeTrain {
						ruleTrainID = levelRuleArr.IncentiveRelationRuleId
						ruleConditionTrain = levelRuleArr.RuleCondition
					}
					if levelRuleArr.IncentiveConfigType == excitationModels.IncentiveConfigTypeOnline {
						ruleOnlineID = levelRuleArr.IncentiveRelationRuleId
						ruleConditionOnline = levelRuleArr.RuleCondition
					}
				}
			}

			if levelRule.ID > 0 {

				// 修改
				var levelRuleWhere = map[string]interface{}{
					"enterprise_id": params.EnterpriseId,
					"id":            levelRule.ID,
				}
				var levelRuleUpdate = map[string]interface{}{
					"level_type":              levelRule.LevelType,
					"level_name":              levelRule.LevelName,
					"Level_icon":              levelRule.LevelIcon,
					"max_value":               levelRule.MaxValue,
					"min_value":               levelRule.MinValue,
					"rule_invite_id":          ruleInviteID,
					"rule_task_id":            ruleTaskID,
					"rule_interact_id":        ruleInteractID,
					"rule_train_id":           ruleTrainID,
					"rule_online_id":          ruleOnlineID,
					"rule_condition_invite":   ruleConditionInvite,
					"rule_condition_task":     ruleConditionTask,
					"rule_condition_interact": ruleConditionInteract,
					"rule_condition_train":    ruleConditionTrain,
					"rule_condition_online":   ruleConditionOnline,
				}
				err := excitationDao.NewLevelRule().WithTx(tx).UpdateData(c, levelRuleWhere, levelRuleUpdate)
				if err != nil {
					return err
				}

			} else {

				var levelRuleId = app.UID()
				levelRules = append(levelRules, excitationModels.ExcitationLevelRuleModel{
					MicroBaseModel:        models.CreateMicroBase(params.EnterpriseId, levelRuleId),
					LevelType:             levelRule.LevelType,
					LevelName:             levelRule.LevelName,
					LevelIcon:             levelRule.LevelIcon,
					MaxValue:              levelRule.MaxValue,
					MinValue:              levelRule.MinValue,
					RuleInviteID:          ruleInviteID,
					RuleTaskID:            ruleTaskID,
					RuleInteractID:        ruleInteractID,
					RuleTrainID:           ruleTrainID,
					RuleOnlineID:          ruleOnlineID,
					RuleConditionInvite:   ruleConditionInvite,
					RuleConditionTask:     ruleConditionTask,
					RuleConditionInteract: ruleConditionInteract,
					RuleConditionTrain:    ruleConditionTrain,
					RuleConditionOnline:   ruleConditionOnline,
				})
			}
		}
		if len(levelRules) > 0 {
			err := excitationDao.NewLevelRule().WithTx(tx).CreateInBatchesData(c, levelRules)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

// getRelationIDs 等级附加条件添加
func (s *LevelRuleService) getRelationIDs(c ctx.Context, tx *gorm.DB, enterpriseId int64, levelRuleAttachArr []excitationReq.LevelRuleAttachReq) (int64, int64, int64, int64, int64, error) {
	var ruleInviteID, ruleTaskID, ruleInteractID, ruleTrainID, ruleOnlineID int64
	if len(levelRuleAttachArr) > 0 {

		// 附加添加添加数据
		var deleteIds []int64
		var levelRuleAttachArray []excitationModels.ExcitationIncentiveRelationRuleModel
		for _, val := range levelRuleAttachArr {
			if val.IncentiveRelationRuleId > 0 {
				deleteIds = append(deleteIds, val.IncentiveRelationRuleId)
			}
			var relationRuleId = app.UID()
			levelRuleAttachArray = append(levelRuleAttachArray, excitationModels.ExcitationIncentiveRelationRuleModel{
				MicroBaseModel:      models.CreateMicroBase(enterpriseId, relationRuleId),
				IncentiveConfigType: val.IncentiveConfigType,
				RuleCondition:       val.RuleCondition,
				RuleGive:            val.RuleGive,
			})

			switch val.IncentiveConfigType {
			case excitationModels.IncentiveConfigTypeInvite:
				ruleInviteID = relationRuleId
				break
			case excitationModels.IncentiveConfigTypeTask:
				ruleTaskID = relationRuleId
				break
			case excitationModels.IncentiveConfigTypeInteract:
				ruleInteractID = relationRuleId
				break
			case excitationModels.IncentiveConfigTypeTrain:
				ruleTrainID = relationRuleId
				break
			case excitationModels.IncentiveConfigTypeOnline:
				ruleOnlineID = relationRuleId
				break
			}
		}
		if len(deleteIds) > 0 {
			// 删除之前规则
			var relationRuleWhere = map[string]interface{}{
				"enterprise_id": enterpriseId,
				"id":            deleteIds,
			}
			err := excitationDao.NewIncentiveRelationRule().WithTx(tx).DeleteData(c, relationRuleWhere)
			if err != nil {
				return ruleInviteID, ruleTaskID, ruleInteractID, ruleTrainID, ruleOnlineID, err
			}
		}
		if len(levelRuleAttachArray) > 0 {
			err := excitationDao.NewIncentiveRelationRule().WithTx(tx).CreateInBatchesData(c, levelRuleAttachArray)
			if err != nil {
				return ruleInviteID, ruleTaskID, ruleInteractID, ruleTrainID, ruleOnlineID, err
			}
		}
	}

	return ruleInviteID, ruleTaskID, ruleInteractID, ruleTrainID, ruleOnlineID, nil
}

// GetDetailList  等级规则
func (s *LevelRuleService) GetDetailList(c ctx.Context, params excitationReq.LevelRuleDetailReq) ([]excitationResp.LevelRuleDetailListResp, error) {
	var result []excitationResp.LevelRuleDetailListResp
	levelRuleRes, err := s.levelRuleDao.GetDetailList(c, params)
	if err != nil {
		return result, err
	}

	var levelRuleIDs []int64
	var levelRuleAttachArray []excitationResp.LevelRuleAttachResp
	for _, item := range levelRuleRes {
		result = append(result, excitationResp.LevelRuleDetailListResp{
			LevelRuleDetailResp: excitationResp.LevelRuleDetailResp{
				EnterpriseId:          item.EnterpriseId,
				ID:                    item.ID,
				LevelName:             item.LevelName,
				LevelType:             item.LevelType,
				LevelIcon:             item.LevelIcon,
				MinValue:              item.MinValue,
				MaxValue:              item.MaxValue,
				RuleConditionInvite:   item.RuleConditionInvite,
				RuleConditionTask:     item.RuleConditionTask,
				RuleConditionInteract: item.RuleConditionInteract,
				RuleConditionTrain:    item.RuleConditionTrain,
				RuleConditionOnline:   item.RuleConditionOnline,
			},
			LevelRuleDetailListIDsResp: excitationResp.LevelRuleDetailListIDsResp{
				RuleInviteID:   item.RuleInviteID,
				RuleTaskID:     item.RuleTaskID,
				RuleInteractID: item.RuleInteractID,
				RuleTrainID:    item.RuleTrainID,
				RuleOnlineID:   item.RuleOnlineID,
			},
			LevelRuleAttachArr: levelRuleAttachArray,
		})

		levelRuleIDs = append(levelRuleIDs, item.RuleInviteID, item.RuleTaskID, item.RuleInteractID, item.RuleTrainID, item.RuleOnlineID)
	}
	// 查询规则信息
	var relationRuleIdsWhere = map[string]interface{}{
		"enterprise_id": params.EnterpriseId,
		"id":            levelRuleIDs,
	}
	levelRules, err := excitationDao.NewIncentiveRelationRule().GetMany(c, relationRuleIdsWhere)

	levelRulesMap := make(map[int64]excitationModels.ExcitationIncentiveRelationRuleModel)
	for _, val := range levelRules {
		levelRulesMap[val.ID] = val
	}

	for key, value := range result {
		if _, ok := levelRulesMap[value.RuleInviteID]; ok {
			temp := s.handleLevelRulesMap(levelRulesMap, value.RuleInviteID)
			temp.RuleCondition = value.RuleConditionInvite
			result[key].LevelRuleAttachArr = append(result[key].LevelRuleAttachArr, temp)
		}
		if _, ok := levelRulesMap[value.RuleTaskID]; ok {
			temp := s.handleLevelRulesMap(levelRulesMap, value.RuleTaskID)
			temp.RuleCondition = value.RuleConditionTask
			result[key].LevelRuleAttachArr = append(result[key].LevelRuleAttachArr, temp)
		}
		if _, ok := levelRulesMap[value.RuleInteractID]; ok {
			temp := s.handleLevelRulesMap(levelRulesMap, value.RuleInteractID)
			temp.RuleCondition = value.RuleConditionInteract
			result[key].LevelRuleAttachArr = append(result[key].LevelRuleAttachArr, temp)
		}
		if _, ok := levelRulesMap[value.RuleTrainID]; ok {
			temp := s.handleLevelRulesMap(levelRulesMap, value.RuleTrainID)
			temp.RuleCondition = value.RuleConditionTrain
			result[key].LevelRuleAttachArr = append(result[key].LevelRuleAttachArr, temp)
		}
		if _, ok := levelRulesMap[value.RuleOnlineID]; ok {
			temp := s.handleLevelRulesMap(levelRulesMap, value.RuleOnlineID)
			temp.RuleCondition = value.RuleConditionOnline
			result[key].LevelRuleAttachArr = append(result[key].LevelRuleAttachArr, temp)
		}
	}

	return result, nil
}

// handleLevelRulesMap 处理等级附加条件匹配
func (s *LevelRuleService) handleLevelRulesMap(levelRulesMap map[int64]excitationModels.ExcitationIncentiveRelationRuleModel, ruleInviteID int64) excitationResp.LevelRuleAttachResp {

	var temp excitationResp.LevelRuleAttachResp
	temp.IncentiveRelationRuleId = levelRulesMap[ruleInviteID].ID
	temp.IncentiveConfigType = levelRulesMap[ruleInviteID].IncentiveConfigType
	temp.RuleCondition = levelRulesMap[ruleInviteID].RuleCondition
	temp.RuleGive = levelRulesMap[ruleInviteID].RuleGive

	return temp
}

// GetList 等级列表
func (s *LevelRuleService) GetList(c app.Ctx, params excitationReq.LevelRuleDetailReq) (result []excitationResp.LevelRuleDetailInfoResp, err error) {

	result, err = s.levelRuleDao.GetDetailList(c, params)
	if err != nil {
		return result, err
	}
	return result, nil
}
