package workwechat

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/modules/enterprise"
	enterpriseRepository "micro-scrm/internal/modules/enterprise/respository/mysql"
	"micro-scrm/internal/pkg/sdk/work"
)

type WorkWeChat struct {
	obj *work.QyWorkImpl
	ctx ctx.Context
}

func NewWorkWeChat(c ctx.Context, enterpriseID int64) (*WorkWeChat, error) {
	// 根据企业ID获取企业信息
	enterpriseService := enterprise.New(enterpriseRepository.New())
	enterpriseData, enterpriseInfoErr := enterpriseService.Detail(c, enterprise.Corp{
		EnterpriseID: enterpriseID,
	})
	if enterpriseInfoErr != nil {
		log.Info(c).Msgf("获取企业信息错误:" + enterpriseInfoErr.Error())
		return nil, enterpriseInfoErr
	}

	return &WorkWeChat{
		obj: new(work.QyWorkImpl).SetDomain(enterpriseData.Domain),
		ctx: c,
	}, nil
}
