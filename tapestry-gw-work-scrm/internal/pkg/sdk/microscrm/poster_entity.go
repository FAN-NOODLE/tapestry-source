package microScrm

// PosterRequest 生成海报请求
type PosterRequest struct {
	Type           int    `json:"type" desc:"二维码类型"`
	OperatorType   int    `json:"operator_type" desc:"操作类型 1 发送 2 下载"`
	MemUserID      string `json:"member_userid" desc:"导购userid"`
	MemberName     string `json:"member_name" desc:"员工姓名"`
	MemberAvatar   string `json:"member_avatar" desc:"员工头像"`
	EnterpriseID   int64  `json:"enterprise_id" desc:"企业ID"`
	EnterpriseHash string `json:"enterprise_hash" desc:"企业ID"`
	SunCode        string `json:"sun_code" desc:"太阳码"`
	Params         string `json:"params" desc:"params海报参数"`
	FileTmpUrl     string `json:"file_tmp_url" desc:"海报临时地址"`
	ShareImg       string `json:"share_img" desc:"分享图片地址"`
}
