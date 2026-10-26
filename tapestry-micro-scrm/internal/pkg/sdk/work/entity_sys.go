package work

// SCallbackAdd 外部联系人回调地址添加
type SCallbackAdd struct {
	CallType int    `json:"callType"`
	CorpID   string `json:"corpId"`
	CreateBy string `json:"createdBy"`
	ID       string `json:"id"`
	Uri      string `json:"uri"`
}

// SCallbackDel 删除外部联系人回调
type SCallbackDel struct {
	ID string `json:"id"`
}

// SBaseResponse UC基础回调
type SBaseResponse struct {
	Code      string `json:"code"`
	Msg       string `json:"msg"`
	RequestID string `json:"requestId"`
}

// SCallbackList 查询是否已经配置了
type SCallbackList struct {
	CorpID       string `json:"corpId"`
	CallbackType int    `json:"callType"`
}

// CallbackItem 外部联系人回调地址
type CallbackItem struct {
	ID           string `json:"id"`
	CorpID       string `json:"corpId"`
	Uri          string `json:"uri"`
	CallbackType int    `json:"callType"`
	CreateBy     string `json:"createdBy"`
	ModifiedBy   string `json:"modifiedBy"`
	CreateOn     int64  `json:"createdOn"`
	ModifiedOn   int64  `json:"modifiedOn"`
	DeletedOn    int64  `json:"deletedOn"`
	IsDel        int    `json:"isDel"`
}
