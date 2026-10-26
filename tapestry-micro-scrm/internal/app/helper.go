package app

import (
	"github.com/go-redis/redis/v8"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/clients"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
	"micro-scrm/internal/config"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/logger"
	pkgRedis "micro-scrm/internal/pkg/redis"
	"micro-scrm/internal/pkg/uid"
)

// DB micro_scrm实例
func DB() *gorm.DB {
	return gormx.GetInstance()
}

// Redis redis 实例
func Redis() *redis.Client {
	return pkgRedis.Client()
}

// UID 分布式Id
func UID() int64 {
	return uid.ChannelId.GetChannelId()
}

// UIDContext 带上下文的分布式Id
func UIDContext(c ctx.Context) int64 {
	return uid.XID.ID(c)
}

// Logger 全局日志
func Logger() *logger.Logger {
	return logger.Log
}

// Client http客户端
func Client(options ...clients.HTTPClientOption) *clients.HTTPClient {
	return clients.NewHTTPClient(options...)
}

// Cfg 全局配置
func Cfg() *config.Config {
	return config.CfgData
}
