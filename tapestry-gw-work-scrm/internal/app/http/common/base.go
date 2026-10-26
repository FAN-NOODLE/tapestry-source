package common

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/errorx"
	"gw-work-scrm/internal/pkg/sdk/boss"
	microAuthority "gw-work-scrm/internal/pkg/sdk/micro_authority"
)

type BaseHandle struct {
}

func NewBaseHandle() BaseHandle {
	return BaseHandle{}
}

func (h BaseHandle) GetMemberThirdsByType(c ctx.Context, enterpriseID int64, mobile string) (thirdIds []int, err error) {

	var args boss.EnterpriseAdminDetailRequest
	args.EnterpriseID = enterpriseID
	args.Mobile = mobile
	var params microAuthority.AccountDataRangeReq
	adminInfo, errs := boss.EnterpriseAdminDetail(ctx.Wrap(c), args)
	if errs != nil {
		return nil, errs
	}
	params.EnterpriseID = enterpriseID
	params.AdminId = adminInfo.ID
	thirdIds, errs = microAuthority.GetAccountDataRangeThirdIds(ctx.Wrap(c), params)
	if errs != nil {
		return nil, errs
	}
	if len(thirdIds) == 0 {
		return nil, errorx.NewErrorMsg(10090, "导购管理权限未设置")
	}

	return thirdIds, nil

}
