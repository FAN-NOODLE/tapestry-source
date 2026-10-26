package notificationReq

import (
	"micro-scrm/internal/app/http/request"
)

type NoticeConfigSaveReq struct {
	request.EnterpriseBaseReq
	NoticeConfigList []NoticeConfigReq `json:"notice_config_list"`
}

type NoticeConfigReq struct {
	InfoType     int    `json:"info_type" `
	DeadlineWeek int    `json:"deadline_week"`
	DeadlineTime string `json:"deadline_time"`
	ReleaseWeek  int    `json:"release_week"`
	ReleaseTime  string `json:"release_time"`
}

type NoticeConfigDetailReq struct {
	request.EnterpriseBaseReq
}
