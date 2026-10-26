package main

import (
	"os"

	_ "github.com/ClickHouse/clickhouse-go"
	_ "github.com/go-sql-driver/mysql"
	"micro-scrm/cmd"
	"micro-scrm/internal/modules/customer"
	userEs "micro-scrm/internal/modules/customer/respository/elastic"
	"micro-scrm/internal/modules/customer/respository/mysql"
)

// Version ...
var Version = "dev"

// Build ...
var Build = "now"

func main() {
	// 解析命令行参数
	cmd.ParseConfigFromCmd(Version, Build, os.Args)
	// 注册配置
	cmd.RegisterConfig(false)
	//同步binlog至es
	cs := customer.CanalServer{EsRepo: userEs.NewEsRepo(), Repo: mysql.NewCustomer()}
	cs.NewCanal()
}
