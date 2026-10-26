package microAuthority

// 返回

type AccountDataRangeReq struct {
	EnterpriseID int64  `json:"enterprise_id"`
	AdminId      int64  `json:"admin_id"`
	Mobile       string `json:"mobile"`
}

type AccountDataRangeRes struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	RequestId string      `json:"request_id"`
	Data      interface{} `json:"data"`
}

type AccountDataRangeInfo struct {
	EnterpriseID int64 `json:"enterprise_id"`
	ID           int64 `json:"id"`
	ObjectId     int   `json:"object_id"`
	Type         int   `json:"type"`
}
