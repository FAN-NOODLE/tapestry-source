package reserveReq

const (
	ReserveOperaType_BA       = 1
	ReserveOperaType_Customer = 2
)

type LogList struct {
	LogInfo
	LogOrderInfo
	LogScoreInfo
	CreateAt    int64  `json:"create_at"`
	CreateAtStr string `json:"create_at_str"`
}

type LogInfo struct {
	MemberName  string `json:"member_name"`
	Event       string `json:"event"`
	OpreateName string `json:"opreate_name"` // 操作人姓名
}

type LogOrderInfo struct {
	ReserveId      int64  `json:"reserve_id"`
	ReserveTime    int    `json:"reserve_time"`
	StoreId        int64  `json:"store_id"`
	StoreName      string `json:"store_name"`
	StoreAddress   string `json:"store_address"`
	UserName       string `json:"user_name"`
	UserUnionid    string `json:"user_unionid"`
	UserExternalId string `json:"user_external_id"`
	ReserveRemark  string `json:"reserve_remark"`
	UserRemark     string `json:"user_remark"`
	ReserveName    string `json:"reserve_name"`
	ReserveTimeStr string `json:"reserve_time_str"`
}

type LogListArr struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type LogScoreInfo struct {
	ReserveNumber  string `json:"reserve_number"`
	ReserveOrderID int64  `json:"reserve_order_id"`
	UserName       string `json:"user_name"`
	UserUnionid    string `json:"user_unionid"`
	UserExternalID string `json:"user_external_id"`
	MemberID       int64  `json:"member_id"`
	MemberName     string `json:"member_name"`
	MemberUserID   string `json:"member_user_id"`
	StoreID        int64  `json:"store_id"`
	StoreName      string `json:"store_name"`
	ServiceScore   int32  `json:"service_score"`
	StoreScore     int32  `json:"store_score"`
	Remark         string `json:"remark"`
	FlowScore      int32  `json:"flow_score"`
	ReserveTime    int    `json:"reserve_time"`
}
