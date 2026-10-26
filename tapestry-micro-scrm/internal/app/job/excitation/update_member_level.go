package excitation

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	"micro-scrm/internal/app"
	membersDao "micro-scrm/internal/app/dao/members"
	microBossDao "micro-scrm/internal/app/dao/micro_boss"
	excitationReq "micro-scrm/internal/app/http/request/excitation"
	appJob "micro-scrm/internal/app/job"
	membersModel "micro-scrm/internal/app/models/members"
	excitationSvc "micro-scrm/internal/app/service/excitation"
	"micro-scrm/internal/pkg/xjob"
)

type UpdateMemberLevel struct {
	name          string
	memberDataSvc *excitationSvc.MemberDataService
}

func NewUpdateMemberLevel() *UpdateMemberLevel {
	return &UpdateMemberLevel{
		name:          "每天更新导购等级数据",
		memberDataSvc: excitationSvc.NewMemberDataService(),
	}
}

func (job *UpdateMemberLevel) Run(c ctx.Context, handler *xjob.JobHandler) error {
	enterpriseIds, err := appJob.GetEnterpriseIds(handler)
	if err != nil {
		return err
	}

	var wg threading.RoutineGroup
	for _, id := range enterpriseIds {
		tmpId := id
		wg.Run(func() {
			err := job.Handle(c, tmpId)
			if err != nil {
				app.Logger().Error(c).Msgf("error : %v", err)
			}
		})
	}

	wg.Wait()

	return nil
}

func (job *UpdateMemberLevel) Handle(c ctx.Context, enterpriseId int64) error {

	// 查询配置的企业
	resConfig, err := microBossDao.NewEnterpriseConfig().
		GetInfo(c, map[string]interface{}{"enterprise_id": enterpriseId, "key": "update_member_level"})
	if err != nil {
		return err
	}
	if resConfig.ID == 0 {
		return nil
	}

	pageLimit := 200
	beginID := int64(0)

	for {
		var rcResult []membersModel.MemberModel

		// 分页批量获取导购列表
		rcResult, err := membersDao.NewMembers().GetManyLimit(c, enterpriseId, pageLimit, beginID)
		if err != nil {
			return err
		}

		if len(rcResult) > 0 {
			for _, row := range rcResult {

				var params excitationReq.UpdateLevelReq
				params.EnterpriseId = enterpriseId
				params.MemberId = row.ID

				errs := job.memberDataSvc.UpdateLevel(c, params)
				if errs != nil {
					return errs
				}

			}
			beginID = rcResult[len(rcResult)-1].ID

		}
		if len(rcResult) < pageLimit {
			break
		}
	}

	return nil

}
