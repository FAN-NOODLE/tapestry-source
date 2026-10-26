package ssoSvc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"micro-scrm/internal/orm/model/micro_scrm"
	"net/http"
	"net/url"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	enterpriseSSODAO "micro-scrm/internal/app/dao/enterprise_sso_config"
	"micro-scrm/internal/app/service"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
)

// MicrosoftOAuthConfig  OAuth 配置结构
type MicrosoftOAuthConfig struct {
	TenantID     string
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

// TokenResponse OAuth token 响应
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

// UserInfoResponse Microsoft Graph API 用户信息响应
type UserInfoResponse struct {
	EmployeeId string `json:"employeeId"`
}

type SSOService struct {
	service.Base
	enterpriseSSODAO *enterpriseSSODAO.EnterpriseSSOConfig
	httpClient       *clients.HTTPClient
}

func NewSSOService() *SSOService {
	return &SSOService{
		enterpriseSSODAO: enterpriseSSODAO.NewEnterpriseSSOConfig(),
		httpClient:       clients.NewHTTPClient(clients.Timeout(30 * time.Second)),
	}
}

// GetMicrosoftConfig 从数据库获取 Microsoft OAuth 配置
func (s *SSOService) GetMicrosoftConfig(c ctx.Context, enterpriseId int64) (micro_scrm.SSOConfigModel, error) {
	return s.enterpriseSSODAO.GetConfigByEnterpriseIdAndType(c, enterpriseId, "microsoft")
}

// GetAccessToken 使用授权码获取访问令牌
func (s *SSOService) GetAccessToken(code string, config micro_scrm.SSOConfigModel) (*TokenResponse, error) {
	tokenURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", config.TenantID)

	data := url.Values{}
	data.Set("client_id", config.ClientID)
	data.Set("client_secret", config.ClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", config.RedirectURI)
	data.Set("grant_type", "authorization_code")
	data.Set("scope", "User.Read")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.PostForm(tokenURL, data)
	if err != nil {
		return nil, fmt.Errorf("请求 token 失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("获取 token 失败, 状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("解析 token 响应失败: %w", err)
	}

	return &tokenResp, nil
}

// GetUserInfo 使用访问令牌获取用户信息
func (s *SSOService) GetUserInfo(accessToken string) (*UserInfoResponse, error) {
	userInfoURL := "https://graph.microsoft.com/v1.0/me?$select=employeeId"

	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求用户信息失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("获取用户信息失败, 状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 读取完整的响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	// 打印完整的 JSON 响应（调试用）
	fmt.Printf("\n=== 完整的用户信息 JSON 响应 ===\n")
	fmt.Println(string(body))
	fmt.Println("================================\n")

	// 使用 map 解析，方便查看所有字段
	var rawData map[string]interface{}
	if err := json.Unmarshal(body, &rawData); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %w", err)
	}

	// 打印解析后的 map（调试用）
	fmt.Printf("=== 解析后的 map 数据 ===\n")
	for key, value := range rawData {
		fmt.Printf("%s: %v\n", key, value)
	}
	fmt.Println("=========================\n")

	// 再解析到结构体
	var userInfo UserInfoResponse
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("解析用户信息响应失败: %w", err)
	}

	return &userInfo, nil
}

// Login SSO 登录，使用授权码获取用户信息
func (s *SSOService) Login(c ctx.Context, code string, enterpriseId int64) (*UserInfoResponse, error) {
	// 1. 从数据库获取配置
	config, err := s.GetMicrosoftConfig(c, enterpriseId)
	if err != nil {
		return nil, err
	}
	if config.ID == 0 {
		return nil, errors.New("SSO未配置，请联系管理员")
	}

	// 2. 获取 token
	tokenResp, err := s.GetAccessToken(code, config)
	if err != nil {
		return nil, err
	}

	// 3. 获取用户信息
	return s.GetUserInfo(tokenResp.AccessToken)

}
