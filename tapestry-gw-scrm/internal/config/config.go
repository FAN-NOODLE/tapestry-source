package config

var Version = "1.0.0"
var Build = "now"

const DefaultConfigFile = "app.yaml"

var CfgData *Config

func init() {
	// 生产默认配置选项
	CfgData = &Config{}
}

// Init 业务数据初始化，在解析命令行参数过后执行
func Init() error {

	return nil
}
