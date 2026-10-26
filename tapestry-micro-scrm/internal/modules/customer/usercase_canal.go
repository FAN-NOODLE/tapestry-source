package customer

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"micro-scrm/internal/pkg/helper"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/fx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"gitlab-ce.k8s.tools.vchangyi.com/go/canalgo/binlog"
	"gitlab-ce.k8s.tools.vchangyi.com/go/canalgo/canal"
	"gitlab-ce.k8s.tools.vchangyi.com/go/canalgo/replication"
	"micro-scrm/internal/config"
	"micro-scrm/internal/modules/customer/respository/redis"
	canal2 "micro-scrm/internal/pkg/canal"
)

var (
	count       = 0
	nowPos      binlog.Position
	nowPosMutex sync.Mutex
)
var indexExit = struct {
	sync.RWMutex
	m map[string]int
}{m: make(map[string]int)}

var docMap = struct {
	sync.RWMutex
	m map[string]map[string]string
}{m: make(map[string]map[string]string)}

type CanalServer struct {
	EsRepo EsRepository
	Repo   Repository
}

type MyEventHandler struct {
	canal.DummyEventHandler
	esRepo EsRepository
}

type WxUnion struct {
	Pid     string `json:"pid"`
	UnionId string `json:"union_id"`
}

type WxThirdCommon struct {
	AppID     string      `json:"app_id"`
	ThirdUser ThirdUserID `json:"third_user"`
}

//监听数据记录
func (h *MyEventHandler) OnRow(ev *canal.RowsEvent) error {
	defer func() {
		if err := recover(); err != nil {
			st := string(debug.Stack())
			log.Error(ctx.New()).Msgf("cancal onrow panic %v,stack %s", err, st)
		}
	}()
	if ev.Table.Name != "enterprise_user" {
		return nil
	}
	if !strings.HasPrefix(ev.Table.Schema, "micro_scrm_enterprise_") {
		return nil
	}

	esIndex := ev.Table.Schema + "." + ev.Table.Name //获取es索引
	prefix := config.CfgData.Env + "."
	esIndex = prefix + esIndex

	for cr := 0; cr < len(ev.Rows); cr++ {
		if "update" == ev.Action && cr%2 == 0 {
			continue
		}
		upDateNeedOpEs := false
		eueMap := make(map[string]interface{})
		for columnIndex, currColumn := range ev.Table.Columns {
			//字段名，字段的索引顺序，字段对应的值
			rowval := ev.Rows[cr][columnIndex]
			//需要判断是否有业务字段有更新，没有则不执行同步
			if "update" == ev.Action && upDateNeedOpEs == false && helper.MarshalToString(ev.Rows[cr][columnIndex]) != helper.MarshalToString(ev.Rows[cr-1][columnIndex]) {
				upDateNeedOpEs = true
			}

			eueMap[currColumn.Name] = rowval
			switch currColumn.Name {
			case "birthday":
				if _, ok := rowval.([]string); !ok {
					delete(eueMap, currColumn.Name)
				} else {
					if _, err := time.Parse("2006-01-02", rowval.(string)); err != nil {
						delete(eueMap, currColumn.Name)
					}
				}
			case "remark":
				if _, ok := rowval.([]byte); ok {
					eueMap[currColumn.Name] = string(rowval.([]byte))
				} else {
					eueMap[currColumn.Name] = ""
				}

			case "wx_union":
				if _, ok := rowval.([]byte); !ok {
					rowval = []byte("{}")
				}

				wxUnionList := []WxUnion{}
				eueMap[currColumn.Name] = wxUnionList
				item := make(map[string]string)
				err := json.Unmarshal(rowval.([]byte), &item)
				if err == nil {
					for k, v := range item {
						wx := WxUnion{
							Pid:     k,
							UnionId: v,
						}
						wxUnionList = append(wxUnionList, wx)
					}
					eueMap[currColumn.Name] = wxUnionList
				}
			case "wx_shop", "wx_open", "wx_external":
				if _, ok := rowval.([]byte); !ok {
					rowval = []byte("{}")
				}
				wtcs := []WxThirdCommon{}
				eueMap[currColumn.Name] = wtcs
				item := make(map[string]ThirdUserID)
				err := json.Unmarshal(rowval.([]byte), &item)
				if err == nil {
					for k, v := range item {
						wx := WxThirdCommon{
							AppID:     k,
							ThirdUser: v,
						}
						wtcs = append(wtcs, wx)
					}
					eueMap[currColumn.Name] = wtcs
				}
			case "groups", "tags":
				if _, ok := rowval.([]byte); !ok {
					rowval = []byte("[]")
				}
				item := []int64{}
				eueMap[currColumn.Name] = item
				err := json.Unmarshal(rowval.([]byte), &item)
				if err == nil && len(item) > 0 {
					eueMap[currColumn.Name] = item
				}
			}
		}

		if "update" == ev.Action && !upDateNeedOpEs {
			continue
		}
		eueMap["es_update"] = helper.GetFormatDateByParams(0)
		eueJson, err := json.Marshal(eueMap)
		if err != nil {
			log.Error(ctx.New()).Msgf("canal eueJson err %v", err)
			return err
		}
		euestr := string(eueJson)

		if _, ok := eueMap["id"].(int64); !ok {
			log.Error(ctx.New()).Msg("canal id not int64 ")
			continue
		}
		id := fmt.Sprint(eueMap["id"].(int64))
		if _, ok := eueMap["bind_to"].(int64); !ok {
			eueMap["bind_to"] = 0
		}
		bindTo := fmt.Sprint(eueMap["bind_to"])
		if ev.Action == "delete" || bindTo != "0" {
			if err = h.esRepo.DeleteDoc(ctx.New(), esIndex, id); err != nil {
				log.Error(ctx.New()).Msgf("canal OnRow elastic.DeleteDoc -> esIndex : %s, id:%s ,ev.Action:%v, bindTo:%v err:%v", esIndex, id, ev.Action, bindTo, err)
			}
		} else {
			docMap.Lock()
			count++
			if _, ok := docMap.m[esIndex]; !ok {
				docMap.m[esIndex] = make(map[string]string)
			}
			docMap.m[esIndex][id] = euestr
			docMap.Unlock()
		}
	}
	return nil
}

