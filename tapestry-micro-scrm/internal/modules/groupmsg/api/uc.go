package api

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/pkg/sdk/uc"
)

type UCApi struct {
}

// GetMemUID 获取UC memUid
func (U *UCApi) GetMemUID(c ctx.Context, userID string, domain string) (string, error) {
	// 从UC 查询企微成员的memUid,等待智慧导购业绩分享兼容企微userid，须得删除所有关于UC代码
	memberInfo, err := uc.MemberByUserId(c, uc.MemberByUserIdRequest{
		MemUserID: userID,
	}, "crs-manu-customer", domain)
	if err != nil {
		log.Info(c).Msgf("获取UC memUid失败, err:%s, userID:%s, domain:%s", err.Error(), userID, domain)
		return "", err
	}
	return memberInfo.MemUid, nil
}
