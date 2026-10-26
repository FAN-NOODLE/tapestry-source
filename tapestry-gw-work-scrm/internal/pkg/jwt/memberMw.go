package jwt

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-work-scrm/internal/config"
	microScrm "gw-work-scrm/internal/pkg/sdk/microscrm"
)

// ViableMembers 验证人员
func ViableMembers(c ctx.Context, path string, payload *CustomClaims) (member *microScrm.Member, err error) {
	// 通过userId 查询信息
	userId, err := payload.MemberUserId()
	if err != nil {
		return
	}

	memberSdk := microScrm.NewMember(config.CfgData.Service.MicroScrm.Domain)
	userInfo, err := memberSdk.DetailByUserid(
		c,
		microScrm.DetailUseridRequest{Userid: userId, EnterpriseID: payload.EnterpriseID},
	)
	if err != nil {
		return
	}

	// 代表没查询到用户信息
	if userInfo == nil || userInfo.ID == 0 {
		return
	}

	return userInfo, nil
}
