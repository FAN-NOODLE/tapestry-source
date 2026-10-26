package work_external_http

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	workexternal "micro-scrm/internal/modules/work-external"
	"micro-scrm/internal/pkg/errs"
)

var workService *workexternal.Service
var enterpriseService workexternal.EnterpriseService

func Calc(c *gin.Context) {
	cc := ctx.Wrap(c)
	// 参数校验
	var param workexternal.SelectKeep
	if err := c.ShouldBind(&param); err != nil {
		ginplus.ResError(c, err)
		return
	}
	if param.Day < 0 {
		ginplus.ResError(c, errs.DayIsMust)
		return
	}

	var paramThan workexternal.SelectKeep
	var realRes workexternal.QwFriend

	// 计算时间
	start, end := workService.GetBetweenByTime(param.Day)
	param.StartAt = start
	param.EndAt = end
	paramThan = param
	// 环比的时间差
	timeSub := param.Day * 24 * 60 * 60
	if param.Day == 0 {
		timeSub = 24 * 60 * 60
	}
	paramThan.StartAt = param.StartAt - timeSub
	paramThan.EndAt = param.EndAt - timeSub

	group := threading.NewRoutineGroup()

	// 留存用户数
	group.RunSafe(func() {
		reRes, err := workService.SelectReUv(cc, param)
		if err != nil {
			ginplus.ResError(c, err)
			return
		}
		realRes.ReUv = reRes
	})
	group.RunSafe(func() {
		reResThan, err := workService.SelectReUv(cc, paramThan)
		if err != nil {
			ginplus.ResError(c, err)
			return
		}
		realRes.ReUvThan = reResThan
	})

	// 新增好友数 - 按客户第一次绑定关系算一次，之后不管跟谁绑定，什么时候绑定，都不算
	group.RunSafe(func() {
		addRes, err := workService.SelectAddUv(cc, param)
		if err != nil {
			ginplus.ResError(c, err)
			return
		}
		realRes.AddUv = addRes
	})
	group.RunSafe(func() {
		addResThan, err := workService.SelectAddUv(cc, paramThan)
		if err != nil {
			ginplus.ResError(c, err)
			return
		}
		realRes.AddUvThan = addResThan
	})

	// 流失好友数
	group.RunSafe(func() {
		lossRes, err := workService.SelectLossUv(cc, param)
		if err != nil {
			ginplus.ResError(c, err)
			return
		}
		realRes.LossUv = lossRes
	})
	group.RunSafe(func() {
		lossResThan, err := workService.SelectLossUv(cc, paramThan)
		if err != nil {
			ginplus.ResError(c, err)
			return
		}

		realRes.LossUvThan = lossResThan
	})

	group.Wait()

	ginplus.ResSuccess(c, realRes)

	return
}
