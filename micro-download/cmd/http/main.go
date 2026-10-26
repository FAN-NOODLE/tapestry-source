package main

import (
	"github.com/gin-gonic/gin/binding"
	_ "github.com/go-sql-driver/mysql"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/helper"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/servers"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/uid"
	"micro-download/cmd"
	"micro-download/cmd/http/api"
	"micro-download/internal/config"
	"micro-download/internal/pkg/cos"
	"micro-download/internal/pkg/db"
	"micro-download/internal/pkg/redis"
	"net/http"
	"os"
	"reflect"
)

// Build 打包号
var Build = "test"

// Version 版本号
var Version = "now"

func init() {
	// 汉化参数验证器
	binding.Validator = ginplus.NewValidator()
}

func main() {
	cmd.ParseConfigFromCmd(Version, Build, os.Args)
	helper.Must(config.Init())
	helper.Must(log.Init(config.CfgData.Log()))
	helper.Must(db.Init(config.CfgData.DB))
	helper.Must(cos.Init(config.CfgData.Cos))
	helper.Must(redis.Init(config.CfgData.Redis))
	uid.Init(uid.Config{Mode: "xid", URL: config.CfgData.Hosts.UID})

	serverGroup := servers.Group(&http.Server{
		Addr:    config.CfgData.Restful.Addr,
		Handler: api.Api(config.CfgData.Mode),
	})

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
