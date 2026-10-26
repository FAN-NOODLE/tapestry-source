package activityCampaign

import "micro-scrm/internal/orm/model/micro_scrm_enterprise"

type UsersignInResponse struct {
	UnionId       string       `json:"union_id"`
	ActivitieDate string       `json:"activitie_date"`
	Message       string       `json:"message"`
	SingInDates   []SingInDate `json:"sing_in_dates"`
}
type SingInDate struct {
	Date            string `json:"date"`
	FormatDate      string `json:"format_date"`
	ActivitieStatus int    `json:"activitie_status"`
	Remark          string `json:"remark"`
}

type UsersignInInfoResponse struct {
	UnionId         string                                           `json:"union_id"`
	ActivitieDate   string                                           `json:"activitie_date"`
	FormatDate      string                                           `json:"format_date"`
	ActivitieStatus int                                              `json:"activitie_status"`
	Remark          string                                           `json:"remark"`
	SingInDate      string                                           `json:"sing_in_date"`
	SingInInfo      micro_scrm_enterprise.OperationActivitySignInLog `json:"sing_in_info"`
}
