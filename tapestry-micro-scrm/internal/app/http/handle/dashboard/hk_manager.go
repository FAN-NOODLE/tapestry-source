package dashboard

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	"micro-scrm/internal/app/http/request"
	dashboardReq "micro-scrm/internal/app/http/request/dashboard"
	dashboardSvc "micro-scrm/internal/app/service/dashboard"
)

type HkManagerHandle struct {
	handle.BaseHandle
	managerSvc *dashboardSvc.HkManagerService
}

func NewHkManagerHandle() HkManagerHandle {
	return HkManagerHandle{
		managerSvc: dashboardSvc.NewHkManagerService(),
	}
}

// StoreList 获取香港门店列表
func (s HkManagerHandle) StoreList(c *gin.Context) {
	var req dashboardReq.HkManagerStoreList
	if err := c.ShouldBind(&req); err != nil {
		s.ResError(c, err)
		return
	}
	list, err := s.managerSvc.StoreList(s.WarpContext(c), req)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, list)
}

// MemberAuthList 获取香港管理者权限列表
func (s HkManagerHandle) MemberAuthList(c *gin.Context) {
	var req dashboardReq.HkManagerAuthList
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
}

// MemberAuthDetail 获取香港管理者权限详情
func (s HkManagerHandle) MemberAuthDetail(c *gin.Context) {
	var req dashboardReq.HkManagerAuthDetail
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
}

// SaveManagerAuth 创建/更新香港管理者权限
func (s HkManagerHandle) SaveManagerAuth(c *gin.Context) {
	var req dashboardReq.HkManagerAuthSave
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
}

// DeleteManagerAuth 删除香港管理者权限
func (s HkManagerHandle) DeleteManagerAuth(c *gin.Context) {
	var req dashboardReq.HkManagerRemove
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
}

// GetAllManagerIds 所有已选香港管理者memberIds
func (s HkManagerHandle) GetAllManagerIds(c *gin.Context) {
	var req request.EnterpriseBaseReq
	if err := c.ShouldBind(&req); err != nil {
		s.ResError(c, err)
		return
	}
	ids, err := s.managerSvc.GetAllManagerMemberIds(s.WarpContext(c), req.EnterpriseId)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, ids)
}
