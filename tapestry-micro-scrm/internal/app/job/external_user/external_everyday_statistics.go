package externalUserJob

import (
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	"micro-scrm/internal/app"
	microScrmDao "micro-scrm/internal/app/dao/micro_scrm"
	appJob "micro-scrm/internal/app/job"
	"micro-scrm/internal/app/models"
	microScrmModel "micro-scrm/internal/app/models/micro_scrm"
	membersSvc "micro-scrm/internal/app/service/member"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/xjob"
	"time"
)

type ExternalEverydayStatistics struct {
	name                          string
	membersSvc                    *membersSvc.MembersService
	userRelationCk                *microScrmDao.UserRelation
	externalEverydayStatisticsDao *microScrmDao.ExternalEverydayStatisticsDao
}

func NewExternalEverydayStatistics() *ExternalEverydayStatistics {
	return &ExternalEverydayStatistics{
		name:                          "累计数据获得积分任务",
		membersSvc:                    membersSvc.NewMembersService(),
		userRelationCk:                microScrmDao.NewUserRelation(),
		externalEverydayStatisticsDao: microScrmDao.NewExternalEverydayStatisticsDao(),
	}
}

func (job *ExternalEverydayStatistics) Run(c ctx.Context, handler *xjob.JobHandler) error {
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

func (job *ExternalEverydayStatistics) Handle(c ctx.Context, enterpriseId int64) error {
	// 获取导购最新一条统计数据
	dataList, err := job.externalEverydayStatisticsDao.GetLastData(c, enterpriseId)
	if err != nil {
		fmt.Printf("查询最新数据错误（没有建表的企业不用执行）:%s", err.Error())
		return err
	}
	var memberToData = make(map[int64]microScrmModel.ExternalEverydayStatisticsModel)
	for _, data := range dataList {
		memberToData[data.MemberID] = data
	}
	// 查询当前所有导购的三方会员信息数量
	statisticsData, err := job.userRelationCk.ExternalStatistics(c, enterpriseId)
	if err != nil {
		return err
	}
	// 循环处理导购相关客户的三方会员信息统计
	var addData []microScrmModel.ExternalEverydayStatisticsModel
	for _, statistic := range statisticsData {
		//如果目前不存在数据，直接不插入新记录
		//默认如果没有老数据则直接全部已新数据为标准
		tmp := microScrmModel.ExternalEverydayStatisticsModel{
			BaseModel:           models.CreateBase(app.UIDContext(c)),
			MemberID:            statistic.MemberID,
			UnixTime:            int(helper.GetZeroTimeByTime(time.Now()).Unix()),
			Date:                helper.GetNowDate(),
			NewCustomer:         int(statistic.NewCustomers),
			OldCustomer:         int(statistic.OldCustomers),
			SleepCustomer:       int(statistic.SleepCustomers),
			VipCustomer:         int(statistic.VipCustomers),
			NewCustomerChange:   int(statistic.NewCustomers),
			OldCustomerChange:   int(statistic.OldCustomers),
			SleepCustomerChange: int(statistic.SleepCustomers),
			VipCustomerChange:   int(statistic.VipCustomers),
		}
		//如果存在老数据则需要计算变化量
		if data, ok := memberToData[statistic.MemberID]; ok {
			tmp.NewCustomerChange = int(statistic.NewCustomers) - data.NewCustomer
			tmp.OldCustomerChange = int(statistic.OldCustomers) - data.OldCustomer
			tmp.SleepCustomerChange = int(statistic.SleepCustomers) - data.SleepCustomer
			tmp.VipCustomerChange = int(statistic.VipCustomers) - data.VipCustomer
		}
		addData = append(addData, tmp)
	}
	err = job.externalEverydayStatisticsDao.CreateInBatchesData(c, enterpriseId, addData)
	if err != nil {
		return err
	}
	return nil
}
