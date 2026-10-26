package microScrm

// ResponseBase 返回
type ResponseBase struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	RequestID string      `json:"request_id"`
}

// Response 返回
type Response struct {
	ResponseBase
	Data      interface{} `json:"data"`
}

// ResponseList 返回列表
type ResponseList struct {
	ResponseBase
	Data      List   `json:"data"`
}

// List 列表结果
type List struct {
	ListBase
	List  interface{} `json:"list"`
}

type ListBase struct {
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}
type ID struct {
	ID int64 `json:"id"`
}

type Rows struct {
	Rows int64 `json:"rows"`
}

type RequestBase struct {
	EnterpriseID int64 `json:"enterprise_id"`
}