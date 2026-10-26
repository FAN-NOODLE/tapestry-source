package planReq

import (
	"micro-scrm/internal/app/models/micro_scrm"
)

type FileDetection struct {
	EnterpriseID int64  `json:"enterprise_id"`
	FileName     string `json:"file_name"`
}

type CountInfo struct {
	Count int64 `json:"id"`
}
type ImportFileReq struct {
	EnterpriseID int64  `json:"enterprise_id"`
	FileName     string `json:"file_name"`
	FileUrl      string `json:"file_url"`
	TaskType     int    `json:"task_type"`
}

type ImportFileUpdateReq struct {
	EnterpriseID int64 `json:"enterprise_id"`
	ID           int64 `json:"id"`
}

type ImportFileRes struct {
	ID int64 `json:"id"`
}

type UnionIDData struct {
	UnionID    string `json:"unionid"`
	CustomerId string `json:"customer_id"`
	SaCode     string `json:"sa_code"`
}

type ImportInsertData struct {
	UnionID            string `json:"union_id" db:"union_id"`
	UserID             string `json:"user_id" db:"user_id"`
	Name               string `json:"name" db:"name"`
	WorkExternalUserId string `json:"work_external_user_id" db:"work_external_user_id"`
	UserName           string `json:"user_name" db:"user_name"`
}

type ImportInsertDataLog struct {
	ImportInsertData
	InsertLog string `json:"insert_log"`
}

type ConsumerList struct {
	ConsumerFileID int64 `json:"consumer_file_id"`
	EnterpriseID   int64 `json:"enterprise_id"`
	Page           int   `json:"page"`
	Size           int   `json:"size"`
	Total          int64 `json:"total"`
}

type ConsumerListRes struct {
	Total int64 `json:"total"`
	List  []microScrmModel.PlanTaskImportDataLogModel
}

type BirthdayUserInfo struct {
	MemberId int64  `json:"member_id"`
	Unionid  string `json:"unionid"`
	Openid   string `json:"openid"`
}

type ImportDataQueueData struct {
	EnterpriseId   int64  `json:"enterprise_id"`
	ConsumerFileId int64  `json:"consumer_file_id"`
	SendType       int    `json:"send_type"`
	TaskName       string `json:"task_name"`
	TaskId         int64  `json:"task_id"`
}

type TagTaskQueueData struct {
	EnterpriseId int64  `json:"enterprise_id"`
	TaskId       int64  `json:"task_id"`
	TaskName     string `json:"task_name"`
}
