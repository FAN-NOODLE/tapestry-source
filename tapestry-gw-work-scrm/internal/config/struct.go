package config

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gw-work-scrm/internal/pkg/redis"
	"time"
)

// Config 配置组合
type Config struct {
	App         string `json:"app"`
	Mode        string `json:"mode"`
	QYNoticeKey string `json:"qy_notice_key"`
	// log config
	Logger              LoggerCfg    `json:"logger"`
	Redis               redis.Config `json:"redis"`
	Restful             RestfulCfg   `json:"restful"`
	Monitor             MonitorCfg   `json:"monitor"`
	Xid                 string       `json:"xid"`
	Service             ServiceCfg   `json:"service"`
	JWT                 JwtCfg       `json:"jwt"`
	WeLearningTokenSalt string       `json:"we_learning_token_salt"`
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
	MicroScrm      ServiceConfig `yaml:"micro-scrm"`
	Boss           ServiceConfig `yaml:"boss"`
	Work           ServiceConfig `yaml:"work"`
	MicroAuthority ServiceConfig `yaml:"micro-authority"`
}

// ServiceConfig SDK配置
type ServiceConfig struct {
	Domain string `yaml:"domain"`
}

// WorkCfg 企业微信配置
type WorkCfg struct {
	Domain string `yaml:"domain"`
	ID     string `yaml:"id"`
	Secret string `yaml:"secret"`
}

// JwtCfg JWT配置
type JwtCfg struct {
	Secret     string `yaml:"secret"`
	Ttl        int64  `yaml:"ttl"`
	RefreshTtl int64  `yaml:"refresh_ttl"`
}
