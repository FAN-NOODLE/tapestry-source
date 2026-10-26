package qyWork

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

const (
	// MemberList 成员列表
	memberList = "cgi-bin/user/list"
	// MemberGet 成员详情
	memberGet = "cgi-bin/user/get"
)

// MemberList 成员列表
func (w *QyWork) MemberList(c ctx.Context, req MemberListReq) ([]MemberDetailResp, error) {
	var resp []MemberDetailResp
	// 获取Token
	token, err := w.getToken(c, Contact)
	if err != nil {
		return resp, err
	}

	err = w.requestPostJson(c, w.getTokenUrl(memberList, token), req, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// MemberGet 读取成员
func (w *QyWork) MemberGet(c ctx.Context, req ContactCreateRequest) (MemberDetailResp, error) {
	var resp MemberDetailResp
	// 获取Token
	token, err := w.getToken(c, Contact)
	if err != nil {
		return resp, err
	}

	err = w.requestPostJson(c, w.getTokenUrl(memberGet, token), req, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
