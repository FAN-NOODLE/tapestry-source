package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/args"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/helper"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/monitor"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/servers"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/uid"
	"gw-work-scrm/cmd/http/api"
	"gw-work-scrm/internal/config"
	"gw-work-scrm/internal/pkg/redis"
	"net/http"
	"os"
	"runtime"
)

// Version ...
var Version = "dev"

// Build ...
var Build = "now"

// parseConfigFromCmd
func parseConfigFromCmd(arguments []string) (err error) {
	config.Version = Version
	config.Build = Build
	app := args.New("micro-tracker", args.Version(config.Version),
		args.Store(config.CfgData),
		args.Usage("micro-tracker http server cmd line."),
		args.FileConfigEnabled("config", config.DefaultConfigFile, false,
			"config file, support YAML/JSON/TOML. example: --config=tracker.yaml"),
		args.HelpExit(0),
	)

	err = app.Run(arguments)
	return err
}

func main() {
	// 解析命令行参数
	if err := parseConfigFromCmd(os.Args); err != nil {
		fmt.Println("error:", err.Error())
		os.Exit(1)
	}

	helper.Must(config.Init())
	helper.Must(log.Init(config.CfgData.Log()))

	helper.Must(os.Setenv("app", config.CfgData.App))
	helper.Panic(os.Setenv("qy_notice_key", config.CfgData.QYNoticeKey))

	uid.Init(uid.Config{
		URL:  config.CfgData.Xid,
		Mode: config.CfgData.Mode,
	})

	helper.Must(redis.Init(config.CfgData.Redis))

	// 启动 restful 服务
	serverGroup := servers.Group(
		&http.Server{Addr: config.CfgData.Restful.Addr, Handler: api.Api(config.CfgData.Mode)})
	if config.CfgData.Monitor.Enable {
		monitorServer := monitor.HTTP(config.CfgData.Restful.BasePath, config.CfgData.Monitor.Addr)
		serverGroup.Add(monitorServer)
	}

	fmt.Println(fmt.Sprintf("Server      Name:     %s", config.CfgData.App))
	fmt.Println(fmt.Sprintf("System      Name:     %s", runtime.GOOS))
	fmt.Println(fmt.Sprintf("Go          Version:  %s", runtime.Version()))
	fmt.Println(fmt.Sprintf("Gin         Version:  %s", gin.Version))
	fmt.Println(fmt.Sprintf("Listen      Address:  %s", "http://"+config.CfgData.Restful.Addr))

	serverGroup.ListenAndServe(func(c ctx.Context) {
		redis.Client().Close()
	})
}
