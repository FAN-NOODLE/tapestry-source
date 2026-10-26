package qyWork

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

const (
	// OpenUseridToUserid 将企微代开发应用的open_userid通过自建应用转成userid
	OpenUseridToUserid = "cgi-bin/batch/openuserid_to_userid"
	// OpenExternalUseridToExternalUserid 将企微代开发应用的open_external_userid通过自建应用转成external_userid
	OpenExternalUseridToExternalUserid = "cgi-bin/externalcontact/from_service_external_userid"
)

// OpenUseridToUserid 将企微代开发应用的open_userid通过自建应用转成userid
func (w *QyWork) OpenUseridToUserid(c ctx.Context, req OpenUseridToUseridReq) (OpenUseridToUseridResp, error) {
	var resp OpenUseridToUseridResp
	// 获取Token
	token, err := w.getToken(c, ManuCustomer)
	if err != nil {
		return resp, err
	}

	err = w.requestPostJson(c, w.getTokenUrl(OpenUseridToUserid, token), req, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// OpenExternalUseridToExternalUserid 将企微代开发应用的open_external_userid通过自建应用转成external_userid
func (w *QyWork) OpenExternalUseridToExternalUserid(c ctx.Context, req OpenExternalUseridToExternalUseridReq) (OpenExternalUseridToExternalUseridResp, error) {
	var resp OpenExternalUseridToExternalUseridResp
	// 获取Token
	token, err := w.getToken(c, ManuCustomer)
	if err != nil {
		return resp, err
	}

	err = w.requestPostJson(c, w.getTokenUrl(OpenExternalUseridToExternalUserid, token), req, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
