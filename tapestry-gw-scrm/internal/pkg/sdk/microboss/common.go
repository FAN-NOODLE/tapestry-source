package microboss

import (
	"fmt"
	"gw-scrm/internal/config"
	"gw-scrm/internal/pkg/sdk"
)

var urlEnum = struct {
	AuthContactUrl            string
	WorkPluginUrl             string
	CallBackUrl               string
	AuthContactUpdateUrl      string
	SetWxqyTokenUrl           string
	WorkPluginAddUrl          string
	WorkPluginUpdateUrl       string
	EnterpriseInfoUrl         string
	EnterpriseInfoSetUrl      string
	EnterpriseInfoGetUrl      string
	AnnouncementUrl           string
	SessionArchivingAddUrl    string
	SessionArchivingUpdateUrl string
	BuyNotice                 string
	AdminDetailUrl            string
	AdminDetailListUrl        string
}{
	// 获取外部联系人凭证
	AuthContactUrl: "%s/service/v2/auth/auth-contact",
	// 获取自建应用信息
	WorkPluginUrl: "%s/service/v2/auth/get-work-plugin",
	// 获取客户联系信息
	CallBackUrl: "%s/service/v2/auth/uc-callback",
	// 调用凭证更新
	AuthContactUpdateUrl: "%s/service/v2/auth/work-setting",
	// 客户联系修改token
	SetWxqyTokenUrl: "%s/service/v2/auth/set-wxqy-token",
	// 自建应用添加
	WorkPluginAddUrl: "%s/service/v2/auth/install-app",
	// 自建应用信息更新
	WorkPluginUpdateUrl: "%s/service/v2/auth/update-app",
	// 获取企业信息
	EnterpriseInfoUrl: "%s/service/v1/enterprise/detail",
	// 读取企业配置信息
	EnterpriseInfoGetUrl: "%s/service/v1/enterprise-config/info",
	// 设置企业配置信息
	EnterpriseInfoSetUrl: "%s/service/v1/enterprise-config/save",
	// 公告通知
	AnnouncementUrl: "%s/service/v1/announcement/list",
	// 会话存档自建应用添加
	SessionArchivingAddUrl: "%s/service/v2/auth/install-session-archiving-app",
	// 会话存档自建应用信息更新
	SessionArchivingUpdateUrl: "%s/service/v2/auth/update-session-archiving-app",
	// 续费通知
	BuyNotice: "%s/service/v1/enterprise/buy-notice",
	// 管理员详情
	AdminDetailUrl: "%s/service/v1/enterprise_admin/detail",
	// 管理员详情列表
	AdminDetailListUrl: "%s/service/v1/enterprise_admin/list",
}

// GetUrl 获取请求链接
func GetUrl(uri string) string {
	return fmt.Sprintf(uri, config.CfgData.Service.Boss.Domain)
}

var httpClient = sdk.NewHttpClient(sdk.ParseJson{})

type NoCheckResponse struct {
}

func (e *NoCheckResponse) CheckCode() error {
	return nil
}
