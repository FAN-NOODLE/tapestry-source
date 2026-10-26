package cmd

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/args"
	"micro-download/internal/config"
	"os"
)

// ParseConfigFromCmd
func ParseConfigFromCmd(v, b string, arguments []string) {
	config.Version = v
	config.Build = b
	app := args.New("micro-tracker", args.Version(config.Version),
		args.Store(config.CfgData),
		args.Usage("micro-tracker http server cmd line."),
		args.FileConfigEnabled("config", config.DefaultConfigFile, false,
			"config file, support YAML/JSON/TOML. example: --config=tracker.yaml"),
		args.HelpExit(0),
	)

	err := app.Run(arguments)
	// 解析命令行参数
	if err != nil {
		fmt.Printf("\nError: ")
		fmt.Printf("\n  %s", err.Error())
		fmt.Printf("\n\n")
		os.Exit(1)
	}
}
