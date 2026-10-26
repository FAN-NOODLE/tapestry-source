package job

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/go/cy-sdk-go/exception"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/sqlplus"
	"micro-scrm/internal/config"
	"micro-scrm/internal/modules/enterprise"
	enterpriseMysql "micro-scrm/internal/modules/enterprise/respository/mysql"
	"micro-scrm/internal/pkg/db"
	"micro-scrm/internal/pkg/xjob"
)

type Table struct {
	DBdsn    *mysql.Config
	Clolumns [][]string
	TBname   string
}

func tables(c ctx.Context) ([]*Table, error) {

	var tblist []*Table
	dbdsn, err := mysql.ParseDSN(config.CfgData.DB.MasterDSN)
	if err != nil {
		exception.SendError(err, map[string]interface{}{"dsn": config.CfgData.DB.MasterDSN})
		return nil, err
	}
	enterpriseService := enterprise.New(enterpriseMysql.New())
	_, enterprises, err := enterpriseService.ListEnterprise(c, nil, 1, 0)
	if err != nil {
		return nil, err
	}
	// 外部联系人跟员工的关系表
	tblist = append(tblist, &Table{
		TBname: "work_external_user_relation",
		DBdsn:  &mysql.Config{Addr: dbdsn.Addr, DBName: "micro_scrm", User: dbdsn.User, Passwd: dbdsn.Passwd},
		Clolumns: [][]string{
			{"id", "Int64"},
			{"enterprise_id", "Int64"},
			{"work_external_user_id", "Int64"},
			{"member_id", "Int64"},
			{"remark", "String"},
			{"description", "String"},
			{"createtime", "Int64"},
			{"remark_corp_name", "String"},
			{"remark_mobiles", "String"},
			{"add_way", "Int64"},
			{"state", "String"},
			{"oper_userid", "String"},
			{"relation_status", "Int8"},
			{"member_deleted_time", "Int64"},
			{"external_deleted_time", "Int64"},
			{"created_at", "Int64"},
			{"updated_at", "Int64"},
			{"deleted_at", "Int64"},
		},
	})
	// 外部联系人跟员工的关系表
	tblist = append(tblist, &Table{
		TBname: "work_group_chat_member",
		DBdsn:  &mysql.Config{Addr: dbdsn.Addr, DBName: "micro_scrm", User: dbdsn.User, Passwd: dbdsn.Passwd},
		Clolumns: [][]string{
			{"id", "Int64"},
			{"enterprise_id", "Int64"},
			{"group_chat_id", "Int64"},
			{"member_id", "Int64"},
			{"type", "Int8"},
			{"userid", "String"},
			{"unionid", "String"},
			{"status", "Int8"},
			{"join_time", "Int64"},
			{"leave_time", "Int64"},
			{"join_scene", "Int64"},
			{"nickname", "String"},
			{"name", "String"},
			{"created_at", "Int64"},
		},
	})

	for _, vItem := range *enterprises {
		tblist = append(tblist, &Table{
			TBname: "wx_user_corp_external",
			DBdsn:  &mysql.Config{Addr: dbdsn.Addr, DBName: fmt.Sprintf("micro_scrm_enterprise_%d", vItem.EnterpriseID), User: dbdsn.User, Passwd: dbdsn.Passwd},
			Clolumns: [][]string{
				{"id", "Int64"},
				{"appid", "String"},
				{"openid", "String"},
				{"unionid", "String"},
				{"nick", "String"},
				{"name", "String"},
				{"position", "String"},
				{"avatar", "String"},
				{"birthday", "String"},
				{"union_code", "String"},
				{"point", "String"},
				{"grade", "String"},
				{"store_code", "String"},
				{"ba_code", "String"},
				{"corp_name", "String"},
				{"corp_full_name", "String"},
				{"type", "Int8"},
				{"sex", "Int8"},
				{"euid", "Nullable(Int64)"},
				{"external_profile", "String"},
				{"created_at", "Int64"},
				{"updated_at", "Int64"},
				{"deleted_at", "Int64"},
			},
		})

		tblist = append(tblist, &Table{
			TBname: "campaign_tool_cron_data",
			DBdsn:  &mysql.Config{Addr: dbdsn.Addr, DBName: fmt.Sprintf("micro_scrm_enterprise_%d", vItem.EnterpriseID), User: dbdsn.User, Passwd: dbdsn.Passwd},
			Clolumns: [][]string{
				{"id", "Int64"},
				{"union_id", "String"},
				{"pin_pai_buy_last_time", "String"},
				{"online_buy_last_time", "String"},
				{"offline_buy_time", "String"},
				{"tian_mao_buy_last_time", "String"},
				{"jing_dong_buy_last_time", "String"},
				{"mini_buy_last_time", "String"},
				{"created_at", "Int64"},
			},
		})
	}
	return tblist, nil
}

func syncSocialGroupTable(c ctx.Context) error {

	tbs, err := tables(c)
	if err != nil {
		exception.SendError(err, map[string]interface{}{})
	}
	ckdb := db.ClickHouse()
	for _, item := range tbs {
		if err = initTable(c, ckdb, item.DBdsn, item.TBname, item.TBname, item.Clolumns); err != nil {
			log.Error(c).Msgf("同步ck 失败 table: %s,err:%v", item.TBname, err)
		}
	}
	return nil
}

