package dashboard

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	dashboardReq "micro-scrm/internal/app/http/request/dashboard"
)

type PlanHandle struct {
	handle.BaseHandle
}

func NewPlanHandle() PlanHandle {
	return PlanHandle{}
}

// StoreData 门店数据
func (s PlanHandle) StoreData(c *gin.Context) {
	var request dashboardReq.PlanStatistics
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}

	s.ResSuccess(c, "")

	return
}

// PlanData 我的数据
func (s PlanHandle) PlanData(c *gin.Context) {
	var request dashboardReq.PlanStatistics
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}

	s.ResSuccess(c, "")
	return
}

// PlanList 我的任务列表
func (s PlanHandle) PlanList(c *gin.Context) {
	var request dashboardReq.PlanStatistics
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}

	s.ResSuccess(c, "")
	return
}

// PlanDetail 我的任务详情
func (s PlanHandle) PlanDetail(c *gin.Context) {
	var request dashboardReq.PlanStatistics
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}

	s.ResSuccess(c, "")
	return
}

// SowTaskList 我的播种列表
func (s PlanHandle) SowTaskList(c *gin.Context) {
	var request dashboardReq.PlanStatistics
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}

	s.ResSuccess(c, "")
	return
}

// SowTaskDetail 我的播种详情
func (s PlanHandle) SowTaskDetail(c *gin.Context) {
	var request dashboardReq.PlanStatistics
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}

	s.ResSuccess(c, "")
	return
}
