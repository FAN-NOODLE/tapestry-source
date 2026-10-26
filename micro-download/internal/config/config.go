package config

// Version 版本号
var Version = "1.0.0"

// Build 打包号
var Build string

// DefaultConfigFile 默认配置文件名称
const DefaultConfigFile = "app.yaml"

// CfgData 配置变量
var CfgData *Config

func init() {
	// 生产默认配置选项
	CfgData = &Config{}
}

// Init 业务数据初始化，在解析命令行参数过后执行
func Init() error {

	return nil
}