func initTable(c ctx.Context, ckdb *sqlplus.ClusterDB, dbdsn *mysql.Config, table string, source string, columns [][]string) error {
	var ckTable = fmt.Sprintf("%s.%s", dbdsn.DBName, table)
	var ckTableLocal = fmt.Sprintf("%s_local", ckTable)

	if _, err := ckdb.DB.Exec(c, fmt.Sprintf("create database if not exists %s on cluster %s", dbdsn.DBName, config.CfgData.ClickHouse.Cluster)); err != nil {
		exception.SendError(err, map[string]interface{}{"dsn": config.CfgData.ClickHouse.MasterDSN})
		return err
	}
	if _, err := ckdb.DB.Exec(c, fmt.Sprintf("drop table if exists %s on cluster %s", ckTable, config.CfgData.ClickHouse.Cluster)); err != nil {
		exception.SendError(err, map[string]interface{}{"dsn": config.CfgData.ClickHouse.MasterDSN})
		return err
	}
	if _, err := ckdb.DB.Exec(c, fmt.Sprintf("drop table if exists %s on cluster %s", ckTableLocal, config.CfgData.ClickHouse.Cluster)); err != nil {
		exception.SendError(err, map[string]interface{}{"dsn": config.CfgData.ClickHouse.MasterDSN})
		return err
	}

	var columnsWithType []string
	var columnsName []string
	for _, item := range columns {
		columnsWithType = append(columnsWithType, item[0]+" "+item[1])
		columnsName = append(columnsName, item[0])
	}

	if _, err := ckdb.DB.Exec(c, fmt.Sprintf(`create table %s on cluster %s (%s)
		ENGINE = ReplicatedMergeTree('/clickhouse/tables/%s/{shard}', '{replica}') ORDER BY %s SETTINGS index_granularity = 8192;`, ckTableLocal, config.CfgData.ClickHouse.Cluster, strings.Join(columnsWithType, ","), ckTable, columns[0][0])); err != nil {
		exception.SendError(err, map[string]interface{}{"dsn": config.CfgData.ClickHouse.MasterDSN})
		return err
	}
	query := fmt.Sprintf(`create table %s on cluster %s as %s ENGINE = Distributed('%s', '%s', '%s_local', %s);`, ckTable, config.CfgData.ClickHouse.Cluster, ckTableLocal, config.CfgData.ClickHouse.Cluster, dbdsn.DBName, table, columns[0][0])
	if _, err := ckdb.DB.Exec(c, query); err != nil {
		exception.SendError(err, map[string]interface{}{"dsn": config.CfgData.ClickHouse.MasterDSN, "query": query})
		return err
	}

	if _, err := ckdb.DB.Exec(c, fmt.Sprintf("insert into %s select %s from mysql('%s', '%s', '%s', '%s', '%s')", ckTable, strings.Join(columnsName, ","), dbdsn.Addr, dbdsn.DBName, source, dbdsn.User, dbdsn.Passwd)); err != nil {
		exception.SendError(err, map[string]interface{}{"dsn": dbdsn})
		return err
	}
	return nil
}

func insertGroupMemberLog(c ctx.Context, unixtime int64) error {
	log.Debug(c).Msgf("====================同步社群昨日日志=======================")
	ntime := time.Now().AddDate(0, 0, -1)
	ys := time.Date(ntime.Year(), ntime.Month(), ntime.Day(), 0, 0, 0, 0, ntime.Location()).Unix()
	if unixtime > 0 {
		ys = unixtime
	}

	ckdb := db.ClickHouse()
	inst := "insert into micro_scrm.group_chat_member_log select enterprise_id,group_chat_id,member_id,userid as user_id,unionid as union_id,if(status=2,1,1) as status,type,join_scene,join_time as opera_time,nickname from micro_scrm.work_group_chat_member where ((join_time>=? and join_time<=?) or (leave_time>=? and leave_time<=?))"
	if _, err := ckdb.DB.Exec(c, inst, ys, ys+86399, ys, ys+86399); err != nil {
		log.Error(c).Msgf("写入进群数据失败 err: %+v", err)
		return err
	}
	instdel := "insert into micro_scrm.group_chat_member_log select enterprise_id,group_chat_id,member_id,userid as user_id,unionid as union_id,status,type,join_scene,leave_time as opera_time,nickname from micro_scrm.work_group_chat_member where ((join_time>=? and join_time<=?) or (leave_time>=? and leave_time<=?)) and status=2"
	if _, err := ckdb.DB.Exec(c, instdel, ys, ys+86399, ys, ys+86399); err != nil {
		log.Error(c).Msgf("写入退群数据失败 err: %+v", err)
		return err
	}
	return nil
}

func SocialDrawJob(c ctx.Context, jobHandler *xjob.JobHandler) error {
	specsArgs := jobHandler.Runner().Specs()
	if len(specsArgs) == 0 {
		fmt.Println("参数有误，格式为--job.spec=subcmd,params")
		return nil
	}
	subcmd := specsArgs[0]
	switch subcmd {
	case "SyncSocialGroupTable":
		// 同步社群画像相关表
		return syncSocialGroupTable(c)
	case "InsertGroupMemberLog":
		// 插入昨日入群数据
		var ys int64
		if len(specsArgs) == 2 {
			ys, _ = strconv.ParseInt(specsArgs[1], 10, 64)
		}
		return insertGroupMemberLog(c, ys)
	default:
		return nil
	}
}
