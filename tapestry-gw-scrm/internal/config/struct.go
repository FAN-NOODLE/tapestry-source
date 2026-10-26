package config

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gw-scrm/internal/pkg/redis"
	"time"
)

// Config 配置组合
type Config struct {
	App         string `json:"app"`
	Mode        string `json:"mode"`
	Env         string `json:"env"`
	Key         string `json:"key"`
	QYNoticeKey string `json:"qy_notice_key"`
	// log config
	Logger    LoggerCfg    `json:"logger"`
	Redis     redis.Config `json:"redis"`
	Restful   RestfulCfg   `json:"restful"`
	Monitor   MonitorCfg   `json:"monitor"`
	Service   ServiceCfg   `json:"service"`
	Path      PathCfg      `json:"path"`
	Storage   StorageCfg   `json:"storage"`
	Tencent   Tencent      `json:"tencent"`
	Chatrobot Chatrobot    `json:"chatrobot"`
}

// Log 返回日志配置对象
func (c Config) Log() log.Config {
	return log.Config{
		App:        c.App,
		Version:    Version,
		Build:      Build,
		Level:      c.Logger.Level,
		Format:     c.Logger.Format,
		TargetType: c.Logger.TargetType,
		Target:     c.Logger.Target,
	}
}

// LoggerCfg 日志配置
type LoggerCfg struct {
	Level      string `json:"level" usage:"日志级别"`
	Format     string `json:"format" usage:"日志输出格式【json/raw】"`
	TargetType string `json:"target_type"  usage:"输出目标类型【console/file/tcp/udp】"`
	Target     string `json:"target" usage:"日志输出格式【<file path>/<network address>】"`
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
	ShopToken      WorkCfg       `yaml:"shop-token"`
	Shop           ServiceConfig `yaml:"shop"`
	Boss           ServiceConfig `yaml:"boss"`
	Console        ServiceConfig `yaml:"console"`
	Open           ServiceConfig `yaml:"open"`
	MicroPay       ServiceConfig `yaml:"micro-pay"`
	MicroGoods     ServiceConfig `yaml:"micro-goods"`
	MicroActivity  ServiceConfig `yaml:"micro-activity"`
	MicroAuthority ServiceConfig `yaml:"micro-authority"`
	MicroUser      ServiceConfig `yaml:"micro-user"`
	GoMicroUser    ServiceConfig `yaml:"go-micro-user"`
	MicroScrm      ServiceConfig `yaml:"micro-scrm"`
	Cdp            ServiceConfig `yaml:"cdp"`
	UcCallBack     ServiceConfig `yaml:"uc-call-back"`
	WebH5          ServiceConfig `yaml:"web-h5"`
	ShortUrl       ServiceConfig `yaml:"short-url"`
	WxApiUrl       ServiceConfig `yaml:"wx-api-url"`
	DownloadCenter ServiceConfig `yaml:"download-center"`
}

// ServiceConfig SDK配置
type ServiceConfig struct {
	Domain string `yaml:"domain"`
	Secret string `yaml:"secret"`
}

// WorkCfg 企业微信配置
type WorkCfg struct {
	Domain string `yaml:"domain"`
	ID     string `yaml:"id"`
	Secret string `yaml:"secret"`
}

// PathCfg 路径配置
type PathCfg struct {
	Upload string `yaml:"upload"`
}

// 存储配置
type StorageCfg struct {
	TmpDir string       `yaml:"TmpDir"`
	COS    CloudStorage `yaml:"cos"`
	Blob   BlobStorage  `yaml:"blob"`
}

// 云存储结构体
type CloudStorage struct {
	Host          string `yaml:"host"`
	CdnHost       string `yaml:"cdnHost"`
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

// 腾讯云配置
type Tencent struct {
	Tms Tms `yaml:"tms"`
}

// 文本检测配置
type Tms struct {
	SecretId  string `yaml:"secretId"`
	SecretKey string `yaml:"secretKey"`
	Region    string `yaml:"region"`
}

type Chatrobot struct {
	Url  string `yaml:"url"`
	Time int    `yaml:"time"`
}
