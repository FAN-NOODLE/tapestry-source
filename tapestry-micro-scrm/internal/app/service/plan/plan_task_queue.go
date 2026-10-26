package planService

import (
	"encoding/json"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	planReq "micro-scrm/internal/app/http/request/plan"
	pkgKafka "micro-scrm/internal/pkg/kafka"
)

type PlanTaskQueueSvc struct {
}

func NewPlanTaskQueueSvc() *PlanTaskQueueSvc {
	return &PlanTaskQueueSvc{}
}

// ImportDataProduce 任务人群包的导入，匹配人群包功能进入队列
func (p *PlanTaskQueueSvc) ImportDataProduce(c ctx.Context, params planReq.ImportDataQueueData) {
	str, _ := json.Marshal(params)
	err := pkgKafka.GetAppKafka().Producer(string(str), "import_data")
	if err != nil {
		log.Info(c).Msgf("任务人群包数据进入队列失败,错误原因：%+v", err)
	}

}

// TagDataProduce 标签群发任务类型，任务维度为标签的用户，需要先筛出用户然后再打标签
func (p *PlanTaskQueueSvc) TagDataProduce(c ctx.Context, data planReq.TagTaskQueueData) {
	str, _ := json.Marshal(data)
	err := pkgKafka.GetAppKafka().Producer(string(str), "tags_data")
	if err != nil {
		log.Info(c).Msgf("任务标签用户打标签进入队列失败,错误原因：%+v", err)
	}

}

// MomentDelTagProduce 导入人群包类型的企业群发朋友圈类型，删除标签
func (p *PlanTaskQueueSvc) MomentDelTagProduce(c ctx.Context, params planReq.MomentDelTagQueueData) {
	str, _ := json.Marshal(params)
	err := pkgKafka.GetAppKafka().Producer(string(str), "moment_del_tag")
	if err != nil {
		log.Info(c).Msgf("任务人群包数据朋友圈标签删除队列失败,错误原因：%+v", err)
	}

}

func (p *PlanTaskQueueSvc) MaPushFinishProduce(c ctx.Context, params planReq.MaPushFinishQueueData) {
	str, _ := json.Marshal(params)
	err := pkgKafka.GetAppKafka().Producer(string(str), "ma_push_finish")
	if err != nil {
		log.Info(c).Msgf("MA任务完成结果推送队列失败,错误原因：%+v", err)
	}

}
