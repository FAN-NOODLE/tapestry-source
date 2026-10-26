package appJob

import (
	"fmt"
	"micro-scrm/internal/app"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/pkg/notice"
	"micro-scrm/internal/pkg/xjob"
	"strconv"
)

func GetEnterpriseIds(jobHandler *xjob.JobHandler) ([]int64, error) {
	specsArgs := jobHandler.Runner().Specs()
	enterpriseIds := make([]int64, 0)
	if len(specsArgs) > 0 {
		for _, arg := range specsArgs {
			enterpriseId, _ := strconv.ParseInt(arg, 10, 64)
			enterpriseIds = append(enterpriseIds, enterpriseId)
		}
	} else {
		enterprises, err := orm.BossService.Enterprise.GetAllInfo()
		if err != nil || len(enterprises) == 0 {
			return nil, err
		}
		for _, val := range enterprises {
			enterpriseIds = append(enterpriseIds, val.ID)
		}
	}
	return enterpriseIds, nil
}

type JobData struct {
	EnterpriseName string `json:"enterprise_name" comment:"企业名称"`
	JobName        string `json:"job_name" comment:"job名称"`
	Content        string `json:"content" comment:"执行结果"`
}

// JobNotice   job通知
func JobNotice(c app.Ctx, data JobData) {
	content := fmt.Sprintf("脚本执行结果通知：\n"+
		">企业名称:<font color='comment'>%s</font>\n"+
		">脚本名称:<font color='comment'>%s</font>\n"+
		">执行结果：<font color='red'>%s</font>\n"+
		">RequestId：<font color='comment'>%v</font>\n",
		data.EnterpriseName,
		data.JobName,
		data.Content,
		c.ID(),
	)
	notice.NewWechat().SendTaskMarkdownMessage(c, content)
}
