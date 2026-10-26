package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"runtime"
	"time"
)

var client *redis.Client

const Nil = redis.Nil

// Config redis 配置参数
type Config struct {
	Addr         string        `json:"addr"`
	User         string        `json:"user"`
	Pwd          string        `json:"pwd"`
	DB           int           `json:"db"`
	DialTimeout  time.Duration `json:"dial_timeout"`
	CmdTimeout   time.Duration `json:"cmd_timeout"`
	PoolSize     int           `json:"pool_size"`      // runtime.NumCPU() * 10
	MinIdleConns int           `json:"min_idle_conns"` // runtime.NumCPU()
}

// Init redis 初始化
func Init(cfg Config) error {
	if cfg.PoolSize < runtime.NumCPU() {
		cfg.PoolSize = runtime.NumCPU() * 10
	}
	if cfg.MinIdleConns < runtime.NumCPU() {
		cfg.MinIdleConns = runtime.NumCPU()
	}

	client = redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     cfg.Addr,
		Username: cfg.User,
		Password: cfg.Pwd,
		DB:       cfg.DB,
		OnConnect: func(c context.Context, cn *redis.Conn) error {
			log.Debug(ctx.Wrap(c)).Msg("新建redis链接")
			return nil
		},
		MaxRetries:      0,
		MinRetryBackoff: 8,
		MaxRetryBackoff: 512 * time.Millisecond,
		DialTimeout:     cfg.DialTimeout * time.Second,
		ReadTimeout:     cfg.CmdTimeout * time.Second,
		WriteTimeout:    cfg.CmdTimeout * time.Second,
		PoolSize:        cfg.PoolSize,
		MinIdleConns:    cfg.MinIdleConns,
	})

	_, err := client.Ping(ctx.New()).Result()
	return err
}

// Client 返回 redis 客户端
func Client() *redis.Client {
	return client
}
