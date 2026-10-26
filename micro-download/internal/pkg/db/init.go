package db

import (
	"reflect"
	"time"

	"github.com/rs/zerolog/log"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/sqlplus"
)

var (
	microScrmDB *sqlplus.ClusterDB
)

func DB(flag string) *sqlplus.ClusterDB {
	db := microScrmDB
	if flag == "microScrm" {
		db = microScrmDB
	}
	return db
}

// Init 初始化
func Init(cfg Collect) error {
	var opts []sqlplus.Option
	keys := reflect.TypeOf(cfg)
	values := reflect.ValueOf(cfg)
	for k := 0; k < keys.NumField(); k++ {
		cfgData := values.Field(k).Interface().(Config)
		if len(cfgData.Slaves.DSNList) > 0 {
			opts = append(opts, sqlplus.Slaves(cfgData.Slaves.DSNList...))
			opts = append(opts, sqlplus.LoadStrategy(slavesLoadStrategy(cfgData.Slaves)))
		}
		opts = append(opts, sqlplus.Pool(cfgData.Pool.MaxIdleConns, cfgData.Pool.MaxOpenConns, cfgData.Pool.ConnMaxLifetime*time.Second))
		if cfgData.TemplateDir != "" {
			log.Info().Msgf("load sql template from %s", cfgData.TemplateDir)
			opts = append(opts, sqlplus.Template(sqlplus.ScanSQLTemplate(cfgData.TemplateDir)))
		}
		if cfgData.Log {
			opts = append(opts, sqlplus.Listeners(sqlplus.Log{}))
		}
		name := keys.Field(k).Name
		if name == "MicroScrm" {
			microScrmDB = sqlplus.MustOpen(cfgData.Driver, cfgData.MasterDSN, opts...)
		}
	}
	return nil
}

// slavesLoadStrategy 负载均衡策略配置
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

// PoolCfg 连接池配置
type PoolCfg struct {
	MaxIdleConns    int           `yaml:"maxidleconns"`
	MaxOpenConns    int           `yaml:"maxopenconns"`
	ConnMaxLifetime time.Duration `yaml:"connmaxlifetime" usage:"单位：秒"`
}

// SlavesCfg 从库配置
type SlavesCfg struct {
	DSNList  []string `yaml:"dsnlist"`
	Strategy string   `yaml:"strategy"`
	Weights  []int    `yaml:"weights"`
}

// Config 主配置
type Config struct {
	Driver      string `yaml:"driver"`
	MasterDSN   string `yaml:"masterdsn"`
	TemplateDir string `yaml:"templatedir"`
	Log         bool   `yaml:"log"`

	Slaves SlavesCfg `yaml:"slaves"`
	Pool   PoolCfg   `yaml:"pool"`
}

// Collect 多数据库实例
type Collect struct {
	MicroScrm Config `yaml:"microScrm"`
}
