package main

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/helper"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/servers"
	"micro-scrm/cmd"
	"micro-scrm/internal/app"
	queueBinlog "micro-scrm/internal/app/queue/binlog"
	"micro-scrm/internal/config"
	"micro-scrm/internal/pkg/db"
	pkgKafka "micro-scrm/internal/pkg/kafka"
	"net/http"
	_ "net/http/pprof"
	"os"
)

// Version ...
var Version = "dev"

// Build ...
var Build = "now"

func main() {
	// 解析命令行参数
	cmd.ParseConfigFromCmd(Version, Build, os.Args)
	// 注册配置
	cmd.RegisterConfig(true)

	// 启动pprof
	go func() {
		app.Logger().Debug(ctx.New()).Msgf("pprof err: %v", http.ListenAndServe(config.CfgData.Monitor.Addr, nil))
	}()

	helper.Must(pkgKafka.InitAppKafka(config.CfgData.Kafka))

	c := ctx.New()

	var serverGroup pkgKafka.ServerGroup

	cusServers := make([]servers.Server, 0)

	// 消费binlog日志micro_scrm
	microCrmProduce := queueBinlog.NewMicroCrmProduce()
	microCrmProduceQueue := pkgKafka.NewGroupServer(c, microCrmProduce.GroupID(), microCrmProduce.Run, microCrmProduce.Topics(), 1, 1, false)

	cusServers = append(cusServers, microCrmProduceQueue)

	for _, s := range cusServers {
		serverGroup.Register(s)
	}

	groups := servers.Group(serverGroup.Servers()...)
	groups.ListenAndServe(func(context ctx.Context) {
		_ = db.DB().Close()
		microCrmProduce.CronChan <- true
	})
}
