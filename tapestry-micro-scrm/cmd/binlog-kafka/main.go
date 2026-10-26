package main

import (
	"encoding/json"
	_ "github.com/ClickHouse/clickhouse-go"
	_ "github.com/go-sql-driver/mysql"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/go/canalgo/binlog"
	"gitlab-ce.k8s.tools.vchangyi.com/go/canalgo/canal"
	"micro-scrm/cmd"
	appBinlog "micro-scrm/internal/app/service/binlog"
	"micro-scrm/internal/config"
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
	cmd.RegisterConfig(false)

	cfg := canal.NewDefaultConfig()
	cfg.DSN = config.CfgData.DB.MasterDSN

	// 过滤数据库,表
	cfg.IncludeTableRegex = append(cfg.IncludeTableRegex, "micro_scrm\\.excitation_medal$")
	cfg.IncludeTableRegex = append(cfg.IncludeTableRegex, "micro_scrm\\.members$")
	cfg.IncludeTableRegex = append(cfg.IncludeTableRegex, "micro_scrm\\.work_external_user_relation$")
	cfg.IncludeTableRegex = append(cfg.IncludeTableRegex, "micro_scrm\\.work_group_chat$")
	cfg.IncludeTableRegex = append(cfg.IncludeTableRegex, "micro_scrm\\.work_group_chat_member$")
	cfg.IncludeTableRegex = append(cfg.IncludeTableRegex, "micro_scrm\\.work_group_chat_member_log$")
	cfg.IncludeTableRegex = append(cfg.IncludeTableRegex, "micro_scrm\\.work_external_user_tag$")
	cfg.IncludeTableRegex = append(cfg.IncludeTableRegex, "micro_scrm\\.impression_tag_relation$")
	cfg.IncludeTableRegex = append(cfg.IncludeTableRegex, "micro_scrm_enterprise_.*\\.wx_user_corp_external$")
	cfg.IncludeTableRegex = append(cfg.IncludeTableRegex, "micro_scrm_enterprise_.*\\.reserve_order$")
	cfg.IncludeTableRegex = append(cfg.IncludeTableRegex, "micro_scrm_enterprise_.*\\.reserve_order_sales$")
	cfg.IncludeTableRegex = append(cfg.IncludeTableRegex, "micro_scrm_enterprise_.*\\.reserve_order_goods$")
	cfg.IncludeTableRegex = append(cfg.IncludeTableRegex, "micro_scrm_enterprise_.*\\.task_condition_data_\\d{16}$")

	cc := ctx.New()

	c, err := canal.NewCanal(cfg)
	if err != nil {
		log.Error(cc).Msgf("create canal err %v", err)
		return
	}

	// 初始化，设置事件处理器
	eventHandler := appBinlog.NewKafkaEventHandler()
	c.SetEventHandler(eventHandler)
	mPos, err := c.GetMasterPos()
	if err != nil {
		log.Error(cc).Msgf("canal fail to GetMasterPos %v ", err)
		return
	}

	// 获取上一个binlog节点
	canalPos, err := eventHandler.GetPos(cc)
	if canalPos == "" {
		mPosByte, _ := json.Marshal(mPos)
		_, err = eventHandler.SetPos(cc, string(mPosByte))
		if err != nil {
			log.Error(cc).Msgf("canal fail to SetCanalPos %v ", err)
		}
	}

	if len(canalPos) > 0 {
		cPos := binlog.Position{}
		err := json.Unmarshal([]byte(canalPos), &cPos)
		if err != nil {
			log.Error(cc).Msgf("canal fail to Unmarshal to Position %v", err)
		} else {
			mPos.Name = cPos.Name
			mPos.Pos = cPos.Pos
		}
	}

	err = c.RunFrom(mPos)
	if err != nil {
		log.Error(cc).Msgf("canal fail to RunFrom  %v,mPos = %v", err, mPos)
	} else {
		log.Info(cc).Msgf("canal start to RunFrom %v", mPos)
	}
}
