package excitationSvc

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	excitationDao "micro-scrm/internal/app/dao/excitation"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	"micro-scrm/internal/app/models"
	excitationModels "micro-scrm/internal/app/models/excitation"
	"micro-scrm/internal/app/service"
)

type IncentiveRelationService struct {
	service.Base
	incentiveRelationDao     *excitationDao.IncentiveRelation
	incentiveRelationRuleDao *excitationDao.IncentiveRelationRule
}

func NewIncentiveRelationService() *IncentiveRelationService {
	return &IncentiveRelationService{
		incentiveRelationDao:     excitationDao.NewIncentiveRelation(),
		incentiveRelationRuleDao: excitationDao.NewIncentiveRelationRule(),
	}
}

// RelationCreate 成长激励关联能力图谱规则创建
func (s *IncentiveRelationService) RelationCreate(c ctx.Context, params excitationReq.IncentiveRelationSaveReq) error {

	err := s.Base.DB().Transaction(func(tx *gorm.DB) error {
		relationData := excitationModels.ExcitationIncentiveRelationModel{
			MicroBaseModel:      models.CreateMicroBase(params.EnterpriseId, app.UID()),
			IncentiveConfigType: params.IncentiveConfigType,
			Status:              params.Status,
			AbilityIndex:        params.AbilityIndex,
			TimeRule:            params.TimeRule,
			IsDiffStore:         params.IsDiffStore,
		}
		err := excitationDao.NewIncentiveRelation().WithTx(tx).Create(c, relationData)
		if err != nil {
			return err
		}

		var relationRules []excitationModels.ExcitationIncentiveRelationRuleModel
		for _, rule := range params.RelationRule {
			relationRules = append(relationRules, excitationModels.ExcitationIncentiveRelationRuleModel{
				MicroBaseModel:      models.CreateMicroBase(params.EnterpriseId, app.UID()),
				IncentiveRelationID: relationData.ID,
				IncentiveConfigType: relationData.IncentiveConfigType,
				RuleCondition:       rule.RuleCondition,
				RuleValue:           rule.RuleValue,
				RuleGive:            rule.RuleGive,
				StoreType:           rule.StoreType,
			})
		}
		// 创建指标配置规则
		err = excitationDao.NewIncentiveRelationRule().WithTx(tx).CreateInBatchesData(c, relationRules)
		return err
	})

	return err
}

// RelationUpdate 成长激励关联能力图谱规则修改
func (s *IncentiveRelationService) RelationUpdate(c ctx.Context, params excitationReq.IncentiveRelationSaveReq) error {

	err := s.Base.DB().Transaction(func(tx *gorm.DB) error {

		var relationWhere = map[string]interface{}{
			"enterprise_id": params.EnterpriseId,
			"id":            params.ID,
		}
		var relationUpdateData = map[string]interface{}{
			"incentive_config_type": params.IncentiveConfigType,
			"status":                params.Status,
			"ability_index":         params.AbilityIndex,
			"time_rule":             params.TimeRule,
			"is_diff_store":         params.IsDiffStore,
		}
		err := excitationDao.NewIncentiveRelation().WithTx(tx).UpdateData(c, relationWhere, relationUpdateData)
		if err != nil {
			return err
		}

		// 删除之前规则
		var relationRuleWhere = map[string]interface{}{
			"enterprise_id":         params.EnterpriseId,
			"incentive_relation_id": params.ID,
		}
		err = excitationDao.NewIncentiveRelationRule().WithTx(tx).DeleteData(c, relationRuleWhere)

		var relationRules []excitationModels.ExcitationIncentiveRelationRuleModel
		for _, rule := range params.RelationRule {
			relationRules = append(relationRules, excitationModels.ExcitationIncentiveRelationRuleModel{
				MicroBaseModel:      models.CreateMicroBase(params.EnterpriseId, app.UID()),
				IncentiveRelationID: params.ID,
				IncentiveConfigType: params.IncentiveConfigType,
				RuleCondition:       rule.RuleCondition,
				RuleValue:           rule.RuleValue,
				RuleGive:            rule.RuleGive,
				StoreType:           rule.StoreType,
			})
		}
		// 创建指标配置规则
		err = excitationDao.NewIncentiveRelationRule().WithTx(tx).CreateInBatchesData(c, relationRules)
		return err
	})

	return err
}
