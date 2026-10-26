package dashboard

import (
	"github.com/gin-gonic/gin"
	"gw-scrm/internal/app/http/handle"
	req "gw-scrm/internal/app/http/request"
	dashboardReq "gw-scrm/internal/app/http/request/dashboard"
	"gw-scrm/internal/pkg/common"
	"gw-scrm/internal/pkg/sdk/microscrm"
)

type ManagerHandle struct {
	handle.BaseHandle
	api microscrm.DashboardApi
}

func NewManagerHandle() ManagerHandle {
	return ManagerHandle{
		api: microscrm.NewDashboardApi(),
	}
}

// AuthSave   更新用户dashboard管理权限
func (h ManagerHandle) AuthSave(c *gin.Context) {
	var request dashboardReq.ManagerAuthSave
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseID
	res, err := h.api.ManagerAuthSave(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, res)

	return
}

// AuthDelete    删除用户dashboard管理权限
func (h ManagerHandle) AuthDelete(c *gin.Context) {
	var request dashboardReq.ManagerRemove
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseID
	res, err := h.api.ManagerAuthDelete(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, res)

	return
}

// AuthDetail    更新用户dashboard管理权限
func (h ManagerHandle) AuthDetail(c *gin.Context) {
	var request dashboardReq.ManagerAuthDetail
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseID
	res, err := h.api.ManagerAuthDetail(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, res)

	return
}

// AuthList   管理权限
func (h ManagerHandle) AuthList(c *gin.Context) {
	var request dashboardReq.ManagerAuthList
	if err := c.ShouldBind(&request); err != nil {
		h.ResError(c, err)
		return
	}
	if request.Page == 0 {
		request.Page = 1
	}

	if request.Size == 0 {
		request.Size = 10
	}
	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseID
	res, err := h.api.ManagerAuthList(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, res)

	return
}

// AuthManagerMemberIds     所有已配置管理者id
func (h ManagerHandle) AuthManagerMemberIds(c *gin.Context) {
	var request req.EnterpriseBaseReq
	enterpriseID, _ := common.GetEnterprise(c)
	request.EnterpriseId = enterpriseID
	res, err := h.api.ManagerAllMemberIds(h.WarpContext(c), request)
	if err != nil {
		h.ResError(c, err)
		return
	}

	h.ResSuccess(c, res)

	return
}
