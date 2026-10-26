package excitationSvc

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	excitationDao "micro-scrm/internal/app/dao/excitation"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	excitationResp "micro-scrm/internal/app/http/response/excitation"
	excitationModels "micro-scrm/internal/app/models/excitation"
	"micro-scrm/internal/app/service"
)

type IntegralRuleDescriptionService struct {
	service.Base
	integralRuleDescriptionDao *excitationDao.IntegralRuleDescription
}

func NewIntegralRuleDescriptionService() *IntegralRuleDescriptionService {
	return &IntegralRuleDescriptionService{
		integralRuleDescriptionDao: excitationDao.NewIntegralRuleDescription(),
	}
}

func (s *IntegralRuleDescriptionService) CreateOUpdate(c ctx.Context, params excitationReq.DescriptionCreateOrUpdateReq) error {

	ruleDetailData := excitationModels.ExcitationIntegralRuleDescriptionModel{
		Content:      params.Content,
		EnterpriseID: params.EnterpriseId,
	}
	err := s.integralRuleDescriptionDao.CreateOUpdate(c, ruleDetailData)
	if err != nil {
		return err
	}

	return nil
}

func (s *IntegralRuleDescriptionService) Detail(c ctx.Context, params excitationReq.DescriptionReq) (excitationResp.DescriptionResp, error) {

	var result excitationResp.DescriptionResp
	resp, err := s.integralRuleDescriptionDao.Detail(c, params)
	if err != nil {
		return result, err
	}
	return resp, err
}
