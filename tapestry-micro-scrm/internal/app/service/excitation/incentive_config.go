package excitationSvc

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	excitationDao "micro-scrm/internal/app/dao/excitation"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	excitationResp "micro-scrm/internal/app/http/response/excitation"
	"micro-scrm/internal/app/models"
	excitationModels "micro-scrm/internal/app/models/excitation"
	"micro-scrm/internal/app/service"
)

type IncentiveConfigService struct {
	service.Base
	incentiveConfigDao *excitationDao.IncentiveConfig
}

func NewIncentiveConfigService() *IncentiveConfigService {
	return &IncentiveConfigService{
		incentiveConfigDao: excitationDao.NewIncentiveConfig(),
	}
}

// GetList 维度配置列表关联能力图谱
func (s *IncentiveConfigService) GetList(c ctx.Context, params excitationReq.IncentiveConfigListReq) ([]excitationModels.ExcitationIncentiveConfigModel, int64, error) {
	result, total, err := s.incentiveConfigDao.GetMany(c, params)
	if err != nil {
		return result, total, err
	}

	return result, total, err
}

// GetConfigList 维度配置列表关联表数据列表
func (s *IncentiveConfigService) GetConfigList(c ctx.Context, params excitationReq.IncentiveConfigListReq) (configList []excitationResp.IncentiveConfigListResp, total int64, err error) {
	result, total, err := s.incentiveConfigDao.GetMany(c, params)
	if err != nil {
		return
	}

	var ruleParams excitationReq.RuleMany
	ruleParams.EnterpriseId = params.EnterpriseId
	ruleParams.Types = append(ruleParams.Types, excitationModels.IncentiveConfigTypeInvite, excitationModels.IncentiveConfigTypeTask, excitationModels.IncentiveConfigTypeInteract, excitationModels.IncentiveConfigTypeTrain, excitationModels.IncentiveConfigTypeOnline)

	// 查询维度规则关联
	integralRuleList, _, err1 := excitationDao.NewIntegralRule().GetMany(c, ruleParams)
	if err1 != nil {
		return
	}
	var integralRuleListMap = make(map[int64]excitationResp.IntegralRuleResp)
	for _, rules := range integralRuleList {
		var integralRuleTemp excitationResp.IntegralRuleResp
		integralRuleTemp.ID = rules.ID
		integralRuleTemp.Type = rules.Type
		integralRuleTemp.RuleName = rules.RuleName
		integralRuleListMap[rules.Type] = integralRuleTemp
	}

	for _, item := range result {
		var integralRule excitationResp.IntegralRuleResp
		if _, ok := integralRuleListMap[int64(item.Type)]; ok {
			integralRule = integralRuleListMap[int64(item.Type)]
		}

		configList = append(configList, excitationResp.IncentiveConfigListResp{
			EnterpriseId: item.EnterpriseID,
			Id:           item.ID,
			RuleName:     item.RuleName,
			Status:       item.Status,
			Type:         item.Type,
			IncentiveRelationResp: excitationResp.IncentiveRelationResp{
				Id:                  item.IncentiveRelation.ID,
				IncentiveConfigType: item.IncentiveRelation.IncentiveConfigType,
				Status:              item.IncentiveRelation.Status,
				AbilityIndex:        item.IncentiveRelation.AbilityIndex,
				TimeRule:            item.IncentiveRelation.TimeRule,
			},
			IntegralRuleResp: integralRule,
		})
	}

	return configList, total, err
}

// IncentiveConfigInit 初始化维度配置
func (s *IncentiveConfigService) IncentiveConfigInit(c ctx.Context, params excitationReq.IncentiveConfigListReq) (err error) {

	_, total, err := s.GetList(c, params)
	if err != nil {
		return err
	}
	if total > 0 {
		return nil
	}

	var ruleName = []string{"邀约维度", "任务维度", "客户互动维度", "培训维度", "线上能力"}

	var incentiveConfigModels []excitationModels.ExcitationIncentiveConfigModel
	for k, v := range ruleName {
		incentiveConfigModels = append(incentiveConfigModels, excitationModels.ExcitationIncentiveConfigModel{
			MicroBaseModel: models.MicroBaseModel{
				EnterpriseID: params.EnterpriseId,
			},
			RuleName: v,
			Status:   excitationModels.IncentiveConfigStatusOff,
			Type:     int32(k + 1),
		})
	}

	err = s.incentiveConfigDao.CreateInBatchesData(c, incentiveConfigModels)
	if err != nil {
		return err
	}
	return
}

func (s *IncentiveConfigService) RelationDetail(c ctx.Context, params excitationReq.IncentiveRelationDetailReq) (relationRes excitationResp.RelationDetailResp, err error) {

	var relationIdsWhere = map[string]interface{}{
		"enterprise_id": params.EnterpriseId,
		"id":            params.ID,
	}
	relationRes.RelationDetail, err = excitationDao.NewIncentiveRelation().GetOne(c, relationIdsWhere)
	if err != nil {
		return relationRes, err
	}

	// 查询规则信息
	var relationRuleIdsWhere = map[string]interface{}{
		"enterprise_id":         params.EnterpriseId,
		"incentive_relation_id": params.ID,
	}
	relationRuleDetails, err := excitationDao.NewIncentiveRelationRule().GetMany(c, relationRuleIdsWhere)

	if err != nil {
		return relationRes, err
	}

	for _, val := range relationRuleDetails {
		relationRes.RelationRuleDetails = append(relationRes.RelationRuleDetails,
			excitationResp.RelationRuleDetailResp{
				Id:                  val.ID,
				IncentiveConfigType: val.IncentiveConfigType,
				RuleCondition:       val.RuleCondition,
				RuleValue:           val.RuleValue,
				RuleGive:            val.RuleGive,
				StoreType:           val.StoreType,
			},
		)
	}

	return relationRes, err
}

func (s *IncentiveConfigService) RelationDetailList(c ctx.Context, params excitationReq.IncentiveRelationDetailListReq) (relationRes []excitationModels.ExcitationIncentiveRelationModel, err error) {

	var relationIdsWhere = map[string]interface{}{
		"enterprise_id": params.EnterpriseId,
		"status":        excitationModels.StatusTypeOn,
	}
	var ruleConditions []int
	ruleConditions = append(ruleConditions, 1, 2, 3, 4)
	relationRes, err = excitationDao.NewIncentiveRelation().GetDetailList(c, params.EnterpriseId, ruleConditions, relationIdsWhere)
	if err != nil {
		return relationRes, err
	}

	return relationRes, err

}
