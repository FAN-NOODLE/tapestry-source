package dashboard

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	"micro-scrm/internal/app/http/request"
	dashboardReq "micro-scrm/internal/app/http/request/dashboard"
	dashboardSvc "micro-scrm/internal/app/service/dashboard"
)

type ManagerHandle struct {
	handle.BaseHandle
	managerSvc *dashboardSvc.ManagerService
}

func NewManagerHandle() ManagerHandle {
	return ManagerHandle{
		managerSvc: dashboardSvc.NewMemberService(),
	}
}

// MemberAuthList  获取管理者dashboard权限列表
func (s ManagerHandle) MemberAuthList(c *gin.Context) {
	var req dashboardReq.ManagerAuthList
	if err := c.ShouldBind(&req); err != nil {
		s.ResError(c, err)
		return
	}
	list, total, err := s.managerSvc.ManagerAuthList(s.WarpContext(c), req)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResList(c, total, req.Page, req.Size, list)

	return
}

// MemberAuthDetail 获取导购dashboard权限
func (s ManagerHandle) MemberAuthDetail(c *gin.Context) {
	var req dashboardReq.ManagerAuthDetail
	if err := c.ShouldBind(&req); err != nil {
		s.ResError(c, err)
		return
	}
	detail, err := s.managerSvc.ManagerAuthDetail(s.WarpContext(c), req)

	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, detail)

	return
}

// SaveManagerAuth  创建/更新管理者dashboard权限
func (s ManagerHandle) SaveManagerAuth(c *gin.Context) {
	var req dashboardReq.ManagerAuthSave
	if err := c.ShouldBind(&req); err != nil {
		s.ResError(c, err)
		return
	}
	err := s.managerSvc.SaveManagerAuth(s.WarpContext(c), req)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, "")

	return
}

// DeleteManagerAuth    删除管理者dashboard权限
func (s ManagerHandle) DeleteManagerAuth(c *gin.Context) {
	var req dashboardReq.ManagerRemove
	if err := c.ShouldBind(&req); err != nil {
		s.ResError(c, err)
		return
	}
	err := s.managerSvc.DeleteManagerAuth(s.WarpContext(c), req)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, "")

	return
}

// GetAllManagerIds  所有已选管理者memberIds
func (s ManagerHandle) GetAllManagerIds(c *gin.Context) {
	var req request.EnterpriseBaseReq
	if err := c.ShouldBind(&req); err != nil {
		s.ResError(c, err)
		return
	}
	ids, err := s.managerSvc.GetAllManagerMemberIds(s.WarpContext(c), req.EnterpriseId)
	if err != nil {
		return
	}
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, ids)

	return
}
