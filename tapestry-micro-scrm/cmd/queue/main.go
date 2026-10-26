package main

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/helper"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/servers"
	"micro-scrm/cmd"
	"micro-scrm/internal/app"
	"micro-scrm/internal/app/queue/excitation"
	"micro-scrm/internal/app/queue/plan_task"
	"micro-scrm/internal/app/queue/store_visit"

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

	// 生产导购积分异步队列
	pointsProduce := excitation.NewPointsProduceQueue()
	cusServers = append(cusServers, pkgKafka.NewGroupServer(c, pointsProduce.GroupID(), pointsProduce.Run, pointsProduce.Topics(), 1, 1, false))

	// 导购金币过期异步队列
	coinsExpired := excitation.NewCoinsExpiredQueue()
	cusServers = append(cusServers, pkgKafka.NewGroupServer(c, coinsExpired.GroupID(), coinsExpired.Run, coinsExpired.Topics(), 1, 1, true))

	// 计划任务的人群包匹配队列
	planTaskImportProduce := plan_task.NewImportDataProduce()
	planTaskImportProduceQueue := pkgKafka.NewGroupServer(c, planTaskImportProduce.GroupID(), planTaskImportProduce.Run, planTaskImportProduce.Topics(), 1, 1, false)
	cusServers = append(cusServers, planTaskImportProduceQueue)

	// 批量打标队列（同步队列，避免并发打标）
	batchMarkingTag := plan_task.NewBatchMarkingTagQueue()
	batchMarkingTagQueue := pkgKafka.NewGroupServer(c, batchMarkingTag.GroupID(), batchMarkingTag.Run, batchMarkingTag.Topics(), 1, 1, false)
	cusServers = append(cusServers, batchMarkingTagQueue)

	// 消息队列处理企业群发朋友圈的人群包标签
	planTaskMomentDelTagProduce := plan_task.NewMomentDelTagProduce()
	planTaskMomentDelTagQueue := pkgKafka.NewGroupServer(c, planTaskMomentDelTagProduce.GroupID(), planTaskMomentDelTagProduce.Run, planTaskMomentDelTagProduce.Topics(), 1, 1, false)
	cusServers = append(cusServers, planTaskMomentDelTagQueue)
	// 消息队列处理推送MA任务的企业群发完成数据推送
	planMaFinishPushProduce := plan_task.NewMaPushFinishProduce()
	planMaFinishPushQueue := pkgKafka.NewGroupServer(c, planMaFinishPushProduce.GroupID(), planMaFinishPushProduce.Run, planMaFinishPushProduce.Topics(), 1, 1, false)
	cusServers = append(cusServers, planMaFinishPushQueue)
	storeVisitProduce := store_visit.NewCreateStoreVisitPDFProduce()
	storeVisitQueue := pkgKafka.NewGroupServer(c, storeVisitProduce.GroupID(), storeVisitProduce.Run, storeVisitProduce.Topics(), 1, 1, false)
	cusServers = append(cusServers, storeVisitQueue)
	for _, s := range cusServers {
		serverGroup.Register(s)
	}

	groups := servers.Group(serverGroup.Servers()...)
	groups.ListenAndServe(func(context ctx.Context) {
		_ = db.DB().Close()
	})
}
