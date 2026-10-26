package config

import (
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"

	"micro-download/internal/pkg/db"
	"micro-download/internal/pkg/redis"
)

// Config 配置组合
type Config struct {
	App   string  `json:"app"`
	Mode  string  `json:"mode"`
	Env   string  `json:"env"`
	Cos   QCos    `json:"cos"`
	Blob  QBlob   `json:"blob"`
	Hosts HostCfg `json:"hosts"`
	// log config
	Host       Host         `json:"host"`
	Logger     LoggerCfg    `json:"logger"`
	DB         db.Collect   `json:"db"`
	Redis      redis.Config `json:"redis"`
	Restful    RestfulCfg   `json:"restful"`
	Downloader Downloader   `json:"downloader"`
}

// HostCfg host配置
type HostCfg struct {
	Work      string `json:"work"`
	Boss      string `json:"boss"`
	Console   string `json:"console"`
	Open      string `json:"open"`
	Cdp       string `json:"cdp"`
	UID       string `json:"uid"`
	MicroScrm string `json:"microScrm"`
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

type Host struct {
	Work      string `json:"work"`
	Boss      string `json:"boss"`
	Console   string `json:"console"`
	Open      string `json:"open"`
	Cdp       string `json:"cdp"`
	UID       string `json:"uid"`
	MicroScrm string `json:"micro-scrm"`
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

type QCos struct {
	Region    string `yaml:"region"`
	Schema    string `yaml:"schema"`
	SecretId  string `yaml:"secret_id"`
	SecretKey string `yaml:"secret_key"`
	Bucket    string `yaml:"bucket"`
	HostUrl   string `yaml:"host_url"`
}

type QBlob struct {
	Domain      string `yaml:"domain"`
	Bucket      string `yaml:"bucket"`
	AccountName string `yaml:"accountName"`
	AccountKey  string `yaml:"accountKey"`
}

type Downloader struct {
	// ParallelNumber 数据获取并行数数
	ParallelNumber int    `yaml:"parallel_number"`
	TmpFilePath    string `yaml:"tmp_file_path"`
	Timeout        int    `yaml:"timeout"`
	RetrySleep     []int  `yaml:"retry_sleep"`
}
