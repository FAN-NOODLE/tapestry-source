package boss

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/json"
	"gw-work-scrm/internal/config"
	"gw-work-scrm/internal/pkg/http"
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
	BrandCode         string `json:"brand_code"`
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

type EnterpriseAdminDetailRequest struct {
	EnterpriseID int64  `json:"enterprise_id"`
	Mobile       string `json:"mobile"`
}

type EnterpriseAdminDetailResponse struct {
	http.ResponseBase
	Data EnterpriseAdminDetailRes `json:"data"`
}

type EnterpriseAdminDetailRes struct {
	EnterpriseID int64  `json:"enterprise_id"`
	ID           int64  `json:"id"`
	Mobile       string `json:"mobile"`
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

// EnterpriseAdminDetail 查询admin信息
func EnterpriseAdminDetail(ctx ctx.Context, req EnterpriseAdminDetailRequest) (enterpriseAdmin Enterprise, err error) {
	rs, err := http.ClientObject.Post(ctx, GetURL(urlEnum.EnterpriseAdminDetail), clients.JSONBody(req))
	if err != nil {
		return enterpriseAdmin, err
	}
	resStr, rErr := rs.Read(clients.Text)
	if rErr != nil {
		return enterpriseAdmin, rErr
	}
	var response = EnterpriseResponse{}
	jErr := json.Unmarshal([]byte(resStr.(string)), &response)
	if jErr != nil {
		return enterpriseAdmin, jErr
	}
	return response.Data, nil
}
