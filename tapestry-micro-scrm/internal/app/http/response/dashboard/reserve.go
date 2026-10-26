package dashboardResp

type DataKanban struct {
	SuccessMoney           int64 `json:"success_money" common:"成功金额"`
	SuccessMoneyProportion int64 `json:"success_money_proportion" common:"成功金额比例"`
	SendNumber             int64 `json:"send_number" common:"发送单数"`
	StoreNumber            int64 `json:"store_number" common:"到店单数"`
	StoreProportion        int64 `json:"store_proportion" common:"到店率"`
	SuccessNumber          int64 `json:"success_number" common:"成功单数"`
	SuccessProportion      int64 `json:"success_proportion" common:"成功率"`
	OnlineNumber           int64 `json:"online_number" common:"线上比例"`
	OfflineNumber          int64 `json:"offline_number" common:"线下比例"`
}

type DataList struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	SuccessMoney  int64  `json:"success_money" common:"成功金额"`
	OnlineNumber  int64  `json:"online_number" common:"线上比例"`
	OfflineNumber int64  `json:"offline_number" common:"线下比例"`
}
