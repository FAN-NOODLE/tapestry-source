package member

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	memberReq "micro-scrm/internal/app/http/request/member"
	membersSvc "micro-scrm/internal/app/service/member"
)

type PerformanceListHandle struct {
	handle.BaseHandle
	PerformanceListService *membersSvc.PerformanceListService
}

func NewPerformanceListHandle() PerformanceListHandle {
	return PerformanceListHandle{
		PerformanceListService: membersSvc.NewNewPerformanceListService(),
	}
}

//绩效榜活跃榜顾客榜单

func (s PerformanceListHandle) ActiveCustomers(c *gin.Context) {
	var request memberReq.PerformanceListReq
	if request.Page <= 0 {
		request.Page = 1
	}
	if request.Size <= 0 {
		request.Size = 10
	}

	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}

	result, total, err := s.PerformanceListService.ActiveCustomers(s.WarpContext(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResList(c, total, request.Page, request.Size, result)
	return
}

//邀约榜

func (s PerformanceListHandle) Invitation(c *gin.Context) {

	var request memberReq.PerformanceListInvitationReq
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}

	result, total, err := s.PerformanceListService.Invitation(s.WarpContext(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResList(c, total, request.Page, request.Size, result)
	return
}

//Achievement 业绩榜单
func (s PerformanceListHandle) Achievement(c *gin.Context) {
	var request memberReq.PerformanceListAchievementReq
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}
	result, total, err := s.PerformanceListService.Achievement(s.WarpContext(c), request, request.Start)
	if err != nil {
		s.ResError(c, err)
		return
	}
	if request.Start > 0 {
		s.ResSuccess(c, result)
	} else {
		s.ResList(c, total, request.Page, request.Size, result)
	}
	return
}

//SwCustomers sw顾客榜单
func (s PerformanceListHandle) SwCustomers(c *gin.Context) {
	var request memberReq.SwCustomersListReq
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}
	result, total, err := s.PerformanceListService.SwCustomers(s.WarpContext(c), request, request.Start)
	if err != nil {
		s.ResError(c, err)
		return
	}
	if request.Start > 0 {
		s.ResSuccess(c, result)
	} else {
		s.ResList(c, total, request.Page, request.Size, result)
	}
	return
}

//SwActive sw活跃榜单
func (s PerformanceListHandle) SwActive(c *gin.Context) {
	var request memberReq.SwActiveListReqReq
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}
	result, total, err := s.PerformanceListService.SwActive(s.WarpContext(c), request, request.Start)
	if err != nil {
		s.ResError(c, err)
		return
	}
	if request.Start > 0 {
		s.ResSuccess(c, result)
	} else {
		s.ResList(c, total, request.Page, request.Size, result)
	}

	return
}

//明星榜-活跃榜

func (s PerformanceListHandle) StarActiveCustomers(c *gin.Context) {
	var request memberReq.PerformanceListReq
	if request.Page <= 0 {
		request.Page = 1
	}
	if request.Size <= 0 {
		request.Size = 10
	}
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}

	result, err := s.PerformanceListService.StarActiveCustomers(s.WarpContext(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, result)
	return
}

//明星榜-邀约榜

func (s PerformanceListHandle) StarInvitation(c *gin.Context) {

	var request memberReq.PerformanceListInvitationReq
	if err := c.ShouldBind(&request); err != nil {
		s.ResError(c, err)
		return
	}

	result, err := s.PerformanceListService.StarInvitation(s.WarpContext(c), request)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, result)
	return
}
