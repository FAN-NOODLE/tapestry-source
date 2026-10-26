package queueBinlog

import (
	"encoding/json"
	"fmt"
	"github.com/jasonlvhit/gocron"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gorm.io/gorm"
	"micro-scrm/internal/app"
	svcBinlog "micro-scrm/internal/app/service/binlog"
	"micro-scrm/internal/pkg/db/gormx"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/notice"
	"strings"
	"time"
)

type MicroCrmProduce struct {
	groupId string
	topics  []string
	// 单个表提交时间hash
	timeKey string
	// 单个表提交数据list
	dataKey string
	// 单个表处理锁
	lockKey string
	// 定时任务管道
	CronChan chan bool
}

func NewMicroCrmProduce() *MicroCrmProduce {
	instance := &MicroCrmProduce{
		groupId: "micro_scrm",
		topics: []string{
			"members",
			"work_external_user_relation",
			"work_group_chat",
			"work_group_chat_member",
			"work_group_chat_member_log",
			"work_external_user_tag",
			"impression_tag_relation",

			"wx_user_corp_external",
			"reserve_order",
			"reserve_order_sales",
			"reserve_order_goods",
			"task_condition_data",
		},
		timeKey: "syncTable:time",
		dataKey: "syncTable:data:%s",
		lockKey: "syncTable:lock:%s",
	}

	s := gocron.NewScheduler()
	if err := s.Every(1).Minute().Do(instance.cronRun); err != nil {
		panic(err)
	}

	instance.CronChan = s.Start()

	return instance
}

func (q *MicroCrmProduce) Topics() []string {
	return q.topics
}

func (q *MicroCrmProduce) GroupID() string {
	return q.groupId
}

func (q *MicroCrmProduce) Run(c ctx.Context, topic string, value *string) error {
	var request svcBinlog.KafkaReq
	err := json.Unmarshal([]byte(*value), &request)
	if err != nil {
		app.Logger().Error(c).Msgf("[%s] 队列消费,转换结构错误:%+v", topic, err)
		return err
	}

	table := fmt.Sprintf("%s_ck.%s_local", request.Database, topic)
	// 存放table数据
	app.Redis().RPush(c, fmt.Sprintf(q.dataKey, table), helper.MarshalToString(request.Data))

	err = q.handle(c, table)
	if err != nil {
		app.Logger().Error(c).Msgf("[%s] 队列消费 错误:%+v", table, err)
		return err
	}
	return nil
}

// 定时执行
func (q *MicroCrmProduce) cronRun() {
	cc := ctx.New()
	// 取出redis里面table
	client := app.Redis().WithContext(cc)

	tableMap, err := client.HGetAll(cc, q.timeKey).Result()
	if tableMap == nil || err != nil {
		return
	}

	for table := range tableMap {
		if err = q.handle(cc, table); err != nil {
			app.Logger().Error(cc).Msgf("[%s] 定时执行,错误:%+v", table, err)
		}
	}
}

func (q *MicroCrmProduce) handle(c ctx.Context, table string) error {
	redisClient := app.Redis()

	// 判断redis是否存在表提交时间
	ok, err := redisClient.HExists(c, q.timeKey, table).Result()
	if !ok || err != nil {
		redisClient.HSet(c, q.timeKey, table, time.Now().Unix())
		return nil
	}

	tableCommitTime, err := redisClient.HGet(c, q.timeKey, table).Result()
	if err != nil {
		app.Logger().Error(c).Msgf("取[%s]表提交时间 错误:%+v", table, err)
		return err
	}

	tableCommitTimeInt := helper.StringToInt(tableCommitTime)
	// 将时间戳转换为time.Time类型
	t := time.Unix(int64(tableCommitTimeInt), 0)
	// 计算时间差值
	diff := time.Now().Sub(t)
	// 判断时间差值是否大于2分钟,大于2分钟执行一次
	if diff.Minutes() < 2 {
		app.Logger().Debug(c).Msgf("[%s] 插入时间差值小于2分钟", table)
		return nil
	}

	// 多进程消费防止针对一个表重复操作
	if !redisClient.SetNX(c, fmt.Sprintf(q.lockKey, table), time.Now().Unix(), time.Minute).Val() {
		app.Logger().Error(c).Msgf("[%s] 正在处理")
		return nil
	}

	// 获取redis中次表的数据长度
	listLen := redisClient.LLen(c, fmt.Sprintf(q.dataKey, table)).Val()
	if listLen == 0 {
		return nil
	}
	var mapData []map[string]interface{}
	// 遍历 list 中的每个元素，将其转换为 JSON 字符串
	for i := int64(0); i < listLen; i++ {
		var data map[string]interface{}
		val := redisClient.LPop(c, fmt.Sprintf(q.dataKey, table)).Val()
		_ = json.Unmarshal([]byte(val), &data)
		mapData = append(mapData, data)
	}

	// 去除重复的数据保留最新的一条
	unique := make(map[int64]map[string]interface{})
	for _, data := range mapData {
		id := int64(data["id"].(float64))
		unique[id] = data
	}
	var uniqueSlice []map[string]interface{}
	for _, data := range unique {
		// 将数据中的任何浮点数值转换为整数
		for key, val := range data {
			if floatVal, exits := val.(float64); exits {
				data[key] = int(floatVal)
			}
		}
		uniqueSlice = append(uniqueSlice, data)
	}
	// 先批量删除,后批量添加
	err = gormx.GetCkInstance().WithContext(c).Transaction(func(tx *gorm.DB) error {
		err = tx.Debug().Table(table).CreateInBatches(uniqueSlice, len(uniqueSlice)).Error
		if err != nil {
			return err
		}
		return nil
	})
	// 针对于task_condition_data表不合并
	if !strings.Contains(table, "task_condition_data") {
		reflectErr := gormx.GetCkInstance().WithContext(c).Debug().Exec(fmt.Sprintf("OPTIMIZE TABLE %s FINAL DEDUPLICATE", table)).Error
		if reflectErr != nil {
			notice.NewWechat().SendTextMessage(c, fmt.Sprintf("[%s] 合并,错误:%+v", table, reflectErr))
		}
	}
	if err != nil {
		q.rollbackList(c, table, uniqueSlice)
		app.Logger().Error(c).Msgf("[%s] 插入失败,错误:%+v", table, err)
		notice.NewWechat().SendTextMessage(c, fmt.Sprintf("[%s] 插入失败,错误:%+v", table, err))
		return err
	}

	// 更新表的插入时间
	redisClient.HSet(c, q.timeKey, table, time.Now().Unix())

	return nil
}

func (q *MicroCrmProduce) rollbackList(c ctx.Context, table string, data []map[string]interface{}) {
	redisClient := app.Redis()
	for _, val := range data {
		redisClient.LPush(c, fmt.Sprintf(q.dataKey, table), helper.MarshalToString(val))
	}
}
