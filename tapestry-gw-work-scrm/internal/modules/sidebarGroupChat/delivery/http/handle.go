package sidebarGroupChathttp

import (
	"github.com/gin-gonic/gin"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ginplus"
	microScrm "gw-work-scrm/internal/pkg/sdk/microscrm"
	"net/http"
)

var sidebarGroupChatSdk *microScrm.SidebarGroupChatApi

// groupChatInfo 获取社群侧边栏基础信息
func groupChatInfo(c *gin.Context) {
	var req microScrm.GroupChatInfoRequest
	jsonError := c.BindJSON(&req)
	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}
	req.EnterpriseId = c.Value("enterpriseId").(int64)
	list, err := sidebarGroupChatSdk.GroupChatInfo(ctx.Wrap(c), req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, list)
}

// groupChatMember 侧边栏社群成员列表
func groupChatMember(c *gin.Context) {
	var req microScrm.GroupChatInfoRequest
	jsonError := c.BindJSON(&req)
	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}
	req.EnterpriseId = c.Value("enterpriseId").(int64)
	list, err := sidebarGroupChatSdk.GroupChatMember(ctx.Wrap(c), req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, list)
}

// groupChatMember 侧边栏社群成员列表
func groupUsersAddEmployees(c *gin.Context) {
	req := microScrm.GroupUsersAddEmployeesRequest{}
	jsonError := c.BindJSON(&req)
	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}
	req.EnterpriseId = c.Value("enterpriseId").(int64)
	list, err := sidebarGroupChatSdk.GroupUsersAddEmployeesRequest(ctx.Wrap(c), req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, list)
}

// groupChatMember 侧边栏社群成员列表
func groupUsersJoinActivity(c *gin.Context) {
	req := microScrm.GroupUsersJoinActivityRequest{}
	jsonError := c.BindJSON(&req)
	if jsonError != nil {
		ginplus.ResError(c, jsonError, http.StatusOK)
		return
	}
	req.EnterpriseId = c.Value("enterpriseId").(int64)
	list, err := sidebarGroupChatSdk.GroupUsersJoinActivityRequest(ctx.Wrap(c), req)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, list)
}
