package config

import (
	"os"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/pkg/db"
	"micro-scrm/internal/pkg/redis"
)

// Config 配置组合
type Config struct {
	App         string           `json:"app"`
	Mode        string           `json:"mode"`
	Env         string           `json:"env"`
	QYNoticeKey string           `yaml:"qy_notice_key"`
	NoticeKeys  NoticeKeyCollect `yaml:"notice_keys"`
	// log config
	Logger     LoggerCfg    `json:"logger"`
	DB         db.Config    `json:"db"`
	ClickHouse db.Config    `json:"clickhouse"`
	Redis      redis.Config `json:"redis"`
	Restful    RestfulCfg   `json:"restful"`
	Monitor    MonitorCfg   `json:"monitor"`
	Xid        string       `json:"xid"`
	XIds       string       `json:"xIds"  yaml:"xIds"`
	Service    ServiceCfg   `json:"service"`
	Cmq        Cmq          `json:"cmq"`
	KMS        KMSConfig    `json:"kms" yaml:"kms"`
	Mp         MpConfig     `json:"mp"`
	UploadPath string       `json:"upload_path" yaml:"upload_path"`
	TrustedIps []string     `json:"trustedIps" yaml:"trustedIps"`
	Job        JobCfg       `json:"job" yaml:"job" usage:"定时任务相关配置"`
	Storage    StorageCfg   `json:"storage"`
	WebHtmlUrl WebHtmlUrl   `json:"webhtmlurl"`
	Elastic    Elastic      `json:"elastic"`
	Kafka      KafkaCfg     `json:"kafka"`

	CampaignToolService CampaignToolService `json:"campaign_tool_service" yaml:"campaign_tool_service"`

	JWT           JwtCfg              `json:"jwt"`
	Inventory     InventoryConfig     `yaml:"inventory"`
	InventorySw   InventorySwConfig   `yaml:"inventory_sw"`
	VchangyiEmail VchangyiEmailConfig `yaml:"vchangyi_email"`
	Baozun        BaozunConfig        `yaml:"baozun"`
	BaiQiu        BaiQiuConfig        `yaml:"bai_qiu"`
	CoachSaTag    CoachSaTag          `yaml:"coach_sa_tag"`
	CoachKOS      CoachKOSConfig      `yaml:"coach_kos"`
	Talkie        TalkieConfig        `yaml:"talkie"`
	CgeMa         CgeMaConfig         `yaml:"cge-ma"`
}

// 存储配置
type WebHtmlUrl struct {
	FissionActivityUrl string `json:"fissionactivityurl" yaml:"fissionactivityurl"`
}

type Elastic struct {
	EsAddress  []string `json:"es_address" yaml:"es_address"`
	EsUser     string   `json:"es_user" yaml:"es_user"`
	EsPassword string   `json:"es_password" yaml:"es_password"`
}
type JobCfg struct {
	SubCmd     string `json:"cmd" yaml:"cmd" usage:"执行子命令"`
	Spec       string `json:"spec" yaml:"spec"  usage:"可运行的特定标识。如果为空表示运行所有"`
	DelCmdSpec string `json:"delCmdSpec" yaml:"delCmdSpec"  usage:"删除灰度企业缓存"`
}

// Cmq cmq配置
type Cmq struct {
	SecretId  string `yaml:"secretId"`
	SecretKey string `yaml:"secretKey"`
	Endpoint  string `yaml:"endpoint"`
}

// KMSConfig 腾讯云密钥管理服务配置。
type KMSConfig struct {
	Region    string `yaml:"region"`
	KeyID     string `yaml:"key_id"`
	SecretID  string `yaml:"secret_id"`
	SecretKey string `yaml:"secret_key"`
	Endpoint  string `yaml:"endpoint"`
}

// Log 返回日志配置对象
func (c Config) Log(logToStderr bool) log.Config {
	if logToStderr {
		return log.Config{
			App:    c.App,
			Level:  c.Logger.Level,
			Format: c.Logger.Format,
			Output: os.Stderr,
		}
	}
	return log.Config{
		App:        c.App,
		Level:      c.Logger.Level,
		TargetType: c.Logger.TargetType,
		Target:     c.Logger.Target,
		Format:     c.Logger.Format,
	}
}

// LoggerCfg 日志配置
type LoggerCfg struct {
	Level      string `json:"level" usage:"日志级别"`
	TargetType string `json:"target_type" usage:"输出目标类型【console/file/tcp/udp】"`
	Target     string `json:"target" usage:"日志输出格式【<file path>/<network address>】"`
	Format     string `json:"format" usage:"日志输出格式【json/raw】"`
}

// Cors 跨域配置
type Cors struct {
	Enable           bool          `yaml:"enable"`
	AllowAllOrigins  bool          `yaml:"AllowAllOrigins"`
	AllowOrigins     []string      `yaml:"AllowOrigins"`
	AllowMethods     []string      `yaml:"AllowMethods"`
	AllowHeaders     []string      `yaml:"AllowHeaders"`
	AllowCredentials bool          `yaml:"AllowCredentials"`
	ExposeHeaders    []string      `yaml:"ExposeHeaders"`
	MaxAge           time.Duration `yaml:"MaxAge"`
	AllowWebSockets  bool          `yaml:"AllowWebSockets"`
	AllowFiles       bool          `yaml:"AllowFiles"`
}

// Doc swagger 文档配置
type Doc struct {
	Title         string `yaml:"title"`
	Description   string `yaml:"description"`
	TermOfService string `yaml:"termOfService"`
	Contact       string `yaml:"contact"`
}

