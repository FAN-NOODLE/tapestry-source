package workwechat

import (
	workSDK "micro-scrm/internal/pkg/sdk/work"
)

// 对外提供BA详情
func (w *WorkWeChat) BADetail(userId string) (*workSDK.MemberDetail, error) {
	data, err := w.obj.WorkMemberDetail(w.ctx, workSDK.MemberDetailRequest{
		QyGwRequest: workSDK.QyGwRequest{ReqPL: workSDK.QWCommunalContact},
		Userid:      userId})
	//data, err := w.obj.WorkMemberDetail(w.ctx, workSDK.MemberDetailRequest{Userid: userId})
	if err != nil {
		return nil, err
	}
	return data, nil
}

// DepartmentEmpInfo 查询部门下的员工信息
func (w *WorkWeChat) DepartmentEmpInfo(DepartmentID int64) ([]workSDK.MemberDetail, error) {
	data, err := w.obj.WorkMembers(w.ctx, workSDK.WorkMemberRequest{
		DepartmentID: DepartmentID,
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetWorkExternalUserInfo 根据work_external_user_id，获取用户的信息
func (w *WorkWeChat) GetWorkExternalUserInfo(externalUserID string) (*workSDK.ExternalContactResponse, error) {
	data, err := w.obj.ExternalDetail(w.ctx, workSDK.ExternalContactRequest{
		ExternalUserID: externalUserID,
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}
