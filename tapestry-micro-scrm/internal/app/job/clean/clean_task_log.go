package clean

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/xjob"
	"time"
)

type CleanTaskLog struct {
	name string
}

func NewCleanTaskLog() *CleanTaskLog {
	return &CleanTaskLog{name: "清理gcorn库task_log一个月前的日志数据"}
}

func (job CleanTaskLog) Run(ctx ctx.Context, handler *xjob.JobHandler) error {
	// 获取品牌的SFTP配置列表
	startTimeUnix, _ := helper.YesterdayPeriodOfTimestamp(30)
	startTime := time.Unix(startTimeUnix, 0).Format("2006-01-02 15:00:00")
	res := app.DB().WithContext(ctx).Exec("DELETE FROM `gocron`.`task_log` WHERE `start_time` < ? ", startTime)
	if res.Error != nil {
		fmt.Println(res.Error)
	}
	fmt.Println(res.RowsAffected)
	return nil
}
