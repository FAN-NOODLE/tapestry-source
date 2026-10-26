package user

import "gw-scrm/internal/pkg/common"

// AccountDefaultLimit 账号默认数量
const AccountDefaultLimit = 5

// BusinessSourceDefault 注册来源默认（企业微信注册）
const BusinessSourceDefault = 1

// AvailableSystemDefaultWorker 默认开启智慧导购
const AvailableSystemDefaultWorker = 1

// AvailableSystemDefaultShopping 默认开启商城
const AvailableSystemDefaultShopping = 2

// TemplateDefaultId 智慧导购默认模板ID
const TemplateDefaultId = "tplc994927c018bef30"

// TemplateDefaultName 智慧导购名称
const TemplateDefaultName = "天鲟智慧导购"

// TemplateDefaultName 智慧导购名称
const ProviderCorpId = "ww0dc6c1288b9b294d"

// 返回参数结构
type GwResponse struct {
	Code      int64       `json:"code"`
	Message   string      `json:"message"`
	RequestId string      `json:"request_id"`
	Data      interface{} `json:"data"`
}

type UcResponse struct {
	Code      string      `json:"code"`
	Msg       string      `json:"msg"`
	RequestId string      `json:"requestId"`
	Data      interface{} `json:"data"`
}

// 返回参数结构
type GwTmpResponse struct {
	Code      int32  `json:"code"`
	Message   string `json:"message"`
	RequestId string `json:"request_id"`
	Data      []Tmp  `json:"data"`
}

// 登录返回数据 AdminResponse->List->Admin
type AdminResponse struct {
	Code      int64  `json:"code"`
	Message   string `json:"message"`
	RequestId string `json:"request_id"`
	Data      List   `json:"data"`
}

type List struct {
	EnterpriseAdminList []Admin     `json:"enterprise_admin_list"`
	EnterpriseList      []Admin     `json:"enterprise_list"`
	Token               string      `json:"token"`
	Version             int64       `json:"version"`
	EnterpriseTotal     interface{} `json:"enterprise_total"`
	Hash                string      `json:"hash"`
	Data                interface{} `json:"data"`
	EnterpriseId        int64       `json:"enterprise_id"`
	IsFirstLogin        bool        `json:"is_first_login"`
	BrandCode           string      `json:"brand_code"`
}
type Admin struct {
	AdminId              int64       `json:"admin_id"`
	Version              int64       `json:"version"`
	Hash                 string      `json:"hash"`
	Account              string      `json:"account"`
	AccountEmail         string      `json:"account_email"`
	AccountExpireAt      int64       `json:"account_expire_at"`
	AccountGlobalRoaming string      `json:"account_global_roaming"`
	AccountMobile        string      `json:"account_mobile"`
	AccountName          string      `json:"account_name"`
	AccountStartAt       int64       `json:"account_start_at"`
	AccountStatus        int8        `json:"account_status"`
	Address              string      `json:"address"`
	AdminAvailableSystem string      `json:"admin_available_system"`
	ContactMobile        string      `json:"contact_mobile"`
	ContactName          string      `json:"contact_name"`
	ContactStatus        int8        `json:"contact_status"`
	CorpId               string      `json:"corp_id"`
	CrsStatus            int8        `json:"crs_status"`
	Domain               string      `json:"domain"`
	EnterpriseId         int64       `json:"enterprise_id"`
	EnterpriseRoleList   interface{} `json:"enterprise_role_list"`
	GuidePermissions     interface{} `json:"guide_permissions"`
	GuideRoles           interface{} `json:"guide_roles"`
	Identifier           string      `json:"identifier"`
	IsSplitFlow          int8        `json:"is_split_flow"`
	IsSuper              int8        `json:"is_super"`
	Logo                 string      `json:"logo"`
	Name                 string      `json:"name"`
	RangeType            int8        `json:"range_type"`
	Role                 int64       `json:"role"`
	ServeType            int8        `json:"serve_type"`
	ShopName             string      `json:"shop_name"`
	StoreStatus          int8        `json:"store_status"`
	VersionTag           string      `json:"version_tag"`
	QrCode               string      `json:"qr_code"`
}

// 获取验证码传递参数
type GetCodeRequest struct {
	Account      string `json:"account"`
	Mobile       string `json:"mobile"`
	Client       string `json:"client"`
	Type         int8   `json:"type"`
	EnterpriseId int64  `json:"enterprise_id"`
}
type GetCodeByEPRequest struct {
	Account      string `json:"account"`
	Mobile       string `json:"mobile"`
	Client       string `json:"client"`
	Type         int8   `json:"type"`
	EnterpriseId int64  `json:"enterprise_id"`
}

// 获取绑定验证码传递参数
type GetBindCodeRequest struct {
	Account      string `json:"account"`
	Event        string `json:"event"`
	Mobile       string `json:"mobile"`
	Length       int8   `json:"length"`
	Type         int8   `json:"type"`
	Client       string `json:"client"`
	EnterpriseId int64  `json:"enterprise_id"`
}

// 登录传递参数
type LoginRequest struct {
	Account      string `json:"account"`
	Password     string `json:"password"`
	LoginType    int8   `json:"login_type"`
	Client       string `json:"client"`
	Code         string `json:"code"`
	CodeKey      string `json:"code_key"`
	EnterpriseId int64  `json:"enterprise_id"`
}

// 企业微信二维码传递参数
type QyCodeRequest struct {
	Ref string `json:"ref"`
}

// 企业微信登录传递参数
type QyLoginRequest struct {
	Type                 string      `json:"type"`
	AuthCode             string      `json:"auth_code"`
	Token                string      `json:"token"`
	Corpid               string      `json:"corpid"`
	WorkWechatLoginToken string      `json:"work_wechat_login_token"`
	RegisterEnterprise   interface{} `json:"register_enterprise"`
	BindAccount          interface{} `json:"bind_account"`
}

