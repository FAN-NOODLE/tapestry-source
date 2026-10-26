package workwechat

import workSDK "micro-scrm/internal/pkg/sdk/work"

// GetExternalUser 对外提供BA详情
func (w *WorkWeChat) GetExternalUser(workExternalUserID string) (*workSDK.ExternalContactResponse, error) {

	data, err := w.obj.ExternalDetail(w.ctx, workSDK.ExternalContactRequest{
		ExternalUserID: workExternalUserID,
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}
