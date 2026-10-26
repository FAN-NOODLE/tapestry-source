package microboss

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gw-scrm/internal/pkg/sdk"
)

// AuthContactResponse 外部联系人调用凭证信息
type AuthContactResponse struct {
	sdk.ResponseBase
	Data struct {
		ExternalSecret string `json:"external_secret"`
	}
}

// WorkPluginResponse 自建应用信息
type WorkPluginResponse struct {
	sdk.ResponseBase
	Data struct {
		AgentId                int64  `json:"agent_id"`
		Secret                 string `json:"secret"`
		Domain                 string `json:"domain"`
		RedPacketsUrl          string `json:"red_packets_url"`
		CustomerInformationUrl string `json:"customer_information_url"`
		MaterialLibraryUrl     string `json:"material_library_url"`
		MarketingPromotionUrl  string `json:"marketing_promotion_url"`
		CheckFile              int8   `json:"check_file"`
		FileUrl                string `json:"file_url"`
	}
}

// CallBackResponse 回调返回信息
type CallBackResponse struct {
	sdk.ResponseBase
	Data struct {
		Key           string `json:"key"`
		Token         string `json:"token"`
		UcCallbackUrl string `json:"uc_callback_url"`
		IsSuper       int    `json:"is_super"`
	} `json:"data"`
}

type BuyNoticeRes struct {
	UserCount int64 `json:"user_count"`
	Days      int64 `json:"days"`
}

type AnnouncementRes struct {
	List []struct {
		Title      string `json:"title"`
		Content    string `json:"content"`
		ReleasedAt int64  `json:"released_at"`
	} `json:"list"`
	Page  int64 `json:"page"`
	Size  int64 `json:"size"`
	Total int64 `json:"total"`
}

func (ann *AnnouncementRes) CheckCode() error {
	return nil
}

func (b BuyNoticeRes) CheckCode() error {
	return nil
}

// SettingsUpdateResponse 设置更新返回数据
type SettingsUpdateResponse struct {
	sdk.ResponseBase
	Data struct{}
}

// WorkPluginResponse 自建应用信息
type AdminDetailResponse struct {
	Data struct {
		Id            int64  `json:"id"`
		EnterpriseId  int64  `json:"enterprise_id"`
		Username      string `json:"username"`
		Account       string `json:"account"`
		Password      string `json:"password"`
		Salt          string `json:"salt"`
		LastLoginTime int64  `json:"last_login_time"`
		Role          int64  `json:"role"`
		Status        int8   `json:"status"`
		IsSuper       int64  `json:"is_super"`
		Mobile        string `json:"mobile"`
		Enterprise    struct {
			Id            int64  `json:"id"`
			Hash          string `json:"hash"`
			Name          string `json:"name"`
			Logo          string `json:"logo"`
			Contact       string `json:"contact"`
			Mobile        string `json:"mobile"`
			Address       string `json:"address"`
			CrsStatus     int8   `json:"crs_status"`
			StoreStatus   int8   `json:"store_status"`
			ContactStatus int8   `json:"contact_status"`
			CorpId        string `json:"corp_id"`
			Domain        string `json:"domain"`
			Status        int8   `json:"status"`
			IsSplitFlow   int8   `json:"is_split_flow"`
		}
	}
}

func (a AdminDetailResponse) CheckCode() error {
	return nil
}

// GetAuthContact 获取外部联系人凭证
func GetAuthContact(c ctx.Context, request interface{}) (*AuthContactResponse, error) {
	var response AuthContactResponse
	err := httpClient.Post(c, GetUrl(urlEnum.AuthContactUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetWorkPlugin 获取自建应用信息
func GetWorkPlugin(c ctx.Context, request interface{}) (*WorkPluginResponse, error) {
	var response WorkPluginResponse
	err := httpClient.Post(c, GetUrl(urlEnum.WorkPluginUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetCallBack 获取客户联系信息
func GetCallBack(c ctx.Context, request interface{}) (*CallBackResponse, error) {
	var response CallBackResponse
	err := httpClient.Post(c, GetUrl(urlEnum.CallBackUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// 公告列表
func AnnouncementList(c ctx.Context, request interface{}) (*AnnouncementRes, error) {
	var response AnnouncementRes
	err := httpClient.Post(c, GetUrl(urlEnum.AnnouncementUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// AuthContactUpdate 调用凭证更新
func AuthContactUpdate(c ctx.Context, request interface{}) error {
	return httpClient.Post(c, GetUrl(urlEnum.AuthContactUpdateUrl), request, &SettingsUpdateResponse{})
}

// SetWxqyToken
func SetWxqyToken(c ctx.Context, request interface{}) error {
	return httpClient.Post(c, GetUrl(urlEnum.SetWxqyTokenUrl), request, &SettingsUpdateResponse{})
}

// WorkPluginUpdate 自建应用信息更新
func WorkPluginUpdate(c ctx.Context, request interface{}) error {
	return httpClient.Post(c, GetUrl(urlEnum.WorkPluginUpdateUrl), request, &SettingsUpdateResponse{})
}

// WorkPluginAdd 自建应用新增
func WorkPluginAdd(c ctx.Context, request interface{}) error {
	return httpClient.Post(c, GetUrl(urlEnum.WorkPluginAddUrl), request, &SettingsUpdateResponse{})
}

// SessionArchivingUpdate 会话存档自建应用信息更新
func SessionArchivingUpdate(c ctx.Context, request interface{}) error {
	return httpClient.Post(c, GetUrl(urlEnum.SessionArchivingUpdateUrl), request, &SettingsUpdateResponse{})
}

// SessionArchivingAdd 会话存档自建应用新增
func SessionArchivingAdd(c ctx.Context, request interface{}) error {
	return httpClient.Post(c, GetUrl(urlEnum.SessionArchivingAddUrl), request, &SettingsUpdateResponse{})
}

// 付费通知
func BuyNotice(c ctx.Context, request interface{}) (*BuyNoticeRes, error) {
	var response BuyNoticeRes
	err := httpClient.Post(c, GetUrl(urlEnum.BuyNotice), request, &response)
	if err != nil {
		return nil, err
	}
	return &response, err
}

// GetAdminDetail 获取管理员信息
func GetAdminDetail(c ctx.Context, request interface{}) (*AdminDetailResponse, error) {
	var response AdminDetailResponse
	err := httpClient.Post(c, GetUrl(urlEnum.AdminDetailUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

type AdminDetailListResponse struct {
	List []struct {
		Id           int64  `json:"id"`
		EnterpriseId int64  `json:"enterprise_id"`
		Username     string `json:"username"`
		Account      string `json:"account"`
	} `json:"data"`
	Page  int64 `json:"page"`
	Size  int64 `json:"size"`
	Total int64 `json:"total"`
}

func (a AdminDetailListResponse) CheckCode() error {
	return nil
}

// GetAdminDetailList 获取管理员信息列表
func GetAdminDetailList(c ctx.Context, request interface{}) (*AdminDetailListResponse, error) {
	var response AdminDetailListResponse
	err := httpClient.Post(c, GetUrl(urlEnum.AdminDetailListUrl), request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
