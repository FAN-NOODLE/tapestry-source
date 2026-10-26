package binlog

import (
	"encoding/json"
	"fmt"
	"micro-scrm/internal/app"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/app/service"
	"micro-scrm/internal/config"
	"micro-scrm/internal/pkg/helper"
	pkgKafka "micro-scrm/internal/pkg/kafka"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/go/canalgo/binlog"
	"gitlab-ce.k8s.tools.vchangyi.com/go/canalgo/canal"
	"gitlab-ce.k8s.tools.vchangyi.com/go/canalgo/replication"
)

type KafkaEventHandler struct {
	service.Base
	canal.DummyEventHandler
	nowPos      binlog.Position
	nowPosMutex sync.Mutex

	redisKey string

	taskDao *microScrmDao.PlanTaskDao
}

// KafkaReq kafka消息体
type KafkaReq struct {
	Database string `json:"database"`
	Action   string `json:"action"`
	Data     any    `json:"data"`
}

func NewKafkaEventHandler() *KafkaEventHandler {
	return &KafkaEventHandler{
		redisKey: fmt.Sprintf("binlog-kafka:pos:%s:%s", config.CfgData.Env, helper.Md5V(config.CfgData.DB.MasterDSN)),
		taskDao:  microScrmDao.NewPlanTaskDao(),
	}
}

func (s *KafkaEventHandler) OnRotate(r *replication.RotateEvent) error {
	record := fmt.Sprintf("canal On Rotate %v %v \n", r.Position, r.NextLogName)
	s.Log().Info(ctx.New()).Msg(record)
	return nil
}

// OnTableChanged 创建、更改、重命名或删除表时触发，通常会需要清除与表相关的数据
func (s *KafkaEventHandler) OnTableChanged(schema string, table string) error {
	record := fmt.Sprintf("OnTableChanged -> %s %s \n", schema, table)
	s.Log().Info(ctx.New()).Msg(record)
	return nil
}

// OnDDL create alter drop truncate
func (s *KafkaEventHandler) OnDDL(nextPos binlog.Position, queryEvent *replication.QueryEvent) error {
	// binlog日志的变化文件与记录的位置
	record := fmt.Sprintf("canal OnDDL -> %v %v\n", nextPos.Name, nextPos.Pos)
	app.Logger().Info(ctx.New()).Msg(record)
	database := string(queryEvent.Schema)
	execQuery := string(queryEvent.Query)
	app.Logger().Info(ctx.New()).Msgf("canal OnDDL -> database [%s] query [%s]", database, execQuery)
	return nil
}

func (s *KafkaEventHandler) OnRow(ev *canal.RowsEvent) error {
	s.Log().Info(ctx.New()).Msgf("database [%s] table [%s]", ev.Table.Schema, ev.Table.Name)
	cc := ctx.New()
	defer func() {
		if err := recover(); err != nil {
			st := string(debug.Stack())
			s.Log().Error(cc).Msgf("OnRow panic %v,stack %s", err, st)
		}
	}()

	for i := 0; i < len(ev.Rows); i++ {
		if canal.UpdateAction == ev.Action && i%2 == 0 {
			continue
		}
		upNeed := false
		eueMap := make(map[string]interface{})

		for columnIndex, currColumn := range ev.Table.Columns {
			// 字段名，字段的索引顺序，字段对应的值
			rowVal := ev.Rows[i][columnIndex]
			eueMap[currColumn.Name] = rowVal
			// 需要判断是否有业务字段有更新，没有则不执行同步
			if canal.UpdateAction == ev.Action && upNeed == false && helper.MarshalToString(rowVal) != helper.MarshalToString(ev.Rows[i-1][columnIndex]) {
				upNeed = true
			}
		}

		if canal.UpdateAction == ev.Action && !upNeed {
			continue
		}
		if eueMap != nil {
			// 生产数据
			tableName := ev.Table.Name
			if strings.Contains(ev.Table.Name, "task_condition_data") { // 如果是task_condition_data表数据，判断任务是不是朋友圈子任务，是的话，则不进入kafka消费
				tableName = "task_condition_data"
				taskId := eueMap["task_id"].(int64)
				re := regexp.MustCompile(`_(\d+)$`) // 获取企业id
				matches := re.FindStringSubmatch(ev.Table.Schema)
				enterpriseId, _ := strconv.ParseInt(matches[1], 10, 64)
				taskIds := make([]int64, 0)
				isExist := s.Redis().Exists(cc, "task_moment_ids").Val()
				if isExist == 0 {
					startTime, _ := helper.YesterdayPeriodOfTimestamp(7)
					//过滤掉非人群包的朋友圈，非人群包的朋友圈基本都是全量人群，不走kafka
					list, err := s.taskDao.GetMany(cc, enterpriseId, "task_category = ? and task_start_time >= ? and crowd_id='' and deleted_at=0", microScrmModel.PLANTASKMOMENT, startTime)
					if err != nil {
						continue
					}
					for _, model := range list {
						taskIds = append(taskIds, model.ID)
					}
					s.Redis().Set(cc, "task_moment_ids", helper.MarshalToString(taskIds), time.Minute)
				} else {
					redisTaskIds := s.Redis().Get(cc, "task_moment_ids").Val()
					_ = json.Unmarshal([]byte(redisTaskIds), &taskIds)
				}
				if helper.InArray(taskId, taskIds) {
					continue
				}
			}
			s.Producer(cc, ev.Table.Schema, tableName, ev.Action, eueMap)
		}
	}

	return nil
}

// OnPosSynced 监听binlog日志的变化文件与记录的位置
func (s *KafkaEventHandler) OnPosSynced(pos binlog.Position, set binlog.GTIDSet, force bool) error {
	s.nowPosMutex.Lock()
	s.nowPos = pos
	s.nowPosMutex.Unlock()
	if s.nowPos.Pos > 0 && s.nowPos.String() != "" {
		opsBytes, err := json.Marshal(s.nowPos)
		if err != nil {
			s.Log().Error(ctx.New()).Msgf("canal OnPosSynced fail Marshal(pos) %v", err)
		} else {
			_, err = s.SetPos(ctx.New(), string(opsBytes))
			if err != nil {
				s.Log().Error(ctx.New()).Msgf("canal now pos %s,%v,%s", string(opsBytes), err, helper.GetNowDate())
			}
		}
	}
	return nil
}

func (s *KafkaEventHandler) SetPos(c ctx.Context, val string) (string, error) {
	return s.Redis().Set(c, s.redisKey, val, 24*time.Hour).Result()
}

func (s *KafkaEventHandler) GetPos(c ctx.Context) (string, error) {
	return s.Redis().Get(c, s.redisKey).Result()
}

func (s *KafkaEventHandler) Producer(c ctx.Context, database, tableName, action string, data any) {
	var request KafkaReq

	request.Database = database
	request.Action = action
	request.Data = data

	req := helper.MarshalToString(request)

	app.Logger().Info(c).Msgf("topic: [%s] data: [%s]", tableName, req)

	err := pkgKafka.GetAppKafka().Producer(req, tableName)
	if err != nil {
		s.Log().Error(c).Msgf("binlog Producer error: %v", err)
	}
}