//创建、更改、重命名或删除表时触发，通常会需要清除与表相关的数据，如缓存。It will be called before OnDDL.
func (h *MyEventHandler) OnTableChanged(schema string, table string) error {
	//库，表
	record := fmt.Sprintf("OnTableChanged -> %s %s \n", schema, table)
	log.Info(ctx.New()).Msg(record)
	return nil
}

//监听binlog日志的变化文件与记录的位置
func (h *MyEventHandler) OnPosSynced(pos binlog.Position, set binlog.GTIDSet, force bool) error {
	//源码：当force为true，立即同步位置
	/*record := fmt.Sprintf("OnPosSynced %v %v ,set = %v, force = %v \n", pos.Name, pos.Pos, set, force)
	log.Info(ctx.New()).Msg(record)*/
	nowPosMutex.Lock()
	nowPos = pos
	nowPosMutex.Unlock()
	return nil
}

//当产生新的binlog日志后触发(在达到内存的使用限制后（默认为 1GB），会开启另一个文件，每个新文件的名称后都会有一个增量。)
func (h *MyEventHandler) OnRotate(r *replication.RotateEvent) error {
	//binlog的记录位置，新binlog的文件名
	record := fmt.Sprintf("canal On Rotate %v %v \n", r.Position, r.NextLogName)
	log.Info(ctx.New()).Msg(record)
	return nil

}

// create alter drop truncate(删除当前表再新建一个一模一样的表结构)
func (h *MyEventHandler) OnDDL(nextPos binlog.Position, queryEvent *replication.QueryEvent) error {
	//binlog日志的变化文件与记录的位置
	//record := fmt.Sprintf("canal OnDDL -> %v %v\n", nextPos.Name, nextPos.Pos)
	//query_event := fmt.Sprintf("%v\n %v\n %v\n %v\n %v\n",
	//	queryEvent.ExecutionTime,         //猜是执行时间，但测试显示0
	//	string(queryEvent.Schema),        //库名
	//	string(queryEvent.Query),         //变更的sql语句
	//	string(queryEvent.StatusVars[:]), //测试显示乱码
	//	queryEvent.SlaveProxyID)          //从库代理ID？
	//log.Info(ctx.New()).Msg(record + " " + query_event)
	return nil
}

