package openResp

type Success struct {
	StoreCode       string `json:"store_code" comment:"门店编码"`
	ReserveNumberId int64  `json:"reserve_number_id" comment:"预约单号"`
	MemberId        int64  `json:"member_id" comment:"员工id"`
	WechatUserId    string `json:"wechat_user_id" comment:"企微用户id"`
	MemberName      string `json:"member_name" comment:"员工姓名"`
	ContactQrcode   string `json:"contact_qrcode" comment:"联系我二维码"`
}

type AssignGuideSuccessResp struct {
	StoreCode     string `json:"store_code" comment:"门店编码"`
	Type          int    `json:"type" comment:"二维码场景：0-门店导购码，1-服务预约"`
	MemberId      int64  `json:"member_id" comment:"员工id"`
	WechatUserId  string `json:"wechat_user_id" comment:"企微用户id"`
	MemberName    string `json:"member_name" comment:"员工姓名"`
	ContactQrcode string `json:"contact_qrcode" comment:"联系我二维码"`
}
