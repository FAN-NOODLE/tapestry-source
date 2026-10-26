package dashboard

import (
	"github.com/gin-gonic/gin"
	"micro-scrm/internal/app/http/handle"
	"micro-scrm/internal/app/http/request"
	dashboardReq "micro-scrm/internal/app/http/request/dashboard"
	dashboardSvc "micro-scrm/internal/app/service/dashboard"
)

type AreaHandle struct {
	handle.BaseHandle
	areaSvc *dashboardSvc.AreaService
}

func NewAreaHandle() AreaHandle {
	return AreaHandle{
		areaSvc: dashboardSvc.NewAreaService(),
	}
}

// NotRelationStoreList 待绑定门店列表
func (s AreaHandle) NotRelationStoreList(c *gin.Context) {
	var params request.EnterpriseBaseReq
	if err := c.ShouldBind(&params); err != nil {
		s.ResError(c, err)
		return
	}
	list, err := s.areaSvc.NotRelationStoreList(s.WarpContext(c), params.EnterpriseId)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, list)

	return
}

// Tree 组织架构树
func (s AreaHandle) Tree(c *gin.Context) {
	var params dashboardReq.AreaTree
	if err := c.ShouldBind(&params); err != nil {
		s.ResError(c, err)
		return
	}
	list, err := s.areaSvc.GetAreaList(s.WarpContext(c), params)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, list)

	return
}

// Save 插入更新组织架构树节点
func (s AreaHandle) Save(c *gin.Context) {
	var params dashboardReq.AreaSave
	if err := c.ShouldBind(&params); err != nil {
		s.ResError(c, err)
		return
	}

	err := s.areaSvc.SaveData(s.WarpContext(c), params)
	if err != nil {
		s.ResError(c, err)
		return
	}

	s.ResSuccess(c, "")
	return
}

// Remove 移除组织架构树节点
func (s AreaHandle) Remove(c *gin.Context) {
	var params dashboardReq.AreaRemove
	if err := c.ShouldBind(&params); err != nil {
		s.ResError(c, err)
		return
	}
	err := s.areaSvc.RemoveArea(s.WarpContext(c), params)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, "")
	return
}

// ManagerAuth  管理者节点权限
func (s AreaHandle) ManagerAuth(c *gin.Context) {
	var params dashboardReq.ManagerAuthAreaIds
	if err := c.ShouldBind(&params); err != nil {
		s.ResError(c, err)
		return
	}
	areaIds, err := s.areaSvc.GetManagerAuthAreaIds(s.WarpContext(c), params)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, areaIds)

	return
}

// StoreCloseOrOpen  门店闭店或者重新开店
func (s AreaHandle) StoreCloseOrOpen(c *gin.Context) {
	var params dashboardReq.StoreCloseOrOpen
	if err := c.ShouldBind(&params); err != nil {
		s.ResError(c, err)
		return
	}
	err := s.areaSvc.StoreCloseOrOpen(s.WarpContext(c), params)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, "")
	return
}

// SearchStore  门店搜索
func (s AreaHandle) SearchStore(c *gin.Context) {
	var params dashboardReq.SearchStoreReq
	if err := c.ShouldBind(&params); err != nil {
		s.ResError(c, err)
		return
	}
	store, err := s.areaSvc.SearchStore(s.WarpContext(c), params)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, store)
	return
}

func (s AreaHandle) AreaList(c *gin.Context) {
	var params request.EnterpriseBaseReq
	if err := c.ShouldBind(&params); err != nil {
		s.ResError(c, err)
		return
	}
	store, err := s.areaSvc.AreaList(s.WarpContext(c), params.EnterpriseId)
	if err != nil {
		s.ResError(c, err)
		return
	}
	s.ResSuccess(c, store)
	return
}