func (cs *CanalServer) NewCanal() {
	c := canal2.ShareCanal
	if c == nil {
		err := canal2.InitCanal()
		if err != nil {
			log.Error(ctx.New()).Msgf("canal fail to NewCanal,err %v", err)
		}
		c = canal2.ShareCanal
	}
	if c == nil {
		log.Error(ctx.New()).Msg("canal fail to NewCanal ")
		return
	}

	go fx.RunSafe(func() error {
		for {
			time.Sleep(1000 * time.Millisecond)
			esupdate := time.Now().Format("2006-01-02 15:04:05")
			var body strings.Builder
			docMap.RLock()
			if len(docMap.m) > 0 {
				for index, dmap := range docMap.m {

					//判断索引是否已经存在，如果不存在，则创建
					_, ok := indexExit.m[index]
					if !ok {
						if !cs.EsRepo.ExistsIndex(ctx.New(), index) {
							if CreateIndex(cs.EsRepo, index) {
								indexExit.m[index] = 1
							}
						} else {
							indexExit.m[index] = 1
						}
					}

					if len(dmap) > 0 {
						//log.Info(ctx.New()).Msgf("canal esupdate = %s,count = %d index = %s,%d", esupdate, count, index, len(dmap))
						for id, doc := range dmap {
							body.WriteString(`{"index" : { "_index" : "` + index + `", "_type" : "_doc", "_id" : "` + id + `" }}`)
							body.WriteString("\n" + doc + "\n")
						}
						mapres, err := cs.EsRepo.AddDocsByBulk(ctx.New(), body.String())
						//log.Info(ctx.New()).Interface("dmap", dmap).Msg("AddDocsByBulk")
						if err != nil {
							log.Error(ctx.New()).Msg("canal fail cs.EsRepo.AddDocsByBulk error  " + err.Error())
						}
						if len(mapres) > 0 {
							for k, v := range mapres {
								if len(v) > 0 {
									idList := make(map[string][]string)
									for m, _ := range v {
										delete(docMap.m[k], m)
										idList[k] = append(idList[k], m)
									}
									if len(idList[k]) > 0 {
										ids := idList[k]
										prefix := config.CfgData.Env + "."
										kindex := strings.TrimLeft(k, prefix)
										go func(kindex, esupdate string, ids []string) {
											defer func() {
												if err := recover(); err != nil {
													st := string(debug.Stack())
													log.Error(ctx.New()).Msgf("canal fail onrow panic UpdateEsResult %v,stack %s", err, st)
												}
											}()
										}(kindex, esupdate, ids)
									}
								}
							}
						}
					}
				}
			}
			docMap.RUnlock()

			nowPosMutex.Lock()
			tempNowPos := nowPos
			nowPosMutex.Unlock()
			if tempNowPos.Pos > 0 && tempNowPos.String() != "" {
				opsByts, err := json.Marshal(tempNowPos)
				if err != nil {
					log.Error(ctx.New()).Msgf("canal OnPosSynced fail Marshal(pos) %v", err)
				} else {
					if _, err := redis.SetCanalPos(ctx.New(), string(opsByts)); err != nil {
						log.Error(ctx.New()).Msgf("canal now pos %s,%v,%s", string(opsByts), err, time.Now().Format("2006-01-02 15:04:05"))
					}
				}
			}

		}
	})

	c.SetEventHandler(&MyEventHandler{
		DummyEventHandler: canal.DummyEventHandler{},
		esRepo:            cs.EsRepo,
	})

	mpos, err := c.GetMasterPos()
	if err != nil {
		log.Error(ctx.New()).Msgf("canal fail to GetMasterPos %v ", err)
		return
	}
	canalpos, err := redis.GetCanalPos(ctx.New())
	if canalpos == "" {
		mposByte, _ := json.Marshal(mpos)
		_, err = redis.SetCanalPos(ctx.New(), string(mposByte))
		if err != nil {
			log.Error(ctx.New()).Msgf("canal fail to SetCanalPos %v ", err)
		}
	}
	if len(canalpos) > 0 {
		cpos := binlog.Position{}
		err := json.Unmarshal([]byte(canalpos), &cpos)
		if err != nil {
			log.Error(ctx.New()).Msgf("canal fail to Unmarshal to Position %v", err)
		} else {
			mpos.Name = cpos.Name
			mpos.Pos = cpos.Pos
		}
	}

	err = c.RunFrom(mpos)
	if err != nil {
		//redis.DelCanPos(ctx.New())
		log.Error(ctx.New()).Msgf("canal fail to RunFrom  %v,mpos = %v", err, mpos)
	} else {
		log.Info(ctx.New()).Msgf("canal start to RunFrom %v", mpos)
	}
}

