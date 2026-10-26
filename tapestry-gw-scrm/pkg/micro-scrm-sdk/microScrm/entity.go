package microScrm

// Response 返回
type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	RequestID string      `json:"request_id"`
	Data      interface{} `json:"data"`
}

// ResponseList 返回列表
type ResponseList struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Data      List   `json:"data"`
}

// List 列表结果
type List struct {
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
	List  interface{} `json:"list"`
}

type ID struct {
	ID int64 `json:"id"`
}

type Rows struct {
	Rows int64 `json:"rows"`
}
