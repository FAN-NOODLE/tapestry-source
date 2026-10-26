package dashboardJob

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app"
	membersDao "micro-scrm/internal/app/dao/members"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	appJob "micro-scrm/internal/app/job"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/xjob"
)

type ReserveData struct {
	name      string
	startTime int64
	endTime   int64
}

func NewReserveData() *ReserveData {
	var startTime, endTime int64 = helper.GetBetweenByTime(1)
	return &ReserveData{
		name:      "dashboard邀约数据",
		startTime: startTime,
		endTime:   endTime,
	}
}

func (job *ReserveData) Run(c ctx.Context, handler *xjob.JobHandler) error {
	enterpriseIds, err := appJob.GetEnterpriseIds(handler)
	if err != nil {
		return err
	}

	// wg := threading.NewRoutineGroup()

	for _, id := range enterpriseIds {
		enterpriseId := id
		// wg.Run(func() {
		err = job.handle(c, enterpriseId)
		if err != nil {
			app.Logger().Error(c).Msgf("error : %v", err)
		}
		// })
	}

	// wg.Wait()
	return nil
}

func (job *ReserveData) handle(c ctx.Context, enterpriseId int64) error {
	memberList, err := membersDao.NewMembers().GetMany(c, map[string]any{"enterprise_id": enterpriseId})
	if err != nil {
		return err
	}

	for _, member := range memberList {
		var saveData microScrmModel.DashboardReserveStatisticsModel
		saveData.MemberID = uint64(member.ID)
		microScrmDao.NewReserveStatistics().CreateData(c, saveData)
	}

	return nil
}
