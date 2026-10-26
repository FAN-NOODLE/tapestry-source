package groupmember_http

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	groupmember "micro-scrm/internal/modules/group-member"
	workexternal "micro-scrm/internal/modules/work-external"
	"micro-scrm/internal/pkg/errs"
)

var workService *workexternal.Service
var groupService *groupmember.Service

func Calc(c *gin.Context) {
	var param groupmember.GroupParam
	err := c.BindJSON(&param)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	if param.Day < 0 {
		ginplus.ResError(c, errs.DayIsMust)
		return
	}

	var paramThan groupmember.GroupParam
	var realRes groupmember.GroupRes

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

	wg := threading.NewRoutineGroup()

	// 群总数
	wg.RunSafe(func() {
		numRes, err := groupService.SelectGroupNum(ctx.Wrap(c), param)
		if err != nil {
			ginplus.ResError(c, err)
			return
		}
		realRes.GroupNum = numRes
	})
	wg.RunSafe(func() {
		numResThan, err := groupService.SelectGroupNum(ctx.Wrap(c), paramThan)
		if err != nil {
			ginplus.ResError(c, err)
			return
		}
		realRes.GroupNumThan = numResThan
	})

	// 留存
	wg.RunSafe(func() {
		reRes, err := groupService.SelectReGroupUv(ctx.Wrap(c), param)
		if err != nil {
			ginplus.ResError(c, err)
			return
		}
		realRes.ReUv = reRes
	})
	wg.RunSafe(func() {
		reResThan, err := groupService.SelectReGroupUv(ctx.Wrap(c), paramThan)
		if err != nil {
			ginplus.ResError(c, err)
			return
		}
		realRes.ReUvThan = reResThan
	})

	// 新增
	wg.RunSafe(func() {
		addRes, err := groupService.SelectAddGroupUv(ctx.Wrap(c), param)
		if err != nil {
			ginplus.ResError(c, err)
			return
		}
		realRes.AddUv = addRes
	})
	wg.RunSafe(func() {
		addResThan, err := groupService.SelectAddGroupUv(ctx.Wrap(c), paramThan)
		if err != nil {
			ginplus.ResError(c, err)
			return
		}
		realRes.AddUvThan = addResThan
	})

	// 流失
	wg.RunSafe(func() {
		lossRes, err := groupService.SelectLossGroupUv(ctx.Wrap(c), param)
		if err != nil {
			ginplus.ResError(c, err)
			return
		}
		realRes.LossUv = lossRes
	})
	wg.RunSafe(func() {
		lossResThan, err := groupService.SelectLossGroupUv(ctx.Wrap(c), paramThan)
		if err != nil {
			ginplus.ResError(c, err)
			return
		}
		realRes.LossUvThan = lossResThan
	})

	wg.Wait()

	ginplus.ResSuccess(c, realRes)

	return
}
