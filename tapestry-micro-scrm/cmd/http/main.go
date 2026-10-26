package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/monitor"
	"micro-scrm/internal/app"
	"net/http"
	"os"
	"runtime"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/helper"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/servers"
	"micro-scrm/cmd"
	"micro-scrm/cmd/http/api"
	"micro-scrm/internal/config"
)

// Version ...
var Version = "dev"

// Build ...
var Build = "now"

func init() {
	// 汉化参数验证器
	binding.Validator = ginplus.NewValidator()
}

func main() {
	start := time.Now() // 获取当前时间

	// 解析命令行参数
	cmd.ParseConfigFromCmd(Version, Build, os.Args)
	// 注册配置
	cmd.RegisterConfig(false)

	serverGroup := servers.Group(&http.Server{
		Addr:    config.CfgData.Restful.Addr,
		Handler: api.Api(config.CfgData.Mode),
	})
	if config.CfgData.Monitor.Enable {
		serverGroup = serverGroup.Add(monitor.HTTP(config.CfgData.Restful.BasePath, config.CfgData.Monitor.Addr))
	}

	helper.Panic(os.Setenv("qy_notice_key", config.CfgData.QYNoticeKey))
	helper.Panic(os.Setenv("app", config.CfgData.App))

	elapsed := time.Since(start)

	app.Logger().Info(ctx.New()).Msgf("服务启动用时：%+v", elapsed)

	fmt.Println(fmt.Sprintf("Server      Name:     %s", config.CfgData.App))
	fmt.Println(fmt.Sprintf("System      Name:     %s", runtime.GOOS))
	fmt.Println(fmt.Sprintf("Go          Version:  %s", runtime.Version()))
	fmt.Println(fmt.Sprintf("Gin         Version:  %s", gin.Version))
	fmt.Println(fmt.Sprintf("Listen      Address:  %s", "http://"+config.CfgData.Restful.Addr))

	serverGroup.ListenAndServe(func(context ctx.Context) {
		cmd.ExitHandle()
	})
}