func CreateIndex(esRepo EsRepository, indexName string) bool {
	body := `{
	"settings": {
		"index": {
			"max_result_window": 10000000
		}
	},
	"mappings": {
		"properties": {
			"avatar": {
				"type": "text",
				"fields": {
					"keyword": {
						"type": "keyword",
						"ignore_above": 256
					}
				}
			},
			"bind_to": {
				"type": "long"
			},
			"level_id": {
				"type": "long"
			},
			"birthday": {
				"type": "date"
			},
			"union_code": {
				"type": "text",
				"fields": {
					"keyword": {
						"type": "keyword",
						"ignore_above": 256
					}
				}
			},
			"point": {
				"type": "text"
			},
			"grade": {
				"type": "text",
				"fields": {
					"keyword": {
						"type": "keyword",
						"ignore_above": 256
					}
				}
			},
			"ba_code": {
				"type": "text"
			},
			"store_code": {
				"type": "text"
			},
			"es_update": {
				"type": "text"
			},
			"dt_update": {
				"type": "text"
			},
			"city": {
				"type": "text",
				"fields": {
					"keyword": {
						"type": "keyword",
						"ignore_above": 256
					}
				}
			},
			"country": {
				"type": "text",
				"fields": {
					"keyword": {
						"type": "keyword",
						"ignore_above": 256
					}
				}
			},
			"created_at": {
				"type": "long"
			},
			"groups": {
				"type": "long"
			},
			"id": {
				"type": "long"
			},
			"language": {
				"type": "text",
				"fields": {
					"keyword": {
						"type": "keyword",
						"ignore_above": 256
					}
				}
			},
			"latest_market_time": {
				"type": "long"
			},
			"market_count": {
				"type": "long"
			},
			"mobile": {
				"type": "text",
				"fields": {
					"keyword": {
						"type": "keyword",
						"ignore_above": 256
					}
				}
			},
			"delivery_mobile": {
				"type": "text",
				"fields": {
					"keyword": {
						"type": "keyword",
						"ignore_above": 256
					}
				}
			},
			"name": {
				"type": "text",
				"fields": {
					"keyword": {
						"type": "keyword",
						"ignore_above": 256
					}
				}
			},
			"nick": {
				"type": "text",
				"fields": {
					"keyword": {
						"type": "keyword",
						"ignore_above": 256
					}
				}

			},
			"pay_order": {
				"type": "long"
			},
			"per_customer_transaction": {
				"type": "long"
			},
			"province": {
				"type": "text",
				"fields": {
					"keyword": {
						"type": "keyword",
						"ignore_above": 256
					}
				}
			},
			"refund_fee": {
				"type": "long"
			},
			"refund_order": {
				"type": "long"
			},
			"remark": {
				"type": "text",
				"fields": {
					"keyword": {
						"type": "keyword",
						"ignore_above": 256
					}
				}
			},
			"sex": {
				"type": "long"
			},
			"tags": {
				"type": "long"
			},
			"total_fee": {
				"type": "long"
			},
			"updated_at": {
				"type": "long"
			},
			"version": {
				"type": "long"
			},
			"wx_external": {
				"properties": {
					"app_id": {
						"type": "text",
						"fields": {
							"keyword": {
								"type": "keyword",
								"ignore_above": 256
							}
						}
					},
					"third_user": {
						"properties": {
							"id": {
								"type": "long"
							},
							"third_id": {
								"type": "text",
								"fields": {
									"keyword": {
										"type": "keyword",
										"ignore_above": 256
									}
								}
							}
						}
					}
				}
			},
			"wx_open": {
				"properties": {
					"app_id": {
						"type": "text",
						"fields": {
							"keyword": {
								"type": "keyword",
								"ignore_above": 256
							}
						}
					},
					"third_user": {
						"properties": {
							"id": {
								"type": "long"
							},
							"third_id": {
								"type": "text",
								"fields": {
									"keyword": {
										"type": "keyword",
										"ignore_above": 256
									}
								}
							}
						}
					}
				}
			},
			"wx_shop": {
				"properties": {
					"app_id": {
						"type": "text",
						"fields": {
							"keyword": {
								"type": "keyword",
								"ignore_above": 256
							}
						}
					},
					"third_user": {
						"properties": {
							"id": {
								"type": "long"
							},
							"third_id": {
								"type": "text",
								"fields": {
									"keyword": {
										"type": "keyword",
										"ignore_above": 256
									}
								}
							}
						}
					}
				}
			},
			"wx_union": {
				"properties": {
					"pid": {
						"type": "text",
						"fields": {
							"keyword": {
								"type": "keyword",
								"ignore_above": 256
							}
						}
					},
					"union_id": {
						"type": "text",
						"fields": {
							"keyword": {
								"type": "keyword",
								"ignore_above": 256
							}
						}
					}
				}
			}
		}
	}
}`
	err := esRepo.CreateIndex(ctx.New(), indexName, body)
	if err != nil {
		return false
	}
	return true

}
