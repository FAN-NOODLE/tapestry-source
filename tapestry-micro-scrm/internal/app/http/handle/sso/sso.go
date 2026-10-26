package ssoHandle

import (
	"errors"
	"micro-scrm/internal/app/http/handle"
	ssoReq "micro-scrm/internal/app/http/request/sso"
	ssoResp "micro-scrm/internal/app/http/response/sso"
	memberSvc "micro-scrm/internal/app/service/member"
	ssoSvc "micro-scrm/internal/app/service/sso"
	"strings"

	"github.com/gin-gonic/gin"
)

type SSOHandle struct {
	handle.BaseHandle
	ssoService *ssoSvc.SSOService
}

func NewSSOHandle() SSOHandle {
	return SSOHandle{
		ssoService: ssoSvc.NewSSOService(),
	}
}

// Login SSO登录
func (h SSOHandle) Login(c *gin.Context) {
	var req ssoReq.SSOLoginReq
	if err := c.ShouldBind(&req); err != nil {
		h.ResError(c, err)
		return
	}

	userInfo, err := h.ssoService.Login(h.WarpContext(c), req.Code, req.EnterpriseId)
	if err != nil {
		h.ResError(c, err)
		return
	}
	userId := userInfo.EmployeeId
	if strings.HasPrefix(userId, "000") {
		userId = userId[3:]
	}
	membersService := memberSvc.NewMembersService()
	info, _ := membersService.GetInfoByUserId(h.WarpContext(c), req.EnterpriseId, userId)
	if info.Mobile == "" {
		info, _ = membersService.GetInfoByOpenUserId(h.WarpContext(c), req.EnterpriseId, userId)
		if info.Mobile == "" {
			h.ResError(c, errors.New("用户不存在,请联系管理员"))
			return
		}
	}
	response := ssoResp.SSOLoginResp{
		Mobile: info.Mobile,
	}

	h.ResSuccess(c, response)
}
