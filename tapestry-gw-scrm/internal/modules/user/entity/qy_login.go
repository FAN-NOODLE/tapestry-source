package entity

type WorkWechatLoginResponse struct {
	Code                 string `json:"code"`
	Desc                 string `json:"desc"`
	AuthEnterprise       *BossEnterpriseDetailResponse
	Admin                struct{} `json:"admin"`
	TempToken            string   `json:"temp_token"`
	WorkWechatLoginToken string   `json:"work_wechat_login_token"`
}

type WorkWechatLoginRedisKey struct {
	AuthEffective bool     `json:"auth_effective"`
	IsExpire      bool     `json:"is_expire"`
	Enterprise    struct{} `json:"enterprise"`
	UcEnterprise  struct{} `json:"uc_enterprise"`
	CorpId        string   `json:"corp_id"`
	UcMemberUid   string   `json:"uc_member_uid"`
}

type UcEnSettingQuery struct {
	TimeFormat            string `json:"timeformat"`
	DateFormat            string `json:"dateformat"`
	SiteName              string `json:"sitename"`
	WxMchid               string `json:"wxMchid"`
	WxMchkey              string `json:"wxMchkey"`
	WxQyPkcs12            string `json:"wxqyPkcs12"`
	WxMpToken             string `json:"wxmpToken"`
	WxMpEncodingAeskey    string `json:"wxmpEncodingAeskey"`
	WxMpAppId             string `json:"wxmpAppid"`
	WxMpAppSecret         string `json:"wxmpAppsecret"`
	WxMpQrCode            string `json:"wxmpQrcode"`
	WxQyEncodingAeskey    string `json:"wxqyEncodingAeskey"`
	WxQyCorpId            string `json:"wxqyCorpid"`
	WxQyCorpSecret        string `json:"wxqyCorpsecret"`
	WxQyQrCode            string `json:"wxqyQrcode"`
	WxQyToken             string `json:"wxqyToken"`
	ContactIdentifier     string `json:"contactIdentifier"`
	StudyCenterIdentifier string `json:"studycenterIdentifier"`
	ResourcesDomain       string `json:"resourcesDomain"`
	WxExternalSecret      string `json:"wxExternalSecret"`
	WxExternalCallback    string `json:"wxExternalCallback"`
	MemberExtSyncTag      string `json:"memberExtSyncTag"`
}

type AccountBaseRes struct {
	Data struct {
		Id      int64  `json:"id"`
		Account string `json:"account"`
	} `json:"data"`
}

/*type ucAdminDetail struct {
	memUid string `json:"memUid"`
}

type ucAdminDetail struct {
	memUid string `json:"memUid"`
}*/
