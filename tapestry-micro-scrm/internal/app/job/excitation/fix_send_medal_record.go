package excitation

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"log"
	"micro-scrm/internal/app"
	excitationDao "micro-scrm/internal/app/dao/excitation"
	appJob "micro-scrm/internal/app/job"
	"micro-scrm/internal/pkg/xjob"
)

type FixSendMedalRecord struct {
	name string
}

func NewFixSendMedalRecord() *FixSendMedalRecord {
	return &FixSendMedalRecord{
		name: "修复发送奖章记录",
	}
}

func (job *FixSendMedalRecord) Run(c ctx.Context, handler *xjob.JobHandler) error {
	enterpriseIds, err := appJob.GetEnterpriseIds(handler)
	if err != nil {
		return err
	}

	// var wg threading.RoutineGroup

	for _, id := range enterpriseIds {
		tmpId := id
		// wg.Run(func() {
		err = job.Handle(c, tmpId)
		if err != nil {
			app.Logger().Error(c).Msgf("error : %v", err)
		}
		// })
	}

	// wg.Wait()

	return nil
}

func (job *FixSendMedalRecord) Handle(c ctx.Context, enterpriseId int64) error {
	var nineIds []int64
	var tenIds []int64
	// 获取当前企业的奖章信息
	nineMany, err := excitationDao.NewMedalStyle().GetMany(c, "time_name = ?", "P9")
	if err != nil {
		return err
	}
	for _, res := range nineMany {
		nineIds = append(nineIds, res.ID)
	}

	tenMany, err := excitationDao.NewMedalStyle().GetMany(c, "time_name = ?", "P10")
	if err != nil {
		return err
	}
	for _, res := range tenMany {
		tenIds = append(tenIds, res.ID)
	}

	for i, id := range nineIds {
		where := map[string]any{
			"enterprise_id":  enterpriseId,
			"medal_style_id": id,
			"medal_id":       nineMany[i].MedalID,
		}
		data := map[string]any{"medal_style_id": tenIds[i]}
		err = excitationDao.NewMedalRecord().UpdateData(c, where, data)
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}
