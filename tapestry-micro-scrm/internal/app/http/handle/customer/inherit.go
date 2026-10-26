package customer

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"micro-scrm/internal/app/http/handle"
	customerReq "micro-scrm/internal/app/http/request/customer"
	memberReq "micro-scrm/internal/app/http/request/member"
	customerModels "micro-scrm/internal/app/models/customer"
	customerSvc "micro-scrm/internal/app/service/customer"
	"micro-scrm/internal/modules/transfer_customer/delivery/job"
	"micro-scrm/internal/pkg/redis"
	"time"
)

type InheritHandle struct {
	handle.BaseHandle
	inheritService *customerSvc.Inherit
}

func NewInheritHandle() InheritHandle {
	return InheritHandle{
		inheritService: customerSvc.NewInherit(),
	}
}

// List 离职员工查询列表
func (h InheritHandle) List(c *gin.Context) {

	var request customerReq.CustomerListReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	if request.Status != customerModels.InheritStatusAssigned && request.Status != customerModels.InheritStatusUnassigned {
		h.ResError(c, errors.New("参数有误"))
		return
	}
	// 列表
	result, total, err := h.inheritService.GetInheritList(ctx.Wrap(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	h.ResList(c, total, request.Page, request.Size, result)

	return
}

// GetGroupChatAllotList 获取离职人员的待分配群列表
func (h InheritHandle) GetGroupChatAllotList(c *gin.Context) {
	var request customerReq.GroupChatList
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	var result customerReq.GroupChatRes
	list, errs := h.inheritService.GetGroupChatAllotList(ctx.Wrap(c), request)
	if errs != nil {
		h.ResError(c, errs)
		return
	}
	total := len(list)
	result.Total = int64(total)
	start := (request.Page - 1) * request.Size
	end := request.Page * request.Size
	var lists []customerReq.GroupChatInfo
	if end >= total {
		lists = list[start:total]
	} else {
		lists = list[start:end]
	}
	result.List = lists
	h.ResSuccess(c, result)
}

// GetExternalUserAllotList 获取离职人员的待分配人员列表
func (h InheritHandle) GetExternalUserAllotList(c *gin.Context) {
	var request customerReq.ExternalListReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	result, total, listErr := h.inheritService.GetExternalAllotList(h.WarpContext(c), request)
	if listErr != nil {
		h.ResError(c, listErr)
		return
	}

	h.ResList(c, total, request.Page, request.Size, result)

}

// GetTakeOverList 已分配的导购列表
func (h InheritHandle) GetTakeOverList(c *gin.Context) {
	var request customerReq.GroupChatList
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	var result customerReq.TakeOverList
	total, totalErr := h.inheritService.GetTakeOverCount(h.WarpContext(c), request.EnterpriseID, request.HandoverUserId)
	if totalErr != nil {
		h.ResError(c, totalErr)
		return
	}
	result.Total = total
	result.List = []customerReq.TakeOverInfo{}
	if total > 0 {
		takeOverList, errs := h.inheritService.GetTakeOverList(h.WarpContext(c), request)
		if errs != nil {
			h.ResError(c, errs)
			return
		}
		result.List = takeOverList
	}
	h.ResSuccess(c, result)
}

func (h InheritHandle) TransferCustomer(c *gin.Context) {
	var request memberReq.ResignedTransferCustomerReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	externalUserIdLength := len(request.ExternalUserid)
	if externalUserIdLength == 0 {
		h.ResError(c, errors.New("请选择要分配的用户"))
		return
	}
	var successUserRes []string
	var errorUserRes []string
	if externalUserIdLength > 100 { //企微接口，离职分配每次只能处理100个外部联系人，这里做分批处理
		externalUserIds := request.ExternalUserid
		page := externalUserIdLength / 100
		if (externalUserIdLength % 100) > 0 {
			page += 1
		}
		for i := 0; i < page; i++ {
			start := i * 100
			end := (i + 1) * 100
			if end > externalUserIdLength {
				end = externalUserIdLength
			}
			request.ExternalUserid = externalUserIds[start:end]
			successUser, errorUser, err := h.inheritService.AllocationCustomer(h.WarpContext(c), request)
			if err != nil {
				h.ResError(c, err)
				return
			}
			successUserRes = append(successUserRes, successUser...)
			errorUserRes = append(errorUserRes, errorUser...)
		}
	} else {
		successUser, errorUser, err := h.inheritService.AllocationCustomer(h.WarpContext(c), request)
		if err != nil {
			h.ResError(c, err)
			return
		}
		successUserRes = append(successUserRes, successUser...)
		errorUserRes = append(errorUserRes, errorUser...)
	}

	h.ResSuccess(c, gin.H{
		"success": successUserRes,
		"error":   errorUserRes,
	})
	return
}

func (h InheritHandle) TransferCustomerGroup(c *gin.Context) {
	var request memberReq.ResignedTransferCustomerGroupReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}

	successIds, errorIds, err := h.inheritService.AllocationGroup(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, gin.H{
		"success": successIds,
		"error":   errorIds,
	})

	return
}

// GetTakeOverExternalList 获取接替的外部联系人列表
func (h InheritHandle) GetTakeOverExternalList(c *gin.Context) {
	var request customerReq.ExternalListReq
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	var result customerReq.TakeOverExternalRes
	total, list, err := h.inheritService.GetTakeOverExternalList(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	result.Total = total
	result.List = list
	h.ResSuccess(c, result)
}

func (h InheritHandle) GetTakeOverGroupList(c *gin.Context) {
	var request customerReq.GroupChatList
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	var result customerReq.TakeOverGroupRes
	total, list, err := h.inheritService.GetTakeOverGroupList(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}
	result.Total = total
	result.List = list
	h.ResSuccess(c, result)
}

func (h InheritHandle) SyncData(c *gin.Context) {
	var request customerReq.SyncCustomDataRes
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	syncKey := "sync_inherit_time_" + fmt.Sprintf("%d", request.EnterpriseID)
	lastSyncTime := redis.Client().Get(c, syncKey).Val()
	now := time.Now().Unix()
	if request.ExeSync != true {
		if lastSyncTime == "" {
			h.ResSuccess(c, now)
			return
		} else {
			h.ResSuccess(c, lastSyncTime)
			return
		}
	}
	if lastSyncTime != "" {
		h.ResError(c, errors.New("企微数据同步正在进行中，请稍后再试"))
		return
	}
	go job.NewSaveCustomerInheritData().Exec(h.WarpContext(c), request.EnterpriseID)
	h.ResSuccess(c, now)
}
