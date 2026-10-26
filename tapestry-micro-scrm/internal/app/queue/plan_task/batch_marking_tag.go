package plan_task

import (
	"encoding/json"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-scrm/internal/app"
	tagsReq "micro-scrm/internal/app/http/request/tags"
	workTagSvc "micro-scrm/internal/app/service/tags"
	"micro-scrm/internal/pkg/notice"
)

type BatchMarkingTag struct {
	groupId string
	topics  []string
}

func NewBatchMarkingTagQueue() *BatchMarkingTag {
	return &BatchMarkingTag{
		groupId: "plan_task",
		topics:  []string{"batch_marking_tag"},
	}
}

func (q *BatchMarkingTag) Topics() []string {
	return q.topics
}

func (q *BatchMarkingTag) GroupID() string {
	return q.groupId
}

func (q *BatchMarkingTag) Run(c ctx.Context, key string, value *string) error {
	var params tagsReq.BatchMarkingTagReq
	err := json.Unmarshal([]byte(*value), &params)
	if err != nil {
		log.Info(c).Msgf("%s导入人群打标签数据解析失败,错误原因：%+v", key, err)
		return err
	}

	err = workTagSvc.NewWorkTags().BatchMarking(c, "计划任务人群包", params.EnterpriseID, params.ExternalUserTagList, 1)
	if err != nil {
		app.Logger().Error(c).Msgf("导入人群数据打标签失败,task_id:%d,错误原因:%s", params.TaskID, err.Error())
		notice.NewWechat().SendTextMessage(c, fmt.Sprintf("计划任务人群包-批量打标异常,task_id:%d,,错误:[%+v]", params.TaskID, err))
		return err
	}
	return nil
}
