package db

import (
	"github.com/rs/zerolog/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/sqlplus"
	"time"
)

var db *sqlplus.ClusterDB

func DB() *sqlplus.ClusterDB {
	if db == nil {
		panic("db 模块没有初始化")
	}
	return db
}

//Init db 初始化
func Init(cfg Config) error {
	var opts []sqlplus.Option

	if len(cfg.Slaves.DSNList) > 0 {
		opts = append(opts, sqlplus.Slaves(cfg.Slaves.DSNList...))
		opts = append(opts, sqlplus.LoadStrategy(slavesLoadStrategy(cfg.Slaves)))
	}

	opts = append(opts, sqlplus.Pool(cfg.Pool.MaxIdleConns, cfg.Pool.MaxOpenConns, cfg.Pool.ConnMaxLifetime*time.Second))
	if cfg.TemplateDir != "" {
		log.Info().Msgf("load sql template from %s", cfg.TemplateDir)
		opts = append(opts, sqlplus.Template(sqlplus.ScanSQLTemplate(cfg.TemplateDir)))
	}

	if cfg.Log {
		opts = append(opts, sqlplus.Listeners(sqlplus.Log{}))
	}

	db = sqlplus.MustOpen(cfg.Driver, cfg.MasterDSN, opts...)
	return nil
}

//slavesLoadStrategy 负载均衡策略配置
func slavesLoadStrategy(cfg SlavesCfg) sqlplus.SlavesLoadStrategy {
	switch cfg.Strategy {
	case "weight":
		return sqlplus.WeightPollingStrategy(cfg.Weights)
	case "random":
		return &sqlplus.RandomStrategy{}
	default:
		return &sqlplus.PollingStrategy{}
	}
}

//PoolCfg 连接池配置
type PoolCfg struct {
	MaxIdleConns    int           `json:"max_idle_conns"`
	MaxOpenConns    int           `json:"max_open_conns"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime" usage:"单位：秒"`
}

//SlavesCfg 从库配置
type SlavesCfg struct {
	DSNList  []string `json:"dsn_list"`
	Strategy string   `json:"strategy"`
	Weights  []int    `json:"weights"`
}

//Config 主配置
type Config struct {
	Driver      string `json:"driver"`
	MasterDSN   string `json:"master_dsn"`
	TemplateDir string `json:"template_dir"`
	Log         bool   `json:"log"`

	Slaves SlavesCfg `json:"slaves"`
	Pool   PoolCfg   `json:"pool"`
}
