package api

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/modules/members"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type MemberApi struct {
}

func (m MemberApi) MemberDetail(c ctx.Context, req *members.MemberDetailRequest) interface{} {
	res, err := microscrm.GetMemberDetail(c, req)
	if err != nil {
		return nil
	}
	return res
}

func (m MemberApi) MemberList(c ctx.Context, req *members.MemberListRequest) interface{} {
	res, err := microscrm.GetMemberList(c, req)
	if err != nil {
		return nil
	}
	return res
}

func (m MemberApi) MomentList(c ctx.Context, req *members.MomentListRequest) interface{} {
	res, err := microscrm.MomentList(c, req)
	if err != nil {
		return nil
	}
	return res

}
func (m MemberApi) MomentBaList(c ctx.Context, req *members.MomentListRequest) interface{} {
	res, err := microscrm.MomentBaList(c, req)
	if err != nil {
		return nil
	}
	return res

}

func (m MemberApi) MemberListByThirdIds(c ctx.Context, req *members.MemberListByThirdIdsReq) interface{} {
	res, err := microscrm.MemberListByThirdIds(c, req)
	if err != nil {
		return nil
	}
	return res

}
