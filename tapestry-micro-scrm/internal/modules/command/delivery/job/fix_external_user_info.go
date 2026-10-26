package job

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/workwechat"
	"micro-scrm/internal/pkg/xjob"
	"strconv"
	"sync"
)

// 根据用户的work_external_user_id，修复用户的信息

const PAGE_SIZE = 500

type fixExternalUserInfoJob struct {
	name string
}

func NewFixExternalUserInfoJob() *fixExternalUserInfoJob {
	return &fixExternalUserInfoJob{
		name: "修复外部联系人的信息Job",
	}
}

func (job fixExternalUserInfoJob) Run(ctx ctx.Context, jobHandler *xjob.JobHandler) error {
	// 获取命令行参数配置的企业
	specsArgs := jobHandler.Runner().Specs()
	enterpriseIDs := []int64{}
	for _, param := range specsArgs {
		if enterpriseID, err := strconv.ParseInt(param, 10, 64); err == nil {
			enterpriseIDs = append(enterpriseIDs, enterpriseID)
		}
	}

	var wg sync.WaitGroup
	for _, id := range enterpriseIDs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			newFixExternalUser(job.name, id).exec(ctx)
		}()
	}
	wg.Wait()
	return nil
}

// ============================================================================
type fixExternalUser struct {
	name                  string
	enterpriseID          int64
	wxUserCorpExternalDao *microScrmDao.WxUserCorpExternal
}

func newFixExternalUser(name string, enterpriseID int64) *fixExternalUser {
	return &fixExternalUser{
		name:                  fmt.Sprintf("%s::企业ID:%d", name, enterpriseID),
		enterpriseID:          enterpriseID,
		wxUserCorpExternalDao: microScrmDao.NewWxUserCorpExternal(),
	}
}

func (fix fixExternalUser) exec(ctx ctx.Context) {

	list, _ := fix.wxUserCorpExternalDao.GetMany(ctx, fix.enterpriseID, "nick = '' AND unionid!=''")
	newClient, _ := workwechat.NewWorkWeChat(ctx, fix.enterpriseID)

	// 设置并发度为 10
	concurrency := 30
	ch := make(chan struct{}, concurrency)
	for _, data := range list {
		ch <- struct{}{} // 发送一个信号，表示占用一个并发槽
		go func(data microScrmModel.WxUserCorpExternalModel) {
			defer func() { <-ch }() // 释放一个并发槽
			userInfo, userInfoErr := newClient.GetWorkExternalUserInfo(data.Openid)
			if userInfoErr != nil {
				fmt.Println(userInfoErr)
				return
			}
			if userInfo.Base.Name != "" {
				cond := map[string]interface{}{"openid": data.Openid}
				updateData := map[string]interface{}{
					"nick": userInfo.Base.Name,
				}
				if updateErr := fix.wxUserCorpExternalDao.UpdateData(ctx, fix.enterpriseID, cond, updateData); updateErr != nil {
					log.Info(ctx).Msgf("%s::根据外部联系人ID:%d, 更新数据库失败:%v", fix.name, data.Openid, updateErr)
				} else {
					log.Info(ctx).Msgf("%s::根据外部联系人ID:%d, 更新数据库成功:%v", fix.name, data.Openid, updateErr)
				}
			}
		}(data)
	}

	// 等待所有 goroutine 执行完毕
	for i := 0; i < concurrency; i++ {
		ch <- struct{}{}
	}

}
