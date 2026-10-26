package boss

import (
	"fmt"
	"micro-scrm/internal/app"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/json"
	"micro-scrm/internal/config"
	"micro-scrm/internal/pkg/http"
)

const (
	EnterpriseConfigRootDepartmentIdKey    = "root_department_id"    // 组织架构root部门third_id
	EnterpriseConfigSendWelcomeKey         = "send_welcome"          // 发送欢迎语 开关
	EnterpriseConfigStoreNameKey           = "store_name"            // 调用crm的企业识别标识
	EnterpriseConfigThirdAssignmentPushKey = "third_assignment_push" // 蔻驰外部联系人关系推送
)

// EnterpriseDetailRequest 企业详情请求
type EnterpriseDetailRequest struct {
	EnterpriseID int64  `json:"enterprise_id"`
	Hash         string `json:"hash"`
	CorpID       string `json:"corp_id"`
	Domain       string `json:"domain"`
}

// Enterprise Boss里企业信息返回结果
type Enterprise struct {
	ID                int64  `json:"id"`
	EnterpriseID      int64  `json:"enterprise_id"`
	Hash              string `json:"hash"`
	Name              string `json:"name"`
	Logo              string `json:"logo"`
	Contact           string `json:"contact"`
	Mobile            string `json:"mobile"`
	Address           string `json:"address"`
	IndustryID        int64  `json:"industry_id"`
	SeparateSwitch    int64  `json:"separate_switch"`
	LockTime          int32  `json:"lock_time"`
	Proportion        int64  `json:"proportion"`
	StartFee          int64  `json:"start_fee"`
	CrsStatus         int    `json:"crs_status"`
	ContactStatus     int    `json:"contact_status"`
	CorpID            string `json:"corp_id"`
	Domain            string `json:"domain"`
	ExternalSecret    string `json:"external_secret"`
	Status            int    `json:"status"`
	BusinessSource    int    `json:"business_source"`
	Version           int    `json:"version"`
	VersionTag        string `json:"version_tag"`
	AccountExpireAt   int64  `json:"account_expire_at"`
	AccountStartAt    int64  `json:"account_start_at"`
	AccountLimit      int32  `json:"account_limit"`
	ErpStatus         int    `json:"erp_status"`
	CrmStatus         int    `json:"crm_status"`
	ReturnAddress     string `json:"return_address"`
	IsSplitFlow       int    `json:"is_split_flow"`
	AvailableSystem   string `json:"available_system"`
	ServiceMerchantID int64  `json:"service_merchant_id"`
	EditAt            string `json:"edit_at"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         int64  `json:"updated_at"`
	DeletedAt         int64  `json:"deleted_at"`
	Identifier        string `json:"identifier"`
}

// EnterpriseListRequest 企业列表请求
type EnterpriseListRequest struct {
	Page  int64 `json:"page"`
	Limit int64 `json:"limit"`
}

// EnterpriseListResponse 企业列表响应
type EnterpriseListResponse struct {
	Data []*Enterprise `json:"data"`
	Meta Meta          `json:"meta"`
}

type Meta struct {
	Total int64 `json:"total"`
	Size  int64 `json:"size"`
	Page  int64 `json:"page"`
}

// WorkPlugin boss 返回的应用列表
type WorkPlugin struct {
	ID           int64  `json:"id"`
	EnterpriseID int64  `json:"enterprise_id"`
	CorpID       string `json:"corp_id"`
	AgentID      int64  `json:"agent_id"`
	Secret       string `json:"secret"`
	Identifier   string `json:"identifier"`
	PluginID     string `json:"plugin_id"`
}

// EnterpriseResponse 企业返回企业
type EnterpriseResponse struct {
	http.ResponseBase
	Data Enterprise `json:"data"`
}

// WorkPluginResponse 自建应用信息
type WorkPluginResponse struct {
	http.ResponseBase
	Data struct {
		AgentID    int64  `json:"agent_id"`
		Secret     string `json:"secret"`
		Identifier string `json:"identifier"`
		CorpID     string `json:"corp_id"`
	}
}

// EnterpriseConfigRequest 企业配置入参
type EnterpriseConfigRequest struct {
	EnterpriseId int64  `json:"enterprise_id"`
	Key          string `json:"key"`
}

// EnterpriseConfigResponse 企业配置
type EnterpriseConfigResponse struct {
	http.ResponseBase
	Data []EnterpriseConfigList `json:"data"`
}

type EnterpriseConfigList struct {
	Type  int64  `json:"type"`
	Key   string `json:"key"`
	Value int64  `json:"value"`
	Name  string `json:"name"`
}

// EnterpriseDepleteSmsRequest 消费企业短信
type EnterpriseDepleteSmsRequest struct {
	EnterpriseID int64 `json:"enterprise_id"`
	Number       int64 `json:"number"`
}

// SmsTotalResponse 获取短信数据返回值
type SmsTotalResponse struct {
	SmsNumber      int64 `json:"sms_number"`
	SmsTotalNumber int64 `json:"sms_total_number"`
}

// AuthDetailRequest 小程序授权详情请求参数
type AuthDetailRequest struct {
	EnterpriseID int64 `json:"enterprise_id"`
	AuthType     int   `json:"auth_type"`
}

// AuthDetailResponse 小程序授权详情返回值
type AuthDetailResponse struct {
	Data struct {
		Name     string `json:"name"`
		Logo     string `json:"logo"`
		AppID    string `json:"app_id"`
		OriginId string `json:"origin_id"`
	} `json:"data"`
}

type AdminListResponse struct {
	Id            int64       `json:"id"`
	EnterpriseId  int64       `json:"enterprise_id"`
	Username      string      `json:"username"`
	Account       string      `json:"account"`
	LastLoginTime int         `json:"last_login_time"`
	Status        int         `json:"status"`
	Role          int         `json:"role"`
	OperatorId    int         `json:"operator_id"`
	OperatorName  interface{} `json:"operator_name"`
	CreatedAt     string      `json:"created_at"`
	UpdatedAt     string      `json:"updated_at"`
}

type AdminDetailResponse struct {
	Data struct {
		ID            int64  `json:"id"`
		EnterpriseID  int64  `json:"enterprise_id"`
		Username      string `json:"username"`
		Account       string `json:"account"`
		Password      string `json:"password"`
		Mobile        string `json:"mobile"`
		Salt          string `json:"salt"`
		LastLoginTime int    `json:"last_login_time"`
		Role          int    `json:"role"`
		Status        int    `json:"status"`
		IsSuper       int    `json:"is_super"`
		Enterprise    struct {
			ID            int64         `json:"id"`
			Hash          string        `json:"hash"`
			Name          string        `json:"name"`
			Logo          string        `json:"logo"`
			Contact       string        `json:"contact"`
			Mobile        string        `json:"mobile"`
			Address       string        `json:"address"`
			CrsStatus     int           `json:"crs_status"`
			StoreStatus   int           `json:"store_status"`
			ContactStatus int           `json:"contact_status"`
			CorpID        string        `json:"corp_id"`
			Domain        string        `json:"domain"`
			Status        int           `json:"status"`
			IsSplitFlow   int           `json:"is_split_flow"`
			AppList       []interface{} `json:"app_list"`
		}
	}
}

// EnterpriseDetail 查询企业信息
func EnterpriseDetail(ctx ctx.Context, req *EnterpriseDetailRequest) (*Enterprise, error) {
	rs, err := http.ClientObject.Post(ctx, GetURL(urlEnum.BossEnterpriseDetail), clients.JSONBody(req))
	if err != nil {
		return nil, err
	}
	resStr, rErr := rs.Read(clients.Text)
	if rErr != nil {
		return nil, rErr
	}
	var response = EnterpriseResponse{}
	jErr := json.Unmarshal([]byte(resStr.(string)), &response)
	if jErr != nil {
		return nil, jErr
	}
	return &response.Data, nil
}

// EnterpriseList 查询企业信息
func EnterpriseList(ctx ctx.Context, req *EnterpriseListRequest) (*EnterpriseListResponse, error) {
	rs, err := http.ClientObject.Post(ctx, GetURL(urlEnum.BossEnterpriseList), clients.JSONBody(req))
	if err != nil {
		return nil, err
	}
	resStr, rErr := rs.Read(clients.Text)
	if rErr != nil {
		return nil, rErr
	}
	var response = EnterpriseListResponse{}
	jErr := json.Unmarshal([]byte(resStr.(string)), &response)
	if jErr != nil {
		return nil, jErr
	}
	return &response, nil
}

// WorkPluginList 查询应用列表
func WorkPluginList(ctx ctx.Context, req interface{}) ([]WorkPlugin, error) {
	url := fmt.Sprintf(EnterpriseWorkPluginList, config.CfgData.Service.Boss.Domain)
	res, errP := http.ClientObject.Post(ctx, url, clients.JSONBody(req))
	if errP != nil {
		return nil, errP
	}
	rt, errR := res.Read(clients.Text)
	if errR != nil {
		return nil, errR
	}

	var response = &[]WorkPlugin{}
	errJ := json.Unmarshal([]byte(rt.(string)), response)
	if errJ != nil {
		return nil, errJ
	}
	return *response, nil
}

// WorkPluginInfo 获取自建应用信息
func WorkPluginInfo(ctx ctx.Context, req interface{}) (*WorkPluginResponse, error) {
	url := fmt.Sprintf(WorkPluginUrl, config.CfgData.Service.Boss.Domain)
	res, errP := http.ClientObject.Post(ctx, url, clients.JSONBody(req))
	if errP != nil {
		return nil, errP
	}
	rt, errR := res.Read(clients.Text)
	if errR != nil {
		return nil, errR
	}

	var response = &WorkPluginResponse{}
	errJ := json.Unmarshal([]byte(rt.(string)), response)
	if errJ != nil {
		return nil, errJ
	}
	return response, nil
}

// EnterpriseDetail 查询企业信息
func EnterpriseSmsTotal(ctx ctx.Context, req *EnterpriseDetailRequest) (*SmsTotalResponse, error) {
	rs, err := http.ClientObject.Post(ctx, GetURL(urlEnum.BossEnterpriseSmsTotal), clients.JSONBody(req))
	if err != nil {
		return nil, err
	}
	resStr, rErr := rs.Read(clients.Text)
	if rErr != nil {
		return nil, rErr
	}
	var response = SmsTotalResponse{}
	jErr := json.Unmarshal([]byte(resStr.(string)), &response)
	if jErr != nil {
		return nil, jErr
	}
	return &response, nil
}

// EnterpriseDeplete 更新企业短信数量
func EnterpriseDeplete(ctx ctx.Context, req *EnterpriseDepleteSmsRequest) error {
	_, err := http.ClientObject.Post(ctx, GetURL(urlEnum.EnterpriseSmsDeplete), clients.JSONBody(req))
	if err != nil {
		return err
	}
	return nil
}

// AuthDetail 获取小程序授权信息
func AuthDetail(ctx ctx.Context, req *AuthDetailRequest) (*AuthDetailResponse, error) {
	url := fmt.Sprintf(authDetailUrl, config.CfgData.Service.Boss.Domain)
	res, errP := http.ClientObject.Get(ctx, url, clients.JSONBody(req))
	if errP != nil {
		return nil, errP
	}
	rt, errR := res.Read(clients.Text)
	if errR != nil {
		return nil, errR
	}

	var response = &AuthDetailResponse{}
	errJ := json.Unmarshal([]byte(rt.(string)), response)
	if errJ != nil {
		return nil, errJ
	}
	return response, nil
}

// EnterpriseConfig 获取自定义企业配置
func EnterpriseConfig(ctx ctx.Context, req interface{}) (*EnterpriseConfigResponse, error) {
	url := fmt.Sprintf(EnterpriseConfigInfo, config.CfgData.Service.Boss.Domain)
	var cacheString string
	redisKey := json.MarshalToString(req)

	if redisKey != "" {
		cacheStruct := EnterpriseConfigRequest{}
		_ = json.Unmarshal([]byte(redisKey), &cacheStruct)
		cacheKey := fmt.Sprintf("%d:%s", cacheStruct.EnterpriseId, cacheStruct.Key)
		err := app.Redis().Get(ctx, cacheKey).Scan(&cacheString)
		if err != nil {
			res, errP := http.ClientObject.Post(ctx, url, clients.JSONBody(req))
			if errP != nil {
				return nil, errP
			}
			rt, errR := res.Read(clients.Text)
			if errR != nil {
				return nil, errR
			}

			var response = &EnterpriseConfigResponse{}
			errJ := json.Unmarshal([]byte(rt.(string)), response)
			if errJ != nil {
				return nil, errJ
			}
			app.Redis().Set(ctx, cacheKey, rt.(string), 24*time.Hour)
			return response, nil
		} else {
			var cacheResponse = &EnterpriseConfigResponse{}
			errJson := json.Unmarshal([]byte(cacheString), cacheResponse)
			if errJson != nil {
				return nil, errJson
			}
			return cacheResponse, nil
		}

	} else {
		res, errP := http.ClientObject.Post(ctx, url, clients.JSONBody(req))
		if errP != nil {
			return nil, errP
		}
		rt, errR := res.Read(clients.Text)
		if errR != nil {
			return nil, errR
		}

		var response = &EnterpriseConfigResponse{}
		errJ := json.Unmarshal([]byte(rt.(string)), response)
		if errJ != nil {
			return nil, errJ
		}
		return response, nil
	}
	return nil, nil
}

// GetAdminDetail 获取管理员信息
func GetAdminDetail(ctx ctx.Context, request interface{}) (*AdminDetailResponse, error) {
	url := fmt.Sprintf(AdminDetailUrl, config.CfgData.Service.Boss.Domain)
	res, errP := http.ClientObject.Post(ctx, url, clients.JSONBody(request))
	if errP != nil {
		return nil, errP
	}
	rt, errR := res.Read(clients.Text)
	if errR != nil {
		return nil, errR
	}

	var response AdminDetailResponse
	errJ := json.Unmarshal([]byte(rt.(string)), &response)
	if errJ != nil {
		return nil, errJ
	}
	return &response, nil
}

// GetAdminList 获取管理员列表
func GetAdminList(ctx ctx.Context, request interface{}) ([]AdminListResponse, error) {
	url := fmt.Sprintf(AdminListUrl, config.CfgData.Service.Boss.Domain)
	res, errP := http.ClientObject.Post(ctx, url, clients.JSONBody(request))
	if errP != nil {
		return nil, errP
	}
	rt, errR := res.Read(clients.Text)
	if errR != nil {
		return nil, errR
	}
	var response struct {
		Data []AdminListResponse `json:"data"`
	}

	errJ := json.Unmarshal([]byte(rt.(string)), &response)
	if errJ != nil {
		return nil, errJ
	}
	return response.Data, nil
}
