package config

// Version 版本号
var Version = "1.0.0"

// Build ...
var Build = "now"

// DefaultConfigFile ...
const DefaultConfigFile = "app.yaml"

// CfgData ...
var CfgData *Config

// 初始化配置
func init() {
	// 生产默认配置选项
	CfgData = &Config{}
}

// Init 业务数据初始化，在解析命令行参数过后执行
func Init() error {

	return nil
}
