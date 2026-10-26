package open

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/lock"
	"micro-scrm/internal/app/http/handle"
	"micro-scrm/internal/app/http/middleware"
	openReq "micro-scrm/internal/app/http/request/open"
	excitationSvc "micro-scrm/internal/app/service/excitation"
	"time"
)

type ExcitationHandle struct {
	handle.BaseHandle
}

func NewExcitationHandle() ExcitationHandle {
	return ExcitationHandle{}
}

// MemberInfo 获取导购金币
func (h ExcitationHandle) MemberInfo(c *gin.Context) {
	var request openReq.MemberInfoReq
	request.EnterpriseId = c.GetInt64(middleware.EnterpriseId)
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	res, err := excitationSvc.NewMemberDataService().MemberDataInfo(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, res)
}

// AddOrder  积分商城下单
func (h ExcitationHandle) AddOrder(c *gin.Context) {
	var request openReq.AddOrderReq
	request.EnterpriseId = c.GetInt64(middleware.EnterpriseId)
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	//针对每个导购生成独立的下单并发锁
	lockKey := fmt.Sprintf("lock:open:excitation_order:%d:%s", request.EnterpriseId, request.UserId)
	orderLock, lErr := lock.Mutex(
		h.WarpContext(c),
		lock.Key(lockKey),
		lock.TTL(3*time.Second),
		lock.Retry(lock.NoRetry()),
	)
	if lErr != nil {
		h.ResError(c, errors.New("订单操作并发锁生成失败"))
		return
	}
	if err := orderLock.Lock(); err != nil {
		h.ResError(c, errors.New("请求过于频繁，请稍后再试"))
		return
	}
	defer func() {
		_ = orderLock.Unlock()
	}()
	err := excitationSvc.NewCoinsService().AddOrder(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResOk(c)
}

// CancelOrder AddOrder  积分商城退单
func (h ExcitationHandle) CancelOrder(c *gin.Context) {
	var request openReq.CancelOrderReq
	request.EnterpriseId = c.GetInt64(middleware.EnterpriseId)
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	//针对每个导购生成独立的下单并发锁
	lockKey := fmt.Sprintf("lock:open:excitation_order:%d:%s", request.EnterpriseId, request.UserId)
	orderLock, lErr := lock.Mutex(
		h.WarpContext(c),
		lock.Key(lockKey),
		lock.TTL(3*time.Second),
		lock.Retry(lock.NoRetry()),
	)
	if lErr != nil {
		h.ResError(c, errors.New("订单操作并发锁生成失败"))
		return
	}
	if err := orderLock.Lock(); err != nil {
		h.ResError(c, errors.New("请求过于频繁，请稍后再试"))
		return
	}
	defer func() {
		_ = orderLock.Unlock()
	}()
	err := excitationSvc.NewCoinsService().CancelOrder(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResOk(c)
}

// OrderInfo 订单状态查询
func (h ExcitationHandle) OrderInfo(c *gin.Context) {
	var request openReq.OrderInfoReq
	request.EnterpriseId = c.GetInt64(middleware.EnterpriseId)
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	res, err := excitationSvc.NewCoinsService().OrderInfo(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResSuccess(c, res)
}