// RestfulCfg api 接口配置
type RestfulCfg struct {
	Host     string `yaml:"host"`
	BasePath string `yaml:"basePath"`
	Addr     string `yaml:"addr"`
	Cors     Cors   `yaml:"cors"`
	Doc      Doc    `yaml:"doc"`
}

// MonitorCfg 监控配置
type MonitorCfg struct {
	Enable bool   `yaml:"enable"`
	Addr   string `yaml:"addr"`
}

// ServiceCfg SDK 配置
type ServiceCfg struct {
	Work           WorkCfg       `yaml:"work"`
	Shop           ServiceConfig `yaml:"shop"`
	Boss           ServiceConfig `yaml:"boss"`
	Console        ServiceConfig `yaml:"console"`
	Open           ServiceConfig `yaml:"open"`
	MicroPay       ServiceConfig `yaml:"micro-pay"`
	MicroGoods     ServiceConfig `yaml:"micro-goods"`
	MicroActivity  ServiceConfig `yaml:"micro-activity"`
	MicroAuthority ServiceConfig `yaml:"micro-authority"`
	MicroUser      ServiceConfig `yaml:"micro-user"`
	Cdp            ServiceConfig `yaml:"cdp"`
	SendCloud      ServiceConfig `yaml:"send-cloud"`
	WebH5          ServiceConfig `yaml:"web-h5"`
	MicroScrm      ServiceConfig `yaml:"micro-scrm"`
	MicroWork      ServiceConfig `yaml:"micro-work"`
	SortLinks      ServiceConfig `yaml:"sort-links"`
	MicroChatScrm  ServiceConfig `yaml:"micro-chat-scrm"`
}

// ServiceConfig SDK配置
type ServiceConfig struct {
	Domain string `yaml:"domain"`
}

// CgeMaConfig 深演 MA 回调配置
type CgeMaConfig struct {
	Domain    string `yaml:"domain"`
	AppKey    string `yaml:"app-key"`
	AppSecret string `yaml:"app-secret"`
}

// WorkCfg 企业微信配置
type WorkCfg struct {
	Domain string `yaml:"domain"`
	ID     string `yaml:"id"`
	Secret string `yaml:"secret"`
}

type MpConfig struct {
	AppID  string `yaml:"app_id"`
	Secret string `yaml:"secret"`
}

// 存储配置
type StorageCfg struct {
	TmpDir string       `yaml:"TmpDir"`
	COS    CloudStorage `yaml:"cos"`
	Blob   BlobStorage  `yaml:"blob"`
}

// 云存储结构体
type CloudStorage struct {
	CdnHost       string `yaml:"cdnHost"`
	Host          string `yaml:"host"`
	Schema        string `yaml:"schema"`
	SecretID      string `yaml:"SecretID"`
	SecretKey     string `yaml:"SecretKey"`
	Bucket        string `yaml:"bucket"`
	PrivateBucket string `yaml:"privateBucket"`
	Origin        string `yaml:"origin"`
	AppID         string `yaml:"appID"`
}

// 微软云存储配置
type BlobStorage struct {
	Domain      string `yaml:"domain"`
	Bucket      string `yaml:"bucket"`
	AccountName string `yaml:"accountName"`
	AccountKey  string `yaml:"accountKey"`
}

// Kafka 配置
type KafkaCfg struct {
	Addr    []string `yaml:"addr"`
	Version string   `yaml:"version"`
}

// CampaignToolService Campaign-Tool服务
type CampaignToolService struct {
	AES_CBC_256 AES_CBC_256 `yaml:"aes_cbc_256"`
}

type AES_CBC_256 struct {
	Key string `yaml:"key"`
	Iv  string `yaml:"iv"`
}

// JwtCfg JWT配置
type JwtCfg struct {
	Secret string `yaml:"secret"`
	Ttl    int64  `yaml:"ttl"`
}

type NoticeKeyCollect struct {
	Task string `yaml:"task"`
	Mark string `yaml:"mark"`
}

// InventoryConfig 三方库存接口配置
type InventoryConfig struct {
	Url      string `yaml:"url"`
	DclUrl   string `yaml:"dcl_url"`
	ClientId string `yaml:"X-IBM-Client-Id"`
}

// InventorySwConfig  Sw三方库存接口配置
type InventorySwConfig struct {
	Url    string `yaml:"url"`
	PLocNo string `yaml:"p_loc_no"`
	Key    string `yaml:"key"`
}

// VchangyiEmailConfig  发送邮件配置
type VchangyiEmailConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// BaozunConfig PIM/DAM商品图片接口配置
type BaozunConfig struct {
	Domain       string `yaml:"domain"`
	ClientId     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
}

// BaiQiuConfig  百秋SW商品图片接口配置
type BaiQiuConfig struct {
	BaseURL    string `yaml:"base_url"`
	Secret     string `yaml:"secret"`
	CustomerId string `yaml:"customer_id"`
}

type CoachSaTag struct {
	TagId int `yaml:"tag_id"`
}

// CoachKOSConfig COACH KOS帖子打标接口配置
type CoachKOSConfig struct {
	BaseURL  string `yaml:"base_url"`
	Token    string `yaml:"token"`
	PageSize int    `yaml:"page_size"`
}

// TalkieConfig Talkie API配置
type TalkieConfig struct {
	BaseURL   string `yaml:"base_url"`
	ApiKey    string `yaml:"api_key"`
	ApiSecret string `yaml:"api_secret"`
}
