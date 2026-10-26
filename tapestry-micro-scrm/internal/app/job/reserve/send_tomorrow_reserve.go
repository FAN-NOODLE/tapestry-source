package reserveJob

import (
	"github.com/panjf2000/ants/v2"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"math"
	"micro-scrm/internal/app/dao/micro_scrm"
	reserveReq "micro-scrm/internal/app/http/request/reserve"
	cacheSvc "micro-scrm/internal/app/service/cache"
	"micro-scrm/internal/orm"
	"micro-scrm/internal/orm/model/boss"
	"micro-scrm/internal/pkg/sdk/microwork"
	"micro-scrm/internal/pkg/xjob"
	"sync"
	"time"
)

type SendTomorrowReserveData struct {
	name string
}

func NewSendTomorrowReserveData() *SendTomorrowReserveData {
	return &SendTomorrowReserveData{
		name: "推送明天的邀约订阅消息",
	}
}

func (job *SendTomorrowReserveData) Run(ctx ctx.Context, handler *xjob.JobHandler) error {

	var (
		err            error
		enterpriseList []boss.EnterpriseModel
	)
	enterpriseList, err = orm.BossService.Enterprise.GetAllInfo()
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, enterprise := range enterpriseList {
		wg.Add(1)
		go func(enterpriseID int64) {
			defer wg.Done()
			job.getTask(ctx, enterpriseID)
		}(enterprise.ID)
	}
	wg.Wait()
	return nil

}

func (job *SendTomorrowReserveData) getTask(c ctx.Context, enterpriseId int64) error {

	config, err := cacheSvc.NewRedisData().GetEnterpriseConfig(c, enterpriseId, "store_name")
	if err != nil {
		return err
	}
	var request reserveReq.ReserveOrderListReq
	request.EnterpriseId = enterpriseId
	startTime := time.Now().Unix()
	request.ReserveStartTime = startTime
	endTime := startTime + 1800
	if config.Name == "coach" { // coach，根据到店时间提前半个小时推送模板消息
		request.ReserveEndTime = endTime
	} else { // sw和KS的企业，根据到店时间提前三个小时推送模板消息
		endTime = startTime + 3*60*60
		request.ReserveEndTime = endTime
	}
	size := 100
	results, total, err := microScrmDao.NewReservesOrder().GetReserveOrderBaseList(c, enterpriseId, 1, size, request, -1)
	if err != nil {
		return err
	}
	if int(total) > size {
		maxpage := int(math.Ceil(float64(int(total) / size)))
		maxpage += 1
		for i := 2; i <= maxpage; i++ {
			lists, _, err := microScrmDao.NewReservesOrder().GetReserveOrderBaseList(c, enterpriseId, 1, size, request, -1)
			if err != nil {
				continue
			}
			for _, list := range lists {
				results = append(results, list)
			}
		}
	}
	if len(results) == 0 {
		return nil
	}
	p, _ := ants.NewPool(10)
	defer p.Release()
	wg := sync.WaitGroup{}
	for _, info := range results {
		temp := info
		wg.Add(1)
		_ = p.Submit(func() {
			defer wg.Done()
			if temp.IsCommit != 1 {
				return
			}
			// 获取服务名称
			reserveServiceRes, err := microScrmDao.NewReserveServiceDao().GetInfoById(c, enterpriseId, temp.ReserveID)
			if err != nil {
				log.Error(c).Msgf("获取服务名称,异常：%v", err)
				return
			}
			if reserveServiceRes.ID == 0 {
				log.Error(c).Msgf("获取服务名称-记录不存在")
				return
			}
			// 查询用户的openid
			openidRelationInfo, err := microScrmDao.NewOpenidRelation().GetInfo(c, enterpriseId, map[string]interface{}{
				"unionid": temp.UserUnionid,
			})
			if err != nil {
				log.Error(c).Msgf("查询用户的openid,异常：%v", err)
				return
			}
			if openidRelationInfo.WxOpenid != "" {
				microwork.NewSubscribeMessage().Send(c, microwork.SubscribeMessageSendReq{
					EnterpriseId: enterpriseId,
					OpenId:       openidRelationInfo.WxOpenid,
					Type:         1,
					Data: microwork.SubscribeMessageSendDataReq{
						ReserveOrderId: temp.ID,
						ServiceName:    reserveServiceRes.ReserveName,
						StoreName:      temp.StoreName,
						ReserveTime:    temp.ReserveTime,
						ReserveEndTime: temp.ReserveEndTime,
					},
				})
			}
			return
		})
	}
	wg.Wait()
	return nil
}
