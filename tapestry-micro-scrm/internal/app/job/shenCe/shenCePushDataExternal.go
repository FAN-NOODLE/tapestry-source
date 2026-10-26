package shenCeJob

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	"micro-scrm/internal/app"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	appJob "micro-scrm/internal/app/job"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/sdk/microwork"
	"micro-scrm/internal/pkg/xjob"
)

type ShenCePushDataExternal struct {
	name string

	shenCeClient microwork.ShenCe
	startTime    int64
	endTime      int64
}

func NewShenCePushDataExternal() *ShenCePushDataExternal {
	var startTime, endTime int64 = helper.GetBetweenByTime(1)
	return &ShenCePushDataExternal{
		name:         "神策数据推送",
		shenCeClient: microwork.NewShenCe(),
		startTime:    startTime,
		endTime:      endTime,
	}
}

func (job *ShenCePushDataExternal) Run(c ctx.Context, handler *xjob.JobHandler) error {
	enterpriseIds, err := appJob.GetEnterpriseIds(handler)
	if err != nil {
		return err
	}

	wg := threading.NewRoutineGroup()

	for _, enterpriseId := range enterpriseIds {
		id := enterpriseId
		wg.RunSafe(func() {
			err = job.handle(c, id)
			if err != nil {
				app.Logger().Error(c).Msgf("error : %v", err)
			}
		})
	}

	wg.Wait()

	return nil
}

func (job *ShenCePushDataExternal) handle(c ctx.Context, enterpriseId int64) error {

	// 外部联系人
	if err := job.userCorpExternal(c, enterpriseId); err != nil {
		app.Logger().Error(c).Msgf("外部联系人数据推送 error : %v")
	}

	return nil
}

// 外部联系人数据推送
func (job *ShenCePushDataExternal) userCorpExternal(c ctx.Context, enterpriseId int64) error {
	id := int64(0)
	for {
		wxUserCorpExternals, err := microScrmDao.NewWxUserCorpExternalCk().GetManyByPage(c, enterpriseId, 1000, " id >=? and unionid <> ''", id)
		if err != nil {
			return err
		}
		index := len(wxUserCorpExternals)
		if index == 0 {
			break
		}
		id = wxUserCorpExternals[index-1].ID
		for _, user := range wxUserCorpExternals {
			properties := map[string]any{
				"unionid":       user.Unionid,
				"nick":          user.Nick,
				"sex":           user.Sex,
				"avatar":        user.Avatar,
				"union_code":    user.UnionCode,
				"clt_type":      user.CltType,
				"customer_code": user.CustomerCode,
				"card_level":    user.CardLevel,
				"created_at":    helper.GetFormatDateByParams(user.CreatedAt),
			}
			_, err = job.shenCeClient.External(c, enterpriseId, user.Openid, properties)
		}

	}

	return nil
}
