package uc

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/config"
)

const (
	CommunalContact = "communalContact"
)

// MemberByUserId 成员列表
func MemberByUserId(ctx ctx.Context, request MemberByUserIdRequest, identifier string, domain string) (*Member, error) {
	type member struct {
		Data Member `json:"data"`
	}
	response := &member{}
	url := fmt.Sprintf(urlEnum.MemberByUserId, config.CfgData.Service.Work.Domain, domain, identifier)

	err := Post(ctx, url, request, response)
	if err != nil {
		return nil, err
	}
	return &response.Data, nil
}
