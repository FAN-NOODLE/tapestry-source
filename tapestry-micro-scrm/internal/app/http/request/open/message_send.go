package openReq

type SendWaitPaymentRequest struct {
	OrderNo      string `json:"orderno" binding:"required" binding:"订单号必填"`
	UserId       string `json:"userid" binding:"required" binding:"员工编号必填"`
	Content      string `json:"content" binding:"required" binding:"通知内容必填"`
	EnterpriseID int64  `json:"enterprise_id"`
}

type SendWaitPaymentOrderRes struct {
	PayeeUserid string `json:"payee_userid"`
}