// 获取Token传递参数
type GetTokenRequest struct {
	AdminId      int64  `json:"admin_id"`
	EnterpriseId int64  `json:"enterprise_id"`
	Client       string `json:"client"`
	AppCode      string `json:"app_code"`
}

// 获取token返回数据
type TokenResponse struct {
	Code      int32       `json:"code"`
	Message   string      `json:"message"`
	RequestId string      `json:"request_id"`
	Data      TokenDetail `json:"data"`
}
type TokenDetail struct {
	EnterpriseAdminInfo Admin                `json:"enterprise_admin_info"`
	PermissionsList     interface{}          `json:"permissions_list"`
	PermissionsTree     interface{}          `json:"permissions_tree"`
	Token               string               `json:"token"`
	TokenBytes          common.TokenInfo     `json:"tokenBytes"`
	AdminList           map[int8]interface{} `json:"admin_list"`
	Version             int64                `json:"version"`
	Total               interface{}          `json:"total"`
	Hash                string               `json:"hash"`
	Data                interface{}          `json:"data"`
	JsApiConfig         interface{}          `json:"js_api_config"`
	RoleList            interface{}          `json:"role_list"`
}

// 获取js api
type JsApiRequest struct {
	Domain        string `json:"domain"`
	AppIdentifier string `json:"appIdentifier"`
}

// 获取旧版本Token
type GetTokenRequestV1 struct {
	EnterpriseId  int64  `json:"enterprise_id"`
	ServeType     int8   `json:"serve_type"`
	Role          int64  `json:"role"`
	Account       string `json:"account"`
	AccountMobile string `json:"account_mobile"`
	AdminId       int64  `json:"admin_id"`
}

/**
  行业列表传参
*/
type GetIndustry struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

/**
  验证邮箱是否已存在
*/
type CheckEmail struct {
	Email string `json:"email"`
}

type CheckMobileCode struct {
	Mobile           string `json:"mobile"`
	Client           string `json:"client"`
	CerificationCode string `json:"verification_code"`
	Code             string `json:"code"`
}

type GetRegisterCode struct {
	TemplateName   string `json:"templateName"`
	TemplateId     string `json:"templateId"`
	ProviderCorpId string `json:"providerCorpId"`
}

type AutoAccessParams struct {
	Version          int    `json:"version"`
	CorpId           string `json:"corp_id"`
	Name             string `json:"name"`
	IndustryId       int    `json:"industry_id"`
	Email            string `json:"email"`
	Contact          string `json:"contact"`
	Mobile           string `json:"mobile"`
	VerificationCode string `json:"verification_code"`
	MemUserId        string `json:"memUserid"`
}

/**
  自助开户
*/
type AutoAccess struct {
	Version         int           `json:"version"`
	CorpId          string        `json:"corp_id"`
	Name            string        `json:"name"`
	Contact         string        `json:"contact"`
	Mobile          string        `json:"mobile"`
	Email           string        `json:"email"`
	IndustryId      int           `json:"industry_id"`
	VersionTag      string        `json:"version_tag"`
	MemUserId       string        `json:"memUserId"`
	BusinessSource  int           `json:"business_source"`
	AccountLimit    int           `json:"account_limit"`
	AccountStartAt  int64         `json:"account_start_at"`
	AccountExpireAt int64         `json:"account_expire_at"`
	AermissionList  []int         `json:"permission_list"`
	AvailableSystem map[int8]int8 `json:"available_system"`
}

/**
  开通附加功能
*/
type AdditionalList struct {
	Version int `json:"version"`
}

type AdditionalLists struct {
}

type InstallThirdApp struct {
	BackUrl string `json:"back_url"`
}

// 企业微信登录
type WorkWechatLoginRequest struct {
	Type                 string `json:"type"`
	AuthCode             string `json:"auth_code"`
	Token                string `json:"token"`
	CorpId               string `json:"corp_id"`
	WorkWechatLoginToken string `json:"work_wechat_login_token"`
	RegisterEnterprise   string `json:"register_enterprise"`
	BindAccount          struct {
		Mobile string `json:"mobile"`
		Code   string `json:"code"`
	} `json:"bind_account"`
}

// 企业微信登录
type WorkWechatLoginRequestData struct {
	BindAccount          interface{} `json:"bind_account"`
	Type                 string      `json:"type"`
	AuthCode             string      `json:"auth_code"`
	Token                string      `json:"token"`
	CorpId               string      `json:"corpid"`
	WorkWechatLoginToken string      `json:"work_wechat_login_token"`
	RegisterEnterprise   interface{} `json:"register_enterprise"`
}

// 企业微信AuthInfo返回结果
type UcAuthInfo struct {
	EnterpriseInfo struct {
		EpId string `json:"epId"`
	} `json:"enterpriseInfo"`
	AdminInfo struct {
		EaId string `json:"eaId"`
	} `json:"adminInfo"`
}

// 企业微信ConfigInfo返回结果
type UcConfigInfo struct {
}

type UcAdminDetail struct {
	EaId string `json:"eaId"`
}

type Tmp struct {
	Id int `json:"id"`
}

type TmpToken struct {
	AdminId int `json:"admin_id"`
}

/**
  验证邮箱是否已存在
*/
type WorkLogin struct {
	Account   string `json:"account"`
	Password  string `json:"password"`
	LoginType int    `json:"login_type"`
}

type EmptyRequest struct {
}

type (
	OauthData struct {
		Token string `json:"token"`
	}
)
