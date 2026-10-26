package cmd

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/args"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/helper"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/lock"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/config"
	"micro-scrm/internal/events"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/pkg/canal"
	"micro-scrm/internal/pkg/db"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/elastic"
	pkgKafka "micro-scrm/internal/pkg/kafka"
	"micro-scrm/internal/pkg/logger"
	"micro-scrm/internal/pkg/redis"
	pkgUid "micro-scrm/internal/pkg/uid"
	"os"

	_ "github.com/ClickHouse/clickhouse-go"
	_ "github.com/go-sql-driver/mysql"
)

func ParseConfigFromCmd(v, b string, arguments []string) {
	config.Version = v
	config.Build = b
	appArgs := args.New("micro-tracker", args.Version(config.Version),
		args.Store(config.CfgData),
		args.FileConfigEnabled("config", config.DefaultConfigFile, false, ""),
		args.HelpExit(0),
	)

	err := appArgs.Run(arguments)
	// 解析命令行参数
	if err != nil {
		fmt.Println("args error:", err.Error())
		os.Exit(1)
	}
}

// RegisterConfig 加载配置
func RegisterConfig(logToStderr bool) {
	helper.Must(config.Init())
	helper.Must(log.Init(config.CfgData.Log(logToStderr)))
	logger.NewLogger()

	// redis
	helper.Must(redis.Init(config.CfgData.Redis))
	helper.Must(lock.Init(lock.NewRedisLockerFactory(redis.Client())))

	helper.Must(pkgKafka.InitAppKafka(config.CfgData.Kafka))
	helper.Must(elastic.NewElasticClient(config.CfgData.Elastic))

	helper.Must(db.InitDB(config.CfgData.DB))
	helper.Must(db.InitClickHouse(config.CfgData.ClickHouse))

	// gorm mysql初始化
	helper.Must(gormx.InitConnection(config.CfgData.DB))
	orm.Init()
	// 初始化事件
	events.EventInit()

	// gorm ck初始化
	helper.Must(gormx.InitCKGorm(config.CfgData.ClickHouse))

	// 分布式id
	pkgUid.InitUID()
	// 初始化全局分布式id管道生产
	pkgUid.NewChannel().GenChannelId()
}

// ExitHandle 处理系统退出信号
func ExitHandle() {
	if db.DB() != nil {
		_ = db.DB().Close()
	}
	if db.ClickHouse() != nil {
		_ = db.ClickHouse().Close()
	}
	if pkgKafka.GetAppKafka() != nil {
		_ = pkgKafka.GetAppKafka().Close()
	}
	canal.Close()
	// 初始化微信SDK证书等存储的redis
}
