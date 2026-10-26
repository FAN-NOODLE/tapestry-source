package main

import (
	_ "github.com/go-sql-driver/mysql"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/helper"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/lock"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/servers"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/uid"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/xmq"
	redismsgrepo "gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/xmq/redis"
	"micro-download/cmd"
	"micro-download/internal/config"
	"micro-download/internal/pkg/cos"
	"micro-download/internal/pkg/db"
	"micro-download/internal/pkg/redis"
	"micro-download/internal/queue"
	"micro-download/internal/queue/export"
	"os"
	"reflect"
)

var (
	Version = "prod"
	Build   = "now"
)

func main() {
	cmd.ParseConfigFromCmd(Version, Build, os.Args)
	helper.Must(config.Init())
	helper.Must(log.Init(config.CfgData.Log()))
	helper.Must(db.Init(config.CfgData.DB))
	helper.Must(cos.Init(config.CfgData.Cos))
	helper.Must(redis.Init(config.CfgData.Redis))
	helper.Must(lock.Init(lock.NewRedisLockerFactory(redis.Client())))
	uid.Init(uid.Config{Mode: "xid", URL: config.CfgData.Hosts.UID})
	con := ctx.New()
	mq := redismsgrepo.NewRedisMsgStatusTempRepository(redis.Client())
	redis.Client().Conn(con)
	broker := xmq.NewBroker(queue.NewRedisListener(redis.Client()), mq)

	// 队列注册
	// 注册文件生成队列
	broker.Register(queue.QueueFileGenerate, &export.SourceData{})

	// 服务启动
	serverGroup := servers.Group(xmq.NewServer(broker))
	serverGroup.ListenAndServe(func(c ctx.Context) {
		// 关闭数据库
		dbs := reflect.TypeOf(config.CfgData.DB)
		for i := 0; i < dbs.NumField(); i++ {
			if err := db.DB(dbs.Field(i).Name).Close(); err != nil {
				log.Err(c, err).Msgf("DB Close Error: %v", err)
			}
		}
	})
}
