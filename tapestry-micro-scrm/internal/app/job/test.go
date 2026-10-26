package appJob

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/pkg/xjob"
)

// Microsoft OAuth 配置
const (
	tenantID     = "74da624f-e27a-4b03-8702-f5e0a389e009"
	clientID     = "030b24ab-8fbb-4a01-8f32-3f791807ba56"
	clientSecret = "REDACTED_FOR_SOURCE_SHARING"
	redirectURI  = "https://tapestry-web-uat.coach.com.cn/web-admin/callback/sso"
)

// TokenResponse OAuth token 响应
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

// UserInfoResponse Microsoft Graph API 用户信息响应
type UserInfoResponse struct {
	ID                string `json:"id"`
	DisplayName       string `json:"displayName"`
	GivenName         string `json:"givenName"`
	Surname           string `json:"surname"`
	UserPrincipalName string `json:"userPrincipalName"`
	Mail              string `json:"mail"`
	JobTitle          string `json:"jobTitle"`
	Department        string `json:"department"`
	OfficeLocation    string `json:"officeLocation"`
}

func Test(c ctx.Context, jobHandler *xjob.JobHandler) error {
	// 示例：使用授权码获取 token 和用户信息
	code := ""
	// 1. 使用授权码获取 token
	tokenResp, err := GetAccessToken(code)
	if err != nil {
		fmt.Printf("获取 token 失败: %v\n", err)
		return err
	}
	fmt.Println("获取 token 成功")
	fmt.Printf("Token Type: %s, Expires In: %d\n", tokenResp.TokenType, tokenResp.ExpiresIn)

	// 2. 使用 access_token 获取用户信息
	userInfo, err := GetUserInfo(tokenResp.AccessToken)
	if err != nil {
		fmt.Printf("获取用户信息失败: %v\n", err)
		return err
	}
	fmt.Println("获取用户信息成功")
	fmt.Printf("用户ID: %s\n", userInfo.ID)
	fmt.Printf("用户名: %s\n", userInfo.DisplayName)
	fmt.Printf("邮箱: %s\n", userInfo.Mail)
	fmt.Printf("部门: %s\n", userInfo.Department)
	fmt.Printf("职位: %s\n", userInfo.JobTitle)

	return nil
}

// GetAccessToken 使用授权码获取访问令牌
func GetAccessToken(code string) (*TokenResponse, error) {
	tokenURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID)

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)
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
func GetUserInfo(accessToken string) (*UserInfoResponse, error) {
	userInfoURL := "https://graph.microsoft.com/v1.0/me"

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
