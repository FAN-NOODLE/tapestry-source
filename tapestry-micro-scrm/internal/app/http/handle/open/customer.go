package open

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	"micro-scrm/internal/app/http/middleware"
	openReq "micro-scrm/internal/app/http/request/open"
	cacheSvc "micro-scrm/internal/app/service/cache"
	customerSvc "micro-scrm/internal/app/service/customer"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/redis"
	"time"
)

type Customer struct {
	handle.BaseHandle
}

func NewCustomer() Customer {
	return Customer{}
}

// BatchUpdateCouponStatus   顾客优惠券状态批量更新
func (h Customer) BatchUpdateCouponStatus(c *gin.Context) {
	var request openReq.BatchUpdateCouponReq

	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	if len(request.DataList) > 10000 {
		h.ResError(c, errors.New("入参数量过大"))
		return
	}
	request.EnterpriseId = c.GetInt64(middleware.EnterpriseId)
	config, _ := cacheSvc.NewRedisData().GetEnterpriseConfig(h.WarpContext(c), request.EnterpriseId, "store_name")

	if config.Name != "coach" {
		h.ResError(c, errors.New("账号无权限"))
		return
	}
	//每一分钟更新一次数据
	redisClient := redis.Client()

	lockKey := fmt.Sprintf("%d:batch_update_coupon_status_lock", request.EnterpriseId)
	dataKey := fmt.Sprintf("%d:batch_update_coupon_status_data", request.EnterpriseId)
	// 将数据存入Redis list
	for _, val := range request.DataList {
		_ = redisClient.LPush(c, dataKey, helper.MarshalToString(val))
	}

	// 尝试获取锁(1分钟)
	acquired, err := redisClient.SetNX(c, lockKey, "1", 1*time.Minute).Result()
	if err != nil {
		h.ResError(c, err)
		return
	}

	if acquired {
		// 获取到锁,说明是第一个请求,启动1分钟后的批量更新任务
		// 捕获 gin.Context 和 Customer,避免 goroutine 中使用已失效的变量
		capturedC := c
		capturedH := h
		capturedEnterpriseId := request.EnterpriseId

		go func() {
			time.Sleep(1 * time.Minute)

			// 从Redis list中取出所有数据
			var dataList []openReq.UpdateCouponParams
			for {
				val, err := redisClient.RPop(capturedC, dataKey).Result()
				if err != nil {
					break
				}
				var data openReq.UpdateCouponParams
				if err := json.Unmarshal([]byte(val), &data); err != nil {
					continue
				}
				dataList = append(dataList, data)
			}

			// 执行批量更新
			if len(dataList) > 0 {
				batchRequest := openReq.BatchUpdateCouponReq{
					EnterpriseId: capturedEnterpriseId,
					DataList:     dataList,
				}
				_ = customerSvc.NewExternalSearch().BatchSaveCouponExtraData(capturedH.WarpContext(capturedC), batchRequest)
			}

			// 释放锁
			redisClient.Del(capturedC, lockKey)
		}()
	}

	h.ResSuccess(c, nil)
	return
}

// BatchUpdateActiveStatus    顾客优惠券状态批量更新
func (h Customer) BatchUpdateActiveStatus(c *gin.Context) {
	var request openReq.BatchUpdateActiveReq

	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	if len(request.DataList) == 0 {
		h.ResError(c, errors.New("数据不能为空"))
		return
	}
	if len(request.DataList) > 3000 {
		h.ResError(c, errors.New("入参数量过大"))
		return
	}
	request.EnterpriseId = c.GetInt64(middleware.EnterpriseId)
	err := customerSvc.NewExternalSearch().BatchSaveActiveData(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, nil)
	return
}
